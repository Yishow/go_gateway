import { useCallback, useLayoutEffect, useMemo, useSyncExternalStore } from 'react';

export type SchemaOperationPhase = 'idle' | 'previewing' | 'creating' | 'unknown';
export type SchemaOperationKind = Exclude<SchemaOperationPhase, 'idle' | 'unknown'>;

export interface SchemaOperationSnapshot {
  phase: SchemaOperationPhase;
  scopeSignature: string;
  generation: number;
}

export interface SchemaOperationToken {
  workspaceKey: string;
  scopeSignature: string;
  generation: number;
  phase: SchemaOperationKind;
}

interface SchemaOperationStore {
  snapshot: SchemaOperationSnapshot;
  currentScopeSignature: string;
  generation: number;
  listeners: Set<() => void>;
}

const stores = new Map<string, SchemaOperationStore>();
const serverSnapshot: SchemaOperationSnapshot = {
  phase: 'idle',
  scopeSignature: '',
  generation: 0,
};

function storeFor(workspaceKey: string): SchemaOperationStore {
  const existing = stores.get(workspaceKey);
  if (existing) {
    return existing;
  }

  const created: SchemaOperationStore = {
    snapshot: serverSnapshot,
    currentScopeSignature: '',
    generation: 0,
    listeners: new Set(),
  };
  stores.set(workspaceKey, created);
  return created;
}

function notify(store: SchemaOperationStore): void {
  store.listeners.forEach((listener) => listener());
}

function setSnapshot(store: SchemaOperationStore, snapshot: SchemaOperationSnapshot): void {
  store.snapshot = snapshot;
  notify(store);
}

function isBusy(phase: SchemaOperationPhase): boolean {
  return phase === 'creating' || phase === 'unknown';
}

/** 讓同一工作區的短暫 remount 仍看得到進行中的 schema 操作。 */
export function useSchemaOperation(
  workspaceKey: string,
  scopeSignature: string,
): SchemaOperationSnapshot {
  const store = useMemo(() => storeFor(workspaceKey), [workspaceKey]);
  const subscribe = useCallback((listener: () => void) => {
    store.listeners.add(listener);
    return () => store.listeners.delete(listener);
  }, [store]);
  const getSnapshot = useCallback(() => store.snapshot, [store]);

  const snapshot = useSyncExternalStore(subscribe, getSnapshot, () => serverSnapshot);

  useLayoutEffect(() => {
    store.currentScopeSignature = scopeSignature;
  }, [scopeSignature, store]);

  return snapshot;
}

/** 開始一次預覽或建表；同一工作區的進行中建表／未知結果不可重送。 */
export function beginSchemaOperation(
  workspaceKey: string,
  scopeSignature: string,
  phase: SchemaOperationKind,
): SchemaOperationToken | null {
  const store = storeFor(workspaceKey);
  const current = store.snapshot;

  if (isBusy(current.phase)) {
    return null;
  }
  if (phase === 'creating' && current.phase === 'previewing') {
    return null;
  }
  if (phase === 'previewing' && current.phase === 'previewing' && current.scopeSignature === scopeSignature) {
    return null;
  }

  const generation = store.generation + 1;
  store.generation = generation;
  const token: SchemaOperationToken = {
    workspaceKey,
    scopeSignature,
    generation,
    phase,
  };
  setSnapshot(store, {
    phase,
    scopeSignature,
    generation,
  });
  return token;
}

/** 只接受仍屬於目前 generation 的回覆；舊回覆不會改寫新操作。 */
export function settleSchemaOperation(
  token: SchemaOperationToken,
  phase: 'idle' | 'unknown',
): boolean {
  const store = stores.get(token.workspaceKey);
  if (!store || store.snapshot.generation !== token.generation || store.snapshot.phase !== token.phase) {
    return false;
  }

  setSnapshot(store, {
    phase,
    scopeSignature: token.scopeSignature,
    generation: token.generation,
  });
  return true;
}

/** 判斷回覆是否仍對目前畫面的 scope 有效；不代表後端已取消舊操作。 */
export function isSchemaOperationScopeCurrent(token: SchemaOperationToken): boolean {
  const store = stores.get(token.workspaceKey);
  return store?.currentScopeSignature === token.scopeSignature;
}

/** 僅用於將可判定的 HTTP 錯誤與網路／未知結果分開。 */
export function isExplicitSchemaHTTPFailure(error: unknown): boolean {
  if (typeof error !== 'object' || error === null) {
    return false;
  }
  const response = (error as { response?: unknown }).response;
  const responseStatus = typeof response === 'object' && response !== null
    ? (response as { status?: unknown }).status
    : undefined;
  const directStatus = (error as { status?: unknown }).status;
  const status = typeof responseStatus === 'number' ? responseStatus : directStatus;
  return typeof status === 'number' && status >= 400 && status < 600;
}
