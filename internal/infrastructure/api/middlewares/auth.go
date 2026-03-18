package middlewares

import (
	"context"
	"net/http"
	"strings"

	authpb "github.com/fercho/school-tracking/proto/gen/auth/v1"
	"github.com/fercho/school-tracking/services/gateway/pkg/env"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

// RequireAuth is a middleware that validates JWT tokens and injects claims into the request headers and context.
func RequireAuth(cfg *env.Config, authClient authpb.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authorization header missing", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate the token via Auth gRPC Service
			resp, err := authClient.ValidateToken(r.Context(), &authpb.ValidateTokenRequest{
				AccessToken: tokenString,
			})

			if err != nil || !resp.IsValid {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			userID := resp.UserId
			role := resp.Role

			if userID == "" || role == "" {
				http.Error(w, "missing required claims in token", http.StatusUnauthorized)
				return
			}

			// Inject into request headers for downstream services to consume
			r.Header.Set("X-User-ID", userID)
			r.Header.Set("X-User-Role", role)

			// Inject into context for gateway logic if needed
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
