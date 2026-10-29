import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeRecordingGroupCard } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingGroupCard';
import type { CurrentDelivery, DeliveryEvidence, EvidenceReason, RecordingGroupView } from '../../../src/types/runtimeRuleEvidence';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';

function fixture(): RecordingGroupView {
  const group = ruleEvidenceFixture().groups[0];
  group.delivery = { current: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: false,
    stages: { collecting: 4, queued: 3, retrying: 2, blocked: 1, quarantined: 5, unknown: 6, sql_committed: 0, skipped: 7 }, last_commit: null, sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null },
    backlog: [{ binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'A0' }, destination: { connector_id: 'C1', connector_revision: '7', alias: null, kind: 'unavailable', schema: 'original', table: 'old_table' }, pending: 3, reasons: ['recording-ownership-unavailable'] }],
    counter_replay_quarantines: [{ workspace_id: 'W', group_id: 'G', group_revision: 'A0', pending_samples: 4, reason: 'counter-replay-proof-missing', action: 'preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance' }], oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, test_write: null };
  return group;
}
function delivery(group: RecordingGroupView): DeliveryEvidence { if (!group.delivery) throw new Error('Expected delivery fixture'); return group.delivery; }
function current(group: RecordingGroupView): CurrentDelivery { const value = delivery(group).current; if (!value) throw new Error('Expected current fixture'); return value; }
function fact(region: HTMLElement, label: string) { return within(region).getByText(label, { selector: 'dt' }).nextElementSibling; }
function mount(group = fixture()) { return render(<RuntimeRecordingGroupCard group={group} deviceId="D" />); }
function region(name: string) { return screen.getByRole('region', { name }); }
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);

