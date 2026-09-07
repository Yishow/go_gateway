import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MeasurementTemplateSelector } from '@/features/datalink/workbench-v2/steps/step3/MeasurementTemplateSelector';

const mutateAsyncPreviewMock = vi.fn();
const mutateAsyncApplyMock = vi.fn();

vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2MeasurementTemplatesQuery: () => ({
    data: [
      {
        id: 'three_phase_power_meter_v1',
        name: '三相多功能電表',
        version: 'v1',
        description: '標準三相多功能電表',
        layout_mode: 'mixed',
        items: [],
      },
    ],
    isLoading: false,
  }),
  usePreviewMeasurementTemplateMutation: () => ({
    mutateAsync: mutateAsyncPreviewMock,
    isPending: false,
  }),
  useApplyMeasurementTemplateMutation: () => ({
    mutateAsync: mutateAsyncApplyMock,
    isPending: false,
  }),
}));

describe('MeasurementTemplateSelector Component', () => {
  beforeEach(() => {
    mutateAsyncPreviewMock.mockReset();
    mutateAsyncApplyMock.mockReset();
  });

  it('renders templates, triggers preview and applies template', async () => {
    const onApplied = vi.fn();
    render(<MeasurementTemplateSelector deviceId="dev-1" onApplied={onApplied} />);

    expect(screen.getByTestId('measurement-template-selector')).toBeInTheDocument();

    const select = screen.getByTestId('template-select');
    fireEvent.change(select, { target: { value: 'three_phase_power_meter_v1' } });

    mutateAsyncPreviewMock.mockResolvedValueOnce({
      device_id: 'dev-1',
      device_name: 'Power Meter 1',
      template_id: 'three_phase_power_meter_v1',
      proposed_items: [
        {
          item_id: 'v_a',
          name: 'Phase A Voltage',
          semantic_kind: 'gauge',
          unit: 'V',
          data_type: 'float32',
          register_offset: 0,
        },
      ],
      definitions: [],
      confirmed: true,
      needs_review: false,
    });

    const previewBtn = screen.getByTestId('preview-template-btn');
    fireEvent.click(previewBtn);

    await waitFor(() => {
      expect(screen.getByTestId('template-preview-summary')).toBeInTheDocument();
      expect(screen.getByText('Phase A Voltage')).toBeInTheDocument();
    });

    mutateAsyncApplyMock.mockResolvedValueOnce({ applied: true, count: 1 });
    const applyBtn = screen.getByTestId('apply-template-btn');
    fireEvent.click(applyBtn);

    await waitFor(() => {
      expect(mutateAsyncApplyMock).toHaveBeenCalled();
      expect(onApplied).toHaveBeenCalled();
    });
  });
});
