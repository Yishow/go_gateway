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
import { normalizeTypedEnvelope } from '../../../../../../utils/safeJson';
import { getSafeErrorMessage } from '../../../../../../utils/typedErrors';

/** Only a typed answer other than "result unknown" proves nothing was written; anything else must be checked. */
function isUnconfirmed(error: unknown): boolean {
  const code = normalizeTypedEnvelope(error).code;
  return !code || code === 'WRITE_GROUP_TEST_WRITE_RESULT_UNKNOWN';
}

export interface GroupTestWritePanelProps {
  groupId: string;
  /** Changes whenever the saved group does; an earlier preview or result is dropped with it. */
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
export const GroupTestWritePanel: React.FC<GroupTestWritePanelProps> = ({ groupId, scopeKey, enabled, disabledReason }) => {
  const { t } = useTranslation('workbench-v2');
  const previewMutation = useWriteGroupTestWritePreviewMutation();
  const confirmMutation = useWriteGroupTestWriteMutation();
  const [preview, setPreview] = useState<WriteGroupTestWritePreview | null>(null);
  const [operation, setOperation] = useState<WriteGroupTestWriteOperation | null>(null);
  const [error, setError] = useState<unknown>(null);
  const [unconfirmed, setUnconfirmed] = useState(false);
  const scopeRef = useRef(scopeKey);
  const inFlight = useRef(false);
  const operationQuery = useWriteGroupTestWriteOperationQuery(operation?.operation_id ?? preview?.operation_id);

  // A reply for an earlier group or revision must never land on the current one.
  useEffect(() => {
    scopeRef.current = scopeKey;
    setPreview(null);
    setOperation(null);
    setError(null);
    setUnconfirmed(false);
    inFlight.current = false;
  }, [scopeKey]);

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
    const scope = scopeRef.current;
    const result = await previewMutation.mutateAsync(groupId);
    if (scopeRef.current !== scope) return;
    setPreview(result);
    setOperation(null);
    setUnconfirmed(false);
  });

  const handleConfirm = () => run(async () => {
    if (!preview) return;
    const scope = scopeRef.current;
    try {
      const result = await confirmMutation.mutateAsync({ groupId, confirmation: { token: preview.token, operation_id: preview.operation_id } });
      if (scopeRef.current !== scope) return;
      setOperation(result);
    } catch (cause) {
      // A lost or untyped reply may hide a write that did happen: check the operation, never resend blindly.
      if (scopeRef.current === scope && isUnconfirmed(cause)) setUnconfirmed(true);
      throw cause;
    }
  });

  const handleCheck = () => run(async () => {
    const scope = scopeRef.current;
    const result = await operationQuery.refetch();
    if (scopeRef.current !== scope) return;
    if (result.data) {
      setOperation(result.data);
      setUnconfirmed(false);
    }
  });

  const busy = previewMutation.isPending || confirmMutation.isPending;
  const finished = operation !== null && operation.status !== 'pending' && operation.status !== 'running';

  return (
    <section className="space-y-3 rounded-lg border border-slate-800 p-4" aria-labelledby={`group-test-write-${groupId}`} data-testid="group-test-write">
      <h4 id={`group-test-write-${groupId}`} className="text-sm font-semibold text-slate-200">{t('step4.group.test_write.title')}</h4>
      <p className="text-xs text-slate-400">{t('step4.group.test_write.intro')}</p>

      {!enabled && disabledReason && <p className="text-xs text-amber-300" data-testid="group-test-write-blocked">{disabledReason}</p>}

      <button
        type="button" onClick={handlePreview} disabled={!enabled || busy}
        className="rounded border border-slate-600 px-3 py-1 text-xs text-slate-100 disabled:opacity-40" data-testid="group-test-write-preview"
      >
        {t('step4.group.test_write.preview')}
      </button>

      {preview && !operation && (
        <div className="space-y-2 text-xs text-slate-300" data-testid="group-test-write-preview-card">
          <p>{t('step4.group.test_write.target', { table: preview.target.table })}</p>
          <ul className="list-disc pl-5">
            {preview.values.map((entry) => <li key={entry.column}><span className="font-mono">{entry.column}</span> = {entry.value}</li>)}
          </ul>
          <p>{t('step4.group.test_write.cleanup_note')}</p>
          <button
            type="button" onClick={handleConfirm} disabled={!enabled || busy || unconfirmed}
            className="rounded bg-blue-600 px-3 py-1 text-xs text-white disabled:opacity-40" data-testid="group-test-write-confirm"
          >
            {t('step4.group.test_write.confirm')}
          </button>
        </div>
      )}

      {operation && (
        <dl className="grid grid-cols-[auto,1fr] gap-x-3 gap-y-1 text-xs text-slate-300" data-testid="group-test-write-result">
          <dt>{t('step4.group.test_write.operation')}</dt><dd className="font-mono">{operation.operation_id}</dd>
          {!finished ? (
            <><dt>{t('step4.group.test_write.status')}</dt><dd data-testid="group-test-write-running">{t('step4.group.test_write.running')}</dd></>
          ) : (
            <>
              <dt>{t('step4.group.test_write.write_outcome')}</dt>
              <dd data-testid="group-test-write-outcome">{t(`step4.group.test_write.outcome_${operation.write_outcome}`)}{operation.reason ? ` (${operation.reason})` : ''}</dd>
              <dt>{t('step4.group.test_write.cleanup_status')}</dt>
              <dd data-testid="group-test-write-cleanup">{t(`step4.group.test_write.cleanup_${operation.cleanup_status}`)}{operation.cleanup_reason ? ` (${operation.cleanup_reason})` : ''}</dd>
            </>
          )}
        </dl>
      )}
      {operation && !finished && (
        <button type="button" onClick={handleCheck} className="text-xs underline" data-testid="group-test-write-check">{t('step4.group.test_write.check')}</button>
      )}

      {unconfirmed && (
        <div role="alert" className="text-xs text-amber-300" data-testid="group-test-write-unconfirmed">
          {t('step4.group.test_write.unconfirmed')}
          <button type="button" onClick={handleCheck} disabled={!preview && !operation} className="ml-2 underline" data-testid="group-test-write-check-unconfirmed">{t('step4.group.test_write.check')}</button>
        </div>
      )}
      {error !== null && !unconfirmed && (
        <p role="alert" className="text-xs text-red-300" data-testid="group-test-write-error">{getSafeErrorMessage(error, t).message}</p>
      )}
    </section>
  );
};
