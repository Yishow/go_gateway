import { expect, it } from 'vitest';
import { hasOwnedWideDestination, ownedWideCatalogTables } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/recordingCatalogOwnership';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

it('blocks only other known wide owners for the same workspace and destination', () => {
  const connector = { ...savedGroupState().db.connector, database: 'db', schema: 'public' };
  const group = savedGroup();
  const other = savedGroup({ id: 'other', destination: { ...group.destination, database: 'db', table_schema: 'public', table_name: 'gw_owned', storage_strategy: 'managed' } });
  const data = { workspace_id: 'ws-1', workspace_revision: 'W', groups: [other, { ...other, id: 'current', destination: { ...other.destination, table_name: 'gw_current' } },
    { ...other, id: 'long', layout_version: 'managed_long_v1' as const, destination: { ...other.destination, table_name: 'gw_shared' } },
    { ...other, id: 'elsewhere', destination: { ...other.destination, database: 'other_db', table_name: 'other_db_table' } }] };
  const observation = { data, loading: false, failed: false, retry: () => undefined };
  expect(ownedWideCatalogTables(observation, { workspaceId: 'ws-1', workspaceRevision: 'W', connector, groupId: 'current' }, 'public')).toEqual([{ schema: 'public', name: 'gw_owned' }]);
  expect(hasOwnedWideDestination(observation, { workspaceId: 'ws-1', workspaceRevision: 'W', connector, groupId: 'current' }, 'public', 'gw_owned')).toBe(true);
  expect(hasOwnedWideDestination(observation, { workspaceId: 'ws-1', workspaceRevision: 'W', connector, groupId: 'current' }, ' public ', ' gw_owned ')).toBe(true);
  expect(hasOwnedWideDestination(observation, { workspaceId: 'ws-1', workspaceRevision: 'W', connector, groupId: 'current' }, 'public', 'gw_current')).toBe(false);
  expect(ownedWideCatalogTables(observation, { workspaceId: 'ws-1', workspaceRevision: 'stale', connector, groupId: 'current' }, 'public')).toEqual([]);
});
