import { beforeEach, vi } from 'vitest';

import type { StudioV2WorkspaceDeviceRecord } from '../../../src/services/studioV2WorkspaceDevices';
import type {
  SourceRuleRecord,
  StudioV2WorkspaceMappingRecord,
} from '../../../src/types/datalink';
import type {
  StudioV2RuntimeApplyStatus,
  StudioV2RuntimeAppliedRecord,
} from '../../../src/types/studioV2RuntimeApply';
import type { StudioV2Workspace as WorkspaceType } from '../../../src/types/studioV2Workspace';

export interface WorkbenchV2ShellMockProps<S> {
  state: S;
  actions: { dispatch: (action: unknown) => void };
}

export const studioV2ServiceMocks = {
  workspace: { get: vi.fn() },
  devices: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    updateAvailability: vi.fn(),
    remove: vi.fn(),
    updateOrder: vi.fn(),
  },
  rules: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
  mappings: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
  database: {
    getConfig: vi.fn(),
    updateConfig: vi.fn(),
    listTargets: vi.fn(),
    upsertTarget: vi.fn(),
  },
};

export function workspaceFixture(overrides: Partial<WorkspaceType> = {}): WorkspaceType {
  return {
    id: 'workspace-1',
    kind: 'single',
    status: 'ready',
    ordered_device_ids: ['dev-A'],
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function deviceFixture(
  overrides: Partial<StudioV2WorkspaceDeviceRecord> = {},
): StudioV2WorkspaceDeviceRecord {
  return {
    id: 'dev-A',
    name: 'Line A PLC',
    description: '',
    protocol: 'modbus_tcp',
    status: 'draft',
    connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
    last_test_at: null,
    last_test_success: null,
    last_test_error: '',
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function ruleFixture(overrides: Partial<SourceRuleRecord> = {}): SourceRuleRecord {
  return {
    id: 'rule-A',
    device_id: 'dev-A',
    workspace_id: 'workspace-1',
    start_address: '40001',
    count: 1,
    data_type: 'int16',
    naming_prefix: 'A_',
    enabled: true,
    locked: false,
    origin: 'manual',
    skipped_addresses: [],
    revision_id: 'rev-A',
    scale_multiplier: 1,
    scale_offset: 0,
    data_format: '',
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function mappingFixture(
  overrides: Partial<StudioV2WorkspaceMappingRecord> = {},
): StudioV2WorkspaceMappingRecord {
  return {
    id: 'mapping-A',
    workspace_id: 'workspace-1',
    point_id: 'point-A',
    rule_id: 'rule-A',
    device_id: 'dev-A',
    address: '40001',
    tag_id: 'tag-A',
    tag_key: 'line.a.persisted',
    display_name: 'Line A Temp',
    unit: 'C',
    target_type: 'float64',
    scale: 1,
    offset: 0,
    enabled: true,
    save_state: 'saved',
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function runtimeApplied<T extends object>(
  record: T,
  status: StudioV2RuntimeApplyStatus = 'not_running',
  message?: string,
): StudioV2RuntimeAppliedRecord<T> {
  return { ...record, runtime_apply_status: status, runtime_apply_message: message };
}

export const mappingPoints = [
  {
    id: 'rule-A-p-0',
    device_id: 'dev-A',
    rule_id: 'rule-A',
    rule_name: 'Line A Registers',
    name: 'A_0',
    address: '40001',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
  {
    id: 'rule-B-p-0',
    device_id: 'dev-B',
    rule_id: 'rule-B',
    rule_name: 'Line B Registers',
    name: 'B_0',
    address: '40011',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
];

export const shiftedMappingPoints = [
  {
    ...mappingPoints[0],
    address: '40021',
    name: 'A_20',
  },
  mappingPoints[1],
];

export function registerMappingAutosaveFixtures() {
  beforeEach(() => {
    vi.clearAllMocks();
    studioV2ServiceMocks.workspace.get.mockResolvedValue(
      workspaceFixture({ ordered_device_ids: ['dev-B', 'dev-A'] }),
    );
    studioV2ServiceMocks.devices.list.mockResolvedValue([
      deviceFixture({ id: 'dev-B', name: 'Line B PLC', connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}' }),
      deviceFixture(),
    ]);
    studioV2ServiceMocks.rules.list.mockResolvedValue([
      ruleFixture({
        id: 'rule-B',
        device_id: 'dev-B',
        start_address: '40011',
        naming_prefix: 'B_',
        revision_id: 'rev-B',
      }),
      ruleFixture(),
    ]);
    studioV2ServiceMocks.mappings.list.mockResolvedValue([]);
    studioV2ServiceMocks.database.getConfig.mockResolvedValue(null);
    studioV2ServiceMocks.database.listTargets.mockResolvedValue([]);
  });
}
