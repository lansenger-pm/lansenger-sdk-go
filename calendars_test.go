package lansenger

import (
	"context"
	"testing"
)

func TestFetchPrimaryCalendar(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/primary", 0, "ok", map[string]interface{}{
			"calendarId":  "cal001",
			"summary":     "My Calendar",
			"description": "Work calendar",
			"permissions": "owner",
			"color":       "#FF0000",
			"type":        "primary",
			"role":        "owner",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchPrimaryCalendar(context.Background(), "utok1", "uid1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.CalendarID != "cal001" {
		t.Errorf("expected CalendarID=cal001, got %s", result.CalendarID)
	}
	if result.Summary != "My Calendar" {
		t.Errorf("expected Summary=My Calendar, got %s", result.Summary)
	}
	if result.Role != "owner" {
		t.Errorf("expected Role=owner, got %s", result.Role)
	}
}

func TestFetchPrimaryCalendarNoToken(t *testing.T) {
	c := NewClient("id", "secret")
	_, err := c.FetchPrimaryCalendar(context.Background(), "", "uid1")
	if err == nil {
		t.Error("expected error for missing userToken")
	}
}

func TestCreateSchedule(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal001/schedules/create", 0, "ok", map[string]interface{}{
			"scheduleId": "sch001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.CreateSchedule(context.Background(), "cal001",
		"Team Meeting", "2024-01-15T09:00", "2024-01-15T10:00",
		nil, "", false, "", nil, 0, 0, 0, "utok1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.ScheduleID != "sch001" {
		t.Errorf("expected ScheduleID=sch001, got %s", result.ScheduleID)
	}
}

func TestCreateScheduleNoToken(t *testing.T) {
	c := NewClient("id", "secret")
	_, err := c.CreateSchedule(context.Background(), "cal001",
		"Meeting", "2024-01-15T09:00", "2024-01-15T10:00",
		nil, "", false, "", nil, 0, 0, 0, "")
	if err == nil {
		t.Error("expected error for missing userToken")
	}
}

func TestFetchSchedule(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal001/schedules/sch001/fetch", 0, "ok", map[string]interface{}{
			"scheduleId": "sch001",
			"summary":    "Team Meeting",
			"startTime":  "2024-01-15T09:00",
			"endTime":    "2024-01-15T10:00",
			"allDay":     false,
			"creator":    "s001",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchSchedule(context.Background(), "cal001", "sch001", "utok1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.ScheduleID != "sch001" {
		t.Errorf("expected ScheduleID=sch001, got %s", result.ScheduleID)
	}
	if result.Summary != "Team Meeting" {
		t.Errorf("expected Summary=Team Meeting, got %s", result.Summary)
	}
	if result.AllDay {
		t.Error("expected AllDay=false")
	}
}

func TestDeleteSchedule(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal001/schedules/sch001/delete", 0, "ok", map[string]interface{}{}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DeleteSchedule(context.Background(), "cal001", "sch001", 1, "", "", "utok1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestFetchScheduleList(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal001/schedules/fetch", 0, "ok", map[string]interface{}{
			"scheduleList": []map[string]interface{}{
				{"scheduleId": "sch001", "summary": "Meeting1"},
				{"scheduleId": "sch002", "summary": "Meeting2"},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchScheduleList(context.Background(), "cal001",
		"2024-01-15T00:00", "2024-01-15T23:59", "utok1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if len(result.ScheduleList) != 2 {
		t.Errorf("expected 2 schedules, got %d", len(result.ScheduleList))
	}
}

func TestFetchScheduleAttendees(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal001/schedules/sch001/members/fetch", 0, "ok", map[string]interface{}{
			"total":     3,
			"attendees": []map[string]interface{}{},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchScheduleAttendees(context.Background(), "cal001", "sch001", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Total != 3 {
		t.Errorf("expected Total=3, got %d", result.Total)
	}
}

func TestAddScheduleAttendees(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal001/schedules/sch001/members/create", 0, "ok", map[string]interface{}{}).
		build()
	defer server.Close()

	c := newTestClient(server)
	attendees := []map[string]interface{}{
		{"staffId": "s002"},
	}
	result, err := c.AddScheduleAttendees(context.Background(), "cal001", "sch001", attendees, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestDeleteScheduleAttendees(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal001/schedules/sch001/members/delete", 0, "ok", map[string]interface{}{}).
		build()
	defer server.Close()

	c := newTestClient(server)
	attendees := []map[string]interface{}{
		{"staffId": "s002"},
	}
	result, err := c.DeleteScheduleAttendees(context.Background(), "cal001", "sch001", attendees, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestFetchPrimaryCalendarAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/primary", 61100, "calendar error", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchPrimaryCalendar(context.Background(), "utok1", "uid1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}
