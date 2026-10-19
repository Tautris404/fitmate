package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"fitmate/backend/middleware"
	"fitmate/backend/services"
)

type AthletesHandler struct {
	athletesService *services.AthletesService
}

func NewAthletesHandler(athletesService *services.AthletesService) *AthletesHandler {
	return &AthletesHandler{athletesService: athletesService}
}

func (h *AthletesHandler) GetAthletes(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()

	var afterID uint
	if query.Has("after_id") {
		value, err := strconv.ParseUint(query.Get("after_id"), 10, 31)
		if err != nil {
			http.Error(w, services.ErrInvalidAthleteAfterID.Error(), http.StatusUnprocessableEntity)
			return
		}
		afterID = uint(value)
	}

	response, err := h.athletesService.GetAthletes(r.Context(), userID, query.Get("seed"), afterID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAthleteSeed), errors.Is(err, services.ErrInvalidAthleteAfterID):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, services.ErrProfileNotFound):
			http.Error(w, messageProfileNotFound, http.StatusNotFound)
		default:
			log.Printf(logGetAthletesFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logGetAthletesEncodeFailed, err)
	}
}
