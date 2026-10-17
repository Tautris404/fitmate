package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"fitmate/backend/models"
	"fitmate/backend/services"
)

type GymsHandler struct {
	gymsService *services.GymsService
}

func NewGymsHandler(gymsService *services.GymsService) *GymsHandler {
	return &GymsHandler{gymsService: gymsService}
}

func (h *GymsHandler) GetGyms(w http.ResponseWriter, r *http.Request) {
	gyms, err := h.gymsService.GetGyms(r.Context())
	if err != nil {
		log.Printf(logGetGymsFailed, err)
		http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(gyms); err != nil {
		log.Printf(logGetGymsEncodeFailed, err)
	}
}

func (h *GymsHandler) CreateGym(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, messageRequestTooLarge, http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, messageInvalidGymData, http.StatusBadRequest)
		return
	}

	var req models.GymRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, messageInvalidGymData, http.StatusBadRequest)
		return
	}

	result, err := h.gymsService.CreateGym(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrGymExists):
			http.Error(w, messageGymExists, http.StatusConflict)

		case isGymValidationError(err):
			http.Error(w, err.Error(), http.StatusBadRequest)

		default:
			log.Printf(logPostGymFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusCreated)

	response := models.GymResponse{
		ID:      result,
		Message: messageGymCreationSuccess,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logPostGymEncodeFailed, err)
	}
}

func isGymValidationError(err error) bool {
	return errors.Is(err, services.ErrGymNameEmpty) ||
		errors.Is(err, services.ErrGymNameTooLong) ||
		errors.Is(err, services.ErrGymAddressEmpty) ||
		errors.Is(err, services.ErrGymAddressTooLong) ||
		errors.Is(err, services.ErrGymCityEmpty) ||
		errors.Is(err, services.ErrGymCityTooLong)
}
