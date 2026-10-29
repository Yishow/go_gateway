package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

// Unlike a raw legacy source, this source was counter-bound at its first
// public Apply. A second real destination makes the I target a separate legal
// ownership scope; no saved/applied snapshot is rewritten to manufacture it.
func newConvergingScopeAPIRig(t *testing.T) *scopeAPIRig {
	t.Helper()
	f := newManagedDestinationFixture(t, "sqlite", "uint64")
	r := &cutoverRig{t: t, f: f, kind: "sqlite", anchor: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)}
	r.clock.Store(r.at(-time.Second).UnixNano())
	f.groups.WithClock(func() time.Time { return time.Unix(0, r.clock.Load()).UTC() })
	definition := map[string]any{"id": "meas-energy", "workspace_id": f.record.ID, "device_id": "device-A", "point_id": "point-A", "tag_id": "tag-A", "equipment_id": "main", "name": "Energy", "quantity": "energy", "unit": "kWh", "semantic_kind": "counter", "numeric_encoding": "uint64", "definition_revision": "def-1", "source_binding_revision": "bind-1", "series_epoch": "epoch-1"}
	response := performJSONRequest(t, f.router, http.MethodPost, measurementsPath, definition)
	require.Equal(t, http.StatusCreated, response.Code, response.Body.String())
	response = performJSONRequest(t, f.router, http.MethodPost, "/api/v1/datalink/db-targets/connectors", map[string]any{
		"name": "convergence-target", "kind": "sqlite", "enabled": true,
		"connection_config": map[string]any{"path": filepath.Join(t.TempDir(), "target.db"), "schema": "main"},
	})
	require.Equal(t, http.StatusCreated, response.Code, response.Body.String())
	connector := decodeJSONBody(t, response)["data"].(map[string]any)
	require.NotEqual(t, "connector-A", connector["id"])
	create := func(kind, circuit, layout, connectorID, connectorRevision string) map[string]any {
		interval, unit := 15, "second"
		if kind == "cumulative" {
			interval, unit = 1, "minute"
		}
		request := profileRequest(t, f.writeGroupRouterFixture, kind, interval, unit, circuit, "")
		request["expected_connector_revision"] = connectorRevision
		group := request["group"].(map[string]any)
		group["layout_version"] = layout
		if layout == "managed_long_v1" {
			group["storage_mode"] = "plain"
		}
		member := group["members"].([]map[string]any)[0]
		member["measurement_id"] = "meas-energy"
		if kind == "instantaneous" {
			member["value_role"] = "counter_raw"
		}
		destination := destinationOf(group)
		destination["connector_id"], destination["connector_revision"], destination["table_schema"] = connectorID, connectorRevision, "main"
		response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		_, saved := groupSaveData(t, decodeJSONBody(t, response))
		convergenceAPIConfirm(t, f, saved["id"].(string), connectorRevision)
		return getGroup(t, f, saved["id"].(string))
	}
	r.wide = create("instantaneous", "source", "managed_wide_profile_v1", "connector-A", "connector-1")
	wideID := r.wide["id"].(string)
	response = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+wideID+"/apply", scopeBody(t, f, wideID))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	r.wide = getGroup(t, f, wideID)
	for _, kind := range []string{"instantaneous", "cumulative"} {
		r.targets = append(r.targets, create(kind, "main", "managed_long_v1", connector["id"].(string), connector["identity_revision"].(string)))
	}
	r.long = r.targets[0]
	require.NotEqual(t, destinationOf(r.wide)["database"], destinationOf(r.long)["database"], "ownership scopes use genuinely different SQLite files")
	snap, err := f.groups.ResolveAppliedAt(t.Context(), wideID, r.at(time.Hour))
	require.NoError(t, err)
	r.wideEffective, r.member = snap.EffectiveAt, snap.Group.Members[0]
	require.Equal(t, r.anchor, r.wideEffective)
	r.dest = (&recoveryRig{t: t, f: f, kind: "sqlite"}).openDestination()
	return wireScopeAPIRig(t, r)
}

