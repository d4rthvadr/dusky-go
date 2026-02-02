package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type api struct {
	addr string
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = []User{}

func NewApi(addr string) *api {
	return &api{addr: addr}
}

func (s *api) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// encode users to JSON and write to response
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

}

func generateUserID() string {
	return fmt.Sprintf("user_%d", len(users)+1)
}

func (s *api) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload User
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser := User{
		ID: payload.ID,

		Name: payload.Name,
	}

	err = insertUser(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func insertUser(u User) error {

	if u.Name == "" {
		return errors.New("user name cannot be empty")
	}

	if u.ID == "" {
		u.ID = generateUserID()
	}
	users = append(users, u)
	return nil
}
