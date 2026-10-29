package workspace

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSuccessorActivationReplayRejectsCorruptEventDetail(t *testing.T) {
	for _, mutation := range []string{"moved", "prior"} {
		t.Run(mutation, func(t *testing.T) {
			h, id, req := atomicActivationFixture(t)
			event, err := h.service.ConfirmSuccessorActivation(t.Context(), id, req)
			require.NoError(t, err)
			if mutation == "prior" {
				_, err = h.db.Exec(`UPDATE write_group_cutovers SET prior_cutover_id='unexpected-prior' WHERE id=?`, event.ID)
			} else {
				detail, marshalErr := json.Marshal(cutoverDetail{Ended: event.Ended, Started: event.Started, Moved: nil, OwnershipFenceRevision: event.OwnershipFenceRevision})
				require.NoError(t, marshalErr)
				_, err = h.db.Exec(`UPDATE write_group_cutovers SET detail=? WHERE id=?`, string(detail), event.ID)
			}
			require.NoError(t, err)
			_, err = h.service.ConfirmSuccessorActivation(t.Context(), id, req)
			require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
		})
	}
}
