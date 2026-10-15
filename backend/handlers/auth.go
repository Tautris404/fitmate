package handlers

import (
	"encoding/json"
	"errors"
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

// Handleris gauna servisą su registracijos logika.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Ribojame visą užklausą, įskaitant JSON ir nuotraukas.
	const maxRequestSize = 50 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	// Virš RAM ribos esantys failų duomenys laikomi laikinuose failuose.
	const maxMemory = 20 * 1024 * 1024
	err := r.ParseMultipartForm(maxMemory)

	// Išvalome multipart laikinus failus, bet ne išsaugotas nuotraukas.
	if r.MultipartForm != nil {
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				log.Printf("register: failed to remove temporary multipart files: %v", err)
			}
		}()
	}
	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "request is too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	// Registracijos JSON tikimės multipart lauke „data“.
	data := r.FormValue("data")
	if data == "" {
		http.Error(w, "data is required", http.StatusBadRequest)
		return
	}
	var req models.RegisterRequest
	if err := json.Unmarshal([]byte(data), &req); err != nil {
		http.Error(w, "invalid registration data", http.StatusBadRequest)
		return
	}

	// Multipart failų rodykles konvertuojame į serviso priimamą sąrašą.
	files := r.MultipartForm.File["images"]
	images := make([]multipart.FileHeader, 0, len(files))
	for _, file := range files {
		images = append(images, *file)
	}

	// Duomenų validaciją ir registraciją atlieka servisas.
	_, err = h.authService.Register(r.Context(), req, images)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailExists):
			http.Error(w, "email already exists", http.StatusConflict)
		case isRegistrationValidationError(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			// Vidinių klaidų detalių vartotojui nerodome.
			log.Printf("register: service failed: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Sėkmės atveju grąžiname 201 ir patvirtinimo pranešimą.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "user registered successfully",
	}); err != nil {
		log.Printf("register: failed to encode response: %v", err)
		return
	}
}

// Be bendros serviso validacijos žymos atpažįstame esamas klaidas.
// Nežinomos klaidos laikomos vidinėmis, kad neatskleistume jų detalių.
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
