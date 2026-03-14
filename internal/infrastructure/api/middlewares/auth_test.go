package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fercho/school-tracking/services/gateway/pkg/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestRequireAuth(t *testing.T) {
	cfg := &env.Config{
		JWTSecret: "test-secret",
	}

	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "12345",
		"role": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	validTokenString, _ := validToken.SignedString([]byte(cfg.JWTSecret))

	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "12345",
		"role": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	invalidTokenString, _ := invalidToken.SignedString([]byte("wrong-secret"))

	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "12345",
		"role": "admin",
		"exp":  time.Now().Add(-time.Hour).Unix(),
	})
	expiredTokenString, _ := expiredToken.SignedString([]byte(cfg.JWTSecret))

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectHeaders  bool
	}{
		{
			name:           "Valid Token",
			authHeader:     "Bearer " + validTokenString,
			expectedStatus: http.StatusOK,
			expectHeaders:  true,
		},
		{
			name:           "Missing Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectHeaders:  false,
		},
		{
			name:           "Invalid Format",
			authHeader:     "Basic some-token",
			expectedStatus: http.StatusUnauthorized,
			expectHeaders:  false,
		},
		{
			name:           "Invalid Signature",
			authHeader:     "Bearer " + invalidTokenString,
			expectedStatus: http.StatusUnauthorized,
			expectHeaders:  false,
		},
		{
			name:           "Expired Token",
			authHeader:     "Bearer " + expiredTokenString,
			expectedStatus: http.StatusUnauthorized,
			expectHeaders:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequireAuth(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectHeaders {
					assert.Equal(t, "12345", r.Header.Get("X-User-ID"))
					assert.Equal(t, "admin", r.Header.Get("X-User-Role"))
					
					ctxUserID := r.Context().Value(UserIDKey)
					ctxUserRole := r.Context().Value(UserRoleKey)
					assert.Equal(t, "12345", ctxUserID)
					assert.Equal(t, "admin", ctxUserRole)
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
