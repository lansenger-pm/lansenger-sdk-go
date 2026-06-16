package lansenger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type CredentialStore struct {
	path     string
	profile  string
	mu       sync.Mutex
	migrated bool
}

func NewCredentialStore(path string, profile string) (*CredentialStore, error) {
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home directory: %w", err)
		}
		path = filepath.Join(homeDir, ".lansenger", DefaultStateFile)
	}
	if profile == "" {
		profile = DefaultProfile
	}
	return &CredentialStore{
		path:    path,
		profile: profile,
	}, nil
}

type storeData struct {
	Profiles      map[string]profileData `json:"profiles"`
	ActiveProfile string                 `json:"active_profile"`
}

type userTokenEntry struct {
	UserToken             string `json:"user_token"`
	RefreshToken          string `json:"refresh_token"`
	UserTokenExpiresAt    int64  `json:"user_token_expiry"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expiry"`
}

type profileData struct {
	AppID                 string                    `json:"app_id"`
	AppSecret             string                    `json:"app_secret"`
	APIGatewayURL         string                    `json:"api_gateway_url"`
	PassportURL           string                    `json:"passport_url"`
	RedirectURI           string                    `json:"redirect_uri"`
	EncodingKey           string                    `json:"encoding_key"`
	CallbackToken         string                    `json:"callback_token"`
	AppToken              string                    `json:"app_token"`
	TokenExpiresAt        int64                     `json:"app_token_expiry"`
	UserToken             string                    `json:"user_token"`
	RefreshToken          string                    `json:"refresh_token"`
	UserTokenExpiresAt    int64                     `json:"user_token_expiry"`
	RefreshTokenExpiresAt int64                     `json:"refresh_token_expiry"`
	StaffID               string                    `json:"staff_id"`
	UserTokens            map[string]userTokenEntry `json:"user_tokens,omitempty"`
}

func (p *profileData) UnmarshalJSON(data []byte) error {
	type rawPD struct {
		AppID                 string                    `json:"app_id"`
		AppSecret             string                    `json:"app_secret"`
		APIGatewayURL         string                    `json:"api_gateway_url"`
		PassportURL           string                    `json:"passport_url"`
		RedirectURI           string                    `json:"redirect_uri"`
		EncodingKey           string                    `json:"encoding_key"`
		CallbackToken         string                    `json:"callback_token"`
		AppToken              string                    `json:"app_token"`
		TokenExpiresAt        int64                     `json:"app_token_expiry"`
		ATokenExpiresAtCompat *int64                    `json:"token_expires_at"`
		UserToken             string                    `json:"user_token"`
		RefreshToken          string                    `json:"refresh_token"`
		StaffID               string                    `json:"staff_id"`
		UserTokens            map[string]userTokenEntry `json:"user_tokens"`
		UserTokenExpiry       *int64                    `json:"user_token_expiry"`
		UserTokenExpiresAt    *int64                    `json:"user_token_expires_at"`
		RTokenExpiry          *int64                    `json:"refresh_token_expiry"`
		RTokenExpiresAt       *int64                    `json:"refresh_token_expires_at"`
	}
	var raw rawPD
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.AppID = raw.AppID
	p.AppSecret = raw.AppSecret
	p.APIGatewayURL = raw.APIGatewayURL
	p.PassportURL = raw.PassportURL
	p.RedirectURI = raw.RedirectURI
	p.EncodingKey = raw.EncodingKey
	p.CallbackToken = raw.CallbackToken
	p.AppToken = raw.AppToken
	p.TokenExpiresAt = raw.TokenExpiresAt
	if raw.ATokenExpiresAtCompat != nil && p.TokenExpiresAt == 0 {
		p.TokenExpiresAt = *raw.ATokenExpiresAtCompat
	}
	p.UserToken = raw.UserToken
	p.RefreshToken = raw.RefreshToken
	p.StaffID = raw.StaffID
	if raw.UserTokens != nil {
		p.UserTokens = raw.UserTokens
	}

	if raw.UserTokenExpiry != nil {
		p.UserTokenExpiresAt = *raw.UserTokenExpiry
	} else if raw.UserTokenExpiresAt != nil {
		p.UserTokenExpiresAt = *raw.UserTokenExpiresAt
	}
	if raw.RTokenExpiry != nil {
		p.RefreshTokenExpiresAt = *raw.RTokenExpiry
	} else if raw.RTokenExpiresAt != nil {
		p.RefreshTokenExpiresAt = *raw.RTokenExpiresAt
	}
	return nil
}

