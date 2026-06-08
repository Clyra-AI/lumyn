package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Clyra-AI/lumyn/internal/exitcode"
	"github.com/Clyra-AI/lumyn/internal/result"
)

func main() {
	started := time.Now()
	command := "help"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	status := "pass"
	exitCode := exitcode.Success
	errors := []result.CommandError{}
	if !isKnownCommand(command) {
		status = "fail"
		exitCode = exitcode.InvalidUsageOrInput
		errors = append(errors, result.CommandError{
			Code:    "unknown_command",
			Message: fmt.Sprintf("unknown command %q", command),
		})
	}

	payload := result.CommandResult{
		ObjectType:      "lumyn.command_result",
		SchemaVersion:   "1.0",
		Command:         command,
		Status:          status,
		Mode:            command,
		Warnings:        []string{},
		Errors:          errors,
		Artifacts:       []result.ArtifactRef{},
		DurationMS:      time.Since(started).Milliseconds(),
		RedactionStatus: "not_applicable",
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(payload); err != nil {
		fmt.Fprintf(os.Stderr, "encode command result: %v\n", err)
		os.Exit(exitcode.InternalError)
	}
	os.Exit(exitCode)
}

func isKnownCommand(command string) bool {
	switch command {
	case "help", "version", "init", "check", "record", "verify", "demo", "share", "eval":
		return true
	default:
		return false
	}
}
