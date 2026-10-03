import * as React from 'react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  useWriteGroupSchemaApplyMutation,
  useWriteGroupSchemaOperationQuery,
  useWriteGroupSchemaPreviewMutation,
} from '../../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { SchemaOperation, SchemaPreviewToken } from '../../../../../../types/recordingPlan';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';
import { getSafeErrorMessage } from '../../../../../../utils/typedErrors';

export interface GroupSchemaPanelProps {
  /** Canonical saved group; absent means a new draft which cannot preview yet. */
  group: WriteGroup | undefined;
  workspaceId: string;
  workspaceRevision: string;
  /** Whether the draft still points at the current saved connector identity. */
  destinationMatches: boolean;
  dirty: boolean;
  readonly: boolean;
  /** Refetch canonical group, metadata and readiness after a verified apply. */
  onApplied?: () => void;
}

function isTerminal(operation: SchemaOperation | null): boolean {
  return operation !== null && !['pending', 'running'].includes(operation.status);
}

function scopeGroupId(scopeKey: string): string {
  return scopeKey.split('|', 1)[0] ?? '';
}

function operationStatus(operation: SchemaOperation | null, t: (key: string) => string): string {
  return operation ? t(`step4.group.schema.operation_status.${operation.status}`) : t('step4.group.schema.operation_unconfirmed');
}

function operationDetail(
  operation: SchemaOperation,
  t: (key: string, options?: Record<string, unknown>) => string,
): string | undefined {
  const code = operation.reason === 'permission_denied' ? 'RECORDING_SCHEMA_PERMISSION_DENIED' : operation.reason;
  const safe = getSafeErrorMessage({ code, action: operation.next_action }, t);
  if (safe.code) return safe.message;
  if (safe.action) return safe.action;
  return operation.status === 'succeeded' ? undefined : t('step4.group.schema.operation_attention');
}

/**
 * Managed schema preparation uses the saved group as its only input. Preview
 * is read-only; confirmation sends only the server-issued token and operation
 * identity, and unresolved outcomes are checked through that same operation.
 */
