package handlers

const (
	logGetGymsFailed         = "get gyms: service failed: %v"
	logGetGymsEncodeFailed   = "get gyms: failed to encode response: %v"
	logGetGymFailed          = "get gym: service failed: %v"
	logGetGymEncodeFailed    = "get gym: failed to encode response: %v"
	logPostGymFailed         = "post gym: service failed: %v"
	logPostGymEncodeFailed   = "post gym: failed to encode response: %v"
	logPutGymFailed          = "put gym: service failed: %v"
	logPutGymEncodeFailed    = "put gym: failed to encode response: %v"
	logDeleteGymFailed       = "delete gym: service failed: %v"
	logDeleteGymEncodeFailed = "delete gym: failed to encode response: %v"
)

const (
	messageGymCreationSuccess = "gym created successfully"
	messageGymUpdateSuccess   = "updated gym data"
	messageGymDeletionSuccess = "gym deleted successfully"
	messageInvalidGymData     = "invalid gym data"
	messageGymExists          = "gym already exists"
	messageInvalidGymID       = "invalid gym id"
	messageGymNotFound        = "gym not found"
)
