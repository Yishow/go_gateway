import { afterEach, describe, expect, it, vi } from 'vitest';
import { runtimeLogTransport, parseRuntimeLogFrame } from '@/services/runtimeLogs';
import { RuntimeLogSession } from '@/features/runtime-logs/runtimeLogSession';
import { createRuntimeLogsMock } from '../../e2e/fixtures/runtime-logs-api';
import { domainRecord, runtimeLogsSnapshot, runtimeLogsFrame, RUNTIME_LOGS_PATH } from '../../e2e/fixtures/runtime-logs-state';
import { installRuntimeLogsStream, type RuntimeLogsStreamController } from '../../e2e/fixtures/runtime-logs-stream';

const filters = { level: 'debug' as const, source: '', q: '' };
let stream: RuntimeLogsStreamController | undefined;
let session: RuntimeLogSession | undefined;
afterEach(() => { session?.stop(); session = undefined; stream?.restore(); stream = undefined; vi.unstubAllGlobals(); vi.useRealTimers(); });
function fetchFrom(mock = createRuntimeLogsMock()) {
  const fetch = vi.fn<typeof globalThis.fetch>(async (input, init) => {
    const reply = mock.handle({ method: init?.method ?? 'GET', url: String(input) });
    return new Response(JSON.stringify(reply.body), { status: reply.status, headers: { 'content-type': 'application/json' } });
  });
  vi.stubGlobal('fetch', fetch); return { mock, fetch };
}

