package web

import (
	"net/http"
)

func loadRoutes() {

	// Handle built in pages
	http.HandleFunc("/login", loginHandler)
	http.Handle("/settings", loggingMiddleware(requireLogin(settingsPageHandler)))
	http.Handle("/logout", loggingMiddleware(requireLogin(logoutHandler)))

	// Register handlers with middleware applied manually
	http.Handle("/", loggingMiddleware(requireLogin(homePageHandler)))

	http.Handle("/bastille/quickaction", loggingMiddleware(requireLogin(homePageActionHandler)))

	http.Handle("/bastille/", loggingMiddleware(requireLogin(bastilleWebHandler)))

	http.Handle("/nodes", loggingMiddleware(requireLogin(nodePageHandler)))
	http.Handle("/api/v1/node/add", loggingMiddleware(requireLogin(nodeAddHandler)))
	http.Handle("/api/v1/node/delete", loggingMiddleware(requireLogin(nodeDeleteHandler)))
	http.Handle("/api/v1/node/select", loggingMiddleware(requireLogin(nodeSelectHandler)))

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))
}