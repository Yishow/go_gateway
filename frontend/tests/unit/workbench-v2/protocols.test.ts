import { describe, it, expect } from 'vitest';
import { getStagesForProtocol, getDefaultConfig, PROTOCOLS } from '../../../src/features/datalink/workbench-v2/state/protocols';
import { isStudioV2DeviceValid, toStudioV2DeviceUpdateRequest } from '../../../src/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import type { Device } from '../../../src/features/datalink/workbench-v2/state/types';

function device(protocol: Device['protocol'], config: Device['config']): Device {
  return {
    id: 'dev-01',
    name: 'PLC',
    description: '',
    protocol,
    config,
    status: 'draft',
    test: null,
  };
}

describe('Protocols State Helpers', () => {
  it('應正確匯出 6 種支援的協議資料列表', () => {
    expect(PROTOCOLS).toHaveLength(6);
    const ids = PROTOCOLS.map((p) => p.id);
    expect(ids).toContain('modbus_tcp');
    expect(ids).toContain('modbus_rtu');
    expect(ids).toContain('modbus_udp');
    expect(ids).toContain('fatek_fbs');
    expect(ids).toContain('mc_3e');
    expect(ids).toContain('mqtt');
  });

  describe('getStagesForProtocol', () => {
    it('Modbus TCP 應包含 resolve (connect), connect (connect), probe (probe) 三個步驟', () => {
      const stages = getStagesForProtocol('modbus_tcp');
      expect(stages).toEqual([
        { id: 'resolve', label: 'step1.stages.resolve', group: 'connect', status: 'pending' },
        { id: 'connect', label: 'step1.stages.connect', group: 'connect', status: 'pending' },
        { id: 'probe', label: 'step1.stages.probe', group: 'probe', status: 'pending' },
      ]);
    });

    it('Modbus RTU 應包含 open_port (connect), handshake (connect), probe (probe) 三個步驟', () => {
      const stages = getStagesForProtocol('modbus_rtu');
      expect(stages).toEqual([
        { id: 'open_port', label: 'step1.stages.open_port', group: 'connect', status: 'pending' },
        { id: 'handshake', label: 'step1.stages.handshake', group: 'connect', status: 'pending' },
        { id: 'probe', label: 'step1.stages.probe', group: 'probe', status: 'pending' },
      ]);
    });

    it('MQTT 應包含 resolve (connect), connect (connect), subscribe (probe) 三個步驟', () => {
      const stages = getStagesForProtocol('mqtt');
      expect(stages).toEqual([
        { id: 'resolve', label: 'step1.stages.resolve', group: 'connect', status: 'pending' },
        { id: 'connect', label: 'step1.stages.connect', group: 'connect', status: 'pending' },
        { id: 'subscribe', label: 'step1.stages.subscribe', group: 'probe', status: 'pending' },
      ]);
    });

    it('MC 3E / Fatek FBS 應與 TCP 套用相同步驟', () => {
      expect(getStagesForProtocol('mc_3e')).toEqual(getStagesForProtocol('modbus_tcp'));
      expect(getStagesForProtocol('fatek_fbs')).toEqual(getStagesForProtocol('modbus_tcp'));
    });
  });

  describe('getDefaultConfig', () => {
    it('應能正確回傳不同通訊協定的預設配置', () => {
      // Modbus TCP 預設配置
      const tcpConfig = getDefaultConfig('modbus_tcp');
      expect(tcpConfig).toHaveProperty('host', '192.168.1.100');
      expect(tcpConfig).toHaveProperty('port', 502);
      expect(tcpConfig).toHaveProperty('slave_id', 1);

      // Modbus UDP 預設配置
      const udpConfig = getDefaultConfig('modbus_udp');
      expect(udpConfig).toHaveProperty('host', '192.168.1.100');
      expect(udpConfig).toHaveProperty('port', 502);
      expect(udpConfig).toHaveProperty('slave_id', 1);

      // Modbus RTU 預設配置
      const rtuConfig = getDefaultConfig('modbus_rtu');
      expect(rtuConfig).toHaveProperty('port', '/dev/ttyUSB0');
      expect(rtuConfig).toHaveProperty('baud', 9600);
      expect(rtuConfig).toHaveProperty('parity', 'N');

      // MC 3E 預設配置 (Port 6000, Station 0)
      const mcConfig = getDefaultConfig('mc_3e');
      expect(mcConfig).toHaveProperty('host', '192.168.1.100');
      expect(mcConfig).toHaveProperty('port', 6000);
      expect(mcConfig).toHaveProperty('station', 0);

      // FATEK FBS 預設配置 (Port 500, Station 1)
      const fatekConfig = getDefaultConfig('fatek_fbs');
      expect(fatekConfig).toHaveProperty('host', '192.168.1.100');
      expect(fatekConfig).toHaveProperty('port', 500);
      expect(fatekConfig).toHaveProperty('station', 1);

      // MQTT 預設配置
      const mqttConfig = getDefaultConfig('mqtt');
      expect(mqttConfig).toHaveProperty('broker', 'mqtt://127.0.0.1:1883');
      expect(mqttConfig).toHaveProperty('client_id', 'gw-01');
    });
  });

  describe('backend connection config serialization', () => {
    it('serializes all six protocols with the backend registry field names', () => {
      expect(toStudioV2DeviceUpdateRequest(device('modbus_tcp', {
        host: 'plc-tcp', port: 502, slave_id: 1, timeout: 5,
      })).connection_config).toEqual({
        host: 'plc-tcp', port: 502, slave_id: 1, timeout: 5,
      });
      expect(toStudioV2DeviceUpdateRequest(device('modbus_udp', {
        host: 'plc-udp', port: 502, slave_id: 2, timeout: 5,
      })).connection_config).toEqual({
        host: 'plc-udp', port: 502, slave_id: 2, timeout: 5,
      });
      expect(toStudioV2DeviceUpdateRequest(device('modbus_rtu', {
        port: 'COM1', baud: 19200, parity: 'E', slave_id: 3,
      })).connection_config).toEqual({
        serial_port: 'COM1', baud_rate: 19200, parity: 'even', slave_id: 3,
      });
      expect(toStudioV2DeviceUpdateRequest(device('fatek_fbs', {
        host: 'plc-fatek', port: 500, station: 4, timeout: 5,
      })).connection_config).toEqual({
        mode: 'tcp', host: 'plc-fatek', port: 500, station_no: 4, timeout: 5,
      });
      expect(toStudioV2DeviceUpdateRequest(device('mc_3e', {
        host: 'plc-mc', port: 5000, station: 5, module_io_no: 1023, timeout: 5,
      })).connection_config).toEqual({
        host: 'plc-mc', port: 5000, station_no: 5, io_no: 1023, timeout: 5,
      });
      expect(toStudioV2DeviceUpdateRequest(device('mqtt', {
        broker: 'mqtt://broker', client_id: 'gw-01', username: '', topics: ['line/#'],
      })).connection_config).toEqual({
        broker_url: 'mqtt://broker', client_id: 'gw-01', username: '', topics: ['line/#'],
      });
    });

    it('validates the canonical fields required by each backend connector', () => {
      expect(isStudioV2DeviceValid(device('modbus_tcp', {
        host: 'plc', port: 502, slave_id: 1, timeout: 5,
      }))).toBe(true);
      expect(isStudioV2DeviceValid(device('modbus_udp', {
        host: 'plc', port: 502, slave_id: 1, timeout: 5,
      }))).toBe(true);
      expect(isStudioV2DeviceValid(device('modbus_rtu', {
        serial_port: 'COM1', baud_rate: 9600, parity: 'none', slave_id: 1,
      }))).toBe(true);
      expect(isStudioV2DeviceValid(device('fatek_fbs', {
        mode: 'tcp', host: 'plc', port: 500, station_no: 1, timeout: 5,
      }))).toBe(true);
      expect(isStudioV2DeviceValid(device('mc_3e', {
        host: 'plc', port: 5000, station_no: 0, io_no: 1023, timeout: 5,
      }))).toBe(true);
      expect(isStudioV2DeviceValid(device('mqtt', {
        broker_url: 'mqtt://broker', client_id: 'gw-01', topics: [],
      }))).toBe(true);
      expect(isStudioV2DeviceValid(device('mqtt', {
        broker_url: 'mqtt://broker', client_id: 'gw-01', topics: 'line/#',
      }))).toBe(false);
      expect(isStudioV2DeviceValid(device('modbus_rtu', {
        serial_port: 'COM1', baud_rate: 9600, parity: 'invalid', slave_id: 1,
      }))).toBe(false);
    });
  });
});
