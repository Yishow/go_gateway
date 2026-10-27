package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go-gateway/internal/datalink/counterproof"
)

type legacyProofAssociation struct {
	WorkspaceID, GroupID, GroupRevision, SampleID       string
	SampleSeq                                           int64
	SampleDigest, AcceptedAt, MemberKey, ProofMemberKey string
	ProofVersion                                        int
	ProofDigest, RuntimeDigest, EvidenceKind            string
	AssociationVersion                                  int
	AssociatedAt                                        string
}

func insertLegacyAssociation(ctx context.Context, tx *sql.Tx, want legacyProofAssociation) error {
	rows, err := tx.QueryContext(ctx, `SELECT workspace_id,group_id,group_revision,sample_id,sample_seq,sample_digest,accepted_at,member_key,proof_member_key,proof_version,proof_digest,runtime_digest,evidence_kind,association_version,associated_at FROM wg_counter_proof_associations WHERE sample_seq=? OR (workspace_id=? AND group_id=? AND group_revision=? AND sample_id=?)`, want.SampleSeq, want.WorkspaceID, want.GroupID, want.GroupRevision, want.SampleID)
	if err != nil {
		return err
	}
	exists := false
	for rows.Next() {
		var got legacyProofAssociation
		if err = rows.Scan(&got.WorkspaceID, &got.GroupID, &got.GroupRevision, &got.SampleID, &got.SampleSeq, &got.SampleDigest, &got.AcceptedAt, &got.MemberKey, &got.ProofMemberKey, &got.ProofVersion, &got.ProofDigest, &got.RuntimeDigest, &got.EvidenceKind, &got.AssociationVersion, &got.AssociatedAt); err != nil {
			break
		}
		if got.AssociationVersion != legacyAssociationVersion {
			err = &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofUnsupported, Version: got.AssociationVersion}
			break
		}
		if _, err = time.Parse(time.RFC3339Nano, got.AssociatedAt); err != nil {
			err = sampleProofError(counterproof.ErrCounterBindingProofConflict)
			break
		}
		got.AssociatedAt = want.AssociatedAt
		if exists || got != want {
			err = sampleProofError(counterproof.ErrCounterBindingProofConflict)
			break
		}
		exists = true
	}
	err = errors.Join(err, rows.Err(), rows.Close())
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO wg_counter_proof_associations (workspace_id,group_id,group_revision,sample_id,sample_seq,sample_digest,accepted_at,member_key,proof_member_key,proof_version,proof_digest,runtime_digest,evidence_kind,association_version,associated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, want.WorkspaceID, want.GroupID, want.GroupRevision, want.SampleID, want.SampleSeq, want.SampleDigest, want.AcceptedAt, want.MemberKey, want.ProofMemberKey, want.ProofVersion, want.ProofDigest, want.RuntimeDigest, want.EvidenceKind, want.AssociationVersion, want.AssociatedAt)
	return err
}
