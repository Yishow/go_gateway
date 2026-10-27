package grouptestwrite

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// sqlServerMixedGaugeCounterStateFixture creates a verified SQL Server test-write fixture
// configured with an instantaneous profile under the mean reduction, containing
// a gauge member, a counter_raw member, and a state member.
func sqlServerMixedGaugeCounterStateFixture(t *testing.T) (*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout, map[string]schema.DataType) {
	t.Helper()
	group, layout, _ := sqlServerTestWriteFixture(t)

	// Instantaneous profile with mean reduction.
	group.RecordingProfile = &workspace.WriteGroupRecordingProfile{
		Kind:            workspace.WriteGroupProfileInstantaneous,
		Interval:        15,
		Unit:            workspace.WriteGroupIntervalUnitSecond,
		SampleReduction: workspace.WriteGroupSampleReductionMean,
	}

	measurementID := "meas-counter-1"
	group.Members = []workspace.WriteGroupMember{
		{
			DeviceID:        "device-1",
			PointID:         "point-1",
			TagID:           "tag-gauge",
			TargetColumn:    "value_decimal",
			ValueRole:       workspace.WriteGroupValueRoleGauge,
			SourceValueType: schema.DataTypeInt64,
			ResultValueType: measurement.ExactDecimal,
		},
		{
			DeviceID:        "device-1",
			PointID:         "point-2",
			TagID:           "tag-counter",
			TargetColumn:    "value_int",
			ValueRole:       workspace.WriteGroupValueRoleCounterRaw,
			SourceValueType: schema.DataTypeInt64,
			ResultValueType: measurement.ExactInt64,
			MeasurementID:   &measurementID,
		},
		{
			DeviceID:        "device-1",
			PointID:         "point-3",
			TagID:           "tag-state",
			TargetColumn:    "value_bool",
			ValueRole:       workspace.WriteGroupValueRoleState,
			SourceValueType: schema.DataTypeBool,
			ResultValueType: measurement.ExactBool,
		},
	}

	tagTypes := map[string]schema.DataType{
		"tag-gauge":   schema.DataTypeInt64,
		"tag-counter": schema.DataTypeInt64,
		"tag-state":   schema.DataTypeBool,
	}
	layout.TagTypes = tagTypes

	return group, layout, tagTypes
}

// TestSQLServerTestWritePreviewMixedGaugeAndCounterRaw verifies that when a write group
// configured with instantaneous mean reduction contains gauge, counter_raw, and state members,
// preview generation succeeds without being blocked by ReasonSQLValueBlocked (422 test-value-blocked).
// Only ValueRoleGauge produces mean reduction; counter_raw and state preserve their latest exact type.
func TestSQLServerTestWritePreviewMixedGaugeAndCounterRaw(t *testing.T) {
	group, layout, tagTypes := sqlServerMixedGaugeCounterStateFixture(t)
	require.NoError(t, workspace.ValidateSQLServerRuntimeLayout(group, layout), "fixture must satisfy SQL Server runtime layout contract")

	groups := &sqlServerTestWriteGroups{
		fakeGroups:                     &fakeGroups{group: group},
		layout:                         layout,
		testOnlySupportsGroupTestWrite: true,
	}
	service := New(Dependencies{
		Groups:       groups,
		Tags:         fakeTags(tagTypes),
		Destinations: sqlServerTestWriteDestinations{},
		Ledger:       recordingplan.NewService(recordingplan.NewMemoryRepository()),
	}, Config{})

	preview, err := service.Preview(t.Context(), Workspace{ID: group.WorkspaceID, Revision: testRevision}, group.ID)
	require.NoError(t, err, "Preview must succeed for instantaneous mean group with mixed gauge, counter_raw and state members")
	require.NotNil(t, preview)
	require.Equal(t, string(dbtarget.SQLDialectSQLServer), preview.Target.Dialect)
	require.Len(t, preview.RecordIDs, 3, "Preview must produce exact record IDs for gauge, counter and state members")

	// Plan inspection to verify encoded rows and provenance contracts directly.
	planObj, err := service.buildPlan(t.Context(), group.ID)
	require.NoError(t, err)
	require.NotNil(t, planObj)

	rows, err := planObj.sqlserverRows("test-owner", time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err, "sqlserverRows must encode rows without type mismatch or value blocking")
	require.Len(t, rows, 3)

	rowMap := func(r dbtarget.EncodedRow) map[string]any {
		m := make(map[string]any, len(r.Cells))
		for _, cell := range r.Cells {
			m[cell.Column] = cell.Value
		}
		return m
	}

	// Verify all rows have non-zero window bounds matching profile interval.
	expectedInterval := time.Duration(group.RecordingProfile.IntervalSeconds()) * time.Second
	for i, r := range rows {
		cells := rowMap(r)
		bucketStart, okStart := cells["bucket_start"].(time.Time)
		bucketEnd, okEnd := cells["bucket_end"].(time.Time)
		require.True(t, okStart, "row %d must contain bucket_start as time.Time", i)
		require.True(t, okEnd, "row %d must contain bucket_end as time.Time", i)
		require.True(t, bucketEnd.After(bucketStart), "row %d bucket_end must be strictly after bucket_start (non-zero window)", i)
		require.Equal(t, expectedInterval, bucketEnd.Sub(bucketStart), "row %d window length must equal profile interval", i)
	}

	// 1. Gauge member row: stores mean reduced value in value_decimal, provenance includes mean.
	gaugeCells := rowMap(rows[0])
	require.NotEmpty(t, gaugeCells["value_decimal"])
	var gaugeProv map[string]any
	require.NoError(t, json.Unmarshal([]byte(gaugeCells["provenance"].(string)), &gaugeProv))
	require.Equal(t, "decimal", gaugeProv["type"])
	require.Contains(t, gaugeProv, "mean", "gauge row under mean profile must include mean provenance")

	// 2. Counter raw member row: stores latest exact value in value_int, provenance reflects counter_raw role without mean.
	counterCells := rowMap(rows[1])
	require.NotEmpty(t, counterCells["value_int"])
	var counterProv map[string]any
	require.NoError(t, json.Unmarshal([]byte(counterCells["provenance"].(string)), &counterProv))
	require.Equal(t, "int64", counterProv["type"])
	require.Equal(t, string(workspace.WriteGroupValueRoleCounterRaw), counterProv["role"])
	require.NotContains(t, counterProv, "mean", "counter_raw member must never produce mean reduction provenance")

	// 3. State member row: stores latest exact value in value_bool, provenance reflects state role without mean.
	stateCells := rowMap(rows[2])
	require.NotNil(t, stateCells["value_bool"])
	var stateProv map[string]any
	require.NoError(t, json.Unmarshal([]byte(stateCells["provenance"].(string)), &stateProv))
	require.Equal(t, "bool", stateProv["type"])
	require.Equal(t, string(workspace.WriteGroupValueRoleState), stateProv["role"])
	require.NotContains(t, stateProv, "mean", "state member must never produce mean reduction provenance")
}
