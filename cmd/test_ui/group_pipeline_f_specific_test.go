//go:build f_write_group_fixture

package main

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type fProductionFixture struct {
	target *lifecycleTarget
	env    *outageEnv
	clock  *testClock
	faults *fixtureFaultController
	group  *workspace.WriteGroup
	pipe   *grouppipeline.Pipeline
}

type fOutboxRow struct {
	EffectKey     string
	RecordID      string
	Revision      string
	Payload       string
	PayloadDigest string
	State         string
	RetryCount    int
	LastError     string
}

func newFProductionFixture(t *testing.T, kind, owner string) *fProductionFixture {
	t.Helper()
	target := newLifecycleTarget(t, kind)
	env := target.env
	if kind == "sqlite" {
		file := env.targets["connector-A"]
		destination, err := sql.Open("sqlite", file+"?_pragma=busy_timeout(5000)")
		require.NoError(t, err)
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), destination, schema.DatabaseConnectorKindSQLite, ""))
		require.NoError(t, destination.Close())
		env.dedupe = groupdelivery.DedupeReceipt
	}
	clock := &testClock{}
	clock.set(lifecycleBase.Add(-time.Second))
	env.services.writeGroups.WithClock(clock.now)
	group := env.createAndApply(t, "A")
	clock.set(lifecycleBase.Add(time.Second))
	faults := newFixtureFaultController(env.db, &env.services, clock.now)
	require.NoError(t, faults.install())
	t.Cleanup(func() { require.NoError(t, faults.close(context.Background())) })
	return &fProductionFixture{
		target: target,
		env:    env,
		clock:  clock,
		faults: faults,
		group:  group,
		pipe:   env.pipeline(clock, owner),
	}
}

func (f *fProductionFixture) start(t *testing.T) {
	t.Helper()
	require.NoError(t, f.pipe.Start(t.Context()))
	t.Cleanup(func() { require.NoError(t, f.pipe.Stop(3*time.Second)) })
}

func (f *fProductionFixture) configureFault(t *testing.T, kind string, enabled bool) {
	t.Helper()
	_, err := f.faults.configure(t.Context(), fixtureFaultRequest{Kind: kind, GroupID: f.group.ID, Enabled: enabled})
	require.NoError(t, err)
}

func (f *fProductionFixture) holding() bool {
	for _, state := range f.faults.snapshot() {
		if state.GroupID == f.group.ID {
			return state.Holding
		}
	}
	return false
}

const (
	fAttemptBudget      = 2 * time.Second
	fSettlementBudget   = 200 * time.Millisecond
	fHeldCommitDuration = 600 * time.Millisecond
)

func waitFCommitHold(t *testing.T) {
	t.Helper()
	started := time.Now()
	time.Sleep(fHeldCommitDuration)
	elapsed := time.Since(started)
	require.Greater(t, elapsed, fSettlementBudget, "the held target response must outlast local settlement")
	require.Less(t, elapsed, fAttemptBudget, "the held target response must remain within its attempt budget")
}

func newFSlowPipeline(env *outageEnv, clock *testClock, owner string) *grouppipeline.Pipeline {
	immediate := func(int) (string, time.Time) { return groupdelivery.StateRetrying, clock.now().Add(-time.Hour) }
	return grouppipeline.New(grouppipeline.Dependencies{
		Groups: env.services.writeGroups, Tags: env.services.tag, Destinations: env.services.dbTarget,
		Inspector: dbtarget.NewReadOnlyTableInspector(env.services.dbTarget), Store: groupdelivery.NewStore(env.db),
	}, grouppipeline.Config{
		NodeID: "node-1", Owner: owner, Now: clock.now,
		TickInterval: time.Hour, ReconcileInterval: time.Hour, WorkerInterval: 10 * time.Millisecond,
		Sender: groupdelivery.SenderConfig{
			Backoff: immediate, DeliveryTimeout: fAttemptBudget, SettleTimeout: fSettlementBudget,
		},
	})
}

