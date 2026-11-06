/** Acquisition anomaly kinds recorded by the runtime detector. */
export const ACQUISITION_ANOMALY_KINDS = ['disconnected', 'stale', 'frozen'] as const;
export type AcquisitionAnomalyKind = (typeof ACQUISITION_ANOMALY_KINDS)[number];

/** Active events have no end time yet; ended events carry end time and duration. */
export const ACQUISITION_ANOMALY_STATES = ['active', 'ended'] as const;
export type AcquisitionAnomalyState = (typeof ACQUISITION_ANOMALY_STATES)[number];

/** Stable reason codes stored by the backend; anything else is shown as `unknown`. */
export const ACQUISITION_ANOMALY_DETAIL_CODES = [
  'connection-failed', 'circuit-open', 'breaker-open', 'no-successful-read',
  'value-unchanged', 'monitoring-stopped', 'frozen-detection-disabled',
] as const;
export type AcquisitionAnomalyDetailCode = (typeof ACQUISITION_ANOMALY_DETAIL_CODES)[number] | 'unknown';

/** One persisted event. `point_*` are null for device-level events or deleted points. */
export interface AcquisitionAnomalyItem {
  id: number;
  device_id: string;
  point_id: string | null;
  point_name: string | null;
  point_address: string | null;
  kind: AcquisitionAnomalyKind;
  state: AcquisitionAnomalyState;
  started_at: string;
  ended_at: string | null;
  duration_ms: number | null;
  last_value: string | null;
  detail_codes: AcquisitionAnomalyDetailCode[];
}

export interface AcquisitionAnomalyPage {
  items: AcquisitionAnomalyItem[];
  total: number;
  limit: number;
  offset: number;
}

/** Query filters; omitted fields use backend defaults (limit 50, offset 0). */
export interface AcquisitionAnomalyFilter {
  state?: AcquisitionAnomalyState;
  kind?: AcquisitionAnomalyKind;
  limit?: number;
  offset?: number;
}

export interface AcquisitionAnomalyDeviceCount {
  device_id: string;
  active_count: number;
}

export interface AcquisitionAnomalySummary {
  items: AcquisitionAnomalyDeviceCount[];
  total_active: number;
}

/** Allowlisted failure classes; raw backend text is never retained. */
export type AcquisitionAnomalyErrorCode =
  | 'invalid_filter'
  | 'device_not_found'
  | 'capability_unavailable'
  | 'unavailable'
  | 'malformed_response'
  | 'unknown';

export type AcquisitionAnomalyFilterField = 'state' | 'kind' | 'limit' | 'offset' | 'device_id';
