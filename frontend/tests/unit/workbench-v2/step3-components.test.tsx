import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { PipelineSteps } from '../../../src/features/datalink/workbench-v2/steps/step3/PipelineSteps';
import { PayloadPreview } from '../../../src/features/datalink/workbench-v2/steps/step3/PayloadPreview';
import { MappingRow } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingRow';
import { MappingTable } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingTable';
import { DeviceListContext } from '../../../src/features/datalink/workbench-v2/state/deviceColors';
import type { Point, Mapping, Device } from '../../../src/features/datalink/workbench-v2/state/types';
import { pointAPI, mappingAPI } from '../../../src/services/datalink';

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

describe('Step 3 UI components', () => {
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

  it('formats float64, int16, and bool pipeline values correctly', () => {
    const mappingFloat: Mapping = {
      ...mockMappings['p-01'],
      target_type: 'float64',
      scale: 0.1,
      offset: 5,
    };
    const { rerender } = render(
      <PipelineSteps point={mockPoints[0]} mapping={mappingFloat} rawValue={243} />,
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
    rerender(<PipelineSteps point={mockPoints[0]} mapping={mappingInt} rawValue={243} />);
    expect(screen.getByTestId('step-cast')).toHaveTextContent('30');
    expect(screen.getByTestId('step-final')).toHaveTextContent('30');

    const mappingBool: Mapping = {
      ...mockMappings['p-01'],
      target_type: 'bool',
      scale: 0,
      offset: 0,
    };
    rerender(<PipelineSteps point={mockPoints[0]} mapping={mappingBool} rawValue={243} />);
    expect(screen.getByTestId('step-cast')).toHaveTextContent('false');
    expect(screen.getByTestId('step-final')).toHaveTextContent('false');
  });

  it('renders payload preview JSON without calling fetch', () => {
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

  it('keeps row selection separate from inline editing and shows preview columns', () => {
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
              liveValue={243}
              onSelect={onSelect}
              dispatch={dispatch}
            />
          </tbody>
        </table>
      </DeviceListContext.Provider>,
    );

    expect(screen.getByTestId('device-live-value-p-01')).toHaveTextContent('243');
    expect(screen.getByTestId('preview-scale-p-01')).toHaveTextContent('24.30');
    expect(screen.getByTestId('preview-cast-p-01')).toHaveTextContent('24.30');
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('24.30');

    const inputTagKey = screen.getByTestId('input-tag-key-p-01');
    fireEvent.click(inputTagKey);
    expect(onSelect).not.toHaveBeenCalled();

    fireEvent.change(inputTagKey, { target: { value: 'new.tag' } });
    expect(dispatch).toHaveBeenCalledWith({
      type: 'updateMapping',
      pointId: 'p-01',
      patch: { tag_key: 'new.tag' },
    });

    fireEvent.click(screen.getByTestId('mapping-row-p-01'));
    expect(onSelect).toHaveBeenCalled();
  });

  it('shows bulk apply toolbar state and opens payload modal', async () => {
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
          rawValues={{}}
          connectionByDevice={{}}
          dispatch={dispatch}
        />
      </DeviceListContext.Provider>,
    );

    expect(screen.getByTestId('btn-bulk-apply-all-disabled')).toHaveTextContent(
      'step3.table.bulkApplyPlaceholder',
    );

    rerender(
      <DeviceListContext.Provider value={[mockDevice]}>
        <MappingTable
          points={mockPoints}
          mappings={mockMappings}
          selectedIdx={0}
          setSelectedIdx={setSelectedIdx}
          devices={[mockDevice]}
          rawValues={{}}
          connectionByDevice={{}}
          dispatch={dispatch}
        />
      </DeviceListContext.Provider>,
    );

    expect(screen.getByTestId('btn-bulk-apply-all')).toHaveTextContent(
      'step3.table.bulkApplyAll_tag_line01.temp.inlet',
    );

    fireEvent.click(screen.getByTestId('btn-bulk-apply-all'));
    expect(dispatch).toHaveBeenCalledWith({
      type: 'bulkApplyTransform',
      fromPointId: 'p-01',
      fields: ['scale', 'offset', 'target_type', 'unit'],
    });

    fireEvent.click(screen.getByTestId('payload-modal-trigger-p-01'));
    expect(await screen.findByTestId('payload-modal-content-p-01')).toBeInTheDocument();
  });
});
