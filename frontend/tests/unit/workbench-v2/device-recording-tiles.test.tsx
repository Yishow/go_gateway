import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) =>
      options && 'count' in options ? `${key}:${options.count}` : key,
  }),
}));

const groupsQuery = {
  data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] },
  isLoading: false,
  isError: false,
  refetch: vi.fn(),
};
const ensureAsync = vi.fn();
const startAsync = vi.fn();
const updateAsync = vi.fn();
const ensureCircuitAsync = vi.fn();

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useUpdateWriteGroupMutation: () => ({ mutateAsync: updateAsync, isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensureCircuitAsync, isPending: false }),
  useWriteGroupsQuery: () => groupsQuery,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useWriteGroupSchemaPreviewMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useWriteGroupSchemaApplyMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({
  GroupSchemaPanel: (props: { group: WriteGroup }) => <div data-testid="circuit-schema-panel">{props.group.id}</div>,
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => ({ data: [], isLoading: false, isError: false, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: ensureAsync, isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: startAsync, isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, isLoading: false, isError: false, refetch: vi.fn() }),
}));

// A and B share the same display name and address; only the persisted ID differs.
// B also gets its own saved sources so both devices are recordable.
function twoSharingDevices() {
  const state = savedGroupState();
  const first = state.devices[0];
  const deviceB = { ...first, id: 'dev-B' };
  const points = [...state.points];
  const mappings = { ...state.mappings };
  for (const point of state.points) {
    const pointB = { ...point, id: `${point.id}-b`, device_id: 'dev-B' };
    points.push(pointB);
    const mapping = state.mappings[point.id];
    mappings[pointB.id] = {
      ...mapping, point_id: pointB.id, device_id: 'dev-B',
      persisted_point_id: `${mapping.persisted_point_id}-b`, tag_id: `${mapping.tag_id}-b`,
    };
  }
  return savedGroupState({ devices: [{ ...first }, deviceB], points, mappings });
}

function renderPanel(state = twoSharingDevices()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <BasicRecordingPanel state={state} workspaceId="ws-1" readonly={false} />
    </QueryClientProvider>,
  );
}

