import type { ModbusShareStatus } from '../../../src/types/modbusShare';
import type { StudioV2ShareActivationContext } from '../../../src/features/datalink/workbench-v2/state/studioV2ShareActivation';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

export function mcShareState(overrides: Partial<WorkbenchV2State> = {}): WorkbenchV2State {
  const points = Array.from({ length: 4 }, (_, index) => ({
    id: `mc-rule-p-${index}`,
    device_id: 'dev-mc',
    rule_id: 'mc-rule',
    rule_name: 'MC D registers',
    name: `D_${index}`,
    address: `D${index}`,
    data_type: 'int16',
    function: 'D (Word)',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  }));
  const mappings = Object.fromEntries(points.map((point, index) => [point.id, {
    point_id: point.id,
    workspace_id: 'workspace-1',
    rule_id: point.rule_id,
    device_id: point.device_id,
    address: point.address,
    persisted_point_id: `persisted-point-${index}`,
    tag_key: `mc.d${index}`,
    display_name: `D${index}`,
    unit: '',
    target_type: 'int16' as const,
    scale: 1,
    offset: 0,
    enabled: true,
    tag_id: `tag-d${index}`,
    persisted: true,
  }]));

  return {
    ...INITIAL_STATE,
    devices: [{ ...INITIAL_STATE.devices[0], id: 'dev-mc', protocol: 'mc_3e', name: 'MC PLC' }],
    rules: [{
      ...INITIAL_STATE.rules[0],
      id: 'mc-rule',
      device_id: 'dev-mc',
      workspace_id: 'workspace-1',
      persisted: true,
      revision_id: 'mc-rule-revision-1',
      name: 'MC D registers',
      start_address: 'D0',
      count: 4,
      share_enabled: true,
      share_start_register: 40001,
      share_stride: 1,
    }],
    points,
    mappings,
    db: { ...INITIAL_STATE.db, targets: {} },
    ...overrides,
  };
}

const canonicalDesiredMappings = Array.from({ length: 4 }, (_, index) => ({
  workspace_id: 'workspace-1',
  source_rule_id: 'mc-rule',
  source_rule_revision: 'mc-rule-revision-1',
  tag_id: `tag-d${index}`,
  tag_key: `mc.d${index}`,
  display_name: `D${index}`,
  data_type: 'int16' as const,
  share_start_register: 40001 + index,
  zero_based_register: index,
  span_registers: 1,
  stride_registers: 1,
}));

export const hydratedShareStatus: ModbusShareStatus = {
  workspace_id: 'workspace-1',
  enabled: true,
  port: 15020,
  address: '127.0.0.1:15020',
  bind_address: '127.0.0.1',
  bind_state: 'pass',
  lifecycle_state: 'running',
  mapping_count: 4,
  hydration_state: 'ready',
  readiness: true,
  workspace_revision: 'workspace-revision-1',
  settings_revision: 'settings-revision-1',
  readiness_token: 'readiness-token-1',
  canonical_desired_mappings: canonicalDesiredMappings,
};

export const shareContext: StudioV2ShareActivationContext = {
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-revision-1',
  settings_revision: 'settings-revision-1',
  readiness_token: 'readiness-token-1',
  canonical_plan: {
    workspace_id: 'workspace-1',
    workspace_revision: 'workspace-revision-1',
    settings_revision: 'settings-revision-1',
    signature: 'canonical-plan-signature-1',
    desired_mappings: canonicalDesiredMappings,
  },
};
