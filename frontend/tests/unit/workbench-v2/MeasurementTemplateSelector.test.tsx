import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MeasurementTemplateSelector } from '@/features/datalink/workbench-v2/steps/step3/MeasurementTemplateSelector';

const mutateAsyncPreviewMock = vi.fn();
const mutateAsyncApplyMock = vi.fn();
const queryState = { isLoading: false, isError: false, error: null as unknown, refetch: vi.fn() };
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

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
    ...queryState,
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
    queryState.isLoading = false; queryState.isError = false; queryState.error = null; queryState.refetch.mockReset();
  });

  it('keeps measurement semantics collapsed as an optional advanced section', () => {
    render(<MeasurementTemplateSelector deviceId="dev-1" />);
    const advanced = screen.getByTestId('measurement-advanced');
    expect(advanced.tagName).toBe('DETAILS');
    expect(advanced).not.toHaveAttribute('open');
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
      confirmed: false,
      needs_review: true,
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
      expect(mutateAsyncApplyMock).toHaveBeenCalledWith(expect.objectContaining({ confirmed: true }));
      expect(onApplied).toHaveBeenCalled();
    });
  });
  const preview = { device_id: 'dev-1', device_name: 'Meter', template_id: 'three_phase_power_meter_v1', proposed_items: [], definitions: [], confirmed: false, needs_review: true };
  async function selectAndPreview() {
    fireEvent.change(screen.getByTestId('template-select'), { target: { value: preview.template_id } });
    fireEvent.click(screen.getByTestId('preview-template-btn'));
    await screen.findByTestId('template-preview-summary');
  }

  it('invalidates a preview immediately on device change and ignores a late preview reply', async () => {
    let resolve!: (value: typeof preview) => void;
    mutateAsyncPreviewMock.mockReturnValueOnce(new Promise(r => { resolve = r; }));
    const view = render(<MeasurementTemplateSelector deviceId="dev-1" />);
    fireEvent.change(screen.getByTestId('template-select'), { target: { value: preview.template_id } });
    fireEvent.click(screen.getByTestId('preview-template-btn'));
    view.rerender(<MeasurementTemplateSelector deviceId="dev-2" />);
    resolve(preview);
    await waitFor(() => expect(screen.queryByTestId('apply-template-btn')).not.toBeInTheDocument());
    expect(mutateAsyncApplyMock).not.toHaveBeenCalled();
  });

  it('clears an existing preview on device or template change', async () => {
    mutateAsyncPreviewMock.mockResolvedValue(preview);
    const view = render(<MeasurementTemplateSelector deviceId="dev-1" />);
    await selectAndPreview();
    view.rerender(<MeasurementTemplateSelector deviceId="dev-2" />);
    expect(screen.queryByTestId('apply-template-btn')).not.toBeInTheDocument();
    view.rerender(<MeasurementTemplateSelector deviceId="dev-1" />);
    await selectAndPreview();
    fireEvent.change(screen.getByTestId('template-select'), { target: { value: '' } });
    expect(screen.queryByTestId('apply-template-btn')).not.toBeInTheDocument();
  });

  it.each(['offline', '409', '500'])('shows safe actionable preview failure for %s', async failure => {
    mutateAsyncPreviewMock.mockRejectedValueOnce(new Error(`secret endpoint password ${failure}`)).mockResolvedValueOnce(preview);
    render(<MeasurementTemplateSelector deviceId="dev-1" />);
    fireEvent.change(screen.getByTestId('template-select'), { target: { value: preview.template_id } });
    fireEvent.click(screen.getByTestId('preview-template-btn'));
    const error = await screen.findByTestId('template-error');
    expect(error).not.toHaveTextContent('password');
    fireEvent.click(screen.getByTestId('preview-template-btn'));
    await screen.findByTestId('template-preview-summary');
  });

  it.each(['offline', '409', '500'])('shows safe actionable apply failure for %s and permits re-preview', async failure => {
    mutateAsyncPreviewMock.mockResolvedValue(preview);
    mutateAsyncApplyMock.mockRejectedValueOnce(new Error(`secret endpoint password ${failure}`));
    render(<MeasurementTemplateSelector deviceId="dev-1" />);
    await selectAndPreview();
    fireEvent.click(screen.getByTestId('apply-template-btn'));
    expect(await screen.findByTestId('template-error')).not.toHaveTextContent('password');
    fireEvent.click(screen.getByTestId('preview-template-btn'));
    await screen.findByTestId('template-preview-summary');
  });

  it('distinguishes template load failure and offers reload', () => {
    queryState.isError = true; queryState.error = new Error('private endpoint');
    render(<MeasurementTemplateSelector deviceId="dev-1" />);
    expect(screen.getByTestId('template-load-error')).not.toHaveTextContent('private');
    fireEvent.click(screen.getByTestId('template-load-retry'));
    expect(queryState.refetch).toHaveBeenCalledOnce();
  });

  it('keeps the opened advanced section during template reload', () => {
    queryState.isError = true;
    const view = render(<MeasurementTemplateSelector deviceId="dev-1" />);
    const details = screen.getByTestId('measurement-advanced') as HTMLDetailsElement;
    details.open = true;
    queryState.isError = false; queryState.isLoading = true;
    view.rerender(<MeasurementTemplateSelector deviceId="dev-1" />);
    expect(screen.getByTestId('measurement-advanced')).toHaveAttribute('open');
    queryState.isLoading = false;
    view.rerender(<MeasurementTemplateSelector deviceId="dev-1" />);
    expect(screen.getByTestId('measurement-advanced')).toHaveAttribute('open');
  });

});
