import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeRecordingGroupCard } from '../../../src/features/datalink/runtime-dashboard/RuntimeRecordingGroupCard';
import type { RecordingGroupFacts, RecordingGroupView } from '../../../src/types/runtimeRuleEvidence';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';

function group(): RecordingGroupView { return structuredClone(ruleEvidenceFixture().groups[0]); }
function fact(section: HTMLElement, label: string) { return within(section).getByText(label, { selector: 'dt' }).nextElementSibling; }
function sections() { return { saved: screen.getByRole('region', { name: 'Saved configuration' }), effective: screen.getByRole('region', { name: 'Effective recording snapshot' }), requested: screen.getByRole('region', { name: 'Latest requested revision' }) }; }
function facts(view: RecordingGroupView): { saved: RecordingGroupFacts; effective: RecordingGroupFacts } {
  if (!view.saved || !view.effective) throw new Error('Expected complete test facts');
  return { saved: view.saved, effective: view.effective };
}
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);

describe('revision-bound recording group card', () => {
  it('keeps saved B meter_new / 30 seconds separate from effective A meter_old / 15 seconds', () => {
    const view = group(); const { saved } = facts(view);
    view.saved_revision = 'B'; saved.binding.entity_revision = 'B'; saved.destination.table = 'meter_new'; saved.interval_seconds = 30; saved.recording_profile = { kind: 'instantaneous', interval: 30, unit: 'second' };
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    const sectionsView = sections();
    expect(fact(sectionsView.saved, 'Fact revision')).toHaveTextContent(/^B$/);
    expect(fact(sectionsView.saved, 'Table')).toHaveTextContent(/^meter_new$/);
    expect(fact(sectionsView.saved, 'Recording interval')).toHaveTextContent(/^30 seconds$/);
    expect(fact(sectionsView.effective, 'Fact revision')).toHaveTextContent(/^A$/);
    expect(fact(sectionsView.effective, 'Table')).toHaveTextContent(/^meter_old$/);
    expect(fact(sectionsView.effective, 'Recording interval')).toHaveTextContent(/^15 seconds$/);
    expect(sectionsView.effective).not.toHaveTextContent('meter_new');
  });
  it('shows C saved, B requested and A effective without inventing B facts or a boundary', () => {
    const view = group(); view.applied_effective_at = null;
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    const { saved, effective, requested } = sections();
    expect(fact(saved, 'Fact revision')).toHaveTextContent(/^C$/);
    expect(fact(effective, 'Fact revision')).toHaveTextContent(/^A$/);
    expect(fact(requested, 'Requested revision')).toHaveTextContent(/^B$/);
    expect(fact(requested, 'Requested effective boundary')).toHaveTextContent(/^Unavailable$/);
    expect(requested).not.toHaveTextContent(/meter_draft|meter_old|15 seconds|60 seconds/);
    expect(requested).toHaveTextContent('A separate requested-revision snapshot is not supplied');
    expect(screen.getByRole('article', { name: 'Recording group G' })).toHaveTextContent('2026-10-27T10:00:00Z');
  });
  it('updates effective B only from the supplied refreshed snapshot while saved C remains unchanged', () => {
    const initial = group(); initial.applied_effective_at = null;
    const view = render(<RuntimeRecordingGroupCard group={initial} deviceId="D" />);
    const refreshed = ruleEvidenceFixture(undefined, 'B').groups[0];
    view.rerender(<RuntimeRecordingGroupCard group={refreshed} deviceId="D" />);
    const { saved, effective, requested } = sections();
    expect(fact(saved, 'Table')).toHaveTextContent(/^meter_draft$/);
    expect(fact(saved, 'Fact revision')).toHaveTextContent(/^C$/);
    expect(fact(effective, 'Table')).toHaveTextContent(/^meter_new$/);
    expect(fact(effective, 'Fact revision')).toHaveTextContent(/^B$/);
    expect(fact(effective, 'Effective boundary')).toHaveTextContent('2026-10-27T10:00:01Z');
    expect(fact(requested, 'Requested effective boundary')).toHaveTextContent('2026-10-27T10:00:01Z');
  });
  it('renders every safe destination, layout and policy from its own bound revision', () => {
    const view = group(); const { saved, effective } = facts(view);
    Object.assign(saved.destination, { connector_id: 'saved-connector', connector_revision: 'CR-C', alias: 'Draft store', kind: 'postgres', schema: 'draft_schema' });
    Object.assign(effective.destination, { connector_id: 'old-connector', connector_revision: 'CR-A', alias: 'Frozen store', kind: 'sqlite', schema: 'old_schema' });
    Object.assign(saved, { layout: 'managed_long_v1', storage_mode: 'timescale', interval_seconds: 300, reduction: 'cumulative_delta', write_mode: 'latest', incomplete_policy: 'skip_row', allowed_lateness_seconds: 7, recording_profile: { kind: 'cumulative', interval: 5, unit: 'minute' } });
    Object.assign(effective, { layout: 'managed_wide_profile_v1', storage_mode: 'native_monthly', reduction: 'mean', write_mode: 'append', incomplete_policy: 'partial', allowed_lateness_seconds: 2 });
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    const panels = sections();
    for (const [label, oldValue, newValue] of [
      ['Connector ID', 'old-connector', 'saved-connector'], ['Connector revision', 'CR-A', 'CR-C'],
      ['Destination alias', 'Frozen store', 'Draft store'], ['Destination kind', 'sqlite', 'postgres'], ['Schema', 'old_schema', 'draft_schema'],
      ['Layout', 'managed_wide_profile_v1', 'managed_long_v1'], ['Storage mode', 'native_monthly', 'timescale'],
      ['Reduction', 'mean', 'cumulative_delta'], ['Write mode', 'append', 'latest'], ['Incomplete row policy', 'partial', 'skip_row'],
      ['Allowed lateness', '2 seconds', '7 seconds'], ['Configured enabled state', 'Enabled in this configuration', 'Enabled in this configuration'],
      ['Recording profile', 'instantaneous · 15 seconds', 'cumulative · 5 minutes'],
    ]) {
      expect(fact(panels.effective, label)).toHaveTextContent(oldValue);
      expect(fact(panels.saved, label)).toHaveTextContent(newValue);
    }
  });
  it.each(['saved-only', 'unavailable'] as const)('does not promote saved or unproven facts for %s', (state) => {
    const view = group(); view.state = state; view.proof = 'unavailable';
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    const { saved, effective } = sections();
    expect(saved).toHaveTextContent('meter_draft');
    expect(effective).toHaveTextContent('Effective recording facts unavailable');
    expect(effective).not.toHaveTextContent(/meter_old|meter_draft/);
  });
  it('does not display an effective facts object without the verified frozen proof', () => {
    const view = group(); view.proof = 'unavailable';
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    expect(sections().effective).not.toHaveTextContent('meter_old');
  });
  it('does not label an unproven retained revision ID as effective in the header', () => {
    const view = group(); view.proof = 'unavailable'; view.effective = null;
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    const article = screen.getByRole('article', { name: 'Recording group G' });
    expect(fact(article, 'Effective revision')).toHaveTextContent(/^Unavailable$/);
    expect(fact(article, 'Saved revision')).toHaveTextContent(/^C$/);
    expect(fact(sections().requested, 'Requested revision')).toHaveTextContent(/^B$/);
  });
  it('keeps complete other-device membership in an explicit shared disclosure without relabeling it', () => {
    const view = group(); const { saved, effective } = facts(view);
    saved.members.push({ device_id: 'OTHER-C', point_id: 'PC', tag_id: 'TC', column: 'draft_column', source_value_type: 'int32', circuit_id: 'circuit-C' });
    effective.members.push({ device_id: 'OTHER-A', point_id: 'PA', tag_id: 'TA', column: 'old_column', source_value_type: 'float64', circuit_id: 'circuit-A' });
    const before = structuredClone(view);
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    const { saved: savedSection, effective: effectiveSection } = sections();
    expect(savedSection).toHaveTextContent('Group-wide configuration; not per-device or per-Tag delivery success.');
    for (const [section, device, point, tag, column, circuit] of [
      [savedSection, 'OTHER-C', 'PC', 'TC', 'draft_column', 'circuit-C'],
      [effectiveSection, 'OTHER-A', 'PA', 'TA', 'old_column', 'circuit-A'],
    ] as const) {
      const selected = within(section).getByRole('region', { name: 'Members on selected device D' });
      expect(selected).toHaveTextContent('P'); expect(selected).not.toHaveTextContent(device);
      const summary = within(section).getByText('Members on other devices (1)', { selector: 'summary' });
      summary.focus(); expect(summary).toHaveFocus(); fireEvent.click(summary);
      expect(summary.closest('details')).toHaveAttribute('open');
      for (const identity of [device, point, tag, column, circuit]) expect(summary.closest('details')).toHaveTextContent(identity);
    }
    expect(view).toEqual(before);
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
  });
  it('does not move saved-only membership into effective selected-device membership', () => {
    const view = group(); const { effective } = facts(view); effective.members[0].device_id = 'OTHER';
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    expect(within(sections().effective).getByRole('region', { name: 'Members on selected device D' })).toHaveTextContent('No members on this device in this revision.');
    expect(sections().saved).not.toHaveTextContent('No members on this device in this revision.');
  });
  it('marks nullable and unavailable facts explicitly without synthesized defaults or countdowns', () => {
    const view = group(); const { saved } = facts(view);
    Object.assign(saved, { recording_profile: null, interval_seconds: null, allowed_lateness_seconds: null, layout: 'unavailable', storage_mode: 'unavailable', reduction: 'unavailable', write_mode: 'unavailable', incomplete_policy: 'unavailable' });
    Object.assign(saved.destination, { alias: null, kind: 'unavailable', schema: null });
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    for (const label of ['Destination alias', 'Destination kind', 'Schema', 'Recording profile', 'Recording interval', 'Allowed lateness', 'Layout', 'Storage mode', 'Reduction', 'Write mode', 'Incomplete row policy', 'Snapshot ends at']) expect(fact(sections().saved, label)).toHaveTextContent(/^Unavailable$/);
    expect(sections().saved).not.toHaveTextContent(/next due|countdown|0 seconds/i);
  });
  it.each(['2026-10-27T09:00:00Z', '2026-10-27T11:00:00Z'])('retains ended disabled snapshot facts at %s without claiming intake or SQL success', (endsAt) => {
    const view = group(); const { effective } = facts(view); effective.enabled = false; effective.ends_at = endsAt;
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    expect(fact(sections().effective, 'Configured enabled state')).toHaveTextContent('Disabled in this configuration');
    expect(fact(sections().effective, 'Snapshot ends at')).toHaveTextContent(endsAt);
    expect(sections().effective).toHaveTextContent('meter_old');
    expect(screen.getByText('Verified configuration is not proof of current intake or SQL delivery.')).toBeInTheDocument();
    expect(screen.queryByText(/^(Running|Intake active|SQL committed|Success|LIVE)$/)).not.toBeInTheDocument();
  });
  it('does not present mismatched effective fact revisions as proven', () => {
    const view = group(); facts(view).effective.binding.entity_revision = 'wrong-revision';
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    expect(sections().effective).toHaveTextContent('Effective recording facts unavailable');
    expect(sections().effective).not.toHaveTextContent('meter_old');
  });
  it('renders null facts honestly and localizes the read-only comparison in Traditional Chinese', async () => {
    await i18n.changeLanguage('zh-TW'); const view = group(); view.saved = null; view.effective = null;
    render(<RuntimeRecordingGroupCard group={view} deviceId="D" />);
    expect(screen.getByRole('region', { name: '已儲存設定' })).toHaveTextContent('尚無證據');
    expect(screen.getByRole('region', { name: '生效記錄快照' })).toHaveTextContent('生效記錄設定尚無證據');
    expect(screen.getByRole('region', { name: '最新請求版本' })).toHaveTextContent('B');
  });
});
