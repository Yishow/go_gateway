package modbusshare

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

type emptyWorkspaceRevisionStore struct{}

func (emptyWorkspaceRevisionStore) GetRevision(context.Context, string) (revision string, isDirty bool, err error) {
	return "", false, nil
}

func (emptyWorkspaceRevisionStore) UpdateRevision(context.Context, string, string, string) error {
	return nil
}

func (emptyWorkspaceRevisionStore) MarkDirty(context.Context, string) error {
	return nil
}

func (emptyWorkspaceRevisionStore) GetDesiredMappings(context.Context, string) ([]DesiredMapping, error) {
	return nil, nil
}

func (emptyWorkspaceRevisionStore) CommitDesiredMappings(context.Context, string, string, string, []DesiredMapping) error {
	return nil
}

func TestReconciler_EmptyExpectedRevisionFailsClosedBeforeMutation(t *testing.T) {
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	svc := NewService(tagSvc, 128)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-empty-revision", Readiness: true})
	reconciler := NewReconciler(svc, emptyWorkspaceRevisionStore{}).WithOwnershipValidator(func(context.Context, DesiredMapping) error {
		return nil
	})
	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "empty-revision", DataType: schema.DataTypeInt16})
	require.NoError(t, err)

	for _, testCase := range []struct {
		name                      string
		expectedWorkspaceRevision string
		expectedSettingsRevision  string
	}{
		{name: "workspace", expectedSettingsRevision: "settings-1"},
		{name: "settings", expectedWorkspaceRevision: "rev-1"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			outcome, err := reconciler.Reconcile(ctx, ReconcileRequest{
				WorkspaceID:               "ws-empty-revision",
				ExpectedWorkspaceRevision: testCase.expectedWorkspaceRevision,
				ExpectedSettingsRevision:  testCase.expectedSettingsRevision,
				DesiredMappings: []DesiredMapping{{
					WorkspaceID: "ws-empty-revision", TagID: tagRecord.ID, DataType: schema.DataTypeInt16,
					ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1,
				}},
			})
			require.Error(t, err)
			var shareErr *Error
			require.ErrorAs(t, err, &shareErr)
			require.Equal(t, ErrCodeRevisionRequired, shareErr.Code)
			require.Equal(t, HydrationStateFailed, outcome.Outcome)
			require.Empty(t, svc.ListMappings())
		})
	}
}
