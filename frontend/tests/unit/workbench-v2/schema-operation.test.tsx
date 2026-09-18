import { beforeEach, describe, expect, it, vi } from 'vitest';
import { StrictMode } from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { SchemaSetupSection } from '../../../src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';
import type { DbConnector, Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: string | { defaultValue?: string }) => (
      typeof options === 'string' ? options : options?.defaultValue ?? key
    ),
  }),
}));

vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: {
    generateSchema: vi.fn(),
  },
}));

vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/WorkspaceRecordingPlanSetupSection', () => ({
  WorkspaceRecordingPlanSetupSection: ({ disabled }: { disabled?: boolean }) => (
    <div data-testid="recording-plan-disabled" data-disabled={disabled ? 'true' : 'false'} />
  ),
}));

interface Deferred<T> {
  promise: Promise<T>;
  resolve: (value: T) => void;
}

function deferred<T>(): Deferred<T> {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>((nextResolve) => {
    resolve = nextResolve;
  });
  return { promise, resolve };
}

function connector(workspaceId: string, overrides: Partial<DbConnector> = {}): DbConnector {
  return {
    ...INITIAL_STATE.db.connector,
    kind: 'postgres',
    workspace_id: workspaceId,
    connector_id: 'connector-1',
    identity_revision: 'identity-1',
    persisted: true,
    save_state: 'saved',
    ...overrides,
  };
}

const schemaResult = {
  connector_id: 'connector-1',
  dry_run: true,
  statements: ['CREATE TABLE old_scope (value integer);'],
  executed: 0,
};

const point: Point = {
  id: 'point-1',
  device_id: 'device-1',
  rule_id: 'rule-1',
  rule_name: 'Holding Registers',
  name: 'temperature',
  address: '40001',
  data_type: 'int16',
  function: 'holding_register',
  width: 1,
  enabled: true,
  skipped: false,
  _rule_scale: 1,
  _rule_offset: 0,
};

const mapping: Mapping = {
  point_id: point.id,
  tag_key: 'line1.temperature',
  display_name: 'Temperature',
  unit: 'C',
  target_type: 'float64',
  scale: 1,
  offset: 0,
  enabled: true,
};

function step4State(currentConnector: DbConnector = connector('workspace-step4')) {
  return {
    ...INITIAL_STATE,
    points: [point],
    mappings: { [point.id]: mapping },
    db: {
      ...INITIAL_STATE.db,
      connector: currentConnector,
      targets: {
        [point.id]: {
          tag_id: 'tag.line1.temperature',
          column_name: 'temperature',
          enabled: true,
          save_state: 'saved' as const,
        },
      },
    },
  };
}

