package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// AdminAuth validates the admin bearer token.
func AdminAuth(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token == "" {
				http.Error(w, `{"error":"admin token not configured"}`, http.StatusInternalServerError)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"error":"missing authorization"}`, http.StatusUnauthorized)
				return
			}

			provided := strings.TrimPrefix(auth, "Bearer ")
			if subtle.ConstantTimeCompare([]byte(provided), []byte(token)) != 1 {
				http.Error(w, `{"error":"invalid token"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// DeviceAuth validates device bearer token (device_id:token format stored in Redis).
func DeviceAuth(validateFn func(token string) (string, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"error":"missing authorization"}`, http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			deviceID, err := validateFn(token)
			if err != nil {
				http.Error(w, `{"error":"invalid device token"}`, http.StatusForbidden)
				return
			}

			r.Header.Set("X-Device-ID", deviceID)
			next.ServeHTTP(w, r)
		})
	}
}
