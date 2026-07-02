package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Clyra-AI/lumyn/internal/config"
)

const (
	schemaVersion = "1.0"

	intakeReportPath = "runs/source-checks/source-intake.json"
	checkReportPath  = "runs/source-checks/source-check.json"
)

var (
	httpMethods    = map[string]bool{"get": true, "put": true, "post": true, "delete": true, "patch": true, "head": true, "options": true, "trace": true}
	mutatingMethod = map[string]bool{"put": true, "post": true, "delete": true, "patch": true}
)

type InitOptions struct {
	ConfigPath  string
	OpenAPIPath string
	DocsPath    string
}

type ProjectConfig struct {
	Version int          `json:"version"`
	Sources SourceConfig `json:"sources"`
	Env     EnvConfig    `json:"env"`
	Redact  RedactConfig `json:"redaction"`
	Verify  VerifyConfig `json:"verify"`
}

type SourceConfig struct {
	OpenAPI []SourceEntry `json:"openapi"`
	Docs    []SourceEntry `json:"docs"`
}

type SourceEntry struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type EnvConfig struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}

type RedactConfig struct {
	Headers []string `json:"headers"`
	Fields  []string `json:"fields"`
}

type VerifyConfig struct {
	DefaultStrategy string       `json:"default_strategy"`
	Replay          ReplayConfig `json:"replay"`
}

type ReplayConfig struct {
	AllowNetwork bool `json:"allow_network"`
}

type Report struct {
	ObjectType         string      `json:"object_type"`
	SchemaVersion      string      `json:"schema_version"`
	Status             string      `json:"status"`
	ConfigPath         string      `json:"config_path"`
	ReportPath         string      `json:"report_path"`
	SourceRefs         []SourceRef `json:"source_refs"`
	Findings           []Finding   `json:"findings"`
	SurfaceFingerprint string      `json:"surface_fingerprint"`
}

type SourceRef struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Path string `json:"path"`
	Hash string `json:"hash,omitempty"`
}

type Finding struct {
	Kind              string    `json:"kind"`
	Severity          string    `json:"severity"`
	Message           string    `json:"message"`
	Reference         Reference `json:"reference"`
	FixTarget         string    `json:"fix_target"`
	WorkflowRelevance string    `json:"workflow_relevance"`
}

type Reference struct {
	Path        string `json:"path"`
	JSONPointer string `json:"json_pointer,omitempty"`
	Line        int    `json:"line,omitempty"`
	Object      string `json:"object,omitempty"`
}

type openAPIOperation struct {
	Path              string
	Method            string
	OperationID       string
	Summary           string
	Description       string
	Deprecated        bool
	Pointer           string
	Line              int
	HasRequestBody    bool
	HasRequestSchema  bool
	HasResponseSchema bool
	ReplacementHint   bool
	Parameters        []openAPIParameter
}

type openAPIParameter struct {
	Name        string
	In          string
	Description string
	Pointer     string
	Line        int
}

func InitProject(options InitOptions) (Report, error) {
	options = normalizedInitOptions(options)
	root := projectRoot(options.ConfigPath)
	configPathWasRelative := !filepath.IsAbs(options.ConfigPath)
	projectConfig := ProjectConfig{
		Version: 1,
		Sources: SourceConfig{
			OpenAPI: []SourceEntry{{ID: "public_api", Path: sourcePathForConfigRoot(root, options.OpenAPIPath, configPathWasRelative)}},
			Docs:    []SourceEntry{{ID: "docs", Path: sourcePathForConfigRoot(root, options.DocsPath, configPathWasRelative)}},
		},
		Env: EnvConfig{
			BaseURL: "${API_BASE_URL}",
			APIKey:  "${LUMYN_API_KEY}",
		},
		Redact: RedactConfig{
			Headers: []string{"authorization", "x-api-key"},
			Fields:  []string{"access_token", "refresh_token", "secret", "password"},
		},
		Verify: VerifyConfig{
			DefaultStrategy: "replay",
			Replay:          ReplayConfig{AllowNetwork: false},
		},
	}

	if err := writeProjectConfig(options.ConfigPath, projectConfig); err != nil {
		return Report{}, err
	}

	report := newReport(cleanSlashPath(options.ConfigPath), intakeReportPath, "pass")
	report.SourceRefs = sourceRefs(root, projectConfig)
	report.SurfaceFingerprint = surfaceFingerprint(report.SourceRefs)
	if err := writeReport(reportRoot(options.ConfigPath), report); err != nil {
		return Report{}, err
	}
	return report, nil
}

func CheckProject(configPath string) (Report, error) {
	if configPath == "" {
		configPath = config.DefaultConfigPath
	}
	projectConfig, err := ReadProjectConfig(configPath)
	if err != nil {
		return Report{}, err
	}

	root := projectRoot(configPath)
	report := newReport(cleanSlashPath(configPath), checkReportPath, "pass")
	report.SourceRefs = sourceRefs(root, projectConfig)
	for _, entry := range projectConfig.Sources.OpenAPI {
		report.Findings = append(report.Findings, checkOpenAPI(root, entry)...)
	}
	for _, entry := range projectConfig.Sources.Docs {
		report.Findings = append(report.Findings, checkDocs(root, entry)...)
	}
	report.Status = statusFromFindings(report.Findings)
	report.SurfaceFingerprint = surfaceFingerprint(report.SourceRefs)
	if err := writeReport(reportRoot(configPath), report); err != nil {
		return Report{}, err
	}
	return report, nil
}

func ReadProjectConfig(path string) (ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ProjectConfig{}, fmt.Errorf("read config %s: %w", cleanSlashPath(path), err)
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return ProjectConfig{}, fmt.Errorf("read config %s: empty config", cleanSlashPath(path))
	}
	if trimmed[0] == '{' {
		var parsed ProjectConfig
		if err := json.Unmarshal(trimmed, &parsed); err != nil {
			return ProjectConfig{}, fmt.Errorf("parse config %s: %w", cleanSlashPath(path), err)
		}
		return validateProjectConfig(path, parsed)
	}
	return parseProjectConfigYAML(path, data)
}

func newReport(configPath, reportPath, status string) Report {
	return Report{
		ObjectType:    "lumyn.source_check",
		SchemaVersion: schemaVersion,
		Status:        status,
		ConfigPath:    configPath,
		ReportPath:    cleanSlashPath(reportPath),
		SourceRefs:    []SourceRef{},
		Findings:      []Finding{},
	}
}

func normalizedInitOptions(options InitOptions) InitOptions {
	if options.ConfigPath == "" {
		options.ConfigPath = config.DefaultConfigPath
	}
	if options.OpenAPIPath == "" {
		options.OpenAPIPath = "./openapi.yaml"
	}
	if options.DocsPath == "" {
		options.DocsPath = "./docs"
	}
	return options
}

