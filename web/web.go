package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"bufio"
	"os/exec"
)

// --- Data structure for templates ---
type PageData struct {
	Title  string
	Output string
	Error  string
	Jails []Jails
}

type JailSettings struct {
	JID     string
	Name    string
	Boot    string
	Prio    string
	State   string
	Type    string
	IP      string
	Ports   string
	Release string
	Tags    string
}

type Jails struct {
	Jail JailSettings
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

func homePageActionHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	params := map[string]string{
    	"target": r.FormValue("target"),
    	"action": r.FormValue("action"),
	}

	// Call the API
	out, err := callBastilleAPI(apiPath, target, action)
	data.Output = out
	if err != nil {
		data.Error = err.Error()
	}

    // Redirect back to main page so the table reloads
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// --- Home Page ---
func homePageHandler(w http.ResponseWriter, r *http.Request) {

	var jails []Jails

	data := PageData {
		Title: "Bastille WebUI",
	}

	params := map[string]string{
		"item":    "",
		"options": "",
	}

	// Call API
	out, err := callBastilleAPI("/api/v1/bastille/list", params)
	if err != nil {
		data.Error = err.Error()
		render(w, "home", data)
		return
	}

	data.Output = out // optional: keep raw output

	scanner := bufio.NewScanner(strings.NewReader(out))
	firstLine := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Skip header line
		if firstLine {
			firstLine = false
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

	jails = append(jails, Jails {
		Jail: JailSettings {
			JID:     fields[0],
			Name:    fields[1],
			Boot:    fields[2],
			Prio:    fields[3],
			State:   fields[4],
			Type:    fields[5],
			IP:      fields[6],
			Ports:   fields[7],
			Release: fields[8],
			Tags:    fields[9],
		},
	})
	}

	if err := scanner.Err(); err != nil {
		data.Error = err.Error()
	}
	data.Jails = jails
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

    http.HandleFunc("/bastille/quickaction", homePageActionHandler)

	// Catch-all for any /bastille/<subcommand>
	http.HandleFunc("/bastille/", bastilleWebHandler)

	// Serve log from web/static
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))
}

