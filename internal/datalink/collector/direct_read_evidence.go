package collector

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// DirectReadEvidence is ephemeral execution evidence for a no-slot preview.
// It is never persisted or accepted from HTTP input. Historical values are not
// re-decoded: a changed current identity only causes rejection.
type DirectReadEvidence struct {
	ConnectionRevision string
	Protocol           schema.ProtocolType
	Point              schema.Point
	Profile            connector.ReadProfile
	Mutation           PointMutationLease
	Scope              string
	Guarded            bool
}
