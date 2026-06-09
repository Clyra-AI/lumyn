package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func TestCLIEmitsCommandResultEnvelope(t *testing.T) {
	binary := buildTestBinary(t)
	command := exec.Command(binary, "check")
	output, err := command.Output()
	if err != nil {
		t.Fatalf("run lumyn command: %v", err)
	}

	payload := decodeCommandResult(t, output)
	validateCommandResultSchema(t, payload)

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
	metadata, ok := payload["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata = %T, want object", payload["metadata"])
	}
	if metadata["lumyn_version"] != "0.0.0" {
		t.Fatalf("metadata.lumyn_version = %v, want 0.0.0", metadata["lumyn_version"])
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

	payload := decodeCommandResult(t, output)
	validateCommandResultSchema(t, payload)
	if payload["status"] != "fail" {
		t.Fatalf("status = %v, want fail", payload["status"])
	}
	if payload["command"] != "unknown-command" {
		t.Fatalf("command = %v, want unknown-command", payload["command"])
	}
}

func TestCLIAcceptsRequiredTraceCommand(t *testing.T) {
	binary := buildTestBinary(t)
	command := exec.Command(binary, "trace")
	output, err := command.Output()
	if err != nil {
		t.Fatalf("run lumyn trace: %v", err)
	}

	payload := decodeCommandResult(t, output)
	validateCommandResultSchema(t, payload)
	if payload["command"] != "trace" {
		t.Fatalf("command = %v, want trace", payload["command"])
	}
	if payload["status"] != "pass" {
		t.Fatalf("status = %v, want pass", payload["status"])
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

func decodeCommandResult(t *testing.T, output []byte) map[string]any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(output))
	decoder.UseNumber()
	var payload map[string]any
	if err := decoder.Decode(&payload); err != nil {
		t.Fatalf("decode command result: %v\n%s", err, output)
	}
	return payload
}

func validateCommandResultSchema(t *testing.T, payload map[string]any) {
	t.Helper()
	schema, err := jsonschema.Compile(filepath.Join("..", "..", "schemas", "command-result.schema.json"))
	if err != nil {
		t.Fatalf("compile command-result schema: %v", err)
	}
	if err := schema.Validate(payload); err != nil {
		t.Fatalf("command-result schema validation failed: %v", err)
	}
}
