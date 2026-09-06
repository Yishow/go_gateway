import type {
  CreateDeviceRequest,
  UpdateDeviceRequest,
} from '../../../../types/datalink';
import type { StudioV2WorkspaceDeviceRecord } from '../../../../services/studioV2WorkspaceDevices';
import type { Device, DeviceTest, ProtocolId } from './types';
import { getDefaultConfig } from './protocols';

function asNonEmptyString(value: unknown): string | null {
  return typeof value === 'string' && value.trim() !== '' ? value.trim() : null;
}

function asPositiveNumber(value: unknown): number | null {
  return typeof value === 'number' && Number.isFinite(value) && value > 0 ? value : null;
}

function asNonNegativeNumber(value: unknown): number | null {
  return typeof value === 'number' && Number.isFinite(value) && value >= 0 ? value : null;
}

function isParity(value: unknown): boolean {
  return value === 'N' || value === 'E' || value === 'O' || value === 'none' || value === 'even' || value === 'odd';
}

/**
 * 補齊載入時缺漏的站號 / 從站 ID，讓表單顯示值與 isStudioV2DeviceValid 判定
 * 使用同一份設定資料；預設值來源沿用 protocols.ts 的協議預設表。
 */
function backfillStationDefaults(
  protocol: ProtocolId,
  config: Record<string, unknown>,
): Record<string, unknown> {
  const defaults = getDefaultConfig(protocol);
  if (protocol === 'mc_3e' || protocol === 'fatek_fbs') {
    const stationNo = config.station_no ?? config.station ?? config.slave_id ?? config.module_station_no;
    if (stationNo === undefined && typeof defaults.station === 'number') {
      return { ...config, station_no: defaults.station };
    }
    return config;
  }
  if (protocol === 'modbus_tcp' || protocol === 'modbus_udp') {
    if (config.slave_id === undefined && typeof defaults.slave_id === 'number') {
      return { ...config, slave_id: defaults.slave_id };
    }
  }
  return config;
}

function normalizeParity(value: unknown): unknown {
  switch (value) {
    case 'N':
    case 'n':
    case 'none':
      return 'none';
    case 'E':
    case 'e':
    case 'even':
      return 'even';
    case 'O':
    case 'o':
    case 'odd':
      return 'odd';
    default:
      return value;
  }
}

/** Convert the form's display aliases into the canonical connector registry shape. */
export function toBackendConnectionConfig(device: Device): Record<string, unknown> {
  const config = { ...device.config };

  switch (device.protocol) {
    case 'modbus_rtu': {
      const serialPort = config.serial_port ?? config.port;
      const baudRate = config.baud_rate ?? config.baud;
      delete config.port;
      delete config.baud;
      if (serialPort !== undefined) config.serial_port = serialPort;
      if (baudRate !== undefined) config.baud_rate = baudRate;
      config.parity = normalizeParity(config.parity);
      return config;
    }
    case 'fatek_fbs': {
      const stationNo = config.station_no ?? config.station ?? config.slave_id;
      delete config.station;
      delete config.slave_id;
      config.mode = config.mode ?? 'tcp';
      if (stationNo !== undefined) config.station_no = stationNo;
      return config;
    }
    case 'mc_3e': {
      const stationNo = config.station_no ?? config.station ?? config.module_station_no;
      const ioNo = config.io_no ?? config.module_io_no;
      delete config.station;
      delete config.slave_id;
      delete config.module_io_no;
      delete config.module_station_no;
      if (stationNo !== undefined) config.station_no = stationNo;
      if (ioNo !== undefined) config.io_no = ioNo;
      return config;
    }
    case 'mqtt': {
      const brokerURL = config.broker_url ?? config.broker;
      delete config.broker;
      if (brokerURL !== undefined) config.broker_url = brokerURL;
      if (!Array.isArray(config.topics)) config.topics = [];
      return config;
    }
    default:
      return config;
  }
}

