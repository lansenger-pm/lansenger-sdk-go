package lansenger

import (
	"context"
	"testing"
)

func TestCreateStreamMessage(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/sse/msg/create", 0, "ok", map[string]interface{}{
			"msgId": "sse001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.CreateStreamMessage(context.Background(), "s001", "staff", "stream123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.MessageID != "sse001" {
		t.Errorf("expected MessageID=sse001, got %s", result.MessageID)
	}
}

func TestFetchStreamMessage(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/sse/msg/fetch", 0, "ok", map[string]interface{}{
			"msgId": "sse001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchStreamMessage(context.Background(), "sse001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.MessageID != "sse001" {
		t.Errorf("expected MessageID=sse001, got %s", result.MessageID)
	}
}

func TestCreateStreamMessageAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/sse/msg/create", 59000, "bot error", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.CreateStreamMessage(context.Background(), "s001", "staff", "stream123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}
