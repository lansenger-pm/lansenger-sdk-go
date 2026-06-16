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

func TestCredentialStoreHasFullConfig(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	if store.HasFullConfig() {
		t.Error("expected HasFullConfig=false before saving")
	}

	// Missing api_gateway_url
	store.SaveCredentials("app1", "secret1", "", "", "")
	if store.HasFullConfig() {
		t.Error("expected HasFullConfig=false without gateway URL")
	}

	// Full config
	store.SaveCredentials("app1", "secret1", "https://gateway.example.com", "", "")
	if !store.HasFullConfig() {
		t.Error("expected HasFullConfig=true after saving full config")
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
	if !store.DeleteProfileByName("default") {
		t.Fatal("expected DeleteProfileByName to return true")
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
	if store.DeleteProfileByName("ghost") {
		t.Error("expected false for nonexistent profile, got true")
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

	if !storeA.DeleteProfileByName("alpha") {
		t.Fatal("expected DeleteProfileByName to return true")
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

	if !store.DeleteProfileByName("staging") {
		t.Fatal("expected DeleteProfileByName to return true")
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

// ── ListUserTokens ──────────────────────────────────────────────

func TestListUserTokensEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")

	users, err := store.ListUserTokens()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected empty list, got %d", len(users))
	}
}

func TestListUserTokensSingleUser(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	store.SaveCredentials("app1", "secret1", "", "", "")
	store.SaveUserToken("token1", "rt1", 7200, 0, "staff1")

	users, err := store.ListUserTokens()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 || users[0] != "staff1" {
		t.Errorf("expected [staff1], got %v", users)
	}
}

func TestListUserTokensMultipleUsers(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	store.SaveCredentials("app1", "secret1", "", "", "")
	store.SaveUserToken("token1", "rt1", 7200, 0, "staff1")
	store.SaveUserToken("token2", "rt2", 7200, 0, "staff2")
	store.SaveUserToken("token3", "rt3", 7200, 0, "staff3")

	users, err := store.ListUserTokens()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}

	found := make(map[string]bool)
	for _, u := range users {
		found[u] = true
	}
	if !found["staff1"] || !found["staff2"] || !found["staff3"] {
		t.Errorf("expected staff1, staff2, staff3, got %v", users)
	}
}

func TestListUserTokensProfileIsolation(t *testing.T) {
	tmpDir := t.TempDir()
	storeAlpha, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "alpha")
	storeBeta, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "beta")

	storeAlpha.SaveCredentials("appA", "secA", "", "", "")
	storeBeta.SaveCredentials("appB", "secB", "", "", "")

	storeAlpha.SaveUserToken("t1", "rt1", 7200, 0, "staff-a")
	storeBeta.SaveUserToken("t2", "rt2", 7200, 0, "staff-b")

	alphaUsers, _ := storeAlpha.ListUserTokens()
	betaUsers, _ := storeBeta.ListUserTokens()

	if len(alphaUsers) != 1 || alphaUsers[0] != "staff-a" {
		t.Errorf("alpha: expected [staff-a], got %v", alphaUsers)
	}
	if len(betaUsers) != 1 || betaUsers[0] != "staff-b" {
		t.Errorf("beta: expected [staff-b], got %v", betaUsers)
	}
}

func TestUserTokenAutoMigrationOnSave(t *testing.T) {
	// Verify that saving with staff_id triggers auto-migration of flat fields.
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test_state.json")
	now := time.Now().Unix()

	// 1. Write legacy flat format
	raw := map[string]interface{}{
		"profiles": map[string]interface{}{
			"default": map[string]interface{}{
				"app_id":               "app1",
				"app_secret":           "secret1",
				"user_token":           "legacy-ut",
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

	// 2. Flat is readable before migration
	got, _ := store.LoadUserToken("")
	if got["user_token"] != "legacy-ut" {
		t.Errorf("expected legacy-ut, got %s", got["user_token"])
	}

	// 3. Save with staff_id for a *different* user — triggers migration
	store.SaveUserToken("nested-ut", "nested-rt", 7200, 2592000, "nested-staff")

	// 4. After migration, flat fields should be gone
	sd, _ := store.LoadState()
	profile := sd.Profiles["default"]
	if profile.UserToken != "" {
		t.Error("flat user_token should be cleaned after migration")
	}
	if profile.StaffID != "" {
		t.Error("flat staff_id should be cleaned after migration")
	}

	// 5. Legacy user accessible via nested
	legacy, _ := store.LoadUserToken("legacy-staff")
	if legacy["user_token"] != "legacy-ut" {
		t.Errorf("legacy: expected legacy-ut, got %s", legacy["user_token"])
	}

	// 6. New user accessible via nested
	nested, _ := store.LoadUserToken("nested-staff")
	if nested["user_token"] != "nested-ut" {
		t.Errorf("nested: expected nested-ut, got %s", nested["user_token"])
	}
}

func TestUserTokenNoStaffIDFallback(t *testing.T) {
	// loadUserToken("") returns first available user from nested store
	// when flat fields are empty (post-migration scenario).
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	store.SaveCredentials("app1", "secret1", "", "", "")

	store.SaveUserToken("t1", "r1", 7200, 0, "staff1")
	store.SaveUserToken("t2", "r2", 7200, 0, "staff2")

	// No staff_id → falls back to first entry from nested
	fallback, _ := store.LoadUserToken("")
	if fallback["user_token"] != "t1" && fallback["user_token"] != "t2" {
		t.Errorf("expected t1 or t2, got %s", fallback["user_token"])
	}

	// With exact staff_id, we get the specific one
	one, _ := store.LoadUserToken("staff1")
	two, _ := store.LoadUserToken("staff2")
	if one["user_token"] != "t1" {
		t.Errorf("staff1: expected t1, got %s", one["user_token"])
	}
	if two["user_token"] != "t2" {
		t.Errorf("staff2: expected t2, got %s", two["user_token"])
	}
}

func TestUserTokenNonexistentStaffID(t *testing.T) {
	// loadUserToken with a non-existent staff_id falls back to available tokens.
	tmpDir := t.TempDir()
	store, _ := NewCredentialStore(filepath.Join(tmpDir, "test_state.json"), "default")
	store.SaveCredentials("app1", "secret1", "", "", "")

	store.SaveUserToken("t1", "", 7200, 0, "staff1")
	got, _ := store.LoadUserToken("ghost-staff")
	// Fallback returns first available user (or empty)
	if got["user_token"] != "" && got["user_token"] != "t1" {
		t.Errorf("expected '' or 't1', got %s", got["user_token"])
	}
}
