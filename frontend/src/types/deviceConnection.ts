import type { RuntimeCollectorStatus } from './datalink';

export type DeviceConnectionState = 'unavailable' | 'disconnected' | 'idle' | 'degraded' | 'collecting';
export type DeviceCollector = RuntimeCollectorStatus & { effective_interval_ms?: number };
export interface DeviceConnectionObservation {
  state: DeviceConnectionState;
  collector?: DeviceCollector;
}