func (cs *CredentialStore) loadUnlocked() (*storeData, error) {
	data, err := os.ReadFile(cs.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &storeData{Profiles: map[string]profileData{}, ActiveProfile: DefaultProfile}, nil
		}
		return nil, fmt.Errorf("reading credential store: %w", err)
	}

	var sd storeData
	if err := json.Unmarshal(data, &sd); err != nil {
		return nil, fmt.Errorf("decoding credential store: %w", err)
	}

	if sd.Profiles == nil {
		sd.Profiles = map[string]profileData{}
	}
	if sd.ActiveProfile == "" {
		sd.ActiveProfile = DefaultProfile
	}

	return &sd, nil
}

func (cs *CredentialStore) saveUnlocked(sd *storeData) error {
	dir := filepath.Dir(cs.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating credential directory: %w", err)
	}

	data, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding credential store: %w", err)
	}

	if err := os.WriteFile(cs.path, data, 0600); err != nil {
		return fmt.Errorf("writing credential store: %w", err)
	}

	return nil
}

func (cs *CredentialStore) ensureMigrated() {
	if cs.migrated {
		return
	}
	cs.migrated = true

	data, err := os.ReadFile(cs.path)
	if err != nil {
		return
	}

	// Phase 1: migrate flat top-level format → profiles["default"]
	var sd storeData
	if json.Unmarshal(data, &sd) == nil && sd.Profiles != nil {
		// Phase 2: migrate flat userToken fields → user_tokens[staff_id] per profile
		for name, profile := range sd.Profiles {
			if profile.UserToken != "" && profile.StaffID != "" {
				if profile.UserTokens == nil {
					profile.UserTokens = map[string]userTokenEntry{}
				}
				// Always merge flat into nested — old SDK may rewrite flat after migration
				profile.UserTokens[profile.StaffID] = userTokenEntry{
					UserToken:             profile.UserToken,
					RefreshToken:          profile.RefreshToken,
					UserTokenExpiresAt:    profile.UserTokenExpiresAt,
					RefreshTokenExpiresAt: profile.RefreshTokenExpiresAt,
				}
				profile.UserToken = ""
				profile.RefreshToken = ""
				profile.UserTokenExpiresAt = 0
				profile.RefreshTokenExpiresAt = 0
				profile.StaffID = ""
				sd.Profiles[name] = profile
			}
		}
		cs.saveUnlocked(&sd)
		return
	}

	var flat map[string]interface{}
	if json.Unmarshal(data, &flat) != nil {
		return
	}
	if _, hasProfiles := flat["profiles"]; hasProfiles {
		return
	}

	profile := profileData{}
	if v, ok := flat["app_id"].(string); ok {
		profile.AppID = v
	}
	if v, ok := flat["app_secret"].(string); ok {
		profile.AppSecret = v
	}
	if v, ok := flat["api_gateway_url"].(string); ok {
		profile.APIGatewayURL = v
	}
	if v, ok := flat["passport_url"].(string); ok {
		profile.PassportURL = v
	}
	if v, ok := flat["redirect_uri"].(string); ok {
		profile.RedirectURI = v
	}
	if v, ok := flat["encoding_key"].(string); ok {
		profile.EncodingKey = v
	}
	if v, ok := flat["callback_token"].(string); ok {
		profile.CallbackToken = v
	}
	if v, ok := flat["app_token"].(string); ok {
		profile.AppToken = v
	}
	if v, ok := flat["app_token_expiry"].(float64); ok {
		profile.TokenExpiresAt = int64(v)
	} else if v, ok := flat["token_expires_at"].(float64); ok {
		profile.TokenExpiresAt = int64(v)
	}
	if v, ok := flat["user_token"].(string); ok {
		profile.UserToken = v
	}
	if v, ok := flat["refresh_token"].(string); ok {
		profile.RefreshToken = v
	}
	if v, ok := flat["user_token_expiry"].(float64); ok {
		profile.UserTokenExpiresAt = int64(v)
	} else if v, ok := flat["user_token_expires_at"].(float64); ok {
		profile.UserTokenExpiresAt = int64(v)
	}
	if v, ok := flat["refresh_token_expiry"].(float64); ok {
		profile.RefreshTokenExpiresAt = int64(v)
	} else if v, ok := flat["refresh_token_expires_at"].(float64); ok {
		profile.RefreshTokenExpiresAt = int64(v)
	}
	// Migrate flat userToken to nested if present
	if profile.UserToken != "" && profile.StaffID != "" {
		if profile.UserTokens == nil {
			profile.UserTokens = map[string]userTokenEntry{}
		}
		profile.UserTokens[profile.StaffID] = userTokenEntry{
			UserToken:             profile.UserToken,
			RefreshToken:          profile.RefreshToken,
			UserTokenExpiresAt:    profile.UserTokenExpiresAt,
			RefreshTokenExpiresAt: profile.RefreshTokenExpiresAt,
		}
		profile.UserToken = ""
		profile.RefreshToken = ""
		profile.UserTokenExpiresAt = 0
		profile.RefreshTokenExpiresAt = 0
	}

	newSD := &storeData{
		Profiles:      map[string]profileData{DefaultProfile: profile},
		ActiveProfile: DefaultProfile,
	}
	cs.saveUnlocked(newSD)
}

