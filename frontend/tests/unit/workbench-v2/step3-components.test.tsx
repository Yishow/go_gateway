import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { PipelineSteps } from '../../../src/features/datalink/workbench-v2/steps/step3/PipelineSteps';
import { PayloadPreview } from '../../../src/features/datalink/workbench-v2/steps/step3/PayloadPreview';
import { MappingRow } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingRow';
import { MappingTable } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingTable';
import { DeviceListContext } from '../../../src/features/datalink/workbench-v2/state/deviceColors';
import type { Point, Mapping, Device, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { pointAPI, mappingAPI } from '../../../src/services/datalink';
import { INITIAL_STATE, type WorkbenchV2Action } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { mappingReducer } from '../../../src/features/datalink/workbench-v2/state/mappingReducer';

vi.mock('react-i18next', () => ({
  useTranslation: (namespace?: string) => ({
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
      return namespace === 'mapping-errors' ? `${namespace}:${key}` : key;
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
    expect(screen.getByTestId('preview-scale-p-01')).toHaveTextContent('--');
    expect(screen.getByTestId('preview-cast-p-01')).toHaveTextContent('--');
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('--');
    expect(screen.queryByTestId('input-tag-key-p-01')).not.toBeInTheDocument();
    const inputDisplayName = screen.getByTestId('input-display-name-p-01');
    fireEvent.click(inputDisplayName);
    expect(onSelect).not.toHaveBeenCalled();

    fireEvent.change(inputDisplayName, { target: { value: '新的顯示名稱' } });
    expect(dispatch).toHaveBeenCalledWith({
      type: 'updateMapping',
      pointId: 'p-01',
      patch: { display_name: '新的顯示名稱' },
    });

    fireEvent.click(screen.getByTestId('mapping-row-p-01'));
    expect(onSelect).toHaveBeenCalled();
  });

  it('shows typed runtime recovery metadata without exposing backend detail', () => {
    render(
      <DeviceListContext.Provider value={[mockDevice]}>
        <table>
          <tbody>
            <MappingRow
              point={mockPoints[0]}
              mapping={mockMappings['p-01']}
              isSelected={false}
              devices={[mockDevice]}
              liveValue={undefined}
              connectionState="degraded"
              streamRecovery={{
                code: 'runtime_stream_unavailable',
                action: 'retry the runtime stream',
                requestId: 'runtime-request-9',
                retryable: true,
              }}
              onSelect={vi.fn()}
              dispatch={vi.fn()}
            />
          </tbody>
        </table>
      </DeviceListContext.Provider>,
    );

    expect(screen.getByTestId('runtime-recovery-p-01')).toHaveTextContent(
      'errors.runtime_stream_unavailable',
    );
    expect(screen.getByTestId('runtime-recovery-p-01')).toHaveTextContent(
      'step3.recovery.retryAction',
    );
    expect(screen.getByTestId('runtime-recovery-request-id-p-01')).toHaveTextContent('runtime-request-9');
    expect(screen.getByTestId('runtime-recovery-p-01')).not.toHaveTextContent('retry the runtime stream');
  });

  it('localizes mapping row labels and hides raw save errors', () => {
    render(
      <DeviceListContext.Provider value={[]}>
        <table>
          <tbody>
            <MappingRow
              point={mockPoints[0]}
              mapping={{ ...mockMappings['p-01'], save_state: 'save-error', save_error: 'secret backend details' }}
              isSelected={false}
              devices={[]}
              onSelect={vi.fn()}
              dispatch={vi.fn()}
            />
          </tbody>
        </table>
      </DeviceListContext.Provider>,
    );

    expect(screen.getByText('step3.unknownDevice')).toBeInTheDocument();
    expect(screen.queryByPlaceholderText('step3.placeholders.tagKey')).not.toBeInTheDocument();
    expect(screen.getByPlaceholderText('step3.placeholders.displayName')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('step3.placeholders.unit')).toBeInTheDocument();
    expect(screen.getByTestId('mapping-save-state-p-01')).toHaveTextContent('step3.saveStates.save-error');
    expect(screen.getByTestId('mapping-save-state-p-01')).toHaveTextContent('mapping-errors:generic');
    expect(screen.getByTestId('mapping-save-state-p-01')).not.toHaveTextContent('secret backend details');
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
      'step3.table.bulkApplyAll',
    );

    fireEvent.click(screen.getByTestId('bulk-target-p-02'));
    fireEvent.click(screen.getByTestId('btn-bulk-apply-all'));
    expect(dispatch).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('bulk-confirm'));
    expect(dispatch).toHaveBeenCalledWith({
      type: 'bulkApplyTransform',
      fromPointId: 'p-01',
      toPointIds: ['p-02'],
      fields: ['scale', 'offset', 'target_type', 'unit'],
    });

    fireEvent.click(screen.getByTestId('mapping-advanced-toggle'));
    fireEvent.click(screen.getByTestId('payload-modal-trigger-p-01'));
    expect(await screen.findByTestId('payload-modal-content-p-01')).toBeInTheDocument();
  });
});


