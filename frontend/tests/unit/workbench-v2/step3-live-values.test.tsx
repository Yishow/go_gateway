import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useStep3LiveValues } from '../../../src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues';
import { pointAPI, runtimeAPI } from '../../../src/services/datalink';
import type { Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';
import type { Point as PersistedPoint, RuntimeValueEvent } from '../../../src/types/datalink';

vi.mock('../../../src/services/datalink', () => ({
  runtimeAPI: {
    getStreamUrl: vi.fn((deviceId: string, pointIds?: string[]) => {
      const suffix = pointIds?.join(',') ?? '';
      return `stream://${deviceId}?points=${suffix}`;
    }),
  },
  pointAPI: {
    list: vi.fn(),
    batchPoll: vi.fn(),
  },
}));

class MockEventSource {
  static instances: MockEventSource[] = [];

  readonly url: string;
  readonly listeners = new Map<string, Set<(event: Event) => void>>();
  onopen: (() => void) | null = null;
  onerror: (() => void) | null = null;
  closed = false;

  constructor(url: string) {
    this.url = url;
    MockEventSource.instances.push(this);
  }

  addEventListener(type: string, listener: (event: Event) => void) {
    const set = this.listeners.get(type) ?? new Set<(event: Event) => void>();
    set.add(listener);
    this.listeners.set(type, set);
  }

  removeEventListener(type: string, listener: (event: Event) => void) {
    this.listeners.get(type)?.delete(listener);
  }

  close() {
    this.closed = true;
  }

  emitOpen() {
    this.onopen?.();
  }

  emitValue(payload: RuntimeValueEvent) {
    const event = {
      data: JSON.stringify(payload),
    } as MessageEvent<string>;
    this.listeners.get('value')?.forEach((listener) => listener(event as unknown as Event));
  }
}

const points: Point[] = [
  {
    id: 'p-01',
    device_id: 'dev-01',
    rule_id: 'rule-01',
    rule_name: 'Holding Registers',
    name: 'SENSOR_1',
    address: '40001',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 0.1,
    _rule_offset: 0,
  },
  {
    id: 'p-02',
    device_id: 'dev-02',
    rule_id: 'rule-02',
    rule_name: 'Input Registers',
    name: 'SENSOR_2',
    address: '30001',
    data_type: 'int16',
    function: 'input_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
];

const mappings: Record<string, Mapping> = {
  'p-01': {
    point_id: 'p-01',
    tag_key: 'line01.temp.inlet',
    display_name: '進水溫度',
    unit: '°C',
    target_type: 'float64',
    scale: 0.1,
    offset: 0,
    enabled: true,
    persisted_point_id: 'persisted-point-01',
  },
  'p-02': {
    point_id: 'p-02',
    tag_key: 'line01.temp.outlet',
    display_name: '出水溫度',
    unit: '°C',
    target_type: 'float64',
    scale: 0.1,
    offset: 0,
    enabled: true,
    persisted_point_id: 'persisted-point-02',
  },
};

const mappingWithoutPersistedPointId: Mapping = {
  point_id: 'p-03',
  tag_key: 'line01.flow.q1',
  display_name: '流量 Q1',
  unit: 'L/min',
  target_type: 'int16',
  scale: 1,
  offset: 0,
  enabled: true,
};

const persistedPointsByDevice: Record<string, PersistedPoint[]> = {
  'dev-01': [
    {
      id: 'persisted-point-01',
      device_id: 'dev-01',
      name: 'SENSOR_1',
      description: '',
      data_type: 'int16',
      address: '40001',
      enabled: true,
      polling_group_id: '',
      last_value: 111,
      last_read_at: '2026-05-31T00:00:00Z',
      last_error: '',
      error_count: 0,
      created_at: '',
      updated_at: '',
    },
  ],
  'dev-02': [
    {
      id: 'persisted-point-02',
      device_id: 'dev-02',
      name: 'SENSOR_2',
      description: '',
      data_type: 'int16',
      address: '30001',
      enabled: true,
      polling_group_id: '',
      last_value: 222,
      last_read_at: '2026-05-31T00:00:00Z',
      last_error: '',
      error_count: 0,
      created_at: '',
      updated_at: '',
    },
  ],
};

const createdQueryClients: QueryClient[] = [];

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0 },
      mutations: { retry: false },
    },
  });
  createdQueryClients.push(queryClient);

  return function Wrapper({ children }: { children: ReactNode }) {
    return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
  };
}

