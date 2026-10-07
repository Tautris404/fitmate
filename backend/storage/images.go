package storage

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"

	xdraw "golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

const (
	maxUploadImageSize = 10 * 1024 * 1024
	maxImageWidth      = 1600
	maxImageHeight     = 1600
	maxSavedImageSize  = 1 * 1024 * 1024
	maxImagePixels     = 60 * 1000 * 1000
)

type ImageStorage struct {
	baseDirectory string
}

func NewImageStorage(baseDirectory string) *ImageStorage {
	return &ImageStorage{baseDirectory: baseDirectory}
}

func (s *ImageStorage) SaveUserImage(userID uint, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > maxUploadImageSize {
		return "", errors.New("image is too large")
	}

	img, err := decodeUploadedImage(fileHeader)
	if err != nil {
		return "", err
	}

	imageData, err := prepareImageForStorage(img)
	if err != nil {
		return "", err
	}

	return s.writeUserImage(userID, imageData)
}

func decodeUploadedImage(fileHeader *multipart.FileHeader) (image.Image, error) {
	if err := validateImageHeader(fileHeader); err != nil {
		return nil, err
	}

	source, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer source.Close()

	img, format, err := image.Decode(source)
	if err != nil {
		return nil, errors.New("invalid image")
	}

	if !isSupportedImageFormat(format) {
		return nil, errors.New("unsupported image format")
	}

	return img, nil
}

func validateImageHeader(fileHeader *multipart.FileHeader) error {
	source, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer source.Close()

	config, format, err := image.DecodeConfig(source)
	if err != nil {
		return errors.New("invalid image")
	}

	if !isSupportedImageFormat(format) {
		return errors.New("unsupported image format")
	}

	if config.Width <= 0 || config.Height <= 0 {
		return errors.New("invalid image dimensions")
	}

	if int64(config.Width)*int64(config.Height) > maxImagePixels {
		return errors.New("image dimensions are too large")
	}

	return nil
}

func isSupportedImageFormat(format string) bool {
	if format != "jpeg" &&
		format != "png" &&
		format != "webp" {
		return false
	}

	return true
}

func (s *ImageStorage) writeUserImage(userID uint, compressedImg []byte) (string, error) {
	userDirectory := filepath.Join(s.baseDirectory, "users", fmt.Sprintf("%d", userID))

	if err := os.MkdirAll(userDirectory, 0755); err != nil {
		return "", err
	}

	fileName, err := generateRandomFileName()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(userDirectory, fileName+".jpg")

	destination, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	if _, err = destination.Write(compressedImg); err != nil {
		return "", err
	}

	imageURL := fmt.Sprintf("/uploads/users/%d/%s.jpg", userID, fileName)
	return imageURL, nil
}

func prepareImageForStorage(img image.Image) ([]byte, error) {
	currentImg := resizeImage(img, maxImageWidth, maxImageHeight)

	for {
		for quality := 85; quality >= 60; quality -= 5 {
			var buf bytes.Buffer

			err := jpeg.Encode(&buf, currentImg, &jpeg.Options{Quality: quality})
			if err != nil {
				return nil, err
			}

			if buf.Len() <= maxSavedImageSize {
				return buf.Bytes(), nil
			}
		}

		bounds := currentImg.Bounds()

		newWidth := int(float64(bounds.Dx()) * 0.85)
		newHeight := int(float64(bounds.Dy()) * 0.85)

		if newWidth < 400 && newHeight < 400 {
			return nil, errors.New("failed to compress image below 1 MB")
		}

		currentImg = resizeImage(currentImg, newWidth, newHeight)
	}
}

func resizeImage(img image.Image, maxWidth int, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= maxWidth && height <= maxHeight {
		return img
	}

	widthRatio := float64(maxWidth) / float64(width)
	heightRatio := float64(maxHeight) / float64(height)
	ratio := min(widthRatio, heightRatio)

	newWidth := int(float64(width) * ratio)
	newHeight := int(float64(height) * ratio)

	resizedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	xdraw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, bounds, draw.Src, nil)

	return resizedImg
}

func generateRandomFileName() (string, error) {
	data := make([]byte, 16)

	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(data), nil
}
