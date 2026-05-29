import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ShareSection } from '../../../src/features/datalink/workbench-v2/steps/step2/ShareSection';

describe('ShareSection Component', () => {
  it('當全域未啟用時，應渲染「全域未啟用」標籤', () => {
    const dispatch = vi.fn();
    render(
      <ShareSection
        ruleId="rule-1"
        shareEnabled={false}
        shareStartRegister={null}
        shareStride={null}
        globalShareEnabled={false}
        dispatch={dispatch}
      />
    );

    const chip = screen.getByTestId('global-disabled-chip');
    expect(chip).toBeInTheDocument();
    expect(chip.textContent).toBe('全域未啟用');
  });

  it('當全域啟用時，不應渲染「全域未啟用」標籤', () => {
    const dispatch = vi.fn();
    render(
      <ShareSection
        ruleId="rule-1"
        shareEnabled={false}
        shareStartRegister={null}
        shareStride={null}
        globalShareEnabled={true}
        dispatch={dispatch}
      />
    );

    const chip = screen.queryByTestId('global-disabled-chip');
    expect(chip).toBeNull();
  });

  it('點擊 Toggle 應 dispatch toggleRuleShareEnabled', () => {
    const dispatch = vi.fn();
    render(
      <ShareSection
        ruleId="rule-1"
        shareEnabled={false}
        shareStartRegister={null}
        shareStride={null}
        globalShareEnabled={true}
        dispatch={dispatch}
      />
    );

    const toggle = screen.getByRole('switch');
    fireEvent.click(toggle);

    expect(dispatch).toHaveBeenCalledWith({
      type: 'toggleRuleShareEnabled',
      ruleId: 'rule-1',
    });
  });

  it('手動設定起點暫存器時，應 dispatch updateRuleShareStart；若清空則傳遞 null', () => {
    const dispatch = vi.fn();
    render(
      <ShareSection
        ruleId="rule-1"
        shareEnabled={true}
        shareStartRegister={40005}
        shareStride={null}
        globalShareEnabled={true}
        dispatch={dispatch}
      />
    );

    const startInput = screen.getByTestId('share-start-input') as HTMLInputElement;
    expect(startInput.value).toBe('40005');

    // 修改值
    fireEvent.change(startInput, { target: { value: '40010' } });
    expect(dispatch).toHaveBeenLastCalledWith({
      type: 'updateRuleShareStart',
      ruleId: 'rule-1',
      shareStart: 40010,
    });

    // 清空輸入
    fireEvent.change(startInput, { target: { value: '' } });
    expect(dispatch).toHaveBeenLastCalledWith({
      type: 'updateRuleShareStart',
      ruleId: 'rule-1',
      shareStart: null,
    });
  });

  it('點擊「改回自動分配」連結，應將起點重置為 null', () => {
    const dispatch = vi.fn();
    render(
      <ShareSection
        ruleId="rule-1"
        shareEnabled={true}
        shareStartRegister={40005}
        shareStride={null}
        globalShareEnabled={true}
        dispatch={dispatch}
      />
    );

    const resetLink = screen.getByTestId('share-reset-link');
    expect(resetLink).toBeInTheDocument();

    fireEvent.click(resetLink);
    expect(dispatch).toHaveBeenCalledWith({
      type: 'updateRuleShareStart',
      ruleId: 'rule-1',
      shareStart: null,
    });
  });
});
