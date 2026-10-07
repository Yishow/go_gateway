import { describe, expect, it } from 'vitest';
import { appendRuntimeLogs, runtimeLogBytes } from '../../../src/features/runtime-logs/runtimeLogBuffer';
import { record } from './fixtures';

describe('runtime log browser limits', () => {
  it('evicts oldest rows at the independent 1000-record limit', () => {
    const result = appendRuntimeLogs([], Array.from({ length: 1100 }, (_, i) => record(String(i + 1))));
    expect(result.records).toHaveLength(1000);
    expect(result.records[0].sequence).toBe('101');
    expect(result.truncated).toBe(true);
  });
  it('enforces 2 MiB UTF-8 bytes independently of row count', () => {
    const result = appendRuntimeLogs([], Array.from({ length: 700 }, (_, i) => record(String(i + 1), '界'.repeat(2000))));
    expect(result.records.length).toBeLessThan(700);
    expect(result.records.reduce((n, r) => n + runtimeLogBytes(r), 0)).toBeLessThanOrEqual(2 * 1024 * 1024);
    expect(result.records.at(-1)?.sequence).toBe('700');
  });
  it('deduplicates without coercing 64-bit string identities to numbers', () => {
    const a = record('9007199254740992'); const b = record('9007199254740993');
    expect(appendRuntimeLogs([a], [a, b]).records).toEqual([a, b]);
  });
});
