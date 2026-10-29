import { describe, expect, it } from 'vitest';
import { createRuntimeEntityHistory, readRuntimeEntityHistory } from '../../../src/features/datalink/runtime-dashboard/runtimeEntityHistory';
const scope = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'D' };
const selection = { rule_id: 'R', rule_revision: 'R1', point_id: 'P', mapping_id: 'M' };
describe('identity-only runtime history', () => {
  it('stores only scoped explicit selection and restores a detached identity', () => {
    const state = createRuntimeEntityHistory(scope, selection);
    expect(state).toEqual({ runtimeEntityHistory: { version: 1, ...scope, selection } });
    const restored = readRuntimeEntityHistory(state, scope);
    expect(restored).toEqual(selection); expect(restored).not.toBe(selection);
  });
  it.each(['workspace_id', 'workspace_revision', 'device_id'] as const)('does not restore after %s changes', (key) => {
    expect(readRuntimeEntityHistory(createRuntimeEntityHistory(scope, selection), { ...scope, [key]: 'other' })).toBeNull();
  });
  it.each(['https://private.example', '10.0.0.1', '../private', 'token=secret', ''])('rejects unsafe identity %s before storing it', (id) => {
    expect(createRuntimeEntityHistory({ ...scope, device_id: id }, selection)).toBeNull();
    expect(createRuntimeEntityHistory(scope, { ...selection, mapping_id: id })).toBeNull();
  });
  it('rejects facts, unknown keys, arrays and accessors rather than copying them', () => {
    const valid = { version: 1, ...scope, selection };
    for (const value of [null, [], { runtimeEntityHistory: { ...valid, endpoint: 'private.example' } }, { runtimeEntityHistory: { ...valid, selection: { ...selection, address: '40001' } } }]) expect(readRuntimeEntityHistory(value, scope)).toBeNull();
    let reads = 0;
    expect(readRuntimeEntityHistory({ get runtimeEntityHistory() { reads++; return valid; } }, scope)).toBeNull(); expect(reads).toBe(0);
  });
  it('supports a point-only selection and explicit clear without a first mapping fallback', () => {
    const point = { rule_id: 'R', rule_revision: 'R1', point_id: 'P' };
    expect(readRuntimeEntityHistory(createRuntimeEntityHistory(scope, point), scope)).toEqual(point);
    expect(createRuntimeEntityHistory(scope, null)).toEqual({ runtimeEntityHistory: { version: 1, ...scope, selection: null } });
    expect(readRuntimeEntityHistory(createRuntimeEntityHistory(scope, null), scope)).toBeNull();
  });
});
