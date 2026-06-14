package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Clyra-AI/lumyn/internal/config"
	"github.com/Clyra-AI/lumyn/internal/exitcode"
	"github.com/Clyra-AI/lumyn/internal/result"
	"github.com/Clyra-AI/lumyn/internal/source"
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
	args = trimGlobalFlags(args)
	command := commandFromArgs(args)
	if command == "init" {
		return initCommandResult(args[1:], started)
	}
	if command == "check" {
		return checkCommandResult(args[1:], started)
	}

	payload := baseCommandResult(command, started)
	exitCode := exitcode.Success
	if !isKnownCommand(command) {
		payload.Status = "fail"
		payload.FindingKind = "command_error"
		payload.FixTarget = "command"
		exitCode = exitcode.InvalidUsageOrInput
		payload.Errors = append(payload.Errors, result.CommandError{
			Code:    "unknown_command",
			Message: fmt.Sprintf("unknown command %q", command),
		})
	}
	payload.DurationMS = time.Since(started).Milliseconds()
	return payload, exitCode
}

func initCommandResult(args []string, started time.Time) (result.CommandResult, int) {
	payload := baseCommandResult("init", started)
	options, err := parseInitFlags(args)
	if err != nil {
		return commandInputError(payload, started, "invalid_init_args", err)
	}

	report, err := source.InitProject(options)
	if err != nil {
		return commandInputError(payload, started, "init_failed", err)
	}
	payload.RedactionStatus = "applied"
	payload.EvalMode = "surface_only"
	payload.SurfaceFingerprint = report.SurfaceFingerprint
	payload.Artifacts = append(payload.Artifacts,
		result.ArtifactRef{Path: options.ConfigPath, Type: "config"},
		result.ArtifactRef{Path: report.ReportPath, Type: "source_intake"},
	)
	attachSourceReportMetadata(&payload, report, false)
	payload.DurationMS = time.Since(started).Milliseconds()
	return payload, exitcode.Success
}

func checkCommandResult(args []string, started time.Time) (result.CommandResult, int) {
	payload := baseCommandResult("check", started)
	configPath, strict, err := parseCheckFlags(args)
	if err != nil {
		return commandInputError(payload, started, "invalid_check_args", err)
	}

	report, err := source.CheckProject(configPath)
	if err != nil {
		return commandInputError(payload, started, "config_error", err)
	}

	payload.RedactionStatus = "applied"
	payload.EvalMode = "surface_only"
	payload.SurfaceFingerprint = report.SurfaceFingerprint
	payload.Artifacts = append(payload.Artifacts, result.ArtifactRef{Path: report.ReportPath, Type: "source_check"})
	attachSourceReportMetadata(&payload, report, strict)

	exitCode := exitcode.Success
	if first, ok := source.FirstFinding(report.Findings); ok {
		payload.Status = "warning"
		payload.FindingKind = first.Kind
		payload.FixTarget = first.FixTarget
		if source.ContainsProofGap(report.Findings) {
			payload.ProofStrength = "gap"
		}
		if source.ContainsSecurityRelevantFinding(report.Findings) {
			payload.SecurityRelevance = "safety_relevant"
		}
		if source.HasErrorFindings(report.Findings) {
			payload.Status = "fail"
			exitCode = exitcode.InvalidUsageOrInput
		} else if strict {
			payload.Status = "fail"
			exitCode = exitcode.SourceCompletenessFailure
		}
	}
	payload.DurationMS = time.Since(started).Milliseconds()
	return payload, exitCode
}

func parseInitFlags(args []string) (source.InitOptions, error) {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.Bool("json", false, "emit JSON command result")
	configPath := flags.String("config", config.DefaultConfigPath, "path to lumyn config")
	openAPIPath := flags.String("openapi", "./openapi.yaml", "path to OpenAPI source")
	docsPath := flags.String("docs", "./docs", "path to local docs source")
	if err := flags.Parse(args); err != nil {
		return source.InitOptions{}, err
	}
	if flags.NArg() != 0 {
		return source.InitOptions{}, fmt.Errorf("unexpected init argument %q", flags.Arg(0))
	}
	return source.InitOptions{
		ConfigPath:  *configPath,
		OpenAPIPath: *openAPIPath,
		DocsPath:    *docsPath,
	}, nil
}

func parseCheckFlags(args []string) (string, bool, error) {
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.Bool("json", false, "emit JSON command result")
	configPath := flags.String("config", config.DefaultConfigPath, "path to lumyn config")
	strict := flags.Bool("strict", false, "fail when source findings are present")
	if err := flags.Parse(args); err != nil {
		return "", false, err
	}
	if flags.NArg() != 0 {
		return "", false, fmt.Errorf("unexpected check argument %q", flags.Arg(0))
	}
	return *configPath, *strict, nil
}

func commandInputError(payload result.CommandResult, started time.Time, code string, err error) (result.CommandResult, int) {
	payload.Status = "fail"
	payload.FindingKind = "command_error"
	payload.FixTarget = "command"
	payload.Errors = append(payload.Errors, result.CommandError{
		Code:    code,
		Message: err.Error(),
	})
	payload.DurationMS = time.Since(started).Milliseconds()
	return payload, exitcode.InvalidUsageOrInput
}

func attachSourceReportMetadata(payload *result.CommandResult, report source.Report, strict bool) {
	payload.Metadata["source_report_path"] = report.ReportPath
	payload.Metadata["source_refs"] = report.SourceRefs
	payload.Metadata["source_findings"] = report.Findings
	payload.Metadata["source_finding_count"] = len(report.Findings)
	payload.Metadata["strict"] = strict
}

func baseCommandResult(command string, started time.Time) result.CommandResult {
	status := "pass"
	findingKind := "none"
	fixTarget := "not_applicable"
	errors := []result.CommandError{}

	return result.CommandResult{
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
}

func commandFromArgs(args []string) string {
	if len(args) == 0 {
		return "help"
	}
	return args[0]
}

func trimGlobalFlags(args []string) []string {
	for len(args) > 0 && args[0] == "--json" {
		args = args[1:]
	}
	return args
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
