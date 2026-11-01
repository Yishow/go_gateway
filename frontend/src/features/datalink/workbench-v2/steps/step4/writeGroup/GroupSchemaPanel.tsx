import { RecordingSchemaImpact } from '../../../recording-editor/RecordingSchemaImpact';
import type { GuidedSchemaPresentation, SchemaPreparationBlocker, SchemaPreparationPresentation } from './schemaPreparationPresentation';
import { inspectSchemaOperation, readSchemaOperation, retainSchemaOperation } from './schemaOperationRecovery';
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
import { destinationBlockerCode, type DestinationBlockerCode } from '../../../state/circuitProfileBlockers';
import { isTerminal, operationDetail, operationStatus, schemaIssueCode, scopeGroupId } from './schemaPanelHelpers';

export interface GroupSchemaPanelProps {
  /** Canonical saved group; absent means a new draft which cannot preview yet. */
  group: WriteGroup | undefined;
  autoPreview?: boolean;
  /** Guided navigation never starts a preview POST; only an explicit owner command does. */
  guided?: GuidedSchemaPresentation;
  previewRequest?: number;
  workspaceId: string;
  workspaceRevision: string;
  /** Whether the draft still points at the current saved connector identity. */
  destinationMatches: boolean;
  dirty: boolean;
  readonly: boolean;
  /** Refetch canonical group, metadata and readiness after a verified apply. */
  onApplied?: () => void;
  /** Reports a recognised destination cause instead of showing it here, so one cause is explained once by the owner. */
  onBlocker?: (groupId: string, code: DestinationBlockerCode | undefined) => void;
}

/**
 * Managed schema preparation uses the saved group as its only input. Preview
 * creates a ledger/token with a POST; confirmation sends only the server-issued token and operation
 * identity, and unresolved outcomes are checked through that same operation.
 */
