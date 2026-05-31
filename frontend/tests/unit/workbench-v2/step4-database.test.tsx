import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { KindSelector } from '../../../src/features/datalink/workbench-v2/steps/step4/KindSelector';
import { WriteStrategy } from '../../../src/features/datalink/workbench-v2/steps/step4/WriteStrategy';
import { ConnectorSection } from '../../../src/features/datalink/workbench-v2/steps/step4/ConnectorSection';
import { TargetMappingTable } from '../../../src/features/datalink/workbench-v2/steps/step4/TargetMappingTable';
import { CommitSummary } from '../../../src/features/datalink/workbench-v2/steps/step4/CommitSummary';
import { getColumnsFor } from '../../../src/features/datalink/workbench-v2/state/dbSchemas';
import type { Point, Mapping, DbTarget, DbConnector } from '../../../src/features/datalink/workbench-v2/state/types';

// Mock react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: any) => {
      if (options && options.interval !== undefined) {
        return `${key}_interval_${options.interval}`;
      }
      return key;
    }
  })
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
  });

  describe('KindSelector', () => {
    it('應渲染 4 個 db kind 卡片，且點擊時呼叫 onChange', () => {
      const onChange = vi.fn();
      render(<KindSelector value="postgres" onChange={onChange} />);

      expect(screen.getByText('step4.kind_postgres')).toBeInTheDocument();
      expect(screen.getByText('step4.kind_sqlite')).toBeInTheDocument();
      expect(screen.getByText('step4.kind_mysql')).toBeInTheDocument();
      expect(screen.getByText('step4.kind_sqlserver')).toBeInTheDocument();

      const sqliteBtn = screen.getByText('step4.kind_sqlite').closest('button');
      expect(sqliteBtn).not.toBeNull();
      fireEvent.click(sqliteBtn!);
      expect(onChange).toHaveBeenCalledWith('sqlite');
    });

    it('在 disabled 為 true 時應禁用所有按鈕', () => {
      render(<KindSelector value="postgres" onChange={() => {}} disabled={true} />);
      const buttons = screen.getAllByRole('button');
      buttons.forEach(btn => {
        expect(btn).toBeDisabled();
      });
    });
  });

  describe('WriteStrategy', () => {
    it('應能渲染兩種策略 radio 與 interval 輸入，點擊時呼叫對應 callback', () => {
      const onWriteModeChange = vi.fn();
      const onWriteIntervalChange = vi.fn();

      render(
        <WriteStrategy
          writeMode="insert"
          writeIntervalSeconds={5}
          onWriteModeChange={onWriteModeChange}
          onWriteIntervalChange={onWriteIntervalChange}
        />
      );

      const upsertRadio = screen.getByLabelText(/step4.write_mode_upsert/);
      fireEvent.click(upsertRadio);
      expect(onWriteModeChange).toHaveBeenCalledWith('upsert');

      const intervalInput = screen.getByRole('spinbutton');
      fireEvent.change(intervalInput, { target: { value: '10' } });
      expect(onWriteIntervalChange).toHaveBeenCalledWith(10);
    });
  });

  describe('ConnectorSection', () => {
    it('應依 kind 隱藏或顯示特定連線欄位（如 SQLite 隱藏 host/port）', () => {
      const onUpdateConnector = vi.fn();
      const { rerender } = render(
        <ConnectorSection
          connector={mockConnector}
          onUpdateConnector={onUpdateConnector}
          onKindChange={() => {}}
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
          onKindChange={() => {}}
        />
      );

      // SQLite 應隱藏 Host / Port
      expect(screen.queryByText('step4.field_host')).not.toBeInTheDocument();
      expect(screen.queryByText('step4.field_port')).not.toBeInTheDocument();
    });
  });

  describe('TargetMappingTable', () => {
    const mockPoints: Point[] = [
      { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 'SENSOR_1', address: '40001', data_type: 'int16', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
      { id: 'p-2', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 'SENSOR_2', address: '40002', data_type: 'int16', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 }
    ];

    const mockMappings: Record<string, Mapping> = {
      'p-1': { point_id: 'p-1', tag_key: 'line1.temp_in', display_name: 'Temp In', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true },
      'p-2': { point_id: 'p-2', tag_key: 'line1.temp_out', display_name: 'Temp Out', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true }
    };

    const mockTargets: Record<string, DbTarget> = {
      'p-1': { tag_id: 'tag.line1.temp_in', column_name: 'temp_in_c', enabled: true },
      'p-2': { tag_id: 'tag.line1.temp_out', column_name: 'temp_out_c', enabled: true }
    };

    const mockColumns = getColumnsFor('postgres');

    it('應能正確顯示每一列的 point 與對應 select，且在沒有衝突時不渲染 warning banner', () => {
      render(
        <TargetMappingTable
          points={mockPoints}
          mappings={mockMappings}
          targets={mockTargets}
          columns={mockColumns}
          onUpdateTarget={() => {}}
        />
      );

      expect(screen.getByText('line1.temp_in')).toBeInTheDocument();
      expect(screen.getByText('line1.temp_out')).toBeInTheDocument();
      expect(screen.queryByText('step4.conflict_banner_title')).not.toBeInTheDocument();
    });

    it('當兩點對應到同一個欄位時，應標示衝突紅框與 alert icon，並顯示表尾 warning banner', () => {
      const conflictingTargets = {
        'p-1': { tag_id: 'tag.line1.temp_in', column_name: 'temp_in_c', enabled: true },
        // 刻意讓 p-2 也對應到 temp_in_c
        'p-2': { tag_id: 'tag.line1.temp_out', column_name: 'temp_in_c', enabled: true }
      };

      render(
        <TargetMappingTable
          points={mockPoints}
          mappings={mockMappings}
          targets={conflictingTargets}
          columns={mockColumns}
          onUpdateTarget={() => {}}
        />
      );

      // 表尾 banner 應渲染
      expect(screen.getByText('step4.conflict_banner_title')).toBeInTheDocument();
      // 兩個 select 下拉都有顯示警告 emoji 標記
      const warningIcons = screen.getAllByTitle('step4.conflict_tooltip');
      expect(warningIcons.length).toBe(2);
    });

    it('當其中一列 target 被 toggle 停用後，衝突應隨之消失', () => {
      const conflictingTargets = {
        'p-1': { tag_id: 'tag.line1.temp_in', column_name: 'temp_in_c', enabled: true },
        'p-2': { tag_id: 'tag.line1.temp_out', column_name: 'temp_in_c', enabled: false } // p-2 停用
      };

      render(
        <TargetMappingTable
          points={mockPoints}
          mappings={mockMappings}
          targets={conflictingTargets}
          columns={mockColumns}
          onUpdateTarget={() => {}}
        />
      );

      expect(screen.queryByText('step4.conflict_banner_title')).not.toBeInTheDocument();
    });
  });

  describe('CommitSummary', () => {
    it('應渲染 5 列摘要值且 Commit 按鈕啟用；有衝突時按鈕應 disabled', () => {
      const onActivate = vi.fn();
      const { rerender } = render(
        <CommitSummary
          deviceCount={2}
          ruleCount={4}
          pointCount={8}
          mappingCount={8}
          connector={mockConnector}
          enabledTargetCount={8}
          hasConflict={false}
          onActivate={onActivate}
        />
      );

      const btn = screen.getByRole('button');
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
          enabledTargetCount={8}
          hasConflict={true}
          onActivate={onActivate}
        />
      );

      expect(screen.getByRole('button')).toBeDisabled();
    });
  });
});
