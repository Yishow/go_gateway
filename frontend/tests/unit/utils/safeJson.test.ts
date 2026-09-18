import { describe, expect, it } from 'vitest';
import {
  MAX_SAFE_JSON_BYTES,
  parseBoundedJson,
  parseMessageEventRecord,
  parseRuntimeValueRecord,
  normalizeTypedEnvelope,
  parseMappingPreviewResponse,
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

  it('keeps only the exact recording operation codes and wait action', () => {
    for (const code of ['RECORDING_SCHEMA_NOT_IMPLEMENTED', 'RECORDING_TEST_WRITE_NOT_IMPLEMENTED']) {
      expect(normalizeTypedEnvelope({
        response: {
          data: {
            success: false,
            error: {
              code,
              action: 'wait_for_supported_operation',
              retryable: false,
              request_id: 'req-recording-full',
              message: 'raw backend diagnostic',
            },
          },
        },
      })).toEqual({
        code,
        action: 'wait_for_supported_operation',
        requestId: 'req-recording-full',
        retryable: false,
      });
    }

    expect(normalizeTypedEnvelope({
      code: 'RECORDING_TEST_WRITE_NOT_IMPLEMENTED_EXTRA',
      action: 'wait_for_supported_operation --details',
      request_id: 'req-recording-unknown',
      retryable: false,
    })).toEqual({
      code: undefined,
      action: undefined,
      requestId: 'req-recording-unknown',
      retryable: false,
    });
  });

  it('preserves a bounded existing operation id without retaining arbitrary fields', () => {
    expect(normalizeTypedEnvelope({
      error: {
        code: 'RECORDING_TEST_WRITE_NOT_IMPLEMENTED',
        action: 'wait_for_supported_operation',
        retryable: false,
        request_id: 'req-recording-operation',
      },
      operation_id: 'op-recording-1',
      diagnostics: { secret: 'must not survive' },
    })).toEqual({
      code: 'RECORDING_TEST_WRITE_NOT_IMPLEMENTED',
      action: 'wait_for_supported_operation',
      requestId: 'req-recording-operation',
      retryable: false,
      operationId: 'op-recording-1',
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

  it('parses real backend mapping preview response with input and output keys', () => {
    const backendPayload = {
      raw_value: 243,
      final_value: 24.3,
      step_results: [
        { step_index: 0, step_type: 'decode', input: 243, output: 243 },
        { step_index: 1, step_type: 'scale', input: 243, output: 24.3 },
        { step_index: 2, step_type: 'cast', input: 24.3, output: 24.3 },
      ],
    };

    const parsed = parseMappingPreviewResponse(backendPayload);
    expect(parsed).not.toBeNull();
    expect(parsed?.raw_value).toBe(243);
    expect(parsed?.final_value).toBe(24.3);
    expect(parsed?.step_results).toHaveLength(3);
    expect(parsed?.step_results[0].input_value).toBe(243);
    expect(parsed?.step_results[0].output_value).toBe(243);
    expect(parsed?.step_results[1].output_value).toBe(24.3);
  });

  it('parses mapping preview response with wrapped data envelope and input_value keys', () => {
    const wrappedPayload = {
      success: true,
      data: {
        raw_value: 100,
        final_value: 200,
        step_results: [
          { step_index: 0, step_type: 'scale', input_value: 100, output_value: 200 },
        ],
      },
    };

    const parsed = parseMappingPreviewResponse(wrappedPayload);
    expect(parsed).not.toBeNull();
    expect(parsed?.raw_value).toBe(100);
    expect(parsed?.final_value).toBe(200);
    expect(parsed?.step_results[0].output_value).toBe(200);
  });
});
