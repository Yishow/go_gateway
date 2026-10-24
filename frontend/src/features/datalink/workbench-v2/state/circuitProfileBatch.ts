import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import {
  type CircuitProfileLocalDraft,
  groupCircuitId,
  groupCircuitName,
  groupProfileKind,
  loadCircuitProfileLocalDraft,
  saveCircuitProfileLocalDraft,
} from './circuitProfileBasic';
import { groupProfileInterval } from './circuitProfileBatchOutcome';
import type { GroupCandidate } from './writeGroup/candidates';

export interface BatchCircuitItem {
  circuit_id: string;
  circuit_name: string;
  instantaneousGroup?: WriteGroup;
  cumulativeGroup?: WriteGroup;
  currentInstantaneousInterval?: number;
  currentCumulativeInterval?: number;
  isOverrideInstantaneous: boolean;
  isOverrideCumulative: boolean;
}

export interface BatchConfig {
  applyInstantaneous: boolean;
  instantaneousInterval: number;
  applyCumulative: boolean;
  cumulativeInterval: number;
}

export interface BatchProfileChange {
  isNew: boolean;
  oldInterval?: number;
  newInterval: number;
  /** The saved writer already has the reviewed interval, so nothing is sent for it. */
  unchanged: boolean;
}

export interface BatchPreviewItem {
  circuit_id: string;
  circuit_name: string;
  selected: boolean;
  hasOverride: boolean;
  instantaneous?: BatchProfileChange;
  cumulative?: BatchProfileChange;
}

export interface DeviceCircuitOption {
  circuit_id: string;
  circuit_name: string;
  groupCount: number;
}

function activeCircuitStorageKey(workspaceId: string, deviceId: string): string {
  return `wbv2.circuit-profile-active-cid.${encodeURIComponent(`${workspaceId.trim()}/${deviceId.trim()}`)}`;
}

export function saveActiveCircuitId(workspaceId: string, deviceId: string, circuitId: string | undefined): void {
  if (typeof window === 'undefined' || !workspaceId.trim() || !deviceId.trim()) return;
  try {
    const key = activeCircuitStorageKey(workspaceId, deviceId);
    if (circuitId && circuitId.trim()) {
      window.sessionStorage.setItem(key, circuitId.trim());
    } else {
      window.sessionStorage.removeItem(key);
    }
  } catch { /* ignore */ }
}

export function loadActiveCircuitId(workspaceId: string, deviceId: string): string | undefined {
  if (typeof window === 'undefined' || !workspaceId.trim() || !deviceId.trim()) return undefined;
  try {
    return window.sessionStorage.getItem(activeCircuitStorageKey(workspaceId, deviceId)) ?? undefined;
  } catch {
    return undefined;
  }
}

export function extractDeviceCircuits(savedLongGroups: WriteGroup[]): DeviceCircuitOption[] {
  const map = new Map<string, { circuit_name: string; count: number }>();
  for (const group of savedLongGroups) {
    const cid = groupCircuitId(group);
    if (!cid) continue;
    const cname = groupCircuitName(group);
    const existing = map.get(cid);
    if (!existing) {
      map.set(cid, { circuit_name: cname, count: 1 });
    } else {
      if (!existing.circuit_name && cname) existing.circuit_name = cname;
      existing.count += 1;
    }
  }
  return Array.from(map.entries()).map(([circuit_id, item]) => ({
    circuit_id, circuit_name: item.circuit_name, groupCount: item.count,
  }));
}

export function buildBatchCircuitItems(savedLongGroups: WriteGroup[]): BatchCircuitItem[] {
  const map = new Map<string, BatchCircuitItem>();
  for (const group of savedLongGroups) {
    const cid = groupCircuitId(group);
    if (!cid) continue;
    const cname = groupCircuitName(group);
    const kind = groupProfileKind(group);
    if (!map.has(cid)) {
      map.set(cid, {
        circuit_id: cid,
        circuit_name: cname,
        isOverrideInstantaneous: false,
        isOverrideCumulative: false,
      });
    }
    const item = map.get(cid)!;
    if (!item.circuit_name && cname) item.circuit_name = cname;
    if (kind === 'instantaneous') {
      const interval = groupProfileInterval(group, kind);
      item.instantaneousGroup = group;
      item.currentInstantaneousInterval = interval;
      item.isOverrideInstantaneous = interval !== 15;
    } else if (kind === 'cumulative') {
      const interval = groupProfileInterval(group, kind);
      item.cumulativeGroup = group;
      item.currentCumulativeInterval = interval;
      item.isOverrideCumulative = interval !== 10;
    }
  }
  return Array.from(map.values());
}

