export type PointFormatErrorCode = 'point_format_scope' | 'point_format_conflict' | 'point_format_invalid' | 'point_format_unavailable' | 'point_format_operation_not_found' | 'point_format_unknown';
export class PointFormatError extends Error {
  readonly name = 'PointFormatError';
  constructor(readonly code: PointFormatErrorCode, readonly outcome: 'rejected' | 'unknown', readonly requestId?: string) {
    super(code);
  }
}
const codes: Record<number, PointFormatErrorCode> = {
  403: 'point_format_scope', 409: 'point_format_conflict', 422: 'point_format_invalid',
  503: 'point_format_unavailable', 404: 'point_format_operation_not_found',
};
export function record(value: unknown): Record<string, unknown> | undefined {
  return value !== null && typeof value === 'object' && !Array.isArray(value) ? value as Record<string, unknown> : undefined;
}
/** Copy only allowlisted classification; never retain Axios response/config/cause. */
export function safePointFormatError(error: unknown, requestId?: string): PointFormatError {
  if (error instanceof PointFormatError) return new PointFormatError(error.code, error.outcome, requestId ?? error.requestId);
  const response = record(record(error)?.response);
  const status = response?.status;
  const body = record(response?.data);
  const code = record(body?.error)?.code;
  if (typeof status === 'number' && codes[status] === code) {
    return new PointFormatError(codes[status], status >= 500 || status === 404 ? 'unknown' : 'rejected', requestId);
  }
  return new PointFormatError('point_format_unknown', 'unknown', requestId);
}
export const pointFormatErrorKey = (error: PointFormatError): string => `point-format:errors.${error.code}`;
