package lansenger

import (
	"testing"
)

func TestMediaTypeConstants(t *testing.T) {
	if MediaTypeVideo != 1 {
		t.Errorf("expected MediaTypeVideo=1, got %d", MediaTypeVideo)
	}
	if MediaTypeImage != 2 {
		t.Errorf("expected MediaTypeImage=2, got %d", MediaTypeImage)
	}
	if MediaTypeAudio != 3 {
		t.Errorf("expected MediaTypeAudio=3, got %d", MediaTypeAudio)
	}
}

func TestGuessMediaTypeImage(t *testing.T) {
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".webp", ".gif"} {
		result := GuessMediaType("photo" + ext)
		if result != MediaTypeImage {
			t.Errorf("expected MediaTypeImage for %s, got %d", ext, result)
		}
	}
}

func TestGuessMediaTypeVideo(t *testing.T) {
	for _, ext := range []string{".mp4", ".mov", ".avi", ".mkv"} {
		result := GuessMediaType("video" + ext)
		if result != MediaTypeVideo {
			t.Errorf("expected MediaTypeVideo for %s, got %d", ext, result)
		}
	}
}

func TestGuessMediaTypeFile(t *testing.T) {
	// Non-media files return 0 — callers should fall back to their own default
	result := GuessMediaType("document.pdf")
	if result != 0 {
		t.Errorf("expected 0 for .pdf, got %d", result)
	}
	result = GuessMediaType("data.xlsx")
	if result != 0 {
		t.Errorf("expected 0 for .xlsx, got %d", result)
	}
}

func TestAPIEndpointsStructure(t *testing.T) {
	categories := []string{"auth", "oauth", "users", "staffs", "org", "departments",
		"groups", "messages", "bot", "sse", "medias", "chats",
		"calendars", "websocket"}
	for _, cat := range categories {
		if APIEndpoints[cat] == nil {
			t.Errorf("expected APIEndpoints category '%s' to exist", cat)
		}
	}
}

func TestStaffsEndpointPaths(t *testing.T) {
	tests := map[string]string{
		"basic_info_fetch":           "/v1/staffs/{staff_id}/fetch",
		"detail_fetch":               "/v1/staffs/{staff_id}/infor/fetch",
		"department_ancestors_fetch": "/v1/staffs/{staff_id}/departmentancestors/fetch",
		"id_mapping_fetch":           "/v2/staffs/id_mapping/fetch",
		"search":                     "/v2/staffs/search",
	}
	for key, expected := range tests {
		actual := APIEndpoints["staffs"][key]
		if actual != expected {
			t.Errorf("expected staffs.%s=%s, got %s", key, expected, actual)
		}
	}
}

func TestGroupsEndpointPaths(t *testing.T) {
	tests := map[string]string{
		"create":         "/v2/groups/create",
		"info_fetch":     "/v2/groups/{group_id}/info/fetch",
		"members_fetch":  "/v2/groups/{group_id}/members/fetch",
		"list_fetch":     "/v2/groups/fetch",
		"is_in_group":    "/v2/groups/{group_id}/members/is_in_group",
		"info_update":    "/v2/groups/{group_id}/info/update",
		"members_update": "/v2/groups/{group_id}/members/update",
	}
	for key, expected := range tests {
		actual := APIEndpoints["groups"][key]
		if actual != expected {
			t.Errorf("expected groups.%s=%s, got %s", key, expected, actual)
		}
	}
}

func TestCallbackEventTypes(t *testing.T) {
	if len(CallbackEventTypes) == 0 {
		t.Error("expected CallbackEventTypes to have entries")
	}
	tests := map[string]string{
		"account_subscribe":    "public_account",
		"staff_modify":         "staff",
		"dept_create":          "department",
		"bot_private_message":  "bot",
		"group_create_approve": "group",
	}
	for eventType, category := range tests {
		if CallbackEventTypes[eventType] != category {
			t.Errorf("expected CallbackEventTypes[%s]=%s, got %s",
				eventType, category, CallbackEventTypes[eventType])
		}
	}
}
