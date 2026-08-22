import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import {
  mockCreateDeviceMutation,
  mockDevices,
  mockTestDraftConnectionMutation,
  renderApp,
  resetFoundationMocks,
} from './DatalinkWorkbenchFoundation.testHarness';

describe('DatalinkWorkbench foundation route', () => {
  beforeEach(() => {
    resetFoundationMocks();
  });

  it('opens the create device form inline when there are no devices yet', async () => {
    mockDevices.splice(0, mockDevices.length);

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.create' }));

    expect(screen.getByTestId('device-inline-editor')).toBeInTheDocument();
    expect(screen.queryByTestId('device-panel-overlay')).not.toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'workbench.device.panel.createTitle' }),
    ).toBeInTheDocument();
  });

  it('localizes protocol-specific connection option labels in the create drawer', async () => {
    mockDevices.splice(0, mockDevices.length);

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.create' }));

    fireEvent.change(screen.getByLabelText('workbench.device.fields.protocol'), {
      target: { value: 'modbus_rtu' },
    });

    const paritySelect = screen.getByLabelText('workbench.device.connection.parity');
    expect(within(paritySelect).getByRole('option', { name: 'device.parityNone' })).toBeInTheDocument();
    expect(within(paritySelect).getByRole('option', { name: 'device.parityEven' })).toBeInTheDocument();
    expect(within(paritySelect).getByRole('option', { name: 'device.parityOdd' })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('workbench.device.fields.protocol'), {
      target: { value: 'mc_3e' },
    });

    const dataFormatSelect = screen.getByLabelText('workbench.device.connection.dataFormat');
    expect(
      within(dataFormatSelect).getByRole('option', {
        name: 'workbench.device.connection.dataFormats.abcd',
      }),
    ).toBeInTheDocument();
    expect(
      within(dataFormatSelect).getByRole('option', {
        name: 'workbench.device.connection.dataFormats.cdab',
      }),
    ).toBeInTheDocument();
  });

  it('shows the inline editor in the detail column while keeping the device list visible', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.device.actions.edit' }),
    );

    expect(screen.getByTestId('device-inline-editor')).toBeInTheDocument();
    expect(screen.queryByTestId('device-panel-overlay')).not.toBeInTheDocument();
    expect(screen.queryByTestId('device-detail-panel')).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Mixer PLC' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Backup PLC' })).toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'workbench.device.panel.editTitle' }),
    ).toBeInTheDocument();
  });

  it('tests the current draft config from the inline editor and explains that the backend host runs the dial', async () => {
    mockTestDraftConnectionMutation.mutateAsync.mockResolvedValueOnce({
      success: true,
      error: '',
      latency_ms: 18,
      can_activate: true,
      can_collect: true,
      connect: {
        status: 'success',
        message: 'connect ok',
        latency_ms: 7,
      },
      probe: {
        status: 'success',
        message: 'probe ok',
        latency_ms: 11,
      },
    });

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.device.actions.edit' }),
    );

    expect(
      screen.getByText('workbench.device.connection.backendHostHint'),
    ).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('workbench.device.connection.host'), {
      target: { value: '10.0.0.77' },
    });
    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.device.actions.testDraftConnection',
      }),
    );

    await waitFor(() => {
      expect(mockTestDraftConnectionMutation.mutateAsync).toHaveBeenCalledWith({
        protocol: 'modbus_tcp',
        connection_config: {
          host: '10.0.0.77',
          port: 502,
          slave_id: 1,
          timeout: 5,
        },
      });
    });
  });

  it('opens a clone drawer with connection defaults but requires a new device name', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-inspector-panel')).getByRole('button', {
        name: 'workbench.device.actions.clone',
      }),
    );

    expect(screen.getByTestId('device-inline-editor')).toBeInTheDocument();
    expect(screen.queryByTestId('device-panel-overlay')).not.toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'workbench.device.panel.cloneTitle' }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText('workbench.device.fields.name')).toHaveValue('');
    expect(screen.getByLabelText('workbench.device.connection.host')).toHaveValue('192.168.1.10');
    expect(screen.getByLabelText('workbench.device.connection.port')).toHaveValue('502');
    expect(screen.getByLabelText('workbench.device.connection.slaveId')).toHaveValue('1');
  });

  it('prevents saving a clone with the same source device name', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-inspector-panel')).getByRole('button', {
        name: 'workbench.device.actions.clone',
      }),
    );

    fireEvent.change(screen.getByLabelText('workbench.device.fields.name'), {
      target: { value: 'Mixer PLC' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.save' }));

    await waitFor(() => {
      expect(
        screen.getByText('workbench.device.validation.cloneNameDistinct'),
      ).toBeInTheDocument();
    });
    expect(mockCreateDeviceMutation.mutateAsync).not.toHaveBeenCalled();
  });
});
