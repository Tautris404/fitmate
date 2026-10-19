package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"fitmate/backend/middleware"
	"fitmate/backend/models"
	"fitmate/backend/services"
)

type SwipesHandler struct {
	swipesService *services.SwipesService
}

func NewSwipesHandler(swipesService *services.SwipesService) *SwipesHandler {
	return &SwipesHandler{swipesService: swipesService}
}

func (h *SwipesHandler) CreateSwipe(w http.ResponseWriter, r *http.Request) {
	swiperID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || swiperID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
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
		http.Error(w, messageInvalidSwipeData, http.StatusBadRequest)
		return
	}
	var req models.SwipePostRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, messageInvalidSwipeData, jsonErrorStatus(err))
		return
	}

	matched, err := h.swipesService.CreateSwipe(r.Context(), swiperID, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidSwipeTargetID), errors.Is(err, services.ErrSwipeSelf), errors.Is(err, services.ErrSwipeActionInvalid):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, services.ErrSwipeTargetNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, services.ErrSwipeExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Printf(logPostSwipeFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusCreated)
	response := models.SwipeResponse{Message: messageSwipeCreationSuccess, Match: matched}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logPostSwipeEncodeFailed, err)
	}
}
