package source

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// ArtifactKind is a non-executable provider-evidence class accepted by M3.
type ArtifactKind string

const (
	ArtifactOpenAPI        ArtifactKind = "openapi"
	ArtifactDocumentation  ArtifactKind = "documentation"
	ArtifactSDKTypes       ArtifactKind = "sdk_types"
	ArtifactMigrationGuide ArtifactKind = "migration_guidance"
)

// VersionRole identifies which side of a migration an artifact describes.
type VersionRole string

const (
	VersionSource   VersionRole = "source"
	VersionTarget   VersionRole = "target"
	VersionGuidance VersionRole = "guidance"
)

// PinnedArtifact carries immutable source bytes and their provenance. Data is
// intentionally untrusted and is never executed by this package.
type PinnedArtifact struct {
	ID          string
	Kind        ArtifactKind
	VersionRole VersionRole
	Version     string
	Location    string
	License     string
	MediaType   string
	Digest      string
	RetrievedAt time.Time
	FreshUntil  time.Time
	Data        []byte
}

// PinArtifact validates provenance and returns a defensive copy with a
// computed digest. expectedDigest may be empty on first pin; when present it
// must match the retrieved bytes.
func PinArtifact(artifact PinnedArtifact, expectedDigest string, now time.Time) (PinnedArtifact, error) {
	if strings.TrimSpace(artifact.ID) == "" || strings.TrimSpace(artifact.Version) == "" {
		return PinnedArtifact{}, errors.New("source artifact id and version are required")
	}
	if !validArtifactKind(artifact.Kind) || !validVersionRole(artifact.VersionRole) {
		return PinnedArtifact{}, errors.New("source artifact kind and version role must be supported")
	}
	parsed, err := url.Parse(artifact.Location)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil ||
		parsed.RawQuery != "" || parsed.Fragment != "" {
		return PinnedArtifact{}, errors.New("source artifact location must be a credential-free pinned HTTPS URL")
	}
	if strings.TrimSpace(artifact.License) == "" {
		return PinnedArtifact{}, errors.New("source artifact license posture is required")
	}
	if len(artifact.Data) == 0 || now.IsZero() || artifact.RetrievedAt.IsZero() || artifact.RetrievedAt.After(now) {
		return PinnedArtifact{}, errors.New("source artifact bytes and valid retrieval time are required")
	}
	if !artifact.FreshUntil.IsZero() && !artifact.FreshUntil.After(now) {
		return PinnedArtifact{}, errors.New("source artifact is stale")
	}
	digest := DigestBytes(artifact.Data)
	if expectedDigest != "" && expectedDigest != digest {
		return PinnedArtifact{}, fmt.Errorf("source artifact digest mismatch: got %s", digest)
	}
	artifact.Data = append([]byte(nil), artifact.Data...)
	artifact.Digest = digest
	return artifact, nil
}

func validArtifactKind(kind ArtifactKind) bool {
	switch kind {
	case ArtifactOpenAPI, ArtifactDocumentation, ArtifactSDKTypes, ArtifactMigrationGuide:
		return true
	default:
		return false
	}
}

func validVersionRole(role VersionRole) bool {
	switch role {
	case VersionSource, VersionTarget, VersionGuidance:
		return true
	default:
		return false
	}
}

// DigestBytes returns the repository-wide SHA-256 wire representation.
func DigestBytes(data []byte) string {
	digest := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(digest[:])
}

// OpenAPISnapshot is the deterministic semantic subset used to compare source
// and target descriptions. Raw formatting and JSON/YAML syntax are excluded.
type OpenAPISnapshot struct {
	Operations []OpenAPIOperation `json:"operations"`
}

type OpenAPIOperation struct {
	Method            string   `json:"method"`
	Path              string   `json:"path"`
	OperationID       string   `json:"operation_id"`
	Parameters        []string `json:"parameters"`
	HasRequestSchema  bool     `json:"has_request_schema"`
	HasResponseSchema bool     `json:"has_response_schema"`
	Deprecated        bool     `json:"deprecated"`
}

// SnapshotOpenAPI parses the existing supported OpenAPI JSON/YAML subset and
// returns one format-independent, sorted semantic view.
func SnapshotOpenAPI(data []byte) (OpenAPISnapshot, error) {
	operations, _, err := parseOpenAPI(data)
	if err != nil {
		return OpenAPISnapshot{}, err
	}
	result := OpenAPISnapshot{Operations: make([]OpenAPIOperation, 0, len(operations))}
	for _, operation := range operations {
		parameters := make([]string, 0, len(operation.Parameters))
		for _, parameter := range operation.Parameters {
			parameters = append(parameters, parameter.In+":"+parameter.Name)
		}
		sort.Strings(parameters)
		result.Operations = append(result.Operations, OpenAPIOperation{
			Method: operation.Method, Path: operation.Path, OperationID: operation.OperationID,
			Parameters: parameters, HasRequestSchema: operation.HasRequestSchema,
			HasResponseSchema: operation.HasResponseSchema, Deprecated: operation.Deprecated,
		})
	}
	sort.Slice(result.Operations, func(i, j int) bool {
		left, right := result.Operations[i], result.Operations[j]
		return left.Method+"|"+left.Path+"|"+left.OperationID < right.Method+"|"+right.Path+"|"+right.OperationID
	})
	return result, nil
}
