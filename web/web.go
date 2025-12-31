package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
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
		"web/static/default.html",
		"web/static/sidebar.html",
		"web/static/navbar.html",
		"web/static/parse-options.html",
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

// --- Home Page ---
func homePageHandler(w http.ResponseWriter, r *http.Request) {
    // Prepare the data struct
    data := PageData{
        Title: "Bastille WebUI",
    }

    // Automatically fetch the list from the API
    params := map[string]string{
        "item":    "", // default
        "options": "",
    }

    out, err := callBastilleAPI("/api/v1/bastille/list", params)
    data.Output = out
    if err != nil {
        data.Error = err.Error()
    }

    // Render the home template
    render(w, "home", data)
}

// --- Handle Submitted Forms ---
func bastilleWebHandler(w http.ResponseWriter, r *http.Request) {
	// Extract subcommand from the URL
	// Example: /bastille/list -> subcommand = "list"
	subcommand := r.URL.Path[len("/bastille/"):]
	apiPath := "/api/v1/bastille/" + subcommand

	// Set default page header
	data := PageData{
		Title: "Bastille " + subcommand,
	}

	// Parse submitted form
	if r.Method == http.MethodPost {
		r.ParseForm()
		// Collect all form parameters dynamically
		params := map[string]string{}

		for key, values := range r.PostForm {
			if len(values) > 0 {
				// If the field is "options", join multiple values into a single string
				if key == "options" {
					params[key] = strings.Join(values, " ")
				} else {
					params[key] = values[0]
				}
			}
		}
		// Call the API
		out, err := callBastilleAPI(apiPath, params)
		data.Output = out
		if err != nil {
			data.Error = err.Error()
		}
	}

    // Render the corresponding template
    render(w, subcommand, data)
}

// --- Main function ---
func Start() {

	// Set test api key for now
	apiKey = "testingkey"
	if apiKey == "" {
		log.Fatal("BASTILLE_API_KEY not set")
	}

	// Handle home page
	http.HandleFunc("/", homePageHandler)

	// Catch-all for any /bastille/<subcommand>
	http.HandleFunc("/bastille/", bastilleWebHandler)

	// Serve log from web/static
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))
}