describe('device recording tiles', () => {
  beforeEach(() => {
    sessionStorage.clear();
    groupsQuery.data.groups = [];
    ensureAsync.mockReset();
    startAsync.mockReset();
    updateAsync.mockReset();
    ensureCircuitAsync.mockReset();
    groupsQuery.refetch.mockReset();
  });

  it('renders one ID-bound native device tile per device and no dropdown', () => {
    renderPanel();

    expect(screen.queryByTestId('basic-recording-device')).toBeNull();
    // A legal non-device select (the circuit chooser) stays.
    expect(screen.getByTestId('circuit-profile-circuit-chooser').tagName).toBe('SELECT');

    const tiles = screen.getAllByTestId('basic-recording-device-tile');
    expect(tiles).toHaveLength(2);
    expect(tiles.map((tile) => tile.getAttribute('data-device-id'))).toEqual(['dev-1', 'dev-B']);
    for (const tile of tiles) {
      expect(tile.tagName).toBe('BUTTON');
      expect(tile).toHaveAttribute('type', 'button');
    }
    expect(tiles[0]).toHaveAttribute('aria-pressed', 'true');
    expect(tiles[1]).toHaveAttribute('aria-pressed', 'false');
  });

  it('keeps the legacy completeness select on the legacy device only', () => {
    groupsQuery.data.groups = [savedGroup({
      id: 'legacy-wide-a', name: 'Legacy A', status: 'ready',
      basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1',
    })];
    renderPanel();

    expect(screen.queryByTestId('basic-recording-device')).toBeNull();
    expect(screen.getAllByTestId('basic-recording-device-tile')).toHaveLength(2);
    expect(screen.getByTestId('basic-recording-completeness').tagName).toBe('SELECT');
  });

  it('lets device B use the new branch while legacy wide A stays on its own scope', () => {
    groupsQuery.data.groups = [savedGroup({
      id: 'legacy-wide-a', name: 'Legacy A', status: 'ready',
      basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1',
    })];
    renderPanel();

    fireEvent.click(screen.getAllByTestId('basic-recording-device-tile')[1]);
    expect(screen.queryByTestId('basic-recording-device')).toBeNull();
    expect(screen.queryByTestId('basic-recording-completeness')).toBeNull();
    expect(screen.getByTestId('circuit-profile-source-section')).toBeTruthy();
  });

  it('round-trips legacy A to new B and back to A with A\'s branch restored', () => {
    groupsQuery.data.groups = [savedGroup({
      id: 'legacy-wide-a', name: 'Legacy A', status: 'ready',
      basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1',
    })];
    renderPanel();
    fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '47' } });
    const tileB = () => screen.getAllByTestId('basic-recording-device-tile')[1];
    fireEvent.click(tileB());
    // B is not dragged onto A's legacy branch: the new circuit subtree is shown.
    expect(screen.queryByTestId('basic-recording-completeness')).toBeNull();
    expect(screen.getByTestId('circuit-profile-source-section')).toBeTruthy();
    expect(tileB()).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).not.toHaveValue('47');

    fireEvent.click(screen.getAllByTestId('basic-recording-device-tile')[0]);
    expect(screen.getByTestId('basic-recording-completeness').tagName).toBe('SELECT');
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(47);
    expect(screen.getAllByTestId('basic-recording-device-tile')[0]).toHaveAttribute('aria-pressed', 'true');
    expect(ensureAsync).not.toHaveBeenCalled();
  });

  it('keeps A\'s unsaved legacy interval across an unmount and remount', () => {
    groupsQuery.data.groups = [savedGroup({
      id: 'legacy-wide-a', name: 'Legacy A', status: 'ready',
      basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1',
    })];
    const first = renderPanel();
    fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '47' } });
    first.unmount();

    renderPanel();
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(47);
  });

  it('does not carry an unsaved legacy interval into another workspace', () => {
    groupsQuery.data.groups = [savedGroup({
      id: 'legacy-wide-a', name: 'Legacy A', status: 'ready',
      basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1',
    })];
    const view = renderPanel();
    fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '47' } });

    view.rerender(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
        <BasicRecordingPanel state={twoSharingDevices()} workspaceId="ws-2" readonly={false} />
      </QueryClientProvider>,
    );
    expect(screen.getByTestId('basic-recording-interval')).not.toHaveValue(47);
  });

  it('keeps B and A in their own scopes when A\'s delayed save reply arrives after switching to B', async () => {
    let releaseA: (reply: unknown) => void = () => {};
    ensureCircuitAsync.mockImplementation(() => new Promise((resolve) => { releaseA = resolve; }));
    renderPanel();
    const tile = (index: number) => screen.getAllByTestId('basic-recording-device-tile')[index];
    const interval = () => screen.getByTestId('circuit-profile-interval-instantaneous');

    fireEvent.click(tile(1));
    fireEvent.change(interval(), { target: { value: '33' } });
    fireEvent.click(tile(0));
    fireEvent.change(interval(), { target: { value: '45' } });
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureCircuitAsync).toHaveBeenCalledTimes(1));
    const request = ensureCircuitAsync.mock.calls[0][0].request;
    expect(request.group.members.every((member: { device_id: string }) => member.device_id === 'dev-1')).toBe(true);

    // A is still saving when the operator moves to B, whose draft is its own.
    fireEvent.click(tile(1));
    expect(tile(1)).toHaveAttribute('aria-pressed', 'true');
    expect(interval()).toHaveValue('33');

    await act(async () => {
      releaseA({ workspace_revision: 'wrev-2', group: savedGroup({ ...request.group, id: 'grp-A', revision: 'rev-A-1', status: 'draft' }) });
    });
    expect(interval()).toHaveValue('33');
    expect(screen.queryByTestId('circuit-profile-error')).toBeNull();
    expect(screen.queryByTestId('circuit-schema-panel')).toBeNull();
    expect(tile(1)).toHaveAttribute('aria-pressed', 'true');

    // Back on A: its own draft and identity, with nothing created, started or applied for B.
    fireEvent.click(tile(0));
    expect(interval()).toHaveValue('45');
    expect(screen.queryByTestId('circuit-profile-error')).toBeNull();
    expect(ensureCircuitAsync).toHaveBeenCalledTimes(1);
    expect(startAsync).not.toHaveBeenCalled();
    expect(updateAsync).not.toHaveBeenCalled();
  });

  it('keeps a cross-device custom group whole when devices are selected', () => {
    const custom = savedGroup({
      id: 'custom-ab', name: 'A and B custom', status: 'ready',
      destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: '/d.db', table_schema: '', table_name: 'custom_ab', storage_strategy: 'custom' },
      members: [
        { device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', source_revision: 's0', mapping_revision: 'm0', target_column: 'a', required: true },
        { device_id: 'dev-B', point_id: 'pt-0-b', tag_id: 'tag-0-b', source_revision: 's0', mapping_revision: 'm0', target_column: 'b', required: true },
      ],
    });
    groupsQuery.data.groups = [custom];
    renderPanel();
    const before = JSON.stringify(groupsQuery.data.groups);
    fireEvent.click(screen.getAllByTestId('basic-recording-device-tile')[1]);
    fireEvent.click(screen.getAllByTestId('basic-recording-device-tile')[0]);
    expect(JSON.stringify(groupsQuery.data.groups)).toBe(before);
    expect(custom.id).toBe('custom-ab');
    expect(custom.members.map((member) => member.device_id)).toEqual(['dev-1', 'dev-B']);
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync).not.toHaveBeenCalled();
    expect(updateAsync).not.toHaveBeenCalled();
  });

  it('restores the selected device and its own scope after a reload', () => {
    const first = renderPanel();
    fireEvent.click(screen.getAllByTestId('basic-recording-device-tile')[1]);
    first.unmount();

    renderPanel();
    const tiles = screen.getAllByTestId('basic-recording-device-tile');
    expect(tiles[0]).toHaveAttribute('aria-pressed', 'false');
    expect(tiles[1]).toHaveAttribute('aria-pressed', 'true');
    expect(ensureAsync).not.toHaveBeenCalled();
  });

  it('keeps A\'s unsaved completeness choice across an unmount and remount', () => {
    groupsQuery.data.groups = [savedGroup({
      id: 'legacy-wide-a', name: 'Legacy A', status: 'ready',
      basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1',
    })];
    const first = renderPanel();
    fireEvent.change(screen.getByTestId('basic-recording-completeness'), { target: { value: 'partial' } });
    first.unmount();

    renderPanel();
    expect(screen.getByTestId('basic-recording-completeness')).toHaveValue('partial');
    expect(ensureAsync).not.toHaveBeenCalled();
  });

  it('restores a saved selection only in its own workspace and ignores an unknown device ID', () => {
    const first = renderPanel();
    fireEvent.click(screen.getAllByTestId('basic-recording-device-tile')[1]);
    first.unmount();

    const other = render(
      <QueryClientProvider client={new QueryClient()}>
        <BasicRecordingPanel state={twoSharingDevices()} workspaceId="ws-2" readonly={false} />
      </QueryClientProvider>,
    );
    expect(screen.getAllByTestId('basic-recording-device-tile')[0]).toHaveAttribute('aria-pressed', 'true');
    other.unmount();

    sessionStorage.setItem('wbv2.device-scope.ws-1', 'dev-ghost');
    renderPanel();
    expect(screen.getAllByTestId('basic-recording-device-tile')[0]).toHaveAttribute('aria-pressed', 'true');
    expect(ensureCircuitAsync).not.toHaveBeenCalled();
    expect(updateAsync).not.toHaveBeenCalled();
  });

  describe('selection authority', () => {
    const element = (client: QueryClient, state: ReturnType<typeof twoSharingDevices>, workspaceId?: string) => (
      <QueryClientProvider client={client}>
        <BasicRecordingPanel state={state} workspaceId={workspaceId} readonly={false} />
      </QueryClientProvider>
    );
    const pressed = () => screen.getAllByTestId('basic-recording-device-tile')
      .filter((tile) => tile.getAttribute('aria-pressed') === 'true').map((tile) => tile.getAttribute('data-device-id'));
    const noMutation = () => {
      expect(ensureAsync).not.toHaveBeenCalled();
      expect(ensureCircuitAsync).not.toHaveBeenCalled();
      expect(startAsync).not.toHaveBeenCalled();
      expect(updateAsync).not.toHaveBeenCalled();
    };

    it('restores the saved device when the workspace arrives after the first render', () => {
      const first = renderPanel();
      fireEvent.click(screen.getAllByTestId('basic-recording-device-tile')[1]);
      first.unmount();

      const client = new QueryClient();
      const view = render(element(client, twoSharingDevices(), undefined));
      view.rerender(element(client, twoSharingDevices(), 'ws-1'));
      expect(pressed()).toEqual(['dev-B']);
      noMutation();
    });

    it('restores the saved device when the devices hydrate after the first render, and keeps workspaces apart', () => {
      const first = renderPanel();
      fireEvent.click(screen.getAllByTestId('basic-recording-device-tile')[1]);
      first.unmount();

      const client = new QueryClient();
      const empty = savedGroupState({ devices: [], points: [], mappings: {} });
      const view = render(element(client, empty, 'ws-1'));
      view.rerender(element(client, twoSharingDevices(), 'ws-1'));
      expect(pressed()).toEqual(['dev-B']);

      view.rerender(element(client, twoSharingDevices(), 'ws-2'));
      expect(pressed()).toEqual(['dev-1']);
      view.rerender(element(client, twoSharingDevices(), 'ws-1'));
      expect(pressed()).toEqual(['dev-B']);
      noMutation();
    });

    it('keeps a persisted device that has no recordable sources selected', () => {
      const withSources = twoSharingDevices();
      const devices = withSources.devices;
      const onlyA = savedGroupState({
        devices,
        points: withSources.points.filter((point) => point.device_id !== 'dev-B'),
        mappings: Object.fromEntries(Object.entries(withSources.mappings).filter(([, mapping]) => mapping.device_id !== 'dev-B')),
      });
      renderPanel(onlyA);

      const tiles = screen.getAllByTestId('basic-recording-device-tile');
      expect(tiles.map((tile) => tile.getAttribute('data-device-id'))).toEqual(['dev-1', 'dev-B']);
      fireEvent.click(tiles[1]);
      expect(pressed()).toEqual(['dev-B']);
      expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
      noMutation();
    });
  });

  it('selects the other same-named device by its persisted ID without any recording mutation', () => {
    renderPanel();

    const tiles = screen.getAllByTestId('basic-recording-device-tile');
    const second = tiles.find((tile) => tile.getAttribute('data-device-id') === 'dev-B');
    expect(second).toBeDefined();
    expect(second!.tagName).toBe('BUTTON');
    expect(second).toHaveAttribute('type', 'button');
    fireEvent.click(second!);

    // The device-scoped subtree is rebuilt on switch, so re-read the current DOM.
    const pressed = () => screen.getAllByTestId('basic-recording-device-tile')
      .find((tile) => tile.getAttribute('data-device-id') === 'dev-B');
    expect(pressed()).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getAllByTestId('basic-recording-device-tile')[0]).toHaveAttribute('aria-pressed', 'false');
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync).not.toHaveBeenCalled();
    expect(updateAsync).not.toHaveBeenCalled();
  });
});
