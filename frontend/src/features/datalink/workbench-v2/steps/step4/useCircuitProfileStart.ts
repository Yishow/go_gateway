import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useRecordingStartMutation, useRecordingStartOperationQuery } from '../../../../../hooks/datalink/useStudioV2RecordingStart';
import type { RecordingStartOperation, RecordingStartRequest } from '../../../../../types/studioV2RecordingStart';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { basicRecordingIntentKey, createBasicRecordingRequestId, latestRecordingStartOperation, loadBasicRecordingIntent, saveBasicRecordingIntent, type BasicRecordingIntentRecord } from '../../state/basicRecording';
import { safeBasicStartError } from '../../state/basicRecordingMessages';

interface StartView {
  scope: string;
  intent?: BasicRecordingIntentRecord;
  submitted?: RecordingStartOperation;
  error?: string;
  queryFailed?: boolean;
}

interface StartInput {
  scopeKey: string;
  canStart: boolean;
  readonly: boolean;
  workspaceId?: string;
  deviceId: string;
  connectorRevision: string;
  readinessToken?: string;
  settingsRevision?: string;
  shareWorkspaceRevision?: string;
  prepare: () => Promise<{ groups: WriteGroup[]; workspaceRevision: string } | undefined>;
  onActivated?: (deviceIds: string[]) => Promise<void> | void;
  onOperationActiveChange?: (active: boolean) => void;
  /** Target preparation only: a recording start is never read back, restored or sent from this context. */
  prepareOnly?: boolean;
}

export function uncertainTransport(error: unknown): boolean {
  if (!error || typeof error !== 'object') return false;
  const record = error as Record<string, unknown>;
  if (record.outcome === 'unconfirmed') return true;
  if (['ERR_NETWORK', 'ECONNABORTED', 'ETIMEDOUT', 'ERR_CANCELED'].includes(String(record.code))) return true;
  const response = record.response as { status?: number } | undefined;
  if (typeof response?.status === 'number' && response.status >= 500) return true;
  if (typeof record.status === 'number' && record.status >= 500) return true;
  if (typeof record.code === 'string' && record.code.trim()) return false;
  return error instanceof TypeError || /fetch|network|timeout/i.test(String(record.message ?? ''));
}

/** Preserve the original request before posting: a missing reply is not a new intent. */
export function useCircuitProfileStart(input: StartInput) {
  const { t } = useTranslation('workbench-v2');
  const startMutation = useRecordingStartMutation();
  const [view, setView] = React.useState<StartView>(() => ({ scope: input.scopeKey, intent: loadBasicRecordingIntent(input.scopeKey) }));
  const [busy, setBusy] = React.useState(false);
  const inFlight = React.useRef(false);
  const mounted = React.useRef(true);
  React.useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const currentScope = React.useRef(input.scopeKey);
  currentScope.current = input.scopeKey;
  const scopedView = view.scope === input.scopeKey ? view : { scope: input.scopeKey, intent: loadBasicRecordingIntent(input.scopeKey) };
  const record = input.prepareOnly ? undefined : scopedView.intent;
  const operationId = record?.operation_id;
  const operationQuery = useRecordingStartOperationQuery(operationId);
  const operation = latestRecordingStartOperation(operationId, operationQuery.data, scopedView.submitted);
  const status = operation?.status ?? (record ? 'unknown' : undefined);
  const unresolved = Boolean(record && (!operation || ['pending', 'running', 'unknown'].includes(operation.status) ||
    operation.status !== 'succeeded' && operation.next_action === 'retry_same_request'));
  const onOperationActiveChange = input.onOperationActiveChange;

  React.useEffect(() => {
    setView({ scope: input.scopeKey, intent: loadBasicRecordingIntent(input.scopeKey) });
  }, [input.scopeKey]);
  React.useEffect(() => {
    onOperationActiveChange?.(busy || unresolved);
    return () => onOperationActiveChange?.(false);
  }, [onOperationActiveChange, busy, unresolved]);

  const run = async (action: (isCurrent: () => boolean) => Promise<void>) => {
    if (inFlight.current || input.readonly) return;
    const scope = input.scopeKey;
    inFlight.current = true;
    setBusy(true);
    try { await action(() => mounted.current && currentScope.current === scope); }
    finally { inFlight.current = false; if (mounted.current) setBusy(false); }
  };

  const start = async (retry = false) => {
    if (input.prepareOnly || (!retry && (!input.canStart || unresolved)) || (retry && !record)) return;
    await run(async (isCurrent) => {
      const scope = input.scopeKey;
      let intent = retry ? record : undefined;
      if (intent?.operation_id) {
        const result = await operationQuery.refetch();
        if (!isCurrent()) return;
        if (result.isError || result.error || !result.data) {
          setView((prev) => ({ ...prev, queryFailed: true }));
          return;
        }
        setView({ scope, intent, submitted: result.data });
        if (result.data.status === 'succeeded' || result.data.status === 'failed' && result.data.next_action !== 'retry_same_request') return;
      }
      if (!intent) {
        const resolved = await input.prepare();
        if (!resolved || !isCurrent() || !input.workspaceId) return;
        const request: RecordingStartRequest = {
          request_id: createBasicRecordingRequestId(), workspace_id: input.workspaceId,
          expected_workspace_revision: resolved.workspaceRevision,
          ...(input.readinessToken && input.settingsRevision && input.shareWorkspaceRevision ? {
            readiness_token: input.readinessToken, settings_revision: input.settingsRevision, workspace_revision: input.shareWorkspaceRevision,
          } : {}),
          device_ids: [input.deviceId],
          groups: resolved.groups.map((group) => ({ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: input.connectorRevision })),
        };
        intent = { scope_key: scope, request };
      }
      const submittedIntent = intent;
      saveBasicRecordingIntent(submittedIntent);
      setView({ scope, intent: submittedIntent });
      try {
        const nextOp = await startMutation.mutateAsync(submittedIntent.request);
        const savedIntent = { ...submittedIntent, operation_id: nextOp.operation_id };
        saveBasicRecordingIntent(savedIntent);
        if (!isCurrent()) return;
        setView({ scope, intent: savedIntent, submitted: nextOp });
        const activatedIds = nextOp.devices.filter((device) => device.activated).map((device) => device.device_id);
        if (activatedIds.length) await input.onActivated?.(activatedIds);
      } catch (error) {
        if (uncertainTransport(error)) return; // The persisted request remains unknown, including after reload.
        try { window.sessionStorage.removeItem(basicRecordingIntentKey(scope)); } catch { /* component state still records the failure */ }
        if (isCurrent()) setView({ scope, error: safeBasicStartError(error, t) });
      }
    });
  };

  const prepare = async () => {
    if (!input.canStart || unresolved) return;
    await run(async () => { await input.prepare(); });
  };
  const check = async () => {
    if (!operationId) return;
    await run(async (isCurrent) => {
      const result = await operationQuery.refetch();
      if (!isCurrent()) return;
      setView({ ...scopedView, submitted: result.data ?? scopedView.submitted, queryFailed: Boolean(result.isError || result.error || !result.data) });
    });
  };
  return { busy, unresolved, operation, operationId, status, queryFailed: scopedView.queryFailed,
    error: scopedView.error, canRetrySameRequest: unresolved && Boolean(record), start, prepare, check };
}
