import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ConnectorSection } from '../../../src/features/datalink/workbench-v2/steps/step4/ConnectorSection';
import { CommitSummary } from '../../../src/features/datalink/workbench-v2/steps/step4/CommitSummary';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { getColumnsFor } from '../../../src/features/datalink/workbench-v2/state/dbSchemas';
import type { Point, Mapping, DbTarget, DbConnector } from '../../../src/features/datalink/workbench-v2/state/types';
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

    it('資料庫 target 設定變更後應清掉舊的 DDL 預覽', async () => {
      vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockResolvedValue({
        connector_id: 'db-1',
        dry_run: true,
        statements: ['ALTER TABLE sensor_readings ADD COLUMN temp_in_c REAL;'],
        executed: 0,
      });

      const dispatch = vi.fn();
      const initialState = {
        ...INITIAL_STATE,
        points: [
          { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_1', address: '40001', data_type: 'int16', function: 'holding_register' as const, width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
        ],
        mappings: {
          'p-1': { point_id: 'p-1', tag_key: 'line1.t_1', display_name: 'T1', unit: 'C', target_type: 'float64' as const, scale: 1, offset: 0, enabled: true },
        },
        db: {
          connector: {
            ...INITIAL_STATE.db.connector,
            kind: 'postgres' as const,
            save_state: 'saved' as const,
          },
          targets: {
            'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true, save_state: 'saved' as const },
          },
        },
      };

      const { rerender } = render(
        <Step4Database
          state={initialState}
          dispatch={dispatch}
          onCommit={() => { }}
        />
      );

      fireEvent.click(screen.getByText('step4.schema_preview_btn'));

      await waitFor(() => {
        expect(screen.getByText('ALTER TABLE sensor_readings ADD COLUMN temp_in_c REAL;')).toBeInTheDocument();
      });

      rerender(
        <Step4Database
          state={{
            ...initialState,
            db: {
              ...initialState.db,
              targets: {
                'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_out_c', enabled: true, save_state: 'saved' as const },
              },
            },
          }}
          dispatch={dispatch}
          onCommit={() => { }}
        />
      );

      await waitFor(() => {
        expect(screen.queryByText('ALTER TABLE sensor_readings ADD COLUMN temp_in_c REAL;')).not.toBeInTheDocument();
      });
    });

    it('connector 或 target autosave 未完成時應禁用 schema 按鈕', () => {
      const dispatch = vi.fn();
      const state = {
        ...INITIAL_STATE,
        db: {
          connector: {
            ...INITIAL_STATE.db.connector,
            kind: 'postgres' as const,
            save_state: 'saving' as const,
          },
          targets: {
            'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true, save_state: 'saved' as const },
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

      expect(screen.getByText('step4.schema_preview_btn')).toBeDisabled();
      expect(screen.getByText('step4.schema_create_btn')).toBeDisabled();
    });

    it('target 正在停用 autosave 時也應禁用 schema 按鈕', () => {
      const dispatch = vi.fn();
      const state = {
        ...INITIAL_STATE,
        db: {
          connector: {
            ...INITIAL_STATE.db.connector,
            kind: 'postgres' as const,
            save_state: 'saved' as const,
          },
          targets: {
            'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: false, save_state: 'saving' as const },
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

      expect(screen.getByText('step4.schema_preview_btn')).toBeDisabled();
      expect(screen.getByText('step4.schema_create_btn')).toBeDisabled();
    });

    it('只修改 password 後也應清掉舊的 DDL 預覽', async () => {
      vi.mocked(studioV2WorkspaceDatabaseAPI.generateSchema).mockResolvedValue({
        connector_id: 'db-1',
        dry_run: true,
        statements: ['ALTER TABLE sensor_readings ADD COLUMN temp_in_c REAL;'],
        executed: 0,
      });

      const dispatch = vi.fn();
      const initialState = {
        ...INITIAL_STATE,
        db: {
          connector: {
            ...INITIAL_STATE.db.connector,
            kind: 'postgres' as const,
            password: 'secret-a',
            save_state: 'saved' as const,
          },
          targets: {
            'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true, save_state: 'saved' as const },
          },
        },
      };

      const { rerender } = render(
        <Step4Database
          state={initialState}
          dispatch={dispatch}
          onCommit={() => { }}
        />
      );

      fireEvent.click(screen.getByText('step4.schema_preview_btn'));

      await waitFor(() => {
        expect(screen.getByText('ALTER TABLE sensor_readings ADD COLUMN temp_in_c REAL;')).toBeInTheDocument();
      });

      rerender(
        <Step4Database
          state={{
            ...initialState,
            db: {
              ...initialState.db,
              connector: {
                ...initialState.db.connector,
                password: 'secret-b',
              },
            },
          }}
          dispatch={dispatch}
          onCommit={() => { }}
        />
      );

      await waitFor(() => {
        expect(screen.queryByText('ALTER TABLE sensor_readings ADD COLUMN temp_in_c REAL;')).not.toBeInTheDocument();
      });
    });

    it('在 Step 4 選擇 Connector Pool 既有連線時應自動套用連線設定', () => {
      const dispatch = vi.fn();
      const state = {
        ...INITIAL_STATE,
        settings: {
          ...INITIAL_STATE.settings,
          connectors: [
            {
              id: 'conn-mysql-custom',
              name: 'My Custom MySQL',
              kind: 'mysql' as const,
              host: '10.0.0.99',
              port: 3306,
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
      fireEvent.change(poolSelect, { target: { value: 'conn-mysql-custom' } });

      expect(dispatch).toHaveBeenCalledWith(expect.objectContaining({
        type: 'updateDbConnector',
        patch: expect.objectContaining({
          kind: 'mysql',
          name: 'My Custom MySQL',
          host: '10.0.0.99',
          port: 3306,
          database: 'custom_db',
          username: 'custom_root',
        }),
      }));
    });
  });
});
