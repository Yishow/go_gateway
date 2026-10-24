import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useStudioV2WorkspaceMeasurementsQuery } from '../../../../../hooks/datalink/useStudioV2WorkspaceMeasurements';
import { useEnsureCircuitProfileMutation, useUpdateWriteGroupMutation, useWriteGroupsQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { MeasurementDefinition } from '../../../../../types/measurement';
import type { WriteGroup, WriteGroupProfileKind } from '../../../../../types/studioV2WriteGroup';
import { basicRecordingScopeKey, createBasicRecordingRequestId } from '../../state/basicRecording';
import { safeBasicStartError } from '../../state/basicRecordingMessages';
import { buildCircuitProfileGroupDraft, candidateKey, canonicalProfileGroupMatchesDraft, defaultCircuitProfileLocalDraft, deriveCircuitId, evaluateCandidateCounterConfirmation, groupCircuitId, isManagedLongProfileGroupForDevice, loadCircuitProfileLocalDraft, parseStrictProfileIntervalInput, saveCircuitProfileLocalDraft } from '../../state/circuitProfileBasic';
import type { WorkbenchV2State } from '../../state/types';
import { buildGroupCandidates } from '../../state/writeGroup/candidates';
import { CircuitProfileBasicActions, CircuitProfileStartStatus } from './CircuitProfileBasicActions';
import { CircuitProfileBasicHeader } from './CircuitProfileBasicHeader';
import { destinationIsSaved } from './basicRecordingPanelHelpers';
import { ensureOrUpdateCircuitProfileGroups } from './circuitProfileActions';
import { CircuitProfileBatchEditor } from './CircuitProfileBatchEditor';
import { CircuitProfileConfigFields } from './CircuitProfileConfigFields';
import { CircuitProfileSchemaPreparation } from './CircuitProfileSchemaPreparation';
import { extractDeviceCircuits, loadActiveCircuitId, saveActiveCircuitId } from '../../state/circuitProfileBatch';
import { supportsWriteGroupKind } from './writeGroupCapabilities';
import { CircuitProfileStatusList } from './CircuitProfileStatusList';
import { CircuitProfileRoles } from './CircuitProfileRoles';
import { focusDestinationSettings } from './step4FocusTargets';
import { useCircuitProfileBatchHost } from './useCircuitProfileBatchHost';
import { useReportOperationActive } from './useReportOperationActive';
import { useCircuitProfileStart } from './useCircuitProfileStart';

function readMeasurementsList(data: unknown): MeasurementDefinition[] {
  return Array.isArray(data) ? (data as MeasurementDefinition[]) : (data && typeof data === 'object' && Array.isArray((data as { measurements?: unknown }).measurements) ? (data as { measurements: MeasurementDefinition[] }).measurements : []);
}

export interface CircuitProfileBasicPanelProps {
  state: WorkbenchV2State; workspaceId?: string; readonly: boolean; readinessToken?: string;
  settingsRevision?: string; shareWorkspaceRevision?: string; workspaceRevision?: string; onActivated?: (deviceIds: string[]) => Promise<void> | void;
  onOperationActiveChange?: (active: boolean) => void; onContinueToRuntime?: (deviceIds: string[]) => void; onOpenAdvanced?: (groupId?: string) => void;
  /** Target preparation for a reviewed cutover: no recording start, apply or stop is offered or restored here. */
  prepareOnly?: boolean; onClose?: () => void;
}

export const CircuitProfileBasicPanel: React.FC<CircuitProfileBasicPanelProps> = ({
  state, workspaceId, readonly, readinessToken, settingsRevision, shareWorkspaceRevision, workspaceRevision: propWorkspaceRevision,
  onActivated, onContinueToRuntime, onOpenAdvanced, onOperationActiveChange, prepareOnly = false, onClose,
}) => {
  const { t } = useTranslation('workbench-v2');
  const groupsQuery = useWriteGroupsQuery(Boolean(workspaceId));
  const measurementsQuery = useStudioV2WorkspaceMeasurementsQuery(Boolean(workspaceId));
  const ensureProfile = useEnsureCircuitProfileMutation();
  const updateGroup = useUpdateWriteGroupMutation();

  const savedDevices = React.useMemo(
    () => state.devices.filter((d) => Boolean(d.persisted && d.save_state !== 'saving' && d.save_state !== 'save-error')),
    [state.devices],
  );
  const [selectedDeviceId, setSelectedDeviceId] = React.useState(() => savedDevices[0]?.id ?? '');
  const activeDeviceId = savedDevices.some((d) => d.id === selectedDeviceId) ? selectedDeviceId : (savedDevices[0]?.id ?? '');
  const activeDevice = savedDevices.find((d) => d.id === activeDeviceId);
  const destination = React.useMemo(() => (
    destinationIsSaved(state) ? {
      connector_id: state.db.connector.connector_id ?? '', connector_revision: state.db.connector.identity_revision ?? '',
      table_schema: state.db.connector.schema ?? '', table_name: '', storage_strategy: 'managed' as const,
    } : undefined
  ), [state]);
  const kindStr = String(state.db.connector.kind);
  const isPostgres = kindStr === 'postgres' || kindStr === 'postgresql';
  const connectorKindSupported = supportsWriteGroupKind(state.db.connector.kind) || isPostgres;

  const { candidates: allCandidates } = React.useMemo(() => buildGroupCandidates(state), [state]);
  const deviceCandidates = React.useMemo(() => allCandidates.filter((c) => c.device_id === activeDeviceId), [allCandidates, activeDeviceId]);
  const liveGroups = React.useMemo(() => (groupsQuery.data?.groups ?? []).filter((g) => g.status !== 'deleted'), [groupsQuery.data?.groups]);
  const savedLongGroups = React.useMemo(() => liveGroups.filter((g) => isManagedLongProfileGroupForDevice(g, activeDeviceId)), [liveGroups, activeDeviceId]);
  const advancedGroups = React.useMemo(() => liveGroups.filter((g) => !isManagedLongProfileGroupForDevice(g, activeDeviceId)), [liveGroups, activeDeviceId]);

  const savedCircuits = React.useMemo(() => extractDeviceCircuits(savedLongGroups), [savedLongGroups]);
  const [selectedCircuitId, setSelectedCircuitId] = React.useState<string | undefined>(() => {
    if (!workspaceId || !activeDeviceId) return undefined;
    return loadActiveCircuitId(workspaceId, activeDeviceId);
  });
  const activeCircuitId = selectedCircuitId !== undefined ? selectedCircuitId : savedCircuits[0]?.circuit_id;

  const [localDraft, setLocalDraft] = React.useState(() => loadCircuitProfileLocalDraft(workspaceId ?? '', activeDeviceId, deviceCandidates, savedLongGroups, activeCircuitId, prepareOnly));
  const [preparedGroups, setPreparedGroups] = React.useState<WriteGroup[]>([]);
  const [revisionOverride, setRevisionOverride] = React.useState<{ base: string; revision: string }>();
  const [errorMessage, setErrorMessage] = React.useState<string>();
  const [batchOpen, setBatchOpen] = React.useState(false);

  React.useEffect(() => {
    if (!workspaceId || !activeDeviceId) return;
    const stored = loadActiveCircuitId(workspaceId, activeDeviceId);
    setSelectedCircuitId(stored);
  }, [workspaceId, activeDeviceId]);

  React.useEffect(() => {
    setLocalDraft(loadCircuitProfileLocalDraft(workspaceId ?? '', activeDeviceId, deviceCandidates, savedLongGroups, activeCircuitId, prepareOnly));
    setPreparedGroups([]);
  }, [workspaceId, activeDeviceId, deviceCandidates, savedLongGroups, activeCircuitId, prepareOnly]);

  const handleSelectCircuit = (targetId: string) => {
    if (targetId === '__new__') {
      const newId = `cir-${createBasicRecordingRequestId()}`;
      setSelectedCircuitId(newId);
      if (workspaceId && activeDeviceId) saveActiveCircuitId(workspaceId, activeDeviceId, newId);
      const draft = { ...defaultCircuitProfileLocalDraft(deviceCandidates, savedLongGroups, newId, prepareOnly), circuit_id: newId, circuit_name: '' };
      setLocalDraft(draft);
      if (workspaceId && activeDeviceId) saveCircuitProfileLocalDraft(workspaceId, activeDeviceId, draft);
    } else {
      setSelectedCircuitId(targetId);
      if (workspaceId && activeDeviceId) saveActiveCircuitId(workspaceId, activeDeviceId, targetId);
    }
  };

  const updateLocalDraft = React.useCallback(
    (updater: (prev: typeof localDraft) => typeof localDraft) => {
      setLocalDraft((prev) => {
        const next = updater(prev);
        if (workspaceId && activeDeviceId) {
          saveCircuitProfileLocalDraft(workspaceId, activeDeviceId, next);
          if (next.circuit_id) saveActiveCircuitId(workspaceId, activeDeviceId, next.circuit_id);
        }
        return next;
      });
      setPreparedGroups([]);
    },
    [workspaceId, activeDeviceId],
  );

  const queriedRevision = groupsQuery.data?.workspace_revision || propWorkspaceRevision || '';
  const effectiveWorkspaceRevision = revisionOverride?.base === queriedRevision ? revisionOverride.revision : queriedRevision;

  const instantaneousCandidates = React.useMemo(() => {
    const keySet = new Set(localDraft.instantaneous_keys ?? localDraft.selected_keys);
    return deviceCandidates.filter((c) => keySet.has(candidateKey(c)));
  }, [deviceCandidates, localDraft.instantaneous_keys, localDraft.selected_keys]);

  const cumulativeCandidates = React.useMemo(() => {
    const keySet = new Set(localDraft.cumulative_keys ?? localDraft.selected_keys);
    return deviceCandidates.filter((c) => keySet.has(candidateKey(c)));
  }, [deviceCandidates, localDraft.cumulative_keys, localDraft.selected_keys]);

  const measurements = React.useMemo(() => readMeasurementsList(measurementsQuery.data), [measurementsQuery.data]);
  const counterStatuses = React.useMemo(() => cumulativeCandidates.map((c) => evaluateCandidateCounterConfirmation(workspaceId ?? '', c, measurements)), [workspaceId, cumulativeCandidates, measurements]);
  const unconfirmedCounterStatuses = React.useMemo(() => counterStatuses.filter((s) => !s.confirmed), [counterStatuses]);
  const confirmedMeasurementsByPointKey = React.useMemo(() => {
    const map = new Map<string, MeasurementDefinition>();
    for (const status of counterStatuses) if (status.confirmed && status.measurement) map.set(candidateKey(status.candidate), status.measurement);
    return map;
  }, [counterStatuses]);

  const instParsed = parseStrictProfileIntervalInput('instantaneous', localDraft.instantaneous_interval_input);
  const cumParsed = parseStrictProfileIntervalInput('cumulative', localDraft.cumulative_interval_input);
  const circuitId = deriveCircuitId(localDraft.circuit_name, localDraft.circuit_id);
  const currentGroups = savedLongGroups.filter((group) => groupCircuitId(group) === circuitId);
  const storageLocked = currentGroups.length > 0;
  const effectiveStorageMode = storageLocked ? (currentGroups[0]?.storage_mode === 'timescale' ? 'timescale' : 'plain') : (isPostgres ? localDraft.storage_mode : 'plain');

  const enabledProfiles = React.useMemo(() => {
    const list: Array<{ profile: WriteGroupProfileKind; intervalValue: number; intervalSeconds: number; candidates: typeof deviceCandidates }> = [];
    if (localDraft.instantaneous_enabled && instParsed.valid && instParsed.intervalValue && instParsed.intervalSeconds && instantaneousCandidates.length > 0) {
      list.push({ profile: 'instantaneous', intervalValue: instParsed.intervalValue, intervalSeconds: instParsed.intervalSeconds, candidates: instantaneousCandidates });
    }
    if (localDraft.cumulative_enabled && cumParsed.valid && cumParsed.intervalValue && cumParsed.intervalSeconds && cumulativeCandidates.length > 0) {
      list.push({ profile: 'cumulative', intervalValue: cumParsed.intervalValue, intervalSeconds: cumParsed.intervalSeconds, candidates: cumulativeCandidates });
    }
    return list;
  }, [localDraft.instantaneous_enabled, localDraft.cumulative_enabled, instParsed, cumParsed, instantaneousCandidates, cumulativeCandidates]);

  const scopeKey = workspaceId && activeDeviceId ? basicRecordingScopeKey(workspaceId, activeDeviceId,
    encodeURIComponent(JSON.stringify([circuitId, destination?.connector_id, destination?.connector_revision, destination?.table_schema, ...(prepareOnly ? ['cutover-prepare'] : [])]))) : '';
  const activeScope = React.useRef(scopeKey);
  activeScope.current = scopeKey;
  const batch = useCircuitProfileBatchHost({ scopeKey, workspaceId, deviceId: activeDeviceId, candidates: deviceCandidates, savedLongGroups,
    activeCircuitId, queriedRevision, refetchGroups: groupsQuery.refetch, setRevisionOverride, setLocalDraft });
  const batchBusy = batch.busy;
  React.useEffect(() => { setPreparedGroups([]); setRevisionOverride(undefined); setErrorMessage(undefined); setBatchOpen(false); }, [scopeKey]);
  const counterDefinitionsReady = !measurementsQuery.isLoading && !measurementsQuery.isError && Array.isArray(measurementsQuery.data);
  const instantaneousSourcesMissing = localDraft.instantaneous_enabled && instantaneousCandidates.length === 0;
  const cumulativeSourcesMissing = localDraft.cumulative_enabled && cumulativeCandidates.length === 0;
  const hasValidIntervals = (!localDraft.instantaneous_enabled || instParsed.valid) && (!localDraft.cumulative_enabled || cumParsed.valid);
  const cumulativeCircuitMissing = localDraft.cumulative_enabled && !circuitId;
  const cumulativeCounterBlocked = localDraft.cumulative_enabled && (!counterDefinitionsReady || unconfirmedCounterStatuses.length > 0);
  const canPrepareOrStart = !readonly && !batchBusy && !batchOpen &&
    Boolean(workspaceId && effectiveWorkspaceRevision && activeDeviceId && destination && connectorKindSupported) &&
    enabledProfiles.length > 0 && hasValidIntervals && !instantaneousSourcesMissing && !cumulativeSourcesMissing &&
    !cumulativeCircuitMissing && !cumulativeCounterBlocked;

  const ensureCurrentGroups = async (): Promise<{ groups: WriteGroup[]; workspaceRevision: string } | undefined> => {
    if (!canPrepareOrStart || !workspaceId || !destination) return undefined;
    const requestedScope = scopeKey;
    setErrorMessage(undefined);
    try {
      const result = await ensureOrUpdateCircuitProfileGroups({
        workspaceId, initialWorkspaceRevision: effectiveWorkspaceRevision, deviceName: activeDevice?.name ?? activeDeviceId,
        enabledProfiles, storageMode: effectiveStorageMode, recreateMissing: localDraft.recreate_missing,
        circuitId, circuitName: localDraft.circuit_name, destination, connectorKind: state.db.connector.kind,
        measurementsByPointKey: confirmedMeasurementsByPointKey, savedLongGroups, preparedGroups,
        ensureCircuitProfile: (vars) => ensureProfile.mutateAsync(vars), updateWriteGroup: (vars) => updateGroup.mutateAsync(vars),
      });
      if (activeScope.current !== requestedScope) return undefined;
      setRevisionOverride({ base: queriedRevision, revision: result.workspaceRevision });
      setPreparedGroups(result.groups);
      await groupsQuery.refetch();
      return activeScope.current === requestedScope ? result : undefined;
    } catch (error) {
      if (activeScope.current === requestedScope) setErrorMessage(safeBasicStartError(error, t));
      return undefined;
    }
  };

  const flow = useCircuitProfileStart({ scopeKey, canStart: canPrepareOrStart, readonly, workspaceId, deviceId: activeDeviceId,
    connectorRevision: destination?.connector_revision ?? '', readinessToken, settingsRevision, shareWorkspaceRevision,
    prepare: ensureCurrentGroups, onActivated, prepareOnly });
  useReportOperationActive(flow.busy || flow.unresolved || batchBusy, onOperationActiveChange); // one merged page signal
  const isEffectiveLocked = readonly || flow.busy || flow.unresolved || batchBusy;
  const schemaGroups = preparedGroups.length > 0 ? preparedGroups : currentGroups;
  const dirtyGroupIds = new Set(schemaGroups.filter((group) => {
    const profile = enabledProfiles.find((item) => item.profile === group.recording_profile?.kind);
    if (!profile || !destination) return true;
    const draft = buildCircuitProfileGroupDraft({ workspaceId: workspaceId ?? '', deviceName: activeDevice?.name ?? activeDeviceId,
      profile: profile.profile, intervalValue: profile.intervalValue, intervalSeconds: profile.intervalSeconds,
      storageMode: effectiveStorageMode, recreateMissing: localDraft.recreate_missing, circuitId, circuitName: localDraft.circuit_name,
      destination, selectedCandidates: profile.candidates, measurementsByPointKey: confirmedMeasurementsByPointKey });
    return !canonicalProfileGroupMatchesDraft(group, draft, state.db.connector.kind);
  }).map((group) => group.id));
  const refreshSchema = async () => {
    const requestedScope = scopeKey;
    const refreshed = await groupsQuery.refetch();
    if (activeScope.current !== requestedScope) return;
    setRevisionOverride(undefined);
    if (refreshed.data) {
      setRevisionOverride({ base: queriedRevision, revision: refreshed.data.workspace_revision });
      setPreparedGroups(refreshed.data.groups.filter((group) => schemaGroups.some((saved) => saved.id === group.id)));
    }
  };

  return (
    <section data-testid="circuit-profile-basic-panel" aria-labelledby="circuit-profile-basic-title" className="space-y-4 rounded-xl border border-cyan-900/70 bg-slate-950/40 p-5">
      <CircuitProfileBasicHeader prepareOnly={prepareOnly} />
      <div data-testid="circuit-profile-destination-summary" className="rounded border border-slate-800/80 bg-slate-950/30 p-3 text-xs text-slate-300">
        <p className="font-semibold text-slate-200">{t('step4.basic.circuit.step_destination')}</p>
        <p className="mt-1 text-slate-400">{t('step4.basic.destination', { name: state.db.connector.name || state.db.connector.kind })} · {state.db.connector.name || state.db.connector.kind} ({state.db.connector.kind})</p>
      </div>
      <CircuitProfileConfigFields
        workspaceId={workspaceId} workspaceRevision={effectiveWorkspaceRevision} readonly={isEffectiveLocked}
        savedDevices={savedDevices} activeDeviceId={activeDeviceId} onSelectDevice={setSelectedDeviceId}
        deviceCandidates={deviceCandidates} localDraft={localDraft} updateLocalDraft={updateLocalDraft}
        instParsed={instParsed} cumParsed={cumParsed} isPostgres={isPostgres} storageLocked={storageLocked}
        effectiveStorageMode={effectiveStorageMode} cumulativeCircuitMissing={cumulativeCircuitMissing} circuitId={circuitId}
        unconfirmedCounterStatuses={counterDefinitionsReady ? unconfirmedCounterStatuses : []} onCounterConfirmed={async () => { await measurementsQuery.refetch(); }}
        savedCircuits={savedCircuits} onSelectCircuit={handleSelectCircuit} onOpenBatch={() => setBatchOpen(true)}
      />
      {batchOpen && (
        <CircuitProfileBatchEditor
          workspaceId={workspaceId ?? ''} workspaceRevision={effectiveWorkspaceRevision} scopeKey={scopeKey}
          connectorRevision={destination?.connector_revision} savedLongGroups={savedLongGroups}
          readonly={isEffectiveLocked} onClose={() => setBatchOpen(false)}
          onSuccess={batch.onSaved} onVerify={batch.reload}
          updateWriteGroup={(vars) => updateGroup.mutateAsync(vars)}
          ensureCircuitProfile={(vars) => ensureProfile.mutateAsync(vars)}
          onBusyChange={batch.setHeld}
        />
      )}
      {localDraft.cumulative_enabled && !counterDefinitionsReady && <p role="status" data-testid={measurementsQuery.isError ? 'circuit-profile-measurements-error' : 'circuit-profile-measurements-loading'} className="text-xs text-amber-300">
        {t(measurementsQuery.isError ? 'step4.basic.circuit.measurements_error' : 'step4.basic.circuit.measurements_loading')}
        {measurementsQuery.isError && <button type="button" onClick={() => void measurementsQuery.refetch()}>{t('step4.group.retry')}</button>}
      </p>}
      <CircuitProfileSchemaPreparation
        workspaceId={workspaceId} workspaceRevision={effectiveWorkspaceRevision}
        preparedGroups={schemaGroups} dirtyGroupIds={dirtyGroupIds}
        disabled={!canPrepareOrStart || flow.unresolved} readonly={isEffectiveLocked} isPreparing={flow.busy || batchBusy || ensureProfile.isPending || updateGroup.isPending}
        destinationMatches={Boolean(destination)} onPrepare={() => void flow.prepare()} onApplied={() => void refreshSchema()} onRepairDestination={focusDestinationSettings} onOpenAdvanced={onOpenAdvanced}
      />
      {advancedGroups.length > 0 && <p data-testid="basic-recording-advanced-note" className="text-xs text-slate-400">{t('step4.basic.advanced_note', { count: advancedGroups.length })}</p>}
      {(errorMessage || flow.error) && <p role="alert" data-testid="circuit-profile-error" className="text-xs text-rose-300">{errorMessage || flow.error}</p>}
      <CircuitProfileStartStatus flow={flow} />
      <CircuitProfileStatusList groups={currentGroups} candidates={deviceCandidates} onOpenAdvanced={onOpenAdvanced} draftEnabled={{ instantaneous: localDraft.instantaneous_enabled, cumulative: localDraft.cumulative_enabled }} /><CircuitProfileRoles onOpenAdvanced={onOpenAdvanced} />
      <CircuitProfileBasicActions flow={flow} canStart={canPrepareOrStart} readonly={readonly} batchBusy={batchBusy} prepareOnly={prepareOnly} onClose={onClose} onContinueToRuntime={onContinueToRuntime} onOpenAdvanced={onOpenAdvanced} />
    </section>
  );
};
