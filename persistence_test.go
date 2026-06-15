package lansenger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCredentialStoreInit(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.profile != "default" {
		t.Errorf("expected profile=default, got %s", store.profile)
	}
}

func TestCredentialStoreSaveAndLoadCredentials(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.SaveCredentials("app1", "secret1", "https://gateway.example.com", "https://passport.example.com", "")
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
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.HasCredentials() {
		t.Error("expected HasCredentials=false before saving")
	}

	store.SaveCredentials("app1", "secret1", "", "", "")
	if !store.HasCredentials() {
		t.Error("expected HasCredentials=true after saving")
	}
}

func TestCredentialStoreAppToken(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.SaveCredentials("app1", "secret1", "", "", "")

	err = store.SaveAppToken("test_token", 7200)
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
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.SaveCredentials("app1", "secret1", "", "", "")
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
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.SaveCredentials("app1", "secret1", "", "", "")
	store.SaveUserToken("utok1", "rtok1", 7200, 2592000, "staff1")

	tokens, err := store.LoadUserToken("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens["user_token"] != "utok1" {
		t.Errorf("expected user_token=utok1, got %s", tokens["user_token"])
	}
	if tokens["refresh_token"] != "rtok1" {
		t.Errorf("expected refresh_token=rtok1, got %s", tokens["refresh_token"])
	}
	if tokens["staff_id"] != "staff1" {
		t.Errorf("expected staff_id=staff1, got %s", tokens["staff_id"])
	}
}

func TestCredentialStoreClear(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.SaveCredentials("app1", "secret1", "", "", "")
	store.Clear()

	if store.HasCredentials() {
		t.Error("expected HasCredentials=false after clear")
	}
}

func TestCredentialStoreClearProfile(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.SaveCredentials("app1", "secret1", "", "", "")
	store.ClearProfile()

	if store.HasCredentials() {
		t.Error("expected HasCredentials=false after clearing profile")
	}
}

func TestCredentialStoreListProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.SaveCredentials("app1", "secret1", "", "", "")

	profiles, err := store.ListProfiles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Errorf("expected 1 profile, got %d", len(profiles))
	}
}

func TestCredentialStoreDefaultProfile(t *testing.T) {
	store, err := NewCredentialStore("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.profile != DefaultProfile {
		t.Errorf("expected profile=%s, got %s", DefaultProfile, store.profile)
	}
}

func TestCredentialStorePreservesState(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.SaveCredentials("app1", "secret1", "https://gateway.com", "", "")
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

func TestCredentialStoreDeleteProfileByName(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.SaveCredentials("app1", "secret1", "", "", "")
	err = store.DeleteProfileByName("default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	profiles, err := store.ListProfiles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected 0 profiles after delete, got %d", len(profiles))
	}
}

func TestCredentialStoreDeleteProfileByNameNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.SaveCredentials("app1", "secret1", "", "", "")
	err = store.DeleteProfileByName("ghost")
	if err == nil {
		t.Error("expected error for nonexistent profile, got nil")
	}
}

func TestCredentialStoreDeleteProfileByNamePreservesOthers(t *testing.T) {
	tmpDir := t.TempDir()
	storeA, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "alpha")
	storeB, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "beta")
	storeA.SaveCredentials("appA", "secA", "", "", "")
	storeB.SaveCredentials("appB", "secB", "", "", "")

	profiles, _ := storeA.ListProfiles()
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles before delete, got %d", len(profiles))
	}

	err := storeA.DeleteProfileByName("alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	profiles, _ = storeB.ListProfiles()
	if len(profiles) != 1 {
		t.Errorf("expected 1 profile after delete, got %d", len(profiles))
	}
	if profiles[0] != "beta" {
		t.Errorf("expected beta remaining, got %s", profiles[0])
	}
}

func TestCredentialStoreDeleteProfileByNameActiveFallback(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	store.SetActiveProfile("staging")

	stagingStore, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "staging")
	stagingStore.SaveCredentials("appX", "secX", "", "", "")

	if store.GetActiveProfile() != "staging" {
		t.Fatalf("expected active=staging before delete")
	}

	err := store.DeleteProfileByName("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.GetActiveProfile() != DefaultProfile {
		t.Errorf("expected active to fall back to %s, got %s", DefaultProfile, store.GetActiveProfile())
	}
}

// ── Multi-user userToken isolation ──────────────────────────────────────

func TestUserTokenMultiUserIsolation(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	store.SaveCredentials("app1", "secret1", "", "", "")

	store.SaveUserToken("token-a", "rt-a", 7200, 2592000, "staff-a")
	store.SaveUserToken("token-b", "rt-b", 7200, 2592000, "staff-b")

	a, _ := store.LoadUserToken("staff-a")
	b, _ := store.LoadUserToken("staff-b")

	if a["user_token"] != "token-a" {
		t.Errorf("staff-a user_token: expected token-a, got %s", a["user_token"])
	}
	if a["staff_id"] != "staff-a" {
		t.Errorf("staff-a staff_id: expected staff-a, got %s", a["staff_id"])
	}
	if b["user_token"] != "token-b" {
		t.Errorf("staff-b user_token: expected token-b, got %s", b["user_token"])
	}
}

func TestUserTokenIsolationPreventsOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	store.SaveCredentials("app1", "secret1", "", "", "")

	store.SaveUserToken("token-a", "rt-a", 7200, 0, "staff-a")
	store.SaveUserToken("token-b", "rt-b", 7200, 0, "staff-b")

	a, _ := store.LoadUserToken("staff-a")
	if a["user_token"] != "token-a" {
		t.Error("staff-a should still have its own token after staff-b save")
	}
}

