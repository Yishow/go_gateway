import type { Device, SourceRuleRecord, StudioV2WorkspaceMappingRecord } from '@/types/datalink';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
import { savedGroup } from '../../fixtures/writeGroupState';

const at = '2026-10-29T00:00:00Z';
export const nativeDevice = (id = 'D'): Device => ({ id, name: `Device ${id}`, description: '', protocol: 'modbus_tcp', status: 'draft',
  connection_config: '{"host":"device.example","port":502,"slave_id":1}', last_test_at: null, last_test_success: null,
  last_test_error: '', created_at: at, updated_at: at });
export const nativeRule = (id = 'R', device = 'D'): SourceRuleRecord => ({ id, workspace_id: 'W', device_id: device,
  revision_id: 'R2', start_address: '40001', count: 2, data_type: 'int16', naming_prefix: `Rule ${id}`, enabled: true,
  skipped_addresses: [], locked: false, origin: 'manual', created_at: at, updated_at: at });
export const nativeMapping = (id = 'M'): StudioV2WorkspaceMappingRecord => ({ id, workspace_id: 'W', device_id: 'D', rule_id: 'R',
  point_id: 'P-native', tag_id: 'T', address: '40002', tag_key: 'temperature', display_name: `Mapping ${id}`, unit: 'C',
  target_type: 'float64', scale: 0.1, offset: 0, enabled: true, created_at: at, updated_at: at });
export function locatorRead(step: 1 | 2 | 3 | 4 = 1, id = step === 4 ? 'G' : 'D'): StudioLocatorEntryRead {
  const group = savedGroup({ id, workspace_id: 'W', revision: 'G3', name: 'Shared recorder' });
  const device = step === 4 ? undefined : id;
  return { locator: { workspace_id: 'W', step, device_id: device ?? null, rule_id: step === 2 ? 'R' : null,
    mapping_id: step === 3 ? 'M' : null, group_id: step === 4 ? id : null, group_revision: step === 4 ? 'G1' : null },
    workspace: { id: 'W', kind: 'single', status: 'ready', ordered_device_ids: ['D', 'B'], created_at: at, updated_at: at },
    entry: { workspace_revision: 'W2', hydration_scope: { complete_workspace: false, device_ids: device ? [device] : [],
      rule_ids: step === 2 || step === 3 ? ['R'] : [], mapping_ids: step === 3 ? ['M'] : [], group_ids: step === 4 ? [id] : [] },
      ...(device ? { devices: [nativeDevice(device)] } : {}),
      ...(step === 2 || step === 3 ? { source_rules: [nativeRule('R', device)] } : {}),
      ...(step === 3 ? { mappings: [{ ...nativeMapping(), device_id: device! }] } : {}),
      ...(step === 4 ? { write_groups: { workspace_id: 'W', workspace_revision: 'W2', groups: [group] } } : {}),
      rule_heads: step === 2 || step === 3 ? { R: 'R2' } : {}, group_heads: step === 4 ? { [id]: 'G3' } : {},
      revalidated_mapping_ids: step === 3 ? ['M'] : [] },
    selection: { state: 'selected', target: { workspaceId: 'W', workspaceRevision: 'W2', step,
      ...(device ? { deviceId: device } : {}), ...(step === 2 || step === 3 ? { ruleId: 'R', ruleRevision: 'R2' } : {}),
      ...(step === 3 ? { mappingId: 'M' } : {}), ...(step === 4 ? { groupId: id, groupRevision: 'G3' } : {}) },
      ...(step === 4 ? { history: { groupId: id, revisionHint: 'G1', relation: 'different' } } : {}) } };
}
