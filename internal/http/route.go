package http

import (
	"fmt"
	"strings"

	"github.com/d4rthvadr/dusky-go/internal/http/handlers"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// MountV1Routes sets up the API routes for version 1 of the API. It includes endpoints for health checks, Swagger documentation, and CRUD operations for posts and users.
func MountV1Routes(r chi.Router, handler *handlers.Handler, apiURL string) {
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", handler.HealthCheck)

		normalizedAPIURL := strings.TrimRight(apiURL, "/")
		docsURL := fmt.Sprintf("%s/swagger/doc.json", normalizedAPIURL)

		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", handler.CreatePost)

			r.Route("/{postID}", func(r chi.Router) {
				r.Get("/", handler.GetPost)
				r.Delete("/", handler.DeletePost)
				r.Patch("/", handler.UpdatePost)
			})
		})

		r.Route("/users", func(r chi.Router) {

			r.Put("/activate/{token}", handler.ActivateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(handler.UserContextMiddleware)
				r.Get("/", handler.GetUser)
				r.Put("/follow", handler.FollowUser)
				r.Put("/unfollow", handler.UnfollowUser)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", handler.GetUserFeed)
			})

		})

		// Public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handler.RegisterUser)

		})
	})
}
