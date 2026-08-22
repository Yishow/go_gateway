import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  getSourceStepMocks,
  registerSourceStepFixtures,
  renderPage,
} from './DatalinkWorkbenchSourceStep.testHarness';
const { mockPoints, mockCreatePointMutation, mockCreateSourceRuleMutation } = getSourceStepMocks();

describe('DatalinkWorkbench source step', () => {
  registerSourceStepFixtures();

  it('batch persists eligible source rules from the planned address range', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.actions.createRulePoints' }),
    );

    await waitFor(() => {
      expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith(
      expect.not.objectContaining({ id: expect.anything() }),
    );
    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith(
      expect.objectContaining({
        device_id: 'device-1',
        start_address: '40001',
        count: 2,
        data_type: 'float32',
        naming_prefix: 'MBT',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
      }),
    );
    expect(mockCreatePointMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('allows batch create for safe spans even when conflicts exist elsewhere', () => {
    mockPoints[0].address = '40002';

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));

    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-status', 'conflict');
    // The float32 at 40001 occupies 40001-40002, and 40002 is a conflict, so the
    // only planned span is itself unsafe — button should be disabled.
    expect(
      screen.getByRole('button', { name: 'workbench.source.actions.createRulePoints' }),
    ).toBeDisabled();
  });

  it('shows a selection toolbar when a planned address cell is clicked', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.queryByTestId('source-selection-toolbar')).not.toBeInTheDocument();

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    const toolbar = screen.getByTestId('source-selection-toolbar');
    expect(toolbar).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('button', {
        name: 'workbench.source.selectionToolbar.createPoints',
      }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('button', {
        name: 'workbench.source.selectionToolbar.addToRule',
      }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('button', {
        name: 'workbench.source.selectionToolbar.skip',
      }),
    ).toBeInTheDocument();
  });

  it('creates a point via the selection toolbar create action', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    const toolbar = screen.getByTestId('source-selection-toolbar');
    fireEvent.click(
      within(toolbar).getByRole('button', {
        name: 'workbench.source.selectionToolbar.createPoints',
      }),
    );

    await waitFor(() => {
      expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledWith({
      device_id: 'device-1',
      address: '40001',
      data_type: 'int16',
      name: 'MBT_40001',
    });
  });

  it('clears the selected address when skip is clicked in the selection toolbar', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));
    expect(screen.getByTestId('source-selection-toolbar')).toBeInTheDocument();

    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.selectionToolbar.skip' }),
    );

    expect(screen.queryByTestId('source-selection-toolbar')).not.toBeInTheDocument();
  });
});