func writeProjectConfig(path string, projectConfig ProjectConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return fmt.Errorf("create config directory %s: %w", cleanSlashPath(filepath.Dir(path)), err)
	}
	var builder strings.Builder
	builder.WriteString("version: 1\n\n")
	builder.WriteString("sources:\n")
	builder.WriteString("  openapi:\n")
	for _, entry := range projectConfig.Sources.OpenAPI {
		builder.WriteString(fmt.Sprintf("    - id: %s\n", yamlScalar(entry.ID)))
		builder.WriteString(fmt.Sprintf("      path: %s\n", yamlScalar(entry.Path)))
	}
	builder.WriteString("  docs:\n")
	for _, entry := range projectConfig.Sources.Docs {
		builder.WriteString(fmt.Sprintf("    - id: %s\n", yamlScalar(entry.ID)))
		builder.WriteString(fmt.Sprintf("      path: %s\n", yamlScalar(entry.Path)))
	}
	builder.WriteString("\n")
	builder.WriteString("env:\n")
	builder.WriteString(fmt.Sprintf("  base_url: %s\n", yamlScalar(projectConfig.Env.BaseURL)))
	builder.WriteString(fmt.Sprintf("  api_key: %s\n\n", yamlScalar(projectConfig.Env.APIKey)))
	builder.WriteString("redaction:\n")
	builder.WriteString("  headers:\n")
	for _, header := range projectConfig.Redact.Headers {
		builder.WriteString(fmt.Sprintf("    - %s\n", yamlScalar(header)))
	}
	builder.WriteString("  fields:\n")
	for _, field := range projectConfig.Redact.Fields {
		builder.WriteString(fmt.Sprintf("    - %s\n", yamlScalar(field)))
	}
	builder.WriteString("\n")
	builder.WriteString("verify:\n")
	builder.WriteString(fmt.Sprintf("  default_strategy: %s\n", yamlScalar(projectConfig.Verify.DefaultStrategy)))
	builder.WriteString("  replay:\n")
	builder.WriteString(fmt.Sprintf("    allow_network: %t\n", projectConfig.Verify.Replay.AllowNetwork))
	return os.WriteFile(path, []byte(builder.String()), 0o644)
}

func parseProjectConfigYAML(path string, data []byte) (ProjectConfig, error) {
	parsed := ProjectConfig{Version: 1}
	var section string
	var current *SourceEntry
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		switch trimmed {
		case "sources:":
			section = "sources"
			current = nil
			continue
		case "openapi:":
			if section == "sources" || section == "docs" {
				section = "openapi"
			}
			current = nil
			continue
		case "docs:":
			if section == "sources" || section == "openapi" {
				section = "docs"
			}
			current = nil
			continue
		case "env:", "redaction:", "verify:":
			section = ""
			current = nil
			continue
		}
		if strings.HasPrefix(trimmed, "version:") {
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			key, value, ok := yamlKeyValue(strings.TrimSpace(strings.TrimPrefix(trimmed, "- ")))
			if !ok || (key != "id" && key != "path") {
				current = nil
				continue
			}
			entry := SourceEntry{}
			switch key {
			case "id":
				entry.ID = value
			case "path":
				entry.Path = value
			}
			switch section {
			case "openapi":
				parsed.Sources.OpenAPI = append(parsed.Sources.OpenAPI, entry)
				current = &parsed.Sources.OpenAPI[len(parsed.Sources.OpenAPI)-1]
			case "docs":
				parsed.Sources.Docs = append(parsed.Sources.Docs, entry)
				current = &parsed.Sources.Docs[len(parsed.Sources.Docs)-1]
			default:
				current = nil
			}
			continue
		}
		if strings.HasPrefix(trimmed, "id:") && current != nil {
			current.ID = parseYAMLValue(strings.TrimPrefix(trimmed, "id:"))
			continue
		}
		if strings.HasPrefix(trimmed, "path:") && current != nil {
			current.Path = parseYAMLValue(strings.TrimPrefix(trimmed, "path:"))
		}
	}
	if err := scanner.Err(); err != nil {
		return ProjectConfig{}, fmt.Errorf("parse config %s: %w", cleanSlashPath(path), err)
	}
	return validateProjectConfig(path, parsed)
}

func validateProjectConfig(path string, projectConfig ProjectConfig) (ProjectConfig, error) {
	if projectConfig.Version == 0 {
		projectConfig.Version = 1
	}
	if len(projectConfig.Sources.OpenAPI) == 0 {
		return ProjectConfig{}, fmt.Errorf("parse config %s: missing sources.openapi", cleanSlashPath(path))
	}
	if len(projectConfig.Sources.Docs) == 0 {
		return ProjectConfig{}, fmt.Errorf("parse config %s: missing sources.docs", cleanSlashPath(path))
	}
	for _, entry := range append(projectConfig.Sources.OpenAPI, projectConfig.Sources.Docs...) {
		if entry.ID == "" || entry.Path == "" {
			return ProjectConfig{}, fmt.Errorf("parse config %s: source entries require id and path", cleanSlashPath(path))
		}
	}
	return projectConfig, nil
}

func sourceRefs(root string, projectConfig ProjectConfig) []SourceRef {
	refs := make([]SourceRef, 0, len(projectConfig.Sources.OpenAPI)+len(projectConfig.Sources.Docs))
	for _, entry := range projectConfig.Sources.OpenAPI {
		refs = append(refs, SourceRef{
			ID:   entry.ID,
			Kind: "openapi",
			Path: cleanSlashPath(entry.Path),
			Hash: hashPath(root, entry.Path),
		})
	}
	for _, entry := range projectConfig.Sources.Docs {
		refs = append(refs, SourceRef{
			ID:   entry.ID,
			Kind: "docs",
			Path: cleanSlashPath(entry.Path),
			Hash: hashDocsPath(root, entry.Path),
		})
	}
	return refs
}

func checkOpenAPI(root string, entry SourceEntry) []Finding {
	sourcePath := cleanSlashPath(entry.Path)
	data, err := os.ReadFile(resolveProjectPath(root, entry.Path))
	if err != nil {
		return []Finding{{
			Kind:     "context_missing",
			Severity: "warning",
			Message:  fmt.Sprintf("OpenAPI source %s could not be read", sourcePath),
			Reference: Reference{
				Path:   sourcePath,
				Object: "sources.openapi." + entry.ID,
			},
			FixTarget:         "sources.openapi.path",
			WorkflowRelevance: "Agents cannot choose tools or validators from an unreadable OpenAPI source.",
		}}
	}

	operations, hasSecuritySchemes, parseErr := parseOpenAPI(data)
	if parseErr != nil {
		return []Finding{{
			Kind:     "command_error",
			Severity: "error",
			Message:  fmt.Sprintf("OpenAPI source %s did not parse: %v", sourcePath, parseErr),
			Reference: Reference{
				Path:   sourcePath,
				Object: "openapi",
			},
			FixTarget:         "sources.openapi",
			WorkflowRelevance: "Agent probes need a parseable OpenAPI surface before workflow checks can be trusted.",
		}}
	}

	findings := findingsForOperations(sourcePath, operations)
	if !hasSecuritySchemes {
		findings = append(findings, Finding{
			Kind:     "auth_confusion",
			Severity: "warning",
			Message:  "OpenAPI source does not document components.securitySchemes",
			Reference: Reference{
				Path:        sourcePath,
				JSONPointer: "#/components/securitySchemes",
				Object:      "components.securitySchemes",
			},
			FixTarget:         "auth_schemes",
			WorkflowRelevance: "Agents need explicit auth schemes to avoid guessing credentials or headers.",
		})
	}
	findings = append(findings, authScopeDescriptionFindings(sourcePath, data)...)
	return findings
}

func parseOpenAPI(data []byte) ([]openAPIOperation, bool, error) {
	var raw map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&raw); err == nil {
		var trailing any
		if trailingErr := decoder.Decode(&trailing); trailingErr != io.EOF {
			return nil, false, errors.New("JSON OpenAPI source contains trailing data")
		}
		return parseOpenAPIJSON(raw)
	}
	return parseOpenAPIYAML(data)
}

