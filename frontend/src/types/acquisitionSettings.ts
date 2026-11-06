/** Saved device-level acquisition settings as read back from the Studio v2 device API. */
export interface DeviceAcquisitionSettings {
  /** Interval owned by the device's polling group; null when the backend reported none. */
  interval_ms: number | null;
  auto_interval: boolean;
  /** null means the backend default (interval x5, at least 10 seconds). */
  stale_threshold_ms: number | null;
  /** 0 disables frozen detection. */
  frozen_threshold_ms: number;
}

/** Raw operator text; parsed only when the existing submit adapter runs. */
export interface DeviceAcquisitionInput {
  interval_ms: string;
  auto_interval: boolean;
  stale_threshold_ms: string;
  frozen_threshold_ms: string;
  /** The backend cannot clear a saved stale threshold back to its default. */
  stale_persisted: boolean;
}

export type AcquisitionProfileId =
  | 'safe-v1'
  | 'mitsubishi-fx5-cpu-3e-binary-v1'
  | 'mitsubishi-q-iqr-3e-binary-v1'
  | 'fatek-fbs-ascii-v1';

/** UI choice: `undeclared` sends no `connection_config.acquisition` at all. */
export type AcquisitionProfileChoice = 'undeclared' | AcquisitionProfileId;

export type AcquisitionField = 'interval' | 'stale' | 'frozen' | 'profile' | 'randomPoints';
export type AcquisitionIssue = 'intervalRange' | 'staleInvalid' | 'staleRequired' | 'frozenInvalid' | 'randomRange' | 'fx5Route';

/** Runtime collector evidence; a configured value never stands in for an effective one. */
export type RuntimeIntervalReadout =
  | { status: 'loading' }
  | { status: 'unavailable' }
  | { status: 'not-running' }
  | { status: 'ready'; configured_ms: number; effective_ms: number };

/** Studio v2 device response fields; absent on backends without acquisition settings. */
export interface DeviceAcquisitionRecordFields {
  acquisition_interval_ms?: number;
  acquisition_auto_interval?: boolean;
  stale_threshold_ms?: number | null;
  frozen_threshold_ms?: number;
}

/** Device update fields; an omitted stale threshold keeps the saved value (it cannot be cleared). */
export interface DeviceAcquisitionUpdateFields {
  acquisition_interval_ms?: number;
  acquisition_auto_interval?: boolean;
  stale_threshold_ms?: number;
  frozen_threshold_ms?: number;
}

/** Runtime collector intervals from GET /runtime/status. */
export interface RuntimeCollectorIntervals {
  configured_interval_ms?: number;
  effective_interval_ms?: number;
}
