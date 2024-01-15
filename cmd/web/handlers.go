package main

import (
	"html/template"
	"net/http"
)

func (a *application) getAllNews(w http.ResponseWriter, r *http.Request) {
	newsList, err := a.news.Latest()
	if err != nil {
		a.errorLog.Println("Error:", err)
		http.Error(w, "Error retrieving news", http.StatusInternalServerError)
		return
	}

	templatePath := "ui/html/news.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{"News": newsList})
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
func (a *application) home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/html/index.html")
}
func (a *application) contact(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/html/contact.html")
}

func (a *application) addNewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	audience := r.Form.Get("audience")
	author := r.Form.Get("author")
	title := r.Form.Get("title")
	description := r.Form.Get("description")
	content := r.Form.Get("content")

	_, err = a.news.Insert(audience, author, title, description, content)
	if err != nil {
		http.Error(w, "Error adding news", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (a *application) createNewsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/html/create.html")
}
