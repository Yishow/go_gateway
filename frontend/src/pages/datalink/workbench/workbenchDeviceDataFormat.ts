import type { ProtocolType } from '../../../types/datalink';

const MODBUS_PROTOCOLS = new Set<ProtocolType>([
  'modbus_tcp',
  'modbus_udp',
  'modbus_rtu',
]);
const CONNECTION_DATA_FORMATS = new Set(['ABCD', 'BADC', 'CDAB', 'DCBA']);

export const MODBUS_DEFAULT_DATA_FORMAT = 'ABCD';
export const MC3E_DEFAULT_DATA_FORMAT = 'CDAB';

export function supportsConnectionDataFormat(protocol: ProtocolType): boolean {
  return protocol === 'mc_3e' || MODBUS_PROTOCOLS.has(protocol);
}

export function normalizeConnectionDataFormatValue(
  protocol: ProtocolType,
  value: unknown,
): string {
  if (typeof value === 'string') {
    const upper = value.trim().toUpperCase();
    if (CONNECTION_DATA_FORMATS.has(upper)) {
      return upper;
    }
  }

  if (protocol === 'mc_3e') {
    return MC3E_DEFAULT_DATA_FORMAT;
  }
  if (MODBUS_PROTOCOLS.has(protocol)) {
    return MODBUS_DEFAULT_DATA_FORMAT;
  }
  return '';
}

export function buildConnectionDataFormatOptions(
  t: (key: string) => string,
): Array<{ value: string; label: string }> {
  return [
    {
      value: 'ABCD',
      label: t('workbench.device.connection.dataFormats.abcd'),
    },
    {
      value: 'BADC',
      label: t('workbench.device.connection.dataFormats.badc'),
    },
    {
      value: 'CDAB',
      label: t('workbench.device.connection.dataFormats.cdab'),
    },
    {
      value: 'DCBA',
      label: t('workbench.device.connection.dataFormats.dcba'),
    },
  ];
}
