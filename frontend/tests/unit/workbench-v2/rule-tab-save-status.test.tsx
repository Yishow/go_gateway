import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RuleTabRail } from '../../../src/features/datalink/workbench-v2/steps/step2/RuleTabRail';
import type { Rule } from '../../../src/features/datalink/workbench-v2/state/types';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const rule: Rule = { id: 'R', device_id: 'D', name: 'Range', start_address: '40009', count: 2, data_type: 'int16', naming_prefix: 'R', enabled: false,
  scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null, persisted: true };
function show(save_state?: Rule['save_state']) {
 const dispatch = vi.fn();
 render(<RuleTabRail rules={[{ ...rule, save_state }]} devices={[]} selectedRuleId="R" dispatch={dispatch} />);
 return dispatch;
}
describe('rule toggle save evidence', () => {
 it.each(['saving', 'save-error', 'draft-invalid'] as const)('shows %s instead of implying the disabled draft was persisted', state => {
  const dispatch = show(state);
  const status = screen.getByTestId('rule-save-status-R');
  expect(status).toHaveTextContent(`step3.saveStates.${state}`);
  expect(screen.getByTestId('rule-tab-R')).not.toHaveClass('opacity-60');
  expect(screen.getByTestId('rule-toggle-enabled-R')).toHaveAttribute('aria-describedby', status.id);
  fireEvent.click(screen.getByTestId('rule-toggle-enabled-R'));
  expect(dispatch).toHaveBeenCalledWith({ type: 'toggleRuleEnabled', ruleId: 'R' });
 });
 it('only presents saved disabled styling after explicit saved state', () => {
  show('saved');
  expect(screen.getByTestId('rule-save-status-R')).toHaveTextContent('step3.saveStates.saved');
  expect(screen.getByTestId('rule-tab-R')).toHaveClass('opacity-60');
 });
 it.each([undefined, 'idle'] as const)('does not infer persisted toggle state from missing/idle save evidence', state => {
  show(state);
  expect(screen.getByTestId('rule-save-status-R')).toHaveTextContent('status.draft');
  expect(screen.getByTestId('rule-tab-R')).not.toHaveClass('opacity-60');
 });
});
