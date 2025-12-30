package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// --- Data structure for templates ---
type PageData struct {
	Title  string
	Output string
	Error  string
}

// --- Global API key ---
var apiKey string

func callBastilleAPI(path string, params map[string]string) (string, error) {
	base := "http://localhost:8080"
	u, _ := url.Parse(base + path)

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

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
	tmpl, err := template.ParseFiles(
		"web/static/layout.html",
		"web/static/"+page+".html",
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// --- Handlers ---
func home(w http.ResponseWriter, r *http.Request) {
	render(w, "home", PageData{Title: "Bastille Web UI"})
}

func bastilleHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: r.URL.Path,
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		// Extract subcommand from the URL
		// Example: /bastille/start -> apiPath = /api/v1/bastille/start
		subcommand := r.URL.Path[len("/bastille/"):]
		apiPath := "/api/v1/bastille/" + subcommand

		// Collect form values as params
		params := map[string]string{}
		for k, v := range r.PostForm {
			params[k] = v[0]
		}

		out, err := callBastilleAPI(apiPath, params)
		data.Output = out
		if err != nil {
			data.Error = err.Error()
		}
	}

	render(w, "result", data)
}


// --- Main function ---
func Start() {
	apiKey = "testingkey"
	if apiKey == "" {
		log.Fatal("BASTILLE_API_KEY not set")
	}

	http.HandleFunc("/", home)

	// Catch-all for any /bastille/<subcommand>
	http.HandleFunc("/bastille/", bastilleHandler)

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
}

