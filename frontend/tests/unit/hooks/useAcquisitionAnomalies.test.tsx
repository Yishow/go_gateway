import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { acquisitionAnomalyKeys } from '@/hooks/datalink/keys';
import { useAcquisitionAnomalySummaryQuery, useDeviceAcquisitionAnomaliesQuery } from '@/hooks/datalink/useAcquisitionAnomalies';
import { AcquisitionAnomalyError, getAcquisitionAnomalySummary, listDeviceAcquisitionAnomalies } from '@/services/acquisitionAnomalies';

vi.mock('@/services/acquisitionAnomalies', async () => {
  const actual = await vi.importActual<typeof import('@/services/acquisitionAnomalies')>('@/services/acquisitionAnomalies');
  return { ...actual, getAcquisitionAnomalySummary: vi.fn(), listDeviceAcquisitionAnomalies: vi.fn() };
});

function wrapper() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: 3 } } });
  return { client, Wrapper: ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider> };
}

describe('acquisition anomaly hooks', () => {
  beforeEach(() => vi.clearAllMocks());

  it('surfaces a summary failure as an error, not as zero anomalies', async () => {
    vi.mocked(getAcquisitionAnomalySummary).mockRejectedValue(new AcquisitionAnomalyError('unavailable', true));
    const { Wrapper } = wrapper();
    const { result } = renderHook(() => useAcquisitionAnomalySummaryQuery(), { wrapper: Wrapper });
    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.data).toBeUndefined();
    expect(result.current.error).toMatchObject({ code: 'unavailable' });
    expect(getAcquisitionAnomalySummary).toHaveBeenCalledTimes(1);
  });

  it('keys the device list by device and filter and does not reuse another filter page', async () => {
    vi.mocked(listDeviceAcquisitionAnomalies).mockResolvedValue({ items: [], total: 0, limit: 50, offset: 0 });
    const { Wrapper, client } = wrapper();
    const { result, rerender } = renderHook(({ kind }: { kind?: 'stale' | 'frozen' }) => useDeviceAcquisitionAnomaliesQuery('mc-1', { kind }), { wrapper: Wrapper, initialProps: { kind: 'stale' } });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(listDeviceAcquisitionAnomalies).toHaveBeenCalledWith('mc-1', { kind: 'stale' });
    expect(client.getQueryData(acquisitionAnomalyKeys.deviceList('mc-1', { kind: 'stale' }))).toBeDefined();
    vi.mocked(listDeviceAcquisitionAnomalies).mockReturnValue(new Promise(() => undefined));
    rerender({ kind: 'frozen' });
    await waitFor(() => expect(result.current.isPending).toBe(true));
    expect(result.current.data).toBeUndefined();
  });

  it('does not read while disabled or without a device', () => {
    const { Wrapper } = wrapper();
    renderHook(() => useDeviceAcquisitionAnomaliesQuery('mc-1', {}, false), { wrapper: Wrapper });
    renderHook(() => useDeviceAcquisitionAnomaliesQuery('', {}), { wrapper: Wrapper });
    expect(listDeviceAcquisitionAnomalies).not.toHaveBeenCalled();
  });

  it('device list failure stays an error and is not retried silently', async () => {
    vi.mocked(listDeviceAcquisitionAnomalies).mockRejectedValue(new AcquisitionAnomalyError('unavailable', true));
    const { Wrapper } = wrapper();
    const { result } = renderHook(() => useDeviceAcquisitionAnomaliesQuery('mc-1', {}), { wrapper: Wrapper });
    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.data).toBeUndefined();
    expect(listDeviceAcquisitionAnomalies).toHaveBeenCalledTimes(1);
  });
});
