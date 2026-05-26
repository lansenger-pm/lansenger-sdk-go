package lansenger

import (
	"context"
	"testing"
)

func TestSendAccountMessage(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/messages/create", 0, "ok", map[string]interface{}{
			"msgId":             "msg001",
			"invalidStaff":      []string{},
			"invalidDepartment": []string{},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SendAccountMessage(context.Background(), "text",
		map[string]interface{}{"content": "hello"},
		[]string{"s001"}, nil, "acct1", "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.MessageID != "msg001" {
		t.Errorf("expected MessageID=msg001, got %s", result.MessageID)
	}
}

func TestSendUserMessage(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/messages/chat/create", 0, "ok", map[string]interface{}{
			"msgId": "msg002",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SendUserMessage(context.Background(), "r001", "text",
		map[string]interface{}{"content": "hello"}, "utok1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.MessageID != "msg002" {
		t.Errorf("expected MessageID=msg002, got %s", result.MessageID)
	}
}

func TestSendUserMessageNoToken(t *testing.T) {
	c := NewClient("id", "secret")
	_, err := c.SendUserMessage(context.Background(), "r001", "text",
		map[string]interface{}{"content": "hello"}, "", "")
	if err == nil {
		t.Error("expected error for missing userToken")
	}
}

func TestSendGroupMessage(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/messages/group/create", 0, "ok", map[string]interface{}{
			"msgId": "msg003",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SendGroupMessage(context.Background(), "grp001", "text",
		map[string]interface{}{"content": "hello"}, "utok1", "s001",
		"", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.MessageID != "msg003" {
		t.Errorf("expected MessageID=msg003, got %s", result.MessageID)
	}
	if result.MsgType != "text" {
		t.Errorf("expected MsgType=text, got %s", result.MsgType)
	}
}

func TestSendBotMessage(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/bot/messages/create", 0, "ok", map[string]interface{}{
			"msgId":             "msg004",
			"invalidStaff":      []string{},
			"invalidDepartment": []string{},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SendBotMessage(context.Background(), "text",
		map[string]interface{}{"content": "hello"},
		[]string{"s001"}, nil, "", "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.MessageID != "msg004" {
		t.Errorf("expected MessageID=msg004, got %s", result.MessageID)
	}
}

func TestSendGroupMessageAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/messages/group/create", 51000, "group error", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SendGroupMessage(context.Background(), "grp001", "text",
		map[string]interface{}{"content": "hello"}, "", "",
		"", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}
