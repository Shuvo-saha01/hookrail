package router

import (
	"fmt"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is running")
}

func HealthRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /", healthCheck)
}