func (cs *CredentialStore) LoadCredentials() (map[string]string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return nil, err
	}

	profile, ok := sd.Profiles[cs.profile]
	if !ok {
		return map[string]string{}, nil
	}

	return map[string]string{
		"app_id":          profile.AppID,
		"app_secret":      profile.AppSecret,
		"api_gateway_url": profile.APIGatewayURL,
		"passport_url":    profile.PassportURL,
		"redirect_uri":    profile.RedirectURI,
		"encoding_key":    profile.EncodingKey,
		"callback_token":  profile.CallbackToken,
	}, nil
}

func (cs *CredentialStore) SaveCredentials(appID, appSecret, apiGatewayURL, passportURL, redirectURI string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		sd = &storeData{Profiles: map[string]profileData{}, ActiveProfile: DefaultProfile}
	}

	profile := sd.Profiles[cs.profile]
	profile.AppID = appID
	profile.AppSecret = appSecret
	profile.APIGatewayURL = apiGatewayURL
	profile.PassportURL = passportURL
	profile.RedirectURI = redirectURI
	sd.Profiles[cs.profile] = profile
	sd.ActiveProfile = cs.profile

	return cs.saveUnlocked(sd)
}

func (cs *CredentialStore) SaveCallbackConfig(encodingKey, callbackToken string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		sd = &storeData{Profiles: map[string]profileData{}, ActiveProfile: DefaultProfile}
	}

	profile := sd.Profiles[cs.profile]
	profile.EncodingKey = encodingKey
	profile.CallbackToken = callbackToken
	sd.Profiles[cs.profile] = profile
	sd.ActiveProfile = cs.profile

	return cs.saveUnlocked(sd)
}

func (cs *CredentialStore) LoadAppToken() (string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return "", err
	}

	profile, ok := sd.Profiles[cs.profile]
	if !ok {
		return "", nil
	}

	if profile.AppToken == "" {
		return "", nil
	}

	if time.Now().Unix() >= profile.TokenExpiresAt {
		return "", nil
	}

	return profile.AppToken, nil
}

func (cs *CredentialStore) SaveAppToken(token string, expiresIn int) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		sd = &storeData{Profiles: map[string]profileData{}, ActiveProfile: DefaultProfile}
	}

	profile := sd.Profiles[cs.profile]
	profile.AppToken = token
	margin := TokenRefreshMargin
	if expiresIn < margin*2 {
		margin = expiresIn / 2
	}
	profile.TokenExpiresAt = time.Now().Add(time.Duration(expiresIn-margin) * time.Second).Unix()
	sd.Profiles[cs.profile] = profile

	return cs.saveUnlocked(sd)
}

func (cs *CredentialStore) LoadUserToken(staffID string) (map[string]string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return nil, err
	}

	profile, ok := sd.Profiles[cs.profile]
	if !ok {
		return map[string]string{}, nil
	}

	if staffID != "" {
		if entry, ok := profile.UserTokens[staffID]; ok {
			return map[string]string{
				"user_token":           entry.UserToken,
				"refresh_token":        entry.RefreshToken,
				"user_token_expiry":    strconv.FormatInt(entry.UserTokenExpiresAt, 10),
				"refresh_token_expiry": strconv.FormatInt(entry.RefreshTokenExpiresAt, 10),
				"staff_id":             staffID,
			}, nil
		}
	}

	// Fallback: legacy flat fields
	flat := map[string]string{
		"user_token":           profile.UserToken,
		"refresh_token":        profile.RefreshToken,
		"user_token_expiry":    strconv.FormatInt(profile.UserTokenExpiresAt, 10),
		"refresh_token_expiry": strconv.FormatInt(profile.RefreshTokenExpiresAt, 10),
		"staff_id":             profile.StaffID,
	}
	if profile.UserToken != "" && profile.StaffID != "" {
		return flat, nil
	}

	// Post-migration: try first entry from user_tokens
	for sid, entry := range profile.UserTokens {
		return map[string]string{
			"user_token":           entry.UserToken,
			"refresh_token":        entry.RefreshToken,
			"user_token_expiry":    strconv.FormatInt(entry.UserTokenExpiresAt, 10),
			"refresh_token_expiry": strconv.FormatInt(entry.RefreshTokenExpiresAt, 10),
			"staff_id":             sid,
		}, nil
	}

	return flat, nil
}

