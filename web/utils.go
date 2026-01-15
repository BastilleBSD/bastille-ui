package web

import (
	"html/template"
	"net/http"
	"path/filepath"
	"os"
)

func render(w http.ResponseWriter, page string, data PageData) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Get all subcommand files
	templatePages, err := filepath.Glob("web/static/templates/*.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	partialsPages, err := filepath.Glob("web/static/partials/*.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	mainPages := []string{
        "web/static/home.html",
        "web/static/login.html",
        "web/static/settings.html",
        "web/static/nodes.html",
		"web/static/bastille-logo.png",
	}

    requestedPage := page
    if !filepath.IsAbs(requestedPage) {
        requestedPage = filepath.Join("web/static", requestedPage+".html")
    }

    // Ensure it exists
    if _, err := os.Stat(requestedPage); os.IsNotExist(err) {
        http.Error(w, "Page not found: "+requestedPage, 404)
        return
    }

	// Combine all files
	files := []string{} 
	files = append(files, templatePages...)
	files = append(files, partialsPages...)
	files = append(files, mainPages...)
	files = append(files, requestedPage)

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "default", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}