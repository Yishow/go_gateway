import { RecordingEditorSummary } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingEditorSummary';
import { newDraft } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { DbColumn } from '@/features/datalink/workbench-v2/state/dbSchemas';
import { dbTargetAPI } from '@/services/datalink';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { GroupEditor, type GroupEditorProps } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import type { WriteGroupCalculationCapabilities } from '../../../src/types/studioV2WriteGroupCalculation';
const observedColumns = vi.hoisted((): { value?: DbColumn[] } => ({}));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: observedColumns.value ?? TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  readiness: vi.fn().mockResolvedValue(null), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
} }));
vi.mock('../../../src/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn().mockResolvedValue([]) } }));
const capabilities: WriteGroupCalculationCapabilities = { version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false, reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true }, schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [] };
function mount(overrides: Partial<GroupEditorProps> = {}) {
  const state = savedGroupState();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const props: GroupEditorProps = { presentation: 'simple', group: undefined, calculationCapabilities: capabilities, workspaceId: 'ws-1', workspaceRevision: 'db-setup-9', connector: state.db.connector, candidates: buildGroupCandidates(state).candidates, excluded: [], readonly: false, onSaved: vi.fn(), onClose: vi.fn(), onReload: vi.fn(), fetchLatest: vi.fn().mockResolvedValue({ workspace_revision: 'db-fresh-10', groups: [] }), ...overrides };
  const view = (value: GroupEditorProps) => <QueryClientProvider client={client}><GroupEditor {...value} /></QueryClientProvider>;
  const result = render(view(props));
  return { ...result, props, update: (patch: Partial<GroupEditorProps>) => { Object.assign(props, patch); result.rerender(view(props)); } };
}
beforeEach(() => {
  vi.clearAllMocks(); observedColumns.value = undefined;
  vi.spyOn(dbTargetAPI, 'listConnectors').mockResolvedValue([{ id: 'conn-1', identity_revision: 'crev-1', name: 'Line DB', kind: 'sqlite', enabled: true, status: 'ready', connection_config: { path: savedGroupState().db.connector.database }, last_check_error: '', default_write_interval_seconds: 5, created_at: '', updated_at: '' }]);
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mockImplementation(async request => {
    const canonical = savedGroup();
    const members = request.group.members.map(member => {
      const source = canonical.members.find(value => value.device_id === member.device_id && value.point_id === member.point_id && value.tag_id === member.tag_id);
      if (!source) throw new Error('Fixture source identity is unavailable');
      return { ...member, source_revision: source.source_revision, mapping_revision: source.mapping_revision };
    });
    return { workspace_revision: 'created-1', group: savedGroup({ ...request.group, members,
      destination: { ...canonical.destination, ...request.group.destination }, id: 'new-group', revision: 'new-1' }) };
  });
});
it('starts an explicit simple recording in per-point mode and saves through the original single owner', async () => {
  mount();
  expect(screen.queryByRole('button', { name: 'calculation.chooseMixed' })).not.toBeInTheDocument();
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Simple intent' } });
  fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
  const row = await screen.findByTestId('member-calculation-pt-0');
  expect(within(row).getByLabelText('calculation.interval')).toHaveValue('15');
  fireEvent.change(within(row).getByLabelText('calculation.interval'), { target: { value: '15' } });
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled());
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce());
  const wire = vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mock.calls[0][0];
  expect(wire.expected_workspace_revision).toBe('db-fresh-10');
  expect(wire.group).toMatchObject({ name: 'Simple intent', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, members: [{ calculation: 'latest', cadence: { interval: 15, unit: 'second' } }] });
  expect(wire.group.row_policy).not.toHaveProperty('interval_seconds');
  expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});
