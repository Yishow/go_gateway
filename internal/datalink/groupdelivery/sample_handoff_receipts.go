package groupdelivery

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"go-gateway/internal/datalink/counterproof"
)

type handoffSample struct {
	Kind, HandoffID, Address                                          string
	Origin                                                            GroupKey
	Seq                                                               int64
	Key                                                               GroupKey
	SampleID, Member, ObservedAt, Bucket, Payload, Digest, AcceptedAt string
	Consumed                                                          bool
}
type handoffReceipt struct {
	Targets                                                  []GroupKey
	Version                                                  int
	ID, Fence, State, Reason                                 string
	Source, Target, Origin                                   GroupKey
	SourceSeq, TargetSeq, OriginSeq                          int64
	SampleID, Member, ObservedAt, SourceBucket, TargetBucket string
	Digest, RawDigest, AcceptedAt                            string
	SampleProof                                              counterproof.CounterBindingReference `json:",omitzero"`
	OriginalProof                                            counterproof.CounterBindingReference `json:",omitzero"`
	TargetProof                                              counterproof.CounterBindingReference `json:",omitzero"`
	OriginRuntimeDigest, TargetRuntimeDigest                 string
	// Aliases point directly to the immutable audit named by the target row.
	CanonicalTransferID string `json:",omitempty"`
}

func handoffHash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
func handoffID(r handoffReceipt) string {
	raw, _ := json.Marshal([]any{1, r.Fence, r.Source, r.SourceSeq, r.Target, r.Origin, r.OriginSeq, r.SampleID, r.Member, r.ObservedAt, r.Digest, r.RawDigest, r.AcceptedAt, r.Targets})
	return handoffHash(string(raw))
}
func newHandoffReceipt(row, origin handoffSample, target GroupKey, fence string) handoffReceipt {
	r := handoffReceipt{Version: 1, Fence: fence, Source: row.Key, Target: target, Origin: origin.Key, SourceSeq: row.Seq, OriginSeq: origin.Seq, SampleID: row.SampleID, Member: row.Member, ObservedAt: row.ObservedAt, SourceBucket: row.Bucket, Digest: row.Digest, RawDigest: handoffHash(row.Payload), AcceptedAt: row.AcceptedAt}
	sample, _ := decodeVerifiedSample(row.Payload, row.Digest)
	r.SampleProof = sample.CounterProof
	r.ID = handoffID(r)
	return r
}
func readHandoffSample(ctx context.Context, q ownershipQuery, seq int64) (handoffSample, error) {
	var r handoffSample
	err := q.QueryRowContext(ctx, `SELECT seq,workspace_id,group_id,group_revision,COALESCE(NULLIF(origin_sample_id,''),sample_id),member_key,observed_at,bucket_start,payload,payload_digest,created_at,consumed,acceptance_kind,handoff_id,sample_id,origin_workspace_id,origin_group_id,origin_revision FROM wg_delivery_samples WHERE seq=?`, seq).Scan(&r.Seq, &r.Key.WorkspaceID, &r.Key.GroupID, &r.Key.GroupRevision, &r.SampleID, &r.Member, &r.ObservedAt, &r.Bucket, &r.Payload, &r.Digest, &r.AcceptedAt, &r.Consumed, &r.Kind, &r.HandoffID, &r.Address, &r.Origin.WorkspaceID, &r.Origin.GroupID, &r.Origin.GroupRevision)
	if err != nil {
		return r, &handoffStorageError{cause: err}
	}
	if r.Kind == "native" {
		if r.HandoffID != "" || r.Origin != (GroupKey{}) && r.Origin != r.Key {
			return r, ErrHandoffProofConflict
		}
		r.Origin = r.Key
	} else if r.Kind != "handoff" || r.HandoffID == "" || !r.Origin.Valid() {
		return r, ErrHandoffProofConflict
	}
	sample, err := decodeVerifiedSample(r.Payload, r.Digest)
	observed, e := time.Parse(time.RFC3339Nano, r.ObservedAt)
	_, acceptedErr := time.Parse(time.RFC3339Nano, r.AcceptedAt)
	if err != nil {
		return r, err
	}
	if e != nil || acceptedErr != nil || r.SampleID != sample.SampleID || r.Member != sample.MemberKey || !observed.Equal(sample.ObservedAt) || !r.Key.Valid() {
		return r, ErrHandoffProofConflict
	}
	return r, nil
}

