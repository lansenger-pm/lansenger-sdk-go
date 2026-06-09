package lansenger

import (
	"context"
	"testing"
)

func TestFetchChatList(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/chats/fetch", 0, "ok", map[string]interface{}{
			"staffIdInfos": []map[string]interface{}{
				{"staffId": "s001", "staffName": "Alice"},
			},
			"groupIdInfos": []map[string]interface{}{
				{"groupId": "g001", "groupName": "DevGroup"},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchChatList(context.Background(), "utok1", 1, "", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if len(result.StaffInfos) != 1 {
		t.Errorf("expected 1 staff info, got %d", len(result.StaffInfos))
	}
	if result.StaffInfos[0].StaffID != "s001" {
		t.Errorf("expected StaffID=s001, got %s", result.StaffInfos[0].StaffID)
	}
	if len(result.GroupInfos) != 1 {
		t.Errorf("expected 1 group info, got %d", len(result.GroupInfos))
	}
	if result.GroupInfos[0].GroupID != "g001" {
		t.Errorf("expected GroupID=g001, got %s", result.GroupInfos[0].GroupID)
	}
}

func TestFetchChatListNoToken(t *testing.T) {
	c := NewClient("id", "secret")
	_, err := c.FetchChatList(context.Background(), "", 1, "", 0, 0)
	if err == nil {
		t.Error("expected error for missing userToken")
	}
}

func TestFetchChatMessages(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/messages/fetch", 0, "ok", map[string]interface{}{
			"hasMore":     false,
			"total":       2,
			"lastVersion": "v100",
			"name":        "Alice",
			"chatType":    "private",
			"messageList": []map[string]interface{}{
				{
					"name":     "Alice",
					"chatType": "private",
					"messageInfo": map[string]interface{}{
						"sendTime":    "2024-01-01",
						"sender":      "Alice",
						"type":        "text",
						"content":     map[string]interface{}{"text": "hello"},
					},
				},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchChatMessages(context.Background(), "utok1", 10, "", "s001", "", 0, 0, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Total != 2 {
		t.Errorf("expected Total=2, got %d", result.Total)
	}
	if result.ChatType != "private" {
		t.Errorf("expected ChatType=private, got %s", result.ChatType)
	}
	if len(result.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(result.Messages))
	}
	if result.Messages[0].Sender != "Alice" {
		t.Errorf("expected Sender=Alice, got %s", result.Messages[0].Sender)
	}
}

func TestFetchChatMessagesWithPluralMessageInfos(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/messages/fetch", 0, "ok", map[string]interface{}{
			"hasMore":     false,
			"total":       3,
			"lastVersion": "v200",
			"name":        "Bob",
			"chatType":    "group",
			"messageList": []map[string]interface{}{
				{
					"name":      "Bob",
					"chatType":  "group",
					"messageInfos": []map[string]interface{}{
						{
							"sendTime": "2024-02-01",
							"sender":   "Bob",
							"type":     "image",
							"content":  map[string]interface{}{"url": "http://img.example.com/1.jpg"},
						},
					},
				},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchChatMessages(context.Background(), "utok1", 10, "", "", "g001", 0, 0, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Total != 3 {
		t.Errorf("expected Total=3, got %d", result.Total)
	}
	if len(result.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(result.Messages))
	}
	if result.Messages[0].Sender != "Bob" {
		t.Errorf("expected Sender=Bob, got %s", result.Messages[0].Sender)
	}
	if result.Messages[0].MessageType != "image" {
		t.Errorf("expected MessageType=image, got %s", result.Messages[0].MessageType)
	}
}

func TestFetchChatMessagesWithMessageTypeField(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/messages/fetch", 0, "ok", map[string]interface{}{
			"hasMore":     false,
			"total":       1,
			"lastVersion": "v300",
			"name":        "Carol",
			"chatType":    "group",
			"messageList": []map[string]interface{}{
				{
					"name":     "Carol",
					"chatType": "group",
					"messageInfo": map[string]interface{}{
						"sendTime":    "2024-03-01",
						"sender":      "Carol",
						"messageType": "staff_card",
						"content":     map[string]interface{}{"name": "Dave"},
					},
				},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchChatMessages(context.Background(), "utok1", 10, "", "", "g002", 0, 0, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if len(result.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(result.Messages))
	}
	if result.Messages[0].MessageType != "staff_card" {
		t.Errorf("expected MessageType=staff_card, got %s", result.Messages[0].MessageType)
	}
}

func TestFetchChatMessagesNoToken(t *testing.T) {
	c := NewClient("id", "secret")
	_, err := c.FetchChatMessages(context.Background(), "", 10, "", "s001", "", 0, 0, "")
	if err == nil {
		t.Error("expected error for missing userToken")
	}
}

func TestFetchChatListAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/chats/fetch", 56008, "rate limit", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchChatList(context.Background(), "utok1", 1, "", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}
