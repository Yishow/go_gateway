import { RecordingPlanResponseError } from '../utils/recordingPlanJson';
import { WriteGroupResponseError } from '../utils/studioV2WriteGroupJson';
import { isRecord } from '../utils/studioV2WriteGroupJsonPrimitives';
import { normalizeTypedEnvelope } from '../utils/safeJson';

/** Detach every error from Axios configuration, endpoint, token and raw response. */
export async function safeLifecycleRequest<T>(run: () => Promise<T>): Promise<T> {
  try { return await run(); }
  catch (error) {
    if (error instanceof WriteGroupResponseError || error instanceof RecordingPlanResponseError) throw error;
    const response = isRecord(error) && isRecord(error.response) ? error.response : undefined;
    const status = response?.status;
    const body = response?.data;
    const typed = normalizeTypedEnvelope(error);
    const rejected = typeof status === 'number' && status >= 400 && status < 500 && status !== 408 &&
      isRecord(body) && body.success === false && Boolean(typed.code);
    throw new WriteGroupResponseError('recording operation', error, rejected ? 'failed' : 'unconfirmed');
  }
}
export function requireLifecycleMatch<T>(result: T, matches: boolean): T {
  if (!matches) throw new WriteGroupResponseError('recording operation');
  return result;
}
export function isWriteGroupSuccessorUnconfirmed(error: unknown): boolean {
  return !(error instanceof WriteGroupResponseError || error instanceof RecordingPlanResponseError) || error.outcome === 'unconfirmed';
}

/** Property ordering is irrelevant; every reviewed field and array order is not. */
export function sameLifecycleJSON(a: unknown, b: unknown): boolean {
  if (a === b) return true;
  if (Array.isArray(a) || Array.isArray(b)) return Array.isArray(a) && Array.isArray(b) && a.length === b.length && a.every((item, i) => sameLifecycleJSON(item, b[i]));
  if (!isRecord(a) || !isRecord(b)) return false;
  const keys = Object.keys(a);
  return keys.length === Object.keys(b).length && keys.every(key => Object.prototype.hasOwnProperty.call(b, key) && sameLifecycleJSON(a[key], b[key]));
}
