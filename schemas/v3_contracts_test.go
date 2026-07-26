package schemas_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var v3ArtifactSchemas = []string{
	"provider-change-event",
	"migration-pack",
	"consumer-installation",
	"event-authorization",
	"consumer-execution-manifest",
	"managed-credential-grant",
	"integration-graph",
	"impact-report",
	"migration-plan",
	"candidate-manifest",
	"agent-runner-manifest",
	"agent-runner-conformance-result",
	"agent-attempt",
	"migration-verification",
	"export-result",
	"campaign-summary",
	"provider-status-projection",
	"remediation-outcome",
}

func TestV3ArtifactSchemasCompileAndValidateFixtures(t *testing.T) {
	root := repoRoot(t)
	for _, family := range v3ArtifactSchemas {
		family := family
		t.Run(family, func(t *testing.T) {
			schemaPath := filepath.Join(root, "schemas", family+".schema.json")
			schemaDocument := readJSONFixture(t, schemaPath).(map[string]any)
			if schemaDocument["$schema"] != "https://json-schema.org/draft/2020-12/schema" {
				t.Fatalf("schema draft = %v, want draft 2020-12", schemaDocument["$schema"])
			}
			if schemaDocument["additionalProperties"] != false {
				t.Fatal("top-level contract must reject unknown control fields")
			}
			properties := schemaDocument["properties"].(map[string]any)
			objectType := properties["object_type"].(map[string]any)["const"].(string)
			expectedID := "https://schemas.lumyn.local/" + objectType + ".schema.json"
			if schemaDocument["$id"] != expectedID {
				t.Fatalf("schema id = %v, want %s", schemaDocument["$id"], expectedID)
			}
			schema, err := jsonschema.Compile(schemaPath)
			if err != nil {
				t.Fatalf("compile schema: %v", err)
			}

			fixtureRoot := filepath.Join(root, "tests", "fixtures", "contracts", family)
			validPaths, err := filepath.Glob(filepath.Join(fixtureRoot, "valid*.json"))
			if err != nil || len(validPaths) == 0 {
				t.Fatalf("discover valid fixtures: paths=%v err=%v", validPaths, err)
			}
			for _, path := range validPaths {
				valid := readJSONFixture(t, path)
				validObject := valid.(map[string]any)
				if validObject["object_type"] != objectType || validObject["schema_version"] != "1.0" {
					t.Fatalf("fixture %s identity/version does not match schema", filepath.Base(path))
				}
				if err := schema.Validate(valid); err != nil {
					t.Fatalf("valid fixture %s rejected: %v", filepath.Base(path), err)
				}
			}

			invalidPaths, err := filepath.Glob(filepath.Join(fixtureRoot, "invalid*.json"))
			if err != nil || len(invalidPaths) == 0 {
				t.Fatalf("discover invalid fixtures: paths=%v err=%v", invalidPaths, err)
			}
			for _, path := range invalidPaths {
				invalid := readJSONFixture(t, path)
				if err := schema.Validate(invalid); err == nil {
					t.Fatalf("invalid trust-boundary fixture %s unexpectedly validated", filepath.Base(path))
				}
			}
		})
	}
}

func readJSONFixture(t *testing.T, path string) any {
	t.Helper()
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	var value any
	if err := json.Unmarshal(payload, &value); err != nil {
		t.Fatalf("decode fixture %s: %v", path, err)
	}
	return value
}
