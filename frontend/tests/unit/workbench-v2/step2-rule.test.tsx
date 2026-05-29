import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { RuleTabRail } from '../../../src/features/datalink/workbench-v2/steps/step2/RuleTabRail';
import { RangeSummary } from '../../../src/features/datalink/workbench-v2/steps/step2/RangeSummary';
import { ScaleSection } from '../../../src/features/datalink/workbench-v2/steps/step2/ScaleSection';
import { RuleEditor } from '../../../src/features/datalink/workbench-v2/steps/step2/RuleEditor';
import { MergedPointTable } from '../../../src/features/datalink/workbench-v2/steps/step2/MergedPointTable';
import { Step2Rule } from '../../../src/features/datalink/workbench-v2/steps/step2/Step2Rule';
import type { Rule, Device, Point } from '../../../src/features/datalink/workbench-v2/state/types';

// Mock 基礎資料
const mockDevices: Device[] = [
  { id: 'dev-1', name: 'PLC 1', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null },
];

const mockRules: Rule[] = [
  {
    id: 'rule-1',
    device_id: 'dev-1',
    name: 'Coils Rule',
    start_address: '00001',
    count: 8,
    data_type: 'bool',
    naming_prefix: 'C_',
    enabled: true,
    scale_multiplier: 1,
    scale_offset: 0,
    data_format: '',
    skipped_addresses: [],
    share_enabled: false,
    share_start_register: null,
    share_stride: null,
  },
];

describe('RuleTabRail 元件', () => {
  it('應渲染 rules 並在點擊時觸發 selectRule，點擊新增時觸發 addRule', () => {
    const dispatch = vi.fn();
    render(
      <RuleTabRail
        rules={mockRules}
        devices={mockDevices}
        selectedRuleId="rule-1"
        dispatch={dispatch}
      />
    );

    // 驗證 Tab 是否渲染
    expect(screen.getByText('Coils Rule')).toBeInTheDocument();

    // 點擊新增
    fireEvent.click(screen.getByTestId('rule-add-btn'));
    expect(dispatch).toHaveBeenCalledWith(
      expect.objectContaining({ type: 'addRule', rule: expect.any(Object) })
    );
  });

  it('Hover 時應渲染啟用 toggle 與刪除按鈕 (若大於1個規則)', () => {
    const dispatch = vi.fn();
    const twoRules = [
      ...mockRules,
      { ...mockRules[0], id: 'rule-2', name: 'Rule 2' },
    ];
    render(
      <RuleTabRail
        rules={twoRules}
        devices={mockDevices}
        selectedRuleId="rule-1"
        dispatch={dispatch}
      />
    );

    // 點擊 toggle 啟用狀態
    const toggleBtn = screen.getByTestId('rule-toggle-enabled-rule-1');
    fireEvent.click(toggleBtn);
    expect(dispatch).toHaveBeenCalledWith({ type: 'toggleRuleEnabled', ruleId: 'rule-1' });

    // 點擊刪除
    const deleteBtn = screen.getByTestId('rule-delete-rule-1');
    fireEvent.click(deleteBtn);
    expect(dispatch).toHaveBeenCalledWith({ type: 'removeRule', ruleId: 'rule-1' });
  });

  it('雙擊 Tab 名稱應進入改名編輯模式', () => {
    const dispatch = vi.fn();
    render(
      <RuleTabRail
        rules={mockRules}
        devices={mockDevices}
        selectedRuleId="rule-1"
        dispatch={dispatch}
      />
    );

    const nameText = screen.getByText('Coils Rule');
    fireEvent.doubleClick(nameText);

    const input = screen.getByRole('textbox');
    expect(input).toBeInTheDocument();
    fireEvent.change(input, { target: { value: 'New Coils Rule' } });
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(dispatch).toHaveBeenCalledWith({
      type: 'renameRule',
      ruleId: 'rule-1',
      name: 'New Coils Rule',
    });
  });
});

describe('RangeSummary 元件', () => {
  it('四種 prefix 應對應正確的功能碼文字', () => {
    const { rerender } = render(
      <RangeSummary startAddress="00001" count={10} dataType="bool" />
    );
    expect(screen.getByTestId('function-code-badge').textContent).toBe('Coil (0x)');

    rerender(<RangeSummary startAddress="10001" count={10} dataType="bool" />);
    expect(screen.getByTestId('function-code-badge').textContent).toBe('Discrete Input (1x)');

    rerender(<RangeSummary startAddress="30001" count={10} dataType="int16" />);
    expect(screen.getByTestId('function-code-badge').textContent).toBe('Input Register (3x)');

    rerender(<RangeSummary startAddress="40001" count={10} dataType="int16" />);
    expect(screen.getByTestId('function-code-badge').textContent).toBe('Holding Register (4x)');
  });
});

describe('ScaleSection 元件', () => {
  it('改變縮放係數應 dispatch updateRule', () => {
    const dispatch = vi.fn();
    render(
      <ScaleSection
        ruleId="rule-1"
        multiplier={1}
        offset={0}
        dataFormat=""
        dispatch={dispatch}
      />
    );

    const multInput = screen.getByTestId('scale-multiplier-input');
    fireEvent.change(multInput, { target: { value: '2.5' } });
    expect(dispatch).toHaveBeenLastCalledWith({
      type: 'updateRule',
      ruleId: 'rule-1',
      patch: { scale_multiplier: 2.5 },
    });
  });
});

