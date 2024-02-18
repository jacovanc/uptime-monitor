package webserver

import (
	"net/http"
)

func Start() {
    http.HandleFunc("/", dashboardHandler)
    http.ListenAndServe(":8080", nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    // Serve dashboard page
}
