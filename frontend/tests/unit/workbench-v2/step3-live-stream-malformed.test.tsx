import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook } from '@testing-library/react';
import type { ReactNode } from 'react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useStep3LiveValues } from '../../../src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues';
import type { Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';
import { pointAPI } from '../../../src/services/datalink';

vi.mock('../../../src/services/datalink', () => ({
  runtimeAPI: { getStreamUrl: vi.fn((deviceId: string) => `stream://${deviceId}`) },
  pointAPI: { list: vi.fn().mockResolvedValue([]), batchPoll: vi.fn().mockResolvedValue([]) },
}));

class MockEventSource {
  static instances: MockEventSource[] = [];
  readonly listeners = new Map<string, (event: Event) => void>();
  closed = false;
  onopen: (() => void) | null = null;
  onerror: (() => void) | null = null;

  constructor(public readonly url: string) {
    MockEventSource.instances.push(this);
  }

  addEventListener(type: string, listener: (event: Event) => void) {
    this.listeners.set(type, listener);
  }

  removeEventListener(type: string) {
    this.listeners.delete(type);
  }

  close() {
    this.closed = true;
  }

  emitRaw(type: string, data: string) {
    this.listeners.get(type)?.({ data } as MessageEvent<string> as Event);
  }
}

const point: Point = {
  id: 'p-01', device_id: 'dev-01', rule_id: 'rule-01', rule_name: 'Holding Registers',
  name: 'SENSOR_1', address: '40001', data_type: 'int16', function: 'holding_register',
  width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0,
};
const mapping: Mapping = {
  point_id: 'p-01', tag_key: 'line.temp', display_name: 'Temperature', unit: 'C',
  target_type: 'float64', scale: 1, offset: 0, enabled: true, persisted_point_id: 'persisted-01',
};

function wrapper({ children }: { children: ReactNode }) {
  return <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>{children}</QueryClientProvider>;
}

describe('Step 3 live stream malformed payload recovery', () => {
  beforeEach(() => {
    MockEventSource.instances = [];
    vi.useFakeTimers();
    vi.stubGlobal('EventSource', MockEventSource);
    vi.mocked(pointAPI.list).mockResolvedValue([]);
    vi.mocked(pointAPI.batchPoll).mockResolvedValue([]);
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
  });

  it('does not throw on malformed JSON and enters bounded reconnecting state', async () => {
    const { result } = renderHook(() => useStep3LiveValues([point], { 'p-01': mapping }), { wrapper });
    expect(MockEventSource.instances).toHaveLength(1);

    act(() => MockEventSource.instances[0].emitRaw('value', '{bad-json'));
    expect(result.current.connectionByDevice['dev-01']).toBe('reconnecting');

    await act(async () => vi.advanceTimersByTimeAsync(1000));
    expect(MockEventSource.instances).toHaveLength(2);
    expect(result.current.connectionByDevice['dev-01']).toBe('reconnecting');
  });

  it('keeps the last value truthful as stale when a valid stream later sends a bad shape', async () => {
    const { result } = renderHook(() => useStep3LiveValues([point], { 'p-01': mapping }), { wrapper });
    expect(MockEventSource.instances).toHaveLength(1);
    const source = MockEventSource.instances[0];

    act(() => source.emitRaw('value', JSON.stringify({
      device_id: 'dev-01', point_id: 'persisted-01', address: '40001', raw_value: 12,
      transformed_value: 12, quality: 'good', stale: false, timestamp: '2026-08-25T00:00:00Z',
    })));
    expect(result.current.liveValues['p-01']?.raw_value).toBe(12);

    act(() => source.emitRaw('value', JSON.stringify({ point_id: 'persisted-01' })));
    expect(result.current.connectionByDevice['dev-01']).toBe('stale');
    expect(result.current.rawValues['p-01']).toBe(12);
  });
});
