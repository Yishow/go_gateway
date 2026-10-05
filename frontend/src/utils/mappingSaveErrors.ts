import { normalizeTypedEnvelope } from './safeJson';
import type { MappingSaveError } from '../features/datalink/workbench-v2/state/types';

/** Retains bounded operator metadata, never backend diagnostic text. */
export function mappingSaveError(error: unknown): MappingSaveError {
  const envelope = normalizeTypedEnvelope(error);
  const root = typeof error === 'object' && error !== null ? error as Record<string, unknown> : {};
  const response = typeof root.response === 'object' && root.response !== null ? root.response as Record<string, unknown> : {};
  const status = response.status ?? root.status;
  return {
    code: envelope.code,
    action: envelope.action,
    requestId: envelope.requestId && /^[A-Za-z0-9_-]{1,128}$/.test(envelope.requestId) ? envelope.requestId : undefined,
    status: typeof status === 'number' && Number.isInteger(status) && status >= 400 && status <= 599 ? status : undefined,
  };
}
