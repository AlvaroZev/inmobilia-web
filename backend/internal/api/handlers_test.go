package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/api"
	"github.com/inmobilia/inmobilia-web/backend/internal/services/ai"
)

func TestParseAIEndpoint(t *testing.T) {
	router := api.NewRouter(ai.NewMockParser(), api.ServiceInfo{AIProvider: "mock"})

	body, _ := json.Marshal(map[string]string{
		"description": "Ropero empotrado con cajones y repisas",
		"name":        "Ropero cliente",
	})
	req := httptest.NewRequest(http.MethodPost, "/ai/parse", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.Code, resp.Body.String())
	}

	var furniture map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &furniture); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if furniture["name"] != "Ropero cliente" {
		t.Fatalf("name = %v", furniture["name"])
	}
}