func parseOpenAPIJSON(raw map[string]any) ([]openAPIOperation, bool, error) {
	if _, ok := raw["openapi"]; !ok {
		if _, ok := raw["swagger"]; ok {
			return nil, false, errors.New("swagger 2.0 is not supported; use OpenAPI 3.x")
		}
		return nil, false, errors.New("missing openapi version")
	}
	paths, ok := raw["paths"].(map[string]any)
	if !ok || len(paths) == 0 {
		return nil, false, errors.New("missing paths object")
	}

	operations := []openAPIOperation{}
	pathNames := sortedKeys(paths)
	for _, pathName := range pathNames {
		pathItem, ok := paths[pathName].(map[string]any)
		if !ok {
			continue
		}
		pathParameters := parseJSONParameters(pathItem["parameters"], "#/paths/"+escapeJSONPointer(pathName)+"/parameters", raw)
		for _, method := range sortedHTTPMethods(pathItem) {
			operationValue, ok := pathItem[method].(map[string]any)
			if !ok {
				continue
			}
			pointer := "#/paths/" + escapeJSONPointer(pathName) + "/" + method
			requestBody, hasRequestBody := operationValue["requestBody"]
			operation := openAPIOperation{
				Path:              pathName,
				Method:            method,
				OperationID:       stringValue(operationValue["operationId"]),
				Summary:           stringValue(operationValue["summary"]),
				Description:       stringValue(operationValue["description"]),
				Deprecated:        boolValue(operationValue["deprecated"]),
				Pointer:           pointer,
				HasRequestBody:    hasRequestBody,
				HasRequestSchema:  hasRequestContentSchema(requestBody, raw),
				HasResponseSchema: has2xxResponseSchema(operationValue["responses"], raw),
				ReplacementHint:   hasReplacementHint(operationValue),
				Parameters:        append([]openAPIParameter{}, pathParameters...),
			}
			operation.Parameters = mergeOpenAPIParameterOverrides(
				operation.Parameters,
				parseJSONParameters(operationValue["parameters"], pointer+"/parameters", raw),
			)
			operations = append(operations, operation)
		}
	}
	if len(operations) == 0 {
		return nil, false, errors.New("paths object contains no HTTP operations")
	}
	return operations, hasSecuritySchemesJSON(raw), nil
}

