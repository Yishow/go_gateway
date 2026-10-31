package groupdelivery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedCalculatedValueNeverAppearsAsLegacyRawEvidence(t *testing.T) {
	store, _, _ := newTestStore(t)
	encoded, _ := (longSpec{source: sourceA, profile: "mixed", status: "ok", slots: map[string]any{"value_decimal": "8"}}).row(t)
	encoded.DueAt = at(15)
	payload, digest, err := encodeRowPayload(encoded)
	require.NoError(t, err)
	mixedScope := scope
	mixedScope.ProfileKind = "mixed"
	row := sourceRow{workspaceID: testKey.WorkspaceID, groupID: testKey.GroupID, groupRevision: testKey.GroupRevision, connectorID: scope.ConnectorID, connectorRevision: scope.ConnectorRevision, recordID: encoded.RecordID, effectKey: encoded.EffectKey, entityKey: encoded.EntityKey, bucketStart: timestamp(t0), state: StatePending, payload: string(payload), digest: digest}
	result, err := store.verifySource(t.Context(), mixedScope, sourceA, row)
	require.NoError(t, err)
	require.False(t, result.Available, "mixed selected delta is not raw input evidence")
	require.Nil(t, result.VerifiedSourceRecord)
	require.Equal(t, SourceCauseCalculationUnsupported, result.Cause)
}
