import { describe, expect, it } from 'vitest';
import { BACKEND_ERROR_CODES, getSafeErrorMessage, getSafeErrorStateMessage } from '../../../src/utils/typedErrors';

describe('typedErrors utility', () => {
  const mockT = (key: string, options?: Record<string, unknown>) => {
    if (key === 'errors.preview_invalid_request') return 'Invalid preview request parameters.';
    if (key === 'errors.runtime_device_not_found') return 'Device context not found.';
    if (key === 'errors.modbus_share_revision_conflict') return 'Workspace configuration revision conflict.';
    if (key === 'errors.RECORDING_SCHEMA_NOT_IMPLEMENTED') return 'Managed schema apply is unavailable.';
    if (key === 'errors.RECORDING_TEST_WRITE_NOT_IMPLEMENTED') return 'Recording test write is unavailable.';
    if (key === 'errors.wait_for_supported_operation') return 'Wait for a supported operation.';
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

  it('preserves exact recording code, translates the wait action, and keeps retryable false', () => {
    const result = getSafeErrorMessage({
      response: {
        status: 501,
        data: {
          success: false,
          error: {
            code: 'RECORDING_SCHEMA_NOT_IMPLEMENTED',
            message: 'managed schema apply is not implemented: secret detail',
            retryable: false,
            action: 'wait_for_supported_operation',
            request_id: 'req-recording-schema-full',
          },
        },
      },
    }, mockT);

    expect(result).toEqual({
      code: 'RECORDING_SCHEMA_NOT_IMPLEMENTED',
      title: 'RECORDING_SCHEMA_NOT_IMPLEMENTED',
      message: 'Managed schema apply is unavailable.',
      action: 'Wait for a supported operation.',
      requestId: 'req-recording-schema-full',
      retryable: false,
    });
    expect(result.message).not.toContain('secret detail');
  });

  it('drops unknown recording code and action after axios normalization', () => {
    const result = getSafeErrorMessage({
      response: {
        data: {
          error: {
            code: 'RECORDING_UNKNOWN',
            action: 'run raw command --token secret',
            request_id: 'req-recording-unknown',
            retryable: false,
            message: 'raw exception secret',
          },
        },
      },
    }, mockT);

    expect(result.code).toBeUndefined();
    expect(result.action).toBeUndefined();
    expect(result.requestId).toBe('req-recording-unknown');
    expect(result.retryable).toBe(false);
    expect(result.message).not.toContain('raw exception');
  });

  it('preserves the recording membership code with a localized repair message', () => {
    const t = (key: string, options?: Record<string, unknown>) => (
      key === 'errors.RECORDING_PLAN_MEMBERS_INVALID' ? '請先儲存每個測量項目，再調整記錄方案。' : mockT(key, options)
    );
    const result = getSafeErrorMessage({
      response: {
        status: 422,
        data: {
          success: false,
          error: {
            code: 'RECORDING_PLAN_MEMBERS_INVALID',
            message: 'recording plan members must be saved measurements: meas-secret',
            retryable: false,
            action: 'save each measurement in the workspace, then reload the plan',
            request_id: 'req-recording-members',
          },
        },
      },
    }, t);

    expect(result.code).toBe('RECORDING_PLAN_MEMBERS_INVALID');
    expect(result.message).toBe('請先儲存每個測量項目，再調整記錄方案。');
    expect(result.requestId).toBe('req-recording-members');
    expect(result.retryable).toBe(false);
    expect(JSON.stringify(result)).not.toContain('meas-secret');
  });

  it('does not show retry action when a known action is non-retryable', () => {
    const result = getSafeErrorMessage({
      code: 'RECORDING_SCHEMA_NOT_IMPLEMENTED',
      action: 'retry',
      retryable: false,
      request_id: 'req-recording-no-retry',
    }, mockT);

    expect(result.action).toBeUndefined();
    expect(result.retryable).toBe(false);
  });

  it('retains the backend-emitted preview_stream_closed code', () => {
    expect(BACKEND_ERROR_CODES).toContain('preview_stream_closed');

    const result = getSafeErrorMessage({
      code: 'preview_stream_closed',
      message: 'raw stream diagnostic with secret details',
      retryable: true,
      request_id: 'req-stream-closed-1',
    }, mockT);

    expect(result.code).toBe('preview_stream_closed');
    expect(result.message).not.toContain('secret details');
  });
});
