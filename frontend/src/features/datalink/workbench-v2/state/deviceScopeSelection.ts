/** The selected persisted device of one workspace, kept across reload/Back-Next. */
const PREFIX = 'wbv2.device-scope.';

function storage(): Storage | undefined {
  if (typeof window === 'undefined') return undefined;
  try {
    return window.sessionStorage;
  } catch {
    return undefined;
  }
}

function key(workspaceId: string): string {
  return `${PREFIX}${encodeURIComponent(workspaceId)}`;
}

/** Returns the saved selection only when it is still a legal device ID. */
export function loadSelectedDeviceId(workspaceId: string, allowed: string[]): string | undefined {
  if (!workspaceId.trim()) return undefined;
  try {
    const value = storage()?.getItem(key(workspaceId));
    return value && allowed.includes(value) ? value : undefined;
  } catch {
    return undefined;
  }
}

export function saveSelectedDeviceId(workspaceId: string, deviceId: string): void {
  if (!workspaceId.trim() || !deviceId) return;
  try {
    storage()?.setItem(key(workspaceId), deviceId);
  } catch {
    // A storage failure only loses the remembered view, never the draft.
  }
}
