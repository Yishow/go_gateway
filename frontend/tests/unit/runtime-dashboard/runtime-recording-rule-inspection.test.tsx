import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeRecordingRulesPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingRulesPanel';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import type { EvidenceScope } from '../../../src/types/runtimeRuleEvidence';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const clients: QueryClient[] = [];
function evidence(scope = evidenceScope, revision = 'A') {
  const view = ruleEvidenceFixture(scope, revision);
  view.groups[0].applied_effective_at = revision === 'A' ? null : view.observed_at;
  view.groups[0].saved!.interval_seconds = 45; view.groups[0].saved!.recording_profile!.interval = 45;
  return view;
}
const response = (view = evidence()) => ({ data: { success: true, data: evidenceContext(view, view.workspace_id) } });
function deferred() {
  let resolve!: (value: ReturnType<typeof response>) => void;
  const promise = new Promise<ReturnType<typeof response>>((done) => { resolve = done; });
  return { promise, resolve };
}
function mount(scope: EvidenceScope | null = evidenceScope, selectedDeviceId: string | null = scope?.device_id ?? null) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const content = (nextScope: EvidenceScope | null, deviceId: string | null) => <QueryClientProvider client={client}><RuntimeRecordingRulesPanel scope={nextScope} selectedDeviceId={deviceId} /></QueryClientProvider>;
  const result = render(content(scope, selectedDeviceId));
  return { ...result, switchTo(next: EvidenceScope | null, deviceId: string | null = next?.device_id ?? null) { result.rerender(content(next, deviceId)); } };
}
function fact(region: HTMLElement, label: string) { return within(region).getByText(label, { selector: 'dt' }).parentElement?.querySelector('dd')?.textContent; }
beforeEach(async () => { await i18n.changeLanguage('en'); vi.mocked(studioV2DatalinkApi.get).mockReset(); });
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); vi.restoreAllMocks(); });

