import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { HistoryReportsPanel } from '@/features/datalink/runtime-dashboard/HistoryReportsPanel';
import type { HistoryReport } from '@/types/historyReport';

const mutateAsyncExportMock = vi.fn();
const refetchMock = vi.fn();

let mockPlans: unknown[] = [];
let mockHistoryReport: HistoryReport | null = null;
let mockIsLoadingHistory = false;

vi.mock('@/hooks/datalink/useStudioV2WorkspaceRecordingPlans', () => ({
  useStudioV2WorkspaceRecordingPlansQuery: () => ({
    data: mockPlans,
    isLoading: false,
  }),
}));

vi.mock('@/hooks/datalink/useStudioV2WorkspaceHistory', () => ({
  useStudioV2WorkspaceHistoryQuery: () => ({
    data: mockHistoryReport,
    isLoading: mockIsLoadingHistory,
    refetch: refetchMock,
  }),
  useExportHistoryCSVMutation: () => ({
    mutateAsync: mutateAsyncExportMock,
    isPending: false,
  }),
}));

describe('HistoryReportsPanel', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockPlans = [];
    mockHistoryReport = null;
    mockIsLoadingHistory = false;
  });

  it('renders empty prompt when no recording plan exists for device', () => {
    mockPlans = [];
    render(<HistoryReportsPanel selectedDeviceId="dev-1" />);

    expect(screen.getByTestId('history-reports-empty')).toBeInTheDocument();
    expect(screen.getByText(/尚未建立本設備之記錄方案/)).toBeInTheDocument();
  });

  it('renders history summary cards and table when plan and points exist', async () => {
    mockPlans = [
      {
        id: 'plan-energy-1',
        name: '產線電力記錄方案',
      },
    ];
    mockHistoryReport = {
      workspace_id: 'ws-1',
      plan_id: 'plan-energy-1',
      points: [
        {
          measurement_id: 'meas-kw-1',
          observed_at: '2026-09-07T10:00:00Z',
          value_numeric: 120.5,
          time_weighted_mean: 120.5,
          sampled_min: 110.0,
          sampled_max: 130.0,
          usage_delta: 2.5,
          quality: 'good',
          coverage_ratio: 0.98,
          is_estimated: false,
          is_provisional: false,
        },
      ],
      total_count: 1,
      generated_at: '2026-09-07T10:05:00Z',
    };

    render(<HistoryReportsPanel selectedDeviceId="dev-1" />);

    expect(screen.getByTestId('history-reports-panel')).toBeInTheDocument();
    expect(screen.getByText(/方案: 產線電力記錄方案/)).toBeInTheDocument();
    expect(screen.getByTestId('summary-count')).toHaveTextContent('1');
    expect(screen.getByTestId('summary-avg')).toHaveTextContent('120.50');
    expect(screen.getByTestId('summary-usage')).toHaveTextContent('2.50');
    expect(screen.getByTestId('summary-coverage')).toHaveTextContent('98.0%');

    // Table checks
    expect(screen.getByText('meas-kw-1')).toBeInTheDocument();
    expect(screen.getByText('+2.50')).toBeInTheDocument();
  });

  it('triggers CSV export when clicking export button', async () => {
    mockPlans = [
      {
        id: 'plan-energy-1',
        name: '產線電力記錄方案',
      },
    ];
    mockHistoryReport = {
      workspace_id: 'ws-1',
      plan_id: 'plan-energy-1',
      points: [],
      total_count: 0,
      generated_at: '2026-09-07T10:05:00Z',
    };

    mutateAsyncExportMock.mockResolvedValueOnce(new Blob(['test,csv'], { type: 'text/csv' }));
    // Mock URL.createObjectURL and revokeObjectURL
    window.URL.createObjectURL = vi.fn().mockReturnValue('blob:test');
    window.URL.revokeObjectURL = vi.fn();

    render(<HistoryReportsPanel selectedDeviceId="dev-1" />);

    const exportBtn = screen.getByTestId('history-export-csv-btn');
    fireEvent.click(exportBtn);

    await waitFor(() => {
      expect(mutateAsyncExportMock).toHaveBeenCalledTimes(1);
    });
  });
});
