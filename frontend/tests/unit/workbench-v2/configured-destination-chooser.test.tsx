import type { WriteGroupCalculationCapabilities } from '../../../src/types/studioV2WriteGroupCalculation';
import { studioV2RecordingStartAPI } from '../../../src/services/studioV2RecordingStart';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { ConfiguredDestinationChooser, type ConfiguredDestinationObservation } from '../../../src/features/datalink/workbench-v2/recording-editor/ConfiguredDestinationChooser';
import { GroupEditor } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { WriteGroupSection } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { groupToDraft } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
vi.mock('../../../src/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), get: vi.fn() } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
 readiness: vi.fn().mockResolvedValue(null), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));
vi.mock('../../../src/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI), state = batchPanelState(), connector = state.db.connector;
const original = longProfileGroup({ id: 'long-1', circuitId: 'c1', circuitName: 'Original', interval: 15, points: [0, 1] });
const group = { ...original, destination: { ...original.destination, database: connector.database, table_schema: connector.schema } };
const other = { ...group, id: 'long-2', name: 'Other recording', revision: 'other-rev' };
const data = { workspace_id: 'ws-1', workspace_revision: 'w1', groups: [group, other] };
const context = { workspaceId: 'ws-1', connector, draft: groupToDraft(group) };
const client = () => new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
beforeEach(() => vi.clearAllMocks());
function noWrites() { expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled(); expect(api.testWrite).not.toHaveBeenCalled(); expect(api.create).not.toHaveBeenCalled(); expect(api.update).not.toHaveBeenCalled(); expect(api.apply).not.toHaveBeenCalled(); expect(api.schemaApply).not.toHaveBeenCalled(); }
it.each(['loading', 'error', 'stale', 'empty'] as const)('shows honest %s configuration state without implying physical table absence', status => {
 const observation: ConfiguredDestinationObservation = { data: status === 'empty' ? { ...data, groups: [] } : status === 'stale' ? { ...data, workspace_revision: 'old' } : data, loading: status === 'loading', failed: status === 'error', retry: vi.fn() };
 render(<ConfiguredDestinationChooser observation={observation} context={context} workspaceRevision="w1" disabled={false} onSelect={vi.fn()} onReset={vi.fn()} />);
 expect(screen.getByText(`destinations.${status}`)).toBeInTheDocument();
 expect(screen.getByText('destinations.configuredOnly')).toBeInTheDocument();
 for (const button of screen.queryAllByRole('button', { name: 'destinations.select', hidden: true })) expect(button).toBeDisabled(); noWrites();
});
it('shows same table labels separately, explains foreign and wide owners, and selection only emits intent', () => {
 const foreign = { ...other, id: 'foreign', destination: { ...other.destination, database: 'other-db' } };
 const wide = { ...other, id: 'wide', layout_version: 'managed_wide_profile_v1' as const };
 const select = vi.fn();
 render(<ConfiguredDestinationChooser observation={{ data: { ...data, groups: [group, foreign, wide] }, loading: false, failed: false, retry: vi.fn() }} context={context} workspaceRevision="w1" disabled={false} onSelect={select} onReset={vi.fn()} />);
 fireEvent.click(screen.getByText('destinations.title'));
 expect(screen.getAllByText(group.destination.table_name)).toHaveLength(3);
 expect(within(screen.getByTestId('configured-destination-foreign')).getByText('destinations.reason.destination')).toBeInTheDocument();
 expect(within(screen.getByTestId('configured-destination-wide')).getByText('destinations.reason.wideOwner')).toBeInTheDocument();
 expect(within(screen.getByTestId('configured-destination-wide')).getByRole('link')).toHaveAttribute('href', expect.stringContaining('group_id=wide'));
 fireEvent.click(within(screen.getByTestId('configured-destination-long-1')).getByRole('button', { name: 'destinations.select' }));
 expect(select).toHaveBeenCalledExactlyOnceWith(group); noWrites();
 const row = within(screen.getByTestId('configured-destination-long-1'));
 expect(row.getByRole('heading', { name: group.name })).toBeInTheDocument();
 expect(row.getByRole('button', { name: 'destinations.select' })).toHaveClass('inline-flex', 'h-9', 'border-[var(--gw-control)]', 'bg-[var(--gw-surface)]', 'focus-visible:outline-[var(--gw-focus)]');
 expect(row.getByRole('link')).toHaveClass('inline-flex', 'h-9', 'border-[var(--gw-control)]', 'bg-[var(--gw-surface)]', 'focus-visible:outline-[var(--gw-focus)]');
 expect(row.getByText('destinations.technical').closest('details')).not.toHaveAttribute('open');
 expect(row.getByText(group.revision)).toBeInTheDocument();
 expect(screen.getByText('destinations.reason.wideOwner').closest('details')).toBe(screen.getByTestId('configured-destinations'));

});
it('the original editor refuses Save when fresh readback changes the selected reference and copies no proof', async () => {
 render(<QueryClientProvider client={client()}><GroupEditor group={group} workspaceId="ws-1" workspaceRevision="w1" connector={connector} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false}
  configuredDestinations={{ data, loading: false, failed: false, retry: vi.fn() }} onClose={vi.fn()} onSaved={vi.fn()} onReload={vi.fn()}
  fetchLatest={async () => ({ ...data, groups: [group, { ...other, revision: 'changed-after-selection' }] })} /></QueryClientProvider>);
 fireEvent.click(screen.getByText('destinations.title'));
 fireEvent.click(within(screen.getByTestId('configured-destination-long-2')).getByRole('button', { name: 'destinations.select' }));
 expect(screen.getByTestId('group-name')).toHaveValue(group.name);
 fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'My intent' } });
 expect(screen.getByTestId('group-save')).not.toBeDisabled(); fireEvent.click(screen.getByTestId('group-save'));
 await waitFor(() => expect(screen.getByTestId('group-save-error')).toBeInTheDocument()); noWrites();
});
it('opening another wide owner uses the same existing draft guard and never mutates either recording', async () => {
 const wide = { ...other, id: 'wide', layout_version: 'managed_wide_profile_v1' as const };
 api.list.mockResolvedValue({ ...data, groups: [group, wide] });
 render(<QueryClientProvider client={client()}><RecordingDraftNavigation><WriteGroupSection state={state} workspaceId="ws-1" readonly={false} /></RecordingDraftNavigation></QueryClientProvider>);
 fireEvent.click(await screen.findByTestId('group-open-long-1'));
 fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Keep me' } });
 fireEvent.click(screen.getByText('destinations.title'));
 fireEvent.click(within(screen.getByTestId('configured-destination-wide')).getByRole('button', { name: 'destinations.editOriginal' }));
 expect(await screen.findByRole('alertdialog')).toBeInTheDocument(); fireEvent.click(screen.getByText('navigation.keep'));
 expect(screen.getByTestId('group-name')).toHaveValue('Keep me');
 fireEvent.click(within(screen.getByTestId('configured-destination-wide')).getByRole('button', { name: 'destinations.editOriginal' }));
 fireEvent.click(await screen.findByText('navigation.discard'));
 expect(await screen.findByTestId('group-name')).toHaveValue(wide.name); noWrites();
});
it('an explicit Save after selection updates only this recording and carries no borrowed schema or group identity', async () => {
 const borrowed = { ...other, destination: { ...other.destination, schema_revision: 'foreign-proof', schema_digest: 'foreign-digest' } };
 const list = { ...data, groups: [group, borrowed] };
 api.update.mockResolvedValue({ workspace_revision: 'w2', group: { ...group, revision: 'new-own-revision' } });
 render(<QueryClientProvider client={client()}><GroupEditor group={group} workspaceId="ws-1" workspaceRevision="w1" connector={connector} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false}
  configuredDestinations={{ data: list, loading: false, failed: false, retry: vi.fn() }} onClose={vi.fn()} onSaved={vi.fn()} onReload={vi.fn()} fetchLatest={async () => list} /></QueryClientProvider>);
 fireEvent.click(screen.getByText('destinations.title'));
 fireEvent.click(within(screen.getByTestId('configured-destination-long-2')).getByRole('button', { name: 'destinations.select' }));
 noWrites(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'My own recording' } }); fireEvent.click(screen.getByTestId('group-save'));
 await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
 const [id, request] = api.update.mock.calls[0]; expect(id).toBe(group.id);
 expect(request.expected_group_revision).toBe(group.revision); expect(request.group.name).toBe('My own recording');
 expect(request.group).not.toHaveProperty('id'); expect(request.group).not.toHaveProperty('revision');
 expect(request.group.destination).not.toHaveProperty('schema_revision'); expect(request.group.destination).not.toHaveProperty('schema_digest');
 expect(api.create).not.toHaveBeenCalled(); expect(api.apply).not.toHaveBeenCalled(); expect(api.schemaApply).not.toHaveBeenCalled();
});

