import * as React from 'react';
import '../../../src/i18n/config';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { expect } from 'vitest';
import { studioV2MappingsAPI } from '../../../src/services/studioV2Mappings';
import { vi } from 'vitest';
import { useStudioV2MappingAutosave } from '../../../src/pages/datalink/workbench-v2/useStudioV2MappingAutosave';
import { INITIAL_STATE, workbenchV2Reducer, type WorkbenchV2Action } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { DEFAULT_RULE } from '../../../src/features/datalink/workbench-v2/state/defaults';
import { hydrateStudioV2Mapping } from '../../../src/features/datalink/workbench-v2/state/studioV2MappingAutosave';
import { mappingFixture, mappingPoints, runtimeApplied } from './mapping-autosave-page.testHarness';
import type { Point } from '../../../src/features/datalink/workbench-v2/state/types';
import type { StudioV2WorkspaceMappingRequest } from '../../../src/services/studioV2Mappings';
export const api = vi.mocked(studioV2MappingsAPI);
export const points: Point[] = [...Array.from({ length: 8 }, (_, i) => ({ ...mappingPoints[0], id: `a${i}`, address: String(40001 + i) })), { ...mappingPoints[1], id: 'b0' }];
export const records = points.map((point) => mappingFixture({ id: `mapping-${point.id}`, rule_id: point.rule_id, device_id: point.device_id, address: point.address, tag_key: `tag.${point.id}`, target_type: 'int16' }));
export function deferred<T>() {
  let resolve!: (value: T) => void;
  let reject!: (error: unknown) => void;
  const promise = new Promise<T>((res, rej) => { resolve = res; reject = rej; });
  return { promise, resolve, reject };
}
export function response(id: string, request: StudioV2WorkspaceMappingRequest) { return runtimeApplied(mappingFixture({ ...request, id })); }
export async function mountQueue(customPoints = points, customRecords = records, mutationRetries = 0) {
  let readConfirmedRecords = () => customRecords;
  api.list.mockImplementation(async () => readConfirmedRecords());
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: mutationRetries, retryDelay: 0 } } });
  const initial = { ...INITIAL_STATE, points: customPoints, rules: ['rule-A', 'rule-B'].map((id) => ({ ...DEFAULT_RULE, id, persisted: true, save_state: 'saved' as const })), mappings: Object.fromEntries(customPoints.map((point, i) => [point.id, hydrateStudioV2Mapping(point, customRecords[i])])) };
  const wrapper = ({ children }: { children: React.ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const hook = renderHook(() => {
    const [state, dispatch] = React.useReducer(workbenchV2Reducer, initial);
    const stateRef = React.useRef(state);
    readConfirmedRecords = () => customRecords.map((record, index) => {
      const mapping = stateRef.current.mappings[customPoints[index].id];
      return { ...record, ...(mapping?.persisted_value ?? {}), id: mapping?.mapping_id ?? record.id };
    }).filter((record) => Object.values(stateRef.current.mappings).some((mapping) => mapping.mapping_id === record.id));
    const autosave = useStudioV2MappingAutosave({ dispatch }, stateRef, true);
    return { state, query: autosave.mappingsQuery, dispatch(action: WorkbenchV2Action) {
      const before = stateRef.current;
      stateRef.current = workbenchV2Reducer(before, action);
      dispatch(action);
      autosave.afterMappingAction(before, action, stateRef.current);
    } };
  }, { wrapper });
  await waitFor(() => { expect(hook.result.current.query.isSuccess).toBe(true); });
  const send = (action: WorkbenchV2Action) => act(() => hook.result.current.dispatch(action));
  return { ...hook, send, client };
}
