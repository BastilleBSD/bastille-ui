package web

import (
	"fmt"
	"bufio"
	"net/http"
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

	data := PageData{
		Title:      "Bastion",
		Config:     cfg,
		Nodes:      cfg.Nodes,
		ActiveNode: getActiveNode(),
	}

	options := ""
	item := ""
	params := map[string]string{
		"options": options,
		"item":    item,
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

		jails = append(jails, Jails{
			Jail: JailSettings{
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

	data := PageData{
		Title:      "Settings",
		Config:     cfg,
		Nodes:      cfg.Nodes,
		ActiveNode: getActiveNode(),
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		cfg.User = r.FormValue("user")
		cfg.Password = r.FormValue("password")
		cfg.Host = r.FormValue("host")
		cfg.Port = r.FormValue("port")
		if err := saveConfig(cfg); err != nil {
			data.Error = err.Error()
		} else {
			setConfig(cfg)
			data.Config = cfg
		}
	}

	render(w, "settings", data)
}

func bastilleWebHandler(w http.ResponseWriter, r *http.Request) {
	// Extract subcommand from the URL
	subcommand := r.URL.Path[len("/bastille/"):]
	subcommandpath := "templates/" + subcommand
	apiPath := "/api/v1/bastille/" + subcommand

	// Default page data
	data := PageData{
		Title:      "Bastille " + subcommand,
		Nodes:      cfg.Nodes,
		ActiveNode: getActiveNode(),
	}

	if r.Method == http.MethodGet {
        render(w, subcommandpath, data)
        return
    }

	r.ParseForm()
	params := map[string]string{}

	for key, values := range r.Form {
		if len(values) > 0 {
			if key == "options" {
				params[key] = strings.Join(values, " ")
			} else {
				params[key] = values[0]
			}
		}
	}

	if subcommand == "console" || subcommand == "top" {
		// This call returns a URL instead of normal output
		urlPath, err := callBastilleAPILive(apiPath, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the URL directly as plain text for JS iframe
		url := fmt.Sprintf("http://%s", urlPath)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(url))
		return
	}

	// Normal Bastille API call
	out, err := callBastilleAPI(apiPath, params)
	data.Output = out
	if err != nil {
		data.Error = err.Error()
	}

	render(w, subcommandpath, data)
}