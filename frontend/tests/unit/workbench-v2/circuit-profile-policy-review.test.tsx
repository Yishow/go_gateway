import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { CounterPolicyEditor } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/CounterPolicyEditor';
import { studioV2WorkspaceMeasurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
import type { MeasurementDefinition } from '@/types/measurement';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({
  studioV2WorkspaceMeasurementsAPI: { list: vi.fn(), update: vi.fn(), create: vi.fn() },
}));
const api = vi.mocked(studioV2WorkspaceMeasurementsAPI);
const def = (over: Partial<MeasurementDefinition> = {}): MeasurementDefinition => ({
  id: 'measurement-bound', workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1', equipment_id: 'eq-1',
  definition_revision: 'def-1', source_binding_revision: 'bind-1', series_epoch: 'epoch-1',
  name: 'Energy', quantity: 'energy', unit: 'kWh', semantic_kind: 'counter', numeric_encoding: 'uint64',
  counter_policy: { allow_negative: false, rollover_modulus: '65536', max_rate_per_second: '100' },
  created_at: '2026-10-01T00:00:00Z', updated_at: '2026-10-01T00:00:00Z', ...over,
});
const member = { device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1', measurement_id: 'measurement-bound' };
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(<QueryClientProvider client={client}><CounterPolicyEditor workspaceId="ws-1" members={[member]} readonly={false} /></QueryClientProvider>);
}
beforeEach(() => {
  vi.clearAllMocks();
  api.list.mockResolvedValue([def()]);
});

it('edits the definition actually referenced by the saved member, not the first definition of the source', async () => {
  api.list.mockResolvedValue([def({ id: 'measurement-other' }), def({ counter_policy: { allow_negative: false, rollover_modulus: '131072', max_rate_per_second: '100' } })]);
  api.update.mockResolvedValue({ measurement: def({ definition_revision: 'def-2' }), epoch_transitioned: true });
  mount();
  await screen.findByTestId('counter-policy-source-tag-1');
  expect(screen.getByTestId('counter-policy-modulus-tag-1')).toHaveValue('131072');
  fireEvent.change(screen.getByTestId('counter-policy-modulus-tag-1'), { target: { value: '262144' } });
  fireEvent.click(screen.getByTestId('counter-policy-save-tag-1'));
  await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
  expect(api.update.mock.calls[0][0]).toBe('measurement-bound');
});

it('does not substitute an unrelated definition when the saved measurement reference is missing', async () => {
  api.list.mockResolvedValue([def({ id: 'measurement-other' })]);
  mount();
  expect(await screen.findByTestId('counter-policy-missing-tag-1')).toBeInTheDocument();
  expect(screen.queryByTestId('counter-policy-source-tag-1')).not.toBeInTheDocument();
  expect(api.update).not.toHaveBeenCalled();
});

it('keeps a counter without a unit on the confirmation repair path, rather than calling it confirmed', async () => {
  api.list.mockResolvedValue([def({ unit: '' })]);
  mount();
  expect(await screen.findByTestId('counter-policy-open-basic-tag-1')).toBeInTheDocument();
  expect(screen.queryByTestId('counter-policy-source-tag-1')).not.toBeInTheDocument();
  expect(api.update).not.toHaveBeenCalled();
});

it('locks policy inputs during the actual save so later typing cannot disappear on revision refresh', async () => {
  let finish!: (value: { measurement: MeasurementDefinition; epoch_transitioned: boolean }) => void;
  api.update.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  mount();
  await screen.findByTestId('counter-policy-source-tag-1');
  fireEvent.change(screen.getByTestId('counter-policy-modulus-tag-1'), { target: { value: '131072' } });
  fireEvent.click(screen.getByTestId('counter-policy-save-tag-1'));
  await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
  try {
    for (const name of ['signed', 'modulus', 'rate', 'save']) expect(screen.getByTestId(`counter-policy-${name}-tag-1`), name).toBeDisabled();
  } finally {
    finish({ measurement: def(), epoch_transitioned: false });
  }
});
