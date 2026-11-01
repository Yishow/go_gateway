import type { ProtocolId } from '../state/types';
import { toBackendConnectionConfig } from '../state/studioV2DeviceAutosave';
import type { DeviceSourceBuffer, DeviceSourceOutcome, DeviceSourceRequest } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';

export interface RecentConnection {
  protocol: ProtocolId;
  config: Record<string, string | number>;
}
export interface RecentConnectionPort {
  read(workspaceId: string, protocol?: ProtocolId, subtype?: string): RecentConnection | undefined;
  clear(workspaceId: string, protocol: ProtocolId): void;
}
const network = ['host', 'port', 'timeout'];
const serial = ['serial_port', 'baud_rate', 'parity', 'data_bits', 'stop_bits', 'timeout'];
/** Canonical aliases first, then positive allowlisting. Never copy a whole config. */
export function connectionWhitelist(buffer: Pick<DeviceSourceBuffer, 'protocol' | 'config'>): RecentConnection | undefined {
  const { protocol } = buffer;
  if (protocol === 'mqtt') return;
  const config = toBackendConnectionConfig({ ...buffer, id: '', name: '', description: '', status: 'draft', test: null });
  let keys: readonly string[];
  switch (protocol) {
    case 'modbus_tcp': case 'modbus_udp': keys = [...network, 'slave_id', 'data_format']; break;
    case 'modbus_rtu': keys = [...serial, 'slave_id', 'data_format']; break;
    case 'mc_3e': keys = [...network, 'network_no', 'pc_no', 'io_no', 'station_no', 'data_format']; break;
    case 'fatek_fbs':
      if (config.mode !== 'tcp' && config.mode !== 'serial') return;
      keys = [...(config.mode === 'tcp' ? network : serial), 'mode', 'station_no']; break;
    default: return;
  }
  const eligible: RecentConnection['config'] = {};
  for (const key of keys) {
    const value = config[key];
    if (typeof value === 'string' || typeof value === 'number' && Number.isFinite(value)) eligible[key] = value;
  }
  return { protocol, config: eligible };
}
const keyOf = (protocol: ProtocolId, subtype?: string) => `${protocol}:${protocol === 'fatek_fbs' ? subtype ?? 'tcp' : ''}`;

/** One mounted page's memory only. Tickets observe source outcomes, not write authority. */
export function createRecentConnectionCache(initialWorkspace?: string) {
  let workspaceId = initialWorkspace, epoch = {}, sequence = 0;
  const generations = new Map<ProtocolId, number>();
  const values = new Map<string, { value: RecentConnection; sequence: number }>();
  const resetScope = (next?: string) => {
    if (workspaceId === next) return;
    workspaceId = next; epoch = {}; generations.clear(); values.clear();
  };
  const port: RecentConnectionPort = {
    read: (scope, protocol, subtype) => {
      if (!scope || scope !== workspaceId) return;
      const entry = protocol ? values.get(keyOf(protocol, subtype)) : [...values.values()].sort((a, b) => b.sequence - a.sequence)[0];
      return entry ? structuredClone(entry.value) : undefined;
    },
    clear: (scope, protocol) => {
      if (!scope || scope !== workspaceId) return;
      generations.set(protocol, (generations.get(protocol) ?? 0) + 1);
      for (const [key, entry] of values) if (entry.value.protocol === protocol) values.delete(key);
    },
  };
  return {
    port, resetScope,
    dispose: () => { epoch = {}; values.clear(); generations.clear(); },
    capture: (input: DeviceSourceRequest) => {
      // Keep causal identifiers and eligible values, not another copy of names,
      // credentials or unknown config already owned by the source queue.
      const { workspaceId: scope, sessionId, requestId, generation: requestGeneration, deviceId } = input;
      const protocol = input.buffer.protocol, capturedEpoch = epoch;
      const generation = generations.get(protocol) ?? 0, order = ++sequence;
      const value = connectionWhitelist(input.buffer);
      let accepted = false;
      return (outcome: DeviceSourceOutcome) => {
        if (accepted || deviceId || !value || !scope || scope !== workspaceId || capturedEpoch !== epoch ||
          generation !== (generations.get(protocol) ?? 0) || outcome.status !== 'confirmed' || !outcome.device?.id ||
          outcome.device.runtime_apply_status === 'apply_failed' || outcome.device.protocol !== protocol || outcome.request.deviceId ||
          outcome.request.workspaceId !== scope || outcome.request.sessionId !== sessionId || outcome.request.requestId !== requestId ||
          outcome.request.generation !== requestGeneration || JSON.stringify(connectionWhitelist(outcome.request.buffer)) !== JSON.stringify(value)) return;
        accepted = true;
        const key = keyOf(value.protocol, String(value.config.mode ?? 'tcp'));
        if ((values.get(key)?.sequence ?? -1) < order) values.set(key, { value, sequence: order });
      };
    },
  };
}
