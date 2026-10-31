import type { DeviceSourceBuffer } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { getDefaultConfig } from '../state/protocols';
import { getProtocolSupport } from '../state/protocolSupport';
import { isStudioV2DeviceValid } from '../state/studioV2DeviceAutosave';

export function newDeviceSourceBuffer(): DeviceSourceBuffer {
  return { name: '', description: '', protocol: 'modbus_tcp', config: getDefaultConfig('modbus_tcp') };
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
  const parsed = { ...buffer, config };
  return isStudioV2DeviceValid({ ...parsed, id: '', status: 'draft', test: null }) ? parsed : undefined;
}