export const GroupSchemaPanel: React.FC<GroupSchemaPanelProps> = ({
  group, workspaceId, workspaceRevision, destinationMatches, dirty, readonly, onApplied, onBlocker, autoPreview = false, previewRequest = 0, guided,
}) => {
  const { t } = useTranslation('workbench-v2');
  const titleId = React.useId();
  const previewMutation = useWriteGroupSchemaPreviewMutation();
  const applyMutation = useWriteGroupSchemaApplyMutation();
  const scopeKey = useMemo(() => group
    ? [
      group.id,
      workspaceId,
      group.workspace_id,
      group.destination.connector_id,
      group.revision,
      group.destination.connector_revision,
      group.destination.schema_revision ?? '',
      group.destination.schema_digest ?? '',
      workspaceRevision,
    ].join('|')
    : `unsaved|${workspaceId}|${workspaceRevision}`, [group, workspaceId, workspaceRevision]);
  const ownerIdentity = JSON.stringify([workspaceId, workspaceRevision, group?.workspace_id, group?.id, group?.revision,
    group?.destination.connector_id, group?.destination.connector_revision, group?.destination.schema_revision, group?.destination.schema_digest, dirty, destinationMatches, readonly]);
  const ownerScope = useRef({ identity: ownerIdentity, epoch: 0 });
  if (ownerScope.current.identity !== ownerIdentity) ownerScope.current = { identity: ownerIdentity, epoch: ownerScope.current.epoch + 1 };
  const previewEpoch = useRef<number | undefined>(undefined);
  const parentConsent = useRef(guided?.consentEpoch); parentConsent.current = guided?.consentEpoch;
  const previewConsent = useRef<number | undefined>(undefined);
  const operationId = useRef<string | undefined>(readSchemaOperation(workspaceId, group?.id));
  const scopeRef = useRef(scopeKey);
  scopeRef.current = scopeKey;
  const previousScopeKey = useRef(scopeKey);
  const currentWorkspace = useRef(workspaceId); currentWorkspace.current = workspaceId;
  const previousWorkspace = useRef(workspaceId);
  const inFlight = useRef(false);
  const operationRef = useRef<SchemaOperation | null>(null);
  const [preview, setPreview] = useState<SchemaPreviewToken | null>(null);
  const [operation, setOperation] = useState<SchemaOperation | null>(null);
  const [lastOperation, setLastOperation] = useState<SchemaOperation | null>(null);
  const [error, setError] = useState<unknown>(null);
  const [previewFailureEpoch, setPreviewFailureEpoch] = useState<number>();
  const [needsCheck, setNeedsCheck] = useState(Boolean(operationId.current));
  const [, refreshPresentation] = useState(0);
  const recoveredForQuery = guided ? inspectSchemaOperation(workspaceId, group?.id) : undefined;
  const queryOperationId = recoveredForQuery?.state === 'retained' ? recoveredForQuery.operationId : operationId.current;
  const operationQuery = useWriteGroupSchemaOperationQuery(queryOperationId);
  useEffect(() => {
    const terminal = operationRef.current;
    const sameGroup = previousWorkspace.current === workspaceId && scopeGroupId(previousScopeKey.current) === scopeGroupId(scopeKey);
    previousWorkspace.current = workspaceId;
    if (sameGroup && terminal !== null && isTerminal(terminal)) {
      setLastOperation(terminal);
      operationId.current = terminal.operation_id;
    } else if (!sameGroup) {
      setLastOperation(null);
    }
    previousScopeKey.current = scopeKey;
    previewEpoch.current = undefined;
    setPreview(null);
    setOperation(null);
    operationRef.current = null;
    if (!sameGroup) operationId.current = readSchemaOperation(workspaceId, group?.id);
    setError(null);
    setPreviewFailureEpoch(undefined);
    if (!sameGroup) { setNeedsCheck(Boolean(operationId.current)); refreshPresentation(value => value + 1); }
    inFlight.current = false;
  }, [scopeKey, workspaceId, group?.id]);

  // Editing and reverting is a new consent context, even at the same saved revision.
  // Only the preview is invalidated; a submitted operation remains recoverable.
  const consentEpoch = ownerScope.current.epoch;
  useEffect(() => {
    if (previewEpoch.current !== consentEpoch) { setPreview(null); setPreviewFailureEpoch(undefined); }
  }, [consentEpoch]);
  // Parent-native changes retire only preview consent, not a known failed-preview retry or submitted operation.
  useEffect(() => { if (previewConsent.current !== guided?.consentEpoch) setPreview(null); }, [guided?.consentEpoch]);

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
  const canApply = enabled && preview !== null && previewEpoch.current === ownerScope.current.epoch && previewConsent.current === parentConsent.current && operation === null && !needsCheck;
  const effectiveOperation = operation ?? lastOperation;
  const shouldOfferCheck = needsCheck || (effectiveOperation !== null && effectiveOperation.status !== 'succeeded');
  const blocksNewPreview = needsCheck || ['pending', 'running', 'unknown', 'partial'].includes(effectiveOperation?.status ?? '');
  const detail = effectiveOperation ? operationDetail(effectiveOperation, t) : undefined;
  const blocker = onBlocker ? destinationBlockerCode(schemaIssueCode(error, effectiveOperation, t)) : undefined;
  const reportBlocker = useRef(onBlocker); reportBlocker.current = onBlocker;
  const groupId = group?.id;
  useEffect(() => {
    if (!groupId) return undefined;
    reportBlocker.current?.(groupId, blocker);
    return () => reportBlocker.current?.(groupId, undefined);
  }, [groupId, blocker]);

  const run = async (action: () => Promise<void>, origin?: 'preview') => {
    if (inFlight.current) return;
    inFlight.current = true;
    const requestedScope = scopeRef.current;
    const requestedEpoch = ownerScope.current.epoch, requestedConsent = parentConsent.current;
    setError(null);
    setPreviewFailureEpoch(undefined);
    try {
      await action();
    } catch (cause) {
      if (scopeRef.current === requestedScope && ownerScope.current.epoch === requestedEpoch && (origin !== 'preview' || parentConsent.current === requestedConsent)) {
        setError(cause);
        if (origin === 'preview') setPreviewFailureEpoch(requestedEpoch);
      }
    } finally {
      if (scopeRef.current === requestedScope) inFlight.current = false;
    }
  };

  const handlePreview = () => run(async () => {
    if (!group || !enabled || blocksNewPreview) return;
    const retained = readSchemaOperation(workspaceId, group.id);
    if (retained) { operationId.current = retained; setNeedsCheck(true); return; }
    const requestedScope = scopeRef.current;
    const requestedEpoch = ownerScope.current.epoch, requestedConsent = parentConsent.current;
    const result = await previewMutation.mutateAsync({
      id: group.id,
      request: {
        workspace_id: workspaceId,
        expected_workspace_revision: workspaceRevision,
        expected_group_revision: group.revision,
        expected_connector_revision: group.destination.connector_revision,
      },
    });
    if (scopeRef.current !== requestedScope || ownerScope.current.epoch !== requestedEpoch || parentConsent.current !== requestedConsent) return;
    previewConsent.current = requestedConsent;
    operationId.current = result.operation_id;
    previewEpoch.current = requestedEpoch;
    setPreview(result);
    setOperation(null);
    setLastOperation(null);
    setNeedsCheck(false);
  }, 'preview');

  const handleApply = () => run(async () => {
    if (!group || !preview || !enabled || !canApply || previewEpoch.current !== ownerScope.current.epoch) return;
    const retained = readSchemaOperation(workspaceId, group.id);
    if (retained) { operationId.current = retained; setNeedsCheck(true); return; }
    const requestedScope = scopeRef.current;
    retainSchemaOperation(workspaceId, group.id, preview.operation_id);
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
      const sameGroup = currentWorkspace.current === workspaceId && scopeGroupId(scopeRef.current) === scopeGroupId(requestedScope);
      if (sameGroup) {
        operationId.current = result.operation_id;
        operationRef.current = result;
        if (sameScope) setOperation(result);
        else if (isTerminal(result)) setLastOperation(result);
      }
      if (sameGroup) setNeedsCheck(result.status === 'unknown' || result.status === 'pending' || result.status === 'running');
      if (result.status === 'failed') retainSchemaOperation(workspaceId, group.id, undefined, result.operation_id);
      if (result.status === 'succeeded') { retainSchemaOperation(workspaceId, group.id, undefined, result.operation_id); onApplied?.(); }
    } catch (cause) {
      // A lost reply never authorizes another DDL batch. Check this same op.
      const sameGroup = currentWorkspace.current === workspaceId && scopeGroupId(scopeRef.current) === scopeGroupId(requestedScope);
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
    if (result.data.status === 'failed' && group) retainSchemaOperation(workspaceId, group.id, undefined, result.data.operation_id);
    if (result.data.status === 'succeeded') { if (group) retainSchemaOperation(workspaceId, group.id, undefined, result.data.operation_id); onApplied?.(); }
  });

  const autoStarted = useRef<number | undefined>(undefined);
  const previewAction = useRef(handlePreview); previewAction.current = handlePreview;
  useEffect(() => {
    if (!guided && autoPreview && enabled && autoStarted.current !== previewRequest) {
      autoStarted.current = previewRequest;
      if (!needsCheck) void previewAction.current();
    }
  }, [autoPreview, enabled, needsCheck, previewRequest, guided]);

  const bridgeEpoch = useRef({ identity: '', guard: '', preview, epoch: 0 });
  const guidedRef = useRef(guided); guidedRef.current = guided;
  const readNewCommandBlocked = () => {
    try { return Boolean(guidedRef.current?.newOperationBlocked || guidedRef.current?.isNewOperationBlocked?.()); } catch { return true; }
  };
  const readRetryBlocked = () => {
    try { return guidedRef.current?.isPreviewRetryBlocked?.() ?? true; } catch { return true; }
  };
  const previewFailure = error !== null && previewFailureEpoch === ownerScope.current.epoch;
  const newCommandBlocked = readNewCommandBlocked();
  const bridgeGuard = JSON.stringify([readonly, dirty, destinationMatches, guided?.online, guided?.consentEpoch, busy, needsCheck, effectiveOperation?.status, newCommandBlocked, previewFailure]);
  if (bridgeEpoch.current.identity !== ownerIdentity || bridgeEpoch.current.guard !== bridgeGuard || bridgeEpoch.current.preview !== preview) {
    bridgeEpoch.current = { identity: ownerIdentity, guard: bridgeGuard, preview, epoch: bridgeEpoch.current.epoch + 1 };
  }
  const live = useRef(true);
  useEffect(() => { live.current = true; return () => { live.current = false; }; }, []);
  const latestCommands = useRef({ preview: handlePreview, confirm: handleApply, check: handleCheck });
  latestCommands.current = { preview: handlePreview, confirm: handleApply, check: handleCheck };
  const inspection = inspectSchemaOperation(workspaceId, group?.id);
  const missingScope = !workspaceId || !workspaceRevision || !group?.id || !group.revision || !group.destination.connector_id || !group.destination.connector_revision;
  const foreignScope = group !== undefined && (group.workspace_id !== workspaceId || previousWorkspace.current !== workspaceId ||
    scopeGroupId(previousScopeKey.current) !== group.id);
  const retainedId = inspection.state === 'retained' ? inspection.operationId : undefined;
  const unresolvedId = retainedId ?? (needsCheck ? operationId.current : undefined);
  const observation: SchemaPreparationPresentation['observation'] = missingScope || foreignScope
    ? { owner: 'schema', state: 'unavailable', reason: 'scope_mismatch' }
    : unresolvedId || busy || (effectiveOperation && !['succeeded', 'failed'].includes(effectiveOperation.status))
      ? { owner: 'schema', state: 'unresolved', operationId: unresolvedId ?? (applyMutation.isPending ? operationId.current : effectiveOperation?.operation_id),
        reason: operationQuery.isError ? 'query_failed' : effectiveOperation?.status === 'running' ? 'running' : busy ? 'pending' : 'unknown' }
      : effectiveOperation && ['succeeded', 'failed'].includes(effectiveOperation.status)
        ? { owner: 'schema', state: 'observed', result: 'resolved', status: effectiveOperation.status === 'succeeded' ? 'succeeded' : 'failed', reason: 'terminal', operationId: effectiveOperation.operation_id }
        : inspection.state === 'unavailable' || error !== null
          ? { owner: 'schema', state: 'unavailable', reason: 'unreadable' }
          : { owner: 'schema', state: 'observed', result: 'none' };
  const commandBlocker: SchemaPreparationBlocker | undefined = missingScope ? 'missing_scope' : foreignScope ? 'scope_mismatch'
    : !guided?.online ? 'offline' : busy ? 'busy' : inspection.state === 'unavailable' && !unresolvedId ? 'unreadable'
      : readonly || newCommandBlocked ? 'readonly' : dirty ? 'dirty' : !destinationMatches ? 'destination' : group?.status === 'deleted' ? 'deleted'
        : observation.state === 'unresolved' ? 'operation_unresolved' : observation.state === 'unavailable' ? 'unreadable' : undefined;
  const readyPreview = commandBlocker === undefined && canApply && previewEpoch.current === ownerScope.current.epoch;
  const canRetryPreview = Boolean(guided && previewFailure && guided.online && !missingScope && !foreignScope && enabled && !busy &&
    !blocksNewPreview && inspection.state === 'empty' && !readRetryBlocked());
  const commandState = useRef({ commandBlocker, observation, canApply, canRetryPreview, groupId: group?.id, workspaceId, queryOperationId });
  commandState.current = { commandBlocker, observation, canApply, canRetryPreview, groupId: group?.id, workspaceId, queryOperationId };
  const epoch = bridgeEpoch.current.epoch;
  const observationKey = JSON.stringify(observation);
  const presentation = useMemo<SchemaPreparationPresentation>(() => {
    // Keep this publication's parent lease; a stale command must not borrow a later parent closure.
    const retryLease = guided?.isPreviewRetryBlocked;
    const invoke = async (command: 'preview' | 'confirm' | 'check' | 'retryPreview') => {
      if (!live.current || bridgeEpoch.current.epoch !== epoch || inFlight.current) return;
      const current = commandState.current;
      const recovered = inspectSchemaOperation(current.workspaceId, current.groupId);
      const held = recovered.state === 'retained' ? recovered.operationId : operationId.current;
      if (recovered.state === 'retained' && recovered.operationId !== operationId.current) {
        operationId.current = recovered.operationId; setNeedsCheck(true); refreshPresentation(value => value + 1);
        if (command !== 'check' || current.queryOperationId !== recovered.operationId) return;
      }
      if (command === 'check') {
        if (!held || ['offline', 'missing_scope', 'scope_mismatch', 'busy'].includes(current.commandBlocker ?? '')) return;
      } else if (command === 'retryPreview') {
        if (!current.canRetryPreview || recovered.state !== 'empty' || readRetryBlocked()) return;
        try { if (!retryLease || retryLease()) return; } catch { return; }
      } else if (current.commandBlocker || readNewCommandBlocked() || recovered.state !== 'empty' || (command === 'confirm' && (!current.canApply || previewEpoch.current !== ownerScope.current.epoch))) return;
      await latestCommands.current[command === 'retryPreview' ? 'preview' : command]();
    };
    const phase = observation.state === 'unavailable' ? 'unavailable' : observation.state === 'unresolved' ? 'unresolved'
      : observation.result === 'resolved' ? 'resolved' : readyPreview ? 'previewed' : 'idle';
    const canCheck = Boolean((unresolvedId || shouldOfferCheck) && operationId.current && guided?.online && !missingScope && !foreignScope && !busy);
    return { scope: { workspaceId, workspaceRevision, groupId: group?.id, groupRevision: group?.revision,
      connectorId: group?.destination.connector_id, connectorRevision: group?.destination.connector_revision, epoch },
      observation, phase, blocker: commandBlocker, preview: readyPreview ? 'ready' : 'none', ...(previewFailure ? { previewFailure: true } : {}),
      previewOperationId: readyPreview ? preview?.operation_id : undefined,
      commands: { preview: commandBlocker === undefined && !blocksNewPreview ? () => invoke('preview') : undefined,
        retryPreview: canRetryPreview ? () => invoke('retryPreview') : undefined,
        confirm: readyPreview ? () => invoke('confirm') : undefined, check: canCheck ? () => invoke('check') : undefined } };
    // The key is the complete token-free observation; owner handlers are read from current refs at invocation.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [epoch, observationKey, commandBlocker, readyPreview, blocksNewPreview, shouldOfferCheck, unresolvedId, canRetryPreview, previewFailure]);
  const isGuided = guided !== undefined;
  const notify = useRef(guided?.onPresentation); notify.current = guided?.onPresentation;
  useEffect(() => {
    if (!isGuided) return;
    const report = notify.current; report?.(presentation);
    return () => report?.(null);
  }, [isGuided, presentation]);
  const suppressActions = guided?.suppressDefaultActions || guided?.renderActions !== undefined;

  return (
    <section
      aria-labelledby={titleId}
      className="space-y-3 rounded-lg border border-[var(--gw-border)] p-4"
      data-testid="group-schema-panel"
    >
      <div>
        <h4 id={titleId} className="text-sm font-semibold text-[var(--gw-text)]">{t('step4.group.schema.title')}</h4>
        <p className="mt-1 text-xs text-[var(--gw-muted)]">{t('step4.group.schema.intro')}</p>
      </div>

      {blockedReason && <p role="status" className="text-xs text-[var(--gw-warning-fg)]" data-testid="group-schema-blocked">{blockedReason}</p>}

      {guided && !guided.renderActions && previewFailure && <button type="button" onClick={presentation.commands.retryPreview} disabled={!presentation.commands.retryPreview}
        className="rounded border border-[var(--gw-control)] px-3 py-1 text-xs text-[var(--gw-text)] disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]" data-testid="group-schema-preview">
        {t('step4.group.schema.preview')}
      </button>}
      {!suppressActions && <div className="flex flex-wrap gap-2" role="group" aria-label={t('step4.group.schema.actions')}>
        {!(guided && previewFailure) && <button
          type="button" onClick={guided ? presentation.commands.preview : handlePreview} disabled={guided ? !presentation.commands.preview : !enabled || busy || blocksNewPreview}
          className="rounded border border-[var(--gw-control)] px-3 py-1 text-xs text-[var(--gw-text)] disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]"
          data-testid="group-schema-preview"
        >
          {previewMutation.isPending ? t('step4.group.schema.previewing') : t('step4.group.schema.preview')}
        </button>}
        {preview && (
          <button
            type="button" onClick={guided ? presentation.commands.confirm : handleApply} disabled={guided ? !presentation.commands.confirm : !canApply || busy}
            className="rounded bg-[var(--gw-primary)] px-3 py-1 text-xs text-[var(--gw-on-primary)] disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]"
            data-testid="group-schema-apply"
          >
            {applyMutation.isPending ? t('step4.group.schema.applying') : t('step4.group.schema.confirm')}
          </button>
        )}
      </div>}

      {preview && <RecordingSchemaImpact preview={preview} layout={group?.layout_version} />}

      {shouldOfferCheck && (
        <div className="flex flex-wrap items-center gap-2" role="status" data-testid="group-schema-operation-pending">
          <span className="text-xs text-[var(--gw-warning-fg)]">{t('step4.group.schema.check_required')}</span>
          {!suppressActions && <button type="button" onClick={guided ? presentation.commands.check : handleCheck} disabled={guided ? !presentation.commands.check : busy} className="text-xs underline disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]" data-testid="group-schema-check">
            {t('step4.group.schema.check')}
          </button>}
        </div>
      )}

      {effectiveOperation && (
        <div className="space-y-1 text-xs text-[var(--gw-text)]" role="status" data-testid="group-schema-operation">
          <details><summary>{t('step4.basic.diagnostics')}</summary><p>{t('step4.group.schema.operation')}: <span className="font-mono">{effectiveOperation.operation_id}</span></p></details>
          <p>{t('step4.group.schema.status')}: <span className="font-mono">{operationStatus(effectiveOperation, t)}</span></p>
          {detail && !blocker && <p>{t('step4.group.schema.reason')}: {detail}</p>}
        </div>
      )}

      {error !== null && !blocker && (
        <p role="alert" className="text-xs text-[var(--gw-danger-fg)]" data-testid="group-schema-error">
          {getSafeErrorMessage(error, t).message}
          {needsCheck && <span className="ml-1">{t('step4.group.schema.error_check_same_operation')}</span>}
        </p>
      )}
      {guided?.renderActions?.(presentation)}
    </section>
  );
};
