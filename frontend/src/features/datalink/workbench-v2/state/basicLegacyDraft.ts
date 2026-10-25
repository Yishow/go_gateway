/**
 * The unsaved interval and completeness choice of a legacy Basic wide device.
 * It is a view draft only: it lives in this tab's session storage under the
 * workspace/device/purpose scope key of the Basic intent, so reload and Back/Next
 * restore it for the same device and never lend it to another device or
 * workspace. It is never read as a saved group.
 */
const PREFIX = 'wbv2.basic-legacy-draft.';

export interface BasicLegacyDraft {
  intervalSeconds?: number;
  incompletePolicy?: string;
}

function storage(): Storage | undefined {
  if (typeof window === 'undefined') return undefined;
  try {
    return window.sessionStorage;
  } catch {
    return undefined;
  }
}

export function loadBasicLegacyDraft(scopeKey: string): BasicLegacyDraft | undefined {
  if (!scopeKey) return undefined;
  try {
    const raw = storage()?.getItem(`${PREFIX}${scopeKey}`);
    if (!raw) return undefined;
    const parsed = JSON.parse(raw) as unknown;
    if (!parsed || typeof parsed !== 'object') return undefined;
    const { intervalSeconds, incompletePolicy } = parsed as Record<string, unknown>;
    return {
      ...(typeof intervalSeconds === 'number' && Number.isInteger(intervalSeconds) && intervalSeconds >= 1 ? { intervalSeconds } : {}),
      ...(incompletePolicy === 'skip_row' || incompletePolicy === 'partial' ? { incompletePolicy } : {}),
    };
  } catch {
    return undefined;
  }
}

export function saveBasicLegacyDraft(scopeKey: string, draft: BasicLegacyDraft): void {
  if (!scopeKey) return;
  try {
    storage()?.setItem(`${PREFIX}${scopeKey}`, JSON.stringify({ ...loadBasicLegacyDraft(scopeKey), ...draft }));
  } catch {
    // A storage failure only loses the remembered view, never the draft in memory.
  }
}

/** Forgets a draft once the saved group already holds exactly what it said. */
export function clearBasicLegacyDraft(scopeKey: string): void {
  if (!scopeKey) return;
  try {
    storage()?.removeItem(`${PREFIX}${scopeKey}`);
  } catch {
    // Nothing to clear when storage is unavailable.
  }
}
