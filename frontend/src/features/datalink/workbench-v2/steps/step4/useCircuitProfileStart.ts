import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useRecordingStartMutation, useRecordingStartOperationQuery } from '../../../../../hooks/datalink/useStudioV2RecordingStart';
import type { RecordingStartOperation, RecordingStartRequest } from '../../../../../types/studioV2RecordingStart';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { clearBasicRecordingIntentIfCurrent, completeBasicRecordingIntent, sameBasicRecordingIntent, createBasicRecordingRequestId, latestRecordingStartOperation, loadBasicRecordingIntent, saveBasicRecordingIntent, type BasicRecordingIntentRecord } from '../../state/basicRecording';
import { isOriginalRequestRetry, recordingStartRecoveryBlocked } from './recordingStartRecovery';
import { recordingCompletion } from './recordingCompletion';
import { useBasicRecordingIntentObservation } from './useBasicRecordingIntentObservation';
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
  /** Denial only for current non-start owners; never grants a retry. */
  retryBlocked?: () => boolean;
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
  const [online, setOnline] = React.useState(() => navigator.onLine);
  React.useEffect(() => {
    const update = () => setOnline(navigator.onLine);
    window.addEventListener('online', update); window.addEventListener('offline', update);
    return () => { window.removeEventListener('online', update); window.removeEventListener('offline', update); };
  }, []);
  const inFlight = React.useRef(false);
  const mounted = React.useRef(true);
  React.useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const currentScope = React.useRef(input.scopeKey);
  currentScope.current = input.scopeKey;
  const scopedView = view.scope === input.scopeKey ? view : { scope: input.scopeKey, intent: loadBasicRecordingIntent(input.scopeKey) };
  const record = input.prepareOnly ? undefined : scopedView.intent;
  const operationId = record?.operation_id;
  const retryOwner = React.useRef({ scope: input.scopeKey, epoch: 0, record });
  if (retryOwner.current.scope !== input.scopeKey) retryOwner.current = { scope: input.scopeKey, epoch: retryOwner.current.epoch + 1, record };
  else retryOwner.current.record = record;
  const renderRetryEpoch = retryOwner.current.epoch;
  const currentInput = React.useRef(input); currentInput.current = input;
  const recoveryBlocked = () => {
    try { return !navigator.onLine || currentInput.current.readonly || Boolean(currentInput.current.retryBlocked?.()) ||
      recordingStartRecoveryBlocked(input.scopeKey, input.workspaceId ?? '', record?.request.groups.map(group => group.group_id) ?? [], operationId); }
    catch { return true; }
  };
  const operationQuery = useRecordingStartOperationQuery(operationId);
  const operation = latestRecordingStartOperation(operationId, operationQuery.data, scopedView.submitted);
  const completion = recordingCompletion(record, operation, input.workspaceId, input.deviceId);
  const status = operation?.status ?? (record ? 'unknown' : undefined);
  const unresolved = Boolean(record && (!operation || completion.state === 'unavailable' || ['pending', 'running', 'partial', 'unknown'].includes(operation.status) ||
    operation.status !== 'succeeded' && operation.next_action === 'retry_same_request'));
  const onOperationActiveChange = input.onOperationActiveChange;

  React.useEffect(() => {
    setView({ scope: input.scopeKey, intent: loadBasicRecordingIntent(input.scopeKey) });
  }, [input.scopeKey]);
  useBasicRecordingIntentObservation(input.scopeKey, record, intent => setView(previous =>
    previous.scope === input.scopeKey && sameBasicRecordingIntent(previous.intent, intent) ? { ...previous, intent } : previous));
  React.useEffect(() => {
    onOperationActiveChange?.(busy || unresolved);
    return () => onOperationActiveChange?.(false);
  }, [onOperationActiveChange, busy, unresolved]);

  const run = async (action: (isCurrent: () => boolean) => Promise<void>, readonlyQuery = false) => {
    if (inFlight.current || input.readonly && !readonlyQuery) return;
    const scope = input.scopeKey;
    inFlight.current = true;
    setBusy(true);
    try { await action(() => mounted.current && currentScope.current === scope); }
    finally { inFlight.current = false; if (mounted.current) setBusy(false); }
  };

  const start = async (retry = false) => {
    if (input.prepareOnly || (!retry && (!input.canStart || unresolved)) || (retry && (!record || !operationId || retryOwner.current.epoch !== renderRetryEpoch || retryOwner.current.record !== record || recoveryBlocked()))) return;
    const retryEpoch = renderRetryEpoch;
    const acceptedRetryGuard = input.retryBlocked;
    const acceptedRetryBlocked = () => { try { return Boolean(acceptedRetryGuard?.()); } catch { return true; } };
    await run(async (isCurrent) => {
      const scope = input.scopeKey;
      let intent = retry ? record : undefined;
      if (intent?.operation_id) {
        const result = await operationQuery.refetch();
        if (!isCurrent() || retryOwner.current.epoch !== retryEpoch || retryOwner.current.record !== record) return;
        if (result.isError || result.error || !result.data) {
          setView((prev) => ({ ...prev, queryFailed: true }));
          return;
        }
        setView({ scope, intent, submitted: result.data });
        if (!isOriginalRequestRetry(result.data, operationId, intent.request) || recoveryBlocked() || acceptedRetryBlocked()) return;
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
        if (completeBasicRecordingIntent(submittedIntent, nextOp.operation_id) === 'superseded' || !isCurrent()) return;
        setView({ scope, intent: savedIntent, submitted: nextOp });
        const activatedIds = nextOp.devices.filter((device) => device.activated).map((device) => device.device_id);
        if (activatedIds.length) await input.onActivated?.(activatedIds);
      } catch (error) {
        if (uncertainTransport(error)) return; // The persisted request remains unknown, including after reload.
        clearBasicRecordingIntentIfCurrent(submittedIntent);
        if (isCurrent()) setView({ scope, error: safeBasicStartError(error, t) });
      }
    });
  };

  const prepare = async () => {
    if (!input.canStart || unresolved) return;
    await run(async () => { await input.prepare(); });
  };
  const check = async () => {
    if (!operationId || !navigator.onLine || retryOwner.current.epoch !== renderRetryEpoch || retryOwner.current.record !== record) return;
    await run(async (isCurrent) => {
      const result = await operationQuery.refetch();
      if (!isCurrent() || retryOwner.current.epoch !== renderRetryEpoch || retryOwner.current.record !== record) return;
      setView({ ...scopedView, submitted: result.data ?? scopedView.submitted, queryFailed: Boolean(result.isError || result.error || !result.data) });
    }, true);
  };
  return { busy, unresolved, operation, operationId, status, completion, queryFailed: scopedView.queryFailed,
    error: scopedView.error, canRetrySameRequest: online && !input.prepareOnly && !operationQuery.isError && !operationQuery.isFetching &&
      isOriginalRequestRetry(operationQuery.data, operationId, record?.request) && !recoveryBlocked(), start, prepare, check };
}
