package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/Clyra-AI/lumyn/internal/exitcode"
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
		"object_type":            "lumyn.command_result",
		"schema_version":         "1.0",
		"command":                "check",
		"status":                 "pass",
		"mode":                   "check",
		"redaction_status":       "not_applicable",
		"finding_kind":           "none",
		"proof_strength":         "unknown",
		"action_boundary_status": "not_configured",
		"security_relevance":     "none",
		"fix_target":             "not_applicable",
		"surface_fingerprint":    "not_applicable",
		"eval_mode":              "not_applicable",
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
	if payload["corpus_eligible"] != false {
		t.Fatalf("corpus_eligible = %v, want false", payload["corpus_eligible"])
	}
	providerMetadata, ok := payload["provider_metadata"].(map[string]any)
	if !ok {
		t.Fatalf("provider_metadata = %T, want object", payload["provider_metadata"])
	}
	if providerMetadata["applicable"] != false {
		t.Fatalf("provider_metadata.applicable = %v, want false", providerMetadata["applicable"])
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
	if payload["finding_kind"] != "command_error" {
		t.Fatalf("finding_kind = %v, want command_error", payload["finding_kind"])
	}
	if payload["corpus_eligible"] != false {
		t.Fatalf("corpus_eligible = %v, want false", payload["corpus_eligible"])
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

func TestCommandResultForArgsDefaultsToHelp(t *testing.T) {
	payload, code := commandResultForArgs(nil, time.Now())
	if code != exitcode.Success {
		t.Fatalf("exit code = %d, want %d", code, exitcode.Success)
	}
	if payload.Command != "help" {
		t.Fatalf("command = %q, want help", payload.Command)
	}
	if payload.Status != "pass" {
		t.Fatalf("status = %q, want pass", payload.Status)
	}
	if payload.Metadata["runtime"] != "go" {
		t.Fatalf("metadata.runtime = %v, want go", payload.Metadata["runtime"])
	}
}

func TestCommandResultForArgsRejectsUnknownCommand(t *testing.T) {
	payload, code := commandResultForArgs([]string{"nope"}, time.Now())
	if code != exitcode.InvalidUsageOrInput {
		t.Fatalf("exit code = %d, want %d", code, exitcode.InvalidUsageOrInput)
	}
	if payload.Status != "fail" {
		t.Fatalf("status = %q, want fail", payload.Status)
	}
	if len(payload.Errors) != 1 || payload.Errors[0].Code != "unknown_command" {
		t.Fatalf("errors = %#v, want unknown_command", payload.Errors)
	}
}

func TestCommandFromArgsUsesFirstArg(t *testing.T) {
	if got := commandFromArgs([]string{"verify", "ignored"}); got != "verify" {
		t.Fatalf("commandFromArgs = %q, want verify", got)
	}
}

func TestRunWritesEnvelopeAndReturnsExitCode(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run([]string{"check"}, &stdout, &stderr, time.Now())
	if code != exitcode.Success {
		t.Fatalf("exit code = %d, want %d", code, exitcode.Success)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	payload := decodeCommandResult(t, stdout.Bytes())
	if payload["command"] != "check" {
		t.Fatalf("command = %v, want check", payload["command"])
	}
}

func TestRunReturnsInternalErrorWhenEncodeFails(t *testing.T) {
	var stderr bytes.Buffer
	code := run([]string{"check"}, errWriter{}, &stderr, time.Now())
	if code != exitcode.InternalError {
		t.Fatalf("exit code = %d, want %d", code, exitcode.InternalError)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("encode command result")) {
		t.Fatalf("stderr = %q, want encode error", stderr.String())
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
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
