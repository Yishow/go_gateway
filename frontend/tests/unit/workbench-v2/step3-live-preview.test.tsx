import { act, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { TransformPreview } from '../../../src/features/datalink/workbench-v2/steps/step3/TransformPreview';
import { mappingAPI } from '../../../src/services/datalink';
import type { Point, Mapping } from '../../../src/features/datalink/workbench-v2/state/types';
import type { MappingPreviewResponse } from '../../../src/types/datalink';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: { type?: string }) => {
      if (options?.type) {
        return `${key}_${options.type}`;
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

const point: Point = {
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
};

const mapping: Mapping = {
  point_id: 'p-01',
  tag_key: 'line01.temp.inlet',
  display_name: '進水溫度',
  unit: '°C',
  target_type: 'float64',
  scale: 0.1,
  offset: 0,
  enabled: true,
};

function createPreviewResponse(
  overrides: Partial<MappingPreviewResponse> = {},
): MappingPreviewResponse {
  return {
    raw_value: 243,
    final_value: 24.3,
    step_results: [
      {
        step_index: 1,
        step_type: 'decode',
        input_value: 243,
        output_value: 243,
        error: '',
      },
      {
        step_index: 2,
        step_type: 'scale',
        input_value: 243,
        output_value: 24.3,
        error: '',
      },
      {
        step_index: 3,
        step_type: 'cast',
        input_value: 24.3,
        output_value: 24.3,
        error: '',
      },
    ],
    ...overrides,
  };
}

function deferred<T>() {
  let resolve!: (value: T) => void;
  let reject!: (reason?: unknown) => void;
  const promise = new Promise<T>((res, rej) => {
    resolve = res;
    reject = rej;
  });
  return { promise, resolve, reject };
}

describe('Step 3 live preview', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('sends preview request for the selected row after debounce and renders backend result', async () => {
    vi.mocked(mappingAPI.preview).mockResolvedValue(
      createPreviewResponse({
        final_value: 987.65,
        step_results: [
          {
            step_index: 1,
            step_type: 'decode',
            input_value: 243,
            output_value: 243,
            error: '',
          },
          {
            step_index: 2,
            step_type: 'scale',
            input_value: 243,
            output_value: 987.65,
            error: '',
          },
          {
            step_index: 3,
            step_type: 'cast',
            input_value: 987.65,
            output_value: 987.65,
            error: '',
          },
        ],
      }),
    );

    render(<TransformPreview point={point} mapping={mapping} rawValue={243} />);

    expect(mappingAPI.preview).not.toHaveBeenCalled();

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(mappingAPI.preview).toHaveBeenCalledTimes(1);
    expect(mappingAPI.preview).toHaveBeenCalledWith({
      raw_value: 243,
      transform_pipeline: expect.any(Array),
    });

    expect(screen.getByTestId('step-final')).toHaveTextContent('987.65');
  });

  it('shows explicit loading state while preview request is in flight', async () => {
    const pending = deferred<MappingPreviewResponse>();
    vi.mocked(mappingAPI.preview).mockReturnValue(pending.promise);

    render(<TransformPreview point={point} mapping={mapping} rawValue={243} />);

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(screen.getByTestId('preview-loading')).toBeInTheDocument();

    pending.resolve(
      createPreviewResponse({
        final_value: 987.65,
      }),
    );
    await act(async () => {
      await Promise.resolve();
    });
    expect(screen.getByTestId('step-final')).toHaveTextContent('987.65');
  });

  it('shows actionable error when preview request fails', async () => {
    vi.mocked(mappingAPI.preview).mockRejectedValue(new Error('preview failed from backend'));

    render(<TransformPreview point={point} mapping={mapping} rawValue={243} />);

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(screen.getByTestId('preview-error')).toHaveTextContent(
      'preview failed from backend',
    );
  });

  it('keeps the latest preview result when an older request resolves later', async () => {
    const first = deferred<MappingPreviewResponse>();
    const second = deferred<MappingPreviewResponse>();
    vi.mocked(mappingAPI.preview)
      .mockReturnValueOnce(first.promise)
      .mockReturnValueOnce(second.promise);

    const { rerender } = render(
      <TransformPreview point={point} mapping={mapping} rawValue={243} />,
    );

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    rerender(
      <TransformPreview
        point={point}
        mapping={{ ...mapping, scale: 0.2 }}
        rawValue={243}
      />,
    );

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    second.resolve(createPreviewResponse({ final_value: 654.32 }));
    await act(async () => {
      await Promise.resolve();
    });
    expect(screen.getByTestId('step-final')).toHaveTextContent('654.32');

    first.resolve(createPreviewResponse({ final_value: 24.3 }));

    await act(async () => {
      await Promise.resolve();
    });
    expect(screen.getByTestId('step-final')).toHaveTextContent('654.32');
  });

  it('does not send preview request when there is no selected row', async () => {
    render(<TransformPreview point={null} mapping={null} rawValue={null} />);

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(mappingAPI.preview).not.toHaveBeenCalled();
    expect(screen.getByTestId('preview-empty')).toBeInTheDocument();
  });

  it('waits for live device value instead of sending preview request with a mock seed', async () => {
    render(
      <TransformPreview
        point={point}
        mapping={mapping}
        rawValue={null}
        connectionState="connected"
      />,
    );

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(mappingAPI.preview).not.toHaveBeenCalled();
    expect(screen.getByTestId('preview-waiting')).toHaveTextContent('step3.preview.waiting');
  });
});
