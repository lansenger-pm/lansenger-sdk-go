package lansenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LansengerClient struct {
	config       *Config
	httpClient   *http.Client
	tokenMgr     *TokenManager
	userTokenMgr *UserTokenManager
}

func NewClient(appID, appSecret string) *LansengerClient {
	cfg := NewConfig(appID, appSecret)
	return NewClientWithConfig(cfg)
}

func NewClientWithConfig(cfg *Config) *LansengerClient {
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.HTTPTimeout * float64(time.Second)),
	}
	return &LansengerClient{
		config:     cfg,
		httpClient: httpClient,
		tokenMgr:   NewTokenManager(cfg, httpClient),
	}
}

func NewClientFromEnv() (*LansengerClient, error) {
	cfg, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewClientWithConfig(cfg), nil
}

func NewClientFromStore(store *CredentialStore) (*LansengerClient, error) {
	creds, err := store.LoadCredentials()
	if err != nil {
		return nil, fmt.Errorf("loading credentials from store: %w", err)
	}
	if creds["app_id"] == "" || creds["app_secret"] == "" {
		return nil, NewConfigError("no credentials found in store for profile '" + store.profile + "'")
	}
	cfg := NewConfig(creds["app_id"], creds["app_secret"])
	if creds["api_gateway_url"] != "" {
		cfg.APIGatewayURL = creds["api_gateway_url"]
	}
	if creds["passport_url"] != "" {
		cfg.PassportURL = creds["passport_url"]
	}
	if creds["encoding_key"] != "" {
		cfg.EncodingKey = creds["encoding_key"]
	}
	if creds["callback_token"] != "" {
		cfg.CallbackToken = creds["callback_token"]
	}
	return NewClientWithConfig(cfg), nil
}

func (c *LansengerClient) GetToken(ctx context.Context) (string, error) {
	return c.tokenMgr.GetToken(ctx)
}

func (c *LansengerClient) InvalidateToken() {
	c.tokenMgr.Invalidate()
}

func (c *LansengerClient) HealthCheck(ctx context.Context) bool {
	token, err := c.GetToken(ctx)
	if err != nil {
		return false
	}
	return token != ""
}

func (c *LansengerClient) doGet(ctx context.Context, url string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GET request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewNetworkError("GET request failed: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decoding response JSON: %w", err)
	}

	errCode, _ := result["errCode"].(float64)
	if errCode != 0 {
		errMsg, _ := result["errMsg"].(string)
		return nil, NewAPIError(errMsg, int(errCode))
	}

	return result, nil
}

func (c *LansengerClient) doGetRaw(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GET request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewNetworkError("GET request failed: " + err.Error())
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *LansengerClient) doPost(ctx context.Context, url string, body interface{}) (map[string]interface{}, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewNetworkError("POST request failed: " + err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decoding response JSON: %w", err)
	}

	errCode, _ := result["errCode"].(float64)
	if errCode != 0 {
		errMsg, _ := result["errMsg"].(string)
		return nil, NewAPIError(errMsg, int(errCode))
	}

	return result, nil
}

func (c *LansengerClient) doPostMultipart(ctx context.Context, url string, filePath string, mediaType int) (map[string]interface{}, error) {
	return uploadMediaInternal(ctx, c.httpClient, url, filePath)
}

func extractData(result map[string]interface{}) map[string]interface{} {
	if data, ok := result["data"].(map[string]interface{}); ok {
		return data
	}
	return nil
}

func extractDataArray(result map[string]interface{}) []interface{} {
	if data, ok := result["data"].([]interface{}); ok {
		return data
	}
	return nil
}

func strFromMap(m map[string]interface{}, key string) string {
	v, _ := m[key].(string)
	return v
}

func intFromMap(m map[string]interface{}, key string) int {
	v, _ := m[key].(float64)
	return int(v)
}

func floatFromMap(m map[string]interface{}, key string) float64 {
	v, _ := m[key].(float64)
	return v
}

func boolFromMap(m map[string]interface{}, key string) bool {
	v, _ := m[key].(bool)
	return v
}

func mapFromMap(m map[string]interface{}, key string) map[string]interface{} {
	v, _ := m[key].(map[string]interface{})
	return v
}

func arrayFromMap(m map[string]interface{}, key string) []interface{} {
	v, _ := m[key].([]interface{})
	return v
}

func stringArrayFromMap(m map[string]interface{}, key string) []string {
	arr := arrayFromMap(m, key)
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
