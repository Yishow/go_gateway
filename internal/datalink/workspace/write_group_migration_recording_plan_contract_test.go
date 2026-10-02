package workspace

import (
	"fmt"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func recordingPlanMigrationContractStream(mode recordingplan.StreamMode, rawPolicy recordingplan.RawPolicy) recordingplan.PlanStream {
	deadband := 0.25
	heartbeat := 45
	summaryInterval := 60
	usageInterval := 90
	batchMember := "member-A"
	batchTimeout := 30
	return recordingplan.PlanStream{
		StreamID:               "stream-contract",
		MeasurementID:          "measurement-A",
		EquipmentID:            "equipment-A",
		Mode:                   mode,
		RawPolicy:              rawPolicy,
		OnChangeDeadband:       &deadband,
		MaxHeartbeatSeconds:    &heartbeat,
		SummaryIntervalSeconds: &summaryInterval,
		UsageIntervalSeconds:   &usageInterval,
		BatchTriggerMemberID:   &batchMember,
		BatchTimeoutSeconds:    &batchTimeout,
		DestinationIDs:         []string{"destination-A"},
	}
}

func TestPreviewRecordingPlanMigrationBlocksAllLegacyStreamPoliciesAndPreservesPayload(t *testing.T) {
	cases := []struct {
		name      string
		mode      recordingplan.StreamMode
		rawPolicy recordingplan.RawPolicy
		issueCode string
	}{
		{name: "raw-on-change", mode: recordingplan.StreamModeRawHistory, rawPolicy: recordingplan.RawPolicyOnChange, issueCode: "recording-plan-raw-policy-unsupported"},
		{name: "raw-sampled", mode: recordingplan.StreamModeRawHistory, rawPolicy: recordingplan.RawPolicySampled, issueCode: "recording-plan-raw-policy-unsupported"},
		{name: "raw-empty-policy", mode: recordingplan.StreamModeRawHistory, issueCode: "recording-plan-raw-policy-unsupported"},
		{name: "window-summary", mode: recordingplan.StreamModeWindowSummary, issueCode: "recording-plan-stream-mode-unsupported"},
		{name: "usage-interval", mode: recordingplan.StreamModeUsageInterval, issueCode: "recording-plan-stream-mode-unsupported"},
		{name: "state-changes", mode: recordingplan.StreamModeStateChanges, issueCode: "recording-plan-stream-mode-unsupported"},
		{name: "event-log", mode: recordingplan.StreamModeEventLog, issueCode: "recording-plan-stream-mode-unsupported"},
		{name: "batch-snapshot", mode: recordingplan.StreamModeBatchSnapshot, issueCode: "recording-plan-stream-mode-unsupported"},
		{name: "latest-only", mode: recordingplan.StreamModeLatestOnly, issueCode: "recording-plan-stream-mode-unsupported"},
		{name: "unknown-mode", mode: recordingplan.StreamMode("unknown_mode"), issueCode: "recording-plan-stream-mode-unsupported"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			fixture := newRecordingPlanMigrationFixture(ctx, t)
			defer fixture.close()

			plan := *fixture.plan
			plan.ID = "legacy-plan-" + tc.name
			plan.Name = "Legacy contract " + tc.name
			plan.Streams = []recordingplan.PlanStream{recordingPlanMigrationContractStream(tc.mode, tc.rawPolicy)}
			require.NoError(t, recordingplan.NewSQLRepository(fixture.db).CreatePlan(ctx, &plan))
			before, err := recordingplan.NewSQLRepository(fixture.db).GetPlanByWorkspace(ctx, plan.ID, fixture.base.workspaceID)
			require.NoError(t, err)

			preview, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{plan.ID})
			require.NoError(t, err)
			require.Len(t, preview.Items, 1)
			item := preview.Items[0]
			require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, item.Status)
			require.Nil(t, item.CandidateGroup)
			require.True(t, findingCodes(item.Issues)[tc.issueCode])
			require.NotNil(t, item.BeforeRecordingPlanIntent)
			require.Equal(t, *before, item.BeforeRecordingPlanIntent.Plan)

			_, err = fixture.service.ReviewRecordingPlanMigration(ctx, recordingPlanReviewRequest(preview))
			require.ErrorIs(t, err, ErrWriteGroupPlanMigrationBlocked)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			after, err := recordingplan.NewSQLRepository(fixture.db).GetPlanByWorkspace(ctx, plan.ID, fixture.base.workspaceID)
			require.NoError(t, err)
			require.Equal(t, *before, *after)
		})
	}
}

