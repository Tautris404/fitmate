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
	mux.HandleFunc("POST /users", authMiddleware.RequireAuthAndRole(handler.CreateUser, 3))
	mux.HandleFunc("GET /users", authMiddleware.RequireAuthAndRole(handler.GetUsers, 3))
	mux.HandleFunc("GET /users/{id}", authMiddleware.RequireAuthAndRole(handler.GetUserByID, 3))
	mux.HandleFunc("PUT /users/{id}", authMiddleware.RequireAuthAndRole(handler.UpdateUser, 3))
	mux.HandleFunc("DELETE /users/{id}", authMiddleware.RequireAuthAndRole(handler.DeleteUser, 3))
	mux.HandleFunc("POST /users/{id}/images", authMiddleware.RequireAuthAndRole(handler.AddUserImage, 3))
	mux.HandleFunc("DELETE /users/{id}/images/{imageId}", authMiddleware.RequireAuthAndRole(handler.DeleteUserImage, 3))
	mux.HandleFunc("GET /users/me", authMiddleware.RequireAuthAndRole(handler.GetMyProfile, 1, 2))
	mux.HandleFunc("PUT /users/me", authMiddleware.RequireAuthAndRole(handler.UpdateMyProfile, 1, 2))
	mux.HandleFunc("DELETE /users/me", authMiddleware.RequireAuthAndRole(handler.DeleteMyProfile, 1, 2))
	mux.HandleFunc("POST /users/me/images", authMiddleware.RequireAuthAndRole(handler.AddMyImage, 1, 2))
	mux.HandleFunc("DELETE /users/me/images/{imageId}", authMiddleware.RequireAuthAndRole(handler.DeleteMyImage, 1, 2))
}
