package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (a *application) routes() http.Handler {
	dynamicMiddleware := alice.New(a.session.Enable)
	publicMiddleware := alice.New()
	mux := pat.New()
	mux.Get("/", publicMiddleware.ThenFunc(a.home))
	mux.Get("/news", dynamicMiddleware.Append(a.isAuth).ThenFunc(a.getAllNews))
	mux.Get("/contact", publicMiddleware.Append(a.isAuth).ThenFunc(a.contact))
	mux.Get("/deps", publicMiddleware.Append(a.isAuth).ThenFunc(a.getAllDeps))
	mux.Get("/deps/create", dynamicMiddleware.Append(a.isAuth).ThenFunc(a.showDepartmentForm))
	mux.Get("/create", dynamicMiddleware.Append(a.isAuth).ThenFunc(a.showCreateForm))
	mux.Get("/for", dynamicMiddleware.ThenFunc(a.byAudience))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(a.signupUserForm))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(a.loginUserForm))
	mux.Get("/creator", dynamicMiddleware.Append(a.isAuth).ThenFunc(a.creatorPage))
	mux.Get("/reader", dynamicMiddleware.Append(a.isAuth).ThenFunc(a.readerPage))
	mux.Get("/admin", dynamicMiddleware.Append(a.isAuth).ThenFunc(a.adminPage))
	mux.Post("/user/logout", dynamicMiddleware.Append(a.isAuth).ThenFunc(a.logoutUser))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(a.signupUser))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(a.loginUser))
	mux.Post("/create/add", dynamicMiddleware.Append(a.isAuth).ThenFunc(a.addNewsHandler))
	mux.Post("/deps/create/add", dynamicMiddleware.Append(a.isAuth).ThenFunc(a.fillDep))
	mux.Post("/changeRole", dynamicMiddleware.ThenFunc(a.changeRoleHandler))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
