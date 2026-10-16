package handlers

const (
	messageRequestTooLarge         = "request is too large"
	messageInvalidMultipartForm    = "invalid multipart form"
	messageDataRequired            = "data is required"
	messageInvalidRegistrationData = "invalid registration data"
	messageEmailExists             = "email already exists"
	messageInternalServerError     = "internal server error"
	messageInvalidLoginData        = "invalid login data"
	messageInvalidCredentials      = "invalid email or password"
)

const (
	messageRegistrationSuccess = "user registered successfully"
	logRegisterCleanupFailed   = "register: failed to remove temporary multipart files: %v"
	logRegisterServiceFailed   = "register: service failed: %v"
	logRegisterEncodeFailed    = "register: failed to encode response: %v"
	logLoginServiceFailed      = "login: service failed: %v"
	logLoginEncodeFailed       = "login: failed to encode response: %v"
)
