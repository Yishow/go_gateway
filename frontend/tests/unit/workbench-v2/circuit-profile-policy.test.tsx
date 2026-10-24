import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CounterPolicyEditor } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/CounterPolicyEditor';
import { studioV2WorkspaceMeasurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
import type { MeasurementDefinition } from '@/types/measurement';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({
  studioV2WorkspaceMeasurementsAPI: { list: vi.fn(), update: vi.fn(), create: vi.fn() },
}));

const api = vi.mocked(studioV2WorkspaceMeasurementsAPI);
const definition = (over: Partial<MeasurementDefinition> = {}): MeasurementDefinition => ({
  id: 'meas-1', workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', equipment_id: 'eq-A',
  definition_revision: 'def-5', source_binding_revision: 'bind-5', series_epoch: 'epoch-3', name: 'Energy', quantity: 'energy', unit: 'kWh',
  semantic_kind: 'counter', numeric_encoding: 'uint64',
  counter_policy: { allow_negative: false, max_roll_over: 123, reset_threshold: 7, rollover_modulus: '65536', max_rate_per_second: '0.5' },
  created_at: '2026-10-01T00:00:00Z', updated_at: '2026-10-02T00:00:00Z', ...over,
});
const second = (): MeasurementDefinition => definition({ id: 'meas-2', point_id: 'pt-1', tag_id: 'tag-1', name: 'Energy 2', counter_policy: { allow_negative: false } });
// A saved cumulative member is bound to its measurement definition by the persisted ID; the editor works on exactly that definition.
const MEMBERS = [{ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', measurement_id: 'meas-1' }, { device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1', measurement_id: 'meas-2' }];
const typed = (code: string) => Object.assign(new Error('PUT failed: INSERT INTO measurement_definitions postgres://u:secret@db/x'), { response: { status: 422, data: { error: { code } } } });

function mount(readonly = false, members = MEMBERS) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(<QueryClientProvider client={client}><CounterPolicyEditor workspaceId="ws-1" members={members} labels={{ 'tag-0': 'Energy A', 'tag-1': 'Energy B' }} readonly={readonly} /></QueryClientProvider>);
}
const field = (name: string, tag = 'tag-0') => screen.getByTestId(`counter-policy-${name}-${tag}`);

beforeEach(() => {
  vi.clearAllMocks();
  api.list.mockResolvedValue([definition(), second()]);
});

