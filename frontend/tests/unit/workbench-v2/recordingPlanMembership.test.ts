import { describe, expect, it } from 'vitest';
import { resolveRecordingPlanScope } from '@/features/datalink/workbench-v2/steps/step4/recordingPlanMembership';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { MeasurementDefinition } from '@/types/measurement';

function fixture() {
  const state: WorkbenchV2State = structuredClone(INITIAL_STATE);
  state.devices = ['a', 'b'].map(id => ({ id, name: id, description: '', protocol: 'modbus_tcp', config: {},
    status: 'draft', test: null, persisted: true, save_state: 'saved' }));
  state.points = ['a', 'b'].map(id => ({ id: `point-${id}`, device_id: id, rule_id: `rule-${id}`, rule_name: id,
    name: id, address: '40001', data_type: 'float32', function: 'holding_register', width: 2,
    enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 }));
  const definitions: MeasurementDefinition[] = ['a', 'b'].map(id => ({
    id: `measurement-${id}`, workspace_id: 'workspace-a', device_id: id, point_id: `point-${id}`,
    equipment_id: `equipment-${id}`, definition_revision: 'definition-1', source_binding_revision: 'binding-1',
    series_epoch: 'epoch-1', name: id, quantity: 'temperature', semantic_kind: 'gauge',
  }));
  return { state, definitions };
}

describe('persisted recording membership join', () => {
  it('keeps each saved measurement and equipment rather than using point ids', () => {
    const { state, definitions } = fixture();
    expect(resolveRecordingPlanScope(state, 'workspace-a', definitions).members).toEqual([
      { member_id: 'member-measurement-a', measurement_id: 'measurement-a', equipment_id: 'equipment-a', name: 'a' },
      { member_id: 'member-measurement-b', measurement_id: 'measurement-b', equipment_id: 'equipment-b', name: 'b' },
    ]);
  });
  it('blocks incomplete, foreign or ambiguous definitions without fabricating members', () => {
    const { state, definitions } = fixture();
    for (const list of [[], [definitions[0]], definitions.map(definition => ({ ...definition, workspace_id: 'foreign' })), [...definitions, { ...definitions[0], id: 'duplicate' }]]) {
      expect(resolveRecordingPlanScope(state, 'workspace-a', list).members).toEqual([]);
    }
  });
  it('excludes disabled and removed points and rejects unsaved devices', () => {
    const { state, definitions } = fixture();
    state.points[1].enabled = false;
    expect(resolveRecordingPlanScope(state, 'workspace-a', definitions).members.map(member => member.measurement_id)).toEqual(['measurement-a']);
    state.points = [];
    expect(resolveRecordingPlanScope(state, 'workspace-a', definitions).members).toEqual([]);
    const next = fixture(); next.state.devices[0].save_state = 'saving';
    expect(resolveRecordingPlanScope(next.state, 'workspace-a', next.definitions).members).toEqual([]);
  });
  it('changes the scope version when a saved definition or source binding changes', () => {
    const { state, definitions } = fixture();
    const before = resolveRecordingPlanScope(state, 'workspace-a', definitions).fingerprint;
    definitions[0].definition_revision = 'definition-2';
    expect(resolveRecordingPlanScope(state, 'workspace-a', definitions).fingerprint).not.toEqual(before);
  });
});