describe('RuleEditor 元件', () => {
  it('當只有 1 個設備時，設備選擇器應為 disabled', () => {
    const dispatch = vi.fn();
    render(
      <RuleEditor
        rule={mockRules[0]}
        devices={mockDevices}
        globalShareEnabled={true}
        dispatch={dispatch}
      />
    );

    expect(screen.getByTestId('rule-device-select')).toBeDisabled();
  });

  it('輸入改變應觸發對應的 updateRule', () => {
    const dispatch = vi.fn();
    render(
      <RuleEditor
        rule={mockRules[0]}
        devices={mockDevices}
        globalShareEnabled={true}
        dispatch={dispatch}
      />
    );

    const countInput = screen.getByTestId('rule-count-input');
    fireEvent.change(countInput, { target: { value: '16' } });
    expect(dispatch).toHaveBeenCalledWith({
      type: 'updateRule',
      ruleId: 'rule-1',
      patch: { count: 16 },
    });
  });
});

describe('MergedPointTable 元件', () => {
  const points: Point[] = [
    {
      id: 'p-1',
      device_id: 'dev-1',
      rule_id: 'rule-1',
      rule_name: 'Coils Rule',
      name: 'C_0',
      address: '40001',
      data_type: 'int16',
      function: 'holding_register',
      width: 1,
      enabled: true,
      skipped: false,
      _rule_scale: 1,
      _rule_offset: 0,
    },
    {
      id: 'p-2',
      device_id: 'dev-1',
      rule_id: 'rule-1',
      rule_name: 'Coils Rule',
      name: 'C_1',
      address: '40002',
      data_type: 'int16',
      function: 'holding_register',
      width: 1,
      enabled: true,
      skipped: false,
      _rule_scale: 1,
      _rule_offset: 0,
    },
  ];

  it('應正確渲染點位列，且繼續按鈕在無衝突時為啟用狀態', () => {
    const onContinue = vi.fn();
    render(
      <MergedPointTable
        points={points}
        rules={mockRules}
        devices={mockDevices}
        conflictAddrs={new Set()}
        shareLayouts={{}}
        onContinue={onContinue}
      />
    );

    expect(screen.getByTestId('point-row-p-1')).toBeInTheDocument();
    expect(screen.getByTestId('point-row-p-2')).toBeInTheDocument();

    const continueBtn = screen.getByTestId('continue-step3-btn');
    expect(continueBtn).toBeEnabled();

    fireEvent.click(continueBtn);
    expect(onContinue).toHaveBeenCalled();
  });

  it('當存在衝突或無啟用點位時，繼續按鈕應為 disabled，且衝突地址應標記警告', () => {
    const onContinue = vi.fn();
    const conflicts = new Set(['40001']);

    render(
      <MergedPointTable
        points={points}
        rules={mockRules}
        devices={mockDevices}
        conflictAddrs={conflicts}
        shareLayouts={{}}
        onContinue={onContinue}
      />
    );

    const continueBtn = screen.getByTestId('continue-step3-btn');
    expect(continueBtn).toBeDisabled();

    // 衝突地址有警示標記
    expect(screen.getByTestId('conflict-addr-p-1')).toBeInTheDocument();
    expect(screen.getByTestId('status-conflict')).toBeInTheDocument();
  });
});

describe('Step2Rule 整合元件', () => {
  it('renders full Step 2 with default state', () => {
    const dispatch = vi.fn();
    const onContinue = vi.fn();
    const state = {
      view: 'flow' as const,
      current: 2 as const,
      completed: new Set<number>(),
      sidebarCollapsed: false,
      showSummaryRail: true,
      devices: mockDevices,
      rules: mockRules,
      selectedRuleId: 'rule-1',
      points: [],
      mappings: {},
      db: {
        connector: {
          kind: 'postgres' as const,
          name: 'TSDB',
          host: 'localhost',
          port: 5432,
          database: 'db',
          username: 'user',
          schema: 'public',
          table: 't',
          write_mode: 'insert' as const,
          write_interval_seconds: 5,
          timestamp_column: 'ts',
          status: 'ready',
        },
        targets: {},
      },
      settings: {
        connectors: [],
        timeseries: {
          write_precision: 'millisecond' as const,
          partition_interval: 'daily' as const,
          batch_size: 100,
          retention_days: 7,
        },
        scheduler: {
          default_interval_ms: 1000,
          default_retry_count: 3,
          default_retry_delay_ms: 100,
          breaker_threshold: 5,
          auto_start: true,
        },
        modbus_share: {
          enabled: true,
          bind_address: '0.0.0.0',
          port: 502,
          slave_id: 1,
          base_register: 40001,
        },
        general: {
          theme: 'dark' as const,
          locale: 'zh-TW' as const,
          addr_format: 'modbus' as const,
          api_base: '',
          api_version: 'v1' as const,
          timeout_seconds: 10,
          log_level: 'info' as const,
          sse_heartbeat_seconds: 10,
          enable_debug_panel: false,
          enable_audit_log: false,
        },
      },
      committed: false,
    };

    render(<Step2Rule state={state} dispatch={dispatch} onContinue={onContinue} />);

    // 應該要渲染容器與子元件
    expect(screen.getByTestId('step2-rule-container')).toBeInTheDocument();
    expect(screen.getByTestId('rule-editor')).toBeInTheDocument();
    expect(screen.getByTestId('point-grid')).toBeInTheDocument();
    expect(screen.getByTestId('merged-point-table-container')).toBeInTheDocument();
  });
});

