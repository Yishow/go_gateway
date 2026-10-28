import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
import { savedGroup } from './writeGroupState';

/** Local lifecycle facts stay separate from runtime/delivery counters and mutable form state. */
export function savedLifecycle(group: WriteGroup = savedGroup(), workspaceRevision = 'wrev-1'): WriteGroupLifecycle {
  const destination = { connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision, destination_scope: 'opaque-scope' };
  return {
    workspace_id: group.workspace_id, workspace_revision: workspaceRevision, group_id: group.id, group_revision: group.revision,
    applied_revision: group.applied_revision, status: group.status, connector_id: destination.connector_id, connector_revision: destination.connector_revision,
    actions: ['edit', 'stop', 'remove', 'recreate', 'transition'].map((action) => ({
      action: action as WriteGroupLifecycle['actions'][number]['action'], allowed: group.status === 'deleted' ? action === 'recreate' : action !== 'recreate',
      reason_code: '', next_action: '', requires_review: action === 'recreate' || action === 'transition',
    })),
    history: { canonical: true, lineage_id: 'lineage-1', lineage_revision: 1, current_group_id: group.id, generations: [{
      group_id: group.id, generation: 1, current: true, reserved: false, status: group.status, group_revision: group.revision,
      applied_revision: group.applied_revision, saved_destination: destination, revisions: [{
        group_revision: group.applied_revision || group.revision, effective_at: '2026-10-01T00:00:00Z', sources: [], destination,
        pending: { samples: 1, outbox: { outstanding: 2, total: 6, queued: 2, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 4, skipped: 0 }, handoffs: [] },
      }],
    }] },
  };
}
