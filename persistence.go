package lansenger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CredentialStore struct {
	path    string
	profile string
	mu      sync.RWMutex
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
	AppID              string `json:"app_id,omitempty"`
	AppSecret          string `json:"app_secret,omitempty"`
	APIGatewayURL      string `json:"api_gateway_url,omitempty"`
	PassportURL        string `json:"passport_url,omitempty"`
	AppToken           string `json:"app_token,omitempty"`
	TokenExpiresAt     int64  `json:"token_expires_at,omitempty"`
	UserToken          string `json:"user_token,omitempty"`
	RefreshToken       string `json:"refresh_token,omitempty"`
	UserTokenExpiresAt int64  `json:"user_token_expires_at,omitempty"`
}

func (cs *CredentialStore) load() (*storeData, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

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

func (cs *CredentialStore) save(sd *storeData) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

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
	sd, err := cs.load()
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
	}, nil
}

func (cs *CredentialStore) SaveCredentials(appID, appSecret, apiGatewayURL, passportURL string) error {
	sd, err := cs.load()
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

	return cs.save(sd)
}

func (cs *CredentialStore) LoadAppToken() (string, error) {
	sd, err := cs.load()
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
	sd, err := cs.load()
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

	return cs.save(sd)
}

func (cs *CredentialStore) LoadUserToken() (map[string]string, error) {
	sd, err := cs.load()
	if err != nil {
		return nil, err
	}

	profile, ok := sd.Profiles[cs.profile]
	if !ok {
		return map[string]string{}, nil
	}

	return map[string]string{
		"user_token":    profile.UserToken,
		"refresh_token": profile.RefreshToken,
	}, nil
}

func (cs *CredentialStore) SaveUserToken(userToken, refreshToken string, expiresIn int) error {
	sd, err := cs.load()
	if err != nil {
		sd = &storeData{Profiles: map[string]profileData{}, ActiveProfile: DefaultProfile}
	}

	profile := sd.Profiles[cs.profile]
	profile.UserToken = userToken
	profile.RefreshToken = refreshToken
	if expiresIn > 0 {
		profile.UserTokenExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()
	}
	sd.Profiles[cs.profile] = profile

	return cs.save(sd)
}

func (cs *CredentialStore) HasCredentials() bool {
	creds, err := cs.LoadCredentials()
	if err != nil {
		return false
	}
	return creds["app_id"] != "" && creds["app_secret"] != ""
}

func (cs *CredentialStore) ListProfiles() ([]string, error) {
	sd, err := cs.load()
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
	sd, err := cs.load()
	if err != nil {
		return err
	}

	delete(sd.Profiles, cs.profile)
	return cs.save(sd)
}

func (cs *CredentialStore) Clear() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	return os.Remove(cs.path)
}
