import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { MeasurementDefinition } from '@/types/measurement';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';

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
const measurementsQuery = {
  data: [] as MeasurementDefinition[],
  isLoading: false,
  isError: false,
  refetch: vi.fn(),
};
const ensureCircuitProfileAsync = vi.fn();
const startAsync = vi.fn();

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => groupsQuery,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensureCircuitProfileAsync, isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: startAsync, isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, isLoading: false, isError: false, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurementsQuery,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

function makeConfirmedCounterMeasurement(pointId: string, tagId: string): MeasurementDefinition {
  return {
    id: `meas-${tagId}`,
    workspace_id: 'ws-1',
    definition_revision: 'mrev-1',
    name: `Energy ${tagId}`,
    equipment_id: 'eq-1',
    device_id: 'dev-1',
    point_id: pointId,
    tag_id: tagId,
    source_binding_revision: 's1',
    quantity: 'energy',
    semantic_kind: 'counter',
    unit: 'kWh',
    series_epoch: '1',
    counter_policy: { allow_negative: false },
    created_at: '2026-10-20T00:00:00Z',
    updated_at: '2026-10-20T00:00:00Z',
  };
}

function renderPanel() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  client.setQueryData(studioV2WorkspaceKeys.writeGroups(), groupsQuery.data);
  return render(
    <QueryClientProvider client={client}>
      <BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} />
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.refetch.mockImplementation(async () => ({ data: groupsQuery.data }));
  measurementsQuery.data = [];
  measurementsQuery.refetch.mockImplementation(async () => ({ data: measurementsQuery.data }));
});

describe('CircuitProfile a11y labels and empty source protection', () => {
  it('provides accessible names containing tag name and profile type and verifies focus and native checkbox selection', async () => {
    renderPanel();
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    expandSources();
    const instCheckbox = screen.getByRole('checkbox', {
      name: `line.temperature step4.basic.circuit.instantaneous_label`,
    });
    const cumCheckbox = screen.getByRole('checkbox', {
      name: `line.temperature step4.basic.circuit.cumulative_label`,
    });

    expect(instCheckbox).toBeInTheDocument();
    expect(cumCheckbox).toBeInTheDocument();
    expect(instCheckbox).toBeChecked();
    expect(cumCheckbox).not.toBeChecked();
    fireEvent.click(cumCheckbox);
    expect(cumCheckbox).toBeChecked();

    // Verify focus accessibility
    instCheckbox.focus();
    expect(instCheckbox).toHaveFocus();

    // Native checkbox toggle: toggling instantaneous source off leaves cumulative checked
    fireEvent.click(instCheckbox);
    expect(instCheckbox).not.toBeChecked();
    expect(cumCheckbox).toBeChecked();

    // Toggle instantaneous source back on
    fireEvent.click(instCheckbox);
    expect(instCheckbox).toBeChecked();
    expect(cumCheckbox).toBeChecked();
  });

  it('displays actionable error and blocks prepare/start when instantaneous is enabled but has 0 sources', async () => {
    measurementsQuery.data = [makeConfirmedCounterMeasurement('pt-1', 'tag-1')];
    renderPanel();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Feeder A' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    expandSources();
    // Select tag-1 for cumulative only; uncheck all instantaneous sources
    fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-1'));
    const instTag0 = screen.getByTestId('circuit-profile-member-inst-tag-0');
    const instTag1 = screen.getByTestId('circuit-profile-member-inst-tag-1');
    const instTag2 = screen.getByTestId('circuit-profile-member-inst-tag-2');
    if ((instTag0 as HTMLInputElement).checked) fireEvent.click(instTag0);
    if ((instTag1 as HTMLInputElement).checked) fireEvent.click(instTag1);
    if ((instTag2 as HTMLInputElement).checked) fireEvent.click(instTag2);

    // Instantaneous interval is valid (15s), but has 0 sources selected
    expect(screen.getByText('step4.basic.circuit.instantaneous_sources_required')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toHaveAttribute('data-primary', 'true');
    const repair = screen.getByRole('button', { name: 'next.repair' });
    expect(repair).toBeEnabled(); fireEvent.click(repair);
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveFocus();
    expect(ensureCircuitProfileAsync).not.toHaveBeenCalled(); expect(startAsync).not.toHaveBeenCalled();
  });

  it('displays actionable error and blocks prepare/start when cumulative is enabled but has 0 sources', async () => {
    renderPanel();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Feeder B' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    expandSources();
    // Uncheck all cumulative sources
    const cumTag0 = screen.getByTestId('circuit-profile-member-cum-tag-0');
    const cumTag1 = screen.getByTestId('circuit-profile-member-cum-tag-1');
    const cumTag2 = screen.getByTestId('circuit-profile-member-cum-tag-2');
    if ((cumTag0 as HTMLInputElement).checked) fireEvent.click(cumTag0);
    if ((cumTag1 as HTMLInputElement).checked) fireEvent.click(cumTag1);
    if ((cumTag2 as HTMLInputElement).checked) fireEvent.click(cumTag2);

    // Cumulative interval is valid (10m), but has 0 sources selected
    expect(screen.getByText('step4.basic.circuit.cumulative_sources_required')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toHaveAttribute('data-primary', 'true');
    const repair = screen.getByRole('button', { name: 'next.repair' });
    expect(repair).toBeEnabled(); fireEvent.click(repair);
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveFocus();
    expect(ensureCircuitProfileAsync).not.toHaveBeenCalled(); expect(startAsync).not.toHaveBeenCalled();
  });
});