it('invalidates a bulk preview after a source edit and allows cancellation without dispatch', () => {
  const dispatch = vi.fn();
  const props = { points: mockPoints, mappings: mockMappings, selectedIdx: 0, setSelectedIdx: vi.fn(), devices: [mockDevice], rawValues: {}, connectionByDevice: {}, dispatch };
  const { rerender } = render(<MappingTable {...props} />);
  fireEvent.click(screen.getByTestId('bulk-target-p-02'));
  fireEvent.click(screen.getByTestId('btn-bulk-apply-unit'));
  expect(screen.getByTestId('bulk-preview')).toHaveTextContent('40002');
  rerender(<MappingTable {...props} mappings={{ ...mockMappings, 'p-01': { ...mockMappings['p-01'], unit: 'bar' } }} />);
  expect(screen.queryByTestId('bulk-confirm')).not.toBeInTheDocument();
  expect(dispatch).not.toHaveBeenCalled();
});

it('shows five basic columns and keeps editor nodes and values when expanding Advanced', () => {
  render(<MappingTable points={mockPoints} mappings={mockMappings} selectedIdx={0} setSelectedIdx={vi.fn()} devices={[mockDevice]} rawValues={{}} connectionByDevice={{}} dispatch={vi.fn()} />);
  expect(screen.getAllByRole('columnheader')).toHaveLength(5);
  const scale = screen.getByTestId('input-scale-p-01');
  expect(scale).not.toBeVisible();
  fireEvent.click(screen.getByTestId('mapping-advanced-toggle'));
  expect(screen.getAllByRole('columnheader')).toHaveLength(12);
  expect(scale).toBeVisible();
  expect(scale).toHaveValue(0.1);
  fireEvent.click(screen.getByTestId('mapping-advanced-toggle'));
  expect(scale).not.toBeVisible();
  expect(screen.getByTestId('input-scale-p-01')).toBe(scale);
  expect(screen.getByTestId('preview-final-p-01')).toBeVisible();
});


it('keeps the explicit bulk source when enabled or target checkboxes receive native focus', () => {
  const setSelectedIdx = vi.fn();
  const dispatch = vi.fn();
  render(<MappingTable points={mockPoints} mappings={mockMappings} selectedIdx={0} setSelectedIdx={setSelectedIdx} devices={[mockDevice]} rawValues={{}} connectionByDevice={{}} dispatch={dispatch} />);
  screen.getByTestId('input-display-name-p-02').focus();
  expect(setSelectedIdx).toHaveBeenCalledWith(1);
  setSelectedIdx.mockClear();
  screen.getByTestId('bulk-target-p-02').focus();
  expect(setSelectedIdx).not.toHaveBeenCalled();
  const enabled = screen.getByTestId('checkbox-enabled-p-02');
  enabled.focus();
  expect(setSelectedIdx).not.toHaveBeenCalled();
  fireEvent.click(enabled);
  expect(dispatch).toHaveBeenCalledWith({ type: 'toggleMappingEnabled', pointId: 'p-02' });
  expect(setSelectedIdx).not.toHaveBeenCalled();
});

it('reviews and confirms target type alone for checked targets without changing other fields', () => {
  const mappings: Record<string, Mapping> = {
    ...mockMappings,
    'p-01': { ...mockMappings['p-01'], target_type: 'int16', unit: '°C', scale: 0.5, offset: 2 },
    'p-02': { ...mockMappings['p-02'], target_type: 'uint64', unit: 'kWh', scale: 3, offset: 9 },
  };
  const original: WorkbenchV2State = { ...INITIAL_STATE, points: mockPoints, mappings };
  let result = original;
  const dispatch = vi.fn((action: WorkbenchV2Action) => { result = mappingReducer(result, action); });
  render(<MappingTable points={mockPoints} mappings={mappings} selectedIdx={0} setSelectedIdx={vi.fn()} devices={[mockDevice]} rawValues={{}} connectionByDevice={{}} dispatch={dispatch} />);
  fireEvent.click(screen.getByTestId('bulk-target-p-02'));
  fireEvent.click(screen.getByTestId('btn-bulk-apply-target-type'));
  expect(dispatch).not.toHaveBeenCalled();
  expect(screen.getByTestId('bulk-preview')).toHaveTextContent('target_type: uint64 → int16');
  expect(screen.getByTestId('bulk-preview')).not.toHaveTextContent('unit:');
  fireEvent.click(screen.getByTestId('bulk-confirm'));
  expect(dispatch).toHaveBeenCalledExactlyOnceWith({
    type: 'bulkApplyTransform', fromPointId: 'p-01', toPointIds: ['p-02'], fields: ['target_type'],
  });
  expect(result.mappings['p-01']).toBe(original.mappings['p-01']);
  expect(result.mappings['p-02']).toMatchObject({ ...original.mappings['p-02'], target_type: 'int16' });
});
