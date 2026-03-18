package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authpb "github.com/fercho/school-tracking/proto/gen/auth/v1"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/clients/mocks"
	"github.com/fercho/school-tracking/services/gateway/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequireAuth(t *testing.T) {
	cfg := &env.Config{}

	tests := []struct {
		name           string
		authHeader     string
		setupMock      func(m *mocks.AuthServiceClient)
		expectedStatus int
		expectHeaders  bool
	}{
		{
			name:           "Missing Header",
			authHeader:     "",
			setupMock:      func(m *mocks.AuthServiceClient) {},
			expectedStatus: http.StatusUnauthorized,
			expectHeaders:  false,
		},
		{
			name:           "Invalid Format",
			authHeader:     "BearerToken123",
			setupMock:      func(m *mocks.AuthServiceClient) {},
			expectedStatus: http.StatusUnauthorized,
			expectHeaders:  false,
		},
		{
			name:       "Valid Token",
			authHeader: "Bearer valid-token",
			setupMock: func(m *mocks.AuthServiceClient) {
				m.On("ValidateToken", mock.Anything, &authpb.ValidateTokenRequest{
					AccessToken: "valid-token",
				}).Return(&authpb.ValidateTokenResponse{
					IsValid: true,
					UserId:  "12345",
					Role:    "admin",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectHeaders:  true,
		},
		{
			name:       "Invalid Token",
			authHeader: "Bearer invalid-token",
			setupMock: func(m *mocks.AuthServiceClient) {
				m.On("ValidateToken", mock.Anything, &authpb.ValidateTokenRequest{
					AccessToken: "invalid-token",
				}).Return(&authpb.ValidateTokenResponse{
					IsValid: false,
				}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectHeaders:  false,
		},
		{
			name:       "Service Error",
			authHeader: "Bearer error-token",
			setupMock: func(m *mocks.AuthServiceClient) {
				m.On("ValidateToken", mock.Anything, &authpb.ValidateTokenRequest{
					AccessToken: "error-token",
				}).Return((*authpb.ValidateTokenResponse)(nil), errors.New("internal error"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectHeaders:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewAuthServiceClient(t)
			tt.setupMock(mockClient)

			handler := RequireAuth(cfg, mockClient)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectHeaders {
					assert.Equal(t, "12345", r.Header.Get("X-User-ID"))
					assert.Equal(t, "admin", r.Header.Get("X-User-Role"))
					assert.Equal(t, "12345", r.Context().Value(UserIDKey))
					assert.Equal(t, "admin", r.Context().Value(UserRoleKey))
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
