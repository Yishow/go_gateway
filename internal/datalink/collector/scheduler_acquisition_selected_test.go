package collector

import (
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionSelectedGeometryCannotPromoteRetainedUnknownScope(t *testing.T) {
	for _, kind := range []string{"matching", "workspace-changed", "source-changed", "unknown"} {
		t.Run(kind, func(t *testing.T) {
			s := seedAcquisitionGeometry(t)
			s.acquisitionGeometry["d"].publication.workspaceID = "w"
			s.acquisitionGeometry["d"].publication.sourceRevision = "head"
			if kind == "unknown" {
				s.acquisitionGeometry["d"].publication.revision = ""
			}
			definition := AcquisitionGeometryDefinition{WorkspaceID: "w", SourceRevision: "head", Revision: "semantic"}
			if kind == "workspace-changed" {
				definition.WorkspaceID = "foreign"
			}
			if kind == "source-changed" {
				definition.SourceRevision = "new-head"
			}
			token := s.BeginAcquisitionGeometryMutation("d")
			update := s.PrepareSelectedAcquisitionGeometry(token, schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`}, definition, map[string]bool{"new": true})
			require.NotNil(t, update)
			require.Equal(t, kind == "matching", update.candidate.revision != "")
			require.Equal(t, token.token.before.slots, update.candidate.slots, "retained slots are original receipts, not recaptured")
			s.FinishAcquisitionGeometryMutation(token)
		})
	}
}
