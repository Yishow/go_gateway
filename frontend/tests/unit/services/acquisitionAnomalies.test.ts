import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  AcquisitionAnomalyError,
  acquisitionAnomalyErrorKey,
  activeAnomalyCount,
  getAcquisitionAnomalySummary,
  listDeviceAcquisitionAnomalies,
  parseAcquisitionAnomalyPage,
  parseAcquisitionAnomalySummary,
} from '../../../src/services/acquisitionAnomalies';
import { api } from '../../../src/services/datalinkClient';
import en from '../../../src/i18n/locales/en/recording-home.json';
import zhTW from '../../../src/i18n/locales/zh-TW/recording-home.json';
import type { AcquisitionAnomalyErrorCode } from '../../../src/types/acquisitionAnomaly';

vi.mock('../../../src/services/datalinkClient', () => ({ api: { get: vi.fn() } }));

const ended = {
  id: 7, device_id: 'mc-1', point_id: 'p-d5', point_name: 'Temp', point_address: 'D5', kind: 'stale', state: 'ended',
  started_at: '2026-11-05T01:00:00Z', ended_at: '2026-11-05T01:00:12.5Z', duration_ms: 12500, last_value: '42', detail: 'no-successful-read',
};
const active = {
  id: 8, device_id: 'mc-1', point_id: null, point_name: null, point_address: null, kind: 'disconnected', state: 'active',
  started_at: '2026-11-05T02:00:00Z', ended_at: null, duration_ms: null, last_value: null, detail: 'connection-failed; monitoring-stopped',
};
const page = (items: unknown[], extra: Record<string, unknown> = {}) => ({ success: true, data: { items, total: items.length, limit: 50, offset: 0, ...extra } });
const axiosError = (status: number, code: string, field?: string) => ({
  isAxiosError: true, message: 'Request failed: dial tcp 10.0.0.5:5432 password=secret',
  response: { status, data: { success: false, error: { code, message: 'internal: sql: database is locked at /var/lib/gw.db', ...(field ? { field } : {}) } } },
});

describe('acquisition anomaly list', () => {
  beforeEach(() => vi.clearAllMocks());

  it('sends only normalized filters and parses the typed page', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: page([active, ended], { total: 9, limit: 25, offset: 25 }) });
    const result = await listDeviceAcquisitionAnomalies('mc-1', { state: 'active', kind: 'stale', limit: 25, offset: 25 });
    expect(api.get).toHaveBeenCalledWith('/devices/mc-1/acquisition-anomalies', { params: { state: 'active', kind: 'stale', limit: 25, offset: 25 } });
    expect(result.total).toBe(9);
    expect(result.items[0]).toMatchObject({ id: 8, kind: 'disconnected', state: 'active', point_address: null, ended_at: null, duration_ms: null, detail_codes: ['connection-failed', 'monitoring-stopped'] });
    expect(result.items[1]).toMatchObject({ point_address: 'D5', last_value: '42', duration_ms: 12500, detail_codes: ['no-successful-read'] });
  });

  it('omits empty filters and encodes the device id path segment', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: page([], { limit: 50 }) });
    await listDeviceAcquisitionAnomalies('dev/../x', {});
    expect(api.get).toHaveBeenCalledWith('/devices/dev%2F..%2Fx/acquisition-anomalies', { params: { limit: 50, offset: 0 } });
  });

  it('rejects invalid filters locally instead of clamping them', async () => {
    await expect(listDeviceAcquisitionAnomalies('mc-1', { limit: 501 })).rejects.toMatchObject({ code: 'invalid_filter', field: 'limit' });
    await expect(listDeviceAcquisitionAnomalies('mc-1', { offset: -1 })).rejects.toMatchObject({ code: 'invalid_filter', field: 'offset' });
    await expect(listDeviceAcquisitionAnomalies('mc-1', { kind: 'bogus' as never })).rejects.toMatchObject({ code: 'invalid_filter', field: 'kind' });
    await expect(listDeviceAcquisitionAnomalies('', {})).rejects.toMatchObject({ code: 'invalid_filter', field: 'device_id' });
    expect(api.get).not.toHaveBeenCalled();
  });

  it('turns a failed read into a typed error, never an empty list, without raw backend text', async () => {
    vi.mocked(api.get).mockRejectedValue(axiosError(500, 'acquisition_anomaly_unavailable'));
    const failure = await listDeviceAcquisitionAnomalies('mc-1', {}).catch((error: unknown) => error);
    expect(failure).toBeInstanceOf(AcquisitionAnomalyError);
    expect(failure).toMatchObject({ code: 'unavailable', retryable: true });
    expect(JSON.stringify(failure)).not.toMatch(/10\.0\.0\.5|secret|database is locked|gw\.db/);
    expect((failure as Error).message).not.toMatch(/secret|locked/);
    expect(acquisitionAnomalyErrorKey(failure as AcquisitionAnomalyError)).toBe('recording-home:acquisitionAnomalies.errors.unavailable');
  });

  it.each([
    [400, 'acquisition_anomaly_filter_invalid', 'invalid_filter', false],
    [404, 'acquisition_anomaly_device_not_found', 'device_not_found', false],
    [503, 'acquisition_anomaly_capability_unavailable', 'capability_unavailable', false],
    [500, 'something_else', 'unknown', true],
    [502, 'acquisition_anomaly_unavailable', 'unknown', true],
  ])('classifies HTTP %i %s as %s', async (status, code, expected, retryable) => {
    vi.mocked(api.get).mockRejectedValue(axiosError(status, code, status === 400 ? 'kind' : undefined));
    await expect(listDeviceAcquisitionAnomalies('mc-1', {})).rejects.toMatchObject({ code: expected, retryable });
  });

  it('keeps an allowlisted rejected field only', async () => {
    vi.mocked(api.get).mockRejectedValueOnce(axiosError(400, 'acquisition_anomaly_filter_invalid', 'kind'));
    await expect(listDeviceAcquisitionAnomalies('mc-1', {})).rejects.toMatchObject({ field: 'kind' });
    vi.mocked(api.get).mockRejectedValueOnce(axiosError(400, 'acquisition_anomaly_filter_invalid', '<script>'));
    const failure = await listDeviceAcquisitionAnomalies('mc-1', {}).catch((error: unknown) => error) as AcquisitionAnomalyError;
    expect(failure.field).toBeUndefined();
  });

  it('treats a network failure without response as retryable unknown', async () => {
    vi.mocked(api.get).mockRejectedValue(new Error('Network Error'));
    await expect(listDeviceAcquisitionAnomalies('mc-1', {})).rejects.toMatchObject({ code: 'unknown', retryable: true });
  });

  it.each([
    ['missing data', { success: true }],
    ['success false', { success: false, data: page([]).data }],
    ['items not a list', { success: true, data: { items: null, total: 0, limit: 50, offset: 0 } }],
    ['negative total', page([], { total: -1 })],
    ['page larger than limit', page(Array.from({ length: 51 }, (_, i) => ({ ...ended, id: i + 1 })), { total: 51 })],
    ['echo limit mismatch', page([], { limit: 10 })],
    ['foreign device', page([{ ...ended, device_id: 'other' }])],
    ['unknown kind', page([{ ...ended, kind: 'offline' }])],
    ['ended without end time', page([{ ...ended, ended_at: null }])],
    ['active with duration', page([{ ...active, duration_ms: 5 }])],
    ['bad timestamp', page([{ ...ended, started_at: 'yesterday' }])],
    ['string id', page([{ ...ended, id: '7' }])],
    ['numeric last value', page([{ ...ended, last_value: 42 }])],
  ])('rejects a malformed page (%s) instead of showing it as empty', (_name, body) => {
    expect(() => parseAcquisitionAnomalyPage(body, 'mc-1', { limit: 50, offset: 0 })).toThrow(AcquisitionAnomalyError);
    try { parseAcquisitionAnomalyPage(body, 'mc-1', { limit: 50, offset: 0 }); } catch (error) {
      expect(error).toMatchObject({ code: 'malformed_response', retryable: true });
    }
  });

  it('maps unrecognised detail tokens to unknown rather than echoing them', () => {
    const result = parseAcquisitionAnomalyPage(page([{ ...ended, detail: 'no-successful-read; dial tcp 10.0.0.5' }]), 'mc-1', { limit: 50, offset: 0 });
    expect(result.items[0]?.detail_codes).toEqual(['no-successful-read', 'unknown']);
  });
});

