package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is running")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", healthCheck)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server is running on localhost:8080")

	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("Server failed to run %+v", err)
	}

}
