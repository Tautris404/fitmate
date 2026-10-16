package middleware

import "context"

type authContextKey struct{}

type authIdentity struct {
	UserID uint
	RoleID uint
}

func UserIDFromContext(ctx context.Context) (uint, bool) {
	identity, ok := ctx.Value(authContextKey{}).(authIdentity)
	return identity.UserID, ok
}

func RoleIDFromContext(ctx context.Context) (uint, bool) {
	identity, ok := ctx.Value(authContextKey{}).(authIdentity)
	return identity.RoleID, ok
}
