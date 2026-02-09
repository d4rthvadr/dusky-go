package main

import (
	"net/http"

	"github.com/d4rthvadr/dusky-go/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	query := store.NewPaginatedFeedQuery()
	if err := query.Parse(r); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	posts, err := app.store.Posts.GetUserFeed(r.Context(), 31, query) // TODO: replace with actual user ID from context
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	type PostResponse struct {
		ID           int64    `json:"id"`
		UserID       int64    `json:"userId"`
		Username     string   `json:"username"`
		CommentCount int      `json:"commentCount"`
		Title        string   `json:"title"`
		Tags         []string `json:"tags"`
		Content      string   `json:"content"`
		CreatedAt    string   `json:"createdAt"`
		UpdatedAt    string   `json:"updatedAt"`
	}

	// Map the posts to the response format

	var response []PostResponse
	for _, post := range posts {
		response = append(response, PostResponse{
			ID:           post.ID,
			Title:        post.Title,
			UserID:       post.UserID,
			Username:     post.Username,
			Tags:         post.Tags,
			CommentCount: post.CommentCount,
			Content:      post.Content,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
		})
	}

	if err := writeResponse(w, http.StatusOK, response); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
