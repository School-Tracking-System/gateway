package handlers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewReverseProxy creates a reverse proxy to the target URL.
func NewReverseProxy(target string, prefixToStrip string) (http.HandlerFunc, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify the director to strip a specific prefix if needed
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		
		// If we are proxying to a service that already mounts under /api/v1/...,
		// we might not need to strip the prefix. But if the downstream service
		// expects traffic at the root level, we could strip it here.
		// For the current setup (auth service handles /api/v1/auth), we keep the path intact.
		
		// Ensure Host header matches the target
		req.Host = targetURL.Host

		// Add gateway-specific headers
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Gateway-Proxy", "true")

		// Remove Hop-by-hop headers if necessary, but SingleHostReverseProxy handles most natively.
	}

	// Optional: You can implement proxy.ModifyResponse or proxy.ErrorHandler here for logging/custom errors.

	return func(w http.ResponseWriter, r *http.Request) {
		// Example: strip prefix manually from request path if downstream service doesn't expect it
		if prefixToStrip != "" && strings.HasPrefix(r.URL.Path, prefixToStrip) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefixToStrip)
		}

		proxy.ServeHTTP(w, r)
	}, nil
}
