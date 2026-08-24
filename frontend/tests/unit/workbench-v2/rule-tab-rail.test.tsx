import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { RuleTabRail } from '../../../src/features/datalink/workbench-v2/steps/step2/RuleTabRail';
import type { Rule, Device } from '../../../src/features/datalink/workbench-v2/state/types';

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

    expect(screen.getByText('Coils Rule')).toBeInTheDocument();

    fireEvent.click(screen.getByTestId('rule-add-btn'));
    expect(dispatch).toHaveBeenCalledWith(
      expect.objectContaining({ type: 'addRule', rule: expect.any(Object) })
    );
  });

  it('新增第二條規則時應給唯一的預設點位名稱前綴，避免與既有規則重名', () => {
    const dispatch = vi.fn();

    render(
      <RuleTabRail
        rules={mockRules}
        devices={mockDevices}
        selectedRuleId="rule-1"
        dispatch={dispatch}
      />
    );

    fireEvent.click(screen.getByTestId('rule-add-btn'));

    expect(dispatch).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'addRule',
        rule: expect.objectContaining({
          name: '規則 2',
          naming_prefix: 'BLOCK2_',
        }),
      }),
    );
  });

  it('新增 MC Protocol 設備規則時起始位址應預設為 D0，Modbus 為 40001', () => {
    const dispatch = vi.fn();
    const mcDevices: Device[] = [
      { id: 'dev-mc', name: 'MC PLC', description: '', protocol: 'mc_3e', config: {}, status: 'draft', test: null },
    ];
    const mcRules: Rule[] = [
      { ...mockRules[0], id: 'rule-mc', device_id: 'dev-mc' },
    ];

    render(
      <RuleTabRail
        rules={mcRules}
        devices={mcDevices}
        selectedRuleId="rule-mc"
        dispatch={dispatch}
      />
    );

    fireEvent.click(screen.getByTestId('rule-add-btn'));

    expect(dispatch).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'addRule',
        rule: expect.objectContaining({
          start_address: 'D0',
          device_id: 'dev-mc',
          name: '規則 2',
          count: 8,
        }),
      }),
    );
  });

  it('選取孤兒規則時新增規則應改用目前第一台設備及其協議預設位址', () => {
    const dispatch = vi.fn();
    const devices: Device[] = [
      { id: 'dev-current', name: 'Current MC PLC', description: '', protocol: 'mc_3e', config: {}, status: 'draft', test: null },
      ...mockDevices,
    ];
    const orphanRules: Rule[] = [
      { ...mockRules[0], device_id: 'deleted-device' },
    ];

    render(
      <RuleTabRail
        rules={orphanRules}
        devices={devices}
        selectedRuleId="rule-1"
        dispatch={dispatch}
      />
    );

    fireEvent.click(screen.getByTestId('rule-add-btn'));

    expect(dispatch).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'addRule',
        rule: expect.objectContaining({
          device_id: 'dev-current',
          start_address: 'D0',
        }),
      }),
    );
  });

  it('沒有設備時新增規則不得製造 phantom device', () => {
    const dispatch = vi.fn();

    render(
      <RuleTabRail
        rules={[]}
        devices={[]}
        selectedRuleId={null}
        dispatch={dispatch}
      />
    );

    fireEvent.click(screen.getByTestId('rule-add-btn'));

    expect(dispatch).not.toHaveBeenCalled();
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

    const toggleBtn = screen.getByTestId('rule-toggle-enabled-rule-1');
    fireEvent.click(toggleBtn);
    expect(dispatch).toHaveBeenCalledWith({ type: 'toggleRuleEnabled', ruleId: 'rule-1' });

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

  it('每個 rule tab 永遠顯示 rule color dot 與 enabled count，兩台設備時才顯示 owning device row', () => {
    const dispatch = vi.fn();
    const rules = [
      mockRules[0],
      { ...mockRules[0], id: 'rule-2', device_id: 'dev-2', name: 'Second Rule', enabled: false },
    ];
    const devices: Device[] = [
      ...mockDevices,
      { id: 'dev-2', name: 'PLC 2', description: '', protocol: 'mc_3e', config: {}, status: 'draft', test: null },
    ];

    render(
      <RuleTabRail
        rules={rules}
        devices={devices}
        selectedRuleId="rule-1"
        dispatch={dispatch}
      />
    );

    expect(screen.getByTestId('rule-color-dot-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('rule-color-dot-rule-2')).toBeInTheDocument();
    expect(screen.getByTestId('rule-enabled-count-rule-1')).toHaveTextContent('8/8');
    expect(screen.getByTestId('rule-enabled-count-rule-2')).toHaveTextContent('0/8');
    expect(screen.getByTestId('rule-device-row-rule-1')).toHaveTextContent('PLC 1');
    expect(screen.getByTestId('rule-device-row-rule-2')).toHaveTextContent('PLC 2');
  });
});
