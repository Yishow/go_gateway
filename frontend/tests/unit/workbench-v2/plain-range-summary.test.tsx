import { fireEvent, render, screen } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RangeSummary } from '@/features/datalink/workbench-v2/steps/step2/RangeSummary';
import { RuleEditor } from '@/features/datalink/workbench-v2/steps/step2/RuleEditor';
import { savedGroupState } from '../../fixtures/writeGroupState';
it('distinguishes four float32 values from eight words and bit addressing in both locales', async () => {
  for (const language of ['zh-TW', 'en']) {
    await i18n.changeLanguage(language);
    const view = render(<RangeSummary startAddress="40001" count={4} dataType="float32" />);
    expect(screen.getByTestId('range-summary')).toHaveTextContent(/4.*float32.*8 words/);
    view.rerender(<RangeSummary startAddress="00001" count={4} dataType="bool" />);
    expect(screen.getByTestId('range-summary')).toHaveTextContent(/4.*bool.*4 bits/);
    view.unmount();
  }
});
it('suggests a descriptive editable name only on explicit request and preserves the required type selector', async () => {
  await i18n.changeLanguage('zh-TW');
  const dispatch = vi.fn(), state = savedGroupState();
  render(<RuleEditor rule={{ ...state.rules[0], name: 'Custom', count: 4, data_type: 'float32', start_address: '40001' }} devices={state.devices} globalShareEnabled={false} dispatch={dispatch} />);
  expect(screen.getByTestId('rule-name-input')).toHaveValue('Custom');
  expect(screen.getByTestId('rule-datatype-select')).toBeVisible();
  expect(dispatch).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('rule-suggest-name'));
  expect(dispatch).toHaveBeenCalledWith({ type: 'updateRule', ruleId: 'rule-1', patch: { name: 'PLC A 40001 · 4 float32' } });
});
