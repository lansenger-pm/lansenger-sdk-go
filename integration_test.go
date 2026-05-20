package lansenger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func loadCredentialsFromStore() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home dir: %w", err)
	}
	path := filepath.Join(homeDir, ".lansenger", "sdk_state.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading credential store: %w", err)
	}

	var store struct {
		Profiles map[string]struct {
			AppID         string `json:"app_id"`
			AppSecret     string `json:"app_secret"`
			APIGatewayURL string `json:"api_gateway_url"`
			PassportURL   string `json:"passport_url"`
		} `json:"profiles"`
		ActiveProfile string `json:"active_profile"`
	}

	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("decoding credential store: %w", err)
	}

	profile := store.Profiles[store.ActiveProfile]
	if profile.AppID == "" || profile.AppSecret == "" {
		return nil, fmt.Errorf("no credentials found in profile '%s'", store.ActiveProfile)
	}

	cfg := NewConfig(profile.AppID, profile.AppSecret)
	if profile.APIGatewayURL != "" {
		cfg.APIGatewayURL = profile.APIGatewayURL
	}
	if profile.PassportURL != "" {
		cfg.PassportURL = profile.PassportURL
	}

	return cfg, nil
}

func newIntegrationClient(t *testing.T) *LansengerClient {
	t.Helper()

	cfg, err := loadCredentialsFromStore()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}

	t.Logf("Using AppID=%s, Gateway=%s", cfg.AppID, cfg.APIGatewayURL)
	return NewClientWithConfig(cfg)
}

func TestIntegration_GetToken(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	token, err := c.GetToken(ctx)
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("GetToken returned empty token")
	}
	t.Logf("Got appToken: %s (length=%d)", token[:8]+"...", len(token))
}

func TestIntegration_HealthCheck(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if !c.HealthCheck(ctx) {
		t.Fatal("HealthCheck failed")
	}
	t.Log("HealthCheck passed")
}

func TestIntegration_FetchOrgInfo(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := c.FetchOrgInfo(ctx, "1", "")
	if err != nil {
		t.Fatalf("FetchOrgInfo error: %v", err)
	}
	if result.Success {
		t.Logf("OrgInfo: orgId=%s, orgName=%s", result.OrgID, result.OrgName)
	} else {
		t.Logf("FetchOrgInfo returned unsuccessful (may need userToken): error=%s", result.Error)
	}
}

func TestIntegration_FetchDepartmentChildren(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := c.FetchDepartmentChildren(ctx, "1", "")
	if err != nil {
		t.Fatalf("FetchDepartmentChildren error: %v", err)
	}
	t.Logf("DepartmentChildren: success=%v, departments count=%d", result.Success, len(result.Departments))
}

func TestIntegration_FetchDepartmentStaffs(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := c.FetchDepartmentStaffs(ctx, "1", "", 1, 5)
	if err != nil {
		t.Fatalf("FetchDepartmentStaffs error: %v", err)
	}
	t.Logf("DepartmentStaffs: success=%v, hasMore=%v, total=%d", result.Success, result.HasMore, result.Total)
}

func TestIntegration_FetchGroupList(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := c.FetchGroupList(ctx, "", 0, 5)
	if err != nil {
		t.Fatalf("FetchGroupList error: %v", err)
	}
	t.Logf("GroupList: success=%v, totalGroupIds=%d", result.Success, result.TotalGroupIDs)
}

func TestIntegration_QueryGroups(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := c.QueryGroups(ctx, 1, 5)
	if err != nil {
		t.Fatalf("QueryGroups error: %v", err)
	}
	t.Logf("QueryGroups: success=%v, totalGroupIds=%d, groupIds=%v", result.Success, result.TotalGroupIDs, result.GroupIDs)
}

func TestIntegration_BuildAuthorizeURL(t *testing.T) {
	c := newIntegrationClient(t)

	if c.config.PassportURL == "" {
		t.Skip("Skipping: no passport URL configured")
	}

	authURL := c.BuildAuthorizeURL("https://localhost:8080/callback", "", "test_state_123")
	t.Logf("Authorize URL: %s", authURL)
	if authURL == "" {
		t.Fatal("BuildAuthorizeURL returned empty URL")
	}
}

func TestIntegration_CallbackEventTypes(t *testing.T) {
	types := GetCallbackEventTypes()
	t.Logf("Callback event types count: %d", len(types))
	for eventType, category := range types {
		t.Logf("  %s -> %s", eventType, category)
	}
}

func TestIntegration_TokenInvalidateAndRefresh(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	token1, err := c.GetToken(ctx)
	if err != nil {
		t.Fatalf("first GetToken failed: %v", err)
	}

	c.InvalidateToken()

	token2, err := c.GetToken(ctx)
	if err != nil {
		t.Fatalf("second GetToken after invalidation failed: %v", err)
	}

	t.Logf("Token1=%s..., Token2=%s..., different=%v", token1[:8], token2[:8], token1 != token2)

	if token1 == token2 {
		t.Log("Warning: tokens are the same (server may have cached)")
	}
}