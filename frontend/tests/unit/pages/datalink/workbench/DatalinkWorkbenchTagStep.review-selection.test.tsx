import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { Tag } from '@/types/datalink';
import { mockDeletePointMutation, mockMappings, mockPoints, mockTags, renderPage, resetTagStepMocks } from './DatalinkWorkbenchTagStep.testHarness';

describe('DatalinkWorkbench tag step', () => {
  beforeEach(() => {
    resetTagStepMocks();
  });

  it('replaces preview keys with existing-tag selectors in existing flow rows', () => {
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

    const row = screen.getByTestId('tag-candidate-point-1');

    expect(within(row).queryByTestId('tag-preview-point-1')).not.toBeInTheDocument();
    expect(within(row).getByTestId('existing-tag-select-point-1')).toBeInTheDocument();
  });

  it('shows batch actions only after one or more candidate rows are selected', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.queryByTestId('batch-diff-preview')).not.toBeInTheDocument();
    expect(
      screen.queryByRole('button', { name: 'workbench.tag.actions.bind' }),
    ).not.toBeInTheDocument();

    fireEvent.click(screen.getByLabelText('Flow Sensor'));

    expect(screen.getByTestId('batch-diff-preview')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'workbench.tag.actions.bind' })).toBeInTheDocument();
  });

  it('deletes selected points from the toolbar after confirm', async () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const deleteBtn = screen.getByTestId('tag-delete-selected-points');
    expect(deleteBtn).toBeDisabled();

    fireEvent.click(screen.getByLabelText('Flow Sensor'));

    expect(deleteBtn).not.toBeDisabled();

    fireEvent.click(deleteBtn);

    await waitFor(() => {
      expect(mockDeletePointMutation.mutateAsync).toHaveBeenCalledWith('point-1');
    });

    confirmSpy.mockRestore();
  });

  it('keeps conflict detail out of the row and surfaces it in the inspector', () => {
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
    } as Tag);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const row = screen.getByTestId('tag-candidate-point-1');
    expect(within(row).queryByText('workbench.tag.badges.existingKey')).not.toBeInTheDocument();

    fireEvent.click(row);

    expect(screen.getByTestId('tag-inspector-panel')).toHaveTextContent(
      'workbench.tag.inspector.conflict.existing-key',
    );
  });

  it('shows tag inspector details when a candidate row is selected', () => {
    mockPoints[0].last_value = 123;

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.click(screen.getByTestId('tag-candidate-point-1'));

    expect(screen.getByTestId('tag-inspector-panel')).toBeInTheDocument();
    expect(screen.getByTestId('tag-inspector-tag-key')).toHaveTextContent('TAG_40001');
    expect(screen.getByTestId('tag-inspector-point-address')).toHaveTextContent('40001');
  });

  it('filters tag candidates by keyword and status', () => {
    mockMappings.push({
      id: 'mapping-2',
      point_id: 'point-2',
      tag_id: 'tag-2',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.tag.board.search'), {
      target: { value: 'Pressure' },
    });
    fireEvent.change(screen.getByLabelText('workbench.tag.board.statusFilter'), {
      target: { value: 'bound' },
    });

    expect(screen.queryByTestId('tag-candidate-point-1')).not.toBeInTheDocument();
    expect(screen.getByTestId('tag-candidate-point-2')).toBeInTheDocument();
  });

  it('shows a tag master overview with linked and unused counts', () => {
    mockTags.push(
      {
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
      },
      {
        id: 'tag-bound',
        key: 'BOUND_TAG',
        display_name: 'Bound Tag',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'active',
        created_at: '',
        updated_at: '',
      },
    );
    mockMappings.push({
      id: 'mapping-bound-tag',
      point_id: 'point-1',
      tag_id: 'tag-bound',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('tag-master-surface')).toBeInTheDocument();
    expect(screen.getByTestId('tag-master-total')).toHaveTextContent('2');
    expect(screen.getByTestId('tag-master-linked')).toHaveTextContent('1');
    expect(screen.getByTestId('tag-master-unused')).toHaveTextContent('1');
    expect(screen.getByTestId('tag-master-delete-tag-bound')).toBeDisabled();
    expect(screen.getByTestId('tag-master-delete-tag-free')).toBeEnabled();
  });
});
