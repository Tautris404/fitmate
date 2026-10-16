package storage

import "errors"

var (
	ErrImageTooLarge           = errors.New("image is too large")
	ErrImageInvalid            = errors.New("invalid image")
	ErrImageFormat             = errors.New("unsupported image format")
	ErrImageDimensionsInvalid  = errors.New("invalid image dimensions")
	ErrImageDimensionsTooLarge = errors.New("image dimensions are too large")
	ErrImageCompressionFailed  = errors.New("failed to compress image below 1 MB")
)
