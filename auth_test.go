package lansenger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokenManagerGetToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/apptoken/create" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errCode": 0,
				"errMsg":  "ok",
				"data": map[string]interface{}{
					"appToken":  "test_token_123",
					"expiresIn": 7200,
				},
			})
		}
	}))
	defer server.Close()

	cfg := NewConfig("test_app", "test_secret")
	cfg.APIGatewayURL = server.URL
	httpClient := server.Client()

	tm := NewTokenManager(cfg, httpClient)
	token, err := tm.GetToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "test_token_123" {
		t.Errorf("expected token=test_token_123, got %s", token)
	}
}

func TestTokenManagerRefresh(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errCode": 0,
			"data": map[string]interface{}{
				"appToken":  fmt.Sprintf("token_%d", callCount),
				"expiresIn": 7200,
			},
		})
	}))
	defer server.Close()

	cfg := NewConfig("test_app", "test_secret")
	cfg.APIGatewayURL = server.URL
	httpClient := server.Client()

	tm := NewTokenManager(cfg, httpClient)

	token1, err := tm.GetToken(context.Background())
	if err != nil {
		t.Fatalf("first token request: %v", err)
	}

	tm.Invalidate()

	token2, err := tm.GetToken(context.Background())
	if err != nil {
		t.Fatalf("second token request after invalidation: %v", err)
	}

	if token1 == token2 {
		t.Error("expected different token after invalidation, but got same")
	}
}

func TestTokenManagerCachedToken(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errCode": 0,
			"data": map[string]interface{}{
				"appToken":  fmt.Sprintf("token_%d", callCount),
				"expiresIn": 7200,
			},
		})
	}))
	defer server.Close()

	cfg := NewConfig("test_app", "test_secret")
	cfg.APIGatewayURL = server.URL

	tm := NewTokenManager(cfg, server.Client())
	token1, _ := tm.GetToken(context.Background())
	token2, _ := tm.GetToken(context.Background())

	if callCount != 1 {
		t.Errorf("expected 1 API call for cached token, got %d", callCount)
	}
	if token1 != token2 {
		t.Errorf("expected same cached token, got %s and %s", token1, token2)
	}
}

func TestTokenManagerAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errCode": 40013,
			"errMsg":  "invalid appid",
		})
	}))
	defer server.Close()

	cfg := NewConfig("bad_app", "bad_secret")
	cfg.APIGatewayURL = server.URL

	tm := NewTokenManager(cfg, server.Client())
	_, err := tm.GetToken(context.Background())
	if err == nil {
		t.Error("expected error for API error response")
	}
}

func TestTokenManagerExternalMode(t *testing.T) {
	// External mode: when AppToken is set, GetToken returns it directly
	// without calling the API at all.
	apiCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalled = true
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errCode": 0,
			"data":    map[string]interface{}{"appToken": "api_token", "expiresIn": 7200},
		})
	}))
	defer server.Close()

	cfg := NewConfig("app", "secret")
	cfg.APIGatewayURL = server.URL
	cfg.AppToken = "external_direct_token"

	tm := NewTokenManager(cfg, server.Client())
	token, err := tm.GetToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error in external mode: %v", err)
	}
	if token != "external_direct_token" {
		t.Errorf("expected token=external_direct_token, got %s", token)
	}
	if apiCalled {
		t.Error("external mode should NOT call the token API")
	}

	// Second call should also return the same token directly
	token2, err := tm.GetToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}
	if token2 != "external_direct_token" {
		t.Errorf("expected second call to return external_direct_token, got %s", token2)
	}
}
