package lansenger

import (
	"context"
	"testing"
)

func TestFetchUserInfo(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/users/fetch", 0, "ok", map[string]interface{}{
			"staffId":        "u001",
			"name":           "Bob",
			"orgId":          "org1",
			"orgName":        "TestOrg",
			"avatar":         "https://avatar.example.com/bob.png",
			"email":          "bob@example.com",
			"employeeNumber": "EMP001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchUserInfo(context.Background(), "utok123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.StaffID != "u001" {
		t.Errorf("expected StaffID=u001, got %s", result.StaffID)
	}
	if result.Name != "Bob" {
		t.Errorf("expected Name=Bob, got %s", result.Name)
	}
	if result.OrgName != "TestOrg" {
		t.Errorf("expected OrgName=TestOrg, got %s", result.OrgName)
	}
	if result.AvatarURL != "https://avatar.example.com/bob.png" {
		t.Errorf("expected AvatarURL=https://avatar.example.com/bob.png, got %s", result.AvatarURL)
	}
}

func TestFetchUserInfoNoToken(t *testing.T) {
	c := NewClient("id", "secret")
	_, err := c.FetchUserInfo(context.Background(), "")
	if err == nil {
		t.Error("expected error for missing userToken")
	}
}