describe('acquisition anomaly summary', () => {
  beforeEach(() => vi.clearAllMocks());
  const summary = { success: true, data: { items: [{ device_id: 'a', active_count: 0 }, { device_id: 'mc-1', active_count: 2 }], total_active: 2 } };

  it('reads per-device active counts and keeps zero distinct from unknown', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: summary });
    const result = await getAcquisitionAnomalySummary();
    expect(api.get).toHaveBeenCalledWith('/acquisition-anomalies/summary');
    expect(activeAnomalyCount(result, 'mc-1')).toBe(2);
    expect(activeAnomalyCount(result, 'a')).toBe(0);
    expect(activeAnomalyCount(result, 'missing')).toBeNull();
    expect(activeAnomalyCount(undefined, 'mc-1')).toBeNull();
  });

  it('fails instead of returning an empty summary', async () => {
    vi.mocked(api.get).mockRejectedValue(axiosError(503, 'acquisition_anomaly_capability_unavailable'));
    await expect(getAcquisitionAnomalySummary()).rejects.toMatchObject({ code: 'capability_unavailable' });
  });

  it.each([
    ['no items', { success: true, data: { total_active: 0 } }],
    ['total mismatch', { success: true, data: { ...summary.data, total_active: 3 } }],
    ['duplicate device', { success: true, data: { items: [{ device_id: 'a', active_count: 1 }, { device_id: 'a', active_count: 1 }], total_active: 2 } }],
    ['negative count', { success: true, data: { items: [{ device_id: 'a', active_count: -1 }], total_active: -1 } }],
    ['fractional count', { success: true, data: { items: [{ device_id: 'a', active_count: 1.5 }], total_active: 1.5 } }],
  ])('rejects a malformed summary (%s)', (_name, body) => {
    expect(() => parseAcquisitionAnomalySummary(body)).toThrow(expect.objectContaining({ code: 'malformed_response' }));
  });
});

describe('acquisition anomaly error copy', () => {
  const codes: AcquisitionAnomalyErrorCode[] = ['invalid_filter', 'device_not_found', 'capability_unavailable', 'unavailable', 'malformed_response', 'unknown'];
  it.each([['en', en], ['zh-TW', zhTW]])('%s has an actionable message for every error code', (_locale, copy) => {
    for (const code of codes) {
      const key = acquisitionAnomalyErrorKey(new AcquisitionAnomalyError(code, false)).replace('recording-home:', '').split('.');
      const message = key.reduce<unknown>((node, part) => (node as Record<string, unknown> | undefined)?.[part], copy);
      expect(typeof message === 'string' && message.length > 0, code).toBe(true);
    }
  });
});
