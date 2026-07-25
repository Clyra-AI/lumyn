package exitcode

import "testing"

func TestStableExitCodesMatchPRD(t *testing.T) {
	expected := map[int]string{
		0: "success",
		1: "general or internal error",
		2: "invalid usage, invalid input, parse error, or local configuration error",
		3: "source completeness failure in strict mode",
		4: "workflow contract validation failure",
		5: "workflow verification failure",
		6: "live agent eval failed an explicitly configured regression or threshold gate",
		7: "credential, auth, or environment error",
		8: "dependency, model provider, or network error",
		9: "trace, cassette, or replay integrity failure",
	}
	if len(Stable) != len(expected) {
		t.Fatalf("expected %d stable exit codes, got %d", len(expected), len(Stable))
	}
	for code, meaning := range expected {
		if Stable[code] != meaning {
			t.Fatalf("exit code %d = %q, want %q", code, Stable[code], meaning)
		}
	}
}

func TestLiveEvalGateExitCodeRemainsReserved(t *testing.T) {
	if LiveEvalGateFailure != 6 {
		t.Fatalf("LiveEvalGateFailure = %d, want reserved code 6", LiveEvalGateFailure)
	}
}
