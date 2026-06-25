package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChaosScenariosEndpoint(t *testing.T) {
	mux := newMux(Config{AuthToken: "test-token"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/chaos", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var response struct {
		Scenarios []ChaosScenario `json:"scenarios"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(response.Scenarios) != 4 {
		t.Fatalf("len(scenarios) = %d, want 4", len(response.Scenarios))
	}
}

func TestAuthedEndpointRejectsMissingTokenWithJSON(t *testing.T) {
	mux := newMux(Config{AuthToken: "test-token"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/chaos", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rec.Body.String())
	}
	if response["error"] != "unauthorized" {
		t.Fatalf("error = %q, want unauthorized", response["error"])
	}
}
