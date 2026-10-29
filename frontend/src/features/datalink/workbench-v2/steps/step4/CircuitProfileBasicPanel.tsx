import { safeRuntimeLabel } from '../../../../../utils/runtimeOperatorSafety';
import { circuitRecordingScope, type RenderRecordingScopeSummary } from './recordingScopeSummary';
import type { SavedRecordingActionHost } from './writeGroup/lifecycle/SavedRecordingActions';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useStudioV2WorkspaceMeasurementsQuery } from '../../../../../hooks/datalink/useStudioV2WorkspaceMeasurements';
import { useStudioV2MySQLRecordingCapabilityQuery, useStudioV2SQLServerRecordingCapabilityQuery } from '../../../../../hooks/datalink/useStudioV2WorkspaceDatabase';
import { useEnsureCircuitProfileMutation, useUpdateWriteGroupMutation, useWriteGroupsQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { MeasurementDefinition } from '../../../../../types/measurement';
import type { StudioV2MySQLRecordingCapability } from '../../../../../types/studioV2MySQLRecordingCapability';
import type { StudioV2SQLServerRecordingCapability } from '../../../../../types/studioV2SQLServerRecordingCapability';
import type { WriteGroup, WriteGroupProfileKind } from '../../../../../types/studioV2WriteGroup';
import { basicRecordingScopeKey, createBasicRecordingRequestId } from '../../state/basicRecording';
import { safeBasicStartError } from '../../state/basicRecordingMessages';
import { candidateKey, defaultCircuitProfileLocalDraft, defaultMemberRole, deriveCircuitId, evaluateCandidateCounterConfirmation, groupCircuitId, groupProfileKind, isManagedProfileGroupForDevice, loadCircuitProfileLocalDraft, parseStrictProfileIntervalInput, purposeOfLayout, saveCircuitProfileLocalDraft } from '../../state/circuitProfileBasic';
import { circuitProfileLayoutBlockReason } from '../../state/circuitProfileCapabilities';
import type { WorkbenchV2State } from '../../state/types';
import { buildGroupCandidates } from '../../state/writeGroup/candidates';
import { CircuitProfileBasicActions } from './CircuitProfileBasicActions';
import { CircuitProfileBasicHeader } from './CircuitProfileBasicHeader';
import { ExistingDeviceRecording } from './ExistingDeviceRecording';
import { destinationIsSaved } from './basicRecordingPanelHelpers';
import { ensureOrUpdateCircuitProfileGroups } from './circuitProfileActions';
import { CircuitProfileBatchEditor } from './CircuitProfileBatchEditor';
import { CircuitProfileConfigFields } from './CircuitProfileConfigFields';
import { MySQLRecordingCapabilityPanel } from './MySQLRecordingCapabilityPanel';
import { SQLServerRecordingCapabilityPanel } from './SQLServerRecordingCapabilityPanel';
import { CircuitRecordingPreparation } from './CircuitRecordingPreparation';
import { circuitPreparationDrafts } from './circuitPreparationDrafts';
import { extractDeviceCircuits, loadActiveCircuitId, saveActiveCircuitId } from '../../state/circuitProfileBatch';
import { CircuitProfileStatusList } from './CircuitProfileStatusList';
import { CircuitProfileRoles } from './CircuitProfileRoles';
import { focusDestinationSettings } from './step4FocusTargets';
import { useCircuitProfileBatchHost } from './useCircuitProfileBatchHost';
import { useReportOperationActive } from './useReportOperationActive';
import { loadSelectedDeviceId, saveSelectedDeviceId } from '../../state/deviceScopeSelection';
import { useRecordingWorkspaceGuard } from './useRecordingWorkspaceGuard';
import { useCircuitRecordingCommandFence } from './useCircuitRecordingCommandFence';
import { useCircuitProfileStart } from './useCircuitProfileStart';

function readMeasurementsList(data: unknown): MeasurementDefinition[] {
  return Array.isArray(data) ? (data as MeasurementDefinition[]) : (data && typeof data === 'object' && Array.isArray((data as { measurements?: unknown }).measurements) ? (data as { measurements: MeasurementDefinition[] }).measurements : []);
}

export interface CircuitProfileBasicPanelProps {
  lifecycleHost?: SavedRecordingActionHost;
  renderScopeSummary?: RenderRecordingScopeSummary;
  legacyGroup?: WriteGroup; onEditLegacy?: () => void;
  state: WorkbenchV2State; workspaceId?: string; readonly: boolean; readinessToken?: string;
  settingsRevision?: string; shareWorkspaceRevision?: string; workspaceRevision?: string; onActivated?: (deviceIds: string[]) => Promise<void> | void;
  /** The parent's viewed device. A controlled value keeps the legacy and new branches on one authority. */
  deviceId?: string; onSelectDevice?: (deviceId: string) => void;
  onOperationActiveChange?: (active: boolean) => void; onContinueToRuntime?: (deviceIds: string[]) => void; onOpenAdvanced?: (groupId?: string) => void;
  /** Target preparation for a reviewed cutover: no recording start, apply or stop is offered or restored here. */
  prepareOnly?: boolean; onClose?: () => void;
  /**
   * Shows one saved circuit of its workspace and device, exactly as choosing it would; it starts and saves nothing.
   * The parent owns the request and drops it when this panel reports it done.
   */
  circuitRequest?: { workspaceId: string; deviceId: string; circuitId: string };
  onCircuitRequestDone?: () => void;
}

interface MySQLCapabilityQueryState {
  data?: StudioV2MySQLRecordingCapability;
  isLoading: boolean;
  isFetching: boolean;
  isError: boolean;
  refetch: () => Promise<unknown>;
}

const EMPTY_MYSQL_CAPABILITY_QUERY: MySQLCapabilityQueryState = {
  data: undefined, isLoading: false, isFetching: false, isError: false, refetch: async () => undefined,
};

interface SQLServerCapabilityQueryState {
  data?: StudioV2SQLServerRecordingCapability;
  isLoading: boolean;
  isFetching: boolean;
  isError: boolean;
  refetch: () => Promise<unknown>;
}

const EMPTY_SQLSERVER_CAPABILITY_QUERY: SQLServerCapabilityQueryState = {
  data: undefined, isLoading: false, isFetching: false, isError: false, refetch: async () => undefined,
};

export const CircuitProfileBasicPanel: React.FC<CircuitProfileBasicPanelProps> = (props) => {
  const connector = props.state.db.connector;
  const kind = String(connector.kind).toLowerCase();
  if (kind === 'mysql' && destinationIsSaved(props.state)) {
    return <SavedMySQLCircuitProfileBasicPanel {...props} />;
  }
  if (kind === 'sqlserver' && destinationIsSaved(props.state)) {
    return <SavedSQLServerCircuitProfileBasicPanel {...props} />;
  }
  return <CircuitProfileBasicPanelContent
    {...props}
    mysqlCapabilityQuery={EMPTY_MYSQL_CAPABILITY_QUERY}
    sqlserverCapabilityQuery={EMPTY_SQLSERVER_CAPABILITY_QUERY}
  />;
};

const SavedMySQLCircuitProfileBasicPanel: React.FC<CircuitProfileBasicPanelProps> = (props) => {
  const connector = props.state.db.connector;
  const capabilityQuery = useStudioV2MySQLRecordingCapabilityQuery(
    connector.connector_id ?? '', connector.identity_revision ?? '', true,
  );
  return <CircuitProfileBasicPanelContent
    {...props}
    mysqlCapabilityQuery={capabilityQuery}
    sqlserverCapabilityQuery={EMPTY_SQLSERVER_CAPABILITY_QUERY}
  />;
};

const SavedSQLServerCircuitProfileBasicPanel: React.FC<CircuitProfileBasicPanelProps> = (props) => {
  const connector = props.state.db.connector;
  const capabilityQuery = useStudioV2SQLServerRecordingCapabilityQuery(
    connector.connector_id ?? '', connector.identity_revision ?? '', true,
  );
  return <CircuitProfileBasicPanelContent
    {...props}
    mysqlCapabilityQuery={EMPTY_MYSQL_CAPABILITY_QUERY}
    sqlserverCapabilityQuery={capabilityQuery}
  />;
};

type CircuitProfileBasicPanelContentProps = CircuitProfileBasicPanelProps & {
  mysqlCapabilityQuery: MySQLCapabilityQueryState;
  sqlserverCapabilityQuery: SQLServerCapabilityQueryState;
};

const CircuitProfileBasicPanelContent: React.FC<CircuitProfileBasicPanelContentProps> = ({
  state, workspaceId, readonly, lifecycleHost, readinessToken, settingsRevision, shareWorkspaceRevision, workspaceRevision: propWorkspaceRevision,
  onActivated, onContinueToRuntime, onOpenAdvanced, onOperationActiveChange, prepareOnly = false, onClose, legacyGroup, onEditLegacy,
  deviceId, onSelectDevice, renderScopeSummary, circuitRequest, onCircuitRequestDone, mysqlCapabilityQuery, sqlserverCapabilityQuery,
}) => {
  const { t } = useTranslation('workbench-v2');
  const groupsQuery = useWriteGroupsQuery(Boolean(workspaceId));
  const workspaceBlocked = useRecordingWorkspaceGuard(workspaceId, groupsQuery.data);
  const queriedRevision = groupsQuery.data?.workspace_revision || propWorkspaceRevision || '';
  const ensureProfile = useEnsureCircuitProfileMutation();
  const updateGroup = useUpdateWriteGroupMutation();

  const savedDevices = React.useMemo(
    () => state.devices.filter((d) => Boolean(d.persisted && d.save_state !== 'saving' && d.save_state !== 'save-error')),
    [state.devices],
  );
  const [ownDeviceId, setOwnDeviceId] = React.useState(() => {
    const ids = savedDevices.map((device) => device.id);
    return loadSelectedDeviceId(workspaceId ?? '', ids) ?? ids[0] ?? '';
  });
  const selectedDeviceId = deviceId ?? ownDeviceId;
  const activeDeviceId = savedDevices.some((d) => d.id === selectedDeviceId) ? selectedDeviceId : (savedDevices[0]?.id ?? '');
  const measurementsQuery = useStudioV2WorkspaceMeasurementsQuery(Boolean(workspaceId && activeDeviceId && queriedRevision), activeDeviceId,
    { workspaceId: workspaceId ?? '', workspaceRevision: queriedRevision });
  const selectDevice = React.useCallback((nextDeviceId: string) => {
    saveSelectedDeviceId(workspaceId ?? '', nextDeviceId);
    setOwnDeviceId(nextDeviceId);
    onSelectDevice?.(nextDeviceId);
  }, [onSelectDevice, workspaceId]);
  const activeDevice = savedDevices.find((d) => d.id === activeDeviceId);
  const destination = React.useMemo(() => (
    destinationIsSaved(state) ? {
      connector_id: state.db.connector.connector_id ?? '', connector_revision: state.db.connector.identity_revision ?? '',
      table_schema: state.db.connector.schema ?? '', table_name: '', storage_strategy: 'managed' as const,
    } : undefined
  ), [state]);
  const connectorAlias = safeRuntimeLabel(state.db.connector.name) || state.db.connector.kind;
  const kindStr = String(state.db.connector.kind).toLowerCase();
  const isPostgres = kindStr === 'postgres' || kindStr === 'postgresql';
  const isMySQL = kindStr === 'mysql';
  const isSQLServer = kindStr === 'sqlserver';
  const mysqlConnectorId = state.db.connector.connector_id ?? '';
  const mysqlConnectorRevision = state.db.connector.identity_revision ?? '';
  const mysqlCapability = !destination || mysqlCapabilityQuery.isError || mysqlCapabilityQuery.isFetching
    ? undefined : mysqlCapabilityQuery.data;
  const mysqlCapabilityScopeMatches = Boolean(mysqlCapability && mysqlCapability.connector_id === mysqlConnectorId &&
    mysqlCapability.connector_revision === mysqlConnectorRevision);
  const mysqlMetadataVerified = Boolean(mysqlCapabilityScopeMatches && mysqlCapability?.connected &&
    mysqlCapability.metadata_status === 'verified');
  const mysqlServerVerified = Boolean(mysqlMetadataVerified && mysqlCapability?.transactional_engine);
  const mysqlManagedLongAvailable = Boolean(mysqlServerVerified && mysqlCapability?.supports_managed_long);
  const mysqlNativeMonthlyAvailable = Boolean(mysqlManagedLongAvailable && mysqlCapability?.supports_native_partition);
  const mysqlTestWriteAvailable = Boolean(mysqlMetadataVerified && mysqlCapability?.supports_group_test_write);

  const sqlserverConnectorId = state.db.connector.connector_id ?? '';
  const sqlserverConnectorRevision = state.db.connector.identity_revision ?? '';
  const sqlserverCapability = !destination || sqlserverCapabilityQuery.isError || sqlserverCapabilityQuery.isFetching
    ? undefined : sqlserverCapabilityQuery.data;
  const sqlserverCapabilityScopeMatches = Boolean(sqlserverCapability && sqlserverCapability.connector_id === sqlserverConnectorId &&
    sqlserverCapability.connector_revision === sqlserverConnectorRevision);
  const sqlserverMetadataVerified = Boolean(sqlserverCapabilityScopeMatches && sqlserverCapability?.connected &&
    sqlserverCapability.metadata_status === 'verified');
  const sqlserverManagedLongAvailable = Boolean(sqlserverMetadataVerified && sqlserverCapability?.supports_managed_long);
  const sqlserverNativeMonthlyAvailable = Boolean(sqlserverManagedLongAvailable && sqlserverCapability?.supports_native_partition);
  const sqlserverTestWriteAvailable = Boolean(sqlserverMetadataVerified && sqlserverCapability?.supports_group_test_write);

  const nativeMonthlyAvailable = (isMySQL && mysqlNativeMonthlyAvailable) || (isSQLServer && sqlserverNativeMonthlyAvailable);

  const { candidates: allCandidates } = React.useMemo(() => buildGroupCandidates(state), [state]);
  const deviceCandidates = React.useMemo(() => allCandidates.filter((c) => c.device_id === activeDeviceId), [allCandidates, activeDeviceId]);
  const liveGroups = React.useMemo(() => (groupsQuery.data?.groups ?? []).filter((g) => g.status !== 'deleted'), [groupsQuery.data?.groups]);
  const savedLongGroups = React.useMemo(() => liveGroups.filter((g) => isManagedProfileGroupForDevice(g, activeDeviceId)), [liveGroups, activeDeviceId]);
  const advancedGroups = React.useMemo(() => liveGroups.filter((g) => !isManagedProfileGroupForDevice(g, activeDeviceId)), [liveGroups, activeDeviceId]);

  const savedCircuits = React.useMemo(() => extractDeviceCircuits(savedLongGroups), [savedLongGroups]);
  // The batch tools write long-layout groups only. A wide profile circuit is prepared one circuit at a time, so it is
  // neither offered to nor reachable from a batch.
  const batchGroups = React.useMemo(() => savedLongGroups.filter((group) => group.layout_version === 'managed_long_v1'), [savedLongGroups]);
  const batchCircuits = React.useMemo(() => extractDeviceCircuits(batchGroups), [batchGroups]);
  const [selectedCircuitId, setSelectedCircuitId] = React.useState<string | undefined>(() => {
    if (!workspaceId || !activeDeviceId) return undefined;
    return loadActiveCircuitId(workspaceId, activeDeviceId);
  });
  const activeCircuitId = selectedCircuitId !== undefined ? selectedCircuitId : savedCircuits[0]?.circuit_id;

  const [localDraft, setLocalDraft] = React.useState(() => loadCircuitProfileLocalDraft(workspaceId ?? '', activeDeviceId, deviceCandidates, savedLongGroups, activeCircuitId, prepareOnly));
  const nativeEdits = React.useRef(new Map<string, number>());
  const [preparedGroups, setPreparedGroups] = React.useState<WriteGroup[]>([]);
  const [revisionOverride, setRevisionOverride] = React.useState<{ base: string; revision: string }>();
  const [errorMessage, setErrorMessage] = React.useState<string>();
  const [batchOpen, setBatchOpen] = React.useState(false);

  React.useEffect(() => {
    if (!workspaceId || !activeDeviceId) return;
    const stored = loadActiveCircuitId(workspaceId, activeDeviceId);
    setSelectedCircuitId(stored);
  }, [workspaceId, activeDeviceId]);

  const handleSelectCircuit = (targetId: string) => {
    onCircuitRequestDone?.(); // choosing a circuit by hand retires a pending routing request
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

  React.useEffect(() => {
    if (!circuitRequest || !workspaceId || circuitRequest.workspaceId !== workspaceId || circuitRequest.deviceId !== activeDeviceId) return;
    setSelectedCircuitId(circuitRequest.circuitId);
    saveActiveCircuitId(workspaceId, activeDeviceId, circuitRequest.circuitId);
    onCircuitRequestDone?.();
  }, [circuitRequest, workspaceId, activeDeviceId, onCircuitRequestDone]);
  const updateLocalDraft = React.useCallback(
    (updater: (prev: typeof localDraft) => typeof localDraft) => {
      setLocalDraft((prev) => {
        const next = updater(prev);
        const key = JSON.stringify([workspaceId ?? '', activeDeviceId, next.circuit_id ?? prev.circuit_id ?? '']); nativeEdits.current.set(key, (nativeEdits.current.get(key) ?? 0) + 1);
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

  const instParsed = parseStrictProfileIntervalInput('instantaneous', localDraft.instantaneous_interval_input);
  const cumParsed = parseStrictProfileIntervalInput('cumulative', localDraft.cumulative_interval_input);
  const circuitId = deriveCircuitId(localDraft.circuit_name, localDraft.circuit_id);
  const currentGroups = savedLongGroups.filter((group) => groupCircuitId(group) === circuitId);
  const storageLocked = currentGroups.length > 0;
  const effectiveStorageMode = storageLocked ? (currentGroups[0]?.storage_mode ?? 'plain') : localDraft.storage_mode;
  // A saved circuit keeps the layout it was created with; only a new one lets the operator choose the purpose.
  const savedPurpose = purposeOfLayout(currentGroups[0]?.layout_version);
  const purpose = savedPurpose ?? localDraft.purpose;
  const capabilities = { mysql: mysqlCapabilityQuery, sqlserver: sqlserverCapabilityQuery };
  const capabilityContext = { connectorKind: kindStr, destination, capabilities };
  const layoutBlockedReason = circuitProfileLayoutBlockReason(capabilityContext, purpose, effectiveStorageMode);
  const layoutOptions = {
    wide: circuitProfileLayoutBlockReason(capabilityContext, 'wide', effectiveStorageMode),
    long: circuitProfileLayoutBlockReason(capabilityContext, 'long', 'plain'),
  };
  // A saved member that no current source matches is an exception to repair in its group, never to hide or drop.
  const lostSources = new Map<string, string>();
  const availableSources = new Set(deviceCandidates.map((c) => `${c.device_id}\0${c.point_id}\0${c.tag_id}`));
  for (const group of currentGroups) {
    for (const member of group.members) {
      const key = `${member.device_id}\0${member.point_id}\0${member.tag_id}`;
      if (!availableSources.has(key) && !lostSources.has(key)) lostSources.set(key, group.id);
    }
  }
  const unavailableSources = { count: lostSources.size, groupId: lostSources.values().next().value as string | undefined };
  // A wide circuit states how its rows treat a missing source; the saved required flags are read, never rewritten.
  const wideGroups = currentGroups.filter((group) => purposeOfLayout(group.layout_version) === 'wide');
  const requiredPolicy = wideGroups.length > 0
    ? {
        count: wideGroups.reduce((total, group) => total + group.members.filter((member) => member.required).length, 0),
        groupId: (wideGroups.find((group) => group.members.some((member) => member.required)) ?? wideGroups[0]).id,
      }
    : undefined;
  const savedInstantaneous = currentGroups.find((group) => groupProfileKind(group) === 'instantaneous');
  // A new instantaneous profile is an explicit mean; a saved one keeps the reduction it already has (none = latest).
  const instantaneousReduction = savedInstantaneous ? (savedInstantaneous.recording_profile?.sample_reduction ?? 'latest') : 'mean';
  const savedRoleKeys = React.useMemo(
    () => new Set((savedInstantaneous?.members ?? []).map((m) => `${m.device_id}\0${m.point_id}\0${m.tag_id}`)),
    [savedInstantaneous],
  );
  const roleOf = React.useCallback((candidate: typeof deviceCandidates[number]) => localDraft.member_roles[candidateKey(candidate)] ?? defaultMemberRole(candidate), [localDraft.member_roles]);
  // A counter is whatever a saved, confirmed definition says: cumulative sources and counter_raw roles are checked alike.
  const counterTargets = React.useMemo(() => {
    const targets = new Map<string, typeof deviceCandidates[number]>();
    if (localDraft.cumulative_enabled) for (const c of cumulativeCandidates) targets.set(candidateKey(c), c);
    if (localDraft.instantaneous_enabled) {
      for (const c of instantaneousCandidates) if (roleOf(c) === 'counter_raw') targets.set(candidateKey(c), c);
    }
    return [...targets.values()];
  }, [localDraft.cumulative_enabled, localDraft.instantaneous_enabled, cumulativeCandidates, instantaneousCandidates, roleOf]);
  const counterStatuses = React.useMemo(() => counterTargets.map((c) => evaluateCandidateCounterConfirmation(workspaceId ?? '', c, measurements)), [workspaceId, counterTargets, measurements]);
  const unconfirmedCounterStatuses = React.useMemo(() => counterStatuses.filter((s) => !s.confirmed), [counterStatuses]);
  const confirmedMeasurementsByPointKey = React.useMemo(() => {
    const map = new Map<string, MeasurementDefinition>();
    for (const status of counterStatuses) if (status.confirmed && status.measurement) map.set(candidateKey(status.candidate), status.measurement);
    return map;
  }, [counterStatuses]);

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
  const batch = useCircuitProfileBatchHost({ scopeKey, workspaceId, nativeEdits: nativeEdits.current, nativeConsent: JSON.stringify([localDraft, deviceCandidates, settingsRevision, shareWorkspaceRevision]), deviceId: activeDeviceId, candidates: deviceCandidates, savedLongGroups: batchGroups,
    activeCircuitId, queriedRevision, refetchGroups: groupsQuery.refetch, setRevisionOverride, setLocalDraft });
  const batchBusy = batch.busy;
  const commandFence = useCircuitRecordingCommandFence({ ownerKey: scopeKey, scopeKey, workspaceId: workspaceId ?? '',
    deviceId: activeDeviceId, circuitId, consentKey: JSON.stringify([localDraft, deviceCandidates, settingsRevision, shareWorkspaceRevision]), batchGroups, updateWriteGroup: vars => updateGroup.mutateAsync(vars),
    counterBlocked: batchOpen || batchBusy, counterRepair: { available: !readonly && !measurementsQuery.isLoading && !measurementsQuery.isFetching && !measurementsQuery.isError, statuses: unconfirmedCounterStatuses },
    ensureCircuitProfile: vars => ensureProfile.mutateAsync(vars) });
  const hydrationScope = JSON.stringify([workspaceId, activeDeviceId, activeCircuitId, deviceCandidates, destination, prepareOnly]); const hydratedScope = React.useRef(hydrationScope);
  React.useEffect(() => {
    if (batchBusy && hydratedScope.current === hydrationScope) return; // Hold only own same-context canonical reseeds, not actual viewing-scope changes.
    hydratedScope.current = hydrationScope; setLocalDraft(loadCircuitProfileLocalDraft(workspaceId ?? '', activeDeviceId, deviceCandidates, savedLongGroups, activeCircuitId, prepareOnly));
    setPreparedGroups([]);
  }, [workspaceId, activeDeviceId, deviceCandidates, savedLongGroups, activeCircuitId, prepareOnly, batchBusy, hydrationScope]);
  React.useEffect(() => { setPreparedGroups([]); setRevisionOverride(undefined); setErrorMessage(undefined); setBatchOpen(false); }, [scopeKey]);
  const counterDefinitionsReady = !measurementsQuery.isLoading && !measurementsQuery.isFetching && !measurementsQuery.isError && Array.isArray(measurementsQuery.data);
  const requiresCounterConfirmation = counterTargets.length > 0;
  const counterNameMissing = counterDefinitionsReady && requiresCounterConfirmation && !localDraft.circuit_name.trim() &&
    unconfirmedCounterStatuses.some((status) => !status.measurement && status.issueCode !== 'counter-type-unsupported');
  const instantaneousSourcesMissing = localDraft.instantaneous_enabled && instantaneousCandidates.length === 0;
  const cumulativeSourcesMissing = localDraft.cumulative_enabled && cumulativeCandidates.length === 0;
  const hasValidIntervals = (!localDraft.instantaneous_enabled || instParsed.valid) && (!localDraft.cumulative_enabled || cumParsed.valid);
  const cumulativeCircuitMissing = localDraft.cumulative_enabled && !circuitId;
  const counterConfirmationBlocked = requiresCounterConfirmation && (!counterDefinitionsReady || unconfirmedCounterStatuses.length > 0);
  const canonicalListReady = Boolean(workspaceId && groupsQuery.data?.workspace_id === workspaceId && groupsQuery.data?.workspace_revision &&
    !groupsQuery.isLoading && !groupsQuery.isFetching && !groupsQuery.isError);
  const canPrepareOrStart = canonicalListReady && !readonly && !batchBusy && !batchOpen &&
    Boolean(workspaceId && effectiveWorkspaceRevision && activeDeviceId && destination) &&
    enabledProfiles.length > 0 && hasValidIntervals && !instantaneousSourcesMissing && !cumulativeSourcesMissing &&
    !cumulativeCircuitMissing && !counterNameMissing && !counterConfirmationBlocked && !layoutBlockedReason;

  const prepareAllowed = React.useRef(canPrepareOrStart); prepareAllowed.current = canPrepareOrStart;
  const ensureCurrentGroups = async (): Promise<{ groups: WriteGroup[]; workspaceRevision: string } | undefined> => {
    if (workspaceBlocked() || !prepareAllowed.current || activeScope.current !== scopeKey || !canPrepareOrStart || !workspaceId || !destination) return undefined;
    const requestedScope = scopeKey;
    setErrorMessage(undefined);
    try {
      const result = await ensureOrUpdateCircuitProfileGroups({
        workspaceId, initialWorkspaceRevision: effectiveWorkspaceRevision, deviceName: activeDevice?.name ?? activeDeviceId,
        enabledProfiles, storageMode: effectiveStorageMode, recreateMissing: localDraft.recreate_missing, purpose, memberRoles: localDraft.member_roles, newReduction: instantaneousReduction,
        circuitId, circuitName: localDraft.circuit_name, destination, connectorKind: state.db.connector.kind, capabilities,
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
    prepare: ensureCurrentGroups, onActivated, prepareOnly, retryBlocked: () => workspaceBlocked() || batchBusy || batchOpen || Boolean(layoutBlockedReason) });
  useReportOperationActive(flow.busy || flow.unresolved || batchBusy, onOperationActiveChange); // one merged page signal
  const isEffectiveLocked = readonly || flow.busy || flow.unresolved || batchBusy;
  const schemaGroups = preparedGroups.length > 0 ? preparedGroups : currentGroups;
  const { drafts: recordingDrafts, dirtyGroupIds } = circuitPreparationDrafts(schemaGroups, enabledProfiles, {
    workspaceId: workspaceId ?? '', deviceName: activeDevice?.name ?? activeDeviceId,
    storageMode: effectiveStorageMode, recreateMissing: localDraft.recreate_missing, circuitId, circuitName: localDraft.circuit_name,
    destination, measurementsByPointKey: confirmedMeasurementsByPointKey,
    purpose, reduction: instantaneousReduction, memberRoles: localDraft.member_roles }, state.db.connector.kind);
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
      {!prepareOnly && legacyGroup && <ExistingDeviceRecording lifecycle={{ workspaceId, workspaceRevision: effectiveWorkspaceRevision, lifecycleHost, readonly }} group={legacyGroup} disabled={flow.busy || flow.unresolved || batchBusy} onEdit={onEditLegacy} onOpenAdvanced={onOpenAdvanced} />}
      <div data-testid="circuit-profile-destination-summary" className="rounded border border-slate-800/80 bg-slate-950/30 p-3 text-xs text-slate-300">
        <p className="font-semibold text-slate-200">{t('step4.basic.circuit.step_destination')}</p>
        <p className="mt-1 text-slate-400">{t('step4.basic.destination', { name: connectorAlias })} · {connectorAlias} ({state.db.connector.kind})</p>
      </div>
      {isMySQL && <MySQLRecordingCapabilityPanel
        capability={mysqlCapability}
        isLoading={mysqlCapabilityQuery.isLoading || mysqlCapabilityQuery.isFetching}
        isError={mysqlCapabilityQuery.isError}
        onRetry={() => void mysqlCapabilityQuery.refetch()}
      />}
      {isSQLServer && <SQLServerRecordingCapabilityPanel
        capability={sqlserverCapability}
        isLoading={sqlserverCapabilityQuery.isLoading || sqlserverCapabilityQuery.isFetching}
        isError={sqlserverCapabilityQuery.isError}
        onRetry={() => void sqlserverCapabilityQuery.refetch()}
      />}
      <CircuitProfileConfigFields newOperationBlocked={commandFence.counterCommandBlocked} isNewOperationBlocked={commandFence.isCounterCommandBlocked}
        workspaceId={workspaceId} workspaceRevision={effectiveWorkspaceRevision} readonly={isEffectiveLocked}
        savedDevices={savedDevices} activeDeviceId={activeDeviceId} onSelectDevice={selectDevice}
        deviceCandidates={deviceCandidates} localDraft={localDraft} updateLocalDraft={updateLocalDraft}
        instParsed={instParsed} cumParsed={cumParsed} isPostgres={isPostgres} isMySQL={isMySQL} isSQLServer={isSQLServer}
        nativeMonthlyAvailable={nativeMonthlyAvailable} storageLocked={storageLocked}
        effectiveStorageMode={effectiveStorageMode} cumulativeCircuitMissing={cumulativeCircuitMissing} counterNameMissing={counterNameMissing} circuitId={circuitId}
        unconfirmedCounterStatuses={counterDefinitionsReady ? unconfirmedCounterStatuses : []} onCounterConfirmed={async () => { await measurementsQuery.refetch(); }}
        purpose={purpose} purposeLocked={Boolean(savedPurpose)} layoutOptions={layoutOptions} layoutBlockedReason={layoutBlockedReason} showPurpose={true} unavailableSources={unavailableSources} requiredPolicy={requiredPolicy} onOpenAdvanced={onOpenAdvanced} instantaneousReduction={instantaneousReduction} savedRoleKeys={savedRoleKeys} memberRole={roleOf}
        savedCircuits={savedCircuits} onSelectCircuit={handleSelectCircuit} onOpenBatch={batchCircuits.length > 1 && purpose !== 'wide' ? () => setBatchOpen(true) : undefined}
      />
      {batchOpen && (
        <CircuitProfileBatchEditor
          workspaceId={workspaceId ?? ''} workspaceRevision={effectiveWorkspaceRevision} scopeKey={scopeKey}
          connectorRevision={destination?.connector_revision} capabilityContext={capabilityContext} savedLongGroups={batchGroups}
          readonly={isEffectiveLocked} onClose={() => setBatchOpen(false)}
          onSuccess={batch.onSaved} onVerify={batch.reload}
          updateWriteGroup={commandFence.updateWriteGroup}
          ensureCircuitProfile={commandFence.ensureCircuitProfile}
          onBusyChange={batch.setHeld}
        />
      )}
      {requiresCounterConfirmation && !counterDefinitionsReady && <p role="status" data-testid={measurementsQuery.isError ? 'circuit-profile-measurements-error' : 'circuit-profile-measurements-loading'} className="text-xs text-amber-300">
        {t(measurementsQuery.isError ? 'step4.basic.circuit.measurements_error' : 'step4.basic.circuit.measurements_loading')}
        {measurementsQuery.isError && <button type="button" disabled={measurementsQuery.isFetching}
          onClick={() => void measurementsQuery.refetch({ cancelRefetch: false })}>{t('step4.group.retry')}</button>}
      </p>}
      <CircuitRecordingPreparation onContinueToRuntime={onContinueToRuntime} deviceNames={Object.fromEntries(savedDevices.map(device => [device.id, device.name]))} groupNames={Object.fromEntries(currentGroups.map(group => [group.id, group.name]))} isNewCommandBlocked={workspaceBlocked} ownerKey={scopeKey} fenceKey={scopeKey} onCommandFence={commandFence.onCommandFence} flow={flow} groups={schemaGroups} drafts={recordingDrafts}
        batchHandoff={{ consentKey: JSON.stringify([localDraft, deviceCandidates, settingsRevision, shareWorkspaceRevision]), blocked: readonly || Boolean(layoutBlockedReason) }}
        onRepair={!destination ? focusDestinationSettings : undefined}
        dirtyGroupIds={dirtyGroupIds} counterStatuses={counterStatuses} canPrepare={canPrepareOrStart}
        readonly={readonly || batchBusy || Boolean(layoutBlockedReason)} prepareOnly={prepareOnly} summary={renderScopeSummary}
        consentKey={JSON.stringify([localDraft, recordingDrafts, settingsRevision, shareWorkspaceRevision])}
        scope={circuitRecordingScope({ workspaceId, workspaceRevision: effectiveWorkspaceRevision,
          deviceId: activeDeviceId, circuitId, groups: currentGroups, draft: localDraft, purpose,
          storageMode: effectiveStorageMode, reduction: instantaneousReduction,
          instantaneous: instantaneousCandidates, cumulative: cumulativeCandidates })}
        schema={{ workspaceId, workspaceRevision: effectiveWorkspaceRevision, preparedGroups: schemaGroups, dirtyGroupIds,
          disabled: !canPrepareOrStart || flow.unresolved, readonly: isEffectiveLocked || Boolean(layoutBlockedReason),
          isPreparing: flow.busy || batchBusy || ensureProfile.isPending || updateGroup.isPending,
          destinationMatches: Boolean(destination), onPrepare: () => void flow.prepare(), onApplied: () => void refreshSchema(),
          onRepairDestination: focusDestinationSettings, onOpenAdvanced, isMySQL: isMySQL && !prepareOnly,
          mysqlTestWriteAvailable, isSQLServer: isSQLServer && !prepareOnly, sqlserverTestWriteAvailable }} />
      {advancedGroups.length > 0 && <p data-testid="basic-recording-advanced-note" className="text-xs text-slate-400">{t('step4.basic.advanced_note', { count: advancedGroups.length })}</p>}
      {(errorMessage || flow.error) && <p role="alert" data-testid="circuit-profile-error" className="text-xs text-rose-300">{errorMessage || flow.error}</p>}
      <CircuitProfileStatusList lifecycle={{ workspaceId, workspaceRevision: effectiveWorkspaceRevision, lifecycleHost, readonly: isEffectiveLocked }} groups={currentGroups} candidates={deviceCandidates} onOpenAdvanced={onOpenAdvanced} draftEnabled={{ instantaneous: localDraft.instantaneous_enabled, cumulative: localDraft.cumulative_enabled }} /><CircuitProfileRoles onOpenAdvanced={onOpenAdvanced} />
      <CircuitProfileBasicActions guided flow={flow} canStart={canPrepareOrStart} readonly={readonly || Boolean(layoutBlockedReason)} batchBusy={batchBusy} prepareOnly={prepareOnly} onClose={onClose} onContinueToRuntime={onContinueToRuntime} onOpenAdvanced={onOpenAdvanced} />
    </section>
  );
};
