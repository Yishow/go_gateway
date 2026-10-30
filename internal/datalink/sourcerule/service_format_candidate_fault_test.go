package sourcerule

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
	"testing"
)

type committedCandidateError struct {
	*SQLRepository
	fail      bool
	attempted string
}

func (r *committedCandidateError) ReplaceCandidateSnapshotsAtRevision(ctx context.Context, snapshots []*schema.SourceRuleCandidateSnapshot, revision string) error {
	err := r.SQLRepository.ReplaceCandidateSnapshotsAtRevision(ctx, snapshots, revision)
	if err == nil && r.fail {
		r.fail = false
		r.attempted = revision
		return errors.New("candidate storage error after commit")
	}
	return err
}
func TestRuleFormatCandidateCommitErrorRestoresWholeConfiguration(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	t.Cleanup(func() { _ = db.Close() })
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	repo := &committedCandidateError{SQLRepository: NewSQLRepository(db)}
	svc := NewService(repo, device.NewService(dr, nil), ps, nil)
	rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "BEFORE"})
	require.NoError(t, err)
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	before, err := ps.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	before.FormatMode = schema.PointFormatOverride
	before.FormatOverride = "CDAB"
	before.DataFormat = "CDAB"
	require.NoError(t, ps.Restore(ctx, before))
	oldCandidates, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	repo.fail = true
	_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("FAILED")})
	require.Error(t, err)
	got, err := ps.GetByID(ctx, before.ID)
	require.NoError(t, err)
	require.Equal(t, before.RevisionID, got.RevisionID)
	require.Equal(t, before.FormatOverride, got.FormatOverride)
	current, err := svc.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	require.Equal(t, rule.RevisionID, current.RevisionID)
	candidates, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	require.Equal(t, oldCandidates, candidates)
	failedCandidates, err := repo.ListCandidateSnapshots(ctx, rule.ID, repo.attempted)
	require.NoError(t, err)
	require.Empty(t, failedCandidates)
}

type uncertainCandidateError struct {
	*SQLRepository
	fail bool
}

func (r *uncertainCandidateError) ReplaceCandidateSnapshotsAtRevision(ctx context.Context, snapshots []*schema.SourceRuleCandidateSnapshot, revision string) error {
	err := r.SQLRepository.ReplaceCandidateSnapshotsAtRevision(ctx, snapshots, revision)
	if err == nil && r.fail {
		r.fail = false
		return errors.Join(sqltx.ErrCommitOutcomeUnknown, errors.New("private database diagnostic"))
	}
	return err
}
func TestRuleUnknownCommitDoesNotAttemptInverseCompensation(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	defer db.Close()
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	repo := &uncertainCandidateError{SQLRepository: NewSQLRepository(db)}
	runtime := &formatRevisionRuntime{installed: map[string]string{}}
	svc := NewService(repo, device.NewService(dr, nil), ps, runtime)
	svc.SetShareRuntimeReconciler(&rollbackFailCloser{})
	rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "BEFORE"})
	require.NoError(t, err)
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	repo.fail = true
	_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("ATTEMPTED")})
	require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
	require.NotContains(t, err.Error(), "private database")
	current, err := svc.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	require.NotEqual(t, rule.RevisionID, current.RevisionID)
	require.Equal(t, "ATTEMPTED", current.NamingPrefix)
	require.NotContains(t, runtime.installed, links[0].PointID)
}

func TestShareUnknownCommitDoesNotAttemptInverseCompensation(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	defer db.Close()
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	repo := &uncertainCandidateError{SQLRepository: NewSQLRepository(db)}
	runtime := &formatRevisionRuntime{installed: map[string]string{}}
	svc := NewService(repo, device.NewService(dr, nil), ps, runtime)
	svc.SetShareRuntimeReconciler(&rollbackFailCloser{})
	rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "BEFORE"})
	require.NoError(t, err)
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	svc.SetShareRuntimeReconciler(ShareRuntimeReconcilerFunc(func(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
		return RuntimeReconcileOutcome{Status: RuntimeReconcileStatusFailed}, errors.Join(sqltx.ErrCommitOutcomeUnknown, errors.New("share commit ack lost"))
	}))
	_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("ATTEMPTED")})
	require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
	require.NotContains(t, err.Error(), "private database")
	current, err := svc.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	require.NotEqual(t, rule.RevisionID, current.RevisionID)
	require.Equal(t, "ATTEMPTED", current.NamingPrefix)
	require.NotContains(t, runtime.installed, links[0].PointID)
}

func TestShareUnknownCommitCreateDoesNotDeleteOwnedConfiguration(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	defer db.Close()
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	runtime := &formatRevisionRuntime{installed: map[string]string{}}
	svc := NewService(NewSQLRepository(db), device.NewService(dr, nil), ps, runtime)
	svc.SetShareRuntimeReconciler(ShareRuntimeReconcilerFunc(func(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
		return RuntimeReconcileOutcome{Status: RuntimeReconcileStatusFailed}, sqltx.ErrCommitOutcomeUnknown
	}))
	created, _, err := svc.CreateWithRuntimeReconcile(ctx, CreateRuleRequest{ID: "create-unknown", DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "ATTEMPTED"})
	require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	require.NotNil(t, created)
	stored, err := svc.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, "ATTEMPTED", stored.NamingPrefix)
	links, err := svc.ListLinks(ctx, stored.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.NotContains(t, runtime.installed, links[0].PointID)
}

func TestShareTypedDirtyUnknownDoesNotAttemptInverseCompensation(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	defer db.Close()
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	repo := &uncertainCandidateError{SQLRepository: NewSQLRepository(db)}
	runtime := &formatRevisionRuntime{installed: map[string]string{}}
	svc := NewService(repo, device.NewService(dr, nil), ps, runtime)
	svc.SetShareRuntimeReconciler(&rollbackFailCloser{})
	rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "BEFORE"})
	require.NoError(t, err)
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	svc.SetShareRuntimeReconciler(ShareRuntimeReconcilerFunc(func(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
		return RuntimeReconcileOutcome{Status: RuntimeReconcileStatusFailed}, &modbusshare.Error{Code: modbusshare.ErrCodeDirtyUnknown, Message: "private typed dirty cause"}
	}))
	_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("ATTEMPTED")})
	require.Error(t, err)
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
	mutationErr := err
	current, err := svc.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	require.NotEqual(t, rule.RevisionID, current.RevisionID)
	require.Equal(t, "ATTEMPTED", current.NamingPrefix)
	require.NotContains(t, runtime.installed, links[0].PointID)
	require.NotContains(t, mutationErr.Error(), "private typed dirty cause")
}
