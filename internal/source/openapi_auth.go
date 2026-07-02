package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

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
