import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import {
  mockDevices,
  mockUpdateDeviceMutation,
  renderApp,
  resetFoundationMocks,
} from './DatalinkWorkbenchFoundation.testHarness';

describe('DatalinkWorkbench foundation route', () => {
  beforeEach(() => {
    resetFoundationMocks();
  });

  it('renders compact device rows with endpoint, health, and two primary capability hints', async () => {
    await renderApp();

    const row = screen.getByTestId('device-row-device-1');

    expect(within(row).getByText('192.168.1.10:502')).toBeInTheDocument();
    expect(within(row).getByTestId('device-health-device-1')).toHaveTextContent(
      'workbench.device.card.testPassed',
    );
    expect(
      within(row).getByText('workbench.device.capability.labels.unitId'),
    ).toBeInTheDocument();
    expect(
      within(row).getByText('workbench.device.capability.labels.addressBase'),
    ).toBeInTheDocument();
    expect(
      within(row).queryByText('workbench.device.capability.labels.wordOrder'),
    ).not.toBeInTheDocument();
    expect(
      within(row).queryByText('workbench.device.capability.labels.protocolTraits'),
    ).not.toBeInTheDocument();
    expect(within(row).queryByText('Main line')).not.toBeInTheDocument();
  });

  it('uses protocol-specific connection targets for compact row endpoints', async () => {
    mockDevices.push(
      {
        id: 'device-3',
        name: 'Telemetry Broker',
        description: 'MQTT edge',
        protocol: 'mqtt',
        status: 'active',
        connection_config:
          '{"broker_url":"mqtt://broker.internal:1883","client_id":"edge-gateway","topics":["plant/telemetry"],"use_tls":true,"qos":1}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'device-4',
        name: 'Packaging PLC',
        description: 'MC 3E line',
        protocol: 'mc_3e',
        status: 'active',
        connection_config:
          '{"host":"10.0.0.20","port":5000,"network_no":1,"station_no":2,"data_format":"binary"}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
    );

    await renderApp();

    const mqttRow = screen.getByTestId('device-row-device-3');
    expect(within(mqttRow).getByText('mqtt://broker.internal:1883')).toBeInTheDocument();
    expect(within(mqttRow).queryByText(/QoS 1/)).not.toBeInTheDocument();

    const mcRow = screen.getByTestId('device-row-device-4');
    expect(within(mcRow).getByText('10.0.0.20:5000')).toBeInTheDocument();
    expect(within(mcRow).queryByText('1/2 · binary')).not.toBeInTheDocument();
  });

  it('selects a device from the device step and advances to source planning', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-context-bar')).getByRole('button', {
        name: 'workbench.contextBar.actions.gotoSource',
      }),
    );

    expect(screen.getByRole('tab', { name: /workbench\.steps\.source/ })).toHaveAttribute(
      'aria-selected',
      'true',
    );
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.source.planner.startAddress')).toBeInTheDocument();
      expect(screen.getByLabelText('workbench.source.planner.count')).toBeInTheDocument();
    });
  });

  it('shows a master-detail device panel after a device is selected', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const detailPanel = screen.getByTestId('device-detail-panel');
    expect(detailPanel).toBeInTheDocument();
    expect(within(detailPanel).getByText('Mixer PLC')).toBeInTheDocument();
    expect(within(detailPanel).getByTestId('device-detail-endpoint')).toHaveTextContent(
      '192.168.1.10:502',
    );
    expect(
      within(detailPanel).getByRole('button', {
        name: 'workbench.device.actions.continue',
      }),
    ).toBeInTheDocument();
  });

  it('context bar includes step rail plus primary CTA after selecting a device on step 1', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const contextBar = screen.getByTestId('workbench-context-bar');
    const tabs = within(contextBar).getAllByRole('tab');
    expect(tabs.length).toBe(4);
    expect(
      within(contextBar).getByRole('button', {
        name: 'workbench.contextBar.actions.gotoSource',
      }),
    ).toBeInTheDocument();
  });

  it('keeps the created device selected while the list refreshes', async () => {
    mockDevices.splice(0, mockDevices.length);

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.create' }));
    fireEvent.change(screen.getByLabelText('workbench.device.fields.name'), {
      target: { value: 'Browser Smoke PLC' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.host'), {
      target: { value: '127.0.0.1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.port'), {
      target: { value: '502' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.slaveId'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.timeout'), {
      target: { value: '5' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.save' }));

    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }),
      ).toBeEnabled();
    });
  });

  it('allows clearing device description when editing', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.device.actions.edit' }),
    );
    fireEvent.change(
      screen.getByLabelText('workbench.device.fields.description'),
      { target: { value: '' } },
    );
    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.save' }));

    await waitFor(() => {
      expect(mockUpdateDeviceMutation.mutateAsync).toHaveBeenCalledWith(
        expect.objectContaining({
          id: 'device-1',
          data: expect.objectContaining({
            name: 'Mixer PLC',
            description: '',
          }),
        }),
      );
    });
  });
});
