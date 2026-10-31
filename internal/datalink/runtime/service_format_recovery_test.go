package runtime

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestPreparedFormatRecoveryConsumesOriginalSlotProof(t *testing.T) {
	for _, race := range []string{"normal", "foreign_same_revision", "absent_invalidated", "bad_profile"} {
		t.Run(race, func(t *testing.T) {
			f := newAcquisitionInstallFixture(t, false)
			p := f.projection.Points[0]
			p.FormatMode = schema.PointFormatInherit
			p.FormatVersion = 1
			p.RevisionID = "p-rev"
			p.DataType = schema.DataTypeUint32
			f.projection.RuleLinks[0].Address = p.Address
			require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
			profile, err := f.service.ResolvePointReadProfile(t.Context(), p)
			require.NoError(t, err)
			if race == "absent_invalidated" {
				require.True(t, f.service.SuspendPointRevision(p.ID, p.RevisionID))
			}
			proof, ok := f.service.CaptureOrdinaryPoint(p.ID)
			require.True(t, ok)
			foreign := uint64(0)
			switch race {
			case "foreign_same_revision":
				foreign = f.service.InstallPointForMutation(p)
				require.NotZero(t, foreign)
			case "absent_invalidated":
				require.True(t, f.service.SuspendPointRevision(p.ID, p.RevisionID))
			case "bad_profile":
				profile = connector.ReadProfile{}
			}
			installed, ok := f.service.InstallPreparedFormatRecovery(p, proof, profile)
			if race == "normal" {
				require.True(t, ok)
				require.NotZero(t, installed)
				captured, valid := f.service.CapturePointReadProfile(p.ID, installed)
				require.True(t, valid)
				require.Equal(t, profile, captured)
				require.Empty(t, f.service.scheduler.PollNowContext(t.Context(), []string{p.ID}))
				require.True(t, f.service.ActivatePreparedFormatRecovery(p.ID, installed, profile))
				require.Len(t, f.service.scheduler.PollNowContext(t.Context(), []string{p.ID}), 1)
				require.False(t, f.service.ActivatePreparedFormatRecovery(p.ID, installed, profile), "activation itself is not replay; successful replay must be read-only at the owning operation")
			} else {
				require.False(t, ok)
				require.Zero(t, installed)
				if foreign != 0 {
					require.Equal(t, foreign, f.service.CapturePointLease(p.ID, p.RevisionID))
				}
				if race == "absent_invalidated" {
					require.Zero(t, f.service.CapturePointLease(p.ID, p.RevisionID))
				}
			}
		})
	}
}
