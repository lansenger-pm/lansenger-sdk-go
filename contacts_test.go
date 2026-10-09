package lansenger

import (
	"context"
	"testing"
)

func TestFetchStaffBasicInfo(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/staffs/s001/fetch", 0, "ok", map[string]interface{}{
			"orgId":   "org1",
			"orgName": "TestOrg",
			"name":    "Alice",
			"gender":  "female",
			"avatar":  "https://avatar.example.com/alice.png",
			"status":  "active",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchStaffBasicInfo(context.Background(), "s001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Name != "Alice" {
		t.Errorf("expected Name=Alice, got %s", result.Name)
	}
	if result.OrgID != "org1" {
		t.Errorf("expected OrgID=org1, got %s", result.OrgID)
	}
	if result.OrgName != "TestOrg" {
		t.Errorf("expected OrgName=TestOrg, got %s", result.OrgName)
	}
}

func TestFetchStaffBasicInfoAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/staffs/s002/fetch", 40042, "staff not found", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchStaffBasicInfo(context.Background(), "s002", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
	if result.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestFetchStaffIdMapping(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/staffs/id_mapping/fetch", 0, "ok", map[string]interface{}{
			"staffId": "mapped_staff_001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchStaffIdMapping(context.Background(), "org1", "phone", "13800138000", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.StaffID != "mapped_staff_001" {
		t.Errorf("expected StaffID=mapped_staff_001, got %s", result.StaffID)
	}
}

func TestSearchStaff(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/staffs/search", 0, "ok", map[string]interface{}{
			"hasMore": false,
			"total":   1,
			"staffInfo": []map[string]interface{}{
				{"staffId": "s001", "name": "Alice"},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SearchStaff(context.Background(), "Alice", "", "", true, nil, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Total != 1 {
		t.Errorf("expected Total=1, got %d", result.Total)
	}
	if result.HasMore {
		t.Error("expected HasMore=false")
	}
}

func TestFetchOrgInfo(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/org/org1/fetch", 0, "ok", map[string]interface{}{
			"orgId":   "org1",
			"orgName": "MyOrganization",
			"iconUrl": "https://icon.example.com/org.png",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchOrgInfo(context.Background(), "org1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.OrgID != "org1" {
		t.Errorf("expected OrgID=org1, got %s", result.OrgID)
	}
	if result.OrgName != "MyOrganization" {
		t.Errorf("expected OrgName=MyOrganization, got %s", result.OrgName)
	}
}

// 真实响应把 data 作为数组返回（每个元素带 ancestorDepartments 对象数组），
// 元素内的部门标识是 departmentId 而不是 id —— 旧的 map-only 解析会报 "no data in response"。
func TestFetchDepartmentAncestorsArrayShape(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/staffs/s001/departmentancestors/fetch", 0, "ok", []interface{}{
			map[string]interface{}{
				"ancestorDepartments": []interface{}{
					map[string]interface{}{"departmentId": "d-root", "externalId": "ext-1", "name": "总部"},
					map[string]interface{}{"departmentId": "d-sub", "name": "研发中心"},
				},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchDepartmentAncestors(context.Background(), "s001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected Success=true, got error %q", result.Error)
	}
	if len(result.AncestorGroups) != 1 {
		t.Fatalf("expected 1 ancestor group, got %d", len(result.AncestorGroups))
	}
	group := result.AncestorGroups[0]
	if len(group) != 2 {
		t.Fatalf("expected 2 ancestors, got %d", len(group))
	}
	if group[0]["id"] != "d-root" || group[0]["departmentId"] != "d-root" || group[0]["name"] != "总部" || group[0]["externalId"] != "ext-1" {
		t.Errorf("unexpected first ancestor: %#v", group[0])
	}
	if group[1]["departmentId"] != "d-sub" {
		t.Errorf("expected departmentId=d-sub, got %#v", group[1])
	}
}

// 兼容旧形状：data 是对象、ancestorDepartments 是数组的数组。
func TestFetchDepartmentAncestorsLegacyShape(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/staffs/s002/departmentancestors/fetch", 0, "ok", map[string]interface{}{
			"ancestorDepartments": []interface{}{
				[]interface{}{
					map[string]interface{}{"id": "d-1", "name": "一层"},
				},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchDepartmentAncestors(context.Background(), "s002", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || len(result.AncestorGroups) != 1 || len(result.AncestorGroups[0]) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if result.AncestorGroups[0][0]["id"] != "d-1" {
		t.Errorf("unexpected entry: %#v", result.AncestorGroups[0][0])
	}
}
