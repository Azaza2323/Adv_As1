package main

import (
	"net/http"
)

func (a *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/news", a.getAllNews)
	mux.HandleFunc("/contact", a.contact)
	mux.HandleFunc("/create", a.showCreateForm).Methods("GET")

	// Handle the POST request for creating a post
	mux.HandleFunc("/addNews", a.addNewsHandler).Methods("POST")
	return mux
}
