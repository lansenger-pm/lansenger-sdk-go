package lansenger

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func personalTodoStringHandler(value string) http.HandlerFunc {
	return personalTodoStringHandlerWithCode(0, value)
}

func personalTodoStringHandlerWithCode(code int, value string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"errCode": code, "errMsg": "ok", "data": value})
	}
}

func TestSavePersonalTodoSuccess(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/xtra/tdtask/server/openapi/v3/taskopt/savePersonalTask", personalTodoStringHandler("TASK001"))
	server := b.build()
	defer server.Close()

	result, err := newTestClient(server).SavePersonalTodo(context.Background(), &PersonalTodoSaveParams{
		Subject: "完成项目方案", StartTime: 100, DueTime: 200, Priority: PersonalTodoPriorityNormal,
		CreateUserID: "u1", OrgID: "org1", AppID: "app1", FinishTime: 0,
		Executors: []map[string]interface{}{{"staffId": "u1", "opt": 1}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || result.TodoCode != "TASK001" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSavePersonalTodoValidation(t *testing.T) {
	c := newTestClient(nil)
	result, err := c.SavePersonalTodo(context.Background(), &PersonalTodoSaveParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success || result.Error != "subject is required" {
		t.Fatalf("unexpected validation result: %+v", result)
	}
}

func TestSavePersonalTodoWorkingRequestShape(t *testing.T) {
	var body map[string]interface{}
	var rawURL string
	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/xtra/tdtask/server/openapi/v3/taskopt/savePersonalTask", func(w http.ResponseWriter, r *http.Request) {
		rawURL = r.URL.String()
		_ = json.NewDecoder(r.Body).Decode(&body)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"errCode": 0, "errMsg": "ok", "data": "TASK001"})
	})
	server := b.build()
	defer server.Close()

	result, err := newTestClient(server).SavePersonalTodo(context.Background(), &PersonalTodoSaveParams{
		Subject: "测试待办", StartTime: 1789725600000, DueTime: 1789812000000,
		FinishTime: 0, Priority: 1,
		CreateUserID: "user-001", OrgID: "org-001", AppID: "app-001",
		Executors: []map[string]interface{}{{"staffId": "user-001", "opt": 1}},
		UserToken: "ut1",
	})
	if err != nil || !result.Success {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if body["finishTime"] != float64(0) || body["type"] != float64(1) {
		t.Fatalf("unexpected body: %+v", body)
	}
	if _, ok := body["finishTime"]; !ok {
		t.Fatalf("finishTime must be present: %+v", body)
	}
	if !strings.Contains(rawURL, "app_token=tok1") || !strings.Contains(rawURL, "user_token=ut1") {
		t.Fatalf("unexpected query: %s", rawURL)
	}
}

func TestSavePersonalTodoAcceptsLegacySuccessCode(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc(
		"/xtra/tdtask/server/openapi/v3/taskopt/savePersonalTask",
		personalTodoStringHandlerWithCode(200, "TASK200"),
	)
	server := b.build()
	defer server.Close()

	result, err := newTestClient(server).SavePersonalTodo(context.Background(), &PersonalTodoSaveParams{
		Subject: "兼容旧环境", StartTime: 100, DueTime: 200, Priority: PersonalTodoPriorityNormal,
		CreateUserID: "u1", OrgID: "org1", AppID: "app1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || result.TodoCode != "TASK200" {
		t.Fatalf("expected errCode=200 to be treated as success, got %+v", result)
	}
}

func TestUpdatePersonalTodoTopLevelOrgID(t *testing.T) {
	var body map[string]interface{}
	server := newMuxBuilder().handleToken("tok1")
	server.mux.HandleFunc("/xtra/tdtask/server/openapi/v3/taskopt/updatePersonalTask", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"errCode": 0, "errMsg": "ok", "data": "TASK001"})
	})
	s := server.build()
	defer s.Close()

	result, err := newTestClient(s).UpdatePersonalTodo(context.Background(), &PersonalTodoUpdateParams{
		TodoCode: "TASK001", OrgID: "org1", UpdateFields: []string{"subject"},
		Subject: "新主题", CreateUserID: "u1", AppID: "app1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Fatalf("unexpected result: %+v", result)
	}
	if body["orgId"] != "org1" {
		t.Fatalf("expected top-level orgId, got %+v", body)
	}
	content, _ := body["updateContent"].(map[string]interface{})
	if content["subject"] != "新主题" || content["orgId"] != nil {
		t.Fatalf("unexpected update content: %+v", content)
	}
	if content["createUserId"] != "u1" || content["appid"] != "app1" {
		t.Fatalf("identity fields missing from update content: %+v", content)
	}
}

func TestFetchPersonalTodoListSuccess(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/tdtask/server/openapi/v3/user/list", 0, "ok", map[string]interface{}{
			"pageNo": 2, "pageSize": 10, "pages": 3, "total": 21, "hasNextPage": true,
			"result": []interface{}{map[string]interface{}{"taskCode": "TASK001", "status": 0}},
		}).
		build()
	defer server.Close()

	status := PersonalTodoStatusUnfinished
	result, err := newTestClient(server).FetchPersonalTodoList(
		context.Background(), "org1", "u1", 2, 10, &status, "", "", "",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || result.PageNo != 2 || result.Total != 21 || !result.HasMore {
		t.Fatalf("unexpected result: %+v", result)
	}
	if len(result.Items) != 1 || result.Items[0]["taskCode"] != "TASK001" {
		t.Fatalf("unexpected items: %+v", result.Items)
	}
}

