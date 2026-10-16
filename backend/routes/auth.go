package routes

import (
	"net/http"

	"fitmate/backend/handlers"
	"fitmate/backend/middleware"
)

func RegisterAuthRoutes(
	mux *http.ServeMux,
	handler *handlers.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	mux.HandleFunc("POST /auth/register", handler.Register)
	mux.HandleFunc("POST /auth/login", handler.Login)

}
