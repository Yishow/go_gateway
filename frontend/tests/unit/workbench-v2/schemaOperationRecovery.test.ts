import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { readSchemaOperation, retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';

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
