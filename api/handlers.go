package api

import (
    "fmt"
    "net/http"
    "strings"
)

func getParam(r *http.Request, key string) string {
    return r.URL.Query().Get(key)
}

func BastilleBootstrapHandler(w http.ResponseWriter, r *http.Request) {

	args := []string{"bootstrap"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	arch := getParam(r, "arch")

	if options != "" {
		args = append(args, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	args = append(args, target)
	if arch != "" {
		args = append(args, arch)
	}

	output, err := BastilleCommand(args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Success: %s", output)
}

func BastilleCreateHandler(w http.ResponseWriter, r *http.Request) {

	args := []string{"create"}

	options := getParam(r, "options")
	name := getParam(r, "name")
	release := getParam(r, "release")
	ip := getParam(r, "ip")
	iface :=  getParam(r, "iface")

	if options != "" {
		args = append(args, strings.Fields(options)...)
	}
	if name == "" {
		http.Error(w, "[ERROR]: Missing name paramerter.", http.StatusBadRequest)
		return
	}
	if release == "" {
		http.Error(w, "[ERROR]: Missing release parameter", http.StatusBadRequest)
		return
	}
	if ip == "" {
		http.Error(w, "[ERROR]: Missing ip parameter", http.StatusBadRequest)
		return
	}

	args = append(args, name, release, ip)

	if iface != "" {
		args = append(args, iface)
	}

	output, err := BastilleCommand(args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Success: %s", output)
}

func BastilleDestroyHandler(w http.ResponseWriter, r *http.Request) {

	args := []string{"destroy"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		options = options + " -ay"
		args = append(args, strings.Fields(options)...)
	} else {
		options = "-ay"
		args = append(args, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}

	args = append(args, target)

	output, err := BastilleCommand(args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Success: %s", output)
}

func BastilleStartHandler(w http.ResponseWriter, r *http.Request) {

	args := []string{"start"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		args = append(args, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	args = append(args, target)

	output, err := BastilleCommand(args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Success: %s", output)
}

func BastilleStopHandler(w http.ResponseWriter, r *http.Request) {

	args := []string{"stop"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		args = append(args, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	args = append(args, target)

	output, err := BastilleCommand(args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Success: %s", output)
}

func BastilleRenameHandler(w http.ResponseWriter, r *http.Request) {

	args := []string{"rename"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	new_name := getParam(r, "new_name")

	if options != "" {
		options = options + " -a"
		args = append(args, strings.Fields(options)...)
	} else {
		options = "-a"
		args = append(args, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	if new_name == "" {
		http.Error(w, "[ERROR]: Missing new_name parameter", http.StatusBadRequest)
		return
	}
	args = append(args, target, new_name)

	output, err := BastilleCommand(args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Success: %s", output)
}

func BastilleRestartHandler(w http.ResponseWriter, r *http.Request) {

	args := []string{"restart"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		args = append(args, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	args = append(args, target)

	output, err := BastilleCommand(args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Success: %s", output)
}
