import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { inspectSchemaOperation, readSchemaOperation, retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';

beforeEach(() => {
  vi.unstubAllGlobals();
  sessionStorage.clear();
});
afterEach(() => { vi.unstubAllGlobals(); });

describe.each(['session storage', 'memory fallback'])('schema recovery with %s', (storage) => {
  function useStorage() {
    if (storage === 'memory fallback') vi.stubGlobal('sessionStorage', {
      getItem: () => { throw new Error('storage unavailable'); },
      setItem: () => { throw new Error('storage unavailable'); },
      removeItem: () => { throw new Error('storage unavailable'); },
    });
  }

  it('retains the newer operation when a terminal result belongs to an older operation', () => {
    useStorage();
    const groupId = `newer-${storage}`;
    retainSchemaOperation('workspace-1', groupId, 'operation-B');
    retainSchemaOperation('workspace-1', groupId, undefined, 'operation-A');
    expect(readSchemaOperation('workspace-1', groupId)).toBe('operation-B');
    if (storage === 'session storage') {
      expect(sessionStorage.length).toBe(1);
      expect(sessionStorage.getItem(sessionStorage.key(0)!)).toBe('operation-B');
    }
    retainSchemaOperation('workspace-1', groupId, undefined, 'operation-B');
    expect(readSchemaOperation('workspace-1', groupId)).toBeUndefined();
  });

  it('releases only the matching workspace and group operation', () => {
    useStorage();
    const groupId = `scoped-${storage}`;
    retainSchemaOperation('workspace-1', groupId, 'operation-A');
    retainSchemaOperation('workspace-2', groupId, 'operation-B');
    retainSchemaOperation('workspace-1', groupId, undefined, 'operation-A');
    expect(readSchemaOperation('workspace-1', groupId)).toBeUndefined();
    expect(readSchemaOperation('workspace-2', groupId)).toBe('operation-B');
    retainSchemaOperation('workspace-2', groupId, undefined, 'operation-B');
  });
});

it('keeps the newer memory fence when storage reads work but a quota error rejects writes', () => {
  const storage = sessionStorage;
  const groupId = 'quota-read-fallback';
  retainSchemaOperation('workspace-1', groupId, 'operation-A');
  vi.stubGlobal('sessionStorage', {
    getItem: (key: string) => storage.getItem(key),
    setItem: () => { throw new Error('quota exceeded'); },
    removeItem: (key: string) => storage.removeItem(key),
  });
  retainSchemaOperation('workspace-1', groupId, 'operation-B');
  retainSchemaOperation('workspace-1', groupId, undefined, 'operation-A');
  expect(readSchemaOperation('workspace-1', groupId)).toBe('operation-B');
  retainSchemaOperation('workspace-1', groupId, undefined, 'operation-B');
  expect(readSchemaOperation('workspace-1', groupId)).toBeUndefined();
});


describe('read-only owner-local schema inspection', () => {
  it('distinguishes readable empty, missing scope, and storage failure without changing the compatible reader', () => {
    expect(inspectSchemaOperation('workspace-1', 'inspection-empty')).toEqual({ state: 'empty' });
    expect(inspectSchemaOperation('', 'inspection-empty')).toEqual({ state: 'unavailable', reason: 'missing_scope' });
    expect(inspectSchemaOperation('workspace-1')).toEqual({ state: 'unavailable', reason: 'missing_scope' });
    vi.stubGlobal('sessionStorage', { getItem: () => { throw new Error('private storage error'); } });
    expect(inspectSchemaOperation('workspace-1', 'inspection-empty')).toEqual({ state: 'unavailable', reason: 'unreadable' });
    expect(readSchemaOperation('workspace-1', 'inspection-empty')).toBeUndefined();
  });
  it.each(['opaque-O1', '   ', '{not-json}', 'https://untrusted.example/op/1'])('retains every nonempty stored opaque identity unchanged: %s', (operationId) => {
    const groupId = 'inspection-opaque';
    sessionStorage.setItem(`wbv2.schema-operation.${JSON.stringify(['workspace-1', groupId])}`, operationId);
    expect(inspectSchemaOperation('workspace-1', groupId)).toEqual({ state: 'retained', operationId });
    expect(sessionStorage.getItem(sessionStorage.key(0)!)).toBe(operationId);
  });
  it.each(['empty', 'error'])('keeps the live memory fence ahead of %s storage', (storage) => {
    const groupId = `inspection-memory-${storage}`;
    retainSchemaOperation('workspace-1', groupId, 'O1');
    vi.stubGlobal('sessionStorage', { getItem: () => { if (storage === 'error') throw new Error('unavailable'); return null; },
      setItem: vi.fn(), removeItem: vi.fn() });
    expect(inspectSchemaOperation('workspace-1', groupId)).toEqual({ state: 'retained', operationId: 'O1' });
    retainSchemaOperation('workspace-1', groupId, undefined, 'O1');
  });
});
