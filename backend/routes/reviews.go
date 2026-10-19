package routes

import (
	"net/http"

	"fitmate/backend/handlers"
	"fitmate/backend/middleware"
)

func RegisterReviewsRoutes(
	mux *http.ServeMux,
	handler *handlers.ReviewsHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	mux.HandleFunc("GET /reviews/{id}", authMiddleware.RequireAuthAndRole(handler.GetReviewByID, 1, 2))
	mux.HandleFunc("POST /reviews", authMiddleware.RequireAuthAndRole(handler.CreateReview, 1))
	mux.HandleFunc("PUT /reviews/{id}", authMiddleware.RequireAuthAndRole(handler.UpdateReview, 1))
	mux.HandleFunc("DELETE /reviews/{id}", authMiddleware.RequireAuthAndRole(handler.DeleteReview, 1))
}
