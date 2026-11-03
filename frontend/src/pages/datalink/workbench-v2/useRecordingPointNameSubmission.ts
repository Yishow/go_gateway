import { useEffect, useMemo, useRef, type MutableRefObject } from 'react';
import type { WorkbenchV2State, Mapping, Point } from '../../../features/datalink/workbench-v2/state/types';
import type { WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { createBasicRecordingRequestId } from '../../../features/datalink/workbench-v2/state/basicRecording';
import { mappingValueFromState, toStudioV2MappingUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2MappingAutosave';
import { sameMappingValue } from './studioV2MappingQueueHelpers';
import { prepareRecordingPointName, type PointNamePreparation } from './prepareRecordingPointName';
import type { PointNameOutcome, RecordingPointNamePort } from './recordingPointNameSubmission';
import { isKnownMappingSourceRejection, type MappingQueueObservation } from './sourceQueueObservation';
import { studioV2MappingsAPI, type StudioV2MappingSaveResult, type StudioV2WorkspaceMappingRequest } from '../../../services/studioV2Mappings';
interface Options {
  workspaceId?: string; enabled: boolean; scopeToken: unknown; stateRef: MutableRefObject<WorkbenchV2State>;
  captureRead?: ReturnType<typeof import('./studioV2LocatorHydration').useStudioV2LocatorHydration>['capturePointNameRead'];
  captureSave?: ReturnType<typeof import('./studioV2LocatorHydration').useStudioV2LocatorHydration>['capturePointNameSave'];
  apply(action: WorkbenchV2Action): void;
  queue(previous: WorkbenchV2State, action: WorkbenchV2Action, next: WorkbenchV2State): void;
  canMutateRule(id: string): boolean;
  canWriteRule(id: string): boolean;
  canQueueRule(id: string): boolean;
}
interface Entry { cancelRejected?: (original: Mapping) => boolean; acceptSave?: (saved: StudioV2MappingSaveResult) => void; flight?: boolean; owner: object; preparation: PointNamePreparation; outcome: PointNameOutcome; submitted?: StudioV2WorkspaceMappingRequest; resolve?: (value: PointNameOutcome) => void }
const equal = (expected: object, actual: object) => Object.entries(expected).every(([key, value]) => value === (actual as Record<string, unknown>)[key]);
export function useRecordingPointNameSubmission(options: Options) {
  const current = useRef(options); current.current = options;
  const owner = useRef({ workspaceId: options.workspaceId, enabled: options.enabled, scopeToken: options.scopeToken, token: {} });
  const entries = useRef(new Map<string, Entry>());
  const mounted = useRef(true);
  if (owner.current.workspaceId !== options.workspaceId || owner.current.enabled !== options.enabled || owner.current.scopeToken !== options.scopeToken) {
    for (const entry of entries.current.values()) if (entry.outcome.status === 'pending') { entry.outcome = { ...entry.outcome, status: 'unknown' }; entry.resolve?.(entry.outcome); entry.resolve = undefined; }
    owner.current = { workspaceId: options.workspaceId, enabled: options.enabled, scopeToken: options.scopeToken, token: {} };
  }
  useEffect(() => { mounted.current = true; const activeEntries = entries.current; return () => { mounted.current = false; for (const entry of activeEntries.values()) { entry.resolve?.({ ...entry.outcome, status: 'unknown' }); entry.resolve = undefined; } }; }, []);
  return useMemo(() => {
    const active = (entry: Entry) => entries.current.get(entry.outcome.session.id) === entry && mounted.current && current.current.enabled && entry.owner === owner.current.token && entry.preparation.workspaceId === current.current.workspaceId;
    const finish = (entry: Entry, status: PointNameOutcome['status'], currentName?: string) => {
      entry.outcome = { ...entry.outcome, status, ...(currentName === undefined ? {} : { currentName }) }; entry.resolve?.(entry.outcome); entry.resolve = undefined; return entry.outcome;
    };
    const matched = (point: Point, mapping: Mapping) => [...entries.current.values()].find(entry => active(entry) && entry.outcome.status === 'pending' && entry.preparation.point.id === point.id && entry.preparation.mappingId === mapping.mapping_id && entry.submitted && equal(entry.submitted, toStudioV2MappingUpdateRequest(point, mapping)));
    const liveMatches = (entry: Entry) => {
      const state = current.current.stateRef.current, point = state.points.find(value => value.id === entry.preparation.point.id), mapping = state.mappings[entry.preparation.point.id];
      return !!point && !!mapping && point.device_id === entry.preparation.identity.device_id && point.rule_id === entry.preparation.ruleId && mapping.mapping_id === entry.preparation.mappingId && mapping.tag_id === entry.preparation.identity.tag_id && mapping.persisted_point_id === entry.preparation.identity.point_id &&
        state.rules.find(value => value.id === entry.preparation.ruleId)?.revision_id === entry.preparation.ruleRevision && !!entry.submitted && equal(entry.submitted, toStudioV2MappingUpdateRequest(point, mapping));
    };
    const acknowledged = (entry: Entry, saved: StudioV2MappingSaveResult) => saved.id === entry.preparation.mappingId && saved.workspace_id === entry.preparation.workspaceId && saved.device_id === entry.preparation.identity.device_id && saved.point_id === entry.preparation.identity.point_id && saved.tag_id === entry.preparation.identity.tag_id && !!entry.submitted && equal(entry.submitted, saved);
    const port: RecordingPointNamePort = {
      prepare: async identity => {
        const captured = owner.current.token, workspaceId = current.current.workspaceId;
        if (!mounted.current || !current.current.enabled || !workspaceId) return;
        const prior = [...entries.current.values()].find(entry => active(entry) && entry.preparation.identity.device_id === identity.device_id && entry.preparation.identity.point_id === identity.point_id && entry.preparation.identity.tag_id === identity.tag_id && ['pending', 'unknown'].includes(entry.outcome.status));
        if (prior) return prior.outcome;
        const result = await prepareRecordingPointName(workspaceId, identity, () => current.current.stateRef.current, () => mounted.current && current.current.enabled && owner.current.token === captured);
        if (!result || !mounted.current || owner.current.token !== captured || current.current.stateRef.current !== result.expectedState) return;
        const session = { id: createBasicRecordingRequestId(), identity: { ...identity }, displayName: result.preparation.mapping.display_name };
        const entry: Entry = { owner: captured, preparation: result.preparation, outcome: { session, status: 'ready' } };
        entries.current.set(session.id, entry);
        if (result.merge) {
          const acceptRead = current.current.captureRead?.(result.expectedState, result.preparation, result.merge);
          current.current.apply({ type: 'SET_STATE', payload: result.merge }); acceptRead?.();
        }
        if (!current.current.canQueueRule(result.preparation.ruleId) || !current.current.canWriteRule(result.preparation.ruleId) || (!current.current.canMutateRule(result.preparation.ruleId) && !compatibleMC(result.preparation.ruleId))) return finish(entry, 'blocked');
        return entry.outcome;
      },
      release: id => { const entry = entries.current.get(id); if (entry && entry.outcome.status === 'ready') entries.current.delete(id); },
      submit: async (id, name) => {
        const entry = entries.current.get(id);
        if (!entry || !active(entry) || entry.outcome.status !== 'ready' || entry.flight || !name.trim()) return;
        entry.flight = true;
        try {
        const fresh = await prepareRecordingPointName(entry.preparation.workspaceId, entry.preparation.identity, () => current.current.stateRef.current, () => active(entry));
        if (!active(entry) || entry.outcome.status !== 'ready' || !fresh || current.current.stateRef.current !== fresh.expectedState || fresh.preparation.mappingId !== entry.preparation.mappingId || fresh.preparation.ruleRevision !== entry.preparation.ruleRevision ||
          !sameMappingValue(mappingValueFromState(fresh.preparation.mapping), mappingValueFromState(entry.preparation.mapping)) || !current.current.canQueueRule(entry.preparation.ruleId) || !current.current.canWriteRule(entry.preparation.ruleId) || (!current.current.canMutateRule(entry.preparation.ruleId) && !compatibleMC(entry.preparation.ruleId))) return finish(entry, 'blocked');
        entry.preparation = { ...entry.preparation, serverRecord: fresh.preparation.serverRecord };
        const previous = current.current.stateRef.current, action: WorkbenchV2Action = { type: 'updateMapping', pointId: entry.preparation.point.id, patch: { display_name: name.trim() } };
        entry.submitted = { ...toStudioV2MappingUpdateRequest(entry.preparation.point, entry.preparation.mapping), display_name: name.trim() };
        entry.acceptSave = current.current.captureSave?.(fresh.preparation, entry.submitted);
        entry.outcome = { ...entry.outcome, status: 'pending', submittedName: name.trim() };
        const promise = new Promise<PointNameOutcome>(resolve => { entry.resolve = resolve; });
        current.current.apply(action); current.current.queue(previous, action, current.current.stateRef.current);
        return promise;
        } finally { entry.flight = false; }
      },
      discardRejected: async id => {
        const entry = entries.current.get(id);
        if (!entry || !active(entry) || entry.flight || entry.outcome.status !== 'unknown' || entry.outcome.reason !== 'source-review' || !entry.cancelRejected || !liveMatches(entry)) return false;
        entry.flight = true;
        try {
          const rows = await studioV2MappingsAPI.list();
          const matches = rows.filter(row => row.id === entry.preparation.mappingId);
          const record = matches[0], prepared = entry.preparation;
          if (!active(entry) || entry.outcome.status !== 'unknown' || !liveMatches(entry) || matches.length !== 1 ||
            record.workspace_id !== prepared.workspaceId || record.device_id !== prepared.identity.device_id || record.point_id !== prepared.identity.point_id || record.tag_id !== prepared.identity.tag_id ||
            !prepared.serverRecord || prepared.serverRecord.created_at !== record.created_at || prepared.serverRecord.updated_at !== record.updated_at || prepared.serverRecord.save_state !== record.save_state || !equal(toStudioV2MappingUpdateRequest(prepared.point, prepared.mapping), record) || !entry.cancelRejected(prepared.mapping)) return false;
          entries.current.delete(id);
          return true;
        } catch { return false; } finally { entry.flight = false; }
      },
      inspect: async id => {
        const entry = entries.current.get(id); if (!entry || !active(entry) || entry.flight || entry.outcome.status !== 'unknown') return;
        entry.flight = true;
        try {
          const rows = (await studioV2MappingsAPI.list()).filter(row => row.workspace_id === entry.preparation.workspaceId && row.id === entry.preparation.mappingId && row.device_id === entry.preparation.identity.device_id && row.point_id === entry.preparation.identity.point_id && row.tag_id === entry.preparation.identity.tag_id);
          // A current read does not prove the uncertain write and never enables retry.
          return active(entry) ? finish(entry, 'unknown', rows.length === 1 ? rows[0].display_name : undefined) : undefined;
        } finally { entry.flight = false; }
      },
    };
    const observation: MappingQueueObservation = {
      acknowledged: (point, mapping, saved) => { const entry = matched(point, mapping); if (!entry) return true; if (!current.current.canWriteRule(entry.preparation.ruleId) || !acknowledged(entry, saved) || !liveMatches(entry)) { finish(entry, 'unknown'); return false; } return true; },
      saved: (point, mapping, saved) => {
        const entry = matched(point, mapping); if (!entry) return;
        if (!acknowledged(entry, saved) || !liveMatches(entry) || saved.runtime_apply_status === 'apply_failed' || saved.cleanup_status === 'failed') { finish(entry, 'unknown'); return; }
        void prepareRecordingPointName(entry.preparation.workspaceId, entry.preparation.identity, () => current.current.stateRef.current, () => active(entry)).then(result => {
          if (!active(entry)) { finish(entry, 'unknown'); return; }
          const confirmed = result && result.preparation.mappingId === entry.preparation.mappingId && entry.submitted && equal(entry.submitted, toStudioV2MappingUpdateRequest(result.preparation.point, result.preparation.mapping)) && liveMatches(entry);
          if (confirmed) entry.acceptSave?.(saved);
          finish(entry, confirmed ? 'confirmed' : 'unknown', result?.preparation.mapping.display_name);
        }).catch(() => finish(entry, 'unknown'));
      },
      failed: (point, error) => { for (const entry of entries.current.values()) if (active(entry) && entry.outcome.status === 'pending' && entry.preparation.point.id === point.id) {
        if (isKnownMappingSourceRejection(error) && liveMatches(entry)) entry.outcome = { ...entry.outcome, reason: 'source-review' };
        finish(entry, 'unknown');
      } },
      rejectedName: (point, mapping, cancel) => {
        for (const entry of entries.current.values()) if (active(entry) && entry.outcome.status === 'unknown' && entry.outcome.reason === 'source-review' &&
          entry.preparation.point.id === point.id && entry.preparation.mappingId === mapping.mapping_id && entry.submitted && equal(entry.submitted, toStudioV2MappingUpdateRequest(point, mapping)) && liveMatches(entry)) {
          entry.cancelRejected = cancel;
          // The resolved outcome is the same notification object consumed by the awaiting buffer.
          entry.outcome.canDiscardRejected = true;
        }
      },
    };
    const compatibleMC = (ruleId: string) => [...entries.current.values()].some(entry => active(entry) && entry.preparation.ruleId === ruleId && entry.preparation.mcEvidence?.status === 'compatible' &&
      current.current.stateRef.current.rules.find(rule => rule.id === ruleId)?.revision_id === entry.preparation.ruleRevision && current.current.stateRef.current.rules.find(rule => rule.id === ruleId)?.mc_persisted_points === entry.preparation.mcEvidence && !['blocked', 'unknown'].includes(entry.outcome.status));
    const allowPreparedMapping = (point: Point, mapping: Mapping) => {
      const entry = matched(point, mapping);
      return !!entry && compatibleMC(point.rule_id) && liveMatches(entry) && point.device_id === entry.preparation.identity.device_id && point.address === entry.preparation.point.address && point.data_type === entry.preparation.point.data_type && mapping.workspace_id === entry.preparation.workspaceId;
    };
    return { port, observation, compatibleMC, allowPreparedMapping };
  }, []);
}
