import * as React from 'react';
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  useWriteGroupTestWriteMutation,
  useWriteGroupTestWriteOperationQuery,
  useWriteGroupTestWritePreviewMutation,
} from '../../../../../../hooks/datalink/useStudioV2WriteGroupTestWrite';
import type {
  WriteGroupTestWriteOperation,
  WriteGroupTestWritePreview,
} from '../../../../../../types/studioV2WriteGroupTestWrite';
import { recordingTestWriteRecovery } from '../../../state/recordingTestWriteRecovery';
import { useRecordingTestWriteRecovery } from '../../../../../../hooks/datalink/useRecordingTestWriteRecovery';
import { RecordingTestWriteOperationResult } from '../RecordingTestWriteOperationResult';
import { RecoveryEntries } from '../RecordingTestWriteRecoveryPanel';
import { getSafeErrorMessage } from '../../../../../../utils/typedErrors';

export interface GroupTestWritePanelProps {
  workspaceId: string;
  groupId: string;
  /** Changes whenever the saved group does; an earlier preview is dropped, but recovery is retained. */
  scopeKey: string;
  /** False while the group has unsaved edits, an unsaved destination or a read-only session. */
  enabled: boolean;
  disabledReason?: string;
}

/**
 * Runs the explicit, confirmed test write of one saved group. Write
 * verification and cleanup are separate facts, an unconfirmed result is
 * checked rather than repeated, and nothing is inferred from a missing field.
 */
