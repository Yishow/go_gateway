package runtime

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionProfileFailureCannotAdoptForeignCleanupLease(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	original := *f.projection.Points[0]
	ctx, finish := f.service.beginFullAcquisitionProjection(t.Context(), f.projection)
	defer finish()
	foreign := f.service.scheduler.AddPointWithObservationLease(&original)
	candidate := original
	candidate.FormatMode = schema.PointFormatOverride
	candidate.FormatOverride = "invalid-format"
	candidate.FormatVersion = 1
	f.projection.Points[0] = &candidate
	_, err := f.service.prepareProjectionReadProfiles(ctx, f.projection)
	require.Error(t, err)
	require.Equal(t, foreign, f.service.scheduler.CapturePointLease(original.ID, original.RevisionID), "failed original candidate must not adopt a later cleanup lease")
}
