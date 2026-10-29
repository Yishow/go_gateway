import { describe, expect, it } from 'vitest';
import { studioV2LocatorHydration, studioLocatorStateMatches } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { locatorRead } from './studio-locator-fixtures';

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