const handoffColumns = `transfer_id,workspace_id,source_group_id,source_revision,source_seq,target_group_id,target_revision,target_seq,origin_seq,state,record`

func queryHandoffs(ctx context.Context, q ownershipQuery, where string, args ...any) ([]handoffReceipt, error) {
	rows, err := q.QueryContext(ctx, `SELECT `+handoffColumns+` FROM wg_delivery_sample_handoffs WHERE `+where, args...)
	if err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	defer rows.Close()
	var receipts []handoffReceipt
	for rows.Next() {
		values := make([]any, 11)
		dest := make([]any, 11)
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, &handoffStorageError{cause: err}
		}
		id, a := values[0].(string)
		workspace, b := values[1].(string)
		sg, c := values[2].(string)
		sr, d := values[3].(string)
		seq, e := values[4].(int64)
		tg, f := values[5].(string)
		tr, g := values[6].(string)
		origin, h := values[8].(int64)
		state, i := values[9].(string)
		raw, j := values[10].(string)
		target := sql.NullInt64{}
		if values[7] != nil {
			var ok bool
			target.Int64, ok = values[7].(int64)
			if !ok {
				return nil, ErrOwnershipUnavailable
			}
			target.Valid = true
		}
		if !(a && b && c && d && e && f && g && h && i && j) {
			return nil, ErrOwnershipUnavailable
		}
		var r handoffReceipt
		if !counterproof.UniqueJSONObject([]byte(raw), false) || json.Unmarshal([]byte(raw), &r) != nil || r.ID != id || r.Source != (GroupKey{workspace, sg, sr}) || r.Target != (GroupKey{workspace, tg, tr}) || r.SourceSeq != seq || r.OriginSeq != origin || r.TargetSeq != target.Int64 || r.State != state || r.Version != 1 || r.ID != handoffID(r) {
			return nil, ErrOwnershipUnavailable
		}
		canonical, marshalErr := json.Marshal(r)
		if marshalErr != nil || string(canonical) != raw {
			return nil, ErrOwnershipUnavailable
		}
		if (state == "accepted" || state == "finalized") != target.Valid {
			return nil, ErrOwnershipUnavailable
		}
		receipts = append(receipts, r)
	}
	if err := rows.Err(); err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	if err := rows.Close(); err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	return receipts, nil
}
func readHandoffID(ctx context.Context, q ownershipQuery, id string) (*handoffReceipt, error) {
	rs, err := queryHandoffs(ctx, q, "transfer_id=?", id)
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 {
		return nil, nil
	}
	return &rs[0], nil
}
func outgoingHandoffs(ctx context.Context, q ownershipQuery, seq int64) ([]handoffReceipt, error) {
	return queryHandoffs(ctx, q, "source_seq=?", seq)
}
func handoffReceiptsForKey(ctx context.Context, q ownershipQuery, key GroupKey) ([]handoffReceipt, error) {
	return queryHandoffs(ctx, q, `workspace_id=? AND ((source_group_id=? AND source_revision=?) OR (target_group_id=? AND target_revision=?))`, key.WorkspaceID, key.GroupID, key.GroupRevision, key.GroupID, key.GroupRevision)
}
func writeHandoffReceipt(ctx context.Context, tx *sql.Tx, r handoffReceipt) error {
	raw, err := json.Marshal(r)
	if err != nil {
		return ErrOwnershipUnavailable
	}
	var target any
	if r.TargetSeq > 0 {
		target = r.TargetSeq
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO wg_delivery_sample_handoffs(`+handoffColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(transfer_id) DO UPDATE SET state=excluded.state,target_seq=excluded.target_seq,record=excluded.record`, r.ID, r.Source.WorkspaceID, r.Source.GroupID, r.Source.GroupRevision, r.SourceSeq, r.Target.GroupID, r.Target.GroupRevision, target, r.OriginSeq, r.State, string(raw))
	return err
}
func handoffFence(ctx context.Context, q ownershipQuery, source, target GroupKey, member OwnershipSource, observed time.Time) (string, error) {
	var fence string
	err := q.QueryRowContext(ctx, `SELECT cutover_id FROM write_group_revision_ends WHERE workspace_id=? AND group_id=? AND group_revision=?`, source.WorkspaceID, source.GroupID, source.GroupRevision).Scan(&fence)
	if err != nil || fence == "" {
		return "", ErrOwnershipUnavailable
	}
	own, err := resolveCutoverOwnership(ctx, q, member, observed, GroupKey{})
	if err != nil {
		return "", err
	}
	if !own.Affected[source] || !own.Affected[target] || !handoffDescendant(own, target, source, map[GroupKey]bool{}) {
		return "", ErrOwnershipConflict
	}
	return fence, nil
}
func receiptMatchesRow(r handoffReceipt, row handoffSample) bool {
	return row.SampleID == r.SampleID && row.Member == r.Member && row.ObservedAt == r.ObservedAt && row.Digest == r.Digest && row.AcceptedAt == r.AcceptedAt && handoffHash(row.Payload) == r.RawDigest
}
func verifyHandoffReceiptEdge(ctx context.Context, q ownershipQuery, r handoffReceipt, seen map[string]bool) error {
	if seen[r.ID] || r.ID != handoffID(r) || !r.Source.Valid() || !r.Target.Valid() || !r.Origin.Valid() || r.Origin.WorkspaceID != r.Source.WorkspaceID || r.Source.WorkspaceID != r.Target.WorkspaceID {
		return ErrOwnershipUnavailable
	}
	seen[r.ID] = true
	defer delete(seen, r.ID)
	if err := verifyHandoffBranch(ctx, q, r); err != nil {
		return err
	}
	if r.State == "finalized" {
		if err := verifyCompactHandoffMeaning(ctx, q, r); err != nil {
			return err
		}
		return verifyHandoffTerminal(ctx, q, r)
	}
	source, err := readHandoffSample(ctx, q, r.SourceSeq)
	if err != nil {
		return err
	}
	if source.Key != r.Source || source.Bucket != r.SourceBucket || !receiptMatchesRow(r, source) {
		return ErrHandoffProofConflict
	}
	origin, err := readHandoffSample(ctx, q, r.OriginSeq)
	if err != nil {
		return err
	}
	if origin.Key != r.Origin || origin.Kind != "native" || !receiptMatchesRow(r, origin) || origin.Payload != source.Payload {
		return ErrHandoffProofConflict
	}
	if r.SourceSeq != r.OriginSeq {
		prior, err := incomingHandoff(ctx, q, r.SourceSeq)
		if err != nil {
			return err
		}
		if prior == nil || prior.OriginSeq != r.OriginSeq || prior.Target != r.Source {
			return ErrOwnershipUnavailable
		}
		if err = verifyHandoffReceipt(ctx, q, *prior, seen); err != nil {
			return err
		}
	}
	switch r.State {
	case "pending", "quarantined":
		if r.TargetSeq != 0 || r.CanonicalTransferID != "" || r.Reason != ErrOwnershipUnavailable.Error() && r.Reason != ErrHandoffProofConflict.Error() {
			return ErrOwnershipUnavailable
		}
		return nil
	case "accepted":
		target, err := readHandoffSample(ctx, q, r.TargetSeq)
		if err != nil {
			return err
		}
		if target.Key != r.Target || target.Kind != "handoff" || target.HandoffID != canonicalTransferID(r) || target.Origin != r.Origin || target.Bucket != r.TargetBucket || !receiptMatchesRow(r, target) || target.Payload != source.Payload {
			return ErrHandoffProofConflict
		}
		evidence, err := validateHandoffProof(ctx, q, handoffProofInput{Origin: r.Origin, Sequence: r.OriginSeq, Payload: origin.Payload, Digest: r.Digest, AcceptedAt: r.AcceptedAt}, r.Target)
		if err != nil {
			return err
		}
		if evidence.OriginalReference != r.OriginalProof || evidence.TargetReference != r.TargetProof || evidence.OriginRuntimeDigest != r.OriginRuntimeDigest || evidence.TargetRuntimeDigest != r.TargetRuntimeDigest {
			return ErrHandoffProofConflict
		}
	default:
		return ErrOwnershipUnavailable
	}
	return nil
}

// Branch facts differ legitimately across aliases, but never lose their own
// persisted fence and fanout authority, including after payload reclamation.
func verifyHandoffBranch(ctx context.Context, q ownershipQuery, r handoffReceipt) error {
	if r.ID != handoffID(r) || !r.Source.Valid() || !r.Target.Valid() || !r.Origin.Valid() || r.Origin.WorkspaceID != r.Source.WorkspaceID || r.Source.WorkspaceID != r.Target.WorkspaceID {
		return ErrOwnershipUnavailable
	}
	observed, err := time.Parse(time.RFC3339Nano, r.ObservedAt)
	if err != nil {
		return ErrOwnershipUnavailable
	}
	var ids []string
	if json.Unmarshal([]byte(r.Member), &ids) != nil || len(ids) != 3 {
		return ErrOwnershipUnavailable
	}
	member := OwnershipSource{r.Source.WorkspaceID, ids[0], ids[1], ids[2]}
	fence, err := handoffFence(ctx, q, r.Source, r.Target, member, observed)
	if err != nil {
		return err
	}
	if fence != r.Fence {
		return ErrOwnershipConflict
	}
	return verifyHandoffFanout(ctx, q, r)
}
func handoffDisposed(ctx context.Context, q ownershipQuery, seq int64) (bool, error) {
	rs, err := outgoingHandoffs(ctx, q, seq)
	if err != nil {
		return false, err
	}
	for _, r := range rs {
		if err = verifyHandoffReceipt(ctx, q, r, map[string]bool{}); err != nil {
			return false, err
		}
	}
	return len(rs) > 0, nil
}

func verifyHandoffFanout(ctx context.Context, q ownershipQuery, r handoffReceipt) error {
	if len(r.Targets) == 0 {
		return ErrOwnershipUnavailable
	}
	rs, err := outgoingHandoffs(ctx, q, r.SourceSeq)
	if err != nil {
		return err
	}
	if len(rs) != len(r.Targets) {
		return ErrOwnershipUnavailable
	}
	want := map[GroupKey]bool{}
	for _, key := range r.Targets {
		if want[key] || !key.Valid() {
			return ErrOwnershipUnavailable
		}
		want[key] = true
	}
	encoded, _ := json.Marshal(r.Targets)
	for _, sibling := range rs {
		set, _ := json.Marshal(sibling.Targets)
		if !want[sibling.Target] || string(encoded) != string(set) || sibling.OriginSeq != r.OriginSeq || sibling.Fence != r.Fence || sibling.Digest != r.Digest || sibling.RawDigest != r.RawDigest || sibling.AcceptedAt != r.AcceptedAt {
			return ErrOwnershipUnavailable
		}
		delete(want, sibling.Target)
	}
	return nil
}

// Every logical native ID remains legal, even if it equals an existing physical
// handoff address. Allocation probes actual occupied addresses under the write
// transaction; no prefix or spelling is reserved from native callers.
func allocateJournalAddress(ctx context.Context, q ownershipQuery, key GroupKey, preferred string) (string, error) {
	candidate := preferred
	for n := int64(0); ; n++ {
		if err := ctx.Err(); err != nil {
			return "", &handoffStorageError{cause: err}
		}
		var occupied int
		err := q.QueryRowContext(ctx, "SELECT 1 FROM wg_delivery_samples WHERE group_id=? AND group_revision=? AND sample_id=?", key.GroupID, key.GroupRevision, candidate).Scan(&occupied)
		if err == sql.ErrNoRows {
			return candidate, nil
		}
		if err != nil {
			return "", &handoffStorageError{cause: err}
		}
		candidate = preferred + "#" + strconv.FormatInt(n, 10)
	}
}
