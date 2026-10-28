package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink"
	"go-gateway/internal/testutil/lostcommit"
)

// The seam intercepts only the real cutover transaction, after all writes but
// before commit. Readiness transactions remain unmodified.
type ownershipFenceCommitRepository struct {
	*SQLRepository
	fail, after              bool
	pendingID, pendingDetail string
	digest                   string
	commits                  int
}

func (r *ownershipFenceCommitRepository) BeginDatabaseSetup(ctx context.Context) (DatabaseSetupTx, error) {
	tx, err := r.SQLRepository.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	return &ownershipFenceCommitTx{DatabaseSetupTx: tx, repo: r}, nil
}

type ownershipFenceCommitTx struct {
	DatabaseSetupTx
	repo *ownershipFenceCommitRepository
}

func (tx *ownershipFenceCommitTx) Commit() error {
	var id, detail string
	err := tx.SQLTx().QueryRow(`SELECT id, detail FROM write_group_cutovers WHERE review_digest=?`, tx.repo.digest).Scan(&id, &detail)
	if errors.Is(err, sql.ErrNoRows) {
		return tx.DatabaseSetupTx.Commit()
	}
	if err != nil {
		return err
	}
	tx.repo.pendingID, tx.repo.pendingDetail = id, detail
	tx.repo.commits++
	if tx.repo.fail && !tx.repo.after {
		return errCounterBeforeCommit
	}
	if err := tx.DatabaseSetupTx.Commit(); err != nil {
		return err
	}
	if tx.repo.fail {
		return lostcommit.ErrResponseLost
	}
	return nil
}
func fenceRevisionInDetail(t *testing.T, detail string) string {
	t.Helper()
	var decoded struct {
		Revision string `json:"ownership_fence_revision"`
	}
	require.NoError(t, json.Unmarshal([]byte(detail), &decoded))
	return decoded.Revision
}

func TestCutoverOwnershipFenceCommitBoundaries(t *testing.T) {
	for _, mode := range []string{"success", "before commit", "committed response lost"} {
		t.Run(mode, func(t *testing.T) {
			f, target, req := counterCutoverFixture(t)
			_, original := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
			repo := &ownershipFenceCommitRepository{SQLRepository: NewSQLRepository(f.db), digest: req.ReviewDigest, fail: mode != "success", after: mode == "committed response lost"}
			f.service.workspaceSvc.repo = repo
			result, err := f.service.ApplyCutover(t.Context(), req)
			switch mode {
			case "success":
				require.NoError(t, err)
				require.NotNil(t, result)
			case "before commit":
				require.ErrorIs(t, err, errCounterBeforeCommit)
				require.Nil(t, result)
			default:
				require.ErrorIs(t, err, lostcommit.ErrResponseLost)
				require.Nil(t, result)
			}
			require.Equal(t, 1, repo.commits)
			require.NotEmpty(t, repo.pendingID)
			require.Equal(t, repo.pendingID, fenceRevisionInDetail(t, repo.pendingDetail), "cutover must carry its durable ownership fence revision")
			var cutovers, ends, scopes int
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_cutovers`).Scan(&cutovers))
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_revision_ends`).Scan(&ends))
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences`).Scan(&scopes))
			loaded, err := f.service.Get(t.Context(), target.ID)
			require.NoError(t, err)
			if mode == "before commit" {
				require.Zero(t, cutovers)
				require.Zero(t, ends)
				require.Zero(t, scopes)
				require.Empty(t, loaded.Group.AppliedRevision)
			} else {
				require.Equal(t, 1, cutovers)
				require.Equal(t, 1, ends)
				require.Equal(t, 2, scopes)
				require.Equal(t, target.Revision, loaded.Group.AppliedRevision)
			}
			_, after := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
			require.Equal(t, original, after, "fence cannot rewrite original proof/history")
			repo.fail = false
			retry, err := f.service.ApplyCutover(t.Context(), req)
			require.NoError(t, err)
			require.Equal(t, retry.ID, retry.OwnershipFenceRevision)
			if mode != "before commit" {
				require.Equal(t, repo.pendingID, retry.ID)
				require.Equal(t, 1, repo.commits)
			}
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences`).Scan(&scopes))
			require.Equal(t, 2, scopes, "retry must not create a second owner epoch")
			var detail string
			require.NoError(t, f.db.QueryRow(`SELECT detail FROM write_group_cutovers WHERE id=?`, retry.ID).Scan(&detail))
			require.Equal(t, retry.ID, fenceRevisionInDetail(t, detail))
		})
	}
}

