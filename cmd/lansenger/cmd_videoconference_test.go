package main

import (
	"testing"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"
)

func TestParseVCMembers(t *testing.T) {
	if _, err := parseVCMembers("", "--member"); err == nil {
		t.Fatal("expected error for empty member flag")
	}
	if _, err := parseVCMembers("not-json", "--member"); err == nil {
		t.Fatal("expected error for non-JSON member flag")
	}
	if _, err := parseVCMembers(`[{"staffId":"1","role":"participant"}]`, "--member"); err == nil {
		t.Fatal("expected error for zero admins")
	}
	if _, err := parseVCMembers(`[{"staffId":"1","role":"admin"},{"staffId":"2","role":"admin"}]`, "--member"); err == nil {
		t.Fatal("expected error for two admins")
	}
	members, err := parseVCMembers(`[{"staffId":"1","role":"admin"},{"staffId":"2","role":"participant"}]`, "--member")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 2 || members[0].Role != lansenger.VCMemberRoleHost {
		t.Fatalf("unexpected members: %+v", members)
	}
}

func TestParseVCVods(t *testing.T) {
	if _, err := parseVCVods(""); err == nil {
		t.Fatal("expected error for empty vods flag")
	}
	if _, err := parseVCVods(`[]`); err == nil {
		t.Fatal("expected error for zero vods")
	}
	if _, err := parseVCVods(`["a","b","c","d"]`); err == nil {
		t.Fatal("expected error for more than 3 vods")
	}
	vods, err := parseVCVods(`["v1",{"vodId":"v2"}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vods) != 2 || vods[0].VodID != "v1" || vods[1].VodID != "v2" {
		t.Fatalf("unexpected vods: %+v", vods)
	}
}

func TestVCOpCodeValid(t *testing.T) {
	if !vcOpCodeValid("muteall") || !vcOpCodeValid("setHost") {
		t.Fatal("expected whitelisted opCodes to be valid")
	}
	if vcOpCodeValid("notAnOp") {
		t.Fatal("expected unknown opCode to be rejected")
	}
	if len(lansenger.VCOpCodes) != 22 {
		t.Fatalf("expected 22 opCodes, got %d", len(lansenger.VCOpCodes))
	}
}
