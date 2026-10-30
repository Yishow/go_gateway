package grouppipeline

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"

	"github.com/stretchr/testify/require"
)

func TestDiagnosticFallbackCounterTickAcceptedAfterNative(t *testing.T) {
	f := ownerCounterFixture(t)
	var out []OwnerObservation
	reads, native, generic, void := 0, 0, 0, 0
	ownerCapture(f.pipe, &out, &reads)
	f.pipe.deps.OwnerObservation.Observe = func(OwnerObservation) { void++ }
	f.pipe.deps.OwnerObservation.ObserveCaptured = func(o OwnerObservation) bool {
		require.Equal(t, 1, native, "native callback must run before diagnostic flush")
		out = append(out, o)
		return o.Available
	}
	f.pipe.config.OnError = func(err error) { native++; require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict) }
	f.pipe.config.OnDiagnosticError = func() { generic++ }
	damageCounterProof(t, f, "counter-replay-proof-conflict")
	before := counterDurableState(t, f.store)
	f.closeUntil(at(1200))
	require.Equal(t, 1, native)
	require.Len(t, out, 1)
	require.True(t, out[0].Available)
	require.Zero(t, void, "capture is an alternative, never a second emission")
	require.Zero(t, generic)
	require.Equal(t, before, counterDurableState(t, f.store))
}

func TestDiagnosticFallbackCounterTickFailClosed(t *testing.T) {
	for _, mode := range []string{"nil", "void", "capture-only", "reject", "panic", "missing-proof", "changed-proof", "drop"} {
		t.Run(mode, func(t *testing.T) {
			f := ownerCounterFixture(t)
			var out []OwnerObservation
			reads, native, generic := 0, 0, 0
			ownerCapture(f.pipe, &out, &reads)
			f.pipe.deps.OwnerObservation.ObserveCaptured = func(o OwnerObservation) bool {
				out = append(out, o)
				if mode == "panic" {
					panic("observer failed")
				}
				return mode != "reject" && o.Available
			}
			switch mode {
			case "nil":
				f.pipe.deps.OwnerObservation.Observe = nil
				f.pipe.deps.OwnerObservation.ObserveCaptured = nil
			case "void":
				f.pipe.deps.OwnerObservation.ObserveCaptured = nil
			case "capture-only":
				f.pipe.deps.OwnerObservation.Observe = nil
			case "missing-proof":
				f.pipe.deps.OwnerObservation.ScopeRevision = nil
			case "changed-proof":
				f.pipe.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) {
					reads++
					if reads == 1 {
						return "before", true
					}
					return "after", true
				}
			case "drop":
				f.pipe.deps.OwnerObservation.RecordLimit = 0
			}
			f.pipe.config.OnError = func(err error) { native++; require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict) }
			f.pipe.config.OnDiagnosticError = func() { generic++ }
			damageCounterProof(t, f, "counter-replay-proof-conflict")
			before := counterDurableState(t, f.store)
			f.closeUntil(at(1200))
			require.Equal(t, 1, native)
			if mode == "capture-only" {
				require.Len(t, out, 1)
				require.Zero(t, generic)
			} else {
				require.Equal(t, 1, generic)
			}
			require.Equal(t, before, counterDurableState(t, f.store))
		})
	}
}

func TestDiagnosticFallbackOperationalReportAndPanic(t *testing.T) {
	expected := errors.New("native operational error")
	native, generic := 0, 0
	p := New(Dependencies{}, Config{
		OnError:           func(err error) { native++; require.Same(t, expected, err) },
		OnDiagnosticError: func() { generic++; panic("diagnostic failed") },
	})
	require.NotPanics(t, func() { p.report(expected) })
	p.report(nil)
	require.Equal(t, 1, native)
	require.Equal(t, 1, generic)
	p.config.OnError = func(error) { panic("native panic remains visible") }
	require.PanicsWithValue(t, "native panic remains visible", func() { p.report(expected) })
}

func TestDiagnosticFallbackOwnerTickExactOriginalError(t *testing.T) {
	f, target, _ := ownerRefusalFixture(t)
	f.pipe.boundaries = map[string]*managed{boundaryKey(target.groupID, target.revision): target}
	var observed []OwnerObservation
	reads, native, generic := 0, 0, 0
	ownerCapture(f.pipe, &observed, &reads)
	f.pipe.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) { return "workspace-revision", true }
	f.pipe.deps.OwnerObservation.ObserveCaptured = func(o OwnerObservation) bool { observed = append(observed, o); return o.Available }
	f.pipe.config.OnError = func(err error) { native++; require.Same(t, groupdelivery.ErrOwnershipConflict, err) }
	f.pipe.config.OnDiagnosticError = func() { generic++ }
	before := counterDurableState(t, f.store)
	f.clock = at(40)
	f.pipe.TickAll(t.Context())
	require.Equal(t, 1, native)
	require.Len(t, observed, 1)
	require.True(t, observed[0].Available)
	require.Zero(t, generic)
	require.Equal(t, before, counterDurableState(t, f.store))
}

