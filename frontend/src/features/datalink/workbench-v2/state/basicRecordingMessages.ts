import type { RecordingStartOperation } from '../../../../types/studioV2RecordingStart';
import { getSafeErrorMessage } from '../../../../utils/typedErrors';

const BASIC_REASON_KEYS: Record<string, string> = {
  preparation_required: 'step4.basic.reason.preparation_required',
  stale_intent: 'step4.basic.reason.stale_intent',
  invalid_scope: 'step4.basic.reason.invalid_scope',
  service_unavailable: 'step4.basic.reason.service_unavailable',
  scope_not_ready: 'step4.basic.reason.scope_not_ready',
  share_not_ready: 'step4.basic.reason.share_not_ready',
  start_failed: 'step4.basic.reason.start_failed',
};

const BASIC_ACTION_KEYS: Record<string, string> = {
  prepare_schema: 'step4.basic.action.prepare_schema',
  revalidate: 'step4.basic.action.revalidate',
  review_selection: 'step4.basic.action.review_selection',
  retry_same_request: 'step4.basic.action.retry_same_request',
  review_device: 'step4.basic.action.review_device',
  refresh_share: 'step4.basic.action.refresh_share',
  reload: 'step4.basic.action.reload',
};

const BASIC_START_ERROR_KEYS: Record<string, string> = {
  RECORDING_START_INVALID: 'step4.basic.error.start_invalid',
  RECORDING_START_INTENT_CHANGED: 'step4.basic.error.start_intent_changed',
  RECORDING_START_NOT_FOUND: 'step4.basic.error.start_not_found',
  RECORDING_START_BUSY: 'step4.basic.error.start_busy',
  BASIC_INTENT_CHANGED: 'step4.basic.error.basic_intent_changed',
  WRITE_GROUP_BASIC_INTENT_CHANGED: 'step4.basic.error.basic_intent_changed',
};

type SafeRecord = Record<string, unknown>;

function asRecord(value: unknown): SafeRecord | undefined {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
    ? value as SafeRecord
    : undefined;
}

function errorSources(error: unknown): SafeRecord[] {
  const root = asRecord(error);
  const response = asRecord(root?.response);
  const responseData = asRecord(response?.data);
  const rootData = asRecord(root?.data);
  return [
    asRecord(responseData?.error), responseData, asRecord(root?.error), rootData, root,
  ].filter((source): source is SafeRecord => Boolean(source));
}

function safeKnownField(error: unknown, field: 'code' | 'action', allowed: Record<string, string>): string | undefined {
  for (const source of errorSources(error)) {
    const value = source[field];
    if (typeof value === 'string' && Object.prototype.hasOwnProperty.call(allowed, value)) return value;
  }
  return undefined;
}

export function safeBasicOperationMessage(
  operation: RecordingStartOperation | undefined,
  t: (key: string, options?: Record<string, unknown>) => string,
): string | undefined {
  if (!operation?.reason) return undefined;
  return t(BASIC_REASON_KEYS[operation.reason] ?? 'step4.basic.reason.unconfirmed');
}

export function safeBasicNextAction(
  operation: RecordingStartOperation | undefined,
  t: (key: string, options?: Record<string, unknown>) => string,
): string | undefined {
  const key = operation?.next_action ? BASIC_ACTION_KEYS[operation.next_action] : undefined;
  return key ? t(key) : undefined;
}

/** Maps only known recording-start codes; raw backend messages never reach the DOM. */
export function safeBasicStartError(
  error: unknown,
  t: (key: string, options?: Record<string, unknown>) => string,
): string {
  const code = safeKnownField(error, 'code', BASIC_START_ERROR_KEYS);
  const safeError = getSafeErrorMessage(error, t);
  if (!code) return safeError.message;
  const message = t(BASIC_START_ERROR_KEYS[code]);
  return safeError.requestId ? `${message} (${t('errors.request_id')}: ${safeError.requestId})` : message;
}
