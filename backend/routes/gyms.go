package routes

import (
	"net/http"

	"fitmate/backend/handlers"
	"fitmate/backend/middleware"
)

func RegisterGymsRoutes(
	mux *http.ServeMux,
	handler *handlers.GymsHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	mux.HandleFunc("GET /gyms", handler.GetGyms)
	mux.HandleFunc("POST /gyms", authMiddleware.RequireAuthAndRole(handler.CreateGym, 3))
}
