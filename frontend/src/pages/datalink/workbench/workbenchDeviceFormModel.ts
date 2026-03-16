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

type TranslateFn = (key: string) => string;

export interface DeviceDraft {
  name: string;
  description: string;
  protocol: ProtocolType;
  connectionConfig: DeviceConnectionConfig;
}

export type DeviceCapabilityItem = {
  id: 'address-base' | 'word-order' | 'unit-id' | 'protocol-traits';
  labelKey: string;
  value: string;
};

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

function readConnectionValueAsString(
  connectionConfig: DeviceConnectionConfig,
  key: string,
): string | null {
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
}

export function buildDeviceCapabilitySummary(
  protocol: ProtocolType,
  connectionConfig: DeviceConnectionConfig,
  t: TranslateFn,
): DeviceCapabilityItem[] {
  const notSpecified = t('workbench.device.capability.values.notSpecified');
  const notApplicable = t('workbench.device.capability.values.notApplicable');

  const unitId = (() => {
    switch (protocol) {
      case 'modbus_tcp':
      case 'modbus_udp':
      case 'modbus_rtu':
        return readConnectionValueAsString(connectionConfig, 'slave_id') ?? notSpecified;
      case 'fatek_fbs':
      case 'mc_3e':
        return readConnectionValueAsString(connectionConfig, 'station_no') ?? notSpecified;
      case 'mqtt':
        return notApplicable;
    }
  })();

  const addressBase = (() => {
    switch (protocol) {
      case 'modbus_tcp':
      case 'modbus_udp':
      case 'modbus_rtu':
        return t('workbench.device.capability.values.modbusRegister');
      case 'mqtt':
        return t('workbench.device.capability.values.topicBased');
      case 'fatek_fbs':
      case 'mc_3e':
        return t('workbench.device.capability.values.protocolNative');
    }
  })();

  const wordOrder = (() => {
    switch (protocol) {
      case 'mc_3e':
        return readConnectionValueAsString(connectionConfig, 'data_format')
          ?? t('workbench.device.capability.values.protocolDefault');
      case 'mqtt':
        return notApplicable;
      case 'modbus_tcp':
      case 'modbus_udp':
      case 'modbus_rtu':
      case 'fatek_fbs':
        return t('workbench.device.capability.values.protocolDefault');
    }
  })();

  const protocolTraits = (() => {
    switch (protocol) {
      case 'modbus_tcp':
      case 'modbus_udp': {
        const host = readConnectionValueAsString(connectionConfig, 'host') ?? notSpecified;
        const port = readConnectionValueAsString(connectionConfig, 'port') ?? notSpecified;
        return `${host}:${port}`;
      }
      case 'modbus_rtu': {
        const serialPort =
          readConnectionValueAsString(connectionConfig, 'serial_port') ?? notSpecified;
        const baudRate =
          readConnectionValueAsString(connectionConfig, 'baud_rate') ?? notSpecified;
        return `${serialPort} · ${baudRate}`;
      }
      case 'fatek_fbs': {
        const mode = readConnectionValueAsString(connectionConfig, 'mode') ?? 'tcp';
        if (mode === 'serial') {
          const serialPort =
            readConnectionValueAsString(connectionConfig, 'serial_port') ?? notSpecified;
          const baudRate =
            readConnectionValueAsString(connectionConfig, 'baud_rate') ?? notSpecified;
          return `${t('workbench.device.mode.serial')} · ${serialPort} · ${baudRate}`;
        }
        const host = readConnectionValueAsString(connectionConfig, 'host') ?? notSpecified;
        const port = readConnectionValueAsString(connectionConfig, 'port') ?? notSpecified;
        return `${t('workbench.device.mode.tcp')} · ${host}:${port}`;
      }
      case 'mc_3e': {
        const networkNo =
          readConnectionValueAsString(connectionConfig, 'network_no') ?? notSpecified;
        const stationNo =
          readConnectionValueAsString(connectionConfig, 'station_no') ?? notSpecified;
        const dataFormat =
          readConnectionValueAsString(connectionConfig, 'data_format')
          ?? t('workbench.device.capability.values.protocolDefault');
        return `${networkNo}/${stationNo} · ${dataFormat}`;
      }
      case 'mqtt': {
        const qos = readConnectionValueAsString(connectionConfig, 'qos') ?? '0';
        const useTLS = connectionConfig.use_tls === true
          ? t('workbench.device.capability.values.tlsEnabled')
          : t('workbench.device.capability.values.tlsDisabled');
        return `QoS ${qos} · ${useTLS}`;
      }
    }
  })();

  return [
    {
      id: 'unit-id',
      labelKey: 'workbench.device.capability.labels.unitId',
      value: unitId,
    },
    {
      id: 'address-base',
      labelKey: 'workbench.device.capability.labels.addressBase',
      value: addressBase,
    },
    {
      id: 'word-order',
      labelKey: 'workbench.device.capability.labels.wordOrder',
      value: wordOrder,
    },
    {
      id: 'protocol-traits',
      labelKey: 'workbench.device.capability.labels.protocolTraits',
      value: protocolTraits,
    },
  ];
}

export function getDeviceTestTimestampLabel(
  rawTimestamp: string | null,
  t: TranslateFn,
) {
  if (
    !rawTimestamp ||
    rawTimestamp.trim() === '' ||
    rawTimestamp === '0001-01-01T00:00:00Z'
  ) {
    return t('workbench.device.inspector.unknownTestTime');
  }

  return rawTimestamp;
}

