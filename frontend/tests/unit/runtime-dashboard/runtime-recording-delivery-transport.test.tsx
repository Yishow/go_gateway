import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeRecordingRulesPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingRulesPanel';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import type { EvidenceScope, RuleEvidenceView } from '../../../src/types/runtimeRuleEvidence';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
const clients: QueryClient[] = [];
function fixture(scope = evidenceScope): RuleEvidenceView {
  const view = ruleEvidenceFixture(scope); const group = view.groups[0];
  group.delivery = { current: { binding: { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: false,
    stages: { collecting: 4, queued: 3, retrying: 0, blocked: 1, quarantined: 2, unknown: 0, sql_committed: 0, skipped: 0 }, last_commit: null, sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null },
    backlog: [{ binding: { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, entity_revision: 'A0' }, destination: { connector_id: 'C1', connector_revision: '7', alias: null, kind: 'unavailable', schema: null, table: 'old_table' }, pending: 3, reasons: ['recording-ownership-conflict'] }],
    counter_replay_quarantines: [{ workspace_id: scope.workspace_id, group_id: 'G', group_revision: 'A0', pending_samples: 4, reason: 'counter-replay-proof-conflict', action: 'preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance' }], oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, test_write: null };
  return view;
}
function delivery(view: RuleEvidenceView) { const value = view.groups[0].delivery; if (!value?.current) throw new Error('Expected current fixture'); return { delivery: value, current: value.current }; }
function response(view = fixture()) { return { data: { success: true, data: evidenceContext(view, view.workspace_id) } }; }
function mount(scope = evidenceScope) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const content = (next: EvidenceScope) => <QueryClientProvider client={client}><RuntimeRecordingRulesPanel scope={next} selectedDeviceId={next.device_id} /></QueryClientProvider>;
  const view = render(content(scope)); return { ...view, switchTo(next: EvidenceScope) { view.rerender(content(next)); } };
}
function fact(region: HTMLElement, label: string) { return within(region).getByText(label, { selector: 'dt' }).nextElementSibling; }
beforeEach(async () => { await i18n.changeLanguage('en'); vi.resetAllMocks(); });
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('delivery projection through unchanged actual service, parser and hook', () => {
  it('uses effective A not requested B and performs one GET for many cards, none on disclosure', async () => {
    const view = fixture(); view.groups = [view.groups[0], ...['G2', 'G3'].map((id) => { const group = structuredClone(view.groups[0]); group.group_id = id; if (group.delivery?.counter_replay_quarantines) group.delivery.counter_replay_quarantines[0].group_id = id; return group; })];
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(view)); mount();
    await screen.findByRole('article', { name: 'Recording group G3' });
    const current = screen.getAllByRole('region', { name: 'Current revision delivery' }); expect(current).toHaveLength(3);
    for (const region of current) expect(fact(region, 'Delivery revision')).toHaveTextContent(/^A$/);
    const summary = screen.getAllByText('Backlog revision A0', { selector: 'summary' })[0]; summary.focus(); fireEvent.click(summary); expect(summary.closest('details')).toHaveAttribute('open');
    fireEvent.click(screen.getAllByText('Replay revision A0', { selector: 'summary' })[0]);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: evidenceScope, signal: expect.any(AbortSignal) });
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  });
  it('allows a matching runtime commit count but keeps receipts out of the stage card and shared totals group-wide', async () => {
    const view = fixture(); const c = delivery(view).current; c.shared = true; c.stages.sql_committed = 10;
    view.groups[0].effective?.members.push({ device_id: 'E', point_id: 'PE', tag_id: 'TE', column: 'other', source_value_type: null, circuit_id: null });
    c.last_commit = { group_revision: 'A', connector_revision: '7', record_id: 'RUNTIME_RECORD', effect_key: 'EFFECT', payload_digest: 'DIGEST', committed_at: '2026-10-27T10:00:00Z', provenance: 'runtime' };
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    const region = await screen.findByRole('region', { name: 'Current revision delivery' });
    expect(fact(region, 'SQL committed')).toHaveTextContent(/^10 output rows$/); expect(region).toHaveTextContent('Shared group totals across devices');
    expect(screen.getByRole('region', { name: 'Delivery evidence' }).textContent).not.toMatch(/RUNTIME_RECORD|EFFECT|DIGEST|readback verified|cleanup completed/);
  });
  it.each(['missing', 'wrong-revision', 'wrong-connector', 'proof-with-zero'] as const)('MissingOrMismatchedCommitProofFailsClosed: %s', async (kind) => {
    const view = fixture(); const c = delivery(view).current; c.stages.sql_committed = kind === 'proof-with-zero' ? 0 : 10;
    if (kind !== 'missing') c.last_commit = { group_revision: kind === 'wrong-revision' ? 'B' : 'A', connector_revision: kind === 'wrong-connector' ? '8' : '7', record_id: 'RUNTIME_RECORD', effect_key: 'EFFECT', payload_digest: 'DIGEST', committed_at: '2026-10-27T10:00:00Z', provenance: 'runtime' };
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    await screen.findByRole('alert'); expect(screen.queryByRole('region', { name: 'Current revision delivery' })).not.toBeInTheDocument();
    expect(document.body.textContent).not.toMatch(/10 output rows|RUNTIME_RECORD|EFFECT|DIGEST/);
  });
  it.each(['workspace', 'revision', 'group', 'intake', 'ended-intake', 'shared'] as const)('rejects mismatched %s evidence instead of showing plausible counts', async (kind) => {
    const view = fixture(); const d = delivery(view);
    if (kind === 'workspace') d.current.binding.workspace_id = 'FOREIGN';
    if (kind === 'revision') d.current.binding.entity_revision = 'B';
    if (kind === 'group') d.delivery.counter_replay_quarantines![0].group_id = 'OTHER';
    if (kind === 'intake') d.current.intake = { group_revision: 'B', state: 'active', reason: null };
    if (kind === 'ended-intake') { if (!view.groups[0].effective) throw new Error('Expected facts'); view.groups[0].effective.ends_at = view.observed_at; d.current.intake = { group_revision: 'A', state: 'active', reason: null }; }
    if (kind === 'shared') d.current.shared = true;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); await screen.findByRole('alert');
    expect(screen.queryByText('4 accepted samples')).not.toBeInTheDocument();
  });
  it('missing effective proof drops current counts but preserves independently bound retained backlog and replay', async () => {
    const view = fixture(); view.groups[0].proof = 'unavailable'; vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    await screen.findByText('Current revision delivery unavailable');
    expect(screen.getByRole('region', { name: 'Retained revision backlog' })).toHaveTextContent('old_table');
    expect(screen.getByRole('region', { name: 'Accepted-input replay quarantine' })).toHaveTextContent('4 accepted samples');
    expect(screen.getByRole('region', { name: 'Current revision delivery' })).not.toHaveTextContent('0 output rows');
  });
  it('retains disabled ended configuration and backlog without inventing intake state', async () => {
    const view = fixture(); const facts = view.groups[0].effective; if (!facts) throw new Error('Expected facts'); facts.enabled = false; facts.ends_at = view.observed_at;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); const current = await screen.findByRole('region', { name: 'Current revision delivery' });
    expect(fact(current, 'Observed intake state')).toHaveTextContent(/^Unavailable$/); expect(current).not.toHaveTextContent('Not running');
    expect(screen.getByRole('region', { name: 'Retained revision backlog' })).toHaveTextContent('3 output rows');
  });
  it.each(['reason', 'raw-field', 'unsafe-destination', 'negative-count', 'negative-replay', 'unknown-replay-reason', 'unsafe-replay-action'] as const)('rejects unsafe or invalid %s without raw diagnostics in DOM', async (kind) => {
    const view = fixture(); const d = delivery(view);
    if (kind === 'reason') Object.assign(d.delivery.backlog[0], { reasons: ['password=secret'] });
    if (kind === 'raw-field') Object.assign(d.current, { raw_error: '/private/driver.log password=secret' });
    if (kind === 'unsafe-destination') Object.assign(d.delivery.backlog[0].destination, { host: 'secret-host', path: '/private/db' });
    if (kind === 'negative-count') d.current.stages.collecting = -1;
    if (kind === 'negative-replay') d.delivery.counter_replay_quarantines![0].pending_samples = -1;
    if (kind === 'unknown-replay-reason') Object.assign(d.delivery.counter_replay_quarantines![0], { reason: 'password=secret' });
    if (kind === 'unsafe-replay-action') Object.assign(d.delivery.counter_replay_quarantines![0], { action: '/private/driver' });
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); const mounted = mount(); await screen.findByRole('alert');
    expect(mounted.container.innerHTML).not.toMatch(/password=secret|secret-host|private|driver/); expect(screen.queryByRole('region', { name: 'Delivery evidence' })).not.toBeInTheDocument();
  });
  it('typed unavailable error is not an empty-success response and refresh remains read-only', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { data: { error: { code: 'unavailable', message: 'password=secret' } } } }).mockResolvedValueOnce(response()); mount();
    expect(await screen.findByRole('alert')).toHaveTextContent('Recording evidence is unavailable'); expect(screen.queryByText(/No retained backlog/)).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Refresh recording evidence' })); await screen.findByRole('region', { name: 'Current revision delivery' });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
  it('ignores delayed old-device delivery counts after a scope switch', async () => {
    let finishOld!: (value: ReturnType<typeof response>) => void;
    const pending = new Promise<ReturnType<typeof response>>((resolve) => { finishOld = resolve; });
    const nextScope = { ...evidenceScope, device_id: 'E' }; const next = fixture(nextScope); delivery(next).current.stages.collecting = 9;
    vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(pending).mockResolvedValueOnce(response(next)); const mounted = mount();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1)); const signal = vi.mocked(studioV2DatalinkApi.get).mock.calls[0][1]?.signal;
    mounted.switchTo(nextScope); await screen.findByRole('region', { name: 'Current revision delivery' }); expect(signal?.aborted).toBe(true);
    await act(async () => { finishOld(response()); });
    expect(fact(screen.getByRole('region', { name: 'Current revision delivery' }), 'Collecting')).toHaveTextContent(/^9 accepted samples$/);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
  });
  it('preserves stale observation labeling rather than presenting retained counts as live', async () => {
    const view = fixture(); view.freshness = 'stale'; vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    await screen.findByRole('region', { name: 'Current revision delivery' }); expect(screen.getByText('Stale observation')).toBeInTheDocument();
    expect(screen.getByRole('region', { name: 'Recording rules' })).not.toHaveTextContent('LIVE');
  });
});
