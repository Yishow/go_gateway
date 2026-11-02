import { useLayoutEffect, useRef, useState, type MutableRefObject } from 'react';
import type { StudioLocatorEntryRead } from '../../../types/studioV2LocatorEntry';
import { INITIAL_STATE, type useWorkbenchV2State } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Device, Rule, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import type { StudioV2WorkspaceDeviceRecord } from '../../../services/studioV2WorkspaceDevices';
import { hydrateStudioV2Device, toStudioV2DeviceUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import type { SourceRuleRecord } from '../../../types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '../../../types/studioV2RuntimeApply';
import { hydrateStudioV2Rule, toStudioV2RuleUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { hydrateStudioV2Mapping } from '../../../features/datalink/workbench-v2/state/studioV2MappingAutosave';
import { hydrateStudioV2DatabaseConnector, hydrateStudioV2DatabaseRowGroups, hydrateStudioV2DatabaseTarget } from '../../../features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import { deriveAllPoints, derivePoints, projectPersistedMCPoint } from '../../../features/datalink/workbench-v2/state/sourceRule';

/** Partial records never enter canonical React Query caches or imply workspace completeness. */
export function studioV2LocatorHydration(read: StudioLocatorEntryRead): WorkbenchV2State | null {
  if (read.selection.state !== 'selected') return null;
  const { target } = read.selection;
  const devices = (read.entry.devices ?? []).map(hydrateStudioV2Device);
  const rules = (read.entry.source_rules ?? []).map(hydrateStudioV2Rule);
  const points = target.mappingId || read.entry.editor_hydrated ? (read.entry.mappings ?? []).flatMap(record => {
    const rule = rules.find(value => value.id === record.rule_id && value.device_id === record.device_id);
    const device = devices.find(value => value.id === record.device_id);
    if (!rule || !device || (target.mappingId !== undefined && record.id !== target.mappingId)) return [];
    // The address supplies display geometry, never identity lookup. The row key is the exact native Point ID.
    const geometry = device.protocol === 'mc_3e' && rule.persisted
      ? projectPersistedMCPoint(rule, { point_id: record.point_id, address: record.address })
      : derivePoints({ ...rule, start_address: record.address, count: 1 }, device.id, new Set(), device.protocol)[0];
    return geometry ? [{ ...geometry, id: record.point_id, name: record.display_name }] : [];
  }) : deriveAllPoints(rules, Object.fromEntries(devices.map(device => [device.id, device.protocol])));
  const mappings = Object.fromEntries((read.entry.mappings ?? []).flatMap(record => {
    const point = points.find(value => value.id === record.point_id && value.rule_id === record.rule_id && value.device_id === record.device_id);
    return point ? [[point.id, hydrateStudioV2Mapping(point, record)]] : [];
  }));
  if (target.mappingId && !Object.values(mappings).some(mapping => mapping.mapping_id === target.mappingId)) return null;
  const config = read.entry.editor_hydrated ? read.entry.database_config : undefined;
  const db = { ...INITIAL_STATE.db, targets: Object.fromEntries((read.entry.database_targets ?? []).map(record =>
    [record.point_id, hydrateStudioV2DatabaseTarget(record)])), row_groups: config ? hydrateStudioV2DatabaseRowGroups(config.row_groups) : [],
    connector: config ? hydrateStudioV2DatabaseConnector(config) : { ...INITIAL_STATE.db.connector, name: '', host: '', database: '', username: '', schema: '', table: '',
      status: 'unknown', persisted: false, save_state: 'idle' as const } };
  return { ...INITIAL_STATE, current: target.step, devices, rules, points, mappings, selectedRuleId: target.ruleId ?? null,
    completed: new Set(), showSummaryRail: false, db,
    settings: { ...INITIAL_STATE.settings, connectors: [], modbus_share: { ...INITIAL_STATE.settings.modbus_share, enabled: false } } };
}
// Same canonical comparison used by the existing source-submission observer;
// object key order in a server JSON connection config is not a user edit.
function comparable(value: unknown): string {
  if (Array.isArray(value)) return `[${value.map(comparable).join(',')}]`;
  if (value && typeof value === 'object') return `{${Object.entries(value).sort(([a], [b]) => a.localeCompare(b)).map(([key, item]) => `${JSON.stringify(key)}:${comparable(item)}`).join(',')}}`;
  return JSON.stringify(value);
}
function draftIdentity(state: WorkbenchV2State): string {
  return JSON.stringify([state.devices, state.rules, state.points, state.mappings, state.db, state.settings]);
}
/** A new viewing request cannot replace local edits or the already hydrated normal workspace. */
export function useStudioV2LocatorHydration(actions: ReturnType<typeof useWorkbenchV2State>, stateRef: MutableRefObject<WorkbenchV2State>,
  nativeHydrated: boolean, entry?: StudioLocatorEntryRead, editorMounted = false) {
  const accepted = useRef<{ entry: StudioLocatorEntryRead; baseline: WorkbenchV2State } | undefined>(undefined);
  const mode = useRef({ nativeHydrated, entry, epoch: 0 });
  // Returning to the retained editor object must not revive an earlier submission.
  if (mode.current.entry !== entry || mode.current.nativeHydrated !== nativeHydrated) {
    mode.current = { nativeHydrated, entry, epoch: mode.current.epoch + 1 };
  }
  const attempted = useRef<{ entry: StudioLocatorEntryRead; editorMounted: boolean } | undefined>(undefined);
  const [result, setResult] = useState<{ entry?: StudioLocatorEntryRead; blocked: boolean }>({ blocked: false });
  useLayoutEffect(() => {
    if (!entry || entry.selection.state !== 'selected' || attempted.current?.entry === entry && attempted.current.editorMounted === editorMounted) return;
    const sameAcceptedEntry = !nativeHydrated && accepted.current?.entry === entry;
    attempted.current = { entry, editorMounted };
    // Opening/closing this accepted editor is not a new viewing request. Keep
    // its live state and dirty baseline; a different entry still uses every guard.
    if (sameAcceptedEntry) { setResult({ entry, blocked: false }); return; }
    const dirty = !nativeHydrated && accepted.current && draftIdentity(accepted.current.baseline) !== draftIdentity(stateRef.current);
    // Preserve the live Basic subtree when a query-only route would leave Step 4.
    if (dirty || editorMounted && accepted.current?.entry !== entry || nativeHydrated && (!studioLocatorStateMatches(entry, stateRef.current) || stateRef.current.current === 4 && entry.selection.target.step !== 4)) {
      setResult({ entry, blocked: true }); return;
    }
    const next = nativeHydrated ? stateRef.current : studioV2LocatorHydration(entry);
    if (!next) { setResult({ entry, blocked: true }); return; }
    if (!nativeHydrated) { stateRef.current = next; actions.dispatch({ type: 'SET_STATE', payload: next }); }
    accepted.current = { entry, baseline: structuredClone(next) };
    setResult({ entry, blocked: false });
  }, [entry, nativeHydrated, actions, stateRef, editorMounted]);
  const captureDeviceSave = (submitted: Device, workspaceId?: string) => {
    const owner = accepted.current, target = owner?.entry.selection, viewing = mode.current.entry, selected = viewing?.selection;
    if (!owner || mode.current.nativeHydrated || !workspaceId || viewing?.workspace.id !== workspaceId || viewing.locator.workspace_id !== workspaceId || owner.entry.workspace.id !== workspaceId || selected?.state !== 'selected' || selected.target.workspaceId !== workspaceId ||
        target?.state !== 'selected' || target.target.step !== 1 || target.target.workspaceId !== workspaceId ||
        selected.target.step !== 1 || selected.target.deviceId !== submitted.id || selected.target.workspaceRevision !== target.target.workspaceRevision ||
        target.target.deviceId !== submitted.id || submitted.persisted !== true || owner.baseline.devices.filter(device => device.id === submitted.id && device.persisted).length !== 1) return;
    const captured = structuredClone(submitted), viewingEpoch = mode.current.epoch;
    const submittedIdentity = comparable(toStudioV2DeviceUpdateRequest(captured));
    let consumed = false;
    return (saved: StudioV2WorkspaceDeviceRecord) => {
      if (consumed) return; consumed = true;
      if (accepted.current !== owner || mode.current.nativeHydrated || mode.current.entry !== viewing || mode.current.epoch !== viewingEpoch || saved.id !== captured.id) return;
      const matches = stateRef.current.devices.filter(device => device.id === captured.id);
      const latest = matches.length === 1 ? matches[0] : undefined;
      if (!latest || latest.save_state !== 'saved' || !latest.persisted || saved.runtime_apply_status === 'apply_failed') return;
      try {
        if (comparable(toStudioV2DeviceUpdateRequest(latest)) !== submittedIdentity ||
            comparable(toStudioV2DeviceUpdateRequest({ ...captured, name: saved.name, description: saved.description, protocol: saved.protocol, config: JSON.parse(saved.connection_config) })) !== submittedIdentity) return;
      } catch { return; }
      // Only this ACK's unchanged persisted fields and its own save metadata move
      // forward. Tests/other devices/rules/points/mappings/settings remain dirty.
      owner.baseline = { ...owner.baseline, devices: owner.baseline.devices.map(device => device.id !== captured.id ? device : {
        ...device, name: captured.name, description: captured.description, protocol: captured.protocol, config: structuredClone(captured.config),
        persisted: latest.persisted, status: latest.status, save_state: latest.save_state, save_error: latest.save_error,
        runtime_apply_status: latest.runtime_apply_status, runtime_apply_message: latest.runtime_apply_message,
        availability_status: latest.availability_status, availability_reason: latest.availability_reason, running: latest.running,
      }) };
    };
  };
  const captureRuleSave = (submitted: Rule, workspaceId?: string) => {
    const owner = accepted.current, target = owner?.entry.selection, viewing = mode.current.entry, selected = viewing?.selection;
    if (!owner || mode.current.nativeHydrated || !workspaceId || submitted.workspace_id !== workspaceId ||
        viewing?.workspace.id !== workspaceId || viewing.locator.workspace_id !== workspaceId || owner.entry.workspace.id !== workspaceId ||
        target?.state !== 'selected' || selected?.state !== 'selected' || target.target.step !== 2 || selected.target.step !== 2 ||
        target.target.workspaceId !== workspaceId || selected.target.workspaceId !== workspaceId || selected.target.workspaceRevision !== target.target.workspaceRevision ||
        target.target.deviceId !== submitted.device_id || selected.target.deviceId !== submitted.device_id ||
        target.target.ruleId !== submitted.id || selected.target.ruleId !== submitted.id || submitted.persisted !== true ||
        owner.baseline.rules.filter(rule => rule.id === submitted.id && rule.device_id === submitted.device_id && rule.persisted).length !== 1 ||
        owner.baseline.rules.find(rule => rule.id === submitted.id)?.revision_id !== submitted.revision_id) return;
    const captured = structuredClone(submitted), viewingEpoch = mode.current.epoch, fields = toStudioV2RuleUpdateRequest(captured), identity = comparable(fields);
    let consumed = false;
    return (saved: StudioV2RuntimeAppliedRecord<SourceRuleRecord>) => {
      if (consumed) return; consumed = true;
      if (accepted.current !== owner || mode.current.nativeHydrated || mode.current.entry !== viewing || mode.current.epoch !== viewingEpoch ||
          saved.id !== captured.id || saved.device_id !== captured.device_id || saved.workspace_id !== workspaceId ||
          !saved.revision_id?.trim() || saved.runtime_apply_status === 'apply_failed') return;
      const matches = stateRef.current.rules.filter(rule => rule.id === captured.id && rule.device_id === captured.device_id);
      const latest = matches.length === 1 ? matches[0] : undefined;
      if (!latest || latest.save_state !== 'saved' || !latest.persisted || latest.revision_id !== saved.revision_id ||
          comparable(toStudioV2RuleUpdateRequest(latest)) !== identity || comparable(toStudioV2RuleUpdateRequest(hydrateStudioV2Rule(saved))) !== identity) return;
      // Only this unchanged, acknowledged rule advances. Local name, MC evidence,
      // every other rule/device/point/mapping/database/settings remain untouched.
      owner.baseline = { ...owner.baseline, rules: owner.baseline.rules.map(rule => rule.id !== captured.id ? rule : {
        ...rule, ...structuredClone(fields), data_format: captured.data_format, workspace_id: saved.workspace_id, revision_id: saved.revision_id,
        persisted: true, save_state: latest.save_state, save_error: latest.save_error,
      }) };
    };
  };
  const retainedTarget = accepted.current?.entry.selection;
  const baselineRules = retainedTarget?.state === 'selected' && retainedTarget.target.ruleId
    ? accepted.current?.baseline.rules.filter(rule => rule.id === retainedTarget.target.ruleId && rule.device_id === retainedTarget.target.deviceId && rule.workspace_id === retainedTarget.target.workspaceId) : undefined;
  // This revision can advance only through the existing checked own-ACK baseline.
  const retainedRuleRevision = baselineRules?.length === 1 && baselineRules[0].revision_id?.trim()
    ? baselineRules[0].revision_id : retainedTarget?.state === 'selected' ? retainedTarget.target.ruleRevision : undefined;
  return { captureDeviceSave, captureRuleSave, retainedRuleRevision, canLeave: Boolean(accepted.current && draftIdentity(accepted.current.baseline) === draftIdentity(stateRef.current)), lastEntry: accepted.current?.entry, ready: Boolean(entry && result.entry === entry && !result.blocked), blocked: Boolean(entry && result.entry === entry && result.blocked) };
}

/** Only exact hydrated canonical bindings can receive a hot locator; drafts remain separate. */
export function studioLocatorStateMatches(read: StudioLocatorEntryRead, state: WorkbenchV2State): boolean {
  if (read.selection.state !== 'selected') return false;
  const target = read.selection.target;
  if (target.deviceId && !state.devices.some(value => value.id === target.deviceId && value.persisted)) return false;
  if (target.ruleId && !state.rules.some(value => value.id === target.ruleId && value.device_id === target.deviceId &&
    value.workspace_id === target.workspaceId && value.revision_id === target.ruleRevision)) return false;
  if (target.mappingId) {
    const native = read.entry.mappings?.find(value => value.id === target.mappingId);
    const rows = Object.values(state.mappings).filter(value => value.mapping_id === target.mappingId &&
      value.workspace_id === target.workspaceId && value.device_id === target.deviceId && value.rule_id === target.ruleId &&
      value.persisted_point_id === native?.point_id);
    if (rows.length !== 1 || !native) return false;
  }
  return true;
}
