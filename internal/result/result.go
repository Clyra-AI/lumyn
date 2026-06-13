package result

type CommandResult struct {
	ObjectType           string           `json:"object_type"`
	SchemaVersion        string           `json:"schema_version"`
	Metadata             map[string]any   `json:"metadata"`
	Command              string           `json:"command"`
	Status               string           `json:"status"`
	Mode                 string           `json:"mode"`
	Warnings             []string         `json:"warnings"`
	Errors               []CommandError   `json:"errors"`
	Artifacts            []ArtifactRef    `json:"artifacts"`
	DurationMS           int64            `json:"duration_ms"`
	RedactionStatus      string           `json:"redaction_status"`
	FindingKind          string           `json:"finding_kind"`
	ProofStrength        string           `json:"proof_strength"`
	ActionBoundaryStatus string           `json:"action_boundary_status"`
	SecurityRelevance    string           `json:"security_relevance"`
	FixTarget            string           `json:"fix_target"`
	SurfaceFingerprint   string           `json:"surface_fingerprint"`
	EvalMode             string           `json:"eval_mode"`
	ProviderMetadata     ProviderMetadata `json:"provider_metadata"`
	CorpusEligible       bool             `json:"corpus_eligible"`
}

type CommandError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ArtifactRef struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type ProviderMetadata struct {
	Applicable    bool   `json:"applicable"`
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ModelSnapshot string `json:"model_snapshot,omitempty"`
	BaseURL       string `json:"base_url,omitempty"`
}
