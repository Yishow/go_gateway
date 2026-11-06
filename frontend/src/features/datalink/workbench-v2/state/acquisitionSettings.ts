import type { DeviceSourceBuffer } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import type { StudioV2WorkspaceDeviceRecord } from '../../../../services/studioV2WorkspaceDevices';
import type { RuntimeCollectorStatus } from '../../../../types/datalink';
import type {
  AcquisitionField, AcquisitionIssue, AcquisitionProfileChoice, DeviceAcquisitionInput, DeviceAcquisitionSettings, DeviceAcquisitionUpdateFields,
  RuntimeCollectorIntervals, RuntimeIntervalReadout,
} from '../../../../types/acquisitionSettings';
import type { ProtocolId } from './types';

/** Mirrors internal/datalink/device/service_crud.go Min/MaxAcquisitionIntervalMs. */
export const ACQUISITION_INTERVAL_MIN_MS = 100;
export const ACQUISITION_INTERVAL_MAX_MS = 3_600_000;
const MC_PROFILES = ['undeclared', 'safe-v1', 'mitsubishi-fx5-cpu-3e-binary-v1', 'mitsubishi-q-iqr-3e-binary-v1'] as const;
const FATEK_PROFILES = ['undeclared', 'safe-v1', 'fatek-fbs-ascii-v1'] as const;

/** Offered choices; null for protocols whose connector aggregates contiguous blocks only. */
export function acquisitionProfileOptions(protocol: ProtocolId): readonly AcquisitionProfileChoice[] | null {
  if (protocol === 'mc_3e') return MC_PROFILES;
  if (protocol === 'fatek_fbs') return FATEK_PROFILES;
  return null;
}

/** Mirrors internal/datalink/connector/acquisition_profile.go randomProfileLimit; 0 = no random reads. */
export function randomProfileLimit(profile: string): number {
  if (profile === 'mitsubishi-fx5-cpu-3e-binary-v1' || profile === 'mitsubishi-q-iqr-3e-binary-v1') return 192;
  return profile === 'fatek-fbs-ascii-v1' ? 64 : 0;
}

const objectOf = (value: unknown): Record<string, unknown> | undefined =>
  value && typeof value === 'object' && !Array.isArray(value) ? value as Record<string, unknown> : undefined;
const wholeNumber = (text: string): number | undefined => {
  const trimmed = text.trim();
  if (!/^\d+$/.test(trimmed)) return undefined;
  const value = Number(trimmed);
  return Number.isSafeInteger(value) ? value : undefined;
};

/** The saved profile string (possibly unrecognized) or `undeclared`, plus a lowered random limit. */
export function readAcquisitionDeclaration(config: Record<string, unknown>): { profile: string; randomPoints?: number } {
  const declaration = objectOf(config.acquisition);
  if (!declaration) return { profile: 'undeclared' };
  const profile = typeof declaration.profile === 'string' ? declaration.profile : '';
  const randomPoints = objectOf(declaration.limits)?.random_points;
  return { profile, ...(typeof randomPoints === 'number' ? { randomPoints } : {}) };
}

/** Writes only `connection_config.acquisition`; other declaration fields are kept. */
export function withAcquisitionProfile(config: Record<string, unknown>, profile: string): Record<string, unknown> {
  const next = { ...config };
  if (profile === 'undeclared') { delete next.acquisition; return next; }
  const limits = { ...objectOf(objectOf(config.acquisition)?.limits) };
  delete limits.random_points;
  const declaration: Record<string, unknown> = { ...objectOf(config.acquisition), version: 1, profile, limits };
  if (!Object.keys(limits).length) delete declaration.limits;
  next.acquisition = declaration;
  return next;
}

/** Text shown in the random limit field: operator text, else the saved limit, else the profile limit. */
export function randomPointsText(buffer: DeviceSourceBuffer): string {
  if (buffer.randomPointsInput !== undefined) return buffer.randomPointsInput;
  const declared = readAcquisitionDeclaration(buffer.config);
  return String(declared.randomPoints ?? randomProfileLimit(declared.profile));
}

/** Mirrors ownFX5CPURoute: the FX5 CPU profile only applies to the own-station route. */
function ownFX5Route(config: Record<string, unknown>): boolean {
  const num = (value: unknown, fallback: number) => value === undefined || value === '' ? fallback : Number(value);
  const station = num(config.station_no ?? config.station ?? config.module_station_no, 0);
  const io = num(config.io_no ?? config.module_io_no, 0), pc = num(config.pc_no, 0);
  return num(config.network_no, 0) === 0 && (pc === 0 || pc === 255) && (io === 0 || io === 0x03ff) && station === 0;
}

/** One source of field errors for the form and for the submit adapter. */
export function acquisitionIssues(buffer: DeviceSourceBuffer): Partial<Record<AcquisitionField, AcquisitionIssue>> {
  const issues: Partial<Record<AcquisitionField, AcquisitionIssue>> = {};
  const input = buffer.acquisitionInput;
  if (input) {
    const interval = wholeNumber(input.interval_ms);
    if (interval === undefined || interval < ACQUISITION_INTERVAL_MIN_MS || interval > ACQUISITION_INTERVAL_MAX_MS) issues.interval = 'intervalRange';
    const stale = input.stale_threshold_ms.trim();
    if (stale === '') { if (input.stale_persisted) issues.stale = 'staleRequired'; } else if (!(wholeNumber(stale)! > 0)) issues.stale = 'staleInvalid';
    if (wholeNumber(input.frozen_threshold_ms) === undefined) issues.frozen = 'frozenInvalid';
  }
  if (acquisitionProfileOptions(buffer.protocol)) {
    const { profile } = readAcquisitionDeclaration(buffer.config);
    if (profile === 'mitsubishi-fx5-cpu-3e-binary-v1' && !ownFX5Route(buffer.config)) issues.profile = 'fx5Route';
    const limit = randomProfileLimit(profile), text = randomPointsText(buffer).trim(), points = wholeNumber(text);
    if (limit > 0 && text !== '' && (points === undefined || points < 1 || points > limit)) issues.randomPoints = 'randomRange';
  }
  return issues;
}

