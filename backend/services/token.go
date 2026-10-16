package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"os"
	"strconv"
	"time"

	"fitmate/backend/models"

	"github.com/golang-jwt/jwt/v5"
)

func createAccessToken(user models.AuthUser, expiresAt time.Time) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return "", ErrTokenGenerationFailed
	}

	now := time.Now()
	if !expiresAt.After(now) {
		return "", ErrTokenGenerationFailed
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     strconv.FormatUint(uint64(user.ID), 10),
		"role_id": user.RoleID,
		"iat":     now.Unix(),
		"exp":     expiresAt.Unix(),
	})

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", ErrTokenGenerationFailed
	}
	return signedToken, nil
}

func createRefreshToken() (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", ErrTokenGenerationFailed
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
