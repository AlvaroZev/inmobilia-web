package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/api"
	"github.com/inmobilia/inmobilia-web/backend/internal/services/ai"
)

func TestHealthEndpoint(t *testing.T) {
	router := api.NewRouter(ai.NewMockParser(), api.ServiceInfo{AIProvider: "mock"})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.Code, resp.Body.String())
	}

	var payload map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload["status"] != "ok" {
		t.Fatalf("status = %q, want ok", payload["status"])
	}
	if payload["ai_provider"] != "mock" {
		t.Fatalf("ai_provider = %q, want mock", payload["ai_provider"])
	}
}
