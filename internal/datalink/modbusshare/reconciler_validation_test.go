package modbusshare

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconciler_NilOwnershipValidatorFailsClosed(t *testing.T) {
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	svc := NewService(tagSvc, 4096)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-a", Readiness: true})
	store := newMockWorkspaceRevisionStore()
	reconciler := NewReconciler(svc, store)
	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "unverified", DataType: schema.DataTypeInt16})
	require.NoError(t, err)

	_, err = reconciler.Reconcile(ctx, ReconcileRequest{
		WorkspaceID: "ws-a", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1",
		DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-a", TagID: tagRecord.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}},
	})
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeWorkspaceScope, shareErr.Code)
	assert.False(t, svc.HasMapping(tagRecord.ID))
}

func TestReconciler_PrevalidationCollisionRejectsBeforeMutation(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)

	t1, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t1", DisplayName: "Tag 1", DataType: schema.DataTypeInt32})
	t2, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t2", DisplayName: "Tag 2", DataType: schema.DataTypeFloat32})

	// Initial valid mapping
	initReq := ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: "rev-1",
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings: []DesiredMapping{
			{
				WorkspaceID:        "ws-1",
				TagID:              t1.ID,
				DataType:           schema.DataTypeInt32,
				ZeroBasedRegister:  0,
				ShareStartRegister: 40001,
				SpanRegisters:      2,
				StrideRegisters:    2,
			},
		},
	}
	out1, err := reconciler.Reconcile(ctx, initReq)
	require.NoError(t, err)
	require.Equal(t, "applied", out1.Outcome)

	// Now try to reconcile with a collision: t2 at register 1 overlaps t1 at [0, 2)
	badReq := ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: out1.NewWorkspaceRevision,
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings: []DesiredMapping{
			{
				WorkspaceID:        "ws-1",
				TagID:              t1.ID,
				DataType:           schema.DataTypeInt32,
				ZeroBasedRegister:  0,
				ShareStartRegister: 40001,
				SpanRegisters:      2,
				StrideRegisters:    2,
			},
			{
				WorkspaceID:        "ws-1",
				TagID:              t2.ID,
				DataType:           schema.DataTypeFloat32,
				ZeroBasedRegister:  1, // Collision!
				ShareStartRegister: 40002,
				SpanRegisters:      2,
				StrideRegisters:    2,
			},
		},
	}

	out2, err := reconciler.Reconcile(ctx, badReq)
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeRangeCollision, shareErr.Code)
	assert.Equal(t, "failed", out2.Outcome)

	// Previous projection MUST remain active! t1 exists, t2 does NOT exist.
	assert.True(t, svc.HasMapping(t1.ID))
	assert.False(t, svc.HasMapping(t2.ID))
}

func TestReconciler_ForeignWorkspaceProjectionIsPreservedAndExcluded(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)
	foreign := TagMirrorMapping{WorkspaceID: "ws-b", SourceRuleID: "rule-b", SourceRuleRevision: "rev-b", TagID: "foreign", Register: 4, ZeroBasedRegister: 4, ShareStartRegister: 40005, SpanRegisters: 1, StrideRegisters: 1, DataType: schema.DataTypeInt16}
	svc.ReplaceMappings(map[string]TagMirrorMapping{foreign.TagID: foreign})
	tagA, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "workspace-a", DataType: schema.DataTypeInt16})
	require.NoError(t, err)

	out, err := reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", SourceRuleID: "rule-a", SourceRuleRevision: "rev-a", TagID: tagA.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}})
	require.NoError(t, err)
	require.Len(t, out.Mappings, 1)
	require.Equal(t, tagA.ID, out.Mappings[0].TagID)
	all := svc.ListMappings()
	require.Len(t, all, 2)
	require.NotNil(t, findMapping(all, foreign.TagID))
}

func TestReconciler_CollisionIncludesForeignWorkspaceSpans(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)
	foreignTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "workspace-b-span", DataType: schema.DataTypeInt32})
	require.NoError(t, err)
	svc.ReplaceMappings(map[string]TagMirrorMapping{foreignTag.ID: {
		WorkspaceID: "ws-b", SourceRuleID: "rule-b", SourceRuleRevision: "rev-b", TagID: foreignTag.ID,
		Register: 0, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 2, StrideRegisters: 2, DataType: schema.DataTypeInt32,
	}})
	localTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "workspace-a-span", DataType: schema.DataTypeInt16})
	require.NoError(t, err)

	_, err = reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{{
		WorkspaceID: "ws-1", TagID: localTag.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1,
	}}})
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeRangeCollision, shareErr.Code)
	assert.True(t, svc.HasMapping(foreignTag.ID))
	assert.False(t, svc.HasMapping(localTag.ID))
}

