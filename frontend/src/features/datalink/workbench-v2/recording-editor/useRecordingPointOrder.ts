import { useEffect, useRef } from 'react';
import { createPointFormatOperation, pointFormatAPI } from '../../../../services/pointFormat';
import { createPointFormatRecovery, pointFormatRecoveryAPI } from '../../../../services/pointFormatRecovery';
import { modbusShareAPI } from '../../../../services/modbusShare';
import { PointFormatError, safePointFormatError } from '../../../../services/pointFormatErrors';
import type { PointFormatIntent, PointFormatMutationResult, PointFormatOperation, PointFormatScope, PointFormatSnapshot } from '../../../../types/pointFormat';
import type { PointFormatRecoveryRequest, PointFormatRecoveryResult } from '../../../../types/pointFormatRecovery';
export interface PointOrderEntry {
  status: 'loading' | 'ready' | 'saving' | 'unknown' | 'saved' | 'inspecting' | 'recovering' | 'recovery-unknown' | 'recovered' | 'error';
  snapshot?: PointFormatSnapshot;
  intent?: PointFormatIntent;
  operation?: PointFormatOperation;
  saved?: PointFormatMutationResult;
  recovery?: PointFormatRecoveryRequest;
  recovered?: PointFormatRecoveryResult;
  error?: PointFormatError;
}
export interface RecordingPointOrderPort {
  available(scope: PointFormatScope): boolean;
  subscribe(listener: () => void): () => void;
  get(scope: PointFormatScope): PointOrderEntry | undefined;
  load(scope: PointFormatScope): Promise<void>;
  choose(scope: PointFormatScope, intent: PointFormatIntent): void;
  save(scope: PointFormatScope): Promise<void>;
  inspect(scope: PointFormatScope): Promise<void>;
  recover(scope: PointFormatScope): Promise<void>;
  restart(scope: PointFormatScope): Promise<void>;
}
const keyOf = (scope: PointFormatScope) => JSON.stringify([scope.workspace_id, scope.device_id, scope.point_id, scope.rule_id ?? null]);
export const canRestartPointOrder = (entry: PointOrderEntry | undefined): boolean => Boolean(entry && (entry.status === 'recovered' ||
  (entry.status === 'error' && (!entry.operation || (!entry.saved && entry.error?.outcome === 'rejected')))));
/** Page lifetime only. Original receipts survive editor closure; never replay or rebase automatically. */
export function useRecordingPointOrder(workspaceId: string | undefined, enabled: boolean): RecordingPointOrderPort {
  const mounted = useRef(true);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const current = useRef({ workspaceId, enabled }); current.current = { workspaceId, enabled };
  const port = useRef<RecordingPointOrderPort | null>(null);
  if (!port.current) {
    const entries = new Map<string, PointOrderEntry>(), listeners = new Set<() => void>();
    const active = (scope: PointFormatScope) => mounted.current && current.current.enabled && scope.workspace_id === current.current.workspaceId;
    const publish = (scope: PointFormatScope, entry: PointOrderEntry) => { entries.set(keyOf(scope), entry); listeners.forEach(listener => listener()); };
    const read = (scope: PointFormatScope) => entries.get(keyOf(scope));
    const inspectOrSave = async (scope: PointFormatScope, inspect: boolean) => {
      const entry = read(scope);
      if (!active(scope) || !entry?.snapshot) return;
      if (inspect ? entry.status !== 'unknown' || !entry.operation : entry.status !== 'ready' || !entry.intent || entry.operation) return;
      let operation = entry.operation;
      try {
        operation ??= createPointFormatOperation(entry.snapshot, entry.intent!, `format-${crypto.randomUUID()}`);
        publish(scope, { ...entry, operation, error: undefined, status: inspect ? 'inspecting' : 'saving' });
        const saved = await (inspect ? pointFormatAPI.inspect(operation) : pointFormatAPI.save(operation));
        publish(scope, { ...entry, operation, saved, status: 'saved' });
      } catch (cause) {
        const error = safePointFormatError(cause, operation?.request_id);
        publish(scope, { ...entry, operation, error, status: inspect || error.outcome === 'unknown' ? 'unknown' : 'error' });
      }
    };
    port.current = {
      available: active,
      subscribe: listener => { listeners.add(listener); return () => { listeners.delete(listener); }; },
      get: scope => active(scope) ? read(scope) : undefined,
      async load(scope) {
        if (!active(scope) || read(scope)) return;
        publish(scope, { status: 'loading' });
        try { const snapshot = structuredClone(await pointFormatAPI.read(scope)); publish(scope, { status: 'ready', snapshot }); }
        catch (cause) { publish(scope, { status: 'error', error: safePointFormatError(cause) }); }
      },
      async restart(scope) {
        if (!active(scope) || !canRestartPointOrder(read(scope))) return;
        entries.delete(keyOf(scope));
        await port.current!.load(scope);
      },
      choose(scope, intent) {
        const entry = read(scope);
        if (!active(scope) || entry?.status !== 'ready' || !entry.snapshot?.editable || entry.operation || (intent.format_mode === 'override' && !entry.snapshot.resolution.applicable)) return;
        publish(scope, { ...entry, intent });
      },
      save: scope => inspectOrSave(scope, false), inspect: scope => inspectOrSave(scope, true),
      async recover(scope) {
        const entry = read(scope);
        if (!active(scope) || !entry?.snapshot || !entry.operation || !entry.saved || !['saved', 'recovery-unknown'].includes(entry.status)) return;
        let recovery = entry.recovery;
        publish(scope, { ...entry, status: 'recovering', error: undefined });
        try {
          if (!recovery) {
            const settings = await modbusShareAPI.status();
            if (settings.workspace_id !== undefined && settings.workspace_id !== scope.workspace_id) throw new PointFormatError('point_format_scope', 'rejected');
            recovery = createPointFormatRecovery(entry.snapshot, entry.operation, entry.saved, settings.settings_revision ?? '');
            publish(scope, { ...entry, recovery, status: 'recovering' });
          }
          if (!active(scope)) {
            publish(scope, { ...entry, recovery, status: 'recovery-unknown', error: new PointFormatError('point_format_scope', 'rejected', entry.operation.request_id) });
            return;
          }
          const recovered = await pointFormatRecoveryAPI.recover(recovery);
          publish(scope, { ...entry, recovery, recovered, status: 'recovered' });
        } catch (cause) {
          const error = safePointFormatError(cause, entry.operation.request_id);
          publish(scope, { ...entry, recovery, error, status: error.outcome === 'unknown' ? 'recovery-unknown' : 'error' });
        }
      },
    };
  }
  return port.current;
}
