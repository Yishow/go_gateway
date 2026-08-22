import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { mockDBTargetAPI, mockModbusShareAPI, renderPage } from './DatalinkWorkbenchOutputStep.testHarness';

describe('DatalinkWorkbench output step', () => {
  describe('output inspector traceability', () => {
    it('shows source→tag→output trace when an output candidate is selected', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

      await waitFor(() => {
        expect(screen.getByTestId('inspector-trace-panel')).toBeInTheDocument();
      });

      expect(screen.getByTestId('trace-source-address')).toHaveTextContent('40001');
      expect(screen.getByTestId('trace-tag-key')).toHaveTextContent('TAG_40001');
      expect(screen.getByTestId('trace-output-modbus')).toHaveTextContent('HR1');
    });

    it('shows readiness reasons for partial output candidate', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);
      mockDBTargetAPI.listMappings.mockResolvedValue([]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

      await waitFor(() => {
        expect(screen.getByTestId('inspector-trace-panel')).toBeInTheDocument();
      });

      expect(screen.getByTestId('trace-readiness')).toHaveTextContent(
        'workbench.output.inspector.readiness.partial',
      );
    });

    it('shows ready readiness when both modbus and database are mapped', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);
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

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

      await waitFor(() => {
        expect(screen.getByTestId('inspector-trace-panel')).toBeInTheDocument();
      });

      expect(screen.getByTestId('trace-readiness')).toHaveTextContent(
        'workbench.output.inspector.readiness.ready',
      );
    });

    it('shows database path in trace when database mapping exists', async () => {
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

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

      await waitFor(() => {
        expect(screen.getByTestId('inspector-trace-panel')).toBeInTheDocument();
      });

      expect(screen.getByTestId('trace-output-database')).toHaveTextContent(
        'main.sensor_values.value',
      );
    });

    it('announces inspector mapping load failures through a polite live region', async () => {
      mockModbusShareAPI.listMappings.mockRejectedValue(new Error('boom'));
      mockDBTargetAPI.listMappings.mockResolvedValue([]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

      const tracePanel = await screen.findByTestId('inspector-trace-panel');
      const loadError = await within(tracePanel).findByText('boom');

      expect(loadError).toHaveAttribute('role', 'status');
      expect(loadError).toHaveAttribute('aria-live', 'polite');
    });
  });
});