func parseOpenAPIYAML(data []byte) ([]openAPIOperation, bool, error) {
	responseComponentSchemas := yamlComponentContentSchemaRefs(data, "responses")
	requestBodyComponentSchemas := yamlComponentContentSchemaRefs(data, "requestBodies")
	parameterComponents := yamlComponentParameters(data)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lineNo := 0
	seenVersion := false
	inPaths := false
	pathsIndent := -1
	currentPath := ""
	pathIndent := -1
	pathItemChildIndent := -1
	pathParameters := []openAPIParameter{}
	pathParametersIndent := -1
	var currentPathParameter *openAPIParameter
	var current *openAPIOperation
	operationIndent := -1
	requestBodyIndent := -1
	requestContentIndent := -1
	requestMediaIndent := -1
	requestMediaChildIndent := -1
	requestSchemaIndent := -1
	responsesIndent := -1
	response2xxIndent := -1
	responseContentIndent := -1
	responseMediaIndent := -1
	responseMediaChildIndent := -1
	responseSchemaIndent := -1
	parametersIndent := -1
	var currentParameter *openAPIParameter
	componentsIndent := -1
	componentsChildIndent := -1
	securitySchemesIndent := -1
	securitySchemeIndent := -1
	securitySchemes := false
	operations := []openAPIOperation{}

	flushOperation := func() {
		if current != nil {
			current.Parameters = effectiveOpenAPIParameters(current.Parameters)
			operations = append(operations, *current)
			current = nil
		}
		operationIndent = -1
		requestBodyIndent = -1
		requestContentIndent = -1
		requestMediaIndent = -1
		requestMediaChildIndent = -1
		requestSchemaIndent = -1
		responsesIndent = -1
		response2xxIndent = -1
		responseContentIndent = -1
		responseMediaIndent = -1
		responseMediaChildIndent = -1
		responseSchemaIndent = -1
		parametersIndent = -1
		currentParameter = nil
	}

	applyPathParameters := func() {
		if currentPath == "" || len(pathParameters) == 0 {
			return
		}
		for index := range operations {
			if operations[index].Path == currentPath {
				operations[index].Parameters = mergeOpenAPIParameters(operations[index].Parameters, pathParameters)
			}
		}
		if current != nil && current.Path == currentPath {
			current.Parameters = mergeOpenAPIParameters(current.Parameters, pathParameters)
		}
	}
	resetPathState := func() {
		currentPath = ""
		pathIndent = -1
		pathItemChildIndent = -1
		pathParameters = nil
		pathParametersIndent = -1
		currentPathParameter = nil
	}
	flushPathItem := func() {
		flushOperation()
		applyPathParameters()
		resetPathState()
	}

	for scanner.Scan() {
		lineNo++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || trimmed == "---" {
			continue
		}
		indent := leadingSpaces(raw)
		key, value, hasKey := yamlKeyValue(trimmed)
		if !hasKey {
			continue
		}
		if key == "swagger" && indent == 0 {
			return nil, false, errors.New("swagger 2.0 is not supported; use OpenAPI 3.x")
		}
		if key == "openapi" && indent == 0 {
			seenVersion = true
		}
		if key == "paths" && indent == 0 {
			flushPathItem()
			inPaths = true
			pathsIndent = indent
			operations = append(operations, yamlInlinePathOperations(value, parameterComponents, requestBodyComponentSchemas, responseComponentSchemas)...)
			continue
		}
		if inPaths && indent <= pathsIndent && key != "paths" {
			flushPathItem()
			inPaths = false
		}
		if componentsIndent >= 0 && indent <= componentsIndent && key != "components" {
			componentsIndent = -1
			componentsChildIndent = -1
			securitySchemesIndent = -1
			securitySchemeIndent = -1
		}
		if securitySchemesIndent >= 0 && indent <= securitySchemesIndent && key != "securitySchemes" {
			securitySchemesIndent = -1
			securitySchemeIndent = -1
		}
		if key == "components" && indent == 0 && !inPaths {
			componentsIndent = indent
			componentsChildIndent = -1
			securitySchemesIndent = -1
			securitySchemeIndent = -1
		}
		if componentsIndent >= 0 && indent > componentsIndent && (componentsChildIndent < 0 || indent <= componentsChildIndent) {
			componentsChildIndent = indent
		}
		if componentsChildIndent >= 0 && indent == componentsChildIndent && key == "securitySchemes" {
			securitySchemesIndent = indent
			securitySchemeIndent = -1
			if yamlInlineSecuritySchemesHaveType(value) {
				securitySchemes = true
			}
		} else if securitySchemesIndent >= 0 && indent > securitySchemesIndent {
			if securitySchemeIndent < 0 || indent <= securitySchemeIndent {
				securitySchemeIndent = indent
				if yamlInlineMapValueNonEmpty(value, "type") {
					securitySchemes = true
				}
			} else if indent > securitySchemeIndent && key == "type" && value != "" {
				securitySchemes = true
			}
		}
		if !inPaths {
			continue
		}
		if strings.HasPrefix(key, "/") && indent > pathsIndent {
			flushPathItem()
			currentPath = key
			pathIndent = indent
			pathItemChildIndent = -1
			pathParameters = nil
			pathParametersIndent = -1
			currentPathParameter = nil
			operations = append(operations, yamlInlinePathItemOperations(currentPath, "#/paths/"+escapeJSONPointer(currentPath), value, nil, parameterComponents, requestBodyComponentSchemas, responseComponentSchemas)...)
			continue
		}
		isPathItemChild := currentPath != "" && indent > pathIndent && (pathItemChildIndent < 0 || indent <= pathItemChildIndent)
		if isPathItemChild && pathItemChildIndent < 0 {
			pathItemChildIndent = indent
		}
		if pathParametersIndent >= 0 && indent <= pathParametersIndent {
			pathParametersIndent = -1
			currentPathParameter = nil
		}
		if isPathItemChild && key == "parameters" {
			flushOperation()
			pathParametersIndent = indent
			currentPathParameter = nil
			inlineParameters := yamlInlineParameterSequence(value, "#/paths/"+escapeJSONPointer(currentPath)+"/parameters", parameterComponents)
			if len(inlineParameters) > 0 {
				pathParameters = mergeOpenAPIParameterOverrides(pathParameters, inlineParameters)
				pathParametersIndent = -1
			}
			continue
		}
		if current == nil && pathParametersIndent >= 0 && indent > pathParametersIndent {
			pathParameters = parseYAMLParameterLine(
				pathParameters,
				&currentPathParameter,
				trimmed,
				key,
				value,
				"#/paths/"+escapeJSONPointer(currentPath)+"/parameters",
				lineNo,
				parameterComponents,
			)
			continue
		}
		if isPathItemChild && httpMethods[strings.ToLower(key)] {
			flushOperation()
			method := strings.ToLower(key)
			current = &openAPIOperation{
				Path:       currentPath,
				Method:     method,
				Pointer:    "#/paths/" + escapeJSONPointer(currentPath) + "/" + method,
				Line:       lineNo,
				Parameters: append([]openAPIParameter{}, pathParameters...),
			}
			applyYAMLInlineOperationValue(current, value, parameterComponents, requestBodyComponentSchemas, responseComponentSchemas)
			operationIndent = indent
			continue
		}
		if current == nil {
			continue
		}
		if indent <= pathIndent || indent <= operationIndent {
			flushOperation()
			continue
		}
		if requestBodyIndent >= 0 && indent <= requestBodyIndent {
			requestBodyIndent = -1
			requestContentIndent = -1
			requestMediaIndent = -1
			requestMediaChildIndent = -1
			requestSchemaIndent = -1
		}
		if requestContentIndent >= 0 && indent <= requestContentIndent {
			requestContentIndent = -1
			requestMediaIndent = -1
			requestMediaChildIndent = -1
			requestSchemaIndent = -1
		}
		if requestMediaIndent >= 0 && indent <= requestMediaIndent {
			requestMediaIndent = -1
			requestMediaChildIndent = -1
			requestSchemaIndent = -1
		}
		if responsesIndent >= 0 && indent <= responsesIndent {
			responsesIndent = -1
			response2xxIndent = -1
			responseContentIndent = -1
			responseMediaIndent = -1
			responseMediaChildIndent = -1
			responseSchemaIndent = -1
		}
		if response2xxIndent >= 0 && indent <= response2xxIndent {
			response2xxIndent = -1
			responseContentIndent = -1
			responseMediaIndent = -1
			responseMediaChildIndent = -1
			responseSchemaIndent = -1
		}
		if responseContentIndent >= 0 && indent <= responseContentIndent {
			responseContentIndent = -1
			responseMediaIndent = -1
			responseMediaChildIndent = -1
			responseSchemaIndent = -1
		}
		if responseMediaIndent >= 0 && indent <= responseMediaIndent {
			responseMediaIndent = -1
			responseMediaChildIndent = -1
			responseSchemaIndent = -1
		}
		if requestSchemaIndent >= 0 && indent <= requestSchemaIndent {
			requestSchemaIndent = -1
		}
		if responseSchemaIndent >= 0 && indent <= responseSchemaIndent {
			responseSchemaIndent = -1
		}
		if parametersIndent >= 0 && indent <= parametersIndent {
			parametersIndent = -1
			currentParameter = nil
		}
		if requestSchemaIndent >= 0 && indent > requestSchemaIndent {
			current.HasRequestSchema = true
		}
		if responseSchemaIndent >= 0 && indent > responseSchemaIndent {
			current.HasResponseSchema = true
		}
		directOperationField := requestBodyIndent < 0 && responsesIndent < 0 && parametersIndent < 0
		if directOperationField {
			switch key {
			case "operationId":
				current.OperationID = value
			case "summary":
				current.Summary = value
			case "description":
				current.Description = value
				if mentionsReplacement(value) {
					current.ReplacementHint = true
				}
			case "deprecated":
				current.Deprecated = value == "true"
			case "x-replacement", "x-replaced-by", "x-deprecated-replacement":
				current.ReplacementHint = value != ""
			case "requestBody":
				current.HasRequestBody = true
				requestBodyIndent = indent
				if refName, ok := yamlLocalComponentRef(value, "requestBodies"); ok && requestBodyComponentSchemas[refName] {
					current.HasRequestSchema = true
				}
				if yamlInlineContentContainerHasSchema(value) {
					current.HasRequestSchema = true
				}
			case "responses":
				responsesIndent = indent
				response2xxIndent = -1
			case "parameters":
				parametersIndent = indent
				inlineParameters := yamlInlineParameterSequence(value, current.Pointer+"/parameters", parameterComponents)
				if len(inlineParameters) > 0 {
					current.Parameters = mergeOpenAPIParameterOverrides(current.Parameters, inlineParameters)
					parametersIndent = -1
				}
			}
		}
		if requestBodyIndent >= 0 && indent > requestBodyIndent && key == "schema" {
			if requestMediaIndent >= 0 && (requestMediaChildIndent < 0 || indent == requestMediaChildIndent) {
				if yamlInlineValueHasEntries(value) {
					current.HasRequestSchema = true
				} else if strings.TrimSpace(value) == "" {
					requestSchemaIndent = indent
				}
			}
		}
		if requestBodyIndent >= 0 && indent > requestBodyIndent && key == "$ref" {
			if refName, ok := yamlLocalComponentRef(value, "requestBodies"); ok && requestBodyComponentSchemas[refName] {
				current.HasRequestSchema = true
			}
		}
		if requestBodyIndent >= 0 && indent > requestBodyIndent && key == "content" {
			requestContentIndent = indent
			requestMediaIndent = -1
		}
		if requestContentIndent >= 0 && indent > requestContentIndent && strings.Contains(key, "/") {
			requestMediaIndent = indent
			requestMediaChildIndent = -1
			if yamlInlineValueHasSchema(value) {
				current.HasRequestSchema = true
			}
		}
		if requestMediaIndent >= 0 && indent > requestMediaIndent && (requestMediaChildIndent < 0 || indent <= requestMediaChildIndent) {
			requestMediaChildIndent = indent
		}
		if responsesIndent >= 0 && indent > responsesIndent && is2xxStatusKey(key) {
			response2xxIndent = indent
			responseContentIndent = -1
			responseMediaIndent = -1
			if refName, ok := yamlLocalComponentRef(value, "responses"); ok && responseComponentSchemas[refName] {
				current.HasResponseSchema = true
			}
			if yamlInlineContentContainerHasSchema(value) {
				current.HasResponseSchema = true
			}
		}
		if response2xxIndent >= 0 && indent > response2xxIndent && key == "content" {
			responseContentIndent = indent
			responseMediaIndent = -1
		}
		if responseContentIndent >= 0 && indent > responseContentIndent && strings.Contains(key, "/") {
			responseMediaIndent = indent
			responseMediaChildIndent = -1
			if yamlInlineValueHasSchema(value) {
				current.HasResponseSchema = true
			}
		}
		if responseMediaIndent >= 0 && indent > responseMediaIndent && (responseMediaChildIndent < 0 || indent <= responseMediaChildIndent) {
			responseMediaChildIndent = indent
		}
		if responseMediaIndent >= 0 && indent > responseMediaIndent && key == "schema" {
			if responseMediaChildIndent < 0 || indent == responseMediaChildIndent {
				if yamlInlineValueHasEntries(value) {
					current.HasResponseSchema = true
				} else if strings.TrimSpace(value) == "" {
					responseSchemaIndent = indent
				}
			}
		}
		if response2xxIndent >= 0 && indent > response2xxIndent && key == "$ref" {
			if refName, ok := yamlLocalComponentRef(value, "responses"); ok && responseComponentSchemas[refName] {
				current.HasResponseSchema = true
			}
		}
		if parametersIndent >= 0 && indent > parametersIndent {
			current.Parameters = parseYAMLParameterLine(
				current.Parameters,
				&currentParameter,
				trimmed,
				key,
				value,
				current.Pointer+"/parameters",
				lineNo,
				parameterComponents,
			)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, false, err
	}
	flushPathItem()
	if !seenVersion {
		return nil, false, errors.New("missing openapi version")
	}
	if len(operations) == 0 {
		return nil, false, errors.New("paths object contains no HTTP operations")
	}
	return operations, securitySchemes, nil
}

func findingsForOperations(sourcePath string, operations []openAPIOperation) []Finding {
	findings := []Finding{}
	seenToolNames := map[string]openAPIOperation{}
	seenOperationIDs := map[string]openAPIOperation{}
	for _, operation := range operations {
		object := strings.ToUpper(operation.Method) + " " + operation.Path
		reference := Reference{Path: sourcePath, JSONPointer: operation.Pointer, Object: object}
		if operation.Line > 0 {
			reference.Line = operation.Line
		}
		if operation.Summary == "" && operation.Description == "" {
			findings = append(findings, Finding{
				Kind:              "docs_api_ambiguity",
				Severity:          "warning",
				Message:           fmt.Sprintf("%s lacks a summary or description for agent tool choice", object),
				Reference:         reference,
				FixTarget:         "operation_summary",
				WorkflowRelevance: "Agents choose tools from operation names and descriptions; missing metadata makes similar endpoints hard to distinguish.",
			})
		}
		if operation.OperationID == "" {
			findings = append(findings, Finding{
				Kind:              "wrong_tool_or_endpoint",
				Severity:          "warning",
				Message:           fmt.Sprintf("%s lacks operationId", object),
				Reference:         reference,
				FixTarget:         "operation_id",
				WorkflowRelevance: "Stable operation IDs help map workflow steps to the intended endpoint.",
			})
		} else {
			operationIDKey := strings.ToLower(strings.TrimSpace(operation.OperationID))
			if previous, ok := seenOperationIDs[operationIDKey]; ok {
				findings = append(findings, Finding{
					Kind:     "docs_api_ambiguity",
					Severity: "warning",
					Message: fmt.Sprintf("%s and %s %s reuse operationId %q",
						object, strings.ToUpper(previous.Method), previous.Path, operation.OperationID),
					Reference:         reference,
					FixTarget:         "operation_id_disambiguation",
					WorkflowRelevance: "Duplicate operation IDs can bind recorded workflow steps to the wrong endpoint.",
				})
			} else {
				seenOperationIDs[operationIDKey] = operation
			}
		}
		if mutatingMethod[operation.Method] && !operation.HasRequestSchema && !(operation.Method == "delete" && !operation.HasRequestBody) {
			findings = append(findings, Finding{
				Kind:              "validator_coverage_gap",
				Severity:          "warning",
				Message:           fmt.Sprintf("%s is mutating but lacks a request schema", object),
				Reference:         reference,
				FixTarget:         "request_schema",
				WorkflowRelevance: "Mutating workflow probes need request schemas to construct safe, bounded writes.",
			})
		}
		if !operation.HasResponseSchema {
			findings = append(findings, Finding{
				Kind:              "proof_gap",
				Severity:          "warning",
				Message:           fmt.Sprintf("%s lacks a 2xx response schema for validator read-back", object),
				Reference:         reference,
				FixTarget:         "response_schema",
				WorkflowRelevance: "Read-back validators need response schemas before Lumyn can claim strong proof.",
			})
		}
		for _, parameter := range operation.Parameters {
			if parameter.Name == "" || parameter.Description != "" {
				continue
			}
			paramReference := Reference{
				Path:        sourcePath,
				JSONPointer: parameter.Pointer,
				Line:        parameter.Line,
				Object:      object + " parameter " + parameter.Name,
			}
			findings = append(findings, Finding{
				Kind:              "source_missing_metadata",
				Severity:          "warning",
				Message:           fmt.Sprintf("%s parameter %q lacks a useful description", object, parameter.Name),
				Reference:         paramReference,
				FixTarget:         "parameter_description",
				WorkflowRelevance: "Parameter descriptions reduce state-binding mistakes in recorded workflows.",
			})
		}
		if operation.Deprecated && !operation.ReplacementHint && !mentionsReplacement(operation.Description) && !mentionsReplacement(operation.Summary) {
			findings = append(findings, Finding{
				Kind:              "docs_api_ambiguity",
				Severity:          "warning",
				Message:           fmt.Sprintf("%s is deprecated without replacement guidance", object),
				Reference:         reference,
				FixTarget:         "deprecated_operation_guidance",
				WorkflowRelevance: "Agents need replacement guidance to avoid recording workflows against stale endpoints.",
			})
		}
		toolKey := strings.ToLower(strings.TrimSpace(operation.Summary))
		if toolKey == "" {
			toolKey = strings.ToLower(strings.TrimSpace(operation.OperationID))
		}
		if toolKey == "" {
			continue
		}
		if previous, ok := seenToolNames[toolKey]; ok {
			findings = append(findings, Finding{
				Kind:     "docs_api_ambiguity",
				Severity: "warning",
				Message: fmt.Sprintf("%s and %s %s use indistinguishable agent-facing names",
					object, strings.ToUpper(previous.Method), previous.Path),
				Reference:         reference,
				FixTarget:         "operation_disambiguation",
				WorkflowRelevance: "Near-duplicate operation names increase wrong-tool risk during agent probing.",
			})
		} else {
			seenToolNames[toolKey] = operation
		}
	}
	return findings
}

func hasContentSchema(value any) bool {
	node, ok := value.(map[string]any)
	if !ok {
		return false
	}
	content, ok := node["content"].(map[string]any)
	if !ok {
		return false
	}
	for _, mediaValue := range content {
		media, ok := mediaValue.(map[string]any)
		if !ok {
			continue
		}
		if schema, ok := media["schema"]; ok && schema != nil {
			return true
		}
	}
	return false
}

func hasRequestContentSchema(value any, root map[string]any) bool {
	return hasContentSchema(resolveLocalComponentRef(value, root, "requestBodies"))
}

func has2xxResponseSchema(value any, root map[string]any) bool {
	responses, ok := value.(map[string]any)
	if !ok {
		return false
	}
	for status, responseValue := range responses {
		if !strings.HasPrefix(status, "2") {
			continue
		}
		if hasContentSchema(resolveLocalComponentRef(responseValue, root, "responses")) {
			return true
		}
	}
	return false
}

func resolveLocalComponentRef(value any, root map[string]any, component string) any {
	return resolveLocalComponentRefSeen(value, root, component, map[string]bool{})
}

func resolveLocalComponentRefSeen(value any, root map[string]any, component string, seen map[string]bool) any {
	node, ok := value.(map[string]any)
	if !ok {
		return value
	}
	ref := stringValue(node["$ref"])
	prefix := "#/components/" + component + "/"
	if !strings.HasPrefix(ref, prefix) {
		return value
	}
	if seen[ref] {
		return value
	}
	seen[ref] = true
	components, ok := root["components"].(map[string]any)
	if !ok {
		return value
	}
	group, ok := components[component].(map[string]any)
	if !ok {
		return value
	}
	key := unescapeJSONPointer(strings.TrimPrefix(ref, prefix))
	if resolved, ok := group[key]; ok {
		return resolveLocalComponentRefSeen(resolved, root, component, seen)
	}
	return value
}

func parseJSONParameters(value any, pointer string, root map[string]any) []openAPIParameter {
	values, ok := value.([]any)
	if !ok {
		return nil
	}
	parameters := make([]openAPIParameter, 0, len(values))
	for index, raw := range values {
		parameterValue, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		parameterPointer := fmt.Sprintf("%s/%d", pointer, index)
		if ref := stringValue(parameterValue["$ref"]); ref != "" {
			if resolved, resolvedPointer, ok := resolveLocalParameterRef(ref, root); ok {
				parameterValue = resolved
				parameterPointer = resolvedPointer
			} else {
				parameters = append(parameters, openAPIParameter{
					Name:    ref,
					Pointer: parameterPointer,
				})
				continue
			}
		}
		parameters = append(parameters, openAPIParameter{
			Name:        stringValue(parameterValue["name"]),
			In:          stringValue(parameterValue["in"]),
			Description: stringValue(parameterValue["description"]),
			Pointer:     parameterPointer,
		})
	}
	return parameters
}

func resolveLocalParameterRef(ref string, root map[string]any) (map[string]any, string, bool) {
	return resolveLocalParameterRefSeen(ref, root, map[string]bool{})
}

func resolveLocalParameterRefSeen(ref string, root map[string]any, seen map[string]bool) (map[string]any, string, bool) {
	prefix := "#/components/parameters/"
	if !strings.HasPrefix(ref, prefix) {
		return nil, "", false
	}
	if seen[ref] {
		return nil, "", false
	}
	seen[ref] = true
	components, ok := root["components"].(map[string]any)
	if !ok {
		return nil, "", false
	}
	parameters, ok := components["parameters"].(map[string]any)
	if !ok {
		return nil, "", false
	}
	key := unescapeJSONPointer(strings.TrimPrefix(ref, prefix))
	parameter, ok := parameters[key].(map[string]any)
	if !ok {
		return nil, "", false
	}
	if nextRef := stringValue(parameter["$ref"]); nextRef != "" {
		return resolveLocalParameterRefSeen(nextRef, root, seen)
	}
	return parameter, "#/components/parameters/" + escapeJSONPointer(key), true
}

func mergeOpenAPIParameters(parameters []openAPIParameter, shared []openAPIParameter) []openAPIParameter {
	if len(shared) == 0 {
		return parameters
	}
	merged := append([]openAPIParameter{}, parameters...)
	for _, candidate := range shared {
		duplicate := false
		for _, existing := range merged {
			if openAPIParameterIdentity(existing) != "" && openAPIParameterIdentity(existing) == openAPIParameterIdentity(candidate) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			merged = append(merged, candidate)
		}
	}
	return merged
}

func mergeOpenAPIParameterOverrides(base []openAPIParameter, overrides []openAPIParameter) []openAPIParameter {
	if len(overrides) == 0 {
		return effectiveOpenAPIParameters(base)
	}
	return effectiveOpenAPIParameters(append(append([]openAPIParameter{}, base...), overrides...))
}

func effectiveOpenAPIParameters(parameters []openAPIParameter) []openAPIParameter {
	effective := []openAPIParameter{}
	indexByIdentity := map[string]int{}
	for _, parameter := range parameters {
		identity := openAPIParameterIdentity(parameter)
		if identity == "" {
			effective = append(effective, parameter)
			continue
		}
		if index, ok := indexByIdentity[identity]; ok {
			effective[index] = parameter
			continue
		}
		indexByIdentity[identity] = len(effective)
		effective = append(effective, parameter)
	}
	return effective
}

func openAPIParameterIdentity(parameter openAPIParameter) string {
	if parameter.Name == "" || parameter.In == "" {
		return ""
	}
	return parameter.Name + "\x00" + parameter.In
}

func hasSecuritySchemesJSON(raw map[string]any) bool {
	components, ok := raw["components"].(map[string]any)
	if !ok {
		return false
	}
	schemes, ok := components["securitySchemes"].(map[string]any)
	if !ok {
		return false
	}
	for _, value := range schemes {
		scheme, ok := value.(map[string]any)
		if !ok {
			continue
		}
		if stringValue(scheme["type"]) != "" {
			return true
		}
	}
	return false
}

func authScopeDescriptionFindings(sourcePath string, data []byte) []Finding {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var raw map[string]any
	if err := decoder.Decode(&raw); err != nil {
		return yamlAuthScopeDescriptionFindings(sourcePath, data)
	}
	components, ok := raw["components"].(map[string]any)
	if !ok {
		return nil
	}
	schemes, ok := components["securitySchemes"].(map[string]any)
	if !ok {
		return nil
	}
	findings := []Finding{}
	for schemeName, schemeValue := range schemes {
		scheme, ok := schemeValue.(map[string]any)
		if !ok || strings.ToLower(stringValue(scheme["type"])) != "oauth2" {
			continue
		}
		flows, ok := scheme["flows"].(map[string]any)
		if !ok {
			continue
		}
		for flowName, flowValue := range flows {
			flow, ok := flowValue.(map[string]any)
			if !ok {
				continue
			}
			scopes, ok := flow["scopes"].(map[string]any)
			if !ok {
				continue
			}
			for scopeName, scopeValue := range scopes {
				if stringValue(scopeValue) != "" {
					continue
				}
				findings = append(findings, Finding{
					Kind:     "auth_confusion",
					Severity: "warning",
					Message:  fmt.Sprintf("OAuth scope %q in security scheme %q lacks a useful description", scopeName, schemeName),
					Reference: Reference{
						Path: sourcePath,
						JSONPointer: "#/components/securitySchemes/" + escapeJSONPointer(schemeName) +
							"/flows/" + escapeJSONPointer(flowName) + "/scopes/" + escapeJSONPointer(scopeName),
						Object: "oauth_scope",
					},
					FixTarget:         "auth_scope_description",
					WorkflowRelevance: "Agents need explicit scope meaning before selecting credentials or protected workflow actions.",
				})
			}
		}
	}
	return findings
}

func yamlAuthScopeDescriptionFindings(sourcePath string, data []byte) []Finding {
	findings := []Finding{}
	oauth2Schemes := yamlOAuth2SecuritySchemeNames(data)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	componentsIndent := -1
	schemesIndent := -1
	schemeIndent := -1
	flowsIndent := -1
	flowIndent := -1
	scopesIndent := -1
	currentScheme := ""
	currentFlow := ""
	currentSchemeOAuth2 := false
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || trimmed == "---" {
			continue
		}
		indent := leadingSpaces(scanner.Text())
		key, value, ok := yamlKeyValue(trimmed)
		if !ok {
			continue
		}
		if componentsIndent >= 0 && indent <= componentsIndent && key != "components" {
			componentsIndent = -1
			schemesIndent = -1
			schemeIndent = -1
			flowsIndent = -1
			flowIndent = -1
			scopesIndent = -1
			currentScheme = ""
			currentFlow = ""
			currentSchemeOAuth2 = false
		}
		if schemesIndent >= 0 && indent <= schemesIndent && key != "securitySchemes" {
			schemesIndent = -1
			schemeIndent = -1
			flowsIndent = -1
			flowIndent = -1
			scopesIndent = -1
			currentScheme = ""
			currentFlow = ""
			currentSchemeOAuth2 = false
		}
		if schemeIndent >= 0 && indent <= schemeIndent && key != currentScheme {
			schemeIndent = -1
			flowsIndent = -1
			flowIndent = -1
			scopesIndent = -1
			currentScheme = ""
			currentFlow = ""
			currentSchemeOAuth2 = false
		}
		if flowsIndent >= 0 && indent <= flowsIndent && key != "flows" {
			flowsIndent = -1
			flowIndent = -1
			scopesIndent = -1
			currentFlow = ""
		}
		if flowIndent >= 0 && indent <= flowIndent && key != currentFlow {
			flowIndent = -1
			scopesIndent = -1
			currentFlow = ""
		}
		if scopesIndent >= 0 && indent <= scopesIndent && key != "scopes" {
			scopesIndent = -1
		}
		switch {
		case key == "components" && indent == 0:
			componentsIndent = indent
			schemesIndent = -1
		case componentsIndent >= 0 && indent > componentsIndent && key == "securitySchemes":
			schemesIndent = indent
			schemeIndent = -1
			currentScheme = ""
			findings = append(findings, yamlInlineSecuritySchemeScopeFindings(sourcePath, value)...)
		case schemesIndent >= 0 && indent > schemesIndent && (schemeIndent < 0 || indent <= schemeIndent):
			currentScheme = key
			schemeIndent = indent
			flowsIndent = -1
			currentFlow = ""
			currentSchemeOAuth2 = oauth2Schemes[currentScheme] || yamlInlineMapValueEquals(value, "type", "oauth2")
			if currentSchemeOAuth2 {
				findings = append(findings, yamlInlineOAuthSchemeScopeFindings(sourcePath, currentScheme, value)...)
			}
		case currentScheme != "" && indent > schemeIndent && key == "type":
			currentSchemeOAuth2 = strings.EqualFold(value, "oauth2")
		case currentScheme != "" && indent > schemeIndent && key == "flows":
			flowsIndent = indent
			flowIndent = -1
			currentFlow = ""
		case flowsIndent >= 0 && indent > flowsIndent && (flowIndent < 0 || indent <= flowIndent):
			currentFlow = key
			flowIndent = indent
			scopesIndent = -1
		case currentFlow != "" && indent > flowIndent && key == "scopes":
			scopesIndent = indent
			if currentSchemeOAuth2 {
				findings = append(findings, yamlInlineAuthScopeDescriptionFindings(sourcePath, currentScheme, currentFlow, value)...)
			}
		case scopesIndent >= 0 && indent > scopesIndent:
			if !currentSchemeOAuth2 {
				continue
			}
			if value != "" {
				continue
			}
			findings = append(findings, Finding{
				Kind:     "auth_confusion",
				Severity: "warning",
				Message:  fmt.Sprintf("OAuth scope %q in security scheme %q lacks a useful description", key, currentScheme),
				Reference: Reference{
					Path: sourcePath,
					JSONPointer: "#/components/securitySchemes/" + escapeJSONPointer(currentScheme) +
						"/flows/" + escapeJSONPointer(currentFlow) + "/scopes/" + escapeJSONPointer(key),
					Object: "oauth_scope",
				},
				FixTarget:         "auth_scope_description",
				WorkflowRelevance: "Agents need explicit scope meaning before selecting credentials or protected workflow actions.",
			})
		}
	}
	return findings
}

