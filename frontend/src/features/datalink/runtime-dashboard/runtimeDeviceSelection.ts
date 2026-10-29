import { safeRuntimeIdentity } from '../../../utils/runtimeOperatorSafety';

export type RuntimeDeviceSelection =
  | { kind: 'absent' }
  | { kind: 'valid'; deviceId: string }
  | { kind: 'invalid' };

/** Parse the raw query before URLSearchParams can replace invalid encodings. */
export function parseRuntimeDeviceSelection(search: string): RuntimeDeviceSelection {
  const raw = search.startsWith('?') ? search.slice(1) : search;
  try {
    for (const pair of raw.split('&')) {
      const equals = pair.indexOf('=');
      const key = equals < 0 ? pair : pair.slice(0, equals);
      const value = equals < 0 ? '' : pair.slice(equals + 1);
      decodeURIComponent(key.replace(/\+/g, ' '));
      decodeURIComponent(value.replace(/\+/g, ' '));
    }
  } catch { return { kind: 'invalid' }; }
  const values = new URLSearchParams(raw).getAll('device_id');
  if (values.length === 0) return { kind: 'absent' };
  if (values.length !== 1 || !values[0] || safeRuntimeIdentity(values[0]) !== values[0]) return { kind: 'invalid' };
  return { kind: 'valid', deviceId: values[0] };
}
