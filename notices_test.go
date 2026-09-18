package lansenger

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestSendNoticePhoneSuccess(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/notice/server/openapi/v1/send", 0, "ok", map[string]interface{}{
			"code":            "NTC001",
			"id":              1001,
			"noticeStatus":    2,
			"confirmStatus":   0,
			"publishUserName": "张三",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	confirm := 1
	result, err := c.SendNotice(context.Background(), &NoticeSendParams{
		Title:         "系统升级通知",
		ContentType:   1,
		AccountCode:   "ACC001",
		UserType:      1,
		Content:       "系统将于本周六进行升级维护",
		ReleasePhones: []string{"13800138000", "13800138001"},
		CCPhones:      []string{"13800138002"},
		CreateMobile:  "13800138000",
		ConfirmFlag:   &confirm,
		UserToken:     "utok",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if result.NoticeCode != "NTC001" {
		t.Errorf("expected NoticeCode=NTC001, got %s", result.NoticeCode)
	}
	if result.NoticeStatus != 2 {
		t.Errorf("expected NoticeStatus=2, got %d", result.NoticeStatus)
	}
	if result.PublishUserName != "张三" {
		t.Errorf("expected PublishUserName=张三, got %s", result.PublishUserName)
	}
}

func TestSendNoticeOpenIDSuccess(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/notice/server/openapi/v1/send", 0, "ok", map[string]interface{}{
			"code": "NTC002",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SendNotice(context.Background(), &NoticeSendParams{
		Title:       "部门通知",
		ContentType: 1,
		AccountCode: "ACC001",
		UserType:    2,
		Content:     "请及时填写本周周报",
		ReleaseRange: []map[string]interface{}{
			{"objId": "dept-1", "objName": "研发部", "objType": 2},
		},
		CCStaffIDs:   []string{"staff-002"},
		CreateUserID: "staff-001",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if result.NoticeCode != "NTC002" {
		t.Errorf("expected NoticeCode=NTC002, got %s", result.NoticeCode)
	}
}

func TestSendNoticeDefaultRemindStatus(t *testing.T) {
	// 实测：缺失 remindStatus 服务端抛 errCode=-1，SDK 必须兜底填 0。
	var gotBody map[string]interface{}
	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/xtra/notice/server/openapi/v1/send", func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		json.Unmarshal(raw, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errCode": 0, "errMsg": "ok",
			"data": map[string]interface{}{"code": "NTC_D"},
		})
	})
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SendNotice(context.Background(), &NoticeSendParams{
		Title: "t", ContentType: 1, AccountCode: "ACC001", UserType: 1,
		Content: "c", ReleasePhones: []string{"13800138000"}, CreateMobile: "13800138000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if gotBody["remindStatus"] != float64(0) {
		t.Errorf("expected remindStatus=0 in body, got %v", gotBody["remindStatus"])
	}
	// 实测：ccRangeList 缺失同样触发服务端 NPE，必须强制下发
	pr, _ := gotBody["phoneUserRange"].(map[string]interface{})
	if pr == nil || pr["ccRangeList"] == nil {
		t.Errorf("expected phoneUserRange.ccRangeList to be present, got %v", gotBody["phoneUserRange"])
	}
}

func TestSendNoticeAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/notice/server/openapi/v1/send", 3123, "人员不存在", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SendNotice(context.Background(), &NoticeSendParams{
		Title:        "t",
		ContentType:  1,
		AccountCode:  "ACC001",
		UserType:     1,
		Content:      "c",
		ReleasePhones: []string{"13800138000"},
		CreateMobile: "13800138000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false on API error")
	}
	if result.Error == "" {
		t.Error("expected error message to be set")
	}
}

func TestFetchNoticeAccountsSuccess(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	// The accounts endpoint returns `data` as a JSON array — the generic
	// handle() helper only accepts map payloads, so register a raw handler.
	b.mux.HandleFunc("/xtra/notice/server/openapi/v1/notice/account", mockJSONArrayHandler(0, "ok", []interface{}{
		map[string]interface{}{"roleName": "行政通知", "officialNumberId": "1001", "code": "ACC001"},
		map[string]interface{}{"roleName": "安全通知", "officialNumberId": "1002", "code": "ACC002"},
	}))
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchNoticeAccounts(context.Background(), "org-001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if result.Total != 2 {
		t.Errorf("expected Total=2, got %d", result.Total)
	}
	if len(result.Accounts) != 2 || result.Accounts[0]["code"] != "ACC001" {
		t.Errorf("unexpected accounts: %+v", result.Accounts)
	}
}
