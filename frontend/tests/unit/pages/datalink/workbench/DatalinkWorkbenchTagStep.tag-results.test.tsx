import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { mockBatchCreate, mockDeleteMappingMutation, mockMappings, mockTags, renderPage, resetTagStepMocks } from './DatalinkWorkbenchTagStep.testHarness';

describe('DatalinkWorkbench tag step', () => {
  beforeEach(() => {
    resetTagStepMocks();
  });

  it('shows a partial failure summary when part of the batch bind fails', async () => {
    // 模擬 batch create 部分失敗（default strategy=address, prefix=TAG）
    mockBatchCreate.mockImplementationOnce(async () => {
      // TAG_40001 成功，TAG_40002 失敗
      mockTags.push({
        id: 'tag-TAG_40001',
        key: 'TAG_40001',
        display_name: 'Flow Sensor',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'draft',
        created_at: '',
        updated_at: '',
      });
      return {
        created: ['tag-TAG_40001'],
        errors: [{ key: 'TAG_40002', error: 'duplicate key' }],
      };
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.actions.bind' }));

    await waitFor(() => {
      expect(mockBatchCreate).toHaveBeenCalledTimes(1);
    });

    await waitFor(() => {
      expect(screen.getByText('workbench.tag.results.partialFailure')).toBeInTheDocument();
    });
    expect(screen.getByText('duplicate key')).toBeInTheDocument();
  });
  it('shows a batch diff preview before executing the bind', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));

    const diffPanel = screen.getByTestId('batch-diff-preview');
    expect(diffPanel).toBeInTheDocument();
    expect(diffPanel).toHaveTextContent('TAG_40001');
    expect(diffPanel).toHaveTextContent('TAG_40002');
    expect(screen.getByTestId('diff-to-create-count')).toHaveTextContent('2');
    expect(screen.getByTestId('diff-skipped-count')).toHaveTextContent('0');
  });
  it('shows skipped items in the diff preview when candidates have conflicts', () => {
    mockTags.push({
      id: 'tag-existing',
      key: 'TAG_40001',
      display_name: 'Existing Flow Sensor',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'draft',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));

    expect(screen.getByTestId('diff-to-create-count')).toHaveTextContent('1');
    expect(screen.getByTestId('diff-skipped-count')).toHaveTextContent('1');
  });
  it('switches diff preview to show toBind items in existing flow', () => {
    mockTags.push({
      id: 'tag-existing',
      key: 'LINEA_FLOW',
      display_name: 'Line A Flow',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.board.flow.existing' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));

    expect(screen.getByTestId('diff-to-bind-count')).toBeInTheDocument();
  });
  it('shows an enhanced result summary with created/linked/skipped/failed counts', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.actions.bind' }));

    await waitFor(() => {
      expect(mockBatchCreate).toHaveBeenCalledTimes(1);
    });

    expect(screen.getByTestId('result-created-count')).toHaveTextContent('2');
    expect(screen.getByTestId('result-failed-count')).toHaveTextContent('0');
  });
  it('shows an unbind button for bound candidates and calls delete mutation after confirm', async () => {
    mockTags.push({
      id: 'tag-bound-1',
      key: 'TAG_40001',
      display_name: 'TAG 40001',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });
    mockMappings.push({
      id: 'mapping-bound-1',
      point_id: 'point-1',
      tag_id: 'tag-bound-1',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });
    mockDeleteMappingMutation.mutateAsync.mockResolvedValue(undefined);
    vi.spyOn(window, 'confirm').mockReturnValue(true);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const candidate = screen.getByTestId('tag-candidate-point-1');
    const unbindButton = within(candidate).getByTestId('tag-unbind-point-1');
    expect(unbindButton).toBeInTheDocument();

    const boundLabel = within(candidate).getByTestId('tag-bound-label-point-1');
    expect(boundLabel).toBeInTheDocument();

    fireEvent.click(unbindButton);

    await waitFor(() => {
      expect(window.confirm).toHaveBeenCalled();
      expect(mockDeleteMappingMutation.mutateAsync).toHaveBeenCalledWith('mapping-bound-1');
    });

    vi.restoreAllMocks();
  });
  it('does not unbind when the user cancels the confirmation dialog', async () => {
    mockTags.push({
      id: 'tag-bound-2',
      key: 'TAG_40001',
      display_name: 'TAG 40001',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });
    mockMappings.push({
      id: 'mapping-bound-2',
      point_id: 'point-1',
      tag_id: 'tag-bound-2',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });
    vi.spyOn(window, 'confirm').mockReturnValue(false);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const candidate = screen.getByTestId('tag-candidate-point-1');
    const unbindButton = within(candidate).getByTestId('tag-unbind-point-1');
    fireEvent.click(unbindButton);

    expect(window.confirm).toHaveBeenCalled();
    expect(mockDeleteMappingMutation.mutateAsync).not.toHaveBeenCalled();

    vi.restoreAllMocks();
  });
  it('shows a batch unbind action for selected bound rows and removes every selected mapping', async () => {
    mockTags.push(
      {
        id: 'tag-bound-1',
        key: 'TAG_40001',
        display_name: 'TAG 40001',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'active',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'tag-bound-2',
        key: 'TAG_40002',
        display_name: 'TAG 40002',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'active',
        created_at: '',
        updated_at: '',
      },
    );
    mockMappings.push(
      {
        id: 'mapping-bound-1',
        point_id: 'point-1',
        tag_id: 'tag-bound-1',
        enabled: true,
        transform_pipeline: '',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'mapping-bound-2',
        point_id: 'point-2',
        tag_id: 'tag-bound-2',
        enabled: true,
        transform_pipeline: '',
        created_at: '',
        updated_at: '',
      },
    );
    mockDeleteMappingMutation.mutateAsync.mockResolvedValue(undefined);
    vi.spyOn(window, 'confirm').mockReturnValue(true);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));

    const batchUnbindButton = await screen.findByTestId('tag-batch-unbind');
    expect(batchUnbindButton).toBeInTheDocument();

    fireEvent.click(batchUnbindButton);

    await waitFor(() => {
      expect(window.confirm).toHaveBeenCalled();
      expect(mockDeleteMappingMutation.mutateAsync).toHaveBeenCalledWith('mapping-bound-1');
      expect(mockDeleteMappingMutation.mutateAsync).toHaveBeenCalledWith('mapping-bound-2');
    });

    vi.restoreAllMocks();
  });
  it('shows the template panel collapsed by default with a toggle button', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const toggle = screen.getByTestId('tag-template-toggle');
    expect(toggle).toBeInTheDocument();

    expect(screen.queryByLabelText('workbench.tag.template.prefix')).not.toBeInTheDocument();

    fireEvent.click(toggle);
    expect(screen.getByLabelText('workbench.tag.template.prefix')).toBeInTheDocument();
  });
});
