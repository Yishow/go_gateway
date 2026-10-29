package api

import (
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

// Use actual accepted samples, a frozen pending row and a confirmed SQL receipt.
// Stop workers before comparing state so lifecycle, not background delivery,
// is the only writer under test.
func lifecyclePendingRig(t *testing.T) *recoveryRig {
	t.Helper()
	r := newRecoveryRig(t, "sqlite")
	r.start()
	r.record(10.5)
	r.committed(1, "the first accepted value has a real SQL receipt")
	require.NoError(t, r.pipe.Stop(time.Second))
	// Reconcile and close the next window explicitly, without a sender loop.
	// A delayed first delivery cycle must not race the queued-state assertions.
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{
		Groups: r.f.groups, Tags: tag.NewService(tag.NewSQLRepository(r.f.db)), Destinations: targets,
		Inspector: dbtarget.NewReadOnlyTableInspector(targets), Store: groupdelivery.NewStore(r.f.db),
	}, grouppipeline.Config{Now: func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }})
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	r.record(11.5)
	r.pipe.TickAll(t.Context())
	g := r.groups[0]
	member := g.snapshot.Group.Members[0]
	observed := r.base.Add(2*r.window + 2*time.Second)
	r.clock.Store(observed.Add(time.Second).UnixNano())
	require.NoError(t, r.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{
		SampleID: "Q1", WorkspaceID: r.f.record.ID, DeviceID: g.device, PointID: g.point, TagID: g.tag,
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ObservedAt: observed, ReceivedAt: observed.Add(time.Millisecond), TimeOrigin: "gateway",
		Value: 12.5, Quality: schema.QualityGood,
	}))
	require.NoError(t, r.pipe.Stop(time.Second))
	var pending int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed=0`).Scan(&pending))
	require.Equal(t, 1, pending)
	require.Equal(t, 1, r.outbox()["pending"])
	require.Len(t, r.receipts(), 1)
	require.Len(t, r.stored(), 1)
	r.f.groups.WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	return r
}

func lifecycleAcceptedState(t *testing.T, r *recoveryRig) map[string][]string {
	t.Helper()
	state := lifecycleAPIState(t, r.f.db)
	delete(state, "system_settings")
	delete(state, "database_connectors")
	delete(state, "write_groups")
	return state
}

func TestWriteGroupStopPreservesAcceptedDataAfterConnectorChange(t *testing.T) {
	for _, mode := range []string{"changed", "removed"} {
		t.Run(mode, func(t *testing.T) {
			r := lifecyclePendingRig(t)
			id := r.groups[0].id
			request := scopeBody(t, r.f, id)
			query := `UPDATE database_connectors SET identity_revision='connector-2' WHERE id='connector-A'`
			if mode == "removed" {
				query = `DELETE FROM database_connectors WHERE id='connector-A'`
			}
			_, err := r.f.db.Exec(query)
			require.NoError(t, err)
			before, rows, receipts := lifecycleAcceptedState(t, r), r.stored(), r.receipts()
			response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", request)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			_, group := groupSaveData(t, decodeJSONBody(t, response))
			require.Equal(t, "disabled", group["status"])
			require.Equal(t, r.groups[0].snapshot.AppliedRevision, group["applied_revision"])
			eligibility, err := r.f.groups.IntakeEligibility(t.Context(), id)
			require.NoError(t, err)
			require.False(t, eligibility.Accepting)
			require.Equal(t, before, lifecycleAcceptedState(t, r))
			require.Equal(t, rows, r.stored())
			require.Equal(t, receipts, r.receipts())
		})
	}
}

func TestWriteGroupRemoveUnknownAcceptedOwnerIsNoOp(t *testing.T) {
	for _, pending := range []string{"sample", "row"} {
		t.Run(pending, func(t *testing.T) {
			r := lifecyclePendingRig(t)
			id := r.groups[0].id
			query := `UPDATE wg_delivery_samples SET group_revision='R0',origin_revision='R0' WHERE sample_id='Q1'`
			if pending == "row" {
				query = `UPDATE wg_delivery_outbox SET group_revision='R0' WHERE state='pending'`
			}
			_, err := r.f.db.Exec(query)
			require.NoError(t, err)
			request := scopeBody(t, r.f, id)
			before, rows, receipts := lifecycleAPIState(t, r.f.db), r.stored(), r.receipts()
			response := performJSONRequest(t, r.f.router, http.MethodDelete, writeGroupsPath+"/"+id, request)
			require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
			body := decodeJSONBody(t, response)
			assertWriteGroupError(t, body, "WRITE_GROUP_LIFECYCLE_BLOCKED")
			envelope := body["error"].(map[string]any)
			require.Equal(t, "ownership_unavailable", envelope["reason_code"])
			require.Equal(t, "stop", envelope["next_action"])
			require.NotContains(t, response.Body.String(), "R0")
			require.Equal(t, before, lifecycleAPIState(t, r.f.db))
			require.Equal(t, rows, r.stored())
			require.Equal(t, receipts, r.receipts())
			// The offered alternative is usable without erasing the unresolved data.
			accepted := lifecycleAcceptedState(t, r)
			response = performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", request)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			require.Equal(t, accepted, lifecycleAcceptedState(t, r))
		})
	}
}

func TestWriteGroupRemoveProvenBacklogKeepsReadableTombstone(t *testing.T) {
	r := lifecyclePendingRig(t)
	id := r.groups[0].id
	before, rows, receipts := lifecycleAcceptedState(t, r), r.stored(), r.receipts()
	response := performJSONRequest(t, r.f.router, http.MethodDelete, writeGroupsPath+"/"+id, scopeBody(t, r.f, id))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	group := getGroup(t, r.f, id)
	require.Equal(t, "deleted", group["status"])
	require.Equal(t, before, lifecycleAcceptedState(t, r))
	require.Equal(t, rows, r.stored())
	require.Equal(t, receipts, r.receipts())
	// Reopening the production worker still delivers the already-frozen row
	// to its original table; the tombstone is not treated as a cancel request.
	frozen := r.outboxDigests()
	require.NoError(t, r.pipe.Start(t.Context()))
	t.Cleanup(func() { _ = r.pipe.Stop(time.Second) })
	r.committed(2, "the tombstoned group's queued row reaches its frozen destination")
	require.NoError(t, r.pipe.Stop(time.Second))
	require.Equal(t, frozen, r.outboxDigests())
	require.Len(t, r.stored(), 2)
	for key, receipt := range receipts {
		require.Equal(t, receipt, r.receipts()[key])
	}
}
