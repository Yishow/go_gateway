import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import * as React from 'react';
import { useState } from 'react';
import { PointGrid } from '../../../src/features/datalink/workbench-v2/steps/step2/PointGrid';
import { PointGridToolbar } from '../../../src/features/datalink/workbench-v2/steps/step2/PointGridToolbar';
import type { Point, ShareLayout } from '../../../src/features/datalink/workbench-v2/state/types';

// Mock 點位資料
const mockPoints: Point[] = [
  {
    id: 'r1-p-0',
    device_id: 'dev-1',
    rule_id: 'rule-1',
    rule_name: 'Rule 1',
    name: 'TAG_0',
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
    id: 'r1-p-1',
    device_id: 'dev-1',
    rule_id: 'rule-1',
    rule_name: 'Rule 1',
    name: 'TAG_1',
    address: '40002',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
  {
    id: 'r1-p-2',
    device_id: 'dev-1',
    rule_id: 'rule-1',
    rule_name: 'Rule 1',
    name: 'TAG_2',
    address: '40003',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: false,
    skipped: true,
    _rule_scale: 1,
    _rule_offset: 0,
  },
];

describe('PointGridToolbar 元件', () => {
  it('正常渲染，且當選取為空時特定按鈕為 disabled', () => {
    const onAction = vi.fn();
    const selection = new Set<string>();

    render(
      <PointGridToolbar
        gridSelection={selection}
        totalCount={8}
        stride={2}
        onAction={onAction}
      />
    );

    // 驗證統計
    expect(screen.getByTestId('grid-summary').textContent).toBe('total 8 · stride 2');

    // 選取為空時，略過選取與啟用選取應為 disabled
    expect(screen.getByTestId('btn-skip-selected')).toBeDisabled();
    expect(screen.getByTestId('btn-enable-selected')).toBeDisabled();
    expect(screen.queryByTestId('btn-clear-selection')).toBeNull();

    // 點擊「全部啟用」
    fireEvent.click(screen.getByTestId('btn-all-enabled'));
    expect(onAction).toHaveBeenCalledWith('all_enabled');
  });

  it('當選取不為空時，按鈕應啟用，且出現清除選取按鈕', () => {
    const onAction = vi.fn();
    const selection = new Set<string>(['r1-p-0']);

    render(
      <PointGridToolbar
        gridSelection={selection}
        totalCount={8}
        stride={2}
        onAction={onAction}
      />
    );

    expect(screen.getByTestId('btn-skip-selected')).toBeEnabled();
    expect(screen.getByTestId('btn-enable-selected')).toBeEnabled();

    const clearBtn = screen.getByTestId('btn-clear-selection');
    expect(clearBtn).toBeInTheDocument();
    expect(clearBtn.textContent).toContain('清除選取 (1)');

    fireEvent.click(clearBtn);
    expect(onAction).toHaveBeenCalledWith('clear_selection');
  });
});

// point grid 包裝器，方便測試 React state 連動
const PointGridTestWrapper: React.FC<{
  ruleId: string;
  points: Point[];
  conflictAddrs?: Set<string>;
  shareLayout?: ShareLayout | null;
  onToggle?: (addr: string) => void;
  onBatchToggle?: (addrs: string[], shouldSkip: boolean) => void;
}> = ({
  ruleId,
  points,
  conflictAddrs = new Set(),
  shareLayout = null,
  onToggle = () => {},
  onBatchToggle = () => {},
}) => {
  const [selection, setSelection] = useState<Set<string>>(new Set());
  return (
    <PointGrid
      ruleId={ruleId}
      points={points}
      conflictAddrs={conflictAddrs}
      shareLayout={shareLayout}
      gridSelection={selection}
      setGridSelection={setSelection}
      onToggleSkipAddress={onToggle}
      onBatchToggleSkipAddresses={onBatchToggle}
    />
  );
};

describe('PointGrid 元件', () => {
  it('普通單擊 cell 應呼叫 onToggleSkipAddress', () => {
    const onToggle = vi.fn();
    render(<PointGridTestWrapper ruleId="rule-1" points={mockPoints} onToggle={onToggle} />);

    const cell0 = screen.getByTestId('point-cell-r1-p-0');
    fireEvent.click(cell0);

    expect(onToggle).toHaveBeenCalledWith('40001');
  });

  it('Ctrl+Click 應 toggle 選取狀態', () => {
    render(<PointGridTestWrapper ruleId="rule-1" points={mockPoints} />);

    const cell0 = screen.getByTestId('point-cell-r1-p-0');
    expect(cell0.getAttribute('data-selected')).toBe('false');

    // Ctrl + click
    fireEvent.click(cell0, { ctrlKey: true });
    expect(cell0.getAttribute('data-selected')).toBe('true');

    // 再次 Ctrl + click
    fireEvent.click(cell0, { ctrlKey: true });
    expect(cell0.getAttribute('data-selected')).toBe('false');
  });

  it('Shift+Click 連續選取應呼叫 onBatchToggleSkipAddresses', () => {
    const onBatchToggle = vi.fn();
    render(<PointGridTestWrapper ruleId="rule-1" points={mockPoints} onBatchToggle={onBatchToggle} />);

    const cell0 = screen.getByTestId('point-cell-r1-p-0');
    const cell2 = screen.getByTestId('point-cell-r1-p-2');

    // 先普通單擊第一個 cell0 建立 lastClickedIdx 基底
    fireEvent.click(cell0);

    // 再 Shift+Click cell2
    fireEvent.click(cell2, { shiftKey: true });

    // mockPoints[2] 是 skipped = true，所以對其進行反轉應期望將 0..2 的 address 批次設為啟用 (shouldSkip = false)
    expect(onBatchToggle).toHaveBeenCalledWith(['40001', '40002', '40003'], false);
  });

  it('轉發與衝突狀態渲染', () => {
    const conflicts = new Set(['40002']);
    const shareLayout: ShareLayout = {
      start: 30001,
      stride: 2,
      end: 30005,
      auto: false,
    };

    render(
      <PointGridTestWrapper
        ruleId="rule-1"
        points={mockPoints}
        conflictAddrs={conflicts}
        shareLayout={shareLayout}
      />
    );

    // cell-r1-p-1 地址為 40002，應顯示衝突
    const cell1 = screen.getByTestId('point-cell-r1-p-1');
    expect(cell1.getAttribute('data-conflict')).toBe('true');

    // 轉發位址標籤。啟用點位 r1-p-0 是第 0 個啟用點位，其轉發位址應為 30001
    const shareLabel0 = screen.getByTestId('share-addr-r1-p-0');
    expect(shareLabel0.textContent).toBe('→30001');

    // 啟用點位 r1-p-1 是第 1 個啟用點位，其轉發位址應為 30001 + 1*2 = 30003
    const shareLabel1 = screen.getByTestId('share-addr-r1-p-1');
    expect(shareLabel1.textContent).toBe('→30003');

    // skipped 點位 r1-p-2 不應有轉發標籤
    const shareLabel2 = screen.queryByTestId('share-addr-r1-p-2');
    expect(shareLabel2).toBeNull();
  });
});
