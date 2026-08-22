import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { mockDBTargetAPI, renderPage } from './DatalinkWorkbenchOutputStep.testHarness';

describe('DatalinkWorkbench output step', () => {
  describe('database studio schema snapshot', () => {
    it('clears stale table scope when the operator switches connectors', async () => {
      mockDBTargetAPI.listConnectors.mockResolvedValue([
        {
          id: 'connector-1',
          name: 'Main SQLite',
          kind: 'sqlite',
          connection_config: { dsn: '/tmp/target.db' },
          status: 'ready',
          last_check_at: '',
          last_check_error: '',
          enabled: true,
          created_at: '',
          updated_at: '',
        },
        {
          id: 'connector-2',
          name: 'Warehouse PostgreSQL',
          kind: 'postgres',
          connection_config: {
            host: '127.0.0.1',
            port: '5432',
            user: 'postgres',
            database: 'warehouse',
            sslmode: 'disable',
          },
          status: 'ready',
          last_check_at: '',
          last_check_error: '',
          enabled: true,
          created_at: '',
          updated_at: '',
        },
      ]);
      mockDBTargetAPI.listTables.mockImplementation(async (connectorId: string) => {
        if (connectorId === 'connector-2') {
          return [
            {
              schema: 'analytics',
              name: 'metrics',
              columns: [
                { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true, unique: true },
                { name: 'reading', data_type: 'real', nullable: false, primary_key: false, unique: false },
              ],
            },
          ];
        }

        return [
          {
            schema: 'main',
            name: 'sensor_values',
            columns: [
              { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true, unique: true },
              { name: 'value', data_type: 'real', nullable: false, primary_key: false, unique: false },
            ],
          },
        ];
      });

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      await waitFor(() => {
        expect(screen.getByLabelText('workbench.output.database.mapping.table')).toHaveValue(
          'main.sensor_values',
        );
      });
      expect(screen.getByTestId('schema-column-value')).toBeInTheDocument();

      fireEvent.click(screen.getByRole('button', { name: /Warehouse PostgreSQL/i }));

      await waitFor(() => {
        expect(screen.getByLabelText('workbench.output.database.mapping.table')).toHaveValue(
          'analytics.metrics',
        );
      });
      expect(screen.getByTestId('schema-column-reading')).toBeInTheDocument();
      expect(screen.queryByTestId('schema-column-value')).not.toBeInTheDocument();
    });

    it('keeps write-mode and timestamp scope synchronized', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const writeModeSelect = await screen.findByLabelText(
        'workbench.output.database.mapping.writeMode',
      );
      fireEvent.change(writeModeSelect, {
        target: { value: 'upsert' },
      });

      await waitFor(() => {
        expect(
          screen.getByLabelText('workbench.output.database.mapping.timestampColumn'),
        ).toHaveValue('ts');
      });

      fireEvent.change(writeModeSelect, {
        target: { value: 'insert' },
      });

      await waitFor(() => {
        expect(
          screen.queryByLabelText('workbench.output.database.mapping.timestampColumn'),
        ).not.toBeInTheDocument();
      });
    });

    it('does not overwrite manual write-mode edits when a mapping already exists', async () => {
      mockDBTargetAPI.listMappings.mockResolvedValue([
        {
          id: 'db-mapping-1',
          tag_id: 'tag-1',
          connector_id: 'connector-1',
          table_schema: 'main',
          table_name: 'sensor_values',
          column_name: 'value',
          write_mode: 'upsert',
          timestamp_column: 'ts',
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

      const writeModeSelect = await screen.findByLabelText(
        'workbench.output.database.mapping.writeMode',
      );
      await waitFor(() => {
        expect(
          screen.getByLabelText('workbench.output.database.mapping.writeMode'),
        ).toHaveValue('upsert');
      });
      expect(
        screen.getByLabelText('workbench.output.database.mapping.timestampColumn'),
      ).toHaveValue('ts');

      fireEvent.change(writeModeSelect, {
        target: { value: 'insert' },
      });

      await waitFor(() => {
        expect(
          screen.getByLabelText('workbench.output.database.mapping.writeMode'),
        ).toHaveValue('insert');
      });
      expect(
        screen.queryByLabelText('workbench.output.database.mapping.timestampColumn'),
      ).not.toBeInTheDocument();
    });

    it('collapses connector fields until the operator expands connector setup', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      expect(
        screen.queryByLabelText('workbench.output.database.connector.name'),
      ).not.toBeInTheDocument();

      fireEvent.click(
        screen.getByRole('button', {
          name: 'workbench.output.database.actions.configureConnector',
        }),
      );

      expect(
        await screen.findByLabelText('workbench.output.database.connector.name'),
      ).toBeInTheDocument();
    });

    it('keeps schema snapshot in supporting panels while the row planner stays on the main surface', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const secondaryPanels = await screen.findByTestId('database-secondary-panels');
      expect(secondaryPanels).toHaveAttribute('data-emphasis', 'supporting');
      expect(await within(secondaryPanels).findByTestId('schema-snapshot')).toBeInTheDocument();
      expect(screen.getByTestId('database-grouped-row-planner')).toBeInTheDocument();
      expect(within(secondaryPanels).queryByTestId('database-grouped-row-planner')).not.toBeInTheDocument();
    });
  });
});
