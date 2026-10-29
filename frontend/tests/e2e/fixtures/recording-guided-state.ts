import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';
import type { RecordingStartOperation, RecordingStartRequest } from '../../../src/types/studioV2RecordingStart';
import { recordingAuthoringHostState, initialAuthoringGroups, AUTHORING_HASH } from './recording-authoring-state';
import { savedGroup } from '../../fixtures/writeGroupState';

// Test data only: these identities establish no SQL, device or effective-runtime evidence.
export const recordingGuidedScenarios = ['new', 'saved-long', 'saved-wide', 'saved-custom', 'saved-different', 'blocked-wide',
  'partial', 'unknown', 'long-sources', 'empty-device'] as const;
export type GuidedScenario = typeof recordingGuidedScenarios[number];
export const GUIDED_PATH = '/api/v1/datalink/studio-v2/workspace';
export function guidedScenarioFrom(value: string | null): GuidedScenario {
  return recordingGuidedScenarios.find(item => item === value) ?? 'new';
}
export function recordingGuidedHostState(scenario: GuidedScenario): WorkbenchV2State {
  const state = structuredClone(recordingAuthoringHostState('new-sqlite'));
  state.db.connector.database = '/mock-private/recordings.sqlite';
  state.settings.connectors[0].database = '/mock-private/recordings.sqlite';
  state.current = 4; state.completed = new Set([1, 2, 3]);
  state.devices.push({ ...state.devices[0], id: 'dev-empty', name: 'Empty device' });
  state.rules.push({ ...state.rules[0], id: 'rule-2', device_id: 'dev-2' });
  state.selectedRuleId = 'rule-2'; // Completion must not fall back from confirmed A to selected-rule B.
  if (scenario === 'long-sources') {
    state.points = Array.from({ length: 40 }, (_, index) => ({ ...state.points[0], id: `rule-1-p-${index}`,
      name: `Source ${index + 1}`, data_type: 'float64', address: String(40001 + index) }));
    const base = Object.values(state.mappings)[0];
    state.mappings = Object.fromEntries(state.points.map((point, index) => [point.id, { ...base, point_id: point.id,
      tag_key: `source.${index + 1}`, display_name: point.name, target_type: 'float64', persisted_point_id: `pt-${index}`, tag_id: `tag-${index}` }]));
  }
  return state;
}
export function initialGuidedGroups(scenario: GuidedScenario): WriteGroup[] {
  if (['new', 'long-sources', 'empty-device'].includes(scenario)) return [];
  if (scenario === 'saved-custom') return [savedGroup({ id: 'G1', revision: 'R1', applied_revision: 'APPLIED1',
    status: 'running', name: 'Saved custom recording', destination: { connector_id: 'conn-1', connector_revision: 'crev-1',
      database: '/mock-private/recordings.sqlite', table_schema: 'main', table_name: 'mock_custom_G1', storage_strategy: 'custom',
      schema_revision: 'R1', schema_digest: AUTHORING_HASH } })];
  const group = initialAuthoringGroups('saved-gauge')[0];
  group.destination.database = '/mock-private/recordings.sqlite';
  if (scenario === 'saved-different') {
    group.revision = 'R2'; group.applied_revision = 'R1';
    group.recording_profile = { kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' };
    group.row_policy.interval_seconds = 30;
  }
  group.destination.schema_revision = group.revision;
  if (scenario === 'saved-wide' || scenario === 'blocked-wide') { group.layout_version = 'managed_wide_profile_v1'; delete group.storage_mode; }
  return [group];
}
/** One canonical state/head. Operation records are immutable request receipts, not another group store. */
export interface GuidedMockState {
  scenario: GuidedScenario;
  workspaceRevision: string;
  groups: WriteGroup[];
  originalStart?: RecordingStartRequest;
  operation?: RecordingStartOperation;
}
export function advanceGuidedHead(state: GuidedMockState): string {
  state.workspaceRevision = `W${Number(state.workspaceRevision.slice(1)) + 1}`;
  return state.workspaceRevision;
}
