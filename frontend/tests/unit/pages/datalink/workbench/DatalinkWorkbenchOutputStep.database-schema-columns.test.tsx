import { fireEvent, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { mockDBTargetAPI, renderPage } from './DatalinkWorkbenchOutputStep.testHarness';

describe('DatalinkWorkbench output step', () => {
  describe('database studio schema snapshot', () => {
    it('renders a schema snapshot with column type badges', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      expect(
        await screen.findByTestId('schema-snapshot'),
      ).toBeInTheDocument();
      expect(screen.getByTestId('schema-column-ts')).toHaveTextContent('datetime');
      expect(screen.getByTestId('schema-column-value')).toHaveTextContent('real');
    });

    it('marks primary key columns with a PK badge', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const tsColumn = await screen.findByTestId('schema-column-ts');
      expect(tsColumn).toHaveTextContent('PK');
    });

    it('renders the grouped row planner surface', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      expect(
        await screen.findByTestId('database-grouped-row-planner'),
      ).toBeInTheDocument();
    });

    it('highlights unmapped required columns in the schema snapshot', async () => {
      mockDBTargetAPI.listTables.mockResolvedValue([
        {
          schema: 'main',
          name: 'sensor_values',
          columns: [
            { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true, unique: true },
            { name: 'value', data_type: 'real', nullable: false, primary_key: false, unique: false },
            { name: 'source_id', data_type: 'text', nullable: false, primary_key: false, unique: false },
          ],
        },
      ]);
      mockDBTargetAPI.listMappings.mockResolvedValue([]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const requiredColumn = await screen.findByTestId('schema-column-source_id');
      expect(requiredColumn).toHaveAttribute('data-required', 'true');
    });
  });
});
