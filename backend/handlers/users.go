package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"fitmate/backend/middleware"
	"fitmate/backend/services"
)

type UsersHandler struct {
	usersService *services.UsersService
}

func NewUsersHandler(usersService *services.UsersService) *UsersHandler {
	return &UsersHandler{usersService: usersService}
}

func (h *UsersHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		return
	}

	profile, err := h.usersService.GetUserProfile(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProfileNotFound):
			http.Error(w, messageProfileNotFound, http.StatusNotFound)
		default:
			log.Printf(logGetProfileFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		log.Printf(logGetProfileEncodeFailed, err)
	}
}
