package lansenger

import (
	"context"
	"testing"
)

func TestCreateBotCommands(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/bot/commands/create", 0, "ok", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	cmds := []map[string]interface{}{
		{"command": "add", "description": "add something"},
	}
	result, err := c.CreateBotCommands(context.Background(), 7, cmds, "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestCreateBotCommandsWithChat(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/bot/commands/create", 0, "ok", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	cmds := []map[string]interface{}{
		{"command": "add"},
	}
	result, err := c.CreateBotCommands(context.Background(), 1, cmds, "524288-xxx", "group", "524288-yyy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestCreateBotCommandsAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/bot/commands/create", 10000, "error", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	cmds := []map[string]interface{}{
		{"command": "test"},
	}
	result, err := c.CreateBotCommands(context.Background(), 7, cmds, "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}

func TestFetchBotCommands(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/bot/commands/fetch", 0, "ok", map[string]interface{}{
			"scopeType": float64(7),
			"chatId":    "c1",
			"commands": []interface{}{
				map[string]interface{}{"command": "add", "description": "desc"},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchBotCommands(context.Background(), 7, "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.ScopeType != 7 {
		t.Errorf("expected ScopeType=7, got %d", result.ScopeType)
	}
	if len(result.Commands) != 1 {
		t.Errorf("expected 1 command, got %d", len(result.Commands))
	}
}

func TestDeleteBotCommands(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/bot/commands/delete", 0, "ok", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DeleteBotCommands(context.Background(), 7, "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}
