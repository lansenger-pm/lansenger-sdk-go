package main

import "testing"

func TestQuestionnaireUserIDFallsBackToAsStaffID(t *testing.T) {
	previous := globalAsStaffID
	defer func() {
		globalAsStaffID = previous
	}()

	globalAsStaffID = "staff-from-as"
	if got := questionnaireUserID(""); got != "staff-from-as" {
		t.Fatalf("expected staff-from-as, got %q", got)
	}
	if got := questionnaireUserID("explicit-staff"); got != "explicit-staff" {
		t.Fatalf("expected explicit-staff, got %q", got)
	}
}
