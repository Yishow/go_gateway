import { useRecordingPointOrder } from '../../../src/features/datalink/workbench-v2/recording-editor/useRecordingPointOrder';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { GroupEditor } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { pointFormatAPI } from '../../../src/services/pointFormat';
import type { PointFormatSnapshot } from '../../../src/types/pointFormat';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/services/pointFormat', async importOriginal => ({ ...await importOriginal<typeof import('../../../src/services/pointFormat')>(), pointFormatAPI: { read: vi.fn(), save: vi.fn(), inspect: vi.fn() } }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  readiness: vi.fn().mockResolvedValue(null), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
} }));
function snapshot(pointId: string): PointFormatSnapshot {
  const format = pointId === 'pt-0' ? 'ABCD' : 'CDAB';
  return { workspace_id: 'ws-1', device_id: 'dev-1', point_id: pointId, rule_id: 'rule-1', workspace_revision: 'source-7', workspace_revision_origin: 'persisted', point_revision: `point-${pointId}`, rule_revision: 'rule-2', connection_revision: 'conn-4', editable: true,
    resolution: { resolved: true, format_mode: 'override', configured_format: format, effective_format: format, effective_source: 'override', effective_policy: 'override_scalar', native_compatibility: false, applicable: true, format_version: 1, point_revision: `point-${pointId}`, workspace_revision: 'source-7', rule_revision: 'rule-2', connection_revision: 'conn-4' } };
}
beforeEach(() => { vi.clearAllMocks(); vi.mocked(pointFormatAPI.read).mockImplementation(async scope => snapshot(scope.point_id)); });
describe('recording editor canonical point order consumer', () => {
  it.each(['guided', 'simple'] as const)('reads independent canonical point orders for two selected points in %s without deriving them from mapping datatype or DB revision', async presentation => {
    const state = savedGroupState(), group = savedGroup();
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    const onSaved = vi.fn();
    function Page() {
      const pointOrder = useRecordingPointOrder('ws-1', true);
      return <GroupEditor presentation={presentation} group={group} workspaceId="ws-1" workspaceRevision="DATABASE-SETUP-99"
      connector={state.db.connector} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false}
      deviceSource={{ pointOrder, observations: () => [], submit: vi.fn(), inspect: vi.fn() }} onSaved={onSaved} onClose={vi.fn()} onReload={vi.fn()} fetchLatest={vi.fn()} />;
    }
    render(<QueryClientProvider client={client}><Page /></QueryClientProvider>);
    const row = await screen.findByTestId('recording-point-order-pt-0');
    if (presentation === 'simple') {
      const order = await within(row).findByLabelText('editor.order');
      const details = order.closest('details');
      expect(details).not.toBeNull(); expect(details).not.toHaveAttribute('open'); expect(order).not.toBeVisible();
      fireEvent.click(within(row).getByText('editor.title', { selector: 'summary' }));
      expect(details).toHaveAttribute('open'); expect(order).toBeVisible();
      const controls = order.closest('.recording-point-order-inline__controls');
      expect(controls).not.toBeNull();
      expect(within(row).queryByText('editor.save')).not.toBeInTheDocument();
      expect(controls).toContainElement(within(row).getByText('editor.reset'));
      const version = within(row).getByText(/editor.version/);
      expect(version.closest('details')).toBe(details); expect(version).toBeVisible();
      expect(within(row).getByText('editor.reset')).toBeVisible();
    } else expect(row.querySelector('summary')).toBeVisible();
    expect(await within(screen.getByTestId('recording-point-order-pt-0')).findByTestId('point-order-effective')).toHaveTextContent('ABCD');
    expect(await within(screen.getByTestId('recording-point-order-pt-1')).findByTestId('point-order-effective')).toHaveTextContent('CDAB');
    expect(pointFormatAPI.read).toHaveBeenCalledWith({ workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0' });
    expect(pointFormatAPI.read).toHaveBeenCalledWith({ workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-1' });
    expect(pointFormatAPI.save).not.toHaveBeenCalled(); expect(onSaved).not.toHaveBeenCalled();
  });

  it('connects the managed simple editor rows to the existing point-name and word-order owners', async () => {
    const state = savedGroupState(), group = savedGroup(); group.destination.storage_strategy = 'managed';
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    let session: { id: string; identity: GroupCandidate; displayName: string };
    const name = { prepare: vi.fn(async (identity: GroupCandidate) => ({ status: 'ready' as const, session: session = { id: 'S', identity, displayName: identity.label } })),
      submit: vi.fn(async (_id: string, currentName: string) => ({ status: 'confirmed' as const, session, currentName })), release: vi.fn(), inspect: vi.fn() };
    const onSaved = vi.fn();
    function Page() {
      const pointOrder = useRecordingPointOrder('ws-1', true);
      return <GroupEditor presentation="simple" group={group} workspaceId="ws-1" workspaceRevision="DATABASE-SETUP-99"
        connector={state.db.connector} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false}
        sourceDevices={[{ id: 'dev-1', name: 'PLC A' }]} deviceSource={{ name, pointOrder, observations: () => [], submit: vi.fn(), inspect: vi.fn() }}
        onSaved={onSaved} onClose={vi.fn()} onReload={vi.fn()} fetchLatest={vi.fn()} />;
    }
    render(<QueryClientProvider client={client}><Page /></QueryClientProvider>);
    const row = await screen.findByTestId('group-member-row-pt-0');
    const order = await within(row).findByLabelText('editor.order');
    expect(order).toBeVisible(); expect(order.closest('details')).toBeNull();
    const nameInput = within(row).getByRole('textbox');
    fireEvent.focus(nameInput); await waitFor(() => expect(nameInput).not.toHaveAttribute('readonly'));
    fireEvent.change(nameInput, { target: { value: 'Motor' } }); fireEvent.blur(nameInput);
    await waitFor(() => expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Motor'));
    await waitFor(() => expect(within(row).getByRole('textbox')).toHaveValue('Motor'));
    const required = within(row).getByRole('button', { name: 'step4.group.members.required_label' });
    fireEvent.click(required); expect(required).toHaveAttribute('aria-pressed', 'false');
    expect(group.members[0].required).toBe(true);
    const count = screen.getByTestId('group-member-scope');
    expect(count).toHaveClass('recording-editor-point-count');
    fireEvent.change(screen.getByTestId('group-member-search'), { target: { value: 'Motor' } });
    await waitFor(() => expect(count).toHaveTextContent('compact.selectedScope'));
    expect(pointFormatAPI.read).toHaveBeenCalledWith({ workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0' });
    expect(pointFormatAPI.save).not.toHaveBeenCalled(); expect(onSaved).not.toHaveBeenCalled();
  });
});
