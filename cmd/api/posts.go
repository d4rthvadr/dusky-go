package main

import (
	"net/http"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

type createPostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags,omitempty"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var userID int64 = 1 // Placeholder for authenticated user ID

	var post createPostPayload
	err := readJSON(r, &post)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	postModel := models.Post{
		Title:   post.Title,
		Content: post.Content,
		UserID:  userID,
		Tags:    post.Tags,
	}

	err = app.store.Posts.Create(r.Context(), &postModel)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, postModel); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
