import { describe, expect, it } from 'vitest';
import type { Point, PollingGroup } from '../../../types/datalink';
import { estimatePollingLoadDelta } from '../pollingLoadEstimate';

const pollingGroups: PollingGroup[] = [
  {
    id: 'pg-fast',
    name: 'fast',
    interval_ms: 500,
    priority: 1,
    enabled: true,
    created_at: '',
    updated_at: '',
  },
  {
    id: 'pg-slow',
    name: 'slow',
    interval_ms: 1000,
    priority: 2,
    enabled: true,
    created_at: '',
    updated_at: '',
  },
];

const points: Point[] = [
  {
    id: 'p1',
    device_id: 'd1',
    name: 'P1',
    description: '',
    data_type: 'int16',
    address: '40001',
    enabled: true,
    polling_group_id: 'pg-fast',
    last_value: null,
    last_read_at: '',
    last_error: '',
    error_count: 0,
    created_at: '',
    updated_at: '',
  },
  {
    id: 'p2',
    device_id: 'd1',
    name: 'P2',
    description: '',
    data_type: 'int16',
    address: '40002',
    enabled: true,
    polling_group_id: 'pg-slow',
    last_value: null,
    last_read_at: '',
    last_error: '',
    error_count: 0,
    created_at: '',
    updated_at: '',
  },
];

describe('pollingLoadEstimate', () => {
  it('estimates baseline, delta, and projected reads per second', () => {
    const result = estimatePollingLoadDelta(points, pollingGroups, 3);
    expect(result).toEqual({
      baselineReadsPerSec: 3,
      deltaReadsPerSec: 6,
      projectedReadsPerSec: 9,
      assumedIntervalMs: 500,
    });
  });

  it('falls back to default interval when no polling groups are enabled', () => {
    const result = estimatePollingLoadDelta(points, [], 2);
    expect(result).toEqual({
      baselineReadsPerSec: 0,
      deltaReadsPerSec: 2,
      projectedReadsPerSec: 2,
      assumedIntervalMs: 1000,
    });
  });

  it('excludes points that are bound to disabled polling groups from baseline', () => {
    const result = estimatePollingLoadDelta(
      [
        ...points,
        {
          ...points[0],
          id: 'p3',
          polling_group_id: 'pg-disabled',
        },
      ],
      [
        ...pollingGroups,
        {
          id: 'pg-disabled',
          name: 'disabled',
          interval_ms: 200,
          priority: 3,
          enabled: false,
          created_at: '',
          updated_at: '',
        },
      ],
      1
    );
    expect(result.baselineReadsPerSec).toBe(3);
  });
});
