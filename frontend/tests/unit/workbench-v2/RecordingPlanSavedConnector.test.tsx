import { savedRecordingScope, savedRecordingPlanScopeFields } from '../../fixtures/recordingScope';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { RecordingPlanSetupSection } from '@/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';

const mocks = vi.hoisted(() => ({ create: vi.fn(), preview: vi.fn(), apply: vi.fn(), capabilities: vi.fn(), plans: [] as unknown[] }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: { defaultValue?: string }) => options?.defaultValue ?? key }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceRecordingPlans', () => ({
  useStudioV2WorkspaceRecordingPlansQuery: () => ({ data: mocks.plans.map(plan => ({ ...savedRecordingPlanScopeFields, destinations: [{ connector_id: 'saved-pg', connector_revision: 'identity-7' }], ...(plan as Record<string, unknown>) })), isLoading: false, refetch: vi.fn() }),
  useStudioV2ConnectorCapabilitiesQuery: mocks.capabilities,
  useCreateRecordingPlanMutation: () => ({ mutateAsync: mocks.create, isPending: false }),
  usePreviewSchemaMutation: () => ({ mutateAsync: mocks.preview, isPending: false }),
  useApplySchemaMutation: () => ({ mutateAsync: mocks.apply, isPending: false }),
  useTestWritePlanMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

const saved = {
  ...INITIAL_STATE.db.connector,
  connector_id: 'saved-pg', identity_revision: 'identity-7', workspace_id: 'workspace-a',
  setup_revision: 'setup-9', persisted: true, save_state: 'saved' as const, kind: 'postgres' as const,
};

describe('SavedConnectorIdentity recording UI', () => {
  beforeEach(() => {
    mocks.plans = [];
    mocks.create.mockReset().mockResolvedValue({ id: 'plan-created' });
    mocks.preview.mockReset().mockResolvedValue({ token: 'preview', statements: [] });
    mocks.apply.mockReset();
    mocks.capabilities.mockReset().mockReturnValue({ data: [] });
  });

  const previewToken = {
    token: 'tok-1', operation_id: 'op-1', action: 'schema_apply', workspace_id: 'workspace-a',
    workspace_revision: 'setup-9', plan_id: 'plan-a', plan_revision: 'plan-1', connector_id: 'saved-pg',
    connector_revision: 'identity-7', dialect: 'postgres', table_prefix: 'gw_record_',
    statements: ['CREATE TABLE IF NOT EXISTS gw_record_samples (id TEXT)'], tables: [],
    digest: 'a'.repeat(64), expires_at: '2026-09-16T01:00:00Z', created_at: '2026-09-16T00:50:00Z',
  };

  function enableManagedSchema() {
    mocks.capabilities.mockReturnValue({
      data: [{
        kind: 'postgres', supported: true, supports_managed_schema: true, supports_transactions: true,
        supports_receipts: true, supports_test_writes: false, supported_modes: ['managed_recording', 'custom_table'],
      }],
    });
  }

  async function previewThenApply() {
    mocks.plans = [{ id: 'plan-a', name: 'Plan A', revision: 'plan-1' }];
    enableManagedSchema();
    mocks.preview.mockResolvedValue(previewToken);
    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={saved} />);
    fireEvent.click(screen.getByTestId('preview-schema-btn'));
    await waitFor(() => expect(screen.getByTestId('apply-schema-btn')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('apply-schema-btn'));
  }

  it('confirms the schema with the issued operation and the revisions it was previewed against', async () => {
    mocks.apply.mockResolvedValue({
      operation_id: 'op-1', status: 'succeeded', executed_statements: 9,
      created_at: '2026-09-16T00:51:00Z', updated_at: '2026-09-16T00:51:00Z',
    });

    await previewThenApply();

    await waitFor(() => expect(mocks.apply).toHaveBeenCalledWith({
      token: 'tok-1', operation_id: 'op-1', expected_workspace_revision: 'setup-9',
      expected_plan_revision: 'plan-1', expected_connector_revision: 'identity-7',
    }));
  });

  it('shows a partial result with its next step and retries the same operation', async () => {
    mocks.apply.mockResolvedValue({
      operation_id: 'op-1', status: 'partial', executed_statements: 3,
      reason: 'statement_failed_partially_applied', next_action: 'preview again to create the missing tables',
      created_at: '2026-09-16T00:51:00Z', updated_at: '2026-09-16T00:51:00Z',
    });

    await previewThenApply();

    const status = await screen.findByTestId('recording-plan-operation-status');
    expect(status).toHaveTextContent('partial');
    expect(status).toHaveTextContent('preview again to create the missing tables');

    fireEvent.click(screen.getByTestId('apply-schema-btn'));
    await waitFor(() => expect(mocks.apply).toHaveBeenCalledTimes(2));
    expect(mocks.apply.mock.calls[1][0]).toMatchObject({ operation_id: 'op-1' });
  });

  it.each(['sqlite', 'postgres'] as const)('creates using the saved %s target and version', async (kind) => {
    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={{ ...saved, kind }} />);
    fireEvent.click(screen.getByTestId('create-plan-btn'));
    await waitFor(() => expect(mocks.create).toHaveBeenCalledOnce());
    expect(mocks.create.mock.calls[0][0].destinations).toEqual([expect.objectContaining({
      connector_id: 'saved-pg', connector_revision: 'identity-7',
    })]);
    expect(mocks.capabilities).toHaveBeenCalledWith(true, kind);
  });

  it('blocks preparation while the destination is unsaved', () => {
    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={{ ...saved, save_state: 'saving' }} />);
    expect(screen.getByTestId('create-plan-btn')).toBeDisabled();
    expect(mocks.create).not.toHaveBeenCalled();
  });

  it('previews with the selected saved target revision', async () => {
    mocks.plans = [{ id: 'plan-a', name: 'Plan A', revision: 'plan-1' }];
    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={saved} />);
    fireEvent.click(screen.getByTestId('preview-schema-btn'));
    await waitFor(() => expect(mocks.preview).toHaveBeenCalledWith(expect.objectContaining({
      plan_id: 'plan-a', connector_id: 'saved-pg', expected_connector_revision: 'identity-7',
      expected_workspace_revision: 'setup-9', expected_plan_revision: 'plan-1', dialect: 'postgres',
    })));
  });

  it('does not preview without the saved workspace setup revision', async () => {
    mocks.plans = [{ id: 'plan-a', name: 'Plan A', revision: 'plan-1' }];
    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={{ ...saved, setup_revision: undefined }} />);
    fireEvent.click(screen.getByTestId('preview-schema-btn'));
    await new Promise(resolve => setTimeout(resolve, 0));
    expect(mocks.preview).not.toHaveBeenCalled();
  });
});
