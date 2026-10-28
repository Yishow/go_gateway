-- Transfer audit supplements, never rewrites, original accepted sample bytes.
-- finalized is accepted ownership compacted only by ordinary local reclaim.
CREATE TABLE IF NOT EXISTS wg_delivery_sample_handoffs (
    transfer_id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL,
    source_group_id TEXT NOT NULL,
    source_revision TEXT NOT NULL,
    source_seq INTEGER NOT NULL,
    target_group_id TEXT NOT NULL,
    target_revision TEXT NOT NULL,
    target_seq INTEGER,
    origin_seq INTEGER NOT NULL,
    state TEXT NOT NULL CHECK(state IN ('pending','quarantined','accepted','finalized')),
    record TEXT NOT NULL CHECK(json_valid(record)),
    UNIQUE(source_seq,target_group_id,target_revision),
    CHECK(source_seq > 0 AND origin_seq > 0),
    CHECK((state IN ('pending','quarantined') AND target_seq IS NULL)
       OR (state IN ('accepted','finalized') AND target_seq > 0))
);
CREATE INDEX IF NOT EXISTS idx_wg_handoff_source
 ON wg_delivery_sample_handoffs(workspace_id,source_group_id,source_revision,source_seq);
CREATE INDEX IF NOT EXISTS idx_wg_handoff_target
 ON wg_delivery_sample_handoffs(workspace_id,target_group_id,target_revision,state);
CREATE UNIQUE INDEX IF NOT EXISTS idx_wg_handoff_target_seq
 ON wg_delivery_sample_handoffs(target_seq) WHERE target_seq IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_wg_handoff_origin ON wg_delivery_sample_handoffs(origin_seq);

-- Logical acceptance identity is distinct from the internal storage address.
-- Existing native rows retain every old column and use their sample_id fallback.
ALTER TABLE wg_delivery_samples ADD COLUMN acceptance_kind TEXT NOT NULL DEFAULT 'native'
 CHECK(acceptance_kind IN ('native','handoff'));
ALTER TABLE wg_delivery_samples ADD COLUMN origin_workspace_id TEXT NOT NULL DEFAULT '';
ALTER TABLE wg_delivery_samples ADD COLUMN origin_group_id TEXT NOT NULL DEFAULT '';
ALTER TABLE wg_delivery_samples ADD COLUMN origin_revision TEXT NOT NULL DEFAULT '';
ALTER TABLE wg_delivery_samples ADD COLUMN origin_sample_id TEXT NOT NULL DEFAULT '';
ALTER TABLE wg_delivery_samples ADD COLUMN handoff_id TEXT NOT NULL DEFAULT ''
 CHECK(
 (acceptance_kind='native' AND handoff_id='' AND
   ((origin_workspace_id='' AND origin_group_id='' AND origin_revision='' AND origin_sample_id='')
    OR (origin_workspace_id=workspace_id AND origin_group_id=group_id AND origin_revision=group_revision AND origin_sample_id<>'')))
 OR
 (acceptance_kind='handoff' AND handoff_id<>'' AND origin_workspace_id<>'' AND origin_group_id<>'' AND origin_revision<>'' AND origin_sample_id<>'')
 );
CREATE UNIQUE INDEX IF NOT EXISTS idx_wg_delivery_native_identity
 ON wg_delivery_samples(group_id,group_revision,COALESCE(NULLIF(origin_sample_id,''),sample_id))
 WHERE acceptance_kind='native';
CREATE UNIQUE INDEX IF NOT EXISTS idx_wg_delivery_handoff_identity
 ON wg_delivery_samples(group_id,group_revision,origin_workspace_id,origin_group_id,origin_revision,origin_sample_id)
 WHERE acceptance_kind='handoff';
CREATE UNIQUE INDEX IF NOT EXISTS idx_wg_delivery_handoff_id
 ON wg_delivery_samples(handoff_id) WHERE acceptance_kind='handoff';
