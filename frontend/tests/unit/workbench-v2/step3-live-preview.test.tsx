import { act, cleanup, fireEvent, render, screen } from '@testing-library/react';
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

function renderPreviewCells(nextMapping: Mapping = mapping, rawValue: unknown = 243, workspaceId?: string) {
  return render(
    <table>
      <tbody>
        <tr>
          <MappingPreviewCells point={point} mapping={nextMapping} rawValue={rawValue} workspaceId={workspaceId} />
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
    cleanup();
    vi.useRealTimers();
  });

  it('sends preview request after debounce and updates the final cell with backend result', async () => {
    vi.mocked(mappingAPI.preview).mockResolvedValue(
      createPreviewResponse({ final_value: 987.65 }),
    );

    renderPreviewCells();

    expect(mappingAPI.preview).not.toHaveBeenCalled();
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('--');

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(mappingAPI.preview).toHaveBeenCalledWith({
      raw_value: 243,
      transform_pipeline: expect.any(Array),
    });
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('987.65');
  });

  it('includes the active workspace id in preview requests', async () => {
    vi.mocked(mappingAPI.preview).mockResolvedValue(createPreviewResponse());

    renderPreviewCells(mapping, 243, 'workspace-42');
    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(mappingAPI.preview).toHaveBeenCalledWith(expect.objectContaining({ workspace_id: 'workspace-42' }));
  });

  it('keeps preview cells empty while the server request is in flight', async () => {
    const pending = deferred<MappingPreviewResponse>();
    vi.mocked(mappingAPI.preview).mockReturnValue(pending.promise);

    renderPreviewCells();

    expect(screen.getByTestId('preview-scale-p-01')).toHaveTextContent('--');

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('--');

    pending.resolve(createPreviewResponse({ final_value: 654.32 }));
    await act(async () => {
      await Promise.resolve();
    });

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('654.32');
  });

  it('fails closed when preview request fails instead of showing local values', async () => {
    vi.mocked(mappingAPI.preview).mockRejectedValue({
      response: {
        data: {
          error: {
            code: 'preview_unavailable',
            request_id: 'preview-request-42',
            retryable: true,
          },
        },
      },
    });

    renderPreviewCells();

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('--');
    expect(screen.getByTestId('preview-state-p-01')).toHaveAttribute('data-state', 'error');
    expect(screen.getByTestId('preview-error-p-01')).toHaveTextContent('preview_unavailable');
    expect(screen.getByTestId('preview-request-id-p-01')).toHaveTextContent('preview-request-42');
    expect(screen.getByTestId('preview-retry-p-01')).toBeEnabled();
  });

  it('retries a failed preview when the operator activates retry', async () => {
    vi.mocked(mappingAPI.preview)
      .mockRejectedValueOnce(new Error('transient failure'))
      .mockResolvedValueOnce(createPreviewResponse({ final_value: 55.5 }));

    renderPreviewCells();
    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    fireEvent.click(screen.getByTestId('preview-retry-p-01'));
    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(mappingAPI.preview).toHaveBeenCalledTimes(2);
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('55.50');
  });

  it('clears the previous snapshot when raw input changes before retrying', async () => {
    const retry = deferred<MappingPreviewResponse>();
    vi.mocked(mappingAPI.preview)
      .mockResolvedValueOnce(createPreviewResponse({ final_value: 12.34 }))
      .mockRejectedValueOnce(new Error('transient preview failure'))
      .mockReturnValueOnce(retry.promise);

    const { rerender } = renderPreviewCells();
    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('12.34');

    rerender(
      <table>
        <tbody>
          <tr>
            <MappingPreviewCells point={point} mapping={mapping} rawValue={244} />
          </tr>
        </tbody>
      </table>,
    );
    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('--');
    expect(screen.getByTestId('preview-state-p-01')).toHaveAttribute('data-state', 'error');

    await act(async () => {
      await vi.advanceTimersByTimeAsync(1000);
    });
    retry.resolve(createPreviewResponse({ final_value: 24.4 }));
    await act(async () => {
      await Promise.resolve();
    });

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('24.40');
    expect(screen.getByTestId('preview-state-p-01')).toHaveAttribute('data-state', 'live');
  });

  it('derives stale preview state from the existing live stream state', async () => {
    vi.mocked(mappingAPI.preview).mockResolvedValue(createPreviewResponse({ final_value: 12.34 }));

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
              mapping={mapping}
              rawValue={243}
              connectionState="stale"
            />
          </tr>
        </tbody>
      </table>,
    );

    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('12.34');
    expect(screen.getByTestId('preview-state-p-01')).toHaveAttribute('data-state', 'stale');
  });

  it('uses only the bounded retry schedule for an initial snapshot failure', async () => {
    vi.mocked(mappingAPI.preview).mockRejectedValue(new Error('preview unavailable'));

    renderPreviewCells();
    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
    });

    for (const delay of [1000, 2000, 4000, 8000, 16000, 30000]) {
      await act(async () => {
        await vi.advanceTimersByTimeAsync(delay);
      });
    }

    expect(mappingAPI.preview).toHaveBeenCalledTimes(7);
    expect(screen.getByTestId('preview-state-p-01')).toHaveAttribute('data-state', 'error');
    expect(screen.getByTestId('preview-final-p-01')).toHaveTextContent('--');

    await act(async () => {
      await vi.advanceTimersByTimeAsync(60000);
    });
    expect(mappingAPI.preview).toHaveBeenCalledTimes(7);
  });

});
