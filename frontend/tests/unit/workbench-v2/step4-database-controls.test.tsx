import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { TargetMappingTable } from '../../../src/features/datalink/workbench-v2/steps/step4/TargetMappingTable';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { getColumnsFor } from '../../../src/features/datalink/workbench-v2/state/dbSchemas';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { DbTarget, Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: { interval?: number }) => {
      if (options && options.interval !== undefined) {
        return `${key}_interval_${options.interval}`;
      }
      return key;
    },
  }),
}));

describe('Step 4 Database controls', () => {
  const points: Point[] = [
    { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_1', address: '40001', data_type: 'int16', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
    { id: 'p-2', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_2', address: '40002', data_type: 'int16', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
  ];
  const mappings: Record<string, Mapping> = {
    'p-1': { point_id: 'p-1', tag_key: 'line1.t_1', display_name: 'T1', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    'p-2': { point_id: 'p-2', tag_key: 'line1.t_2', display_name: 'T2', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true },
  };
  const targets: Record<string, DbTarget> = {
    'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true },
    'p-2': { tag_id: 'tag.line1.t_2', column_name: 'temp_out_c', enabled: false },
  };

  it('支援一鍵啟用或停用所有資料表欄位', () => {
    render(
      <TargetMappingTable
        points={points}
        mappings={mappings}
        targets={targets}
        columns={getColumnsFor('postgres')}
        onUpdateTarget={vi.fn()}
        onSetAllEnabled={() => { }}
      />
    );

    expect(screen.getByTestId('step4-target-count')).toHaveTextContent('2');
  });

  it('寫入欄位可自訂，且離開欄位後才提交更新', () => {
    const onUpdateTarget = vi.fn();
    render(
      <TargetMappingTable
        points={points}
        mappings={mappings}
        targets={targets}
        columns={getColumnsFor('postgres')}
        onUpdateTarget={onUpdateTarget}
      />
    );

    const input = screen.getByTestId('step4-target-column-p-1');
    expect(input).toHaveValue('temp_in_c');

    fireEvent.change(input, { target: { value: 'custom_output_c' } });
    expect(onUpdateTarget).not.toHaveBeenCalled();

    fireEvent.blur(input);
    expect(onUpdateTarget).toHaveBeenCalledWith('p-1', { column_name: 'custom_output_c' });
  });

  it('預設顯示資料庫設定，交付狀態卡片維持收合', () => {
    render(
      <Step4Database
        state={{
          ...INITIAL_STATE,
          points: [points[0]],
          mappings: { 'p-1': mappings['p-1'] },
          db: {
            connector: {
              ...INITIAL_STATE.db.connector,
              kind: 'postgres',
              schema: 'public',
              table: 'sensor_readings',
              write_mode: 'insert',
              write_interval_seconds: 5,
            },
            targets: {
              'p-1': { ...targets['p-1'], save_state: 'saved' },
            },
          },
        }}
        dispatch={vi.fn()}
        onCommit={() => { }}
      />
    );

    expect(screen.getByTestId('step4-destination-overview')).toHaveTextContent('public.sensor_readings');
    expect(screen.getByLabelText('step4.field_host')).toBeInTheDocument();
    expect(screen.queryByTestId('database-delivery-truth')).not.toBeInTheDocument();
  });

  it('交付狀態卡片展開後才顯示 delivery 細節', () => {
    render(
      <Step4Database
        state={{
          ...INITIAL_STATE,
          db: {
            ...INITIAL_STATE.db,
            connector: {
              ...INITIAL_STATE.db.connector,
              kind: 'postgres',
            },
          },
        }}
        dispatch={vi.fn()}
        onCommit={() => { }}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: 'step4.delivery_expand_btn' }));

    expect(screen.getByTestId('database-delivery-truth')).toBeInTheDocument();
  });

  it('支援一鍵啟用或停用所有寫入欄位', () => {
    const dispatch = vi.fn();
    render(
      <Step4Database
        state={{
          ...INITIAL_STATE,
          points,
          mappings,
          db: {
            connector: {
              ...INITIAL_STATE.db.connector,
              kind: 'postgres',
            },
            targets,
          },
        }}
        dispatch={dispatch}
        onCommit={() => { }}
      />
    );

    dispatch.mockClear();
    fireEvent.click(screen.getByTestId('btn-disable-all-db-targets'));
    fireEvent.click(screen.getByTestId('btn-enable-all-db-targets'));

    expect(dispatch).toHaveBeenNthCalledWith(1, {
      type: 'setAllDbTargetsEnabled',
      enabled: false,
    });
    expect(dispatch).toHaveBeenNthCalledWith(2, {
      type: 'setAllDbTargetsEnabled',
      enabled: true,
    });
  });

  it('啟動成功後應將 Step 4 控制項切成唯讀，reset 後再恢復可編輯', async () => {
    const activateWorkspace = vi.fn().mockResolvedValue({
      workspace_id: 'workspace-1',
      results: [
        { device_id: 'd-1', status: 'success', message: 'activated' },
      ],
    });

    render(
      <Step4Database
        state={{
          ...INITIAL_STATE,
          points,
          mappings,
          db: {
            connector: {
              ...INITIAL_STATE.db.connector,
              kind: 'postgres',
              save_state: 'saved',
            },
            row_groups: [],
            targets: {
              'p-1': { ...targets['p-1'], enabled: true, save_state: 'saved' },
              'p-2': { ...targets['p-2'], enabled: true, save_state: 'saved' },
            },
          },
        }}
        dispatch={vi.fn()}
        activateWorkspace={activateWorkspace}
        onCommit={() => { }}
      />
    );

    const hostInput = screen.getByLabelText('step4.field_host');
    const createRowGroupButton = screen.getByTestId('step4-row-group-create');
    const disableAllButton = screen.getByTestId('btn-disable-all-db-targets');

    expect(hostInput).not.toBeDisabled();
    expect(createRowGroupButton).not.toBeDisabled();
    expect(disableAllButton).not.toBeDisabled();

    fireEvent.click(screen.getByText('step4.activate_btn'));

    await waitFor(() => {
      expect(screen.getByText('step4.go_to_dashboard_btn')).toBeInTheDocument();
    });

    expect(screen.getByLabelText('step4.field_host')).toBeDisabled();
    expect(screen.getByTestId('step4-row-group-create')).toBeDisabled();
    expect(screen.getByTestId('btn-disable-all-db-targets')).toBeDisabled();

    fireEvent.click(screen.getByText('step4.reset_activation_btn'));

    await waitFor(() => {
      expect(screen.queryByText('step4.go_to_dashboard_btn')).not.toBeInTheDocument();
    });

    expect(screen.getByLabelText('step4.field_host')).not.toBeDisabled();
    expect(screen.getByTestId('step4-row-group-create')).not.toBeDisabled();
    expect(screen.getByTestId('btn-disable-all-db-targets')).not.toBeDisabled();
  });
});
