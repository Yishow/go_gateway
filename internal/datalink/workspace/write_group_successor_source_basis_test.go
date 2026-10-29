package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecreateRepositoryAppliedMeaningCannotBeReplacedByUnappliedDraft(t *testing.T) {
	for _, change := range []string{"reduction", "interval"} {
		t.Run(change, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			seedRecreateEngineEvidence(t, h, req)
			if change == "reduction" {
				_, err := h.db.Exec("UPDATE write_groups SET sample_reduction='mean',revision='unapplied-R2'; UPDATE write_group_members SET value_role='gauge',source_value_type='float32',result_value_type='float64'")
				require.NoError(t, err)
			} else {
				_, err := h.db.Exec("UPDATE write_groups SET profile_interval=profile_interval+1,revision='unapplied-R2'")
				require.NoError(t, err)
			}
			draft, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
			require.NoError(t, err)
			req.ExpectedPredecessorRevision = draft.Revision
			req.Group = cloneWriteGroup(draft)
			before := recreateRepositorySnapshot(t, h)
			_, err = reserveRepositoryTest(t, h, req)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			var rejected *WriteGroupTransitionError
			require.ErrorAs(t, err, &rejected)
			require.Equal(t, "recording-meaning-change", rejected.ReasonCode)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestRecreateRepositoryCanSelectOriginalAppliedDefinitionOverUnappliedDraft(t *testing.T) {
	for _, change := range []string{"reduction", "destination"} {
		t.Run(change, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			source := seedRecreateEngineEvidence(t, h, req)
			if change == "reduction" {
				_, err := h.db.Exec("UPDATE write_groups SET sample_reduction='mean',revision='unapplied-R2'; UPDATE write_group_members SET value_role='gauge',source_value_type='float32',result_value_type='float64'")
				require.NoError(t, err)
			} else {
				changed := req
				changed.Group = cloneWriteGroup(req.Group)
				recreateRepositoryNewTarget(t, h, &changed)
				_, err := h.db.Exec(`UPDATE write_groups SET revision='unapplied-R2',destination_connector_id='recreate-target',destination_connector_revision='target-rev',destination_database='new-db',destination_table_schema='main',destination_table_name='draft-table'`)
				require.NoError(t, err)
				req.ExpectedPredecessorConnectorRevision = "target-rev"
			}
			req.ExpectedPredecessorRevision = "unapplied-R2"
			req.Group = cloneWriteGroup(source.Group)
			before, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
			require.NoError(t, err)
			result, err := reserveRepositoryTest(t, h, req)
			require.NoError(t, err)
			require.Equal(t, source.Group.RecordingProfile, result.Group.RecordingProfile)
			require.Equal(t, source.Group.Destination.ConnectorID, result.Group.Destination.ConnectorID)
			require.Equal(t, source.Group.Destination.ConnectorRevision, result.Group.Destination.ConnectorRevision)
			after, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
			require.NoError(t, err)
			require.Equal(t, before, after)
		})
	}
}

func TestRecreateRepositoryAppliedSourceDefinitionMustRemainAvailableAndScoped(t *testing.T) {
	for _, query := range []string{"DELETE FROM write_group_versions", "UPDATE write_group_versions SET payload=json_set(payload,'$.workspace_id','foreign')"} {
		t.Run(query, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			seedRecreateEngineEvidence(t, h, req)
			_, err := h.db.Exec(query)
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			_, err = reserveRepositoryTest(t, h, req)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}
