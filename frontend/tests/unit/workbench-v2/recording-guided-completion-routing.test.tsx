import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { MemoryRouter, useLocation, useNavigate } from 'react-router-dom';
import { WorkbenchV2Shell } from '@/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { resolveRuntimeDashboardDevice } from '@/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice';
import { savedGroupState } from '../../fixtures/writeGroupState';

// Isolated Shell argument contract. The separate completion-router test mounts real Step4/native owners.
const handoff = vi.hoisted(() => ({ ids: undefined as string[] | undefined, activate: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4', () => ({
  Step4Database: ({ onCommit }: { onCommit?: (ids?: string[]) => void }) =>
    <button type="button" onClick={() => onCommit?.(handoff.ids)}>Continue exact scope</button>,
}));
const clients: QueryClient[] = [];
function state() {
  const value = savedGroupState({ current: 4, completed: new Set([1, 2, 3]), selectedRuleId: 'rule-B' });
  value.devices.push({ ...value.devices[0], id: 'dev-B', name: 'PLC B' });
  value.rules.push({ ...value.rules[0], id: 'rule-B', device_id: 'dev-B' });
  return value;
}
function Host() {
  const actions = useWorkbenchV2State(state()), location = useLocation(), navigate = useNavigate();
  return <><output data-testid="route">{location.pathname}{location.search}</output>
    <WorkbenchV2Shell state={actions.state} actions={actions} workspaceId="ws-1"
      navigateTo={target => navigate(target)} activateWorkspace={handoff.activate} />
  </>;
}
function show(ids?: string[]) {
  handoff.ids = ids;
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
  render(<QueryClientProvider client={client}><MemoryRouter initialEntries={['/studio/v2']}><Host /></MemoryRouter></QueryClientProvider>);
  fireEvent.click(screen.getByRole('button', { name: 'Continue exact scope' }));
}
beforeEach(() => { localStorage.clear(); sessionStorage.clear(); vi.clearAllMocks(); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
it.each<[string, string[]]>([
  ['empty', []], ['multiple', ['dev-1', 'dev-B']], ['local draft', ['draft-device']],
  ['foreign UUID', ['550e8400-e29b-41d4-a716-446655440000']], ['duplicate', ['dev-1', 'dev-1']],
])('refuses an explicit %s completion scope without selecting B or generic runtime', (_, ids) => {
  expect(resolveRuntimeDashboardDevice(state(), ids)).toBeNull();
  show(ids);
  expect(screen.getByTestId('route')).toHaveTextContent(/^\/studio\/v2$/);
  expect(handoff.activate).not.toHaveBeenCalled();
});
it('routes only the exact confirmed A even when B is selected, without activation', () => {
  show(['dev-1']);
  expect(screen.getByTestId('route')).toHaveTextContent(/^\/studio\/runtime\?device_id=dev-1$/);
  fireEvent.click(screen.getByRole('button', { name: 'Continue exact scope' }));
  expect(screen.getByTestId('route')).toHaveTextContent(/^\/studio\/runtime\?device_id=dev-1$/);
  expect(handoff.activate).not.toHaveBeenCalled();
});
it('keeps unspecified manual runtime navigation separate from explicit completion', () => {
  show();
  expect(screen.getByTestId('route')).toHaveTextContent(/^\/studio\/runtime\?device_id=dev-B$/);
  expect(handoff.activate).not.toHaveBeenCalled();
});
