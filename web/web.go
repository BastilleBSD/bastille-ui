package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// --- Data structure for templates ---
type PageData struct {
	Title  string
	Output string
	Error  string
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
		"web/static/layout.html",
		"web/static/"+page+".html",
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// --- Handlers ---
func home(w http.ResponseWriter, r *http.Request) {
	render(w, "home", PageData{Title: "Bastille Web UI"})
}

func startJail(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Start Jail"}
	if r.Method == http.MethodPost {
		jail := r.FormValue("jail")
		out, err := callBastilleAPI("/bastille/start", map[string]string{
			"target": jail,
		})
		data.Output = out
		if err != nil {
			data.Error = err.Error()
		}
	}
	render(w, "start", data)
}

func stopJail(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Stop Jail"}
	if r.Method == http.MethodPost {
		jail := r.FormValue("jail")
		out, err := callBastilleAPI("/bastille/stop", map[string]string{
			"target": jail,
		})
		data.Output = out
		if err != nil {
			data.Error = err.Error()
		}
	}
	render(w, "stop", data)
}

func createJail(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Create Jail"}

	if r.Method == http.MethodPost {
		// Get form values
		jail := r.FormValue("jail")
		release := r.FormValue("release")
		ip := r.FormValue("ip")
		iface := r.FormValue("iface")     // optional
		options := r.FormValue("options") // checkbox options as "-B -V"

		// Build query parameters for the API
		params := map[string]string{
			"name":    jail,
			"release": release,
			"ip":      ip,
		}

		if options != "" {
			params["options"] = options
		}
		if iface != "" {
			params["iface"] = iface
		}

		// Call API
		out, err := callBastilleAPI("/bastille/create", params)
		data.Output = out
		if err != nil {
			data.Error = err.Error()
		}
	}

	render(w, "create", data)
}


// --- Main function ---
func Start() {
	apiKey = "testingkey" // or os.Getenv("BASTILLE_API_KEY")
	if apiKey == "" {
		log.Fatal("BASTILLE_API_KEY not set")
	}
	http.HandleFunc("/", home)
	http.HandleFunc("/start", startJail)
	http.HandleFunc("/stop", stopJail)
	http.HandleFunc("/create", createJail)

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
}