/** Parsed config/settings, or undefined while any acquisition field is invalid. */
export function parseAcquisition(buffer: DeviceSourceBuffer): { config: Record<string, unknown>; acquisition?: DeviceAcquisitionSettings } | undefined {
  if (Object.keys(acquisitionIssues(buffer)).length) return undefined;
  let config = buffer.config;
  const { profile } = readAcquisitionDeclaration(config), limit = randomProfileLimit(profile);
  const offered: readonly string[] = acquisitionProfileOptions(buffer.protocol) ?? [];
  if (profile !== 'undeclared' && offered.includes(profile)) {
    config = withAcquisitionProfile(config, profile);
    const points = wholeNumber(randomPointsText(buffer));
    if (limit > 0 && points !== undefined && points < limit) {
      const declaration = config.acquisition as Record<string, unknown>;
      config = { ...config, acquisition: { ...declaration, limits: { ...objectOf(declaration.limits), random_points: points } } };
    }
  }
  const input = buffer.acquisitionInput;
  if (!input) return { config };
  const stale = input.stale_threshold_ms.trim();
  return { config, acquisition: { interval_ms: wholeNumber(input.interval_ms)!, auto_interval: input.auto_interval,
    stale_threshold_ms: stale === '' ? null : wholeNumber(stale)!, frozen_threshold_ms: wholeNumber(input.frozen_threshold_ms)! } };
}

export function acquisitionInputFromSettings(settings: DeviceAcquisitionSettings): DeviceAcquisitionInput {
  return { interval_ms: settings.interval_ms === null ? '' : String(settings.interval_ms), auto_interval: settings.auto_interval,
    stale_threshold_ms: settings.stale_threshold_ms === null ? '' : String(settings.stale_threshold_ms),
    frozen_threshold_ms: String(settings.frozen_threshold_ms), stale_persisted: settings.stale_threshold_ms !== null };
}

const positiveInt = (value: unknown): value is number => typeof value === 'number' && Number.isSafeInteger(value) && value > 0;
/** Fails closed: a backend without these fields, or with malformed ones, offers no editable settings. */
export function acquisitionSettingsFromRecord(record: StudioV2WorkspaceDeviceRecord): DeviceAcquisitionSettings | undefined {
  const { acquisition_interval_ms: interval, acquisition_auto_interval: auto, stale_threshold_ms: stale, frozen_threshold_ms: frozen } = record;
  if (typeof auto !== 'boolean' || typeof frozen !== 'number' || !Number.isSafeInteger(frozen) || frozen < 0) return undefined;
  if (interval != null && !positiveInt(interval) || stale != null && !positiveInt(stale)) return undefined;
  return { interval_ms: interval ?? null, auto_interval: auto, stale_threshold_ms: stale ?? null, frozen_threshold_ms: frozen };
}

/** Update fields; an unknown interval or a default stale threshold is omitted, never sent as a guess. */
export function acquisitionUpdateFields(settings: DeviceAcquisitionSettings | undefined): DeviceAcquisitionUpdateFields {
  if (!settings) return {};
  return { ...(settings.interval_ms !== null ? { acquisition_interval_ms: settings.interval_ms } : {}), acquisition_auto_interval: settings.auto_interval,
    ...(settings.stale_threshold_ms !== null ? { stale_threshold_ms: settings.stale_threshold_ms } : {}), frozen_threshold_ms: settings.frozen_threshold_ms };
}

/** Readback proof for a submitted save; absent settings were not part of the request. */
export function acquisitionMatches(record: StudioV2WorkspaceDeviceRecord, settings: DeviceAcquisitionSettings | undefined): boolean {
  if (!settings) return true;
  return (settings.interval_ms === null || record.acquisition_interval_ms === settings.interval_ms) && record.acquisition_auto_interval === settings.auto_interval &&
    (record.stale_threshold_ms ?? null) === settings.stale_threshold_ms && record.frozen_threshold_ms === settings.frozen_threshold_ms;
}

type RuntimeIntervalSource = { collectors: readonly (RuntimeCollectorStatus & RuntimeCollectorIntervals)[] };
export function runtimeIntervalReadout(deviceId: string, query: { data?: RuntimeIntervalSource; isError: boolean; isPending: boolean }): RuntimeIntervalReadout {
  if (query.isError) return { status: 'unavailable' };
  if (query.isPending || !query.data) return { status: 'loading' };
  const collector = query.data.collectors.find(item => item.device_id === deviceId);
  if (!collector?.running || !positiveInt(collector.configured_interval_ms) || !positiveInt(collector.effective_interval_ms)) return { status: 'not-running' };
  return { status: 'ready', configured_ms: collector.configured_interval_ms, effective_ms: collector.effective_interval_ms };
}
