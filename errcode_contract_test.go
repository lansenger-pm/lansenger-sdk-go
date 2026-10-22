package lansenger

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// LXBUGS-128629：失败结果必须带结构化 ErrCode，调用方不再正则解析 error 字符串。

func TestErrCodeOf(t *testing.T) {
	if got := ErrCodeOf(nil); got != 0 {
		t.Errorf("nil error: want 0, got %d", got)
	}
	if got := ErrCodeOf(NewNetworkError("boom")); got != 0 {
		t.Errorf("network error: want 0, got %d", got)
	}
	if got := ErrCodeOf(NewAPIError("nope", 59003)); got != 59003 {
		t.Errorf("api error: want 59003, got %d", got)
	}
	if got := ErrCodeOf(NewFileError("no file")); got != 0 {
		t.Errorf("file error: want 0, got %d", got)
	}
}

func TestFetchUserInfoCarriesErrCode(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.handle("/v1/users/fetch", 10001, "用户不存在", nil)
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	res, err := c.FetchUserInfo(context.Background(), "u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success {
		t.Fatal("expected failure result")
	}
	if res.ErrCode != 10001 {
		t.Errorf("want ErrCode=10001, got %d", res.ErrCode)
	}
	if res.Error == "" {
		t.Error("expected human-readable error kept")
	}
	// 编码后仍是结构化字段
	raw, _ := json.Marshal(res)
	var decoded map[string]interface{}
	_ = json.Unmarshal(raw, &decoded)
	if decoded["err_code"] != float64(10001) {
		t.Errorf("marshalled err_code: %v", decoded["err_code"])
	}
}

func TestLocalValidationHasNoErrCode(t *testing.T) {
	c := NewClientWithConfig(NewConfig("app", "secret"))
	res, err := c.RevokeMessage(context.Background(), []string{}, "bot", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success || res.ErrCode != 0 {
		t.Errorf("want local validation failure without code, got success=%v errCode=%d", res.Success, res.ErrCode)
	}
}

// LXBUGS-128628：撤回的 senderId 契约与逐条 subStatus 语义。

func TestRevokeGroupWithoutSenderIDExplainsBotLimitation(t *testing.T) {
	c := NewClientWithConfig(NewConfig("app", "secret"))
	res, err := c.RevokeMessage(context.Background(), []string{"m1"}, "group", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success {
		t.Fatal("expected failure")
	}
	if res.ErrCode != 0 {
		t.Errorf("local validation must not carry a server code, got %d", res.ErrCode)
	}
	for _, want := range []string{"staffId", "user_token/sender_id"} {
		if !strings.Contains(res.Error, want) {
			t.Errorf("error should mention %q, got: %s", want, res.Error)
		}
	}
}

func TestRevokeRejectsUnknownChatType(t *testing.T) {
	c := NewClientWithConfig(NewConfig("app", "secret"))
	res, err := c.RevokeMessage(context.Background(), []string{"m1"}, "channel", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success || !strings.Contains(res.Error, "unsupported chat_type") {
		t.Errorf("want unsupported chat_type error, got success=%v err=%s", res.Success, res.Error)
	}
}

func TestRevokeTopLevelErrorCarriesErrCode(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.handle("/v1/messages/revoke", 59003, "当前组织禁用机器人", nil)
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	res, err := c.RevokeMessage(context.Background(), []string{"m1"}, "bot", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success {
		t.Fatal("expected failure")
	}
	if res.ErrCode != 59003 {
		t.Errorf("want ErrCode=59003, got %d", res.ErrCode)
	}
}

func TestRevokePartialSubStatusFailure(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.handle("/v1/messages/revoke", 0, "ok", map[string]interface{}{
		"subStatus": []interface{}{
			map[string]interface{}{"errCode": 0, "errMsg": "ok"},
			map[string]interface{}{"errCode": 40060, "errMsg": "消息或格式错误"},
		},
	})
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	res, err := c.RevokeMessage(context.Background(), []string{"m1", "m2"}, "bot", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success {
		t.Fatal("partial failure must not report success")
	}
	if res.ErrCode != 40060 {
		t.Errorf("want first failing code 40060, got %d", res.ErrCode)
	}
	if res.SubStatusFailed != 1 || len(res.SubStatus) != 2 {
		t.Errorf("want 1/2 failed, got failed=%d total=%d", res.SubStatusFailed, len(res.SubStatus))
	}
}

func TestRevokeAllSubStatusOK(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.handle("/v1/messages/revoke", 0, "ok", map[string]interface{}{
		"subStatus": []interface{}{map[string]interface{}{"errCode": 0, "errMsg": "ok"}},
	})
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	res, err := c.RevokeMessage(context.Background(), []string{"m1"}, "bot", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success || res.SubStatusFailed != 0 || len(res.SubStatus) != 1 {
		t.Errorf("want success with 1 sub status, got success=%v failed=%d total=%d", res.Success, res.SubStatusFailed, len(res.SubStatus))
	}
}
