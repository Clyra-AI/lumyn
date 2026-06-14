package source

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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
	markdownLinkRE = regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
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
	projectConfig := ProjectConfig{
		Version: 1,
		Sources: SourceConfig{
			OpenAPI: []SourceEntry{{ID: "public_api", Path: filepath.ToSlash(options.OpenAPIPath)}},
			Docs:    []SourceEntry{{ID: "docs", Path: filepath.ToSlash(options.DocsPath)}},
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
	if err := writeReport(root, report); err != nil {
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
	if err := writeReport(root, report); err != nil {
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
	return findings
}

func parseOpenAPI(data []byte) ([]openAPIOperation, bool, error) {
	var raw map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&raw); err == nil {
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
			operation.Parameters = append(operation.Parameters, parseJSONParameters(operationValue["parameters"], pointer+"/parameters", raw)...)
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
	responsesIndent := -1
	response2xxIndent := -1
	responseContentIndent := -1
	responseMediaIndent := -1
	parametersIndent := -1
	var currentParameter *openAPIParameter
	componentsIndent := -1
	securitySchemesIndent := -1
	securitySchemes := false
	operations := []openAPIOperation{}

	flushOperation := func() {
		if current != nil {
			operations = append(operations, *current)
			current = nil
		}
		operationIndent = -1
		requestBodyIndent = -1
		requestContentIndent = -1
		requestMediaIndent = -1
		responsesIndent = -1
		response2xxIndent = -1
		responseContentIndent = -1
		responseMediaIndent = -1
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
		if key == "swagger" {
			return nil, false, errors.New("swagger 2.0 is not supported; use OpenAPI 3.x")
		}
		if key == "openapi" {
			seenVersion = true
		}
		if key == "paths" {
			flushPathItem()
			inPaths = true
			pathsIndent = indent
			continue
		}
		if inPaths && indent <= pathsIndent && key != "paths" {
			flushPathItem()
			inPaths = false
		}
		if componentsIndent >= 0 && indent <= componentsIndent && key != "components" {
			componentsIndent = -1
			securitySchemesIndent = -1
		}
		if securitySchemesIndent >= 0 && indent <= securitySchemesIndent && key != "securitySchemes" {
			securitySchemesIndent = -1
		}
		if key == "components" && !inPaths {
			componentsIndent = indent
			securitySchemesIndent = -1
		}
		if componentsIndent >= 0 && indent > componentsIndent && key == "securitySchemes" {
			securitySchemesIndent = indent
			if yamlInlineValueHasEntries(value) {
				securitySchemes = true
			}
		} else if securitySchemesIndent >= 0 && indent > securitySchemesIndent {
			securitySchemes = true
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
		}
		if requestContentIndent >= 0 && indent <= requestContentIndent {
			requestContentIndent = -1
			requestMediaIndent = -1
		}
		if requestMediaIndent >= 0 && indent <= requestMediaIndent {
			requestMediaIndent = -1
		}
		if responsesIndent >= 0 && indent <= responsesIndent {
			responsesIndent = -1
			response2xxIndent = -1
			responseContentIndent = -1
			responseMediaIndent = -1
		}
		if response2xxIndent >= 0 && indent <= response2xxIndent {
			response2xxIndent = -1
			responseContentIndent = -1
			responseMediaIndent = -1
		}
		if responseContentIndent >= 0 && indent <= responseContentIndent {
			responseContentIndent = -1
			responseMediaIndent = -1
		}
		if responseMediaIndent >= 0 && indent <= responseMediaIndent {
			responseMediaIndent = -1
		}
		if parametersIndent >= 0 && indent <= parametersIndent {
			parametersIndent = -1
			currentParameter = nil
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
			case "x-replacement", "x-replaced-by":
				current.ReplacementHint = value != ""
			case "requestBody":
				current.HasRequestBody = true
				requestBodyIndent = indent
				if refName, ok := yamlLocalComponentRef(value, "requestBodies"); ok && requestBodyComponentSchemas[refName] {
					current.HasRequestSchema = true
				}
			case "responses":
				responsesIndent = indent
				response2xxIndent = -1
			case "parameters":
				parametersIndent = indent
			}
		}
		if requestBodyIndent >= 0 && indent > requestBodyIndent && key == "schema" {
			if requestMediaIndent >= 0 && indent > requestMediaIndent {
				current.HasRequestSchema = true
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
		}
		if responsesIndent >= 0 && indent > responsesIndent && is2xxStatusKey(key) {
			response2xxIndent = indent
			responseContentIndent = -1
			responseMediaIndent = -1
			if refName, ok := yamlLocalComponentRef(value, "responses"); ok && responseComponentSchemas[refName] {
				current.HasResponseSchema = true
			}
		}
		if response2xxIndent >= 0 && indent > response2xxIndent && key == "content" {
			responseContentIndent = indent
			responseMediaIndent = -1
		}
		if responseContentIndent >= 0 && indent > responseContentIndent && strings.Contains(key, "/") {
			responseMediaIndent = indent
		}
		if responseMediaIndent >= 0 && indent > responseMediaIndent && key == "schema" {
			current.HasResponseSchema = true
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
		if previous, ok := seenToolNames[toolKey]; ok && previous.Path != operation.Path {
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

func checkDocs(root string, entry SourceEntry) []Finding {
	docsPath := cleanSlashPath(entry.Path)
	absPath := resolveProjectPath(root, entry.Path)
	info, err := os.Stat(absPath)
	if err != nil {
		return []Finding{{
			Kind:     "context_missing",
			Severity: "warning",
			Message:  fmt.Sprintf("docs source %s could not be read", docsPath),
			Reference: Reference{
				Path:   docsPath,
				Object: "sources.docs." + entry.ID,
			},
			FixTarget:         "sources.docs.path",
			WorkflowRelevance: "Agents need local docs context for auth, retries, pagination, and proof guidance.",
		}}
	}

	findings := []Finding{}
	readableFiles := 0
	var combined strings.Builder
	visitFile := func(path string, d fs.DirEntry) {
		if d.IsDir() {
			return
		}
		if !isDocsFile(path) {
			return
		}
		relPath := relativePath(root, path)
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			findings = append(findings, Finding{
				Kind:     "context_missing",
				Severity: "warning",
				Message:  fmt.Sprintf("docs file %s could not be read", relPath),
				Reference: Reference{
					Path:   relPath,
					Object: "docs_file",
				},
				FixTarget:         "docs_file",
				WorkflowRelevance: "Unreadable docs can hide auth, retry, or validation instructions from agent probes.",
			})
			return
		}
		readableFiles++
		combined.Write(bytes.ToLower(data))
		combined.WriteByte('\n')
		findings = append(findings, brokenLocalReferenceFindings(root, path, data)...)
	}

	if info.IsDir() {
		walkErr := filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				findings = append(findings, Finding{
					Kind:     "context_missing",
					Severity: "warning",
					Message:  fmt.Sprintf("docs path %s could not be inspected: %v", relativePath(root, path), err),
					Reference: Reference{
						Path:   relativePath(root, path),
						Object: "docs_path",
					},
					FixTarget:         "docs_path",
					WorkflowRelevance: "Unreadable docs can hide workflow constraints from source checks.",
				})
				return nil
			}
			if d.IsDir() && path != absPath && shouldSkipGeneratedSourceDir(d.Name()) {
				return filepath.SkipDir
			}
			visitFile(path, d)
			return nil
		})
		if walkErr != nil {
			findings = append(findings, Finding{
				Kind:     "context_missing",
				Severity: "warning",
				Message:  fmt.Sprintf("docs source %s could not be walked: %v", docsPath, walkErr),
				Reference: Reference{
					Path:   docsPath,
					Object: "sources.docs." + entry.ID,
				},
				FixTarget:         "sources.docs.path",
				WorkflowRelevance: "Agents need readable local docs for workflow-relevant source context.",
			})
		}
	} else {
		visitFile(absPath, fileEntry{name: filepath.Base(absPath)})
	}

	if readableFiles == 0 {
		findings = append(findings, Finding{
			Kind:     "context_missing",
			Severity: "warning",
			Message:  fmt.Sprintf("docs source %s contains no readable docs files", docsPath),
			Reference: Reference{
				Path:   docsPath,
				Object: "sources.docs." + entry.ID,
			},
			FixTarget:         "docs_content",
			WorkflowRelevance: "Agent probes need readable docs for setup, auth, retry, and validation instructions.",
		})
		return findings
	}

	if missingOperationalGuidance(combined.String()) {
		findings = append(findings, Finding{
			Kind:     "docs_api_ambiguity",
			Severity: "warning",
			Message:  "local docs do not mention retry, rate-limit, pagination, or idempotency guidance",
			Reference: Reference{
				Path:   docsPath,
				Object: "docs_guidance",
			},
			FixTarget:         "operational_docs",
			WorkflowRelevance: "Retries, rate limits, pagination, and idempotency affect agent workflow stability and write safety.",
		})
	}
	return findings
}

func brokenLocalReferenceFindings(root, docPath string, data []byte) []Finding {
	findings := []Finding{}
	lines := bytes.Split(data, []byte("\n"))
	for index, line := range lines {
		matches := markdownLinkRE.FindAllSubmatch(line, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			target := cleanMarkdownLinkTarget(string(match[1]))
			if target == "" || isExternalReference(target) {
				continue
			}
			targetPath := strings.SplitN(target, "#", 2)[0]
			if targetPath == "" {
				continue
			}
			resolved := resolveMarkdownLinkTarget(root, docPath, targetPath)
			if _, err := os.Stat(resolved); err == nil {
				continue
			}
			findings = append(findings, Finding{
				Kind:     "context_missing",
				Severity: "warning",
				Message:  fmt.Sprintf("docs file %s links to missing local reference %s", relativePath(root, docPath), target),
				Reference: Reference{
					Path: relativePath(root, docPath),
					Line: index + 1,
				},
				FixTarget:         "docs_local_reference",
				WorkflowRelevance: "Broken local docs links can hide workflow setup, auth, or validation instructions from agents.",
			})
		}
	}
	return findings
}

func writeReport(root string, report Report) error {
	path := resolveProjectPath(root, report.ReportPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create report directory %s: %w", cleanSlashPath(filepath.Dir(path)), err)
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("encode source report %s: %w", report.ReportPath, err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write source report %s: %w", report.ReportPath, err)
	}
	return nil
}

func statusFromFindings(findings []Finding) string {
	status := "pass"
	for _, finding := range findings {
		if finding.Severity == "error" {
			return "fail"
		}
		if finding.Severity == "warning" {
			status = "warning"
		}
	}
	return status
}

func HasErrorFindings(findings []Finding) bool {
	for _, finding := range findings {
		if finding.Severity == "error" {
			return true
		}
	}
	return false
}

func FirstFinding(findings []Finding) (Finding, bool) {
	if len(findings) == 0 {
		return Finding{}, false
	}
	return findings[0], true
}

func sourceFindingHasKind(findings []Finding, kind string) bool {
	for _, finding := range findings {
		if finding.Kind == kind {
			return true
		}
	}
	return false
}

func ContainsProofGap(findings []Finding) bool {
	return sourceFindingHasKind(findings, "proof_gap")
}

func ContainsSecurityRelevantFinding(findings []Finding) bool {
	for _, finding := range findings {
		switch finding.Kind {
		case "auth_confusion", "data_exposure_risk", "forbidden_endpoint_call", "scope_escalation", "unexpected_write_action":
			return true
		}
	}
	return false
}

func hashPath(root, path string) string {
	resolved := resolveProjectPath(root, path)
	info, err := os.Stat(resolved)
	if err != nil {
		return ""
	}
	if !info.IsDir() {
		return hashFile(resolved)
	}
	fileHashes := []string{}
	_ = filepath.WalkDir(resolved, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		hash := hashFile(path)
		if hash == "" {
			return nil
		}
		fileHashes = append(fileHashes, relativePath(root, path)+"="+hash)
		return nil
	})
	sort.Strings(fileHashes)
	digest := sha256.Sum256([]byte(strings.Join(fileHashes, "\n")))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func hashDocsPath(root, path string) string {
	resolved := resolveProjectPath(root, path)
	info, err := os.Stat(resolved)
	if err != nil {
		return ""
	}
	if !info.IsDir() {
		return hashFile(resolved)
	}
	fileHashes := []string{}
	_ = filepath.WalkDir(resolved, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if shouldSkipGeneratedSourceDir(d.Name()) && path != resolved {
				return filepath.SkipDir
			}
			return nil
		}
		if !isDocsFile(path) {
			return nil
		}
		hash := hashFile(path)
		if hash == "" {
			return nil
		}
		fileHashes = append(fileHashes, relativePath(root, path)+"="+hash)
		return nil
	})
	sort.Strings(fileHashes)
	digest := sha256.Sum256([]byte(strings.Join(fileHashes, "\n")))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func shouldSkipGeneratedSourceDir(name string) bool {
	switch name {
	case ".git", ".factory", ".factoryd", "runs":
		return true
	default:
		return false
	}
}

func hashFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	digest := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func surfaceFingerprint(refs []SourceRef) string {
	parts := make([]string, 0, len(refs))
	for _, ref := range refs {
		parts = append(parts, ref.Kind+":"+ref.ID+":"+ref.Path+":"+ref.Hash)
	}
	sort.Strings(parts)
	digest := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return "sha256:" + hex.EncodeToString(digest[:])
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
		if _, ok := media["schema"]; ok {
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
	node, ok := value.(map[string]any)
	if !ok {
		return value
	}
	ref := stringValue(node["$ref"])
	prefix := "#/components/" + component + "/"
	if !strings.HasPrefix(ref, prefix) {
		return value
	}
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
		return resolved
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
	prefix := "#/components/parameters/"
	if !strings.HasPrefix(ref, prefix) {
		return nil, "", false
	}
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
			if existing.Name == candidate.Name && existing.In == candidate.In && existing.Pointer == candidate.Pointer {
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

func hasSecuritySchemesJSON(raw map[string]any) bool {
	components, ok := raw["components"].(map[string]any)
	if !ok {
		return false
	}
	schemes, ok := components["securitySchemes"].(map[string]any)
	return ok && len(schemes) > 0
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

func missingOperationalGuidance(lowerDocs string) bool {
	hasRetry := strings.Contains(lowerDocs, "retry")
	hasRateLimit := strings.Contains(lowerDocs, "rate limit") || strings.Contains(lowerDocs, "rate-limit") || strings.Contains(lowerDocs, "429")
	hasPagination := strings.Contains(lowerDocs, "pagination") || strings.Contains(lowerDocs, "page ")
	hasIdempotency := strings.Contains(lowerDocs, "idempotency") || strings.Contains(lowerDocs, "idempotent")
	return !(hasRetry && hasRateLimit && hasPagination && hasIdempotency)
}

func cleanMarkdownLinkTarget(target string) string {
	target = strings.TrimSpace(target)
	if space := strings.IndexByte(target, ' '); space >= 0 {
		target = target[:space]
	}
	return strings.Trim(target, `"'`)
}

func isExternalReference(target string) bool {
	lower := strings.ToLower(target)
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "mailto:") ||
		strings.HasPrefix(lower, "tel:") ||
		strings.HasPrefix(target, "#")
}

func isDocsFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".mdx", ".txt", ".rst":
		return true
	default:
		return false
	}
}

func projectRoot(configPath string) string {
	dir := filepath.Dir(configPath)
	if dir == "" {
		return "."
	}
	return dir
}

func resolveProjectPath(root, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(root, filepath.FromSlash(path)))
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

func resolveMarkdownLinkTarget(root, docPath, targetPath string) string {
	if strings.HasPrefix(targetPath, "/") {
		return filepath.Clean(filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(targetPath, "/"))))
	}
	return filepath.Clean(filepath.Join(filepath.Dir(docPath), filepath.FromSlash(targetPath)))
}

func yamlScalar(value string) string {
	if value == "" {
		return `""`
	}
	if strings.ContainsAny(value, "\n\r\t#:") || strings.HasPrefix(value, " ") || strings.HasSuffix(value, " ") {
		encoded, _ := json.Marshal(value)
		return string(encoded)
	}
	return value
}

func parseYAMLValue(value string) string {
	value = stripYAMLInlineComment(value)
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	return value
}

func stripYAMLInlineComment(value string) string {
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	for index, char := range value {
		switch {
		case escaped:
			escaped = false
		case char == '\\' && inDoubleQuote:
			escaped = true
		case char == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote
		case char == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
		case char == '#' && !inSingleQuote && !inDoubleQuote:
			if index > 0 && (value[index-1] == ' ' || value[index-1] == '\t') {
				return value[:index]
			}
		}
	}
	return value
}

func yamlInlineValueHasEntries(value string) bool {
	value = strings.TrimSpace(strings.Trim(value, `"'`))
	switch value {
	case "", "{}", "[]", "null", "~":
		return false
	default:
		return true
	}
}

func yamlLocalComponentRef(value, component string) (string, bool) {
	prefix := "#/components/" + component + "/"
	for _, candidate := range []string{parseYAMLValue(value), yamlInlineRefValue(value)} {
		if strings.HasPrefix(candidate, prefix) {
			return unescapeJSONPointer(strings.TrimPrefix(candidate, prefix)), true
		}
	}
	return "", false
}

func yamlInlineRefValue(value string) string {
	value = strings.TrimSpace(stripYAMLInlineComment(value))
	value = strings.Trim(value, "{} ")
	for _, entry := range strings.Split(value, ",") {
		key, refValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if !ok || strings.Trim(key, `"'`) != "$ref" {
			continue
		}
		return parseYAMLValue(strings.Trim(refValue, "{} "))
	}
	return ""
}

func parseYAMLParameterLine(parameters []openAPIParameter, current **openAPIParameter, trimmed, key, value, pointer string, lineNo int, parameterComponents map[string]openAPIParameter) []openAPIParameter {
	if strings.HasPrefix(trimmed, "- ") {
		parameterKey, parameterValue, ok := yamlKeyValue(strings.TrimSpace(strings.TrimPrefix(trimmed, "- ")))
		parameter := openAPIParameter{
			Pointer: pointer,
			Line:    lineNo,
		}
		if ok {
			if parameterKey == "$ref" {
				if refName, ok := yamlLocalComponentRef(parameterValue, "parameters"); ok {
					if referenced, ok := parameterComponents[refName]; ok {
						parameters = append(parameters, referenced)
					} else {
						parameters = append(parameters, openAPIParameter{
							Name:    parameterValue,
							Pointer: pointer,
							Line:    lineNo,
						})
					}
					*current = &parameters[len(parameters)-1]
					return parameters
				}
			}
			switch parameterKey {
			case "name":
				parameter.Name = parameterValue
			case "in":
				parameter.In = parameterValue
			case "description":
				parameter.Description = parameterValue
			}
		}
		parameters = append(parameters, parameter)
		*current = &parameters[len(parameters)-1]
		return parameters
	}
	if *current == nil {
		return parameters
	}
	switch key {
	case "name":
		(*current).Name = value
	case "in":
		(*current).In = value
	case "description":
		(*current).Description = value
	}
	return parameters
}

func yamlComponentParameters(data []byte) map[string]openAPIParameter {
	parameters := map[string]openAPIParameter{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	componentsIndent := -1
	parametersIndent := -1
	itemIndent := -1
	currentItem := ""
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
			parametersIndent = -1
			itemIndent = -1
			currentItem = ""
		}
		if parametersIndent >= 0 && indent <= parametersIndent && key != "parameters" {
			parametersIndent = -1
			itemIndent = -1
			currentItem = ""
		}
		if itemIndent >= 0 && indent <= itemIndent && key != currentItem {
			itemIndent = -1
			currentItem = ""
		}
		if key == "components" {
			componentsIndent = indent
			parametersIndent = -1
			itemIndent = -1
			currentItem = ""
			continue
		}
		if componentsIndent >= 0 && indent > componentsIndent && key == "parameters" {
			parametersIndent = indent
			itemIndent = -1
			currentItem = ""
			continue
		}
		if parametersIndent >= 0 && indent > parametersIndent && (itemIndent < 0 || indent <= itemIndent) {
			currentItem = key
			itemIndent = indent
			parameters[currentItem] = yamlInlineParameter(value, "#/components/parameters/"+escapeJSONPointer(currentItem))
			continue
		}
		if currentItem == "" || indent <= itemIndent {
			continue
		}
		parameter := parameters[currentItem]
		switch key {
		case "name":
			parameter.Name = value
		case "in":
			parameter.In = value
		case "description":
			parameter.Description = value
		}
		parameters[currentItem] = parameter
	}
	return parameters
}

func yamlInlineParameter(value, pointer string) openAPIParameter {
	parameter := openAPIParameter{Pointer: pointer}
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	for _, entry := range strings.Split(value, ",") {
		key, entryValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if !ok {
			continue
		}
		switch key {
		case "name":
			parameter.Name = entryValue
		case "in":
			parameter.In = entryValue
		case "description":
			parameter.Description = entryValue
		}
	}
	return parameter
}

func yamlComponentContentSchemaRefs(data []byte, group string) map[string]bool {
	refs := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	componentsIndent := -1
	groupIndent := -1
	itemIndent := -1
	contentIndent := -1
	mediaIndent := -1
	currentItem := ""
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
			groupIndent = -1
			itemIndent = -1
			contentIndent = -1
			mediaIndent = -1
			currentItem = ""
		}
		if groupIndent >= 0 && indent <= groupIndent && key != group {
			groupIndent = -1
			itemIndent = -1
			contentIndent = -1
			mediaIndent = -1
			currentItem = ""
		}
		if itemIndent >= 0 && indent <= itemIndent && key != currentItem {
			itemIndent = -1
			contentIndent = -1
			mediaIndent = -1
			currentItem = ""
		}
		if contentIndent >= 0 && indent <= contentIndent && key != "content" {
			contentIndent = -1
			mediaIndent = -1
		}
		if mediaIndent >= 0 && indent <= mediaIndent {
			mediaIndent = -1
		}
		if key == "components" {
			componentsIndent = indent
			continue
		}
		if componentsIndent >= 0 && indent > componentsIndent && key == group {
			groupIndent = indent
			continue
		}
		if groupIndent >= 0 && indent > groupIndent && (itemIndent < 0 || indent <= itemIndent) {
			currentItem = key
			itemIndent = indent
			contentIndent = -1
			mediaIndent = -1
			continue
		}
		if currentItem == "" {
			continue
		}
		if itemIndent >= 0 && indent > itemIndent && key == "content" {
			contentIndent = indent
			mediaIndent = -1
			continue
		}
		if contentIndent >= 0 && indent > contentIndent && strings.Contains(key, "/") {
			mediaIndent = indent
			if yamlInlineValueHasSchema(value) {
				refs[currentItem] = true
			}
			continue
		}
		if mediaIndent >= 0 && indent > mediaIndent && key == "schema" {
			refs[currentItem] = true
		}
	}
	return refs
}

