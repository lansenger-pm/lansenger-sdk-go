package lansenger

import (
	"context"
	"testing"
)

func TestCreateGroup(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/groups/create", 0, "ok", map[string]interface{}{
			"groupId":           "grp001",
			"totalMembers":      3,
			"invalidStaff":      []string{},
			"invalidDepartment": []string{},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	info := &GroupCreateInfo{
		Name:        "TestGroup",
		OrgID:       1,
		OwnerID:     "s001",
		StaffIDList: []string{"s001", "s002", "s003"},
	}
	result, err := c.CreateGroup(context.Background(), info, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.GroupID != "grp001" {
		t.Errorf("expected GroupID=grp001, got %s", result.GroupID)
	}
	if result.TotalMembers != 3 {
		t.Errorf("expected TotalMembers=3, got %d", result.TotalMembers)
	}
}

func TestFetchGroupInfo(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/groups/grp001/info/fetch", 0, "ok", map[string]interface{}{
			"name":         "TestGroup",
			"description":  "A test group",
			"owner":        "s001",
			"creator":      "s002",
			"state":        "normal",
			"manageMode":   "owner_manage",
			"totalMembers": 10,
			"avatarUrl":    "https://avatar.example.com/group.png",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchGroupInfo(context.Background(), "grp001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Name != "TestGroup" {
		t.Errorf("expected Name=TestGroup, got %s", result.Name)
	}
	if result.Owner != "s001" {
		t.Errorf("expected Owner=s001, got %s", result.Owner)
	}
	if result.TotalMembers != 10 {
		t.Errorf("expected TotalMembers=10, got %d", result.TotalMembers)
	}
}

func TestFetchGroupMembers(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/groups/grp001/members/fetch", 0, "ok", map[string]interface{}{
			"totalMembers": 3,
			"members": []map[string]interface{}{
				{"staffId": "s001", "name": "Alice"},
				{"staffId": "s002", "name": "Bob"},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchGroupMembers(context.Background(), "grp001", "", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.TotalMembers != 3 {
		t.Errorf("expected TotalMembers=3, got %d", result.TotalMembers)
	}
}

func TestFetchGroupList(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/groups/fetch", 0, "ok", map[string]interface{}{
			"totalGroupIds": 2,
			"groupIds":      []string{"grp001", "grp002"},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchGroupList(context.Background(), "", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.TotalGroupIDs != 2 {
		t.Errorf("expected TotalGroupIDs=2, got %d", result.TotalGroupIDs)
	}
}

func TestCheckIsInGroupTrue(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/groups/grp001/members/is_in_group", 0, "ok", map[string]interface{}{
			"isInGroup": true,
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.CheckIsInGroup(context.Background(), "grp001", "", "s001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if !result.IsInGroup {
		t.Error("expected IsInGroup=true")
	}
}

func TestCheckIsInGroupFalse(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/groups/grp001/members/is_in_group", 0, "ok", map[string]interface{}{
			"isInGroup": false,
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.CheckIsInGroup(context.Background(), "grp001", "", "s999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsInGroup {
		t.Error("expected IsInGroup=false")
	}
}

func TestCreateGroupAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/groups/create", 51000, "group error", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	info := &GroupCreateInfo{Name: "TestGroup", OrgID: 1}
	result, err := c.CreateGroup(context.Background(), info, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}
