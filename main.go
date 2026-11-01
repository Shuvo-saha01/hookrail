package main

import (
	"hookrail/db"
	"hookrail/router"
	"log"
	"net/http"
	"time"
)

func main() {
	pool, err := db.ConnectToDb()
	defer pool.Close()

	mux := http.NewServeMux()

	router.HealthRouter(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server is running on localhost:8080")

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Server failed to run %+v", err)
	}

}
