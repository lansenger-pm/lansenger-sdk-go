package lansenger

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig("app1", "secret1")
	if cfg.AppID != "app1" {
		t.Errorf("expected AppID=app1, got %s", cfg.AppID)
	}
	if cfg.AppSecret != "secret1" {
		t.Errorf("expected AppSecret=secret1, got %s", cfg.AppSecret)
	}
	if cfg.APIGatewayURL != DefaultAPIGatewayURL {
		t.Errorf("expected default APIGatewayURL=%s, got %s", DefaultAPIGatewayURL, cfg.APIGatewayURL)
	}
	if cfg.HTTPTimeout != 30.0 {
		t.Errorf("expected default HTTPTimeout=30.0, got %f", cfg.HTTPTimeout)
	}
}

func TestConfigIsConfigured(t *testing.T) {
	cfg := NewConfig("app1", "secret1")
	if !cfg.IsConfigured() {
		t.Error("expected IsConfigured=true with valid credentials")
	}

	cfg2 := NewConfig("", "")
	if cfg2.IsConfigured() {
		t.Error("expected IsConfigured=false with empty credentials")
	}
}

func TestConfigHasPassportURL(t *testing.T) {
	cfg := NewConfig("app1", "secret1")
	if !cfg.HasPassportURL() {
		t.Error("expected HasPassportURL=true with default passport URL")
	}

	cfg.PassportURL = "https://passport.example.com"
	if !cfg.HasPassportURL() {
		t.Error("expected HasPassportURL=true with passport URL set")
	}
}

func TestConfigFromEnv(t *testing.T) {
	os.Setenv("LANSENGER_APP_ID", "env_app")
	os.Setenv("LANSENGER_APP_SECRET", "env_secret")
	os.Setenv("LANSENGER_API_GATEWAY_URL", "https://custom.example.com")
	os.Setenv("LANSENGER_PASSPORT_URL", "https://passport.example.com")
	os.Setenv("LANSENGER_HTTP_TIMEOUT", "60")
	defer func() {
		os.Unsetenv("LANSENGER_APP_ID")
		os.Unsetenv("LANSENGER_APP_SECRET")
		os.Unsetenv("LANSENGER_API_GATEWAY_URL")
		os.Unsetenv("LANSENGER_PASSPORT_URL")
		os.Unsetenv("LANSENGER_HTTP_TIMEOUT")
	}()

	cfg, err := ConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AppID != "env_app" {
		t.Errorf("expected AppID=env_app, got %s", cfg.AppID)
	}
	if cfg.AppSecret != "env_secret" {
		t.Errorf("expected AppSecret=env_secret, got %s", cfg.AppSecret)
	}
	if cfg.APIGatewayURL != "https://custom.example.com" {
		t.Errorf("expected custom gateway URL, got %s", cfg.APIGatewayURL)
	}
	if cfg.PassportURL != "https://passport.example.com" {
		t.Errorf("expected passport URL, got %s", cfg.PassportURL)
	}
	if cfg.HTTPTimeout != 60.0 {
		t.Errorf("expected HTTPTimeout=60.0, got %f", cfg.HTTPTimeout)
	}
}

func TestConfigFromEnvMissing(t *testing.T) {
	os.Unsetenv("LANSENGER_APP_ID")
	os.Unsetenv("LANSENGER_APP_SECRET")
	_, err := ConfigFromEnv()
	if err == nil {
		t.Error("expected error for missing env vars")
	}
}
