import { normalizeTypedEnvelope } from './safeJson';
import { BACKEND_ERROR_CODES } from './backendErrorCodes';

export { BACKEND_ERROR_CODES } from './backendErrorCodes';

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

/** Returns a bounded state message without exposing exception or backend text. */
export function getSafeErrorStateMessage(error: unknown, fallback: string): string {
  const requestId = normalizeTypedEnvelope(error).requestId;
  return requestId ? `${fallback} (Request ID: ${requestId})` : fallback;
}

const KNOWN_ERROR_CODES = new Set<string>(BACKEND_ERROR_CODES);

function localizeAction(
  action: string | undefined,
  retryable: boolean,
  t: (key: string, options?: Record<string, unknown>) => string,
): string | undefined {
  if (action === 'wait_for_supported_operation') {
    return t('errors.wait_for_supported_operation');
  }
  return retryable ? t('errors.retry') : undefined;
}

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
      title: t('errors.generic_failure'),
      message: t('errors.generic_failure'),
      retryable: false,
    };
  }

  const envelope = normalizeTypedEnvelope(error);
  const code = envelope.code && KNOWN_ERROR_CODES.has(envelope.code)
    ? envelope.code
    : undefined;
  const requestId = envelope.requestId;
  const retryable = envelope.retryable === true;
  const localizedAction = localizeAction(envelope.action, retryable, t);

  if (code) {
    const localizedMessage = t(`errors.${code}`, {
      requestId: requestId || 'unknown',
      defaultValue: t('errors.generic_failure', { requestId: requestId || 'unknown' }),
    });
    return {
      code,
      title: code,
      message: localizedMessage,
      action: localizedAction,
      requestId,
      retryable,
    };
  }

  const genericMsg = t('errors.generic_failure', {
    requestId: requestId || 'unknown',
  });

  return {
    code: undefined,
    title: t('errors.generic_failure'),
    message: genericMsg,
    action: localizedAction,
    requestId,
    retryable,
  };
}
