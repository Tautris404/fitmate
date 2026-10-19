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
	"fitmate/backend/storage"
)

type UsersHandler struct {
	usersService *services.UsersService
	authService  *services.AuthService
}

func NewUsersHandler(usersService *services.UsersService, authService *services.AuthService) *UsersHandler {
	return &UsersHandler{usersService: usersService, authService: authService}
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := createUserFromRequest(w, r, h.authService)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(models.UserResponse{ID: userID, Message: messageUserCreatedSuccess}); err != nil {
		log.Printf(logCreateUserEncodeFailed, err)
	}
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || userID <= 0 {
		http.Error(w, messageInvalidUserID, http.StatusBadRequest)
		return
	}

	if err := h.usersService.DeleteUser(r.Context(), uint(userID)); err != nil {
		switch {
		case errors.Is(err, services.ErrProfileNotFound):
			http.Error(w, messageProfileNotFound, http.StatusNotFound)
		default:
			log.Printf(logDeleteUserFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	if err := json.NewEncoder(w).Encode(models.UserResponse{ID: uint(userID), Message: messageUserDeletedSuccess}); err != nil {
		log.Printf(logDeleteUserEncodeFailed, err)
	}
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || userID <= 0 {
		http.Error(w, messageInvalidUserID, http.StatusBadRequest)
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
		http.Error(w, messageInvalidProfileData, http.StatusBadRequest)
		return
	}

	var req models.UserUpdateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, messageInvalidProfileData, http.StatusBadRequest)
		return
	}
	if err := h.usersService.UpdateUser(r.Context(), uint(userID), &req); err != nil {
		switch {
		case errors.Is(err, services.ErrProfileNotFound):
			http.Error(w, messageProfileNotFound, http.StatusNotFound)
		case isProfileUpdateValidationError(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf(logUpdateUserFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	if err := json.NewEncoder(w).Encode(models.UserResponse{ID: uint(userID), Message: messageUserUpdatedSuccess}); err != nil {
		log.Printf(logUpdateUserEncodeFailed, err)
	}
}

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || userID <= 0 {
		http.Error(w, messageInvalidUserID, http.StatusBadRequest)
		return
	}

	user, err := h.usersService.GetUserByID(r.Context(), uint(userID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProfileNotFound):
			http.Error(w, messageProfileNotFound, http.StatusNotFound)
		default:
			log.Printf(logGetUserByIDFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf(logGetUserByIDEncodeFailed, err)
	}
}

func (h *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	page, limit := 1, 10
	query := r.URL.Query()
	if query.Has("page") {
		value, err := strconv.Atoi(query.Get("page"))
		if err != nil {
			http.Error(w, services.ErrInvalidPage.Error(), http.StatusBadRequest)
			return
		}
		page = value
	}
	if query.Has("limit") {
		value, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			http.Error(w, services.ErrInvalidLimit.Error(), http.StatusBadRequest)
			return
		}
		limit = value
	}

	response, err := h.usersService.GetUsers(r.Context(), page, limit)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPage), errors.Is(err, services.ErrInvalidLimit):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf(logGetUsersFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(logGetUsersEncodeFailed, err)
	}
}

func (h *UsersHandler) DeleteMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		return
	}

	if err := h.usersService.DeleteUser(r.Context(), userID); err != nil {
		switch {
		case errors.Is(err, services.ErrProfileNotFound):
			http.Error(w, messageProfileNotFound, http.StatusNotFound)
		default:
			log.Printf(logDeleteMyProfileFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	if err := json.NewEncoder(w).Encode(models.UserResponse{ID: userID, Message: messageProfileDeletedSuccess}); err != nil {
		log.Printf(logDeleteMyProfileEncodeFailed, err)
	}
}

func (h *UsersHandler) DeleteMyImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		return
	}

	h.deleteUserImage(w, r, userID)
}

func (h *UsersHandler) DeleteUserImage(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || userID <= 0 {
		http.Error(w, messageInvalidUserID, http.StatusBadRequest)
		return
	}
	h.deleteUserImage(w, r, uint(userID))
}

func (h *UsersHandler) deleteUserImage(w http.ResponseWriter, r *http.Request, userID uint) {
	imageID, err := strconv.ParseInt(r.PathValue("imageId"), 10, 32)
	if err != nil || imageID <= 0 {
		http.Error(w, messageInvalidImageID, http.StatusBadRequest)
		return
	}

	if err := h.usersService.DeleteUserImage(r.Context(), userID, uint(imageID)); err != nil {
		switch {
		case errors.Is(err, services.ErrLastImage):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, services.ErrImageNotFound):
			http.Error(w, messageImageNotFound, http.StatusNotFound)
		default:
			log.Printf(logDeleteUserImageFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	if err := json.NewEncoder(w).Encode(models.UserResponse{ID: uint(imageID), Message: messageImageDeletedSuccess}); err != nil {
		log.Printf(logDeleteUserImageEncodeFailed, err)
	}
}

func (h *UsersHandler) AddMyImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		return
	}

	h.addUserImage(w, r, userID)
}

func (h *UsersHandler) AddUserImage(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || userID <= 0 {
		http.Error(w, messageInvalidUserID, http.StatusBadRequest)
		return
	}
	h.addUserImage(w, r, uint(userID))
}

func (h *UsersHandler) addUserImage(w http.ResponseWriter, r *http.Request, userID uint) {
	r.Body = http.MaxBytesReader(w, r.Body, 11*1024*1024)
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if r.MultipartForm != nil {
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				log.Printf(logAddUserImageCleanupFailed, err)
			}
		}()
	}
	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, messageRequestTooLarge, http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, messageInvalidMultipartForm, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) != 1 {
		http.Error(w, services.ErrSingleImageRequired.Error(), http.StatusBadRequest)
		return
	}

	imageID, err := h.usersService.AddUserImage(r.Context(), userID, files[0])
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProfileNotFound):
			http.Error(w, messageProfileNotFound, http.StatusNotFound)
		case isImageUploadValidationError(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf(logAddUserImageFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(models.UserResponse{ID: imageID, Message: messageImageAddedSuccess}); err != nil {
		log.Printf(logAddUserImageEncodeFailed, err)
	}
}

func isImageUploadValidationError(err error) bool {
	return errors.Is(err, services.ErrSingleImageRequired) ||
		errors.Is(err, services.ErrImageLimitExceeded) ||
		errors.Is(err, storage.ErrImageTooLarge) ||
		errors.Is(err, storage.ErrImageInvalid) ||
		errors.Is(err, storage.ErrImageFormat) ||
		errors.Is(err, storage.ErrImageDimensionsInvalid) ||
		errors.Is(err, storage.ErrImageDimensionsTooLarge) ||
		errors.Is(err, storage.ErrImageCompressionFailed)
}

func (h *UsersHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, messageUnauthorized, http.StatusUnauthorized)
		return
	}

	profile, err := h.usersService.GetMyProfile(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProfileNotFound):
			http.Error(w, messageProfileNotFound, http.StatusNotFound)
		default:
			log.Printf(logGetMyProfileFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		log.Printf(logGetMyProfileEncodeFailed, err)
	}
}

func (h *UsersHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
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
		http.Error(w, messageInvalidProfileData, http.StatusBadRequest)
		return
	}

	var req models.UserProfileUpdateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, messageInvalidProfileData, http.StatusBadRequest)
		return
	}

	if err := h.usersService.UpdateMyProfile(r.Context(), userID, &req); err != nil {
		switch {
		case errors.Is(err, services.ErrProfileNotFound):
			http.Error(w, messageProfileNotFound, http.StatusNotFound)
		case isProfileUpdateValidationError(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf(logUpdateMyProfileFailed, err)
			http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	if err := json.NewEncoder(w).Encode(models.UserResponse{Message: messageProfileUpdateSuccess}); err != nil {
		log.Printf(logUpdateMyProfileEncodeFailed, err)
	}
}

func isProfileUpdateValidationError(err error) bool {
	validationErrors := []error{
		services.ErrInvalidRole,
		services.ErrPasswordNotLongEnough,
		services.ErrPasswordTooLong,
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
	}
	for _, validationErr := range validationErrors {
		if errors.Is(err, validationErr) {
			return true
		}
	}
	return false
}
