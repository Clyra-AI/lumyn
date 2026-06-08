package exitcode

const (
	Success                         = 0
	InternalError                   = 1
	InvalidUsageOrInput             = 2
	SourceCompletenessFailure       = 3
	WorkflowContractValidationError = 4
	WorkflowVerificationFailure     = 5
	LiveEvalGateFailure             = 6
	CredentialAuthOrEnvironment     = 7
	DependencyProviderOrNetwork     = 8
	TraceCassetteOrReplayIntegrity  = 9
)

var Stable = map[int]string{
	Success:                         "success",
	InternalError:                   "general or internal error",
	InvalidUsageOrInput:             "invalid usage, invalid input, parse error, or local configuration error",
	SourceCompletenessFailure:       "source completeness failure in strict mode",
	WorkflowContractValidationError: "workflow contract validation failure",
	WorkflowVerificationFailure:     "workflow verification failure",
	LiveEvalGateFailure:             "live agent eval failed an explicitly configured regression or threshold gate",
	CredentialAuthOrEnvironment:     "credential, auth, or environment error",
	DependencyProviderOrNetwork:     "dependency, model provider, or network error",
	TraceCassetteOrReplayIntegrity:  "trace, cassette, or replay integrity failure",
}
