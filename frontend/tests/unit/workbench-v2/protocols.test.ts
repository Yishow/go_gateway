import { describe, it, expect } from 'vitest';
import { getStagesForProtocol, getDefaultConfig, PROTOCOLS } from '../../../src/features/datalink/workbench-v2/state/protocols';

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

      // Modbus RTU 預設配置
      const rtuConfig = getDefaultConfig('modbus_rtu');
      expect(rtuConfig).toHaveProperty('port', '/dev/ttyUSB0');
      expect(rtuConfig).toHaveProperty('baud', 9600);

      // MQTT 預設配置
      const mqttConfig = getDefaultConfig('mqtt');
      expect(mqttConfig).toHaveProperty('broker', 'mqtts://broker.local:8883');
      expect(mqttConfig).toHaveProperty('client_id', 'gw-01');
    });
  });
});
