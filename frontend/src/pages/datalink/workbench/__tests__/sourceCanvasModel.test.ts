import { describe, expect, it } from 'vitest';
import type { Point } from '../../../../types/datalink';
import {
  buildAddressCanvasItems,
  buildPlannedPointAddresses,
} from '../sourceCanvasModel';

function createPoint(overrides: Partial<Point>): Point {
  return {
    id: 'point-1',
    device_id: 'device-1',
    name: 'Existing Point',
    description: '',
    data_type: 'int16',
    address: '40001',
    enabled: true,
    polling_group_id: '',
    last_value: null,
    last_read_at: '',
    last_error: '',
    error_count: 0,
    created_at: '',
    updated_at: '',
    ...overrides,
  };
}

describe('sourceCanvasModel', () => {
  it('builds logical base addresses for wide data types', () => {
    expect(
      buildPlannedPointAddresses({
        startAddress: '40001',
        count: 2,
        dataType: 'float32',
        protocol: 'modbus_tcp',
      }),
    ).toEqual(['40001', '40003']);
  });

  it('expands existing wide points across all occupied cells', () => {
    const items = buildAddressCanvasItems({
      points: [createPoint({ address: '40005', data_type: 'int32' })],
      plannedPointAddresses: ['40009'],
      plannedDataType: 'int16',
      protocol: 'modbus_tcp',
    });

    expect(items.find((item) => item.address === '40005')?.status).toBe('used');
    expect(items.find((item) => item.address === '40006')?.status).toBe('used');
  });

  it('marks overlapping planned cells as conflict', () => {
    const items = buildAddressCanvasItems({
      points: [createPoint({ address: '40002', data_type: 'int16' })],
      plannedPointAddresses: ['40001'],
      plannedDataType: 'float32',
      protocol: 'modbus_tcp',
    });

    expect(items.find((item) => item.address === '40001')?.status).toBe('planned');
    expect(items.find((item) => item.address === '40002')?.status).toBe('conflict');
  });
});
