package api

import (
	"os"
	"strings"
)

func corsAllowedOrigin() string {
	if value := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGIN")); value != "" {
		return value
	}
	return "*"
}

func trustedProxies() []string {
	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	if raw == "" {
		return []string{"127.0.0.1", "::1"}
	}
	if raw == "none" {
		return nil
	}

	parts := strings.Split(raw, ",")
	proxies := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			proxies = append(proxies, trimmed)
		}
	}
	return proxies
}