func yamlInlineValueHasSchema(value string) bool {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return false
	}
	if yamlInlineMapHasKey(value, "schema") {
		return true
	}
	return strings.HasPrefix(value, "schema:") || strings.Contains(value, " schema:") || strings.Contains(value, `"schema"`)
}

func yamlInlineMapHasKey(value, expected string) bool {
	for _, entry := range strings.Split(value, ",") {
		key, _, ok := yamlKeyValue(strings.TrimSpace(entry))
		if ok && strings.Trim(key, `"'`) == expected {
			return true
		}
	}
	return false
}

func yamlKeyValue(trimmed string) (string, string, bool) {
	index := yamlMappingSeparator(trimmed)
	if index < 0 {
		return "", "", false
	}
	key := strings.Trim(strings.TrimSpace(trimmed[:index]), `"'`)
	value := parseYAMLValue(trimmed[index+1:])
	return key, value, true
}

func yamlMappingSeparator(trimmed string) int {
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	for index, char := range trimmed {
		switch {
		case escaped:
			escaped = false
		case char == '\\' && inDoubleQuote:
			escaped = true
		case char == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote
		case char == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
		case char == ':' && !inSingleQuote && !inDoubleQuote:
			if index+1 == len(trimmed) || trimmed[index+1] == ' ' || trimmed[index+1] == '\t' {
				return index
			}
		}
	}
	return -1
}

func leadingSpaces(value string) int {
	return len(value) - len(strings.TrimLeft(value, " "))
}

func is2xxStatusKey(key string) bool {
	key = strings.Trim(key, `"'`)
	return len(key) == 3 && key[0] == '2'
}

func escapeJSONPointer(value string) string {
	value = strings.ReplaceAll(value, "~", "~0")
	value = strings.ReplaceAll(value, "/", "~1")
	return value
}

func unescapeJSONPointer(value string) string {
	value = strings.ReplaceAll(value, "~1", "/")
	value = strings.ReplaceAll(value, "~0", "~")
	return value
}

type fileEntry struct {
	name string
}

func (f fileEntry) Name() string               { return f.name }
func (f fileEntry) IsDir() bool                { return false }
func (f fileEntry) Type() fs.FileMode          { return 0 }
func (f fileEntry) Info() (fs.FileInfo, error) { return nil, errors.New("file info unavailable") }
