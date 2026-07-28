package pack

import (
	"testing"

	"github.com/Clyra-AI/lumyn/internal/source"
)

func TestDeriveOpenAPIChangesRequiresReviewedMappingForRemovedOperation(t *testing.T) {
	before := source.OpenAPISnapshot{Operations: []source.OpenAPIOperation{{
		Method: "post", Path: "/charges", OperationID: "createCharge",
		Parameters: []string{"header:account"}, HasRequestSchema: true, HasResponseSchema: true,
	}}}
	after := source.OpenAPISnapshot{Operations: []source.OpenAPIOperation{{
		Method: "post", Path: "/payment_intents", OperationID: "createPaymentIntent",
		Parameters: []string{"header:account"}, HasRequestSchema: true, HasResponseSchema: true,
	}}}
	declarations, ambiguities, err := DeriveOpenAPIChanges(before, after, "evidence.openapi_v1", "evidence.openapi_v2", nil)
	if err != nil {
		t.Fatalf("derive unmapped change: %v", err)
	}
	if len(declarations) != 0 || len(ambiguities) != 1 {
		t.Fatalf("removed operation must remain ambiguous: declarations=%#v ambiguities=%#v", declarations, ambiguities)
	}

	declarations, ambiguities, err = DeriveOpenAPIChanges(before, after, "evidence.openapi_v1", "evidence.openapi_v2", []OperationMapping{{
		ChangeID: "change.create_payment_intent", SourceOperationID: "createCharge",
		TargetOperationID: "createPaymentIntent", Intent: "Use Payment Intents for charge creation.",
		Applicability: "Direct createCharge calls.", MappingMode: "repository_specific_reasoning",
		DeclarativeSteps: []string{"Preserve the amount and currency expressions in the target request."},
	}})
	if err != nil {
		t.Fatalf("derive reviewed mapping: %v", err)
	}
	if len(ambiguities) != 0 || len(declarations) != 1 ||
		declarations[0].Kind != "rename" ||
		declarations[0].SourceEvidenceIDs[0] != "evidence.openapi_v1" ||
		declarations[0].TargetEvidenceIDs[0] != "evidence.openapi_v2" {
		t.Fatalf("unexpected reviewed mapping: declarations=%#v ambiguities=%#v", declarations, ambiguities)
	}
}

func TestDeriveOpenAPIChangesDetectsStableOperationSignatureChange(t *testing.T) {
	before := source.OpenAPISnapshot{Operations: []source.OpenAPIOperation{{
		Method: "post", Path: "/charges", OperationID: "createCharge",
		Parameters: []string{"header:account"},
	}}}
	after := source.OpenAPISnapshot{Operations: []source.OpenAPIOperation{{
		Method: "post", Path: "/charges", OperationID: "createCharge",
		Parameters: []string{"header:account", "query:idempotency_key"},
	}}}
	declarations, ambiguities, err := DeriveOpenAPIChanges(before, after, "evidence.openapi_v1", "evidence.openapi_v2", nil)
	if err != nil {
		t.Fatalf("derive signature change: %v", err)
	}
	if len(ambiguities) != 0 || len(declarations) != 1 || declarations[0].Kind != "signature_change" {
		t.Fatalf("unexpected signature derivation: declarations=%#v ambiguities=%#v", declarations, ambiguities)
	}
}

func TestDeriveOpenAPIChangesFailsClosedForUnmodeledSchemaDetails(t *testing.T) {
	before := source.OpenAPISnapshot{Operations: []source.OpenAPIOperation{{
		Method: "post", Path: "/charges", OperationID: "createCharge",
		HasRequestSchema: true, HasResponseSchema: true,
	}}}
	after := source.OpenAPISnapshot{Operations: []source.OpenAPIOperation{{
		Method: "post", Path: "/charges", OperationID: "createCharge",
		HasRequestSchema: true, HasResponseSchema: true,
	}}}
	declarations, ambiguities, err := DeriveOpenAPIChanges(before, after, "evidence.openapi_v1", "evidence.openapi_v2", nil)
	if err != nil {
		t.Fatalf("derive schema-bearing operation: %v", err)
	}
	if len(declarations) != 0 || len(ambiguities) != 1 ||
		ambiguities[0].RequiredResolution != "provider_clarification" {
		t.Fatalf("unmodeled schema semantics must fail closed: declarations=%#v ambiguities=%#v", declarations, ambiguities)
	}
}
