import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useStep3LiveValues } from '../../../src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues';
import type { Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';
import { pointAPI } from '../../../src/services/datalink';
import { pointKeys } from '../../../src/hooks/datalink/keys';
import type { Point as PersistedPoint, PollResult } from '../../../src/types/datalink';

vi.mock('../../../src/services/datalink', () => ({
  pointAPI: { list: vi.fn(), batchPoll: vi.fn() },
}));

const mcId = '19b1b25a-56e7-4f95-83af-9793eeed32cc';
const mc: Point = {
  id: mcId, device_id: 'mc-device', rule_id: 'mc-rule', rule_name: 'MC',
  name: 'MC__D0', address: 'D0', data_type: 'uint16', function: '03', width: 1,
  enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0,
};
function mapping(point: Point, persistedId?: string): Mapping {
  return {
    point_id: point.id, tag_key: `tag.${point.id}`, display_name: point.name,
    unit: '', target_type: 'uint16', scale: 1, offset: 0, enabled: true,
    persisted_point_id: persistedId,
  };
}
function snapshot(deviceId = mc.device_id): PersistedPoint {
  return {
    id: mcId, device_id: deviceId, name: mc.name, address: mc.address,
    data_type: 'uint16', description: '', enabled: true, polling_group_id: '',
    last_value: null, last_read_at: '', last_error: '', error_count: 0,
    created_at: '', updated_at: '',
  };
}
function poll(pointId: string, value: number): PollResult {
  return { point_id: pointId, value, transformed_value: value, quality: 192,
    timestamp: '2026-11-01T22:30:00Z', stale: false, error: '' };
}
const clients: QueryClient[] = [];
function setup() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } });
  clients.push(client);
  const wrapper = ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={client}>{children}</QueryClientProvider>
  );
  return { client, wrapper };
}

beforeEach(() => {
  vi.resetAllMocks();
  vi.stubGlobal('EventSource', undefined);
  vi.mocked(pointAPI.list).mockResolvedValue([]);
  vi.mocked(pointAPI.batchPoll).mockResolvedValue([]);
});
afterEach(() => {
  cleanup();
  clients.splice(0).forEach(client => client.clear());
  vi.unstubAllGlobals();
});

describe('Step3 persisted point identity', () => {
  it('polls canonical MC identity once alongside three distinct persisted aliases', async () => {
    const aliases = ['mb', 'fatek', 'mb-second'].map((id, index) => ({
      ...mc, id: `local-${id}`, device_id: id, address: String(index),
    }));
    const points = [mc, ...aliases];
    const mappings = {
      [mcId]: mapping(mc, mcId),
      'local-mb': mapping(aliases[0], 'saved-mb'),
      'local-fatek': mapping(aliases[1], 'saved-fatek'),
      'local-mb-second': mapping(aliases[2], 'saved-mb-second'),
    };
    vi.mocked(pointAPI.batchPoll).mockResolvedValue([
      poll(mcId, 125), poll('saved-mb', 11), poll('saved-fatek', 22), poll('saved-mb-second', 33),
    ]);
    const { result } = renderHook(() => useStep3LiveValues(points, mappings), setup());
    await waitFor(() => expect(pointAPI.batchPoll).toHaveBeenCalledWith([
      mcId, 'saved-mb', 'saved-fatek', 'saved-mb-second',
    ]));
    await waitFor(() => expect(result.current.rawValues).toEqual({
      [mcId]: 125, 'local-mb': 11, 'local-fatek': 22, 'local-mb-second': 33,
    }));
    expect(pointAPI.batchPoll).toHaveBeenCalledTimes(1);
  });

  it('polls a same-ID point only after matching persisted snapshot evidence arrives', async () => {
    vi.mocked(pointAPI.list).mockResolvedValue([snapshot()]);
    vi.mocked(pointAPI.batchPoll).mockResolvedValue([poll(mcId, 125)]);
    const { result } = renderHook(() => useStep3LiveValues([mc], { [mcId]: mapping(mc) }), setup());
    await waitFor(() => expect(result.current.rawValues[mcId]).toBe(125));
    expect(pointAPI.batchPoll).toHaveBeenCalledExactlyOnceWith([mcId]);
  });

  it.each(['absent', 'foreign-device'] as const)('does not poll a bare UUID when persisted evidence is %s', async (evidence) => {
    vi.mocked(pointAPI.list).mockResolvedValue(evidence === 'absent' ? [] : [snapshot('foreign-device')]);
    const { client, wrapper } = setup();
    const { result } = renderHook(() => useStep3LiveValues([mc], { [mcId]: mapping(mc) }), { wrapper });
    await waitFor(() => expect(client.getQueryState(pointKeys.list({ device_id: mc.device_id }))?.status).toBe('success'));
    expect(pointAPI.batchPoll).not.toHaveBeenCalled();
    expect(result.current.rawValues).toEqual({});
  });

  it('tracks persisted evidence arriving and being revoked even when the subscription ID stays equal', async () => {
    vi.mocked(pointAPI.batchPoll).mockResolvedValue([poll(mcId, 125)]);
    const { client, wrapper } = setup();
    const { result, rerender } = renderHook(
      ({ saved }) => useStep3LiveValues([mc], { [mcId]: mapping(mc, saved ? mcId : undefined) }),
      { wrapper, initialProps: { saved: false } },
    );
    await waitFor(() => expect(client.getQueryState(pointKeys.list({ device_id: mc.device_id }))?.status).toBe('success'));
    expect(pointAPI.batchPoll).not.toHaveBeenCalled();
    rerender({ saved: true });
    await waitFor(() => expect(result.current.rawValues[mcId]).toBe(125));
    expect(pointAPI.batchPoll).toHaveBeenCalledExactlyOnceWith([mcId]);
    rerender({ saved: false });
    expect(result.current.rawValues).toEqual({});
    expect(pointAPI.batchPoll).toHaveBeenCalledTimes(1);
  });

  it('does not project a late old-scope poll response into the next selected point', async () => {
    let resolveOld: ((value: PollResult[]) => void) | undefined;
    vi.mocked(pointAPI.batchPoll).mockImplementation(ids => ids.includes('saved-old')
      ? new Promise(resolve => { resolveOld = resolve; })
      : Promise.resolve([poll('saved-new', 200)]));
    const oldPoint = { ...mc, id: 'local-old' };
    const newPoint = { ...mc, id: 'local-new', device_id: 'new-device', address: 'D1' };
    const { result, rerender } = renderHook(
      ({ point, savedId }) => useStep3LiveValues([point], { [point.id]: mapping(point, savedId) }),
      { ...setup(), initialProps: { point: oldPoint, savedId: 'saved-old' } },
    );
    await waitFor(() => expect(pointAPI.batchPoll).toHaveBeenCalledWith(['saved-old']));
    rerender({ point: newPoint, savedId: 'saved-new' });
    await waitFor(() => expect(result.current.rawValues).toEqual({ 'local-new': 200 }));
    await act(async () => { resolveOld?.([poll('saved-old', 999)]); });
    expect(result.current.rawValues).toEqual({ 'local-new': 200 });
    expect(pointAPI.batchPoll).toHaveBeenCalledTimes(2);
  });
});
