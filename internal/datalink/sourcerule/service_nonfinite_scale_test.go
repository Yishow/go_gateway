package sourcerule

import (
	"math"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

func TestNonfiniteRuleScaleIsRejectedBeforeAnyMutation(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		t.Run(operation, func(t *testing.T) {
			for _, kind := range []string{"nan-multiplier", "infinite-offset"} {
				t.Run(kind, func(t *testing.T) {
					faults := &lostcommit.Faults{}
					db := lostcommit.OpenFaults(":memory:", faults)
					db.SetMaxOpenConns(1)
					t.Cleanup(func() { require.NoError(t, db.Close()) })
					require.NoError(t, datalink.NewMigrator().Migrate(db))
					dr := device.NewSQLRepository(db)
					require.NoError(t, seedSQLSourceRuleDevice(t.Context(), dr, "device"))
					require.NoError(t, dr.UpdateTestResult(t.Context(), "device", true, ""))
					ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
					repo := NewSQLRepository(db)
					svc := NewService(repo, device.NewService(dr, nil), ps, nil)
					tags := tag.NewService(tag.NewSQLRepository(db))
					svc.SetTagMappingServices(tags, mapping.NewServiceWithTagResolver(mapping.NewSQLRepository(db), tags.GetByID))
					var ruleID string
					if operation == "update" {
						original, err := svc.Create(t.Context(), CreateRuleRequest{ID: "rule", DeviceID: "device", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt16})
						require.NoError(t, err)
						ruleID = original.ID
						links, err := repo.ListLinks(t.Context(), ruleID)
						require.NoError(t, err)
						require.Len(t, links, 1)
						managed, err := tags.Create(t.Context(), tag.CreateTagRequest{Key: "managed", DataType: schema.DataTypeInt16, Labels: map[string]string{ruleManagedTagLabelSource: ruleManagedTagLabelValue, ruleManagedTagLabelRuleID: ruleID, ruleManagedTagLabelAddress: normalizeAddressKey(links[0].Address)}})
						require.NoError(t, err)
						links[0].TagID = &managed.ID
						require.NoError(t, repo.ReplaceLinks(t.Context(), ruleID, links))
						require.NoError(t, ps.Delete(t.Context(), links[0].PointID))
					}
					runtime := &acquisitionMutationRecorder{reconcileRuntimeSync: &reconcileRuntimeSync{}}
					svc.runtimeSync = runtime
					writes := make([]string, 0)
					observe := func(query string) {
						fields := strings.Fields(strings.ToUpper(query))
						if len(fields) == 0 {
							return
						}
						switch fields[0] {
						case "INSERT", "UPDATE", "DELETE", "REPLACE":
							writes = append(writes, query)
						}
					}
					faults.ObserveStatement.Store(&observe)
					nan, inf := math.NaN(), math.Inf(1)
					var multiplier, offset *float64
					if kind == "nan-multiplier" {
						multiplier = &nan
					} else {
						offset = &inf
					}
					var err error
					if operation == "create" {
						_, err = svc.Create(t.Context(), CreateRuleRequest{ID: "new-rule", DeviceID: "device", StartAddress: "40011", Count: 1, DataType: schema.DataTypeInt16, Enabled: true, ScaleMultiplier: multiplier, ScaleOffset: offset})
					} else {
						enabled := true
						_, err = svc.Update(t.Context(), ruleID, UpdateRuleRequest{Enabled: &enabled, ScaleMultiplier: multiplier, ScaleOffset: offset})
					}
					require.ErrorIs(t, err, ErrValidation)
					field := "scale_multiplier"
					if kind == "infinite-offset" {
						field = "scale_offset"
					}
					require.ErrorContains(t, err, field+" must be finite")
					require.NotContains(t, err.Error(), "NaN")
					require.NotContains(t, err.Error(), "+Inf")
					require.Empty(t, writes, "no SQL mutation to device/rule/point/tag/mapping/link tables")
					require.Zero(t, runtime.begun, "reject before native acquisition invalidation")
					require.Zero(t, runtime.finished)
					require.Empty(t, runtime.upserted)
					require.Empty(t, runtime.removed)
				})
			}
		})
	}
}

func TestUpdateRejectsInvalidStoredRuleBeforeAcquisition(t *testing.T) {
	dr := device.NewMemoryRepository()
	_, err := seedActiveDevice(t.Context(), dr, "device")
	require.NoError(t, err)
	repo := &invalidRuleWriteSpy{MemoryRepository: NewMemoryRepository()}
	invalid := &schema.SourceRule{ID: "invalid-stored", DeviceID: "device", RevisionID: "revision", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt16, CreatedAt: time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)}
	require.NoError(t, repo.MemoryRepository.Create(t.Context(), invalid))
	runtime := &acquisitionMutationRecorder{reconcileRuntimeSync: &reconcileRuntimeSync{}}
	svc := NewService(repo, device.NewService(dr, nil), point.NewService(point.NewMemoryRepository(), nil), runtime)
	_, err = svc.Update(t.Context(), invalid.ID, UpdateRuleRequest{})
	require.ErrorIs(t, err, ErrRuleRevisionConflict)
	require.Zero(t, runtime.begun)
	require.Zero(t, runtime.finished)
	require.Empty(t, runtime.upserted)
	require.Empty(t, runtime.removed)
	require.Zero(t, repo.creates)
	require.Zero(t, repo.updates)
	require.Zero(t, repo.deletes)
	require.Zero(t, repo.links)
}