export const GroupSchemaPanel: React.FC<GroupSchemaPanelProps> = ({
  group, workspaceId, workspaceRevision, destinationMatches, dirty, readonly, onApplied,
}) => {
  const { t } = useTranslation('workbench-v2');
  const previewMutation = useWriteGroupSchemaPreviewMutation();
  const applyMutation = useWriteGroupSchemaApplyMutation();
  const scopeKey = useMemo(() => group
    ? [
      group.id,
      group.revision,
      group.destination.connector_revision,
      group.destination.schema_revision ?? '',
      group.destination.schema_digest ?? '',
      workspaceRevision,
    ].join('|')
    : `unsaved|${workspaceId}|${workspaceRevision}`, [group, workspaceId, workspaceRevision]);
  const operationId = useRef<string | undefined>(undefined);
  const scopeRef = useRef(scopeKey);
  scopeRef.current = scopeKey;
  const previousScopeKey = useRef(scopeKey);
  const inFlight = useRef(false);
  const operationRef = useRef<SchemaOperation | null>(null);
  const [preview, setPreview] = useState<SchemaPreviewToken | null>(null);
  const [operation, setOperation] = useState<SchemaOperation | null>(null);
  const [lastOperation, setLastOperation] = useState<SchemaOperation | null>(null);
  const [error, setError] = useState<unknown>(null);
  const [needsCheck, setNeedsCheck] = useState(false);
  const operationQuery = useWriteGroupSchemaOperationQuery(operationId.current);

  useEffect(() => {
    const terminal = operationRef.current;
    const sameGroup = scopeGroupId(previousScopeKey.current) === scopeGroupId(scopeKey);
    if (sameGroup && terminal !== null && isTerminal(terminal)) {
      setLastOperation(terminal);
      operationId.current = terminal.operation_id;
    } else if (!sameGroup) {
      setLastOperation(null);
    }
    previousScopeKey.current = scopeKey;
    setPreview(null);
    setOperation(null);
    operationRef.current = null;
    if (!sameGroup) operationId.current = undefined;
    setError(null);
    if (!sameGroup) setNeedsCheck(false);
    inFlight.current = false;
  }, [scopeKey]);

  const blockedReason = !group
    ? t('step4.group.schema.blocked_unsaved')
    : readonly
      ? t('step4.group.schema.blocked_readonly')
      : dirty
        ? t('step4.group.schema.blocked_dirty')
        : !destinationMatches
          ? t('step4.group.schema.blocked_destination')
          : group.status === 'deleted'
            ? t('step4.group.schema.blocked_deleted')
            : undefined;
  const enabled = blockedReason === undefined;
  const busy = previewMutation.isPending || applyMutation.isPending || inFlight.current;
  const canApply = enabled && preview !== null && operation === null && !needsCheck;
  const columns = preview?.group_layout?.columns.filter((column) => column.name !== preview.group_layout?.owner_column) ?? [];
  const effectiveOperation = operation ?? lastOperation;
  const shouldOfferCheck = needsCheck || (effectiveOperation !== null && effectiveOperation.status !== 'succeeded');
  const detail = effectiveOperation ? operationDetail(effectiveOperation, t) : undefined;

  const run = async (action: () => Promise<void>) => {
    if (inFlight.current) return;
    inFlight.current = true;
    const requestedScope = scopeRef.current;
    setError(null);
    try {
      await action();
    } catch (cause) {
      if (scopeRef.current === requestedScope) setError(cause);
    } finally {
      if (scopeRef.current === requestedScope) inFlight.current = false;
    }
  };

  const handlePreview = () => run(async () => {
    if (!group || !enabled) return;
    const requestedScope = scopeRef.current;
    const result = await previewMutation.mutateAsync({
      id: group.id,
      request: {
        workspace_id: workspaceId,
        expected_workspace_revision: workspaceRevision,
        expected_group_revision: group.revision,
        expected_connector_revision: group.destination.connector_revision,
      },
    });
    if (scopeRef.current !== requestedScope) return;
    operationId.current = result.operation_id;
    setPreview(result);
    setOperation(null);
    setLastOperation(null);
    setNeedsCheck(false);
  });

  const handleApply = () => run(async () => {
    if (!group || !preview || !enabled) return;
    const requestedScope = scopeRef.current;
    try {
      const result = await applyMutation.mutateAsync({
        id: group.id,
        request: {
          workspace_id: workspaceId,
          expected_workspace_revision: workspaceRevision,
          expected_group_revision: group.revision,
          expected_connector_revision: group.destination.connector_revision,
          token: preview.token,
          operation_id: preview.operation_id,
        },
      });
      const sameScope = scopeRef.current === requestedScope;
      const sameGroup = scopeGroupId(scopeRef.current) === scopeGroupId(requestedScope);
      if (sameGroup) {
        operationId.current = result.operation_id;
        operationRef.current = result;
        if (sameScope) setOperation(result);
        else if (isTerminal(result)) setLastOperation(result);
      }
      setNeedsCheck(false);
      if (result.status === 'succeeded') onApplied?.();
    } catch (cause) {
      // A lost reply never authorizes another DDL batch. Check this same op.
      const sameGroup = scopeGroupId(scopeRef.current) === scopeGroupId(requestedScope);
      if (sameGroup) {
        operationId.current = preview.operation_id;
        setNeedsCheck(true);
        if (scopeRef.current !== requestedScope) {
          setError({ code: 'RECORDING_SCHEMA_OPERATION_UNKNOWN', retryable: false });
        }
      }
      throw cause;
    }
  });

  const handleCheck = () => run(async () => {
    const requestedScope = scopeRef.current;
    const result = await operationQuery.refetch();
    if (scopeRef.current !== requestedScope) return;
    if (result.error || !result.data) {
      setNeedsCheck(true);
      setError(result.error ?? { code: 'RECORDING_SCHEMA_OPERATION_UNKNOWN', retryable: false });
      return;
    }
    operationId.current = result.data.operation_id;
    operationRef.current = result.data;
    setOperation(result.data);
    setNeedsCheck(false);
    if (result.data.status === 'succeeded') onApplied?.();
  });

  return (
    <section
      aria-labelledby="group-schema-title"
      className="space-y-3 rounded-lg border border-slate-800 p-4"
      data-testid="group-schema-panel"
    >
      <div>
        <h4 id="group-schema-title" className="text-sm font-semibold text-slate-200">{t('step4.group.schema.title')}</h4>
        <p className="mt-1 text-xs text-slate-400">{t('step4.group.schema.intro')}</p>
      </div>

      {blockedReason && <p role="status" className="text-xs text-amber-300" data-testid="group-schema-blocked">{blockedReason}</p>}

      <div className="flex flex-wrap gap-2" role="group" aria-label={t('step4.group.schema.actions')}>
        <button
          type="button" onClick={handlePreview} disabled={!enabled || busy}
          className="rounded border border-slate-600 px-3 py-1 text-xs text-slate-100 disabled:opacity-40"
          data-testid="group-schema-preview"
        >
          {previewMutation.isPending ? t('step4.group.schema.previewing') : t('step4.group.schema.preview')}
        </button>
        {preview && (
          <button
            type="button" onClick={handleApply} disabled={!canApply || busy}
            className="rounded bg-blue-600 px-3 py-1 text-xs text-white disabled:opacity-40"
            data-testid="group-schema-apply"
          >
            {applyMutation.isPending ? t('step4.group.schema.applying') : t('step4.group.schema.confirm')}
          </button>
        )}
      </div>

      {preview && (
        <div className="max-h-72 space-y-3 overflow-auto rounded border border-slate-800 p-3 text-xs text-slate-300" data-testid="group-schema-preview-result">
          <div className="grid gap-1 sm:grid-cols-2" data-testid="group-schema-preview-table">
            <span className="text-slate-500">{t('step4.group.schema.table')}</span>
            <span className="font-mono">{preview.group_layout?.table_name ?? preview.table_prefix}</span>
            {preview.no_change_reason && <><span className="text-slate-500">{t('step4.group.schema.change')}</span><span>{t('step4.group.schema.no_change')}</span></>}
          </div>

          {columns.length > 0 && (
            <div className="overflow-x-auto" role="region" aria-label={t('step4.group.schema.columns_region')} tabIndex={0}>
              <table className="w-full min-w-[520px] border-collapse text-left" data-testid="group-schema-columns">
                <thead className="text-slate-500">
                  <tr>
                    <th scope="col" className="p-2">{t('step4.group.schema.column_name')}</th>
                    <th scope="col" className="p-2">{t('step4.group.schema.sql_type')}</th>
                    <th scope="col" className="p-2">{t('step4.group.schema.nullable')}</th>
                    <th scope="col" className="p-2">{t('step4.group.schema.primary_key')}</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800">
                  {columns.map((column) => (
                    <tr key={column.name}>
                      <td className="p-2 font-mono">{column.name}</td>
                      <td className="p-2 font-mono">{column.sql_type}</td>
                      <td className="p-2">{column.nullable ? t('step4.group.schema.yes') : t('step4.group.schema.no')}</td>
                      <td className="p-2">{column.primary_key ? t('step4.group.schema.yes') : t('step4.group.schema.no')}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          <details>
            <summary className="cursor-pointer text-slate-400">{t('step4.group.schema.advanced_sql')}</summary>
            <pre className="mt-2 max-w-full overflow-x-auto whitespace-pre-wrap text-[11px] text-slate-500">{preview.statements.join('\n')}</pre>
          </details>
          <p className="text-slate-500" role="status">{t('step4.group.schema.confirm_note')}</p>
        </div>
      )}

      {shouldOfferCheck && (
        <div className="flex flex-wrap items-center gap-2" role="status" data-testid="group-schema-operation-pending">
          <span className="text-xs text-amber-300">{t('step4.group.schema.check_required')}</span>
          <button type="button" onClick={handleCheck} disabled={busy} className="text-xs underline disabled:opacity-40" data-testid="group-schema-check">
            {t('step4.group.schema.check')}
          </button>
        </div>
      )}

      {effectiveOperation && (
        <div className="space-y-1 text-xs text-slate-300" role="status" data-testid="group-schema-operation">
          <p>{t('step4.group.schema.operation')}: <span className="font-mono">{effectiveOperation.operation_id}</span></p>
          <p>{t('step4.group.schema.status')}: <span className="font-mono">{operationStatus(effectiveOperation, t)}</span></p>
          {detail && <p>{t('step4.group.schema.reason')}: {detail}</p>}
        </div>
      )}

      {error !== null && (
        <p role="alert" className="text-xs text-red-300" data-testid="group-schema-error">
          {getSafeErrorMessage(error, t).message}
          {needsCheck && <span className="ml-1">{t('step4.group.schema.error_check_same_operation')}</span>}
        </p>
      )}
    </section>
  );
};