export function buildDeviceConnectionSummary(
  protocol: ProtocolType,
  connectionConfig: DeviceConnectionConfig,
): Array<{ labelKey: string; value: string }> {
  switch (protocol) {
    case 'modbus_tcp':
    case 'modbus_udp':
      return [
        { labelKey: 'workbench.device.connection.host', value: readConnectionValueAsString(connectionConfig, 'host') ?? '—' },
        { labelKey: 'workbench.device.connection.port', value: readConnectionValueAsString(connectionConfig, 'port') ?? '—' },
        {
          labelKey: 'workbench.device.connection.slaveId',
          value: readConnectionValueAsString(connectionConfig, 'slave_id') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.timeout',
          value: readConnectionValueAsString(connectionConfig, 'timeout') ?? '—',
        },
      ];
    case 'modbus_rtu':
      return [
        {
          labelKey: 'workbench.device.connection.serialPort',
          value: readConnectionValueAsString(connectionConfig, 'serial_port') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.baudRate',
          value: readConnectionValueAsString(connectionConfig, 'baud_rate') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.parity',
          value: readConnectionValueAsString(connectionConfig, 'parity') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.slaveId',
          value: readConnectionValueAsString(connectionConfig, 'slave_id') ?? '—',
        },
      ];
    case 'fatek_fbs':
      return [
        {
          labelKey: 'workbench.device.connection.mode',
          value: readConnectionValueAsString(connectionConfig, 'mode') ?? '—',
        },
        {
          labelKey:
            readConnectionValueAsString(connectionConfig, 'mode') === 'serial'
              ? 'workbench.device.connection.serialPort'
              : 'workbench.device.connection.host',
          value:
            readConnectionValueAsString(connectionConfig, 'mode') === 'serial'
              ? readConnectionValueAsString(connectionConfig, 'serial_port') ?? '—'
              : readConnectionValueAsString(connectionConfig, 'host') ?? '—',
        },
        {
          labelKey:
            readConnectionValueAsString(connectionConfig, 'mode') === 'serial'
              ? 'workbench.device.connection.baudRate'
              : 'workbench.device.connection.port',
          value:
            readConnectionValueAsString(connectionConfig, 'mode') === 'serial'
              ? readConnectionValueAsString(connectionConfig, 'baud_rate') ?? '—'
              : readConnectionValueAsString(connectionConfig, 'port') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.stationNo',
          value: readConnectionValueAsString(connectionConfig, 'station_no') ?? '—',
        },
      ];
    case 'mc_3e':
      return [
        { labelKey: 'workbench.device.connection.host', value: readConnectionValueAsString(connectionConfig, 'host') ?? '—' },
        { labelKey: 'workbench.device.connection.port', value: readConnectionValueAsString(connectionConfig, 'port') ?? '—' },
        {
          labelKey: 'workbench.device.connection.networkNo',
          value: readConnectionValueAsString(connectionConfig, 'network_no') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.stationNo',
          value: readConnectionValueAsString(connectionConfig, 'station_no') ?? '—',
        },
      ];
    case 'mqtt':
      return [
        {
          labelKey: 'workbench.device.connection.brokerUrl',
          value: readConnectionValueAsString(connectionConfig, 'broker_url') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.clientId',
          value: readConnectionValueAsString(connectionConfig, 'client_id') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.topics',
          value: readConnectionValueAsString(connectionConfig, 'topics') ?? '—',
        },
        {
          labelKey: 'workbench.device.connection.useTls',
          value: readConnectionValueAsString(connectionConfig, 'use_tls') ?? 'false',
        },
      ];
  }
}

function joinEndpointParts(
  primary: string | null | undefined,
  secondary?: string | null | undefined,
  separator = ' · ',
) {
  const left = primary && primary.trim() !== '' ? primary : '—';
  const right = secondary && secondary.trim() !== '' ? secondary : null;

  if (!right) {
    return left;
  }

  if (left === '—') {
    return right;
  }

  return `${left}${separator}${right}`;
}

export function buildDeviceEndpointSummary(
  protocol: ProtocolType,
  connectionConfig: DeviceConnectionConfig,
) {
  switch (protocol) {
    case 'modbus_tcp':
    case 'modbus_udp':
    case 'mc_3e':
      return joinEndpointParts(
        readConnectionValueAsString(connectionConfig, 'host'),
        readConnectionValueAsString(connectionConfig, 'port'),
        ':',
      );
    case 'modbus_rtu':
      return joinEndpointParts(
        readConnectionValueAsString(connectionConfig, 'serial_port'),
        readConnectionValueAsString(connectionConfig, 'baud_rate'),
      );
    case 'fatek_fbs':
      if (readConnectionValueAsString(connectionConfig, 'mode') === 'serial') {
        return joinEndpointParts(
          readConnectionValueAsString(connectionConfig, 'serial_port'),
          readConnectionValueAsString(connectionConfig, 'baud_rate'),
        );
      }
      return joinEndpointParts(
        readConnectionValueAsString(connectionConfig, 'host'),
        readConnectionValueAsString(connectionConfig, 'port'),
        ':',
      );
    case 'mqtt':
      return readConnectionValueAsString(connectionConfig, 'broker_url') ?? '—';
  }
}