// Existing single-connector helpers intentionally pin connector-1. This local
// wrapper carries the second connector's actual public identity revision.
func convergenceAPIConfirm(t *testing.T, f managedGroupFixture, id, revision string) {
	t.Helper()
	scope := func() map[string]any {
		body := scopeBody(t, f, id)
		body["expected_connector_revision"] = revision
		return body
	}
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/schema-preview", scope())
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var preview struct {
		Data recordingplan.SchemaPreviewToken
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &preview))
	body := scope()
	body["token"], body["operation_id"] = preview.Data.Token, preview.Data.OperationID
	response = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/schema-apply", body)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var applied struct{ Data recordingplan.SchemaOperation }
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &applied))
	require.Equal(t, recordingplan.SchemaOperationSucceeded, applied.Data.Status)
}

func TestCutoverPublicCompatibleFanoutRollbackConvergesOneOrigin(t *testing.T) {
	r := newConvergingScopeAPIRig(t)
	wideID := r.wide["id"].(string)
	key := groupdelivery.GroupKey{WorkspaceID: r.f.record.ID, GroupID: wideID, GroupRevision: r.wide["applied_revision"].(string)}
	_, originalVersion := counterReplayAPISnapshot(t, r.store, key)
	originalApplied := r.versionPayload(key.GroupID, key.GroupRevision)
	r.advance(r.at(8 * time.Second))
	require.NoError(t, r.offer("future-25", r.member, r.at(25*time.Second), uint64(100)))
	originals := partialAPIOriginals(t, r)
	require.Len(t, originals, 1)
	r.advance(r.at(9 * time.Second))
	cutover := r.cutover(r.at(10 * time.Second))
	require.Len(t, cutover["started"], 2)
	r.reconcile()
	require.Equal(t, 3, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples`))
	require.Equal(t, 2, r.localCount(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE state='accepted'`))
	require.Equal(t, 2, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples WHERE acceptance_kind='handoff' AND origin_group_id=? AND origin_revision=? AND origin_sample_id='future-25'`, key.GroupID, key.GroupRevision))
	convergenceAPIExactCopies(t, r, 2)
	versions := map[groupdelivery.GroupKey][]byte{}
	appliedVersions := map[groupdelivery.GroupKey]string{}
	for _, group := range r.targets {
		target := groupdelivery.GroupKey{WorkspaceID: key.WorkspaceID, GroupID: group["id"].(string), GroupRevision: appliedRevisionOf(t, r.f, group["id"].(string))}
		_, versions[target] = counterReplayAPISnapshot(t, r.store, target)
		appliedVersions[target] = r.versionPayload(target.GroupID, target.GroupRevision)
	}
	convergenceAPIState(t, r, "forward fanout accepted")
	r.advance(r.at(19 * time.Second))
	review := r.rollbackReview(cutover["id"].(string), r.at(20*time.Second))
	require.Equal(t, true, review["applicable"], "%v", review)
	code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.rollbackBody(cutover["id"].(string), r.at(20*time.Second), map[string]any{"review_digest": review["digest"]}))
	require.Equal(t, http.StatusOK, code, "%v", body)
	require.Len(t, body["data"].(map[string]any)["ended"], 2)
	restored := groupdelivery.GroupKey{WorkspaceID: key.WorkspaceID, GroupID: wideID, GroupRevision: appliedRevisionOf(t, r.f, wideID)}
	require.NotEqual(t, key, restored)
	owner, err := r.store.ResolveCutoverOwnership(t.Context(), groupdelivery.OwnershipSource{WorkspaceID: key.WorkspaceID, DeviceID: r.member.DeviceID, PointID: r.member.PointID, TagID: r.member.TagID}, r.at(25*time.Second))
	require.NoError(t, err)
	require.Equal(t, map[groupdelivery.GroupKey]bool{restored: true}, owner.Owners)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	err = r.pipe.Reconcile(ctx)
	convergenceAPIState(t, r, "rollback reconcile result: "+errorText(err))
	require.NoError(t, err, "two compatible paths of one original must converge into one restored owner")
	require.Empty(t, r.reports)
	require.Equal(t, originals, partialAPIOriginals(t, r))
	_, unchanged := counterReplayAPISnapshot(t, r.store, key)
	require.Equal(t, originalVersion, unchanged)
	require.Equal(t, originalApplied, r.versionPayload(key.GroupID, key.GroupRevision))
	for target, before := range versions {
		_, after := counterReplayAPISnapshot(t, r.store, target)
		require.Equal(t, before, after)
		require.Equal(t, appliedVersions[target], r.versionPayload(target.GroupID, target.GroupRevision))
	}
	require.Equal(t, 4, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples`), "one original, two former profile copies, one restored copy")
	require.Equal(t, 1, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=? AND group_revision=? AND origin_sample_id='future-25'`, restored.GroupID, restored.GroupRevision))
	convergenceAPIExactCopies(t, r, 3)
	convergenceAPIEdges(t, r, restored)
	accepted, err := r.store.Restore(t.Context(), restored)
	require.NoError(t, err)
	require.Len(t, accepted.Samples, 1)
	sample := accepted.Samples[0]
	beforeRestart := finalizedAPIState(t, r.cutoverRig)
	handoffs := convergenceAPIReceiptBytes(t, r)
	require.Len(t, handoffs, 4, "no hidden pending or quarantined disposition may remain")
	r = convergenceAPIReopen(t, r)
	require.Equal(t, beforeRestart, finalizedAPIState(t, r.cutoverRig), "physical reopen/reconcile cannot append or consume input")
	require.Equal(t, handoffs, convergenceAPIReceiptBytes(t, r))
	replayed, err := r.store.Restore(t.Context(), restored)
	require.NoError(t, err)
	require.Equal(t, accepted, replayed, "one canonical acceptance survives a cold SQL handle and pipeline")
	report := r.flush(r.at(31 * time.Second))
	convergenceAPIState(t, r, "after closure and delivery")
	// Long instantaneous emits its two empty historical window markers; those
	// are not consumers of obs25. Only restored wide contains the observation.
	require.Equal(t, 3, report.Delivered)
	profileHistory := convergenceAPIEmptyProfileRows(t, r)
	require.Equal(t, 1, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed=1`))
	require.Equal(t, 1, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed=1 AND group_id=? AND group_revision=?`, restored.GroupID, restored.GroupRevision))
	rows := r.wideRowsInOrder()
	require.Len(t, rows, 1)
	require.Equal(t, r.at(20*time.Second), rows[0].Start)
	require.NotNil(t, rows[0].First)
	require.Equal(t, float64(100), *rows[0].First)
	var entries []map[string]any
	require.NoError(t, json.Unmarshal([]byte(rows[0].Provenance), &entries))
	require.Len(t, entries, 1)
	require.Equal(t, "future-25", entries[0]["sample_id"])
	require.Equal(t, r.at(25*time.Second), storedInstant(t, entries[0]["observed_at"].(string)))
	window := windowOfWide(t, rows[0].Provenance)
	require.Equal(t, r.at(20*time.Second), storedInstant(t, window["start"].(string)))
	require.Equal(t, r.at(30*time.Second), storedInstant(t, window["end"].(string)))
	before, frozen, receipts := barrierAPIDeliveryState(t, r.cutoverRig), r.frozenOutbox(), r.receipts()
	r.reconcile()
	require.Zero(t, r.flush(r.at(31*time.Second)).Delivered)
	require.Equal(t, before, barrierAPIDeliveryState(t, r.cutoverRig))
	require.Equal(t, frozen, r.frozenOutbox())
	require.Equal(t, receipts, r.receipts())
	require.Equal(t, profileHistory, convergenceAPIEmptyProfileRows(t, r))
	require.Equal(t, handoffs, convergenceAPIReceiptBytes(t, r), "repeat cannot change either edge's transfer identity")

	// Normal reclamation removes payloads, not the immutable transfer graph or
	// already committed destination SQL. Reopen again to exclude cached proofs.
	closed, err := r.store.Restore(t.Context(), restored)
	require.NoError(t, err)
	require.Empty(t, closed.Samples)
	removed, err := r.store.Reclaim(t.Context(), time.Hour)
	require.NoError(t, err)
	require.EqualValues(t, 4, removed)
	require.Zero(t, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples`))
	require.Equal(t, 4, r.localCount(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE state='finalized'`))
	compactedReceipts := convergenceAPIReceiptBytes(t, r)
	require.Len(t, compactedReceipts, len(handoffs))
	for i, record := range compactedReceipts {
		var before, after map[string]any
		require.NoError(t, json.Unmarshal([]byte(handoffs[i]), &before))
		require.NoError(t, json.Unmarshal([]byte(record), &after))
		require.Equal(t, "finalized", after["State"])
		after["State"] = "accepted"
		require.Equal(t, before, after, "normal GC may compact state, never rewrite transfer identity or proof evidence")
	}
	r = convergenceAPIReopen(t, r)
	compacted, err := r.store.Restore(t.Context(), restored)
	require.NoError(t, err)
	require.Equal(t, closed, compacted, "GC alone cannot stale ownership or handoff authority")
	recognized, err := r.store.RecognizeTransferredSample(t.Context(), restored, sample)
	require.NoError(t, err)
	require.True(t, recognized.Found)
	require.True(t, recognized.Finalized)
	require.Nil(t, recognized.Sample, "no payload may be synthesized after normal reclaim")
	require.Zero(t, r.flush(r.at(31*time.Second)).Delivered)
	require.Equal(t, rows, r.wideRowsInOrder())
	require.Equal(t, compactedReceipts, convergenceAPIReceiptBytes(t, r), "cold recognition and repeated delivery cannot rewrite compact audit")
	require.Equal(t, frozen, r.frozenOutbox())
	require.Equal(t, receipts, r.receipts())
	require.Equal(t, profileHistory, convergenceAPIEmptyProfileRows(t, r))
	for target, before := range versions {
		_, after := counterReplayAPISnapshot(t, r.store, target)
		require.Equal(t, before, after)
		require.Equal(t, appliedVersions[target], r.versionPayload(target.GroupID, target.GroupRevision))
	}
	_, unchanged = counterReplayAPISnapshot(t, r.store, key)
	require.Equal(t, originalVersion, unchanged)
	require.Equal(t, originalApplied, r.versionPayload(key.GroupID, key.GroupRevision))
}

func convergenceAPIEmptyProfileRows(t *testing.T, r *scopeAPIRig) []string {
	t.Helper()
	destination := destinationOf(r.long)
	opened, err := r.connectors.OpenDestination(t.Context(), destination["connector_id"].(string), destination["connector_revision"].(string))
	require.NoError(t, err)
	defer opened.Close()
	table := r.quoted(destination["table_name"].(string))
	var total, missing int
	require.NoError(t, opened.DB.QueryRow(`SELECT COUNT(*),SUM(CASE WHEN profile_kind='instantaneous' AND status='missing'
		AND observed_at IS NULL AND value_int IS NULL AND value_decimal IS NULL AND value_float IS NULL AND value_bool IS NULL AND value_text IS NULL
		AND json_extract(provenance,'$.reason')='no-samples' THEN 1 ELSE 0 END) FROM `+table).Scan(&total, &missing))
	require.Equal(t, []int{2, 2}, []int{total, missing}, "neither former profile may emit the transferred observation")
	rows, err := opened.DB.Query(`SELECT * FROM ` + table + ` ORDER BY bucket_start,record_id`)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	var records []string
	for rows.Next() {
		values, pointers := make([]any, len(columns)), make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		require.NoError(t, rows.Scan(pointers...))
		record := map[string]any{}
		for i, name := range columns {
			record[name] = values[i]
		}
		index := time.Duration(len(records))
		require.Equal(t, r.at((10+index*5)*time.Second), storedInstant(t, record["bucket_start"].(string)))
		require.Equal(t, r.at((15+index*5)*time.Second), storedInstant(t, record["bucket_end"].(string)))
		raw, err := json.Marshal(record)
		require.NoError(t, err)
		records = append(records, string(raw))
	}
	require.NoError(t, rows.Err())
	return records
}

func convergenceAPIReopen(t *testing.T, r *scopeAPIRig) *scopeAPIRig {
	t.Helper()
	require.NoError(t, r.pipe.Stop(time.Second))
	path := sqliteConfigPath(t, r.f.db)
	require.NoError(t, r.f.db.Close())
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	r.f.db = db
	r.f.workspace = workspace.NewService(workspace.NewSQLRepository(db))
	connectors := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(db))
	r.f.groups = workspace.NewWriteGroupService(r.f.workspace, workspace.NewSQLWriteGroupRepository(db)).WithTableInspector(dbtarget.NewReadOnlyTableInspector(connectors))
	r.f.groups.WithManagedTableInspector(r.f.groups.ManagedTableInspector(connectors))
	// The cold four-scope graph exceeded 5s under race (retained diagnostic19).
	// Keep that specific rebuild bounded without changing ordinary/4.1 budgets.
	return wireScopeAPIWithReconcileBudget(t, r.cutoverRig, 15*time.Second)
}

func convergenceAPIEdges(t *testing.T, r *scopeAPIRig, restored groupdelivery.GroupKey) {
	t.Helper()
	require.Equal(t, 4, r.localCount(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE state='accepted'`))
	var edges, targets, origins, sources, canonicals, aliases int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*),COUNT(DISTINCT h.target_seq),COUNT(DISTINCT h.origin_seq),COUNT(DISTINCT h.source_seq),
		SUM(CASE WHEN h.transfer_id=s.handoff_id AND COALESCE(json_extract(h.record,'$.CanonicalTransferID'),'')='' THEN 1 ELSE 0 END),
		SUM(CASE WHEN h.transfer_id<>s.handoff_id AND json_extract(h.record,'$.CanonicalTransferID')=s.handoff_id THEN 1 ELSE 0 END)
		FROM wg_delivery_sample_handoffs h JOIN wg_delivery_samples s ON s.seq=h.target_seq
		WHERE h.target_group_id=? AND h.target_revision=?`, restored.GroupID, restored.GroupRevision).Scan(&edges, &targets, &origins, &sources, &canonicals, &aliases))
	require.Equal(t, []int{2, 1, 1, 2, 1, 1}, []int{edges, targets, origins, sources, canonicals, aliases}, "both independently verified branches retain audit edges to one canonical target acceptance")
}

func convergenceAPIReceiptBytes(t *testing.T, r *scopeAPIRig) []string {
	t.Helper()
	rows, err := r.f.db.Query(`SELECT record FROM wg_delivery_sample_handoffs ORDER BY transfer_id`)
	require.NoError(t, err)
	defer rows.Close()
	var records []string
	for rows.Next() {
		var record string
		require.NoError(t, rows.Scan(&record))
		records = append(records, record)
	}
	require.NoError(t, rows.Err())
	return records
}

func convergenceAPIExactCopies(t *testing.T, r *scopeAPIRig, count int) {
	t.Helper()
	require.Equal(t, count, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples copied JOIN wg_delivery_samples original
		ON original.acceptance_kind='native' AND copied.origin_workspace_id=original.workspace_id
		AND copied.origin_group_id=original.group_id AND copied.origin_revision=original.group_revision AND copied.origin_sample_id=original.sample_id
		WHERE copied.acceptance_kind='handoff' AND copied.payload=original.payload AND copied.payload_digest=original.payload_digest AND copied.created_at=original.created_at`), "copies keep the original payload, proof, digest and acceptance instant")
}

func convergenceAPIState(t *testing.T, r *scopeAPIRig, stage string) {
	t.Helper()
	state := finalizedAPIState(t, r.cutoverRig)
	rows, err := r.f.db.Query(`SELECT record FROM wg_delivery_sample_handoffs ORDER BY source_seq,target_group_id,target_revision`)
	require.NoError(t, err)
	for rows.Next() {
		var record string
		require.NoError(t, rows.Scan(&record))
		state["wg_delivery_sample_handoffs"] = append(state["wg_delivery_sample_handoffs"], record)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	encoded, err := json.Marshal(state)
	require.NoError(t, err)
	t.Logf("%s: %s", stage, encoded)
}

func errorText(err error) string {
	if err == nil {
		return "<nil>"
	}
	return err.Error()
}
