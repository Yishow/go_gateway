import i18n from '../../../../src/i18n/config';
import { fireEvent, render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import QuickFlowPage from '../../../../src/features/datalink/quickflow/QuickFlowPage';
import type { Device, SourceRuleRecord } from '../../../../src/types/datalink';

const mockState = vi.hoisted(() => ({
  devices: [] as Device[],
  sourceRules: [] as SourceRuleRecord[],
  deviceError: null as Error | null,
  sourceRuleError: null as Error | null,
  runtimeState: 'connected' as 'connected' | 'connecting' | 'disconnected' | 'error',
}));

vi.mock('../../../../src/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({
    data: mockState.devices,
    isLoading: false,
    error: mockState.deviceError,
  }),
}));

vi.mock('../../../../src/hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: () => ({
    data: mockState.sourceRules,
    isLoading: false,
    error: mockState.sourceRuleError,
  }),
}));

vi.mock('../../../../src/hooks/datalink/useRuntimeStream', () => ({
  useRuntimeStream: () => ({
    connectionState: mockState.runtimeState,
    liveValues: {
      'point-40021': {
        device_id: 'device-1',
        point_id: 'point-40021',
        address: '40021',
        raw_value: 742,
        transformed_value: 74.2,
        quality: 'good',
        stale: false,
        timestamp: '2026-05-09T09:41:26Z',
      },
    },
  }),
}));

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <QuickFlowPage />
    </QueryClientProvider>,
  );
}

describe('QuickFlowPage', () => {
  beforeEach(() => {
    void i18n.changeLanguage('zh-TW');
    mockState.devices = [
      {
        id: 'device-1',
        name: '混料站 PLC-07',
        description: 'Line 7 mixer',
        protocol: 'modbus_tcp',
        status: 'active',
        connection_config: '{}',
        last_test_at: '2026-05-09T09:41:11Z',
        last_test_success: true,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
    ];
    mockState.sourceRules = [
      {
        id: 'rule-1',
        device_id: 'device-1',
        start_address: '40021',
        count: 12,
        data_type: 'int16',
        naming_prefix: 'MIX',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
        created_at: '',
        updated_at: '',
      },
    ];
    mockState.deviceError = null;
    mockState.sourceRuleError = null;
    mockState.runtimeState = 'connected';
  });

  it('renders the four-step quickflow shell with pipeline preview', () => {
    renderPage();

    expect(screen.getByRole('heading', { name: 'QuickFlow 快速導入' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /設備/ })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /探測/ })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /輸出/ })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /啟動/ })).toBeInTheDocument();
    expect(screen.getByText('即時管線預覽')).toBeInTheDocument();
    expect(screen.getByText('Modbus TCP')).toBeInTheDocument();
    expect(screen.getByText('SQLite')).toBeInTheDocument();
  });

  it('switches wizard content when selecting the points step', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: /探測/ }));

    expect(screen.getByRole('heading', { name: '探測與資料點規劃' })).toBeInTheDocument();
    expect(screen.getAllByText('Holding Register 40021-40032').length).toBeGreaterThan(0);
    expect(screen.getByText('74.2 °C')).toBeInTheDocument();
    expect(screen.getByText('PostgreSQL 欄位待補')).toBeInTheDocument();
  });

  it('shows an actionable fallback when device data fails to load', () => {
    mockState.deviceError = new Error('backend unavailable');

    renderPage();

    expect(screen.getByText('設備清單載入失敗，請重試或檢查 API 狀態。')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /重新整理設備/ })).toBeInTheDocument();
  });
});
