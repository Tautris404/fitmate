package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

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

func (h *GymsHandler) GetGymByID(w http.ResponseWriter, r *http.Request) {
	gymID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || gymID <= 0 {
		http.Error(w, messageInvalidGymID, http.StatusUnprocessableEntity)
		return
	}

	gym, err := h.gymsService.GetGymByID(r.Context(), uint(gymID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrGymNotFound):
			http.Error(w, messageGymNotFound, http.StatusNotFound)
		default:
			log.Printf(logGetGymFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(gym); err != nil {
		log.Printf(logGetGymEncodeFailed, err)
	}
}

func (h *GymsHandler) GetGymTrainers(w http.ResponseWriter, r *http.Request) {
	gymID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || gymID <= 0 {
		http.Error(w, messageInvalidGymID, http.StatusUnprocessableEntity)
		return
	}
	page := 1
	query := r.URL.Query()
	if query.Has("page") {
		value, err := strconv.Atoi(query.Get("page"))
		if err != nil {
			http.Error(w, services.ErrInvalidPage.Error(), http.StatusUnprocessableEntity)
			return
		}
		page = value
	}

	response, err := h.gymsService.GetGymTrainers(r.Context(), uint(gymID), page)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidGymID), errors.Is(err, services.ErrInvalidPage):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, services.ErrGymNotFound):
			http.Error(w, messageGymNotFound, http.StatusNotFound)
		default:
			log.Printf(logGetGymTrainersFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logGetGymTrainersEncodeFailed, err)
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
		http.Error(w, messageInvalidGymData, jsonErrorStatus(err))
		return
	}

	result, err := h.gymsService.CreateGym(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrGymExists):
			http.Error(w, messageGymExists, http.StatusConflict)

		case isGymValidationError(err):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)

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

func (h *GymsHandler) UpdateGym(w http.ResponseWriter, r *http.Request) {
	gymID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || gymID <= 0 {
		http.Error(w, messageInvalidGymID, http.StatusUnprocessableEntity)
		return
	}

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
		http.Error(w, messageInvalidGymData, jsonErrorStatus(err))
		return
	}

	if err := h.gymsService.UpdateGym(r.Context(), uint(gymID), &req); err != nil {
		switch {
		case errors.Is(err, services.ErrGymNotFound):
			http.Error(w, messageGymNotFound, http.StatusNotFound)
		case errors.Is(err, services.ErrGymExists):
			http.Error(w, messageGymExists, http.StatusConflict)
		case isGymValidationError(err):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			log.Printf(logPutGymFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	response := models.GymResponse{
		ID:      uint(gymID),
		Message: messageGymUpdateSuccess,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logPutGymEncodeFailed, err)
	}
}

func (h *GymsHandler) DeleteGym(w http.ResponseWriter, r *http.Request) {
	gymID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || gymID <= 0 {
		http.Error(w, messageInvalidGymID, http.StatusUnprocessableEntity)
		return
	}

	if err := h.gymsService.DeleteGym(r.Context(), uint(gymID)); err != nil {
		switch {
		case errors.Is(err, services.ErrGymNotFound):
			http.Error(w, messageGymNotFound, http.StatusNotFound)
		default:
			log.Printf(logDeleteGymFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	response := models.GymResponse{
		ID:      uint(gymID),
		Message: messageGymDeletionSuccess,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logDeleteGymEncodeFailed, err)
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
