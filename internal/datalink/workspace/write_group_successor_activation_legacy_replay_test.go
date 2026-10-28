package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSuccessorActivationReplaysPersistedV1WithoutKind(t *testing.T) {
	h, id, req := atomicActivationFixture(t)
	event, err := h.service.ConfirmSuccessorActivation(t.Context(), id, req)
	require.NoError(t, err)
	// Model the exact historical v1 wire form independently of the new envelope
	// normalizer. No kind was present in its digest or persisted intent.
	req.Envelope.Version = "successor-activation-review-v1"
	req.Envelope.ActivationKind = ""
	encoded, err := json.Marshal(struct {
		Domain string                                `json:"domain"`
		Value  WriteGroupSuccessorActivationEnvelope `json:"value"`
	}{"go-gateway/successor-activation-review/v1", req.Envelope})
	require.NoError(t, err)
	sum := sha256.Sum256(encoded)
	req.ReviewDigest = hex.EncodeToString(sum[:])
	intent, err := json.Marshal(struct {
		Version            string                                    `json:"version"`
		PredecessorGroupID string                                    `json:"predecessor_group_id"`
		Confirmation       WriteGroupSuccessorActivationConfirmation `json:"confirmation"`
	}{"go-gateway/successor-activation-intent/v1", id, req})
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_group_cutovers SET review_digest=?,intent=? WHERE id=?`, req.ReviewDigest, string(intent), event.ID)
	require.NoError(t, err)
	event.ReviewDigest = req.ReviewDigest
	event.intent = string(intent)
	probe := &reservedProbeCounter{}
	h.service.WithManagedTableInspector(probe)
	replay, err := h.service.ConfirmSuccessorActivation(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, event, replay)
	require.Zero(t, probe.calls)
	_, err = h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.Error(t, err)
	bad := req
	bad.Envelope.ActivationKind = "first_start"
	_, err = h.service.ConfirmSuccessorActivation(t.Context(), id, bad)
	require.Error(t, err)
}
