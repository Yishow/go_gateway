import { describe, expect, it } from 'vitest';
import { BACKEND_ERROR_CODES } from '../../../src/utils/backendErrorCodes';
import { normalizeTypedEnvelope } from '../../../src/utils/safeJson';
import { getSafeErrorMessage } from '../../../src/utils/typedErrors';

const mockT = (key: string) => key;

describe('backendErrorCodes single source', () => {
  it('retains the legacy write conflict and shows the write-group repair action safely', () => {
    const failure = {
      response: {
        data: {
          success: false,
          error: {
            code: 'WRITE_GROUP_LEGACY_WRITE_CONFLICT',
            action: 'open_write_groups',
            request_id: 'legacy-write-conflict-1',
            retryable: false,
            message: 'private DSN diagnostic',
          },
        },
      },
    };
    expect(normalizeTypedEnvelope(failure)).toEqual({
      code: 'WRITE_GROUP_LEGACY_WRITE_CONFLICT',
      action: 'open_write_groups',
      requestId: 'legacy-write-conflict-1',
      retryable: false,
    });
    expect(getSafeErrorMessage(failure, mockT)).toEqual({
      code: 'WRITE_GROUP_LEGACY_WRITE_CONFLICT',
      title: 'WRITE_GROUP_LEGACY_WRITE_CONFLICT',
      message: 'errors.WRITE_GROUP_LEGACY_WRITE_CONFLICT',
      action: 'errors.open_write_groups',
      requestId: 'legacy-write-conflict-1',
      retryable: false,
    });
  });

  it('allowlists every write-group code and action the handlers emit', () => {
    for (const code of ['WRITE_GROUP_LEGACY_READ_CONFLICT', 'WRITE_GROUP_LIFECYCLE_BLOCKED', 'WRITE_GROUP_NOT_FOUND', 'WRITE_GROUP_INVALID', 'WRITE_GROUP_INVALID_REQUEST', 'WRITE_GROUP_UNAVAILABLE']) {
      for (const action of ['review_request', 'review_group', 'disable_group']) {
        const failure = { response: { data: { success: false, error: { code, action, request_id: 'r1', retryable: false, message: 'private detail' } } } };
        expect(normalizeTypedEnvelope(failure)).toEqual({ code, action, requestId: 'r1', retryable: false });
        expect(getSafeErrorMessage(failure, mockT)).toMatchObject({ message: `errors.${code}`, action: `errors.${action}` });
      }
    }
  });

  it('has no duplicate codes', () => {
    expect(new Set(BACKEND_ERROR_CODES).size).toBe(BACKEND_ERROR_CODES.length);
  });

  it('is the allowlist source for normalizeTypedEnvelope', () => {
    for (const code of BACKEND_ERROR_CODES) {
      expect(normalizeTypedEnvelope({ code, retryable: true }).code).toBe(code);
    }
  });

  it('is the allowlist source for getSafeErrorMessage', () => {
    for (const code of BACKEND_ERROR_CODES) {
      expect(getSafeErrorMessage({ code, message: 'raw secret', retryable: false }, mockT).code).toBe(code);
    }
  });
});
