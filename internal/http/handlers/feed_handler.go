package handlers

import (
	"net/http"

	"github.com/d4rthvadr/dusky-go/internal/store"
)

type FeedPostResponse struct {
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

// GetUserFeed godoc
//
//	@Summary		Get the authenticated user's feed
//	@Description	Retrieve a paginated list of posts from users that the authenticated user follows.
//	@Tags			Feed
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number for pagination"
//	@Param			size	query		int	false	"Number of items per page"
//	@Success		200		{array}		FeedPostResponse
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (h *Handler) GetUserFeed(w http.ResponseWriter, r *http.Request) {
	query := store.NewPaginatedFeedQuery()
	if err := query.Parse(r); err != nil {
		h.badRequestError(w, r, err)
		return
	}

	posts, err := h.store.Posts.GetUserFeed(r.Context(), 1, query)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	response := make([]FeedPostResponse, 0, len(posts))
	for _, post := range posts {
		response = append(response, FeedPostResponse{
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
		h.internalServerError(w, r, err)
		return
	}
}