func yamlInlineSecuritySchemeScopeFindings(sourcePath, value string) []Finding {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return nil
	}
	findings := []Finding{}
	for _, schemeEntry := range splitYAMLFlowEntries(value) {
		schemeName, schemeValue, ok := yamlKeyValue(strings.TrimSpace(schemeEntry))
		if !ok {
			continue
		}
		findings = append(findings, yamlInlineOAuthSchemeScopeFindings(sourcePath, schemeName, schemeValue)...)
	}
	return findings
}

func yamlInlineOAuthSchemeScopeFindings(sourcePath, schemeName, schemeValue string) []Finding {
	schemeValue = strings.TrimSpace(strings.Trim(schemeValue, "{} "))
	if !yamlInlineMapValueEquals(schemeValue, "type", "oauth2") {
		return nil
	}
	flowsValue, ok := yamlInlineMapValue(schemeValue, "flows")
	if !ok {
		return nil
	}
	flowsValue = strings.TrimSpace(strings.Trim(flowsValue, "{} "))
	findings := []Finding{}
	for _, flowEntry := range splitYAMLFlowEntries(flowsValue) {
		flowName, flowValue, ok := yamlKeyValue(strings.TrimSpace(flowEntry))
		if !ok {
			continue
		}
		flowValue = strings.TrimSpace(strings.Trim(flowValue, "{} "))
		scopesValue, ok := yamlInlineMapValue(flowValue, "scopes")
		if !ok {
			continue
		}
		findings = append(findings, yamlInlineAuthScopeDescriptionFindings(sourcePath, schemeName, flowName, scopesValue)...)
	}
	return findings
}

