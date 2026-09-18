import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor, act } from '@testing-library/react';
import { RecordingPlanSetupSection } from '@/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { WorkspaceRecordingPlanSetupSection } from '@/features/datalink/workbench-v2/steps/step4/WorkspaceRecordingPlanSetupSection';

const mocks = vi.hoisted(() => ({ create: vi.fn(), preview: vi.fn(), plans: [] as unknown[], isError: false, refetch: vi.fn() }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => ({ data: [], isError: true, isSuccess: false, isLoading: false, refetch: mocks.refetch }),
}));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: { defaultValue?: string }) => options?.defaultValue ?? key }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceRecordingPlans', () => ({
  useStudioV2WorkspaceRecordingPlansQuery: () => ({ data: mocks.plans, isLoading: false, isError: mocks.isError, refetch: mocks.refetch }),
  useStudioV2ConnectorCapabilitiesQuery: () => ({ data: [] }),
  useCreateRecordingPlanMutation: () => ({ mutateAsync: mocks.create, isPending: false }),
  usePreviewSchemaMutation: () => ({ mutateAsync: mocks.preview, isPending: false }),
  useApplySchemaMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestWritePlanMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
const connector = { ...INITIAL_STATE.db.connector, connector_id: 'saved', identity_revision: 'identity-1', setup_revision: 'setup-1', workspace_id: 'workspace-a', persisted: true, save_state: 'saved' as const };
const members = [
  { member_id: 'member-a', measurement_id: 'measurement-a', equipment_id: 'equipment-a', name: 'A temperature' },
  { member_id: 'member-b', measurement_id: 'measurement-b', equipment_id: 'equipment-b', name: 'B pressure' },
];
const scope = { workspaceId: 'workspace-a', members, fingerprint: 'membership-1' };
const plan = (id: string) => ({ id, workspace_id: scope.workspaceId, name: id, revision: 'rev-1', members,
  destinations: [{ connector_id: connector.connector_id, connector_revision: connector.identity_revision }],
  streams: members.map(member => ({ measurement_id: member.measurement_id })) });

describe('MultiDeviceMembershipAndPlanSelection', () => {
  beforeEach(() => { mocks.plans = []; mocks.isError = false; mocks.create.mockReset().mockResolvedValue({ id: 'created' }); mocks.preview.mockReset(); mocks.refetch.mockReset(); });

  it('creates with each persisted measurement and its own equipment identity', async () => {
    render(<RecordingPlanSetupSection connector={connector} scope={scope} />);
    fireEvent.click(screen.getByTestId('create-plan-btn'));
    await waitFor(() => expect(mocks.create).toHaveBeenCalledOnce());
    expect(mocks.create.mock.calls[0][0].members).toEqual(members);
    expect(mocks.create.mock.calls[0][0].streams.map((stream: { measurement_id: string }) => stream.measurement_id)).toEqual(['measurement-a', 'measurement-b']);
  });

  it('never invents a measurement when no saved members are available', () => {
    render(<RecordingPlanSetupSection connector={connector} scope={{ ...scope, members: [] }} />);
    expect(screen.getByTestId('create-plan-btn')).toBeDisabled();
    fireEvent.click(screen.getByTestId('create-plan-btn'));
    expect(mocks.create).not.toHaveBeenCalled();
  });

  it('requires an explicit selection when multiple plans exist', async () => {
    mocks.plans = [plan('plan-a'), plan('plan-b')];
    mocks.preview.mockResolvedValue({ token: 'preview-b', statements: ['DDL B'] });
    render(<RecordingPlanSetupSection connector={connector} scope={scope} />);
    expect(screen.queryByTestId('preview-schema-btn')).not.toBeInTheDocument();
    fireEvent.change(screen.getByTestId('recording-plan-selection'), { target: { value: 'plan-b' } });
    fireEvent.click(screen.getByTestId('preview-schema-btn'));
    await waitFor(() => expect(mocks.preview).toHaveBeenCalledWith(expect.objectContaining({ plan_id: 'plan-b' })));
  });

  it('does not auto-select a foreign or mismatched single plan', () => {
    mocks.plans = [{ ...plan('foreign'), workspace_id: 'workspace-b' }];
    render(<RecordingPlanSetupSection connector={connector} scope={scope} />);
    expect(screen.queryByTestId('preview-schema-btn')).not.toBeInTheDocument();
    expect(screen.queryByTestId('create-plan-btn')).not.toBeInTheDocument();
  });

  it('distinguishes a failed list from an empty list and offers retry', () => {
    mocks.isError = true;
    render(<RecordingPlanSetupSection connector={connector} scope={scope} />);
    expect(screen.queryByTestId('create-plan-btn')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('recording-plan-retry'));
    expect(mocks.refetch).toHaveBeenCalledOnce();
  });

  it('does not show an empty-plan create form when membership loading fails', () => {
    render(<WorkspaceRecordingPlanSetupSection state={{ ...INITIAL_STATE, db: { ...INITIAL_STATE.db, connector } }} workspaceId="workspace-a" />);
    expect(screen.getByRole('alert')).toHaveTextContent('step4.recording_members_load_failed');
    expect(screen.queryByTestId('create-plan-btn')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'step4.recording_retry' }));
    expect(mocks.refetch).toHaveBeenCalledOnce();
  });

  it('ignores a preview that finishes after switching plans', async () => {
    mocks.plans = [plan('plan-a'), plan('plan-b')];
    let resolvePreview!: (result: unknown) => void;
    mocks.preview.mockReturnValue(new Promise(resolve => { resolvePreview = resolve; }));
    render(<RecordingPlanSetupSection connector={connector} scope={scope} />);
    fireEvent.change(screen.getByTestId('recording-plan-selection'), { target: { value: 'plan-a' } });
    fireEvent.click(screen.getByTestId('preview-schema-btn'));
    fireEvent.change(screen.getByTestId('recording-plan-selection'), { target: { value: 'plan-b' } });
    await act(async () => resolvePreview({ token: 'old-token', statements: ['OLD DDL A'] }));
    expect(screen.queryByText('OLD DDL A')).not.toBeInTheDocument();
    expect(screen.queryByTestId('apply-schema-btn')).not.toBeInTheDocument();
  });
});
