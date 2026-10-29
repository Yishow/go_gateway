import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeRecordingGroupCard } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingGroupCard';
import type { CurrentDelivery, RecordingGroupView, RuntimeReceipt, SourceResult } from '../../../src/types/runtimeRuleEvidence';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
function receipt(record: string): RuntimeReceipt { return { group_revision: 'A', connector_revision: '7', record_id: record, effect_key: `effect-${record}`, payload_digest: `digest-${record}`, committed_at: '2026-10-27T09:59:00Z', provenance: 'runtime' }; }
function result(state: SourceResult['state'] = 'partial'): SourceResult { return { state, window_start: '2026-10-27T09:59:30Z', window_end: '2026-10-27T09:59:45Z', coverage: state === 'partial' ? 'partial' : 'unknown', delta_known: null, receipt: null }; }
function fixture(): RecordingGroupView {
  const group = ruleEvidenceFixture().groups[0];
  const committed = result('available'); committed.window_start = '2026-10-27T09:58:30Z'; committed.window_end = '2026-10-27T09:58:45Z'; committed.coverage = 'complete'; committed.receipt = receipt('source-previous');
  group.delivery = { current: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: false, stages: { collecting: 1, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 10, skipped: 0 }, last_commit: receipt('group-latest'), sources: [{ device_id: 'D', point_id: 'P', tag_id: 'T', latest: result(), last_committed: committed }], source_state: 'available', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null }, backlog: [], counter_replay_quarantines: [], oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, test_write: null };
  return group;
}
function current(group: RecordingGroupView): CurrentDelivery { if (!group.delivery?.current) throw new Error('Expected fixture'); return group.delivery.current; }
function fact(region: HTMLElement, label: string) { return within(region).getByText(label, { selector: 'dt' }).nextElementSibling; }
function mount(group = fixture(), deviceId = 'D') { return render(<RuntimeRecordingGroupCard group={group} deviceId={deviceId} />); }
function details() { return screen.getByRole('region', { name: 'Source and receipt details' }); }
function source(region = details()) { const summary = within(region).getByText('Source P · Tag T', { selector: 'summary' }); fireEvent.click(summary); const element = summary.closest('details'); if (!element) throw new Error('Expected details'); return element; }
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);

