package web

import (
	"bastille-ui/api"
	"bastille-ui/config"
	"net/http"
	"bufio"
	"strings"
)

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

func settingsPageHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.LoadConfig()

	data := PageData {
		Title: "Settings",
		Config: &cfg,
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		cfg.Username = r.FormValue("username")
		cfg.Password = r.FormValue("password")
		cfg.Address  = r.FormValue("address")
		cfg.WebPort  = r.FormValue("webPort")
		cfg.APIPort  = r.FormValue("apiPort")
		cfg.APIKey   = r.FormValue("apiKey")

		if err := config.SaveConfig(cfg); err != nil {
			data.Error = err.Error()
		} else {
			SetCredentials(cfg.Username, cfg.Password)
			SetAPIKey(cfg.APIKey)
			api.SetAPIKey(cfg.APIKey)
			data.Config = &cfg
		}
	}

	render(w, "settings", data)
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