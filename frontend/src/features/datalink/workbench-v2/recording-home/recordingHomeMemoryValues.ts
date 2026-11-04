import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { pointAPI, runtimeAPI } from '../../../../services/datalink';
import type { ProtocolType, RuntimeValueEvent } from '../../../../types/datalink';
import type { Point } from '../../../../types/point';
import { addressParser } from '../../../../utils/addressParser';
import { parseMessageEventRecord, parseRuntimeValueRecord } from '../../../../utils/safeJson';

/** The gateway's last read of one persisted point; the browser never reads the device itself. */
export interface MemoryLiveValue { pointId?: string; value?: string; error?: string; readAt?: string; stale?: boolean }
/** `live`: pushed over SSE; `polling`: SSE unavailable, the stored list is re-read. */
export type MemoryValuesStatus = 'loading' | 'error' | 'live' | 'polling';
export interface MemoryValues {
  byCell: ReadonlyMap<string, MemoryLiveValue>; status: MemoryValuesStatus; updatedAt?: number;
  /** The view reports which cell keys are on screen; only those points are streamed. */
  setVisible?: (keys: readonly string[]) => void;
}

// One screen refresh per second: pushed readings are batched per second, and the SSE fallback polls at the same pace.
const FALLBACK_POLL_MS = 1000, VISIBLE_DEBOUNCE_MS = 300, FLUSH_MS = 1000;

/** Area + numeric offset + data type, so alias spellings of one address meet and a value
 * decoded as another type is never shown under this range's type. */
export function memoryValueKey(area: string, offset: number, dataType: string): string {
  return `${area}:${offset}:${dataType}`;
}

export function formatMemoryValue(raw: unknown): string | undefined {
  if (raw === null || raw === undefined || raw === '') return undefined;
  if (typeof raw === 'boolean') return raw ? '1' : '0';
  const text = typeof raw === 'string' ? raw : typeof raw === 'number' ? String(raw) : JSON.stringify(raw);
  const numeric = Number(text);
  // Long float tails (0.30000000000000004) do not fit a cell; integers and short decimals stay exact.
  if (text.trim() !== '' && Number.isFinite(numeric) && !Number.isInteger(numeric) && text.length > 10) return String(Number(numeric.toPrecision(7)));
  return text;
}

/** Indexes persisted points of one device by memory cell. Points whose address the
 * protocol parser rejects are skipped rather than guessed. */
export function indexMemoryValues(points: readonly Point[], deviceId: string, protocol: ProtocolType): Map<string, MemoryLiveValue> {
  const result = new Map<string, MemoryLiveValue>();
  for (const point of points) {
    if (point.device_id !== deviceId || !point.enabled) continue;
    let parsed;
    try { parsed = addressParser.parse(point.address.trim().toUpperCase(), protocol); } catch { continue; }
    result.set(memoryValueKey(parsed.area, parsed.startNumber, point.data_type),
      { pointId: point.id, value: formatMemoryValue(point.last_value), error: point.last_error?.trim() || undefined, readAt: point.last_read_at || undefined });
  }
  return result;
}

/** A pushed reading replaces the stored one; bad quality is a failed read, stale keeps the value but marks it. */
export function applyLiveEvent(base: MemoryLiveValue, event: RuntimeValueEvent): MemoryLiveValue {
  const failed = event.quality === 'bad' || event.quality === 'invalid' || event.quality === 'missing' || event.value_state === 'unsupported';
  const value = formatMemoryValue(event.transformed_value ?? event.raw_value);
  return { ...base, value: failed ? base.value : value, error: failed ? event.quality : undefined, readAt: event.timestamp, stale: event.stale || event.quality === 'stale' };
}

/** Baseline from the stored point list, then SSE pushes for the on-screen cells only. */
export function useHomeMemoryValues(deviceId: string | undefined, protocol: ProtocolType | undefined, enabled: boolean): MemoryValues {
  const [streamState, setStreamState] = useState<'idle' | 'open' | 'failed'>('idle');
  const active = enabled && Boolean(deviceId && protocol);
  const query = useQuery({
    queryKey: ['recording-home-memory-values', deviceId],
    queryFn: () => pointAPI.list({ device_id: deviceId }),
    enabled: active,
    refetchInterval: streamState === 'failed' ? FALLBACK_POLL_MS : false, refetchOnWindowFocus: false, retry: false, gcTime: 0,
  });
  const base = useMemo(() => query.data && deviceId && protocol ? indexMemoryValues(query.data, deviceId, protocol) : new Map<string, MemoryLiveValue>(),
    [query.data, deviceId, protocol]);
  const [visibleIds, setVisibleIds] = useState('');
  const timer = useRef<ReturnType<typeof setTimeout>>(undefined);
  const setVisible = useCallback((keys: readonly string[]) => {
    clearTimeout(timer.current);
    timer.current = setTimeout(() => {
      const ids = [...new Set(keys.map(key => base.get(key)?.pointId).filter((id): id is string => Boolean(id)))].sort().join(',');
      setVisibleIds(current => current === ids ? current : ids);
    }, VISIBLE_DEBOUNCE_MS);
  }, [base]);
  useEffect(() => () => clearTimeout(timer.current), []);

  const [live, setLive] = useState<ReadonlyMap<string, RuntimeValueEvent>>(new Map());
  const [liveAt, setLiveAt] = useState<number>();
  useEffect(() => {
    // An empty point_ids filter would subscribe the whole device, so nothing on screen means no stream.
    if (!active || !deviceId || !visibleIds || typeof EventSource === 'undefined') return undefined;
    const source = new EventSource(runtimeAPI.getStreamUrl(deviceId, visibleIds.split(',')));
    const pending = new Map<string, RuntimeValueEvent>();
    let flush: ReturnType<typeof setTimeout> | undefined;
    const onValue = (event: Event) => {
      const value = parseRuntimeValueRecord(parseMessageEventRecord(event));
      if (!value || value.device_id !== deviceId) return;
      pending.set(value.point_id, value);
      // Many points report in the same poll cycle: one render per batch, not per event.
      flush ??= setTimeout(() => {
        flush = undefined;
        const batch = new Map(pending); pending.clear();
        setLive(current => { const next = new Map(current); batch.forEach((item, id) => next.set(id, item)); return next; });
        setLiveAt(Date.now());
      }, FLUSH_MS);
    };
    source.onopen = () => setStreamState('open');
    // While the browser reconnects (or after the server refuses), the stored list is polled; a reopen stops it.
    source.onerror = () => setStreamState('failed');
    source.addEventListener('value', onValue);
    return () => { clearTimeout(flush); source.removeEventListener('value', onValue); source.close(); };
  }, [active, deviceId, visibleIds]);
  useEffect(() => { if (!active) { setLive(new Map()); setStreamState('idle'); } }, [active]);

  const byCell = useMemo(() => {
    if (!live.size) return base;
    const merged = new Map(base);
    merged.forEach((item, key) => { const event = item.pointId ? live.get(item.pointId) : undefined; if (event) merged.set(key, applyLiveEvent(item, event)); });
    return merged;
  }, [base, live]);
  const status: MemoryValuesStatus = query.isError ? 'error' : !query.data ? 'loading' : streamState === 'failed' ? 'polling' : 'live';
  return { byCell, status, updatedAt: Math.max(liveAt ?? 0, query.dataUpdatedAt || 0) || undefined, setVisible };
}
