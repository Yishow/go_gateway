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
  const view = ruleEvidenceFixture(scope); view.groups[0].delivery = { oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: 125.5, current: { binding: { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: false, stages: { collecting: 0, queued: 3, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, last_commit: null, sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null }, backlog: [], counter_replay_quarantines: [], test_write: null }; return view;
}
function delivery(view: RuleEvidenceView) { const value = view.groups[0].delivery; if (!value) throw new Error('Expected delivery'); return value; }
function response(view = fixture()) { return { data: { success: true, data: evidenceContext(view, view.workspace_id) } }; }
function mount(scope = evidenceScope) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const content = (next: EvidenceScope) => <QueryClientProvider client={client}><RuntimeRecordingRulesPanel scope={next} selectedDeviceId={next.device_id} /></QueryClientProvider>;
  const view = render(content(scope)); return { ...view, switchTo(next: EvidenceScope) { view.rerender(content(next)); } };
}
function fact(region: HTMLElement, label: string) { return within(region).getByText(label, { selector: 'dt' }).nextElementSibling; }
beforeEach(async () => { await i18n.changeLanguage('en'); vi.resetAllMocks(); });
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('observed timing through unchanged service, parser and shared query', () => {
  it('displays positive fractional all-revision age and production-null scheduling without deriving10:00:15', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); mount(); const region = await screen.findByRole('region', { name: 'Observed timing' });
    expect(fact(region, 'Oldest pending age (all group revisions)')).toHaveTextContent(/^125.5 seconds$/);
    for (const label of ['Next due at', 'Next retry at', 'Attempt duration']) expect(fact(region, label)).toHaveTextContent(/^Unavailable$/);
    expect(region).not.toHaveTextContent('10:00:15'); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it('retains observed historical age when disabled/missing applied proof makes current unavailable', async () => {
    const view = fixture(); const group = view.groups[0]; group.state = 'unavailable'; group.proof = 'unavailable'; group.effective = null; group.effective_revision = null; if (group.saved) group.saved.enabled = false; delivery(view).current = null;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); const region = await screen.findByRole('region', { name: 'Observed timing' });
    expect(fact(region, 'Oldest pending age (all group revisions)')).toHaveTextContent(/^125.5 seconds$/); expect(fact(region, 'Next due at')).toHaveTextContent(/^Unavailable$/);
    expect(region).not.toHaveTextContent(/current-only|device age|No pending/);
  });
  it.each(['null-age', 'missing-delivery'] as const)('nullable %s is not zero or a claim that no pending work exists', async (kind) => {
    const view = fixture(); if (kind === 'null-age') delivery(view).oldest_pending_seconds = null; else view.groups[0].delivery = null;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); const region = await screen.findByRole('region', { name: 'Observed timing' });
    expect(fact(region, 'Oldest pending age (all group revisions)')).toHaveTextContent(/^Unavailable$/); expect(region).not.toHaveTextContent(/0 seconds|No pending/);
  });
  it.each([-1, Number.NaN, Number.POSITIVE_INFINITY, Number.MAX_SAFE_INTEGER + 1])('InvalidAgeFailsClosedWithoutClamp: rejects %s in the existing parser', async (age) => {
    const view = fixture(); delivery(view).oldest_pending_seconds = age; vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    await screen.findByRole('alert'); expect(screen.queryByRole('region', { name: 'Observed timing' })).not.toBeInTheDocument(); expect(screen.queryByText('0 seconds')).not.toBeInTheDocument();
  });
  it.each(['wrong-age-scope', 'negative-duration', 'unsafe-due', 'extra-raw-field'] as const)('rejects malformed timing evidence %s without raw text', async (kind) => {
    const view = fixture(); const d = delivery(view); if (!d.current) throw new Error('Expected current');
    if (kind === 'wrong-age-scope') Object.assign(d, { oldest_pending_scope: 'selected-device' });
    if (kind === 'negative-duration') d.current.attempt_duration_ms = -1;
    if (kind === 'unsafe-due') d.current.next_due_at = '/private/clock password=secret';
    if (kind === 'extra-raw-field') Object.assign(d, { raw_timing: '/private/clock password=secret' });
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); const rendered = mount(); await screen.findByRole('alert'); expect(rendered.container.innerHTML).not.toMatch(/private|password=secret/);
  });
  it.each(['stale', 'foreign'] as const)('does not display timing from %s scope and only explicit read-only refresh recovers', async (kind) => {
    const view = fixture(); if (kind === 'stale') view.workspace_revision = 'W2'; else view.workspace_id = 'FOREIGN';
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)).mockResolvedValueOnce(response()); mount(); await screen.findByRole('alert'); expect(screen.queryByRole('region', { name: 'Observed timing' })).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Refresh recording evidence' })); await screen.findByRole('region', { name: 'Observed timing' }); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
  it('TimingUsesOneScopedRead: many cards and receipt disclosure do not fetch again', async () => {
    const view = fixture(); view.groups = [view.groups[0], { ...structuredClone(view.groups[0]), group_id: 'G2' }, { ...structuredClone(view.groups[0]), group_id: 'G3' }]; vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(view)); mount();
    await screen.findByRole('article', { name: 'Recording group G3' }); expect(screen.getAllByRole('region', { name: 'Observed timing' })).toHaveLength(3);
    const summary = screen.getAllByText('Latest runtime receipt (group-wide)', { selector: 'summary' })[0]; summary.focus(); expect(summary).toHaveFocus(); fireEvent.click(summary);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1); expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: evidenceScope, signal: expect.any(AbortSignal) });
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  });
  it('drops late deviceA timing instead of replacing current deviceB observation', async () => {
    let finishOld!: (value: ReturnType<typeof response>) => void; const pending = new Promise<ReturnType<typeof response>>((resolve) => { finishOld = resolve; });
    const nextScope = { ...evidenceScope, device_id: 'E' }; const next = fixture(nextScope); delivery(next).oldest_pending_seconds = 9.25;
    vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(pending).mockResolvedValueOnce(response(next)); const view = mount(); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1)); const signal = vi.mocked(studioV2DatalinkApi.get).mock.calls[0][1]?.signal;
    view.switchTo(nextScope); const region = await screen.findByRole('region', { name: 'Observed timing' }); expect(signal?.aborted).toBe(true); await act(async () => { finishOld(response()); });
    expect(fact(region, 'Oldest pending age (all group revisions)')).toHaveTextContent(/^9.25 seconds$/); expect(region).not.toHaveTextContent('125.5'); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
  });
  it('keeps stale observation visibly stale and does not label retained age live', async () => {
    const view = fixture(); view.freshness = 'stale'; vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); await screen.findByRole('region', { name: 'Observed timing' });
    expect(screen.getByText('Stale observation')).toBeInTheDocument(); expect(screen.getByRole('region', { name: 'Observed timing' })).toHaveTextContent('not a live clock or countdown'); expect(screen.getByRole('region', { name: 'Recording rules' })).not.toHaveTextContent('LIVE');
  });
  it('redacts provider errors rather than turning them into absent or zero timing', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { data: { error: { code: 'unavailable', message: '/private/clock password=secret' } } } }); const view = mount(); await screen.findByRole('alert');
    expect(view.container.innerHTML).not.toMatch(/private|password=secret/); expect(screen.queryByRole('region', { name: 'Observed timing' })).not.toBeInTheDocument();
  });
});
