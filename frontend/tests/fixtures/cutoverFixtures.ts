import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { CutoverGroupView, WriteGroupCutover, WriteGroupCutoverReview } from '@/types/studioV2WriteGroupCutover';
import { longProfileGroup } from './circuitBatchPanel';
import { savedGroup } from './writeGroupState';

export const REVIEWED_BOUNDARY = '2099-01-01T00:00:00.123456Z';
export const DIGEST = 'a'.repeat(64);

export const wide = (over: Partial<WriteGroup> = {}): WriteGroup => savedGroup({
  id: 'wide-1', name: 'PLC A raw', basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1', status: 'running', revision: 'rev-wide', applied_revision: 'applied-wide',
  members: ['temperature', 'pressure', 'batch'].map((column, index) => ({ device_id: 'dev-1', point_id: `pt-${index}`, tag_id: `tag-${index}`, source_revision: `s${index}`, mapping_revision: `m${index}`, target_column: column, required: index !== 2 })),
  destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: '/d.db', table_schema: '', table_name: 'gw_basic_dev_1', storage_strategy: 'managed' },
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0, incomplete_policy: 'skip_row' }, ...over,
});
export const target = (over: Partial<WriteGroup> = {}): WriteGroup => ({ ...longProfileGroup({ id: 'tgt-1', circuitId: 'c1', circuitName: 'Feeder', interval: 30, points: [0, 1] }), status: 'draft', applied_revision: '', ...over });
export const view = (group: WriteGroup, over: Partial<CutoverGroupView> = {}): CutoverGroupView => ({
  group_id: group.id, name: group.name, revision: group.revision, layout: group.layout_version ?? 'managed_wide_v1', interval_seconds: group.row_policy.interval_seconds,
  connector_id: 'conn-1', table_schema: '', table_name: group.destination.table_name, member_count: group.members.length, schema_confirmed: true, ...over,
});
export const member = (index: number) => ({ device_id: 'dev-1', point_id: `pt-${index}`, tag_id: `tag-${index}` });
export const reviewOf = (over: Partial<WriteGroupCutoverReview> = {}): WriteGroupCutoverReview => ({
  kind: 'forward', boundary: REVIEWED_BOUNDARY, digest: DIGEST, applicable: true, issues: [], source: view(wide()), targets: [view(target(), { storage_mode: 'plain', profile_kind: 'instantaneous' })],
  moved: [member(0), member(1)], unselected: [member(2)], remaining: [member(2)], ...over,
});
export const recordOf = (over: Partial<WriteGroupCutover> = {}): WriteGroupCutover => ({
  id: 'cut-1', workspace_id: 'ws-1', kind: 'forward', source_group_id: 'wide-1', boundary: REVIEWED_BOUNDARY, review_digest: DIGEST,
  ended: [{ group_id: 'wide-1', revision: 'applied-wide' }], started: [{ group_id: 'tgt-1', revision: 'rev-tgt-1-1' }], moved: [member(0), member(1)], created_at: '2026-10-24T10:00:00Z', ...over,
});
export const quietDelivery = {
  group_id: 'x', intake: { state: 'active' }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
  last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
};
