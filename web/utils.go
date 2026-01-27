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
	templatePages, err := filepath.Glob(filepath.Join(webDir, "web/static/templates/*.html"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	partialsPages, err := filepath.Glob(filepath.Join(webDir, "web/static/partials/*.html"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	mainPages := []string{
        webDir + "web/static/home.html",
        webDir + "web/static/login.html",
        webDir + "web/static/settings.html",
        webDir + "web/static/nodes.html",
	webDir + "web/static/bastille-logo.png",
	}

    requestedPage := page
    if !filepath.IsAbs(requestedPage) {
        requestedPage = filepath.Join(webDir, "web/static", requestedPage+".html")
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