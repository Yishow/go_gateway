import { describe, expect, expectTypeOf, it } from 'vitest';
import {
  hydrateStudioV2DatabaseRowGroups,
  hydrateStudioV2DatabaseTarget,
  resolveStudioV2DatabaseTargetPointID,
  toStudioV2DatabaseConfigRequest,
  toStudioV2DatabaseTargetRequest,
} from '../../../src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import type { DbConnector, DbRowGroup, DbTarget } from '../../../src/features/datalink/workbench-v2/state/types';

describe('studio v2 database row-group autosave helpers', () => {
  it('keeps row groups separate from target refs at the type level', () => {
    expectTypeOf<DbRowGroup['member_point_ids']>().toEqualTypeOf<string[]>();
    expectTypeOf<DbTarget['row_group_id']>().toEqualTypeOf<string | undefined>();
  });

  it('hydrates row groups from workspace database config', () => {
    const groups = hydrateStudioV2DatabaseRowGroups([
      {
        id: 'group-shared-temp',
        connector_id: 'db-1',
        table_schema: 'main',
        table_name: 'sensor_values',
        member_point_ids: ['point-A', 'point-B'],
        group_key_columns: ['ts', 'line_id'],
      },
    ]);

    expect(groups).toEqual<DbRowGroup[]>([
      {
        id: 'group-shared-temp',
        connector_id: 'db-1',
        table_schema: 'main',
        table_name: 'sensor_values',
        member_point_ids: ['point-A', 'point-B'],
        group_key_columns: ['ts', 'line_id'],
      },
    ]);
  });

  it('preserves row groups and target references in autosave requests', () => {
    const connector: DbConnector = {
      kind: 'sqlite',
      name: 'Line A SQLite',
      host: '',
      port: 0,
      database: '/tmp/target.db',
      username: '',
      schema: 'main',
      table: 'sensor_values',
      write_mode: 'insert',
      write_interval_seconds: 5,
      timestamp_column: 'ts',
      status: 'ready',
    };
    const rowGroups: DbRowGroup[] = [
      {
        id: 'group-shared-temp',
        connector_id: 'db-1',
        table_schema: 'main',
        table_name: 'sensor_values',
        member_point_ids: ['point-A', 'point-B'],
        group_key_columns: ['ts', 'line_id'],
      },
    ];
    const target: DbTarget = {
      tag_id: 'tag-A',
      column_name: 'temperature_c',
      enabled: true,
      row_group_id: 'group-shared-temp',
    };

    expect(toStudioV2DatabaseConfigRequest(connector, rowGroups).row_groups).toEqual(rowGroups);
    expect(toStudioV2DatabaseTargetRequest(target)).toEqual({
      column_name: 'temperature_c',
      enabled: true,
      row_group_id: 'group-shared-temp',
    });
  });

  it('hydrates target row-group references without inventing them for legacy rows', () => {
    expect(hydrateStudioV2DatabaseTarget({
      id: 'row-A',
      workspace_id: 'workspace-1',
      point_id: 'point-A',
      tag_id: 'tag-A',
      column_name: 'temperature_c',
      enabled: true,
      row_group_id: 'group-shared-temp',
      save_state: 'saved',
      created_at: '2026-06-09T00:00:00Z',
      updated_at: '2026-06-09T00:00:00Z',
    }).row_group_id).toBe('group-shared-temp');

    expect(hydrateStudioV2DatabaseTarget({
      id: 'row-legacy',
      workspace_id: 'workspace-1',
      point_id: 'point-legacy',
      tag_id: 'tag-legacy',
      column_name: 'line_a',
      enabled: true,
      save_state: 'saved',
      created_at: '2026-06-09T00:00:00Z',
      updated_at: '2026-06-09T00:00:00Z',
    }).row_group_id).toBeUndefined();
  });

  it('uses the persisted point id for database target saves', () => {
    expect(resolveStudioV2DatabaseTargetPointID('rule-local-p-1', {
      point_id: 'rule-local-p-1',
      persisted_point_id: 'persisted-point-1',
      tag_id: 'tag-1',
      tag_key: 'line.t1',
      display_name: 'T1',
      unit: '',
      target_type: 'float64',
      scale: 1,
      offset: 0,
      enabled: true,
      persisted: true,
    })).toBe('persisted-point-1');

    expect(resolveStudioV2DatabaseTargetPointID('rule-local-p-2', undefined)).toBeUndefined();
  });
});
