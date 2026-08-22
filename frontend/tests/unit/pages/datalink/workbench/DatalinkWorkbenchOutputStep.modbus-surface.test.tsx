import { fireEvent, screen, waitFor } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  mockDBTargetAPI,
  mockMappings,
  mockModbusShareAPI,
  renderPage,
} from './DatalinkWorkbenchOutputStep.testHarness';

describe('DatalinkWorkbench output step', () => {
  it('shows an empty state when the selected device has no linked tags', async () => {
    mockMappings.splice(0, mockMappings.length);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(await screen.findByText('workbench.output.empty.title')).toBeInTheDocument();
  });

  it('starts the local modbus server with the specified port', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByLabelText('workbench.output.mapping.register');
    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: 'workbench.output.actions.startServer' }),
      ).toBeEnabled();
    });

    fireEvent.change(screen.getByLabelText('workbench.output.server.port'), {
      target: { value: '5030' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.output.actions.startServer' }));

    await waitFor(() => {
      expect(mockModbusShareAPI.start).toHaveBeenCalledWith(5030);
    });
  });

  it('renders compact tag chips with status suffixes for the active target', async () => {
    mockModbusShareAPI.listMappings.mockResolvedValue([
      {
        tag_id: 'tag-1',
        register: 12,
        data_type: 'int16',
        updated_at: '',
      },
    ]);
    mockDBTargetAPI.listMappings.mockResolvedValue([
      {
        id: 'db-mapping-2',
        tag_id: 'tag-2',
        connector_id: 'connector-1',
        table_schema: 'main',
        table_name: 'sensor_values',
        column_name: 'value',
        write_mode: 'insert',
        timestamp_column: null,
        enabled: true,
        created_at: '',
        updated_at: '',
      },
    ]);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.modbus',
      }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    ).toBeInTheDocument();

    const tagChips = screen.getByTestId('output-tag-chips');
    expect(tagChips).toBeInTheDocument();
    expect(screen.getByTestId('output-candidate-tag-1')).toBeInTheDocument();

    // Click tag-1 chip and verify active badge appears
    fireEvent.click(screen.getByTestId('output-candidate-tag-1'));
    expect(screen.getByTestId('active-tag-badge')).toBeInTheDocument();

    // Switch to database target
    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    // Select tag-2 via chip click
    fireEvent.click(screen.getByTestId('output-candidate-tag-2'));

    await waitFor(() => {
      expect(screen.getByTestId('active-tag-badge')).toHaveTextContent('TAG_40002');
    });

    expect(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    ).toHaveAttribute('aria-pressed', 'true');
  });

  it('marks modbus operational panels as supporting sections', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(await screen.findByTestId('modbus-secondary-panels')).toHaveAttribute(
      'data-emphasis',
      'supporting',
    );
  });

  it('uses the tag chips as the unified selection surface for both targets', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByTestId('output-tag-chips');

    const tagChips = screen.getByTestId('output-tag-chips');
    expect(tagChips).toBeInTheDocument();
    expect(screen.getByTestId('output-candidate-tag-1')).toBeInTheDocument();

    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    // Tag chips remain visible after switching target
    expect(tagChips).toBeInTheDocument();
    // Database panel shows selected tag context
    await waitFor(() => {
      expect(screen.getByTestId('database-selected-tag')).toBeInTheDocument();
    });
  });
});
