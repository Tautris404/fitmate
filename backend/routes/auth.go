package routes

import (
	"net/http"

	"fitmate/backend/handlers"
)

func RegisterAuthRoutes(mux *http.ServeMux, handler *handlers.AuthHandler) {
	mux.HandleFunc("POST /auth/register", handler.Register)
}
