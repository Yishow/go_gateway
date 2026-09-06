import type { ProtocolId, ReadinessStage } from './types';

/**
 * 支援的六種通訊協議定義列表
 */
export const PROTOCOLS = [
  { id: 'modbus_tcp', nameKey: 'step1.protocols.modbus_tcp', descKey: 'step1.protocols.modbus_tcp_desc' },
  { id: 'modbus_rtu', nameKey: 'step1.protocols.modbus_rtu', descKey: 'step1.protocols.modbus_rtu_desc' },
  { id: 'modbus_udp', nameKey: 'step1.protocols.modbus_udp', descKey: 'step1.protocols.modbus_udp_desc' },
  { id: 'fatek_fbs', nameKey: 'step1.protocols.fatek_fbs', descKey: 'step1.protocols.fatek_fbs_desc' },
  { id: 'mc_3e', nameKey: 'step1.protocols.mc_3e', descKey: 'step1.protocols.mc_3e_desc' },
  { id: 'mqtt', nameKey: 'step1.protocols.mqtt', descKey: 'step1.protocols.mqtt_desc' },
] as const;

/**
 * 依協議取得其診斷步驟（ReadinessStage）的序列
 * 
 * @param protocol 協議 ID
 * @returns 包含 connect 或 probe 群組標記的步驟陣列
 */
export function getStagesForProtocol(protocol: ProtocolId): ReadinessStage[] {
  switch (protocol) {
    case 'modbus_rtu':
      return [
        { id: 'open_port', label: 'step1.stages.open_port', group: 'connect', status: 'pending' },
        { id: 'handshake', label: 'step1.stages.handshake', group: 'connect', status: 'pending' },
        { id: 'probe', label: 'step1.stages.probe', group: 'probe', status: 'pending' },
      ];
    case 'mqtt':
      return [
        { id: 'resolve', label: 'step1.stages.resolve', group: 'connect', status: 'pending' },
        { id: 'connect', label: 'step1.stages.connect', group: 'connect', status: 'pending' },
        { id: 'subscribe', label: 'step1.stages.subscribe', group: 'probe', status: 'pending' },
      ];
    case 'modbus_tcp':
    case 'modbus_udp':
    case 'fatek_fbs':
    case 'mc_3e':
    default:
      return [
        { id: 'resolve', label: 'step1.stages.resolve', group: 'connect', status: 'pending' },
        { id: 'connect', label: 'step1.stages.connect', group: 'connect', status: 'pending' },
        { id: 'probe', label: 'step1.stages.probe', group: 'probe', status: 'pending' },
      ];
  }
}

/**
 * 取得給定協議的預設連線設定參數
 * 
 * @param protocol 協議 ID
 * @returns 預設設定物件
 */
export function getDefaultConfig(protocol: ProtocolId): Record<string, string | number | string[]> {
  switch (protocol) {
    case 'modbus_rtu':
      return {
        port: '/dev/ttyUSB0',
        baud: 9600,
        parity: 'N',
        data_bits: 8,
        stop_bits: 1,
        slave_id: 1,
        timeout: 5,
      };
    case 'mqtt':
      return {
        broker: 'mqtt://127.0.0.1:1883',
        client_id: 'gw-01',
        username: '',
        topics: [],
      };
    case 'mc_3e':
      return {
        host: '192.168.1.100',
        port: 6000,
        station: 0,
        network_no: 0,
        pc_no: 255,
        module_io_no: 1023,
        module_station_no: 0,
        timeout: 5,
      };
    case 'fatek_fbs':
      return {
        mode: 'tcp',
        host: '192.168.1.100',
        port: 500,
        station: 1,
        timeout: 5,
      };
    case 'modbus_tcp':
    case 'modbus_udp':
    default:
      return {
        host: '192.168.1.100',
        port: 502,
        slave_id: 1,
        timeout: 5,
      };
  }
}
