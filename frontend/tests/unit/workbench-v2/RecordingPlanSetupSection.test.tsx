import { savedRecordingScope, savedRecordingPlanScopeFields } from '../../fixtures/recordingScope';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { RecordingPlanSetupSection } from '@/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection';
import { RecordingPlanResponseError } from '@/utils/recordingPlanJson';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
const savedConnector = { ...INITIAL_STATE.db.connector, kind: 'sqlite' as const, workspace_id: 'workspace-a', connector_id: 'saved-sqlite', identity_revision: 'identity-1', setup_revision: 'setup-1', persisted: true, save_state: 'saved' as const };

/** A preview the backend would issue: scope, revisions, operation and digest included. */
function previewToken(token: string, statements: string[]) {
  return {
    token, operation_id: `op-${token}`, action: 'schema_apply', workspace_id: 'workspace-a',
    workspace_revision: 'setup-1', plan_id: 'plan-existing', plan_revision: 'rev-1',
    connector_id: 'saved-sqlite', connector_revision: 'identity-1', dialect: 'sqlite',
    table_prefix: 'gw_record_', statements, tables: [], digest: 'a'.repeat(64),
    expires_at: '2026-09-16T01:00:00Z', created_at: '2026-09-16T00:50:00Z',
  };
}

/** The confirmation the client must send back for that preview. */
function confirmationFor(token: string) {
  return {
    token, operation_id: `op-${token}`, expected_workspace_revision: 'setup-1',
    expected_plan_revision: 'rev-1', expected_connector_revision: 'identity-1',
  };
}

const mutateCreateMock = vi.fn();
const mutatePreviewMock = vi.fn();
const mutateApplyMock = vi.fn();
const mutateTestWriteMock = vi.fn();
const refetchMock = vi.fn();
const operationRefetchMock = vi.fn();
let checkedOperation: Record<string, unknown> | undefined;

let mockPlans: unknown[] = [];
let mockCapabilities: unknown[] = [];
const supportedSqliteCapability = {
  kind: 'sqlite',
  supported: true,
  supports_managed_schema: true,
  supports_transactions: true,
  supports_receipts: true,
  supports_test_writes: true,
  supported_modes: ['managed_recording', 'custom_table'],
};

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) => options?.defaultValue ?? key,
  }),
}));

vi.mock('@/hooks/datalink/useStudioV2WorkspaceRecordingPlans', () => ({
  useStudioV2WorkspaceRecordingPlansQuery: () => ({
    data: mockPlans.map(plan => ({ ...savedRecordingPlanScopeFields, ...(plan as Record<string, unknown>) })),
    isLoading: false,
    refetch: refetchMock,
  }),
  useStudioV2ConnectorCapabilitiesQuery: () => ({
    data: mockCapabilities,
    isLoading: false,
  }),
  useCreateRecordingPlanMutation: () => ({
    mutateAsync: mutateCreateMock,
    isPending: false,
  }),
  usePreviewSchemaMutation: () => ({
    mutateAsync: mutatePreviewMock,
    isPending: false,
  }),
  useApplySchemaMutation: () => ({
    mutateAsync: mutateApplyMock,
    isPending: false,
  }),
  useTestWritePlanMutation: () => ({
    mutateAsync: mutateTestWriteMock,
    isPending: false,
  }),
  useSchemaOperationQuery: () => ({
    data: checkedOperation,
    isError: false,
    isFetching: false,
    refetch: operationRefetchMock,
  }),
}));

