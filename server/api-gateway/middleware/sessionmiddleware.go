package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionHeader := r.Header.Get("Session")
		if sessionHeader == "" {
			http.Error(w, "Session header requried", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(sessionHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			log.Println("validation error: ", err)
			http.Error(w, "Invalid session token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), CtxSessionKeyID, claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
