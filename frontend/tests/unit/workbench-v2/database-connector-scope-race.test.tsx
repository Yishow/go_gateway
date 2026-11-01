import { act, renderHook } from '@testing-library/react';
import type { MutableRefObject } from 'react';
import { expect, it, vi } from 'vitest';
import { useStudioV2DatabaseAutosave } from '@/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave';
import { workbenchV2Reducer, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import { savedGroupState } from '../../fixtures/writeGroupState';

const mocks = vi.hoisted(() => ({ updateConfig: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceDatabase', () => ({
  useStudioV2DatabaseConfigQuery: () => ({ isSuccess: false }),
  useStudioV2DatabaseTargetsQuery: () => ({ isSuccess: false }),
  useUpdateStudioV2DatabaseConfigMutation: () => ({ mutateAsync: mocks.updateConfig }),
  useUpsertStudioV2DatabaseTargetMutation: () => ({ mutateAsync: vi.fn() }),
}));

it('does not replay an unknown connector write when a source gate becomes writable again', async () => {
  let reject!: (error: unknown) => void;
  mocks.updateConfig.mockReturnValue(new Promise((_, no) => { reject = no; }));
  const state = savedGroupState();
  state.rules = [{ ...state.rules[0], mc_persisted_points: {
    status: 'compatible', points: [], scope: { workspace_id: 'ws-1', workspace_revision: 'rev-1', revision_id: 'rule-1' },
  } }];
  state.devices = [{ ...state.devices[0], protocol: 'mc_3e' }];
  const stateRef = { current: state } as MutableRefObject<WorkbenchV2State>;
  const { result, rerender } = renderHook(({ allowed }) => useStudioV2DatabaseAutosave({ dispatch: vi.fn() }, stateRef, true, () => allowed), { initialProps: { allowed: true } });
  const perform = (action: WorkbenchV2Action) => act(() => {
    stateRef.current = workbenchV2Reducer(stateRef.current, action);
    result.current.afterDatabaseAction(action, stateRef.current);
  });
  perform({ type: 'updateDbConnector', patch: { name: 'L' } });
  perform({ type: 'updateDbConnector', patch: { name: 'Linux QA SQLite' } });
  rerender({ allowed: false });
  await act(async () => { reject(new Error('Network Error')); });
  rerender({ allowed: true });
  await act(async () => {});
  expect(mocks.updateConfig).toHaveBeenCalledTimes(1);
  expect(stateRef.current.db.connector).toMatchObject({ name: 'Linux QA SQLite', save_state: 'save-error' });
});
