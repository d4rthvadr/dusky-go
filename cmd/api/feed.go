package main

import (
	"net/http"

	"github.com/d4rthvadr/dusky-go/internal/store"
)

// GetUserFeed godoc
//
//	@Summary		Get the user's feed
//	@Description	Get a paginated list of posts from users that the authenticated user follows.
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//
//	@param			limit	query		int			false	"Number of items per page for pagination (default: 10)"
//	@param			offset	query		int			false	"Page number for pagination (default: 1)"
//	@param			search	query		string		false	"Search term to filter posts by title or content"
//	@param			tags	query		[]string	false	"Comma-separated list of tags to filter posts by (e.g., tag1,tag2,tag3)"
//	@param			sort	query		string		false	"Sort order for posts (e.g., createdAt:desc or createdAt:asc)"
//	@Success		200		{array}		map[string]interface{}
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	query := store.NewPaginatedFeedQuery()
	if err := query.Parse(r); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	posts, err := app.store.Posts.GetUserFeed(r.Context(), 1, query) // TODO: replace with actual user ID from context
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
