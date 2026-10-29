import type { WriteGroupLifecycle, WriteGroupLifecycleHistory } from '@/types/studioV2WriteGroupLifecycle';
import type { WriteGroupSuccessorEnvelope, WriteGroupSuccessorPreparation, WriteGroupSuccessorActivationEnvelope, WriteGroupSuccessorOperationStatus } from '@/types/studioV2WriteGroupSuccessor';
import { wide, DIGEST, REVIEWED_BOUNDARY } from './cutoverFixtures';

export const candidate = () => {
  const group = wide({ id: '', revision: '', applied_revision: '', status: 'draft', migration: {},
    created_at: '0001-01-01T00:00:00Z', updated_at: '0001-01-01T00:00:00Z' });
  group.destination.table_name = '';
  return group;
};
export const history = (): WriteGroupLifecycleHistory => ({
  canonical: true, lineage_id: 'lineage-1', lineage_revision: 1, current_group_id: 'wide-1',
  generations: [{ group_id: 'wide-1', generation: 1, current: true, reserved: false, status: 'deleted',
    group_revision: 'rev-wide', applied_revision: 'applied-wide',
    saved_destination: { connector_id: 'conn-1', connector_revision: 'crev-1', destination_scope: DIGEST },
    revisions: [{ group_revision: 'applied-wide', effective_at: '2026-10-20T00:00:00Z', sources: [],
      destination: { connector_id: 'conn-1', connector_revision: 'crev-1', destination_scope: DIGEST },
      pending: { samples: 3, outbox: { outstanding: 2, total: 5, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 1, sql_committed: 3, skipped: 0 }, handoffs: [] } }],
  }],
});
export const lifecycle = (): WriteGroupLifecycle => ({
  workspace_id: 'ws-1', workspace_revision: 'ws-r1', group_id: 'wide-1', group_revision: 'rev-wide',
  applied_revision: 'applied-wide', status: 'deleted', connector_id: 'conn-1', connector_revision: 'crev-1',
  actions: ['edit', 'stop', 'remove', 'recreate', 'transition'].map(action => ({
    action: action as WriteGroupLifecycle['actions'][number]['action'], allowed: action === 'recreate',
    reason_code: action === 'recreate' ? '' : 'recording_removed', next_action: action === 'recreate' ? '' : 'review_lifecycle',
    requires_review: action === 'recreate' || action === 'transition',
  })), history: history(), successor_seed: { basis_kind: 'applied', basis_revision: 'applied-wide', basis_digest: DIGEST, group: candidate() },
});
export const envelope = (): WriteGroupSuccessorEnvelope => ({
  version: 'successor-review-v1', predecessor_group_id: 'wide-1', generation: 1,
  source_applied_revision: 'applied-wide', source_definition_digest: DIGEST,
  request: { workspace_id: 'ws-1', expected_workspace_revision: 'ws-r1', lineage_id: 'lineage-1', expected_lineage_revision: 1,
    expected_predecessor_revision: 'rev-wide', expected_predecessor_connector_revision: 'crev-1',
    expected_target_connector_revision: 'crev-1', operation_id: 'root-op-1', boundary: REVIEWED_BOUNDARY, group: candidate() },
});
export const preparation = (): WriteGroupSuccessorPreparation => ({
  workspace_revision: 'ws-r2', lineage_id: 'lineage-1', lineage_revision: 2, generation: 2, predecessor_group_id: 'wide-1',
  operation_id: 'root-op-1', intent_digest: DIGEST, group: wide({ id: 'new-2', revision: 'g-r2', applied_revision: '', status: 'draft' }),
});
export const activation = (): WriteGroupSuccessorActivationEnvelope => ({
  version: 'successor-activation-review-v2', activation_kind: 'cutover',
  request: { stage: { envelope: envelope(), review_digest: DIGEST,
    current: { workspace_revision: 'ws-r3', lineage_revision: 2, predecessor_revision: 'rev-wide', predecessor_connector_revision: 'crev-1', group_id: 'new-2', group_revision: 'g-r2', target_connector_revision: 'crev-1' } }, boundary: REVIEWED_BOUNDARY },
  source_applied_revision: 'applied-wide', source_definition_digest: DIGEST, target_definition_digest: DIGEST,
  counter_proof_digest: DIGEST, schema_revision: 'schema-op-1', schema_digest: DIGEST,
});
export const operationStatus = (): WriteGroupSuccessorOperationStatus => ({
  workspace_id: 'ws-1', lineage_id: 'lineage-1', predecessor_group_id: 'wide-1', operation_id: 'root-op-1',
  reserve: { status: 'succeeded', result: preparation() }, schema_operations: [],
  cutover: { status: 'not_submitted' }, first_start: { status: 'not_submitted' },
});
