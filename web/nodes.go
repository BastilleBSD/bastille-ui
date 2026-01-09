package web

import (
	"fmt"
	"bastille-ui/config"
	"net/http"
)

func nodeAddHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/nodes", http.StatusSeeOther)
		return
	}

	data := PageData{
		Config: config.Config,
		Nodes:      config.Config.Nodes,
		ActiveNode: config.GetActiveNode(),
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		name := r.FormValue("node")
		ip := r.FormValue("ip")
		port := r.FormValue("port")

		if name == "" || ip == "" || port == "" {
			data.Error = "Please fill out all fields."
			render(w, "nodes", data)
			return
		}

		// Add node to config
		newNode := config.Node{
			Name: name,
			IP:   ip,
			Port: port,
		}

		config.Config.Nodes = append(config.Config.Nodes, newNode)

		// Save updated config
		if err := config.SaveConfig(config.Config); err != nil {
			data.Error = fmt.Sprintf("Failed to save config: %v", err)
			render(w, "nodes", data)
			return
		}
	}

	http.Redirect(w, r, "/nodes", http.StatusSeeOther)
}

func nodeDeleteHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/nodes", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	name := r.FormValue("node")

	var filtered []config.Node
	for _, n := range config.Config.Nodes {
		if n.Name != name {
			filtered = append(filtered, n)
		}
	}

	config.Config.Nodes = filtered
	_ = config.SaveConfig(config.Config)

	http.Redirect(w, r, "/nodes", http.StatusSeeOther)
}

func nodeSelectHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	nodeName := r.FormValue("node")
	if nodeName == "" {
		http.Error(w, "No node selected", http.StatusBadRequest)
		return
	}

	if err := config.SetActiveNodeByName(nodeName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	referer := r.Referer()
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func nodesPageHandler(w http.ResponseWriter, r *http.Request) {

	data := PageData {
		Config: config.Config,
		Nodes: config.Config.Nodes,
		ActiveNode: config.GetActiveNode(),
	}

	render(w, "nodes", data)
}