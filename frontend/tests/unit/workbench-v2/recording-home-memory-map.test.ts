import { describe, expect, it } from 'vitest';
import { buildMemoryMap, describeMemoryRule, initialMemoryWindow, locateMemoryOffset, memoryRuleInputs, memoryWindow, type MemoryRuleInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';
import type { PersistedMCRuleState } from '@/features/datalink/workbench-v2/state/mcPersistedPoints';

const rule = (id: string, start_address: string, count: number, data_type: MemoryRuleInput['data_type'], extra: Partial<MemoryRuleInput> = {}): MemoryRuleInput =>
  ({ id, label: `${id}_`, start_address, count, data_type, skipped_addresses: [], ...extra });
const mcProof = (start: number, count: number, width = 1, area = 'D') => Array.from({ length: count }, (_, index) => `${area}${start + index * width}`);
const area = (map: ReturnType<typeof buildMemoryMap>, key: string) => map.areas.find(item => item.key === key)!;
const cells = (map: ReturnType<typeof buildMemoryMap>, key: string, from: number, length: number, blocked = false) => memoryWindow(area(map, key), from, length, blocked);

describe('memory map geometry follows the protocol parser and source-rule widths', () => {
  it('shows MC D int16 values as independent one-word cells and float32 as one value across two words', () => {
    const map = buildMemoryMap('mc_3e', [rule('I', 'D0', 2, 'int16', { proven: mcProof(0, 2) }), rule('F', 'D10', 1, 'float32', { proven: ['D10'] })]);
    const shown = cells(map, 'D', 0, 20);
    expect(shown[0]).toMatchObject({ address: 'D0', state: 'used', owners: [{ ruleId: 'I', valueAddress: 'D0', part: 1, width: 1 }] });
    expect(shown[1]).toMatchObject({ address: 'D1', state: 'used', owners: [{ ruleId: 'I', valueAddress: 'D1', part: 1, width: 1 }] });
    expect(shown[10]).toMatchObject({ address: 'D10', owners: [{ ruleId: 'F', valueAddress: 'D10', part: 1, width: 2 }] });
    expect(shown[11]).toMatchObject({ address: 'D11', owners: [{ ruleId: 'F', valueAddress: 'D10', part: 2, width: 2 }] });
    expect(shown[2]).toMatchObject({ address: 'D2', state: 'empty', owners: [] });
    expect(shown[12]).toMatchObject({ address: 'D12', state: 'empty' });
  });

  it('keeps a value that crosses a row recognisable with each word on its own address', () => {
    const map = buildMemoryMap('mc_3e', [rule('F', 'D9', 1, 'float32', { proven: ['D9'] })]);
    const [nine, ten] = cells(map, 'D', 9, 2);
    expect(area(map, 'D').columns).toBe(10);
    expect(nine).toMatchObject({ address: 'D9', owners: [{ ruleId: 'F', valueAddress: 'D9', part: 1, width: 2 }] });
    expect(ten).toMatchObject({ address: 'D10', owners: [{ ruleId: 'F', valueAddress: 'D9', part: 2, width: 2 }] });
  });

  it('isolates MC D/M and Modbus 3x/4x areas and uses the parser unit and radix', () => {
    const mc = buildMemoryMap('mc_3e', [rule('D', 'D0', 1, 'int16', { proven: ['D0'] }), rule('M', 'M0', 1, 'bool', { proven: ['M0'] })]);
    expect(area(mc, 'D')).toMatchObject({ unit: 'word', columns: 10 });
    expect(area(mc, 'M')).toMatchObject({ unit: 'bit', columns: 10 });
    expect(area(mc, 'X')).toMatchObject({ unit: 'bit', columns: 8 });
    expect(area(mc, 'W')).toMatchObject({ unit: 'word', columns: 16 });
    expect(cells(mc, 'M', 0, 2).map(cell => cell.state)).toEqual(['used', 'empty']);
    expect(cells(mc, 'D', 0, 2).map(cell => cell.owners.map(owner => owner.ruleId))).toEqual([['D'], []]);
    expect(cells(mc, 'X', 6, 3).map(cell => cell.address)).toEqual(['X6', 'X7', 'X10']);
    expect(cells(mc, 'W', 9, 2).map(cell => cell.address)).toEqual(['W9', 'WA']);
    const modbus = buildMemoryMap('modbus_tcp', [rule('H', '40001', 2, 'float32'), rule('I', '30001', 1, 'int16')]);
    expect(modbus.areas.map(item => item.key)).toEqual(['HR', 'IR', 'CS', 'IS']);
    expect(area(modbus, 'HR')).toMatchObject({ unit: 'word', base: 1 });
    expect(area(modbus, 'CS')).toMatchObject({ unit: 'bit' });
    expect(cells(modbus, 'HR', 1, 5).map(cell => `${cell.address}:${cell.owners.map(owner => `${owner.valueAddress}/${owner.part}`).join()}`))
      .toEqual(['40001:40001/1', '40002:40001/2', '40003:40003/1', '40004:40003/2', '40005:']);
    expect(cells(modbus, 'IR', 1, 2).map(cell => cell.state)).toEqual(['used', 'empty']);
  });

  it('marks a cell claimed by two ranges as a conflict listing both owners', () => {
    const map = buildMemoryMap('modbus_tcp', [rule('A', '40001', 2, 'int16'), rule('B', '40002', 1, 'float32')]);
    const [, second, third] = cells(map, 'HR', 1, 3);
    expect(second.state).toBe('conflict');
    expect(second.owners.map(owner => owner.ruleId)).toEqual(['A', 'B']);
    expect(third).toMatchObject({ state: 'used', owners: [{ ruleId: 'B', part: 2 }] });
  });

  it('keeps skipped saved addresses owned instead of offering them as free cells', () => {
    const modbus = buildMemoryMap('modbus_tcp', [rule('A', '40001', 3, 'int16', { skipped_addresses: ['40002'] })]);
    expect(cells(modbus, 'HR', 2, 1)[0]).toMatchObject({ state: 'used', owners: [{ ruleId: 'A', skipped: true }] });
    const mc = buildMemoryMap('mc_3e', [rule('A', 'D0', 3, 'int16', { skipped_addresses: ['D1'], proven: ['D0', 'D2'] })]);
    expect(cells(mc, 'D', 0, 4).map(cell => `${cell.state}:${cell.owners[0]?.skipped ?? '-'}`)).toEqual(['used:false', 'used:true', 'used:false', 'empty:-']);
  });
});

describe('memory map fails closed when occupancy cannot be proved', () => {
  it('never derives MC saved points from start plus stride without current persisted proof', () => {
    const map = buildMemoryMap('mc_3e', [rule('A', 'D0', 32, 'int16')]);
    expect(area(map, 'D').issues).toEqual([{ ruleId: 'A', label: 'A_', issue: 'mc-proof' }]);
    expect(cells(map, 'D', 0, 40).every(cell => cell.state === 'unconfirmed' && cell.owners.length === 0)).toBe(true);
    // A valid revision with one missing proof blocks only that range's own area.
    expect(map.global).toEqual([]);
    expect(cells(map, 'M', 0, 1)[0].state).toBe('empty');
    const noScope = buildMemoryMap('mc_3e', [rule('A', 'D0', 32, 'int16', { scopeMissing: true })]);
    expect(noScope.global).toEqual([{ ruleId: 'A', label: 'A_', issue: 'mc-scope' }]);
    expect(cells(noScope, 'M', 0, 3).every(cell => cell.state === 'unconfirmed')).toBe(true);
    const emptyNoScope = buildMemoryMap('mc_3e', [], { listScopeMissing: true });
    expect(emptyNoScope.global).toEqual([{ ruleId: '', label: '', issue: 'list-scope' }]);
    expect(['D', 'M', 'X'].every(key => cells(emptyNoScope, key, 0, 2).every(cell => cell.state === 'unconfirmed'))).toBe(true);
    expect(cells(buildMemoryMap('mc_3e', []), 'D', 0, 1)[0].state).toBe('empty');
    const proven = buildMemoryMap('mc_3e', [rule('A', 'D0', 1, 'int16', { proven: ['D0'] })]);
    expect(proven.global).toEqual([]);
    expect(cells(proven, 'M', 0, 1)[0].state).toBe('empty');
  });

  it('marks a string area unconfirmed because frontend and backend point spans disagree', () => {
    const map = buildMemoryMap('modbus_tcp', [rule('S', '40001', 2, 'string')]);
    expect(area(map, 'HR').issues).toEqual([{ ruleId: 'S', label: 'S_', issue: 'string-width' }]);
    const shown = cells(map, 'HR', 1, 30);
    expect(shown[0]).toMatchObject({ state: 'used', owners: [{ ruleId: 'S', partial: true }] });
    expect(shown.slice(1).every(cell => cell.state === 'unconfirmed')).toBe(true);
    expect(describeMemoryRule('modbus_tcp', rule('S', '40001', 2, 'string'))).toMatchObject({ start: '40001', end: undefined, cells: undefined, issue: 'string-width' });
  });

  it('does not claim multi-word values on a bit area', () => {
    const map = buildMemoryMap('modbus_tcp', [rule('C', '00001', 2, 'float32')]);
    expect(area(map, 'CS').issues).toEqual([{ ruleId: 'C', label: 'C_', issue: 'bit-width' }]);
    expect(cells(map, 'CS', 2, 3).every(cell => cell.state === 'unconfirmed')).toBe(true);
  });

  it('blocks every area when a saved range cannot be placed and when the caller reports blocked evidence', () => {
    const map = buildMemoryMap('modbus_tcp', [rule('BAD', 'D0', 1, 'int16'), rule('A', '40001', 1, 'int16')]);
    expect(map.unplaced).toEqual([{ ruleId: 'BAD', label: 'BAD_', issue: 'invalid-rule' }]);
    expect(cells(map, 'IR', 1, 2).every(cell => cell.state === 'unconfirmed')).toBe(true);
    const ok = buildMemoryMap('modbus_tcp', [rule('A', '40001', 1, 'int16')]);
    expect(cells(ok, 'HR', 1, 3, true).map(cell => cell.state)).toEqual(['used', 'unconfirmed', 'unconfirmed']);
  });

  it('supports every FATEK area the parser expresses with its own bit or word unit', () => {
    const map = buildMemoryMap('fatek_fbs', [rule('A', 'D0', 2, 'int16'), rule('F', 'R10', 1, 'float32'), rule('B', 'M0', 1, 'bool')]);
    expect(map.supported).toBe(true);
    expect(map.areas.map(item => `${item.key}:${item.unit}:${item.columns}`)).toEqual(['D:word:10', 'R:word:10', 'M:bit:10', 'S:bit:10', 'X:bit:10', 'Y:bit:10', 'T:bit:10', 'C:bit:10']);
    expect(cells(map, 'D', 0, 3).map(cell => cell.state)).toEqual(['used', 'used', 'empty']);
    expect(cells(map, 'R', 10, 2).map(cell => `${cell.address}:${cell.owners[0]?.valueAddress}/${cell.owners[0]?.part}`)).toEqual(['R10:R10/1', 'R11:R10/2']);
    expect(cells(map, 'X', 9, 2).map(cell => cell.address)).toEqual(['X9', 'X10']);
    expect(cells(map, 'M', 0, 2).map(cell => cell.state)).toEqual(['used', 'empty']);
  });

  it('marks an area unconfirmed instead of truncating a range past the parser limit or safe integers', () => {
    const tail = buildMemoryMap('modbus_tcp', [rule('T', '199990', 20, 'int16')]);
    expect(area(tail, 'IS').issues).toEqual([{ ruleId: 'T', label: 'T_', issue: 'out-of-range' }]);
    expect(cells(tail, 'IS', 99980, 10).every(cell => cell.state === 'unconfirmed')).toBe(true);
    const fits = buildMemoryMap('modbus_tcp', [rule('T', '199990', 10, 'int16')]);
    expect(area(fits, 'IS').issues).toEqual([]);
    expect(cells(fits, 'IS', 99999, 1)[0]).toMatchObject({ address: '199999', state: 'used' });
    const huge = buildMemoryMap('modbus_tcp', [rule('H', '40001', Number.MAX_SAFE_INTEGER, 'float64')]);
    expect(area(huge, 'HR').issues).toEqual([{ ruleId: 'H', label: 'H_', issue: 'out-of-range' }]);
    const mcEdge = buildMemoryMap('mc_3e', [rule('E', 'DFFFFFF', 1, 'int16', { proven: ['D16777215'] }), rule('F', 'D16777215', 1, 'float32', { proven: ['D16777215'] })]);
    expect(area(mcEdge, 'D').max).toBe(0xffffff);
    expect(area(mcEdge, 'D').issues).toEqual([{ ruleId: 'F', label: 'F_', issue: 'out-of-range' }]);
    expect(cells(mcEdge, 'D', 0, 3).every(cell => cell.state === 'unconfirmed')).toBe(true);
    expect(describeMemoryRule('modbus_tcp', rule('T', '199990', 20, 'int16'))).toMatchObject({ end: undefined, cells: undefined, issue: 'out-of-range' });
  });

  it('bounds FATEK areas by the backend AddressLimit: D/R end at 65535, X/Y/M/S/T/C at 9999', () => {
    const map = buildMemoryMap('fatek_fbs', []);
    expect(map.areas.map(item => `${item.key}:${item.max}`)).toEqual(['D:65535', 'R:65535', 'M:9999', 'S:9999', 'X:9999', 'Y:9999', 'T:9999', 'C:9999']);
    const tail = memoryWindow(area(map, 'D'), 65530, 20, false);
    expect(tail.map(cell => cell.address)).toEqual(['D65530', 'D65531', 'D65532', 'D65533', 'D65534', 'D65535']);
    expect(memoryWindow(area(map, 'X'), 9995, 20, false).at(-1)).toMatchObject({ address: 'X9999', state: 'empty' });
    expect(locateMemoryOffset(area(map, 'D'), 'D65535')).toBe(65535);
    expect(locateMemoryOffset(area(map, 'D'), 'D65536')).toBeUndefined();
    expect(locateMemoryOffset(area(map, 'X'), 'X9999')).toBe(9999);
    expect(locateMemoryOffset(area(map, 'X'), 'X10000')).toBeUndefined();
  });

  it('marks FATEK ranges that cross or start past the protocol limit unconfirmed instead of free', () => {
    const fits = buildMemoryMap('fatek_fbs', [rule('F', 'R65534', 1, 'float32'), rule('B', 'M9999', 1, 'bool')]);
    expect(area(fits, 'R').issues).toEqual([]);
    expect(cells(fits, 'R', 65533, 3).map(cell => `${cell.address}:${cell.state}:${cell.owners[0]?.part ?? '-'}`)).toEqual(['R65533:empty:-', 'R65534:used:1', 'R65535:used:2']);
    expect(describeMemoryRule('fatek_fbs', rule('F', 'R65534', 1, 'float32'))).toMatchObject({ start: 'R65534', end: 'R65535', cells: 2 });
    expect(cells(fits, 'M', 9999, 1)[0]).toMatchObject({ address: 'M9999', state: 'used' });
    const crossing = buildMemoryMap('fatek_fbs', [rule('F', 'R65535', 1, 'float32'), rule('B', 'M9999', 2, 'bool'), rule('P', 'D70000', 1, 'int16')]);
    for (const key of ['R', 'M', 'D']) {
      expect(area(crossing, key).issues.map(item => item.issue)).toEqual(['out-of-range']);
      expect(cells(crossing, key, 0, 3).every(cell => cell.state === 'unconfirmed')).toBe(true);
    }
    expect(cells(crossing, 'R', 65535, 1)[0].state).toBe('unconfirmed');
    expect(describeMemoryRule('fatek_fbs', rule('P', 'D70000', 1, 'int16'))).toMatchObject({ end: undefined, cells: undefined, issue: 'out-of-range' });
  });

  it('reports a protocol the address parser does not support as unsupported', () => {
    expect(buildMemoryMap('opcua' as never, [rule('A', 'D0', 1, 'int16')])).toMatchObject({ supported: false, areas: [] });
  });

  it('accepts MC proof only for the same workspace, current rule revision and compatible status', () => {
    const base = { id: 'R', revision_id: 'R2', start_address: 'D0', count: 2, data_type: 'int16' as const, naming_prefix: 'P_', skipped_addresses: [] };
    const good: PersistedMCRuleState = { status: 'compatible', scope: { workspace_id: 'W', workspace_revision: 'W2', revision_id: 'R2' }, points: [{ point_id: 'a', address: 'D0' }, { point_id: 'b', address: 'D1' }] };
    const listed = { id: 'W', revision: 'W2' };
    expect(memoryRuleInputs([base], 'mc_3e', listed, () => good)[0].proven).toEqual(['D0', 'D1']);
    for (const bad of [{ ...good, status: 'conflict' as const }, { ...good, scope: { ...good.scope!, workspace_id: 'X' } }, { ...good, scope: { ...good.scope!, revision_id: 'OLD' } }, { ...good, points: [] }, undefined])
      expect(memoryRuleInputs([base], 'mc_3e', listed, () => bad)[0].proven).toBeUndefined();
    expect(memoryRuleInputs([base], 'modbus_tcp', listed, () => undefined)[0]).toMatchObject({ label: 'P_', proven: undefined });
  });

  it('rejects W1 evidence for a W2 range list and any proof when the list revision is unknown', () => {
    const base = { id: 'R', revision_id: 'R2', start_address: 'D0', count: 2, data_type: 'int16' as const, naming_prefix: 'P_', skipped_addresses: [] };
    const w1: PersistedMCRuleState = { status: 'compatible', scope: { workspace_id: 'W', workspace_revision: 'W1', revision_id: 'R2' }, points: [{ point_id: 'a', address: 'D0' }, { point_id: 'b', address: 'D1' }] };
    const inputs = memoryRuleInputs([base], 'mc_3e', { id: 'W', revision: 'W2' }, () => w1);
    expect(inputs[0].proven).toBeUndefined();
    expect(memoryRuleInputs([base], 'mc_3e', { id: 'W' }, () => w1)[0]).toMatchObject({ proven: undefined, scopeMissing: true });
    expect(inputs[0].scopeMissing).toBeUndefined();
    const map = buildMemoryMap('mc_3e', inputs);
    expect(area(map, 'D').issues).toEqual([{ ruleId: 'R', label: 'P_', issue: 'mc-proof' }]);
    expect(cells(map, 'D', 0, 4).every(cell => cell.state === 'unconfirmed')).toBe(true);
  });
});

describe('memory windows stay bounded and locatable', () => {
  it('starts MC D at D0 and shows empty cells after the last used register', () => {
    const map = buildMemoryMap('mc_3e', [rule('A', 'D0', 32, 'int16', { proven: mcProof(0, 32) })]);
    const window = initialMemoryWindow(area(map, 'D'));
    expect(window.start).toBe(0);
    const shown = memoryWindow(area(map, 'D'), window.start, window.length, false);
    expect(shown[0].address).toBe('D0');
    expect(shown.find(cell => cell.address === 'D33')).toMatchObject({ state: 'empty' });
  });

  it('never renders millions of cells for a huge range and can locate a far address', () => {
    const map = buildMemoryMap('modbus_tcp', [rule('BIG', '40001', 60000, 'int16')]);
    const hr = area(map, 'HR');
    const window = initialMemoryWindow(hr);
    expect(window.length).toBeLessThanOrEqual(256);
    expect(memoryWindow(hr, window.start, window.length, false).length).toBeLessThanOrEqual(256);
    const located = locateMemoryOffset(hr, '460001');
    expect(located).toBe(60001);
    expect(locateMemoryOffset(hr, '30001')).toBeUndefined();
    expect(locateMemoryOffset(area(buildMemoryMap('mc_3e', []), 'X'), 'X10')).toBe(8);
    const tail = memoryWindow(hr, hr.max - 5, 20, false);
    expect(tail.length).toBe(6);
    expect(tail.at(-1)!.offset).toBe(hr.max);
  });

  it('describes saved ranges with real start, end, item count and register footprint', () => {
    expect(describeMemoryRule('mc_3e', rule('F', 'D0', 4, 'float32', { proven: mcProof(0, 4, 2) }))).toEqual({ area: 'D', unit: 'word', start: 'D0', end: 'D7', cells: 8, issue: undefined });
    expect(describeMemoryRule('mc_3e', rule('F', 'D0', 4, 'float32'))).toMatchObject({ start: 'D0', end: undefined, cells: undefined, issue: 'mc-proof' });
    expect(describeMemoryRule('modbus_tcp', rule('B', '00001', 3, 'bool'))).toEqual({ area: 'CS', unit: 'bit', start: '00001', end: '00003', cells: 3, issue: undefined });
  });
});
