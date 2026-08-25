import { describe, expect, it } from 'vitest';
import { getSafeErrorMessage, getSafeErrorStateMessage } from '../../../src/utils/typedErrors';

describe('typedErrors utility', () => {
  const mockT = (key: string, options?: Record<string, unknown>) => {
    if (key === 'errors.preview_invalid_request') return 'Invalid preview request parameters.';
    if (key === 'errors.runtime_device_not_found') return 'Device context not found.';
    if (key === 'errors.modbus_share_revision_conflict') return 'Workspace configuration revision conflict.';
    if (key === 'errors.generic_failure') {
      return `Generic failure. Request ID: ${options?.requestId || 'unknown'}`;
    }
    return (options?.defaultValue as string) || key;
  };

  it('maps known error code to localized message', () => {
    const error = {
      code: 'preview_invalid_request',
      message: 'Raw exception: invalid transform step index 3',
      retryable: false,
      request_id: 'req-1234',
    };

    const res = getSafeErrorMessage(error, mockT);
    expect(res.title).toBe('preview_invalid_request');
    expect(res.code).toBe('preview_invalid_request');
    expect(res.message).toBe('Invalid preview request parameters.');
    expect(res.message).not.toContain('Raw exception');
    expect(res.requestId).toBe('req-1234');
    expect(res.retryable).toBe(false);
  });

  it('hides raw backend message when error code is unknown', () => {
    const error = {
      code: 'internal_database_panic',
      message: 'FATAL: connect failed at 10.0.0.1:5432 with user admin password secret',
      request_id: 'req-secret-999',
      retryable: true,
    };

    const res = getSafeErrorMessage(error, mockT);
    expect(res.message).toBe('Generic failure. Request ID: req-secret-999');
    expect(res.message).not.toContain('password');
    expect(res.message).not.toContain('FATAL');
    expect(res.requestId).toBe('req-secret-999');
    expect(res.retryable).toBe(true);
    expect(res.code).toBeUndefined();
    expect(res.title).toBe('Generic failure. Request ID: unknown');
  });

  it('handles nested error response structure', () => {
    const errorResponse = {
      success: false,
      error: {
        code: 'runtime_device_not_found',
        message: 'device dev-99 not registered',
        retryable: false,
      },
    };

    const res = getSafeErrorMessage(errorResponse, mockT);
    expect(res.title).toBe('runtime_device_not_found');
    expect(res.code).toBe('runtime_device_not_found');
    expect(res.message).toBe('Device context not found.');
  });

  it('unwraps axios responses and never returns raw action text', () => {
    const res = getSafeErrorMessage(
      {
        response: {
          data: {
            error: {
              code: 'runtime_stream_unavailable',
              message: 'dial tcp 10.0.0.9:502: secret details',
              action: 'run internal command --token secret',
              retryable: true,
              request_id: 'req-42',
            },
          },
        },
      },
      mockT,
    );

    expect(res.message).not.toContain('secret details');
    expect(res.action).not.toContain('internal command');
    expect(res.requestId).toBe('req-42');
    expect(res.retryable).toBe(true);
  });

  it('maps all Modbus Share backend lifecycle codes without exposing diagnostics', () => {
    const codes = [
      'modbus_share_reconcile_failed',
      'modbus_share_dirty_unknown',
      'modbus_share_projection_required',
      'modbus_share_invalid_geometry',
    ];

    for (const code of codes) {
      const result = getSafeErrorMessage({
        code,
        message: 'internal stack and secret details',
        retryable: true,
      }, mockT);
      expect(result.code).toBe(code);
      expect(result.message).not.toContain('internal stack');
      expect(result.message).not.toContain('secret details');
      expect(result.action).toBe('errors.retry');
    }
  });

  it('renders a 409 typed revision conflict through localized safe fields', () => {
    const result = getSafeErrorMessage({
      response: {
        status: 409,
        data: {
          success: false,
          error: {
            code: 'modbus_share_revision_conflict',
            message: 'current settings revision is settings-secret-9',
            action: 'reload internal state at /private/secret',
            retryable: true,
            request_id: 'req-conflict-9',
          },
        },
      },
    }, mockT);

    expect(result.code).toBe('modbus_share_revision_conflict');
    expect(result.message).toBe('Workspace configuration revision conflict.');
    expect(result.action).toBe('errors.retry');
    expect(result.requestId).toBe('req-conflict-9');
    expect(result.message).not.toContain('settings-secret-9');
    expect(result.action).not.toContain('/private/secret');
  });

  it('keeps autosave state bounded while retaining only the request id', () => {
    const result = getSafeErrorStateMessage({
      response: {
        data: {
          error: {
            message: 'secret backend diagnostic',
            request_id: 'req-save-7',
          },
        },
      },
    }, 'Save failed');
    expect(result).toBe('Save failed (Request ID: req-save-7)');
    expect(result).not.toContain('secret backend diagnostic');
  });

  it('drops malformed and oversized typed fields before they reach safe UI state', () => {
    const oversized = 'x'.repeat(129);
    const malformedInputs: unknown[] = [
      [],
      null,
      0,
      '',
      42,
      { error: [] },
      { error: { code: 42, request_id: ['request'], action: { label: 'retry' }, retryable: 'yes' } },
      {
        response: {
          data: {
            error: {
              code: oversized,
              request_id: oversized,
              action: oversized,
              retryable: 1,
            },
          },
        },
      },
    ];

    for (const input of malformedInputs) {
      const result = getSafeErrorMessage(input, mockT);
      expect(result.code).toBeUndefined();
      expect(result.requestId).toBeUndefined();
      expect(result.action).toBeUndefined();
      expect(result.retryable).toBe(false);
      expect(result.message).not.toContain(oversized);
    }
  });

  it('normalizes bounded request ids and rejects invalid field types in nested responses', () => {
    const result = getSafeErrorMessage({
      response: {
        data: {
          error: {
            code: ' runtime_device_not_found ',
            request_id: '  req-42  ',
            action: ' retry runtime stream ',
            retryable: true,
          },
        },
      },
    }, mockT);

    expect(result.code).toBe('runtime_device_not_found');
    expect(result.requestId).toBe('req-42');
    expect(result.action).toBe('errors.retry');
    expect(result.retryable).toBe(true);
  });
});