describe('Counter policy form edits the saved definition without guessing', () => {
  it('shows the saved exact values and revisions truthfully and writes nothing on open', async () => {
    mount();
    expect(await screen.findByTestId('counter-policy-source-tag-0')).toHaveTextContent('Energy A');
    expect(field('modulus')).toHaveValue('65536');
    expect(field('rate')).toHaveValue('0.5');
    expect(field('signed')).not.toBeChecked();
    expect(field('revision')).toHaveTextContent('def-5');
    expect(field('revision')).toHaveTextContent('epoch-3');
    expect(api.update).not.toHaveBeenCalled();
    expect(api.create).not.toHaveBeenCalled();
    expect(field('save')).toBeDisabled();
  });

  it('saves an exact decimal string above 2^53 and keeps every other part of the definition', async () => {
    api.update.mockResolvedValue({ measurement: definition({ definition_revision: 'def-6' }), epoch_transitioned: false });
    mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    fireEvent.change(field('modulus'), { target: { value: '9007199254740993' } });
    fireEvent.click(field('save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    expect(api.update).toHaveBeenCalledWith('meas-1', {
      ...definition(), counter_policy: { allow_negative: false, max_roll_over: 123, reset_threshold: 7, rollover_modulus: '9007199254740993', max_rate_per_second: '0.5' },
    });
    expect(typeof api.update.mock.calls[0][1].counter_policy?.rollover_modulus).toBe('string');
  });

  it('writes only the source that was edited', async () => {
    api.update.mockResolvedValue({ measurement: second(), epoch_transitioned: false });
    mount();
    await screen.findByTestId('counter-policy-source-tag-1');
    fireEvent.click(field('signed', 'tag-1'));
    fireEvent.click(field('save', 'tag-1'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    expect(api.update.mock.calls[0][0]).toBe('meas-2');
  });

  it('needs an explicit signed choice: it turns on negative differences and cannot keep a wrap policy', async () => {
    api.update.mockResolvedValue({ measurement: definition({ semantic_kind: 'signed_counter' }), epoch_transitioned: false });
    mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    fireEvent.click(field('signed'));
    expect(field('issues')).toHaveTextContent('step4.group.policy.issue.rollover_with_signed');
    expect(field('save')).toBeDisabled();
    fireEvent.change(field('modulus'), { target: { value: '' } });
    fireEvent.change(field('rate'), { target: { value: '' } });
    expect(field('save')).toBeEnabled();
    fireEvent.click(field('save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    const sent = api.update.mock.calls[0][1];
    expect(sent).toMatchObject({ semantic_kind: 'signed_counter', counter_policy: { allow_negative: true, max_roll_over: 123, reset_threshold: 7 } });
    expect(sent.counter_policy).not.toHaveProperty('rollover_modulus');
    expect(sent.counter_policy).not.toHaveProperty('max_rate_per_second');
  });

  it.each([['1e3'], ['-5'], ['0'], ['0.0'], ['abc'], ['1'.repeat(65)]])('refuses %s as a modulus before anything is sent', async (value) => {
    mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    fireEvent.change(field('modulus'), { target: { value } });
    expect(field('issues')).toHaveTextContent('step4.group.policy.issue.modulus_invalid');
    expect(field('modulus')).toHaveAttribute('aria-invalid', 'true');
    expect(field('save')).toBeDisabled();
    expect(api.update).not.toHaveBeenCalled();
  });

  it('refuses a rate limit without a modulus and a wrap policy on a float counter', async () => {
    api.list.mockResolvedValue([definition({ counter_policy: { allow_negative: false } })]);
    mount(false, [MEMBERS[0]]);
    await screen.findByTestId('counter-policy-source-tag-0');
    fireEvent.change(field('rate'), { target: { value: '5' } });
    expect(field('issues')).toHaveTextContent('step4.group.policy.issue.rate_needs_modulus');
    expect(field('save')).toBeDisabled();
  });

  it('locks the wrap policy of a float counter with an explanation', async () => {
    api.list.mockResolvedValue([definition({ numeric_encoding: 'float64', counter_policy: { allow_negative: false } })]);
    mount(false, [MEMBERS[0]]);
    await screen.findByTestId('counter-policy-source-tag-0');
    expect(field('modulus')).toBeDisabled();
    expect(field('rate')).toBeDisabled();
    expect(field('issues')).toHaveTextContent('step4.group.policy.issue.rollover_with_float');
  });

  it('offers no creation for a source without a confirmed counter meaning', async () => {
    api.list.mockResolvedValue([definition()]);
    mount();
    expect(await screen.findByTestId('counter-policy-missing-tag-1')).toHaveTextContent('step4.group.policy.missing');
    expect(screen.queryByTestId('counter-policy-source-tag-1')).not.toBeInTheDocument();
    expect(api.create).not.toHaveBeenCalled();
  });

  it('says so, without moving anywhere, when the Basic entry of that source is not on screen', async () => {
    api.list.mockResolvedValue([definition()]);
    mount();
    fireEvent.click(await screen.findByTestId('counter-policy-open-basic-tag-1'));
    expect(await screen.findByTestId('counter-policy-open-basic-unavailable-tag-1')).toHaveTextContent('step4.group.policy.open_basic_unavailable');
  });

  it('shows a safe message when the server refuses and keeps what was typed', async () => {
    api.update.mockRejectedValue(typed('counter_semantics_required'));
    const view = mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    fireEvent.change(field('modulus'), { target: { value: '128' } });
    fireEvent.click(field('save'));
    expect(await within(field('source')).findByRole('alert')).toBeInTheDocument();
    expect(field('modulus')).toHaveValue('128');
    expect(view.container.textContent).not.toMatch(/secret|INSERT INTO|postgres:/);
    expect(field('revision')).toHaveTextContent('def-5');
  });

  it('adopts the revision and epoch the server answers with and says when a new epoch began', async () => {
    api.update.mockResolvedValue({ measurement: definition({ definition_revision: 'def-6', series_epoch: 'epoch-4' }), epoch_transitioned: true });
    mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    fireEvent.change(field('modulus'), { target: { value: '128' } });
    api.list.mockResolvedValue([definition({ definition_revision: 'def-6', series_epoch: 'epoch-4', counter_policy: { allow_negative: false, max_roll_over: 123, reset_threshold: 7, rollover_modulus: '128', max_rate_per_second: '0.5' } }), second()]);
    fireEvent.click(field('save'));
    await waitFor(() => expect(field('revision')).toHaveTextContent('def-6'));
    expect(field('revision')).toHaveTextContent('epoch-4');
    expect(field('result')).toHaveTextContent('step4.group.policy.epoch_started');
  });

  it('locks every field while another operation is in flight', async () => {
    mount(true);
    await screen.findByTestId('counter-policy-source-tag-0');
    for (const name of ['signed', 'modulus', 'rate', 'save']) expect(field(name)).toBeDisabled();
  });
});
