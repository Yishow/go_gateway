import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { PipelineSteps } from '../../../src/features/datalink/workbench-v2/steps/step3/PipelineSteps';
import { PayloadPreview } from '../../../src/features/datalink/workbench-v2/steps/step3/PayloadPreview';
import { TransformPreview } from '../../../src/features/datalink/workbench-v2/steps/step3/TransformPreview';
import { MappingRow } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingRow';
import { MappingTable } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingTable';
import { Step3Mapping } from '../../../src/features/datalink/workbench-v2/steps/step3/Step3Mapping';
import { DeviceListContext } from '../../../src/features/datalink/workbench-v2/state/deviceColors';
import type { Point, Mapping, Rule, Device, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

// Mock react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: any) => {
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

const mockPoints: Point[] = [
  {
    id: 'p-01',
    device_id: 'dev-01',
    rule_id: 'rule-01',
    rule_name: 'Holding Registers',
    name: 'SENSOR_1',
    address: '40001',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 0.1,
    _rule_offset: 0,
  },
  {
    id: 'p-02',
    device_id: 'dev-01',
    rule_id: 'rule-01',
    rule_name: 'Holding Registers',
    name: 'SENSOR_2',
    address: '40002',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 0.1,
    _rule_offset: 0,
  },
];

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
  });

  describe('PipelineSteps', () => {
    it('應能正確格式化 float64、int16、bool 三種 target_type 數值與 steps 顯示', () => {
      const mappingFloat: Mapping = {
        ...mockMappings['p-01'],
        target_type: 'float64',
        scale: 0.1,
        offset: 5,
      };
      const { rerender } = render(
        <PipelineSteps point={mockPoints[0]} mapping={mappingFloat} rawSeed={243} />
      );
      expect(screen.getByTestId('step-decode')).toHaveTextContent('243');
      expect(screen.getByTestId('step-scale')).toHaveTextContent('243 × 0.1 + 5 = 29.30');
      expect(screen.getByTestId('step-final')).toHaveTextContent('29.30');

      const mappingInt: Mapping = {
        ...mockMappings['p-01'],
        target_type: 'int16',
        scale: 0.1,
        offset: 5.6,
      };
      rerender(<PipelineSteps point={mockPoints[0]} mapping={mappingInt} rawSeed={243} />);
      expect(screen.getByTestId('step-cast')).toHaveTextContent('30');
      expect(screen.getByTestId('step-final')).toHaveTextContent('30');

      const mappingBool: Mapping = {
        ...mockMappings['p-01'],
        target_type: 'bool',
        scale: 0,
        offset: 0,
      };
      rerender(<PipelineSteps point={mockPoints[0]} mapping={mappingBool} rawSeed={243} />);
      expect(screen.getByTestId('step-cast')).toHaveTextContent('false');
      expect(screen.getByTestId('step-final')).toHaveTextContent('false');
    });
  });

  describe('PayloadPreview', () => {
    it('應渲染正確的 JSON 結構，且斷言未呼叫 fetch', () => {
      const spyFetch = vi.spyOn(window, 'fetch');
      render(<PayloadPreview point={mockPoints[0]} mapping={mockMappings['p-01']} />);

      const codeElement = screen.getByTestId('payload-preview').querySelector('code');
      expect(codeElement).not.toBeNull();
      const json = JSON.parse(codeElement!.textContent!);
      expect(json.point_id).toBe('p-01');
      expect(json.tag_id).toBe('line01.temp.inlet');
      expect(json.transform_pipeline).toHaveLength(3);
      expect(spyFetch).not.toHaveBeenCalled();
    });
  });

  describe('TransformPreview', () => {
    it('在 point/mapping 為 null 時應渲染 empty state，否則渲染 active preview 且 subtitle 同步', () => {
      const { rerender } = render(
        <TransformPreview point={null} mapping={null} rawSeed={null} />
      );
      expect(screen.getByTestId('preview-empty')).toBeInTheDocument();

      rerender(
        <TransformPreview point={mockPoints[0]} mapping={mockMappings['p-01']} rawSeed={243} />
      );
      expect(screen.getByTestId('preview-active')).toBeInTheDocument();
      expect(screen.getByText('SENSOR_1 @ 40001')).toBeInTheDocument();
    });
  });

  describe('MappingRow', () => {
    it('點擊輸入元素時不觸發 onSelect 列選取，但點擊空白處可選取，且 onChange 時 dispatch updateMapping', () => {
      const onSelect = vi.fn();
      const dispatch = vi.fn();

      render(
        <DeviceListContext.Provider value={[mockDevice]}>
          <table>
            <tbody>
              <MappingRow
                point={mockPoints[0]}
                mapping={mockMappings['p-01']}
                isSelected={false}
                devices={[mockDevice]}
                onSelect={onSelect}
                dispatch={dispatch}
              />
            </tbody>
          </table>
        </DeviceListContext.Provider>
      );

      const inputTagKey = screen.getByTestId('input-tag-key-p-01');
      fireEvent.click(inputTagKey);
      expect(onSelect).not.toHaveBeenCalled();

      fireEvent.change(inputTagKey, { target: { value: 'new.tag' } });
      expect(dispatch).toHaveBeenCalledWith({
        type: 'updateMapping',
        pointId: 'p-01',
        patch: { tag_key: 'new.tag' },
      });

      const row = screen.getByTestId('mapping-row-p-01');
      fireEvent.click(row);
      expect(onSelect).toHaveBeenCalled();
    });
  });

  describe('MappingTable', () => {
    it('在 selectedIdx 為 null 時批次套用連結呈現 disabled，否則顯示正確文字且點擊時 dispatch bulkApplyTransform', () => {
      const dispatch = vi.fn();
      const setSelectedIdx = vi.fn();

      const { rerender } = render(
        <DeviceListContext.Provider value={[mockDevice]}>
          <MappingTable
            points={mockPoints}
            mappings={mockMappings}
            selectedIdx={null}
            setSelectedIdx={setSelectedIdx}
            devices={[mockDevice]}
            dispatch={dispatch}
          />
        </DeviceListContext.Provider>
      );

      expect(screen.getByTestId('btn-bulk-apply-disabled')).toHaveTextContent(
        'step3.table.bulkApplyPlaceholder'
      );

      rerender(
        <DeviceListContext.Provider value={[mockDevice]}>
          <MappingTable
            points={mockPoints}
            mappings={mockMappings}
            selectedIdx={0}
            setSelectedIdx={setSelectedIdx}
            devices={[mockDevice]}
            dispatch={dispatch}
          />
        </DeviceListContext.Provider>
      );

      const btnApply = screen.getByTestId('btn-bulk-apply');
      expect(btnApply).toHaveTextContent('step3.table.bulkApplyText_tag_line01.temp.inlet');

      fireEvent.click(btnApply);
      expect(dispatch).toHaveBeenCalledWith({
        type: 'bulkApplyTransform',
        fromPointId: 'p-01',
        fields: ['scale', 'offset', 'target_type'],
      });
    });
  });

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
        <DeviceListContext.Provider value={[mockDevice]}>
          <Step3Mapping
            state={mockState}
            dispatch={dispatch}
            onContinue={onContinue}
            onBack={onBack}
          />
        </DeviceListContext.Provider>
      );

      expect(screen.getByTestId('step3-mapping-container')).toBeInTheDocument();
      expect(dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'initMappingsForPoints',
        })
      );
    });

    it('繼續按鈕啟用條件: A. 正常 -> 啟用; B. 留空 -> 禁用且 aside warning; C. 無啟用點位 -> 禁用且 aside warning', () => {
      const dispatch = vi.fn();
      const onContinue = vi.fn();
      const onBack = vi.fn();

      const { rerender } = render(
        <DeviceListContext.Provider value={[mockDevice]}>
          <Step3Mapping
            state={mockState}
            dispatch={dispatch}
            onContinue={onContinue}
            onBack={onBack}
          />
        </DeviceListContext.Provider>
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
        <DeviceListContext.Provider value={[mockDevice]}>
          <Step3Mapping
            state={stateWithEmpty}
            dispatch={dispatch}
            onContinue={onContinue}
            onBack={onBack}
          />
        </DeviceListContext.Provider>
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
        <DeviceListContext.Provider value={[mockDevice]}>
          <Step3Mapping
            state={stateWithAllDisabled}
            dispatch={dispatch}
            onContinue={onContinue}
            onBack={onBack}
          />
        </DeviceListContext.Provider>
      );

      expect(btnContinue).toBeDisabled();
      expect(screen.getByTestId('aside-chip-no-enabled')).toBeInTheDocument();
    });
  });
});