func applyYAMLInlineOperationValue(operation *openAPIOperation, value string, parameterComponents map[string]openAPIParameter, requestBodyComponentSchemas, responseComponentSchemas map[string]bool) {
	if operation == nil {
		return
	}
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return
	}
	for _, entry := range splitYAMLFlowEntries(value) {
		key, entryValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if !ok {
			continue
		}
		switch key {
		case "operationId":
			operation.OperationID = entryValue
		case "summary":
			operation.Summary = entryValue
		case "description":
			operation.Description = entryValue
			if mentionsReplacement(entryValue) {
				operation.ReplacementHint = true
			}
		case "deprecated":
			operation.Deprecated = entryValue == "true"
		case "x-replacement", "x-replaced-by", "x-deprecated-replacement":
			operation.ReplacementHint = entryValue != ""
		case "requestBody":
			operation.HasRequestBody = true
			if refName, ok := yamlLocalComponentRef(entryValue, "requestBodies"); ok && requestBodyComponentSchemas[refName] {
				operation.HasRequestSchema = true
			}
			if yamlInlineContentContainerHasSchema(entryValue) {
				operation.HasRequestSchema = true
			}
		case "responses":
			if yamlInlineResponsesHave2xxSchema(entryValue, responseComponentSchemas) {
				operation.HasResponseSchema = true
			}
		case "parameters":
			operation.Parameters = mergeOpenAPIParameterOverrides(
				operation.Parameters,
				yamlInlineParameterSequence(entryValue, operation.Pointer+"/parameters", parameterComponents),
			)
		}
	}
}

