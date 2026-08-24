import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import type { ReactNode } from 'react';
import { Step3Mapping } from '../../../src/features/datalink/workbench-v2/steps/step3/Step3Mapping';
import { DeviceListContext } from '../../../src/features/datalink/workbench-v2/state/deviceColors';
import type { Point, Mapping, Rule, Device, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { pointAPI, mappingAPI } from '../../../src/services/datalink';

// Mock react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: { count?: number; type?: string; tag?: string }) => {
      if (options && options.count !== undefined) {
        return `${key}_count_${options.count}`;
      }
      if (options && options.type !== undefined) {
        return `${key}_type_${options.type}`;
      }
      if (options && options.tag !== undefined) {
        return `${key}_tag_${options.tag}`;
      }
      return key;
    },
  }),
}));

vi.mock('../../../src/services/datalink', () => ({
  pointAPI: {
    list: vi.fn().mockResolvedValue([]),
  },
  mappingAPI: {
    preview: vi.fn().mockResolvedValue({
      raw_value: 243,
      final_value: 24.3,
      step_results: [
        { step_index: 1, step_type: 'decode', input_value: 243, output_value: 243, error: '' },
        { step_index: 2, step_type: 'scale', input_value: 243, output_value: 24.3, error: '' },
        { step_index: 3, step_type: 'cast', input_value: 24.3, output_value: 24.3, error: '' },
      ],
    }),
  },
}));

const mockDevice: Device = {
  id: 'dev-01',
  name: 'PLC-生產線-01',
  description: 'Modbus TCP',
  protocol: 'modbus_tcp',
  config: {},
  status: 'draft',
  test: null,
};

const mockRule: Rule = {
  id: 'rule-01',
  device_id: 'dev-01',
  name: 'Holding Registers',
  start_address: '40001',
  count: 2,
  data_type: 'int16',
  naming_prefix: 'SENSOR_',
  enabled: true,
  scale_multiplier: 0.1,
  scale_offset: 0,
  data_format: '',
  skipped_addresses: [],
  share_enabled: false,
  share_start_register: null,
  share_stride: null,
};

const mockMappings: Record<string, Mapping> = {
  'p-01': {
    point_id: 'p-01',
    tag_key: 'line01.temp.inlet',
    display_name: '進水溫度',
    unit: '°C',
    target_type: 'float64',
    scale: 0.1,
    offset: 0,
    enabled: true,
  },
  'p-02': {
    point_id: 'p-02',
    tag_key: 'line01.temp.outlet',
    display_name: '出水溫度',
    unit: '°C',
    target_type: 'float64',
    scale: 0.1,
    offset: 0,
    enabled: true,
  },
};

