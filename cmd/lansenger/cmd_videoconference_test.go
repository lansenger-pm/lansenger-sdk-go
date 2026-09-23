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

// VCOpCodes 只是「已知值」参考表：客户端不再校验 OP_CODE，服务端才是权威。
func TestVCOpCodesIsReferenceOnly(t *testing.T) {
	if vcMemberControlCmd.PreRunE != nil {
		t.Fatal("member-control must not validate OP_CODE on the client")
	}

	known := make(map[string]bool, len(lansenger.VCOpCodes))
	for _, op := range lansenger.VCOpCodes {
		known[op] = true
	}
	// 实测修正 (2026-09-23)：服务端认 "mute"，不认 "applyAudio"
	if !known["mute"] {
		t.Fatal(`expected "mute" to be listed as a known opCode`)
	}
	if known["applyAudio"] {
		t.Fatal(`"applyAudio" is rejected by the server (errCode 105601) and must not be listed`)
	}
}

// modify 与 create 同契约：都必须能传自动结束时间。
func TestVCModifyAcceptsUserStopTime(t *testing.T) {
	if vcModifyCmd.Flags().Lookup("user-stop-time") == nil {
		t.Fatal("modify must accept --user-stop-time")
	}
	if int64PtrOrNil(-1) != nil {
		t.Fatal("-1 must mean omit")
	}
	if v := int64PtrOrNil(1700003600000); v == nil || *v != 1700003600000 {
		t.Fatal("expected the value to be forwarded")
	}
}
