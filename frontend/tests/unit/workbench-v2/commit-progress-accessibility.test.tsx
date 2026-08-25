import { act, render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CommitProgress } from '../../../src/features/datalink/workbench-v2/steps/step4/CommitProgress';
import { Toggle } from '../../../src/features/datalink/workbench-v2/components/Toggle';
import type { CommitLog } from '../../../src/features/datalink/workbench-v2/state/types';

describe('CommitProgress presentation and accessibility', () => {
  it('truthfully renders success, failed, and skipped commit rows and never renders 200 on failure', () => {
    const logs: CommitLog[] = [
      {
        label: 'Start Scheduler',
        detail: 'Scheduler started successfully',
        status: 'success',
      },
      {
        label: 'Activate Modbus Share',
        detail: 'Port 5020 in use',
        status: 'failed',
        code: 'modbus_share_listener_bind_failed',
        action: 'Choose an available port in settings',
      },
      {
        label: 'Database Sync',
        detail: 'Database connector disabled',
        status: 'skipped',
      },
    ];

    render(<CommitProgress logs={logs} status="failed" />);

    const row0 = screen.getByTestId('commit-log-row-0');
    expect(row0).toHaveTextContent('200');
    expect(row0).toHaveTextContent('Start Scheduler');

    const row1 = screen.getByTestId('commit-log-row-1');
    expect(row1).toHaveTextContent('ERR');
    expect(row1).not.toHaveTextContent('200');
    expect(row1).not.toHaveTextContent('Port 5020 in use');
    expect(row1).not.toHaveTextContent('Choose an available port in settings');
    expect(row1).toHaveTextContent('errors.generic_failure');

    const row2 = screen.getByTestId('commit-log-row-2');
    expect(row2).toHaveTextContent('SKIP');
    expect(row2).not.toHaveTextContent('200');
  });

  it('does not infer success when a result status is missing', () => {
    render(
      <CommitProgress
        logs={[{ label: 'Unknown operation', detail: 'raw backend detail' }]}
        status="failed"
      />,
    );

    const row = screen.getByTestId('commit-log-row-0');
    expect(row).toHaveAttribute('data-status', 'pending');
    expect(row).not.toHaveTextContent('200');
    expect(row).not.toHaveTextContent('raw backend detail');
  });

  it('renders an executable localized retry control with single-flight pending state', async () => {
    let resolveRetry!: () => void;
    const onRetry = vi.fn(
      () => new Promise<void>((resolve) => {
        resolveRetry = resolve;
      }),
    );

    render(
      <CommitProgress
        logs={[{
          label: 'Activate device',
          detail: 'raw backend detail',
          status: 'failed',
          code: 'activation_failed',
        }]}
        status="failed"
        onRetry={onRetry}
      />,
    );

    const retry = screen.getByRole('button', { name: 'step4.retry_action' });
    fireEvent.click(retry);
    fireEvent.click(retry);
    expect(onRetry).toHaveBeenCalledTimes(1);
    expect(retry).toBeDisabled();
    expect(screen.getByTestId('commit-log-row-0')).not.toHaveTextContent('raw backend detail');

    await act(async () => {
      resolveRetry();
      await Promise.resolve();
    });
    expect(retry).not.toBeDisabled();
  });
});

describe('Toggle accessibility and single-flight control', () => {
  it('renders switch role and handles keyboard / single click state', () => {
    const onChange = vi.fn();
    const { rerender } = render(
      <Toggle checked={false} onChange={onChange} label="Share Toggle" />,
    );

    const toggle = screen.getByRole('switch', { name: 'Share Toggle' });
    expect(toggle).toHaveAttribute('aria-checked', 'false');

    fireEvent.click(toggle);
    expect(onChange).toHaveBeenCalledTimes(1);
    expect(onChange).toHaveBeenCalledWith(true);

    // Disabled / pending state
    rerender(
      <Toggle checked={true} onChange={onChange} label="Share Toggle" disabled={true} />,
    );
    expect(toggle).toHaveAttribute('aria-checked', 'true');
    expect(toggle).toBeDisabled();

    fireEvent.click(toggle);
    expect(onChange).toHaveBeenCalledTimes(1); // not called again
  });
});
