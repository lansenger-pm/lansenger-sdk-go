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
		"Team Meeting", map[string]interface{}{"time": 1000}, map[string]interface{}{"time": 2000},
		nil, "", "no", "", nil, "", "no", "", "utok1", "")
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

func TestFetchSchedule(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal001/schedules/sch001/fetch", 0, "ok", map[string]interface{}{
			"scheduleId": "sch001",
			"summary":    "Team Meeting",
			"startTime":  "2024-01-15T09:00",
			"endTime":    "2024-01-15T10:00",
			"allDay":     "no",
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
	if result.AllDay != "no" {
		t.Errorf("expected AllDay=no, got %s", result.AllDay)
	}
}

func TestDeleteSchedule(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal001/schedules/sch001/delete", 0, "ok", map[string]interface{}{}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DeleteSchedule(context.Background(), "cal001", "sch001", "no", "", "", "utok1", "")
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
		1000, 2000, "utok1", "")
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
	result, err := c.FetchScheduleAttendees(context.Background(), "cal001", "sch001", 1, 10, "utok1", "")
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
	attendees := []string{"s002"}
	result, err := c.AddScheduleAttendees(context.Background(), "cal001", "sch001", attendees, "no", "", "", "utok1", "")
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
	attendees := []string{"s002"}
	result, err := c.DeleteScheduleAttendees(context.Background(), "cal001", "sch001", attendees, "no", "", "", "utok1", "")
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

func TestUpdateScheduleAttendees(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/calendars/cal1/schedules/sch1/members/update", 0, "ok", map[string]interface{}{
			"scheduleIds": []interface{}{"s1", "s2"},
			"attendees":   []interface{}{"failed1"},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.UpdateScheduleAttendees(
		context.Background(), "cal1", "sch1",
		[]string{"a1"}, []string{"a2"},
		"yes", "modify_all", 0, "utok1", "uid1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if len(result.ScheduleIDs) != 2 {
		t.Errorf("expected 2 schedule IDs, got %d", len(result.ScheduleIDs))
	}
	if len(result.FailedAttendees) != 1 {
		t.Errorf("expected 1 failed attendee, got %d", len(result.FailedAttendees))
	}
}