export const GroupTestWritePanel: React.FC<GroupTestWritePanelProps> = ({ workspaceId, groupId, scopeKey, enabled, disabledReason }) => {
  const { t } = useTranslation('workbench-v2');
  const recoveryList = useRecordingTestWriteRecovery(workspaceId);
  const previewMutation = useWriteGroupTestWritePreviewMutation();
  const confirmMutation = useWriteGroupTestWriteMutation();
  const [preview, setPreview] = useState<WriteGroupTestWritePreview | null>(null);
  const [operation, setOperation] = useState<WriteGroupTestWriteOperation | null>(null);
  const [error, setError] = useState<unknown>(null);
  const [unconfirmed, setUnconfirmed] = useState(false);
  const identity = `${workspaceId}|${groupId}|${scopeKey}`;
  const scopeRef = useRef({ identity, epoch: 0 });
  if (scopeRef.current.identity !== identity) scopeRef.current = { identity, epoch: scopeRef.current.epoch + 1 };
  const previewScope = useRef<typeof scopeRef.current | null>(null);
  const [recovery, setRecovery] = useState(() => recordingTestWriteRecovery.read());
  const references = recovery.references.filter((ref) => ref.workspace_id === workspaceId);
  const reference = references.find((ref) => ref.group_id === groupId);
  const refreshRecovery = () => setRecovery(recordingTestWriteRecovery.read());
  useEffect(() => {
    refreshRecovery();
    return recordingTestWriteRecovery.subscribe(refreshRecovery);
  }, []);
  useEffect(() => {
    window.addEventListener('storage', refreshRecovery);
    return () => window.removeEventListener('storage', refreshRecovery);
  }, []);
  const inFlight = useRef(false);
  const operationQuery = useWriteGroupTestWriteOperationQuery(reference?.operation_id ?? operation?.operation_id ?? preview?.operation_id, workspaceId);

  // A reply for an earlier group or revision must never land on the current one.
  useEffect(() => {
    setPreview(null);
    previewScope.current = null;
    setOperation(null);
    setError(null);
    setUnconfirmed(Boolean(reference));
    inFlight.current = false;
    // Recover on entry/reload only; retaining an operation before dispatch must not trigger a racing read.
    if (reference) {
      const scope = scopeRef.current;
      void operationQuery.refetch().then((result) => {
        if (scopeRef.current !== scope || result.isError || result.data?.operation_id !== reference.operation_id) return;
        setOperation(result.data);
        setUnconfirmed(false);
        recordingTestWriteRecovery.resolve(result.data, workspaceId);
      });
    }
    // The editor identity, not query state, defines entry into a recovery scope.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [identity]);

  const run = async (action: () => Promise<void>) => {
    if (inFlight.current) return;
    inFlight.current = true;
    const scope = scopeRef.current;
    setError(null);
    try {
      await action();
    } catch (cause) {
      if (scopeRef.current === scope) setError(cause);
    } finally {
      if (scopeRef.current === scope) inFlight.current = false;
    }
  };

  const handlePreview = () => run(async () => {
    if (!enabled || reference || recovery.unavailable || !recoveryList.available) return;
    const scope = scopeRef.current;
    const result = await previewMutation.mutateAsync(groupId);
    if (scopeRef.current !== scope) return;
    if (result.group_id !== groupId) throw new Error('Test-write preview identity was not confirmed.');
    previewScope.current = scope;
    setPreview(result);
    setOperation(null);
    setUnconfirmed(false);
  });

  const handleConfirm = () => run(async () => {
    if (!preview || reference || !enabled || !recoveryList.available) return;
    const scope = scopeRef.current;
    const retained = recordingTestWriteRecovery.retain({
      workspace_id: workspaceId, group_id: groupId, group_revision: preview.group_revision,
      operation_id: preview.operation_id, operation_kind: 'test_write', created_at: new Date().toISOString(),
    });
    setUnconfirmed(true);
    if (!retained) return;
    try {
      const result = await confirmMutation.mutateAsync({ groupId, confirmation: { token: preview.token, operation_id: preview.operation_id } });
      if (result.operation_id !== preview.operation_id) return;
      recordingTestWriteRecovery.resolve(result, workspaceId);
      if (scopeRef.current !== scope) return;
      setOperation(result);
      setUnconfirmed(false);
    } catch (cause) {
      // A lost or untyped reply may hide a write that did happen: check the operation, never resend blindly.
      if (scopeRef.current === scope) setUnconfirmed(true);
      throw cause;
    }
  });

  const handleRepairCleanup = () => run(async () => {
    if (!preview || !operation || !canRepairCleanup || !enabled || unconfirmed) return;
    const scope = scopeRef.current;
    setUnconfirmed(true);
    try {
      const result = await confirmMutation.mutateAsync({
        groupId,
        confirmation: { token: preview.token, operation_id: preview.operation_id, action: 'cleanup_repair' },
      });
      if (result.operation_id !== preview.operation_id) return;
      recordingTestWriteRecovery.resolve(result, workspaceId);
      if (scopeRef.current !== scope) return;
      setOperation(result);
      setUnconfirmed(false);
    } catch (cause) {
      if (scopeRef.current === scope) setUnconfirmed(true);
      throw cause;
    }
  });

  const handleCheck = () => run(async () => {
    const scope = scopeRef.current;
    const result = await operationQuery.refetch();
    if (scopeRef.current !== scope) return;
    if (result.isError || !result.data) {
      setUnconfirmed(true);
      return;
    }
    if (result.data.operation_id === (reference?.operation_id ?? operation?.operation_id ?? preview?.operation_id)) {
      recordingTestWriteRecovery.resolve(result.data, workspaceId);
      setOperation(result.data);
      setUnconfirmed(false);
    }
  });

  const busy = previewMutation.isPending || confirmMutation.isPending;
  const finished = operation !== null && operation.status !== 'pending' && operation.status !== 'running';
  const canRepairCleanup =
    finished &&
    operation !== null &&
    operation.write_outcome === 'written_verified' &&
    (operation.cleanup_status === 'failed' || operation.cleanup_status === 'unknown') &&
    preview !== null &&
    preview.group_id === groupId &&
    preview.operation_id === operation.operation_id &&
    previewScope.current === scopeRef.current &&
    preview.target.dialect === 'sqlserver';

  return (
    <section className="space-y-3 rounded-lg border border-slate-800 p-4" aria-labelledby={`group-test-write-${groupId}`} data-testid="group-test-write">
      <h4 id={`group-test-write-${groupId}`} className="text-sm font-semibold text-slate-200">{t('step4.group.test_write.title')}</h4>
      <p className="text-xs text-slate-400">{t('step4.group.test_write.intro')}</p>

      {!enabled && disabledReason && <p className="text-xs text-amber-300" data-testid="group-test-write-blocked">{disabledReason}</p>}

      <button
        type="button" onClick={handlePreview} disabled={!enabled || busy || Boolean(reference) || recovery.unavailable || !recoveryList.available}
        className="rounded border border-slate-600 px-3 py-1 text-xs text-slate-100 disabled:opacity-40" data-testid="group-test-write-preview"
      >
        {t('step4.group.test_write.preview')}
      </button>

      {!recoveryList.available && <div role="alert" data-testid="group-test-write-list-unavailable">
        {t('step4.group.test_write.recovery_list_unavailable')}
        <button type="button" disabled={recoveryList.isFetching} onClick={() => void recoveryList.refetch()} data-testid="group-test-write-list-retry">
          {t('step4.group.test_write.check')}
        </button>
      </div>}
      {recovery.unavailable && <p role="alert" data-testid="group-test-write-storage-error">{t('step4.group.test_write.recovery_storage_error')}</p>}
      <RecoveryEntries scopeKey={identity} references={references} />
      {preview && !operation && (
        <div className="space-y-2 text-xs text-slate-300" data-testid="group-test-write-preview-card">
          <p>{t('step4.group.test_write.target', { table: preview.target.table })}</p>
          {preview.record_ids && <p data-testid="group-test-write-row-count">{t('step4.group.test_write.long_preview_count', { count: preview.record_ids.length })}</p>}
          <ul className="list-disc pl-5">
            {preview.values.map((entry) => <li key={entry.column}><span className="font-mono">{entry.column}</span> = {entry.value}</li>)}
          </ul>
          <p>{t(preview.record_ids ? 'step4.group.test_write.long_cleanup_note' : 'step4.group.test_write.cleanup_note')}</p>
          <button
            type="button" onClick={handleConfirm} disabled={!enabled || busy || unconfirmed || !recoveryList.available}
            className="rounded bg-blue-600 px-3 py-1 text-xs text-white disabled:opacity-40" data-testid="group-test-write-confirm"
          >
            {t('step4.group.test_write.confirm')}
          </button>
        </div>
      )}

      {operation && <RecordingTestWriteOperationResult operation={operation}
        manualGuidance={finished && !canRepairCleanup && (operation.cleanup_status === 'failed' || operation.cleanup_status === 'unknown')} />}
      {canRepairCleanup && (
        <button
          type="button" onClick={handleRepairCleanup} disabled={!enabled || busy || unconfirmed}
          className="rounded border border-amber-600 px-3 py-1 text-xs text-amber-200 hover:bg-amber-950/40 disabled:opacity-40" data-testid="group-test-write-cleanup-repair"
        >
          {t('step4.group.test_write.cleanup_repair')}
        </button>
      )}
      {operation && reference && (
        <button type="button" onClick={handleCheck} className="text-xs underline" data-testid="group-test-write-check">{t('step4.group.test_write.check')}</button>
      )}

      {unconfirmed && (
        <div role="alert" className="text-xs text-amber-300" data-testid="group-test-write-unconfirmed">
          {t(operation ? 'step4.group.test_write.operation_status_unconfirmed' : 'step4.group.test_write.unconfirmed')}
          <button type="button" onClick={handleCheck} disabled={!reference && !preview && !operation} className="ml-2 underline" data-testid="group-test-write-check-unconfirmed">{t('step4.group.test_write.check')}</button>
        </div>
      )}
      {error !== null && !unconfirmed && (
        <p role="alert" className="text-xs text-red-300" data-testid="group-test-write-error">{getSafeErrorMessage(error, t).message}</p>
      )}
    </section>
  );
};
