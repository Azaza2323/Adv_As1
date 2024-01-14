package main

import (
	models "asik1/pkg"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	news, err := app.news.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, news)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) byCategory(w http.ResponseWriter, r *http.Request) {
	audience := r.URL.Query().Get("audience")
	if audience != "students" && audience != "staff" && audience != "applicants" {
		http.NotFound(w, r)
		return
	}
	news, err := app.news.GetByAudience(audience)
	if err != nil {
		app.serverError(w, err)
		return
	}
	files := []string{
		"./ui/html/for_" + audience + ".page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, news)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
func (app *application) update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	var updatedNews *models.News
	err := json.NewDecoder(r.Body).Decode(&updatedNews)
	if err != nil || id < 1 {
		log.Println(err.Error())
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	result, err := app.news.Update(updatedNews.ID, updatedNews.Audience, updatedNews.Author, updatedNews.Title, updatedNews.Description, updatedNews.Content)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	fmt.Println(result)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	_, err := app.news.Delete(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}
func (app *application) add(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/add.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
func (app *application) create(w http.ResponseWriter, r *http.Request) {
	var createdNews *models.News
	err := json.NewDecoder(r.Body).Decode(&createdNews)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	result, err := app.news.Insert(createdNews.Audience, createdNews.Author, createdNews.Title, createdNews.Description, createdNews.Content)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	fmt.Println(result)
}
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
