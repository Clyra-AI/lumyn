package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Clyra-AI/lumyn/internal/config"
)

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