type diagnosticIsOnlyError struct{}

func (diagnosticIsOnlyError) Error() string { return "unknown" }
func (diagnosticIsOnlyError) Is(error) bool { return true }

type diagnosticCycleError struct{}

func (diagnosticCycleError) Error() string   { return "cycle" }
func (e diagnosticCycleError) Unwrap() error { return e }

func TestDiagnosticFallbackClosedLeafMatching(t *testing.T) {
	key := groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-G", GroupRevision: "rev-1"}
	typed := groupdelivery.ErrOwnershipConflict
	d := ownerDescriptor{key: key, kind: OwnerRecordingRefusal, reason: typed.Reason}
	for _, tc := range []struct {
		name string
		err  error
		want bool
	}{
		{"typed", typed, true},
		{"single-wrapper", fmt.Errorf("native context: %w", typed), true},
		{"joined-one", errors.Join(typed), false},
		{"joined-covered-unknown", errors.Join(typed, errors.New("operational")), false},
		{"joined-all-covered", errors.Join(typed, groupdelivery.ErrOwnershipUnavailable), false},
		{"multi-wrap", fmt.Errorf("two: %w %w", typed, typed), false},
		{"is-only", diagnosticIsOnlyError{}, false},
		{"matching-string", errors.New(typed.Error()), false},
		{"cycle", diagnosticCycleError{}, false},
		{"typed-nil", (*groupdelivery.OwnershipError)(nil), false},
		{"unknown-kind", groupdelivery.ErrIntakeStopped, false},
	} {
		t.Run(tc.name, func(t *testing.T) { require.Equal(t, tc.want, tickDescriptorCovers(tc.err, key, d)) })
	}
	foreign := key
	foreign.GroupRevision = "other"
	require.False(t, tickDescriptorCovers(typed, foreign, d))
	d.reason = groupdelivery.ErrOwnershipUnavailable.Reason
	require.False(t, tickDescriptorCovers(typed, key, d))
	proof := &counterproof.CounterBindingProofError{Reason: errors.New(counterproof.ErrCounterBindingProofConflict.Error())}
	require.False(t, tickDescriptorCovers(proof, key, ownerDescriptor{key: key, kind: OwnerCounterQuarantine, reason: proof.Error()}))
}

func TestDiagnosticFallbackIntervalsCannotBorrowAcceptance(t *testing.T) {
	for _, mode := range []string{"prior-boundary", "joined-unknown", "joined-typed", "duplicate-descriptor", "no-descriptor", "overflow"} {
		t.Run(mode, func(t *testing.T) {
			f, target, _ := ownerRefusalFixture(t)
			p := f.pipe
			reads := 0
			var out []OwnerObservation
			ownerCapture(p, &out, &reads)
			p.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) { return "workspace-revision", true }
			p.deps.OwnerObservation.ObserveCaptured = func(o OwnerObservation) bool { out = append(out, o); return o.Available }
			if mode == "overflow" {
				p.deps.OwnerObservation.RecordLimit = 1
			}
			native, generic := 0, 0
			var expected error = groupdelivery.ErrOwnershipConflict
			p.config.OnError = func(got error) { native++; require.Same(t, expected, got) }
			p.config.OnDiagnosticError = func() { generic++ }
			ctx, finish := p.ownerOperation(t.Context())
			if mode == "prior-boundary" || mode == "overflow" {
				p.reportTick(ctx, target, func() error { return expected })
				if mode == "prior-boundary" {
					source := p.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]
					require.NotNil(t, source)
					require.NotSame(t, target, source)
					target = source
					expected = errors.New("uncovered next boundary")
				}
			}
			switch mode {
			case "joined-unknown":
				expected = errors.Join(expected, errors.New("local failure"))
			case "joined-typed":
				expected = errors.Join(expected, groupdelivery.ErrOwnershipUnavailable)
			case "no-descriptor":
				p.deps.OwnerObservation.Observe = nil
				p.deps.OwnerObservation.ObserveCaptured = nil
			}
			p.reportTick(ctx, target, func() error {
				if mode == "duplicate-descriptor" {
					p.recordOwnerObservation(ctx, ownerDescriptor{key: groupKey(f.target), kind: OwnerRecordingRefusal, reason: groupdelivery.ErrOwnershipConflict.Reason, owner: target.boundary.ReplayState().OwnershipView, ownerKnown: true})
				}
				return expected
			})
			require.GreaterOrEqual(t, native, 1)
			finish()
			require.Equal(t, 1, generic, "only the uncovered exact invocation keeps its generic")
		})
	}
}

