import { fireEvent, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { mockCreateMappingMutation, mockCreateTagMutation, mockDeleteTagMutation, mockTags, renderPage, resetTagStepMocks } from './DatalinkWorkbenchTagStep.testHarness';

describe('DatalinkWorkbench tag step', () => {
  beforeEach(() => {
    resetTagStepMocks();
  });

  it('creates and deletes standalone tags from the tag master surface', async () => {
    mockTags.push({
      id: 'tag-free',
      key: 'FREE_TAG',
      display_name: 'Free Tag',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'draft',
      created_at: '',
      updated_at: '',
    });
    vi.spyOn(window, 'confirm').mockReturnValue(true);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.tag.master.fields.key'), {
      target: { value: 'MANUAL_TAG' },
    });
    fireEvent.change(screen.getByLabelText('workbench.tag.master.fields.displayName'), {
      target: { value: 'Manual Tag' },
    });
    fireEvent.change(screen.getByLabelText('workbench.tag.master.fields.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.click(screen.getByTestId('tag-master-create'));

    await waitFor(() => {
      expect(mockCreateTagMutation.mutateAsync).toHaveBeenCalledWith({
        key: 'MANUAL_TAG',
        display_name: 'Manual Tag',
        data_type: 'float32',
      });
    });

    fireEvent.click(screen.getByTestId('tag-master-delete-tag-free'));

    await waitFor(() => {
      expect(window.confirm).toHaveBeenCalled();
      expect(mockDeleteTagMutation.mutateAsync).toHaveBeenCalledWith('tag-free');
    });

    vi.restoreAllMocks();
  });

  it('supports binding selected points to an existing tag', async () => {
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
    fireEvent.change(screen.getByTestId('existing-tag-select-point-1'), {
      target: { value: 'tag-existing' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.actions.bind' }));

    await waitFor(() => {
      expect(mockCreateMappingMutation.mutateAsync).toHaveBeenCalledWith({
        point_id: 'point-1',
        tag_id: 'tag-existing',
        enabled: true,
      });
    });

    expect(mockCreateTagMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('blocks batch binding when preview keys conflict with existing tags', () => {
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

    expect(screen.getByTestId('tag-preview-point-1')).toHaveAttribute('data-conflict', 'true');
    expect(screen.getByRole('button', { name: 'workbench.tag.actions.bind' })).toBeDisabled();
    expect(screen.getByText('workbench.tag.conflicts.summary')).toBeInTheDocument();
  });
});
