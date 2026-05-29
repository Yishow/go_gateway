import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { dbReducer } from '../../../src/features/datalink/workbench-v2/state/dbReducer';
import type { WorkbenchV2State, CommitLog } from '../../../src/features/datalink/workbench-v2/state/types';
import type { WorkbenchV2Action } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

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

describe('Step 4 Database Commit Flow Integration', () => {
  const mockState: WorkbenchV2State = {
    view: 'flow',
    current: 4,
    completed: new Set([1, 2, 3]),
    sidebarCollapsed: false,
    showSummaryRail: true,
    devices: [
      { id: 'd-1', name: 'PLC-1', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null }
    ],
    rules: [
      { id: 'r-1', device_id: 'd-1', name: 'Holding Registers', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 't_', enabled: true, scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null }
    ],
    selectedRuleId: null,
    points: [
      { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_1', address: '40001', data_type: 'int16', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 }
    ],
    mappings: {
      'p-1': { point_id: 'p-1', tag_key: 'line1.t_1', display_name: 'T1', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true }
    },
    db: {
      connector: {
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
      },
      targets: {
        'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true }
      }
    },
    settings: {} as any,
    committed: false
  };

  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('點擊 Commit 後應進入 committing 狀態，經 10 次 timer 後完成，且 disabled form', () => {
    let state = { ...mockState };
    const dispatch = vi.fn((action: WorkbenchV2Action) => {
      state = dbReducer(state, action);
    });

    const { rerender } = render(
      <Step4Database state={state} dispatch={dispatch} />
    );

    // 點擊提交按鈕
    const commitBtn = screen.getByText('step4.submit_btn');
    fireEvent.click(commitBtn);

    expect(dispatch).toHaveBeenCalledWith({ type: 'startCommit' });

    // 重新用跑完 startCommit 的 state 渲染
    rerender(<Step4Database state={state} dispatch={dispatch} />);

    // 此時狀態為 committing，應顯示 pulse loader 與執行中
    expect(screen.getByTestId('pulse-loader')).toBeInTheDocument();
    expect(screen.getByText('step4.executing_next_command')).toBeInTheDocument();

    // 依序推進 10 次計時器
    for (let i = 0; i < 10; i++) {
      act(() => {
        vi.advanceTimersByTime(280);
      });
      rerender(<Step4Database state={state} dispatch={dispatch} />);
    }

    // 完成後應有 10 個 log
    expect(state.commit?.logs).toHaveLength(10);
    expect(state.committed).toBe(true);
    expect(state.commit?.status).toBe('success');

    // 再次重新渲染，此時應顯示 CommitSuccessCard
    rerender(<Step4Database state={state} dispatch={dispatch} />);
    expect(screen.getByText('step4.success_title')).toBeInTheDocument();
  });

  it('在 commit 進行中斷開並重 mount，應延續原進度繼續執行', () => {
    // 模擬已經執行了 3 筆 logs 的 committing 狀態
    const partialLogs: CommitLog[] = [
      { label: 'POST /devices × 1', detail: 'PLC-1 (modbus_tcp)', status: 'success' },
      { label: 'POST /devices/:id/activate × 1', detail: 'draft → active', status: 'success' },
      { label: 'POST /source-rules × 1', detail: 'Holding Registers', status: 'success' }
    ];

    let state: WorkbenchV2State = {
      ...mockState,
      commit: {
        status: 'committing',
        logs: partialLogs,
        started_at: new Date().toISOString()
      }
    };

    const dispatch = vi.fn((action: WorkbenchV2Action) => {
      state = dbReducer(state, action);
    });

    // 第一次渲染
    const { unmount } = render(
      <Step4Database state={state} dispatch={dispatch} />
    );

    // unmount 模擬元件被卸載
    unmount();

    // 重新 mount 渲染
    render(
      <Step4Database state={state} dispatch={dispatch} />
    );

    // 推進一次 timer
    act(() => {
      vi.advanceTimersByTime(280);
    });

    // 應該會增加第 4 筆 log，其 label 應為 'POST /points × 1'
    expect(dispatch).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'appendCommitLog',
        log: expect.objectContaining({
          label: 'POST /points × 1'
        })
      })
    );
  });

  it('點擊 emerald 完成卡內的「前往 Runtime Dashboard」按鈕，應觸發 onCommit 一次', () => {
    const onCommit = vi.fn();
    const successState: WorkbenchV2State = {
      ...mockState,
      committed: true,
      commit: {
        status: 'success',
        logs: []
      }
    };

    render(
      <Step4Database state={successState} dispatch={() => {}} onCommit={onCommit} />
    );

    const btn = screen.getByText('step4.go_to_dashboard_btn');
    fireEvent.click(btn);

    expect(onCommit).toHaveBeenCalledTimes(1);
  });
});
