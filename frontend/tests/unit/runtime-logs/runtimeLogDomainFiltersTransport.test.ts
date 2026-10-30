import { afterEach, describe, expect, it, vi } from 'vitest';
import { runtimeLogsURL, runtimeLogTransport, RuntimeLogRequestError } from '../../../src/services/runtimeLogs';
import type { EvidenceStage } from '../../../src/types/runtimeRuleEvidence';
import type { RuntimeLogFilters } from '../../../src/types/runtimeLogs';
import { filters, record, snapshot } from './fixtures';

const scoped: RuntimeLogFilters = { ...filters, source: 'domain', q: '告警 literal +[x]',
  workspaceId: 'W_1', deviceId: 'device-A', ruleId: 'rule-1', groupId: 'group-1',
  stage: 'blocked', actionable: true };
const keys = ['workspace_id', 'device_id', 'rule_id', 'group_id', 'stage', 'actionable'];
const signal = () => new AbortController().signal;
const url = (value: string) => new URL(value, 'http://localhost');
afterEach(() => vi.unstubAllGlobals());

describe('domain filters use one strict snapshot and replay/live wire contract', () => {
  it('serializes exact independent filters with legacy literal Unicode search and precise cursor', () => {
    const a = url(runtimeLogsURL(scoped));
    const b = url(runtimeLogsURL(scoped, true, 'instance-a:9007199254740993'));
    expect(Object.fromEntries(keys.map((key) => [key, a.searchParams.get(key)]))).toEqual({
      workspace_id: 'W_1', device_id: 'device-A', rule_id: 'rule-1', group_id: 'group-1',
      stage: 'blocked', actionable: 'true',
    });
    for (const key of [...keys, 'level', 'source', 'q']) {
      expect(b.searchParams.getAll(key)).toEqual(a.searchParams.getAll(key));
    }
    expect(a.searchParams.get('q')).toBe(scoped.q);
    expect(b.searchParams.get('after')).toBe('instance-a:9007199254740993');
    expect(a.searchParams.get('limit')).toBe('200');
    expect(b.searchParams.has('limit')).toBe(false);
  });
  it('omits empty UI identity fields without manufacturing a tuple or narrowing false', () => {
    const group = url(runtimeLogsURL({ ...filters, workspaceId: '', deviceId: '', ruleId: '', groupId: 'G', actionable: false }));
    expect(group.searchParams.get('group_id')).toBe('G');
    expect(group.searchParams.get('actionable')).toBe('false');
    for (const key of ['workspace_id', 'device_id', 'rule_id', 'stage']) expect(group.searchParams.has(key)).toBe(false);
    for (const key of keys) expect(url(runtimeLogsURL(filters)).searchParams.has(key)).toBe(false);
  });
  it.each<EvidenceStage>(['unavailable', 'acquisition', 'accepted', 'collecting', 'queued', 'retrying',
    'blocked', 'quarantined', 'unknown', 'skipped', 'sql-committed', 'readback-verified',
    'cleanup-completed', 'transport-event'])('uses the existing exact stage %s', (stage) => {
    expect(url(runtimeLogsURL({ ...filters, stage })).searchParams.get('stage')).toBe(stage);
  });
  it.each(['workspaceId', 'deviceId', 'ruleId', 'groupId'] as const)('rejects unsafe %s before either fetch with no value in errors', async (key) => {
    const fetch = vi.fn<typeof globalThis.fetch>(); vi.stubGlobal('fetch', fetch);
    for (const invalid of [' device-A', 'device-A ', 'gateway.private', '192.0.2.1', '[::1]',
      '/private/file', 'https://gateway.private', 'password=fixture', 'x'.repeat(129), '<img>', '\n', null, 7]) {
      const candidate: RuntimeLogFilters = { ...filters }; Object.assign(candidate, { [key]: invalid });
      for (const work of [() => runtimeLogTransport.snapshot(candidate, signal()),
        () => runtimeLogTransport.stream(candidate, '', signal(), vi.fn())]) {
        await expect(work()).rejects.toMatchObject({ status: 422, message: 'runtime_log_request_failed' });
      }
    }
    expect(fetch).not.toHaveBeenCalled();
  });
  it.each([{ stage: 'healthy' }, { stage: 'BLOCKED' }, { stage: ' blocked' }, { stage: null },
    { actionable: 'true' }, { actionable: 1 }, { actionable: null }])('rejects invalid typed filter %j without defaulting or fetching', async (invalid) => {
    const fetch = vi.fn<typeof globalThis.fetch>(); vi.stubGlobal('fetch', fetch);
    const candidate: RuntimeLogFilters = { ...filters }; Object.assign(candidate, invalid);
    expect(() => runtimeLogsURL(candidate)).toThrow(RuntimeLogRequestError);
    await expect(runtimeLogTransport.snapshot(candidate, signal())).rejects.toMatchObject({ status: 422 });
    await expect(runtimeLogTransport.stream(candidate, '', signal(), vi.fn())).rejects.toMatchObject({ status: 422 });
    expect(fetch).not.toHaveBeenCalled();
  });
  it('passes identical filters through real fetch, snapshot parser and SSE parser', async () => {
    const value = { ...record('9007199254740993'), code: 'domain.delivery.blocked', source: 'delivery',
      fields: { workspace_id: 'W_1', device_id: 'device-A', stage: 'blocked', action: 'inspect-delivery' } };
    const fetch = vi.fn<typeof globalThis.fetch>()
      .mockResolvedValueOnce(new Response(JSON.stringify({ ...snapshot, records: [value], latest_cursor: 'instance-a:9007199254740993' }),
        { headers: { 'content-type': 'application/json' } }))
      .mockResolvedValueOnce(new Response(`event: log\nid: instance-a:9007199254740993\ndata: ${JSON.stringify(value)}\n\n`,
        { headers: { 'content-type': 'text/event-stream' } }));
    vi.stubGlobal('fetch', fetch);
    const result = await runtimeLogTransport.snapshot(scoped, signal()); const events = vi.fn();
    await runtimeLogTransport.stream(scoped, result.latest_cursor, signal(), events);
    expect(result.records).toEqual([value]);
    expect(events).toHaveBeenCalledWith({ type: 'log', record: value, cursor: 'instance-a:9007199254740993' });
    expect(fetch).toHaveBeenCalledTimes(2);
    for (const call of fetch.mock.calls) {
      const request = url(String(call[0]));
      for (const key of keys) expect(request.searchParams.get(key)).toBe(url(runtimeLogsURL(scoped)).searchParams.get(key));
      expect(call[1]).toMatchObject({ credentials: 'same-origin', cache: 'no-store', redirect: 'error' });
    }
  });
});
