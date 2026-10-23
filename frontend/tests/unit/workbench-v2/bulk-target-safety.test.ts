import { expect, it } from 'vitest';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { mappingReducer } from '@/features/datalink/workbench-v2/state/mappingReducer';
it('copies only the checked fields to checked enabled identities across mixed units and types', () => {
  const state = savedGroupState();
  const [temperature, pressure, energy] = state.points.map(p => p.id);
  state.mappings[temperature] = { ...state.mappings[temperature], unit: 'C', scale: 0.1 };
  state.mappings[pressure] = { ...state.mappings[pressure], unit: 'bar', scale: 2 };
  state.mappings[energy] = { ...state.mappings[energy], unit: 'kWh', target_type: 'uint64', scale: 1 };
  const next = mappingReducer(state, { type: 'bulkApplyTransform', fromPointId: temperature, toPointIds: [pressure], fields: ['scale'] });
  expect(next.mappings[pressure]).toMatchObject({ scale: 0.1, unit: 'bar' });
  expect(next.mappings[energy]).toBe(state.mappings[energy]);
  expect(next.mappings[temperature]).toBe(state.mappings[temperature]);
  const none = mappingReducer(state, { type: 'bulkApplyTransform', fromPointId: temperature, toPointIds: [], fields: ['unit'] });
  expect(none.mappings).toEqual(state.mappings);
});
it('ignores disabled, absent and source rows even if present in the target list', () => {
  const state = savedGroupState(), ids = state.points.map(p => p.id);
  state.mappings[ids[1]].enabled = false;
  const next = mappingReducer(state, { type: 'bulkApplyTransform', fromPointId: ids[0], toPointIds: [ids[0], ids[1], 'foreign'], fields: ['unit'] });
  expect(next.mappings).toEqual(state.mappings);
});
