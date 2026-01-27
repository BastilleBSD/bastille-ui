package web

import (
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
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
	_, err := callBastilleAPI("/api/v1/bastille/"+action, params)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string {
			"error": err.Error(),
		})
        	return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
    data := PageData{
        Title:      "Bastion",
        Config:     cfg,
        Nodes:      cfg.Nodes,
        ActiveNode: getActiveNode(),
    }

    // Simply render the home template
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

	subcommand := r.URL.Path[len("/bastille/"):]
	subcommandpath := "templates/" + subcommand
	apiPath := "/api/v1/bastille/" + subcommand
	apiPathLive := "/api/v1/bastille/live/" + subcommand

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
		urlPath, err := callBastilleAPILive(apiPathLive, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("%s", urlPath)
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

func getJailsJSONHandler(w http.ResponseWriter, r *http.Request) {

    // Call the API, expecting it to return JSON
    out, err := callBastilleAPI("/api/v1/bastille/list?options=--json", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the content type for JSON
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)

    // Write the API output directly to the response
    w.Write([]byte(out))
}
