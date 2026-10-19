package routes

import (
	"net/http"

	"fitmate/backend/handlers"
	"fitmate/backend/middleware"
)

func RegisterTrainersRoutes(mux *http.ServeMux, handler *handlers.TrainersHandler, authMiddleware *middleware.AuthMiddleware) {
	mux.HandleFunc("GET /trainers/reviews/me", authMiddleware.RequireAuthAndRole(handler.GetMyReviews, 2))
	mux.HandleFunc("GET /trainers/{trainerId}/reviews", authMiddleware.RequireAuthAndRole(handler.GetTrainerReviews, 1, 2))
}
