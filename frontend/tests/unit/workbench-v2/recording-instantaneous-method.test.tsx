import { dbTargetAPI } from '@/services/datalink';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import enConstraints from '@/i18n/locales/en/recording-editor-constraints.json';
import zhConstraints from '@/i18n/locales/zh-TW/recording-editor-constraints.json';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { applyInstantaneousReduction } from '@/features/datalink/workbench-v2/recording-editor/instantaneousReduction';
import { groupToDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { GroupEditor, type GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup, WriteGroupDraft } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) => options && 'min' in options ? `${key}:${String(options.min)}-${String(options.max)}` : key,
  }),
}));

vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));

vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    readiness: vi.fn().mockResolvedValue(null), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(),
    apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
  },
}));

vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({
  studioV2WorkspaceMeasurementsAPI: { list: vi.fn().mockResolvedValue([]) },
}));

function profileGroup(overrides: Partial<WriteGroup> = {}): WriteGroup {
  const base = savedGroup();
  return savedGroup({
    id: 'inst-1', name: 'Line A instantaneous', layout_version: 'managed_long_v1', storage_mode: 'plain',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    members: base.members.map(member => ({ ...member, target_column: '', required: false })),
    destination: { ...base.destination, table_name: 'gateway_long', storage_strategy: 'managed' },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
    ...overrides,
  });
}

function mountEditor(group: WriteGroup) {
  const state = savedGroupState();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const props: GroupEditorProps = {
    presentation: 'simple', group, workspaceId: 'ws-1', workspaceRevision: 'db-setup-9', connector: state.db.connector,
    candidates: buildGroupCandidates(state).candidates, excluded: [], readonly: false,
    onSaved: vi.fn(), onClose: vi.fn(), onReload: vi.fn(),
    fetchLatest: vi.fn().mockResolvedValue({ workspace_revision: 'db-fresh-10', groups: [group] }),
  };
  return render(<QueryClientProvider client={client}><GroupEditor {...props} /></QueryClientProvider>);
}

function updateBody(): WriteGroupDraft {
  return vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mock.calls[0][1].group;
}

// The full suite runs many editors in parallel; the save round trip needs more than the default second.
const WAIT = { timeout: 5000 };

afterEach(cleanup);

beforeEach(() => {
  vi.clearAllMocks();
  vi.spyOn(dbTargetAPI, 'listConnectors').mockResolvedValue([{
    id: 'conn-1', identity_revision: 'crev-1', name: 'Line DB', kind: 'sqlite', enabled: true, status: 'ready',
    connection_config: { path: savedGroupState().db.connector.database }, last_check_error: '', default_write_interval_seconds: 5, created_at: '', updated_at: '',
  }]);
});

