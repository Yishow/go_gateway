import type { ProtocolId, Rule } from '../state/types';
import { addressParser } from '../../../../utils/addressParser';
import { dataTypeWidth, getRuleReadinessReason } from '../state/sourceRule';
import { isStudioV2RuleValid } from '../state/studioV2RuleAutosave';
export interface RecordingRangeBuffer { name: string; address: string; count: string; dataType: Rule['data_type']; scale: string; offset: string }
export const recordingRangeDataTypes: readonly Rule['data_type'][] = ['bool', 'int16', 'int32', 'int64', 'uint16', 'uint32', 'uint64', 'float32', 'float64', 'string'];
export function newRecordingRangeBuffer(): RecordingRangeBuffer { return { name: '', address: '', count: '', dataType: 'int16', scale: '1', offset: '0' }; }
// SourceRule service normalizes naming prefixes to ASCII upper-case tokens.
// Normalize the submitted value only; the raw local name remains untouched.
function namingPrefix(text: string): string {
  let result = '';
  for (const raw of text.trim()) {
    const upper = raw.toUpperCase();
    const char = upper.length === 1 ? upper : raw;
    if (/^[A-Z0-9_-]$/.test(char)) result += char;
    else if (!result.endsWith('-')) result += '-';
  }
  return result.replace(/^-+|-+$/g, '') || 'SRC';
}
function numberText(text: string): number | undefined {
  if (!/^[+-]?(?:\d+(?:\.\d*)?|\.\d+)$/.test(text.trim())) return;
  const value = Number(text); return Number.isFinite(value) ? value : undefined;
}
/** Raw text stays untouched. Existing protocol parser/geometry and rule validator
 * are applied only for explicit submission; no PLC capability is inferred. */
export function parseRecordingRangeBuffer(buffer: RecordingRangeBuffer, deviceId: string, protocol: ProtocolId): { rule: Rule; span: number; end: string } | undefined {
  if (!/^\d+$/.test(buffer.count.trim()) || !buffer.name.trim() || !recordingRangeDataTypes.includes(buffer.dataType)) return;
  const count = Number(buffer.count), scale = numberText(buffer.scale), offset = numberText(buffer.offset);
  if (!Number.isSafeInteger(count) || count <= 0 || scale === undefined || offset === undefined) return;
  const rule: Rule = { id: '', device_id: deviceId, name: buffer.name, naming_prefix: namingPrefix(buffer.name), start_address: buffer.address.trim().toUpperCase(),
    count, data_type: buffer.dataType, enabled: true, scale_multiplier: scale, scale_offset: offset, data_format: '',
    skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null };
  if (!isStudioV2RuleValid(rule) || getRuleReadinessReason(rule, { [deviceId]: protocol })) return;
  const span = count * dataTypeWidth(buffer.dataType);
  if (!Number.isSafeInteger(span)) return;
  try { return { rule, span, end: addressParser.offset(buffer.address, span - 1, protocol) }; } catch { return; }
}
