package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"fitmate/backend/middleware"
	"fitmate/backend/models"
	"fitmate/backend/services"
)

type ReviewsHandler struct {
	reviewsService *services.ReviewsService
}

func NewReviewsHandler(reviewsService *services.ReviewsService) *ReviewsHandler {
	return &ReviewsHandler{reviewsService: reviewsService}
}

func (h *ReviewsHandler) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || reviewID <= 0 {
		http.Error(w, messageInvalidReviewID, http.StatusBadRequest)
		return
	}

	review, err := h.reviewsService.GetReviewByID(r.Context(), uint(reviewID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrReviewNotFound):
			http.Error(w, messageReviewNotFound, http.StatusNotFound)
		default:
			log.Printf(logGetReviewFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(review); err != nil {
		log.Printf(logGetReviewEncodeFailed, err)
	}
}

func (h *ReviewsHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	authorID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || authorID == 0 {
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

		http.Error(w, messageInvalidReviewData, http.StatusBadRequest)
		return
	}

	var req models.ReviewPostRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, messageInvalidReviewData, http.StatusBadRequest)
		return
	}

	reviewID, err := h.reviewsService.CreateReview(r.Context(), authorID, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidReviewAuthorID):
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		case errors.Is(err, services.ErrInvalidTrainerID), errors.Is(err, services.ErrReviewRatingInvalid):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, services.ErrTrainerNotFound):
			http.Error(w, messageTrainerNotFound, http.StatusNotFound)
		case errors.Is(err, services.ErrReviewExists):
			http.Error(w, messageReviewExists, http.StatusConflict)
		default:
			log.Printf(logPostReviewFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusCreated)

	response := models.ReviewResponse{
		ID:      reviewID,
		Message: messageReviewCreationSuccess,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logPostReviewEncodeFailed, err)
	}
}

func (h *ReviewsHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	authorID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || authorID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		return
	}

	reviewID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || reviewID <= 0 {
		http.Error(w, messageInvalidReviewID, http.StatusBadRequest)
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

		http.Error(w, messageInvalidReviewData, http.StatusBadRequest)
		return
	}

	var req models.ReviewUpdateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, messageInvalidReviewData, http.StatusBadRequest)
		return
	}

	if err := h.reviewsService.UpdateReview(r.Context(), uint(reviewID), authorID, &req); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidReviewAuthorID):
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		case errors.Is(err, services.ErrInvalidReviewID), errors.Is(err, services.ErrReviewRatingInvalid):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, services.ErrReviewNotFound):
			http.Error(w, messageReviewNotFound, http.StatusNotFound)
		default:
			log.Printf(logPutReviewFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	response := models.ReviewResponse{
		ID:      uint(reviewID),
		Message: messageReviewUpdateSuccess,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logPutReviewEncodeFailed, err)
	}
}

func (h *ReviewsHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	authorID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || authorID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		return
	}

	reviewID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || reviewID <= 0 {
		http.Error(w, messageInvalidReviewID, http.StatusBadRequest)
		return
	}

	if err := h.reviewsService.DeleteReview(r.Context(), uint(reviewID), authorID); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidReviewAuthorID):
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		case errors.Is(err, services.ErrInvalidReviewID):
			http.Error(w, messageInvalidReviewID, http.StatusBadRequest)
		case errors.Is(err, services.ErrReviewNotFound):
			http.Error(w, messageReviewNotFound, http.StatusNotFound)
		default:
			log.Printf(logDeleteReviewFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	response := models.ReviewResponse{
		ID:      uint(reviewID),
		Message: messageReviewDeletionSuccess,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logDeleteReviewEncodeFailed, err)
	}
}
