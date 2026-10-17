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
	mux.HandleFunc("GET /gyms/{id}", handler.GetGymByID)
	mux.HandleFunc("POST /gyms", authMiddleware.RequireAuthAndRole(handler.CreateGym, 3))
	mux.HandleFunc("PUT /gyms/{id}", authMiddleware.RequireAuthAndRole(handler.UpdateGym, 3))
	mux.HandleFunc("DELETE /gyms/{id}", authMiddleware.RequireAuthAndRole(handler.DeleteGym, 3))
}
