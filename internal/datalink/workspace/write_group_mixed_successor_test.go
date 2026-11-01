package workspace

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const mixedSuccessorTestBucketStartColumn = "bucket_start"

func TestMixedSuccessorTransitionPreservesMeaning(t *testing.T) {
	source, types := mixedSchemaFixture()
	for _, engine := range []string{"sqlite", "postgres"} {
		for _, layout := range []WriteGroupLayoutVersion{WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedLongV1} {
			target := cloneWriteGroup(source)
			target.LayoutVersion = layout
			if layout == WriteGroupLayoutManagedLongV1 {
				target.StorageMode = WriteGroupStorageModePlain
				target.RowPolicy.BucketStartColumn = mixedSuccessorTestBucketStartColumn
				target.RowPolicy.DueAtColumn = ""
			}
			target.Destination.TableName = "next_recording"
			require.NoError(t, ValidateWriteGroupTransitionStructure(source, target, engine, engine, types))
		}
	}
	changes := map[string]func(*WriteGroup){
		"record":   func(g *WriteGroup) { g.RecordID = "3f629763-f426-460f-a935-552d3a2e1ba8" },
		"method":   func(g *WriteGroup) { g.Members[0].Calculation = WriteGroupMemberCalculationLatest },
		"cadence":  func(g *WriteGroup) { g.Members[0].Cadence.Interval = 60 },
		"member":   func(g *WriteGroup) { g.Members[0].PointID = "other" },
		"source":   func(g *WriteGroup) { g.Members[0].SourceRevision = "other" },
		"type":     func(g *WriteGroup) { g.Members[0].ResultValueType = "float64" },
		"required": func(g *WriteGroup) { g.Members[0].Required = false },
		"profile":  func(g *WriteGroup) { g.RecordingProfile.Kind = WriteGroupProfileInstantaneous },
	}
	changes["proof"] = func(g *WriteGroup) {
		_, m := exactCounterBindingFixture()
		g.Members[0].ExactRecordingSource = m.ExactRecordingSource
	}
	changes["mapping"] = func(g *WriteGroup) { g.Members[0].MappingRevision = "other" }
	changes["source type"] = func(g *WriteGroup) { g.Members[0].SourceValueType = "float64" }
	changes["mode"] = func(g *WriteGroup) {
		g.LayoutVersion = WriteGroupLayoutManagedLongV1
		g.StorageMode = WriteGroupStorageModeTimescale
	}
	for name, change := range changes {
		t.Run(name, func(t *testing.T) {
			target := cloneWriteGroup(source)
			change(target)
			require.Error(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types))
		})
	}
	for _, engine := range []string{"mysql", "mariadb", "sqlserver"} {
		source.LayoutVersion = WriteGroupLayoutManagedLongV1
		require.Error(t, ValidateWriteGroupTransitionStructure(source, source, engine, engine, types))
	}
	require.Error(t, ValidateWriteGroupTransitionStructure(source, source, "sqlite", "postgres", types))
}

func TestMixedSuccessorReviewUsesExistingAppliedRecord(t *testing.T) {
	for _, delta := range []bool{false, true} {
		t.Run(map[bool]string{false: "latest", true: "delta"}[delta], func(t *testing.T) {
			now := time.Date(2026, 10, 31, 12, 0, 7, 0, time.UTC)
			var h *circuitProfileHarness
			var source *WriteGroup
			if delta {
				h, source = mixedSuccessorDeltaSource(t, &now)
			} else {
				h, source = mixedHistoryFixture(t, &now)
				independent := cloneWriteGroup(source)
				independent.ID, independent.Revision, independent.AppliedRevision = "", "", ""
				independent.RecordID = uuid.NewString()
				independent.Status = WriteGroupStatusDraft
				independent.Name = "Independent record"
				independent.Destination.TableName = "independent_record"
				independent.Destination.SchemaRevision, independent.Destination.SchemaDigest = "", ""
				created, err := h.service.Create(t.Context(), h.mutation(t.Context(), independent))
				require.NoError(t, err)
				other := mixedHistoryConfirm(t, h, created.Group)
				_, err = h.service.Apply(t.Context(), other.ID, mixedHistoryRequest(t, h, other))
				require.NoError(t, err)
			}
			var baseline func(string)
			if delta {
				baseline = mixedSuccessorBaseline(t, h, source, now)
			}

			h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
			lineage := requireLineageGenesis(t, h.db, source)
			target := cloneWriteGroup(source)
			target.LayoutVersion = WriteGroupLayoutManagedLongV1
			target.StorageMode = WriteGroupStorageModePlain
			target.RowPolicy.BucketStartColumn, target.RowPolicy.DueAtColumn = mixedSuccessorTestBucketStartColumn, ""
			target.Destination.TableName = ""
			req := WriteGroupSuccessorReviewRequest{
				WorkspaceID: source.WorkspaceID, ExpectedWorkspaceRevision: h.mutation(t.Context(), nil).ExpectedWorkspaceRevision,
				LineageID: lineage, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: source.Revision,
				ExpectedPredecessorConnectorRevision: source.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: source.Destination.ConnectorRevision,
				OperationID: "mixed-successor", Boundary: now.Add(time.Minute), Group: target,
			}
			review, err := h.service.SuccessorReview(t.Context(), source.ID, req)
			require.NoError(t, err)
			require.Empty(t, review.Blockers)
			reserved, err := h.service.ReserveSuccessor(t.Context(), source.ID, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
			require.NoError(t, err)
			require.Equal(t, source.RecordID, reserved.Group.RecordID)
			require.NotEqual(t, source.ID, reserved.Group.ID)
			if baseline != nil {
				baseline(reserved.Group.ID)
			}
			mixedSuccessorConfirmAndActivate(t.Context(), t, h, source, review, reserved, req.Boundary)
			if baseline != nil {
				baseline(reserved.Group.ID)
			}
		})
	}
}
