package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-api/pkg/jwt"
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
			bufferedWriter := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(bufferedWriter, r)

			if bufferedWriter.buffer.Len() > 0 {
				var response AuthResponse

				if err := json.NewDecoder(&bufferedWriter.buffer).Decode(&response); err == nil {
					if response.Token != "" {
						if valid, data := jwtHandler.Parse(response.Token); valid && data != nil {
							ctx := context.WithValue(r.Context(), PhoneContextKey, data.Phone)
							r = r.WithContext(ctx)
							fmt.Printf("Phone from token: %s\n", data.Phone)
						}
					}
				}
			}
		})
	}
}
