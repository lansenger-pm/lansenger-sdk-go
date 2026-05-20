package lansenger

import (
	"context"
	"net/url"
	"testing"
)

func TestBuildAuthorizeURL(t *testing.T) {
	c := NewClient("myapp", "secret")
	c.config.PassportURL = "https://passport.example.com"

	authURL := c.BuildAuthorizeURL("https://myapp.example.com/callback", "", "state123")
	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("failed to parse authorize URL: %v", err)
	}

	if parsed.Host != "passport.example.com" {
		t.Errorf("expected host passport.example.com, got %s", parsed.Host)
	}
	if parsed.Path != "/oauth2/authorize" {
		t.Errorf("expected path /oauth2/authorize, got %s", parsed.Path)
	}
	if parsed.Query().Get("appid") != "myapp" {
		t.Errorf("expected appid=myapp, got %s", parsed.Query().Get("appid"))
	}
	if parsed.Query().Get("response_type") != "code" {
		t.Errorf("expected response_type=code, got %s", parsed.Query().Get("response_type"))
	}
	if parsed.Query().Get("scope") != OAuth2ScopeBasicUserInfo {
		t.Errorf("expected scope=%s, got %s", OAuth2ScopeBasicUserInfo, parsed.Query().Get("scope"))
	}
	if parsed.Query().Get("state") != "state123" {
		t.Errorf("expected state=state123, got %s", parsed.Query().Get("state"))
	}
	if parsed.Query().Get("redirect_uri") != "https://myapp.example.com/callback" {
		t.Errorf("expected redirect_uri, got %s", parsed.Query().Get("redirect_uri"))
	}
}

func TestBuildAuthorizeURLCustomScope(t *testing.T) {
	c := NewClient("myapp", "secret")
	c.config.PassportURL = "https://passport.example.com"

	authURL := c.BuildAuthorizeURL("https://callback.example.com", "custom_scope", "")
	parsed, _ := url.Parse(authURL)
	if parsed.Query().Get("scope") != "custom_scope" {
		t.Errorf("expected scope=custom_scope, got %s", parsed.Query().Get("scope"))
	}
}

func TestParseAuthorizeCallback(t *testing.T) {
	result, err := ParseAuthorizeCallback("code=abc123&state=state456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["code"] != "abc123" {
		t.Errorf("expected code=abc123, got %s", result["code"])
	}
	if result["state"] != "state456" {
		t.Errorf("expected state=state456, got %s", result["state"])
	}
}

func TestParseAuthorizeCallbackWithError(t *testing.T) {
	result, err := ParseAuthorizeCallback("error=access_denied&error_description=user+refused")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["error"] != "access_denied" {
		t.Errorf("expected error=access_denied, got %s", result["error"])
	}
}

func TestValidateCallbackState(t *testing.T) {
	if !ValidateCallbackState("state123", "state123") {
		t.Error("expected state match to return true")
	}
	if ValidateCallbackState("state123", "state456") {
		t.Error("expected state mismatch to return false")
	}
}

func TestExchangeCode(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/user_token/create", 0, "ok", map[string]interface{}{
			"userToken":        "utok123",
			"expiresIn":        7200,
			"refreshToken":     "rtok123",
			"refreshExpiresIn": 2592000,
			"staffId":          "s001",
			"scope":            "basic_userinfor",
			"state":            "state456",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.ExchangeCode(context.Background(), "code123", "https://callback.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.UserToken != "utok123" {
		t.Errorf("expected UserToken=utok123, got %s", result.UserToken)
	}
	if result.RefreshToken != "rtok123" {
		t.Errorf("expected RefreshToken=rtok123, got %s", result.RefreshToken)
	}
	if result.StaffID != "s001" {
		t.Errorf("expected StaffID=s001, got %s", result.StaffID)
	}
	if result.ExpiresIn != 7200 {
		t.Errorf("expected ExpiresIn=7200, got %d", result.ExpiresIn)
	}
}

func TestRefreshUserToken(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/refresh_token/create", 0, "ok", map[string]interface{}{
			"userToken":        "utok_new",
			"expiresIn":        7200,
			"refreshToken":     "rtok_new",
			"refreshExpiresIn": 2592000,
			"staffId":          "s001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.RefreshUserToken(context.Background(), "rtok123", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.UserToken != "utok_new" {
		t.Errorf("expected UserToken=utok_new, got %s", result.UserToken)
	}
}
