package web

import (
	"html/template"
	"net/http"
	"time"
)

// --- Web login credentials from config ---
var webUser, webPass string

func SetCredentials(user, pass string) {
	webUser = user
	webPass = pass
}

// --- Show login page ---
func loginHandler(w http.ResponseWriter, r *http.Request) {

	data := PageData{}

	if r.Method == http.MethodPost {
		r.ParseForm()
		user := r.FormValue("username")
		pass := r.FormValue("password")

		if user == webUser && pass == webPass {
			// Login success: set session cookie
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

	tmpl := "web/static/login.html"
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Template ERROR: "+err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// --- Logout handler ---
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "bastille-web-auth",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// --- Middleware to enforce login ---
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