describe('independent runtime receipt and selected-source evidence', () => {
  it('LatestQueuedSeparateFromEarlierCommittedWindow: keeps newer partial latest separate from earlier committed receipt', () => {
    mount(); const selected = source(); const latest = within(selected).getByRole('region', { name: 'Latest observed source result' }); const committed = within(selected).getByRole('region', { name: 'Last committed source result' });
    expect(fact(latest, 'Result state')).toHaveTextContent(/^Partial$/); expect(fact(latest, 'Window start')).toHaveTextContent('09:59:30Z'); expect(fact(latest, 'Coverage')).toHaveTextContent(/^Partial$/);
    expect(fact(latest, 'Record ID')).toHaveTextContent(/^Unavailable$/); expect(latest).not.toHaveTextContent(/source-previous|group-latest/);
    expect(fact(committed, 'Window start')).toHaveTextContent('09:58:30Z'); expect(fact(committed, 'Record ID')).toHaveTextContent(/^source-previous$/); expect(fact(committed, 'Coverage')).toHaveTextContent(/^Complete$/);
  });
  it('labels the latest runtime receipt group-wide and does not substitute its identity for source receipts', () => {
    mount(); const summary = within(details()).getByText('Latest runtime receipt (group-wide)', { selector: 'summary' }); summary.focus(); expect(summary).toHaveFocus(); fireEvent.click(summary);
    const groupReceipt = summary.closest('details'); expect(groupReceipt).toHaveAttribute('open'); if (!groupReceipt) throw new Error('Expected receipt');
    expect(groupReceipt).toHaveTextContent('This receipt belongs to the whole group; it does not establish every source or Tag result.');
    for (const [label, value] of [['Record ID', 'group-latest'], ['Effect key', 'effect-group-latest'], ['Payload digest', 'digest-group-latest'], ['Receipt group revision', 'A'], ['Receipt connector revision', '7'], ['Runtime committed at', '2026-10-27T09:59:00Z'], ['Evidence provenance', 'Runtime']]) expect(fact(groupReceipt, label)).toHaveTextContent(value);
    expect(source()).toHaveTextContent('source-previous'); expect(source()).not.toHaveTextContent('group-latest');
  });
  it.each(['unknown', 'partial', 'missing'] as const)('MeanPartialAndUnknownCoverageNeverComplete: preserves %s evidence', (state) => {
    const group = fixture(); const latest = result(state); latest.coverage = state === 'missing' ? null : state; current(group).sources[0].latest = latest; mount(group);
    const region = within(source()).getByRole('region', { name: 'Latest observed source result' });
    expect(fact(region, 'Result state')).toHaveTextContent(state[0].toUpperCase() + state.slice(1));
    expect(fact(region, 'Coverage')).not.toHaveTextContent('Complete'); expect(fact(region, 'Record ID')).toHaveTextContent('Unavailable');
  });
  it('DeltaUnknownNeverZero: false delta remains unknown and is never rendered as a numeric value', () => {
    const group = fixture(); current(group).sources[0].latest = { ...result('unknown'), delta_known: false }; mount(group);
    const region = within(source()).getByRole('region', { name: 'Latest observed source result' });
    expect(fact(region, 'Delta evidence')).toHaveTextContent(/^Unknown; not zero$/); expect(fact(region, 'Result state')).toHaveTextContent(/^Unknown$/);
    expect(region.querySelectorAll('input')).toHaveLength(0); expect(fact(region, 'Delta evidence')).not.toHaveTextContent(/^0$/);
  });
  it('MissingRecordNoZeroWindowOrReceipt: null latest does not borrow last-committed evidence', () => {
    const group = fixture(); current(group).sources[0].latest = null; mount(group);
    const region = within(source()).getByRole('region', { name: 'Latest observed source result' });
    for (const label of ['Result state', 'Window start', 'Window end', 'Coverage', 'Delta evidence', 'Record ID', 'Runtime committed at']) expect(fact(region, label)).toHaveTextContent(/^Unavailable$/);
    expect(region).not.toHaveTextContent(/source-previous|group-latest|1970|09:58:30/);
  });
  it('SameGroupOtherDeviceNeverGetsSelectedSourceSuccess: filters exact selected tuples without borrowing another device', () => {
    const group = fixture(); const c = current(group); c.shared = true;
    group.effective?.members.push({ device_id: 'E', point_id: 'PE', tag_id: 'TE', column: 'other', source_value_type: null, circuit_id: null });
    c.sources.push({ device_id: 'E', point_id: 'PE', tag_id: 'TE', latest: null, last_committed: { ...c.sources[0].last_committed!, receipt: receipt('other-device') } });
    const before = structuredClone(group); mount(group);
    expect(details()).toHaveTextContent('Source P · Tag T'); expect(details()).not.toHaveTextContent(/Source PE|Tag TE|other-device/); expect(group).toEqual(before);
  });
  it('does not allocate ten committed group rows when the selected device has no source result', () => {
    const group = fixture(); current(group).sources = []; mount(group);
    expect(details()).toHaveTextContent('No source result reported for this device in this observation.');
    expect(within(details()).getByRole('region', { name: 'Source evidence for selected device D' })).not.toHaveTextContent(/10|0 successful|Source P/);
  });
  it('LegacyCustomMissingFrozenProofUnavailable: missing delivery or source proof stays unavailable', () => {
    const group = fixture(); current(group).source_state = 'unavailable'; current(group).sources = []; const rendered = mount(group);
    expect(details()).toHaveTextContent('Selected-source evidence unavailable'); expect(details()).not.toHaveTextContent('No source result reported');
    group.delivery = null; rendered.rerender(<RuntimeRecordingGroupCard group={group} deviceId="D" />);
    expect(details()).toHaveTextContent('Selected-source evidence unavailable'); expect(details()).not.toHaveTextContent('group-latest');
  });
  it('TestWriteRuntimeReadbackCleanupIndependent: production null leaves each phase explicitly unavailable', () => {
    mount(); const testWrite = within(details()).getByRole('region', { name: 'Test-write evidence' });
    for (const label of ['Test-write operation ID', 'Test write recorded at', 'Independent readback verified at', 'Owned cleanup completed at']) expect(fact(testWrite, label)).toHaveTextContent(/^Unavailable$/);
    expect(testWrite).not.toHaveTextContent(/group-latest|09:59:00/);
  });
  it.each(['written_at', 'readback_at', 'cleanup_at'] as const)('consumes only the supplied %s test-write phase timestamp, not inferred success', (phase) => {
    const group = fixture(); if (!group.delivery) throw new Error('Expected fixture');
    group.delivery.test_write = { operation_id: 'operation-one', provenance: 'test-write', written_at: null, readback_at: null, cleanup_at: null, [phase]: '2026-10-27T11:00:00Z' }; mount(group);
    const testWrite = within(details()).getByRole('region', { name: 'Test-write evidence' });
    const labels = { written_at: 'Test write recorded at', readback_at: 'Independent readback verified at', cleanup_at: 'Owned cleanup completed at' };
    for (const key of Object.keys(labels) as (keyof typeof labels)[]) expect(fact(testWrite, labels[key])).toHaveTextContent(key === phase ? '2026-10-27T11:00:00Z' : 'Unavailable');
    expect(testWrite).toHaveTextContent('operation-one'); expect(testWrite).not.toHaveTextContent('Runtime committed at');
  });
  it('does not expose raw values, source mutation actions or scheduler timing', () => {
    const group = fixture(); Object.assign(current(group).sources[0].latest!, { value: 987654321, raw_values: 'secret-value' }); current(group).next_retry_at = '2099-01-01T00:00:00Z'; mount(group);
    expect(details().textContent).not.toMatch(/987654321|secret-value|2099|next retry|countdown/i); expect(details().querySelectorAll('button,a,input')).toHaveLength(0);
  });
  it('localizes independent source and test-write evidence in Traditional Chinese', async () => {
    await i18n.changeLanguage('zh-TW'); mount(); const region = screen.getByRole('region', { name: '來源與收據詳情' });
    const summary = within(region).getByText('來源 P · Tag T', { selector: 'summary' }); summary.focus(); expect(summary).toHaveFocus(); fireEvent.click(summary); expect(summary.closest('details')).toHaveAttribute('open');
    expect(region).toHaveTextContent('部分'); expect(within(region).getByRole('region', { name: '測試寫入證據' })).toHaveTextContent('尚無證據');
  });
});
