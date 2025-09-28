package handlers

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"tipjar"
	"tipjar/internal/auth"
	"tipjar/internal/config"
	"tipjar/internal/database"
	"tipjar/internal/models"
	"tipjar/internal/services"
	"tipjar/internal/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	db             *database.DB
	auth           *auth.Service
	cfg            *config.Config
	userService    *services.UserService
	tipJarService  *services.TipJarService
	sessionService *services.SessionService
}

func New(db *database.DB, authService *auth.Service, cfg *config.Config) *Handlers {
	return &Handlers{
		db:             db,
		auth:           authService,
		cfg:            cfg,
		userService:    services.NewUserService(db),
		tipJarService:  services.NewTipJarService(db),
		sessionService: services.NewSessionService(cfg.SessionSecret),
	}
}

func (h *Handlers) RegisterRoutes(e *echo.Echo) {
	// Setup static files
	h.setupStaticFiles(e)

	// Public routes
	e.GET("/", h.handleHome)
	e.GET("/login", h.handleLogin)
	e.GET("/auth/google", h.handleGoogleAuth)
	e.GET("/auth/callback", h.handleAuthCallback)
	e.POST("/logout", h.handleLogout)

	// Protected routes
	protected := e.Group("")
	protected.Use(h.requireAuth)
	protected.GET("/dashboard", h.handleDashboard)
	protected.GET("/jars", h.handleListJars)
	protected.GET("/jars/create", h.handleCreateJarForm)
	protected.POST("/jars", h.handleCreateJar)
	protected.GET("/jars/join", h.handleJoinJarForm)
	protected.POST("/jars/join", h.handleJoinJar)
	protected.GET("/jars/:id", h.handleViewJar)
	protected.POST("/jars/:id/offenses", h.handleReportOffense)
	protected.POST("/offenses/:id/pay", h.handlePayOffense)

	// API routes
	api := e.Group("/api/v1")
	api.Use(h.requireAuth)
	api.GET("/user", h.handleGetUser)
	api.GET("/jars", h.handleAPIListJars)
}

func (h *Handlers) setupStaticFiles(e *echo.Echo) {
	// Serve embedded static files
	staticFS, err := fs.Sub(tipjar.StaticFiles, "static")
	if err != nil {
		panic(fmt.Sprintf("Failed to create static filesystem: %v", err))
	}

	// Setup static file handler
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))))

	// Serve uploads directory from filesystem (not embedded)
	e.Static("/uploads", h.cfg.UploadsDir)
}

func (h *Handlers) handleHome(c echo.Context) error {
	// Check if user is already logged in
	if user := h.getCurrentUser(c); user != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func (h *Handlers) handleLogin(c echo.Context) error {
	// Check if user is already logged in
	if user := h.getCurrentUser(c); user != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
	}

	// Render login template
	return h.renderTemplate(c, templates.Login())
}

