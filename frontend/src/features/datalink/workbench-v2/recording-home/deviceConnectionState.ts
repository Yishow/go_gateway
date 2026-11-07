import type { DeviceCollector, DeviceConnectionState } from '../../../../types/deviceConnection';

/** A missing or untrusted observation can never become a healthy state. */
export function deriveDeviceConnectionState(collector: DeviceCollector | undefined, nowMs: number): DeviceConnectionState {
  if (!collector) return 'unavailable';
  if (collector.breaker_state !== 'closed' || collector.status === 'error') return 'disconnected';
  if (!collector.running || collector.availability_status === 'unavailable' || collector.points_total === 0) return 'idle';
  const readAt = collector.last_read_at ? Date.parse(collector.last_read_at) : NaN;
  const interval = collector.effective_interval_ms;
  const threshold = Math.max(3 * (typeof interval === 'number' && Number.isFinite(interval) && interval > 0 ? interval : 0), 5000);
  if (collector.points_error > 0 || collector.points_stale > 0 || collector.points_healthy !== collector.points_total ||
    collector.status === 'warning' || collector.last_error || !Number.isFinite(readAt) || nowMs - readAt > threshold) return 'degraded';
  return 'collecting';
}
