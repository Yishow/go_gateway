import { fireEvent, render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { WorkbenchProvider, useWorkbench } from '@/pages/datalink/workbench/WorkbenchProvider';
import { TagBindingStudio } from '@/pages/datalink/workbench/TagBindingStudio';
import type { Device, Mapping, Point, SourceRuleRecord, Tag } from '@/types/datalink';
import type { SourceRuleCandidateSnapshotView } from '@/types/sourceRuleCandidates';

const {
  mockDevices,
  mockPoints,
  mockMappings,
  mockTags,
  mockSourceRules,
  mockCandidateViews,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockMappings: [] as Mapping[],
  mockTags: [] as Tag[],
  mockSourceRules: [] as SourceRuleRecord[],
  mockCandidateViews: {} as Record<string, SourceRuleCandidateSnapshotView>,
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('@/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: mockDevices, isLoading: false }),
}));

vi.mock('@/hooks/datalink/usePoints', () => ({
  usePointsQuery: () => ({ data: mockPoints, isLoading: false }),
  useDeletePointMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('@/hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({ data: mockTags, isLoading: false, refetch: vi.fn() }),
  useCreateTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('@/hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mockMappings, isLoading: false }),
  useCreateMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('@/hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: () => ({
    data: mockSourceRules,
    isLoading: false,
    isSuccess: true,
  }),
}));

vi.mock('@/hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: (ruleId?: string | null) => ({
    data: ruleId ? mockCandidateViews[ruleId] ?? null : null,
    isLoading: false,
    isError: false,
    error: null,
    isRefetching: false,
    refetch: vi.fn(),
  }),
}));

vi.mock('@/services/datalink', () => ({
  tagAPI: {
    batchCreate: vi.fn(),
  },
}));

function makeCandidateSnapshot(
  ruleId: string,
  revisionId: string,
  pointId: string,
  address: string,
): SourceRuleCandidateSnapshotView {
  return {
    source_rule_id: ruleId,
    revision_id: revisionId,
    tags: {
      status: 'ready',
      candidates: [
        {
          id: `candidate-${pointId}`,
          identity: {
            source_rule_id: ruleId,
            candidate_type: 'tags',
            candidate_kind: 'tag',
            derived_from_rule_address: address,
          },
          proposed_signature: `signature-${revisionId}-${pointId}`,
          address,
          point_id: pointId,
          tag_key: `TAG_${address}`,
          display_name: `${pointId}-tag`,
          data_type: 'int16',
          status: 'draft',
        },
      ],
    },
    database_outputs: { status: 'deferred', candidates: [] },
    local_modbus_outputs: { status: 'deferred', candidates: [] },
  };
}

function TagStudioHarness() {
  const { crossStepContext, setFocusedRuleId, setSelectedDeviceId } = useWorkbench();
  const [, setRefreshTick] = useState(0);

  return (
    <>
      <button
        type="button"
        onClick={() => {
          setSelectedDeviceId('device-1');
          setFocusedRuleId('rule-1');
        }}
      >
        focus-rule-1
      </button>
      <button
        type="button"
        onClick={() => {
          setSelectedDeviceId('device-1');
          setFocusedRuleId('rule-2');
        }}
      >
        focus-rule-2
      </button>
      <button type="button" onClick={() => setRefreshTick((value) => value + 1)}>
        refresh-candidate-snapshot
      </button>
      <div data-testid="focused-rule-id">{crossStepContext.focusedRuleId ?? ''}</div>
      <TagBindingStudio />
    </>
  );
}

function renderTagStudio() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <WorkbenchProvider>
        <TagStudioHarness />
      </WorkbenchProvider>
    </QueryClientProvider>,
  );
}

describe('TagBindingStudio source-rule scope', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockDevices.splice(0, mockDevices.length, {
      id: 'device-1',
      name: 'Mixer PLC',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    });
    mockPoints.splice(0, mockPoints.length,
      {
        id: 'point-1',
        device_id: 'device-1',
        name: 'Flow Sensor',
        description: '',
        data_type: 'int16',
        address: '40001',
        enabled: true,
        polling_group_id: '',
        last_value: null,
        last_read_at: '',
        last_error: '',
        error_count: 0,
        created_at: '',
        updated_at: '',
      },
      {
        id: 'point-2',
        device_id: 'device-1',
        name: 'Pressure Sensor',
        description: '',
        data_type: 'int16',
        address: '40002',
        enabled: true,
        polling_group_id: '',
        last_value: null,
        last_read_at: '',
        last_error: '',
        error_count: 0,
        created_at: '',
        updated_at: '',
      },
    );
    mockMappings.splice(0, mockMappings.length);
    mockTags.splice(0, mockTags.length);
    mockSourceRules.splice(0, mockSourceRules.length,
      {
        id: 'rule-1',
        device_id: 'device-1',
        start_address: '40001',
        count: 1,
        data_type: 'int16',
        naming_prefix: 'SRC_A',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
        revision_id: 'rev-1',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'rule-2',
        device_id: 'device-1',
        start_address: '40002',
        count: 1,
        data_type: 'int16',
        naming_prefix: 'SRC_B',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
        revision_id: 'rev-2',
        created_at: '',
        updated_at: '',
      },
    );
    mockCandidateViews['rule-1'] = makeCandidateSnapshot('rule-1', 'rev-1', 'point-1', '40001');
    mockCandidateViews['rule-2'] = makeCandidateSnapshot('rule-2', 'rev-2', 'point-2', '40002');
  });

  it('scopes tag candidates to the focused source rule handoff', () => {
    renderTagStudio();
    fireEvent.click(screen.getByRole('button', { name: 'focus-rule-1' }));
    expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('rule-1');
    expect(screen.getByLabelText('Flow Sensor')).toBeInTheDocument();
    expect(screen.queryByLabelText('Pressure Sensor')).not.toBeInTheDocument();
  });

  it('clears stale diff preview when the focused source rule changes', () => {
    renderTagStudio();
    fireEvent.click(screen.getByRole('button', { name: 'focus-rule-1' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    expect(screen.getByTestId('batch-diff-preview')).toHaveTextContent('TAG_40001');

    fireEvent.click(screen.getByRole('button', { name: 'focus-rule-2' }));

    expect(screen.queryByTestId('batch-diff-preview')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Flow Sensor')).not.toBeInTheDocument();
    expect(screen.getByLabelText('Pressure Sensor')).toBeInTheDocument();
  });

  it('invalidates stale diff preview when the active candidate snapshot changes', () => {
    renderTagStudio();
    fireEvent.click(screen.getByRole('button', { name: 'focus-rule-1' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    expect(screen.getByTestId('batch-diff-preview')).toHaveTextContent('TAG_40001');

    mockCandidateViews['rule-1'] = makeCandidateSnapshot('rule-1', 'rev-3', 'point-2', '40002');
    fireEvent.click(screen.getByRole('button', { name: 'refresh-candidate-snapshot' }));

    expect(screen.queryByTestId('batch-diff-preview')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Flow Sensor')).not.toBeInTheDocument();
    expect(screen.getByLabelText('Pressure Sensor')).toBeInTheDocument();
  });
});
