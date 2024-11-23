package middleware

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"order-api/pkg/jwt"
	"strings"
)

type ContextKey string

const PhoneContextKey ContextKey = "phone"

type responseWriter struct {
	http.ResponseWriter
	buffer bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.buffer.Write(b)
	return rw.ResponseWriter.Write(b)
}

type AuthResponse struct {
	Token string `json:"token"`
}

func TokenMiddleware(secret string) func(http.Handler) http.Handler {
	jwtHandler := jwt.NewJWT(secret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/auth/login" || r.URL.Path == "/auth/register" || r.URL.Path == "/auth/session" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("No Authorization header for path: %s", r.URL.Path)
				http.Error(w, "unauthorized: no token provided", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				log.Printf("Invalid token format for path: %s", r.URL.Path)
				http.Error(w, "unauthorized: invalid token format", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]

			if r.URL.Path == "/auth/session" {
				bufferedWriter := &responseWriter{ResponseWriter: w}
				next.ServeHTTP(bufferedWriter, r)
				return
			}

			valid, userData := jwtHandler.Parse(tokenString)
			log.Printf("Token validation: valid=%v, userData=%+v", valid, userData)

			if !valid || userData == nil {
				log.Printf("Invalid token for path: %s", r.URL.Path)
				http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), PhoneContextKey, userData.Phone)
			log.Printf("Added phone to context: %s for path: %s", userData.Phone, r.URL.Path)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
