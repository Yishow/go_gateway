import { useMemo, useRef, useState } from 'react';
import type { MutableRefObject } from 'react';
import type { Device, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import type { WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { isStudioV2DeviceValid, toBackendConnectionConfig } from '../../../features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import type { StudioV2WorkspaceDeviceRecord } from '../../../services/studioV2WorkspaceDevices';
import type { DeviceSourceOutcome, DeviceSourcePort } from './deviceSourceSubmission';

interface Options {
  workspaceId?: string;
  enabled: boolean;
  canSubmit?: () => boolean;
  stateRef: MutableRefObject<WorkbenchV2State>;
  apply: (action: WorkbenchV2Action) => void;
  busy: (id: string) => boolean;
  queue: (id: string) => void;
  resume: (id: string) => void;
  readback: () => Promise<readonly StudioV2WorkspaceDeviceRecord[] | undefined>;
}
interface Entry {
  outcome: DeviceSourceOutcome;
  deviceId: string;
  owner: object;
  snapshot: Device;
  acknowledgement?: StudioV2WorkspaceDeviceRecord;
  promise: Promise<DeviceSourceOutcome>;
  resolve: (outcome: DeviceSourceOutcome) => void;
}
function comparable(value: unknown): string {
  if (Array.isArray(value)) return `[${value.map(comparable).join(',')}]`;
  if (value && typeof value === 'object') return `{${Object.entries(value).sort(([a], [b]) => a.localeCompare(b)).map(([key, item]) => `${JSON.stringify(key)}:${comparable(item)}`).join(',')}}`;
  return JSON.stringify(value);
}
function matches(record: StudioV2WorkspaceDeviceRecord, entry: Entry): boolean {
  if (record.id !== entry.deviceId || record.name !== entry.snapshot.name || record.description !== entry.snapshot.description || record.protocol !== entry.snapshot.protocol) return false;
  try { return comparable(JSON.parse(record.connection_config)) === comparable(toBackendConnectionConfig(entry.snapshot)); }
  catch { return false; }
}

/** Observes the existing device queue; this hook owns no API mutation or second queue. */
export function useDeviceSourceSubmission(options: Options) {
  const current = useRef(options); current.current = options;
  const owner = useRef({ workspaceId: options.workspaceId, token: {} });
  if (owner.current.workspaceId !== options.workspaceId) owner.current = { workspaceId: options.workspaceId, token: {} };
  // Existing autosave state hydrates once. A scope change must remount it before new writes.
  const hydratedWorkspace = useRef<string | undefined>(undefined);
  if (!hydratedWorkspace.current && options.enabled && options.workspaceId) hydratedWorkspace.current = options.workspaceId;
  const entries = useRef(new Map<string, Entry>());
  const [, changed] = useState(0);
  return useMemo(() => {
    const publish = (entry: Entry, status: DeviceSourceOutcome['status'], device?: StudioV2WorkspaceDeviceRecord) => {
      entry.outcome = { request: entry.outcome.request, status, ...(device ? { device } : {}) };
      entry.resolve(structuredClone(entry.outcome)); changed(value => value + 1); return structuredClone(entry.outcome);
    };
    const active = (entry: Entry) => current.current.enabled && owner.current.token === entry.owner &&
      hydratedWorkspace.current === current.current.workspaceId && entry.outcome.request.workspaceId === current.current.workspaceId;
    const byDevice = (id: string) => [...entries.current.values()].reverse().find(entry => entry.deviceId === id);
    const readback = async (entry: Entry): Promise<DeviceSourceOutcome> => {
      if (!active(entry) || !entry.acknowledgement) return publish(entry, 'unknown');
      try {
        const rows = await current.current.readback();
        if (!active(entry)) return publish(entry, 'unknown');
        const exact = rows?.filter(row => row.id === entry.deviceId);
        if (exact?.length !== 1 || !matches(exact[0], entry)) return publish(entry, 'unknown');
        return publish(entry, entry.acknowledgement.runtime_apply_status === 'apply_failed' || exact[0].runtime_apply_status === 'apply_failed'
          ? 'persisted-unavailable' : 'confirmed', exact[0]);
      } catch { return publish(entry, 'unknown'); }
    };
    const port: DeviceSourcePort = {
      observations: () => {
        const scoped = [...entries.current.values()].filter(entry => entry.outcome.request.workspaceId === current.current.workspaceId);
        return scoped.filter((entry, index) => entry.outcome.status !== 'persisted-unavailable' ||
          !scoped.slice(index + 1).some(newer => newer.deviceId === entry.deviceId && newer.outcome.status === 'confirmed'))
          .map(entry => structuredClone(entry.outcome));
      },
      submit: request => {
        const existing = entries.current.get(request.requestId);
        if (existing) return active(existing) && comparable(existing.outcome.request) === comparable(request)
          ? existing.outcome.status === 'pending' ? existing.promise : Promise.resolve(structuredClone(existing.outcome))
          : Promise.resolve({ request, status: 'blocked' });
        const opt = current.current;
        const unresolved = [...entries.current.values()].some(entry => entry.outcome.request.workspaceId === request.workspaceId &&
          (entry.outcome.status === 'pending' || entry.outcome.status === 'unknown') &&
          (entry.outcome.request.sessionId === request.sessionId || request.deviceId === entry.deviceId));
        if (!opt.enabled || opt.canSubmit?.() === false || !request.requestId || !request.sessionId || !Number.isSafeInteger(request.generation) || request.generation < 0 ||
          !request.workspaceId || request.workspaceId !== opt.workspaceId || hydratedWorkspace.current !== opt.workspaceId || unresolved) return Promise.resolve({ request, status: 'blocked' });
        const previous = request.deviceId ? opt.stateRef.current.devices.find(device => device.id === request.deviceId && device.persisted) : undefined;
        if (request.deviceId && (!previous || previous.protocol !== request.buffer.protocol || opt.busy(request.deviceId))) return Promise.resolve({ request, status: 'blocked' });
        const snapshot: Device = previous ? { ...previous, ...structuredClone(request.buffer) }
          : { ...structuredClone(request.buffer), id: `dev-${crypto.randomUUID()}`, status: 'draft', test: null, persisted: false, save_state: 'idle' };
        if (!isStudioV2DeviceValid(snapshot) || snapshot.protocol === 'mqtt') return Promise.resolve({ request, status: 'blocked' });
        let resolve!: Entry['resolve']; const promise = new Promise<DeviceSourceOutcome>(yes => { resolve = yes; });
        const entry: Entry = { outcome: { request: structuredClone(request), status: 'pending' }, snapshot, deviceId: snapshot.id, owner: owner.current.token, promise, resolve };
        entries.current.set(request.requestId, entry); changed(value => value + 1);
        opt.apply(previous ? { type: 'updateDevice', deviceId: snapshot.id, patch: snapshot } : { type: 'addDevice', device: snapshot });
        opt.queue(snapshot.id); return promise;
      },
      inspect: async requestId => {
        const entry = entries.current.get(requestId);
        if (!entry || entry.outcome.request.workspaceId !== current.current.workspaceId) return undefined;
        if (entry.outcome.status === 'unknown' && entry.acknowledgement && active(entry)) {
          const outcome = await readback(entry);
          if (outcome.status === 'confirmed' || outcome.status === 'persisted-unavailable') {
            const latest = current.current.stateRef.current.devices.find(device => device.id === entry.deviceId);
            if (latest && comparable(toBackendConnectionConfig(latest)) === comparable(toBackendConnectionConfig(entry.snapshot)) &&
              latest.name === entry.snapshot.name && latest.description === entry.snapshot.description && latest.protocol === entry.snapshot.protocol) {
              current.current.apply({ type: 'updateDevice', deviceId: entry.deviceId, patch: {
                persisted: true, save_state: outcome.status === 'confirmed' ? 'saved' : 'save-error',
                save_error: null, runtime_apply_status: outcome.device?.runtime_apply_status,
              } });
            }
            current.current.resume(entry.deviceId);
          }
          return outcome;
        }
        return structuredClone(entry.outcome);
      },
    };
    return {
      port,
      canWrite: (id: string) => { const entry = byDevice(id); return !entry || active(entry) && entry.outcome.status !== 'unknown'; },
      canPublish: (id: string) => { const entry = byDevice(id); return !entry || active(entry); },
      hasUnresolved: (id: string) => { const entry = byDevice(id); return entry?.outcome.status === 'pending' || entry?.outcome.status === 'unknown'; },
      acknowledge: async (id: string, device: StudioV2WorkspaceDeviceRecord) => {
        const entry = byDevice(id);
        if (!entry || entry.outcome.status !== 'pending') return undefined;
        if (!active(entry) || !matches(device, entry)) return publish(entry, 'unknown');
        entry.acknowledgement = structuredClone(device); return readback(entry);
      },
      failed: (id: string) => { const entry = byDevice(id); if (entry?.outcome.status === 'pending') publish(entry, 'unknown'); },
    };
  }, []);
}