describe('instantaneous recording method', () => {
  it('has operator text in both languages', () => {
    for (const bundle of [enConstraints, zhConstraints]) {
      expect(Object.keys(bundle.instantaneousMethod).sort()).toEqual([
        'cancelMean', 'intervalInvalid', 'intervalLabel', 'latest', 'mean', 'meanConsentTitle', 'pendingNote', 'reductionLabel', 'rejected',
      ]);
    }
  });

  it('does not attribute a generic save rejection to a particular point or method', () => {
    expect(zhConstraints.instantaneousMethod.rejected).toContain('檢查記錄設定');
    expect(zhConstraints.instantaneousMethod.rejected).not.toMatch(/拒絕此記錄方式|移除該點位/);
    expect(enConstraints.instantaneousMethod.rejected).toContain('Check the recording settings');
    expect(enConstraints.instantaneousMethod.rejected).not.toMatch(/rejected this recording method|remove the point/);
  });

  it('requires the gauge-use confirmation before mean, then saves mean with the new interval and asks for a fresh schema check', async () => {
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockImplementation(async (_id, request) => ({
      workspace_revision: 'db-fresh-11',
      group: profileGroup({ ...request.group, id: 'inst-1', revision: 'rev-2', destination: { ...profileGroup().destination, ...request.group.destination } } as Partial<WriteGroup>),
    }));
    mountEditor(profileGroup());

    expect(screen.getByTestId('recording-instant-reduction-latest')).toBeChecked();
    expect(screen.queryByText('step4.group.profile.note')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('recording-instant-reduction-mean'));
    // Nothing changes before the operator confirms gauge use.
    expect(screen.getByTestId('recording-instant-mean-consent')).toHaveTextContent('compact.gaugeConsent');
    expect(screen.getByTestId('recording-instant-confirm-gauge').className).toContain('border-[var(--gw-control)]');
    expect(screen.getByTestId('recording-instant-cancel-mean').className).toContain('border-[var(--gw-control)]');
    expect(screen.getByTestId('recording-instant-reduction-latest')).toBeChecked();
    expect(screen.queryByTestId('recording-instant-pending-note')).not.toBeInTheDocument();

    fireEvent.click(screen.getByTestId('recording-instant-confirm-gauge'));
    expect(screen.getByTestId('recording-instant-reduction-mean')).toBeChecked();
    expect(screen.queryByTestId('recording-instant-mean-consent')).not.toBeInTheDocument();
    fireEvent.change(screen.getByTestId('recording-instant-interval'), { target: { value: '30' } });
    expect(screen.getByTestId('recording-instant-pending-note')).toHaveTextContent('instantaneousMethod.pendingNote');

    await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled(), WAIT);
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce(), WAIT);
    const [id, request] = vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mock.calls[0];
    expect(id).toBe('inst-1');
    expect(request.expected_group_revision).toBe('rev-1');
    const body = updateBody();
    expect(body.recording_profile).toEqual({ kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' });
    expect(body.row_policy.interval_seconds).toBe(30);
    expect(body.members.map(member => member.value_role)).toEqual(['gauge', 'gauge']);
    expect(body.members.every(member => member.result_value_type === undefined)).toBe(true);

    // The saved revision now needs a fresh schema check before recording uses it.
    expect(await screen.findByTestId('recording-preparation-check', {}, WAIT)).toBeInTheDocument();
    await waitFor(() => expect(screen.getByTestId('recording-preparation-phase')).toHaveTextContent('preparation.stage.saved'), WAIT);
    expect(screen.queryByTestId('recording-instant-pending-note')).not.toBeInTheDocument();
  });

  it('keeps latest when the confirmation is cancelled and only changes the interval', async () => {
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockResolvedValue({ workspace_revision: 'db-fresh-11', group: profileGroup({ revision: 'rev-2' }) });
    mountEditor(profileGroup());
    fireEvent.click(screen.getByTestId('recording-instant-reduction-mean'));
    fireEvent.click(screen.getByTestId('recording-instant-cancel-mean'));
    expect(screen.queryByTestId('recording-instant-mean-consent')).not.toBeInTheDocument();
    expect(screen.getByTestId('recording-instant-reduction-latest')).toBeChecked();

    fireEvent.change(screen.getByTestId('recording-instant-interval'), { target: { value: '20' } });
    await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled(), WAIT);
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce(), WAIT);
    const body = updateBody();
    expect(body.recording_profile).toEqual({ kind: 'instantaneous', interval: 20, unit: 'second' });
    expect(body.members.every(member => member.value_role === undefined)).toBe(true);
  });

  it('rejects an interval outside the instantaneous bounds without saving', () => {
    mountEditor(profileGroup());
    fireEvent.change(screen.getByTestId('recording-instant-interval'), { target: { value: '301' } });
    expect(screen.getByTestId('recording-instant-interval-error')).toHaveTextContent('instantaneousMethod.intervalInvalid:1-300');
    expect(screen.getByTestId('recording-instant-interval')).toHaveAttribute('aria-invalid', 'true');
    expect(screen.getByTestId('group-save')).toBeDisabled();
  });

  it('shows the backend rejection at the reduction field and keeps the draft', async () => {
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockRejectedValue(Object.assign(new Error('invalid'), {
      response: { status: 422, data: { success: false, error: { code: 'WRITE_GROUP_INVALID', message: 'write-group configuration is invalid', retryable: false, request_id: 'r1', action: 'review_group' } } },
    }));
    mountEditor(profileGroup());
    fireEvent.click(screen.getByTestId('recording-instant-reduction-mean'));
    fireEvent.click(screen.getByTestId('recording-instant-confirm-gauge'));
    await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled(), WAIT);
    fireEvent.click(screen.getByTestId('group-save'));
    expect(await screen.findByTestId('recording-instant-rejected', {}, WAIT)).toHaveTextContent('instantaneousMethod.rejected');
    expect(screen.getByTestId('group-save-error')).toBeInTheDocument();
    expect(screen.getByTestId('recording-instant-reduction-mean')).toBeChecked();
    expect(screen.getByTestId('recording-instant-pending-note')).toBeInTheDocument();
  });

  it('never offers mean for a cumulative recording', () => {
    mountEditor(profileGroup({ name: 'Line A cumulative', recording_profile: { kind: 'cumulative', interval: 10, unit: 'minute' }, row_policy: { interval_seconds: 600, allowed_lateness_seconds: 0, incomplete_policy: 'partial' } }));
    expect(screen.getByTestId('group-profile-summary')).toBeInTheDocument();
    expect(screen.queryByTestId('recording-instant-reduction-mean')).not.toBeInTheDocument();
    expect(screen.queryByTestId('recording-instant-method')).not.toBeInTheDocument();
  });
});

describe('applyInstantaneousReduction', () => {
  const state = savedGroupState();
  const candidates = buildGroupCandidates(state).candidates;
  const withBatch = (): WriteGroup => {
    const group = profileGroup();
    return { ...group, members: [...group.members, { device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2', source_revision: 's2', mapping_revision: 'm2', target_column: '', required: false }] };
  };

  it('keeps a non-numeric source as state and a counter binding as counter_raw', () => {
    const group = withBatch();
    group.members[1] = { ...group.members[1], value_role: 'counter_raw', measurement_id: 'meas-1', result_value_type: 'float64' };
    const next = applyInstantaneousReduction(groupToDraft(group), 'mean', group, candidates);
    expect(next.members.map(member => member.value_role)).toEqual(['gauge', 'counter_raw', 'state']);
    expect(next.members.every(member => member.result_value_type === undefined)).toBe(true);
  });

  it('returns to the exact saved shape when mean is undone, and leaves cumulative drafts unchanged', () => {
    const group = withBatch();
    const draft = groupToDraft(group);
    const back = applyInstantaneousReduction(applyInstantaneousReduction(draft, 'mean', group, candidates), 'latest', group, candidates);
    expect(back).toEqual(draft);
    const cumulative = groupToDraft(profileGroup({ recording_profile: { kind: 'cumulative', interval: 10, unit: 'minute' } }));
    expect(applyInstantaneousReduction(cumulative, 'mean', undefined, candidates)).toBe(cumulative);
  });
});
