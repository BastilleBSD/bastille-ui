package web

import (
	"log"
	"net/http"
	"os/exec"
	"fmt"
	"strings"
)

func consoleHandler(w http.ResponseWriter, r *http.Request) {

	host := r.Host[:strings.Index(r.Host, ":")]
	jail := r.FormValue("target")
	if jail == "" {
		http.Error(w, "Missing target parameter", http.StatusBadRequest)
		return
	}

	// Run ttyd to serve the jail console (using bastille as the command)
	cmd := exec.Command("ttyd", "-t", "disableLeaveAlert=true", "-o", "-p", "8182", "-W", "bastille", "console", jail)

	// Start the ttyd process
	if err := cmd.Start(); err != nil {
		log.Println("Error starting ttyd:", err)
		http.Error(w, fmt.Sprintf("Error starting ttyd: %s", err), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("http://%s:%s/", host, "8182")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(url))
}