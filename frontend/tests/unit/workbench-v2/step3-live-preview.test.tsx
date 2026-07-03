import { act, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { MappingPreviewCells } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells';
import { mappingAPI } from '../../../src/services/datalink';
import type { Point, Mapping } from '../../../src/features/datalink/workbench-v2/state/types';
import type { MappingPreviewResponse } from '../../../src/types/datalink';

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
      { step_index: 1, step_type: 'decode', input_value: 243, output_value: 243, error: '' },
      { step_index: 2, step_type: 'scale', input_value: 243, output_value: 24.3, error: '' },
      { step_index: 3, step_type: 'cast', input_value: 24.3, output_value: 24.3, error: '' },
    ],
    ...overrides,
  };
}

function deferred<T>() {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>((res) => {
    resolve = res;
  });
  return { promise, resolve };
}

function renderPreviewCells(nextMapping: Mapping = mapping, rawValue: unknown = 243) {
  return render(
    <table>
      <tbody>
        <tr>
          <MappingPreviewCells point={point} mapping={nextMapping} rawValue={rawValue} />
        </tr>
      </tbody>
    </table>,
  );
}

describe('Step 3 live preview cells', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('sends preview request after debounce and updates the final cell with backend result', async () => {
    vi.mocked(mappingAPI.preview).mockResolvedValue(
      createPreviewResponse({ final_value: 987.65 }),
    );

    renderPreviewCells();

    expect(mappingAPI.preview).not.toHaveBeenCalled();
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('24.30');

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(mappingAPI.preview).toHaveBeenCalledWith({
      raw_value: 243,
      transform_pipeline: expect.any(Array),
    });
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('987.65');
  });

  it('keeps local fallback values visible while preview request is in flight', async () => {
    const pending = deferred<MappingPreviewResponse>();
    vi.mocked(mappingAPI.preview).mockReturnValue(pending.promise);

    renderPreviewCells();

    expect(screen.getByTestId('preview-scale-p-01')).toHaveTextContent('24.30');

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('24.30');

    pending.resolve(createPreviewResponse({ final_value: 654.32 }));
    await act(async () => {
      await Promise.resolve();
    });

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('654.32');
  });

  it('keeps local fallback values when preview request fails', async () => {
    vi.mocked(mappingAPI.preview).mockRejectedValue(new Error('preview failed from backend'));

    renderPreviewCells();

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('24.30');
  });

  it('ignores stale preview responses after the row mapping changes', async () => {
    const first = deferred<MappingPreviewResponse>();
    const second = deferred<MappingPreviewResponse>();
    vi.mocked(mappingAPI.preview)
      .mockReturnValueOnce(first.promise)
      .mockReturnValueOnce(second.promise);

    const { rerender } = renderPreviewCells();

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    rerender(
      <table>
        <tbody>
          <tr>
            <MappingPreviewCells
              point={point}
              mapping={{ ...mapping, scale: 0.2 }}
              rawValue={243}
            />
          </tr>
        </tbody>
      </table>,
    );

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    second.resolve(createPreviewResponse({ final_value: 654.32 }));
    await act(async () => {
      await Promise.resolve();
    });

    first.resolve(createPreviewResponse({ final_value: 24.3 }));
    await act(async () => {
      await Promise.resolve();
    });

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('654.32');
  });

  it('does not send preview request when there is no live raw value yet', async () => {
    renderPreviewCells(mapping, null);

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(mappingAPI.preview).not.toHaveBeenCalled();
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('--');
  });
});
