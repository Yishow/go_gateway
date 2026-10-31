package runtime

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestPointRuntimeInstallationRequiresKnownActiveScope(t *testing.T) {
	for _, state := range []schema.DeviceStatus{schema.DeviceStatusDraft, schema.DeviceStatusActive, schema.DeviceStatusDisabled, "", "unknown"} {
		t.Run(string(state), func(t *testing.T) {
			f := newAcquisitionInstallFixture(t, false)
			p := f.projection.Points[0]
			p.FormatMode = schema.PointFormatInherit
			p.FormatVersion = 1
			p.RevisionID = "p"
			p.DataType = schema.DataTypeUint32
			f.projection.RuleLinks[0].Address = p.Address
			f.projection.Devices[0].Status = state
			required, err := f.service.PointRuntimeInstallationRequired(t.Context(), p)
			if state == "" || state == "unknown" {
				require.Error(t, err)
				require.False(t, required)
				return
			}
			require.NoError(t, err)
			require.Equal(t, state == schema.DeviceStatusActive, required)
			stale := *p
			stale.RevisionID = "foreign"
			_, err = f.service.PointRuntimeInstallationRequired(t.Context(), &stale)
			require.Error(t, err, "draft is not permission to skip canonical scope validation")
			f.projection.DeviceIDs = nil
			_, err = f.service.PointRuntimeInstallationRequired(t.Context(), p)
			require.Error(t, err, "unknown workspace ownership remains fail closed")
		})
	}
}
