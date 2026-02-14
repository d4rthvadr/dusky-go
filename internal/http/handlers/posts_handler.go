package handlers

import (
	"errors"
	"net/http"

	errCustom "github.com/d4rthvadr/dusky-go/internal/errors"
	"github.com/d4rthvadr/dusky-go/internal/models"
)

type createPostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=5000"`
	Tags    []string `json:"tags,omitempty" validate:"dive,required"`
}

const PostIDKey string = "postID"

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with the provided title, content, and tags.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		createPostPayload	true	"Post payload"
//	@Success		201		{object}	models.Post
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/posts [post]
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var userID int64 = 1

	var post createPostPayload
	if err := readJSON(r, &post); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validatorInstance.Struct(post); err != nil {
		writeValidationError(w, err)
		return
	}

	postModel := models.Post{
		Title:   post.Title,
		Content: post.Content,
		UserID:  userID,
		Tags:    post.Tags,
	}

	if err := h.store.Posts.Create(r.Context(), &postModel); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	if err := writeResponse(w, http.StatusCreated, postModel); err != nil {
		h.internalServerError(w, r, err)
		return
	}
}

// getPostID extracts post ID from route params.
func (h *Handler) getPostID(r *http.Request) (int64, error) {
	postID, err := parseIDParam(r, PostIDKey)
	if err != nil {
		return 0, err
	}
	return postID, nil
}

// GetPost godoc
//
//	@Summary		Get a post by ID
//	@Description	Get a post by its ID, including its comments.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int64	true	"Post ID"
//	@Success		200		{object}	models.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/posts/{postID} [get]
func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID, err := h.getPostID(r)
	if err != nil {
		h.badRequestError(w, r, err)
		return
	}

	post, err := h.store.Posts.GetByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, errCustom.ErrResourceNotFound):
			h.notFoundError(w, r, err)
			return
		default:
			h.internalServerError(w, r, err)
			return
		}
	}

	comments, err := h.store.Comments.GetByPostID(ctx, postID)
	if err != nil {
		h.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := writeResponse(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// DeletePost godoc
//
//	@Summary		Delete a post by ID
//	@Description	Delete a post by its ID.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	int64	true	"Post ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/posts/{postID} [delete]
func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID, err := h.getPostID(r)
	if err != nil {
		h.badRequestError(w, r, err)
		return
	}

	if err := h.store.Posts.Delete(ctx, postID); err != nil {
		switch {
		case errors.Is(err, errCustom.ErrResourceNotFound):
			h.notFoundError(w, r, err)
			return
		default:
			h.internalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePost godoc
//
//	@Summary		Update a post by ID
//	@Description	Update a post's title, content, and tags by its ID.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int64				true	"Post ID"
//	@Param			post	body		createPostPayload	true	"Updated post payload"
//	@Success		200		{object}	models.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/posts/{postID} [patch]
func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID, err := h.getPostID(r)
	if err != nil {
		h.badRequestError(w, r, err)
		return
	}

	var payload createPostPayload
	if err := readJSON(r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validatorInstance.Struct(payload); err != nil {
		writeValidationError(w, err)
		return
	}

	postModel := models.Post{
		ID:      postID,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	if err := h.store.Posts.Update(ctx, &postModel); err != nil {
		switch {
		case errors.Is(err, errCustom.ErrResourceNotFound):
			h.notFoundError(w, r, err)
			return
		default:
			h.internalServerError(w, r, err)
			return
		}
	}

	if err := writeResponse(w, http.StatusOK, postModel); err != nil {
		h.internalServerError(w, r, err)
		return
	}
}