it('keeps local name and source intent while destination and calculation permission are unavailable', async () => {
  const state = savedGroupState();
  const view = mount({ calculationCapabilities: undefined, connector: { ...state.db.connector, persisted: false, save_state: 'idle', connector_id: undefined, identity_revision: undefined } });
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep my draft' } });
  fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
  expect(screen.getByTestId('group-save')).toBeDisabled();
  expect(screen.getByTestId('group-name')).toHaveValue('Keep my draft');
  view.update({ calculationCapabilities: capabilities, connector: state.db.connector });
  expect(screen.getByTestId('group-name')).toHaveValue('Keep my draft');
  const row = await screen.findByTestId('member-calculation-pt-0');
  expect(within(row).getByLabelText('calculation.interval')).toHaveValue('15');
  fireEvent.change(within(row).getByLabelText('calculation.interval'), { target: { value: '1e' } });
  view.update({ calculationCapabilities: undefined });
  view.update({ calculationCapabilities: capabilities });
  expect(within(row).getByLabelText('calculation.interval')).toHaveValue('1e');
  expect(screen.getByTestId('group-save')).toBeDisabled();
  expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
});
it('orders visible source and per-point controls before the destination without hiding common policy or safety notes', async () => {
  mount({ destinationContent: <button type="button">Existing database owner</button> });
  fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
  const sources = screen.getByTestId('recording-editor-sources');
  const destination = screen.getByTestId('recording-editor-destination');
  expect(sources.compareDocumentPosition(destination) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
  expect(within(sources).getByTestId('member-calculation-pt-0')).toBeVisible();
  expect(within(destination).getByRole('combobox', { name: 'newDestination.label' })).toBeVisible();
  expect(within(destination).queryByRole('button', { name: 'Existing database owner' })).not.toBeInTheDocument();
  expect(within(destination).getByLabelText('calculation.layoutLabel')).toBeVisible();
  expect(screen.getByTestId('group-lateness')).toBeVisible();
  // Only risks that apply stay visible; the method explanation lives in the section's help button.
  expect(screen.queryByText('calculation.firstBaseline')).not.toBeInTheDocument();
  expect(screen.getByText('simple.requiredRisk')).toBeVisible();
  const method = screen.getByTestId('recording-editor-method-note');
  expect(method).not.toHaveTextContent('calculation.dueUnion');
  expect(within(method).getByTestId('recording-method-groups')).toHaveTextContent('simple.summary.points');
  expect(method.querySelector('[data-parameter-help="record.method"]')).not.toBeNull();
  expect(method.querySelector('details')).toBeNull();
  expect(screen.getByTestId('member-calculation-pt-0').closest('td')).toHaveAttribute('data-testid', 'group-member-calculation-pt-0');
});
it('keeps one repair action inside the padded editor before preparation and in the readonly fallback', () => {
  const lookup = vi.fn();
  const repairContent = <button type="button" data-testid="repair-slot" onClick={lookup}>Original operation lookup</button>;
  const view = mount({ repairContent });
  const action = screen.getByTestId('repair-slot'), editor = screen.getByTestId('group-editor');
  expect(editor).toHaveClass('recording-editor-simple'); expect(editor).toContainElement(action);
  expect(action.compareDocumentPosition(screen.getByTestId('recording-preparation-actions')) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
  expect(screen.getAllByTestId('repair-slot')).toHaveLength(1); expect(lookup).not.toHaveBeenCalled();
  action.focus(); expect(action).toHaveFocus(); fireEvent.click(action); expect(lookup).toHaveBeenCalledOnce();
  expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  view.unmount();
  mount({ repairContent, calculationCapabilities: undefined, group: savedGroup({ recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' } }) });
  expect(screen.getByTestId('mixed-group-editor-readonly')).toContainElement(screen.getByTestId('repair-slot'));
  expect(screen.getAllByTestId('repair-slot')).toHaveLength(1); expect(screen.queryByTestId('recording-preparation-actions')).toBeNull();
});
it('preserves saved legacy policy instead of silently converting it into per-point mode', () => {
  mount({ group: savedGroup({ layout_version: 'managed_wide_profile_v1', recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' } }) });
  expect(screen.queryByTestId('member-calculation-pt-0')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-profile-interval')).toHaveTextContent('15');
  expect(screen.getByText('calculation.legacyUnchanged')).toBeVisible();
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('keeps the original guided explicit mode choice when no simple presentation is requested', () => {
  mount({ presentation: undefined });
  expect(screen.getByRole('button', { name: 'calculation.chooseMixed' })).toBeVisible();
  expect(screen.queryByTestId('recording-editor-sources')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-editor-close')).toBeVisible();
});
it('does not give a readonly simple presentation permission to change or save', () => {
  mount({ readonly: true });
  expect(screen.getByTestId('group-name')).toBeDisabled();
  expect(screen.getByTestId('group-member-include-pt-0')).toBeDisabled();
  expect(screen.getByTestId('group-save')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-save'));
  expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
});

it('explains server-assigned columns for only a new managed draft without asking for manual mapping or claiming SQL creation', async () => {
  mount();
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Server-owned columns' } });
  fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
  expect(screen.queryByTestId('group-member-column-pt-0')).not.toBeInTheDocument();
  expect(screen.queryByTestId('group-column-proposal')).not.toBeInTheDocument();
  expect(screen.getByText('simple.columnsAssignedOnSave')).toBeVisible();
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled());
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce());
  const wire = vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mock.calls[0][0].group;
  expect(wire.members[0].target_column).toBe('');
  expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});
it.each(['custom', 'managed'] as const)('does not replace saved %s column repair with the server-assigned columns presentation', storage => {
  const group = savedGroup();
  mount({ group: { ...group, destination: { ...group.destination, storage_strategy: storage }, members: group.members.map(member => ({ ...member, target_column: '' })) } });
  expect(screen.getByTestId('group-member-column-pt-0')).toBeVisible();
  expect(screen.getByTestId('group-column-proposal')).toBeVisible();
  expect(screen.queryByText('simple.columnsAssignedOnSave')).not.toBeInTheDocument();
});


it('keeps the empty simple picker focused and derives a read-only summary without extra save owners', () => {
  const source: NonNullable<GroupEditorProps['deviceSource']> = { observations: () => [], inspect: vi.fn().mockResolvedValue(undefined),
    submit: vi.fn(async request => ({ request, status: 'blocked' as const })) };
  mount({ candidates: [], deviceSource: source });
  expect(screen.getAllByTestId('recording-points-open')).toHaveLength(1);
  expect(screen.queryByTestId('group-member-search')).not.toBeInTheDocument();
  expect(screen.queryByTestId('group-member-problems-only')).not.toBeInTheDocument();
  expect(within(screen.getByTestId('recording-selected-members')).queryByRole('table')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-member-empty')).toHaveTextContent('simple.noSelectedPoints');
  expect(screen.getByTestId('recording-summary-name')).toHaveTextContent('simple.unset');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'My unsaved name' } });
  expect(screen.getByTestId('recording-summary-name')).toHaveTextContent('My unsaved name');
  // The summary rail is always open: name and table are plain rows, not a disclosure.
  expect(screen.getByTestId('recording-summary-name').closest('details')).toBeNull();
  expect(screen.getByTestId('recording-summary-name')).toBeVisible();
  expect(screen.getByTestId('recording-summary-table')).toBeVisible();
  expect(screen.getByText('compact.unsavedSettings')).toBeVisible();
  expect(screen.getAllByTestId('recording-preparation-actions')).toHaveLength(1);
  expect(screen.getAllByTestId('group-save')).toHaveLength(1);
  expect(screen.getByTestId('group-editor')).toHaveClass('recording-editor-simple');
  expect(screen.queryByTestId('group-editor-close')).not.toBeInTheDocument();
  expect(screen.queryByTestId('recording-calculation-limits')).not.toBeInTheDocument();
  expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  expect(source.submit).not.toHaveBeenCalled();
});
it('locates three sections in the same page without changing history or requiring wizard steps', () => {
  mount();
  const method = screen.getByTestId('recording-editor-method-note');
  const scroll = vi.fn(); Object.defineProperty(method, 'scrollIntoView', { value: scroll, configurable: true });
  const push = vi.spyOn(window.history, 'pushState');
  fireEvent.click(screen.getByRole('button', { name: '2 simple.section.method' }));
  expect(scroll).toHaveBeenCalledOnce(); expect(method).toHaveFocus(); expect(push).not.toHaveBeenCalled();
  expect(screen.getByTestId('recording-editor-sources')).toBeVisible();
  expect(screen.getByTestId('recording-editor-destination')).toBeVisible();
  push.mockRestore();
});
it('keeps relevant counter risks and capability denial visible rather than folding them into help', () => {
  const base = savedGroup();
  const view = mount({ group: { ...base, record_id: '11111111-1111-4111-8111-111111111111', layout_version: 'managed_wide_profile_v1',
    recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, destination: { ...base.destination, storage_strategy: 'managed' },
    row_policy: { allowed_lateness_seconds: 0, incomplete_policy: 'partial' }, members: [{ ...base.members[0], calculation: 'delta',
      cadence: { interval: 60, unit: 'second', alignment: 'utc_epoch' }, source_value_type: 'uint64', value_role: 'counter_raw', measurement_id: 'counter-1' }] } });
  expect(screen.getByText('calculation.firstBaseline')).toBeVisible();
  expect(screen.getByText('simple.requiredRisk')).toBeVisible();
  expect(screen.getByTestId('recording-method-groups')).toHaveTextContent('calculation.delta · 60 calculation.second');
  expect(screen.getByTestId('recording-editor-summary')).toHaveTextContent('60');
  view.update({ calculationCapabilities: undefined });
  expect(screen.getByText('calculation.capability')).toBeVisible();
  expect(screen.getByText('calculation.firstBaseline')).toBeVisible();
  expect(screen.getByTestId('group-save')).toBeDisabled();
});
it('keeps only the outer Back for unsupported simple records while guided retains its original Back', () => {
  const view = mount({ calculationCapabilities: undefined, group: savedGroup({ recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' } }) });
  expect(screen.getByTestId('mixed-group-editor-readonly')).toBeVisible();
  expect(screen.queryByRole('button', { name: 'step4.group.editor.back' })).not.toBeInTheDocument();
  view.update({ presentation: 'guided' });
  expect(screen.getByRole('button', { name: 'step4.group.editor.back' })).toBeVisible();
});

it('shows simple version identifiers in the summary rail as short IDs with full values, keeping status and destination denial visible', () => {
  const state = savedGroupState();
  mount({ group: savedGroup(), scopeDiagnostics: <p>Fresh scoped revision W2</p>, connector: { ...state.db.connector, identity_revision: 'different-connector-revision' } });
  const versions = screen.getByTestId('group-revision-details');
  const values = screen.getByTestId('group-revisions');
  expect(versions.tagName).toBe('SECTION'); expect(versions.querySelector('details')).toBeNull();
  expect(within(versions).getByRole('heading', { name: 'versions.details' })).toBeVisible();
  expect(values).toBeVisible(); expect(values).toHaveTextContent('savedStatuses.ready');
  expect(values).toHaveTextContent('rev-1'); expect(values).toHaveTextContent('crev-1');
  expect(screen.getByTestId('recording-editor-summary')).toContainElement(versions);
  expect(versions.closest('.recording-editor-tools')).toBeNull();
  expect(screen.getByTestId('group-lifecycle-return')).toHaveTextContent(savedGroup().name);
  expect(versions).toContainElement(screen.getByText('Fresh scoped revision W2'));
  expect(screen.getAllByText('versions.details')).toHaveLength(1);
  expect(screen.getByTestId('group-destination-mismatch')).toBeVisible();
});
it('keeps saved guided version identifiers visible without a disclosure wrapper', () => {
  mount({ presentation: 'guided', group: savedGroup() });
  expect(screen.queryByTestId('group-revision-details')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-revisions')).toBeVisible();
  expect(screen.getByTestId('group-revisions')).toHaveTextContent('rev-1');
});

it.each([
  { presentation: 'simple', storage: 'managed', layout: 'managed_long_v1', proposal: false },
  { presentation: 'guided', storage: 'managed', layout: 'managed_long_v1', proposal: false },
  { presentation: 'simple', storage: 'managed', layout: 'managed_wide_profile_v1', proposal: true },
  { presentation: 'simple', storage: 'custom', layout: 'managed_long_v1', proposal: true },
] as const)('keeps typed managed long storage separate from wide column proposals: $presentation $storage $layout', ({ presentation, storage, layout, proposal }) => {
  observedColumns.value = [
    { name: 'value_float', type: 'double precision', nullable: true, primary_key: false },
    { name: 'value_decimal', type: 'numeric', nullable: true, primary_key: false },
    { name: 'bucket_start', type: 'timestamp', nullable: false, primary_key: true },
  ];
  const base = savedGroup();
  const group = savedGroup({ record_id: '11111111-1111-4111-8111-111111111111', layout_version: layout, storage_mode: 'plain',
    recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, destination: { ...base.destination, storage_strategy: storage },
    members: base.members.map((member, index) => ({ ...member, target_column: `v_saved_${index}`, calculation: 'mean',
      cadence: { interval: 5, unit: 'second', alignment: 'utc_epoch' }, source_value_type: 'float64', result_value_type: 'float64' })) });
  const deviceSource: NonNullable<GroupEditorProps['deviceSource']> = { observations: () => [], inspect: vi.fn().mockResolvedValue(undefined),
    submit: vi.fn(async request => ({ request, status: 'blocked' as const })) };
  const before = structuredClone(group); const view = mount({ presentation, group, deviceSource: presentation === 'simple' ? deviceSource : undefined });
  expect(screen.getByTestId('group-editor')).toBeInTheDocument();
  if (proposal) expect(screen.getByTestId('group-column-proposal')).toBeVisible();
  else expect(screen.queryByTestId('group-column-proposal')).not.toBeInTheDocument();
  for (const member of group.members) {
    const row = screen.getByTestId(`group-member-row-${member.point_id}`);
    expect(row).toHaveAttribute('data-problem', proposal ? 'column-missing' : '');
    if (presentation === 'simple' && storage === 'managed') {
      expect(screen.queryByTestId(`group-member-column-${member.point_id}`)).not.toBeInTheDocument();
      expect(within(row).queryByText(member.target_column, { selector: 'code' })).not.toBeInTheDocument();
    } else expect(screen.getByTestId(`group-member-column-${member.point_id}`)).toHaveAttribute('aria-invalid', String(proposal));
    if (proposal) expect(within(row).getByRole('alert')).toHaveTextContent('problem_column_missing');
    else expect(within(row).queryByRole('alert')).not.toBeInTheDocument();
  }
  observedColumns.value = [...observedColumns.value, { name: 'v_saved_0', type: 'text', nullable: true, primary_key: false }];
  view.update({});
  const firstRow = screen.getByTestId('group-member-row-pt-0');
  expect(firstRow).toHaveAttribute('data-problem', proposal ? 'column-incompatible' : '');
  if (presentation === 'simple' && storage === 'managed') {
    expect(screen.queryByTestId('group-member-column-pt-0')).not.toBeInTheDocument();
    expect(within(firstRow).queryByText('v_saved_0', { selector: 'code' })).not.toBeInTheDocument();
  } else expect(screen.getByTestId('group-member-column-pt-0')).toHaveAttribute('aria-invalid', String(proposal));
  if (proposal) expect(within(firstRow).getByRole('alert')).toHaveTextContent('problem_column_incompatible');
  else expect(within(firstRow).queryByRole('alert')).not.toBeInTheDocument();
  if (presentation === 'simple') expect(screen.getByTestId('recording-editor-summary')).not.toHaveTextContent('problem_column');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: '' } });
  expect(screen.getByTestId('group-save')).toBeDisabled();
  expect(screen.queryByText('simple.columnsAssignedOnSave')).not.toBeInTheDocument();
  expect(group).toEqual(before); expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});

it('opens an existing simple recording when insecure HTTP has getRandomValues but no randomUUID', () => {
  const getRandomValues = globalThis.crypto.getRandomValues.bind(globalThis.crypto);
  vi.stubGlobal('crypto', { getRandomValues });
  try {
    const group = savedGroup();
    const before = structuredClone(group);
    const deviceSource = { submit: vi.fn(), inspect: vi.fn(), observations: () => [] };
    const view = mount({ group, deviceSource });
    expect(deviceSource.submit).not.toHaveBeenCalled();
    expect(screen.getByTestId('group-name')).toHaveValue(group.name);
    expect(group).toEqual(before);
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
    expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
    expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
    view.unmount();
  } finally { vi.unstubAllGlobals(); }
});

it('keeps the original saved confirmation owners inside the destination panel and outside save fields', async () => {
  mount({ group: savedGroup() });
  const destination = screen.getByTestId('recording-editor-destination');
  for (const id of ['recording-preparation-actions', 'group-delivery-strip', 'group-lifecycle', 'group-test-write']) {
    const owner = await screen.findByTestId(id);
    expect(destination).toContainElement(owner); expect(owner.closest('fieldset')).toBeNull();
    expect(screen.getAllByTestId(id)).toHaveLength(1);
  }
  expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
  expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});

it('keeps version diagnostics below the everyday summary rather than above source editing', () => {
  mount({ group: savedGroup(), scopeDiagnostics: <p>Exact diagnostic revision</p> });
  expect(screen.getByTestId('recording-editor-summary')).toContainElement(screen.getByTestId('group-revision-details'));
  expect(screen.getByTestId('group-revisions')).toBeVisible();
  expect(screen.getByTestId('recording-editor-summary')).toHaveTextContent('compact.savedSettings');
});
it.each([
  { kind: 'instantaneous' as const, interval: 15, unit: 'second' as const, sample_reduction: 'mean' as const, label: 'calculation.mean' },
  { kind: 'cumulative' as const, interval: 10, unit: 'minute' as const, label: 'simple.summary.cumulative' },
])('summarizes the actual saved whole-group $kind method and interval without converting it', profile => {
  const { label, ...recording_profile } = profile;
  mount({ group: savedGroup({ recording_profile, layout_version: 'managed_wide_profile_v1' }) });
  const summary = screen.getByTestId('recording-editor-summary');
  expect(summary).toHaveTextContent(label); expect(summary).toHaveTextContent(String(profile.interval));
  expect(summary).toHaveTextContent(`calculation.${profile.unit}`);
  if (profile.kind === 'cumulative') expect(summary).not.toHaveTextContent('calculation.mean');
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('shows a known legacy interval without inventing a recording method', () => {
  mount({ group: savedGroup({ recording_profile: undefined, row_policy: { ...savedGroup().row_policy, interval_seconds: 7 } }) });
  const summary = screen.getByTestId('recording-editor-summary');
  expect(summary).toHaveTextContent('7'); expect(summary).toHaveTextContent('simple.summary.methodUnknown');
  expect(summary).not.toHaveTextContent('calculation.latest'); expect(summary).not.toHaveTextContent('calculation.mean');
});

it('does not borrow a legacy interval when an existing profile is invalid', () => {
  const draft = newDraft({ connector_id: 'C', connector_revision: 'C1', table_schema: '', table_name: '' });
  draft.interval_seconds = 42; draft.recording_profile = { kind: 'instantaneous', interval: 0, unit: 'second' };
  render(<RecordingEditorSummary draft={draft} saved dirty={false} />);
  const summary = screen.getByTestId('recording-editor-summary');
  expect(summary).toHaveTextContent('simple.summary.methodUnknown'); expect(summary).toHaveTextContent('simple.summary.intervalUnknown');
  expect(summary).not.toHaveTextContent('42'); expect(summary).not.toHaveTextContent('calculation.latest');
});
