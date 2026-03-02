package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestUser(t *testing.T) {

	t.Run("should not allow unauthenticated users to access user details", func(t *testing.T) {

		// Arrange
		app := newTestApplication(t)
		mux := app.mount()

		userID := int64(1)

		endpoint := fmt.Sprintf("/v1/users/%d", userID)

		// Act
		request, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := executeRequest(mux, request)

		// Assert
		checkResponseCode(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("should allow authenticated users to access user details", func(t *testing.T) {

		// Arrange
		app := newTestApplication(t)
		mux := app.mount()

		userID := int64(1)
		endpoint := fmt.Sprintf("/v1/users/%d", userID)

		// Act
		request, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			t.Fatal(err)
		}

		token, err := generateTokenForUser(userID, app.jwtAuthenticator)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response := executeRequest(mux, request)

		// Assert
		checkResponseCode(t, http.StatusOK, response.Code)
	})

}
