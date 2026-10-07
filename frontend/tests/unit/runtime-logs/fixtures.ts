import type { RuntimeLogRecord, RuntimeLogSnapshot } from '../../../src/types/runtimeLogs';
export const record = (sequence = '1', message = 'Gateway ready'): RuntimeLogRecord => ({
  instance_id: 'instance-a', sequence, timestamp: '2026-10-07T00:00:00Z', level: 'info', source: 'runtime',
  code: 'runtime.started', message, fields: {}, truncated: false,
});
export const snapshot: RuntimeLogSnapshot = { instance_id: 'instance-a', oldest_cursor: 'instance-a:1',
  latest_cursor: 'instance-a:1', records: [record()], capture_dropped: 0, evicted: 0,
  subscriber_overflow: 0, sink_health: 'healthy' };
export const filters = { level: 'debug' as const, source: '', q: '' };
