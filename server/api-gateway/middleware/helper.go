package middleware

import "github.com/golang-jwt/jwt/v4"

var jwtKey = []byte("key")

type Claims struct {
	jwt.RegisteredClaims
}

type contextKey string

const (
	CtxKeyID        contextKey = "id"
	CtxSessionKeyID contextKey = "sessionId"
)
