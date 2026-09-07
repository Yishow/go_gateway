import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { RecordingPlanSetupSection } from '@/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection';

const mutateCreateMock = vi.fn();
const mutatePreviewMock = vi.fn();
const mutateApplyMock = vi.fn();
const mutateTestWriteMock = vi.fn();
const refetchMock = vi.fn();

let mockPlans: unknown[] = [];

vi.mock('@/hooks/datalink/useStudioV2WorkspaceRecordingPlans', () => ({
  useStudioV2WorkspaceRecordingPlansQuery: () => ({
    data: mockPlans,
    isLoading: false,
    refetch: refetchMock,
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
}));

describe('RecordingPlanSetupSection Component', () => {
  beforeEach(() => {
    mutateCreateMock.mockReset();
    mutatePreviewMock.mockReset();
    mutateApplyMock.mockReset();
    mutateTestWriteMock.mockReset();
    refetchMock.mockReset();
    mockPlans = [];
  });

  it('renders create plan form when no plan exists and allows creation', async () => {
    render(<RecordingPlanSetupSection deviceId="dev-1" />);

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
    mockPlans = [
      {
        id: 'plan-existing',
        name: '生產記錄方案',
        revision: 'rev-1',
      },
    ];

    render(<RecordingPlanSetupSection deviceId="dev-1" />);

    // 1. 預覽 DDL
    mutatePreviewMock.mockResolvedValueOnce({
      token: 'tok-1',
      statements: ['CREATE TABLE gw_record_samples ...'],
    });

    const previewBtn = screen.getByTestId('preview-schema-btn');
    fireEvent.click(previewBtn);

    await waitFor(() => {
      expect(screen.getByTestId('apply-schema-btn')).toBeInTheDocument();
      expect(screen.getByText('CREATE TABLE gw_record_samples ...')).toBeInTheDocument();
    });

    // 2. 套用 DDL
    mutateApplyMock.mockResolvedValueOnce({ applied: true });
    fireEvent.click(screen.getByTestId('apply-schema-btn'));

    await waitFor(() => {
      expect(mutateApplyMock).toHaveBeenCalledWith('tok-1');
    });

    // 3. 執行試寫
    mutateTestWriteMock.mockResolvedValueOnce({
      status: 'success',
      table: 'gw_record_samples',
      record_id: 'rec-test-1',
      observed_at: '2026-09-07T10:00:00Z',
    });

    const testWriteBtn = screen.getByTestId('test-write-btn');
    fireEvent.click(testWriteBtn);

    await waitFor(() => {
      expect(screen.getByTestId('test-write-result')).toBeInTheDocument();
      expect(screen.getByText(/成功 \(已驗證回讀且清理乾淨\)/)).toBeInTheDocument();
    });
  });
});