func TestCutoverOwnershipFenceScopeAndRollback(t *testing.T) {
	for _, layout := range []WriteGroupLayoutVersion{WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1} {
		t.Run(string(layout), func(t *testing.T) {
			f, first, req := counterCutoverFixtureWithTarget(t, layout)
			forward, err := f.service.ApplyCutover(t.Context(), req)
			require.NoError(t, err)
			require.Equal(t, forward.ID, forward.OwnershipFenceRevision)
			assertFenceScope(t, f, forward, "ended", f.group)
			assertFenceScope(t, f, forward, "started", first)
			var count int
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision=?`, forward.ID).Scan(&count))
			require.Equal(t, 2, count)
			rollback := reviewCounterRollback(t, f, forward)
			repo := &ownershipFenceCommitRepository{SQLRepository: NewSQLRepository(f.db), digest: rollback.ReviewDigest, fail: true}
			f.service.workspaceSvc.repo = repo
			result, err := f.service.ApplyCutover(t.Context(), rollback)
			require.ErrorIs(t, err, errCounterBeforeCommit)
			require.Nil(t, result)
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences`).Scan(&count))
			require.Equal(t, 2, count, "rollback failure preserves forward owner epoch")
			repo.fail = false
			restored, err := f.service.ApplyCutover(t.Context(), rollback)
			require.NoError(t, err)
			require.NotEqual(t, forward.OwnershipFenceRevision, restored.OwnershipFenceRevision)
			require.Equal(t, restored.ID, restored.OwnershipFenceRevision)
			assertFenceScope(t, f, restored, "ended", first)
			active, err := f.service.Get(t.Context(), f.group.ID)
			require.NoError(t, err)
			assertFenceScope(t, f, restored, "started", active.Group)
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences`).Scan(&count))
			require.Equal(t, 4, count)
		})
	}
}

func assertFenceScope(t *testing.T, f counterApplyFixture, c *WriteGroupCutover, direction string, g *WriteGroup) {
	t.Helper()
	var workspace, device, point, tag, kind, connector, connectorRevision, database, schemaName, table string
	require.NoError(t, f.db.QueryRow(`SELECT workspace_id,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name FROM write_group_ownership_fences WHERE fence_revision=? AND direction=? AND group_id=? AND group_revision=?`, c.ID, direction, g.ID, g.Revision).Scan(&workspace, &device, &point, &tag, &kind, &connector, &connectorRevision, &database, &schemaName, &table))
	require.Equal(t, g.WorkspaceID, workspace)
	require.Equal(t, g.Members[0].DeviceID, device)
	require.Equal(t, g.Members[0].PointID, point)
	require.Equal(t, g.Members[0].TagID, tag)
	require.Equal(t, string(g.RecordingProfile.Kind), kind)
	require.Equal(t, g.Destination.ConnectorID, connector)
	require.Equal(t, g.Destination.ConnectorRevision, connectorRevision)
	require.Equal(t, g.Destination.Database, database)
	require.Equal(t, g.Destination.TableSchema, schemaName)
	require.Equal(t, g.Destination.TableName, table)
}

func TestCutoverOwnershipFenceWriteFailureIsAtomic(t *testing.T) {
	f, target, req := counterCutoverFixture(t)
	// Interrupt after the source scope is inserted, before the target scope.
	_, err := f.db.Exec(`CREATE TRIGGER fail_fence BEFORE INSERT ON write_group_ownership_fences WHEN NEW.direction='started' BEGIN SELECT RAISE(ABORT, 'injected fence persistence failure'); END`)
	require.NoError(t, err)
	result, err := f.service.ApplyCutover(t.Context(), req)
	require.ErrorContains(t, err, "injected fence persistence failure")
	require.Nil(t, result)
	for _, table := range []string{"write_group_ownership_fences", "write_group_cutovers", "write_group_revision_ends"} {
		var count int
		require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM `+table).Scan(&count))
		require.Zero(t, count)
	}
	loaded, err := f.service.Get(t.Context(), target.ID)
	require.NoError(t, err)
	require.Empty(t, loaded.Group.AppliedRevision)
	end, err := f.service.RevisionEnd(t.Context(), f.group.ID, f.group.Revision)
	require.NoError(t, err)
	require.True(t, end.IsZero())
	_, err = f.db.Exec(`DROP TRIGGER fail_fence`)
	require.NoError(t, err)
	result, err = f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, result.ID, result.OwnershipFenceRevision)
}

