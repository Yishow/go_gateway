import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Mapping, Point, WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { WriteGroup } from '@/types/studioV2WriteGroup';

/** A workspace whose device, points, Tags and destination connector are all saved. */
export function savedGroupState(overrides: Partial<WorkbenchV2State> = {}): WorkbenchV2State {
  const device = { ...INITIAL_STATE.devices[0], id: 'dev-1', name: 'PLC A', persisted: true, save_state: 'saved' as const };
  const rule = { ...INITIAL_STATE.rules[0], id: 'rule-1', device_id: 'dev-1', enabled: true, persisted: true };
  const mk = (index: number, tagKey: string, targetType: Mapping['target_type']): [Point, Mapping] => {
    const id = `rule-1-p-${index}`;
    return [
      {
        id, device_id: 'dev-1', rule_id: 'rule-1', rule_name: 'Rule 1', name: tagKey, address: `4000${index + 1}`, data_type: targetType,
        function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0,
      },
      {
        point_id: id, tag_key: tagKey, display_name: tagKey, unit: '', target_type: targetType, scale: 1, offset: 0, enabled: true,
        persisted: true, save_state: 'saved', persisted_point_id: `pt-${index}`, tag_id: `tag-${index}`, device_id: 'dev-1',
      },
    ];
  };
  const entries = [mk(0, 'line.temperature', 'float64'), mk(1, 'line.pressure', 'float64'), mk(2, 'line.batch', 'string')];
  return {
    ...INITIAL_STATE,
    devices: [device],
    rules: [rule],
    points: entries.map(([point]) => point),
    mappings: Object.fromEntries(entries.map(([, mapping]) => [mapping.point_id, mapping])),
    db: {
      ...INITIAL_STATE.db,
      connector: {
        ...INITIAL_STATE.db.connector, kind: 'sqlite', name: 'Line DB', table: 'readings', schema: '', workspace_id: 'ws-1',
        connector_id: 'conn-1', identity_revision: 'crev-1', persisted: true, save_state: 'saved', setup_revision: 'setup-1',
      },
    },
    ...overrides,
  };
}

export function savedGroup(overrides: Partial<WriteGroup> = {}): WriteGroup {
  return {
    id: 'group-1', workspace_id: 'ws-1', revision: 'rev-1', applied_revision: '', name: 'Line A', status: 'ready',
    members: [
      { device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', source_revision: 's0', mapping_revision: 'm0', target_column: 'temperature', required: true },
      { device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1', source_revision: 's1', mapping_revision: 'm1', target_column: 'pressure', required: true },
    ],
    destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: '/d.db', table_schema: '', table_name: 'readings', storage_strategy: 'custom' },
    row_policy: { interval_seconds: 10, allowed_lateness_seconds: 0, incomplete_policy: 'skip_row' },
    write_policy: {}, migration: {}, created_at: '2026-10-02T00:00:00Z', updated_at: '2026-10-02T00:00:00Z',
    ...overrides,
  };
}

export const TABLE_COLUMNS = [
  { name: 'id', type: 'integer', nullable: false, primary_key: true },
  { name: 'temperature', type: 'double precision', nullable: true, primary_key: false },
  { name: 'pressure', type: 'double precision', nullable: true, primary_key: false },
  { name: 'batch', type: 'text', nullable: true, primary_key: false },
  { name: 'note', type: 'text', nullable: true, primary_key: false },
];
