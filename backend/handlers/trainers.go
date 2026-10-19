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

type TrainersHandler struct {
	trainersService *services.TrainersService
}

func NewTrainersHandler(trainersService *services.TrainersService) *TrainersHandler {
	return &TrainersHandler{trainersService: trainersService}
}

func (h *TrainersHandler) GetTrainerReviews(w http.ResponseWriter, r *http.Request) {
	trainerID, err := strconv.ParseInt(r.PathValue("trainerId"), 10, 32)
	if err != nil || trainerID <= 0 {
		http.Error(w, messageInvalidTrainerID, http.StatusUnprocessableEntity)
		return
	}

	reviews, err := h.trainersService.GetTrainerReviews(r.Context(), uint(trainerID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidTrainerID):
			http.Error(w, messageInvalidTrainerID, http.StatusUnprocessableEntity)
		default:
			log.Printf(logGetTrainerReviewsFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		log.Printf(logGetTrainerReviewsEncodeFailed, err)
	}
}

func (h *TrainersHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || trainerID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		return
	}

	reviews, err := h.trainersService.GetTrainerReviews(r.Context(), trainerID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidTrainerID):
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		default:
			log.Printf(logGetMyReviewsFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		log.Printf(logGetMyReviewsEncodeFailed, err)
	}
}
