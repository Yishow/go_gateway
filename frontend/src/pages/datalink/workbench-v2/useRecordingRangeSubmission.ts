import { useEffect, useMemo, useRef, useState, type MutableRefObject } from 'react';
import type { Mapping, Point, Rule, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import type { WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { hydrateStudioV2Device } from '../../../features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { hydrateStudioV2Rule, isStudioV2RuleValid, toStudioV2RuleCreateRequest } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { derivePoints, getRuleReadinessReason } from '../../../features/datalink/workbench-v2/state/sourceRule';
import { mappingValueFromState, toStudioV2MappingUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2MappingAutosave';
import type { RangeMappingReceipt, RangeRuleReceipt, RangeSourceOutcome, RangeSourcePort, RangeSourcePreparation } from './rangeSourceSubmission';
import type { RuleQueueObservation, MappingQueueObservation } from './sourceQueueObservation';
import { readRuleProof, readRangeCompletion, matchesRangeRuleAcknowledgement, type RangeRuleProof } from './rangeSourceReadback';
export interface RangeSubmissionOptions {
  workspaceId?: string; enabled: boolean; stateRef: MutableRefObject<WorkbenchV2State>;
  apply(action: WorkbenchV2Action): void;
  prepareDevice(id: string): Promise<RangeSourcePreparation | undefined>;
  diagnoseDevice(id: string): Promise<boolean>;
  queueRule(action: WorkbenchV2Action): void;
  queueMapping(previous: WorkbenchV2State, action: WorkbenchV2Action, next: WorkbenchV2State): void;
}
interface Entry {
  owner: object; outcome: RangeSourceOutcome; snapshot: Rule; receipt?: RangeRuleReceipt; proof?: RangeRuleProof;
  preparing?: boolean; mappingStarted: boolean; mappingIds: Set<string>; mappings: Map<string, RangeMappingReceipt>; checking: boolean;
  promise: Promise<RangeSourceOutcome>; resolve: (result: RangeSourceOutcome) => void;
}
const same = (left: unknown, right: unknown): boolean => JSON.stringify(left) === JSON.stringify(right);
/** Coordinates explicit intents using the existing queues. Owns no API mutation,
 * autosave engine, durable draft, or fabricated persisted identity. */
export function useRecordingRangeSubmission(options: RangeSubmissionOptions) {
  const current = useRef(options); current.current = options;
  const owner = useRef({ workspaceId: options.workspaceId, enabled: options.enabled, token: {} });
  if (owner.current.workspaceId !== options.workspaceId || owner.current.enabled !== options.enabled) owner.current = { workspaceId: options.workspaceId, enabled: options.enabled, token: {} };
  const entries = useRef(new Map<string, Entry>());
  const mounted = useRef(true);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const [version, changed] = useState(0);
  const controller = useMemo(() => {
    const active = (entry: Entry) => mounted.current && current.current.enabled && entry.owner === owner.current.token && entry.outcome.request.workspaceId === current.current.workspaceId;
    const byRule = (id: string) => [...entries.current.values()].reverse().find(entry => entry.snapshot.id === id);
    const publish = (entry: Entry, status: RangeSourceOutcome['status'], extra: Partial<RangeSourceOutcome> = {}) => {
      entry.outcome = { ...entry.outcome, ...extra, status }; entry.resolve(structuredClone(entry.outcome)); if (mounted.current) changed(value => value + 1); return structuredClone(entry.outcome);
    };
    const apply = (action: WorkbenchV2Action, mapping = false) => {
      const previous = current.current.stateRef.current; current.current.apply(action);
      if (mapping) current.current.queueMapping(previous, action, current.current.stateRef.current);
    };
    const complete = async (entry: Entry) => {
      if (!active(entry) || !entry.proof || entry.mappings.size !== entry.mappingIds.size || entry.mappingIds.size === 0) return publish(entry, 'unknown');
      const candidates = await readRangeCompletion(entry.proof, [...entry.mappings.values()]);
      if (!active(entry) || !candidates) return publish(entry, 'unknown');
      return publish(entry, candidates.inactive ? 'saved-inactive' : 'confirmed', { rule: entry.proof.record, candidates, savedPoints: candidates.savedPoints });
    };
    const advance = async (entry: Entry) => {
      if (entry.checking) return entry.outcome;
      entry.checking = true;
      try {
        if (!active(entry) || !entry.receipt) return publish(entry, 'unknown');
        const proof = await readRuleProof(entry.receipt, entry.outcome.request);
        if (!active(entry) || !proof) return publish(entry, 'unknown');
        entry.proof = proof;
        if (entry.mappingStarted) return complete(entry);
        const rule = { ...hydrateStudioV2Rule(proof.record), ...(proof.mcEvidence ? { mc_persisted_points: proof.mcEvidence } : {}) };
        apply({ type: 'updateRule', ruleId: rule.id, patch: rule });
        const points = derivePoints(rule, proof.device.id, new Set(rule.skipped_addresses), proof.device.protocol).filter(point => !point.skipped);
        if (!points.length || points.length !== proof.snapshot.tags.candidates.length) return publish(entry, 'unknown');
        // A verified known rule may proceed to its never-submitted mapping phase.
        // Do not resolve the submission promise before complete canonical readback.
        entry.outcome = { ...entry.outcome, status: 'pending' };
        entry.mappingStarted = true; entry.mappingIds = new Set(points.map(point => point.id)); changed(value => value + 1);
        // Append only this new range. Never replace unrelated points/mappings or recording drafts.
        apply({ type: 'initMappingsForPoints', points: [...current.current.stateRef.current.points.filter(point => point.rule_id !== rule.id), ...points] }, true);
        for (const point of points) {
          const mapping = current.current.stateRef.current.mappings[point.id];
          if (!mapping) return publish(entry, 'unknown');
          apply({ type: 'updateMapping', pointId: point.id, patch: { ...mappingValueFromState(mapping), enabled: rule.enabled } }, true);
        }
        return entry.outcome;
      } catch { return publish(entry, 'unknown'); }
      finally { entry.checking = false; }
    };
    const unchangedLocal = (source: Entry, proof: RangeRuleProof) => {
      const state = current.current.stateRef.current;
      const rule = state.rules.find(row => row.id === proof.record.id);
      if (!rule || rule.save_state !== 'saved' || !rule.persisted || rule.revision_id !== proof.record.revision_id ||
        !same(toStudioV2RuleCreateRequest(rule), toStudioV2RuleCreateRequest(hydrateStudioV2Rule(proof.record)))) return false;
      return [...source.mappings.values()].every(receipt => {
        const point = state.points.find(row => row.rule_id === proof.record.id && row.address === receipt.submitted.address);
        const mapping = point && state.mappings[point.id];
        return point && mapping?.save_state === 'saved' && same(toStudioV2MappingUpdateRequest(point, mapping), receipt.submitted);
      });
    };
    const port: RangeSourcePort = {
      prepare: async id => {
        const captured = owner.current.token;
        if (!mounted.current || !current.current.enabled) return;
        const prepared = await current.current.prepareDevice(id);
        const device = prepared?.device;
        if (!mounted.current || owner.current.token !== captured || !current.current.enabled || device?.id !== id || device.runtime_apply_status === 'apply_failed') return;
        return prepared;
      },
      canSelect: (deviceId, pointId) => ![...entries.current.values()].some(entry => byRule(entry.snapshot.id) === entry && entry.outcome.request.workspaceId === current.current.workspaceId && entry.outcome.request.deviceId === deviceId && entry.proof?.snapshot.tags.candidates.some(point => point.point_id === pointId) && entry.outcome.status !== 'confirmed'),
      observations: () => [...entries.current.values()].filter(entry => byRule(entry.snapshot.id) === entry && entry.outcome.request.workspaceId === current.current.workspaceId).map(entry => structuredClone(entry.outcome)),
      submit: request => {
        const previous = entries.current.get(request.requestId);
        if (previous) return same(previous.outcome.request, request) && active(previous) ? previous.outcome.status === 'pending' ? previous.promise : Promise.resolve(structuredClone(previous.outcome)) : Promise.resolve({ request, status: 'blocked' });
        if (!mounted.current || !current.current.enabled || request.workspaceId !== current.current.workspaceId || !request.requestId || !request.sessionId || !Number.isSafeInteger(request.generation) || request.generation < 0 ||
          [...entries.current.values()].some(entry => entry.outcome.request.workspaceId === request.workspaceId && ['pending', 'unknown', 'persisted-unavailable'].includes(entry.outcome.status))) return Promise.resolve({ request, status: 'blocked' });
        const snapshot: Rule = { ...structuredClone(request.buffer), id: `rule-${crypto.randomUUID()}`, device_id: request.deviceId, persisted: false, save_state: 'idle' };
        let resolve!: Entry['resolve']; const promise = new Promise<RangeSourceOutcome>(yes => { resolve = yes; });
        const entry: Entry = { owner: owner.current.token, outcome: { request: structuredClone(request), status: 'pending' }, snapshot, mappingStarted: false, mappingIds: new Set(), mappings: new Map(), checking: false, promise, resolve };
        entries.current.set(request.requestId, entry); changed(value => value + 1);
        void (async () => {
          try {
            const prepared = await port.prepare(request.deviceId);
            const device = prepared?.device;
            if (!active(entry) || !device || (snapshot.enabled && !prepared?.activationAllowed) || !isStudioV2RuleValid(snapshot) || !Number.isSafeInteger(snapshot.count) || getRuleReadinessReason(snapshot, { [device.id]: device.protocol })) { publish(entry, 'blocked'); return; }
            const local = current.current.stateRef.current.devices.find(row => row.id === device.id);
            if (local && (!local.persisted || local.save_state !== 'saved' || local.protocol !== device.protocol)) { publish(entry, 'blocked'); return; }
            if (!local) apply({ type: 'addDevice', device: hydrateStudioV2Device(device) });
            const action: WorkbenchV2Action = { type: 'addRule', rule: snapshot }; apply(action); current.current.queueRule(action);
          } catch { publish(entry, 'blocked'); }
        })();
        return promise;
      },
      prepareSaved: async (sourceRequestId, request) => {
        const source = entries.current.get(sourceRequestId);
        const existing = entries.current.get(request.requestId);
        if (existing) return same(existing.outcome.request, request) ? existing.promise : { request, status: 'blocked' };
        if (!source || !active(source) || source.outcome.status !== 'saved-inactive' || !source.receipt || !source.proof ||
          request.workspaceId !== source.outcome.request.workspaceId || request.deviceId !== source.outcome.request.deviceId ||
          !request.requestId || request.requestId === sourceRequestId || !request.sessionId || !Number.isSafeInteger(request.generation) || request.generation < 0 ||
          !same({ ...source.outcome.request.buffer, enabled: true }, request.buffer)) return { request, status: 'blocked' };
        let resolve!: Entry['resolve']; const promise = new Promise<RangeSourceOutcome>(yes => { resolve = yes; });
        const entry: Entry = { ...source, owner: owner.current.token, outcome: { request: structuredClone(request), status: 'pending', rule: source.proof.record, savedPoints: source.outcome.savedPoints },
          snapshot: { ...hydrateStudioV2Rule(source.proof.record), enabled: true }, preparing: true, mappingStarted: false, mappingIds: new Set(source.mappingIds), mappings: new Map(source.mappings), checking: false, promise, resolve };
        entries.current.set(request.requestId, entry); changed(value => value + 1);
        void (async () => {
          let mutationStarted = false;
          const inactive = (preparationError: 'diagnostic-failed' | 'source-changed') => { entry.preparing = false; return publish(entry, 'saved-inactive', { preparationError }); };
          try {
            const before = await readRuleProof(source.receipt!, source.proof!.request);
            const prepared = await port.prepare(request.deviceId);
            if (!active(entry) || !before || !prepared || before.device.protocol !== prepared.device.protocol || before.device.connection_config !== prepared.device.connection_config ||
              !unchangedLocal(source, before) || !await readRangeCompletion(before, [...source.mappings.values()])) { inactive('source-changed'); return; }
            if (!active(entry)) { inactive('source-changed'); return; }
            if (!await current.current.diagnoseDevice(request.deviceId)) { inactive('diagnostic-failed'); return; }
            const after = await readRuleProof(source.receipt!, source.proof!.request);
            if (!active(entry) || !after || !unchangedLocal(source, after) || !await readRangeCompletion(after, [...source.mappings.values()]) || !same(before.record, after.record) || before.device.connection_config !== after.device.connection_config || before.device.protocol !== after.device.protocol) { inactive('source-changed'); return; }
            if (!active(entry) || !unchangedLocal(source, after)) { inactive('source-changed'); return; }
            entry.receipt = undefined; entry.proof = undefined; entry.preparing = false; entry.mappingIds = new Set(); entry.mappings = new Map();
            const action: WorkbenchV2Action = { type: 'updateRule', ruleId: entry.snapshot.id, patch: { enabled: true } };
            mutationStarted = true; apply(action); current.current.queueRule(action);
          } catch { if (mutationStarted) publish(entry, 'unknown'); else inactive('diagnostic-failed'); }
        })();
        return promise;
      },
      inspect: async requestId => {
        const entry = entries.current.get(requestId);
        if (!entry || entry.outcome.request.workspaceId !== current.current.workspaceId) return;
        if (entry.outcome.status === 'unknown' && entry.receipt && active(entry)) {
          if (entry.mappingStarted && entry.mappings.size !== entry.mappingIds.size) return structuredClone(entry.outcome); // No lost mapping ack guessing or blind create retry.
          await advance(entry);
        }
        return structuredClone(entry.outcome);
      },
    };
    const ruleObservation: RuleQueueObservation = {
      acknowledged: (rule, saved) => {
        const entry = byRule(rule.id); if (!entry) return true;
        const match = active(entry) && same(toStudioV2RuleCreateRequest(entry.snapshot), toStudioV2RuleCreateRequest(rule)) &&
          matchesRangeRuleAcknowledgement({ submitted: toStudioV2RuleCreateRequest(rule), acknowledged: saved }, entry.outcome.request);
        if (!match) { publish(entry, 'unknown'); return false; }
        entry.receipt = { submitted: toStudioV2RuleCreateRequest(rule), acknowledged: structuredClone(saved) }; return true;
      },
      saved: (rule, saved) => { const entry = byRule(rule.id); if (!entry) return; if (saved.runtime_apply_status === 'apply_failed') publish(entry, 'persisted-unavailable', { rule: saved }); else void advance(entry); },
      failed: rule => { const entry = byRule(rule.id); if (entry) publish(entry, 'unknown'); },
    };
    const mappingReceipt = (point: Point, mapping: Mapping, saved: RangeMappingReceipt['acknowledged']): RangeMappingReceipt => ({ submitted: toStudioV2MappingUpdateRequest(point, mapping), acknowledged: structuredClone(saved) });
    const mappingObservation: MappingQueueObservation = {
      acknowledged: (point, mapping, saved) => {
        const entry = byRule(point.rule_id); if (!entry) return true;
        const submitted = toStudioV2MappingUpdateRequest(point, mapping);
        const candidate = entry.proof?.snapshot.tags.candidates.find(row => row.address === point.address);
        if (!active(entry) || !entry.mappingIds.has(point.id) || !saved.id || !saved.tag_id || saved.point_id !== candidate?.point_id || saved.workspace_id !== entry.outcome.request.workspaceId || saved.device_id !== entry.outcome.request.deviceId ||
          !Object.entries(submitted).every(([key, value]) => same(value, (saved as unknown as Record<string, unknown>)[key]))) { publish(entry, 'unknown'); return false; }
        return true;
      },
      saved: (point, mapping, saved) => {
        const entry = byRule(point.rule_id); if (!entry) return;
        entry.mappings.set(point.id, mappingReceipt(point, mapping, saved));
        if (saved.runtime_apply_status === 'apply_failed' || saved.cleanup_status === 'failed') { publish(entry, 'persisted-unavailable'); return; }
        if (entry.mappings.size === entry.mappingIds.size) void complete(entry);
      },
      failed: point => { const entry = byRule(point.rule_id); if (entry) publish(entry, 'unknown'); },
    };
    return { port, ruleObservation, mappingObservation,
      canWrite: (id: string) => { const entry = byRule(id); return !entry || active(entry) && !entry.preparing && ['pending', 'confirmed'].includes(entry.outcome.status); },
      compatibleMC: (id: string) => { const entry = byRule(id); const rule = current.current.stateRef.current.rules.find(row => row.id === id); return Boolean(entry && active(entry) && entry.proof?.mcEvidence?.status === 'compatible' && entry.proof.record.revision_id === rule?.revision_id); },
    };
  }, []);
  return { ...controller, version };
}
