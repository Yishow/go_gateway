import { describe, expect, it } from 'vitest';
import {
  MAX_SAFE_JSON_BYTES,
  parseBoundedJson,
  parseMessageEventRecord,
  parseRuntimeValueRecord,
  normalizeTypedEnvelope,
} from '../../../src/utils/safeJson';

describe('bounded external JSON parser', () => {
  it('rejects oversized payloads before parsing', () => {
    expect(parseBoundedJson('x'.repeat(MAX_SAFE_JSON_BYTES + 1))).toBeNull();
  });

  it('rejects deep, oversized-array, and non-finite JSON values', () => {
    expect(parseBoundedJson('{"a":{"b":{"c":{"d":{"e":{"f":{"g":1}}}}}}}')).toBeNull();
    expect(parseBoundedJson(JSON.stringify(Array.from({ length: 257 }, () => 1)))).toBeNull();
    expect(parseBoundedJson('{"value":1e999}')).toBeNull();
  });

  it('rebuilds only bounded typed envelope fields', () => {
    const event = parseMessageEventRecord({
      data: JSON.stringify({
        code: 'unknown-code',
        action: 'raw backend action',
        request_id: 'req-1',
        retryable: true,
        secret: { should_not: 'survive' },
      }),
    } as MessageEvent<string>);

    expect(normalizeTypedEnvelope(event)).toEqual({
      code: undefined,
      action: undefined,
      requestId: 'req-1',
      retryable: true,
    });
  });

  it('rejects overlong typed envelope fields and retains no raw object', () => {
    const event = parseMessageEventRecord({
      data: JSON.stringify({
        code: 'runtime_stream_unavailable',
        action: 'x'.repeat(129),
        request_id: 'r'.repeat(129),
        retryable: false,
      }),
    } as MessageEvent<string>);

    expect(normalizeTypedEnvelope(event)).toEqual({
      code: 'runtime_stream_unavailable',
      action: undefined,
      requestId: undefined,
      retryable: false,
    });
  });

  it('fails closed for long runtime fields and non-finite runtime values', () => {
    expect(parseRuntimeValueRecord({
      device_id: 'device-1', point_id: 'point-1', address: '40001',
      raw_value: 'x'.repeat(257), transformed_value: 1, quality: 'good', stale: false,
      timestamp: '2026-08-25T00:00:00Z',
    })).toBeNull();
    expect(parseRuntimeValueRecord({
      device_id: 'device-1', point_id: 'point-1', address: '40001',
      raw_value: Infinity, transformed_value: 1, quality: 'good', stale: false,
      timestamp: '2026-08-25T00:00:00Z',
    })).toBeNull();
  });
});
