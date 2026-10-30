import { afterEach, describe, expect, expectTypeOf, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { RecordingRuleSummary, type RecordingRuleSummaryLabels, type RecordingRuleSummaryProps,
  type RecordingRuleDetailSlot } from '@/features/datalink/workbench-v2/steps/step4/RecordingRuleSummary';
import { recordingWorkflowView, type RecordingWorkflowScope } from '@/features/datalink/workbench-v2/state/recordingWorkflowView';
import type { RuleEvidenceView } from '@/types/runtimeRuleEvidence';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';

const scope: RecordingWorkflowScope = { workspaceId: 'W1', deviceId: 'A', circuitId: 'C1', groupId: 'G1', workspaceRevision: 'W2', groupRevision: 'R2', appliedRevision: 'R1' };
const labels: RecordingRuleSummaryLabels = {
  heading: 'Recording rules', saved: 'Saved configuration', effective: 'Effective configuration',
  savedDescription: 'Saved is not an activation or SQL receipt.', effectiveDescription: 'Immutable facts at the server observation, not SQL success.',
  unavailable: 'Unavailable', emptySources: 'No sources in this version', normalizedInterval: 'Normalized interval', groupScope: 'Group-level scope, not exclusive device ownership',
  requestedNotEffective: 'The requested revision is not the observed effective revision.', details: 'Open rules and delivery details', detailsTitle: 'Scoped rules and evidence',
  detailsDescription: 'Read-only existing bindings, policy and delivery evidence.', detailsUnavailable: 'Current same-scope details are unavailable.', close: 'Close details',
  fields: { sources: 'Sources and types', layout: 'Layout', storageMode: 'Storage mode', profile: 'Profile', reduction: 'Reduction', interval: 'Interval',
    lateness: 'Allowed lateness', incomplete: 'Incomplete rows', destination: 'Destination table', ownership: 'Ownership context', workspace: 'Workspace', device: 'Selected device',
    circuit: 'Selected circuit', group: 'Group', savedRevision: 'Saved revision', requestedRevision: 'Requested revision', effectiveRevision: 'Effective revision',
    asOf: 'Server observation', effectiveAt: 'Effective from', endsAt: 'Ends at' },
  profiles: { instantaneous: 'Instantaneous', cumulative: 'Cumulative' }, units: { second: 'seconds', minute: 'minutes' },
  layouts: { managed_wide_v1: 'Wide', managed_long_v1: 'Long', managed_wide_profile_v1: 'Profile wide', unavailable: 'Unavailable' },
  storageModes: { plain: 'Plain table', timescale: 'Timescale', native_monthly: 'Native monthly', unavailable: 'Unavailable' },
  reductions: { latest: 'Latest sample', mean: 'Mean', cumulative_delta: 'Cumulative delta', unavailable: 'Unavailable' },
  incompletePolicies: { skip_row: 'Skip incomplete row', partial: 'Retain partial row', unavailable: 'Unavailable' },
  sourceTypes: { bool: 'Boolean', int16: 'Int16', uint16: 'Uint16', int32: 'Int32', uint32: 'Uint32', int64: 'Int64', uint64: 'Uint64', float32: 'Float32', float64: 'Float64', string: 'Text' },
  reasons: { missing: 'Evidence missing', loading: 'Evidence loading', error: 'Evidence error', unknown: 'Evidence unknown',
    scope_mismatch: 'Scope unavailable', revision_mismatch: 'Revision mismatch', effective_snapshot_unavailable: 'Effective facts unavailable',
    operation_unresolved: 'Check unresolved operation', stale: 'Evidence stale', invalid_evidence: 'Evidence invalid' },
};
function saved(): WriteGroup {
  return { id: 'G1', workspace_id: 'W1', revision: 'R2', applied_revision: 'R1', name: 'PRIVATE-NAME', status: 'ready', basic_managed_device_id: 'A',
    members: [{ device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1', source_revision: 'S2', mapping_revision: 'B1', target_column: 'new_value', required: true, source_value_type: 'float64' }],
    destination: { connector_id: 'DB2', connector_revision: 'D2', database: '/private/database.sqlite', table_schema: 'public', table_name: 'new_readings', storage_strategy: 'managed' },
    recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' }, row_policy: { interval_seconds: 30, allowed_lateness_seconds: 10, incomplete_policy: 'partial' },
    write_policy: {}, migration: {}, created_at: '2026-10-01T00:00:00Z', updated_at: '2026-10-02T00:00:00Z' };
}
function evidence(): RuleEvidenceView {
  return { contract_version: 1, workspace_id: 'W1', workspace_revision: 'W2', device_id: 'A', observed_at: '2026-10-02T00:00:00Z', freshness: 'fresh', source: null, locator: null,
    groups: [{ group_id: 'G1', saved_revision: 'R2', applied_revision: 'R1', effective_revision: 'R1', applied_effective_at: null, as_of: '2026-10-02T00:00:00Z', state: 'effective', proof: 'applied-snapshot-and-layout',
      saved: { binding: { workspace_id: 'W1', workspace_revision: 'W2', entity_revision: 'R2' }, enabled: true, effective_at: null, ends_at: null,
        recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second' }, allowed_lateness_seconds: 10,
        destination: { connector_id: 'DB2', connector_revision: 'D2', alias: 'New archive', kind: 'postgres', schema: 'public', table: 'new_readings' },
        members: [{ device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1', column: 'new_value', source_value_type: 'float64' }],
        layout: 'managed_wide_profile_v1', storage_mode: 'plain', interval_seconds: 30, reduction: 'mean', write_mode: 'append', incomplete_policy: 'partial' },
      effective: { binding: { workspace_id: 'W1', workspace_revision: 'W2', entity_revision: 'R1' }, enabled: true, effective_at: '2026-10-01T00:00:00Z', ends_at: null,
        recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' }, allowed_lateness_seconds: 2,
        destination: { connector_id: 'DB1', connector_revision: 'D1', alias: 'Old archive', kind: 'postgres', schema: 'archive', table: 'old_readings' },
        members: [{ device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1', column: 'old_value', source_value_type: 'int16' }],
        layout: 'managed_wide_v1', storage_mode: 'plain', interval_seconds: 15, reduction: 'latest', write_mode: 'append', incomplete_policy: 'skip_row' }, delivery: null }] };
}
function view(proof = evidence(), selected = scope, savedGroup = saved()) {
  return recordingWorkflowView({ scope: selected, savedGroup, runtimeEvidence: { state: 'available', scope: selected, value: proof } });
}
function detailSlot(): RecordingRuleDetailSlot {
  return { scope, effectiveRevision: 'R1', asOf: '2026-10-02T00:00:00Z', content: <p>Source binding B1; applied revision R1; counter epoch E1; existing policy details</p> };
}
function field(region: HTMLElement, label: string) {
  return within(region).getByText(label, { selector: 'dt' }).nextElementSibling;
}
function show(proof = evidence()) {
  render(<RecordingRuleSummary view={view(proof)} scopeLabel="Device A / Circuit C1" labels={labels} />);
  return screen.getByRole('region', { name: 'Effective configuration' });
}
afterEach(() => { cleanup(); vi.restoreAllMocks(); });

describe('unmounted recording rule summary', () => {
  it('uses exact identity display labels without exposing source IDs or SQL columns in the collapsed summary', () => {
    const sourceLabels = { [JSON.stringify(['A', 'P1', 'T1'])]: '供電電流' };
    render(<RecordingRuleSummary {...{ view: view(), scopeLabel: 'Main panel', labels, sourceLabels }} />);
    const sources = field(screen.getByRole('region', { name: labels.effective }), labels.fields.sources)!;
    expect(sources).toHaveTextContent('供電電流');
    expect(sources).not.toHaveTextContent(/T1|P1|old_value/);
  });
  it('does not reuse a same-tag label from a different device or point, or expose unsafe display labels', () => {
    const sourceLabels = { [JSON.stringify(['B', 'P1', 'T1'])]: 'Wrong device',
      [JSON.stringify(['A', 'P2', 'T1'])]: 'Wrong point', [JSON.stringify(['A', 'P1', 'T1'])]: '/private/archive.sqlite' };
    render(<RecordingRuleSummary {...{ view: view(), scopeLabel: 'Main panel', labels, sourceLabels }} />);
    const sources = field(screen.getByRole('region', { name: labels.effective }), labels.fields.sources)!;
    expect(sources).not.toHaveTextContent(/Wrong|private|archive.sqlite|T1|P1|old_value/);
    expect(sources).toHaveTextContent(labels.unavailable);
  });
  it('labels saved R2/30 seconds separately from effective R1/15 seconds with each version’s own types and table', () => {
    render(<RecordingRuleSummary view={view()} scopeLabel="Device A / Circuit C1" labels={labels} />);
    const draft = within(screen.getByRole('region', { name: 'Saved configuration' }));
    const active = within(screen.getByRole('region', { name: 'Effective configuration' }));
    expect(draft.getByText('30 seconds')).toBeInTheDocument(); expect(active.getByText('15 seconds')).toBeInTheDocument();
    expect(draft.getByText(/Float64/)).toBeInTheDocument(); expect(active.getByText(/Int16/)).toBeInTheDocument();
    expect(draft.getByText(/public.new_readings/)).toBeInTheDocument(); expect(active.getByText(/archive.old_readings/)).toBeInTheDocument();
    expect(document.body).not.toHaveTextContent(/PRIVATE-NAME|database.sqlite/);
  });
  it('does not render stale current/effective facts or silently use the editable saved group as summary proof', () => {
    const proof = evidence(); proof.freshness = 'stale';
    render(<RecordingRuleSummary view={view(proof)} scopeLabel="Device A / Circuit C1" labels={labels} />);
    expect(screen.getAllByRole('status').some(node => node.textContent?.includes('Evidence stale'))).toBe(true);
    expect(document.body).not.toHaveTextContent(/30 seconds|15 seconds|new_readings|old_readings|Float64|Int16/);
  });
  it('shows source/type, layout/profile/reduction, lateness/incomplete rules and safe group context from that exact version', () => {
    const active = show();
    expect(field(active, 'Layout')).toHaveTextContent('Wide'); expect(field(active, 'Profile')).toHaveTextContent('Instantaneous');
    expect(field(active, 'Reduction')).toHaveTextContent('Latest sample'); expect(field(active, 'Allowed lateness')).toHaveTextContent('2 seconds');
    expect(field(active, 'Incomplete rows')).toHaveTextContent('Skip incomplete row');
    expect(field(active, 'Ownership context')).toHaveTextContent('Group-level scope, not exclusive device ownership');
    const draft = screen.getByRole('region', { name: 'Saved configuration' });
    expect(field(draft, 'Layout')).toHaveTextContent('Profile wide'); expect(field(draft, 'Reduction')).toHaveTextContent('Mean');
    expect(field(draft, 'Allowed lateness')).toHaveTextContent('10 seconds'); expect(field(draft, 'Incomplete rows')).toHaveTextContent('Retain partial row');
    expect(active).toHaveTextContent('Immutable facts at the server observation, not SQL success.');
    expect(draft).toHaveTextContent('Saved is not an activation or SQL receipt.');
  });
  it('shows the proven cumulative profile period in minutes without claiming a physical measurement unit', () => {
    const proof = evidence(); const facts = proof.groups[0].effective!;
    facts.recording_profile = { kind: 'cumulative', interval: 10, unit: 'minute' }; facts.interval_seconds = 600;
    facts.reduction = 'cumulative_delta'; facts.layout = 'managed_long_v1'; facts.storage_mode = 'timescale'; facts.members[0].source_value_type = 'uint64';
    const active = show(proof);
    expect(field(active, 'Interval')).toHaveTextContent('10 minutes'); expect(field(active, 'Profile')).toHaveTextContent('Cumulative');
    expect(field(active, 'Reduction')).toHaveTextContent('Cumulative delta'); expect(field(active, 'Storage mode')).toHaveTextContent('Timescale');
    expect(active).toHaveTextContent('Uint64'); expect(active).not.toHaveTextContent(/kWh|600 seconds|Float64/);
  });
  it.each(['bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string'] as const)(
    'retains selected source type %s without changing or guessing it', (type) => {
      const proof = evidence(); proof.groups[0].effective!.members[0].source_value_type = type;
      expect(show(proof)).toHaveTextContent(labels.sourceTypes[type]);
    });
  it('does not replace unknown effective profile/type/lateness with newer saved values or zero', () => {
    const proof = evidence(); const facts = proof.groups[0].effective!;
    facts.recording_profile = null; facts.allowed_lateness_seconds = null; facts.members[0].source_value_type = null;
    const active = show(proof);
    expect(field(active, 'Profile')).toHaveTextContent('Unavailable'); expect(field(active, 'Allowed lateness')).toHaveTextContent('Unavailable');
    expect(field(active, 'Sources and types')).toHaveTextContent('Unavailable');
    expect(field(active, 'Interval')).toHaveTextContent('15 seconds'); expect(field(active, 'Interval')).toHaveTextContent('Normalized interval');
    expect(active).not.toHaveTextContent(/Float64|30 seconds|10 seconds|Instantaneous/);
  });
  it('preserves a proven zero lateness instead of presenting it as missing', () => {
    const proof = evidence(); proof.groups[0].effective!.allowed_lateness_seconds = 0;
    expect(field(show(proof), 'Allowed lateness')).toHaveTextContent('0 seconds');
  });
  it('shows an unknown saved interval as unavailable, not zero or an inferred profile default', () => {
    const proof = evidence(); proof.groups[0].saved!.recording_profile = null; proof.groups[0].saved!.interval_seconds = null; show(proof);
    const draft = screen.getByRole('region', { name: labels.saved });
    expect(field(draft, 'Interval')).toHaveTextContent(/^Unavailable$/);
    expect(field(draft, 'Interval')).not.toHaveTextContent(/0 seconds|Normalized/);
  });
  it('does not invent a schema or alias for the safe destination', () => {
    const proof = evidence(); proof.groups[0].effective!.destination.schema = null; proof.groups[0].effective!.destination.alias = null;
    expect(field(show(proof), 'Destination table')).toHaveTextContent(/^old_readings$/);
  });
  it('keeps a shared group’s complete source list without relabeling it as selected A’s exclusive group', () => {
    const proof = evidence(); proof.groups[0].effective!.members.push({ device_id: 'B', circuit_id: 'C2', point_id: 'P2', tag_id: 'T2', column: 'other_device', source_value_type: 'bool' });
    const active = show(proof); const sources = field(active, 'Sources and types');
    expect(sources?.querySelectorAll('li')).toHaveLength(2);
    expect(sources).toHaveTextContent('Int16'); expect(sources).toHaveTextContent('Boolean');
    expect(sources).not.toHaveTextContent(/old_value|other_device/);
    expect(field(active, 'Ownership context')).toHaveTextContent(labels.groupScope);
  });
  it('keeps selected Tag identities distinguishable when a long layout uses the same value column', () => {
    const proof = evidence(); const facts = proof.groups[0].effective!; facts.layout = 'managed_long_v1';
    facts.members = [{ device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1', column: 'value', source_value_type: 'int16' },
      { device_id: 'A', circuit_id: 'C1', point_id: 'P2', tag_id: 'T2', column: 'value', source_value_type: 'int16' }];
    const sources = field(show(proof), 'Sources and types');
    expect(sources).toHaveTextContent('Unavailable 1'); expect(sources).toHaveTextContent('Unavailable 2');
    expect(sources?.querySelectorAll('li')).toHaveLength(2);
  });
  it('distinguishes an empty known source list from absent version evidence', () => {
    const proof = evidence(); proof.groups[0].effective!.members = [];
    expect(field(show(proof), 'Sources and types')).toHaveTextContent('No sources in this version');
  });
  it('keeps available effective facts when safe saved-summary facts are missing, without reading the private editable record', () => {
    const proof = evidence(); proof.groups[0].saved = null; show(proof);
    expect(screen.getByRole('region', { name: 'Saved configuration' })).toHaveTextContent('Unavailable');
    expect(screen.getByRole('region', { name: 'Effective configuration' })).toHaveTextContent('15 seconds');
    expect(document.body).not.toHaveTextContent(/30 seconds|new_readings|PRIVATE-NAME|database.sqlite/);
  });
  it('keeps safe saved facts but labels missing effective snapshot proof unavailable', () => {
    const proof = evidence(); proof.groups[0].proof = 'unavailable'; show(proof);
    expect(screen.getByRole('region', { name: 'Saved configuration' })).toHaveTextContent('30 seconds');
    expect(screen.getByRole('region', { name: 'Effective configuration' })).toHaveTextContent('Effective facts unavailable');
    expect(document.body).not.toHaveTextContent(/15 seconds|old_readings/);
  });
  it.each(['loading', 'error', 'unknown', 'unavailable'] as const)('renders safe %s degradation without substituting draft data', (state) => {
    const current = recordingWorkflowView({ scope, savedGroup: saved(), runtimeEvidence: { state } });
    render(<RecordingRuleSummary view={current} scopeLabel="Device A" labels={labels} />);
    expect(screen.getAllByRole('status').length).toBeGreaterThan(0);
    expect(document.body).not.toHaveTextContent(/new_readings|old_readings|30 seconds|15 seconds/);
  });
  it('shows a safe unavailable view without exposing a supplied foreign detail slot', () => {
    render(<RecordingRuleSummary view={{ state: 'unavailable', reason: 'scope_mismatch' }} scopeLabel="Unavailable selection" labels={labels}
      details={{ ...detailSlot(), content: <p>FOREIGN PRIVATE DETAILS</p> }} />);
    expect(screen.getByRole('status')).toHaveTextContent('Scope unavailable');
    expect(screen.queryByRole('button', { name: labels.details })).not.toBeInTheDocument(); expect(document.body).not.toHaveTextContent('FOREIGN PRIVATE DETAILS');
  });
  it('never forwards raw backend diagnostic fields or guessed labels from rejected runtime evidence', () => {
    const proof = evidence(); Object.assign(proof.groups[0].saved!.destination, { database: '/private/SECRET.sqlite', raw_error: 'SECRET DRIVER FAILURE' });
    show(proof); expect(document.body).not.toHaveTextContent(/SECRET|PRIVATE-NAME|new_readings|database.sqlite/);
    expect(screen.getAllByRole('status').some(node => node.textContent?.includes('Evidence invalid'))).toBe(true);
  });
  it('distinguishes a requested scheduled revision from actual effective facts', () => {
    const proof = evidence(); proof.groups[0].applied_revision = 'R3'; proof.groups[0].applied_effective_at = '2026-10-03T00:00:00Z';
    const selected = { ...scope, appliedRevision: 'R3' };
    render(<RecordingRuleSummary view={view(proof, selected, { ...saved(), applied_revision: 'R3' })} scopeLabel="Device A" labels={labels} />);
    expect(screen.getByRole('region', { name: labels.effective })).toHaveTextContent('15 seconds');
    expect(screen.getByText(labels.requestedNotEffective)).toBeInTheDocument();
    expect(document.body).not.toHaveTextContent(/running R3|R3 active/i);
  });
  it.each(['ended', 'disabled'] as const)('presents retained %s configuration as immutable history, not active intake or start permission', (kind) => {
    const proof = evidence(); if (kind === 'ended') proof.groups[0].effective!.ends_at = proof.observed_at;
    else proof.groups[0].effective!.enabled = false;
    const active = show(proof); expect(active).toHaveTextContent('15 seconds');
    expect(screen.getAllByRole('button').map(button => button.textContent)).toEqual([labels.details]);
    expect(active).not.toHaveTextContent(/running|start allowed|SQL committed|readback verified|cleanup completed/i);
  });
  it.each(['layout', 'storage_mode', 'reduction', 'incomplete_policy'] as const)('does not default an unknown saved %s enum to a supported value', (key) => {
    const proof = evidence(); proof.groups[0].saved![key] = 'unavailable'; show(proof);
    const title = { layout: 'Layout', storage_mode: 'Storage mode', reduction: 'Reduction', incomplete_policy: 'Incomplete rows' }[key];
    expect(field(screen.getByRole('region', { name: labels.saved }), title)).toHaveTextContent('Unavailable');
  });
});

describe('same-scope read-only detail slot', () => {
  it.each(['close', 'Escape'])('keeps B1/R1 and scope accessible, then restores focus with %s', async (dismiss) => {
    render(<RecordingRuleSummary view={view()} scopeLabel="Device A / Circuit C1" labels={labels} details={detailSlot()} />);
    const opener = screen.getByRole('button', { name: labels.details }); opener.focus(); fireEvent.click(opener);
    const dialog = screen.getByRole('dialog', { name: labels.detailsTitle });
    expect(dialog).toHaveTextContent('Source binding B1'); expect(dialog).toHaveTextContent('applied revision R1'); expect(dialog).toHaveTextContent('counter epoch E1');
    expect(field(dialog, 'Workspace')).toHaveTextContent('W1'); expect(field(dialog, 'Group')).toHaveTextContent('G1');
    expect(field(dialog, 'Saved revision')).toHaveTextContent('R2'); expect(field(dialog, 'Effective revision')).toHaveTextContent('R1');
    expect(field(dialog, 'Effective from')).toHaveTextContent('2026-10-01T00:00:00Z');
    if (dismiss === 'close') fireEvent.click(within(dialog).getByRole('button', { name: labels.close }));
    else fireEvent.keyDown(dialog, { key: 'Escape' });
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    await waitFor(() => expect(opener).toHaveFocus());
  });
  it.each(['workspaceId', 'deviceId', 'circuitId', 'groupId', 'workspaceRevision', 'groupRevision', 'appliedRevision', 'effectiveRevision', 'asOf'] as const)(
    'rejects foreign/stale detail content bound by %s', (key) => {
      const slot = detailSlot(); const invalid = key === 'effectiveRevision' || key === 'asOf' ? { ...slot, [key]: 'foreign' }
        : { ...slot, scope: { ...scope, [key]: 'foreign' } };
      render(<RecordingRuleSummary view={view()} scopeLabel="Device A" labels={labels} details={invalid} />);
      fireEvent.click(screen.getByRole('button', { name: labels.details }));
      expect(screen.getByRole('dialog')).toHaveTextContent(labels.detailsUnavailable);
      expect(document.body).not.toHaveTextContent('Source binding B1');
    });
  it('shows explicit unavailable details rather than inventing policies or counter epochs when no slot exists', () => {
    show(); fireEvent.click(screen.getByRole('button', { name: labels.details }));
    expect(screen.getByRole('dialog')).toHaveTextContent(labels.detailsUnavailable); expect(document.body).not.toHaveTextContent(/epoch E1|Source binding B1/);
  });
  it.each([null, undefined, false, ''])('does not silently display an empty supplied detail slot: %j', (content) => {
    render(<RecordingRuleSummary view={view()} scopeLabel="Device A" labels={labels} details={{ ...detailSlot(), content }} />);
    fireEvent.click(screen.getByRole('button', { name: labels.details }));
    expect(screen.getByRole('dialog')).toHaveTextContent(labels.detailsUnavailable);
  });
  it('labels missing requested/circuit evidence without inferring it from historical members and preserves a proven end time', () => {
    const proof = evidence(); proof.groups[0].applied_revision = null; proof.groups[0].effective!.ends_at = '2026-10-02T00:00:00Z';
    const selected = { ...scope, circuitId: undefined, appliedRevision: '' };
    render(<RecordingRuleSummary view={view(proof, selected, { ...saved(), applied_revision: '' })} scopeLabel="Device A" labels={labels} />);
    fireEvent.click(screen.getByRole('button', { name: labels.details })); const dialog = screen.getByRole('dialog');
    expect(dialog).toHaveClass('gw-modal'); expect(dialog).toHaveAttribute('data-size', 'form');
    expect(field(dialog, 'Selected circuit')).toHaveTextContent('Unavailable'); expect(field(dialog, 'Requested revision')).toHaveTextContent('Unavailable');
    expect(field(dialog, 'Ends at')).toHaveTextContent('2026-10-02T00:00:00Z'); expect(field(dialog, 'Effective revision')).toHaveTextContent('R1');
  });
  it('removes old detail content when effective proof becomes unavailable while the selected scope stays the same', () => {
    const rendered = render(<RecordingRuleSummary view={view()} scopeLabel="Device A" labels={labels} details={detailSlot()} />);
    fireEvent.click(screen.getByRole('button', { name: labels.details })); expect(screen.getByRole('dialog')).toHaveTextContent('Source binding B1');
    const proof = evidence(); proof.groups[0].as_of = null;
    rendered.rerender(<RecordingRuleSummary view={view(proof)} scopeLabel="Device A" labels={labels} details={detailSlot()} />);
    expect(screen.getByRole('dialog')).toHaveTextContent(labels.detailsUnavailable);
    expect(screen.getByRole('dialog')).not.toHaveTextContent('Source binding B1');
    expect(field(screen.getByRole('dialog'), 'Effective revision')).toHaveTextContent('Unavailable');
  });
  it('closes on a viewed-scope change and does not reopen or leak old content when returning', () => {
    const rendered = render(<RecordingRuleSummary view={view()} scopeLabel="Device A" labels={labels} details={detailSlot()} />);
    fireEvent.click(screen.getByRole('button', { name: labels.details })); expect(screen.getByRole('dialog')).toBeInTheDocument();
    const other = evidence(); other.device_id = 'B';
    for (const facts of [other.groups[0].saved!, other.groups[0].effective!]) for (const member of facts.members) member.device_id = 'B';
    const group = saved(); group.basic_managed_device_id = 'B'; group.members[0].device_id = 'B';
    rendered.rerender(<RecordingRuleSummary view={view(other, { ...scope, deviceId: 'B' }, group)} scopeLabel="Device B" labels={labels} details={detailSlot()} />);
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); expect(document.body).not.toHaveTextContent('Source binding B1');
    rendered.rerender(<RecordingRuleSummary view={view()} scopeLabel="Device A" labels={labels} details={detailSlot()} />);
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); expect(screen.getByRole('region', { name: labels.saved })).toHaveTextContent('30 seconds');
  });
  it('never saves/applies/starts or writes storage on render, inspect, close, rerender or unmount', async () => {
    const forbidden = () => { throw new Error('Mutation from read-only summary'); };
    const spies = [vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'create').mockImplementation(forbidden),
      vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'update').mockImplementation(forbidden), vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'apply').mockImplementation(forbidden),
      vi.spyOn(studioV2RecordingStartAPI, 'start').mockImplementation(forbidden), vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => undefined),
      vi.spyOn(globalThis, 'fetch').mockImplementation(forbidden)];
    const current = view(); const before = structuredClone(current);
    const rendered = render(<RecordingRuleSummary view={current} scopeLabel="Device A" labels={labels} details={detailSlot()} />);
    fireEvent.click(screen.getByRole('button', { name: labels.details })); fireEvent.click(screen.getByRole('button', { name: labels.close }));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    rendered.rerender(<RecordingRuleSummary view={current} scopeLabel="Device A" labels={labels} details={detailSlot()} />); rendered.unmount();
    for (const spy of spies) expect(spy).not.toHaveBeenCalled(); expect(current).toEqual(before);
    expectTypeOf<Extract<keyof RecordingRuleSummaryProps, 'onSave' | 'onApply' | 'onStart' | 'readCurrent'>>().toEqualTypeOf<never>();
  });
});
