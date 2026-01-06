package web

import (
	"html/template"
	"net/http"
	"time"
)

var username, password string

// Set credentials from config file
func SetCredentials(user, pass string) {
	username = user
	password = pass
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	data := PageData{}

	if r.Method == http.MethodPost {
		r.ParseForm()
		u := r.FormValue("username")
		p := r.FormValue("password")

		if u == username && p == password {
			// Set session cookie for 24 hours
			http.SetCookie(w, &http.Cookie{
				Name:     "bastille-web-auth",
				Value:    "authenticated",
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(24 * time.Hour),
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			data.Error = "Invalid username or password"
		}
	}

	// Login page needs its own template
	tmpl := "web/static/login.html"
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Template ERROR: "+err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "bastille-web-auth",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Enforce login cookie on every request
func requireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("bastille-web-auth")
		if err != nil || cookie.Value != "authenticated" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
