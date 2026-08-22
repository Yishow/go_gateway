import { fireEvent, screen, waitFor } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { mockModbusShareAPI, renderPage } from './DatalinkWorkbenchOutputStep.testHarness';

describe('DatalinkWorkbench output step', () => {
  describe('modbus studio register map canvas', () => {
    it('renders a register map canvas showing allocated register ranges', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
        { tag_id: 'tag-2', register: 2, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(
        await screen.findByTestId('register-map-canvas'),
      ).toBeInTheDocument();
      expect(screen.getByTestId('register-slot-1')).toHaveTextContent('TAG_40001');
      expect(screen.getByTestId('register-slot-3')).toHaveTextContent('TAG_40002');
    });

    it('highlights conflicting register slots', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 10, data_type: 'int16', updated_at: '' },
        { tag_id: 'tag-2', register: 10, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const conflictSlot = await screen.findByTestId('register-slot-11');
      expect(conflictSlot).toHaveAttribute('data-conflict', 'true');
    });

    it('treats overlapping multi-word ranges as conflicting slots', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 10, data_type: 'int32', updated_at: '' },
        { tag_id: 'tag-2', register: 11, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const overlapSlot = await screen.findByTestId('register-slot-12');
      expect(overlapSlot).toHaveAttribute('data-conflict', 'true');
    });

    it('keeps high mapped registers visible in the canvas', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-2', register: 199, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(await screen.findByTestId('register-slot-200')).toHaveTextContent('TAG_40002');
    });

    it('lets the operator unbind HR1 directly from the canvas', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      fireEvent.click(await screen.findByTestId('register-slot-1'));

      await waitFor(() => {
        expect(mockModbusShareAPI.deleteMapping).toHaveBeenCalledWith('tag-1');
      });
      expect(screen.getByRole('status')).toHaveTextContent('workbench.output.results.slotUnbound');
    });

    it('binds the selected tag by clicking an empty register slot', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const targetSlot = await screen.findByTestId('register-slot-5');
      fireEvent.click(targetSlot);

      await waitFor(() => {
        expect(mockModbusShareAPI.upsertMapping).toHaveBeenCalledWith('tag-1', 4);
      });
      expect(screen.getByRole('status')).toHaveTextContent('workbench.output.results.slotBound');
    });

    it('renders auto-map strategy selector with three strategies', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('register-map-canvas');

      expect(
        screen.getByRole('button', { name: 'workbench.output.modbusStudio.autoMap.sequential' }),
      ).toBeInTheDocument();
      expect(
        screen.getByRole('button', { name: 'workbench.output.modbusStudio.autoMap.gapAware' }),
      ).toBeInTheDocument();
      expect(
        screen.getByRole('button', { name: 'workbench.output.modbusStudio.autoMap.aligned' }),
      ).toBeInTheDocument();
    });

    it('sequential auto-map appends after the current occupied range', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('register-map-canvas');

      fireEvent.click(
        screen.getByRole('button', {
          name: 'workbench.output.modbusStudio.autoMap.sequential',
        }),
      );

      await waitFor(() => {
        expect(mockModbusShareAPI.upsertMapping).toHaveBeenCalledWith('tag-2', 1);
      });
    });

    it('renders a dry-run validation button and displays results', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('register-map-canvas');

      const dryRunButton = screen.getByRole('button', {
        name: 'workbench.output.modbusStudio.actions.dryRun',
      });
      expect(dryRunButton).toBeInTheDocument();

      fireEvent.click(dryRunButton);

      await waitFor(() => {
        expect(screen.getByTestId('dry-run-results')).toBeInTheDocument();
      });
    });
  });
});