export function isStudioV2DeviceValid(device: Device): boolean {
  if (!asNonEmptyString(device.name)) {
    return false;
  }

  if (device.protocol === 'modbus_rtu') {
    return Boolean(
      asNonEmptyString(device.config.serial_port ?? device.config.port) &&
      asPositiveNumber(device.config.baud_rate ?? device.config.baud) &&
      isParity(device.config.parity) &&
      asPositiveNumber(device.config.slave_id),
    );
  }

  if (device.protocol === 'fatek_fbs') {
    const mode = asNonEmptyString(device.config.mode) ?? 'tcp';
    const stationNo = device.config.station_no ?? device.config.station ?? device.config.slave_id;
    if (mode === 'serial') {
      return Boolean(
        asNonEmptyString(device.config.serial_port) &&
        asPositiveNumber(device.config.baud_rate) &&
        isParity(device.config.parity) &&
        asNonNegativeNumber(stationNo) !== null,
      );
    }
    if (mode !== 'tcp') return false;
    return Boolean(
      asNonEmptyString(device.config.host) &&
      asPositiveNumber(device.config.port) &&
      asNonNegativeNumber(stationNo) !== null &&
      asPositiveNumber(device.config.timeout),
    );
  }

  if (device.protocol === 'mqtt') {
    return Boolean(
      asNonEmptyString(device.config.broker_url ?? device.config.broker) &&
      asNonEmptyString(device.config.client_id) &&
      (device.config.topics === undefined || Array.isArray(device.config.topics)),
    );
  }

  if (device.protocol === 'mc_3e') {
    return Boolean(
      asNonEmptyString(device.config.host) &&
      asPositiveNumber(device.config.port) &&
      asNonNegativeNumber(device.config.station_no ?? device.config.station) !== null &&
      asPositiveNumber(device.config.timeout),
    );
  }

  return Boolean(
    asNonEmptyString(device.config.host) &&
    asPositiveNumber(device.config.port) &&
    asPositiveNumber(device.config.slave_id) &&
    asPositiveNumber(device.config.timeout),
  );
}

export function toStudioV2DeviceCreateRequest(device: Device): CreateDeviceRequest {
  return {
    id: device.id,
    name: device.name,
    description: device.description,
    protocol: device.protocol,
    connection_config: toBackendConnectionConfig(device),
  };
}

export function toStudioV2DeviceUpdateRequest(
  device: Device,
): UpdateDeviceRequest & { protocol: ProtocolId } {
  return {
    name: device.name,
    description: device.description,
    protocol: device.protocol,
    connection_config: toBackendConnectionConfig(device),
  };
}

export function hydrateStudioV2Device(record: StudioV2WorkspaceDeviceRecord): Device {
  let config: Record<string, unknown> = {};

  try {
    const parsed = JSON.parse(record.connection_config);
    if (parsed && typeof parsed === 'object') {
      config = parsed as Record<string, unknown>;
    }
  } catch {
    config = {};
  }

  const persistedStages: DeviceTest['stages'] = record.last_test_success
    ? {
        connect: { status: 'success' },
        probe: { status: 'success' },
      }
    : {};
  const test: DeviceTest | null = record.last_test_success === null
    ? null
    : {
        status: record.last_test_success ? 'success' : 'failed',
        stages: persistedStages,
        tested_at: record.last_test_at ?? undefined,
      };

  const protocol = record.protocol as ProtocolId;

  return {
    id: record.id,
    name: record.name,
    description: record.description ?? '',
    protocol,
    config: backfillStationDefaults(protocol, config),
    status: record.status === 'active' ? 'active' : record.last_test_success ? 'tested' : 'draft',
    test,
    persisted: true,
    save_state: 'saved',
    save_error: null,
    runtime_apply_status: null,
    runtime_apply_message: null,
    availability_status: record.availability_status ?? 'available',
    availability_reason: record.availability_reason ?? null,
    running: record.running ?? false,
  };
}
