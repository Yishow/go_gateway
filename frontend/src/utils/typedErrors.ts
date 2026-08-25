export interface TypedAPIErrorPayload {
  code?: string;
  message?: string;
  retryable?: boolean;
  request_id?: string;
  action?: string;
  workspace_revision?: string;
  settings_revision?: string;
  dirty_state?: string;
}

export interface SafeErrorMessage {
  code?: string;
  title: string;
  message: string;
  action?: string;
  requestId?: string;
  retryable: boolean;
}

const MAX_SAFE_ERROR_FIELD_LENGTH = 128;

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function normalizeSafeString(value: unknown): string | undefined {
  if (typeof value !== 'string') return undefined;
  const normalized = value.trim();
  return normalized.length > 0 && normalized.length <= MAX_SAFE_ERROR_FIELD_LENGTH
    ? normalized
    : undefined;
}

function extractErrorPayload(error: unknown): Record<string, unknown> | undefined {
  if (!isRecord(error)) return undefined;

  const response = isRecord(error.response) ? error.response : undefined;
  const responseData = response && isRecord(response.data) ? response.data : undefined;
  if (responseData && isRecord(responseData.error)) return responseData.error;
  if (isRecord(error.error)) return error.error;
  return error;
}

/** Returns a bounded state message without exposing exception or backend text. */
export function getSafeErrorStateMessage(error: unknown, fallback = 'Save failed'): string {
  const payload = extractErrorPayload(error);
  const requestId = normalizeSafeString(payload?.request_id);
  return requestId ? `${fallback} (Request ID: ${requestId})` : fallback;
}

export const BACKEND_ERROR_CODES = [
  'preview_invalid_request',
  'preview_unavailable',
  'runtime_device_not_found',
  'runtime_snapshot_unavailable',
  'runtime_stream_unavailable',
  'workspace_not_ready',
  'activation_failed',
  'settings_unavailable',
  'settings_update_failed',
  'settings_invalid',
  'activation_request_invalid',
  'readiness_blocked',
  'modbus_share_save_incomplete',
  'modbus_share_revision_conflict',
  'modbus_share_hydration_required',
  'modbus_share_disabled',
  'modbus_share_range_collision',
  'modbus_share_capacity_exceeded',
  'modbus_share_listener_bind_failed',
  'modbus_share_workspace_scope',
  'modbus_share_reconcile_failed',
  'modbus_share_dirty_unknown',
  'modbus_share_projection_required',
  'modbus_share_invalid_geometry',
  'validation',
  'revision_mismatch',
  'not_found',
  'internal',
] as const;

const KNOWN_ERROR_CODES = new Set<string>(BACKEND_ERROR_CODES);

/**
 * Extracts and maps a typed error payload into a safe localized operator message.
 * Keeps raw exception text, stack traces, and internal strings outside of user-facing copy.
 */
export function getSafeErrorMessage(
  error: unknown,
  t: (key: string, options?: Record<string, unknown>) => string,
): SafeErrorMessage {
  if (error === null || error === undefined) {
    return {
      code: undefined,
      title: t('errors.generic_failure', { defaultValue: 'Operation failed' }),
      message: t('errors.generic_failure', { defaultValue: 'An unexpected error occurred.' }),
      retryable: false,
    };
  }

  // Axios errors keep the server envelope under response.data; unwrap it
  // without copying the backend message into operator-facing state.
  const payload = extractErrorPayload(error);
  const normalizedCode = normalizeSafeString(payload?.code);
  const code = normalizedCode && KNOWN_ERROR_CODES.has(normalizedCode)
    ? normalizedCode
    : undefined;
  const requestId = normalizeSafeString(payload?.request_id);
  const action = normalizeSafeString(payload?.action);
  const retryable = payload?.retryable === true;

  // The server action is advisory only. Keep it bounded and ignore unknown
  // values so untrusted operator-facing text never reaches JSX.
  const localizedRetryAction = retryable && (!action || action.length <= MAX_SAFE_ERROR_FIELD_LENGTH)
    ? t('errors.retry')
    : undefined;

  if (code) {
    const localizedMessage = t(`errors.${code}`, {
      defaultValue: t('errors.generic_failure', { requestId: requestId || 'unknown' }),
    });
    return {
      code,
      title: code,
      message: localizedMessage,
      action: localizedRetryAction,
      requestId,
      retryable,
    };
  }

  const genericMsg = t('errors.generic_failure', {
    requestId: requestId || 'unknown',
    defaultValue: `An unexpected error occurred. Request ID: ${requestId || 'unknown'}`,
  });

  return {
    code: undefined,
    title: t('errors.generic_failure', { defaultValue: 'Operation failed' }),
    message: genericMsg,
    action: localizedRetryAction,
    requestId,
    retryable,
  };
}
