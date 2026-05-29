import { assertType } from 'vitest';
import type { Device, DeviceTest, ProtocolId, ReadinessStage } from '../../../src/features/datalink/workbench-v2/state/types';

/**
 * 測試 ProtocolId 型別的靜態斷言
 */
assertType<ProtocolId>('modbus_tcp');
assertType<ProtocolId>('modbus_rtu');
assertType<ProtocolId>('mqtt');

/**
 * 測試 ReadinessStage 結構的靜態斷言
 */
const dummyStage: ReadinessStage = {
  id: 'resolve',
  label: 'Resolve IP',
  group: 'connect',
  status: 'success',
  latency_ms: 12,
};
assertType<ReadinessStage>(dummyStage);

/**
 * 測試 DeviceTest 結構的靜態斷言，特別是 stages 映射
 */
const dummyTest: DeviceTest = {
  status: 'success',
  latency_ms: 50,
  stages: {
    resolve: { status: 'success', latency_ms: 10 },
    connect: { status: 'success', latency_ms: 15 },
    probe: { status: 'success', latency_ms: 25 },
  },
  tested_at: '2026-05-29T03:00:00.000Z',
};
assertType<DeviceTest>(dummyTest);

/**
 * 測試 Device 結構的靜態斷言
 */
const dummyDevice: Device = {
  id: 'dev-01',
  name: 'PLC-1',
  description: 'Production line 1 PLC',
  protocol: 'modbus_tcp',
  config: { host: '192.168.1.100', port: 502 },
  status: 'tested',
  test: dummyTest,
};
assertType<Device>(dummyDevice);
