package main

import (
	"fmt"
	"net/http"
)

func main() {

	api := NewApi(":8082")

	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	// routes registration
	mux.HandleFunc("GET /users", api.getUsersHandler)
	mux.HandleFunc("POST /users", api.createUserHandler)

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}

}
