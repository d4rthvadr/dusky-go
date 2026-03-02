package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/auth"
	"github.com/d4rthvadr/dusky-go/internal/cache"
	"github.com/d4rthvadr/dusky-go/internal/config"
	"github.com/d4rthvadr/dusky-go/internal/http/handlers"
	"github.com/d4rthvadr/dusky-go/internal/mailer"
	"github.com/d4rthvadr/dusky-go/internal/store"
	"github.com/d4rthvadr/dusky-go/internal/utils/logger"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func newTestApplication(t *testing.T) *application {

	t.Helper()

	//logger := logger.NewLoggerMock()
	logger := logger.NewLogger()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockCache()
	mockMailer := &mailer.MockMailer{}

	// Create JWT authenticator with test secret
	jwtAuthenticator := auth.NewJWTAuthenticator("test-secret-key", "test-audience", "test-issuer", 3600)

	return &application{
		store:            mockStore,
		cache:            mockCache,
		logger:           logger,
		jwtAuthenticator: jwtAuthenticator,
		handler: handlers.New(handlers.HandlerOptions{
			Store:            mockStore,
			Cache:            mockCache,
			Version:          "test",
			Logger:           logger,
			MailConfig:       config.MailConfig{},
			Mailer:           mockMailer,
			JWTAuthenticator: jwtAuthenticator,
			IsProdEnv:        false,
		}),
	}

}

func generateTokenForUser(userID int64, jwtAuthenticator *auth.JWTAuthenticator) (string, error) {

	claims := jwt.MapClaims{
		"sub": userID,
		"aud": jwtAuthenticator.Aud,
		"iss": jwtAuthenticator.Iss,
		"exp": time.Now().Unix() + jwtAuthenticator.Exp,
	}

	return jwtAuthenticator.GenerateToken(claims)

}

func executeRequest(mux *chi.Mux, request *http.Request) *httptest.ResponseRecorder {

	responseRecorder := httptest.NewRecorder()

	mux.ServeHTTP(responseRecorder, request)

	return responseRecorder
}

func checkResponseCode(t *testing.T, expected, actual int) {

	t.Helper()

	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}

}