func yamlInlinePathOperations(value string, parameterComponents map[string]openAPIParameter, requestBodyComponentSchemas, responseComponentSchemas map[string]bool) []openAPIOperation {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return nil
	}
	operations := []openAPIOperation{}
	for _, pathEntry := range splitYAMLFlowEntries(value) {
		pathName, pathValue, ok := yamlKeyValue(strings.TrimSpace(pathEntry))
		if !ok || !strings.HasPrefix(pathName, "/") {
			continue
		}
		pathPointer := "#/paths/" + escapeJSONPointer(pathName)
		operations = append(operations, yamlInlinePathItemOperations(pathName, pathPointer, pathValue, nil, parameterComponents, requestBodyComponentSchemas, responseComponentSchemas)...)
	}
	return operations
}

func yamlInlinePathItemOperations(pathName, pathPointer, value string, inheritedParameters []openAPIParameter, parameterComponents map[string]openAPIParameter, requestBodyComponentSchemas, responseComponentSchemas map[string]bool) []openAPIOperation {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return nil
	}
	pathParameters := append([]openAPIParameter{}, inheritedParameters...)
	if parametersValue, ok := yamlInlineMapValue(value, "parameters"); ok {
		pathParameters = mergeOpenAPIParameterOverrides(pathParameters, yamlInlineParameterSequence(parametersValue, pathPointer+"/parameters", parameterComponents))
	}
	operations := []openAPIOperation{}
	for _, operationEntry := range splitYAMLFlowEntries(value) {
		method, operationValue, ok := yamlKeyValue(strings.TrimSpace(operationEntry))
		if !ok {
			continue
		}
		method = strings.ToLower(method)
		if !httpMethods[method] {
			continue
		}
		operation := openAPIOperation{
			Path:       pathName,
			Method:     method,
			Pointer:    pathPointer + "/" + method,
			Parameters: append([]openAPIParameter{}, pathParameters...),
		}
		applyYAMLInlineOperationValue(&operation, operationValue, parameterComponents, requestBodyComponentSchemas, responseComponentSchemas)
		operation.Parameters = effectiveOpenAPIParameters(operation.Parameters)
		operations = append(operations, operation)
	}
	return operations
}

