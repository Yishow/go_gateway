import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { CircuitCounterConfirmation } from '@/features/datalink/workbench-v2/steps/step4/CircuitCounterConfirmation';
import { evaluateCandidateCounterConfirmation } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { MeasurementDefinition } from '@/types/measurement';
import { savedGroupState } from '../../fixtures/writeGroupState';

const { create, update } = vi.hoisted(() => ({ create: vi.fn(), update: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: create }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: update }),
}));
const candidates = buildGroupCandidates(savedGroupState()).candidates;
function definition(overrides: Partial<MeasurementDefinition> = {}): MeasurementDefinition {
  return {
    id: 'm-1', workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0',
    name: 'Water total', equipment_id: 'circuit-water', quantity: 'volume', unit: '', semantic_kind: 'counter',
    definition_revision: 'rev-9', source_binding_revision: 'binding-4', series_epoch: 'epoch-3', numeric_encoding: 'uint64',
    counter_policy: { allow_negative: false, rollover_modulus: '18446744073709551616', max_rate_per_second: '0.125', reset_threshold: 4 },
    ...overrides,
  };
}
beforeEach(() => vi.clearAllMocks());

it('preserves full saved definition and exact policies when confirming a missing unit through PUT', async () => {
  const def = definition();
  render(<CircuitCounterConfirmation workspaceId="ws-1" workspaceRevision="r-1" circuitId="circuit-water" readonly={false}
    unconfirmedStatuses={[{ candidate: candidates[0], confirmed: false, measurement: def }]} onConfirmed={vi.fn()} />);
  fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: 'm3' } });
  fireEvent.click(screen.getByTestId('circuit-confirm-counter-submit'));
  await waitFor(() => expect(update).toHaveBeenCalledOnce());
  expect(update.mock.calls[0][0]).toEqual({ id: 'm-1', def: { ...def, unit: 'm3' } });
});

it('confirms mixed source units independently without changing another source', async () => {
  const first = definition({ unit: 'm3' });
  const second = definition({ id: 'm-2', point_id: 'pt-1', tag_id: 'tag-1', unit: 'kWh' });
  render(<CircuitCounterConfirmation workspaceId="ws-1" workspaceRevision="r-1" circuitId="circuit-water" readonly={false}
    unconfirmedStatuses={[{ candidate: candidates[0], confirmed: false, measurement: first }, { candidate: candidates[1], confirmed: false, measurement: second }]}
    onConfirmed={vi.fn()} />);
  const source = screen.getByTestId('circuit-counter-confirmation-source-tag-0');
  fireEvent.click(within(source).getByTestId('circuit-confirm-counter-submit'));
  await waitFor(() => expect(update).toHaveBeenCalledOnce());
  expect(update.mock.calls[0][0]).toEqual({ id: 'm-1', def: first });
});

it('requires the operator to enter a unit instead of inventing energy or equipment identity', async () => {
  render(<CircuitCounterConfirmation workspaceId="ws-1" workspaceRevision="r-1" circuitId="explicit-circuit" readonly={false}
    unconfirmedStatuses={[{ candidate: candidates[0], confirmed: false }]} onConfirmed={vi.fn()} />);
  expect(screen.getByTestId('circuit-confirm-counter-unit')).toHaveValue('');
  expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
  fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: 'm3' } });
  fireEvent.click(screen.getByTestId('circuit-confirm-counter-submit'));
  await waitFor(() => expect(create).toHaveBeenCalledOnce());
  expect(create.mock.calls[0][0]).toMatchObject({ equipment_id: 'explicit-circuit', quantity: '', unit: 'm3' });
});

it.each([
  ['unsigned counter without optional policy', definition({ unit: 'm3', counter_policy: undefined }), true],
  ['signed counter without permission', definition({ unit: 'm3', semantic_kind: 'signed_counter', counter_policy: { allow_negative: false } }), false],
  ['float counter with exact rollover', definition({ unit: 'm3', numeric_encoding: 'float64' }), false],
  ['rate without modulus', definition({ unit: 'm3', counter_policy: { allow_negative: false, max_rate_per_second: '0.125' } }), false],
] as const)('matches backend confirmation: %s', (_name, def, confirmed) => {
  expect(evaluateCandidateCounterConfirmation('ws-1', candidates[0], [def]).confirmed).toBe(confirmed);
});

it('does not treat a text Tag as countable merely because it names a saved counter', () => {
  const text = candidates[2];
  const def = definition({ point_id: text.point_id, tag_id: text.tag_id, unit: 'count' });
  expect(evaluateCandidateCounterConfirmation('ws-1', text, [def]).issueCode).toBe('counter-type-unsupported');
});

it('does not present a unit confirmation as a repair for missing source binding evidence', () => {
  const refresh = vi.fn();
  render(<CircuitCounterConfirmation workspaceId="ws-1" workspaceRevision="r-1" circuitId="explicit-circuit" readonly={false}
    unconfirmedStatuses={[{ candidate: candidates[0], confirmed: false, measurement: definition({ unit: 'm3', source_binding_revision: '' }) }]}
    onConfirmed={refresh} />);
  expect(screen.queryByTestId('circuit-confirm-counter-submit')).not.toBeInTheDocument();
  expect(screen.getByTestId('circuit-counter-binding-unconfirmed')).toBeInTheDocument();
  fireEvent.click(screen.getByTestId('circuit-counter-reload-source'));
  expect(refresh).toHaveBeenCalledOnce();
  expect(create).not.toHaveBeenCalled(); expect(update).not.toHaveBeenCalled();
});
