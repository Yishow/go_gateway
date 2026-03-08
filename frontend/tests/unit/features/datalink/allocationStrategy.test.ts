import { describe, expect, it } from 'vitest';
import { findNearestValidContiguousSpan } from '@/features/datalink/allocationStrategy';

describe('findNearestValidContiguousSpan', () => {
  it('returns the same start when span is immediately available', () => {
    const result = findNearestValidContiguousSpan({
      startAddress: '40001',
      spanSize: 2,
      protocol: 'modbus_tcp',
      blockedAddresses: new Set(['40010']),
    });

    expect(result).toEqual({ startAddress: '40001', reason: 'ok' });
  });

  it('skips blocked cells and returns nearest contiguous span', () => {
    const result = findNearestValidContiguousSpan({
      startAddress: '40001',
      spanSize: 3,
      protocol: 'modbus_tcp',
      blockedAddresses: new Set(['40001', '40002', '40004']),
    });

    expect(result).toEqual({ startAddress: '40005', reason: 'ok' });
  });

  it('returns no_contiguous_span when no valid span exists in window', () => {
    const blocked = new Set(
      Array.from({ length: 10 }).map((_, index) => String(40001 + index))
    );
    const result = findNearestValidContiguousSpan({
      startAddress: '40001',
      spanSize: 2,
      protocol: 'modbus_tcp',
      blockedAddresses: blocked,
      candidateWindow: 10,
    });

    expect(result).toEqual({
      startAddress: null,
      reason: 'no_contiguous_span',
    });
  });

  it('returns invalid_start for malformed start address', () => {
    const result = findNearestValidContiguousSpan({
      startAddress: 'INVALID',
      spanSize: 2,
      protocol: 'modbus_tcp',
      blockedAddresses: new Set(),
    });

    expect(result).toEqual({ startAddress: null, reason: 'invalid_start' });
  });
});
