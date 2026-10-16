package main

import (
	"fitmate/backend/database"
	"fitmate/backend/handlers"
	"fitmate/backend/repositories"
	"fitmate/backend/routes"
	"fitmate/backend/services"
	"fitmate/backend/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("successfully connected to the database")
	}
	defer db.Close()

	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal(err)
	}

	authRepository := repositories.NewAuthRepository(db)
	imageStorage := storage.NewImageStorage("uploads")
	authService := services.NewAuthService(authRepository, imageStorage)
	authHandler := handlers.NewAuthHandler(authService)

	mux := http.NewServeMux()
	routes.RegisterAuthRoutes(mux, authHandler)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("server running on :8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
