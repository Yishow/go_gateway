package api

import (
	"database/sql"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// TestNumericStateAuthoringSQLiteAcceptance exercises the public authoring,
// schema and Apply path before production intake and real SQLite readback.
// The same integer observations must yield latest 2 for state, but mean 1.5
// for gauge. No applied descriptor or computed result is fabricated here.
func TestNumericStateAuthoringSQLiteAcceptance(t *testing.T) {
	ctx := t.Context()
	f := newManagedDestinationFixture(t, "sqlite")
	sources := []meanSource{
		{name: "state", tagType: "int16", role: "state"},
		{name: "gauge", tagType: "int16", role: "gauge"},
	}
	seedMeanSources(t, f, sources)
	request := meanLongRequest(t, f, sources)
	authored := request["group"].(map[string]any)["members"].([]map[string]any)
	require.Equal(t, "state", authored[0]["value_role"])
	require.Equal(t, "gauge", authored[1]["value_role"])
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	groupID, table := saved["id"].(string), destinationOf(saved)["table_name"].(string)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, groupID, previewLong(t, f, groupID)).Status)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	snapshot, err := f.groups.ResolveAppliedAt(ctx, groupID, time.Now().Add(time.Minute))
	require.NoError(t, err)
	frozen := map[string]workspace.WriteGroupMember{}
	for _, member := range snapshot.Group.Members {
		frozen[member.TagID] = member
	}
	require.Len(t, frozen, 2)
	for _, source := range sources {
		member := frozen["tag-"+source.name]
		require.Equal(t, workspace.WriteGroupValueRole(source.role), member.ValueRole)
		require.Equal(t, schema.DataTypeInt16, member.SourceValueType)
	}
	require.Equal(t, measurement.ExactInt64, frozen["tag-state"].ResultValueType)
	require.Equal(t, measurement.ExactDecimal, frozen["tag-gauge"].ResultValueType)

	var clock atomic.Int64
	effective := snapshot.EffectiveAt
	clock.Store(effective.Add(3 * time.Second).UnixNano())
	pipe := newMeanPipeline(t, f, &clock, "node/numeric-state-acceptance", false)
	require.Equal(t, grouppipeline.StateActive, pipeStatus(pipe, groupID).State, "%+v", pipe.Status())
	for _, source := range sources {
		member := frozen["tag-"+source.name]
		for _, value := range []int{1, 2} {
			observed := effective.Add(time.Duration(value) * time.Second)
			require.NoError(t, pipe.AcceptSample(ctx, measurement.SampleEnvelope{
				SampleID:    source.name + observed.Format(time.RFC3339Nano),
				WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: member.PointID, TagID: member.TagID,
				SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
				ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond),
				TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
			}))
		}
	}
	clock.Store(effective.Add(16 * time.Second).UnixNano())
	target, err := sql.Open("sqlite", f.file)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	require.Eventually(t, func() bool {
		var count int
		return target.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table).Scan(&count) == nil && count == 2
	}, 15*time.Second, 50*time.Millisecond, "both sources reach the real SQLite destination")
	rows := readMeanRows(t, target, table)
	require.Len(t, rows, 2)
	state, gauge := rows["tag-state"], rows["tag-gauge"]
	require.Equal(t, effective.UTC().Format(meanTimeLayout), state.bucketStart)
	require.Equal(t, state.bucketStart, gauge.bucketStart, "both roles use the same window")
	require.Equal(t, [3]string{"ok", "int64", "2"}, [3]string{state.status, state.valueType, state.integer.String})
	require.True(t, state.integer.Valid)
	require.False(t, state.decimal.Valid || state.float.Valid, "numeric state is integer latest 2, never mean 1.5")
	require.Equal(t, [3]string{"state", "int16", "int64"}, [3]string{state.prov.Role, state.prov.SourceType, state.prov.ResultType})
	require.Nil(t, state.prov.Mean)
	require.Empty(t, state.prov.Counter)
	require.False(t, state.deltaValue.Valid || state.deltaFrom.Valid || state.deltaTo.Valid)
	require.Equal(t, effective.Add(2*time.Second).UTC().Format(meanTimeLayout), state.observed.String)
	require.Equal(t, [3]string{"ok", "decimal", "1.5"}, [3]string{gauge.status, gauge.valueType, gauge.decimal.String})
	require.True(t, gauge.decimal.Valid)
	require.False(t, gauge.integer.Valid || gauge.float.Valid)
	require.NotNil(t, gauge.prov.Mean)
	mean := gauge.prov.Mean
	require.Equal(t, [4]string{"mean", "gauge", "int16", "decimal"}, [4]string{mean.Reduction, mean.Role, mean.SourceType, mean.ResultType})
	require.Equal(t, "arithmetic-mean", mean.Method)
	require.Equal(t, "3", mean.Sum)
	require.Equal(t, int64(2), mean.ValidCount)
	require.Zero(t, mean.RejectedCount)
}
