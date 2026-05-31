import { render, screen, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Step3Mapping } from '../../../src/features/datalink/workbench-v2/steps/step3/Step3Mapping';
import { DeviceListContext } from '../../../src/features/datalink/workbench-v2/state/deviceColors';
import type { Device, Mapping, Rule, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { mappingAPI } from '../../../src/services/datalink';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: { count?: number; type?: string; tag?: string }) => {
      if (options?.count !== undefined) {
        return `${key}_count_${options.count}`;
      }
      if (options?.type) {
        return `${key}_type_${options.type}`;
      }
      if (options?.tag) {
        return `${key}_tag_${options.tag}`;
      }
      return key;
    },
  }),
}));

vi.mock('../../../src/services/datalink', () => ({
  mappingAPI: {
    preview: vi.fn(),
  },
}));

const device: Device = {
  id: 'dev-01',
  name: 'PLC-生產線-01',
  description: 'Modbus TCP',
  protocol: 'modbus_tcp',
  config: {},
  status: 'draft',
  test: null,
};

const rule: Rule = {
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

const mappings: Record<string, Mapping> = {
  'rule-01-p-0': {
    point_id: 'rule-01-p-0',
    tag_key: 'line01.temp.inlet',
    display_name: '進水溫度',
    unit: '°C',
    target_type: 'float64',
    scale: 0.1,
    offset: 0,
    enabled: true,
  },
  'rule-01-p-1': {
    point_id: 'rule-01-p-1',
    tag_key: 'line01.temp.outlet',
    display_name: '出水溫度',
    unit: '°C',
    target_type: 'float64',
    scale: 0.1,
    offset: 0,
    enabled: true,
  },
};

function renderStep3(state: WorkbenchV2State, dispatch = vi.fn()) {
  render(
    <DeviceListContext.Provider value={[device]}>
      <Step3Mapping
        state={state}
        dispatch={dispatch}
        onContinue={vi.fn()}
        onBack={vi.fn()}
      />
    </DeviceListContext.Provider>,
  );
  return dispatch;
}

describe('Step 3 target type quick actions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
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

  it('updates the selected row target_type when a quick-action chip is clicked', () => {
    const dispatch = renderStep3(
      {
        ...INITIAL_STATE,
        devices: [device],
        rules: [rule],
        mappings,
      },
    );

    fireEvent.click(screen.getByTestId('target-type-chip-int16'));

    expect(dispatch).toHaveBeenCalledWith({
      type: 'updateMapping',
      pointId: 'rule-01-p-0',
      patch: { target_type: 'int16' },
    });
  });

  it('applies the selected row target_type to all rows', () => {
    const dispatch = renderStep3(
      {
        ...INITIAL_STATE,
        devices: [device],
        rules: [rule],
        mappings: {
          ...mappings,
          'rule-01-p-0': {
            ...mappings['rule-01-p-0'],
            target_type: 'int16',
          },
        },
      },
    );

    fireEvent.click(screen.getByTestId('btn-apply-target-type-all'));

    expect(dispatch).toHaveBeenCalledWith({
      type: 'bulkApplyTransform',
      fromPointId: 'rule-01-p-0',
      fields: ['target_type'],
    });
  });

  it('does not render active quick actions when there is no selected row', () => {
    renderStep3({
      ...INITIAL_STATE,
      devices: [device],
      rules: [],
      mappings: {},
    });

    expect(screen.getByTestId('preview-empty')).toBeInTheDocument();
    expect(screen.queryByTestId('target-type-quick-actions')).not.toBeInTheDocument();
  });
});
