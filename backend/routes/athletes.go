package routes

import (
	"net/http"

	"fitmate/backend/handlers"
	"fitmate/backend/middleware"
)

func RegisterAthletesRoutes(mux *http.ServeMux, handler *handlers.AthletesHandler, authMiddleware *middleware.AuthMiddleware) {
	mux.HandleFunc("GET /athletes", authMiddleware.RequireAuthAndRole(handler.GetAthletes, 1))
}