func TestReconciler_EmptyHydrationWorkspaceRejectsNormalRequestButRestoreMayBypass(t *testing.T) {
	ctx := context.Background()
	reconciler, _, tagSvc, _ := setupReconcilerFixture(t)
	reconciler.svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	tagA, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "restore-only", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	desired := []DesiredMapping{{WorkspaceID: "ws-1", TagID: tagA.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}
	_, err = reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: desired})
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeHydrationRequired, shareErr.Code)

	out, err := reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: desired, Restore: true})
	require.NoError(t, err)
	require.Equal(t, "applied", out.Outcome)
}

func findMapping(mappings []TagMirrorMapping, tagID string) *TagMirrorMapping {
	for _, mapping := range mappings {
		if mapping.TagID == tagID {
			copyMapping := mapping
			return &copyMapping
		}
	}
	return nil
}

func TestReconcileKeySortsAllDesiredMappingFieldsDeterministically(t *testing.T) {
	left := DesiredMapping{WorkspaceID: "ws-1", SourceRuleID: "rule-a", SourceRuleRevision: "r1", TagID: "tag", DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}
	right := DesiredMapping{WorkspaceID: "ws-1", SourceRuleID: "rule-b", SourceRuleRevision: "r1", TagID: "tag", DataType: schema.DataTypeInt16, ZeroBasedRegister: 1, ShareStartRegister: 40002, SpanRegisters: 1, StrideRegisters: 1}
	first, err := reconcileKey(ReconcileRequest{WorkspaceID: "ws-1", DesiredMappings: []DesiredMapping{left, right}})
	require.NoError(t, err)
	second, err := reconcileKey(ReconcileRequest{WorkspaceID: "ws-1", DesiredMappings: []DesiredMapping{right, left}})
	require.NoError(t, err)
	assert.Equal(t, first, second)
}

func TestReconciler_DuplicateTagIdentityRejectsBeforeRuntimeRevisionMutation(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, revStore := setupReconcilerFixture(t)
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "duplicate-reconcile", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	valid := DesiredMapping{WorkspaceID: "ws-1", TagID: t1.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}
	first, err := reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{valid}})
	require.NoError(t, err)
	duplicate := valid
	duplicate.ZeroBasedRegister = 2
	duplicate.ShareStartRegister = 40003
	_, err = reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: first.NewWorkspaceRevision, ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{valid, duplicate}})
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeInvalidGeometry, shareErr.Code)
	assert.Equal(t, first.NewWorkspaceRevision, revStore.revisions["ws-1"])
	assert.Equal(t, uint16(0), svc.ListMappings()[0].ZeroBasedRegister)
}

func TestReconciler_IdempotencyValidatesCurrentMappingOwnershipBeforeReplay(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "owner-replay", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	reconciler.WithOwnershipValidator(func(_ context.Context, dm DesiredMapping) error {
		if dm.SourceRuleRevision != "r1" {
			return NewError(ErrCodeWorkspaceScope, "mapping ownership changed", true)
		}
		return nil
	})
	req := ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", SourceRuleID: "rule-1", SourceRuleRevision: "r1", TagID: t1.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}}
	_, err = reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	mapping := svc.ListMappings()[0]
	mapping.SourceRuleRevision = "r2"
	svc.ReplaceMappings(map[string]TagMirrorMapping{t1.ID: mapping})

	_, err = reconciler.Reconcile(ctx, req)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeWorkspaceScope, shareErr.Code)
}

func TestReconciler_EquivalentComparisonIncludesOwnerRevisionAndMetadata(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "equivalence", DisplayName: "Equivalent", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	first, err := reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", SourceRuleID: "rule-1", SourceRuleRevision: "r1", TagID: t1.ID, TagKey: "old-key", DisplayName: "Old", DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}})
	require.NoError(t, err)
	second, err := reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: first.NewWorkspaceRevision, ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", SourceRuleID: "rule-1", SourceRuleRevision: "r2", TagID: t1.ID, TagKey: "new-key", DisplayName: "New", DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}})
	require.NoError(t, err)
	assert.Equal(t, "applied", second.Outcome)
	assert.Equal(t, "r2", svc.ListMappings()[0].SourceRuleRevision)
	assert.Equal(t, "new-key", svc.ListMappings()[0].TagKey)
}
