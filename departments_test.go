package lansenger

import (
	"context"
	"testing"
)

func TestFetchDepartmentDetail(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/departments/d001/fetch", 0, "ok", map[string]interface{}{
			"id":            "d001",
			"name":          "Engineering",
			"externalId":    "ext_d001",
			"parentId":      "d000",
			"order":         1,
			"hasChildren":   true,
			"normalMembers": 50,
			"deptType":      "normal",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchDepartmentDetail(context.Background(), "d001", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.ID != "d001" {
		t.Errorf("expected ID=d001, got %s", result.ID)
	}
	if result.Name != "Engineering" {
		t.Errorf("expected Name=Engineering, got %s", result.Name)
	}
	if !result.HasChildren {
		t.Error("expected HasChildren=true")
	}
	if result.NormalMembers != 50 {
		t.Errorf("expected NormalMembers=50, got %d", result.NormalMembers)
	}
}

func TestFetchDepartmentChildren(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/departments/d001/children/fetch", 0, "ok", map[string]interface{}{
			"departments": []map[string]interface{}{
				{"id": "d002", "name": "Frontend"},
				{"id": "d003", "name": "Backend"},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchDepartmentChildren(context.Background(), "d001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestFetchDepartmentStaffs(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/departments/d001/staffs/fetch", 0, "ok", map[string]interface{}{
			"hasMore": true,
			"total":   100,
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchDepartmentStaffs(context.Background(), "d001", "", 1, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if !result.HasMore {
		t.Error("expected HasMore=true")
	}
	if result.Total != 100 {
		t.Errorf("expected Total=100, got %d", result.Total)
	}
}

func TestFetchDepartmentDetailAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/departments/d999/fetch", 52000, "department not found", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchDepartmentDetail(context.Background(), "d999", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}