describe('schema operation scope and locking', () => {
  beforeEach(() => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockReset();
  });

  it('ignores a deferred preview after the selected table scope changes', async () => {
    const request = deferred<typeof schemaResult>();
    vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockReturnValueOnce(request.promise);
    const firstConnector = connector('workspace-table', { table: 'table_a' });
    const secondConnector = connector('workspace-table', { table: 'table_b' });
    const { rerender } = render(
      <SchemaSetupSection connector={firstConnector} schemaPreviewSignature="table-a" />,
    );

    fireEvent.click(screen.getByText('step4.schema_preview_btn'));
    rerender(<SchemaSetupSection connector={secondConnector} schemaPreviewSignature="table-b" />);
    await act(async () => {
      request.resolve(schemaResult);
      await request.promise;
      await Promise.resolve();
    });

    await waitFor(() => {
      expect(screen.queryByText(schemaResult.statements[0])).not.toBeInTheDocument();
    });
  });

  it('ignores a deferred preview after the device membership scope changes', async () => {
    const request = deferred<typeof schemaResult>();
    vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockReturnValueOnce(request.promise);
    const currentConnector = connector('workspace-device');
    const { rerender } = render(
      <SchemaSetupSection connector={currentConnector} schemaPreviewSignature="device-a" />,
    );

    fireEvent.click(screen.getByText('step4.schema_preview_btn'));
    rerender(<SchemaSetupSection connector={currentConnector} schemaPreviewSignature="device-b" />);
    await act(async () => {
      request.resolve(schemaResult);
      await request.promise;
      await Promise.resolve();
    });

    await waitFor(() => {
      expect(screen.queryByText(schemaResult.statements[0])).not.toBeInTheDocument();
    });
  });

  it('keeps a pending create locked across a short unmount and remount', async () => {
    const request = deferred<typeof schemaResult>();
    vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockReturnValueOnce(request.promise);
    const currentConnector = connector('workspace-remount');
    const { unmount } = render(<SchemaSetupSection connector={currentConnector} />);

    fireEvent.click(screen.getByText('step4.schema_create_btn'));
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.generateSchema).toHaveBeenCalledTimes(1));
    unmount();

    render(<SchemaSetupSection connector={currentConnector} />);
    expect(screen.getByTestId('schema-operation-pending')).toBeInTheDocument();
    expect(screen.getByTestId('schema-create-btn')).toBeDisabled();
    fireEvent.click(screen.getByTestId('schema-create-btn'));
    expect(studioV2WorkspaceDatabaseAPI.generateSchema).toHaveBeenCalledTimes(1);

    request.resolve(schemaResult);
    await waitFor(() => expect(screen.queryByTestId('schema-operation-pending')).not.toBeInTheDocument());
  });

  it('blocks schema and configuration controls while activation is in progress', async () => {
    const activation = deferred<{ workspace_id: string; results: Array<{ device_id: string; status: 'success'; message: string }> }>();
    const activateWorkspace = vi.fn(() => activation.promise);
    const dispatch = vi.fn();
    render(
      <Step4Database
        state={step4State()}
        dispatch={dispatch}
        activateWorkspace={activateWorkspace}
      />,
    );

    fireEvent.click(screen.getByText('step4.activate_btn'));
    await waitFor(() => expect(screen.getByText('step4.schema_preview_btn')).toBeDisabled());

    expect(screen.getByLabelText('step4.field_host')).toBeDisabled();
    expect(screen.getByTestId('step4-row-group-create')).toBeDisabled();
    expect(screen.getByTestId('btn-disable-all-db-targets')).toBeDisabled();
    expect(screen.getByTestId('step4-target-column-point-1')).toBeDisabled();
    expect(screen.getByTestId('recording-plan-disabled')).toHaveAttribute('data-disabled', 'true');

    activation.resolve({ workspace_id: 'workspace-step4', results: [{ device_id: 'device-1', status: 'success', message: 'activated' }] });
  });

  it('ignores a deferred Step 4 preview after the enabled device points change', async () => {
    const request = deferred<typeof schemaResult>();
    vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockReturnValueOnce(request.promise);
    const initialState = step4State(connector('workspace-step4-scope'));
    const dispatch = vi.fn();
    const { rerender } = render(
      <Step4Database state={initialState} dispatch={dispatch} workspaceId="workspace-step4-scope" />,
    );

    fireEvent.click(screen.getByTestId('schema-preview-btn'));
    const movedPoint: Point = { ...point, device_id: 'device-2' };
    rerender(
      <Step4Database state={{ ...initialState, points: [movedPoint] }} dispatch={dispatch} workspaceId="workspace-step4-scope" />,
    );
    await act(async () => {
      request.resolve(schemaResult);
      await request.promise;
      await Promise.resolve();
    });

    expect(screen.queryByText(schemaResult.statements[0])).not.toBeInTheDocument();
  });

  it('reports that creating tables needs a confirmed preview instead of claiming success', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockRejectedValueOnce({
      response: {
        status: 409,
        data: {
          success: false,
          error: {
            code: 'SCHEMA_CONFIRMATION_REQUIRED',
            message: 'schema creation requires a confirmed preview',
            retryable: false,
            request_id: 'req-1',
            action: 'preview the managed schema and confirm it before creating tables',
          },
        },
      },
    });

    render(<SchemaSetupSection connector={connector('workspace-confirm')} />);
    fireEvent.click(screen.getByTestId('schema-create-btn'));

    const refusal = await screen.findByTestId('schema-confirmation-required');
    expect(refusal).toBeInTheDocument();
    expect(screen.queryByText(/已建立/)).not.toBeInTheDocument();
    expect(screen.getByTestId('schema-create-btn')).not.toBeDisabled();
  });

  it('shows a preview result when mounted under StrictMode', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockResolvedValueOnce(schemaResult);
    render(
      <StrictMode>
        <SchemaSetupSection connector={connector('workspace-strict')} />
      </StrictMode>,
    );

    fireEvent.click(screen.getByTestId('schema-preview-btn'));

    expect(await screen.findByText(schemaResult.statements[0])).toBeInTheDocument();
    expect(screen.getByTestId('schema-preview-btn')).not.toBeDisabled();
  });
});