func (cs *CredentialStore) SaveUserToken(userToken, refreshToken string, expiresIn int, refreshExpiresIn int, staffID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		sd = &storeData{Profiles: map[string]profileData{}, ActiveProfile: DefaultProfile}
	}

	profile := sd.Profiles[cs.profile]

	if staffID == "" {
		// Legacy flat path — no staff_id to key on
		profile.UserToken = userToken
		profile.RefreshToken = refreshToken
		if expiresIn > 0 {
			margin := UserTokenRefreshMargin
			if expiresIn < margin*2 {
				margin = expiresIn / 2
			}
			profile.UserTokenExpiresAt = time.Now().Add(time.Duration(expiresIn-margin) * time.Second).Unix()
		}
		if refreshExpiresIn > 0 {
			profile.RefreshTokenExpiresAt = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second).Unix()
		}
		profile.StaffID = ""
	} else {
		if profile.UserTokens == nil {
			profile.UserTokens = map[string]userTokenEntry{}
		}
		entry := profile.UserTokens[staffID]
		entry.UserToken = userToken
		entry.RefreshToken = refreshToken
		if expiresIn > 0 {
			margin := UserTokenRefreshMargin
			if expiresIn < margin*2 {
				margin = expiresIn / 2
			}
			entry.UserTokenExpiresAt = time.Now().Add(time.Duration(expiresIn-margin) * time.Second).Unix()
		}
		if refreshExpiresIn > 0 {
			entry.RefreshTokenExpiresAt = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second).Unix()
		}
		profile.UserTokens[staffID] = entry

		// Clean up legacy flat fields after first nested save
		profile.UserToken = ""
		profile.RefreshToken = ""
		profile.UserTokenExpiresAt = 0
		profile.RefreshTokenExpiresAt = 0
		profile.StaffID = ""
	}
	sd.Profiles[cs.profile] = profile

	return cs.saveUnlocked(sd)
}

func (cs *CredentialStore) HasCredentials() bool {
	creds, err := cs.LoadCredentials()
	if err != nil {
		return false
	}
	return creds["app_id"] != "" && creds["app_secret"] != ""
}

// HasFullConfig returns true if app_id, app_secret, and api_gateway_url
// are all non-empty for the current profile.
func (cs *CredentialStore) HasFullConfig() bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return false
	}

	profile, ok := sd.Profiles[cs.profile]
	if !ok {
		return false
	}

	return profile.AppID != "" && profile.AppSecret != "" && profile.APIGatewayURL != ""
}

func (cs *CredentialStore) ListProfiles() ([]string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return nil, err
	}

	profiles := make([]string, 0, len(sd.Profiles))
	for k := range sd.Profiles {
		profiles = append(profiles, k)
	}
	return profiles, nil
}

func (cs *CredentialStore) ListUserTokens() ([]string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return nil, err
	}

	profile, ok := sd.Profiles[cs.profile]
	if !ok || profile.UserTokens == nil {
		return []string{}, nil
	}

	users := make([]string, 0, len(profile.UserTokens))
	for k := range profile.UserTokens {
		users = append(users, k)
	}
	return users, nil
}

func (cs *CredentialStore) ClearProfile() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return err
	}

	delete(sd.Profiles, cs.profile)
	return cs.saveUnlocked(sd)
}

// DeleteProfileByName deletes a profile by name. If the deleted profile
// is the active profile, it falls back to "default".
// Returns true if the profile was found and deleted, false if it did not exist.
func (cs *CredentialStore) DeleteProfileByName(name string) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return false
	}

	if _, ok := sd.Profiles[name]; !ok {
		return false
	}
	delete(sd.Profiles, name)
	if sd.ActiveProfile == name {
		sd.ActiveProfile = DefaultProfile
	}
	if err := cs.saveUnlocked(sd); err != nil {
		return false
	}
	return true
}

func (cs *CredentialStore) Clear() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	return os.Remove(cs.path)
}

func (cs *CredentialStore) LoadState() (*storeData, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	return cs.loadUnlocked()
}

func (cs *CredentialStore) GetActiveProfile() string {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return DefaultProfile
	}
	return sd.ActiveProfile
}

func (cs *CredentialStore) SetActiveProfile(profile string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.ensureMigrated()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return err
	}
	sd.ActiveProfile = profile
	return cs.saveUnlocked(sd)
}

func (cs *CredentialStore) Path() string {
	return cs.path
}
