import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';
import { parseWriteGroupDeliveryData } from '@/utils/studioV2WriteGroupDeliveryJson';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { expandStatus } from '../../fixtures/sourceReview';

const { query, deliveryFor, deliveries } = vi.hoisted(() => ({
  query: { data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] }, isLoading: false, isError: false, refetch: vi.fn() },
  deliveryFor: vi.fn(),
  deliveries: {} as Record<string, unknown>,
}));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', async importOriginal => ({
  ...await importOriginal<typeof import('@/hooks/datalink/useStudioV2WriteGroups')>(),
  useWriteGroupsQuery: () => query,
  useWriteGroupDeliveryQuery: (id: string, enabled: boolean) => { deliveryFor(id, enabled); return { data: enabled ? deliveries[id] : undefined, isError: false, isLoading: false }; },
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useRecordingStartMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => ({ data: [], isLoading: false, isError: false, refetch: vi.fn() }),
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({ GroupSchemaPanel: () => <div /> }));

const APPLIED = 'applied-1';
const state = batchPanelState();
const applied = (id: string, circuit: string, points: [number, number], revision = APPLIED): WriteGroup =>
  ({ ...longProfileGroup({ id, circuitId: circuit, circuitName: circuit, interval: 15, points }), applied_revision: revision });

/** A delivery answer in the exact JSON shape the API sends, read through the real parser. */
function apiAnswer(groupId: string, tag: string, over: Record<string, unknown> = {}, revision = APPLIED): WriteGroupDelivery {
  const stages = { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 1, skipped: 0 };
  const parsed = parseWriteGroupDeliveryData(JSON.parse(JSON.stringify({
    group_id: groupId, intake: { state: 'active' }, stages, last_sql_committed_at: '2026-10-24T10:00:20Z', oldest_pending_seconds: 0, no_data_buckets: 0,
    skipped_buckets: 0, backlog: [], quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
    revision_stages: { group_revision: revision, stages },
    source_evidence: { state: 'available', group_revision: revision, connector_id: 'conn-1', connector_revision: 'crev-1', sources: [{
      circuit_id: groupId === 'grp-a' ? 'cir-a' : 'cir-b', device_id: 'dev-1', point_id: `pt-${tag.slice(4)}`, tag_id: tag,
      latest: { available: true, group_revision: revision, connector_revision: 'crev-1', record_id: `rec-${tag}`, effect_key: `eff-${tag}`, payload_digest: 'c'.repeat(64),
        profile_kind: 'instantaneous', window_start: '2026-10-24T10:00:00Z', window_end: '2026-10-24T10:00:15Z', status: 'ok', raw_present: true, delivery: 'sql_committed', committed_at: '2026-10-24T10:00:20Z' },
      last_sql_committed: null,
    }] },
    ...over,
  })));
  if (!parsed) throw new Error('the API shaped answer must parse');
  return parsed;
}
const panel = () => <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}><CircuitProfileBasicPanel state={state} workspaceId="ws-1" readonly={false} /></QueryClientProvider>;

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  query.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [applied('grp-a', 'cir-a', [0, 1]), applied('grp-b', 'cir-b', [2, 3])] };
  query.refetch.mockResolvedValue({ data: query.data });
  for (const key of Object.keys(deliveries)) delete deliveries[key];
  deliveries['grp-a'] = apiAnswer('grp-a', 'tag-0');
  deliveries['grp-b'] = apiAnswer('grp-b', 'tag-2');
});

describe('Basic shows the status of the selected circuit canonical group only', () => {
  it('renders the per source status under the active circuit, named from the mapped Tag', () => {
    render(panel());
    expect(screen.getByTestId('circuit-status-instantaneous')).toBeInTheDocument();
    expandStatus('grp-a');
    expect(screen.getByTestId('circuit-status-source-tag-0')).toHaveTextContent('sensor.0');
    expect(screen.queryByTestId('circuit-status-source-tag-2')).not.toBeInTheDocument();
    expect(deliveryFor).toHaveBeenCalledWith('grp-a', true);
    expect(deliveryFor).not.toHaveBeenCalledWith('grp-b', expect.anything());
  });

  it('follows the circuit chooser instead of mixing another circuit evidence in', () => {
    render(panel());
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-chooser'), { target: { value: 'cir-b' } });
    expandStatus('grp-b');
    expect(screen.getByTestId('circuit-status-source-tag-2')).toHaveTextContent('sensor.2');
    expect(screen.queryByTestId('circuit-status-source-tag-0')).not.toBeInTheDocument();
  });

  it('does not show evidence frozen to an older applied revision as the current state', () => {
    deliveries['grp-a'] = apiAnswer('grp-a', 'tag-0', {}, 'older-applied');
    render(panel());
    expect(screen.queryByTestId('circuit-status-source-tag-0')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-status-sources-grp-a-toggle')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-status-evidence-note')).toHaveTextContent('step4.basic.status.evidence_not_current');
  });

  it('does not query a profile group that was saved but never applied', () => {
    query.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [applied('grp-a', 'cir-a', [0, 1], '')] };
    render(panel());
    expect(deliveryFor).toHaveBeenCalledWith('grp-a', false);
    expect(screen.getByTestId('circuit-status-instantaneous')).toHaveTextContent('step4.basic.status.not_applied');
  });
});
