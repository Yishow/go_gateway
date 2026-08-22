import { fireEvent, screen, waitFor } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  getSourceStepMocks,
  registerSourceStepFixtures,
  renderPage,
} from './DatalinkWorkbenchSourceStep.testHarness';
const { mockPoints, mockCreatePointMutation, mockCreateSourceRuleMutation } = getSourceStepMocks();

describe('DatalinkWorkbench source step', () => {
  registerSourceStepFixtures();

  it('visually merges 32-bit cells with gridColumn span on root and hides continuations', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const rootCell = screen.getByTestId('address-cell-40001');
    const contCell = screen.getByTestId('address-cell-40002');

    // Root cell should span 2 columns
    expect(rootCell.style.gridColumn).toBe('span 2');
    // Continuation cell should be visually hidden
    expect(contCell).toHaveClass('sr-only');
  });

  it('persists only safe spans when conflicts exist alongside valid planned ranges', async () => {
    // Rule plans int16 at 40001, 40002, 40003.
    // Existing point at 40005 (no conflict).
    // Add a second rule at 40003 to create a rule-overlap conflict at 40003.
    mockPoints.splice(0, mockPoints.length, {
      id: 'point-1',
      device_id: 'device-1',
      name: 'Existing Pressure',
      description: '',
      data_type: 'int16',
      address: '40005',
      enabled: true,
      polling_group_id: '',
      last_value: 12,
      last_read_at: '',
      last_error: '',
      error_count: 0,
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    // Rule 1: int16 at 40001, count=3 → plans 40001, 40002, 40003
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    // Rule 2: int16 at 40003, count=1 → overlap at 40003
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40003' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));

    // 40003 is conflict, but 40001 and 40002 are planned and safe
    expect(screen.getByTestId('address-cell-40003')).toHaveAttribute('data-status', 'conflict');

    // Batch create should still be enabled — 2 safe spans
    const batchButton = screen.getByRole('button', {
      name: /workbench\.source\.actions\.createRulePoints/,
    });
    expect(batchButton).toBeEnabled();

    fireEvent.click(batchButton);

    await waitFor(() => {
      expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith(
      expect.not.objectContaining({ id: expect.anything() }),
    );
    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith(
      expect.objectContaining({
        start_address: '40001',
        skipped_addresses: ['40003'],
      }),
    );
    expect(mockCreatePointMutation.mutateAsync).not.toHaveBeenCalled();
  });
});
