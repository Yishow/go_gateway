import { describe, expect, it } from 'vitest';
import { studioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import type { StudioV2WorkspaceMappingRecord } from '@/types/datalink';
import { locatorRead } from './studio-locator-fixtures';

function mcLocatorRead() {
  const read = locatorRead(3);
  read.entry.devices![0].protocol = 'mc_3e';
  read.entry.source_rules![0].start_address = 'X0';
  read.entry.source_rules![0].count = 2;
  read.entry.source_rules![0].data_type = 'int16';
  return read;
}

describe('MC locator hydration', () => {
  it('projects a cold selected MC mapping from its persisted point ID and address', () => {
    const read = mcLocatorRead();
    read.entry.mappings![0] = {
      ...read.entry.mappings![0],
      point_id: 'P-legacy-x8',
      address: 'X8',
      target_type: 'float64',
    };

    const state = studioV2LocatorHydration(read);

    expect(state).not.toBeNull();
    expect(state?.points).toEqual([expect.objectContaining({
      id: 'P-legacy-x8',
      address: 'X8',
      data_type: 'int16',
      width: 1,
      name: 'Mapping M',
    })]);
    expect(state?.mappings['P-legacy-x8']).toMatchObject({
      mapping_id: 'M',
      persisted_point_id: 'P-legacy-x8',
      target_type: 'float64',
    });
  });

  it('projects every persisted MC mapping row for an editor hydration without a Step 2 snapshot', () => {
    const read = mcLocatorRead();
    const first: StudioV2WorkspaceMappingRecord = { ...read.entry.mappings![0], point_id: 'P-x8', address: 'X8', target_type: 'float64' };
    const second: StudioV2WorkspaceMappingRecord = { ...first, id: 'M-2', point_id: 'P-wa', address: 'WA', tag_id: 'T-2', target_type: 'uint64' };
    read.entry = { ...read.entry, editor_hydrated: true, mappings: [first, second] };
    if (read.selection.state !== 'selected') throw new Error('expected selected locator');
    read.selection = { ...read.selection, target: { ...read.selection.target, mappingId: undefined } };

    const state = studioV2LocatorHydration(read);

    expect(state?.points.map(point => ({ id: point.id, address: point.address, dataType: point.data_type }))).toEqual([
      { id: 'P-x8', address: 'X8', dataType: 'int16' },
      { id: 'P-wa', address: 'WA', dataType: 'int16' },
    ]);
    expect(Object.values(state?.mappings ?? {}).map(mapping => mapping.mapping_id)).toEqual(['M', 'M-2']);
  });

  it('retains non-MC locator geometry from normal rule derivation', () => {
    const state = studioV2LocatorHydration(locatorRead(3));

    expect(state?.points).toEqual([expect.objectContaining({
      id: 'P-native',
      address: '40002',
      data_type: 'int16',
      width: 1,
      name: 'Mapping M',
    })]);
  });
});
