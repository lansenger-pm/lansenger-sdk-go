package lansenger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type TokenManager struct {
	config     *Config
	httpClient *http.Client
	token      string
	expiresAt  time.Time
	mu         sync.RWMutex
}

func NewTokenManager(cfg *Config, httpClient *http.Client) *TokenManager {
	return &TokenManager{
		config:     cfg,
		httpClient: httpClient,
	}
}

func (tm *TokenManager) GetToken(ctx context.Context) (string, error) {
	tm.mu.RLock()
	if tm.token != "" && time.Now().Before(tm.expiresAt) {
		token := tm.token
		tm.mu.RUnlock()
		return token, nil
	}
	tm.mu.RUnlock()

	return tm.refreshToken(ctx)
}

func (tm *TokenManager) refreshToken(ctx context.Context) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.token != "" && time.Now().Before(tm.expiresAt) {
		return tm.token, nil
	}

	url := BuildAPIURL(tm.config, "auth", "app_token_create", "",
		WithGrantType("client_credential"),
		WithAppID(tm.config.AppID),
		WithSecret(tm.config.AppSecret),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating token request: %w", err)
	}

	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return "", NewNetworkError("network error fetching app token: " + err.Error())
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding token response: %w", err)
	}

	errCode, _ := result["errCode"].(float64)
	if errCode != 0 {
		errMsg, _ := result["errMsg"].(string)
		return "", NewAuthError(fmt.Sprintf("token request failed: %s (errCode: %d)", errMsg, int(errCode)))
	}

	data, _ := result["data"].(map[string]interface{})
	appToken, _ := data["appToken"].(string)
	expiresIn, _ := data["expiresIn"].(float64)
	if expiresIn == 0 {
		expiresIn = 7200
	}

	margin := float64(TokenRefreshMargin)
	if expiresIn < margin*2 {
		margin = expiresIn / 2
	}

	tm.token = appToken
	tm.expiresAt = time.Now().Add(time.Duration(expiresIn-margin) * time.Second)

	return tm.token, nil
}

func (tm *TokenManager) Invalidate() {
	tm.mu.Lock()
	tm.token = ""
	tm.expiresAt = time.Time{}
	tm.mu.Unlock()
}
