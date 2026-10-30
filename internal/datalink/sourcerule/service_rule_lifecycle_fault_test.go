package sourcerule

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"testing"
)

type committedRuleCreateError struct{ *SQLRepository }

func (r *committedRuleCreateError) Create(ctx context.Context, p *schema.SourceRule) error {
	if err := r.SQLRepository.Create(ctx, p); err != nil {
		return err
	}
	return errors.New("after committed create")
}
func TestRuleCreateCommittedErrorFailsClosed(t *testing.T) {
	db := setupSQLRepoDB(t)
	defer db.Close()
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(t.Context(), dr, "d"))
	repo := &committedRuleCreateError{NewSQLRepository(db)}
	svc := NewService(repo, device.NewService(dr, nil), point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db)), nil)
	fc := &rollbackFailCloser{}
	svc.SetShareRuntimeReconciler(fc)
	_, err := svc.Create(t.Context(), CreateRuleRequest{ID: "owned", DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "p"})
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
	_, err = repo.GetByID(t.Context(), "owned")
	require.NoError(t, err)
}

type committedRuleDeleteError struct {
	*SQLRepository
	fail bool
}

func (r *committedRuleDeleteError) Delete(ctx context.Context, id string) error {
	err := r.SQLRepository.Delete(ctx, id)
	if err == nil && r.fail {
		r.fail = false
		return errors.New("after committed delete")
	}
	return err
}
func TestRuleDeleteCommittedErrorDoesNotGuessOwnedAbsence(t *testing.T) {
	db := setupSQLRepoDB(t)
	defer db.Close()
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(t.Context(), dr, "d"))
	repo := &committedRuleDeleteError{SQLRepository: NewSQLRepository(db)}
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	runtime := &formatRevisionRuntime{installed: map[string]string{}}
	svc := NewService(repo, device.NewService(dr, nil), ps, runtime)
	fc := &rollbackFailCloser{}
	svc.SetShareRuntimeReconciler(fc)
	rule, err := svc.Create(t.Context(), CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "p"})
	require.NoError(t, err)
	repo.fail = true
	_, err = svc.DeleteWithRuntimeReconcile(t.Context(), rule.ID)
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
	_, err = repo.GetByID(t.Context(), rule.ID)
	require.ErrorIs(t, err, ErrSourceRuleNotFound)
}

func (r *committedRuleDeleteError) DeleteRevision(ctx context.Context, p *schema.SourceRule) error {
	err := r.SQLRepository.DeleteRevision(ctx, p)
	if err == nil && r.fail {
		r.fail = false
		return errors.New("after committed delete")
	}
	return err
}

func TestRuleConditionalMutationRejectsSameRevisionForeignToggle(t *testing.T) {
	for _, operation := range []string{"update", "delete"} {
		t.Run(operation, func(t *testing.T) {
			db := setupSQLRepoDB(t)
			defer db.Close()
			dr := device.NewSQLRepository(db)
			require.NoError(t, seedSQLSourceRuleDevice(t.Context(), dr, "d"))
			repo := NewSQLRepository(db)
			svc := NewService(repo, device.NewService(dr, nil), point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db)), nil)
			before, err := svc.Create(t.Context(), CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "p"})
			require.NoError(t, err)
			before, err = repo.GetByID(t.Context(), before.ID)
			require.NoError(t, err)
			foreign := cloneRule(before)
			foreign.Enabled = !before.Enabled
			require.NoError(t, repo.Update(t.Context(), foreign))
			foreign, err = repo.GetByID(t.Context(), before.ID)
			require.NoError(t, err)
			// Model SQLite's second-resolution collision without a wall-clock race.
			before.UpdatedAt = foreign.UpdatedAt
			if operation == "delete" {
				err = repo.DeleteRevision(t.Context(), before)
			} else {
				next := cloneRule(before)
				next.NamingPrefix = "stale"
				err = repo.UpdateRevision(t.Context(), next, before)
			}
			require.ErrorIs(t, err, ErrRuleRevisionConflict)
			current, err := repo.GetByID(t.Context(), before.ID)
			require.NoError(t, err)
			require.Equal(t, foreign.Enabled, current.Enabled)
			links, err := repo.ListLinks(t.Context(), before.ID)
			require.NoError(t, err)
			require.Len(t, links, 1)
		})
	}
}

func TestRuleSnapshotDoesNotRecreateUnownedAbsence(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil, nil, nil)
	original := &schema.SourceRule{ID: "gone", RevisionID: "old"}
	ctx, release := svc.AcquireRuleMutation(t.Context(), original.ID)
	defer release()
	svc.noteRuleBaseline(ctx, original)
	err := svc.restoreRuleSnapshot(ctx, &sourceRuleSnapshot{rule: original, pointWrites: svc.pointJournal(ctx)})
	require.ErrorIs(t, err, ErrRuleRevisionConflict)
	_, err = repo.GetByID(ctx, original.ID)
	require.ErrorIs(t, err, ErrSourceRuleNotFound)
}
