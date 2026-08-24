import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { RangeSummary } from '../../../src/features/datalink/workbench-v2/steps/step2/RangeSummary';
import { ScaleSection } from '../../../src/features/datalink/workbench-v2/steps/step2/ScaleSection';
import { RuleEditor } from '../../../src/features/datalink/workbench-v2/steps/step2/RuleEditor';
import { MergedPointTable } from '../../../src/features/datalink/workbench-v2/steps/step2/MergedPointTable';
import { Step2Rule } from '../../../src/features/datalink/workbench-v2/steps/step2/Step2Rule';
import type { Rule, Device, Point } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

vi.mock('react-i18next', () => ({ useTranslation: () => ({
  t: (key: string, value?: string | { defaultValue?: string; index?: number }) => typeof value === 'string' ? value : (value?.defaultValue ?? key).replace('{{index}}', String(value?.index ?? '')),
}) }));
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

  it('MC/FATEK 應保留暫存器前綴並顯示協議範圍與 badge', () => {
    const { rerender } = render(
      <RangeSummary startAddress="D0" count={8} dataType="int16" protocol="mc_3e" />
    );
    expect(screen.getByTestId('range-summary')).toHaveTextContent('D0 ~ D7');
    expect(screen.getByTestId('function-code-badge')).toHaveTextContent('MC 3E D (Word)');

    rerender(<RangeSummary startAddress="R0" count={8} dataType="int16" protocol="fatek_fbs" />);
    expect(screen.getByTestId('range-summary')).toHaveTextContent('R0 ~ R7');
    expect(screen.getByTestId('function-code-badge')).toHaveTextContent('FATEK R (Word)');
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

  it('無效位址應顯示 localized error', () => {
    const dispatch = vi.fn();
    const mcRule = { ...mockRules[0], start_address: 'Z999' };
    const mcDevice: Device = { id: 'dev-mc', name: 'MC PLC', description: '', protocol: 'mc_3e', config: {}, status: 'draft', test: null };
    render(<RuleEditor rule={mcRule} devices={[mcDevice]} globalShareEnabled={true} dispatch={dispatch} />);
    expect(screen.getByTestId('rule-start-input')).toHaveAttribute('aria-invalid', 'true');
    expect(screen.getByTestId('rule-start-input')).toHaveAttribute('placeholder', 'D0');
    expect(screen.getByTestId('rule-start-error')).toBeInTheDocument();
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

describe('Step2Rule 元件', () => {
  it('當目前沒有任何規則時仍保留新增入口與 Step 2 骨架', () => {
    const dispatch = vi.fn();

    render(
      <Step2Rule
        state={{
          ...INITIAL_STATE,
          devices: mockDevices,
          rules: [],
          selectedRuleId: 'missing-rule',
        }}
        dispatch={dispatch}
        onContinue={vi.fn()}
      />
    );

    expect(screen.getByText('暫無配置規則，請先新增接入規則。')).toBeInTheDocument();
    expect(screen.getByTestId('rule-tab-rail')).toBeInTheDocument();
    expect(screen.getByTestId('rule-add-btn')).toBeInTheDocument();
    expect(screen.queryByTestId('rule-editor')).not.toBeInTheDocument();
    expect(screen.queryByTestId('point-grid')).not.toBeInTheDocument();
    expect(screen.getByTestId('merged-point-table-container')).toBeInTheDocument();
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

  it('無效 MC 位址不應衍生點位，且繼續到 Step 3 應 disabled', () => {
    const invalidRule = { ...mockRules[0], device_id: 'dev-mc', start_address: 'Z999' };
    const state = {
      ...INITIAL_STATE,
      devices: [{ id: 'dev-mc', name: 'MC PLC', description: '', protocol: 'mc_3e' as const, config: {}, status: 'draft' as const, test: null }],
      rules: [invalidRule],
      selectedRuleId: invalidRule.id,
    };

    render(<Step2Rule state={state} dispatch={vi.fn()} onContinue={vi.fn()} />);

    expect(screen.getByTestId('rule-start-input')).toHaveAttribute('aria-invalid', 'true');
    expect(screen.getByTestId('continue-step3-btn')).toBeDisabled();
  });

  it('任一啟用規則位址無效時，即使其他規則有效也應阻擋並定位該規則', () => {
    const validRule = { ...mockRules[0], id: 'rule-valid', name: 'Valid Rule', start_address: '40001' };
    const invalidRule = { ...mockRules[0], id: 'rule-invalid', name: 'Invalid MC Rule', device_id: 'dev-mc', start_address: 'Z999' };
    const state = {
      ...INITIAL_STATE,
      devices: [
        mockDevices[0],
        { id: 'dev-mc', name: 'MC PLC', description: '', protocol: 'mc_3e' as const, config: {}, status: 'draft' as const, test: null },
      ],
      rules: [validRule, invalidRule],
      selectedRuleId: validRule.id,
    };

    render(<Step2Rule state={state} dispatch={vi.fn()} onContinue={vi.fn()} />);

    expect(screen.getByTestId('continue-step3-btn')).toBeDisabled();
    expect(screen.getByTestId('invalid-rule-warning-rule-invalid')).toHaveTextContent('Invalid MC Rule');
  });
});