func acceptFBucket(ctx context.Context, env *outageEnv, pipe *grouppipeline.Pipeline, group *workspace.WriteGroup, observed time.Time, prefix string, temperature float64, pressure int64) error {
	if err := pipe.AcceptSample(ctx, env.envelope(group, 0, prefix+"-temperature", observed, temperature)); err != nil {
		return err
	}
	return pipe.AcceptSample(ctx, env.envelope(group, 1, prefix+"-pressure", observed, pressure))
}

func readFOutbox(ctx context.Context, db *sql.DB, groupID string) ([]fOutboxRow, error) {
	rows, err := db.QueryContext(ctx, `SELECT effect_key,record_id,group_revision,payload,payload_digest,state,retry_count,last_error_code FROM wg_delivery_outbox WHERE group_id = ? ORDER BY bucket_start, effect_key`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []fOutboxRow
	for rows.Next() {
		var row fOutboxRow
		if err := rows.Scan(&row.EffectKey, &row.RecordID, &row.Revision, &row.Payload, &row.PayloadDigest, &row.State, &row.RetryCount, &row.LastError); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func waitFOutbox(t *testing.T, env *outageEnv, groupID, wantState string) fOutboxRow {
	t.Helper()
	var result fOutboxRow
	require.Eventually(t, func() bool {
		rows, err := readFOutbox(t.Context(), env.db, groupID)
		if err != nil || len(rows) != 1 {
			return false
		}
		result = rows[0]
		return result.State == wantState
	}, 10*time.Second, 10*time.Millisecond)
	return result
}

func waitFCommitted(t *testing.T, env *outageEnv, groupID string, count int) []fOutboxRow {
	t.Helper()
	var result []fOutboxRow
	require.Eventually(t, func() bool {
		rows, err := readFOutbox(t.Context(), env.db, groupID)
		if err != nil || len(rows) != count {
			return false
		}
		for _, row := range rows {
			if row.State != groupdelivery.StateCommitted {
				return false
			}
		}
		result = rows
		return true
	}, 10*time.Second, 10*time.Millisecond)
	return result
}

func fLocalReceiptCount(t *testing.T, env *outageEnv, groupID string) int {
	t.Helper()
	var count int
	require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_receipts r JOIN wg_delivery_outbox o USING (effect_key) WHERE o.group_id = ?`, groupID).Scan(&count))
	return count
}

func fAssertTargetReceipts(t *testing.T, env *outageEnv, group *workspace.WriteGroup, effects []fOutboxRow) {
	t.Helper()
	opened, err := env.services.dbTarget.OpenDestination(t.Context(), group.Destination.ConnectorID, group.Destination.ConnectorRevision)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, opened.Close()) })
	table := quoteFIdentifier(group.Destination.TableSchema) + "." + quoteFIdentifier(dbtarget.EffectReceiptTable)
	rows, err := opened.DB.QueryContext(t.Context(), `SELECT effect_key,payload_digest FROM `+table+` ORDER BY effect_key`)
	require.NoError(t, err)
	defer rows.Close()
	var got []fOutboxRow
	for rows.Next() {
		var row fOutboxRow
		require.NoError(t, rows.Scan(&row.EffectKey, &row.PayloadDigest))
		got = append(got, row)
	}
	require.NoError(t, rows.Err())
	require.Len(t, got, len(effects))
	want := make(map[string]fOutboxRow, len(effects))
	for _, effect := range effects {
		want[effect.EffectKey] = effect
	}
	for _, receipt := range got {
		effect, ok := want[receipt.EffectKey]
		require.True(t, ok, "target receipt must match a local outbox effect")
		require.Equal(t, effect.PayloadDigest, receipt.PayloadDigest)
	}
}

func quoteFIdentifier(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

func fPayloadValues(t *testing.T, payload string) (temperature float64, pressure int64) {
	t.Helper()
	row, err := groupdelivery.DecodeRowPayload([]byte(payload))
	require.NoError(t, err)
	for _, cell := range row.Cells {
		switch cell.Column {
		case "temperature":
			var ok bool
			temperature, ok = cell.Value.(float64)
			require.True(t, ok, "temperature payload keeps its exact SQL float type")
		case "pressure":
			var ok bool
			pressure, ok = cell.Value.(int64)
			require.True(t, ok, "pressure payload keeps its exact SQL integer type")
		}
	}
	return temperature, pressure
}

func fIsOldPair(temperature float64, pressure int64) bool {
	return (temperature == 41 && pressure == 410) || (temperature == 51 && pressure == 510)
}

func fAssertCutoffPayloads(t *testing.T, rows []fOutboxRow, action, oldRevision, newRevision string) {
	t.Helper()
	for _, row := range rows {
		temperature, pressure := fPayloadValues(t, row.Payload)
		switch {
		case action == "replace" && row.Revision == newRevision:
			require.Equal(t, float64(61), temperature)
			require.Equal(t, int64(610), pressure)
		default:
			require.Equal(t, oldRevision, row.Revision, "old buckets remain owned by the original applied revision")
			require.True(t, fIsOldPair(temperature, pressure), "old payload keeps one of the accepted pre-cutoff values")
		}
	}
}

func fAssertCutoffTargetRows(t *testing.T, rows [][2]float64, action string) {
	t.Helper()
	for _, row := range rows {
		if action == "replace" && row == [2]float64{61, 610} {
			continue
		}
		require.True(t, fIsOldPair(row[0], int64(row[1])), "target SQL row keeps the frozen pre-cutoff payload")
	}
	if action == "replace" {
		require.Contains(t, rows, [2]float64{61, 610}, "superseding revision writes its own exact payload once")
	}
}

func TestProductionFSlowSenderSuccessAndFailureDurable(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind+"/success", func(t *testing.T) {
			f := newFProductionFixture(t, kind, "node-1/f-slow-success")
			f.pipe = newFSlowPipeline(f.env, f.clock, "node-1/f-slow-success")
			f.start(t)
			f.configureFault(t, string(dbtarget.FixtureDestinationFaultTargetCommitHold), true)
			require.NoError(t, acceptFBucket(t.Context(), f.env, f.pipe, f.group, lifecycleBase.Add(5*time.Second), "f-slow-success", 21, 210))
			f.clock.set(lifecycleBase.Add(11 * time.Second))
			f.pipe.TickAll(t.Context())
			require.Eventually(t, f.holding, 10*time.Second, 10*time.Millisecond, "the real target commit must reach the hold hook")
			waitFOutbox(t, f.env, f.group.ID, groupdelivery.StateSending)
			require.Len(t, f.target.rows(), 1, "target SQL commit is visible while local settlement is held")
			require.Zero(t, fLocalReceiptCount(t, f.env, f.group.ID), "local receipt waits for the post-attempt settlement budget")
			waitFCommitHold(t)
			require.True(t, f.holding(), "the target response remains held beyond both attempt and settlement budgets")
			waitFOutbox(t, f.env, f.group.ID, groupdelivery.StateSending)
			require.Zero(t, fLocalReceiptCount(t, f.env, f.group.ID))

			f.configureFault(t, string(dbtarget.FixtureDestinationFaultTargetCommitHold), false)
			row := waitFOutbox(t, f.env, f.group.ID, groupdelivery.StateCommitted)
			require.Empty(t, row.LastError)
			require.Equal(t, 1, fLocalReceiptCount(t, f.env, f.group.ID))
			require.Len(t, f.target.rows(), 1)
			fAssertTargetReceipts(t, f.env, f.group, []fOutboxRow{row})
		})

		t.Run(kind+"/failure", func(t *testing.T) {
			f := newFProductionFixture(t, kind, "node-1/f-slow-failure")
			f.pipe = newFSlowPipeline(f.env, f.clock, "node-1/f-slow-failure")
			f.start(t)
			f.configureFault(t, string(dbtarget.FixtureDestinationFaultTargetCommitHold), true)
			require.NoError(t, acceptFBucket(t.Context(), f.env, f.pipe, f.group, lifecycleBase.Add(5*time.Second), "f-slow-failure", 31, 310))
			f.clock.set(lifecycleBase.Add(11 * time.Second))
			f.pipe.TickAll(t.Context())
			require.Eventually(t, f.holding, 10*time.Second, 10*time.Millisecond)
			waitFOutbox(t, f.env, f.group.ID, groupdelivery.StateSending)
			require.Len(t, f.target.rows(), 1)
			require.Zero(t, fLocalReceiptCount(t, f.env, f.group.ID))
			waitFCommitHold(t)
			require.True(t, f.holding(), "the target response remains held beyond both attempt and settlement budgets")
			require.Zero(t, fLocalReceiptCount(t, f.env, f.group.ID))

			// Release the slow attempt as an actual lost response. The target row and
			// target receipt already exist, so a retry must settle the local ledger.
			f.configureFault(t, string(dbtarget.FixtureDestinationFaultTargetCommitResponseLost), true)
			require.Eventually(t, func() bool {
				rows, err := readFOutbox(t.Context(), f.env.db, f.group.ID)
				return err == nil && len(rows) == 1 && rows[0].RetryCount >= 1
			}, 10*time.Second, 10*time.Millisecond, "the held attempt must record a durable ambiguous retry")
			f.configureFault(t, string(dbtarget.FixtureDestinationFaultTargetCommitResponseLost), false)
			row := waitFOutbox(t, f.env, f.group.ID, groupdelivery.StateCommitted)
			require.GreaterOrEqual(t, row.RetryCount, 1)
			require.Empty(t, row.LastError)
			require.Equal(t, 1, fLocalReceiptCount(t, f.env, f.group.ID))
			require.Len(t, f.target.rows(), 1, "ambiguous retry must not duplicate the target row")
			fAssertTargetReceipts(t, f.env, f.group, []fOutboxRow{row})
		})
	}
}

type fCutoffResult struct {
	group *workspace.WriteGroup
	err   error
}

func runFCutoffAction(ctx context.Context, env *outageEnv, original *workspace.WriteGroup, action string) fCutoffResult {
	switch action {
	case "disable":
		mutation, err := lifecycleMutationForF(ctx, env, original.ID)
		if err != nil {
			return fCutoffResult{err: err}
		}
		saved, err := env.services.writeGroups.Disable(ctx, original.ID, mutation)
		if err != nil {
			return fCutoffResult{err: err}
		}
		return fCutoffResult{group: saved.Group}
	case "replace":
		current, err := env.services.writeGroups.Get(ctx, original.ID)
		if err != nil {
			return fCutoffResult{err: err}
		}
		draft := *current.Group
		draft.RowPolicy.IntervalSeconds = 5
		updateMutation, err := lifecycleMutationForF(ctx, env, original.ID)
		if err != nil {
			return fCutoffResult{err: err}
		}
		if _, err := env.services.writeGroups.Update(ctx, original.ID, workspace.WriteGroupMutation{
			WorkspaceID:               updateMutation.WorkspaceID,
			ExpectedWorkspaceRevision: updateMutation.ExpectedWorkspaceRevision,
			ExpectedGroupRevision:     current.Group.Revision,
			ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision,
			Group:                     &draft,
		}); err != nil {
			return fCutoffResult{err: err}
		}
		applyMutation, err := lifecycleMutationForF(ctx, env, original.ID)
		if err != nil {
			return fCutoffResult{err: err}
		}
		applied, err := env.services.writeGroups.Apply(ctx, original.ID, applyMutation)
		if err != nil {
			return fCutoffResult{err: err}
		}
		return fCutoffResult{group: applied.Group}
	default:
		return fCutoffResult{err: errors.New("unknown F cutoff action")}
	}
}

func lifecycleMutationForF(ctx context.Context, env *outageEnv, id string) (workspace.WriteGroupMutation, error) {
	record, err := env.services.workspace.GetOrCreate(ctx)
	if err != nil {
		return workspace.WriteGroupMutation{}, err
	}
	current, err := env.services.writeGroups.Get(ctx, id)
	if err != nil || current == nil || current.Group == nil {
		if err == nil {
			err = errors.New("write group is missing")
		}
		return workspace.WriteGroupMutation{}, err
	}
	return workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision,
		ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision}, nil
}

func TestProductionFParallelCutoffKeepsDurableIdentity(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		for _, action := range []string{"disable", "replace"} {
			t.Run(kind+"/"+action, func(t *testing.T) {
				f := newFProductionFixture(t, kind, "node-1/f-cutoff")
				require.NoError(t, f.pipe.Reconcile(t.Context()))
				tracingAt := lifecycleBase.Add(3 * time.Second)
				observed := lifecycleBase.Add(5 * time.Second)
				f.clock.set(tracingAt)
				start := make(chan struct{})
				const raceSamples = 4
				results := make(chan error, raceSamples)
				cutoff := make(chan fCutoffResult, 1)
				var wg sync.WaitGroup
				for i := 0; i < raceSamples; i++ {
					prefix := "f-cutoff-race-" + strconv.Itoa(i)
					wg.Go(func() {
						<-start
						results <- acceptFBucket(t.Context(), f.env, f.pipe, f.group, observed, prefix, 41, 410)
					})
				}
				wg.Go(func() {
					<-start
					result := runFCutoffAction(t.Context(), f.env, f.group, action)
					if result.err == nil {
						result.err = f.pipe.Reconcile(t.Context())
					}
					cutoff <- result
				})
				close(start)
				wg.Wait()
				for i := 0; i < raceSamples; i++ {
					require.NoError(t, <-results)
				}
				actionResult := <-cutoff
				require.NoError(t, actionResult.err)

				// Ensure the cutoff is installed at the same simulated instant even
				// if the racing reconcile observed the old saved group first.
				f.clock.set(tracingAt)
				require.NoError(t, f.pipe.Reconcile(t.Context()))
				require.NoError(t, acceptFBucket(t.Context(), f.env, f.pipe, f.group, observed, "f-cutoff-before", 51, 510))
				f.clock.set(lifecycleBase.Add(11 * time.Second))
				f.pipe.TickAll(t.Context())

				post := lifecycleBase.Add(15 * time.Second)
				f.clock.set(post)
				require.NoError(t, acceptFBucket(t.Context(), f.env, f.pipe, f.group, post, "f-cutoff-post", 61, 610))
				f.clock.set(lifecycleBase.Add(21 * time.Second))
				f.pipe.TickAll(t.Context())

				var postRevision string
				err := f.env.db.QueryRowContext(t.Context(), `SELECT group_revision FROM wg_delivery_samples WHERE sample_id LIKE 'f-cutoff-post-%' ORDER BY seq LIMIT 1`).Scan(&postRevision)
				if action == "disable" {
					require.ErrorIs(t, err, sql.ErrNoRows, "disabled group must not journal post-cutoff samples")
				} else {
					require.NoError(t, err)
					require.Equal(t, actionResult.group.AppliedRevision, postRevision, "post-cutoff samples belong to the superseding applied revision")
				}

				expected := 1
				if action == "replace" {
					expected = 2
				}
				recovered := f.env.pipeline(f.clock, "node-1/f-cutoff-restart")
				require.NoError(t, recovered.Start(t.Context()))
				t.Cleanup(func() { require.NoError(t, recovered.Stop(3*time.Second)) })
				rows := waitFCommitted(t, f.env, f.group.ID, expected)
				newRevision := ""
				if action == "replace" {
					newRevision = actionResult.group.AppliedRevision
					require.NotEqual(t, f.group.AppliedRevision, newRevision)
				}
				fAssertCutoffPayloads(t, rows, action, f.group.AppliedRevision, newRevision)
				seen := make(map[string]struct{}, len(rows))
				for _, row := range rows {
					require.NotEmpty(t, row.EffectKey)
					require.NotEmpty(t, row.RecordID)
					require.NotEmpty(t, row.PayloadDigest)
					require.NotContains(t, seen, row.EffectKey)
					seen[row.EffectKey] = struct{}{}
				}
				require.Equal(t, expected, fLocalReceiptCount(t, f.env, f.group.ID))
				targetRows := f.target.rows()
				require.Len(t, targetRows, expected)
				fAssertCutoffTargetRows(t, targetRows, action)
				fAssertTargetReceipts(t, f.env, f.group, rows)

				var openSamples int
				require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id = ? AND consumed = 0`, f.group.ID).Scan(&openSamples))
				require.Zero(t, openSamples, "cutoff rows are closed and durable before restart drain")
				var checkpoints int
				require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_checkpoints WHERE group_id = ?`, f.group.ID).Scan(&checkpoints))
				require.GreaterOrEqual(t, checkpoints, expected)
			})
		}
	}
}
