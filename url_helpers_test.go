package lansenger

import (
	"net/url"
	"strings"
	"testing"
)

func TestBuildAPIURLBasic(t *testing.T) {
	cfg := NewConfig("app1", "secret1")
	u := BuildAPIURL(cfg, "staffs", "basic_info_fetch", "tok123",
		WithPathVar("staff_id", "s001"),
	)
	if u == "" {
		t.Fatal("expected non-empty URL")
	}
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}
	if !strings.Contains(parsed.Path, "/v1/staffs/s001/fetch") {
		t.Errorf("expected path to contain /v1/staffs/s001/fetch, got %s", parsed.Path)
	}
	if parsed.Query().Get("app_token") != "tok123" {
		t.Errorf("expected app_token=tok123, got %s", parsed.Query().Get("app_token"))
	}
}

func TestBuildAPIURLWithUserToken(t *testing.T) {
	cfg := NewConfig("app1", "secret1")
	u := BuildAPIURL(cfg, "users", "fetch", "atok",
		WithUserToken("utok"),
	)
	parsed, _ := url.Parse(u)
	if parsed.Query().Get("app_token") != "atok" {
		t.Errorf("expected app_token=atok, got %s", parsed.Query().Get("app_token"))
	}
	if parsed.Query().Get("user_token") != "utok" {
		t.Errorf("expected user_token=utok, got %s", parsed.Query().Get("user_token"))
	}
}

func TestBuildAPIURLWithMultiplePathVars(t *testing.T) {
	cfg := NewConfig("app1", "secret1")
	u := BuildAPIURL(cfg, "calendars", "schedules_fetch", "tok",
		WithPathVar("calendar_id", "cal1"),
		WithPathVar("schedule_id", "sch1"),
		WithUserToken("utok"),
		WithUserID("uid1"),
	)
	parsed, _ := url.Parse(u)
	if !strings.Contains(parsed.Path, "/v1/calendars/cal1/schedules/sch1/fetch") {
		t.Errorf("expected path with substituted vars, got %s", parsed.Path)
	}
	if parsed.Query().Get("user_id") != "uid1" {
		t.Errorf("expected user_id=uid1, got %s", parsed.Query().Get("user_id"))
	}
}

func TestBuildAPIURLWithMediaType(t *testing.T) {
	cfg := NewConfig("app1", "secret1")
	u := BuildAPIURL(cfg, "medias", "create", "tok",
		WithMediaType(MediaTypeImage),
	)
	parsed, _ := url.Parse(u)
	if parsed.Query().Get("type") != "2" {
		t.Errorf("expected type=2, got %s", parsed.Query().Get("type"))
	}
}

func TestBuildAPIURLUnknownEndpoint(t *testing.T) {
	cfg := NewConfig("app1", "secret1")
	u := BuildAPIURL(cfg, "nonexistent", "unknown", "tok")
	if u != "" {
		t.Errorf("expected empty URL for unknown endpoint, got %s", u)
	}
}
