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
	"testing"
)

type committedPointError struct {
	*point.SQLRepository
	operation string
	cancel    context.CancelFunc
}

func (r *committedPointError) finish(operation string, err error) error {
	if err == nil && r.operation == operation {
		r.operation = ""
		if r.cancel != nil {
			r.cancel()
		}
		return errors.New("storage returned error after commit")
	}
	return err
}
func (r *committedPointError) Create(ctx context.Context, p *schema.Point) error {
	return r.finish("create", r.SQLRepository.Create(ctx, p))
}
func (r *committedPointError) UpdateRevision(ctx context.Context, p *schema.Point, expected string) error {
	return r.finish("update", r.SQLRepository.UpdateRevision(ctx, p, expected))
}
func (r *committedPointError) DeleteRevision(ctx context.Context, id, deviceID, expected string) error {
	return r.finish("delete", r.SQLRepository.DeleteRevision(ctx, id, deviceID, expected))
}

type formatRevisionRuntime struct {
	leases    map[string]uint64
	nextLease uint64
	installed map[string]string
	suspended []string
}

func (r *formatRevisionRuntime) UpsertPoint(p *schema.Point) {
	r.installed[p.ID] = p.RevisionID
	if r.leases == nil {
		r.leases = map[string]uint64{}
	}
	r.nextLease++
	r.leases[p.ID] = r.nextLease
}
func (r *formatRevisionRuntime) RemovePoint(id string) { delete(r.installed, id) }
func (r *formatRevisionRuntime) SuspendPointRevision(id, revision string) bool {
	current, ok := r.installed[id]
	if !ok {
		return true
	}
	if current != revision {
		return false
	}
	delete(r.installed, id)
	r.suspended = append(r.suspended, id)
	return true
}

func TestRuleFormatCommittedErrorIsUncertainAndStopsOwnedRevision(t *testing.T) {
	for _, operation := range []string{"create", "update", "delete", "update-cancel"} {
		t.Run(operation, func(t *testing.T) {
			base := t.Context()
			ctx, cancel := context.WithCancel(base)
			defer cancel()
			db := setupSQLRepoDB(t)
			t.Cleanup(func() { _ = db.Close() })
			dr := device.NewSQLRepository(db)
			require.NoError(t, seedSQLSourceRuleDevice(base, dr, "d"))
			repo := &committedPointError{SQLRepository: point.NewSQLRepository(db)}
			ps := point.NewService(repo, pollinggroup.NewSQLRepository(db))
			syncer := &formatRevisionRuntime{installed: map[string]string{}}
			svc := NewService(NewSQLRepository(db), device.NewService(dr, nil), ps, syncer)
			fc := &rollbackFailCloser{}
			svc.SetShareRuntimeReconciler(fc)
			rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "BEFORE"})
			require.NoError(t, err)
			links, err := svc.ListLinks(ctx, rule.ID)
			require.NoError(t, err)
			original, err := ps.GetByID(ctx, links[0].PointID)
			require.NoError(t, err)
			syncer.UpsertPoint(original)
			repo.operation = operation
			if operation == "update-cancel" {
				repo.operation = "update"
				repo.cancel = cancel
			}
			switch operation {
			case "create":
				_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{Count: common.Ptr(2)})
			case "delete":
				err = svc.Delete(ctx, rule.ID)
			default:
				_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("AFTER")})
			}
			var typed *modbusshare.Error
			require.ErrorAs(t, err, &typed)
			require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
			require.Positive(t, fc.failClosedCalls)
			if operation != "create" {
				require.NotContains(t, syncer.installed, original.ID)
				require.Contains(t, syncer.suspended, original.ID)
			}
		})
	}
}

