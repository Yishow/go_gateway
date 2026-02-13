import type { ProtocolType } from '../../types/datalink';
import { addressParser } from '../../utils/addressParser';

export interface FindNearestValidSpanInput {
  startAddress: string;
  spanSize: number;
  protocol: ProtocolType;
  blockedAddresses: Set<string>;
  candidateWindow?: number;
}

export interface FindNearestValidSpanResult {
  startAddress: string | null;
  reason: 'ok' | 'invalid_start' | 'no_contiguous_span';
}

export function findNearestValidContiguousSpan(
  input: FindNearestValidSpanInput
): FindNearestValidSpanResult {
  const {
    startAddress,
    spanSize,
    protocol,
    blockedAddresses,
    candidateWindow = 800,
  } = input;

  if (!startAddress || spanSize <= 0) {
    return { startAddress: null, reason: 'invalid_start' };
  }

  const candidates = addressParser.expand(startAddress, candidateWindow, protocol);
  if (candidates.length === 0) {
    return { startAddress: null, reason: 'invalid_start' };
  }

  const matchedStart = candidates.find((candidateStart) => {
    const span = addressParser.expand(candidateStart, spanSize, protocol);
    if (span.length !== spanSize) return false;
    return span.every((address) => !blockedAddresses.has(address));
  });

  if (!matchedStart) {
    return { startAddress: null, reason: 'no_contiguous_span' };
  }

  return { startAddress: matchedStart, reason: 'ok' };
}
