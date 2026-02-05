package main

import (
	"errors"
	"net/http"

	"github.com/d4rthvadr/dusky-go/internal/models"

	errCustom "github.com/d4rthvadr/dusky-go/internal/errors"
)

type createPostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=5000"`
	Tags    []string `json:"tags,omitempty" validate:"dive,required"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var userID int64 = 1 // Placeholder for authenticated user ID

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

	if err := app.store.Posts.Create(r.Context(), &postModel); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, postModel); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	postID, err := parseIDParam(r, "postID")
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post, err := app.store.Posts.GetByID(ctx, postID)

	if err != nil {

		switch {
		case errors.Is(err, errCustom.ErrResourceNotFound):
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}

	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) listPostsHandler(w http.ResponseWriter, r *http.Request) {

	var posts []models.Post

	err := app.store.Posts.List(r.Context(), &posts)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
	}
}
