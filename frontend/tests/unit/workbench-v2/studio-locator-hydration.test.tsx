import { describe, expect, it } from 'vitest';
import { studioV2LocatorHydration, studioLocatorStateMatches } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { locatorRead } from './studio-locator-fixtures';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';

describe('detached locator hydration', () => {
  it('hydrates only requested device without default samples or observed connector', () => {
    const read = locatorRead(); const state = studioV2LocatorHydration(read)!;
    expect(state.devices.map(device => device.id)).toEqual(['D']);
    expect(state.rules).toEqual([]); expect(state.settings.connectors).toEqual([]);
    expect(state.db.connector).toMatchObject({ name: '', host: '', status: 'unknown', persisted: false });
    expect(read.entry.source_rules).toBeUndefined(); expect(read.entry.hydration_scope.complete_workspace).toBe(false);
  });
  it('selects a native mapping by exact mapping and persisted point IDs, not the synthetic rule index', () => {
    const read = locatorRead(3); const state = studioV2LocatorHydration(read)!;
    expect(state.points.map(point => point.id)).toEqual(['P-native']);
    expect(state.mappings['P-native']).toMatchObject({ mapping_id: 'M', persisted_point_id: 'P-native', scale: 0.1 });
    expect(studioLocatorStateMatches(read, state)).toBe(true);
    state.mappings['P-native'].persisted_point_id = 'same-address-other-point';
    expect(studioLocatorStateMatches(read, state)).toBe(false);
  });
  it('does not hydrate unavailable or mismatched rule ancestry', () => {
    const read = locatorRead(3); read.entry.mappings![0].rule_id = 'missing';
    expect(studioV2LocatorHydration(read)).toBeNull();
    read.selection = { state: 'unavailable', reason: 'target_unavailable' };
    expect(studioV2LocatorHydration(read)).toBeNull();
  });
  it('requires the current rule revision for hot selection without rewriting a draft', () => {
    const read = locatorRead(2); const state = studioV2LocatorHydration(read)!;
    state.rules[0].revision_id = 'R1'; state.rules[0].name = 'My draft';
    expect(studioLocatorStateMatches(read, state)).toBe(false); expect(state.rules[0].name).toBe('My draft');
  });
  it('preserves full shared/orphan group evidence separately from unloaded device state', () => {
    const read = locatorRead(4); const before = JSON.stringify(read.entry.write_groups);
    const state = studioV2LocatorHydration(read)!;
    expect(state.devices).toEqual([]); expect(state.rules).toEqual([]);
    expect(JSON.stringify(read.entry.write_groups)).toBe(before);
  });
});

it('keeps saved-target DTO hydration detached from the workspace default and connector pool', () => {
  const read = locatorRead(4); const original = JSON.stringify(INITIAL_STATE);
  const group = read.entry.write_groups!.groups[0];
  group.destination = { ...group.destination, connector_id: 'saved-target', connector_revision: 'saved-rev', database: 'saved.db', table_schema: 'main', table_name: 'saved_table' };
  read.entry.editor_hydrated = true;
  read.entry.database_config = { id: 'saved-target', identity_revision: 'saved-rev', setup_revision: 'W2', workspace_id: 'W', kind: 'sqlite', name: 'Saved target label', host: '', port: 0, database: 'saved.db', username: '', schema: 'main', table: 'saved_table', write_mode: 'insert', write_interval_seconds: 15, timestamp_column: 'time', status: 'ready', save_state: 'saved', created_at: '', updated_at: '' };
  read.entry.database_targets = [];
  const state = studioV2LocatorHydration(read)!;
  expect(state.db.connector).toMatchObject({ connector_id: group.destination.connector_id, identity_revision: group.destination.connector_revision,
    database: group.destination.database, schema: group.destination.table_schema, table: group.destination.table_name, name: 'Saved target label', persisted: true, save_state: 'saved' });
  expect(state.db.connector.password).toBeUndefined(); expect(state.settings.connectors).toEqual([]);
  expect(JSON.stringify(INITIAL_STATE)).toBe(original); expect(read.entry.database_config.id).toBe('saved-target');
});
