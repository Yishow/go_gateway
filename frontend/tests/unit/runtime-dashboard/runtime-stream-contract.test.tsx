import { act, renderHook, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useRuntimeDashboardStream } from '../../../src/features/datalink/runtime-dashboard/useRuntimeStream';

const { getStreamUrl } = vi.hoisted(() => ({ getStreamUrl: vi.fn(() => '/datalink/runtime/stream?device_id=device-1') }));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (_key: string, fallback: string) => fallback }),
}));

vi.mock('../../../src/services/datalink', () => ({
  runtimeAPI: { getStreamUrl },
}));

class MockEventSource {
  static instances: MockEventSource[] = [];
  onopen: (() => void) | null = null;
  onerror: (() => void) | null = null;
  private listeners = new Map<string, (event: Event) => void>();
  closed = false;

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

  emit(type: string, payload: unknown) {
    this.listeners.get(type)?.({ data: JSON.stringify(payload) } as MessageEvent<string>);
  }
}

describe('useRuntimeDashboardStream SSE contract', () => {
  beforeEach(() => {
    MockEventSource.instances = [];
    vi.stubGlobal('EventSource', MockEventSource);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('accepts named runtime events and ignores malformed payloads without throwing', async () => {
    const { result } = renderHook(() => useRuntimeDashboardStream({
      deviceId: 'device-1',
      pointIds: ['point-1'],
    }));
    const source = MockEventSource.instances[0];

    act(() => {
      source.onopen?.();
      source.emit('value', '{"raw":"not an object"}');
      source.emit('status', 'bad status');
      source.emit('stream_state', { stream_state: { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false } });
      source.emit('value', {
        device_id: 'device-1',
        point_id: 'point-1',
        address: '40001',
        raw_value: 12,
        transformed_value: 12,
        quality: 'good',
        stale: false,
        timestamp: '2026-08-25T00:00:00Z',
      });
    });

    await waitFor(() => expect(result.current.connectionState).toBe('connected'));
    await waitFor(() => expect(result.current.liveValues['point-1']?.raw_value).toBe(12));
    expect(result.current.latestStatus).toBeNull();
  });

  it('does not claim connected on transport open until a complete server event arrives', () => {
    const { result } = renderHook(() => useRuntimeDashboardStream({
      deviceId: 'device-1',
      pointIds: [],
    }));
    const source = MockEventSource.instances[0];

    act(() => source.onopen?.());
    expect(result.current.connectionState).toBe('connecting');
  });

  it('does not replay a value buffered by a previous stream generation', async () => {
    vi.useFakeTimers();
    const { result } = renderHook(() => useRuntimeDashboardStream({
      deviceId: 'device-1',
      pointIds: ['point-1'],
    }));
    const firstSource = MockEventSource.instances[0];
    act(() => firstSource.emit('value', {
      device_id: 'device-1', point_id: 'point-1', address: '40001', raw_value: 7,
      transformed_value: 7, quality: 'good', stale: false, timestamp: '2026-08-25T00:00:00Z',
    }));
    act(() => result.current.reconnect());
    const secondSource = MockEventSource.instances[1];
    await act(async () => vi.advanceTimersByTimeAsync(200));
    expect(result.current.liveValues).toEqual({});
    act(() => secondSource.emit('value', {
      device_id: 'device-1', point_id: 'point-1', address: '40001', raw_value: 8,
      transformed_value: 8, quality: 'good', stale: false, timestamp: '2026-08-25T00:00:01Z',
    }));
    await act(async () => vi.advanceTimersByTimeAsync(200));
    expect(result.current.liveValues['point-1']?.raw_value).toBe(8);
    vi.useRealTimers();
  });

  it('closes the EventSource on unmount and ignores late events', () => {
    const { result, unmount } = renderHook(() => useRuntimeDashboardStream({
      deviceId: 'device-1',
      pointIds: [],
    }));
    const source = MockEventSource.instances[0];
    act(() => source.onopen?.());
    unmount();

    expect(source.closed).toBe(true);
    act(() => {
      source.onopen?.();
      source.emit('stream_state', { stream_state: { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false } });
    });
    expect(result.current.connectionState).toBe('connecting');
  });

  it('retains safe stream recovery metadata and degrades on malformed events', async () => {
    const { result } = renderHook(() => useRuntimeDashboardStream({
      deviceId: 'device-1',
      pointIds: [],
    }));
    const source = MockEventSource.instances[0];

    act(() => {
      source.emit('stream_state', {
        device_id: 'device-1',
        stream_state: { state: 'unavailable', empty: false, degraded: false, unavailable: true, stale: false },
        code: 'runtime_stream_unavailable',
        action: 'retry runtime stream',
        request_id: 'req-runtime-1',
        retryable: true,
      });
    });

    await waitFor(() => expect(result.current.connectionState).toBe('unavailable'));
    expect(result.current.streamRecovery).toEqual({
      code: 'runtime_stream_unavailable',
      action: 'retry runtime stream',
      requestId: 'req-runtime-1',
      retryable: true,
    });

    act(() => source.emit('status', { device_id: 'device-1' }));
    await waitFor(() => expect(result.current.connectionState).toBe('degraded'));
  });

  it('reconnects by rebuilding the EventSource when the stream recovery action is invoked', async () => {
    const { result } = renderHook(() => useRuntimeDashboardStream({
      deviceId: 'device-1',
      pointIds: ['point-1'],
    }));
    const firstSource = MockEventSource.instances[0];

    await act(async () => {
      await result.current.reconnect();
    });

    expect(firstSource.closed).toBe(true);
    expect(MockEventSource.instances).toHaveLength(2);
    expect(MockEventSource.instances[1].url).toBe(firstSource.url);
  });
});
