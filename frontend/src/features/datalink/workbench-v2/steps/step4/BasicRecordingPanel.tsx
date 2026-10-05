import { destinationIsSaved, optionalShareRequest, groupMatchesRequest, operationScopeMatchesCurrent, requestMatchesCurrentControls, groupSnapshotMatchesRequest, operationBelongsToGroupId, unresolvedRequestMatchesCurrentScope } from './basicRecordingPanelHelpers';
import { BasicSchemaPreparation } from './BasicSchemaPreparation';
import { GroupDeliveryRecovery } from './writeGroup/GroupDeliveryRecovery';
import { supportsWriteGroupKind } from './writeGroupCapabilities';
import * as React from 'react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import {
  useEnsureBasicManagedMutation,
  useRecordingStartMutation,
  useRecordingStartOperationQuery,
} from '../../../../../hooks/datalink/useStudioV2RecordingStart';
import { useWriteGroupDeliveryQuery, useWriteGroupsQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { RecordingStartOperation, RecordingStartRequest } from '../../../../../types/studioV2RecordingStart';
import type { WriteGroup, WriteGroupDraft } from '../../../../../types/studioV2WriteGroup';
import type { ModbusShareStatus } from '../../../../../types/modbusShare';
import { isModbusShareConfiguredEnabled } from '../../../../../types/modbusShare';
import { studioV2WorkspaceKeys } from '../../../../../hooks/datalink/keys';
import type { WorkbenchV2State } from '../../state/types';
import { buildGroupCandidates } from '../../state/writeGroup/candidates';
import {
  safeBasicNextAction,
  safeBasicOperationMessage,
  safeBasicStartError,
} from '../../state/basicRecordingMessages';
import { BasicRecordingEvidence } from './BasicRecordingEvidence';
import { BasicRecordingMembers } from './BasicRecordingMembers';
import { BasicRecordingProgress } from './BasicRecordingProgress';
import {
  basicManagedGroupDraft,
  basicRecordingScopeKey,
  buildBasicManagedDrafts,
  createBasicRecordingRequestId,
  isBasicManagedWriteGroup,
  latestRecordingStartOperation,
  loadBasicRecordingIntent,
  saveBasicRecordingIntent,
  type BasicRecordingIntentRecord,
} from '../../state/basicRecording';

export interface BasicRecordingPanelProps {
  state: WorkbenchV2State;
  workspaceId?: string;
  readonly: boolean;
  shareStatus?: ModbusShareStatus | null;
  role?: string;
  onCommit?: (confirmedDeviceIds?: string[]) => void;
  onOpenAdvanced?: () => void;
  onOperationActiveChange?: (active: boolean) => void;
}

const DEFAULT_INTERVAL_SECONDS = 60;
const BASIC_ROLE = 'managed-recording';

/**
 * Basic owns one device and one managed group, including safe schema preparation.
 * Custom columns and other group variants remain in the advanced editor.
 */
export const BasicRecordingPanel: React.FC<BasicRecordingPanelProps> = ({
  state,
  workspaceId,
  readonly,
  shareStatus,
  role = BASIC_ROLE,
  onCommit,
  onOpenAdvanced,
  onOperationActiveChange,
}) => {
  const { t } = useTranslation('workbench-v2');
  const queryClient = useQueryClient();
  const groupsQuery = useWriteGroupsQuery(Boolean(workspaceId));
  const ensureMutation = useEnsureBasicManagedMutation();
  const startMutation = useRecordingStartMutation();
  const { candidates } = useMemo(() => buildGroupCandidates(state), [state]);
  const destination = useMemo(() => ({
    connector_id: state.db.connector.connector_id ?? '',
    connector_revision: state.db.connector.identity_revision ?? '',
    table_schema: state.db.connector.schema ?? '',
    table_name: '',
    storage_strategy: 'managed' as const,
  }), [state.db.connector.connector_id, state.db.connector.identity_revision, state.db.connector.schema]);
  const draftByDevice = useMemo(
    () => buildBasicManagedDrafts(workspaceId ?? '', candidates, destination, role),
    [candidates, destination, role, workspaceId],
  );
  const shareOnly = isModbusShareConfiguredEnabled(shareStatus) && !destinationIsSaved(state);
  const deviceIds = useMemo(
    () => shareOnly
      ? state.devices.filter((device) => device.persisted === true && device.save_state === 'saved').map((device) => device.id)
      : draftByDevice.map((item) => item.device_id),
    [draftByDevice, shareOnly, state.devices],
  );
  const [selectedDeviceId, setSelectedDeviceId] = useState('');
  const [intervalSeconds, setIntervalSeconds] = useState(DEFAULT_INTERVAL_SECONDS);
  const [incompletePolicy, setIncompletePolicy] = useState('skip_row');
  const [localIntent, setLocalIntent] = useState<BasicRecordingIntentRecord | undefined>();
  const [submittedOperation, setSubmittedOperation] = useState<RecordingStartOperation | undefined>();
  const [actionError, setActionError] = useState<unknown>(null);
  const [actionBusy, setActionBusy] = useState(false);
  const actionBusyRef = useRef(false);
  const activeSelectionRef = useRef('');
  const refreshedOperationRef = useRef('');

  useEffect(() => {
    if (!deviceIds.includes(selectedDeviceId)) setSelectedDeviceId(deviceIds[0] ?? '');
  }, [deviceIds, selectedDeviceId]);

  const selectedDraft = draftByDevice.find((item) => item.device_id === selectedDeviceId)?.draft;
  const selectedCandidates = candidates.filter((candidate) => candidate.device_id === selectedDeviceId);
  const groups = groupsQuery.data?.groups ?? [];
  const liveGroups = groups.filter((group) => group.status !== 'deleted');
  const selectedGroup = liveGroups.find((group) => isBasicManagedWriteGroup(group, selectedDeviceId));
  const basicSelectedGroup = shareOnly ? undefined : selectedGroup;
  const selectedGroupId = basicSelectedGroup?.id;
  const selectedGroupRevision = basicSelectedGroup?.revision;
  const selectedGroupAppliedRevision = basicSelectedGroup?.applied_revision;
  const connectorRevision = state.db.connector.identity_revision;
  const selectedGroupInterval = basicSelectedGroup?.row_policy.interval_seconds;
  const selectedGroupIncompletePolicy = basicSelectedGroup?.row_policy.incomplete_policy;
  const selectedGroupMembersMatch = !basicSelectedGroup || basicSelectedGroup.members.length === selectedCandidates.length &&
    basicSelectedGroup.members.every((member) => selectedCandidates.some((candidate) => candidate.point_id === member.point_id && candidate.tag_id === member.tag_id));
  const membersMismatch = Boolean(basicSelectedGroup && !selectedGroupMembersMatch);
  const scopeKey = workspaceId && selectedDeviceId
    ? basicRecordingScopeKey(workspaceId, selectedDeviceId, role)
    : '';

  useEffect(() => {
    actionBusyRef.current = false;
    setActionBusy(false);
    activeSelectionRef.current = scopeKey;
    setSubmittedOperation(undefined);
    setActionError(null);
    refreshedOperationRef.current = '';
    return () => {
      if (activeSelectionRef.current === scopeKey) {
        activeSelectionRef.current = '';
        actionBusyRef.current = false;
      }
    };
  }, [scopeKey]);

  useEffect(() => {
    const stored = scopeKey ? loadBasicRecordingIntent(scopeKey) : undefined;
    const keepTerminalEvidence = stored?.operation_id && operationBelongsToGroupId(stored, selectedGroupId);
    const keepUnresolvedRequest = unresolvedRequestMatchesCurrentScope(stored, selectedGroupId, connectorRevision);
    setLocalIntent(groupSnapshotMatchesRequest(stored, selectedGroupId, selectedGroupRevision, connectorRevision) || keepTerminalEvidence || keepUnresolvedRequest ? stored : undefined);
  }, [scopeKey, selectedGroupId, selectedGroupRevision, selectedGroupAppliedRevision, connectorRevision]);

  useEffect(() => {
    setIntervalSeconds(selectedGroupInterval ?? DEFAULT_INTERVAL_SECONDS);
    setIncompletePolicy(selectedGroupIncompletePolicy === 'partial' ? 'partial' : 'skip_row');
  }, [selectedGroup?.id, selectedGroup?.revision, selectedGroupIncompletePolicy, selectedGroupInterval]);

  const destinationMismatch = Boolean(!shareOnly && selectedGroup && (
    selectedGroup.destination.connector_id !== state.db.connector.connector_id ||
    selectedGroup.destination.connector_revision !== state.db.connector.identity_revision
  ));
  const operationId = localIntent?.operation_id;
  const operationQuery = useRecordingStartOperationQuery(operationId);
  const deliveryQuery = useWriteGroupDeliveryQuery(selectedGroup?.id, !shareOnly && Boolean(selectedGroup) && !destinationMismatch, 5000);
  const currentRevisionStages = selectedGroup && deliveryQuery.data?.group_id === selectedGroup.id &&
    deliveryQuery.data.revision_stages?.group_revision === selectedGroup.applied_revision
    ? deliveryQuery.data.revision_stages.stages
    : undefined;
  const lastCommittedEffect = deliveryQuery.data?.last_sql_committed_effect;
  const currentCommittedEffect = selectedGroup && deliveryQuery.data?.group_id === selectedGroup.id &&
    selectedGroup.applied_revision && lastCommittedEffect &&
    lastCommittedEffect.group_revision === selectedGroup.applied_revision &&
    lastCommittedEffect.connector_revision === selectedGroup.destination.connector_revision
    ? lastCommittedEffect
    : undefined;
  const workspaceSnapshotReady = groupsQuery.data?.workspace_id === workspaceId;
  const workspaceRevision = workspaceSnapshotReady ? groupsQuery.data?.workspace_revision ?? '' : '';
  const operation = latestRecordingStartOperation(operationId, operationQuery.data, submittedOperation);
  const intentMatchesCurrent = groupMatchesRequest(localIntent, basicSelectedGroup, state.db.connector.identity_revision, shareStatus) ||
    operationScopeMatchesCurrent(localIntent, operation, basicSelectedGroup, workspaceId, selectedDeviceId, state.db.connector.identity_revision);
  const requestCanBeReused = requestMatchesCurrentControls(
    localIntent, operation, basicSelectedGroup, workspaceId, selectedDeviceId, state.db.connector.identity_revision,
    intervalSeconds, incompletePolicy, shareStatus, workspaceRevision,
  );
  const destinationSaved = destinationIsSaved(state);
  const advancedGroups = liveGroups.filter((group) => !isBasicManagedWriteGroup(group, selectedDeviceId));
  const operationActive = operation?.status === 'pending' || operation?.status === 'running' || actionBusy;

  useEffect(() => {
    onOperationActiveChange?.(operationActive);
    return () => onOperationActiveChange?.(false);
  }, [onOperationActiveChange, operationActive]);

  useEffect(() => {
    if (operation?.status !== 'succeeded' || refreshedOperationRef.current === operation.operation_id) return;
    refreshedOperationRef.current = operation.operation_id;
    void (async () => {
      await groupsQuery.refetch();
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroups() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseMetadataAll() });
      for (const progress of operation.groups) {
        await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroupReadiness(progress.group_id) });
      }
    })();
  }, [groupsQuery, operation, queryClient]);

  const blockReason = !workspaceId ? 'workspace'
    : readonly ? 'readonly'
      : !shareOnly && candidates.length === 0 ? 'mapping'
        : !shareOnly && !supportsWriteGroupKind(state.db.connector.kind) ? 'destination_kind'
          : !shareOnly && !destinationSaved ? 'destination'
          : !workspaceSnapshotReady || (!shareOnly && !workspaceRevision) ? 'workspace_revision'
            : destinationMismatch ? 'destination_mismatch'
              : membersMismatch ? 'members_changed' : undefined;

  const buildStartRequest = (group: WriteGroup | undefined, expectedWorkspaceRevision: string, requestId: string, draft?: WriteGroupDraft): RecordingStartRequest => ({
    request_id: requestId,
    workspace_id: workspaceId ?? group?.workspace_id ?? '',
    expected_workspace_revision: expectedWorkspaceRevision,
    device_ids: [selectedDeviceId],
    groups: group ? [{
      group_id: group.id,
      expected_group_revision: group.revision,
      expected_connector_revision: state.db.connector.identity_revision ?? group.destination.connector_revision,
      ...(draft ? { draft } : {}),
    }] : [],
    ...optionalShareRequest(shareStatus),
  });

  const ensureGroup = async () => {
    if (basicSelectedGroup) return { group: basicSelectedGroup, workspace_revision: workspaceRevision };
    if (!selectedDraft || !workspaceId || blockReason) throw new Error('basic preparation unavailable');
    return ensureMutation.mutateAsync({ deviceId: selectedDeviceId, request: {
      workspace_id: workspaceId, expected_workspace_revision: workspaceRevision,
      expected_connector_revision: state.db.connector.identity_revision ?? '',
      group: { ...selectedDraft, row_policy: { ...selectedDraft.row_policy, interval_seconds: intervalSeconds, incomplete_policy: incompletePolicy } },
    } });
  };

  const handleStart = async () => {
    if (blockReason || operationActive || actionBusyRef.current || (!shareOnly && !selectedDraft) || !scopeKey) return;
    const scopeAtRequest = scopeKey;
    actionBusyRef.current = true;
    setActionBusy(true);
    setActionError(null);
    try {
      let group = basicSelectedGroup;
      let expectedWorkspaceRevision = workspaceRevision;
      if (!group && !shareOnly) {
        if (!selectedDraft) return;
        const ensured = await ensureGroup();
        if (activeSelectionRef.current !== scopeAtRequest) return;
        group = ensured.group;
        expectedWorkspaceRevision = ensured.workspace_revision;
      }
      if (activeSelectionRef.current !== scopeAtRequest) return;
      const changed = group ? group.row_policy.interval_seconds !== intervalSeconds ||
        (group.row_policy.incomplete_policy === 'partial') !== (incompletePolicy === 'partial') : false;
      const draft = changed && group ? basicManagedGroupDraft(group, intervalSeconds, incompletePolicy) : undefined;
      const request = requestCanBeReused && localIntent
        ? localIntent.request
        : buildStartRequest(group, expectedWorkspaceRevision, createBasicRecordingRequestId(), draft);
      const record: BasicRecordingIntentRecord = { scope_key: scopeAtRequest, request };
      saveBasicRecordingIntent(record);
      setLocalIntent(record);
      const result = await startMutation.mutateAsync(request);
      saveBasicRecordingIntent({ ...record, operation_id: result.operation_id });
      if (activeSelectionRef.current !== scopeAtRequest) return;
      setLocalIntent({ ...record, operation_id: result.operation_id });
      setSubmittedOperation(result);
    } catch (error) {
      if (activeSelectionRef.current === scopeAtRequest) setActionError(error);
    } finally {
      if (activeSelectionRef.current === scopeAtRequest) {
        actionBusyRef.current = false;
        setActionBusy(false);
      }
    }
  };

  const handleRetry = () => {
    const request = localIntent?.request;
    if (!request || !intentMatchesCurrent || blockReason || actionBusy || actionBusyRef.current || activeSelectionRef.current !== scopeKey) return;
    actionBusyRef.current = true;
    setActionBusy(true);
    setActionError(null);
    void startMutation.mutateAsync(request).then((result) => {
      if (activeSelectionRef.current !== scopeKey) return;
      const record = { scope_key: scopeKey, request, operation_id: result.operation_id };
      saveBasicRecordingIntent(record);
      setLocalIntent(record);
      setSubmittedOperation(result);
    }).catch((error) => {
      if (activeSelectionRef.current === scopeKey) setActionError(error);
    }).finally(() => {
      if (activeSelectionRef.current === scopeKey) {
        actionBusyRef.current = false;
        setActionBusy(false);
      }
    });
  };

  const statusReason = safeBasicOperationMessage(operation, t);
  const safeActionError = actionError ? safeBasicStartError(actionError, t) : undefined;
  const safeNextAction = safeBasicNextAction(operation, t);
  const activatedDeviceIds = operation?.status === 'succeeded'
    ? operation.devices.filter((device) => device.activated).map((device) => device.device_id)
    : [];

  return (
    <section className="space-y-4 rounded-xl border border-cyan-900/70 bg-slate-950/40 p-5" aria-labelledby="basic-recording-title" data-testid="basic-recording-panel">
      <div>
        <h3 id="basic-recording-title" className="text-sm font-semibold text-slate-100">{t('step4.basic.title')}</h3>
        <p className="mt-1 text-xs leading-5 text-slate-400">{t('step4.basic.subtitle')}</p>
      </div>

      {deviceIds.length > 0 && (
        <label className="block text-xs text-slate-300">
          <span>{t('step4.basic.device')}</span>
          <select value={selectedDeviceId} onChange={(event) => setSelectedDeviceId(event.target.value)} disabled={operationActive} className="mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-2 text-sm text-slate-100" data-testid="basic-recording-device">
            {deviceIds.map((deviceId) => <option key={deviceId} value={deviceId}>{state.devices.find((device) => device.id === deviceId)?.name || deviceId}</option>)}
          </select>
        </label>
      )}

      <div className="grid gap-3 sm:grid-cols-2">
        <label className="block text-xs text-slate-300">
          <span>{t('step4.basic.interval')}</span>
          <input type="number" min={1} value={intervalSeconds} disabled={operationActive || readonly} onChange={(event) => setIntervalSeconds(Math.max(1, Number(event.target.value) || 1))} className="mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-2 text-sm text-slate-100" data-testid="basic-recording-interval" />
        </label>
        <label className="block text-xs text-slate-300">
          <span>{t('step4.basic.completeness')}</span>
          <select value={incompletePolicy} disabled={operationActive || readonly} onChange={(event) => setIncompletePolicy(event.target.value)} className="mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-2 text-sm text-slate-100" data-testid="basic-recording-completeness">
            <option value="skip_row">{t('step4.basic.skip_row')}</option>
            <option value="partial">{t('step4.basic.partial_row')}</option>
          </select>
        </label>
      </div>

      <p className="text-xs text-slate-400" data-testid="basic-recording-destination">
        {shareOnly
          ? t('step4.basic.share_only_destination')
          : t('step4.basic.destination', { name: state.db.connector.name || state.db.connector.kind })}
      </p>

      {blockReason && <p role="status" className="text-xs text-amber-300" data-testid="basic-recording-blocked">{t(`step4.basic.blocked.${blockReason}`)}</p>}
      {destinationMismatch && <p role="alert" className="text-xs text-amber-300" data-testid="basic-recording-destination-mismatch">{t('step4.basic.blocked.destination_mismatch')}</p>}
      {advancedGroups.length > 0 && <p className="text-xs text-slate-500" data-testid="basic-recording-advanced-note">{t('step4.basic.advanced_note', { count: advancedGroups.length })}</p>}
      {!shareOnly && (selectedGroup?.members.length ?? selectedCandidates.length) > 0 && (
        <>
          <BasicRecordingMembers candidates={selectedCandidates} group={selectedGroup} />
          <p className="text-xs text-slate-500">{t('step4.basic.members', { count: selectedGroup?.members.length ?? selectedCandidates.length })}</p>
        </>
      )}
      {membersMismatch && onOpenAdvanced && (
        <button type="button" onClick={onOpenAdvanced} className="rounded border border-amber-700 px-3 py-2 text-xs text-amber-200" data-testid="basic-recording-members-review">
          {t('step4.basic.open_advanced')}
        </button>
      )}

      {!shareOnly && (
        <BasicRecordingEvidence
          intervalSeconds={intervalSeconds}
          deliveryError={deliveryQuery.isError}
          stages={currentRevisionStages}
          committedEffect={currentCommittedEffect}
          group={selectedGroup}
        />
      )}

      {!shareOnly && <BasicSchemaPreparation key={`${scopeKey}|${connectorRevision}|${selectedGroupMembersMatch}`} group={basicSelectedGroup} workspaceId={workspaceId ?? ''} workspaceRevision={workspaceRevision} disabled={Boolean(blockReason) || operationActive} readonly={readonly || operationActive} destinationMatches={!destinationMismatch} ensureGroup={ensureGroup} onApplied={() => void groupsQuery.refetch()} />}
      {!shareOnly && selectedGroup && !deliveryQuery.isError && deliveryQuery.data?.group_id === selectedGroup.id && <GroupDeliveryRecovery groupId={selectedGroup.id} items={deliveryQuery.data.attention ?? []} readonly={readonly || operationActive} onResolved={() => void deliveryQuery.refetch()} />}

      <div className="flex flex-wrap items-center gap-3">
        <button type="button" disabled={Boolean(blockReason) || operationActive || groupsQuery.isLoading} onClick={() => void handleStart()} className="rounded bg-cyan-600 px-4 py-2 text-xs font-semibold text-white disabled:cursor-not-allowed disabled:opacity-40" data-testid="basic-recording-start">
          {operationActive ? t('step4.basic.starting') : t('step4.basic.start')}
        </button>
        {operation?.status && ['pending', 'running', 'partial', 'failed', 'unknown'].includes(operation.status) && <button type="button" disabled={actionBusy || !intentMatchesCurrent || Boolean(blockReason)} onClick={handleRetry} className="rounded border border-slate-600 px-3 py-2 text-xs text-slate-200 disabled:opacity-40" data-testid="basic-recording-retry">{t('step4.basic.retry')}</button>}
        {operationId && operationQuery.isError && <button type="button" onClick={() => void operationQuery.refetch()} className="rounded border border-slate-600 px-3 py-2 text-xs text-slate-200" data-testid="basic-recording-check">{t('step4.basic.check')}</button>}
      </div>

      {safeActionError && <p role="alert" className="text-xs text-red-300" data-testid="basic-recording-error">{safeActionError}</p>}
      {operation && (
        <BasicRecordingProgress
          groupNames={Object.fromEntries(groups.map((group) => [group.id, group.name]))}
          deviceNames={Object.fromEntries(state.devices.map((device) => [device.id, device.name]))}
          operation={operation}
          intentMatchesCurrent={intentMatchesCurrent}
          activatedDeviceIds={activatedDeviceIds}
          statusReason={statusReason}
          nextAction={safeNextAction}
          onCommit={onCommit}
          onOpenAdvanced={onOpenAdvanced}
        />
      )}
      {operationQuery.isError && operationId && <p role="alert" className="text-xs text-amber-300" data-testid="basic-recording-query-error">{t('step4.basic.query_failed')}</p>}
    </section>
  );
};