func (h *Handlers) handleGoogleAuth(c echo.Context) error {
	state, err := auth.GenerateState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate state")
	}

	// Store state in session cookie for CSRF protection
	cookie := &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	authURL := h.auth.GetAuthURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *Handlers) handleAuthCallback(c echo.Context) error {
	// Verify state parameter for CSRF protection
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing state cookie")
	}

	stateParam := c.QueryParam("state")
	if stateParam != stateCookie.Value {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid state parameter")
	}

	// Clear the state cookie
	cookie := &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	c.SetCookie(cookie)

	// Exchange authorization code for token
	code := c.QueryParam("code")
	if code == "" {
		// Check if there was an OAuth error
		if errorCode := c.QueryParam("error"); errorCode != "" {
			errorDesc := c.QueryParam("error_description")
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("OAuth error: %s - %s", errorCode, errorDesc))
		}
		return echo.NewHTTPError(http.StatusBadRequest, "Missing authorization code")
	}

	token, err := h.auth.ExchangeCode(c.Request().Context(), code)
	if err != nil {
		c.Logger().Error("Failed to exchange OAuth code", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to exchange code for token")
	}

	// Get user info from Google
	googleUser, err := h.auth.GetUserInfo(c.Request().Context(), token)
	if err != nil {
		c.Logger().Error("Failed to get user info from Google", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user info")
	}

	c.Logger().Info("OAuth user info received", "email", googleUser.Email, "name", googleUser.Name, "google_id", googleUser.ID)

	// Check if user exists in our database
	user, err := h.userService.GetUserByGoogleID(c.Request().Context(), googleUser.ID)
	if err != nil {
		c.Logger().Error("Database error when looking up user", "error", err, "google_id", googleUser.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	// Create new user if they don't exist
	if user == nil {
		c.Logger().Info("Creating new user", "email", googleUser.Email, "name", googleUser.Name, "google_id", googleUser.ID)
		user, err = h.userService.CreateUser(
			c.Request().Context(),
			googleUser.Email,
			googleUser.Name,
			googleUser.Picture,
			googleUser.ID,
		)
		if err != nil {
			c.Logger().Error("Failed to create new user", "error", err, "email", googleUser.Email)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
		}
		c.Logger().Info("Successfully created new user", "user_id", user.ID, "email", user.Email)
	} else {
		c.Logger().Info("Existing user found", "user_id", user.ID, "email", user.Email)
	}

	// Create session
	sessionToken, err := h.sessionService.CreateSession(user.ID, user.Email, user.Name)
	if err != nil {
		c.Logger().Error("Failed to create session", "error", err, "user_id", user.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create session")
	}

	// Set session cookie
	h.sessionService.SetSessionCookie(c.Response().Writer, sessionToken)

	c.Logger().Info("User successfully authenticated", "user_id", user.ID, "email", user.Email)

	// Redirect to dashboard
	return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}

func (h *Handlers) handleLogout(c echo.Context) error {
	h.sessionService.ClearSessionCookie(c.Response().Writer)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func (h *Handlers) handleDashboard(c echo.Context) error {
	user := h.getCurrentUser(c)
	
	// Get user's tip jars
	jars, err := h.tipJarService.ListTipJarsForUser(c.Request().Context(), user.ID)
	if err != nil {
		c.Logger().Error("Failed to load user's tip jars", "error", err, "user_id", user.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load tip jars")
	}

	return h.renderTemplate(c, templates.Dashboard(user, jars))
}

func (h *Handlers) handleCreateJarForm(c echo.Context) error {
	user := h.getCurrentUser(c)
	return h.renderTemplate(c, templates.CreateJar(user))
}

func (h *Handlers) handleCreateJar(c echo.Context) error {
	user := h.getCurrentUser(c)
	
	name := strings.TrimSpace(c.FormValue("name"))
	description := strings.TrimSpace(c.FormValue("description"))
	inviteCode := strings.TrimSpace(c.FormValue("invite_code"))

	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Jar name is required")
	}
	
	if inviteCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invite code is required")
	}
	
	if len(inviteCode) != 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invite code must be 8 characters")
	}
	
	existingJar, err := h.tipJarService.GetTipJarByInviteCode(c.Request().Context(), inviteCode)
	if err != nil {
		c.Logger().Error("Failed to check invite code", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to validate invite code")
	}
	
	if existingJar != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invite code already exists. Please generate a new one.")
	}
	
	jar, err := h.tipJarService.CreateTipJarWithInviteCode(c.Request().Context(), name, description, inviteCode, user.ID)
	if err != nil {
		c.Logger().Error("Failed to create tip jar", "error", err, "user_id", user.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create tip jar")
	}
	
	c.Logger().Info("Tip jar created successfully", "jar_id", jar.ID, "name", jar.Name, "created_by", user.ID)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"jar_id":   jar.ID,
		"redirect": fmt.Sprintf("/jars/%d", jar.ID),
	})
}

func (h *Handlers) handleJoinJarForm(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not ready yet")
}

// Placeholder handlers - to be implemented in future iterations
func (h *Handlers) handleListJars(c echo.Context) error {
	return c.String(http.StatusOK, "List jars - not implemented yet")
}


func (h *Handlers) handleJoinJar(c echo.Context) error {
	return c.String(http.StatusOK, "Join jar - not implemented yet")
}

func (h *Handlers) handleViewJar(c echo.Context) error {
	return c.String(http.StatusOK, "View jar - not implemented yet")
}

func (h *Handlers) handleReportOffense(c echo.Context) error {
	return c.String(http.StatusOK, "Report offense - not implemented yet")
}

func (h *Handlers) handlePayOffense(c echo.Context) error {
	return c.String(http.StatusOK, "Pay offense - not implemented yet")
}

func (h *Handlers) handleGetUser(c echo.Context) error {
	user := h.getCurrentUser(c)
	return c.JSON(http.StatusOK, user)
}

func (h *Handlers) handleAPIListJars(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "API list jars - not implemented yet"})
}

func (h *Handlers) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := h.getCurrentUser(c)
		if user == nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
		
		// Store user in context for easy access
		c.Set("user", user)
		return next(c)
	}
}

func (h *Handlers) getCurrentUser(c echo.Context) *models.User {
	// Try to get from context first (set by middleware)
	if user, ok := c.Get("user").(*models.User); ok {
		return user
	}

	// Get session cookie
	cookie, err := c.Cookie("session")
	if err != nil {
		return nil
	}

	// Validate session
	sessionData, err := h.sessionService.ValidateSession(cookie.Value)
	if err != nil {
		return nil
	}

	// Get user from database
	user, err := h.userService.GetUserByID(c.Request().Context(), sessionData.UserID)
	if err != nil {
		return nil
	}

	return user
}

func (h *Handlers) renderTemplate(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response().Writer)
}