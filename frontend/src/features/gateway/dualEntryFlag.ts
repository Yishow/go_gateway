import { settingsAPI } from '../../services/datalink';

export const ENABLE_GATEWAY_DUAL_ENTRY_KEY = 'ENABLE_GATEWAY_DUAL_ENTRY';

export function coerceBooleanFlag(value: unknown, fallback = false): boolean {
  if (typeof value === 'boolean') return value;
  if (typeof value === 'number') return value !== 0;

  if (typeof value === 'string') {
    const normalized = value.trim().toLowerCase();

    if (['1', 'true', 'yes', 'y', 'on', 'enabled'].includes(normalized)) {
      return true;
    }
    if (['0', 'false', 'no', 'n', 'off', 'disabled'].includes(normalized)) {
      return false;
    }
  }

  return fallback;
}

export async function fetchGatewayDualEntryEnabled(): Promise<boolean> {
  const items = await settingsAPI.listItems();
  const match = items.find((item) => item.key.toUpperCase() === ENABLE_GATEWAY_DUAL_ENTRY_KEY);
  return coerceBooleanFlag(match?.value, false);
}

export function resolveGatewayCreateEntryPath(enabled: boolean): string {
  return enabled ? '/gateway/entry' : '/datalink';
}