func TestPreviewRecordingPlanMigrationBlocksMissingMeasurementAndPreservesPlan(t *testing.T) {
	ctx := t.Context()
	fixture := newRecordingPlanMigrationFixture(ctx, t)
	defer fixture.close()

	plan := *fixture.plan
	plan.ID = "legacy-plan-missing-measurement"
	plan.Members[0].MeasurementID = "missing-measurement"
	plan.Streams[0].MeasurementID = "missing-measurement"
	require.NoError(t, recordingplan.NewSQLRepository(fixture.db).CreatePlan(ctx, &plan))
	before, err := recordingplan.NewSQLRepository(fixture.db).GetPlanByWorkspace(ctx, plan.ID, fixture.base.workspaceID)
	require.NoError(t, err)

	preview, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{plan.ID})
	require.NoError(t, err)
	item := preview.Items[0]
	require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, item.Status)
	require.Nil(t, item.CandidateGroup)
	require.True(t, findingCodes(item.Issues)["recording-plan-source-missing"])
	require.Len(t, item.BeforeRecordingPlanIntent.Sources, 1)
	require.Equal(t, recordingPlanSourceStatusBlocked, item.BeforeRecordingPlanIntent.Sources[0].Status)
	require.Equal(t, "missing-measurement", item.BeforeRecordingPlanIntent.Sources[0].MeasurementID)

	_, err = fixture.service.ReviewRecordingPlanMigration(ctx, recordingPlanReviewRequest(preview))
	require.ErrorIs(t, err, ErrWriteGroupPlanMigrationBlocked)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	after, err := recordingplan.NewSQLRepository(fixture.db).GetPlanByWorkspace(ctx, plan.ID, fixture.base.workspaceID)
	require.NoError(t, err)
	require.Equal(t, *before, *after)
}

func TestPreviewRecordingPlanMigrationBlocksIncompleteMeasurementSourceRevisions(t *testing.T) {
	cases := []struct {
		name   string
		column string
	}{
		{name: "definition", column: "definition_revision"},
		{name: "source-binding", column: "source_binding_revision"},
		{name: "series-epoch", column: "series_epoch"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			fixture := newRecordingPlanMigrationFixture(ctx, t)
			defer fixture.close()
			_, err := fixture.db.ExecContext(ctx, fmt.Sprintf("UPDATE measurement_definitions SET %s = ? WHERE id = ?", tc.column), "", fixture.measurementID)
			require.NoError(t, err)

			preview, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{fixture.plan.ID})
			require.NoError(t, err)
			item := preview.Items[0]
			require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, item.Status)
			require.Nil(t, item.CandidateGroup)
			require.True(t, findingCodes(item.Issues)["recording-plan-source-revision-missing"])
			require.Equal(t, recordingPlanSourceStatusBlocked, item.BeforeRecordingPlanIntent.Sources[0].Status)

			_, err = fixture.service.ReviewRecordingPlanMigration(ctx, recordingPlanReviewRequest(preview))
			require.ErrorIs(t, err, ErrWriteGroupPlanMigrationBlocked)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
		})
	}
}

func TestPreviewRecordingPlanMigrationReturnsSafeNotFoundForForeignMeasurementChains(t *testing.T) {
	cases := []struct {
		name             string
		foreignWorkspace bool
	}{
		{name: "foreign-workspace", foreignWorkspace: true},
		{name: "foreign-device"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			fixture := newRecordingPlanMigrationFixture(ctx, t)
			defer fixture.close()
			tagID := fixture.base.tagID
			workspaceID := fixture.base.workspaceID
			if tc.foreignWorkspace {
				workspaceID = "foreign-workspace"
			}
			foreign := &measurement.MeasurementDefinition{
				ID:                    "measurement-" + tc.name,
				WorkspaceID:           workspaceID,
				DeviceID:              "foreign-device",
				PointID:               "foreign-point",
				TagID:                 &tagID,
				EquipmentID:           "foreign-equipment",
				DefinitionRevision:    "foreign-definition-revision",
				SourceBindingRevision: "foreign-binding-revision",
				SeriesEpoch:           "foreign-epoch",
				Name:                  "Foreign temperature",
				Quantity:              "temperature",
				Unit:                  "C",
				SemanticKind:          measurement.SemanticKindGauge,
				NumericEncoding:       "float64",
			}
			require.NoError(t, measurement.NewSQLRepository(fixture.db).Create(ctx, foreign))
			plan := *fixture.plan
			plan.ID = "legacy-plan-" + tc.name
			plan.Members[0].MeasurementID = foreign.ID
			plan.Streams[0].MeasurementID = foreign.ID
			require.NoError(t, recordingplan.NewSQLRepository(fixture.db).CreatePlan(ctx, &plan))

			_, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{plan.ID})
			require.ErrorIs(t, err, ErrWriteGroupNotFound)
		})
	}
}
