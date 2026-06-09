package result

type CommandResult struct {
	ObjectType      string         `json:"object_type"`
	SchemaVersion   string         `json:"schema_version"`
	Metadata        map[string]any `json:"metadata"`
	Command         string         `json:"command"`
	Status          string         `json:"status"`
	Mode            string         `json:"mode"`
	Warnings        []string       `json:"warnings"`
	Errors          []CommandError `json:"errors"`
	Artifacts       []ArtifactRef  `json:"artifacts"`
	DurationMS      int64          `json:"duration_ms"`
	RedactionStatus string         `json:"redaction_status"`
}

type CommandError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ArtifactRef struct {
	Path string `json:"path"`
	Type string `json:"type"`
}
