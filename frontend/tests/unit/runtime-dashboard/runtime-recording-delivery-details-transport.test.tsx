import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeRecordingRulesPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingRulesPanel';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import type { EvidenceScope, RuleEvidenceView, RuntimeReceipt, SourceResult } from '../../../src/types/runtimeRuleEvidence';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
const clients: QueryClient[] = [];
function receipt(record: string): RuntimeReceipt { return { group_revision: 'A', connector_revision: '7', record_id: record, effect_key: `effect-${record}`, payload_digest: `digest-${record}`, committed_at: '2026-10-27T09:59:00Z', provenance: 'runtime' }; }
function fixture(scope = evidenceScope): RuleEvidenceView {
  const view = ruleEvidenceFixture(scope); const group = view.groups[0];
  const latest: SourceResult = { state: 'partial', window_start: '2026-10-27T09:59:30Z', window_end: '2026-10-27T09:59:45Z', delta_known: null, coverage: 'partial', receipt: null };
  const committed: SourceResult = { state: 'available', window_start: '2026-10-27T09:58:30Z', window_end: '2026-10-27T09:58:45Z', delta_known: null, coverage: 'complete', receipt: receipt('source-previous') };
  group.delivery = { current: { binding: { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: false, stages: { collecting: 1, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 10, skipped: 0 }, last_commit: receipt('group-latest'), sources: [{ device_id: scope.device_id, point_id: 'P', tag_id: 'T', latest, last_committed: committed }], source_state: 'available', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null }, backlog: [], counter_replay_quarantines: [], oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, test_write: null };
  return view;
}
function current(view: RuleEvidenceView) { const value = view.groups[0].delivery?.current; if (!value) throw new Error('Expected fixture'); return value; }
function response(view = fixture()) { return { data: { success: true, data: evidenceContext(view, view.workspace_id) } }; }
function mount(scope = evidenceScope) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const content = (next: EvidenceScope) => <QueryClientProvider client={client}><RuntimeRecordingRulesPanel scope={next} selectedDeviceId={next.device_id} /></QueryClientProvider>;
  const rendered = render(content(scope)); return { ...rendered, switchTo(next: EvidenceScope) { rendered.rerender(content(next)); } };
}
function fact(region: HTMLElement, label: string) { return within(region).getByText(label, { selector: 'dt' }).nextElementSibling; }
async function selectedSource() { const summary = await screen.findByText('Source P · Tag T', { selector: 'summary' }); fireEvent.click(summary); const details = summary.closest('details'); if (!details) throw new Error('Expected details'); return details; }
beforeEach(async () => { await i18n.changeLanguage('en'); vi.resetAllMocks(); });
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('source and receipt details through actual service, parser and shared hook', () => {
  it('accepts different valid group/source receipt IDs and preserves latest partial separately from earlier committed window', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); mount(); const selected = await selectedSource();
    const latest = within(selected).getByRole('region', { name: 'Latest observed source result' }); const committed = within(selected).getByRole('region', { name: 'Last committed source result' });
    expect(fact(latest, 'Coverage')).toHaveTextContent(/^Partial$/); expect(fact(latest, 'Window start')).toHaveTextContent('09:59:30Z'); expect(fact(latest, 'Record ID')).toHaveTextContent(/^Unavailable$/);
    expect(fact(committed, 'Record ID')).toHaveTextContent(/^source-previous$/); expect(fact(committed, 'Window start')).toHaveTextContent('09:58:30Z');
    expect(screen.getByRole('region', { name: 'Source and receipt details' })).toHaveTextContent('group-latest'); expect(selected).not.toHaveTextContent('group-latest');
  });
  it('delta unknown stays unknown without changing a partial other source into complete', async () => {
    const view = fixture(); const c = current(view); c.sources[0].latest = { ...c.sources[0].latest!, state: 'unknown', delta_known: false, coverage: 'unknown' };
    view.groups[0].effective?.members.push({ device_id: 'D', point_id: 'P2', tag_id: 'T2', column: 'second', source_value_type: null, circuit_id: null });
    c.sources.push({ device_id: 'D', point_id: 'P2', tag_id: 'T2', latest: { state: 'partial', window_start: '2026-10-27T09:59:30Z', window_end: '2026-10-27T09:59:45Z', coverage: 'partial', delta_known: null, receipt: null }, last_committed: null });
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); const selected = await selectedSource();
    expect(fact(within(selected).getByRole('region', { name: 'Latest observed source result' }), 'Delta evidence')).toHaveTextContent(/^Unknown; not zero$/);
    const other = screen.getByText('Source P2 · Tag T2', { selector: 'summary' }); fireEvent.click(other); const detail = other.closest('details'); if (!detail) throw new Error('Expected details');
    const latest = within(detail).getByRole('region', { name: 'Latest observed source result' }); expect(fact(latest, 'Coverage')).toHaveTextContent(/^Partial$/); expect(fact(latest, 'Record ID')).toHaveTextContent('Unavailable');
  });
  it.each(['connector', 'revision', 'identity-collision', 'missing-group-proof', 'missing-source-receipt'] as const)('SourceReceiptMismatchNeverBorrowsGroupCommit: rejects %s', async (kind) => {
    const view = fixture(); const c = current(view); const result = c.sources[0].last_committed; if (!result?.receipt) throw new Error('Expected receipt');
    if (kind === 'connector') result.receipt.connector_revision = 'other';
    if (kind === 'revision') result.receipt.group_revision = 'B';
    if (kind === 'identity-collision') result.receipt.record_id = 'group-latest';
    if (kind === 'missing-group-proof') c.last_commit = null;
    if (kind === 'missing-source-receipt') result.receipt = null;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); await screen.findByRole('alert');
    expect(screen.queryByRole('region', { name: 'Source and receipt details' })).not.toBeInTheDocument(); expect(document.body.textContent).not.toMatch(/source-previous|group-latest/);
  });
  it.each(['foreign-device', 'foreign-point', 'foreign-tag', 'duplicate', 'raw-value', 'raw-reason', 'invalid-window', 'contradictory-coverage', 'contradictory-delta'] as const)('ForeignDuplicateTupleAndUnsafeProofFailClosed: %s', async (kind) => {
    const view = fixture(); const c = current(view); const source = c.sources[0]; const latest = source.latest; if (!latest) throw new Error('Expected latest');
    if (kind === 'foreign-device') source.device_id = 'FOREIGN'; if (kind === 'foreign-point') source.point_id = 'OTHER'; if (kind === 'foreign-tag') source.tag_id = 'OTHER';
    if (kind === 'duplicate') c.sources.push(structuredClone(source));
    if (kind === 'raw-value') Object.assign(latest, { value: 'password=secret /private/value' });
    if (kind === 'raw-reason') Object.assign(source, { reason: '/private/proof password=secret' });
    if (kind === 'invalid-window') latest.window_end = latest.window_start;
    if (kind === 'contradictory-coverage') latest.state = 'available';
    if (kind === 'contradictory-delta') latest.delta_known = false;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); const rendered = mount(); await screen.findByRole('alert');
    expect(rendered.container.innerHTML).not.toMatch(/password=secret|private|source-previous|group-latest/);
  });
  it('SameGroupOtherDeviceNeverGetsSelectedSourceSuccess: does not render another device receipt or allocate its group count', async () => {
    const view = fixture(); const c = current(view); c.shared = true;
    view.groups[0].effective?.members.push({ device_id: 'E', point_id: 'PE', tag_id: 'TE', column: 'other', source_value_type: null, circuit_id: null });
    c.sources = [{ device_id: 'E', point_id: 'PE', tag_id: 'TE', latest: null, last_committed: { ...c.sources[0].last_committed!, receipt: receipt('other-device') } }];
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); const region = await screen.findByRole('region', { name: 'Source and receipt details' });
    expect(region).toHaveTextContent('No source result reported for this device in this observation.'); expect(region).not.toHaveTextContent(/other-device|Source PE|Tag TE/);
    expect(screen.getByRole('region', { name: 'Current revision delivery' })).toHaveTextContent('10 output rows');
  });
  it('producer source-unavailable keeps independently valid A group receipt but never invents source results', async () => {
    const view = fixture(); current(view).source_state = 'unavailable'; current(view).sources = []; vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    const region = await screen.findByRole('region', { name: 'Source and receipt details' }); expect(region).toHaveTextContent('Selected-source evidence unavailable'); expect(region).toHaveTextContent('group-latest'); expect(region).not.toHaveTextContent('source-previous');
  });
  it('LegacyCustomMissingFrozenProofUnavailable: missing group proof does not show retained source or receipt as current', async () => {
    const view = fixture(); view.groups[0].proof = 'unavailable'; vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    const region = await screen.findByRole('region', { name: 'Source and receipt details' }); expect(region).toHaveTextContent('Selected-source evidence unavailable'); expect(region).not.toHaveTextContent(/source-previous|group-latest/);
  });
  it('TestWriteRuntimeReadbackCleanupIndependent: runtime receipt does not populate production-null test-write phases', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); mount(); const region = await screen.findByRole('region', { name: 'Test-write evidence' });
    for (const label of ['Test-write operation ID', 'Test write recorded at', 'Independent readback verified at', 'Owned cleanup completed at']) expect(fact(region, label)).toHaveTextContent(/^Unavailable$/);
  });
  it('NoExtraGETOnSourceOrReceiptDisclosure: shares one query for many groups and only explicit panel refresh fetches', async () => {
    const view = fixture(); view.groups = [view.groups[0], { ...structuredClone(view.groups[0]), group_id: 'G2' }]; vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(view)); mount();
    await screen.findByRole('article', { name: 'Recording group G2' });
    for (const summary of screen.getAllByText('Latest runtime receipt (group-wide)', { selector: 'summary' })) { summary.focus(); fireEvent.click(summary); expect(summary.closest('details')).toHaveAttribute('open'); }
    for (const summary of screen.getAllByText('Source P · Tag T', { selector: 'summary' })) fireEvent.click(summary);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: 'Refresh recording evidence' })); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
  });
  it('ScopeSwitchDiscardsLateDetails: old A response cannot replace B source receipts', async () => {
    let finishOld!: (value: ReturnType<typeof response>) => void; const pending = new Promise<ReturnType<typeof response>>((resolve) => { finishOld = resolve; });
    const nextScope = { ...evidenceScope, device_id: 'E' }; const next = fixture(nextScope); current(next).sources[0].last_committed!.receipt = receipt('new-device-source');
    vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(pending).mockResolvedValueOnce(response(next)); const mounted = mount(); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    const signal = vi.mocked(studioV2DatalinkApi.get).mock.calls[0][1]?.signal; mounted.switchTo(nextScope); await screen.findByText('new-device-source'); expect(signal?.aborted).toBe(true);
    await act(async () => { finishOld(response()); }); expect(screen.getByRole('region', { name: 'Source and receipt details' })).not.toHaveTextContent('source-previous');
  });
  it('mixed workspace observation fails stale and read-only refresh can recover without displaying partial details', async () => {
    const mixed = fixture(); mixed.workspace_revision = 'W2'; vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(mixed)).mockResolvedValueOnce(response()); mount();
    expect(await screen.findByRole('alert')).toHaveTextContent('Recording configuration changed'); expect(screen.queryByText('source-previous')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Refresh recording evidence' })); await screen.findByText('source-previous'); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
  });
});