export function buildBatchPreview(
  circuits: BatchCircuitItem[],
  selectedCircuitIds: Set<string>,
  config: BatchConfig,
): BatchPreviewItem[] {
  return circuits.map((circuit) => {
    const isSelected = selectedCircuitIds.has(circuit.circuit_id);
    const change = (group: WriteGroup | undefined, oldInterval: number | undefined, newInterval: number): BatchProfileChange => ({
      isNew: !group, oldInterval, newInterval, unchanged: Boolean(group) && oldInterval === newInterval,
    });
    const instantaneous = isSelected && config.applyInstantaneous
      ? change(circuit.instantaneousGroup, circuit.currentInstantaneousInterval, config.instantaneousInterval) : undefined;
    const cumulative = isSelected && config.applyCumulative
      ? change(circuit.cumulativeGroup, circuit.currentCumulativeInterval, config.cumulativeInterval) : undefined;
    const hasOverride = Boolean(
      (instantaneous && !instantaneous.unchanged && circuit.isOverrideInstantaneous) ||
      (cumulative && !cumulative.unchanged && circuit.isOverrideCumulative),
    );

    return {
      circuit_id: circuit.circuit_id,
      circuit_name: circuit.circuit_name,
      selected: isSelected,
      hasOverride,
      instantaneous,
      cumulative,
    };
  });
}

export function hasPlanGroupRevisionDrift(
  frozenTargets: Array<{ instantaneousGroup?: WriteGroup; cumulativeGroup?: WriteGroup }>,
  savedLongGroups: WriteGroup[],
): boolean {
  for (const target of frozenTargets) {
    if (target.instantaneousGroup) {
      const live = savedLongGroups.find((g) => g.id === target.instantaneousGroup!.id);
      if (live && live.revision !== target.instantaneousGroup.revision) return true;
    }
    if (target.cumulativeGroup) {
      const live = savedLongGroups.find((g) => g.id === target.cumulativeGroup!.id);
      if (live && live.revision !== target.cumulativeGroup.revision) return true;
    }
  }
  return false;
}

const memberKey = (member: { device_id: string; point_id: string; tag_id: string }) => `${member.device_id}\0${member.point_id}\0${member.tag_id}`;

/**
 * Brings a local draft in line with groups a batch saved. A group that already existed only gets its reviewed
 * period: the operator's source choices and other unsaved fields stay as they were. Only a group that did not
 * exist before is a new profile, which starts out enabled with the reviewed sources it was created with.
 */
function mergeSavedProfiles(draft: CircuitProfileLocalDraft, saved: WriteGroup[], existingIds: ReadonlySet<string>): CircuitProfileLocalDraft {
  const next = { ...draft };
  for (const group of saved) {
    const kind = groupProfileKind(group);
    const interval = group.recording_profile?.interval;
    const created = !existingIds.has(group.id);
    const keys = (group.members ?? []).map(memberKey);
    if (kind === 'instantaneous') {
      if (interval) next.instantaneous_interval_input = String(interval);
      if (created) { next.instantaneous_enabled = true; next.instantaneous_keys = keys; }
    } else if (kind === 'cumulative') {
      if (interval) next.cumulative_interval_input = String(interval);
      if (created) { next.cumulative_enabled = true; next.cumulative_keys = keys; }
    }
  }
  next.selected_keys = Array.from(new Set([...next.instantaneous_keys, ...next.cumulative_keys]));
  return next;
}

/** Persists the merged draft of every circuit a batch touched; `savedLongGroups` is what existed before the batch. */
export function syncBatchDraftsAfterSave(
  workspaceId: string,
  deviceId: string,
  candidates: GroupCandidate[],
  savedLongGroups: WriteGroup[],
  updatedGroups: WriteGroup[],
): void {
  if (!workspaceId.trim() || !deviceId.trim() || updatedGroups.length === 0) return;
  const existingIds = new Set(savedLongGroups.map((group) => group.id));
  const byCircuit = new Map<string, WriteGroup[]>();
  for (const group of updatedGroups) {
    const cid = groupCircuitId(group);
    if (cid) byCircuit.set(cid, [...(byCircuit.get(cid) ?? []), group]);
  }
  for (const [cid, groups] of byCircuit.entries()) {
    const existingDraft = loadCircuitProfileLocalDraft(workspaceId, deviceId, candidates, savedLongGroups, cid);
    saveCircuitProfileLocalDraft(workspaceId, deviceId, mergeSavedProfiles(existingDraft, groups, existingIds));
  }
}

/** `existingGroups` are the groups that existed before the batch; without them every saved group counts as existing. */
export function applyUpdatedGroupsToActiveDraft(
  prevDraft: CircuitProfileLocalDraft,
  activeCircuitId: string | undefined,
  updatedGroups: WriteGroup[],
  existingGroups: WriteGroup[] = updatedGroups,
): CircuitProfileLocalDraft {
  if (!activeCircuitId) return prevDraft;
  const forActive = updatedGroups.filter((group) => groupCircuitId(group) === activeCircuitId);
  if (forActive.length === 0) return prevDraft;
  return mergeSavedProfiles(prevDraft, forActive, new Set(existingGroups.map((group) => group.id)));
}
