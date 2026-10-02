package workspace

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupWriterOwnershipActivationBarrierDoesNotRunBeforeReadinessOrCAS(t *testing.T) {
	tests := []struct {
		name          string
		mutate        func(context.Context, *testing.T, *WriteGroupService, writeGroupFixture, *WriteGroupSaveResult, string) string
		withInspector bool
	}{
		{
			name:          "unready",
			withInspector: false,
			mutate: func(context.Context, *testing.T, *WriteGroupService, writeGroupFixture, *WriteGroupSaveResult, string) string {
				return ""
			},
		},
		{
			name:          "stale workspace",
			withInspector: true,
			mutate: func(ctx context.Context, t *testing.T, service *WriteGroupService, _ writeGroupFixture, _ *WriteGroupSaveResult, revision string) string {
				_, err := service.workspaceSvc.UpdateDatabaseSetup(ctx, revision, func(context.Context, *sql.Tx, *Record) error { return nil })
				require.NoError(t, err)
				return revision
			},
		},
		{
			name:          "stale group",
			withInspector: true,
			mutate: func(ctx context.Context, t *testing.T, service *WriteGroupService, fixture writeGroupFixture, created *WriteGroupSaveResult, revision string) string {
				edited := cloneWriteGroup(created.Group)
				edited.Name = "stale-group-edit"
				updated, err := service.Update(ctx, created.Group.ID, WriteGroupMutation{
					WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: revision,
					ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
					Group: edited,
				})
				require.NoError(t, err)
				return updated.WorkspaceRevision
			},
		},
		{
			name:          "stale connector",
			withInspector: true,
			mutate: func(ctx context.Context, t *testing.T, service *WriteGroupService, fixture writeGroupFixture, _ *WriteGroupSaveResult, revision string) string {
				_, err := service.repo.db.ExecContext(ctx,
					`UPDATE database_connectors SET identity_revision = ? WHERE id = ?`,
					"connector-revision-2", fixture.connectorID,
				)
				require.NoError(t, err)
				return revision
			},
		},
		{
			name:          "stale source",
			withInspector: true,
			mutate: func(ctx context.Context, t *testing.T, service *WriteGroupService, fixture writeGroupFixture, _ *WriteGroupSaveResult, revision string) string {
				_, err := service.repo.db.ExecContext(ctx, `
					UPDATE mappings SET transform_pipeline = ? WHERE point_id = ? AND tag_id = ?
				`, `[{"kind":"scale","factor":2}]`, fixture.pointID, fixture.tagID)
				require.NoError(t, err)
				return revision
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			service, fixture, created := newReadinessGroup(ctx, t, db)
			if tc.withInspector {
				service.WithTableInspector(writeGroupLifecycleInspector{})
			}
			probe := &writerOwnershipActivationProbe{}
			service.WithWriterOwnershipActivationBarrier(probe)
			workspaceBefore, err := service.workspaceSvc.GetOrCreate(ctx)
			require.NoError(t, err)
			expectedWorkspaceRevision := workspaceBefore.DatabaseSetupRevision
			if tc.mutate != nil {
				expectedWorkspaceRevision = tc.mutate(ctx, t, service, fixture, created, expectedWorkspaceRevision)
			}
			_, err = service.Apply(ctx, created.Group.ID, WriteGroupMutation{
				WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: expectedWorkspaceRevision,
				ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
			})
			require.Error(t, err)
			require.Zero(t, probe.calls)
		})
	}
}
