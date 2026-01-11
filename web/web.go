package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"os"
)

func callBastilleAPI(path string, params map[string]string) (string, error) {

	node := getActiveNode()
	if node == nil {
		return "", fmt.Errorf("no node selected")
	}

	rawurl := fmt.Sprintf("http://%s:%s%s", node.Host, node.Port, path)
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	Key := node.Key

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return string(body), fmt.Errorf("API error: %s", resp.Status)
	}
	return string(body), nil
}

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

func Start() {

       var bindAddr string
	config := loadConfig()
	setConfig(config)

	if Host == "0.0.0.0" || Host == "localhost" || Host == "" {
		bindAddr = "0.0.0.0"
		Host = "localhost"
	} else {
	       bindAddr = Host
	}
	
	addr := fmt.Sprintf("%s:%s", bindAddr, Port)

	loadRoutes()

	log.Println("Starting BastilleBSD WebUI server on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func loadRoutes() {

	// Handle built in pages
	http.HandleFunc("/login", loginHandler)
	http.Handle("/settings", loggingMiddleware(requireLogin(settingsPageHandler)))
	http.Handle("/logout", loggingMiddleware(requireLogin(logoutHandler)))

	// Register handlers with middleware applied manually
	http.Handle("/", loggingMiddleware(requireLogin(homePageHandler)))
	http.Handle("/bastille/quickaction", loggingMiddleware(requireLogin(homePageActionHandler)))
	http.Handle("/bastille/", loggingMiddleware(requireLogin(bastilleWebHandler)))
	http.Handle("/nodes", loggingMiddleware(requireLogin(nodePageHandler)))
	http.Handle("/api/v1/node/add", loggingMiddleware(requireLogin(nodeAddHandler)))
	http.Handle("/api/v1/node/delete", loggingMiddleware(requireLogin(nodeDeleteHandler)))
	http.Handle("/api/v1/node/select", loggingMiddleware(requireLogin(nodeSelectHandler)))
	http.Handle("/api/v1/web-console", loggingMiddleware(requireLogin(consoleHandler)))

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))
}