func TestUserTokenCrossStaffIndependence(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	store.SaveCredentials("app1", "secret1", "", "", "")

	store.SaveUserToken("token-a-v1", "rt-a", 7200, 0, "staff-a")
	store.SaveUserToken("token-b", "rt-b", 7200, 0, "staff-b")

	// Update staff-a
	store.SaveUserToken("token-a-v2", "rt-a-v2", 7200, 0, "staff-a")

	a, _ := store.LoadUserToken("staff-a")
	b, _ := store.LoadUserToken("staff-b")
	if a["user_token"] != "token-a-v2" {
		t.Errorf("expected token-a-v2, got %s", a["user_token"])
	}
	if b["user_token"] != "token-b" {
		t.Error("staff-b must be untouched")
	}
}

func TestUserTokenBackwardCompatLegacyFlat(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test_state.json")

	// Write legacy flat format BEFORE creating the store
	// (store.migrated must be false for migration to run)
	now := time.Now().Unix()
	raw := map[string]interface{}{
		"profiles": map[string]interface{}{
			"default": map[string]interface{}{
				"app_id":               "app1",
				"app_secret":           "secret1",
				"user_token":           "legacy-ut",
				"refresh_token":        "legacy-rt",
				"staff_id":             "legacy-staff",
				"user_token_expiry":    int(now + 7200),
				"refresh_token_expiry": int(now + 2592000),
			},
		},
		"active_profile": "default",
	}
	data, _ := json.MarshalIndent(raw, "", "  ")
	os.WriteFile(path, data, 0600)

	store, _ := NewCredentialStore(path, "default")

	// Load — migration should run and return the legacy token via fallback
	got, _ := store.LoadUserToken("")
	if got["user_token"] != "legacy-ut" {
		t.Errorf("expected legacy-ut, got %s", got["user_token"])
	}

	// After migration, load by exact staff_id should work from nested store
	gotNested, _ := store.LoadUserToken("legacy-staff")
	if gotNested["user_token"] != "legacy-ut" {
		t.Errorf("nested: expected legacy-ut, got %s", gotNested["user_token"])
	}
}

func TestUserTokenRawStateStructure(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	store.SaveCredentials("app1", "secret1", "", "", "")

	store.SaveUserToken("t-a", "r-a", 7200, 0, "staff-a")
	store.SaveUserToken("t-b", "r-b", 7200, 0, "staff-b")

	sd, _ := store.LoadState()
	profile := sd.Profiles["default"]
	if profile.UserTokens == nil {
		t.Fatal("UserTokens should not be nil")
	}

	entryA, okA := profile.UserTokens["staff-a"]
	entryB, okB := profile.UserTokens["staff-b"]
	if !okA || entryA.UserToken != "t-a" {
		t.Errorf("staff-a: expected t-a, got %s", entryA.UserToken)
	}
	if !okB || entryB.UserToken != "t-b" {
		t.Errorf("staff-b: expected t-b, got %s", entryB.UserToken)
	}
}

func TestUserTokenNoStaffIDStillWritesFlat(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	store.SaveUserToken("flat-ut", "flat-rt", 7200, 0, "")
	got, _ := store.LoadUserToken("")
	if got["user_token"] != "flat-ut" {
		t.Errorf("expected flat-ut, got %s", got["user_token"])
	}
}

func TestUserTokenMigrationCleansStaleFlat(t *testing.T) {
	// Issue #2: flat fields written by old SDK after migration are cleaned.
	// Simulation: a file that has BOTH nested user_tokens AND flat user_token/staff_id.
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test_state.json")
	now := time.Now().Unix()

	// Write a file where nested already exists AND flat has (newer) data
	raw := map[string]interface{}{
		"profiles": map[string]interface{}{
			"default": map[string]interface{}{
				"app_id":     "app1",
				"app_secret": "secret1",
				// Nested (old data from previous migration)
				"user_tokens": map[string]interface{}{
					"staff-1": map[string]interface{}{
						"user_token":             "nested-old",
						"refresh_token":          "nested-rt",
						"user_token_expiry":      int(now + 3600),
						"refresh_token_expiry":   int(now + 86400),
					},
				},
				// Flat (written by old SDK) — has NEWER token
				"user_token":           "flat-new",
				"refresh_token":        "flat-rt-new",
				"staff_id":             "staff-1",
				"user_token_expiry":    int(now + 7200),
				"refresh_token_expiry": int(now + 172800),
			},
		},
		"active_profile": "default",
	}
	data, _ := json.MarshalIndent(raw, "", "  ")
	os.WriteFile(path, data, 0600)

	// Create store — migration should merge flat into nested and clean flat
	store, _ := NewCredentialStore(path, "default")
	got, _ := store.LoadUserToken("staff-1")
	if got["user_token"] != "flat-new" {
		t.Errorf("expected flat-new (flat overrides nested), got %s", got["user_token"])
	}

	// Verify file has no flat token values
	rawAfter, _ := os.ReadFile(path)
	var parsed map[string]interface{}
	json.Unmarshal(rawAfter, &parsed)
	profiles := parsed["profiles"].(map[string]interface{})
	profile := profiles["default"].(map[string]interface{})

	ut, _ := profile["user_token"]
	if ut != nil && ut != "" {
		t.Error("flat user_token should be empty after migration")
	}
	sid, _ := profile["staff_id"]
	if sid != nil && sid != "" {
		t.Error("flat staff_id should be empty after migration")
	}
}
