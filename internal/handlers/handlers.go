package handlers

import (
	"net/http"

	"tipjar/internal/auth"
	"tipjar/internal/config"
	"tipjar/internal/database"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	db   *database.DB
	auth *auth.Service
	cfg  *config.Config
}

func New(db *database.DB, authService *auth.Service, cfg *config.Config) *Handlers {
	return &Handlers{
		db:   db,
		auth: authService,
		cfg:  cfg,
	}
}

func (h *Handlers) RegisterRoutes(e *echo.Echo) {
	// Static files
	e.Static("/static", "static")
	e.Static("/uploads", h.cfg.UploadsDir)

	// Public routes
	e.GET("/", h.handleHome)
	e.GET("/login", h.handleLogin)
	e.GET("/auth/callback", h.handleAuthCallback)

	// Protected routes
	protected := e.Group("")
	protected.Use(h.requireAuth)
	protected.GET("/dashboard", h.handleDashboard)
	protected.GET("/jars", h.handleListJars)
	protected.POST("/jars", h.handleCreateJar)
	protected.GET("/jars/:id", h.handleViewJar)
	protected.POST("/jars/:id/join", h.handleJoinJar)
	protected.POST("/jars/:id/offenses", h.handleReportOffense)
	protected.POST("/offenses/:id/pay", h.handlePayOffense)

	// API routes
	api := e.Group("/api/v1")
	api.Use(h.requireAuth)
	api.GET("/user", h.handleGetUser)
	api.GET("/jars", h.handleAPIListJars)
}

func (h *Handlers) handleHome(c echo.Context) error {
	return c.HTML(http.StatusOK, "<h1>Tip Jar - Coming Soon</h1>")
}

func (h *Handlers) handleLogin(c echo.Context) error {
	state, err := auth.GenerateState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate state")
	}

	// Store state in session/cookie for CSRF protection
	authURL := h.auth.GetAuthURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *Handlers) handleAuthCallback(c echo.Context) error {
	// TODO: Implement OAuth callback
	return c.String(http.StatusOK, "Auth callback - not implemented yet")
}

func (h *Handlers) handleDashboard(c echo.Context) error {
	return c.String(http.StatusOK, "Dashboard - not implemented yet")
}

func (h *Handlers) handleListJars(c echo.Context) error {
	return c.String(http.StatusOK, "List jars - not implemented yet")
}

func (h *Handlers) handleCreateJar(c echo.Context) error {
	return c.String(http.StatusOK, "Create jar - not implemented yet")
}

func (h *Handlers) handleViewJar(c echo.Context) error {
	return c.String(http.StatusOK, "View jar - not implemented yet")
}

func (h *Handlers) handleJoinJar(c echo.Context) error {
	return c.String(http.StatusOK, "Join jar - not implemented yet")
}

func (h *Handlers) handleReportOffense(c echo.Context) error {
	return c.String(http.StatusOK, "Report offense - not implemented yet")
}

func (h *Handlers) handlePayOffense(c echo.Context) error {
	return c.String(http.StatusOK, "Pay offense - not implemented yet")
}

func (h *Handlers) handleGetUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Get user - not implemented yet"})
}

func (h *Handlers) handleAPIListJars(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "API list jars - not implemented yet"})
}

func (h *Handlers) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: Implement session-based authentication middleware
		return next(c)
	}
}
