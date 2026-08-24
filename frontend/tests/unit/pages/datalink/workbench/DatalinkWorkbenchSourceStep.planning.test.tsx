import { fireEvent, screen, waitFor } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  getSourceStepMocks,
  registerSourceStepFixtures,
  renderPage,
} from './DatalinkWorkbenchSourceStep.testHarness';
const { mockDevices, mockPoints, mockSourceRules, mockToggleDeviceStatusMutation } =
  getSourceStepMocks();

describe('DatalinkWorkbench source step', () => {
  registerSourceStepFixtures();

  it('gates source planning behind device selection', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));

    expect(screen.getByText('workbench.source.empty.title')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Mixer PLC' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toBeInTheDocument();
    expect(screen.getByLabelText('workbench.source.planner.count')).toHaveValue(4);
    expect(screen.getByLabelText('workbench.source.planner.namingPrefix')).toHaveValue('MBT');
  });

  it('switches rule layer sidebar between planner and added-rules tabs', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const tabPlanner = screen.getByTestId('source-rule-layer-tab-planner');
    const tabRules = screen.getByTestId('source-rule-layer-tab-rules');
    expect(tabPlanner).toHaveAttribute('aria-selected', 'true');
    expect(tabRules).toHaveAttribute('aria-selected', 'false');
    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toBeVisible();

    fireEvent.click(tabRules);
    expect(tabPlanner).toHaveAttribute('aria-selected', 'false');
    expect(tabRules).toHaveAttribute('aria-selected', 'true');
    expect(screen.getByText('workbench.source.ruleLayer.empty')).toBeVisible();

    fireEvent.click(tabPlanner);
    expect(tabPlanner).toHaveAttribute('aria-selected', 'true');
    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toBeVisible();
  });

  it('remembers the last planner start address per device and falls back to protocol defaults', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const startAddressInput = screen.getByLabelText('workbench.source.planner.startAddress');
    fireEvent.change(startAddressInput, { target: { value: '40010' } });
    expect(startAddressInput).toHaveValue('40010');

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.device/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Fatek Cell' }));
    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.source.planner.startAddress')).toHaveValue('D0');
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));
    expect(screen.getByTestId('address-cell-D0')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-D2')).toHaveAttribute('data-status', 'planned');

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: 'D20' },
    });
    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toHaveValue('D20');

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.device/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.source.planner.startAddress')).toHaveValue('40010');
    });
  });

  it('allows clearing planner count and disables add rule until count is valid', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const countInput = screen.getByLabelText('workbench.source.planner.count');
    fireEvent.change(countInput, { target: { value: '' } });
    expect(countInput).toHaveValue(null);

    const addRule = screen.getByRole('button', { name: 'workbench.source.planner.addRule' });
    expect(addRule).toBeDisabled();

    fireEvent.change(countInput, { target: { value: '2' } });
    expect(addRule).not.toBeDisabled();
  });

  it('renders a source rule layer and continuous gap cells after applying a rule', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('source-coverage-overview')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40003')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40004')).toHaveAttribute('data-status', 'gap');
    expect(screen.getByTestId('address-cell-40005')).toHaveAttribute('data-status', 'unmanaged');
  });

  it('loads persisted source rules from backend state for the selected device', () => {
    mockSourceRules.splice(0, mockSourceRules.length, {
      id: 'persisted-rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 2,
      data_type: 'int16',
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      template_name: '',
      skipped_addresses: [],
      created_at: '',
      updated_at: '2026-03-19T00:00:00Z',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('source-rule-persisted-rule-1')).toBeInTheDocument();
    expect(screen.getByText('workbench.source.ruleLayer.persistedBadge')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-status', 'planned');
  });

  it('shows runtime summary counters and collection hint in the source runtime card', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));

    expect(screen.getByTestId('source-runtime-collection-panel')).toBeInTheDocument();
    expect(screen.getByTestId('source-runtime-collection-hint')).toBeInTheDocument();
    expect(mockToggleDeviceStatusMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('preserves draft rules per device when switching the selected device', async () => {
    mockDevices[0].protocol = 'fatek_fbs';
    mockPoints[0].address = 'D15';
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: 'D10' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-D10')).toHaveAttribute('data-status', 'planned');

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.device/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Fatek Cell' }));
    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));

    await waitFor(() => {
      expect(screen.queryByTestId('source-rule-rule-1')).not.toBeInTheDocument();
    });
    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toHaveValue('D0');

    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-D0')).toHaveAttribute('data-status', 'planned');
    expect(screen.queryByTestId('address-cell-D10')).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.device/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));

    await waitFor(() => {
      expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    });
    expect(screen.getByTestId('address-cell-D10')).toHaveAttribute('data-status', 'planned');
    expect(screen.queryByTestId('address-cell-D0')).not.toBeInTheDocument();
  });

  it('keeps each device persisted and draft rules isolated when switching back and forth', async () => {
    mockSourceRules.splice(
      0,
      mockSourceRules.length,
      {
        id: 'shared-rule',
        device_id: 'device-1',
        start_address: '40011',
        count: 2,
        data_type: 'int16',
        naming_prefix: 'SRC_A',
        enabled: true,
        locked: false,
        origin: 'manual',
        template_name: '',
        skipped_addresses: [],
        created_at: '',
        updated_at: '2026-03-19T00:00:00Z',
      },
      {
        id: 'shared-rule',
        device_id: 'device-2',
        start_address: 'D10',
        count: 1,
        data_type: 'int16',
        naming_prefix: 'SRC_B',
        enabled: true,
        locked: false,
        origin: 'manual',
        template_name: '',
        skipped_addresses: [],
        created_at: '',
        updated_at: '2026-03-19T00:00:00Z',
      },
    );

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('address-cell-40011')).toHaveAttribute('data-status', 'planned');

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.device/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Fatek Cell' }));
    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));

    await waitFor(() => {
      expect(screen.getByTestId('address-cell-D10')).toHaveAttribute('data-status', 'planned');
    });
    expect(screen.queryByTestId('address-cell-40011')).not.toBeInTheDocument();
    expect(screen.queryByTestId('address-cell-40001')).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.device/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));

    await waitFor(() => {
      expect(screen.getByTestId('address-cell-40011')).toHaveAttribute('data-status', 'planned');
    });
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.queryByTestId('address-cell-D10')).not.toBeInTheDocument();
  });
});
