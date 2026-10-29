import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { studioV2RecordingStartAPI as startAPI } from '@/services/studioV2RecordingStart';
import { studioV2WorkspaceWriteGroupsAPI as api } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { clickEnabled, guidedLifecycle, guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { expandSources } from '../../fixtures/sourceReview';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(), ensureCircuitProfile: vi.fn(),
    reviewCutover: vi.fn(), applyCutover: vi.fn(), getCutover: vi.fn(), listCutovers: vi.fn(),
  },
}));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn(), ensureBasic: vi.fn() } }));

const legacyWide = (): WriteGroup => savedGroup({
  id: 'wide-1', name: 'PLC A raw', basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1', status: 'running',
  revision: 'rev-wide', applied_revision: 'applied-wide',
  members: ['temperature', 'pressure', 'batch'].map((column, index) => ({
    device_id: 'dev-1', point_id: `pt-${index}`, tag_id: `tag-${index}`, source_revision: `s${index}`, mapping_revision: `m${index}`, target_column: column, required: index !== 2,
  })),
  destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: 'prod_db', table_schema: 'public', table_name: 'gw_basic_dev_1', storage_strategy: 'managed' },
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0, incomplete_policy: 'skip_row' },
});

let groups: WriteGroup[] = [], workspaceRevision = 'wrev-1';
function canonical(id: string) { const group = groups.find(item => item.id === id); if (!group) throw new Error('Missing canonical group'); return group; }
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn(async (id: string) => guidedLifecycle(canonical(id), workspaceRevision)) } }));
const WRITES = [api.create, api.update, api.apply, api.disable, api.remove, startAPI.start, startAPI.ensureBasic];
const expectNoWriterWrites = () => WRITES.forEach((write) => expect(write).not.toHaveBeenCalled());
const typed = (code: string) => Object.assign(new Error(`Request failed with status code 422: ${code}`), { response: { status: 422, data: { error: { code } } } });

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groups = [legacyWide()]; workspaceRevision = 'wrev-1';
  vi.mocked(api.list).mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: workspaceRevision, groups }));
  vi.mocked(api.readiness).mockImplementation(async id => guidedReadiness(canonical(id), workspaceRevision));
  vi.mocked(api.delivery).mockResolvedValue({
    group_id: 'wide-1', intake: { state: 'active' }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
    last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  } as never);
  vi.mocked(api.listCutovers).mockResolvedValue([]);
});

describe('prepareOnly migration mode surfaces Timescale unavailability truthfully without plain fallback', () => {
  it('keeps selected timescale checked, surfaces timescale_unavailable, and never calls start or groupApply', async () => {
    const pgState = savedGroupState({
      db: { ...savedGroupState().db, connector: { ...savedGroupState().db.connector, kind: 'postgres', name: 'Timescale PG', identity_revision: 'crev-1' } },
    });
    vi.mocked(api.ensureCircuitProfile).mockImplementation(async (request) => {
      const kind = request.group.recording_profile?.kind;
      const group = savedGroup({
        ...request.group, id: `group-${kind}`, revision: 'group-rev-1', status: 'draft', applied_revision: '',
        members: request.group.members.map((member) => ({ ...member, source_revision: 's0', mapping_revision: 'm0' })),
        destination: { ...request.group.destination, database: 'prod_db', table_name: 'gw_long_fixture' },
      });
      groups = [...groups, group]; workspaceRevision = 'wrev-2';
      return { workspace_revision: workspaceRevision, group };
    });
    vi.mocked(api.schemaPreview).mockRejectedValue(typed('timescale_unavailable'));

    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <Step4Database state={pgState} dispatch={vi.fn()} workspaceId="ws-1" />
      </QueryClientProvider>,
    );

    fireEvent.click(await screen.findByTestId('device-recording-edit-existing'));
    fireEvent.click(await screen.findByTestId('legacy-migration-open'));
    expect(await screen.findByTestId('circuit-profile-basic-panel')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare-only')).toBeInTheDocument();

    const tsRadio = screen.getByTestId('circuit-profile-storage-timescale');
    expect(tsRadio).not.toBeDisabled();
    fireEvent.click(tsRadio);
    expect((tsRadio as HTMLInputElement).checked).toBe(true);

    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Timescale Cutover Target' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-instantaneous'));
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
    expandSources();
    fireEvent.click(screen.getByTestId('circuit-profile-member-inst-tag-0'));

    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));

    // Guide contract: preview is an explicit command, separate from save and DDL.
    expect(api.schemaPreview).not.toHaveBeenCalled();
    await clickEnabled('group-schema-preview');
    expect(await screen.findByTestId('circuit-blocker-timescale_unavailable-0')).toBeInTheDocument();
    expect((screen.getByTestId('circuit-profile-storage-timescale') as HTMLInputElement).checked).toBe(true);
    expect(vi.mocked(api.ensureCircuitProfile).mock.calls[0][0].group.storage_mode).toBe('timescale');
    expectNoWriterWrites();
  });
});