describe('useStep3LiveValues', () => {
  beforeEach(() => {
    MockEventSource.instances = [];
    vi.clearAllMocks();
    vi.stubGlobal('EventSource', MockEventSource);
    vi.mocked(pointAPI.list).mockImplementation(async ({ device_id } = {}) => {
      return device_id ? persistedPointsByDevice[device_id] ?? [] : [];
    });
    vi.mocked(pointAPI.batchPoll).mockResolvedValue([]);
  });

  afterEach(() => {
    cleanup();
    createdQueryClients.splice(0).forEach((queryClient) => {
      queryClient.clear();
    });
    vi.unstubAllGlobals();
  });

  it('opens one runtime stream per device and merges live values by point id', () => {
    const { result } = renderHook(() => useStep3LiveValues(points, mappings), {
      wrapper: createWrapper(),
    });

    expect(runtimeAPI.getStreamUrl).toHaveBeenCalledTimes(2);
    expect(MockEventSource.instances).toHaveLength(2);
    expect(MockEventSource.instances.map((source) => source.url)).toEqual([
      'stream://dev-01?points=persisted-point-01',
      'stream://dev-02?points=persisted-point-02',
    ]);

    act(() => {
      MockEventSource.instances[0].emitOpen();
      MockEventSource.instances[1].emitOpen();
      MockEventSource.instances[1].emitValue({
        device_id: 'dev-02',
        point_id: 'persisted-point-02',
        address: '30001',
        raw_value: 321,
        transformed_value: 321,
        quality: 'good',
        stale: false,
        timestamp: '2026-05-31T00:00:00Z',
      });
    });

    expect(result.current.connectionByDevice['dev-01']).toBe('connected');
    expect(result.current.connectionByDevice['dev-02']).toBe('connected');
    expect(result.current.liveValues['p-02']?.raw_value).toBe(321);
    expect(result.current.rawValues['p-02']).toBe(321);
  });

  it('does not recreate runtime streams when unrelated mapping fields change', () => {
    const { rerender } = renderHook(
      ({ nextMappings }) => useStep3LiveValues(points, nextMappings),
      {
        wrapper: createWrapper(),
        initialProps: {
          nextMappings: mappings,
        },
      },
    );

    expect(runtimeAPI.getStreamUrl).toHaveBeenCalledTimes(2);
    expect(MockEventSource.instances).toHaveLength(2);

    rerender({
      nextMappings: {
        ...mappings,
        'p-01': {
          ...mappings['p-01'],
          save_state: 'saving',
        },
      },
    });

    expect(runtimeAPI.getStreamUrl).toHaveBeenCalledTimes(2);
    expect(MockEventSource.instances).toHaveLength(2);
    expect(MockEventSource.instances.every((source) => !source.closed)).toBe(true);
  });

  it('falls back to persisted point last_value when SSE has not arrived yet', async () => {
    const { result } = renderHook(() => useStep3LiveValues(points, mappings), {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(result.current.rawValues['p-01']).toBe(111);
      expect(result.current.rawValues['p-02']).toBe(222);
    });
  });

  it('uses batch point polling values when neither SSE nor last_value is available yet', async () => {
    vi.mocked(pointAPI.list).mockResolvedValue([
      {
        id: 'persisted-point-01',
        device_id: 'dev-01',
        name: 'SENSOR_1',
        description: '',
        data_type: 'int16',
        address: '40001',
        enabled: true,
        polling_group_id: '',
        last_value: null,
        last_read_at: '',
        last_error: '',
        error_count: 0,
        created_at: '',
        updated_at: '',
      },
    ]);
    vi.mocked(pointAPI.batchPoll).mockResolvedValue([
      {
        point_id: 'persisted-point-01',
        value: 345,
        transformed_value: 345,
        timestamp: '2026-05-31T00:00:00Z',
        quality: 192,
        stale: false,
        error: '',
      },
    ]);

    const { result } = renderHook(
      () => useStep3LiveValues([points[0]], { 'p-01': mappings['p-01'] }),
      {
        wrapper: createWrapper(),
      },
    );

    await waitFor(() => {
      expect(pointAPI.batchPoll).toHaveBeenCalledWith(['persisted-point-01']);
      expect(result.current.rawValues['p-01']).toBe(345);
    });
  });

  it('resolves persisted point ids from /points snapshot when mapping has not been saved yet', async () => {
    vi.mocked(pointAPI.list).mockResolvedValue([
      {
        id: 'persisted-point-03',
        device_id: 'dev-03',
        name: 'SENSOR_3',
        description: '',
        data_type: 'int16',
        address: '40003',
        enabled: true,
        polling_group_id: '',
        last_value: null,
        last_read_at: '',
        last_error: '',
        error_count: 0,
        created_at: '',
        updated_at: '',
      },
    ]);
    vi.mocked(pointAPI.batchPoll).mockResolvedValue([
      {
        point_id: 'persisted-point-03',
        value: 678,
        transformed_value: 678,
        timestamp: '2026-05-31T00:00:00Z',
        quality: 192,
        stale: false,
        error: '',
      },
    ]);

    const { result } = renderHook(
      () =>
        useStep3LiveValues(
          [
            {
              ...points[0],
              id: 'p-03',
              device_id: 'dev-03',
              address: '40003',
              name: 'SENSOR_3',
            },
          ],
          { 'p-03': mappingWithoutPersistedPointId },
        ),
      {
        wrapper: createWrapper(),
      },
    );

    await waitFor(() => {
      expect(pointAPI.batchPoll).toHaveBeenCalledWith(['persisted-point-03']);
      expect(result.current.rawValues['p-03']).toBe(678);
    });
  });
});