it('changing a selected long layout explains incompatibility rather than claiming a stale workspace', () => {
 const capabilities: WriteGroupCalculationCapabilities = { version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false, reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true }, schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [] };
 render(<QueryClientProvider client={client()}><GroupEditor group={undefined} calculationCapabilities={capabilities} workspaceId="ws-1" workspaceRevision="w1" connector={connector} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false}
  configuredDestinations={{ data, loading: false, failed: false, retry: vi.fn() }} onClose={vi.fn()} onSaved={vi.fn()} onReload={vi.fn()} fetchLatest={async () => data} /></QueryClientProvider>);
 fireEvent.click(screen.getByText('calculation.chooseMixed'));
 fireEvent.change(screen.getByLabelText('calculation.layoutLabel'), { target: { value: 'managed_long_v1' } });
 fireEvent.click(screen.getByText('destinations.title'));
 fireEvent.click(within(screen.getByTestId('configured-destination-long-1')).getByRole('button', { name: 'destinations.select' }));
 fireEvent.change(screen.getByLabelText('calculation.layoutLabel'), { target: { value: 'managed_wide_profile_v1' } });
 expect(screen.queryByText('destinations.stale')).not.toBeInTheDocument();
 expect(screen.getAllByText('destinations.reason.layout').length).toBeGreaterThan(1);
 expect(screen.getByTestId('group-save')).toBeDisabled(); noWrites();
});

it('a styled disabled original-record anchor prevents keyboard-origin activation without nesting a button', () => {
 render(<ConfiguredDestinationChooser observation={{ data: { ...data, groups: [group] }, loading: false, failed: false, retry: vi.fn() }} context={context} workspaceRevision="w1" disabled onSelect={vi.fn()} onReset={vi.fn()} />);
 fireEvent.click(screen.getByText('destinations.title'));
 const link = screen.getByRole('link', { name: 'destinations.editOriginal' });
 expect(link).toHaveAttribute('aria-disabled', 'true'); expect(link).toHaveAttribute('tabindex', '-1');
 expect(link.querySelector('button')).toBeNull(); expect(link.closest('button')).toBeNull();
 // Native Enter activation dispatches a click with zero pointer detail; the route guard remains intact.
 expect(fireEvent.click(link, { detail: 0 })).toBe(false); noWrites();
});
