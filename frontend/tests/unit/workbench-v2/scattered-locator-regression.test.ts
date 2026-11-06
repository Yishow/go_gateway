import { describe, expect, it } from 'vitest';
import { studioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { locatorRead } from './studio-locator-fixtures';

function fatekRead(address = 'R10') {
  const read = locatorRead(3);
  read.entry.devices![0].protocol = 'fatek_fbs';
  Object.assign(read.entry.source_rules![0], { start_address: 'R0', layout_mode: 'mixed', count: 2, mixed_items: [
    { item_id: 'a', name: 'A', address: 'R0', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' },
    { item_id: 'b', name: 'B', address: 'R10', register_offset: 0, data_type: 'uint16', semantic_kind: 'gauge' },
  ] });
  Object.assign(read.entry.mappings![0], { address, point_id: 'point-b' });
  return read;
}

describe('FATEK mixed mapping locator', () => {
  it('hydrates the selected native point with its own address and type', () => {
    const state = studioV2LocatorHydration(fatekRead());
    expect(state?.points).toEqual([expect.objectContaining({ id: 'point-b', address: 'R10', data_type: 'uint16' })]);
    expect(state?.mappings['point-b']).toMatchObject({ mapping_id: 'M', persisted_point_id: 'point-b' });
  });

  it('rejects a selected mapping whose address matches no mixed item', () => {
    expect(studioV2LocatorHydration(fatekRead('R20'))).toBeNull();
  });

  it('does not hydrate a partial editor when a mixed mapping has no matching item', () => {
    const read = fatekRead('R20');
    read.entry.editor_hydrated = true;
    if (read.selection.state !== 'selected') throw new Error('expected selected locator');
    read.selection = { ...read.selection, target: { ...read.selection.target, mappingId: undefined } };
    expect(studioV2LocatorHydration(read)).toBeNull();
  });
});
