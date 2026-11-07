import { describe, expect, it } from 'vitest';
import type { RuntimeCollectorStatus } from '@/types/datalink';
import { deriveDeviceConnectionState } from '@/features/datalink/workbench-v2/recording-home/deviceConnectionState';

type Collector = RuntimeCollectorStatus & { effective_interval_ms?: number };
const now = Date.parse('2026-11-07T12:00:10Z');
const healthy: Collector = {
  device_id: 'MC', device_name: 'MC', protocol: 'mc', status: 'running', availability_status: 'available',
  running: true, points_total: 48, points_healthy: 48, points_stale: 0, points_error: 0,
  last_read_at: '2026-11-07T12:00:09Z', last_error: null, breaker_state: 'closed', effective_interval_ms: 1000,
};

describe('device connection states', () => {
  it('recognizes a running collector with all points healthy and a recent read', () => {
    expect(deriveDeviceConnectionState(healthy, now)).toBe('collecting');
  });
  it.each([
    ['stale points and warning', { points_healthy: 46, points_stale: 2, status: 'warning' }, 'degraded'],
    ['read six seconds ago', { last_read_at: '2026-11-07T12:00:04Z' }, 'degraded'],
    ['read exactly five seconds ago', { last_read_at: '2026-11-07T12:00:05Z' }, 'collecting'],
    ['not running', { running: false, status: 'idle' }, 'idle'],
    ['open breaker even when idle', { breaker_state: 'open', running: false }, 'disconnected'],
    ['error status', { status: 'error' }, 'disconnected'],
    ['unavailable setting', { availability_status: 'unavailable' }, 'idle'],
    ['no points', { points_total: 0 }, 'idle'],
    ['error points', { points_error: 1 }, 'degraded'],
    ['last error', { last_error: 'error' }, 'degraded'],
    ['missing timestamp', { last_read_at: null }, 'degraded'],
    ['invalid timestamp', { last_read_at: 'bad' }, 'degraded'],
    ['longer polling interval within threshold', { effective_interval_ms: 3000, last_read_at: '2026-11-07T12:00:02Z' }, 'collecting'],
    ['longer polling interval beyond threshold', { effective_interval_ms: 3000, last_read_at: '2026-11-07T12:00:00Z' }, 'degraded'],
    ['unaccounted unhealthy points', { points_healthy: 47 }, 'degraded'],
  ] as const)('%s', (_label, change, expected) => {
    expect(deriveDeviceConnectionState({ ...healthy, ...change }, now)).toBe(expected);
  });
  it('treats an absent collector as unavailable', () => {
    expect(deriveDeviceConnectionState(undefined, now)).toBe('unavailable');
  });
});
