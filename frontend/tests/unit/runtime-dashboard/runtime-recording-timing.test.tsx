import { act, cleanup, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeRecordingGroupCard } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingGroupCard';
import type { RecordingGroupView } from '../../../src/types/runtimeRuleEvidence';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
function fixture(): RecordingGroupView {
  const group = ruleEvidenceFixture().groups[0]; group.delivery = { oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: 12.375, current: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: false, stages: { collecting: 1, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, last_commit: null, sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null }, backlog: [], counter_replay_quarantines: [], test_write: null }; return group;
}
function fact(label: string) { return within(screen.getByRole('region', { name: 'Observed timing' })).getByText(label, { selector: 'dt' }).nextElementSibling; }
function mount(group = fixture()) { return render(<RuntimeRecordingGroupCard group={group} deviceId="D" />); }
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(() => { cleanup(); vi.useRealTimers(); });

describe('evidence-backed nullable timing display', () => {
  it('NoSchedulingEvidenceMeansUnavailableNot10_00_15From15SecondInterval', () => {
    vi.useFakeTimers(); vi.setSystemTime(new Date('2026-10-27T10:00:01Z')); mount();
    expect(fact('Oldest pending age (all group revisions)')).toHaveTextContent(/^12.375 seconds$/);
    for (const field of ['Next due at', 'Next retry at', 'Attempt duration']) expect(fact(field)).toHaveTextContent(/^Unavailable$/);
    expect(screen.getByRole('region', { name: 'Observed timing' })).not.toHaveTextContent('10:00:15');
  });
  it('DisabledAndMissingEffectiveProofRetainHistoricalAge without attributing it to the selected device or current revision', () => {
    const group = fixture(); if (!group.delivery || !group.effective) throw new Error('Expected fixture'); group.effective.enabled = false; group.delivery.current = null; group.effective = null; group.effective_revision = null; group.proof = 'unavailable'; group.state = 'unavailable'; mount(group);
    expect(fact('Oldest pending age (all group revisions)')).toHaveTextContent(/^12.375 seconds$/);
    expect(screen.getByRole('region', { name: 'Observed timing' })).toHaveTextContent('Oldest pending age covers all revisions of this group, not only this device or the current revision.');
    expect(fact('Next retry at')).toHaveTextContent(/^Unavailable$/);
  });
  it('SuppliedNullableTimingDisplaysExactValuesWithoutDerivation for existing-DTO fixture evidence', () => {
    const group = fixture(); const current = group.delivery?.current; if (!current) throw new Error('Expected fixture');
    current.next_due_at = '2026-10-27T10:01:17.123456789Z'; current.next_retry_at = '2026-10-27T10:02:19Z'; current.attempt_duration_ms = 27.625; mount(group);
    expect(fact('Next due at')).toHaveTextContent(/^2026-10-27T10:01:17.123456789Z$/); expect(fact('Next retry at')).toHaveTextContent(/^2026-10-27T10:02:19Z$/); expect(fact('Attempt duration')).toHaveTextContent(/^27.625 ms$/);
  });
  it('keeps explicitly supplied zero ages and durations distinct from nullable absence', () => {
    const group = fixture(); if (!group.delivery?.current) throw new Error('Expected fixture'); group.delivery.oldest_pending_seconds = 0; group.delivery.current.attempt_duration_ms = 0; const view = mount(group);
    expect(fact('Oldest pending age (all group revisions)')).toHaveTextContent(/^0 seconds$/); expect(fact('Attempt duration')).toHaveTextContent(/^0 ms$/);
    group.delivery.oldest_pending_seconds = null; group.delivery.current.attempt_duration_ms = null; view.rerender(<RuntimeRecordingGroupCard group={group} deviceId="D" />);
    expect(fact('Oldest pending age (all group revisions)')).toHaveTextContent(/^Unavailable$/); expect(fact('Attempt duration')).toHaveTextContent(/^Unavailable$/);
  });
  it('ObservedAgeDoesNotTickWithBrowserTime and changes only when supplied observation changes', () => {
    vi.useFakeTimers(); vi.setSystemTime(new Date('2026-10-27T10:00:01Z')); const group = fixture(); const view = mount(group);
    act(() => { vi.advanceTimersByTime(120_000); }); expect(fact('Oldest pending age (all group revisions)')).toHaveTextContent(/^12.375 seconds$/);
    if (!group.delivery) throw new Error('Expected fixture'); group.delivery.oldest_pending_seconds = 13.125; view.rerender(<RuntimeRecordingGroupCard group={group} deviceId="D" />);
    expect(fact('Oldest pending age (all group revisions)')).toHaveTextContent(/^13.125 seconds$/);
    expect(screen.getByRole('region', { name: 'Observed timing' })).toHaveTextContent('Observed values from the last read; not a live clock or countdown.');
  });
  it('missing delivery does not imply no pending rows or synthesize timestamps', () => {
    const group = fixture(); group.delivery = null; mount(group);
    for (const field of ['Oldest pending age (all group revisions)', 'Next due at', 'Next retry at', 'Attempt duration']) expect(fact(field)).toHaveTextContent(/^Unavailable$/);
    const region = screen.getByRole('region', { name: 'Observed timing' }); expect(region).not.toHaveTextContent(/0 seconds|No pending|1970/); expect(region.querySelectorAll('button,a,input')).toHaveLength(0);
  });
  it('localizes observed timing and unavailable fields in Traditional Chinese', async () => {
    await i18n.changeLanguage('zh-TW'); mount(); const region = screen.getByRole('region', { name: '觀察到的時間證據' });
    expect(region).toHaveTextContent('12.375 秒'); expect(region).toHaveTextContent('此群組的所有版本'); expect(region).toHaveTextContent('尚無證據');
  });
});