func TestDiagnosticFallbackOwnerTokenChangesAndReentry(t *testing.T) {
	t.Run("token-changed-after-native", func(t *testing.T) {
		f, target, _ := ownerRefusalFixture(t)
		p := f.pipe
		var out []OwnerObservation
		reads, native, generic := 0, 0, 0
		ownerCapture(p, &out, &reads)
		p.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) { return "workspace-revision", true }
		p.deps.OwnerObservation.ObserveCaptured = func(o OwnerObservation) bool { out = append(out, o); return true }
		p.config.OnError = func(err error) {
			native++
			require.Same(t, groupdelivery.ErrOwnershipConflict, err)
			_, changeErr := f.local.ExecContext(t.Context(), `UPDATE write_group_versions SET payload=json_set(payload,'$.row_policy.interval_seconds',120) WHERE group_id=? AND group_revision=?`, f.target.GroupID, f.target.AppliedRevision)
			require.NoError(t, changeErr)
		}
		p.config.OnDiagnosticError = func() { generic++ }
		ctx, finish := p.ownerOperation(t.Context())
		p.reportTick(ctx, target, func() error { return groupdelivery.ErrOwnershipConflict })
		require.Equal(t, 1, native)
		finish()
		requireOwnerUnavailable(t, out)
		require.Equal(t, 1, generic)
	})
	t.Run("capture-reentry-after-native-locks", func(t *testing.T) {
		f := ownerCounterFixture(t)
		m := f.pipe.boundaries[boundaryKey(f.group.ID, "rev-1")]
		var out []OwnerObservation
		reads, native, generic := 0, 0, 0
		ownerCapture(f.pipe, &out, &reads)
		entered := false
		f.pipe.config.OnError = func(err error) { native++; require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict) }
		f.pipe.config.OnDiagnosticError = func() { generic++ }
		f.pipe.deps.OwnerObservation.ObserveCaptured = func(o OwnerObservation) bool {
			require.Equal(t, 1, native)
			if !entered {
				entered = true
				require.True(t, m.proofGate.TryLock())
				m.proofGate.Unlock()
				require.True(t, f.pipe.mu.TryLock())
				f.pipe.mu.Unlock()
				require.True(t, f.pipe.reconcileMu.TryLock())
				f.pipe.reconcileMu.Unlock()
				require.NoError(t, f.pipe.Reconcile(t.Context()))
			}
			return o.Available
		}
		damageCounterProof(t, f, "counter-replay-proof-conflict")
		f.closeUntil(at(1200))
		require.True(t, entered)
		require.Zero(t, generic)
	})
}

func TestDiagnosticFallbackQuotaAcceptRemainsNativeAndVoid(t *testing.T) {
	f, p, out, _ := intakeQuotaFixture(t, groupdelivery.QuotaScopeGroup)
	before := counterDurableState(t, f.store)
	native, generic := 0, 0
	p.config.OnError = func(error) { native++ }
	p.config.OnDiagnosticError = func() { generic++ }
	err := p.AcceptSample(t.Context(), envelope(f.groups.groups[0], "refused", at(2), 2.0))
	var refusal *groupdelivery.QuotaError
	require.ErrorAs(t, err, &refusal)
	require.Equal(t, groupdelivery.ReasonQuotaHardLimit, refusal.Reason)
	require.Len(t, intakeKind(*out, IntakeQuotaRefused), 1)
	require.Equal(t, before, counterDurableState(t, f.store))
	require.Zero(t, native, "AcceptSample still returns, not reports, its original error")
	require.Zero(t, generic)
}

func TestDiagnosticFallbackTickRealDiskFullStaysUnavailable(t *testing.T) {
	f := ownerCounterFixture(t)
	f.local.SetMaxOpenConns(1)
	f.local.SetMaxIdleConns(1)
	_, err := f.local.ExecContext(t.Context(), `CREATE TABLE diagnostic_filler (payload BLOB)`)
	require.NoError(t, err)
	var pages int
	require.NoError(t, f.local.QueryRowContext(t.Context(), `PRAGMA page_count`).Scan(&pages))
	_, err = f.local.ExecContext(t.Context(), fmt.Sprintf(`PRAGMA max_page_count = %d`, pages))
	require.NoError(t, err)
	// Real SQLite file growth is refused; this is not a fabricated quota Tick.
	_, err = f.local.ExecContext(t.Context(), `INSERT INTO diagnostic_filler VALUES (zeroblob(1048576))`)
	require.Error(t, err)
	var out []IntakeQuotaObservation
	reads, native, generic := 0, 0, 0
	intakeCapture(f.pipe, &out, &reads)
	f.pipe.config.OnError = func(err error) {
		native++
		var refusal *groupdelivery.QuotaError
		require.ErrorAs(t, err, &refusal)
		require.Equal(t, groupdelivery.ReasonDiskFull, refusal.Reason)
		require.Zero(t, refusal.MaxBytes)
	}
	f.pipe.config.OnDiagnosticError = func() { generic++ }
	before := counterDurableState(t, f.store)
	f.closeUntil(at(120000))
	require.Equal(t, 1, native)
	require.Equal(t, 1, generic)
	require.Equal(t, []IntakeQuotaObservation{{}}, out)
	require.Equal(t, before, counterDurableState(t, f.store), "failed closure cannot consume accepted inputs")
}
