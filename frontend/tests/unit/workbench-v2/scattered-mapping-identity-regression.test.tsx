import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import '../../../src/i18n/config';
import { DEFAULT_RULE } from '@/features/datalink/workbench-v2/state/defaults';
import { derivePoints } from '@/features/datalink/workbench-v2/state/sourceRule';
import { hydrateStudioV2Mapping } from '@/features/datalink/workbench-v2/state/studioV2MappingAutosave';
import { INITIAL_STATE, workbenchV2Reducer, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { useStudioV2MappingAutosave } from '@/pages/datalink/workbench-v2/useStudioV2MappingAutosave';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { mappingFixture } from './mapping-autosave-page.testHarness';

vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
beforeEach(() => vi.resetAllMocks());

it('keeps the retained FATEK item mapping when an earlier item is removed', async () => {
  const rule = { ...DEFAULT_RULE, id: 'rule-A', device_id: 'dev-A', persisted: true, save_state: 'saved' as const,
    layout_mode: 'mixed' as const, mixed_items: [
      { item_id: 'a', name: 'A', address: 'R0', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' as const },
      { item_id: 'b', name: 'B', address: 'R10', register_offset: 0, data_type: 'uint16', semantic_kind: 'gauge' as const },
    ] };
  const before = derivePoints(rule, rule.device_id, new Set(), 'fatek_fbs');
  const record = mappingFixture({ id: 'mapping-b', rule_id: rule.id, device_id: rule.device_id, address: 'R10', target_type: 'uint16' });
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([record]);
  vi.mocked(studioV2MappingsAPI.remove).mockResolvedValue({ runtime_apply_status: 'not_running' });
  const initial = { ...INITIAL_STATE, rules: [rule], points: before, mappings: { [before[1].id]: hydrateStudioV2Mapping(before[1], record) } };
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const wrapper = ({ children }: { children: React.ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const hook = renderHook(() => {
    const [state, dispatch] = React.useReducer(workbenchV2Reducer, initial);
    const stateRef = React.useRef(state);
    const autosave = useStudioV2MappingAutosave({ dispatch }, stateRef, true);
    return { state, query: autosave.mappingsQuery, send(action: WorkbenchV2Action) {
      const previous = stateRef.current;
      stateRef.current = workbenchV2Reducer(previous, action);
      dispatch(action);
      autosave.afterMappingAction(previous, action, stateRef.current);
    } };
  }, { wrapper });
  await waitFor(() => expect(hook.result.current.query.isSuccess).toBe(true));
  const after = derivePoints({ ...rule, mixed_items: [rule.mixed_items[1]] }, rule.device_id, new Set(), 'fatek_fbs');
  await act(async () => hook.result.current.send({ type: 'initMappingsForPoints', points: after }));
  expect(studioV2MappingsAPI.remove).not.toHaveBeenCalled();
  expect(hook.result.current.state.mappings[after[0].id]).toMatchObject({ mapping_id: 'mapping-b', address: 'R10' });
  hook.unmount();
  client.clear();
});
