import type {
  Device,
  DeviceStatus,
  ProtocolType,
} from '../../../types/datalink';

export type DeviceConnectionValue = string | number | boolean | string[];
export type DeviceConnectionConfig = Record<
  string,
  DeviceConnectionValue | undefined
>;

export interface DeviceDraft {
  name: string;
  description: string;
  protocol: ProtocolType;
  connectionConfig: DeviceConnectionConfig;
}

export const WORKBENCH_PROTOCOLS: readonly ProtocolType[] = [
  'modbus_tcp',
  'modbus_udp',
  'modbus_rtu',
  'fatek_fbs',
  'mc_3e',
  'mqtt',
];

export const WORKBENCH_DEVICE_STATUSES: readonly DeviceStatus[] = [
  'active',
  'draft',
  'disabled',
];

const protocolLabelKeyMap: Record<ProtocolType, string> = {
  modbus_tcp: 'workbench.device.protocols.modbus_tcp',
  modbus_udp: 'workbench.device.protocols.modbus_udp',
  modbus_rtu: 'workbench.device.protocols.modbus_rtu',
  fatek_fbs: 'workbench.device.protocols.fatek_fbs',
  mc_3e: 'workbench.device.protocols.mc_3e',
  mqtt: 'workbench.device.protocols.mqtt',
};

const statusLabelKeyMap: Record<DeviceStatus, string> = {
  active: 'workbench.device.status.active',
  draft: 'workbench.device.status.draft',
  disabled: 'workbench.device.status.disabled',
};

function isConnectionObject(
  value: unknown,
): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function normalizeConnectionConfig(
  value: Record<string, unknown>,
): DeviceConnectionConfig {
  const normalized: DeviceConnectionConfig = {};

  Object.entries(value).forEach(([key, entry]) => {
    if (
      typeof entry === 'string' ||
      typeof entry === 'number' ||
      typeof entry === 'boolean'
    ) {
      normalized[key] = entry;
      return;
    }

    if (Array.isArray(entry) && entry.every((item) => typeof item === 'string')) {
      normalized[key] = entry;
    }
  });

  return normalized;
}

export function getWorkbenchProtocolLabelKey(protocol: ProtocolType): string {
  return protocolLabelKeyMap[protocol];
}

export function getWorkbenchDeviceStatusLabelKey(status: DeviceStatus): string {
  return statusLabelKeyMap[status];
}

export function createDefaultDeviceConnectionConfig(
  protocol: ProtocolType,
): DeviceConnectionConfig {
  switch (protocol) {
    case 'modbus_tcp':
      return {
        host: '',
        port: 502,
        slave_id: 1,
        timeout: 5,
      };
    case 'modbus_udp':
      return {
        host: '',
        port: 502,
        slave_id: 1,
        timeout: 5,
      };
    case 'modbus_rtu':
      return {
        serial_port: '',
        baud_rate: 9600,
        data_bits: 8,
        stop_bits: 1,
        parity: 'none',
        slave_id: 1,
        timeout: 5,
      };
    case 'fatek_fbs':
      return {
        mode: 'tcp',
        host: '',
        port: 500,
        serial_port: '',
        baud_rate: 9600,
        station_no: 1,
        timeout: 5,
      };
    case 'mc_3e':
      return {
        host: '',
        port: 5000,
        network_no: 1,
        pc_no: 255,
        io_no: 1023,
        station_no: 0,
        timeout: 5,
        data_format: 'binary',
      };
    case 'mqtt':
      return {
        broker_url: '',
        client_id: '',
        username: '',
        password: '',
        use_tls: false,
        topics: [],
        qos: 0,
      };
  }
}

export function createEmptyDeviceDraft(
  protocol: ProtocolType = 'modbus_tcp',
): DeviceDraft {
  return {
    name: '',
    description: '',
    protocol,
    connectionConfig: createDefaultDeviceConnectionConfig(protocol),
  };
}

export function parseDeviceConnectionConfig(
  rawConnectionConfig: string,
): DeviceConnectionConfig {
  if (!rawConnectionConfig) {
    return {};
  }

  try {
    const parsed = JSON.parse(rawConnectionConfig);
    if (!isConnectionObject(parsed)) {
      return {};
    }

    return normalizeConnectionConfig(parsed);
  } catch {
    return {};
  }
}

