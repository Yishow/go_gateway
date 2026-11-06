import { act, renderHook } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { useStudioV2DatabaseAutosave } from '@/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave';
import { workbenchV2Reducer, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
const mocks = vi.hoisted(() => ({ update: vi.fn(), config: vi.fn(() => ({ isSuccess: false })) }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceDatabase', () => ({
  useStudioV2DatabaseConfigQuery: () => mocks.config(), useStudioV2DatabaseTargetsQuery: () => ({ isSuccess: false }),
  useUpdateStudioV2DatabaseConfigMutation: () => ({ mutateAsync: mocks.update }), useUpsertStudioV2DatabaseTargetMutation: () => ({ mutateAsync: vi.fn() }),
}));
it.each(['workspace', 'workspace-aba', 'unmount', 'origin', 'origin-aba'] as const)('does not apply a pending ACK after %s changes', async change => {
  const originalBase = studioV2DatalinkApi.defaults.baseURL;
  const stateRef = { current: savedGroupState() }, dispatch = vi.fn();
  let finish!: (record: unknown) => void;
  mocks.update.mockReturnValue(new Promise(resolve => { finish = resolve; }));
  const hook = renderHook(({ owner }) => useStudioV2DatabaseAutosave({ dispatch }, stateRef, true, undefined, owner), { initialProps: { owner: 'ws-1' } });
  const action: WorkbenchV2Action = { type: 'updateDbConnector', patch: { name: 'Submitted' } };
  act(() => { stateRef.current = workbenchV2Reducer(stateRef.current, action); hook.result.current.afterDatabaseAction(action, stateRef.current); });
  if (change === 'workspace' || change === 'workspace-aba') hook.rerender({ owner: 'ws-2' });
  if (change === 'workspace-aba') hook.rerender({ owner: 'ws-1' });
  if (change === 'unmount') hook.unmount();
  if (change === 'origin' || change === 'origin-aba') studioV2DatalinkApi.defaults.baseURL = '/other/api';
  if (change === 'origin-aba') { hook.rerender({ owner: 'ws-1' }); studioV2DatalinkApi.defaults.baseURL = originalBase; hook.rerender({ owner: 'ws-1' }); }
  dispatch.mockClear();
  await act(async () => { finish({ ...stateRef.current.db.connector, id: stateRef.current.db.connector.connector_id, workspace_id: 'ws-1', identity_revision: 'new', setup_revision: 'new' }); });
  expect(dispatch).not.toHaveBeenCalledWith(expect.objectContaining({ type: 'acknowledgeDbConnectorSave' }));
  studioV2DatalinkApi.defaults.baseURL = originalBase;
});

it('does not hydrate invalidated cached config over an already acknowledged connector', async () => {
  const stateRef = { current: savedGroupState() }, dispatch = vi.fn();
  mocks.config.mockReturnValue({ isSuccess: true, isStale: true, data: {
    ...stateRef.current.db.connector, id: stateRef.current.db.connector.connector_id, workspace_id: 'ws-1',
    name: 'Old cached config', identity_revision: 'old', setup_revision: 'old', updated_at: 'old',
  } } as never);
  renderHook(() => useStudioV2DatabaseAutosave({ dispatch }, stateRef, true, undefined, 'ws-1'));
  await act(async () => {});
  expect(stateRef.current.db.connector.name).toBe('Line DB');
  expect(dispatch).not.toHaveBeenCalledWith(expect.objectContaining({ type: 'hydrateDbConnector' }));
  mocks.config.mockReturnValue({ isSuccess: false });
});
