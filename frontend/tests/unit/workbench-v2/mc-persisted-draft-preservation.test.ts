import { describe, expect, it } from 'vitest';
import type {
  StudioV2WorkspaceDatabaseTargetRecord,
  StudioV2WorkspaceMappingRecord,
} from '../../../src/types/datalink';
import {
  mergePersistedMCState,
  type PersistedMCRuleState,
} from '../../../src/features/datalink/workbench-v2/state/mcPersistedPoints';
import type { Point, Rule, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE, workbenchV2Reducer } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

const scope = {
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-rev-1',
  revision_id: 'rule-mc-rev-1',
};

const persistedPointIdentities = [
  { point_id: 'mc-persisted', address: 'W0', tag_id: 'legacy-tag', mapping_id: 'mapping-legacy' },
  { point_id: 'mc-new', address: 'W1' },
];

function persistedEvidence(status: PersistedMCRuleState['status'], points = persistedPointIdentities): PersistedMCRuleState {
  return { status, scope, points, ...(status === 'unavailable' ? { reason: 'candidate_get_failed' } : {}) };
}

function point(id: string, ruleID: string, deviceID: string, address: string): Point {
  return {
    id,
    device_id: deviceID,
    rule_id: ruleID,
    rule_name: ruleID === 'rule-mc' ? 'MC words' : 'Modbus words',
    name: id,
    address,
    data_type: 'int16',
    function: 'W (Word)',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  };
}

function mcRule(): Rule {
  return {
    id: 'rule-mc',
    device_id: 'dev-mc',
    workspace_id: scope.workspace_id,
    revision_id: scope.revision_id,
    name: 'MC words',
    start_address: 'W0',
    count: 2,
    data_type: 'int16',
    naming_prefix: 'MC_',
    enabled: true,
    scale_multiplier: 1,
    scale_offset: 0,
    data_format: '',
    skipped_addresses: [],
    share_enabled: false,
    share_start_register: null,
    share_stride: null,
    persisted: true,
    save_state: 'saved',
    mc_persisted_points: persistedEvidence('compatible'),
  };
}

function modbusRule(): Rule {
  return {
    id: 'rule-modbus',
    device_id: 'dev-modbus',
    name: 'Modbus words',
    start_address: '40001',
    count: 1,
    data_type: 'int16',
    naming_prefix: 'P1_',
    enabled: true,
    scale_multiplier: 1,
    scale_offset: 0,
    data_format: '',
    skipped_addresses: [],
    share_enabled: false,
    share_start_register: null,
    share_stride: null,
    persisted: true,
    save_state: 'saved',
  };
}

function mappingRecord(): StudioV2WorkspaceMappingRecord {
  return {
    id: 'mapping-legacy',
    workspace_id: scope.workspace_id,
    point_id: 'mc-persisted',
    rule_id: 'rule-mc',
    device_id: 'dev-mc',
    address: 'W0',
    tag_id: 'legacy-tag',
    tag_key: 'server.legacy',
    display_name: 'Legacy',
    unit: '',
    target_type: 'int16',
    scale: 1,
    offset: 0,
    enabled: true,
    created_at: '2026-10-30T00:00:00.000Z',
    updated_at: '2026-10-30T00:00:00.000Z',
  };
}

function targetRecord(): StudioV2WorkspaceDatabaseTargetRecord {
  return {
    id: 'target-legacy',
    workspace_id: scope.workspace_id,
    point_id: 'mc-persisted',
    tag_id: 'legacy-tag',
    column_name: 'server_column',
    enabled: true,
    created_at: '2026-10-30T00:00:00.000Z',
    updated_at: '2026-10-30T00:00:00.000Z',
  };
}