func TestPersonalTodoAPIAndHTTPErrors(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/tdtask/server/openapi/v3/user/list", 3124, "日期格式错误", nil).
		build()
	defer server.Close()

	result, err := newTestClient(server).FetchPersonalTodoList(
		context.Background(), "org1", "u1", 1, 10, nil, "", "", "",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success || result.Error == "" {
		t.Fatalf("expected API error, got %+v", result)
	}

	httpServer := newMuxBuilder().
		handleToken("tok1").
		handleError("/xtra/tdtask/server/openapi/v3/user/list").
		build()
	defer httpServer.Close()
	result, err = newTestClient(httpServer).FetchPersonalTodoList(
		context.Background(), "org1", "u1", 1, 10, nil, "", "", "",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Fatalf("expected HTTP error result: %+v", result)
	}
}

func TestPersonalTodoResourceOperations(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/xtra/tdtask/server/openapi/resource/update", mockAPIResponseHandler(0, "ok", map[string]interface{}{
		"fileName": "a.pdf", "mimeType": "application/pdf", "size": 10, "resourceId": "res1",
	}))
	b.mux.HandleFunc("/xtra/tdtask/server/openapi/resource/getResourceDownload", personalTodoStringHandler("https://example.com/download"))
	b.mux.HandleFunc("/xtra/tdtask/server/openapi/resource/getUploadUrl", personalTodoStringHandler("https://example.com/upload"))
	server := b.build()
	defer server.Close()
	c := newTestClient(server)

	upload, err := c.UploadPersonalTodoResource(context.Background(), &PersonalTodoResourceUploadParams{
		AppID: "app1", Size: 10, FileName: "a.pdf", ContentType: "application/pdf",
		FileData: "YWJj", OrgID: "org1",
	})
	if err != nil || !upload.Success || upload.ResourceID != "res1" {
		t.Fatalf("unexpected upload result: %+v err=%v", upload, err)
	}

	download, err := c.FetchPersonalTodoResourceDownloadURL(context.Background(), "res1", "org1", "", "")
	if err != nil || !download.Success || download.URL == "" {
		t.Fatalf("unexpected download result: %+v err=%v", download, err)
	}

	uploadURL, err := c.FetchPersonalTodoResourceUploadURL(context.Background(), "a.pdf", "md5", 10, "org1", "")
	if err != nil || !uploadURL.Success || uploadURL.URL == "" {
		t.Fatalf("unexpected upload URL result: %+v err=%v", uploadURL, err)
	}
}

func TestPersonalTodoResourceGuards(t *testing.T) {
	c := newTestClient(nil)
	result, err := c.UploadPersonalTodoResource(context.Background(), &PersonalTodoResourceUploadParams{
		AppID: "app1", Size: PersonalTodoResourceMaxSize + 1, FileName: "a.pdf",
		ContentType: "application/pdf", FileData: "x", OrgID: "org1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success || result.Error == "" {
		t.Fatalf("expected size guard, got %+v", result)
	}
}

func TestPersonalTodoResourceEntryFromUpload(t *testing.T) {
	// 上传响应用 mimeType/size，挂附件必须叫 fileType/fileSize
	want := map[string]interface{}{
		"fileName":   "a.pdf",
		"resourceId": "res1",
		"fileType":   "application/pdf",
		"fileSize":   int64(10),
		"opt":        1,
	}
	raw := map[string]interface{}{
		"errCode": 0,
		"data": map[string]interface{}{
			"fileName": "a.pdf", "mimeType": "application/pdf",
			"size": float64(10), "resourceId": "res1",
		},
	}

	t.Run("accepts the upload result object", func(t *testing.T) {
		res := &PersonalTodoResourceResult{
			Success: true, FileName: "a.pdf", MimeType: "application/pdf",
			Size: 10, ResourceID: "res1",
		}
		got, err := PersonalTodoResourceEntryFromUpload(res, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	})

	t.Run("accepts the raw response and its inner data", func(t *testing.T) {
		for _, in := range []interface{}{raw, raw["data"]} {
			got, err := PersonalTodoResourceEntryFromUpload(in, 1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v, want %+v", got, want)
			}
		}
	})

	t.Run("agrees with the explicit constructor", func(t *testing.T) {
		want0 := BuildPersonalTodoResourceEntry("res1", "a.pdf", "application/pdf", 10, 0)
		got, err := PersonalTodoResourceEntryFromUpload(&PersonalTodoResourceResult{
			Success: true, FileName: "a.pdf", MimeType: "application/pdf",
			Size: 10, ResourceID: "res1",
		}, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, want0) || got["opt"] != 0 {
			t.Fatalf("got %+v, want %+v", got, want0)
		}
	})

	t.Run("rejects input with no usable payload", func(t *testing.T) {
		for _, in := range []interface{}{nil, 42, "res1", (*PersonalTodoResourceResult)(nil)} {
			if _, err := PersonalTodoResourceEntryFromUpload(in, 1); err == nil {
				t.Fatalf("expected error for %T, got nil", in)
			}
		}
	})
}
