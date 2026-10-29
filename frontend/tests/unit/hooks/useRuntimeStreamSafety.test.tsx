import { act, renderHook } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { useRuntimeStream } from '@/hooks/datalink/useRuntimeStream';
vi.mock('@/services/datalink', () => ({ runtimeAPI: { getStreamUrl: () => '/runtime/stream' } }));
class Stream {
  static latest: Stream;
  listeners = new Map<string, (event: Event) => void>();
  onopen = null;
  onerror = null;
  constructor() { Stream.latest = this; }
  addEventListener(type: string, listener: (event: Event) => void) { this.listeners.set(type, listener); }
  removeEventListener(type: string) { this.listeners.delete(type); }
  close() {}
  emit(raw_value: unknown) {
    this.listeners.get('value')?.(new MessageEvent('value', { data: JSON.stringify({
      device_id: 'device-A', point_id: 'point-1', address: '40001', raw_value, transformed_value: 1,
      quality: 'good', stale: false, timestamp: '2026-10-28T00:00:00Z',
    }) }));
  }
}
afterEach(() => vi.unstubAllGlobals());
describe('legacy runtime stream hook safe parser consumer', () => {
  it.each(['path=/private/db', 'path=C:/private/db', '(/private/db)', '\\\\server\\share\\private', './private/db', '../private/db'])(
    'does not forward private path %s into live operator values', (path) => {
      vi.stubGlobal('EventSource', Stream);
      const { result } = renderHook(() => useRuntimeStream({ deviceId: 'device-A' }));
      act(() => Stream.latest.emit(path));
      expect(result.current.connectionState).toBe('degraded');
      expect(result.current.liveValues).toEqual({});
    },
  );
  it('degrades hostile values and accepts ordinary observations without forwarding secrets', () => {
    vi.stubGlobal('EventSource', Stream);
    const { result } = renderHook(() => useRuntimeStream({ deviceId: 'device-A' }));
    act(() => Stream.latest.emit('postgres://user:secret@host:5432'));
    expect(result.current.connectionState).toBe('degraded');
    expect(result.current.liveValues).toEqual({});
    act(() => Stream.latest.emit('9007199254740993'));
    expect(result.current.connectionState).toBe('connected');
    expect(result.current.liveValues['point-1'].raw_value).toBe('9007199254740993');
  });
});
