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
	path    string
	profile string
	mu      sync.Mutex
}

func NewCredentialStore(path string, profile string) *CredentialStore {
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "~"
		}
		path = filepath.Join(homeDir, ".lansenger", DefaultStateFile)
	}
	if profile == "" {
		profile = DefaultProfile
	}
	return &CredentialStore{
		path:    path,
		profile: profile,
	}
}

type storeData struct {
	Profiles      map[string]profileData `json:"profiles"`
	ActiveProfile string                 `json:"active_profile"`
}

type profileData struct {
	AppID                 string `json:"app_id"`
	AppSecret             string `json:"app_secret"`
	APIGatewayURL         string `json:"api_gateway_url"`
	PassportURL           string `json:"passport_url"`
	EncodingKey           string `json:"encoding_key"`
	CallbackToken         string `json:"callback_token"`
	AppToken              string `json:"app_token"`
	TokenExpiresAt        int64  `json:"token_expires_at"`
	UserToken             string `json:"user_token"`
	RefreshToken          string `json:"refresh_token"`
	UserTokenExpiresAt    int64  `json:"user_token_expires_at"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
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

func (cs *CredentialStore) LoadCredentials() (map[string]string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

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
		"encoding_key":    profile.EncodingKey,
		"callback_token":  profile.CallbackToken,
	}, nil
}

func (cs *CredentialStore) SaveCredentials(appID, appSecret, apiGatewayURL, passportURL string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	sd, err := cs.loadUnlocked()
	if err != nil {
		sd = &storeData{Profiles: map[string]profileData{}, ActiveProfile: DefaultProfile}
	}

	profile := sd.Profiles[cs.profile]
	profile.AppID = appID
	profile.AppSecret = appSecret
	profile.APIGatewayURL = apiGatewayURL
	profile.PassportURL = passportURL
	sd.Profiles[cs.profile] = profile
	sd.ActiveProfile = cs.profile

	return cs.saveUnlocked(sd)
}

func (cs *CredentialStore) SaveCallbackConfig(encodingKey, callbackToken string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

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

func (cs *CredentialStore) LoadUserToken() (map[string]string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return nil, err
	}

	profile, ok := sd.Profiles[cs.profile]
	if !ok {
		return map[string]string{}, nil
	}

	return map[string]string{
		"user_token":          profile.UserToken,
		"refresh_token":       profile.RefreshToken,
		"user_token_expires_at": strconv.FormatInt(profile.UserTokenExpiresAt, 10),
		"refresh_token_expires_at": strconv.FormatInt(profile.RefreshTokenExpiresAt, 10),
	}, nil
}

func (cs *CredentialStore) SaveUserToken(userToken, refreshToken string, expiresIn int, refreshExpiresIn int) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	sd, err := cs.loadUnlocked()
	if err != nil {
		sd = &storeData{Profiles: map[string]profileData{}, ActiveProfile: DefaultProfile}
	}

	profile := sd.Profiles[cs.profile]
	profile.UserToken = userToken
	profile.RefreshToken = refreshToken
	if expiresIn > 0 {
		profile.UserTokenExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()
	}
	if refreshExpiresIn > 0 {
		profile.RefreshTokenExpiresAt = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second).Unix()
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

func (cs *CredentialStore) ListProfiles() ([]string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

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

func (cs *CredentialStore) ClearProfile() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return err
	}

	delete(sd.Profiles, cs.profile)
	return cs.saveUnlocked(sd)
}

func (cs *CredentialStore) Clear() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	return os.Remove(cs.path)
}

func (cs *CredentialStore) LoadState() (*storeData, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	return cs.loadUnlocked()
}

func (cs *CredentialStore) GetActiveProfile() string {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return DefaultProfile
	}
	return sd.ActiveProfile
}

func (cs *CredentialStore) SetActiveProfile(profile string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

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

func (cs *CredentialStore) migrateLegacyFormat() {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	sd, err := cs.loadUnlocked()
	if err != nil {
		return
	}

	if len(sd.Profiles) == 1 && sd.ActiveProfile == DefaultProfile {
		p := sd.Profiles[DefaultProfile]
		if p.AppID != "" && p.APIGatewayURL == "" && p.AppToken != "" && p.AppSecret == "" {
			return
		}
	}

	data, err := os.ReadFile(cs.path)
	if err != nil {
		return
	}

	var flat map[string]interface{}
	if json.Unmarshal(data, &flat) == nil {
		if _, hasProfiles := flat["profiles"]; !hasProfiles {
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
			if v, ok := flat["encoding_key"].(string); ok {
				profile.EncodingKey = v
			}
			if v, ok := flat["callback_token"].(string); ok {
				profile.CallbackToken = v
			}
			if v, ok := flat["app_token"].(string); ok {
				profile.AppToken = v
			}
			if v, ok := flat["token_expires_at"].(float64); ok {
				profile.TokenExpiresAt = int64(v)
			}
			if v, ok := flat["user_token"].(string); ok {
				profile.UserToken = v
			}
			if v, ok := flat["refresh_token"].(string); ok {
				profile.RefreshToken = v
			}
			if v, ok := flat["user_token_expires_at"].(float64); ok {
				profile.UserTokenExpiresAt = int64(v)
			}
			if v, ok := flat["refresh_token_expires_at"].(float64); ok {
				profile.RefreshTokenExpiresAt = int64(v)
			}

			newSD := &storeData{
				Profiles:      map[string]profileData{DefaultProfile: profile},
				ActiveProfile: DefaultProfile,
			}
			cs.saveUnlocked(newSD)
		}
	}
}

func init() {
	store := NewCredentialStore("", DefaultProfile)
	store.migrateLegacyFormat()
}