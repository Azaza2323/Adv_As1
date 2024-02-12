package main

import (
	"golang.org/x/net/context"
	"net/http"
)

import (
	"strings"
)

func (a *application) isAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("userInfo")
		if err != nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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

		if username != "" && role != "" {
			r = r.WithContext(context.WithValue(r.Context(), "username", username))
			r = r.WithContext(context.WithValue(r.Context(), "role", role))
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}
	})
}
