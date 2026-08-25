package modbusshare

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconciler_DeleteClearsCompleteOldSpan(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)

	// Create int64 tag (span=4)
	tInt64, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t_int64", DisplayName: "Int64", DataType: schema.DataTypeInt64})

	// Initial apply
	initReq := ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: "rev-1",
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings: []DesiredMapping{
			{
				WorkspaceID:        "ws-1",
				SourceRuleID:       "rule-1",
				SourceRuleRevision: "rev-1",
				TagID:              tInt64.ID,
				DataType:           schema.DataTypeInt64,
				ZeroBasedRegister:  0,
				ShareStartRegister: 40001,
				SpanRegisters:      4,
				StrideRegisters:    4,
			},
		},
	}
	out1, err := reconciler.Reconcile(ctx, initReq)
	require.NoError(t, err)
	require.Equal(t, 1, out1.AppliedCount)
	require.NoError(t, svc.WriteTagValue(ctx, tInt64.ID, int64(123456)))
	words, err := svc.ReadHoldingWords(0, 4)
	require.NoError(t, err)
	assert.NotEqual(t, []uint16{0, 0, 0, 0}, words)

	// Now reconcile empty desired set (deleting the int64 mapping)
	delReq := ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: out1.NewWorkspaceRevision,
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings:           []DesiredMapping{},
	}
	out2, err := reconciler.Reconcile(ctx, delReq)
	require.NoError(t, err)
	assert.Equal(t, "applied", out2.Outcome)
	assert.Equal(t, 1, out2.RemovedCount)
	require.Len(t, out2.RemovedSpans, 1)
	assert.Equal(t, tInt64.ID, out2.RemovedSpans[0].TagID)
	assert.Equal(t, uint16(0), out2.RemovedSpans[0].Start)
	assert.Equal(t, 4, out2.RemovedSpans[0].Count)
	assert.False(t, svc.HasMapping(tInt64.ID))
	words, err = svc.ReadHoldingWords(0, 4)
	require.NoError(t, err)
	assert.Equal(t, []uint16{0, 0, 0, 0}, words)
}

func TestReconciler_RetypeFromFloat32ToInt64InvalidatesOldSpan(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)

	// Initially float32 (span=2) at start 0 -> [0, 2)
	t1, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "val", DisplayName: "Value", DataType: schema.DataTypeFloat32})

	out1, err := reconciler.Reconcile(ctx, ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: "rev-1",
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings: []DesiredMapping{
			{
				WorkspaceID:        "ws-1",
				SourceRuleID:       "rule-1",
				SourceRuleRevision: "rev-1",
				TagID:              t1.ID,
				DataType:           schema.DataTypeFloat32,
				ZeroBasedRegister:  0,
				ShareStartRegister: 40001,
				SpanRegisters:      2,
				StrideRegisters:    2,
			},
		},
	})
	require.NoError(t, err)

	// Retype to int64 (span=4) at start 0 -> [0, 4)
	out2, err := reconciler.Reconcile(ctx, ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: out1.NewWorkspaceRevision,
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings: []DesiredMapping{
			{
				WorkspaceID:        "ws-1",
				SourceRuleID:       "rule-1",
				SourceRuleRevision: "rev-2",
				TagID:              t1.ID,
				DataType:           schema.DataTypeInt64,
				ZeroBasedRegister:  0,
				ShareStartRegister: 40001,
				SpanRegisters:      4,
				StrideRegisters:    4,
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "applied", out2.Outcome)
	assert.Equal(t, 1, out2.InvalidatedCount)
	require.Len(t, out2.InvalidatedSpans, 1)
	assert.Equal(t, 2, out2.InvalidatedSpans[0].Count) // Old float32 span 2 was invalidated

	mList := svc.ListMappings()
	require.Len(t, mList, 1)
	assert.Equal(t, 4, mList[0].SpanRegisters)
	assert.Equal(t, schema.DataTypeInt64, mList[0].DataType)
}

func TestReconciler_MoveWithRangeCollisionPreservesOldProjection(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)

	t1, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t1", DisplayName: "Tag 1", DataType: schema.DataTypeInt16})
	t2, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t2", DisplayName: "Tag 2", DataType: schema.DataTypeInt32})

	// Initial mappings: t1 at 0 [0, 1), t2 at 10 [10, 12)
	out1, err := reconciler.Reconcile(ctx, ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: "rev-1",
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings: []DesiredMapping{
			{
				WorkspaceID:        "ws-1",
				TagID:              t1.ID,
				DataType:           schema.DataTypeInt16,
				ZeroBasedRegister:  0,
				ShareStartRegister: 40001,
				SpanRegisters:      1,
				StrideRegisters:    1,
			},
			{
				WorkspaceID:        "ws-1",
				TagID:              t2.ID,
				DataType:           schema.DataTypeInt32,
				ZeroBasedRegister:  10,
				ShareStartRegister: 40011,
				SpanRegisters:      2,
				StrideRegisters:    2,
			},
		},
	})
	require.NoError(t, err)

	// Move t2 to 0 -> COLLISION with t1 at [0, 1)
	out2, err := reconciler.Reconcile(ctx, ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: out1.NewWorkspaceRevision,
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings: []DesiredMapping{
			{
				WorkspaceID:        "ws-1",
				TagID:              t1.ID,
				DataType:           schema.DataTypeInt16,
				ZeroBasedRegister:  0,
				ShareStartRegister: 40001,
				SpanRegisters:      1,
				StrideRegisters:    1,
			},
			{
				WorkspaceID:        "ws-1",
				TagID:              t2.ID,
				DataType:           schema.DataTypeInt32,
				ZeroBasedRegister:  0, // Collision!
				ShareStartRegister: 40001,
				SpanRegisters:      2,
				StrideRegisters:    2,
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, "failed", out2.Outcome)

	// Old projection MUST be preserved: t1 at 0, t2 at 10
	mList := svc.ListMappings()
	require.Len(t, mList, 2)
	for _, m := range mList {
		if m.TagID == t1.ID {
			assert.Equal(t, uint16(0), m.ZeroBasedRegister)
		}
		if m.TagID == t2.ID {
			assert.Equal(t, uint16(10), m.ZeroBasedRegister)
		}
	}
}
