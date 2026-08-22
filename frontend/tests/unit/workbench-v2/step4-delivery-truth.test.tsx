import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { DeliveryTruthStrip } from '../../../src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip';
import { hydrateStudioV2DatabaseConnector } from '../../../src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

describe('Step 4 database delivery truth', () => {
  it('shows connector readiness and failed delivery outcome without raw logs', () => {
    const connector = hydrateStudioV2DatabaseConnector({
      id: 'db-main',
      workspace_id: 'workspace-1',
      kind: 'postgres',
      name: 'Line A PG',
      host: '127.0.0.1',
      port: 5432,
      database: 'gateway',
      username: 'gw_writer',
      schema: 'public',
      table: 'sensor_values',
      write_mode: 'insert',
      write_interval_seconds: 5,
      timestamp_column: 'ts',
      status: 'ready',
      last_schema_ensure_at: '2026-03-16T10:03:00Z',
      last_schema_ensure_status: 'success',
      last_write_at: '2026-03-16T10:05:00Z',
      last_write_status: 'failed',
      last_write_error: 'permission denied',
      last_flush_at: '2026-03-16T10:06:00Z',
      last_flush_status: 'failed',
      last_flush_error: 'flush timeout',
      save_state: 'saved',
      created_at: '2026-03-16T10:00:00Z',
      updated_at: '2026-03-16T10:06:00Z',
    });
    connector.save_error = 'autosave failed';

    render(
      <DeliveryTruthStrip
        connector={connector}
      />,
    );

    const truth = screen.getByTestId('database-delivery-truth');
    expect(truth).toHaveTextContent('step4.delivery_truth_title');
    expect(truth).toHaveTextContent('step4.delivery_readiness_ready');
    expect(truth).toHaveTextContent('step4.delivery_schema_success');
    expect(truth).toHaveTextContent('step4.delivery_write_failed');
    expect(truth).toHaveTextContent('permission denied');
    expect(truth).toHaveTextContent('step4.delivery_flush_failed');
    expect(truth).toHaveTextContent('flush timeout');
    expect(truth).not.toHaveTextContent('autosave failed');
  });
});
