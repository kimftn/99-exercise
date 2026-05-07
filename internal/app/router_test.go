package app

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetListingsEndpoint(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("GET", "/listings", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("listings request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestCreatePublicUserEndpoint(t *testing.T) {
	server := NewServer()

	body := `{"name":"Dina Saputra","email":"dina@example.com","role":"seller"}`
	req := httptest.NewRequest("POST", "/public-api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("public user request failed: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["message"] != "user created" {
		t.Fatalf("unexpected message: %v", payload["message"])
	}
}
