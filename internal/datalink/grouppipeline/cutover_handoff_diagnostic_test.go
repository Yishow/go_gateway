package grouppipeline

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

type scopedHandoffDiagnostic interface {
	error
	HandoffScope() groupdelivery.GroupKey
	RepairAction() string
}

func TestCutoverHandoffFailureCarriesSafeScopeAndGuidance(t *testing.T) {
	for _, mode := range []string{"epoch", "missing-authority"} {
		t.Run(mode, func(t *testing.T) {
			f, _, healthy := originalConflictFixture(t, shapeMeanLong, mode)
			f.clock = at(7)
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "future-guidance", at(11), uint64(100))))
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(healthy.Group, "healthy", at(11), uint64(7))))
			original := handoffOriginalRows(t, f)
			f.commit()
			f.source.EndsAt = f.boundary
			f.groups.groups = append(f.groups.groups, f.target.Group)
			f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
			expected := groupdelivery.ErrHandoffProofConflict
			if mode == "missing-authority" {
				_, err := f.local.Exec(`DELETE FROM wg_runtime_versions WHERE group_id=?`, f.target.GroupID)
				require.NoError(t, err)
				expected = groupdelivery.ErrOwnershipUnavailable
			}
			f.clock = at(12)
			for _, operation := range []func() error{
				func() error { return f.pipe.PrepareCutoverHandoffs(t.Context(), f.source.WorkspaceID) },
				func() error { return f.pipe.Reconcile(t.Context()) },
			} {
				err := operation()
				require.ErrorIs(t, err, expected)
				var diagnostic scopedHandoffDiagnostic
				require.True(t, errors.As(err, &diagnostic), "blocked handoff must carry source scope and actionable guidance")
				require.Equal(t, groupKey(f.source), diagnostic.HandoffScope())
				require.Equal(t, "inspect-immutable-evidence", diagnostic.RepairAction())
				require.Equal(t, expected.Error(), diagnostic.Error())
				require.Equal(t, original, handoffOriginalRows(t, f))
			}
			require.NoError(t, f.pipe.boundaries[boundaryKey(healthy.GroupID, healthy.AppliedRevision)].boundary.Tick(t.Context(), at(80)))
			f.deliver()
			var reading string
			require.NoError(t, f.plainTarget().QueryRow(`SELECT reading FROM healthy_proof_readings`).Scan(&reading))
			require.Equal(t, "7", reading)
		})
	}
}

func TestScopedHandoffRefusalClassificationAndErrorChain(t *testing.T) {
	key := groupdelivery.GroupKey{WorkspaceID: "w", GroupID: "g", GroupRevision: "r"}
	for _, reason := range []*groupdelivery.OwnershipError{groupdelivery.ErrHandoffProofConflict, groupdelivery.ErrOwnershipConflict, groupdelivery.ErrOwnershipUnavailable} {
		t.Run(reason.Reason, func(t *testing.T) {
			cause := errors.Join(errors.New("private-endpoint-password-payload"), reason)
			err := scopedHandoffRefusal(key, cause)
			require.Equal(t, reason.Reason, err.Error())
			require.ErrorIs(t, err, cause)
			require.ErrorIs(t, err, reason)
			var original *groupdelivery.OwnershipError
			require.ErrorAs(t, err, &original)
			require.Same(t, reason, original)
			var diagnostic scopedHandoffDiagnostic
			require.ErrorAs(t, err, &diagnostic)
			require.Equal(t, key, diagnostic.HandoffScope())
			require.Equal(t, "inspect-immutable-evidence", diagnostic.RepairAction())
		})
	}
	for _, cause := range []error{nil, errors.New("unknown-private-error"), &groupdelivery.OwnershipError{Reason: "unknown-private-reason"}} {
		require.True(t, cause == scopedHandoffRefusal(key, cause))
	}
	require.Same(t, groupdelivery.ErrOwnershipUnavailable, scopedHandoffRefusal(groupdelivery.GroupKey{}, groupdelivery.ErrOwnershipUnavailable))
}

func TestScopedHandoffRefusalsRetainEachSourceWhenJoined(t *testing.T) {
	first := groupdelivery.GroupKey{WorkspaceID: "w", GroupID: "g1", GroupRevision: "r1"}
	second := groupdelivery.GroupKey{WorkspaceID: "w", GroupID: "g2", GroupRevision: "r2"}
	joined := errors.Join(scopedHandoffRefusal(first, groupdelivery.ErrHandoffProofConflict), scopedHandoffRefusal(second, groupdelivery.ErrOwnershipUnavailable))
	branches := joined.(interface{ Unwrap() []error }).Unwrap()
	require.Len(t, branches, 2)
	for i, expected := range []groupdelivery.GroupKey{first, second} {
		var diagnostic scopedHandoffDiagnostic
		require.ErrorAs(t, branches[i], &diagnostic)
		require.Equal(t, expected, diagnostic.HandoffScope())
		require.Equal(t, "inspect-immutable-evidence", diagnostic.RepairAction())
	}
	require.ErrorIs(t, joined, groupdelivery.ErrHandoffProofConflict)
	require.ErrorIs(t, joined, groupdelivery.ErrOwnershipUnavailable)
}
