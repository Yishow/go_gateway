// Runtime-only projections: do not apply these display rules to editable setup payloads.
const network = /(?:[a-z0-9-]+\.)+[a-z]{2,}|(?:[0-9]{1,3}\.){3}[0-9]{1,3}/i;
// Delimited absolute/relative prefixes are private syntax; UNC needs server and share.
const privatePath = /(?:^|[\s=:'"()[\]{},;\x60])(?:\/|(?:~|\.\.?)[\\/]|[a-z]:[\\/]|\\\\[^\\/\s]+[\\/][^\\/\s]+)/i;
const unsafeText = /(?:[a-z]+:\/\/|\bfile:|\b(?:password|passwd|pwd|uid|token|secret|dsn|username)\s*[=:]|\b(?:SELECT\s+.+\s+FROM|INSERT\s+INTO|UPDATE\s+\S+\s+SET|DELETE\s+FROM|(?:DROP|ALTER|CREATE)\s+TABLE|GRANT\s+.+\s+ON)\b)/i;
/** Reject recognizable private syntax, not ordinary scalar words or precise integer strings. */
export function isSafeRuntimeSensorText(value: string): boolean {
  return value.length <= 256 && !network.test(value) && !privatePath.test(value) && !unsafeText.test(value) && !/\p{Cc}/u.test(value);
}
/** IDs are opaque references, never transport endpoints. */
export function safeRuntimeIdentity(value: unknown): string {
  return typeof value === 'string' && /^[A-Za-z0-9_-][A-Za-z0-9_.-]{0,127}$/.test(value) && !network.test(value) ? value : '';
}
/** Human labels may contain Unicode, combining marks and harmless punctuation. */
export function safeRuntimeLabel(value: unknown): string {
  return typeof value === 'string' && /^[\p{L}\p{M}\p{N}_ .(),°%+·（）[\]［］-]{1,128}$/u.test(value) && isSafeRuntimeSensorText(value) ? value : '';
}
/** Register syntax includes existing decimal, FATEK and MC hex-addressed devices. */
export function safeRuntimeAddress(value: unknown): string {
  if (typeof value !== 'string') return '';
  const mc = /^(?:SW|SB|W|B)([0-9a-f]+)$/i.exec(value);
  if (mc) return mc[1].length <= 8 && Number.parseInt(mc[1], 16) <= 0x7fffffff ? value : '';
  return /^(?:[0-9]{1,6}|0x[0-9a-f]{1,8}|[a-z]{1,4}[0-9]{1,6}(?:\.[0-9]{1,2})?)$/i.test(value) ? value : '';
}
/** Semantic tag keys follow tag/service.go, not opaque identity or DNS syntax.
 * Dotted/slashed keys may resemble hostnames; only recognizable private syntax is rejected.
 */
export function safeRuntimeTagKey(value: unknown): string {
  return typeof value === 'string' && /^[a-zA-Z0-9_\-./]{1,128}$/.test(value) && !unsafeText.test(value) && !privatePath.test(value) && !value.split('/').some((segment) => segment === '.' || segment === '..') ? value : '';
}

/** Known runtime protocol names are shared by collector and workspace projections. */
export function safeRuntimeProtocol(value: unknown): string {
  return ['modbus_tcp', 'modbus_rtu', 'modbus_udp', 'fatek_fbs', 'mc_3e', 'mqtt'].find((item) => item === value) ?? 'unavailable';
}
/** Retain the deployed alignment vocabulary, including the legacy underscore spelling. */
export function safeRuntimeAlignment(value: unknown): string {
  return ['aligned', 'deferred', 'restart-required', 'restart_required', 'stale', 'degraded'].find((item) => item === value) ?? 'unavailable';
}