describe('revision-bound recording rule inspection through the real reader and parser', () => {
  it('separates saved C, requested B metadata and effective A facts without fabricating a B boundary', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); mount();
    const card = await screen.findByRole('article', { name: 'Recording group G' });
    const saved = within(card).getByRole('region', { name: 'Saved configuration' });
    const current = within(card).getByRole('region', { name: 'Effective recording snapshot' });
    const requested = within(card).getByRole('region', { name: 'Latest requested revision' });
    expect(fact(saved, 'Fact revision')).toBe('C'); expect(fact(saved, 'Table')).toBe('meter_draft'); expect(fact(saved, 'Recording interval')).toBe('45 seconds');
    expect(fact(current, 'Fact revision')).toBe('A'); expect(fact(current, 'Table')).toBe('meter_old'); expect(fact(current, 'Recording interval')).toBe('15 seconds');
    expect(fact(requested, 'Requested revision')).toBe('B'); expect(fact(requested, 'Requested effective boundary')).toBe('Unavailable');
    expect(requested).not.toHaveTextContent(/meter_old|meter_draft|meter_new/);
    expect(requested).toHaveTextContent('a requested revision does not prove current execution.');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: evidenceScope, signal: expect.any(AbortSignal) });
  });
  it('preserves the spec example: draft B meter_new / 30 seconds while A remains meter_old / 15 seconds', async () => {
    const view = evidence(); const group = view.groups[0];
    group.saved_revision = 'B'; group.saved!.binding.entity_revision = 'B'; group.saved!.destination.table = 'meter_new';
    group.saved!.interval_seconds = 30; group.saved!.recording_profile!.interval = 30; group.applied_revision = 'A'; group.applied_effective_at = group.effective!.effective_at;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    const saved = await screen.findByRole('region', { name: 'Saved configuration' }); const current = screen.getByRole('region', { name: 'Effective recording snapshot' });
    expect(fact(saved, 'Table')).toBe('meter_new'); expect(fact(saved, 'Recording interval')).toBe('30 seconds');
    expect(fact(current, 'Table')).toBe('meter_old'); expect(fact(current, 'Recording interval')).toBe('15 seconds');
  });
  it('only changes effective A to B on a new server read, leaving saved C unchanged', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()).mockResolvedValueOnce(response(evidence(evidenceScope, 'B')));
    mount(); await screen.findByText('meter_old');
    const refresh = screen.getByRole('button', { name: 'Refresh recording evidence' }); refresh.focus(); expect(refresh).toHaveFocus(); fireEvent.click(refresh);
    await screen.findByText('meter_new'); const current = screen.getByRole('region', { name: 'Effective recording snapshot' });
    expect(fact(current, 'Fact revision')).toBe('B'); expect(fact(current, 'Effective boundary')).toBe('2026-10-27T10:00:01Z');
    expect(fact(screen.getByRole('region', { name: 'Saved configuration' }), 'Fact revision')).toBe('C'); expect(screen.getByText('meter_draft')).toBeInTheDocument();
    expect(screen.queryByText('meter_old')).not.toBeInTheDocument(); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
  });
  it('keeps saved facts but never substitutes them when immutable layout proof is missing', async () => {
    const view = evidence(); view.groups[0].proof = 'unavailable';
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    await screen.findByText('Effective recording facts unavailable'); expect(screen.getByText('meter_draft')).toBeInTheDocument();
    expect(screen.getByRole('region', { name: 'Effective recording snapshot' })).not.toHaveTextContent(/meter_draft|meter_old|45 seconds/);
  });
  it.each(['workspace', 'entity'])('rejects a mixed %s revision and offers only a read-only refresh', async (kind) => {
    const mixed = evidence(); if (kind === 'workspace') mixed.workspace_revision = 'W2'; else mixed.groups[0].effective!.binding.entity_revision = 'B';
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(mixed)).mockResolvedValueOnce(response()); mount();
    expect(await screen.findByRole('alert')).toHaveTextContent('Recording configuration changed'); expect(screen.queryByText('meter_old')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Refresh recording evidence' })); await screen.findByText('meter_old');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
    expect(vi.mocked(studioV2DatalinkApi.get).mock.calls.every(([url]) => url === '/studio-v2/workspace/runtime-context')).toBe(true);
  });
  it('discards late device A results after selecting B without retaining prior cards', async () => {
    const a = deferred(); const b = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(a.promise).mockReturnValueOnce(b.promise);
    const mounted = mount(); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    const signal = vi.mocked(studioV2DatalinkApi.get).mock.calls[0][1]?.signal; const nextScope = { ...evidenceScope, device_id: 'E' }; mounted.switchTo(nextScope);
    expect(screen.queryByRole('article')).not.toBeInTheDocument(); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2)); expect(signal?.aborted).toBe(true);
    const next = evidence(nextScope); next.groups[0].effective!.destination.table = 'other_meter';
    await act(async () => { b.resolve(response(next)); }); await screen.findByText('other_meter'); await act(async () => { a.resolve(response()); });
    expect(screen.queryByText('meter_old')).not.toBeInTheDocument(); expect(screen.getByText('other_meter')).toBeInTheDocument();
  });
  it.each([
    null, { ...evidenceScope, workspace_id: '' }, { ...evidenceScope, workspace_revision: '' }, { ...evidenceScope, device_id: '' },
    { ...evidenceScope, workspace_id: '/private/secret' }, { workspace_revision: 'W1', device_id: 'D' } as EvidenceScope,
  ])('does not mount a query for incomplete or unsafe scope %j', (scope) => {
    mount(scope, 'D'); expect(screen.getByText(/Recording evidence unavailable until/)).toBeInTheDocument(); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
    expect(document.body.textContent).not.toContain('/private/secret');
  });
  it.each([null, 'foreign'])('does not use a valid but differently selected device scope %s', (selected) => {
    mount(evidenceScope, selected); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled(); expect(screen.getByText(/Recording evidence unavailable until/)).toBeInTheDocument();
  });
  it('clears an existing card immediately when the selected scope is removed', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); const mounted = mount(); await screen.findByText('meter_old');
    mounted.switchTo(null); expect(screen.queryByText('meter_old')).not.toBeInTheDocument(); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it('performs one group read for many cards and none for shared-member disclosure', async () => {
    const view = evidence(); view.groups[0].effective!.members.push({ device_id: 'E', point_id: 'PE', tag_id: 'TE', column: 'other_value', source_value_type: 'float64', circuit_id: null });
    view.groups = [view.groups[0], { ...structuredClone(view.groups[0]), group_id: 'G2' }, { ...structuredClone(view.groups[0]), group_id: 'G3' }];
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); await screen.findByRole('article', { name: 'Recording group G3' });
    const summary = screen.getAllByText('Members on other devices (1)', { selector: 'summary' })[0]; summary.focus(); fireEvent.click(summary);
    expect(summary.parentElement).toHaveAttribute('open'); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    expect(screen.getAllByRole('article')).toHaveLength(3);
  });
  it('shows stopped and retained ended configuration without intake or SQL-success claims', async () => {
    const view = evidence(); view.groups[0].effective!.enabled = false; view.groups[0].effective!.ends_at = view.observed_at;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount(); await screen.findByText('Disabled in this configuration');
    expect(screen.getAllByText(/not proof of current intake or SQL delivery/).length).toBeGreaterThan(0);
    expect(fact(screen.getByRole('region', { name: 'Effective recording snapshot' }), 'Snapshot ends at')).toBe('2026-10-27T10:00:00Z');
    expect(screen.queryByRole('button', { name: /start|activate|retry delivery/i })).not.toBeInTheDocument();
  });
  it.each(['stale', 'unavailable'] as const)('presents server freshness %s without manufacturing LIVE status', async (freshness) => {
    const view = evidence(); view.freshness = freshness; if (freshness === 'unavailable') { view.observed_at = null; view.groups = []; }
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    if (freshness === 'stale') { await screen.findByText('Stale observation'); expect(screen.getByText('meter_old')).toBeInTheDocument(); }
    else expect(await screen.findByRole('alert')).toHaveTextContent('Recording evidence is unavailable');
    const panel = screen.getByRole('region', { name: 'Recording rules' }); expect(panel).not.toHaveTextContent('Fresh observation'); expect(panel).not.toHaveTextContent('LIVE');
    expect(fact(panel, 'Observed at')).toBe(freshness === 'stale' ? '2026-10-27T10:00:00Z' : 'Unavailable');
  });
  it.each([
    ['unavailable', null], ['unavailable', '2026-10-27T10:00:00Z'],
    ['stale', '2026-10-27T10:00:00Z'], ['fresh', null],
  ] as const)('does not claim no groups from %s evidence observed at %s', async (freshness, observedAt) => {
    const view = evidence(); view.groups = []; view.freshness = freshness; view.observed_at = observedAt;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); mount();
    expect(await screen.findByRole('alert')).toHaveTextContent('Recording evidence is unavailable');
    expect(screen.queryByText('No recording groups for this device.')).not.toBeInTheDocument();
    expect(screen.queryByRole('article')).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Refresh recording evidence' })).toBeEnabled();
  });
  it('renders loading and fresh observed empty states without fabricated configuration', async () => {
    const pending = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(pending.promise); mount();
    expect(screen.getByText('Reading recording evidence…')).toBeInTheDocument(); expect(screen.getByRole('button', { name: 'Refresh recording evidence' })).toBeDisabled();
    const empty = evidence(); empty.groups = []; await act(async () => { pending.resolve(response(empty)); });
    await screen.findByText('No recording groups for this device.'); expect(screen.queryByRole('article')).not.toBeInTheDocument();
  });
  it.each(['host', 'database', 'raw_error'])('rejects unsafe %s through the actual service/parser without exposing it', async (field) => {
    const view = evidence(); Object.assign(view.groups[0].effective!.destination, { [field]: '/private/secret' });
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(view)); const mounted = mount();
    expect(await screen.findByRole('alert')).toHaveTextContent('Recording evidence is unavailable'); expect(mounted.container.innerHTML).not.toMatch(/private|secret|meter_old/);
  });
  it.each(['stale-reference', 'foreign-scope', 'driver: secret'])('uses safe typed HTTP error %s without raw exception copy', async (code) => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { data: { error: { code, message: 'password=secret /private/driver.log' } } } }); const mounted = mount();
    await screen.findByRole('alert'); expect(mounted.container.innerHTML).not.toMatch(/password|secret|private|driver/); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it('supports Traditional Chinese panel labels and native focusable refresh', async () => {
    await i18n.changeLanguage('zh-TW'); vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); mount();
    await screen.findByText('meter_old'); const refresh = screen.getByRole('button', { name: i18n.t('runtime-dashboard:recordingRules.refresh') }); refresh.focus(); expect(refresh).toHaveFocus();
    expect(screen.getByRole('region', { name: i18n.t('runtime-dashboard:recordingRules.title') })).toBeInTheDocument();
  });
});
