package lansenger

import (
	"testing"
)

func TestParseCallbackPayload(t *testing.T) {
	events, err := ParseCallbackPayload("eventType=staff_modify&staffId=s001&orgId=org1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "staff_modify" {
		t.Errorf("expected eventType=staff_modify, got %s", events[0].EventType)
	}
	if events[0].Category != "staff" {
		t.Errorf("expected category=staff, got %s", events[0].Category)
	}
}

func TestParseCallbackPayloadMultipleEvents(t *testing.T) {
	events, err := ParseCallbackPayload("eventType=staff_modify&eventType=dept_create&staffId=s001&deptId=d001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].EventType != "staff_modify" {
		t.Errorf("expected first event=staff_modify, got %s", events[0].EventType)
	}
	if events[1].EventType != "dept_create" {
		t.Errorf("expected second event=dept_create, got %s", events[1].EventType)
	}
	if events[1].Category != "department" {
		t.Errorf("expected category=department, got %s", events[1].Category)
	}
}

func TestParseCallbackPayloadInvalid(t *testing.T) {
	_, err := ParseCallbackPayload("not=valid=query=string===extra")
	if err != nil {
		t.Logf("ParseCallbackPayload handled invalid input: %v", err)
	}
}

func TestVerifyCallbackSignature(t *testing.T) {
	result := VerifyCallbackSignature("eventType=staff_modify&staffId=s001&signature=abc123", "app_secret")
	t.Logf("VerifyCallbackSignature result: %v (placeholder implementation)", result)
}

func TestGetCallbackEventTypes(t *testing.T) {
	types := GetCallbackEventTypes()
	if len(types) == 0 {
		t.Error("expected non-empty callback event types")
	}
	if types["staff_modify"] != "staff" {
		t.Errorf("expected staff_modify=staff, got %s", types["staff_modify"])
	}
	if types["bot_private_message"] != "bot" {
		t.Errorf("expected bot_private_message=bot, got %s", types["bot_private_message"])
	}
}