describe('formal browser fixtures through actual runtime-log wire parser/session (not browser evidence)', () => {
  it('parses the unwrapped snapshot with exact uint64 string cursor and detached copies', async () => {
    const { mock } = fetchFrom();
    const result = await runtimeLogTransport.snapshot(filters, new AbortController().signal);
    expect(result.records[0]?.sequence).toBe('9007199254740993');
    expect(result.latest_cursor).toBe('instance-a:9007199254740993');
    result.records[0]!.fields.group_revision = 'changed';
    const again = await runtimeLogTransport.snapshot(filters, new AbortController().signal);
    expect(again.records[0]?.fields.group_revision).toBe('R1');
    expect(mock.unexpected).toEqual([]);
  });
  it('serializes actual log/progress/gap/reset envelopes rather than invented event shapes', () => {
    const record = domainRecord('domain.delivery.recovered');
    expect(parseRuntimeLogFrame(runtimeLogsFrame({ type: 'log', record, cursor: `instance-a:${record.sequence}` })))
      .toEqual({ type: 'log', record, cursor: `instance-a:${record.sequence}` });
    const metadata = runtimeLogsSnapshot([]);
    for (const type of ['handshake', 'heartbeat'] as const) {
      expect(parseRuntimeLogFrame(runtimeLogsFrame({ type, metadata, cursor: 'instance-a:9007199254740993' })))
        .toMatchObject({ type, cursor: 'instance-a:9007199254740993' });
    }
    for (const type of ['gap', 'reset'] as const) expect(parseRuntimeLogFrame(runtimeLogsFrame({ type, reason: 'retention', cursor: '' })))
      .toEqual({ type, reason: 'retention', cursor: '' });
  });
  it('matches six typed filters and preserves unfiltered progress without numeric rounding', async () => {
    const { mock } = fetchFrom();
    mock.setSnapshot(runtimeLogsSnapshot([
      domainRecord('domain.source.rejected', { fields: { device_id: 'A', rule_id: 'RULE' } }),
      domainRecord('domain.source.rejected', { sequence: '9007199254740994', fields: { device_id: 'B', rule_id: 'OTHER' } }),
    ]));
    const selected = await runtimeLogTransport.snapshot({ ...filters, workspaceId: 'W', deviceId: 'A', ruleId: 'RULE', stage: 'blocked', actionable: true }, new AbortController().signal);
    expect(selected.records).toHaveLength(1); expect(selected.records[0]?.fields.device_id).toBe('A');
    expect(selected.latest_cursor).toBe('instance-a:9007199254740994');
    const none = await runtimeLogTransport.snapshot({ ...filters, groupId: 'G' }, new AbortController().signal);
    expect(none.records).toEqual([]); expect(mock.unexpected).toEqual([]);
  });
  it.each([
    ['POST', '/api/v1/datalink/studio-v2/workspace/activate'],
    ['DELETE', '/api/v1/datalink/studio-v2/workspace/write-groups/G'],
    ['GET', '/api/v1/datalink/studio-v2/workspace?entry=locator&workspace_id=W&step=4&group_id=G&editor=1'],
    ['GET', '/api/v1/datalink/studio-v2/workspace?entry=locator&workspace_id=W&step=4'],
    ['GET', '/api/v1/datalink/studio-v2/workspace?entry=locator&workspace_id=FOREIGN&step=4&group_id=G'],
    ['GET', '/api/v1/datalink/studio-v2/workspace/write-groups/G/delivery?workspace_id=W'],
    ['GET', `${RUNTIME_LOGS_PATH}?level=debug&device_id=A&device_id=B`],
    ['GET', '/api/v1/unknown'],
  ])('strict ledger rejects %s %s instead of silently providing fallback data', (method, url) => {
    const mock = createRuntimeLogsMock(); expect(mock.handle({ method, url }).status).toBeGreaterThanOrEqual(400);
    expect(mock.unexpected).toHaveLength(1); expect(mock.ledger[0]).toMatchObject({ method, mutation: method !== 'GET' });
  });
  it('native R1 locator responses use current R2 and preserve request-time scenario data', () => {
    const mock = createRuntimeLogsMock();
    const url = '/api/v1/datalink/studio-v2/workspace?entry=locator&workspace_id=W&step=4&group_id=G&group_revision=R1';
    const original = mock.handle({ method: 'GET', url }); mock.setScenario('renamed');
    expect(JSON.stringify(original.body)).toContain('R2'); expect(JSON.stringify(original.body)).not.toContain('R3');
    expect(JSON.stringify(mock.handle({ method: 'GET', url }).body)).toContain('R3');
    expect(mock.unexpected).toEqual([]);
  });
  it('streams real encoded chunks and promptly aborts/cancels without delivering late data', async () => {
    const { fetch } = fetchFrom(); stream = installRuntimeLogsStream();
    const abort = new AbortController(); const events: string[] = [];
    const done = runtimeLogTransport.stream(filters, 'instance-a:9007199254740993', abort.signal, event => events.push(event.type)).catch(error => error);
    await vi.waitFor(() => expect(stream?.requests).toHaveLength(1));
    const frame = runtimeLogsFrame({ type: 'handshake', metadata: runtimeLogsSnapshot(), cursor: 'instance-a:9007199254740993' });
    expect(stream.send(0, frame.slice(0, 10))).toBe(true); expect(events).toEqual([]);
    stream.send(0, frame.slice(10)); await vi.waitFor(() => expect(events).toEqual(['handshake']));
    abort.abort(); await done;
    expect(stream.requests[0]).toMatchObject({ aborted: true, closed: true }); expect(stream.send(0, frame)).toBe(false);
    expect(fetch).not.toHaveBeenCalled();
  });
  it('actual session receives 403 and stops automatic requests without a proxy fallback', async () => {
    vi.useFakeTimers(); const { mock } = fetchFrom(); stream = installRuntimeLogsStream(); stream.denyNext(403);
    session = new RuntimeLogSession(filters); session.start();
    await vi.waitFor(() => expect(session?.getSnapshot().status).toBe('denied'));
    await vi.advanceTimersByTimeAsync(120000);
    expect(stream.requests).toHaveLength(1); expect(mock.ledger).toHaveLength(1); expect(mock.unexpected).toEqual([]);
  });
});
