package main

import (
	"fitmate/backend/database"
	"fitmate/backend/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to the database successfully")
	}
	defer db.Close()

	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", handlers.Register(db))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
