import type { WriteGroup, WriteGroupReadiness } from '../../../src/types/studioV2WriteGroup';
import type { WriteGroupLifecycle, WriteGroupLifecycleHistory } from '../../../src/types/studioV2WriteGroupLifecycle';
import type * as S from '../../../src/types/studioV2WriteGroupSuccessor';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { wide, wideMeanTarget, DIGEST } from '../../fixtures/cutoverFixtures';

// TEST DATA ONLY. These identities represent no database, device, or production backend.
export const lifecycleScenarios = ['draft-wide', 'applied-wide', 'disabled-wide', 'deleted-wide', 'applied-long', 'deleted-long', 'deleted-unapplied', 'deleted-basic', 'blocked-wide', 'review-blocked-wide', 'stale-wide'] as const;
export type LifecycleScenario = typeof lifecycleScenarios[number];
export const FIXTURE_HASH = DIGEST;
export const FIXTURE_TIME = '2026-10-01T00:00:00Z';
export const FIXTURE_PATH = '/api/v1/datalink/studio-v2/workspace';
export function scenarioFrom(value: string | null): LifecycleScenario {
  return lifecycleScenarios.find(item => item === value) ?? 'applied-wide';
}
export function recordingLifecycleHostState() {
  const state = savedGroupState();
  state.devices[0].name = 'Mock PLC A';
  state.db.connector.database = 'mock-only.db';
  state.settings.connectors = ['conn-1', 'conn-2'].map((id, index) => ({
    id, identity_revision: `crev-${index + 1}`, name: `Mock SQLite ${index + 1}`, kind: 'sqlite', host: '', port: 0,
    database: `mock-only-${index + 1}.db`, username: '', schema: '', table: '', enabled: true, status: 'unknown', default_write_interval_seconds: 15,
  }));
  return state;
}
export function initialLifecycleGroup(scenario: LifecycleScenario): WriteGroup {
  const unapplied = scenario === 'draft-wide' || scenario === 'deleted-unapplied';
  const group = (scenario === 'deleted-basic' ? wide : wideMeanTarget)({ id: 'G1', revision: 'R1', applied_revision: unapplied ? '' : 'APPLIED1', name: 'Mock saved recording',
    status: scenario.startsWith('deleted') ? 'deleted' : scenario.startsWith('disabled') ? 'disabled' : unapplied ? 'draft' : 'running' });
  group.destination.database = 'mock-only.db';
  group.destination.table_name = 'mock_recording_G1';
  if (scenario.endsWith('long')) { group.layout_version = 'managed_long_v1'; group.storage_mode = 'plain'; }
  return group;
}
export function successorSeed(group: WriteGroup): WriteGroup {
  const seed = structuredClone(group);
  Object.assign(seed, { id: '', revision: '', applied_revision: '', status: 'draft', migration: {}, created_at: '0001-01-01T00:00:00Z', updated_at: '0001-01-01T00:00:00Z' });
  seed.destination.table_name = '';
  delete seed.destination.schema_revision;
  delete seed.destination.schema_digest;
  return seed;
}
export interface LifecycleMockState {
  scenario: LifecycleScenario;
  source: WriteGroup;
  target?: WriteGroup;
  workspaceRevision: string;
  lineageRevision: number;
  envelope?: S.WriteGroupSuccessorEnvelope;
  preparation?: S.WriteGroupSuccessorPreparation;
  activation?: S.WriteGroupSuccessorActivationReview;
  status?: S.WriteGroupSuccessorOperationStatus;
}
export function fixtureHistory(state: LifecycleMockState): WriteGroupLifecycleHistory {
  const terminal = state.status?.cutover.status === 'succeeded' || state.status?.first_start.status === 'succeeded';
  const boundary = state.status?.cutover.result?.boundary ?? state.status?.first_start.result?.boundary;
  const current = terminal ? 'G2' : 'G1';
  return { canonical: true, lineage_id: 'L', lineage_revision: state.lineageRevision, current_group_id: current,
    ...(state.target && !terminal ? { reserved_group_id: 'G2' } : {}),
    generations: [state.source, ...(state.target ? [state.target] : [])].map((group, index) => ({
      group_id: group.id, generation: index + 1, current: group.id === current, reserved: group.id === 'G2' && !terminal,
      ...(index ? { predecessor_group_id: 'G1' } : {}), status: group.status, group_revision: group.revision, applied_revision: group.applied_revision,
      saved_destination: { connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision, destination_scope: DIGEST },
      revisions: group.applied_revision ? [{ group_revision: group.applied_revision, effective_at: index && boundary ? boundary : FIXTURE_TIME,
        ...(!index && state.status?.cutover.status === 'succeeded' ? { ends_at: boundary } : {}),
        sources: group.members.map(({ device_id, point_id, tag_id }) => ({ device_id, point_id, tag_id })),
        destination: { connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision, destination_scope: DIGEST },
        pending: index ? { samples: 0, outbox: { outstanding: 0, total: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, handoffs: [] }
          : { samples: 3, outbox: { outstanding: 4, total: 7, queued: 1, retrying: 0, blocked: 2, quarantined: 0, unknown: 1, sql_committed: 3, skipped: 0 }, handoffs: [] },
      }] : [],
    })),
  };
}
export function fixtureLifecycle(state: LifecycleMockState, group: WriteGroup): WriteGroupLifecycle {
  const blocked = state.scenario === 'blocked-wide';
  const historical = group.id !== fixtureHistory(state).current_group_id;
  const reserved = Boolean(fixtureHistory(state).reserved_group_id);
  const allowed = (action: string) => (!blocked || action === 'stop') && !historical && !reserved && (group.status === 'deleted' ? action === 'recreate'
    : action === 'edit' || action === 'remove' || action === 'stop' && group.status !== 'disabled' || action === 'transition' && Boolean(group.applied_revision));
  return { workspace_id: 'ws-1', workspace_revision: state.workspaceRevision, group_id: group.id, group_revision: group.revision,
    applied_revision: group.applied_revision, status: group.status, connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision,
    actions: (['edit', 'stop', 'remove', 'recreate', 'transition'] as const).map(action => ({ action, allowed: allowed(action),
      reason_code: allowed(action) ? '' : blocked ? 'connector_identity_stale' : historical ? 'historical_generation' : reserved ? 'successor_reserved'
        : group.status === 'deleted' ? 'recording_removed' : action === 'stop' ? 'already_stopped' : action === 'recreate' ? 'recording_not_removed' : 'source_draft_not_active',
      next_action: allowed(action) ? '' : blocked ? action === 'remove' ? 'stop' : 'review_connector' : 'review_lifecycle', requires_review: action === 'recreate' || action === 'transition' })),
    history: fixtureHistory(state), successor_seed: { basis_kind: group.applied_revision ? 'applied' : 'saved',
      basis_revision: group.applied_revision || group.revision, basis_digest: DIGEST, group: successorSeed(group) },
  };
}
export function fixtureReadiness(state: LifecycleMockState, group: WriteGroup): WriteGroupReadiness {
  return { workspace_id: 'ws-1', workspace_revision: state.workspaceRevision, group_id: group.id, group_revision: group.revision,
    applied_revision: group.applied_revision, config_ready: true, schema_ready: true, schema_digest: DIGEST, ready: true, issues: [] };
}
export function fixtureCurrent(state: LifecycleMockState): S.WriteGroupSuccessorCurrentCAS {
  if (!state.target) throw new Error('Mock target not reserved');
  return { workspace_revision: state.workspaceRevision, lineage_revision: state.lineageRevision, predecessor_revision: state.source.revision,
    predecessor_connector_revision: state.source.destination.connector_revision, group_id: 'G2', group_revision: state.target.revision,
    target_connector_revision: state.target.destination.connector_revision };
}