func (r *formatRevisionRuntime) CapturePointLease(id, revision string) uint64 {
	if r.installed[id] != revision {
		return 0
	}
	return r.leases[id]
}
func (r *formatRevisionRuntime) InstallPointForMutation(p *schema.Point) uint64 {
	r.UpsertPoint(p)
	return r.leases[p.ID]
}
func (r *formatRevisionRuntime) SuspendPointLease(id string, lease uint64) bool {
	if r.leases[id] != lease {
		return false
	}
	delete(r.installed, id)
	delete(r.leases, id)
	r.suspended = append(r.suspended, id)
	return true
}
func (r *formatRevisionRuntime) RestorePointLease(p *schema.Point, lease uint64) bool {
	if lease == 0 || r.leases[p.ID] != lease {
		return false
	}
	r.UpsertPoint(p)
	return true
}

type committedRuleError struct {
	*SQLRepository
	fail   bool
	cancel context.CancelFunc
}

func (r *committedRuleError) UpdateRevision(ctx context.Context, next, expected *schema.SourceRule) error {
	err := r.SQLRepository.UpdateRevision(ctx, next, expected)
	if err == nil && r.fail {
		r.fail = false
		if r.cancel != nil {
			r.cancel()
		}
		return errors.New("rule storage error after commit")
	}
	return err
}
func TestRuleFormatRuleCommitErrorUsesAttemptedRevision(t *testing.T) {
	for _, cancelAfter := range []bool{false, true} {
		t.Run(map[bool]string{false: "restorable", true: "canceled"}[cancelAfter], func(t *testing.T) {
			base := t.Context()
			ctx, cancel := context.WithCancel(base)
			defer cancel()
			db := setupSQLRepoDB(t)
			t.Cleanup(func() { _ = db.Close() })
			dr := device.NewSQLRepository(db)
			require.NoError(t, seedSQLSourceRuleDevice(base, dr, "d"))
			ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
			repo := &committedRuleError{SQLRepository: NewSQLRepository(db)}
			runtime := &formatRevisionRuntime{installed: map[string]string{}}
			svc := NewService(repo, device.NewService(dr, nil), ps, runtime)
			fc := &rollbackFailCloser{}
			svc.SetShareRuntimeReconciler(fc)
			rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "ORIGINAL"})
			require.NoError(t, err)
			links, err := svc.ListLinks(ctx, rule.ID)
			require.NoError(t, err)
			p, err := ps.GetByID(ctx, links[0].PointID)
			require.NoError(t, err)
			runtime.UpsertPoint(p)
			repo.fail = true
			if cancelAfter {
				repo.cancel = cancel
			}
			_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("CHANGED")})
			require.Error(t, err)
			if cancelAfter {
				var typed *modbusshare.Error
				require.ErrorAs(t, err, &typed)
				require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
				require.NotContains(t, runtime.installed, p.ID)
				require.Positive(t, fc.failClosedCalls)
			} else {
				current, e := repo.GetByID(base, rule.ID)
				require.NoError(t, e)
				require.Equal(t, rule.RevisionID, current.RevisionID)
				require.Equal(t, rule.NamingPrefix, current.NamingPrefix)
			}
		})
	}
}

func (r *formatRevisionRuntime) InstallPointForMutationAtLease(p *schema.Point, expected uint64) (uint64, bool) {
	if r.leases[p.ID] != expected {
		return 0, false
	}
	r.UpsertPoint(p)
	return r.leases[p.ID], true
}
func (r *formatRevisionRuntime) RemovePointOwned(id string, expected uint64) (uint64, bool) {
	if r.leases[id] != expected {
		return 0, false
	}
	delete(r.installed, id)
	delete(r.leases, id)
	r.suspended = append(r.suspended, id)
	return r.nextLease, true
}
func (r *formatRevisionRuntime) RestoreRemovedPointOwned(p *schema.Point, removedLease, epoch uint64) (uint64, bool) {
	if removedLease == 0 || r.leases[p.ID] != 0 || r.nextLease != epoch {
		return 0, false
	}
	r.UpsertPoint(p)
	return r.leases[p.ID], true
}
func (r *formatRevisionRuntime) RestorePointLeaseReceipt(p *schema.Point, lease uint64) (uint64, bool) {
	if !r.RestorePointLease(p, lease) {
		return 0, false
	}
	return r.leases[p.ID], true
}
