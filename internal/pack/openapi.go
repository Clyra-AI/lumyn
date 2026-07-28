package pack

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/Clyra-AI/lumyn/internal/source"
)

type OperationMapping struct {
	ChangeID          string
	SourceOperationID string
	TargetOperationID string
	Intent            string
	Applicability     string
	MappingMode       string
	DeclarativeSteps  []string
}

// DeriveOpenAPIChanges compares format-independent source and target
// snapshots. Stable operation IDs are compared automatically; renames and
// removals require an explicit, evidence-reviewed mapping or remain blocked.
func DeriveOpenAPIChanges(
	sourceSnapshot, targetSnapshot source.OpenAPISnapshot,
	sourceEvidenceID, targetEvidenceID string,
	mappings []OperationMapping,
) ([]Declaration, []Ambiguity, error) {
	sourceOperations, err := operationIndex(sourceSnapshot)
	if err != nil {
		return nil, nil, fmt.Errorf("source OpenAPI: %w", err)
	}
	targetOperations, err := operationIndex(targetSnapshot)
	if err != nil {
		return nil, nil, fmt.Errorf("target OpenAPI: %w", err)
	}
	mappingBySource := map[string]OperationMapping{}
	for _, mapping := range mappings {
		if mapping.SourceOperationID == "" || mapping.TargetOperationID == "" {
			return nil, nil, errors.New("OpenAPI mapping requires source and target operation ids")
		}
		if _, exists := mappingBySource[mapping.SourceOperationID]; exists {
			return nil, nil, fmt.Errorf("duplicate OpenAPI mapping for %s", mapping.SourceOperationID)
		}
		if _, exists := sourceOperations[mapping.SourceOperationID]; !exists {
			return nil, nil, fmt.Errorf("mapped source operation %s is absent", mapping.SourceOperationID)
		}
		if _, exists := targetOperations[mapping.TargetOperationID]; !exists {
			return nil, nil, fmt.Errorf("mapped target operation %s is absent", mapping.TargetOperationID)
		}
		mappingBySource[mapping.SourceOperationID] = mapping
	}

	declarations := []Declaration{}
	ambiguities := []Ambiguity{}
	sourceIDs := sortedOperationIDs(sourceOperations)
	for _, operationID := range sourceIDs {
		before := sourceOperations[operationID]
		after, exists := targetOperations[operationID]
		mapping, mapped := mappingBySource[operationID]
		if mapped {
			after = targetOperations[mapping.TargetOperationID]
			exists = true
		}
		if !exists {
			ambiguities = append(ambiguities, Ambiguity{
				AmbiguityID:        "ambiguity.openapi." + stableFragment(operationID),
				Description:        "Source operation " + operationID + " is absent from the target and has no reviewed mapping.",
				RequiredResolution: "provider_clarification",
				AffectedChangeIDs:  []string{},
			})
			continue
		}
		if !mapped && (before.HasRequestSchema || before.HasResponseSchema ||
			after.HasRequestSchema || after.HasResponseSchema) {
			ambiguities = append(ambiguities, Ambiguity{
				AmbiguityID:        "ambiguity.openapi_schema." + stableFragment(operationID),
				Description:        "Operation " + operationID + " carries request or response schemas whose property, required-field, and type semantics are outside the current snapshot.",
				RequiredResolution: "provider_clarification",
				AffectedChangeIDs:  []string{},
			})
			continue
		}
		if !mapped && reflect.DeepEqual(before, after) {
			continue
		}
		changeID := mapping.ChangeID
		intent := mapping.Intent
		applicability := mapping.Applicability
		mode := mapping.MappingMode
		steps := mapping.DeclarativeSteps
		if !mapped {
			changeID = "change.openapi." + stableFragment(operationID)
			intent = "Align the operation with the target OpenAPI description."
			applicability = "Calls bound to OpenAPI operation " + operationID + "."
			mode = "repository_specific_reasoning"
			steps = []string{"Preserve consumer expressions while adapting the call to the target operation signature."}
		}
		kind := "signature_change"
		if before.OperationID != after.OperationID {
			kind = "rename"
		} else if before.Method != after.Method || before.Path != after.Path {
			kind = "endpoint_change"
		}
		targetSymbol := operationSymbol(after)
		declarations = append(declarations, Declaration{SemanticChange: SemanticChange{
			ChangeID: changeID, Kind: kind, SourceSymbol: operationSymbol(before),
			TargetSymbol: &targetSymbol, Intent: intent, Applicability: applicability,
			Mapping:           Mapping{Mode: mode, DeclarativeSteps: append([]string(nil), steps...)},
			SourceEvidenceIDs: []string{sourceEvidenceID},
			TargetEvidenceIDs: []string{targetEvidenceID},
		}})
	}
	sort.Slice(declarations, func(i, j int) bool {
		return declarations[i].ChangeID < declarations[j].ChangeID
	})
	sort.Slice(ambiguities, func(i, j int) bool {
		return ambiguities[i].AmbiguityID < ambiguities[j].AmbiguityID
	})
	return declarations, ambiguities, nil
}

func operationIndex(snapshot source.OpenAPISnapshot) (map[string]source.OpenAPIOperation, error) {
	result := make(map[string]source.OpenAPIOperation, len(snapshot.Operations))
	for _, operation := range snapshot.Operations {
		if operation.OperationID == "" {
			return nil, errors.New("every compared operation requires operationId")
		}
		if _, exists := result[operation.OperationID]; exists {
			return nil, fmt.Errorf("duplicate operationId %s", operation.OperationID)
		}
		result[operation.OperationID] = operation
	}
	return result, nil
}

func sortedOperationIDs(operations map[string]source.OpenAPIOperation) []string {
	ids := make([]string, 0, len(operations))
	for id := range operations {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func operationSymbol(operation source.OpenAPIOperation) string {
	return strings.ToUpper(operation.Method) + " " + operation.Path + " (" + operation.OperationID + ")"
}

var unstableFragment = regexp.MustCompile(`[^a-zA-Z0-9._:-]+`)

func stableFragment(value string) string {
	value = unstableFragment.ReplaceAllString(value, "_")
	value = strings.Trim(value, "_")
	if value == "" {
		return "unknown"
	}
	return value
}
