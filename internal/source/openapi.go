package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

var (
	httpMethods    = map[string]bool{"get": true, "put": true, "post": true, "delete": true, "patch": true, "head": true, "options": true, "trace": true}
	mutatingMethod = map[string]bool{"put": true, "post": true, "delete": true, "patch": true}
)

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
