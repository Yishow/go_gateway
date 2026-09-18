import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Mapping, Point, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: string | { defaultValue?: string }) => (
      typeof options === 'string' ? options : options?.defaultValue ?? key
    ),
  }),
}));
vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: { getMetadata: vi.fn(), generateSchema: vi.fn() },
}));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/WorkspaceRecordingPlanSetupSection', () => ({
  WorkspaceRecordingPlanSetupSection: () => null,
}));

const point: Point = {
  id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 'reactor', address: '40001',
  data_type: 'float32', function: 'holding_register', width: 2, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0,
};
const mapping: Mapping = {
  point_id: 'p-1', tag_key: 'line1.reactor_temp', display_name: 'Reactor', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true,
};

function savedState(identityRevision: string): WorkbenchV2State {
  return {
    ...INITIAL_STATE,
    points: [point],
    mappings: { [point.id]: mapping },
    db: {
      ...INITIAL_STATE.db,
      connector: {
        ...INITIAL_STATE.db.connector, kind: 'sqlite', name: 'Line A', database: '/tmp/line-a.db', schema: 'main',
        table: 'sensor_values', connector_id: 'db-1', identity_revision: identityRevision, workspace_id: 'ws-1',
        persisted: true, save_state: 'saved',
      },
      targets: { [point.id]: { tag_id: 'tag.line1.reactor_temp', column_name: 'reactor_temp', enabled: true, save_state: 'saved' } },
    },
  };
}

function metadata(identityRevision: string, overrides: Record<string, unknown> = {}) {
  return {
    workspace_id: 'ws-1', connector_id: 'db-1', connector_revision: identityRevision, database: '/tmp/line-a.db',
    schema: 'main', table: 'sensor_values', inspection_status: 'exists',
    columns: [
      { name: 'recorded_at', data_type: 'TEXT', nullable: false, primary_key: true },
      { name: 'reactor_temp', data_type: 'REAL', nullable: false, primary_key: false },
    ],
    ...overrides,
  };
}

describe('Step 4 uses real target metadata', () => {
  beforeEach(() => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.getMetadata).mockReset();
  });

  it('offers the inspected columns instead of sample columns', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.getMetadata).mockResolvedValue(metadata('identity-real') as never);
    render(<Step4Database state={savedState('identity-real')} dispatch={vi.fn()} workspaceId="ws-1" />);

    await waitFor(() => expect(document.querySelector('option[value="reactor_temp"]')).not.toBeNull());
    expect(document.querySelector('option[value="temp_in_c"]')).toBeNull();
    expect(studioV2WorkspaceDatabaseAPI.getMetadata).toHaveBeenCalledWith('identity-real');
    expect(screen.getByTestId('step4-target-metadata-status')).toHaveAttribute('data-inspection-status', 'exists');
  });

  it('reports a failed inspection without offering sample columns', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.getMetadata).mockResolvedValue(
      metadata('identity-failed', { inspection_status: 'failed', reason: 'connection_failed', columns: [] }) as never,
    );
    render(<Step4Database state={savedState('identity-failed')} dispatch={vi.fn()} workspaceId="ws-1" />);

    await waitFor(() => expect(screen.getByTestId('step4-target-metadata-status')).toHaveAttribute('data-inspection-status', 'failed'));
    expect(document.querySelector('option[value="temp_in_c"]')).toBeNull();
  });
});
