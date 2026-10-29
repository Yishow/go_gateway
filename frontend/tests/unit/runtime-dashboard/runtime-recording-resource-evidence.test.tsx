import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeRecordingRulesPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingRulesPanel';
import { RuntimeRecordingDeliveryCard } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingDeliveryCard';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { parseRuntimeRuleEvidence } from '../../../src/services/runtimeRuleEvidenceParser';
import type { QuotaEvidence, RuleEvidenceView, StorageEvidence } from '../../../src/types/runtimeRuleEvidence';
import { createRuntimeRecordingMock } from '../../e2e/fixtures/runtime-recording-api';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
const clients: QueryClient[] = [];
const title = 'Quota and storage evidence';
const quota = (overrides: Partial<QuotaEvidence> = {}): QuotaEvidence => ({ scope: 'group', state: 'warning', used_bytes: 80, max_bytes: 100, intake_refused: false, ...overrides });
function fixture(q: QuotaEvidence | null = quota(), storage: StorageEvidence | null = null): RuleEvidenceView {
  const value = ruleEvidenceFixture(); const group = value.groups[0];
  for (const facts of [group.saved, group.effective]) if (facts) facts.members.push({ device_id: 'E', point_id: 'PE', tag_id: 'TE', column: 'other', source_value_type: 'float64', circuit_id: null });
  group.delivery = { current: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: true,
    stages: { collecting: 2, queued: 3, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, last_commit: null, sources: [], source_state: 'unavailable',
    next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: q, storage }, backlog: [], counter_replay_quarantines: [], oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, test_write: null };
  return value;
}
function mount(value: unknown = fixture()) {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: evidenceContext(value) } });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  return render(<QueryClientProvider client={client}><RuntimeRecordingRulesPanel scope={evidenceScope} selectedDeviceId="D" /></QueryClientProvider>);
}
function pure(value: RuleEvidenceView) {
  const parsed = parseRuntimeRuleEvidence(value, evidenceScope);
  expect(parsed.ok).toBe(true); if (!parsed.ok) throw new Error('Expected existing typed resource contract');
  return render(<RuntimeRecordingDeliveryCard group={parsed.value.groups[0]} />);
}
function fact(region: HTMLElement, label: string) { return within(region).getByText(label, { selector: 'dt' }).nextElementSibling; }
beforeEach(async () => { await i18n.changeLanguage('en'); vi.resetAllMocks(); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

describe('existing quota and frozen-storage evidence presentation', () => {
  it.each(['quota-group', 'quota-global'] as const)('formal fixture supplies its exact quota observation: %s', scenario => {
    const mock = createRuntimeRecordingMock(scenario);
    const reply = mock.handle({ method: 'GET', url: '/api/v1/datalink/studio-v2/workspace/runtime-context?workspace_id=W&workspace_revision=W1&device_id=A' });
    expect(reply.status).toBe(200);
    expect(reply.body).toMatchObject({ data: { rule_evidence: { groups: [{ delivery: { current: { shared: true, storage: null,
      quota: scenario === 'quota-group' ? { scope: 'group', state: 'warning', used_bytes: 80, max_bytes: 100, intake_refused: false }
        : { scope: 'global', state: 'hard_limit', used_bytes: 125, max_bytes: 100, intake_refused: true } } } }] } } });
  });

  it.each([
    quota(), quota({ scope: 'global', state: 'hard_limit', used_bytes: 125, intake_refused: true }),
    quota({ state: 'ok', used_bytes: 0 }), quota({ state: 'ok', used_bytes: 125 }),
    quota({ state: 'unconfigured', used_bytes: 0, max_bytes: 0 }), quota({ state: 'unavailable' }),
  ])('renders exact observed quota without thresholds, ratios or selected-device allocation: %j', async observed => {
    mount(fixture(observed)); const region = await screen.findByRole('region', { name: title });
    expect(fact(region, 'Quota scope')).toHaveTextContent(observed.scope === 'global' ? 'Global — all groups and revisions' : 'This group — all revisions');
    expect(fact(region, 'Used durable bytes')).toHaveTextContent(new RegExp(`^${observed.used_bytes}$`));
    expect(fact(region, 'Configured maximum bytes')).toHaveTextContent(new RegExp(`^${observed.max_bytes}$`));
    expect(fact(region, 'New intake refused by quota')).toHaveTextContent(observed.intake_refused ? /^Yes$/ : /^No$/);
    const states = { warning: 'Warning', hard_limit: 'Hard limit', ok: 'OK', unconfigured: 'Unconfigured', unavailable: 'Unavailable' };
    expect(fact(region, 'Quota state')).toHaveTextContent(new RegExp(`^${states[observed.state]}$`));
    expect(region).toHaveTextContent('Quota bytes are not selected-device or Tag totals and do not prove SQL delivery or intake permission.');
    expect(region).not.toHaveTextContent(/\d+%|next retry|SQL committed/); expect(region.querySelector('[style*="width"]')).not.toBeInTheDocument();
    expect(within(region).queryByRole('button')).not.toBeInTheDocument(); expect(within(region).queryByRole('link')).not.toBeInTheDocument();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  });
  it.each(['quota-null', 'current-null', 'delivery-null'] as const)('keeps missing/unconfigured quota and frozen storage explicitly unavailable for %s', kind => {
    const value = fixture(null); const group = value.groups[0];
    if (kind === 'delivery-null') group.delivery = null;
    else if (kind === 'current-null' && group.delivery) group.delivery.current = null;
    pure(value); const region = screen.getByRole('region', { name: title });
    for (const label of ['Quota scope', 'Quota state', 'Used durable bytes', 'Configured maximum bytes', 'New intake refused by quota']) expect(fact(region, label)).toHaveTextContent(/^Unavailable$/);
    const storage = within(region).getByRole('region', { name: 'Frozen-destination storage history' });
    expect(storage).toHaveTextContent('Frozen-destination storage history and events unavailable.');
    expect(storage).not.toHaveTextContent(/No storage events|^0$|healthy/);
    expect(within(storage).queryByRole('list')).not.toBeInTheDocument(); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it('renders only existing validated nonnull storage facts/events, without repair actions or inferred history restoration', () => {
    pure(fixture(null, { generation: 2, history_may_be_lost: true, reconstructions: 3, last_reconstructed_at: '2026-10-29T09:59:00.123456789Z', repair_state: 'blocked',
      events: [{ kind: 'confirmed', at: '2026-10-29T09:58:00Z' }, { kind: 'recreated_by_confirmation', at: null }, { kind: 'reconstructed', at: '2026-10-29T09:59:00Z' }, { kind: 'repair_blocked', at: null }, { kind: 'unavailable', at: null }] }));
    const region = screen.getByRole('region', { name: 'Frozen-destination storage history' });
    expect(region).toHaveTextContent('Bound to the verified frozen destination and revision.');
    expect(fact(region, 'Storage generation')).toHaveTextContent(/^2$/); expect(fact(region, 'Earlier destination history may be lost')).toHaveTextContent(/^Yes$/);
    expect(fact(region, 'Observed reconstructions')).toHaveTextContent(/^3$/); expect(fact(region, 'Observed repair state')).toHaveTextContent(/^Blocked$/);
    expect(fact(region, 'Last reconstruction time')).toHaveTextContent('2026-10-29T09:59:00.123456789Z');
    const events = within(region).getAllByRole('listitem'); expect(events).toHaveLength(5);
    expect(events[0]).toHaveTextContent('Confirmed'); expect(events[1]).toHaveTextContent('Recreated by confirmation'); expect(events[1]).toHaveTextContent('Unavailable');
    expect(events[2]).toHaveTextContent('Reconstructed'); expect(events[3]).toHaveTextContent('Repair blocked'); expect(events[4]).toHaveTextContent('Unavailable');
    expect(region).toHaveTextContent('Storage evidence does not grant repair permission or prove restored history.');
    expect(within(region).queryByRole('button')).not.toBeInTheDocument(); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it('distinguishes observed zero/false/empty storage from missing history', () => {
    pure(fixture(null, { generation: 0, history_may_be_lost: false, reconstructions: 0, last_reconstructed_at: null, repair_state: 'none', events: [] }));
    const region = screen.getByRole('region', { name: 'Frozen-destination storage history' });
    expect(fact(region, 'Storage generation')).toHaveTextContent(/^0$/); expect(fact(region, 'Observed reconstructions')).toHaveTextContent(/^0$/);
    expect(fact(region, 'Earlier destination history may be lost')).toHaveTextContent(/^No$/); expect(fact(region, 'Last reconstruction time')).toHaveTextContent(/^Unavailable$/);
    expect(region).toHaveTextContent('No storage events in this verified observation.');
    expect(region).not.toHaveTextContent('Frozen-destination storage history and events unavailable.');
  });
  it.each(['quota-extra', 'storage-extra', 'event-extra', 'negative-bytes', 'unknown-state'] as const)('rejects unsafe or invalid resource DTOs through actual service/parser: %s', kind => {
    const value = fixture(); const current = value.groups[0].delivery?.current; if (!current || !current.quota) throw new Error('Expected fixture current');
    if (kind === 'quota-extra') Object.assign(current.quota, { password: 'secret-quota /private/quota.sqlite' });
    if (kind === 'negative-bytes') current.quota.used_bytes = -1;
    if (kind === 'unknown-state') Object.assign(current.quota, { state: 'driver-secret' });
    if (kind === 'storage-extra' || kind === 'event-extra') {
      current.storage = { generation: 1, history_may_be_lost: false, reconstructions: 0, last_reconstructed_at: null, repair_state: 'none', events: [] };
      if (kind === 'storage-extra') Object.assign(current.storage, { raw_error: 'secret-quota /private/quota.sqlite' });
      else current.storage.events.push(Object.assign({ kind: 'confirmed' as const, at: null }, { sql: 'SELECT secret FROM credentials' }));
    }
    mount(value);
    return waitFor(() => { expect(screen.getByRole('alert')).toHaveTextContent('Recording evidence is unavailable');
      expect(document.body).not.toHaveTextContent(/secret-quota|private\/quota|SELECT secret|driver-secret/); expect(screen.queryByRole('region', { name: title })).not.toBeInTheDocument(); });
  });
  it('localizes quota/storage availability and scope in Traditional Chinese', async () => {
    await i18n.changeLanguage('zh-TW'); mount(fixture(quota({ scope: 'global' })));
    const region = await screen.findByRole('region', { name: '配額與儲存證據' });
    expect(region).toHaveTextContent('全域 — 所有群組與版本'); expect(region).toHaveTextContent('警告');
    expect(region).toHaveTextContent('凍結目的地的儲存歷史與事件無法取得。');
  });
});