describe('Step 3 UI Components & Integration', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    vi.mocked(pointAPI.list).mockResolvedValue([]);
    vi.mocked(mappingAPI.preview).mockResolvedValue({
      raw_value: 243,
      final_value: 24.3,
      step_results: [
        { step_index: 1, step_type: 'decode', input_value: 243, output_value: 243, error: '' },
        { step_index: 2, step_type: 'scale', input_value: 243, output_value: 24.3, error: '' },
        { step_index: 3, step_type: 'cast', input_value: 24.3, output_value: 24.3, error: '' },
      ],
    });
  });

  function createWrapper() {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });

    return function Wrapper({ children }: { children: ReactNode }) {
      return (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
      );
    };
  }

  describe('Step3Mapping Integration', () => {
    const mockState: WorkbenchV2State = {
      ...INITIAL_STATE,
      devices: [mockDevice],
      rules: [mockRule],
      mappings: mockMappings,
    };

    it('應能完整 render Step 3，並於 mount 時 dispatch initMappingsForPoints', () => {
      const dispatch = vi.fn();
      const onContinue = vi.fn();
      const onBack = vi.fn();

      render(
        <QueryClientProvider
          client={
            new QueryClient({
              defaultOptions: {
                queries: { retry: false },
                mutations: { retry: false },
              },
            })
          }
        >
          <DeviceListContext.Provider value={[mockDevice]}>
            <Step3Mapping
              state={mockState}
              dispatch={dispatch}
              onContinue={onContinue}
              onBack={onBack}
            />
          </DeviceListContext.Provider>
        </QueryClientProvider>
      );

      expect(screen.getByTestId('step3-mapping-container')).toBeInTheDocument();
      expect(screen.queryByTestId('preview-active')).not.toBeInTheDocument();
      expect(screen.queryByTestId('preview-empty')).not.toBeInTheDocument();
      expect(dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'initMappingsForPoints',
        })
      );
    });

    it('同一 point id 的協議位址改變時，應再次初始化 mappings 並傳遞新位址', () => {
      const dispatch = vi.fn();
      const onContinue = vi.fn();
      const onBack = vi.fn();
      const queryClient = new QueryClient({
        defaultOptions: {
          queries: { retry: false },
          mutations: { retry: false },
        },
      });
      const movedDevice: Device = {
        ...mockDevice,
        protocol: 'mc_3e',
        description: 'Mitsubishi MC 3E',
      };
      const initialState: WorkbenchV2State = {
        ...mockState,
        devices: [movedDevice],
        rules: [{ ...mockRule, start_address: 'D0' }],
      };
      const movedState: WorkbenchV2State = {
        ...initialState,
        devices: [movedDevice],
        rules: [{ ...mockRule, start_address: 'D100' }],
      };

      const { rerender } = render(
        <QueryClientProvider client={queryClient}>
          <DeviceListContext.Provider value={[movedDevice]}>
            <Step3Mapping
              state={initialState}
              dispatch={dispatch}
              onContinue={onContinue}
              onBack={onBack}
            />
          </DeviceListContext.Provider>
        </QueryClientProvider>,
      );
      dispatch.mockClear();

      rerender(
        <QueryClientProvider client={queryClient}>
          <DeviceListContext.Provider value={[movedDevice]}>
            <Step3Mapping
              state={movedState}
              dispatch={dispatch}
              onContinue={onContinue}
              onBack={onBack}
            />
          </DeviceListContext.Provider>
        </QueryClientProvider>,
      );

      const initCalls = dispatch.mock.calls.filter(
        ([action]) => action?.type === 'initMappingsForPoints',
      );
      expect(initCalls).toHaveLength(1);
      expect(initCalls[0][0].points).toEqual(
        expect.arrayContaining([
          expect.objectContaining({ id: 'rule-01-p-0', device_id: 'dev-01', rule_id: 'rule-01', address: 'D100' }),
          expect.objectContaining({ id: 'rule-01-p-1', device_id: 'dev-01', rule_id: 'rule-01', address: 'D101' }),
        ]),
      );
    });

    it('只應為 enabled points 初始化 Point → Tag rows', () => {
      const dispatch = vi.fn();
      const onContinue = vi.fn();
      const onBack = vi.fn();

      const disabledRule: Rule = {
        ...mockRule,
        id: 'rule-02',
        device_id: 'dev-01',
        start_address: '40101',
        count: 4,
        naming_prefix: 'DISABLED_',
        enabled: false,
      };
      const enabledRule: Rule = {
        ...mockRule,
        id: 'rule-01',
        device_id: 'dev-01',
        start_address: '40001',
        count: 4,
        naming_prefix: 'ENABLED_',
        enabled: true,
      };
      const state: WorkbenchV2State = {
        ...INITIAL_STATE,
        devices: [mockDevice],
        rules: [enabledRule, disabledRule],
        mappings: {},
      };

      render(
        <QueryClientProvider
          client={
            new QueryClient({
              defaultOptions: {
                queries: { retry: false },
                mutations: { retry: false },
              },
            })
          }
        >
          <DeviceListContext.Provider value={[mockDevice]}>
            <Step3Mapping
              state={state}
              dispatch={dispatch}
              onContinue={onContinue}
              onBack={onBack}
            />
          </DeviceListContext.Provider>
        </QueryClientProvider>
      );

      expect(dispatch).toHaveBeenCalledWith({
        type: 'initMappingsForPoints',
        points: expect.arrayContaining([
          expect.objectContaining({ rule_id: 'rule-01', enabled: true, address: '40001' }),
          expect.objectContaining({ rule_id: 'rule-01', enabled: true, address: '40002' }),
          expect.objectContaining({ rule_id: 'rule-01', enabled: true, address: '40003' }),
          expect.objectContaining({ rule_id: 'rule-01', enabled: true, address: '40004' }),
        ]),
      });

      const initCall = dispatch.mock.calls.find(
        ([action]) => action?.type === 'initMappingsForPoints',
      )?.[0];
      expect(initCall.points).toHaveLength(4);
      expect(initCall.points.every((point: Point) => point.enabled)).toBe(true);
      expect(initCall.points.some((point: Point) => point.rule_id === 'rule-02')).toBe(false);
    });

    it('會為多組規則產生穩定且不重複的預設 Tag Key', () => {
      const dispatch = vi.fn();

      const state: WorkbenchV2State = {
        ...INITIAL_STATE,
        devices: [mockDevice],
        rules: [
          {
            ...mockRule,
            id: 'rule-01',
            count: 8,
            naming_prefix: 'LINE_A_',
            start_address: '40001',
          },
          {
            ...mockRule,
            id: 'rule-02',
            count: 8,
            naming_prefix: 'LINE_B_',
            start_address: '40101',
          },
        ],
        mappings: {},
      };

      render(
        <QueryClientProvider
          client={
            new QueryClient({
              defaultOptions: {
                queries: { retry: false },
                mutations: { retry: false },
              },
            })
          }
        >
          <DeviceListContext.Provider value={[mockDevice]}>
            <Step3Mapping
              state={state}
              dispatch={dispatch}
              onContinue={vi.fn()}
              onBack={vi.fn()}
            />
          </DeviceListContext.Provider>
        </QueryClientProvider>
      );

      const initCall = dispatch.mock.calls.find(
        ([action]) => action?.type === 'initMappingsForPoints',
      )?.[0];

      expect(initCall.points).toHaveLength(16);
      expect(initCall.points[0]).toEqual(
        expect.objectContaining({ address: '40001', name: 'LINE_A_0' }),
      );
      expect(initCall.points[8]).toEqual(
        expect.objectContaining({ address: '40101', name: 'LINE_B_0' }),
      );
    });

    it('繼續按鈕啟用條件: A. 正常 -> 啟用; B. 留空 -> 禁用且 aside warning; C. 無啟用點位 -> 禁用且 aside warning', () => {
      const dispatch = vi.fn();
      const onContinue = vi.fn();
      const onBack = vi.fn();

      const Wrapper = createWrapper();

      const { rerender } = render(
        <Wrapper>
          <DeviceListContext.Provider value={[mockDevice]}>
            <Step3Mapping
              state={mockState}
              dispatch={dispatch}
              onContinue={onContinue}
              onBack={onBack}
            />
          </DeviceListContext.Provider>
        </Wrapper>
      );

      const btnContinue = screen.getByTestId('btn-continue');
      expect(btnContinue).not.toBeDisabled();
      expect(screen.getByTestId('aside-chip-success')).toBeInTheDocument();

      const stateWithEmpty: WorkbenchV2State = {
        ...mockState,
        mappings: {
          ...mockMappings,
          'p-01': {
            ...mockMappings['p-01'],
            tag_key: '  ',
          },
        },
      };

      rerender(
        <Wrapper>
          <DeviceListContext.Provider value={[mockDevice]}>
            <Step3Mapping
              state={stateWithEmpty}
              dispatch={dispatch}
              onContinue={onContinue}
              onBack={onBack}
            />
          </DeviceListContext.Provider>
        </Wrapper>
      );

      expect(btnContinue).toBeDisabled();
      expect(screen.getByTestId('aside-chip-empty-tag')).toHaveTextContent(
        'step3.footer.emptyTagWarning_count_1'
      );

      const stateWithAllDisabled: WorkbenchV2State = {
        ...mockState,
        mappings: {
          'p-01': { ...mockMappings['p-01'], enabled: false },
          'p-02': { ...mockMappings['p-02'], enabled: false },
        },
      };

      rerender(
        <Wrapper>
          <DeviceListContext.Provider value={[mockDevice]}>
            <Step3Mapping
              state={stateWithAllDisabled}
              dispatch={dispatch}
              onContinue={onContinue}
              onBack={onBack}
            />
          </DeviceListContext.Provider>
        </Wrapper>
      );

      expect(btnContinue).toBeDisabled();
      expect(screen.getByTestId('aside-chip-no-enabled')).toBeInTheDocument();
    });

    it('支援一鍵啟用或停用所有映射點位', () => {
      const dispatch = vi.fn();

      render(
        <QueryClientProvider
          client={
            new QueryClient({
              defaultOptions: {
                queries: { retry: false },
                mutations: { retry: false },
              },
            })
          }
        >
          <DeviceListContext.Provider value={[mockDevice]}>
            <Step3Mapping
              state={mockState}
              dispatch={dispatch}
              onContinue={vi.fn()}
              onBack={vi.fn()}
            />
          </DeviceListContext.Provider>
        </QueryClientProvider>
      );

      dispatch.mockClear();
      fireEvent.click(screen.getByTestId('btn-disable-all-mappings'));
      fireEvent.click(screen.getByTestId('btn-enable-all-mappings'));

      expect(dispatch).toHaveBeenNthCalledWith(1, {
        type: 'setAllMappingsEnabled',
        enabled: false,
      });
      expect(dispatch).toHaveBeenNthCalledWith(2, {
        type: 'setAllMappingsEnabled',
        enabled: true,
      });
    });
  });
});
