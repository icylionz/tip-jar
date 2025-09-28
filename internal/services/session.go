package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type SessionService struct {
	secret []byte
}

type SessionData struct {
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewSessionService(secret string) *SessionService {
	return &SessionService{
		secret: []byte(secret),
	}
}

func (s *SessionService) CreateSession(userID int, email, name string) (string, error) {
	sessionData := SessionData{
		UserID:    userID,
		Email:     email,
		Name:      name,
		ExpiresAt: time.Now().Add(24 * time.Hour * 7), // 7 days
	}

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(jsonData)
	signature := s.sign(encoded)

	return fmt.Sprintf("%s.%s", encoded, signature), nil
}

func (s *SessionService) ValidateSession(sessionToken string) (*SessionData, error) {
	parts := strings.Split(sessionToken, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid session format")
	}

	encoded, signature := parts[0], parts[1]

	// Verify signature
	expectedSignature := s.sign(encoded)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, fmt.Errorf("invalid session signature")
	}

	// Decode session data
	jsonData, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode session: %w", err)
	}

	var sessionData SessionData
	if err := json.Unmarshal(jsonData, &sessionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Check expiration
	if time.Now().After(sessionData.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	return &sessionData, nil
}

func (s *SessionService) SetSessionCookie(w http.ResponseWriter, sessionToken string) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   int((24 * time.Hour * 7).Seconds()), // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func (s *SessionService) ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func (s *SessionService) sign(data string) string {
	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}