package main

import (
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLIEmitsCommandResultEnvelope(t *testing.T) {
	binary := buildTestBinary(t)
	command := exec.Command(binary, "check")
	output, err := command.Output()
	if err != nil {
		t.Fatalf("run lumyn command: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("decode command result: %v\n%s", err, output)
	}

	expected := map[string]string{
		"object_type":      "lumyn.command_result",
		"schema_version":   "1.0",
		"command":          "check",
		"status":           "pass",
		"mode":             "check",
		"redaction_status": "not_applicable",
	}
	for key, value := range expected {
		if payload[key] != value {
			t.Fatalf("%s = %v, want %q", key, payload[key], value)
		}
	}
}

func TestCLIRejectsUnknownCommandWithJSONEnvelope(t *testing.T) {
	binary := buildTestBinary(t)
	command := exec.Command(binary, "unknown-command")
	output, err := command.Output()
	if err == nil {
		t.Fatal("unknown command should fail")
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("unexpected error type: %v", err)
	}
	if exitErr.ExitCode() != 2 {
		t.Fatalf("exit code = %d, want 2", exitErr.ExitCode())
	}

	var payload map[string]any
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("decode command result: %v\n%s", err, output)
	}
	if payload["status"] != "fail" {
		t.Fatalf("status = %v, want fail", payload["status"])
	}
	if payload["command"] != "unknown-command" {
		t.Fatalf("command = %v, want unknown-command", payload["command"])
	}
}

func buildTestBinary(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "lumyn")
	command := exec.Command("go", "build", "-o", binary, ".")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build test binary: %v\n%s", err, output)
	}
	return binary
}
