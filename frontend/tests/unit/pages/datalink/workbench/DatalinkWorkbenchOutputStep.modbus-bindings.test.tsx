import { fireEvent, screen, waitFor } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { mockModbusShareAPI, renderPage } from './DatalinkWorkbenchOutputStep.testHarness';

describe('DatalinkWorkbench output step', () => {
  it('binds a linked tag to a local modbus register', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByLabelText('workbench.output.mapping.register');
    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: 'workbench.output.actions.bind' }),
      ).toBeEnabled();
    });

    fireEvent.change(screen.getByLabelText('workbench.output.mapping.register'), {
      target: { value: '12' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.output.actions.bind' }));

    await waitFor(() => {
      expect(mockModbusShareAPI.upsertMapping).toHaveBeenCalledWith('tag-1', 11);
    });
  });

  it('keeps tag chip selection and register input in sync', async () => {
    mockModbusShareAPI.listMappings.mockResolvedValue([
      {
        tag_id: 'tag-1',
        register: 12,
        data_type: 'int16',
        updated_at: '',
      },
      {
        tag_id: 'tag-2',
        register: 24,
        data_type: 'int16',
        updated_at: '',
      },
    ]);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByLabelText('workbench.output.mapping.register');
    await waitFor(() => {
      expect(screen.getByDisplayValue('13')).toBeInTheDocument();
    });

    const tag1Chip = screen.getByTestId('output-candidate-tag-1');
    expect(tag1Chip).toHaveAttribute('aria-pressed', 'true');

    fireEvent.click(screen.getByTestId('output-candidate-tag-2'));

    await waitFor(() => {
      expect(screen.getByDisplayValue('25')).toBeInTheDocument();
    });
    expect(screen.getByTestId('output-candidate-tag-2')).toHaveAttribute('aria-pressed', 'true');
    expect(tag1Chip).toHaveAttribute('aria-pressed', 'false');
  });

  it('blocks sync when register conflicts exist', async () => {
    mockModbusShareAPI.status.mockResolvedValue({
      enabled: true,
      port: 5020,
      address: '127.0.0.1:5020',
      bind_state: 'pass',
      mapping_count: 2,
    });
    mockModbusShareAPI.listMappings.mockResolvedValue([
      {
        tag_id: 'tag-1',
        register: 12,
        data_type: 'int16',
        updated_at: '',
      },
      {
        tag_id: 'tag-2',
        register: 12,
        data_type: 'int16',
        updated_at: '',
      },
    ]);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(await screen.findByRole('button', { name: 'workbench.output.actions.sync' })).toBeDisabled();
    expect(screen.getByRole('status')).toHaveTextContent('workbench.output.conflicts.summary');
  });
});
