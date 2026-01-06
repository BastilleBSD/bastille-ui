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
)

// --- Global API key ---
var apiKey string
var apiUrl string

func SetAPIKey(key string) {
	apiKey = key
}

func SetAPIUrl(ip, port string) {
	apiUrl = fmt.Sprintf("http://%s:%s", ip, port)
}

func callBastilleAPI(path string, params map[string]string) (string, error) {
	u, _ := url.Parse(apiUrl + path)

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
		"web/static/login.html",
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

	action := r.FormValue("action")
	target := r.FormValue("target")
	params := map[string]string{
		"target": target,
	}

	// Call the API
	callBastilleAPI("/api/v1/bastille/"+action, params)

    // Redirect back to main page so the table reloads
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// --- Home Page ---
func homePageHandler(w http.ResponseWriter, r *http.Request) {

	var jails []Jails

	data := PageData {
		Title: "Bastille WebUI",
	}
	options := ""
	item := ""
	params := map[string]string{
		"options": options,
		"item": item,
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
// Start the web server
// --- Main web.Start ---
func Start(addr string) {
	// Register handlers with middleware applied manually
	http.Handle("/", loggingMiddleware(requireLogin(homePageHandler)))
	http.Handle("/bastille/quickaction", loggingMiddleware(requireLogin(homePageActionHandler)))
	http.Handle("/bastille/", loggingMiddleware(requireLogin(bastilleWebHandler)))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))

	log.Println("Starting BastilleBSD WebUI server on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}


