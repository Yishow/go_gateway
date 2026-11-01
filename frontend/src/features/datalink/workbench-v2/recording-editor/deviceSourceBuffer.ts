import { connectionWhitelist, type RecentConnection } from './recentConnectionDefaults';
import type { ProtocolId } from '../state/types';
import { MODBUS_RTU_BAUD_RATES } from '../state/protocolConnectionOptions';
import type { DeviceSourceBuffer } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { getDefaultConfig } from '../state/protocols';
import { getProtocolSupport } from '../state/protocolSupport';
import { isStudioV2DeviceValid } from '../state/studioV2DeviceAutosave';

/** New-device defaults only: existing devices retain inactive fields in the shared form. */
export function newConnectionDefaults(protocol: ProtocolId, subtype?: string): Record<string, unknown> {
  const defaults = protocol === 'fatek_fbs' && subtype === 'serial'
    ? { mode: 'serial', station_no: 1, timeout: 5, baud_rate: 9600, parity: 'even', data_bits: 7, stop_bits: 1 }
    : getDefaultConfig(protocol);
  return connectionWhitelist({ protocol, config: defaults })?.config ?? {};
}
export function newDeviceSourceBuffer(recent?: RecentConnection): DeviceSourceBuffer {
  const protocol = recent?.protocol ?? 'modbus_tcp';
  return { name: '', description: '', protocol, config: {
    ...newConnectionDefaults(protocol, String(recent?.config.mode ?? 'tcp')), ...recent?.config,
  } };
}
/** Parse only on explicit submit; raw fields and unrendered config remain untouched. */
export function parsedDeviceSourceBuffer(buffer: DeviceSourceBuffer): DeviceSourceBuffer | undefined {
  if (!getProtocolSupport(buffer.protocol).available) return undefined;
  const config = { ...buffer.config };
  const numeric = ['timeout', 'slave_id', 'station', 'station_no', 'baud', 'baud_rate', 'network_no', 'pc_no', 'io_no', 'data_bits', 'stop_bits'];
  if (buffer.protocol !== 'modbus_rtu') numeric.push('port');
  for (const key of numeric) {
    const value = config[key];
    if (typeof value !== 'string') continue;
    if (!/^[+-]?(?:\d+(?:\.\d*)?|\.\d+)$/.test(value.trim())) return undefined;
    const parsed = Number(value);
    if (!Number.isFinite(parsed)) return undefined;
    config[key] = parsed;
  }
  if (buffer.protocol === 'modbus_rtu' && !MODBUS_RTU_BAUD_RATES.includes(Number(config.baud_rate ?? config.baud))) return undefined;
  const parsed = { ...buffer, config };
  return isStudioV2DeviceValid({ ...parsed, id: '', status: 'draft', test: null }) ? parsed : undefined;
}
