import type { RuntimeLogEvent, RuntimeLogMetadata, RuntimeLogRecord, RuntimeLogSnapshot } from '../../../src/types/runtimeLogs';
import { domainInstant, domainRecord } from '../../unit/runtime-logs/domain-event-fixtures';
export { domainRecord } from '../../unit/runtime-logs/domain-event-fixtures';
export { legacyRecord } from '../../unit/runtime-logs/legacy-event-fixtures';
export type RuntimeLogsScenario = 'normal' | 'hostile' | 'revisions' | 'empty' | 'deleted' | 'renamed' | 'foreign' | 'unavailable';
export const RUNTIME_LOGS_PATH = '/api/v1/system/logs';
export const RUNTIME_LOGS_HOST = '/tests/e2e/fixtures/runtime-logs-host.html';
export const RUNTIME_LOGS_TIME = domainInstant;

/** Mock wire data only; never a producer/receipt or backend correctness claim. */
export function runtimeLogsSnapshot(records: RuntimeLogRecord[] = [domainRecord('domain.delivery.retrying')], metadata: Partial<RuntimeLogMetadata> = {}): RuntimeLogSnapshot {
  if (records.length > 500) throw new Error('Snapshot fixture exceeds the existing replay record bound.');
  const instance = metadata.instance_id ?? records[0]?.instance_id ?? 'instance-a';
  return structuredClone({ instance_id: instance, oldest_cursor: records.length ? `${instance}:${records[0].sequence}` : '',
    latest_cursor: `${instance}:${records.at(-1)?.sequence ?? '0'}`, capture_dropped: 0, evicted: 0, subscriber_overflow: 0,
    file_dropped: 0, disk_errors: 0, sink_health: 'healthy', ...metadata, records });
}

/** Exactly the production parser's existing SSE envelope, including progress_cursor. */
export function runtimeLogsFrame(event: RuntimeLogEvent): string {
  if (event.type === 'log') return `event: log\nid: ${event.cursor}\ndata: ${JSON.stringify(event.record)}\n\n`;
  if ('metadata' in event) return `event: ${event.type}\ndata: ${JSON.stringify({ ...event.metadata, progress_cursor: event.cursor })}\n\n`;
  return `event: ${event.type}\ndata: ${JSON.stringify({ reason: event.reason })}\n\n`;
}

export function runtimeLogsScenarioSnapshot(scenario: RuntimeLogsScenario): RuntimeLogSnapshot {
  if (scenario === 'empty') return runtimeLogsSnapshot([]);
  if (scenario === 'deleted') return runtimeLogsSnapshot([domainRecord('domain.delivery.retrying', { fields: { group_id: 'G1' } })]);
  if (scenario === 'hostile') return runtimeLogsSnapshot([
    domainRecord('domain.delivery.blocked', { fields: { raw_error: 'password_fixture private.example SELECT secret' }, message: '<script>credential</script>' }),
    domainRecord('domain.delivery.blocked', { sequence: '9007199254740994', code: 'domain.unknown', message: 'password_fixture' }),
  ]);
  if (scenario === 'revisions') return runtimeLogsSnapshot([
    domainRecord('domain.delivery.retrying', { fields: { first_seen_at: RUNTIME_LOGS_TIME, last_seen_at: RUNTIME_LOGS_TIME, occurrences: '3' } }),
    domainRecord('domain.delivery.retrying', { sequence: '9007199254740994', fields: { entity_revision: 'R2', group_revision: 'R2', reference_id: 'REF2' } }),
  ]);
  return runtimeLogsSnapshot();
}