describe('RecordingPlanSetupSection Component', () => {
  beforeEach(() => {
    mutateCreateMock.mockReset();
    mutatePreviewMock.mockReset();
    mutateApplyMock.mockReset();
    mutateTestWriteMock.mockReset();
    refetchMock.mockReset();
    operationRefetchMock.mockReset();
    checkedOperation = undefined;
    mockPlans = [];
    mockCapabilities = [{
      kind: 'sqlite',
      supported: true,
      supports_managed_schema: false,
      supports_transactions: true,
      supports_receipts: true,
      supports_test_writes: false,
      supported_modes: ['managed_recording', 'custom_table'],
    }];
  });

  it('renders create plan form when no plan exists and allows creation', async () => {
    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);

    expect(screen.getByTestId('recording-plan-setup-section')).toBeInTheDocument();
    const createBtn = screen.getByTestId('create-plan-btn');

    mutateCreateMock.mockResolvedValueOnce({ id: 'plan-new' });
    fireEvent.click(createBtn);

    await waitFor(() => {
      expect(mutateCreateMock).toHaveBeenCalledWith(
        expect.objectContaining({
          timezone: 'Asia/Taipei',
        })
      );
    });
  });

  it('previews schema, applies schema and executes test write when plan exists', async () => {
    mockCapabilities = [supportedSqliteCapability];
    mockPlans = [
      {
        id: 'plan-existing',
        name: '生產記錄方案',
        revision: 'rev-1',
      },
    ];

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);

    // 1. 預覽 DDL
    mutatePreviewMock.mockResolvedValueOnce(previewToken('tok-1', ['CREATE TABLE gw_record_samples ...']));

    const previewBtn = screen.getByTestId('preview-schema-btn');
    fireEvent.click(previewBtn);

    await waitFor(() => {
      expect(screen.getByTestId('apply-schema-btn')).toBeInTheDocument();
      expect(screen.getByText('CREATE TABLE gw_record_samples ...')).toBeInTheDocument();
    });

    // 2. 套用 DDL
    mutateApplyMock.mockResolvedValueOnce({
      operation_id: 'op-tok-1', status: 'succeeded', executed_statements: 1,
      created_at: '2026-09-16T00:51:00Z', updated_at: '2026-09-16T00:51:00Z',
    });
    fireEvent.click(screen.getByTestId('apply-schema-btn'));

    await waitFor(() => {
      expect(mutateApplyMock).toHaveBeenCalledWith(confirmationFor('tok-1'));
    });

    // 3. 執行試寫
    mutateTestWriteMock.mockResolvedValueOnce({
      status: 'written_verified',
      table: 'gw_record_samples',
      record_id: 'rec-test-1',
      observed_at: '2026-09-07T10:00:00Z',
      delivered_at: '2026-09-07T10:00:00.012Z',
      message: 'raw backend diagnostic must stay out of the UI',
    });

    const testWriteBtn = screen.getByTestId('test-write-btn');
    fireEvent.click(testWriteBtn);

    await waitFor(() => {
      expect(screen.getByTestId('test-write-result')).toBeInTheDocument();
      expect(screen.getByText('step4.recording_written_verified')).toBeInTheDocument();
    });
    expect(screen.queryByText(/清理/)).not.toBeInTheDocument();
    expect(screen.queryByText(/raw backend diagnostic/)).not.toBeInTheDocument();
  });

  it('keeps preview usable while capability flags block apply and test write', async () => {
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutatePreviewMock.mockResolvedValueOnce(previewToken('tok-capability', ['CREATE TABLE gw_record_samples ...']));

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    fireEvent.click(screen.getByTestId('preview-schema-btn'));

    await waitFor(() => {
      expect(screen.getByText('CREATE TABLE gw_record_samples ...')).toBeInTheDocument();
    });
    expect(screen.getByTestId('apply-schema-btn')).toBeDisabled();
    expect(screen.getByTestId('test-write-btn')).toBeDisabled();
    expect(screen.getByTestId('recording-plan-capability-warning')).toHaveTextContent('step4.recording_operations_unavailable');
  });

  it.each([
    [
      'managed schema only',
      { supports_managed_schema: true, supports_test_writes: false },
      false,
      true,
    ],
    [
      'test writes only',
      { supports_managed_schema: false, supports_test_writes: true },
      true,
      false,
    ],
  ] as const)('gates each recording operation independently for %s', async (
    _label,
    flags,
    applyDisabled,
    testWriteDisabled,
  ) => {
    mockCapabilities = [{ ...supportedSqliteCapability, ...flags }];
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutatePreviewMock.mockResolvedValueOnce(previewToken('tok-independent-capability', ['CREATE TABLE gw_record_samples ...']));

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    fireEvent.click(screen.getByTestId('preview-schema-btn'));

    await waitFor(() => expect(screen.getByTestId('apply-schema-btn')).toBeInTheDocument());
    expect(screen.getByTestId('apply-schema-btn')).toHaveProperty('disabled', applyDisabled);
    expect(screen.getByTestId('test-write-btn')).toHaveProperty('disabled', testWriteDisabled);
  });

  it('reports a failed operation with its next step and keeps the preview', async () => {
    mockCapabilities = [supportedSqliteCapability];
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutatePreviewMock.mockResolvedValueOnce(previewToken('tok-apply-failed', ['CREATE TABLE x']));
    mutateApplyMock.mockResolvedValueOnce({
      operation_id: 'op-tok-apply-failed', status: 'failed', executed_statements: 0,
      reason: 'statement_failed_rolled_back', next_action: 'fix the reported cause, then preview again',
      created_at: '2026-09-16T00:51:00Z', updated_at: '2026-09-16T00:51:00Z',
    });

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    fireEvent.click(screen.getByTestId('preview-schema-btn'));
    await waitFor(() => expect(screen.getByTestId('apply-schema-btn')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('apply-schema-btn'));

    const status = await screen.findByTestId('recording-plan-operation-status');
    expect(status).toHaveTextContent('failed');
    expect(status).toHaveTextContent('fix the reported cause, then preview again');
    expect(screen.getByText('CREATE TABLE x')).toBeInTheDocument();
    expect(refetchMock).not.toHaveBeenCalled();
  });

  it('preserves preview after a typed 501 apply failure and does not expose diagnostics', async () => {
    mockCapabilities = [supportedSqliteCapability];
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutatePreviewMock.mockResolvedValueOnce(previewToken('tok-501', ['CREATE TABLE gw_record_samples ...']));
    mutateApplyMock.mockRejectedValueOnce({
      response: {
        status: 501,
        data: {
          success: false,
          error: {
            code: 'RECORDING_SCHEMA_NOT_IMPLEMENTED',
            message: 'raw managed schema diagnostic secret',
            action: 'wait_for_supported_operation',
            retryable: false,
            request_id: 'req-schema-501',
          },
        },
      },
    });

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    fireEvent.click(screen.getByTestId('preview-schema-btn'));
    await waitFor(() => expect(screen.getByTestId('apply-schema-btn')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('apply-schema-btn'));

    await waitFor(() => {
      expect(mutateApplyMock).toHaveBeenCalledWith(confirmationFor('tok-501'));
      expect(screen.getByText('CREATE TABLE gw_record_samples ...')).toBeInTheDocument();
      expect(screen.getByTestId('recording-plan-action-error')).toHaveTextContent('errors.generic_failure');
    });
    expect(screen.getByTestId('recording-plan-action-error')).not.toHaveTextContent('raw managed schema diagnostic secret');
    expect(refetchMock).not.toHaveBeenCalled();
  });

  it('lets an unconfirmed apply be checked instead of retried blindly', async () => {
    mockCapabilities = [supportedSqliteCapability];
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutatePreviewMock.mockResolvedValueOnce(previewToken('tok-unconfirmed-check', ['CREATE TABLE x']));
    mutateApplyMock.mockRejectedValueOnce({
      response: {
        status: 503,
        data: {
          success: false,
          error: {
            code: 'RECORDING_SCHEMA_OPERATION_UNKNOWN',
            message: 'raw diagnostic must stay out of the UI',
            retryable: true,
            operation_id: 'op-tok-unconfirmed-check',
          },
        },
      },
    });

    const { rerender } = render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    fireEvent.click(screen.getByTestId('preview-schema-btn'));
    await waitFor(() => expect(screen.getByTestId('apply-schema-btn')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('apply-schema-btn'));
    await waitFor(() => expect(screen.getByTestId('recording-plan-operation-unconfirmed')).toBeInTheDocument());

    fireEvent.click(screen.getByTestId('recording-plan-operation-unconfirmed-check'));
    await waitFor(() => expect(operationRefetchMock).toHaveBeenCalledTimes(1));

    checkedOperation = {
      operation_id: 'op-tok-unconfirmed-check', status: 'succeeded', executed_statements: 2,
      created_at: '2026-09-16T01:20:00Z', updated_at: '2026-09-16T01:20:05Z',
    };
    rerender(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    expect(screen.getByTestId('recording-plan-operation-unconfirmed-result')).toHaveTextContent('succeeded');
  });

  it('turns malformed apply responses into unconfirmed while preserving preview', async () => {
    mockCapabilities = [supportedSqliteCapability];
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutatePreviewMock.mockResolvedValueOnce(previewToken('tok-malformed-apply', ['CREATE TABLE x']));
    mutateApplyMock.mockRejectedValueOnce(new RecordingPlanResponseError('invalid apply response'));

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    fireEvent.click(screen.getByTestId('preview-schema-btn'));
    await waitFor(() => expect(screen.getByTestId('apply-schema-btn')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('apply-schema-btn'));

    await waitFor(() => {
      expect(screen.getByTestId('recording-plan-operation-unconfirmed')).toHaveTextContent('step4.recording_operation_unconfirmed');
      expect(screen.getByText('CREATE TABLE x')).toBeInTheDocument();
    });
    expect(refetchMock).not.toHaveBeenCalled();
  });

  it('shows an explicit unconfirmed result after transport loss and never resubmits', async () => {
    mockCapabilities = [supportedSqliteCapability];
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutateTestWriteMock.mockRejectedValueOnce({ message: 'network secret diagnostic', operation_id: 'op-transport-1' });

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    fireEvent.click(screen.getByTestId('test-write-btn'));

    await waitFor(() => {
      expect(screen.getByTestId('recording-plan-test-unconfirmed')).toHaveTextContent('step4.recording_test_unconfirmed');
      expect(screen.getByTestId('recording-plan-test-unconfirmed')).toHaveTextContent('op-transport-1');
    });
    expect(screen.getByTestId('recording-plan-test-unconfirmed')).not.toHaveTextContent('network secret diagnostic');
    expect(mutateTestWriteMock).toHaveBeenCalledTimes(1);
  });

  it('does not display green success for unfamiliar or incomplete nominal results', async () => {
    mockCapabilities = [supportedSqliteCapability];
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutateTestWriteMock.mockResolvedValueOnce({ status: 'success', record_id: 'rec-1' });

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    fireEvent.click(screen.getByTestId('test-write-btn'));

    await waitFor(() => expect(screen.getByTestId('recording-plan-test-unconfirmed')).toBeInTheDocument());
    expect(screen.queryByText('step4.recording_written_verified')).not.toBeInTheDocument();
    expect(screen.queryByTestId('test-write-result')).not.toBeInTheDocument();
  });

  it.each([
    ['written_verified', 'step4.recording_written_verified'],
    ['written_unverified', 'step4.recording_written_unverified'],
    ['failed', 'step4.recording_write_failed'],
    ['unknown', 'step4.recording_write_unknown'],
  ] as const)('renders truthful copy for %s results', async (status, expectedCopy) => {
    mockCapabilities = [supportedSqliteCapability];
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutateTestWriteMock.mockResolvedValueOnce({
      status,
      record_id: 'rec-result-1',
      table: 'gw_record_samples',
      observed_at: '2026-09-07T10:00:00Z',
      delivered_at: '2026-09-07T10:00:00.012Z',
    });

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    fireEvent.click(screen.getByTestId('test-write-btn'));

    await waitFor(() => expect(screen.getByTestId('test-write-result')).toHaveTextContent(expectedCopy));
    expect(screen.getByTestId('test-write-result')).not.toHaveTextContent('清理乾淨');
    if (status === 'unknown') {
      expect(screen.getByTestId('test-write-result')).not.toHaveTextContent('step4.recording_write_failed');
      expect(screen.getByTestId('test-write-result')).not.toHaveTextContent('已寫入且回讀已驗證');
    }
  });

  it('clears a stale verified result before an uncertain follow-up attempt', async () => {
    mockCapabilities = [supportedSqliteCapability];
    mockPlans = [{ id: 'plan-existing', name: '生產記錄方案', revision: 'rev-1' }];
    mutateTestWriteMock
      .mockResolvedValueOnce({
        status: 'written_verified',
        record_id: 'rec-first',
        table: 'gw_record_samples',
        observed_at: '2026-09-07T10:00:00Z',
        delivered_at: '2026-09-07T10:00:00.012Z',
      })
      .mockRejectedValueOnce({ message: 'transport secret', operation_id: 'op-follow-up' });

    render(<RecordingPlanSetupSection scope={savedRecordingScope} connector={savedConnector} />);
    const testButton = screen.getByTestId('test-write-btn');
    fireEvent.click(testButton);
    await waitFor(() => expect(screen.getByText('step4.recording_written_verified')).toBeInTheDocument());
    fireEvent.click(testButton);

    await waitFor(() => expect(screen.getByTestId('recording-plan-test-unconfirmed')).toBeInTheDocument());
    expect(screen.queryByText('step4.recording_written_verified')).not.toBeInTheDocument();
  });
});
