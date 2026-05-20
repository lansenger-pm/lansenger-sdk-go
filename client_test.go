package lansenger

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClientFromEnv(t *testing.T) {
	os.Setenv("LANSENGER_APP_ID", "env_id")
	os.Setenv("LANSENGER_APP_SECRET", "env_secret")
	defer os.Unsetenv("LANSENGER_APP_ID")
	defer os.Unsetenv("LANSENGER_APP_SECRET")

	c, err := NewClientFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.config.AppID != "env_id" {
		t.Errorf("expected AppID=env_id, got %s", c.config.AppID)
	}
}

func TestClientFromEnvMissing(t *testing.T) {
	os.Unsetenv("LANSENGER_APP_ID")
	os.Unsetenv("LANSENGER_APP_SECRET")
	_, err := NewClientFromEnv()
	if err == nil {
		t.Error("expected error for missing env vars")
	}
}

func TestClientHealthCheck(t *testing.T) {
	server := httptest.NewServer(mockAppTokenHandler("test_token"))
	defer server.Close()

	c := newTestClient(server)
	if !c.HealthCheck(context.Background()) {
		t.Error("expected health check to pass")
	}
}

func TestClientInvalidateToken(t *testing.T) {
	server := httptest.NewServer(mockAppTokenHandler("test_token"))
	defer server.Close()

	c := newTestClient(server)
	c.InvalidateToken()

	token, err := c.GetToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "test_token" {
		t.Errorf("expected fresh token, got %s", token)
	}
}

func TestRevokeMessage(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/messages/revoke", 0, "ok", map[string]interface{}{}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.RevokeMessage(context.Background(), []string{"msg1", "msg2"}, "bot", "sender1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Operation != "revoke_message" {
		t.Errorf("expected Operation=revoke_message, got %s", result.Operation)
	}
}

func TestRevokeMessageAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/messages/revoke", 50080, "message error", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.RevokeMessage(context.Background(), []string{"msg1"}, "bot", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}

func TestQueryGroups(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/groups/fetch", 0, "ok", map[string]interface{}{
			"totalGroupIds": 3,
			"groupIds":      []string{"g1", "g2", "g3"},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.QueryGroups(context.Background(), 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.TotalGroupIDs != 3 {
		t.Errorf("expected TotalGroupIDs=3, got %d", result.TotalGroupIDs)
	}
	if result.Operation != "query_groups" {
		t.Errorf("expected Operation=query_groups, got %s", result.Operation)
	}
}
