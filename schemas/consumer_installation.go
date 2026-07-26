// Package schemas exposes the executable JSON contracts that production
// decoders must apply before deriving product authority.
package schemas

import (
	_ "embed"
	"fmt"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed consumer-installation.schema.json
var consumerInstallationSchemaSource string

//go:embed event-authorization.schema.json
var eventAuthorizationSchemaSource string

var (
	consumerInstallationSchemaOnce sync.Once
	consumerInstallationSchema     *jsonschema.Schema
	consumerInstallationSchemaErr  error
	eventAuthorizationSchemaOnce   sync.Once
	eventAuthorizationSchema       *jsonschema.Schema
	eventAuthorizationSchemaErr    error
)

// ValidateConsumerInstallation validates a decoded JSON value against the
// canonical embedded Consumer Installation contract. Embedding keeps runtime
// validation bound to the same schema that contract tests exercise, without a
// working-directory-dependent schema lookup.
func ValidateConsumerInstallation(value any) error {
	consumerInstallationSchemaOnce.Do(func() {
		consumerInstallationSchema, consumerInstallationSchemaErr = jsonschema.CompileString(
			"consumer-installation.schema.json",
			consumerInstallationSchemaSource,
		)
	})
	if consumerInstallationSchemaErr != nil {
		return fmt.Errorf("compile consumer installation schema: %w", consumerInstallationSchemaErr)
	}
	if err := consumerInstallationSchema.Validate(value); err != nil {
		return fmt.Errorf("validate consumer installation schema: %w", err)
	}
	return nil
}

// ValidateEventAuthorization validates a decoded JSON value against the
// canonical embedded Event Authorization contract.
func ValidateEventAuthorization(value any) error {
	eventAuthorizationSchemaOnce.Do(func() {
		eventAuthorizationSchema, eventAuthorizationSchemaErr = jsonschema.CompileString(
			"event-authorization.schema.json",
			eventAuthorizationSchemaSource,
		)
	})
	if eventAuthorizationSchemaErr != nil {
		return fmt.Errorf("compile event authorization schema: %w", eventAuthorizationSchemaErr)
	}
	if err := eventAuthorizationSchema.Validate(value); err != nil {
		return fmt.Errorf("validate event authorization schema: %w", err)
	}
	return nil
}
