import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import {
  mockDevices,
  mockTestConnectionMutation,
  renderApp,
  resetFoundationMocks,
} from './DatalinkWorkbenchFoundation.testHarness';

describe('DatalinkWorkbench foundation route', () => {
  beforeEach(() => {
    resetFoundationMocks();
  });

  it('hides zero-value test timestamps in the inspector', async () => {
    mockDevices[0].last_test_at = '0001-01-01T00:00:00Z';

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(
      screen.getByText('workbench.device.inspector.unknownTestTime'),
    ).toBeInTheDocument();
    expect(screen.queryByText('0001-01-01T00:00:00Z')).not.toBeInTheDocument();
  });

  it('renders selected device details in the right-side inspector', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const inspector = screen.getByTestId('workbench-inspector-panel');
    expect(within(inspector).getByText('Mixer PLC')).toBeInTheDocument();
    expect(within(inspector).getByText('workbench.device.inspector.connectionSummary')).toBeInTheDocument();
    expect(within(inspector).getByRole('button', { name: 'workbench.device.actions.edit' })).toBeInTheDocument();
    expect(within(inspector).getByRole('button', { name: 'workbench.device.actions.testConnection' })).toBeInTheDocument();
    expect(within(inspector).getByRole('button', { name: 'workbench.device.actions.clone' })).toBeInTheDocument();
  });

  it('disables inspector test action while a connection test is already pending', async () => {
    mockTestConnectionMutation.isPending = true;

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    const inspector = screen.getByTestId('workbench-inspector-panel');
    expect(
      within(inspector).getByRole('button', {
        name: 'workbench.device.actions.testing',
      }),
    ).toBeDisabled();
  });

  it('keeps the latest session test result in the inspector while the context bar stays compact', async () => {
    mockTestConnectionMutation.mutateAsync
      .mockReset()
      .mockResolvedValueOnce({ success: false, error: 'timeout-latest', latency_ms: 0 });

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-inspector-panel')).getByRole('button', {
        name: 'workbench.device.actions.testConnection',
      }),
    );

    await waitFor(() => {
      const inspectorEntry = within(screen.getByTestId('workbench-inspector-panel')).getByText(
        'timeout-latest',
      );
      expect(inspectorEntry).toBeInTheDocument();
      expect(inspectorEntry).toHaveClass('text-rose-300');
      expect(screen.queryByTestId('context-bar-test-status')).not.toBeInTheDocument();
    });
  });

  it('keeps only the three most recent connection tests in the inspector timeline', async () => {
    mockTestConnectionMutation.mutateAsync
      .mockReset()
      .mockResolvedValueOnce({ success: false, error: 'timeout-1', latency_ms: 0 })
      .mockResolvedValueOnce({ success: true, error: '', latency_ms: 14 })
      .mockResolvedValueOnce({ success: false, error: 'crc-2', latency_ms: 0 })
      .mockResolvedValueOnce({ success: false, error: 'offline-3', latency_ms: 0 });

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    const inspector = screen.getByTestId('workbench-inspector-panel');
    const testButton = within(inspector).getByRole('button', {
      name: 'workbench.device.actions.testConnection',
    });

    fireEvent.click(testButton);
    await waitFor(() => {
      expect(within(inspector).getByText('timeout-1')).toBeInTheDocument();
    });

    fireEvent.click(testButton);
    await waitFor(() => {
      expect(
        within(inspector).getByText('workbench.device.messages.testSuccess'),
      ).toBeInTheDocument();
    });

    fireEvent.click(testButton);
    await waitFor(() => {
      expect(within(inspector).getByText('crc-2')).toBeInTheDocument();
    });

    fireEvent.click(testButton);
    await waitFor(() => {
      expect(within(inspector).getByText('offline-3')).toBeInTheDocument();
      expect(within(inspector).getAllByRole('listitem')).toHaveLength(3);
    });

    expect(within(inspector).queryByText('timeout-1')).not.toBeInTheDocument();
  });

  it('shows connect and probe phase diagnostics separately in the inspector timeline', async () => {
    mockTestConnectionMutation.mutateAsync
      .mockReset()
      .mockResolvedValueOnce({
        success: false,
        error: 'probe failed',
        latency_ms: 21,
        can_activate: false,
        can_collect: false,
        connect: {
          status: 'success',
          message: 'TCP ready',
          latency_ms: 8,
        },
        probe: {
          status: 'failed',
          error: 'CRC mismatch',
          latency_ms: 13,
        },
      });

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    const inspector = screen.getByTestId('workbench-inspector-panel');

    fireEvent.click(
      within(inspector).getByRole('button', {
        name: 'workbench.device.actions.testConnection',
      }),
    );

    await waitFor(() => {
      const timelineEntry = within(inspector).getByText('probe failed').closest('li');
      expect(timelineEntry).not.toBeNull();
      expect(within(timelineEntry as HTMLElement).getByText('workbench.device.inspector.phases.connect')).toBeInTheDocument();
      expect(within(timelineEntry as HTMLElement).getByText('TCP ready')).toBeInTheDocument();
      expect(within(timelineEntry as HTMLElement).getByText('workbench.device.inspector.phases.probe')).toBeInTheDocument();
      expect(within(timelineEntry as HTMLElement).getByText('CRC mismatch')).toBeInTheDocument();
      expect(within(timelineEntry as HTMLElement).getByText('workbench.device.inspector.activationBlocked')).toBeInTheDocument();
    });
  });

  it('announces device notices through a polite live region', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.refresh' }));

    const statusRegions = await screen.findAllByRole('status');
    const liveNotice = statusRegions.find((element) =>
      element.textContent?.includes('workbench.device.messages.refreshed'),
    );

    expect(liveNotice).toHaveAttribute('aria-live', 'polite');
  });
});
