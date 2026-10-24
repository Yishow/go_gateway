import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CounterPolicyEditor, type CounterPolicySource } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/CounterPolicyEditor';
import { studioV2WorkspaceMeasurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
import type { MeasurementDefinition } from '@/types/measurement';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({
  studioV2WorkspaceMeasurementsAPI: { list: vi.fn(), update: vi.fn(), create: vi.fn() },
}));

const api = vi.mocked(studioV2WorkspaceMeasurementsAPI);
const def = (over: Partial<MeasurementDefinition> = {}): MeasurementDefinition => ({
  id: 'meas-1', workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', equipment_id: 'eq-A',
  definition_revision: 'def-5', source_binding_revision: 'bind-5', series_epoch: 'epoch-3', name: 'Energy', quantity: 'energy', unit: 'kWh',
  semantic_kind: 'counter', numeric_encoding: 'uint64', counter_policy: { allow_negative: false, max_roll_over: 123 },
  created_at: '2026-10-01T00:00:00Z', updated_at: '2026-10-02T00:00:00Z', ...over,
});
const member = { device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', measurement_id: 'meas-1' };
const field = (name: string) => screen.getByTestId(`counter-policy-${name}-tag-0`);

function mount(members: CounterPolicySource[] = [member], workspaceId = 'ws-1', client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })) {
  const element = (id: string, list: CounterPolicySource[]) => <QueryClientProvider client={client}><CounterPolicyEditor workspaceId={id} members={list} readonly={false} /></QueryClientProvider>;
  const view = render(element(workspaceId, members));
  return { ...view, switchTo: (id: string, list: CounterPolicySource[]) => view.rerender(element(id, list)) };
}

beforeEach(() => {
  vi.clearAllMocks();
  api.list.mockResolvedValue([def()]);
});

describe('The editor works only on the definition the saved member is bound to', () => {
  it.each([
    ['belongs to another source', def({ point_id: 'pt-9', tag_id: 'tag-9' })],
    ['belongs to another workspace', def({ workspace_id: 'ws-2' })],
    ['is not a counter', def({ semantic_kind: 'gauge' })],
    ['has no definition revision', def({ definition_revision: '' })],
    ['has no series epoch', def({ series_epoch: ' ' })],
    ['has no source binding revision', def({ source_binding_revision: '' })],
  ])('sends the operator to Basic when the bound definition %s', async (_name, bound) => {
    api.list.mockResolvedValue([bound]);
    mount();
    expect(await screen.findByTestId('counter-policy-missing-tag-0')).toBeInTheDocument();
    expect(screen.getByTestId('counter-policy-open-basic-tag-0')).toBeInTheDocument();
    expect(screen.queryByTestId('counter-policy-source-tag-0')).not.toBeInTheDocument();
    expect(api.update).not.toHaveBeenCalled();
    expect(api.create).not.toHaveBeenCalled();
  });

  it('never binds a source on its own: a member without a persisted measurement ID gets no form even if a counter exists', async () => {
    mount([{ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0' }]);
    expect(await screen.findByTestId('counter-policy-missing-tag-0')).toBeInTheDocument();
    expect(screen.queryByTestId('counter-policy-source-tag-0')).not.toBeInTheDocument();
    expect(api.update).not.toHaveBeenCalled();
  });
});

describe('A policy that is not confirmed is never presented as the confirmed one', () => {
  it('shows a signed counter whose allow-negative flag disagrees as unconfirmed and repairs it on an explicit save', async () => {
    api.list.mockResolvedValue([def({ semantic_kind: 'signed_counter', counter_policy: { allow_negative: false, max_roll_over: 123 } })]);
    api.update.mockResolvedValue({ measurement: def({ semantic_kind: 'signed_counter', counter_policy: { allow_negative: true } }), epoch_transitioned: false });
    mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    expect(screen.getByTestId('counter-policy-unconfirmed-tag-0')).toHaveTextContent('step4.group.policy.unconfirmed');
    expect(field('signed')).toBeChecked();
    expect(field('save')).toBeEnabled();
    fireEvent.click(field('save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    expect(api.update.mock.calls[0][1]).toMatchObject({ semantic_kind: 'signed_counter', unit: 'kWh', counter_policy: { allow_negative: true, max_roll_over: 123 } });
  });

  it('shows a correct signed counter as confirmed with nothing to save', async () => {
    api.list.mockResolvedValue([def({ semantic_kind: 'signed_counter', counter_policy: { allow_negative: true } })]);
    mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    expect(screen.queryByTestId('counter-policy-unconfirmed-tag-0')).not.toBeInTheDocument();
    expect(field('signed')).toBeChecked();
    expect(field('save')).toBeDisabled();
  });

  it('lets an invalid wrap policy on a float counter be cleared, and only then be saved', async () => {
    api.list.mockResolvedValue([def({ numeric_encoding: 'float64', counter_policy: { allow_negative: false, rollover_modulus: '65536' } })]);
    api.update.mockResolvedValue({ measurement: def({ numeric_encoding: 'float64' }), epoch_transitioned: false });
    mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    expect(screen.getByTestId('counter-policy-unconfirmed-tag-0')).toBeInTheDocument();
    expect(field('modulus')).toBeEnabled();
    expect(field('save')).toBeDisabled();
    fireEvent.change(field('modulus'), { target: { value: '' } });
    expect(field('save')).toBeEnabled();
    fireEvent.click(field('save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    expect(api.update.mock.calls[0][1].counter_policy).not.toHaveProperty('rollover_modulus');
    expect(api.update.mock.calls[0][1].numeric_encoding).toBe('float64');
  });

  it('keeps a correct float counter confirmed with its wrap fields locked', async () => {
    api.list.mockResolvedValue([def({ numeric_encoding: 'float64', counter_policy: { allow_negative: false } })]);
    mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    expect(screen.queryByTestId('counter-policy-unconfirmed-tag-0')).not.toBeInTheDocument();
    expect(field('modulus')).toBeDisabled();
    expect(field('rate')).toBeDisabled();
  });
});

describe('A save in flight stays one request for one scope', () => {
  it('ignores a second click while the first request is pending and keeps the typed value', async () => {
    let finish!: (value: { measurement: MeasurementDefinition; epoch_transitioned: boolean }) => void;
    api.update.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
    mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    fireEvent.change(field('modulus'), { target: { value: '128' } });
    fireEvent.click(field('save'));
    fireEvent.click(field('save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    expect(field('modulus')).toHaveValue('128');
    finish({ measurement: def(), epoch_transitioned: false });
  });

  it('does not show the result of a request from another workspace once the scope changed', async () => {
    let finish!: (value: { measurement: MeasurementDefinition; epoch_transitioned: boolean }) => void;
    api.update.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
    // The other workspace has its own definition of a source with the same IDs, so its form is on screen when the late answer arrives.
    api.list.mockResolvedValue([def(), def({ id: 'meas-9', workspace_id: 'ws-2' })]);
    const view = mount();
    await screen.findByTestId('counter-policy-source-tag-0');
    fireEvent.change(field('modulus'), { target: { value: '128' } });
    fireEvent.click(field('save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    view.switchTo('ws-2', [{ ...member, measurement_id: 'meas-9' }]);
    await waitFor(() => expect(field('modulus')).toHaveValue(''));
    finish({ measurement: def(), epoch_transitioned: true });
    await waitFor(() => expect(api.list).toHaveBeenCalled());
    await new Promise((resolve) => setTimeout(resolve, 0));
    expect(screen.queryByTestId('counter-policy-result-tag-0')).not.toBeInTheDocument();
  });
});
