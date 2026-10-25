import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { BasicRecordingMembers } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingMembers';
import { CircuitCounterConfirmation } from '@/features/datalink/workbench-v2/steps/step4/CircuitCounterConfirmation';
import { CircuitProfileStatusList } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusList';
import { allocationRowsFromCanonicalMappings } from '@/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface';
import { CutoverReviewView } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverReviewView';
import type { GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { ModbusShareDesiredMapping } from '@/types/modbusShare';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { reviewOf, wide } from '../../fixtures/cutoverFixtures';
import { savedGroup } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false, isLoading: false }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus', () => ({
  CircuitProfileStatus: ({ sourceLabels }: { sourceLabels?: Record<string, string> }) => <output data-testid="status-source-label">{sourceLabels?.['tag-1']}</output>,
}));

const candidate = (overrides: Partial<GroupCandidate> = {}): GroupCandidate => ({
  key: 'pt-1|tag-1', device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1',
  tag_key: 'internal.tag.1', label: '進水溫度', device_name: 'PLC A', address: '40002', target_type: 'float64',
  ...overrides,
});

describe('Step 4 source names', () => {
  it('passes candidate labels to the saved profile status instead of Tag Keys', () => {
    const group = longProfileGroup({ id: 'group-1', circuitId: 'circuit-1', circuitName: 'Feeder A', points: [0, 1], interval: 15 });
    render(<CircuitProfileStatusList groups={[group]} candidates={[candidate()]} />);
    expect(screen.getByTestId('status-source-label')).toHaveTextContent('進水溫度');
    expect(screen.getByTestId('status-source-label')).not.toHaveTextContent('internal.tag.1');
  });

  it('uses a human source name in selected members and an address when the name is absent', () => {
    render(<BasicRecordingMembers candidates={[candidate(), candidate({ key: 'pt-2|tag-2', point_id: 'pt-2', tag_id: 'tag-2', label: '', tag_key: 'internal.tag.2', address: '40003' })]} />);
    expandSources('basic-recording-sources');
    expect(screen.getByText('進水溫度')).toBeInTheDocument();
    expect(screen.getByTestId('basic-recording-selected-members')).toHaveTextContent('40003');
    expect(screen.queryByText('internal.tag.1')).not.toBeInTheDocument();
    expect(screen.queryByText('internal.tag.2')).not.toBeInTheDocument();
  });

  it('uses the mapped display name for share allocation rows, then the source address', () => {
    const mappings = [
      { workspace_id: 'ws-1', source_rule_revision: 'rule-rev-1', tag_id: 'tag-1', tag_key: 'internal.tag.1', display_name: '累積電量', source_rule_id: 'rule-1', share_start_register: 1, zero_based_register: 0, span_registers: 1, stride_registers: 1, data_type: 'float64' },
      { workspace_id: 'ws-1', source_rule_revision: 'rule-rev-1', tag_id: 'tag-2', tag_key: 'internal.tag.2', display_name: '', source_rule_id: 'rule-1', share_start_register: 2, zero_based_register: 1, span_registers: 1, stride_registers: 1, data_type: 'float64' },
    ] satisfies ModbusShareDesiredMapping[];
    const rows = allocationRowsFromCanonicalMappings(mappings, {
      rules: [{ id: 'rule-1', name: 'Rule 1' } as never],
      points: [{ id: 'pt-1', address: '40002' }, { id: 'pt-2', address: '40003' }] as never,
      mappings: {
        'pt-1': { tag_id: 'tag-1', point_id: 'pt-1', display_name: '累積電量' } as never,
        'pt-2': { tag_id: 'tag-2', point_id: 'pt-2', display_name: '' } as never,
      },
    });
    expect(rows.map((row) => row.sourceLabel)).toEqual(['累積電量', '40003']);
  });

  it('renders the counter confirmation legend with the human source name without changing the source id', () => {
    render(<CircuitCounterConfirmation workspaceId="ws-1" workspaceRevision="rev-1" circuitId="circuit-1" readonly={false}
      unconfirmedStatuses={[{ candidate: candidate(), confirmed: false }]} onConfirmed={vi.fn()} />);
    expect(screen.getByText('進水溫度')).toBeInTheDocument();
    expect(screen.queryByText('internal.tag.1')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-counter-confirmation-source-tag-1')).toBeInTheDocument();
  });

  it('renders cutover moved and remaining sources by candidate label, not Tag Key', () => {
    const review = reviewOf({ source: { ...reviewOf().source, ...{ group_id: wide().id } } });
    render(<CutoverReviewView review={review} stale={false} candidates={[candidate({ point_id: 'pt-0', tag_id: 'tag-0', key: 'pt-0|tag-0' })]} sourceGroup={savedGroup()} />);
    expect(screen.getByTestId('cutover-moved')).toHaveTextContent('進水溫度');
    expect(screen.queryByText('internal.tag.1')).not.toBeInTheDocument();
  });
});
