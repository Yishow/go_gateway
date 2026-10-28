package workspace

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLifecycleSuccessorCanonicalActionsAndSafeSeed(t *testing.T) {
	for _, kind := range []string{"basic", "profile", "active-wide", "active-long", "disabled"} {
		t.Run(kind, func(t *testing.T) {
			var h *circuitProfileHarness
			var id string
			if kind == "basic" {
				var req writeGroupRecreateReservationRequest
				h, req = recreateBasicFixture(t, "sqlite")
				id = req.PredecessorGroupID
			} else if kind == "profile" {
				h, id, _ = successorReviewFixture(t)
			} else if kind == "active-long" {
				h, id, _ = successorActiveWideFixture(t, WriteGroupLayoutManagedLongV1)
			} else {
				h, id, _ = successorActiveWideFixture(t)
			}
			if kind == "disabled" {
				_, err := h.db.Exec(`UPDATE write_groups SET status='disabled' WHERE id=?`, id)
				require.NoError(t, err)
			}
			before := recreateRepositorySnapshot(t, h)
			probe := &reservedProbeCounter{}
			h.service.WithManagedTableInspector(probe)
			view, err := h.service.Lifecycle(t.Context(), id)
			require.NoError(t, err)
			require.Equal(t, kind == "basic" || kind == "profile", view.Actions[3].Allowed)
			require.Equal(t, kind == "active-wide" || kind == "active-long" || kind == "disabled", view.Actions[4].Allowed)
			raw, err := json.Marshal(view)
			require.NoError(t, err)
			var dto map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(raw, &dto))
			require.Contains(t, dto, "successor_seed")
			var seed struct {
				BasisKind     string      `json:"basis_kind"`
				BasisRevision string      `json:"basis_revision"`
				BasisDigest   string      `json:"basis_digest"`
				Group         *WriteGroup `json:"group"`
			}
			require.NoError(t, json.Unmarshal(dto["successor_seed"], &seed))
			require.NotNil(t, seed.Group)
			require.NotEmpty(t, seed.BasisRevision)
			require.Len(t, seed.BasisDigest, 64)
			require.Empty(t, seed.Group.ID)
			require.Empty(t, seed.Group.Revision)
			require.Empty(t, seed.Group.AppliedRevision)
			require.Empty(t, seed.Group.Destination.SchemaRevision)
			require.Empty(t, seed.Group.Destination.SchemaDigest)
			require.Equal(t, WriteGroupStatusDraft, seed.Group.Status)
			require.Zero(t, probe.calls)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}
