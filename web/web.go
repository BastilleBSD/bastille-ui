package web

import (
	"bastille-ui/config"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func callBastilleAPI(path string, params map[string]string) (string, error) {

	node := config.GetActiveNode()
	if node == nil {
		return "", fmt.Errorf("no active node selected")
	}

	rawurl := fmt.Sprintf("http://%s:%s%s", node.IP, node.Port, path)
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	apiKey := node.APIKey
	if apiKey == "" {
		apiKey = config.Config.APIKey
	}

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)

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

// --- Template renderer ---
func render(w http.ResponseWriter, page string, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles(
		"web/static/partials/default.html",
		"web/static/partials/sidebar.html",
		"web/static/partials/navbar.html",
		"web/static/partials/options.html",
		"web/static/home.html",
		"web/static/login.html",
		"web/static/settings.html",
		"web/static/nodes.html",
		"web/static/"+page+".html",
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = tmpl.ExecuteTemplate(w, "default", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func Start(addr string) {
	loadRoutes()
	log.Println("Starting BastilleBSD WebUI server on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Handle routes
func loadRoutes() {

	// Handle built in pages
	http.HandleFunc("/login", loginHandler)
	http.Handle("/settings", loggingMiddleware(requireLogin(settingsPageHandler)))
	http.Handle("/logout", loggingMiddleware(requireLogin(logoutHandler)))

	// Register handlers with middleware applied manually
	http.Handle("/", loggingMiddleware(requireLogin(homePageHandler)))
	http.Handle("/bastille/quickaction", loggingMiddleware(requireLogin(homePageActionHandler)))
	http.Handle("/bastille/", loggingMiddleware(requireLogin(bastilleWebHandler)))
	http.Handle("/nodes", loggingMiddleware(requireLogin(nodesPageHandler)))
	http.Handle("/api/v1/node/add", loggingMiddleware(requireLogin(nodeAddHandler)))
	http.Handle("/api/v1/node/delete", loggingMiddleware(requireLogin(nodeDeleteHandler)))
	http.Handle("/api/v1/node/select", loggingMiddleware(requireLogin(nodeSelectHandler)))

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))
}
