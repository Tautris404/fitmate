package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type accessClaims struct {
	RoleID uint `json:"role_id"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	key []byte
}

func NewAuthMiddleware(secret string) (*AuthMiddleware, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT secret must contain at least 32 bytes")
	}
	return &AuthMiddleware{key: []byte(secret)}, nil
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := strings.Fields(r.Header.Get("Authorization"))
		if len(header) != 2 || !strings.EqualFold(header[0], "Bearer") {
			unauthorized(w)
			return
		}

		claims := &accessClaims{}
		token, err := jwt.ParseWithClaims(
			header[1],
			claims,
			func(token *jwt.Token) (any, error) {
				return m.key, nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
			jwt.WithIssuedAt(),
		)
		if err != nil || !token.Valid {
			unauthorized(w)
			return
		}

		userID, err := strconv.ParseUint(claims.Subject, 10, strconv.IntSize)
		if err != nil || userID == 0 {
			unauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), authContextKey{}, authIdentity{
			UserID: uint(userID),
			RoleID: claims.RoleID,
		})
		next(w, r.WithContext(ctx))
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func requireRole(next http.HandlerFunc, allowedRoles ...uint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleID, ok := RoleIDFromContext(r.Context())
		if !ok {
			unauthorized(w)
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleID != 0 && roleID == allowedRole {
				next(w, r)
				return
			}
		}

		http.Error(w, "forbidden", http.StatusForbidden)
	}
}

func (m *AuthMiddleware) RequireAuthAndRole(next http.HandlerFunc, allowedRoles ...uint) http.HandlerFunc {
	return m.RequireAuth(requireRole(next, allowedRoles...))
}