export function buildDeviceDraftFromDevice(device: Device): DeviceDraft {
  return {
    name: device.name,
    description: device.description ?? '',
    protocol: device.protocol,
    connectionConfig: {
      ...createDefaultDeviceConnectionConfig(device.protocol),
      ...parseDeviceConnectionConfig(device.connection_config),
    },
  };
}

export function sanitizeDeviceConnectionConfig(
  config: DeviceConnectionConfig,
): Record<string, unknown> {
  const sanitized: Record<string, unknown> = {};

  Object.entries(config).forEach(([key, value]) => {
    if (value === undefined) {
      return;
    }

    if (typeof value === 'string') {
      const trimmed = value.trim();
      if (trimmed === '') {
        return;
      }
      sanitized[key] = trimmed;
      return;
    }

    if (Array.isArray(value)) {
      const items = value
        .map((item) => item.trim())
        .filter((item) => item.length > 0);
      if (items.length === 0) {
        return;
      }
      sanitized[key] = items;
      return;
    }

    sanitized[key] = value;
  });

  return sanitized;
}

export function buildDeviceConnectionSummary(
  protocol: ProtocolType,
  connectionConfig: DeviceConnectionConfig,
): Array<{ labelKey: string; value: string }> {
  const stringValue = (key: string): string | null => {
    const value = connectionConfig[key];
    if (typeof value === 'string' && value.trim() !== '') {
      return value.trim();
    }
    if (typeof value === 'number') {
      return String(value);
    }
    if (typeof value === 'boolean') {
      return value ? 'true' : 'false';
    }
    if (Array.isArray(value) && value.length > 0) {
      return value.join(', ');
    }

    return null;
  };

  switch (protocol) {
    case 'modbus_tcp':
    case 'modbus_udp':
      return [
        { labelKey: 'workbench.device.connection.host', value: stringValue('host') ?? '—' },
        { labelKey: 'workbench.device.connection.port', value: stringValue('port') ?? '—' },
        {
          labelKey: 'workbench.device.connection.slaveId',
          value: stringValue('slave_id') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.timeout',
          value: stringValue('timeout') ?? '—',
        },
      ];
    case 'modbus_rtu':
      return [
        {
          labelKey: 'workbench.device.connection.serialPort',
          value: stringValue('serial_port') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.baudRate',
          value: stringValue('baud_rate') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.parity',
          value: stringValue('parity') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.slaveId',
          value: stringValue('slave_id') ?? '—',
        },
      ];
    case 'fatek_fbs':
      return [
        {
          labelKey: 'workbench.device.connection.mode',
          value: stringValue('mode') ?? '—',
        },
        {
          labelKey:
            stringValue('mode') === 'serial'
              ? 'workbench.device.connection.serialPort'
              : 'workbench.device.connection.host',
          value:
            stringValue('mode') === 'serial'
              ? stringValue('serial_port') ?? '—'
              : stringValue('host') ?? '—',
        },
        {
          labelKey:
            stringValue('mode') === 'serial'
              ? 'workbench.device.connection.baudRate'
              : 'workbench.device.connection.port',
          value:
            stringValue('mode') === 'serial'
              ? stringValue('baud_rate') ?? '—'
              : stringValue('port') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.stationNo',
          value: stringValue('station_no') ?? '—',
        },
      ];
    case 'mc_3e':
      return [
        { labelKey: 'workbench.device.connection.host', value: stringValue('host') ?? '—' },
        { labelKey: 'workbench.device.connection.port', value: stringValue('port') ?? '—' },
        {
          labelKey: 'workbench.device.connection.networkNo',
          value: stringValue('network_no') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.stationNo',
          value: stringValue('station_no') ?? '—',
        },
      ];
    case 'mqtt':
      return [
        {
          labelKey: 'workbench.device.connection.brokerUrl',
          value: stringValue('broker_url') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.clientId',
          value: stringValue('client_id') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.topics',
          value: stringValue('topics') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.useTls',
          value: stringValue('use_tls') ?? 'false',
        },
      ];
  }
}