func yamlInlineResponsesHave2xxSchema(value string, responseComponentSchemas map[string]bool) bool {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return false
	}
	for _, entry := range splitYAMLFlowEntries(value) {
		key, entryValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if !ok || !is2xxStatusKey(key) {
			continue
		}
		if refName, ok := yamlLocalComponentRef(entryValue, "responses"); ok && responseComponentSchemas[refName] {
			return true
		}
		if yamlInlineContentContainerHasSchema(entryValue) {
			return true
		}
	}
	return false
}

func yamlInlineAuthScopeDescriptionFindings(sourcePath, schemeName, flowName, value string) []Finding {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return nil
	}
	findings := []Finding{}
	for _, entry := range splitYAMLFlowEntries(value) {
		scopeName, scopeValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if !ok || scopeValue != "" {
			continue
		}
		findings = append(findings, Finding{
			Kind:     "auth_confusion",
			Severity: "warning",
			Message:  fmt.Sprintf("OAuth scope %q in security scheme %q lacks a useful description", scopeName, schemeName),
			Reference: Reference{
				Path: sourcePath,
				JSONPointer: "#/components/securitySchemes/" + escapeJSONPointer(schemeName) +
					"/flows/" + escapeJSONPointer(flowName) + "/scopes/" + escapeJSONPointer(scopeName),
				Object: "oauth_scope",
			},
			FixTarget:         "auth_scope_description",
			WorkflowRelevance: "Agents need explicit scope meaning before selecting credentials or protected workflow actions.",
		})
	}
	return findings
}

func yamlOAuth2SecuritySchemeNames(data []byte) map[string]bool {
	names := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	componentsIndent := -1
	schemesIndent := -1
	schemeIndent := -1
	currentScheme := ""
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || trimmed == "---" {
			continue
		}
		indent := leadingSpaces(scanner.Text())
		key, value, ok := yamlKeyValue(trimmed)
		if !ok {
			continue
		}
		if componentsIndent >= 0 && indent <= componentsIndent && key != "components" {
			componentsIndent = -1
			schemesIndent = -1
			schemeIndent = -1
			currentScheme = ""
		}
		if schemesIndent >= 0 && indent <= schemesIndent && key != "securitySchemes" {
			schemesIndent = -1
			schemeIndent = -1
			currentScheme = ""
		}
		if schemeIndent >= 0 && indent <= schemeIndent && key != currentScheme {
			schemeIndent = -1
			currentScheme = ""
		}
		switch {
		case key == "components" && indent == 0:
			componentsIndent = indent
			schemesIndent = -1
		case componentsIndent >= 0 && indent > componentsIndent && key == "securitySchemes":
			schemesIndent = indent
			schemeIndent = -1
			currentScheme = ""
		case schemesIndent >= 0 && indent > schemesIndent && (schemeIndent < 0 || indent <= schemeIndent):
			currentScheme = key
			schemeIndent = indent
		case currentScheme != "" && indent > schemeIndent && key == "type" && strings.EqualFold(value, "oauth2"):
			names[currentScheme] = true
		}
	}
	return names
}

func hasReplacementHint(operation map[string]any) bool {
	for _, key := range []string{"x-replacement", "x-replaced-by", "x-deprecated-replacement"} {
		if stringValue(operation[key]) != "" {
			return true
		}
	}
	return mentionsReplacement(stringValue(operation["description"])) || mentionsReplacement(stringValue(operation["summary"]))
}

func sortedHTTPMethods(pathItem map[string]any) []string {
	methods := []string{}
	for key := range pathItem {
		lower := strings.ToLower(key)
		if httpMethods[lower] {
			methods = append(methods, lower)
		}
	}
	sort.Strings(methods)
	return methods
}

func sortedKeys(input map[string]any) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case json.Number:
		return typed.String()
	default:
		return ""
	}
}

func boolValue(value any) bool {
	typed, ok := value.(bool)
	return ok && typed
}

func mentionsReplacement(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "replace") || strings.Contains(lower, "use ") || strings.Contains(lower, "instead")
}

func projectRoot(configPath string) string {
	dir := filepath.Dir(configPath)
	if dir == "" {
		return "."
	}
	return dir
}

func reportRoot(configPath string) string {
	if configPath == "" {
		return "."
	}
	if !filepath.IsAbs(configPath) {
		return "."
	}
	dir := filepath.Dir(configPath)
	if strings.HasPrefix(filepath.Base(dir), ".") {
		parent := filepath.Dir(dir)
		if parent == "" {
			return "."
		}
		return parent
	}
	return dir
}

func resolveProjectPath(root, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(root, filepath.FromSlash(path)))
}

func sourcePathForConfigRoot(root, path string, preferWorkingDirectory bool) string {
	if path == "" {
		return ""
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return cleanSlashPath(path)
	}
	sourcePath := filepath.FromSlash(filepath.Clean(path))
	sourceAbs := sourcePath
	if !filepath.IsAbs(sourceAbs) {
		cwd, cwdErr := os.Getwd()
		cwdCandidate := sourcePath
		if cwdErr == nil {
			cwdCandidate = filepath.Join(cwd, sourcePath)
		}
		rootCandidate := filepath.Join(rootAbs, sourcePath)
		switch {
		case preferWorkingDirectory && cwdErr == nil && pathExists(cwdCandidate):
			sourceAbs = cwdCandidate
		case pathExists(rootCandidate):
			sourceAbs = rootCandidate
		case cwdErr == nil:
			sourceAbs = cwdCandidate
		default:
			return cleanSlashPath(path)
		}
	}
	rel, err := filepath.Rel(rootAbs, sourceAbs)
	if err != nil {
		return cleanSlashPath(path)
	}
	return cleanSlashPath(rel)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func relativePath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return cleanSlashPath(path)
	}
	return cleanSlashPath(rel)
}

func cleanSlashPath(path string) string {
	if path == "" {
		return ""
	}
	return filepath.ToSlash(filepath.Clean(path))
}
