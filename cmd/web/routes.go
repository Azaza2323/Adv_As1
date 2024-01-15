package main

import (
	"net/http"
)

func (a *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/news", a.getAllNews)
	mux.HandleFunc("/contact", a.contact)
	mux.HandleFunc("/create/add", a.addNewsHandler)
	mux.HandleFunc("/create", a.createNewsHandler)
	return mux
}
