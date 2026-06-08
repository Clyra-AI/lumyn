package exitcode

import "testing"

func TestStableExitCodesMatchPRD(t *testing.T) {
	for code := Success; code <= TraceCassetteOrReplayIntegrity; code++ {
		if Stable[code] == "" {
			t.Fatalf("missing stable meaning for exit code %d", code)
		}
	}
	if len(Stable) != 10 {
		t.Fatalf("expected 10 stable exit codes, got %d", len(Stable))
	}
}
