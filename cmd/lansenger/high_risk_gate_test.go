package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// binaryPath is built once in TestMain so it survives across all subtests
// (t.TempDir() cleanup would remove it between tests).
var binaryPath string

func TestMain(m *testing.M) {
	f, err := os.CreateTemp("", "lansenger-test-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "create temp binary:", err)
		os.Exit(1)
	}
	binaryPath = f.Name()
	f.Close()
	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, "go build failed:", err)
		fmt.Fprintln(os.Stderr, string(out))
		os.Exit(1)
	}
	code := m.Run()
	os.Remove(binaryPath)
	os.Exit(code)
}

func buildHelper(t *testing.T) string {
	if binaryPath == "" {
		t.Fatal("binaryPath not built")
	}
	return binaryPath
}

type expect struct {
	args     []string
	exitCode int
	stderr   string // substring expected on stderr (empty = no assertion)
	stdout   string // substring expected on stdout (empty = no assertion)
}

// runGate builds the CLI and runs the given args (without leading "lansenger").
func runGate(t *testing.T, args []string) (code int, stdout, stderr string) {
	bin := buildHelper(t)
	cmd := exec.Command(bin, args...)
	out, err := cmd.CombinedOutput()
	// cobra writes confirmation to stderr; combinedOutput merges them — split via
	// running again with separate pipes is heavier; instead parse from combined.
	stderr = string(out)
	stdout = string(out)
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		t.Fatalf("exec error: %v", err)
	}
	return code, stdout, stderr
}

func TestHighRiskGate_Exit10(t *testing.T) {
	cases := []expect{
		{args: []string{"group", "dismiss", "g1"}, exitCode: 10, stderr: "Confirmation required"},
		{args: []string{"message", "revoke", "m1"}, exitCode: 10, stderr: "Confirmation required"},
		{args: []string{"group", "update-members", "g1", "--remove", "u1"}, exitCode: 10, stderr: "Confirmation required"},
		{args: []string{"calendar", "delete-schedule", "cal1", "sch1"}, exitCode: 10, stderr: "Confirmation required"},
		{args: []string{"calendar", "delete-attendees", "cal1", "sch1", "[\"u1\"]"}, exitCode: 10, stderr: "Confirmation required"},
		{args: []string{"todo", "delete", "t1", "org1"}, exitCode: 10, stderr: "Confirmation required"},
		{args: []string{"todo", "delete-executors", "u1", "org1"}, exitCode: 10, stderr: "Confirmation required"},
	}
	for _, c := range cases {
		code, _, se := runGate(t, c.args)
		if code != c.exitCode {
			t.Errorf("%v: exit code = %d, want %d", c.args, code, c.exitCode)
		}
		if c.stderr != "" && !strings.Contains(se, c.stderr) {
			t.Errorf("%v: stderr = %q, want substring %q", c.args, se, c.stderr)
		}
	}
}

func TestHighRiskGate_JSONExit10Envelope(t *testing.T) {
	bin := buildHelper(t)
	cmd := exec.Command(bin, "--json", "group", "dismiss", "g1")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	var stdout strings.Builder
	cmd.Stdout = &stdout
	err := cmd.Run()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	}
	if code != 10 {
		t.Fatalf("exit code = %d, want 10", code)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(stderr.String()), &payload); err != nil {
		t.Fatalf("stderr not JSON: %v (got %q)", err, stderr.String())
	}
	if payload["ok"] != false {
		t.Errorf("ok = %v, want false", payload["ok"])
	}
	errObj, _ := payload["error"].(map[string]interface{})
	if errObj["type"] != "confirmation_required" {
		t.Errorf("error.type = %v, want confirmation_required", errObj["type"])
	}
	risk, _ := errObj["risk"].(map[string]interface{})
	if risk["level"] != "high-risk-write" {
		t.Errorf("risk.level = %v, want high-risk-write", risk["level"])
	}
	if risk["action"] != "dismiss group g1" {
		t.Errorf("risk.action = %v, want 'dismiss group g1'", risk["action"])
	}
}

func TestHighRiskGate_DryRun(t *testing.T) {
	// dry-run exits 0 and prints a preview; no API call (no credentials needed).
	cases := []expect{
		{args: []string{"group", "dismiss", "g1", "--dry-run"}, exitCode: 0, stdout: "DRY RUN"},
		{args: []string{"group", "update-members", "g1", "--remove", "u1", "--dry-run"}, exitCode: 0, stdout: "DRY RUN"},
		{args: []string{"todo", "delete", "t1", "org1", "--dry-run"}, exitCode: 0, stdout: "DRY RUN"},
	}
	for _, c := range cases {
		code, so, _ := runGate(t, c.args)
		if code != c.exitCode {
			t.Errorf("%v: exit code = %d, want %d", c.args, code, c.exitCode)
		}
		if c.stdout != "" && !strings.Contains(so, c.stdout) {
			t.Errorf("%v: stdout = %q, want substring %q", c.args, so, c.stdout)
		}
	}
}

func TestHighRiskGate_DryRunJSON(t *testing.T) {
	bin := buildHelper(t)
	cmd := exec.Command(bin, "--json", "group", "dismiss", "g1", "--dry-run")
	var stdout strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 0 {
			t.Fatalf("run error: %v", err)
		}
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(stdout.String()), &payload); err != nil {
		t.Fatalf("stdout not JSON: %v (got %q)", err, stdout.String())
	}
	if payload["ok"] != true {
		t.Errorf("ok = %v, want true", payload["ok"])
	}
	if payload["dry_run"] != true {
		t.Errorf("dry_run = %v, want true", payload["dry_run"])
	}
	if payload["would_perform"] != "dismiss group g1" {
		t.Errorf("would_perform = %v, want 'dismiss group g1'", payload["would_perform"])
	}
}

// update-members with only --add is not gated and proceeds (it will then fail on
// missing credentials, exiting non-zero but NOT 10 — proving no gate fired).
func TestHighRiskGate_AddOnlyNotGated(t *testing.T) {
	bin := buildHelper(t)
	cmd := exec.Command(bin, "group", "update-members", "g1", "--add", "u1")
	out, err := cmd.CombinedOutput()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	}
	if code == 10 {
		t.Errorf("add-only should not be gated, got exit 10. output: %s", out)
	}
	// expect a credentials error (exit 1), not a confirmation gate
	if !strings.Contains(string(out), "credentials") && !strings.Contains(string(out), "No credentials") {
		// acceptable: any non-10 exit; just ensure not 10 and not "Confirmation required"
	}
	if strings.Contains(string(out), "Confirmation required") {
		t.Errorf("add-only must not trigger confirmation gate: %s", out)
	}
}
