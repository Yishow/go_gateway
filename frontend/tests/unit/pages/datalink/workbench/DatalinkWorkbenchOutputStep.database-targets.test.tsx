import { fireEvent, screen, waitFor } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { mockDBTargetAPI, renderPage } from './DatalinkWorkbenchOutputStep.testHarness';

describe('DatalinkWorkbench output step', () => {
  it('binds a linked tag to a database column via direct surface click', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    // Wait for schema columns to load
    await screen.findByTestId('schema-column-surface');

    // Select a table first
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.output.database.mapping.table')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText('workbench.output.database.mapping.table'), {
      target: { value: 'main.sensor_values' },
    });

    // Click on the 'value' column to bind tag-1
    await waitFor(() => {
      expect(screen.getByTestId('schema-column-value')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('schema-column-value'));

    await waitFor(() => {
      expect(mockDBTargetAPI.createMapping).toHaveBeenCalledWith(
        expect.objectContaining({
          tag_id: 'tag-1',
          connector_id: 'connector-1',
          table_name: 'sensor_values',
          column_name: 'value',
        }),
      );
    });
    expect(screen.getByRole('status')).toHaveTextContent(
      'workbench.output.database.results.mappingSaved',
    );
  });

  it('removes a database mapping with inline feedback when clicking a bound column', async () => {
    mockDBTargetAPI.listMappings.mockResolvedValue([
      {
        id: 'db-mapping-1',
        tag_id: 'tag-1',
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
    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    await screen.findByTestId('schema-column-value');
    fireEvent.click(screen.getByTestId('schema-column-value'));

    await waitFor(() => {
      expect(mockDBTargetAPI.deleteMapping).toHaveBeenCalledWith('db-mapping-1');
    });
    expect(screen.getByRole('status')).toHaveTextContent(
      'workbench.output.database.results.mappingDeleted',
    );
  });

  it('keeps the database selected-tag display aligned with the tag chips', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    // Select tag-2 via chip
    fireEvent.click(screen.getByTestId('output-candidate-tag-2'));

    await waitFor(() => {
      expect(screen.getByTestId('database-selected-tag')).toHaveTextContent('TAG_40002');
    });
  });
});
