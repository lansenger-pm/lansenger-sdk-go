package lansenger

import (
	"path/filepath"
	"testing"
	"time"
)

func TestCredentialStoreInit(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if store.profile != "default" {
		t.Errorf("expected profile=default, got %s", store.profile)
	}
}

func TestCredentialStoreSaveAndLoadCredentials(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	err := store.SaveCredentials("app1", "secret1", "https://gateway.example.com", "https://passport.example.com")
	if err != nil {
		t.Fatalf("unexpected error saving credentials: %v", err)
	}

	creds, err := store.LoadCredentials()
	if err != nil {
		t.Fatalf("unexpected error loading credentials: %v", err)
	}
	if creds["app_id"] != "app1" {
		t.Errorf("expected app_id=app1, got %s", creds["app_id"])
	}
	if creds["app_secret"] != "secret1" {
		t.Errorf("expected app_secret=secret1, got %s", creds["app_secret"])
	}
	if creds["api_gateway_url"] != "https://gateway.example.com" {
		t.Errorf("expected custom gateway URL, got %s", creds["api_gateway_url"])
	}
	if creds["passport_url"] != "https://passport.example.com" {
		t.Errorf("expected passport URL, got %s", creds["passport_url"])
	}
}

func TestCredentialStoreHasCredentials(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	if store.HasCredentials() {
		t.Error("expected HasCredentials=false before saving")
	}

	store.SaveCredentials("app1", "secret1", "", "")
	if !store.HasCredentials() {
		t.Error("expected HasCredentials=true after saving")
	}
}

func TestCredentialStoreAppToken(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	store.SaveCredentials("app1", "secret1", "", "")

	err := store.SaveAppToken("test_token", 7200)
	if err != nil {
		t.Fatalf("unexpected error saving app token: %v", err)
	}

	token, err := store.LoadAppToken()
	if err != nil {
		t.Fatalf("unexpected error loading app token: %v", err)
	}
	if token != "test_token" {
		t.Errorf("expected token=test_token, got %s", token)
	}
}

func TestCredentialStoreAppTokenExpired(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	store.SaveCredentials("app1", "secret1", "", "")
	store.SaveAppToken("expired_token", 1)

	time.Sleep(2 * time.Second)

	token, err := store.LoadAppToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "" {
		t.Errorf("expected empty token for expired, got %s", token)
	}
}

func TestCredentialStoreUserToken(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	store.SaveCredentials("app1", "secret1", "", "")
	store.SaveUserToken("utok1", "rtok1", 7200, 2592000)

	tokens, err := store.LoadUserToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens["user_token"] != "utok1" {
		t.Errorf("expected user_token=utok1, got %s", tokens["user_token"])
	}
	if tokens["refresh_token"] != "rtok1" {
		t.Errorf("expected refresh_token=rtok1, got %s", tokens["refresh_token"])
	}
}

func TestCredentialStoreClear(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	store.SaveCredentials("app1", "secret1", "", "")
	store.Clear()

	if store.HasCredentials() {
		t.Error("expected HasCredentials=false after clear")
	}
}

func TestCredentialStoreClearProfile(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	store.SaveCredentials("app1", "secret1", "", "")
	store.ClearProfile()

	if store.HasCredentials() {
		t.Error("expected HasCredentials=false after clearing profile")
	}
}

func TestCredentialStoreListProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	store.SaveCredentials("app1", "secret1", "", "")

	profiles, err := store.ListProfiles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Errorf("expected 1 profile, got %d", len(profiles))
	}
}

func TestCredentialStoreDefaultProfile(t *testing.T) {
	store := NewCredentialStore("", "")
	if store.profile != DefaultProfile {
		t.Errorf("expected profile=%s, got %s", DefaultProfile, store.profile)
	}
}

func TestCredentialStorePreservesState(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	store.SaveCredentials("app1", "secret1", "https://gateway.com", "")
	store.SaveAppToken("tok1", 7200)

	creds, _ := store.LoadCredentials()
	if creds["app_id"] != "app1" {
		t.Errorf("expected app_id preserved, got %s", creds["app_id"])
	}

	tok, _ := store.LoadAppToken()
	if tok != "tok1" {
		t.Errorf("expected token preserved, got %s", tok)
	}
}
