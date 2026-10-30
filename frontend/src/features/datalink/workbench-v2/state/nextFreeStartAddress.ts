import { addressParser } from '../../../../utils/addressParser';
import type { ProtocolType } from '../../../../types/datalink';
import type { Rule } from './types';

const REGISTERS_PER_VALUE: Record<Rule['data_type'], number> = {
  bool: 1, int16: 1, uint16: 1, int32: 2, uint32: 2, float32: 2, int64: 4, uint64: 4, float64: 4, string: 1,
};

function nextMCAddress(deviceRules: Rule[], fallback: string): string {
  let best: { offset: number; address: string } | null = null;
  try {
    for (const rule of deviceRules) {
      const width = REGISTERS_PER_VALUE[rule.data_type] ?? 1;
      const identities = rule.mc_persisted_points?.points;
      if (rule.persisted && !identities?.length) return '';
      const spans = rule.persisted && identities?.length
        ? identities.map((point) => ({ address: point.address, width }))
        : [{ address: rule.start_address, width: Math.max(1, rule.count) * width }];
      for (const span of spans) {
        const address = addressParser.offset(span.address, span.width, 'mc_3e');
        const offset = addressParser.parse(address, 'mc_3e').startNumber;
        if (!best || offset > best.offset) best = { offset, address };
      }
    }
  } catch {
    return deviceRules.some((rule) => rule.persisted) ? '' : fallback;
  }
  return best?.address ?? fallback;
}

/** Chooses the first free address using MC radix and persisted point identities. */
export function nextFreeStartAddress(deviceRules: Rule[], protocol: ProtocolType, fallback: string): string {
  if (protocol === 'mc_3e') return nextMCAddress(deviceRules, fallback);
  let best: { end: number; prefix: string; width: number } | null = null;
  for (const rule of deviceRules) {
    const text = rule.start_address.trim().toUpperCase();
    const span = Math.max(1, rule.count) * (REGISTERS_PER_VALUE[rule.data_type] ?? 1);
    const modbus = /^([0134])(\d{4,5})$/.exec(text);
    const lettered = /^([A-Z]+)(\d+)$/.exec(text);
    const match = protocol.startsWith('modbus') ? modbus : lettered;
    if (!match) return fallback;
    const end = parseInt(match[2], 10) + span;
    if (!best || end > best.end) best = { end, prefix: match[1], width: match[2].length };
  }
  if (!best) return fallback;
  return `${best.prefix}${String(best.end).padStart(best.width, '0')}`;
}
