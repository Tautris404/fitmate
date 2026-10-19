package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"fitmate/backend/models"
	"fitmate/backend/services"
	"fitmate/backend/storage"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if _, ok := createUserFromRequest(w, r, h.authService); !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(models.UserResponse{Message: messageRegistrationSuccess}); err != nil {
		log.Printf(logRegisterEncodeFailed, err)
	}
}

func createUserFromRequest(w http.ResponseWriter, r *http.Request, authService *services.AuthService) (uint, bool) {
	const maxRequestSize = 50 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	const maxMemory = 20 * 1024 * 1024
	err := r.ParseMultipartForm(maxMemory)

	if r.MultipartForm != nil {
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				log.Printf(logRegisterCleanupFailed, err)
			}
		}()
	}
	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, messageRequestTooLarge, http.StatusRequestEntityTooLarge)
			return 0, false
		}
		http.Error(w, messageInvalidMultipartForm, http.StatusBadRequest)
		return 0, false
	}

	data := r.FormValue("data")
	if data == "" {
		http.Error(w, messageDataRequired, http.StatusBadRequest)
		return 0, false
	}
	var req models.RegisterRequest
	if err := json.Unmarshal([]byte(data), &req); err != nil {
		http.Error(w, messageInvalidRegistrationData, http.StatusBadRequest)
		return 0, false
	}

	files := r.MultipartForm.File["images"]
	images := make([]multipart.FileHeader, 0, len(files))
	for _, file := range files {
		images = append(images, *file)
	}

	userID, err := authService.Register(r.Context(), req, images)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailExists):
			http.Error(w, messageEmailExists, http.StatusConflict)
		case isRegistrationValidationError(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf(logRegisterServiceFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return 0, false
	}

	return userID, true
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, messageRequestTooLarge, http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, messageInvalidLoginData, http.StatusBadRequest)
		return
	}

	var req models.LoginRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, messageInvalidLoginData, http.StatusBadRequest)
		return
	}

	result, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			http.Error(w, messageInvalidCredentials, http.StatusUnauthorized)

		case isLoginValidationError(err):
			http.Error(w, err.Error(), http.StatusBadRequest)

		default:
			log.Printf(logLoginServiceFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf(logLoginEncodeFailed, err)
	}
}

func isLoginValidationError(err error) bool {
	return errors.Is(err, services.ErrEmailEmpty) ||
		errors.Is(err, services.ErrEmailTooLong) ||
		errors.Is(err, services.ErrEmailInvalid) ||
		errors.Is(err, services.ErrPasswordNotLongEnough)
}

func isRegistrationValidationError(err error) bool {
	validationErrors := []error{
		services.ErrInvalidRole,
		services.ErrEmailEmpty,
		services.ErrEmailTooLong,
		services.ErrEmailInvalid,
		services.ErrPasswordNotLongEnough,
		services.ErrFirstNameEmpty,
		services.ErrFirstNameTooLong,
		services.ErrGenderInvalid,
		services.ErrGenderEmpty,
		services.ErrBioTooLong,
		services.ErrGymIDsInvalid,
		services.ErrWeekdayInvalid,
		services.ErrStartTimeInvalid,
		services.ErrEndTimeInvalid,
		services.ErrStartTimeAfterEndTime,
		services.ErrWorkoutTimesOverlap,
		services.ErrBirthDateInvalid,
		services.ErrAthleteProfileEmpty,
		services.ErrTrainerProfileEmpty,
		services.ErrTrainerProfileNotAllowed,
		services.ErrAthleteProfileNotAllowed,
		services.ErrEducationEmpty,
		services.ErrEducationTooLong,
		services.ErrSpecializationIDsInvalid,
		services.ErrHeightInvalid,
		services.ErrWeightInvalid,
		services.ErrExperienceYearsInvalid,
		services.ErrImageCountInvalid,
		storage.ErrImageTooLarge,
		storage.ErrImageInvalid,
		storage.ErrImageFormat,
		storage.ErrImageDimensionsInvalid,
		storage.ErrImageDimensionsTooLarge,
		storage.ErrImageCompressionFailed,
	}
	for _, validationErr := range validationErrors {
		if errors.Is(err, validationErr) {
			return true
		}
	}
	return false
}