describe('recording delivery stages, retained backlog and replay inputs', () => {
  it('CollectingQueuedNeverSQLCommitted: renders accepted samples separately from output rows without adding them', () => {
    mount(); const currentRegion = region('Current revision delivery');
    expect(fact(currentRegion, 'Delivery revision')).toHaveTextContent(/^A$/);
    expect(fact(currentRegion, 'Collecting')).toHaveTextContent(/^4 accepted samples$/);
    for (const [label, count] of [['Queued', 3], ['Retrying', 2], ['Blocked', 1], ['Finalized output quarantined', 5], ['Unknown output state', 6], ['SQL committed', 0], ['Skipped', 7]] as const) expect(fact(currentRegion, label)).toHaveTextContent(`${count} output rows`);
    expect(currentRegion).toHaveTextContent('Accepted samples and output rows use different units and are not added.');
    expect(currentRegion).not.toHaveTextContent(/Total pending|Successfully delivered/);
    expect(fact(currentRegion, 'Current output reason')).toHaveTextContent(/^Unavailable$/);
  });
  it('DisabledHistoricalDestinationNotRelabeled: retains C1 revision7 old_table / three rows despite disabled C2 draft', () => {
    const group = fixture(); if (!group.saved || !group.effective) throw new Error('Expected facts');
    group.saved.destination.connector_id = 'C2'; group.saved.destination.connector_revision = '8'; group.saved.destination.table = 'new_table';
    group.effective.enabled = false; current(group).intake = { group_revision: 'A', state: 'not_running', reason: null };
    mount(group); const backlog = region('Retained revision backlog');
    fireEvent.click(within(backlog).getByText('Backlog revision A0', { selector: 'summary' }));
    for (const [label, value] of [['Connector ID', 'C1'], ['Connector revision', '7'], ['Table', 'old_table'], ['Schema', 'original'], ['Pending output rows', '3 output rows'], ['Destination alias', 'Unavailable'], ['Destination kind', 'Unavailable']]) expect(fact(backlog, label)).toHaveTextContent(value);
    expect(backlog).not.toHaveTextContent(/C2|new_table/);
    expect(fact(region('Current revision delivery'), 'Observed intake state')).toHaveTextContent('Not running');
  });
  it('keeps separate frozen destinations within a retained revision without moving disclosure state', () => {
    const group = fixture(); const other = structuredClone(delivery(group).backlog[0]); other.destination.table = 'other_table'; other.destination.connector_revision = '8'; delivery(group).backlog.push(other);
    const view = mount(group); const summaries = within(region('Retained revision backlog')).getAllByText('Backlog revision A0', { selector: 'summary' });
    fireEvent.click(summaries[1]); expect(summaries[1].closest('details')).toHaveAttribute('open');
    delivery(group).backlog.shift(); view.rerender(<RuntimeRecordingGroupCard group={group} deviceId="D" />);
    const retained = within(region('Retained revision backlog')).getByText('Backlog revision A0', { selector: 'summary' });
    expect(retained.closest('details')).toHaveAttribute('open'); expect(retained.closest('details')).toHaveTextContent('other_table');
    expect(region('Retained revision backlog')).not.toHaveTextContent('old_table');
  });
  it('SharedFrozenMemberTotalsNotPerTag: labels ten committed rows as shared group totals without source success', () => {
    const group = fixture(); const c = current(group); c.shared = true; c.stages.sql_committed = 10;
    c.last_commit = { group_revision: 'A', connector_revision: '7', record_id: 'record', effect_key: 'effect', payload_digest: 'digest', committed_at: '2026-10-27T10:00:00Z', provenance: 'runtime' };
    group.effective?.members.push({ device_id: 'E', point_id: 'PE', tag_id: 'TE', column: 'other', source_value_type: null, circuit_id: null });
    mount(group); const currentRegion = region('Current revision delivery');
    expect(currentRegion).toHaveTextContent('Shared group totals across devices; not device or Tag success counts.');
    expect(fact(currentRegion, 'SQL committed')).toHaveTextContent(/^10 output rows$/);
    expect(currentRegion).not.toHaveTextContent(/record|digest|effect|readback|cleanup/i);
    expect(currentRegion.querySelectorAll('button,a')).toHaveLength(0);
  });
  it('uses frozen current shared flag rather than new saved members', () => {
    const group = fixture(); group.saved?.members.push({ device_id: 'E', point_id: 'PE', tag_id: 'TE', column: 'other', source_value_type: null, circuit_id: null });
    mount(group); expect(region('Current revision delivery')).toHaveTextContent('Group totals; not device or Tag success counts.');
    expect(region('Current revision delivery')).not.toHaveTextContent('Shared group totals');
  });
  it('ReplayInputsSeparateFromFinalizedQuarantine: preserves overlapping sample counts and finalized rows independently', () => {
    const group = fixture(); const before = structuredClone(group); mount(group);
    const replay = region('Accepted-input replay quarantine');
    expect(replay).toHaveTextContent('Replay samples may also appear in collecting; do not add them to finalized output rows.');
    const summary = within(replay).getByText('Replay revision A0', { selector: 'summary' }); summary.focus(); expect(summary).toHaveFocus(); fireEvent.click(summary);
    expect(summary.closest('details')).toHaveAttribute('open');
    expect(fact(replay, 'Pending accepted samples')).toHaveTextContent(/^4 accepted samples$/);
    expect(fact(region('Current revision delivery'), 'Finalized output quarantined')).toHaveTextContent(/^5 output rows$/);
    expect(replay).toHaveTextContent('Preserve accepted inputs and restore the original immutable revision evidence, or contact maintenance.');
    expect(group).toEqual(before); expect(region('Delivery evidence').querySelectorAll('button,a')).toHaveLength(0);
  });
  it.each([
    ['recording-ownership-unavailable', 'Recording ownership evidence unavailable'],
    ['recording-ownership-conflict', 'Recording ownership evidence conflicts'],
    ['recording-handoff-proof-conflict', 'Recording handoff evidence conflicts'],
    ['recording-cutover-finalized-overlap', 'Finalized recording windows overlap at cutover'],
    ['counter-replay-proof-missing', 'Counter replay evidence missing'],
    ['counter-replay-proof-conflict', 'Counter replay evidence conflicts'],
    ['counter-replay-proof-unsupported', 'Counter replay evidence unsupported'],
  ] as const)('HistoricalTypedOwnershipHandoffOverlapReasons: shows safe %s guidance', (reason, text) => {
    const group = fixture(); delivery(group).backlog[0].reasons = [reason]; mount(group);
    const backlog = region('Retained revision backlog'); expect(backlog).toHaveTextContent(text);
    expect(backlog).toHaveTextContent('Inspect the original immutable evidence; do not change or resend retained data.');
    expect(region('Current revision delivery')).not.toHaveTextContent(text);
  });
  it('maps an unexpected reason to safe unavailable without reflecting raw diagnostic text', () => {
    const group = fixture(); delivery(group).backlog[0].reasons = ['password=secret /private/driver' as EvidenceReason]; mount(group);
    expect(region('Retained revision backlog')).toHaveTextContent('Reason evidence unavailable');
    expect(document.body.textContent).not.toMatch(/password=secret|private\/driver/);
  });
  it('missing delivery is unavailable rather than zero or an empty-backlog success', () => {
    const group = fixture(); group.delivery = null; mount(group);
    expect(region('Delivery evidence')).toHaveTextContent('Delivery evidence unavailable');
    expect(region('Delivery evidence')).not.toHaveTextContent(/0 output rows|No retained backlog|No accepted-input/);
  });
  it('current null retains backlog and replay without inventing current zero counts or known active ordering', () => {
    const group = fixture(); delivery(group).current = null; group.state = 'unavailable'; group.proof = 'unavailable'; group.effective = null; group.effective_revision = null; mount(group);
    expect(region('Current revision delivery')).toHaveTextContent('Current revision delivery unavailable');
    expect(region('Current revision delivery')).not.toHaveTextContent('0 output rows');
    expect(region('Retained revision backlog')).toHaveTextContent('old_table');
    expect(region('Retained revision backlog')).not.toHaveTextContent(/older than active|previous active/i);
    expect(region('Accepted-input replay quarantine')).toHaveTextContent('4 accepted samples');
  });
  it('distinguishes unavailable replay evidence from an observed empty replay list', () => {
    const group = fixture(); delivery(group).counter_replay_quarantines = null; const view = mount(group);
    expect(region('Accepted-input replay quarantine')).toHaveTextContent('Replay quarantine evidence unavailable');
    delivery(group).counter_replay_quarantines = []; delivery(group).backlog = []; view.rerender(<RuntimeRecordingGroupCard group={group} deviceId="D" />);
    expect(region('Accepted-input replay quarantine')).toHaveTextContent('No accepted-input replay quarantine reported in this observation.');
    expect(region('Retained revision backlog')).toHaveTextContent('No retained backlog reported in this observation.');
    expect(region('Delivery evidence')).not.toHaveTextContent(/All delivered|Success/);
  });
  it('does not infer intake from disabled or ended snapshot facts or borrow intake reason for output', () => {
    const group = fixture(); if (!group.effective) throw new Error('Expected facts'); group.effective.enabled = false; group.effective.ends_at = group.as_of; mount(group);
    expect(fact(region('Current revision delivery'), 'Observed intake state')).toHaveTextContent(/^Unavailable$/);
    expect(region('Current revision delivery')).not.toHaveTextContent('Not running');
  });
  it('keeps exact intake reason distinct from unavailable current output reasons', () => {
    const group = fixture(); current(group).intake = { group_revision: 'A', state: 'blocked', reason: 'recording-ownership-conflict' }; mount(group);
    expect(fact(region('Current revision delivery'), 'Observed intake state')).toHaveTextContent('Blocked');
    expect(region('Current revision delivery')).toHaveTextContent('Recording ownership evidence conflicts');
    expect(fact(region('Current revision delivery'), 'Current output reason')).toHaveTextContent(/^Unavailable$/);
  });
  it('does not expose receipt, timing, test-write or per-source details in this slice', () => {
    const group = fixture(); current(group).next_retry_at = '2026-10-27T11:00:00Z'; delivery(group).oldest_pending_seconds = 123;
    delivery(group).test_write = { operation_id: 'test-operation', provenance: 'test-write', written_at: '2026-10-27T10:00:00Z', readback_at: null, cleanup_at: null }; mount(group);
    expect(region('Delivery evidence')).not.toHaveTextContent(/11:00:00|123|test-operation|readback|cleanup|next retry/i);
    const resources = region('Quota and storage evidence');
    for (const label of ['Quota scope', 'Quota state', 'Used durable bytes', 'Configured maximum bytes', 'New intake refused by quota']) expect(fact(resources, label)).toHaveTextContent(/^Unavailable$/);
    const storage = within(resources).getByRole('region', { name: 'Frozen-destination storage history' });
    expect(storage).toHaveTextContent('Frozen-destination storage history and events unavailable.');
    expect(within(storage).queryByRole('list')).not.toBeInTheDocument();
  });
  it('localizes delivery distinctions and native read-only disclosure in Traditional Chinese', async () => {
    await i18n.changeLanguage('zh-TW'); mount();
    expect(screen.getByRole('region', { name: '交付證據' })).toHaveTextContent('4 筆已接收樣本');
    const summary = screen.getByText('待交付版本 A0', { selector: 'summary' }); summary.focus(); expect(summary).toHaveFocus(); fireEvent.click(summary); expect(summary.closest('details')).toHaveAttribute('open');
    expect(screen.getByRole('region', { name: '已接收輸入的重播隔離' })).toHaveTextContent('保留已接收輸入');
  });
});
