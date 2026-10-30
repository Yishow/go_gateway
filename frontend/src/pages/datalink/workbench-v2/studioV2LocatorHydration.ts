import { useLayoutEffect, useRef, useState, type MutableRefObject } from 'react';
import type { StudioLocatorEntryRead } from '../../../types/studioV2LocatorEntry';
import { INITIAL_STATE, type useWorkbenchV2State } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import { hydrateStudioV2Device } from '../../../features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { hydrateStudioV2Rule } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
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
function draftIdentity(state: WorkbenchV2State): string {
  return JSON.stringify([state.devices, state.rules, state.points, state.mappings, state.db, state.settings]);
}
/** A new viewing request cannot replace local edits or the already hydrated normal workspace. */
export function useStudioV2LocatorHydration(actions: ReturnType<typeof useWorkbenchV2State>, stateRef: MutableRefObject<WorkbenchV2State>,
  nativeHydrated: boolean, entry?: StudioLocatorEntryRead, editorMounted = false) {
  const accepted = useRef<{ entry: StudioLocatorEntryRead; baseline: string } | undefined>(undefined);
  const attempted = useRef<{ entry: StudioLocatorEntryRead; editorMounted: boolean } | undefined>(undefined);
  const [result, setResult] = useState<{ entry?: StudioLocatorEntryRead; blocked: boolean }>({ blocked: false });
  useLayoutEffect(() => {
    if (!entry || entry.selection.state !== 'selected' || attempted.current?.entry === entry && attempted.current.editorMounted === editorMounted) return;
    attempted.current = { entry, editorMounted };
    const dirty = !nativeHydrated && accepted.current && accepted.current.baseline !== draftIdentity(stateRef.current);
    // Preserve the live Basic subtree when a query-only route would leave Step 4.
    if (dirty || editorMounted && accepted.current?.entry !== entry || nativeHydrated && (!studioLocatorStateMatches(entry, stateRef.current) || stateRef.current.current === 4 && entry.selection.target.step !== 4)) {
      setResult({ entry, blocked: true }); return;
    }
    const next = nativeHydrated ? stateRef.current : studioV2LocatorHydration(entry);
    if (!next) { setResult({ entry, blocked: true }); return; }
    if (!nativeHydrated) { stateRef.current = next; actions.dispatch({ type: 'SET_STATE', payload: next }); }
    accepted.current = { entry, baseline: draftIdentity(next) };
    setResult({ entry, blocked: false });
  }, [entry, nativeHydrated, actions, stateRef, editorMounted]);
  return { lastEntry: accepted.current?.entry, ready: Boolean(entry && result.entry === entry && !result.blocked), blocked: Boolean(entry && result.entry === entry && result.blocked) };
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
