-- Additive local evidence: the original sample and runtime descriptor remain
-- immutable. Each append-only row is also the audit of the proven association.
CREATE TABLE IF NOT EXISTS wg_counter_proof_associations (
    workspace_id        TEXT NOT NULL,
    group_id            TEXT NOT NULL,
    group_revision      TEXT NOT NULL,
    sample_id           TEXT NOT NULL,
    sample_seq          INTEGER NOT NULL UNIQUE,
    sample_digest       TEXT NOT NULL,
    accepted_at         TEXT NOT NULL,
    member_key          TEXT NOT NULL,
    proof_member_key    TEXT NOT NULL,
    proof_version       INTEGER NOT NULL,
    proof_digest        TEXT NOT NULL,
    runtime_digest      TEXT NOT NULL,
    evidence_kind       TEXT NOT NULL,
    association_version INTEGER NOT NULL,
    associated_at       TEXT NOT NULL,
    PRIMARY KEY (workspace_id, group_id, group_revision, sample_id)
);
