package routes

import (
	"net/http"

	"fitmate/backend/handlers"
	"fitmate/backend/middleware"
)

func RegisterSwipesRoutes(mux *http.ServeMux, handler *handlers.SwipesHandler, authMiddleware *middleware.AuthMiddleware) {
	mux.HandleFunc("POST /swipes", authMiddleware.RequireAuthAndRole(handler.CreateSwipe, 1))
}
