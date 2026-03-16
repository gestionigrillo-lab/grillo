package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	DeviceTokenKey contextKey = "device_token"
	AdminUserKey   contextKey = "admin_user"
)

func AdminAuth(adminToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			token := parts[1]
			if token != adminToken {
				http.Error(w, `{"error":"invalid admin token"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), AdminUserKey, "admin")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func DeviceAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Device-Token")
			if token == "" {
				http.Error(w, `{"error":"missing device token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), DeviceTokenKey, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
