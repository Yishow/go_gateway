import { ParameterHelp } from '../../components/ParameterHelp';
import { LegacyRecordingReadError } from './LegacyRecordingReadError';
import { RecordingDraftBoundary, useRecordingNavigation } from '../../recording-editor/RecordingDraftNavigation';
import { useBasicLegacyDraftNavigation, useShareOnlyPresentation } from './useLegacyRecordingDraftNavigation';
import { safeRuntimeLabel } from '../../../../../utils/runtimeOperatorSafety';
import type { RecordingViewedDevice, RenderRecordingScopeSummary } from './recordingScopeModel';
import { RemovedRecordingManagement, SavedRecordingManagement } from './SavedRecordingManagement';
import type { SavedRecordingActionHost } from './writeGroup/lifecycle/SavedRecordingActions';
import { destinationIsSaved, groupMatchesRequest, operationScopeMatchesCurrent, requestMatchesCurrentControls, groupSnapshotMatchesRequest, operationBelongsToGroupId, unresolvedRequestMatchesCurrentScope } from './basicRecordingPanelHelpers';
import { BasicMemberUpdate } from './BasicMemberUpdate';
import { buildBasicStartRequest } from './basicStartRequest';
import { LegacyMigrationEntry } from './LegacyMigrationEntry';
import { BasicRecordingPreparation } from './BasicRecordingPreparation';
import { BasicRecordingRetry } from './BasicRecordingRetry';
import { useBasicRecordingIntentObservation } from './useBasicRecordingIntentObservation';
import { useRecordingWorkspaceGuard } from './useRecordingWorkspaceGuard';
import { isOriginalRequestRetry, recordingStartRecoveryBlocked } from './recordingStartRecovery';
import type { RecordingPreparationCommandFence } from './recordingPreparationPresentation';
import { GroupDeliveryRecovery } from './writeGroup/GroupDeliveryRecovery';
import { supportsWriteGroupKind } from './writeGroupCapabilities';
import * as React from 'react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useEnsureBasicManagedMutation, useRecordingStartMutation, useRecordingStartOperationQuery } from '../../../../../hooks/datalink/useStudioV2RecordingStart';
import * as writeGroupHooks from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import { useWriteGroupDeliveryQuery, useWriteGroupsQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { RecordingStartOperation } from '../../../../../types/studioV2RecordingStart';
import type { ModbusShareStatus } from '../../../../../types/modbusShare';
import { isModbusShareConfiguredEnabled } from '../../../../../types/modbusShare';
import { studioV2WorkspaceKeys } from '../../../../../hooks/datalink/keys';
import type { WorkbenchV2State } from '../../state/types';
import { buildGroupCandidates } from '../../state/writeGroup/candidates';
import { safeBasicNextAction, safeBasicOperationMessage, safeBasicStartError } from '../../state/basicRecordingMessages';
import { BasicRecordingEvidence } from './BasicRecordingEvidence';
import { BasicRecordingMembers } from './BasicRecordingMembers';
import { recordingCompletion } from './recordingCompletion';
import { BasicRecordingProgress } from './BasicRecordingProgress';
import { CircuitProfileBasicPanel } from './CircuitProfileBasicPanel';
import { DeviceScopeTiles } from './DeviceScopeTiles';
import { loadSelectedDeviceId, saveSelectedDeviceId } from '../../state/deviceScopeSelection';
import { groupCircuitId, managedProfileDeviceId } from '../../state/circuitProfileBasic';
import { saveBasicLegacyDraft } from '../../state/basicLegacyDraft';
import { basicManagedGroupDraft, completeBasicRecordingIntent, basicRecordingScopeKey, buildBasicManagedDrafts, createBasicRecordingRequestId, isBasicManagedWriteGroup, latestRecordingStartOperation, loadBasicRecordingIntent, saveBasicRecordingIntent, type BasicRecordingIntentRecord } from '../../state/basicRecording';

export interface BasicRecordingPanelProps {
  state: WorkbenchV2State;
  workspaceId?: string;
  readonly: boolean;
  lifecycleHost?: SavedRecordingActionHost;
  shareStatus?: ModbusShareStatus | null;
  role?: string;
  onCommit?: (confirmedDeviceIds?: string[]) => void;
  onOpenAdvanced?: (groupId?: string) => void;
  onOpenMemberSettings?: (groupId: string) => void;
  onOperationActiveChange?: (active: boolean) => void;
  onViewedDeviceChange?: (scope: RecordingViewedDevice) => void;
  renderScopeSummary?: RenderRecordingScopeSummary;
  /**
   * Shows a managed profile group's own device and circuit here for review. It belongs to one workspace, is shown once,
   * lapses when the operator chooses a device or circuit themselves, and never saves, applies or starts anything.
   */
  reviewRequest?: { workspaceId: string; groupId: string };
}

const DEFAULT_INTERVAL_SECONDS = 60;
const BASIC_ROLE = 'managed-recording';

/**
 * Basic owns one device and one managed group, including safe schema preparation.
 * Custom columns and other group variants remain in the advanced editor.
 */
export const BasicRecordingPanel: React.FC<BasicRecordingPanelProps> = props => <RecordingDraftBoundary><BasicRecordingPanelContent {...props} /></RecordingDraftBoundary>;
const BasicRecordingPanelContent: React.FC<BasicRecordingPanelProps> = ({
  state,
  workspaceId,
  readonly: operatorReadonly,
  lifecycleHost,
  shareStatus,
  role = BASIC_ROLE,
  onCommit,
  onOpenAdvanced,
  onOpenMemberSettings,
  onOperationActiveChange,
  reviewRequest, onViewedDeviceChange, renderScopeSummary,
}) => {
  const { t } = useTranslation(['workbench-v2', 'studio-review']);
  const { request } = useRecordingNavigation();
  const queryClient = useQueryClient();
  const groupsQuery = useWriteGroupsQuery(Boolean(workspaceId));
  const readonly = operatorReadonly || groupsQuery.isError;
  const readError = groupsQuery.isError && <LegacyRecordingReadError retry={groupsQuery.refetch} />;
  const workspaceBlocked = useRecordingWorkspaceGuard(workspaceId, groupsQuery.data);
  const ensureMutation = useEnsureBasicManagedMutation();
  const startMutation = useRecordingStartMutation();
  const { candidates, excluded } = useMemo(() => buildGroupCandidates(state), [state]);
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
  const shareOnly = useShareOnlyPresentation(workspaceId ?? '', Boolean(state.db.connector.persisted && state.db.connector.connector_id && state.db.connector.identity_revision),
    Boolean(groupsQuery.data) && !groupsQuery.isLoading && !groupsQuery.isError, isModbusShareConfiguredEnabled(shareStatus));
  // Every persisted device is a legal view, whether or not it has sources to record
  // yet; a device without sources shows its own blocked state instead of being
  // replaced by another one. The Circuit subtree lists the same devices.
  const deviceIds = useMemo(() => {
    if (shareOnly) return state.devices.filter((device) => device.persisted === true && device.save_state === 'saved').map((device) => device.id);
    const recordable = new Set(draftByDevice.map((item) => item.device_id));
    return state.devices
      .filter((device) => recordable.has(device.id) || Boolean(device.persisted && device.save_state !== 'saving' && device.save_state !== 'save-error'))
      .map((device) => device.id);
  }, [draftByDevice, shareOnly, state.devices]);
  // The viewed device is derived, never copied once: an explicit choice of this
  // workspace wins while it is still legal, otherwise the saved selection of the
  // current workspace, otherwise the first device. A workspace or device list that
  // arrives after the first render therefore restores the saved selection too.
  const [choice, setChoice] = useState<{ workspaceId: string; deviceId: string }>();
  const [legacyView, setLegacyView] = useState<{ scope: string; mode: 'edit' | 'prepare' }>();
  const restoredDeviceId = useMemo(() => loadSelectedDeviceId(workspaceId ?? '', deviceIds), [workspaceId, deviceIds]);
  const selectedDeviceId = choice && choice.workspaceId === (workspaceId ?? '') && deviceIds.includes(choice.deviceId)
    ? choice.deviceId
    : restoredDeviceId ?? deviceIds[0] ?? '';
  useEffect(() => {
    onViewedDeviceChange?.({ workspaceId: workspaceId ?? '', deviceId: selectedDeviceId });
  }, [onViewedDeviceChange, workspaceId, selectedDeviceId]);
  const selectDevice = useCallback((deviceId: string) => {
    saveSelectedDeviceId(workspaceId ?? '', deviceId);
    setChoice({ workspaceId: workspaceId ?? '', deviceId });
  }, [workspaceId]);
  // A routed group is shown once, on its own device and circuit, the same way choosing them by hand would.
  const [circuitRequest, setCircuitRequest] = useState<{ workspaceId: string; deviceId: string; circuitId: string }>();
  const consumedReview = useRef<typeof reviewRequest>(undefined);
  const reviewGroups = groupsQuery.data?.groups;
  const clearCircuitRequest = useCallback(() => setCircuitRequest(undefined), []);
  const selectByOperator = useCallback((deviceId: string) => request(() => {
    consumedReview.current = reviewRequest;
    setCircuitRequest(undefined);
    setLegacyView(undefined);
    selectDevice(deviceId);
  }), [request, reviewRequest, selectDevice]);
  useEffect(() => {
    if (!reviewRequest || reviewRequest === consumedReview.current || reviewRequest.workspaceId !== (workspaceId ?? '')) return;
    const group = reviewGroups?.find((candidate) => candidate.id === reviewRequest.groupId && candidate.status !== 'deleted');
    const reviewDeviceId = group ? managedProfileDeviceId(group) : undefined;
    if (!group || !reviewDeviceId || !deviceIds.includes(reviewDeviceId)) return;
    consumedReview.current = reviewRequest;
    request(() => { setLegacyView(undefined); selectDevice(reviewDeviceId);
      setCircuitRequest({ workspaceId: reviewRequest.workspaceId, deviceId: reviewDeviceId, circuitId: groupCircuitId(group) }); });
  }, [reviewRequest, reviewGroups, deviceIds, selectDevice, workspaceId, request]);
  const [intervalSeconds, setIntervalSeconds] = useState(DEFAULT_INTERVAL_SECONDS);
  const [incompletePolicy, setIncompletePolicy] = useState('skip_row');
  const [localIntent, setLocalIntent] = useState<BasicRecordingIntentRecord | undefined>();
  const [submittedOperation, setSubmittedOperation] = useState<RecordingStartOperation | undefined>();
  const [actionError, setActionError] = useState<unknown>(null);
  const [actionBusy, setActionBusy] = useState(false);
  const [preparationFence, setPreparationFence] = useState<RecordingPreparationCommandFence | null>(null);
  const preparationScope = useRef('');
  const preparationFenceRef = useRef(preparationFence); preparationFenceRef.current = preparationFence;
  const isPreparationBlocked = useCallback((groupId?: string) => {
    const fence = preparationFenceRef.current;
    return !fence || fence.scopeKey !== preparationScope.current || fence.isBlocked(groupId);
  }, []);
  const receivePreparationFence = useCallback((fence: RecordingPreparationCommandFence | null) => {
    if (!fence || fence.scopeKey === preparationScope.current) setPreparationFence(fence);
  }, []);
  const actionBusyRef = useRef(false);
  const activeSelectionRef = useRef('');
  const refreshedOperationRef = useRef('');

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
  useBasicRecordingIntentObservation(scopeKey, localIntent, setLocalIntent);
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

  const supportsCircuitProfile = 'useEnsureCircuitProfileMutation' in writeGroupHooks;
  const legacyGroup = liveGroups.find((g) => g.basic_managed_device_id === selectedDeviceId && (g.layout_version ?? 'managed_wide_v1') === 'managed_wide_v1');
  const hasLegacyBasicWide = Boolean(legacyGroup);
  const legacyScope = JSON.stringify([workspaceId, selectedDeviceId]);
  const legacyMode = legacyView?.scope === legacyScope ? legacyView.mode : undefined;
  useBasicLegacyDraftNavigation({ scopeKey, group: basicSelectedGroup, visible: !supportsCircuitProfile || shareOnly || hasLegacyBasicWide && legacyMode === 'edit',
    pending: operationActive || Boolean(localIntent && (!operation || ['pending', 'running', 'partial', 'unknown'].includes(operation.status))),
    interval: intervalSeconds, policy: incompletePolicy, setInterval: setIntervalSeconds, setPolicy: setIncompletePolicy });
  if (supportsCircuitProfile && !shareOnly) {
    if (groupsQuery.isLoading) return <div data-testid="circuit-profile-groups-loading" className="rounded border border-[var(--gw-control)] p-4 text-xs text-[var(--gw-muted)]">{t('step4.basic.circuit.groups_loading')}</div>;
    if (groupsQuery.isError && !groupsQuery.data) return readError;
    if (!hasLegacyBasicWide || legacyMode !== 'edit') return <>{readError}<CircuitProfileBasicPanel renderScopeSummary={renderScopeSummary} lifecycleHost={lifecycleHost} key={JSON.stringify([workspaceId, selectedDeviceId, legacyMode])} prepareOnly={legacyMode === 'prepare'} onClose={() => request(() => setLegacyView({ scope: legacyScope, mode: 'edit' }))} legacyGroup={legacyGroup} onEditLegacy={() => request(() => setLegacyView({ scope: legacyScope, mode: 'edit' }))} deviceId={selectedDeviceId} onSelectDevice={selectByOperator} state={state} workspaceId={workspaceId} readonly={readonly} workspaceRevision={workspaceRevision} readinessToken={shareStatus?.readiness_token} settingsRevision={shareStatus?.settings_revision} shareWorkspaceRevision={shareStatus?.workspace_revision} onOpenAdvanced={onOpenAdvanced} onContinueToRuntime={(ids) => onCommit?.(ids)} onOperationActiveChange={onOperationActiveChange} circuitRequest={circuitRequest} onCircuitRequestDone={clearCircuitRequest} /><RemovedRecordingManagement groups={groups} deviceId={selectedDeviceId} workspaceId={workspaceId} workspaceRevision={workspaceRevision} readonly={readonly} lifecycleHost={lifecycleHost} onOpenAdvanced={onOpenAdvanced} /></>;
  }

  const blockReason = !workspaceId ? 'workspace'
    : readonly ? 'readonly'
      : !shareOnly && candidates.length === 0 ? 'mapping'
        : !shareOnly && !supportsWriteGroupKind(state.db.connector.kind) ? 'destination_kind'
          : !shareOnly && !destinationSaved ? 'destination'
          : !workspaceSnapshotReady || (!shareOnly && !workspaceRevision) ? 'workspace_revision'
            : destinationMismatch ? 'destination_mismatch'
              : membersMismatch ? 'members_changed' : undefined;

  const preparationKey = JSON.stringify([scopeKey, workspaceRevision, basicSelectedGroup?.id, basicSelectedGroup?.revision,
    basicSelectedGroup?.applied_revision, destination.connector_id, connectorRevision, intervalSeconds, incompletePolicy,
    selectedCandidates.map(candidate => [candidate.device_id, candidate.point_id, candidate.tag_id]), readonly, blockReason]);
  preparationScope.current = preparationKey;

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
        : buildBasicStartRequest({ workspaceId, deviceId: selectedDeviceId, connectorRevision: state.db.connector.identity_revision, shareStatus }, group, expectedWorkspaceRevision, createBasicRecordingRequestId(), draft);
      const record: BasicRecordingIntentRecord = { scope_key: scopeAtRequest, request };
      saveBasicRecordingIntent(record);
      setLocalIntent(record);
      const result = await startMutation.mutateAsync(request);
      if (completeBasicRecordingIntent(record, result.operation_id) === 'superseded' || activeSelectionRef.current !== scopeAtRequest) return;
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
      if (completeBasicRecordingIntent({ scope_key: scopeKey, request, operation_id: localIntent?.operation_id }, result.operation_id) === 'superseded') return;
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
  const completion = recordingCompletion(localIntent, operation, workspaceId, selectedDeviceId);
  const renderStartResult = (primary: boolean) => operation ? <BasicRecordingProgress
          groupNames={Object.fromEntries(groups.map((group) => [group.id, group.name]))}
          deviceNames={Object.fromEntries(state.devices.map((device) => [device.id, device.name]))}
          operation={operation}
          intentMatchesCurrent={intentMatchesCurrent}
          completion={completion} primary={primary}
          statusReason={statusReason}
          nextAction={safeNextAction}
          onCommit={onCommit}
          onOpenAdvanced={() => document.getElementById(`basic-schema-${selectedDeviceId}`)?.focus()}
        /> : undefined;

  return (
    <section className="space-y-4 rounded-xl border border-[var(--gw-focus)] bg-[var(--gw-surface)] p-5" aria-labelledby="basic-recording-title" data-testid="basic-recording-panel">
      {readError}
      <div>
        <h3 id="basic-recording-title" className="text-sm font-semibold text-[var(--gw-text)]">{t('step4.basic.title')}</h3>
        <p className="mt-1 text-xs leading-5 text-[var(--gw-muted)]">{t('step4.basic.subtitle')}</p>
        {supportsCircuitProfile && !shareOnly && <button type="button" data-testid="device-recording-back-setup" disabled={operationActive} onClick={() => request(() => setLegacyView(undefined))} className="mt-2 text-xs text-[var(--gw-accent)] underline disabled:opacity-40">{t('step4.basic.existing.back')}</button>}
      </div>

      {deviceIds.length > 0 && (
        <div className="block text-xs text-[var(--gw-text)]">
          <span>{t('step4.basic.device')}</span>
          <div className="mt-1">
            <DeviceScopeTiles
              devices={deviceIds.map((deviceId) => {
                const device = state.devices.find((item) => item.id === deviceId);
                return { id: deviceId, name: device?.name };
              })}
              selectedId={selectedDeviceId}
              onSelect={selectByOperator}
              label={t('step4.basic.device')}
            />
          </div>
        </div>
      )}

      <div className="grid gap-3 sm:grid-cols-2">
        <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block text-xs text-[var(--gw-text)]">
          <span>{t('step4.basic.interval')}</span>
          <input type="number" min={1} value={intervalSeconds} disabled={operationActive || readonly} onChange={(event) => {
            const next = Math.max(1, Number(event.target.value) || 1);
            saveBasicLegacyDraft(scopeKey, { intervalSeconds: next });
            setIntervalSeconds(next);
          }} className="mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-2 text-sm text-[var(--gw-text)]" data-testid="basic-recording-interval" />
        </label><ParameterHelp help="policy.interval" /></div>
        <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block text-xs text-[var(--gw-text)]">
          <span>{t('step4.basic.completeness')}</span>
          <select value={incompletePolicy} disabled={operationActive || readonly} onChange={(event) => {
            saveBasicLegacyDraft(scopeKey, { incompletePolicy: event.target.value });
            setIncompletePolicy(event.target.value);
          }} className="mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-2 text-sm text-[var(--gw-text)]" data-testid="basic-recording-completeness">
            <option value="skip_row">{t('step4.basic.skip_row')}</option>
            <option value="partial">{t('step4.basic.partial_row')}</option>
          </select>
        </label><ParameterHelp help="policy.incomplete" /></div>
      </div>

      {!shareOnly && <div className="space-y-1 text-xs text-[var(--gw-muted)]" data-testid="basic-required-policy">
        <p>{t('basic.required_hint', { ns: 'studio-review' })}</p>
        <p>{t('basic.required_count', { ns: 'studio-review', count: (basicSelectedGroup?.members ?? selectedDraft?.members ?? []).filter(member => member.required).length })}</p>
        {basicSelectedGroup && onOpenMemberSettings
          ? <button type="button" disabled={operationActive} className="text-[var(--gw-accent)] underline" data-testid="basic-member-settings" onClick={() => onOpenMemberSettings(basicSelectedGroup.id)}>{t('basic.member_settings', { ns: 'studio-review' })}</button>
          : <p>{t('basic.prepare_first', { ns: 'studio-review' })}</p>}
      </div>}

      <p className="text-xs text-[var(--gw-muted)]" data-testid="basic-recording-destination">
        {shareOnly
          ? t('step4.basic.share_only_destination')
          : t('step4.basic.destination', { name: safeRuntimeLabel(state.db.connector.name) || state.db.connector.kind })}
      </p>

      {blockReason && <p role="status" className="text-xs text-[var(--gw-warning-fg)]" data-testid="basic-recording-blocked">{t(`step4.basic.blocked.${blockReason}`)}</p>}
      {destinationMismatch && <p role="alert" className="text-xs text-[var(--gw-warning-fg)]" data-testid="basic-recording-destination-mismatch">{t('step4.basic.blocked.destination_mismatch')}</p>}
      {advancedGroups.length > 0 && <p className="text-xs text-[var(--gw-muted)]" data-testid="basic-recording-advanced-note">{t('step4.basic.advanced_note', { count: advancedGroups.length })}</p>}
      {!shareOnly && (selectedGroup?.members.length ?? selectedCandidates.length) > 0 && (
        <>
          <BasicRecordingMembers candidates={selectedCandidates} group={selectedGroup} onRepair={onOpenAdvanced} />
          <p className="text-xs text-[var(--gw-muted)]">{t('step4.basic.members', { count: selectedGroup?.members.length ?? selectedCandidates.length })}</p>
        </>
      )}
      {membersMismatch && basicSelectedGroup && <BasicMemberUpdate group={basicSelectedGroup} deviceId={selectedDeviceId} candidates={selectedCandidates} workspaceRevision={workspaceRevision} onOpenAdvanced={onOpenAdvanced} isBlocked={() => isPreparationBlocked(basicSelectedGroup.id)} disabled={excluded.length > 0 || readonly || operationActive || destinationMismatch || !workspaceSnapshotReady || operation?.status === 'unknown' || Boolean(localIntent && !operation)} onUpdated={() => void groupsQuery.refetch()} />}

      {!shareOnly && (
        <BasicRecordingEvidence
          intervalSeconds={intervalSeconds}
          deliveryError={deliveryQuery.isError}
          stages={currentRevisionStages}
          committedEffect={currentCommittedEffect}
          group={selectedGroup}
        />
      )}

      {!shareOnly && <div id={`basic-schema-${selectedDeviceId}`} tabIndex={-1}><BasicRecordingPreparation
        isNewCommandBlocked={workspaceBlocked} workspaceId={workspaceId ?? ''} workspaceRevision={workspaceRevision} deviceId={selectedDeviceId} ownerKey={scopeKey} fenceKey={preparationKey} onCommandFence={receivePreparationFence}
        group={basicSelectedGroup} draft={selectedDraft} candidates={selectedCandidates} intervalSeconds={intervalSeconds} incompletePolicy={incompletePolicy}
        disabled={Boolean(blockReason) || groupsQuery.isLoading || groupsQuery.isFetching || groupsQuery.isError} readonly={readonly}
        destinationMatches={!destinationMismatch} ensureGroup={ensureGroup} onPrepared={() => groupsQuery.refetch()} summary={renderScopeSummary}
        completedOperationId={completion.state === 'available' && completion.handoffDeviceId ? completion.operationId : undefined}
        renderStartResult={renderStartResult}
        start={handleStart} check={async () => { if (operationId && localIntent?.scope_key === scopeKey && activeSelectionRef.current === scopeKey) await operationQuery.refetch(); }}
        live={{ unresolved: Boolean(localIntent?.scope_key === scopeKey && localIntent && (!operation || completion.state === 'unavailable' || !intentMatchesCurrent || ['pending', 'running', 'partial', 'unknown'].includes(operation.status) || operation.status === 'failed' && operation.next_action === 'retry_same_request')),
          operationId: localIntent?.scope_key === scopeKey ? operationId : undefined, status: localIntent?.scope_key === scopeKey ? operation?.status : undefined, queryFailed: operationQuery.isError, busy: actionBusy || operationQuery.isFetching }} /></div>}
      {!shareOnly && selectedGroup && !deliveryQuery.isError && deliveryQuery.data?.group_id === selectedGroup.id && <GroupDeliveryRecovery groupId={selectedGroup.id} items={deliveryQuery.data.attention ?? []} readonly={readonly || operationActive} onResolved={() => void deliveryQuery.refetch()} />}

      {!shareOnly && basicSelectedGroup && lifecycleHost && <SavedRecordingManagement group={basicSelectedGroup} workspaceId={workspaceId} workspaceRevision={workspaceRevision} readonly={readonly || operationActive} lifecycleHost={lifecycleHost} />}
      {!shareOnly && <RemovedRecordingManagement groups={groups} deviceId={selectedDeviceId} workspaceId={workspaceId} workspaceRevision={workspaceRevision} readonly={readonly} lifecycleHost={lifecycleHost} onOpenAdvanced={onOpenAdvanced} />}
      {supportsCircuitProfile && !shareOnly && hasLegacyBasicWide && destinationSaved && <LegacyMigrationEntry disabled={readonly || operationActive} onOpen={() => request(() => setLegacyView({ scope: legacyScope, mode: 'prepare' }))} />}

      {!shareOnly && localIntent && operationId && <BasicRecordingRetry context={JSON.stringify([preparationKey, localIntent.request.request_id, operationId])}
        eligible={isOriginalRequestRetry(operation, operationId, localIntent.request)} onRetry={handleRetry}
        disabled={actionBusy || !intentMatchesCurrent || Boolean(blockReason) || recordingStartRecoveryBlocked(scopeKey, workspaceId ?? '', localIntent.request.groups.map(group => group.group_id), operationId)}
        verify={async () => { if (workspaceBlocked()) return false; const result = await operationQuery.refetch(); const retained = loadBasicRecordingIntent(scopeKey);
          return !workspaceBlocked() && !result.isError && !result.error && retained?.operation_id === operationId && JSON.stringify(retained.request) === JSON.stringify(localIntent.request) &&
            isOriginalRequestRetry(result.data, operationId, localIntent.request) && !recordingStartRecoveryBlocked(scopeKey, workspaceId ?? '', localIntent.request.groups.map(group => group.group_id), operationId); }} />}
      {shareOnly && <div className="flex flex-wrap items-center gap-3">
        <button type="button" disabled={Boolean(blockReason) || operationActive || groupsQuery.isLoading} onClick={() => void handleStart()} className="rounded bg-[var(--gw-active)] px-4 py-2 text-xs font-semibold text-[var(--gw-text)] disabled:cursor-not-allowed disabled:opacity-40" data-testid="basic-recording-start">
          {operationActive ? t('step4.basic.starting') : t('step4.basic.start')}
        </button>
        {operation?.status && ['pending', 'running', 'partial', 'failed', 'unknown'].includes(operation.status) && <button type="button" disabled={actionBusy || !intentMatchesCurrent || Boolean(blockReason)} onClick={handleRetry} className="rounded border border-[var(--gw-control)] px-3 py-2 text-xs text-[var(--gw-text)] disabled:opacity-40" data-testid="basic-recording-retry">{t('step4.basic.retry')}</button>}
        {operationId && operationQuery.isError && <button type="button" onClick={() => void operationQuery.refetch()} className="rounded border border-[var(--gw-control)] px-3 py-2 text-xs text-[var(--gw-text)]" data-testid="basic-recording-check">{t('step4.basic.check')}</button>}
      </div>}

      {safeActionError && <p role="alert" className="text-xs text-[var(--gw-danger-fg)]" data-testid="basic-recording-error">{safeActionError}</p>}
      {shareOnly && renderStartResult(false)}
      {operationQuery.isError && operationId && <p role="alert" className="text-xs text-[var(--gw-warning-fg)]" data-testid="basic-recording-query-error">{t('step4.basic.query_failed')}</p>}
    </section>
  );
};
