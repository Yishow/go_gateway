import { describe, expect, it } from 'vitest';
import { BACKEND_ERROR_CODES } from '../../../src/utils/backendErrorCodes';
import { normalizeTypedEnvelope } from '../../../src/utils/safeJson';
import { getSafeErrorMessage } from '../../../src/utils/typedErrors';

const mockT = (key: string) => key;

describe('backendErrorCodes single source', () => {
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