function stateWithDrafts(): WorkbenchV2State {
  const points = [
    point('mc-persisted', 'rule-mc', 'dev-mc', 'W0'),
    point('mc-new', 'rule-mc', 'dev-mc', 'W1'),
    point('p1', 'rule-modbus', 'dev-modbus', '40001'),
  ];
  let state: WorkbenchV2State = {
    ...INITIAL_STATE,
    devices: [
      { ...INITIAL_STATE.devices[0], id: 'dev-mc', protocol: 'mc_3e', persisted: true },
      { ...INITIAL_STATE.devices[0], id: 'dev-modbus', protocol: 'modbus_tcp', persisted: true },
    ],
    rules: [mcRule(), modbusRule()],
    points,
    mappings: {},
    db: { ...INITIAL_STATE.db, targets: {} },
  };
  state = workbenchV2Reducer(state, { type: 'initMappingsForPoints', points });
  state = workbenchV2Reducer(state, {
    type: 'updateMapping',
    pointId: 'mc-new',
    patch: { tag_key: 'operator.custom', save_state: 'draft-invalid' },
  });
  state = workbenchV2Reducer(state, {
    type: 'updateMapping',
    pointId: 'p1',
    patch: { tag_key: 'p1', save_state: 'draft-invalid' },
  });

  const legacyMapping = state.mappings['mc-persisted'];
  if (!legacyMapping) throw new Error('legacy mapping fixture was not initialized');
  state = {
    ...state,
    mappings: {
      ...state.mappings,
      'mc-persisted': {
        ...legacyMapping,
        mapping_id: 'mapping-legacy',
        persisted_point_id: 'mc-persisted',
        persisted: true,
        tag_key: 'operator.legacy.latest',
        local_value: { ...legacyMapping.local_value!, tag_key: 'operator.legacy.latest' },
        persisted_value: { ...legacyMapping.local_value!, tag_key: 'server.legacy' },
        save_state: 'saving',
      },
    },
  };
  state = workbenchV2Reducer(state, {
    type: 'upsertDbTarget',
    pointId: 'mc-persisted',
    target: {
      tag_id: 'legacy-tag',
      column_name: 'operator_column_latest',
      enabled: true,
      point_id: 'mc-persisted',
      row_id: 'target-legacy',
      persisted: true,
      save_state: 'saving',
    },
  });
  state = workbenchV2Reducer(state, {
    type: 'upsertDbTarget',
    pointId: 'mc-new',
    target: {
      tag_id: 'new-tag',
      column_name: 'operator_column',
      enabled: true,
      point_id: 'mc-new',
      persisted: false,
      save_state: 'draft-invalid',
    },
  });
  state = workbenchV2Reducer(state, {
    type: 'upsertDbTarget',
    pointId: 'p1',
    target: {
      tag_id: 'p1-tag',
      column_name: 'p1',
      enabled: true,
      point_id: 'p1',
      persisted: false,
      save_state: 'draft-invalid',
    },
  });
  return state;
}

type MergedMCState = ReturnType<typeof mergePersistedMCState>;

function expectDrafts(state: MergedMCState): void {
  expect(state.points.map((item) => item.id)).toEqual(['mc-persisted', 'mc-new', 'p1']);
  expect(state.mappings['mc-new']?.tag_key).toBe('operator.custom');
  expect(state.db.targets['mc-new']?.column_name).toBe('operator_column');
  expect(state.mappings['mc-persisted']?.tag_key).toBe('operator.legacy.latest');
  expect(state.mappings['mc-persisted']?.mapping_id).toBe('mapping-legacy');
  expect(state.db.targets['mc-persisted']?.column_name).toBe('operator_column_latest');
  expect(state.db.targets['mc-persisted']?.row_id).toBe('target-legacy');
  expect(state.mappings.p1?.tag_key).toBe('p1');
  expect(state.db.targets.p1?.column_name).toBe('p1');
}

describe('persisted MC draft preservation', () => {
  it('retains local rows without server records through checking and unavailable until compatible evidence returns', () => {
    const records = [mappingRecord()];
    const targets = [targetRecord()];
    let state = stateWithDrafts();
    const checking = mergePersistedMCState(
      state,
      { 'rule-mc': persistedEvidence('checking', []) },
      records,
      targets,
    );
    expect(checking.rules[0]?.mc_persisted_points?.status).toBe('checking');
    expectDrafts(checking);
    state = { ...state, ...checking };

    const unavailable = mergePersistedMCState(
      state,
      { 'rule-mc': persistedEvidence('unavailable', []) },
      records,
      targets,
    );
    expect(unavailable.rules[0]?.mc_persisted_points?.status).toBe('unavailable');
    expectDrafts(unavailable);
    state = { ...state, ...unavailable };

    const compatible = mergePersistedMCState(
      state,
      { 'rule-mc': persistedEvidence('compatible') },
      records,
      targets,
    );
    expect(compatible.rules[0]?.mc_persisted_points?.status).toBe('compatible');
    expectDrafts(compatible);
  });
});
