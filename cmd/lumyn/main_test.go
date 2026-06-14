package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/Clyra-AI/lumyn/internal/exitcode"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func TestCLIEmitsCommandResultEnvelope(t *testing.T) {
	binary := buildTestBinary(t)
	command := exec.Command(binary, "help")
	output, err := command.Output()
	if err != nil {
		t.Fatalf("run lumyn command: %v", err)
	}

	payload := decodeCommandResult(t, output)
	validateCommandResultSchema(t, payload)

	expected := map[string]string{
		"object_type":            "lumyn.command_result",
		"schema_version":         "1.0",
		"command":                "help",
		"status":                 "pass",
		"mode":                   "help",
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

func TestCLIInitWritesConfigAndSourceArtifact(t *testing.T) {
	binary := buildTestBinary(t)
	projectDir := t.TempDir()
	writeOpenAPIFixture(t, filepath.Join(projectDir, "openapi.json"))
	if err := os.Mkdir(filepath.Join(projectDir, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "docs", "auth.md"), []byte("# Auth\nUse the API key header.\n"), 0o644); err != nil {
		t.Fatalf("write docs: %v", err)
	}

	command := exec.Command(binary, "init", "--openapi", "./openapi.json", "--docs", "./docs")
	command.Dir = projectDir
	output, err := command.Output()
	if err != nil {
		t.Fatalf("run lumyn init: %v", err)
	}

	payload := decodeCommandResult(t, output)
	validateCommandResultSchema(t, payload)
	if payload["command"] != "init" {
		t.Fatalf("command = %v, want init", payload["command"])
	}
	if payload["status"] != "pass" {
		t.Fatalf("status = %v, want pass", payload["status"])
	}
	configBytes, err := os.ReadFile(filepath.Join(projectDir, "lumyn.yaml"))
	if err != nil {
		t.Fatalf("read generated config: %v", err)
	}
	if !bytes.Contains(configBytes, []byte("path: ./openapi.json")) {
		t.Fatalf("generated config missing openapi path:\n%s", configBytes)
	}

	artifactPath := firstArtifactPath(t, payload, "source_intake")
	artifactBytes, err := os.ReadFile(filepath.Join(projectDir, artifactPath))
	if err != nil {
		t.Fatalf("read source intake artifact %s: %v", artifactPath, err)
	}
	var artifact map[string]any
	if err := json.Unmarshal(artifactBytes, &artifact); err != nil {
		t.Fatalf("source intake artifact is not valid JSON: %v", err)
	}
	if artifact["object_type"] != "lumyn.source_check" {
		t.Fatalf("artifact object_type = %v, want lumyn.source_check", artifact["object_type"])
	}
}

func TestCLICheckEmitsWorkflowRelevantSourceFindingReference(t *testing.T) {
	binary := buildTestBinary(t)
	projectDir := t.TempDir()
	writeOpenAPIFixture(t, filepath.Join(projectDir, "openapi.json"))
	if err := os.Mkdir(filepath.Join(projectDir, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "docs", "guide.md"), []byte("[missing](missing.md)\n"), 0o644); err != nil {
		t.Fatalf("write docs: %v", err)
	}

	initCommand := exec.Command(binary, "init", "--openapi", "./openapi.json", "--docs", "./docs")
	initCommand.Dir = projectDir
	if output, err := initCommand.CombinedOutput(); err != nil {
		t.Fatalf("run lumyn init: %v\n%s", err, output)
	}

	checkCommand := exec.Command(binary, "check")
	checkCommand.Dir = projectDir
	output, err := checkCommand.Output()
	if err != nil {
		t.Fatalf("run lumyn check: %v", err)
	}

	payload := decodeCommandResult(t, output)
	validateCommandResultSchema(t, payload)
	if payload["command"] != "check" {
		t.Fatalf("command = %v, want check", payload["command"])
	}
	if payload["status"] != "warning" {
		t.Fatalf("status = %v, want warning", payload["status"])
	}
	if payload["finding_kind"] != "docs_api_ambiguity" {
		t.Fatalf("finding_kind = %v, want docs_api_ambiguity", payload["finding_kind"])
	}
	metadata, ok := payload["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata = %T, want object", payload["metadata"])
	}
	findings, ok := metadata["source_findings"].([]any)
	if !ok || len(findings) == 0 {
		t.Fatalf("metadata.source_findings = %#v, want non-empty array", metadata["source_findings"])
	}
	first, ok := findings[0].(map[string]any)
	if !ok {
		t.Fatalf("first finding = %T, want object", findings[0])
	}
	reference, ok := first["reference"].(map[string]any)
	if !ok {
		t.Fatalf("first finding reference = %T, want object", first["reference"])
	}
	if reference["path"] == "" {
		t.Fatalf("first finding reference lacks path: %#v", reference)
	}
	if reference["json_pointer"] == "" && reference["line"] == nil && reference["object"] == "" {
		t.Fatalf("first finding lacks concrete source reference: %#v", reference)
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

func TestCommandResultForArgsRunsInitAndStrictCheck(t *testing.T) {
	projectDir := t.TempDir()
	t.Chdir(projectDir)
	writeOpenAPIFixture(t, filepath.Join(projectDir, "openapi.json"))
	if err := os.Mkdir(filepath.Join(projectDir, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "docs", "guide.md"), []byte("See [missing](missing.md)\n"), 0o644); err != nil {
		t.Fatalf("write docs: %v", err)
	}

	initPayload, initCode := commandResultForArgs([]string{"init", "--openapi", "./openapi.json", "--docs", "./docs"}, time.Now())
	if initCode != exitcode.Success {
		t.Fatalf("init exit code = %d, want %d; errors=%#v", initCode, exitcode.Success, initPayload.Errors)
	}
	if initPayload.Status != "pass" {
		t.Fatalf("init status = %q, want pass", initPayload.Status)
	}
	if initPayload.Metadata["source_finding_count"] != 0 {
		t.Fatalf("init source_finding_count = %v, want 0", initPayload.Metadata["source_finding_count"])
	}

	checkPayload, checkCode := commandResultForArgs([]string{"check", "--strict"}, time.Now())
	if checkCode != exitcode.SourceCompletenessFailure {
		t.Fatalf("strict check exit code = %d, want %d; errors=%#v", checkCode, exitcode.SourceCompletenessFailure, checkPayload.Errors)
	}
	if checkPayload.Status != "fail" {
		t.Fatalf("strict check status = %q, want fail", checkPayload.Status)
	}
	if checkPayload.FindingKind == "none" {
		t.Fatalf("strict check finding kind should identify the source finding")
	}
	if checkPayload.ProofStrength != "gap" {
		t.Fatalf("strict check proof_strength = %q, want gap", checkPayload.ProofStrength)
	}
}

func TestCommandResultForArgsHandlesJSONFlagAndBadCommandFlags(t *testing.T) {
	payload, code := commandResultForArgs([]string{"--json", "check", "--bad-flag"}, time.Now())
	if code != exitcode.InvalidUsageOrInput {
		t.Fatalf("exit code = %d, want %d", code, exitcode.InvalidUsageOrInput)
	}
	if payload.Command != "check" {
		t.Fatalf("command = %q, want check", payload.Command)
	}
	if payload.Status != "fail" {
		t.Fatalf("status = %q, want fail", payload.Status)
	}
	if len(payload.Errors) != 1 || payload.Errors[0].Code != "invalid_check_args" {
		t.Fatalf("errors = %#v, want invalid_check_args", payload.Errors)
	}

	payload, code = commandResultForArgs([]string{"init", "--bad-flag"}, time.Now())
	if code != exitcode.InvalidUsageOrInput {
		t.Fatalf("init exit code = %d, want %d", code, exitcode.InvalidUsageOrInput)
	}
	if len(payload.Errors) != 1 || payload.Errors[0].Code != "invalid_init_args" {
		t.Fatalf("init errors = %#v, want invalid_init_args", payload.Errors)
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
	code := run([]string{"help"}, &stdout, &stderr, time.Now())
	if code != exitcode.Success {
		t.Fatalf("exit code = %d, want %d", code, exitcode.Success)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	payload := decodeCommandResult(t, stdout.Bytes())
	if payload["command"] != "help" {
		t.Fatalf("command = %v, want help", payload["command"])
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

func writeOpenAPIFixture(t *testing.T, path string) {
	t.Helper()
	spec := []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "post": {
        "operationId": "createCustomer",
        "responses": {
          "201": {"description": "created"}
        }
      }
    }
  }
}`)
	if err := os.WriteFile(path, spec, 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
}

func firstArtifactPath(t *testing.T, payload map[string]any, artifactType string) string {
	t.Helper()
	artifacts, ok := payload["artifacts"].([]any)
	if !ok {
		t.Fatalf("artifacts = %T, want array", payload["artifacts"])
	}
	for _, artifact := range artifacts {
		item, ok := artifact.(map[string]any)
		if !ok {
			t.Fatalf("artifact item = %T, want object", artifact)
		}
		if item["type"] == artifactType {
			path, ok := item["path"].(string)
			if !ok || path == "" {
				t.Fatalf("artifact path = %#v, want non-empty string", item["path"])
			}
			return path
		}
	}
	t.Fatalf("artifact type %q not found in %#v", artifactType, artifacts)
	return ""
}
