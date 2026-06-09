package lansenger

import (
	"context"
	"testing"
)

func TestCreateTodoTask(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/create", 0, "ok", map[string]interface{}{
			"todotaskId": "todo001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.CreateTodoTask(context.Background(), "Review doc", TodoTypeNotification,
		"https://link", "https://pclink", []string{"s001", "s002"}, "org1", "", "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.TodotaskID != "todo001" {
		t.Errorf("expected TodotaskID=todo001, got %s", result.TodotaskID)
	}
}

func TestUpdateTodoTaskStatus(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/status/update", 0, "ok", map[string]interface{}{
			"todotaskId": "todo001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.UpdateTodoTaskStatus(context.Background(), "todo001", TodoStatusDone, "org1", "s001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestDeleteTodoTask(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/sender/todotask/delete", 0, "ok", map[string]interface{}{}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DeleteTodoTask(context.Background(), "todo001", "org1", "s001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestFetchTodoTaskList(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/list/fetch", 0, "ok", map[string]interface{}{
			"total":        5,
			"todotaskList": []map[string]interface{}{},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchTodoTaskList(context.Background(), "org1", nil, "", nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Total != 5 {
		t.Errorf("expected Total=5, got %d", result.Total)
	}
}

func TestFetchTodoTaskBySourceID(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/info/fetchbysourceid", 0, "ok", map[string]interface{}{
			"todotaskId": "todo002",
			"sourceId":   "src001",
			"title":      "Task Title",
			"status":     "21",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchTodoTaskBySourceID(context.Background(), "src001", "org1", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.TodotaskID != "todo002" {
		t.Errorf("expected TodotaskID=todo002, got %s", result.TodotaskID)
	}
	if result.SourceID != "src001" {
		t.Errorf("expected SourceID=src001, got %s", result.SourceID)
	}
	if result.Title != "Task Title" {
		t.Errorf("expected Title=Task Title, got %s", result.Title)
	}
}

func TestFetchTodoTaskByID(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/info/fetch", 0, "ok", map[string]interface{}{
			"todotaskId": "todo001",
			"title":      "My Task",
			"status":     "22",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchTodoTaskByID(context.Background(), "todo001", "org1", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.TodotaskID != "todo001" {
		t.Errorf("expected TodotaskID=todo001, got %s", result.TodotaskID)
	}
	if result.Title != "My Task" {
		t.Errorf("expected Title=My Task, got %s", result.Title)
	}
}

func TestFetchTodoTaskStatusCounts(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/status/countList/fetch", 0, "ok", map[string]interface{}{
			"statusCounts": []map[string]interface{}{
				{"status": "21", "count": 3},
				{"status": "22", "count": 5},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchTodoTaskStatusCounts(context.Background(), "s001", "org1", "", nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestAddExecutors(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/executor/create", 0, "ok", map[string]interface{}{
			"todotaskId": "todo001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.AddExecutors(context.Background(), []string{"s003", "s004"}, "org1", "todo001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestDeleteExecutors(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/executor/delete", 0, "ok", map[string]interface{}{}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DeleteExecutors(context.Background(), []string{"s003"}, "org1", "todo001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestFetchExecutorList(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/executor/list/fetch", 0, "ok", map[string]interface{}{
			"total":        2,
			"executorList": []map[string]interface{}{},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchExecutorList(context.Background(), "todo001", "org1", "", nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Total != 2 {
		t.Errorf("expected Total=2, got %d", result.Total)
	}
}

func TestCreateTodoTaskAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/task/unified/v1/todotask/create", 56008, "rate limit", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.CreateTodoTask(context.Background(), "Task", TodoTypeNotification,
		"", "", []string{"s001"}, "org1", "", "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}
