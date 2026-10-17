package handlers

const (
	logGetTrainerReviewsFailed       = "get trainer reviews: service failed: %v"
	logGetTrainerReviewsEncodeFailed = "get trainer reviews: failed to encode response: %v"
	logGetMyReviewsFailed            = "get my reviews: service failed: %v"
	logGetMyReviewsEncodeFailed      = "get my reviews: failed to encode response: %v"
	logGetReviewFailed               = "get review: service failed: %v"
	logGetReviewEncodeFailed         = "get review: failed to encode response: %v"
	logPostReviewFailed              = "post review: service failed: %v"
	logPostReviewEncodeFailed        = "post review: failed to encode response: %v"
	logPutReviewFailed               = "put review: service failed: %v"
	logPutReviewEncodeFailed         = "put review: failed to encode response: %v"
	logDeleteReviewFailed            = "delete review: service failed: %v"
	logDeleteReviewEncodeFailed      = "delete review: failed to encode response: %v"
)

const (
	messageInvalidTrainerID      = "invalid trainer id"
	messageInvalidReviewID       = "invalid review id"
	messageReviewNotFound        = "review not found"
	messageInvalidReviewData     = "invalid review data"
	messageTrainerNotFound       = "trainer not found"
	messageReviewExists          = "review already exists for this trainer"
	messageReviewCreationSuccess = "review created successfully"
	messageReviewUpdateSuccess   = "review updated successfully"
	messageReviewDeletionSuccess = "review deleted successfully"
)
