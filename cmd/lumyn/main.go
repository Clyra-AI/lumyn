package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Clyra-AI/lumyn/internal/exitcode"
	"github.com/Clyra-AI/lumyn/internal/result"
	"github.com/Clyra-AI/lumyn/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, time.Now()))
}

func run(args []string, stdout io.Writer, stderr io.Writer, started time.Time) int {
	payload, exitCode := commandResultForArgs(args, started)

	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(payload); err != nil {
		fmt.Fprintf(stderr, "encode command result: %v\n", err)
		return exitcode.InternalError
	}
	return exitCode
}

func commandResultForArgs(args []string, started time.Time) (result.CommandResult, int) {
	command := commandFromArgs(args)
	status := "pass"
	findingKind := "none"
	fixTarget := "not_applicable"
	exitCode := exitcode.Success
	errors := []result.CommandError{}
	if !isKnownCommand(command) {
		status = "fail"
		findingKind = "command_error"
		fixTarget = "command"
		exitCode = exitcode.InvalidUsageOrInput
		errors = append(errors, result.CommandError{
			Code:    "unknown_command",
			Message: fmt.Sprintf("unknown command %q", command),
		})
	}

	payload := result.CommandResult{
		ObjectType:           "lumyn.command_result",
		SchemaVersion:        "1.0",
		Metadata:             commandMetadata(),
		Command:              command,
		Status:               status,
		Mode:                 command,
		Warnings:             []string{},
		Errors:               errors,
		Artifacts:            []result.ArtifactRef{},
		DurationMS:           time.Since(started).Milliseconds(),
		RedactionStatus:      "not_applicable",
		FindingKind:          findingKind,
		ProofStrength:        "unknown",
		ActionBoundaryStatus: "not_configured",
		SecurityRelevance:    "none",
		FixTarget:            fixTarget,
		SurfaceFingerprint:   "not_applicable",
		EvalMode:             "not_applicable",
		ProviderMetadata: result.ProviderMetadata{
			Applicable: false,
			Provider:   "not_applicable",
			Model:      "not_applicable",
		},
		CorpusEligible: false,
	}

	return payload, exitCode
}

func commandFromArgs(args []string) string {
	if len(args) == 0 {
		return "help"
	}
	return args[0]
}

func commandMetadata() map[string]any {
	return map[string]any{
		"lumyn_version": version.Version,
		"runtime":       "go",
		"source":        "cli",
	}
}

func isKnownCommand(command string) bool {
	switch command {
	case "help", "version", "init", "check", "record", "verify", "trace", "demo", "share", "eval":
		return true
	default:
		return false
	}
}
