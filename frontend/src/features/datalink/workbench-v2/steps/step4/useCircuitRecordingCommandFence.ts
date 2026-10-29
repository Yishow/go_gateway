import { useCallback, useContext, useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import { QueryClientContext } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from '../../../../../hooks/datalink/keys';
import type { WriteGroup, WriteGroupDraft } from '../../../../../types/studioV2WriteGroup';
import type { CandidateCounterStatus } from '../../state/circuitProfileBasic';
import type { ExecuteBatchPlanInput } from './circuitProfileBatchExecutor';
import type { RecordingPreparationCommandFence } from './recordingPreparationPresentation';

export interface CircuitRecordingCommandFenceInput {
  ownerKey: string; scopeKey: string; workspaceId: string; deviceId: string; circuitId: string;
  batchGroups: readonly WriteGroup[];
  counterBlocked?: boolean;
  counterRepair?: { available: boolean; statuses: readonly CandidateCounterStatus[] };
  /** Native controls only: own accepted workspace revisions do not renew this consent. */
  consentKey?: string;
  updateWriteGroup: ExecuteBatchPlanInput['updateWriteGroup'];
  ensureCircuitProfile: NonNullable<ExecuteBatchPlanInput['ensureCircuitProfile']>;
}

const denied = (code: 'RECORDING_START_BUSY' | 'BASIC_INTENT_CHANGED') => ({ outcome: 'failed' as const, code });
function memberCircuit(group: WriteGroupDraft, workspaceId: string, deviceId: string): string | undefined {
  const circuitId = group.members[0]?.circuit_id;
  return group.workspace_id === workspaceId && circuitId && group.members.every(member =>
    member.device_id === deviceId && member.circuit_id === circuitId) ? circuitId : undefined;
}

/** Adds only current-circuit denials; native mutation owners and the executor retain all other guards. */
export function useCircuitRecordingCommandFence(input: CircuitRecordingCommandFenceInput) {
  const { ownerKey, scopeKey, workspaceId, deviceId, circuitId, consentKey } = input;
  const client = useContext(QueryClientContext);
  const strictHandoff = consentKey !== undefined;
  const waiters = useRef(new Set<{ check: () => void; reject: () => void; afterPublication: () => void }>());
  const receipt = useRef<{ epoch: number; previousHead: string; head: string } | undefined>(undefined);
  // Batch busy and successful preceding targets are not new viewing scopes.
  const key = JSON.stringify([ownerKey, scopeKey, workspaceId, deviceId, circuitId, consentKey]);
  const fence = useRef({ key, epoch: 0 });
  if (fence.current.key !== key) fence.current = { key, epoch: fence.current.epoch + 1 };
  const epoch = fence.current.epoch;
  const lease = useMemo(() => ({ active: true, epoch }), [epoch]);
  useLayoutEffect(() => { lease.active = true; const pending = waiters.current; return () => { lease.active = false;
    [...pending].forEach(waiter => waiter.reject());
  }; }, [lease]);
  const latest = useRef(input); latest.current = input;
  const accepted = useRef<{ epoch: number; report: RecordingPreparationCommandFence } | null>(null);
  const [publication, notify] = useState(0);
  // Child report effects publish first. A missing/old report cannot keep a batch pending indefinitely.
  useEffect(() => { [...waiters.current].forEach(waiter => waiter.afterPublication()); }, [publication, epoch]);
  const currentOwner = useCallback(() => lease.active && fence.current.epoch === epoch &&
    Boolean(ownerKey && scopeKey && workspaceId && deviceId), [lease, epoch, ownerKey, scopeKey, workspaceId, deviceId]);
  // Batch identity still requires a selected circuit; saved definition repair does not invent one.
  const current = useCallback(() => currentOwner() && Boolean(circuitId), [currentOwner, circuitId]);
  const canonicalHead = useCallback(() => {
    const state = client?.getQueryState<{ workspace_id?: string; workspace_revision?: string; groups?: WriteGroup[] }>(studioV2WorkspaceKeys.writeGroups());
    if (!state || state.status !== 'success' || state.fetchStatus !== 'idle' || state.data?.workspace_id !== workspaceId ||
      !state.data.workspace_revision || !Array.isArray(state.data.groups)) return undefined;
    return state.data.workspace_revision;
  }, [client, workspaceId]);
  const reportBlocked = useCallback((groupId?: string) => {
    const value = accepted.current;
    if (!currentOwner() || value?.epoch !== epoch || value.report.ownerKey !== ownerKey || value.report.scopeKey !== scopeKey) return true;
    try {
      if (strictHandoff && (!value.report.handoff || value.report.handoff.consentKey !== consentKey ||
        value.report.handoff.workspaceRevision !== canonicalHead() || value.report.handoff.state() !== 'current')) return true;
      // Always aggregate first: group-only clearance must not overlook another selected owner.
      return value.report.isBlocked() || (groupId !== undefined && value.report.isBlocked(groupId));
    } catch { return true; }
  }, [currentOwner, epoch, ownerKey, scopeKey, strictHandoff, consentKey, canonicalHead]);
  const isBlocked = useCallback((groupId?: string) => !current() || reportBlocked(groupId), [current, reportBlocked]);
  const repairKey = JSON.stringify([key, input.counterRepair, input.counterBlocked]);
  const repairFence = useRef({ key: repairKey, epoch: 0 });
  if (repairFence.current.key !== repairKey) repairFence.current = { key: repairKey, epoch: repairFence.current.epoch + 1 };
  const repairEpoch = repairFence.current.epoch;
  const isCounterRepairBlocked = useCallback(() => {
    const repair = latest.current.counterRepair;
    if (!currentOwner() || repairFence.current.epoch !== repairEpoch || !navigator.onLine || !repair?.available) return true;
    // Existing native definition/binding identity only; the source owner still validates and submits the update.
    const saved = repair.statuses.some(({ candidate, measurement: value }) => value && candidate.device_id === deviceId &&
      value.workspace_id === workspaceId && value.device_id === deviceId && value.point_id === candidate.point_id && value.tag_id === candidate.tag_id &&
      [value.id, value.definition_revision, value.source_binding_revision, value.series_epoch].every(part => Boolean(part?.trim())));
    return !saved || reportBlocked();
  }, [currentOwner, repairEpoch, deviceId, workspaceId, reportBlocked]);
  // Counter submit is a different native owner: an open/running batch denies it, never its own sequential progress.
  const isCounterCommandBlocked = useCallback(() => {
    if (repairFence.current.epoch !== repairEpoch || latest.current.counterBlocked ||
      latest.current.counterRepair && !latest.current.counterRepair.available) return true;
    return circuitId ? isBlocked() : isCounterRepairBlocked();
  }, [repairEpoch, circuitId, isBlocked, isCounterRepairBlocked]);
  const onCommandFence = useCallback((report: RecordingPreparationCommandFence | null) => {
    if (!currentOwner()) return;
    if (report && (report.ownerKey !== ownerKey || report.scopeKey !== scopeKey)) {
      [...waiters.current].forEach(waiter => waiter.reject()); return;
    }
    if (accepted.current?.epoch === epoch && accepted.current.report === report || !accepted.current && !report) return;
    accepted.current = report ? { epoch, report } : null;
    [...waiters.current].forEach(waiter => waiter.check());
    notify(value => value + 1);
  }, [currentOwner, epoch, ownerKey, scopeKey]);
  const assertCurrent = useCallback(() => { if (!current()) throw denied('BASIC_INTENT_CHANGED'); }, [current]);
  const beforeDispatch = useCallback(async (expectedHead: string, selected: boolean) => {
    assertCurrent();
    if (!strictHandoff) { if (selected && isBlocked()) throw denied('RECORDING_START_BUSY'); return; }
    const own = receipt.current; receipt.current = undefined;
    if (!selected || !own || own.epoch !== epoch || own.head !== expectedHead) {
      if (canonicalHead() !== expectedHead) throw denied('BASIC_INTENT_CHANGED');
      if (selected && isBlocked()) throw denied('RECORDING_START_BUSY');
      return;
    }
    // The prior native success has already returned to the executor. Only its exact next target may wait.
    await new Promise<void>((resolve, reject) => {
      let finished = false;
      let passes = 0;
      let unsubscribe = () => {};
      const finish = (failed: boolean) => {
        if (finished) return; finished = true; unsubscribe(); waiters.current.delete(waiter);
        if (failed) reject(denied('BASIC_INTENT_CHANGED')); else resolve();
      };
      const waiter = { reject: () => finish(true), afterPublication: () => {
        waiter.check();
        if (finished) return;
        const report = accepted.current?.epoch === epoch ? accepted.current.report : undefined;
        if (report?.handoff?.workspaceRevision !== own.head || ++passes >= 2) { finish(true); return; }
        notify(value => value + 1);
      }, check: () => {
        const report = accepted.current?.epoch === epoch ? accepted.current.report : undefined;
        const head = canonicalHead();
        try {
          if (!current() || head !== own.head || !report?.handoff ||
            report.ownerKey !== ownerKey || report.scopeKey !== scopeKey || report.handoff.consentKey !== consentKey ||
            ![own.previousHead, own.head].includes(report.handoff.workspaceRevision) || report.handoff.state() === 'blocked') { finish(true); return; }
          if (head === own.head && report.handoff.workspaceRevision === own.head && report.handoff.state() === 'current') {
            finish(isBlocked());
          }
        } catch { finish(true); }
      } };
      waiters.current.add(waiter);
      unsubscribe = client?.getQueryCache().subscribe(() => waiter.check()) ?? (() => {});
      waiter.check();
      if (!finished) notify(value => value + 1);
    });
    assertCurrent();
    if (canonicalHead() !== expectedHead || isBlocked()) throw denied('BASIC_INTENT_CHANGED');
  }, [assertCurrent, strictHandoff, isBlocked, epoch, canonicalHead, current, ownerKey, scopeKey, consentKey, client]);
  const rememberSuccess = useCallback((previousHead: string, head: string) => {
    if (strictHandoff && current()) receipt.current = { epoch, previousHead, head };
  }, [strictHandoff, current, epoch]);
  const updateWriteGroup = useCallback<CircuitRecordingCommandFenceInput['updateWriteGroup']>(async vars => {
    assertCurrent();
    const matches = latest.current.batchGroups.filter(group => group.id === vars.id);
    const target = matches.length === 1 ? matches[0] : undefined;
    const targetCircuit = target && target.status !== 'deleted' ? memberCircuit(target, workspaceId, deviceId) : undefined;
    if (!target || !targetCircuit || vars.request.workspace_id !== workspaceId ||
      memberCircuit(vars.request.group, workspaceId, deviceId) !== targetCircuit || vars.request.expected_group_revision !== target.revision) {
      throw denied('BASIC_INTENT_CHANGED');
    }
    await beforeDispatch(vars.request.expected_workspace_revision, targetCircuit === circuitId);
    assertCurrent();
    const now = latest.current.batchGroups.filter(group => group.id === vars.id);
    if (now.length !== 1 || now[0].status === 'deleted' || memberCircuit(now[0], workspaceId, deviceId) !== targetCircuit ||
      now[0].revision !== vars.request.expected_group_revision) throw denied('BASIC_INTENT_CHANGED');
    if (strictHandoff) {
      if (canonicalHead() !== vars.request.expected_workspace_revision) throw denied('BASIC_INTENT_CHANGED');
      const cached = client?.getQueryData<{ groups: WriteGroup[] }>(studioV2WorkspaceKeys.writeGroups())?.groups.filter(group => group.id === vars.id);
      if (cached?.length !== 1 || cached[0].status === 'deleted' || memberCircuit(cached[0], workspaceId, deviceId) !== targetCircuit ||
        cached[0].revision !== vars.request.expected_group_revision) throw denied('BASIC_INTENT_CHANGED');
    }
    if (targetCircuit === circuitId && isBlocked()) throw denied('RECORDING_START_BUSY');
    const result = await latest.current.updateWriteGroup(vars);
    rememberSuccess(vars.request.expected_workspace_revision, result.workspace_revision);
    return result;
  }, [assertCurrent, workspaceId, deviceId, circuitId, beforeDispatch, rememberSuccess, strictHandoff, canonicalHead, isBlocked, client]);
  const ensureCircuitProfile = useCallback<CircuitRecordingCommandFenceInput['ensureCircuitProfile']>(async vars => {
    assertCurrent();
    const targetCircuit = memberCircuit(vars.request.group, workspaceId, deviceId);
    const knownCircuit = targetCircuit === circuitId || latest.current.batchGroups.some(group => group.status !== 'deleted' &&
      memberCircuit(group, workspaceId, deviceId) === targetCircuit);
    if (vars.request.workspace_id !== workspaceId || !targetCircuit || !knownCircuit) throw denied('BASIC_INTENT_CHANGED');
    // A legitimate missing profile has no persisted group ID; its original request supplies the scope.
    await beforeDispatch(vars.request.expected_workspace_revision, targetCircuit === circuitId);
    assertCurrent();
    if (targetCircuit !== circuitId && !latest.current.batchGroups.some(group => group.status !== 'deleted' &&
      memberCircuit(group, workspaceId, deviceId) === targetCircuit)) throw denied('BASIC_INTENT_CHANGED');
    if (strictHandoff && canonicalHead() !== vars.request.expected_workspace_revision) throw denied('BASIC_INTENT_CHANGED');
    if (targetCircuit === circuitId && isBlocked()) throw denied('RECORDING_START_BUSY');
    const result = await latest.current.ensureCircuitProfile(vars);
    rememberSuccess(vars.request.expected_workspace_revision, result.workspace_revision);
    return result;
  }, [assertCurrent, workspaceId, deviceId, circuitId, beforeDispatch, rememberSuccess, strictHandoff, canonicalHead, isBlocked]);
  return { blocked: isBlocked(), isBlocked, counterCommandBlocked: isCounterCommandBlocked(), isCounterCommandBlocked, counterRepairBlocked: isCounterRepairBlocked(), isCounterRepairBlocked, onCommandFence, updateWriteGroup, ensureCircuitProfile };
}
