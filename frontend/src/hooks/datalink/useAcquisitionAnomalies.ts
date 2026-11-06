import { useQuery } from '@tanstack/react-query';
import { getAcquisitionAnomalySummary, listDeviceAcquisitionAnomalies } from '../../services/acquisitionAnomalies';
import type { AcquisitionAnomalyFilter } from '../../types/acquisitionAnomaly';
import { acquisitionAnomalyKeys } from './keys';

/** Summary refresh period (design condition): the home buttons follow new or ended events without a reload. */
export const ACQUISITION_ANOMALY_SUMMARY_REFRESH_MS = 30_000;

/**
 * Per-device active anomaly counts. A failed read stays an error (status unknown)
 * and is never replaced by zero counts; retries are explicit.
 */
export function useAcquisitionAnomalySummaryQuery(enabled = true) {
  return useQuery({
    queryKey: acquisitionAnomalyKeys.summary(),
    queryFn: getAcquisitionAnomalySummary,
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
    refetchInterval: enabled ? ACQUISITION_ANOMALY_SUMMARY_REFRESH_MS : false,
  });
}

/**
 * One device's anomaly page for the given filter. Results from another filter
 * are never shown as placeholders, and failures stay errors with an explicit retry.
 */
export function useDeviceAcquisitionAnomaliesQuery(deviceId: string, filter: AcquisitionAnomalyFilter, enabled = true) {
  return useQuery({
    queryKey: acquisitionAnomalyKeys.deviceList(deviceId, filter),
    queryFn: () => listDeviceAcquisitionAnomalies(deviceId, filter),
    enabled: enabled && deviceId !== '',
    retry: false,
    refetchOnWindowFocus: false,
  });
}
