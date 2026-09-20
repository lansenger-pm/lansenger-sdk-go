package lansenger

import (
	"context"
	"encoding/json"
	"net/http"
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

	finish := int64(0)
	result, err := newTestClient(server).SavePersonalTodo(context.Background(), &PersonalTodoSaveParams{
		Subject: "完成项目方案", StartTime: 100, DueTime: 200, Priority: PersonalTodoPriorityNormal,
		CreateUserID: "u1", OrgID: "org1", AppID: "app1", FinishTime: &finish,
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
		Subject: "新主题",
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
