import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ConnectorSection } from '../../../src/features/datalink/workbench-v2/steps/step4/ConnectorSection';
import { CommitSummary } from '../../../src/features/datalink/workbench-v2/steps/step4/CommitSummary';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import type { DbConnector } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';

// Mock react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: { interval?: number }) => {
      if (options && options.interval !== undefined) {
        return `${key}_interval_${options.interval}`;
      }
      return key;
    }
  })
}));

vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: {
    generateSchema: vi.fn(),
  },
}));

describe('Step 4 Database UI Components & Integration', () => {
  const mockConnector: DbConnector = {
    kind: 'postgres',
    name: 'PostgreSQL Connector',
    host: 'tsdb.internal',
    port: 5432,
    database: 'gateway_metrics',
    username: 'postgres',
    schema: 'public',
    table: 'sensor_readings',
    write_mode: 'insert',
    write_interval_seconds: 5,
    timestamp_column: 'ts',
    status: 'unknown'
  };

  beforeEach(() => {
    vi.restoreAllMocks();
    vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockReset();
  });


  describe('ConnectorSection', () => {
    it('應依 kind 隱藏或顯示特定連線欄位（如 SQLite 隱藏 host/port）', () => {
      const onUpdateConnector = vi.fn();
      const { rerender } = render(
        <ConnectorSection
          connector={mockConnector}
          onUpdateConnector={onUpdateConnector}
          onKindChange={() => { }}
        />
      );

      // Postgres 應顯示 Host / Port 欄位
      expect(screen.getByText('step4.field_host')).toBeInTheDocument();
      expect(screen.getByText('step4.field_port')).toBeInTheDocument();

      // 切換為 sqlite
      const sqliteConnector: DbConnector = {
        ...mockConnector,
        kind: 'sqlite',
        host: '',
        port: 0
      };

      rerender(
        <ConnectorSection
          connector={sqliteConnector}
          onUpdateConnector={onUpdateConnector}
          onKindChange={() => { }}
        />
      );

      // SQLite 應隱藏 Host / Port
      expect(screen.queryByText('step4.field_host')).not.toBeInTheDocument();
      expect(screen.queryByText('step4.field_port')).not.toBeInTheDocument();
    });
  });

  describe('CommitSummary', () => {
    it('Commit 只依後端 readiness：無阻擋時啟用，有阻擋時 disabled，且不以舊 target 數量當門檻', () => {
      const onActivate = vi.fn();
      const { rerender } = render(
        <CommitSummary
          deviceCount={2}
          ruleCount={4}
          pointCount={8}
          mappingCount={8}
          connector={mockConnector}
          groupSummary={{ state: 'ready', total: 1, applied: 0 }}
          onActivate={onActivate}
        />
      );

      const btn = screen.getByRole('button', { name: /step4\.activate_btn/ });
      expect(btn).not.toBeDisabled();
      fireEvent.click(btn);
      expect(onActivate).toHaveBeenCalledTimes(1);

      // 當有衝突時
      rerender(
        <CommitSummary
          deviceCount={2}
          ruleCount={4}
          pointCount={8}
          mappingCount={8}
          connector={mockConnector}
          groupSummary={{ state: 'ready', total: 0, applied: 0 }}
          readinessSummary={{ ready: false, blocking_count: 1, warning_count: 0, issues: [{ code: 'x', severity: 'blocking', step: 'Step 4', scope: 's', message: 'm' }] }}
          onActivate={onActivate}
        />
      );

      expect(screen.getByRole('button', { name: /step4\.activate_btn/ })).toBeDisabled();
    });
  });

  describe('Step4Database', () => {
    it('切換到 sqlite 時應明確清掉 connector password', () => {
      const dispatch = vi.fn();
      const state = {
        ...INITIAL_STATE,
        db: {
          ...INITIAL_STATE.db,
          connector: {
            ...INITIAL_STATE.db.connector,
            kind: 'postgres' as const,
            password: 'stale-secret',
          },
        },
      };

      render(
        <Step4Database
          state={state}
          dispatch={dispatch}
          onCommit={() => { }}
        />
      );

      dispatch.mockClear();
      fireEvent.click(screen.getByText('step4.kind_sqlite'));

      expect(dispatch).toHaveBeenCalledWith(expect.objectContaining({
        type: 'updateDbConnector',
        patch: expect.objectContaining({
          kind: 'sqlite',
          password: undefined,
        }),
      }));
    });

    it('主線不再提供 legacy non-preview DDL 或重複 table/interval 設定', () => {
      const state = { ...INITIAL_STATE, db: { ...INITIAL_STATE.db, connector: { ...INITIAL_STATE.db.connector, kind: 'postgres' as const, save_state: 'saved' as const } } };
      const { rerender } = render(<Step4Database state={state} dispatch={vi.fn()} />);
      expect(screen.queryByText('step4.schema_preview_btn')).not.toBeInTheDocument();
      expect(screen.queryByText('step4.schema_create_btn')).not.toBeInTheDocument();
      expect(screen.queryByLabelText('step4.field_table')).not.toBeInTheDocument();
      rerender(<Step4Database state={{ ...state, db: { ...state.db, connector: { ...state.db.connector, password: 'updated-secret', save_state: 'saving' as const } } }} dispatch={vi.fn()} />);
      expect(screen.queryByText('step4.schema_create_btn')).not.toBeInTheDocument();
      expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument();
    });

    it('在 Step 4 選擇 Connector Pool 既有連線時應自動套用連線設定', () => {
      const dispatch = vi.fn();
      const state = {
        ...INITIAL_STATE,
        settings: {
          ...INITIAL_STATE.settings,
          connectors: [
            {
              id: 'conn-postgres-custom',
              name: 'My Custom PostgreSQL',
              kind: 'postgres' as const,
              host: '10.0.0.99',
              port: 5439,
              database: 'custom_db',
              username: 'custom_root',
              password: 'p',
              schema: '',
              table: 'custom_table',
              enabled: true,
              status: 'ready' as const,
              default_write_interval_seconds: 5,
            },
          ],
        },
      };

      render(
        <Step4Database
          state={state}
          dispatch={dispatch}
          onCommit={() => { }}
        />
      );

      const poolSelect = screen.getByLabelText('step4.load_from_pool');
      fireEvent.change(poolSelect, { target: { value: 'conn-postgres-custom' } });

      expect(dispatch).toHaveBeenCalledWith(expect.objectContaining({
        type: 'updateDbConnector',
        patch: expect.objectContaining({
          kind: 'postgres',
          name: 'My Custom PostgreSQL',
          host: '10.0.0.99',
          port: 5439,
          database: 'custom_db',
          username: 'custom_root',
        }),
      }));
    });
  });
});
