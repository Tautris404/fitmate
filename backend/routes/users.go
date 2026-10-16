package routes

import (
	"net/http"

	"fitmate/backend/handlers"
	"fitmate/backend/middleware"
)

func RegisterUsersRoutes(
	mux *http.ServeMux,
	handler *handlers.UsersHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	mux.HandleFunc("GET /users/me", authMiddleware.RequireAuthAndRole(handler.GetUserProfile, 1, 2))
}
