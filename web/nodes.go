package web

import (
	"fmt"
	"net/http"
)

func nodeAddHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/nodes", http.StatusSeeOther)
		return
	}

	data := PageData{
		Config:     cfg,
		Nodes:      cfg.Nodes,
		ActiveNode: getActiveNode(),
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		name := r.FormValue("node")
		host := r.FormValue("host")
		port := r.FormValue("port")
		key := r.FormValue("key")

		if name == "" || host == "" || port == "" || key == "" {
			data.Error = "Please fill out all fields."
			render(w, "nodes", data)
			return
		}

		// Add node to config
		newNode := Node{
			Name:   name,
			Host:     host,
			Port:   port,
			Key: key,
		}

		cfg.Nodes = append(cfg.Nodes, newNode)

		// Save updated config
		if err := saveConfig(cfg); err != nil {
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

	filtered := []Node{}
	for _, n := range cfg.Nodes {
		if n.Name != name {
			filtered = append(filtered, n)
		}
	}

	cfg.Nodes = filtered
	_ = saveConfig(cfg)

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

	if err := setActiveNode(nodeName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	referer := r.Referer()
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func nodePageHandler(w http.ResponseWriter, r *http.Request) {

	data := PageData{
		Title: "Manage Nodes",
		Config:     cfg,
		Nodes:      cfg.Nodes,
		ActiveNode: getActiveNode(),
	}

	render(w, "nodes", data)
}

func setActiveNode(name string) error {

	activeNodeMu.Lock()
	defer activeNodeMu.Unlock()

	if cfg == nil {
		return fmt.Errorf("config not loaded")
	}

	for i := range cfg.Nodes {
		if cfg.Nodes[i].Name == name {
			activeNode = &cfg.Nodes[i]
			return nil
		}
	}

	return fmt.Errorf("node with name %s not found", name)
}

func getActiveNode() *Node {
	activeNodeMu.RLock()
	defer activeNodeMu.RUnlock()
	return activeNode
}

