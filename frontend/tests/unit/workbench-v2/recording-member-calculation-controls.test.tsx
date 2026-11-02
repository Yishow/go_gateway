import { render, screen, fireEvent } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { RecordingMemberCalculation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingMemberCalculationPanel';
import type { MemberCalculationContext } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingMemberCalculation';
import type { EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const member: EditorMember = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', target_column: '', entity_key: '', required: true, calculation: 'latest', cadence: { interval: 15, unit: 'second', alignment: 'utc_epoch' } };
const context: MemberCalculationContext = { workspaceId: 'W', recordId: 'record-1', definitions: [], definitionsAvailable: true, layout: 'managed_wide_profile_v1', dialect: 'sqlite', capabilities: { version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false, reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true }, schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1'], unsupported_combinations: [] } };
it('keeps cadence per selected member and shows unavailable derived choices without global policy mutation', () => {
 const change = vi.fn(); render(<RecordingMemberCalculation member={member} context={context} disabled={false} onChange={change} />);
 expect(screen.getByRole('option', { name: 'calculation.mean' })).toBeDisabled();
 expect(screen.getByRole('option', { name: 'calculation.delta' })).toBeDisabled();
 fireEvent.change(screen.getByLabelText('calculation.interval'), { target: { value: '60' } });
 expect(change).toHaveBeenLastCalledWith({ ...member, cadence_input: '60', cadence: { interval: 60, unit: 'second', alignment: 'utc_epoch' } });
 expect(member.cadence?.interval).toBe(15);
 expect(screen.getByText('calculation.dueUnion')).toBeInTheDocument();
});
it('retains an invalid cadence as invalid instead of clamping or treating it as the last saved value', () => {
 const change = vi.fn(); render(<RecordingMemberCalculation member={member} context={context} disabled={false} onChange={change} />);
 fireEvent.change(screen.getByLabelText('calculation.interval'), { target: { value: '' } });
 expect(change.mock.calls[0][0].cadence.interval).toBeNaN();
});
it('keeps read-only controls inert even if backend advertises preparation support', () => {
 const change = vi.fn(); render(<RecordingMemberCalculation member={member} context={context} disabled onChange={change} />);
 for (const control of screen.getAllByRole('combobox')) expect(control).toBeDisabled();
 expect(screen.getByRole('textbox')).toBeDisabled(); expect(change).not.toHaveBeenCalled();
});

it('keeps an invalid interval lexeme in the member buffer so picker cancellation or undo cannot silently restore a valid value', () => {
 const change = vi.fn(); const value = { ...member, cadence_input: '1e', cadence: { ...member.cadence!, interval: NaN } };
 render(<RecordingMemberCalculation member={value} context={context} disabled={false} onChange={change} />);
 expect(screen.getByLabelText('calculation.interval')).toHaveValue('1e'); expect(change).not.toHaveBeenCalled();
});

it('offers explicit record-only gauge confirmation only for a fully known unbound source', () => {
 const change = vi.fn(); const source = { ...member, source_value_type: 'float32' as const };
 const view = render(<RecordingMemberCalculation member={source} context={context} disabled={false} onChange={change} />);
 expect(screen.getByRole('option', { name: 'calculation.mean' })).toBeDisabled(); expect(change).not.toHaveBeenCalled();
 fireEvent.click(screen.getByRole('button', { name: 'calculation.confirmGauge' }));
 const confirmed = change.mock.calls[0][0]; expect(confirmed.value_role).toBe('gauge'); expect(confirmed.measurement_id).toBeUndefined();
 view.rerender(<RecordingMemberCalculation member={confirmed} context={context} disabled={false} onChange={change} />);
 expect(screen.getByRole('option', { name: 'calculation.mean' })).not.toBeDisabled();
});

it('renders compact method and cadence in two native cells while preserving explicit consent and accessible controls', () => {
  const change = vi.fn(), source = { ...member, source_value_type: 'float32' as const };
  const view = render(<table><tbody><tr><RecordingMemberCalculation compact tableCells member={source} context={context} disabled={false} onChange={change} /></tr></tbody></table>);
  const cells = screen.getAllByRole('cell'); expect(cells).toHaveLength(2);
  expect(cells[0]).toContainElement(screen.getByLabelText('calculation.title'));
  expect(cells[1]).toContainElement(screen.getByLabelText('calculation.interval'));
  expect(cells[1]).toContainElement(screen.getByLabelText('calculation.unit'));
  expect(view.container.querySelector('tr > fieldset')).toBeNull();
  const consent = screen.getByText('compact.gaugeConsent'); expect(consent).toBeVisible(); expect(consent.closest('details')).toBeNull();
  expect(screen.getByRole('option', { name: 'calculation.mean' })).toBeDisabled(); expect(change).not.toHaveBeenCalled();
  expect(screen.getByText(/calculation.counter_required/).closest('details')).toBe(screen.getByTestId('member-calculation-help-P'));
  fireEvent.click(screen.getByRole('button', { name: 'compact.confirmGauge' }));
  expect(change.mock.calls[0][0].value_role).toBe('gauge'); expect(change.mock.calls[0][0].calculation).toBe('latest');
});
it('keeps an already selected invalid calculation visible outside optional help and never repairs it on render', () => {
  const change = vi.fn(); render(<table><tbody><tr><RecordingMemberCalculation compact tableCells member={{ ...member, calculation: 'delta' }} context={context} disabled={false} onChange={change} /></tr></tbody></table>);
  const problem = screen.getByRole('alert'); expect(problem).toHaveTextContent('calculation.counter_required'); expect(problem.closest('details')).toBeNull();
  expect(screen.getByLabelText('calculation.title')).toHaveValue('delta'); expect(change).not.toHaveBeenCalled();
});
