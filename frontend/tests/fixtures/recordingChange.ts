import { vi } from 'vitest';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
import type * as S from '@/types/studioV2WriteGroupSuccessor';
import type { SchemaPreviewToken } from '@/types/recordingPlan';
import type { studioV2WriteGroupSuccessorAPI } from '@/services/studioV2WriteGroupSuccessor';
import type { RecordingChangeDependencies } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController';
export const HASH = 'a'.repeat(64);
export const BOUNDARY = '2099-01-01T00:00:00.000Z';
export function fixture(workspace = crypto.randomUUID()) {
  const group: WriteGroup = { id: 'G1', workspace_id: workspace, revision: 'R1', applied_revision: '', name: 'Saved recording', status: 'deleted',
    members: [{ device_id: 'D', point_id: 'P', tag_id: 'T', source_revision: 'S', mapping_revision: 'M', target_column: 'value', required: false }],
    destination: { connector_id: 'C', connector_revision: 'C1', database: 'private.db', table_schema: 'main', table_name: 'recording_G1', storage_strategy: 'managed' },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0 }, write_policy: { mode: 'append', dedupe_capability: 'receipt' }, migration: {},
    layout_version: 'managed_wide_v1', basic_managed_device_id: 'D', created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' };
  const seed: WriteGroup = { ...group, id: '', revision: '', applied_revision: '', status: 'draft', created_at: '0001-01-01T00:00:00Z', updated_at: '0001-01-01T00:00:00Z', destination: { ...group.destination, table_name: '' } };
  const lifecycle: WriteGroupLifecycle = { workspace_id: workspace, workspace_revision: 'W1', group_id: 'G1', group_revision: 'R1', applied_revision: '', status: 'deleted', connector_id: 'C', connector_revision: 'C1',
    actions: [{ action: 'recreate', allowed: true, reason_code: '', next_action: '', requires_review: true }],
    history: { canonical: true, lineage_id: 'L', lineage_revision: 1, current_group_id: 'G1', generations: [{ group_id: 'G1', generation: 1, current: true, reserved: false, status: 'deleted', group_revision: 'R1', applied_revision: '', saved_destination: { connector_id: 'C', connector_revision: 'C1', destination_scope: HASH }, revisions: [] }] },
    successor_seed: { basis_kind: 'saved', basis_revision: 'R1', basis_digest: HASH, group: seed } };
  const envelope: S.WriteGroupSuccessorEnvelope = { version: 'successor-review-v1', predecessor_group_id: 'G1', generation: 1, source_applied_revision: '', source_definition_digest: HASH,
    request: { workspace_id: workspace, expected_workspace_revision: 'W1', lineage_id: 'L', expected_lineage_revision: 1, expected_predecessor_revision: 'R1', expected_predecessor_connector_revision: 'C1', expected_target_connector_revision: 'C1', operation_id: 'ROOT', boundary: BOUNDARY, group: seed } };
  const review: S.WriteGroupSuccessorReview = { envelope, review_digest: HASH, history: lifecycle.history, blockers: [], required_checks: [] };
  const reserved = { ...group, id: 'G2', revision: 'R2', applied_revision: '', status: 'draft' as const, destination: { ...group.destination, table_name: 'recording_G2' } };
  const preparation: S.WriteGroupSuccessorPreparation = { workspace_revision: 'W2', lineage_id: 'L', lineage_revision: 2, generation: 2, predecessor_group_id: 'G1', operation_id: 'ROOT', intent_digest: HASH, group: reserved };
  const current = { workspace_revision: 'W2', lineage_revision: 2, predecessor_revision: 'R1', predecessor_connector_revision: 'C1', group_id: 'G2', group_revision: 'R2', target_connector_revision: 'C1' };
  const stage: S.WriteGroupSuccessorStageRequest = { envelope, review_digest: HASH, current };
  const preview: SchemaPreviewToken = { token: 'TOKEN', operation_id: 'SCHEMA', action: 'schema_apply', workspace_id: workspace, workspace_revision: 'W2', plan_id: 'G2', plan_revision: 'R2', connector_id: 'C', connector_revision: 'C1', table_prefix: '', statements: ['PRIVATE SQL'], tables: [{ name: 'recording_G2', action: 'create', columns: ['value'] }], group_layout: { table_name: 'recording_G2', owner_column: '__owner', columns: [{ name: 'value', sql_type: 'REAL', nullable: true, primary_key: false }] }, digest: HASH, expires_at: BOUNDARY, created_at: '2026-01-01T00:00:00Z' };
  const activation: S.WriteGroupSuccessorActivationReadiness = { review: { envelope: { version: 'successor-activation-review-v2', activation_kind: 'first_start', request: { stage, boundary: BOUNDARY }, source_applied_revision: '', source_definition_digest: HASH, target_definition_digest: HASH, counter_proof_digest: HASH, schema_revision: 'SCHEMA', schema_digest: HASH }, review_digest: HASH, history: lifecycle.history, blockers: [] }, readiness: { workspace_id: workspace, workspace_revision: 'W2', group_id: 'G2', group_revision: 'R2', applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] }, readiness_digest: HASH };
  const status: S.WriteGroupSuccessorOperationStatus = { workspace_id: workspace, lineage_id: 'L', predecessor_group_id: 'G1', operation_id: 'ROOT', reserve: { status: 'succeeded', result: preparation }, schema_operations: [{ operation_id: 'SCHEMA', status: 'preview_not_submitted' }], cutover: { status: 'not_submitted' }, first_start: { status: 'not_submitted' } };
  const api = {
    review: vi.fn<typeof studioV2WriteGroupSuccessorAPI.review>().mockResolvedValue(review), reserve: vi.fn<typeof studioV2WriteGroupSuccessorAPI.reserve>().mockResolvedValue(preparation),
    schemaPreview: vi.fn<typeof studioV2WriteGroupSuccessorAPI.schemaPreview>().mockResolvedValue(preview), schemaConfirm: vi.fn<typeof studioV2WriteGroupSuccessorAPI.schemaConfirm>().mockResolvedValue({ operation_id: 'SCHEMA', status: 'succeeded', executed_statements: 1, created_at: BOUNDARY, updated_at: BOUNDARY }),
    activationReview: vi.fn<typeof studioV2WriteGroupSuccessorAPI.activationReview>().mockResolvedValue(activation.review), activationReadiness: vi.fn<typeof studioV2WriteGroupSuccessorAPI.activationReadiness>().mockResolvedValue(activation),
    confirmCutover: vi.fn<typeof studioV2WriteGroupSuccessorAPI.confirmCutover>(), confirmFirstStart: vi.fn<typeof studioV2WriteGroupSuccessorAPI.confirmFirstStart>().mockResolvedValue({ operation_id: 'FIRST', root_operation_id: 'ROOT', workspace_revision: 'W3', lineage_id: 'L', lineage_revision: 3, predecessor_group_id: 'G1', group_id: 'G2', applied_revision: 'R2', boundary: BOUNDARY, generation: 2, status: 'succeeded' }),
    operationStatus: vi.fn<typeof studioV2WriteGroupSuccessorAPI.operationStatus>().mockResolvedValue(status),
  };
  const lifecycleRead = vi.fn<RecordingChangeDependencies['lifecycle']>().mockImplementation(async id => {
    if (!api.reserve.mock.calls.length) return structuredClone(lifecycle);
    return { ...lifecycle, workspace_revision: 'W2', group_id: id, group_revision: id === 'G2' ? 'R2' : 'R1', history: { ...lifecycle.history, lineage_revision: 2, reserved_group_id: 'G2' } };
  });
  const deps: RecordingChangeDependencies = { api, lifecycle: lifecycleRead, getGroup: vi.fn().mockResolvedValue({ workspace_revision: 'W2', group: reserved }), onChanged: vi.fn() };
  const intent = { action: 'recreate' as const, workspaceId: workspace, groupId: 'G1', workspaceRevision: 'W1', groupRevision: 'R1' };
  return { workspace, group, seed, lifecycle, envelope, review, preparation, stage, preview, activation, status, api, deps, intent, lifecycleRead };
}
