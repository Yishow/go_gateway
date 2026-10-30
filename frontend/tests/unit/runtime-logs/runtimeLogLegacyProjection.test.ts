import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { describe, expect, it } from 'vitest';
import type { RuntimeLogRecord } from '@/types/runtimeLogs';
import { projectRuntimeLogLegacy } from '@/services/runtimeLogLegacyProjection';
import { legacyCases, legacyInstant, legacyRecord, legacyRequestID } from './legacy-event-fixtures';

function suppressed(value: RuntimeLogRecord) { expect(projectRuntimeLogLegacy(value).kind).toBe('suppressed'); }

describe('closed legacy runtime log presentation', () => {
  it('matches all 33 actual Go legacy template codes, severity and source', () => {
    const go = readFileSync(resolve(process.cwd(), '../internal/diagnostics/projection.go'), 'utf8').split('func project')[0];
    const actual = [...go.matchAll(/"((?:startup|runtime|shutdown|http)\.[^"]+)":\s*\{level(\w+),\s*source(\w+),/g)]
      .map(([, code, level, source]) => [code, level.toLowerCase(), source.toLowerCase()]);
    actual.push(['raw.suppressed', 'info', 'application']);
    expect(actual).toEqual(legacyCases); expect(actual).toHaveLength(33);
  });
  it.each(legacyCases)('projects known %s without trusting candidate message', (code, level, source) => {
    const record = legacyRecord(code, { message: 'password_fixture <script>private.example</script>' });
    const view = projectRuntimeLogLegacy(record);
    expect(view).toMatchObject({ kind: 'legacy', code, level, source, timestamp: legacyInstant, cursor: 'instance-a:9007199254740993', truncated: false,
      category: code === 'raw.suppressed' ? 'suppressed' : code.startsWith('http.') ? 'http' : 'process' });
    expect(JSON.stringify(view)).not.toMatch(/password_fixture|script|private\.example|Candidate/);
    expect(view).not.toHaveProperty('message'); expect(view).not.toHaveProperty('fields');
  });
  it.each(['1', '9007199254740993', '18446744073709551615'])('keeps sequence %s exact and never numeric', sequence => {
    expect(projectRuntimeLogLegacy(legacyRecord('runtime.started', { sequence }))).toMatchObject({ cursor: `instance-a:${sequence}` });
  });
  it.each(['0', '00', '01', '-1', '1.0', '1e3', ' 1', '18446744073709551616', '9'.repeat(200)])('rejects invalid canonical sequence %s', sequence => {
    const view = projectRuntimeLogLegacy(legacyRecord('runtime.started', { sequence })); expect(view).toMatchObject({ kind: 'suppressed', cursor: null });
  });
  it.each(['', 'unsafe.example', '127.0.0.1', '[::1]', '../path', 'A:B', 'A B', 'A'.repeat(129)])('rejects unsafe transport instance %s', instance_id => {
    expect(projectRuntimeLogLegacy(legacyRecord('runtime.started', { instance_id }))).toMatchObject({ kind: 'suppressed', cursor: null });
  });
  it.each(['A_1-x', '0123456789abcdef0123456789abcdef', 'A'.repeat(128)])('preserves allowed opaque transport instance %s', instance_id => {
    expect(projectRuntimeLogLegacy(legacyRecord('runtime.started', { instance_id }))).toMatchObject({ cursor: `${instance_id}:9007199254740993` });
  });
  it.each(['2026-02-30T10:00:00Z', '2026-10-29T24:00:00Z', '2026-10-29T10:00:60Z', '2026-10-29T10:00:00+00:00',
    '2026-10-29T10:00:00Z password_fixture', '2026-10-29T10:00:00.1234567890Z'])('rejects noncanonical UTC timestamp %s', timestamp => {
    expect(projectRuntimeLogLegacy(legacyRecord('runtime.started', { timestamp }))).toMatchObject({ kind: 'suppressed', timestamp: null });
  });
  it.each(['2024-02-29T10:00:00Z', '0001-01-01T00:00:00.000000001Z', '9999-12-31T23:59:59.999999999Z', legacyInstant])('preserves exact valid UTC bytes %s', timestamp => {
    expect(projectRuntimeLogLegacy(legacyRecord('runtime.started', { timestamp }))).toMatchObject({ kind: 'legacy', timestamp });
  });
  it.each(['GET', 'HEAD', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'CONNECT', 'TRACE'])('accepts producer method %s', method => {
    expect(projectRuntimeLogLegacy(legacyRecord('http.access', { fields: { method } }))).toMatchObject({ method });
  });
  it.each<RuntimeLogRecord['fields']>([{ method: 'get' }, { method: 'FETCH' }, { method: 'GET\nsecret' }, { status: 99 }, { status: 600 }, { status: 200.1 },
    { status: '200' }, { status: NaN }, { duration_ms: -1 }, { duration_ms: Infinity }, { duration_ms: 86400001 }, { duration_ms: '0' },
    { request_id: 'not-an-opaque-id' }, { request_id: 'g'.repeat(32) }, { request_id: '0123456789abcdef0123456789abcdef00---0' }])('rejects unsupported HTTP scalar %j', fields => {
    suppressed(legacyRecord('http.access', { fields }));
  });
  it.each([100, 599])('keeps boundary status %i', status => {
    expect(projectRuntimeLogLegacy(legacyRecord('http.request', { fields: { status } }))).toMatchObject({ status });
  });
  it.each([0, 0.001, 86400000])('retains finite producer duration %s without fabricated zero', duration_ms => {
    expect(projectRuntimeLogLegacy(legacyRecord('http.access', { fields: { duration_ms } }))).toMatchObject({ durationMs: duration_ms });
  });
  it.each(['0123456789abcdef0123456789abcdef', '0123456789ABCDEF0123456789ABCDEF', legacyRequestID])('keeps exact opaque request reference %s', request_id => {
    expect(projectRuntimeLogLegacy(legacyRecord('http.access', { fields: { request_id } }))).toMatchObject({ requestID: request_id });
  });
  it('keeps absent HTTP scalar fields null independently of observed zero duration', () => {
    const record = legacyRecord('http.access'); record.fields = {};
    expect(projectRuntimeLogLegacy(record)).toMatchObject({ method: null, status: null, durationMs: null, requestID: null });
    record.fields.duration_ms = 0;
    expect(projectRuntimeLogLegacy(record)).toMatchObject({ method: null, status: null, durationMs: 0, requestID: null });
  });
  it.each(['/api/v1/devices/:id', '/private/password_fixture', 'https://private.example/db?token=secret', '界'.repeat(341)])('never retains route scalar %s', route => {
    const view = projectRuntimeLogLegacy(legacyRecord('http.access', { fields: { route } }));
    expect(view).toMatchObject({ kind: 'legacy', method: 'GET', status: 200 }); expect(view).not.toHaveProperty('route');
    expect(JSON.stringify(view)).not.toContain(route);
  });
  it.each<RuntimeLogRecord['fields']>([{ route: '界'.repeat(342) }, { route: 'x'.repeat(1025) }, { route: 7 }])('rejects out-of-contract route structure %j', fields => {
    suppressed(legacyRecord('http.access', { fields }));
  });
  it.each(['startup', 'http', 'runtime', 'pipeline', 'share', 'connections', 'database', 'diagnostics', 'tray', 'owner'])('allows exact shutdown phase %s', phase => {
    expect(projectRuntimeLogLegacy(legacyRecord('shutdown.finalizing', { fields: { phase } }))).toMatchObject({ phase });
  });
  it('does not invent shutdown phase and rejects unknown phase', () => {
    const record = legacyRecord('shutdown.failed'); delete record.fields.phase;
    expect(projectRuntimeLogLegacy(record)).toMatchObject({ phase: null });
    suppressed(legacyRecord('shutdown.failed', { fields: { phase: 'sql-committed' } }));
  });
  it.each(['application', 'standard', 'slog', 'http', 'runtime'])('permits only producer suppression source %s', source => {
    expect(projectRuntimeLogLegacy(legacyRecord('raw.suppressed', { source }))).toMatchObject({ kind: 'legacy', category: 'suppressed', source });
  });
  it.each(['startup', 'shutdown', 'private.example', '<img src=x>'])('does not invent suppressed source %s', source => suppressed(legacyRecord('raw.suppressed', { source })));
  it('rejects wrong fixed source and severity rather than borrowing envelope meaning', () => {
    suppressed(legacyRecord('runtime.started', { source: 'startup' })); suppressed(legacyRecord('startup.ready', { level: 'warn' }));
    suppressed(legacyRecord('http.recovered', { level: 'info' })); suppressed(legacyRecord('raw.suppressed', { level: 'error' }));
  });
  it.each(['startup.ready', 'runtime.started', 'raw.suppressed'] as const)('does not attribute arbitrary fields to %s', code => {
    for (const fields of [{ phase: 'database' }, { method: 'GET' }, { workspace_id: 'W' }, { group_id: 'G' }, { value: 4 }] as RuntimeLogRecord['fields'][])
      suppressed(legacyRecord(code, { fields }));
  });
  it('rejects unknown branch fields instead of serializing candidate JSON', () => {
    suppressed(legacyRecord('http.access', { fields: { raw_error: 'secret', connector: 'private.example' } }));
    suppressed(legacyRecord('shutdown.complete', { fields: { request_id: legacyRequestID } }));
  });
  it('does not execute getters or accept symbols, prototypes, nested payloads or unsafe outer cells', () => {
    let reads = 0; const getter = () => { reads++; throw new Error('password_fixture'); };
    for (const key of ['code', 'source', 'timestamp', 'fields']) {
      const record = legacyRecord(); Object.defineProperty(record, key, { get: getter }); suppressed(record);
    }
    const inner = legacyRecord('http.access'); Object.defineProperty(inner.fields, 'route', { get: getter }); suppressed(inner);
    expect(reads).toBe(0);
    const symbol = legacyRecord(); Object.defineProperty(symbol.fields, Symbol('secret'), { value: 'secret' }); suppressed(symbol);
    const prototype = legacyRecord(); Object.setPrototypeOf(prototype.fields, { private: 'secret' }); suppressed(prototype);
    const inherited = legacyRecord(); Object.setPrototypeOf(inherited, { private: 'secret' }); suppressed(inherited);
    const nested = legacyRecord('http.access'); Object.defineProperty(nested.fields, 'request_id', { value: { password: 'secret' } }); suppressed(nested);
    const outer = legacyRecord(); Object.defineProperty(outer, 'private', { value: 'secret' }); suppressed(outer);
    const wrong = legacyRecord(); Object.defineProperty(wrong, 'truncated', { value: 'false' }); suppressed(wrong);
  });
  it('detaches retained values and never retains route/message candidates', () => {
    const record = legacyRecord('http.access', { truncated: true, fields: { route: '/private/password_fixture' } });
    const view = projectRuntimeLogLegacy(record); record.fields.method = 'DELETE'; record.fields.request_id = 'other'; record.sequence = '9';
    expect(view).toMatchObject({ kind: 'legacy', method: 'GET', requestID: legacyRequestID, cursor: 'instance-a:9007199254740993', truncated: true });
    expect(JSON.stringify(view)).not.toMatch(/password_fixture|Candidate/);
  });
  it.each(['unknown.code', 'password_fixture <script>', 'domain.unavailable', 'domain.delivery.first-confirmed'])('unknown or nonlegacy code %s has only fixed fallback metadata', code => {
    const view = projectRuntimeLogLegacy(legacyRecord('runtime.started', { code, message: 'password_fixture' }));
    expect(view).toEqual({ kind: 'suppressed', timestamp: legacyInstant, cursor: 'instance-a:9007199254740993', truncated: false });
  });
  it('does not turn transport/process codes into acquisition, receipt or action authority', () => {
    for (const code of ['runtime.started', 'startup.ready', 'http.access', 'http.recovered'] as const) {
      const view = projectRuntimeLogLegacy(legacyRecord(code)); expect(view.kind).toBe('legacy');
      for (const property of ['stage', 'action', 'receipt', 'locator', 'group_id', 'device_id', 'healthy']) expect(view).not.toHaveProperty(property);
    }
  });
});