// Repository fixture isolates scope serialization; the existing API partial
// cutover regression exercises selection of these ended/started revisions.
func TestCutoverOwnershipFenceIncludesRetainedSuccessor(t *testing.T) {
	f, target, req := counterCutoverFixture(t)
	full := cloneWriteGroup(f.group)
	full.Revision, full.AppliedRevision = "full-revision", "full-revision"
	second := full.Members[0]
	second.PointID, second.TagID = "retained-point", "retained-tag"
	full.Members = append(full.Members, second)
	remaining := cloneWriteGroup(full)
	remaining.Revision, remaining.AppliedRevision = "remaining-revision", "remaining-revision"
	remaining.Members = []WriteGroupMember{second}
	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	for _, g := range []*WriteGroup{full, remaining, target} {
		require.NoError(t, f.service.repo.insertVersionInTx(t.Context(), tx, g, req.Boundary))
	}
	cut := &WriteGroupCutover{ID: "partial-fence", OwnershipFenceRevision: "partial-fence", WorkspaceID: f.group.WorkspaceID, Kind: WriteGroupCutoverForward, SourceGroupID: f.group.ID, Boundary: req.Boundary, ReviewDigest: "partial-review", CreatedAt: req.Boundary,
		Ended: []CutoverEnd{{GroupID: full.ID, Revision: full.Revision}}, Started: []CutoverStart{{GroupID: target.ID, Revision: target.Revision}, {GroupID: remaining.ID, Revision: remaining.Revision}}, Moved: []CutoverMemberRef{memberRefOf(full.Members[0])}}
	require.NoError(t, f.service.repo.insertRevisionEndInTx(t.Context(), tx, cut.WorkspaceID, full.ID, full.Revision, cut.ID, cut.Boundary))
	require.NoError(t, f.service.repo.insertCutoverInTx(t.Context(), tx, cut))
	require.NoError(t, f.service.repo.insertOwnershipFenceInTx(t.Context(), tx, cut))
	require.NoError(t, tx.Commit())
	var total, retainedEnded, retainedStarted, movedRemaining int
	require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision='partial-fence'`).Scan(&total))
	require.Equal(t, 4, total)
	require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision='partial-fence' AND tag_id='retained-tag' AND direction='ended'`).Scan(&retainedEnded))
	require.Equal(t, 1, retainedEnded)
	require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision='partial-fence' AND tag_id='retained-tag' AND group_revision='remaining-revision' AND direction='started'`).Scan(&retainedStarted))
	require.Equal(t, 1, retainedStarted)
	require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision='partial-fence' AND tag_id=? AND group_revision='remaining-revision'`, full.Members[0].TagID).Scan(&movedRemaining))
	require.Zero(t, movedRemaining)
}

func TestCutoverOwnershipFenceSchemaKeepsIndependentKinds(t *testing.T) {
	f, _, req := counterCutoverFixture(t)
	cut, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	copyScope := func(kind, direction string) error {
		_, err := tx.Exec(`INSERT INTO write_group_ownership_fences
   (fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name)
   SELECT fence_revision,workspace_id,?,'independent-owner','independent-revision',device_id,point_id,tag_id,?,connector_id,connector_revision,database_name,table_schema,table_name
   FROM write_group_ownership_fences WHERE fence_revision=? AND direction='started' AND group_id<>'independent-owner'`, direction, kind, cut.ID)
		return err
	}
	require.Error(t, copyScope("instantaneous", "started"), "same-kind duplicate ownership is rejected even for a distinct group")
	require.Error(t, copyScope("invalid", "started"), "unknown kind is not coerced to a supported kind")
	require.Error(t, copyScope("cumulative", "invalid"), "unknown direction is not an owner")
	require.NoError(t, copyScope("cumulative", "started"), "independent kind must not be deduplicated")
	var count int
	require.NoError(t, tx.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision=? AND direction='started'`, cut.ID).Scan(&count))
	require.Equal(t, 2, count)
}

func TestCutoverOwnershipFenceUpgradeDoesNotInventLegacyAuthority(t *testing.T) {
	f, _, req := counterCutoverFixture(t)
	cut, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	// Model an old installation's unchanged intent before the additive upgrade.
	_, err = f.db.Exec(`DROP TABLE write_group_ownership_fences`)
	require.NoError(t, err)
	legacyBytes, err := json.Marshal(cutoverDetail{Ended: cut.Ended, Started: cut.Started, Moved: cut.Moved})
	require.NoError(t, err)
	legacy := string(legacyBytes)
	_, err = f.db.Exec(`UPDATE write_group_cutovers SET detail=? WHERE id=?`, legacy, cut.ID)
	require.NoError(t, err)
	for range 2 {
		require.NoError(t, datalink.NewMigrator().Migrate(f.db))
	}
	var detail string
	require.NoError(t, f.db.QueryRow(`SELECT detail FROM write_group_cutovers WHERE id=?`, cut.ID).Scan(&detail))
	require.Equal(t, legacy, detail)
	saved, err := f.service.GetCutover(t.Context(), f.group.ID, cut.ID)
	require.NoError(t, err)
	require.Empty(t, saved.OwnershipFenceRevision)
	var count int
	require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_ownership_fences`).Scan(&count))
	require.Zero(t, count, "unknown legacy ownership must not be fabricated")
}
