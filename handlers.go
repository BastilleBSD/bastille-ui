package main

import (
	"fmt"
	"net/http"
)

func getParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	name, release, ip := getParam(r, "name"), getParam(r, "release"), getParam(r, "ip")
	iface := getParam(r, "iface")
	if name == "" || release == "" || ip == "" {
		http.Error(w, "Missing name, release, or ip", 400)
		return
	}
	if err := BastilleCreate(name, release, ip, iface); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "✅ Created jail: %s", name)
}

func destroyHandler(w http.ResponseWriter, r *http.Request) {
	name := getParam(r, "name")
	if name == "" {
		http.Error(w, "Missing name", 400)
		return
	}
	if err := BastilleDestroy(name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "🗑️ Destroyed jail: %s", name)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	name := getParam(r, "name")
	if name == "" {
		http.Error(w, "Missing name", 400)
		return
	}
	if err := BastilleStart(name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "🚀 Started jail: %s", name)
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	name := getParam(r, "name")
	if name == "" {
		http.Error(w, "Missing name", 400)
		return
	}
	if err := BastilleStop(name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "🛑 Stopped jail: %s", name)
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	name := getParam(r, "name")
	if name == "" {
		http.Error(w, "Missing name", 400)
		return
	}
	if err := BastilleRestart(name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "🔄 Restarted jail: %s", name)
}

func renameHandler(w http.ResponseWriter, r *http.Request) {
	oldName := getParam(r, "old")
	newName := getParam(r, "new")
	if oldName == "" || newName == "" {
		http.Error(w, "Missing old or new name", 400)
		return
	}
	if err := BastilleRename(oldName, newName); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "✏️ Renamed jail: %s → %s", oldName, newName)
}

func upgradeHandler(w http.ResponseWriter, r *http.Request) {
	name := getParam(r, "name")
	if name == "" {
		http.Error(w, "Missing name", 400)
		return
	}
	if err := BastilleUpgrade(name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "⬆️ Upgraded jail: %s", name)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	output, err := BastilleList()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Your jails:\%s", output)
}
