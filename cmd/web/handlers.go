package main

import (
	models "asik1/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (a *application) getAllNews(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("userInfo")
	if err != nil {
		a.render(w, "ui/html/index.html", nil)
	}

	values := strings.Split(cookie.Value, "&")
	var username, role string
	for _, v := range values {
		if strings.HasPrefix(v, "username=") {
			username = strings.TrimPrefix(v, "username=")
		} else if strings.HasPrefix(v, "role=") {
			role = strings.TrimPrefix(v, "role=")
		}
	}

	newsList, err := a.news.Latest()
	if err != nil {
		a.errorLog.Println("Error:", err)
		http.Error(w, "Error retrieving news", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Username": username,
		"Role":     role,
		"News":     newsList,
	}

	a.render(w, "ui/html/news.html", data)
}

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("userInfo")
	if err != nil {
		a.render(w, "ui/html/index.html", nil)
		return
	}

	values := strings.Split(cookie.Value, "&")
	var username, role string
	for _, v := range values {
		if strings.HasPrefix(v, "username=") {
			username = strings.TrimPrefix(v, "username=")
		} else if strings.HasPrefix(v, "role=") {
			role = strings.TrimPrefix(v, "role=")
		}
	}

	var page string
	switch role {
	case "creator":
		page = "ui/html/creator.html"
	case "reader":
		page = "ui/html/reader.html"
	case "admin":
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	default:
		page = "ui/html/index.html"
	}
	data := map[string]interface{}{
		"Username": username,
		"Role":     role,
	}

	a.render(w, page, data)
}

func (a *application) contact(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/html/contact.html")
}
func (a *application) addNewsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.render(w, "ui/html/create.html", nil)
	case http.MethodPost:
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
		a.session.Put(r, "flash", "Snippet successfully created!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (a *application) showCreateForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, "ui/html/create.html", nil)
}
func (a *application) showDepartmentForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, "ui/html/dep.html", nil)
}
func (a *application) byAudience(w http.ResponseWriter, r *http.Request) {
	audience := r.URL.Query().Get("audience")
	a.infoLog.Println("Audience:", audience)
	if audience != "students" && audience != "staff" && audience != "applications" {
		http.NotFound(w, r)
		return
	}
	list, err := a.news.GetByAudience(audience)
	if err != nil {
		a.errorLog.Println("Error:", err)
		http.Error(w, "Error retrieving news", http.StatusInternalServerError)
		return
	}
	templatePath := "ui/html/audience.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	data := struct {
		Audience string
		NewsList []*models.News
	}{
		Audience: audience,
		NewsList: list,
	}
	a.infoLog.Println("Rendering template...")
	err = tmpl.Execute(w, data)
	if err != nil {
		a.errorLog.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
func (a *application) fillDep(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form: "+err.Error(), http.StatusInternalServerError)
			return
		}

		idStr := r.Form.Get("id")
		depIDStr := r.Form.Get("depID")
		staffQStr := r.Form.Get("staff_q")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Error converting 'id' to integer: "+err.Error(), http.StatusBadRequest)
			return
		}

		depID, err := strconv.Atoi(depIDStr)
		if err != nil {
			http.Error(w, "Error converting 'depID' to integer: "+err.Error(), http.StatusBadRequest)
			return
		}

		staffQ, err := strconv.Atoi(staffQStr)
		if err != nil {
			http.Error(w, "Error converting 'staff_q' to integer: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := a.departments.DB.Ping(); err != nil {
			a.errorLog.Println("Error: Database connection is not alive")
			http.Error(w, "Error adding department", http.StatusInternalServerError)
			return
		}
		_, err = a.departments.InsertDepo(id, depID, staffQ)
		if err != nil {
			a.errorLog.Println("Error adding department:", err)
			http.Error(w, "Error adding department", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (a *application) getAllDeps(w http.ResponseWriter, r *http.Request) {
	newsList, err := a.departments.Deps()
	if err != nil {
		a.errorLog.Println("Error:", err)
		http.Error(w, "Error retrieving news", http.StatusInternalServerError)
		return
	}
	a.render(w, "ui/html/deps.html", map[string]interface{}{"News": newsList})
}

func (a *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, "ui/html/login.html", nil)
}
func (a *application) loginUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, err := a.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			a.serverError(w, err)
		}
		return
	}
	a.session.Put(r, "authenticatedUserID", id)
	a.infoLog.Printf("User ID %d has been stored in the session", id)
	role, err := a.users.GetRoleByEmail(email)
	if err != nil {
		a.serverError(w, err)
		return
	}
	name, err := a.users.GetNameByEmail(email)
	if err != nil {
		a.serverError(w, err)
		return
	}
	//data := struct {
	//	UserName string
	//}{
	//	UserName: name,
	//}
	cookie := http.Cookie{
		Name:  "userInfo",
		Value: fmt.Sprintf("role=%s&username=%s", role, name),
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	switch role {
	case "creator":
		http.Redirect(w, r, "/creator", http.StatusSeeOther)
	case "reader":
		http.Redirect(w, r, "/reader", http.StatusSeeOther)

	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (a *application) signupUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

	if name == "" || email == "" || password == "" || role == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	existingUser, err := a.users.GetByEmail(email)
	if err != nil {
		http.Error(w, "Error checking existing email", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "Email already registered", http.StatusBadRequest)
		return
	}

	err = a.users.Insert(name, email, password, role)
	if err != nil {
		http.Error(w, "Error signing up user", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, "ui/html/signup.html", nil)
}
func (a *application) creatorPage(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("userInfo")
	if err != nil {
		a.render(w, "ui/html/index.html", nil)
	}

	values := strings.Split(cookie.Value, "&")
	var username, role string
	for _, v := range values {
		if strings.HasPrefix(v, "username=") {
			username = strings.TrimPrefix(v, "username=")
		} else if strings.HasPrefix(v, "role=") {
			role = strings.TrimPrefix(v, "role=")
		}
	}
	data := map[string]interface{}{
		"Username": username,
		"Role":     role,
	}
	a.render(w, "ui/html/creator.html", data)
}

func (a *application) readerPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, "ui/html/reader.html", nil)
}
func (a *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	a.session.Remove(r, "authenticatedUserID")
	cookie := http.Cookie{
		Name:     "userInfo",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (a *application) adminPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("userInfo")
	if err != nil {
		a.render(w, "ui/html/index.html", nil)
	}

	values := strings.Split(cookie.Value, "&")
	var username, role string
	for _, v := range values {
		if strings.HasPrefix(v, "username=") {
			username = strings.TrimPrefix(v, "username=")
		} else if strings.HasPrefix(v, "role=") {
			role = strings.TrimPrefix(v, "role=")
		}
	}
	list, err := a.users.GetAllUsers()
	if err != nil {
		a.errorLog.Println("Error:", err)
		http.Error(w, "Error retrieving news", http.StatusInternalServerError)
		return
	}
	templatePath := "ui/html/admin.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	data := struct {
		Username string
		Role     string
		Users    []*models.User
	}{
		Username: username,
		Role:     role,
		Users:    list,
	}
	a.infoLog.Println("Rendering template...")
	err = tmpl.Execute(w, data)
	if err != nil {
		a.errorLog.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
func (a *application) changeRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestBody struct {
		UserID  string `json:"userId"`
		NewRole string `json:"newRole"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err := a.users.ChangeUserRole(requestBody.UserID, requestBody.NewRole)
	if err != nil {
		a.errorLog.Println("Error updating user role:", err)
		http.Error(w, "Error updating user role", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Role updated successfully"})
}
