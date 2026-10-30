import type { QueryClient } from '@tanstack/react-query';
import { settingsAPI } from '../services/datalink';
import { DATALINK_BASE } from '../services/datalinkClient';
import { settingsKeys } from '../hooks/datalink/keys';
export type ThemePreference = 'system' | 'light' | 'dark';
export type Theme = 'light' | 'dark';
export type ThemeSaveState = 'idle' | 'unsaved' | 'pending' | 'saved' | 'error';
export const THEME_PREFERENCE_KEY = `gateway.theme.preference.v1:${DATALINK_BASE}`;
export function normalizeTheme(value: unknown): ThemePreference | null {
  return value === 'auto' || value === 'system' ? 'system' : value === 'light' || value === 'dark' ? value : null;
}
export function themeWireValue(value: ThemePreference): 'auto' | 'light' | 'dark' { return value === 'system' ? 'auto' : value; }
function localPreference(): ThemePreference {
  try {
    const stored = localStorage.getItem(THEME_PREFERENCE_KEY);
    if (stored) {
      const value: unknown = JSON.parse(stored);
      if (typeof value === 'object' && value !== null && 'version' in value && value.version === 1 && 'preference' in value) {
        const preference = normalizeTheme(value.preference); if (preference) return preference;
      }
    }
  } catch { /* Invalid or inaccessible versioned storage does not block legacy fallback. */ }
  try { const legacy = localStorage.getItem('theme'); if (legacy === 'light' || legacy === 'dark') return legacy; } catch { /* Optional browser storage. */ }
  return 'system';
}
function media(): MediaQueryList | null {
  try { return typeof window.matchMedia === 'function' ? window.matchMedia('(prefers-color-scheme: dark)') : null; } catch { return null; }
}
const writeQueues = new WeakMap<object, Promise<unknown>>();
/** One controller per actual QueryClient/app context; never scoped to a workspace. */
export function createThemePreference(client: QueryClient) {
  let preference = localPreference();
  let state = { preference, theme: (preference === 'system' ? media()?.matches ? 'dark' : 'light' : preference) as Theme, saveState: 'idle' as ThemeSaveState };
  const listeners = new Set<() => void>();
  let generation = 0, choice = 0, save = 0, reset = 0, active = false;
  const publish = (patch: Partial<typeof state>) => { state = { ...state, ...patch }; listeners.forEach(listener => listener()); };
  const resolve = () => publish({ preference, theme: preference === 'system' ? media()?.matches ? 'dark' : 'light' : preference });
  const persist = (selected?: ThemePreference) => {
    const capturedGeneration = generation, capturedChoice = choice, capturedSave = ++save, capturedReset = reset;
    publish({ saveState: 'pending' });
    const task = (writeQueues.get(settingsAPI) ?? Promise.resolve()).catch(() => undefined).then(async () => {
      if (!active || capturedGeneration !== generation || capturedReset !== reset || (selected !== undefined && capturedChoice !== choice)) return;
      // Whole-settings Save calls without a selection: read the live choice now,
      // not a stale draft captured before scheduler/Share writes completed.
      const value = themeWireValue(selected ?? preference);
      await settingsAPI.updateKey('theme', value);
      if (active && capturedGeneration === generation && capturedChoice === choice && capturedSave === save) {
        try { localStorage.setItem(THEME_PREFERENCE_KEY, JSON.stringify({ version: 1, preference })); } catch { /* Server success is independent of local cache. */ }
        publish({ saveState: 'saved' });
      }
    }).catch(error => {
      if (active && capturedGeneration === generation && capturedChoice === choice && capturedSave === save) publish({ saveState: 'error' });
      throw error;
    });
    writeQueues.set(settingsAPI, task); return task;
  };
  return {
    subscribe: (listener: () => void) => { listeners.add(listener); return () => { listeners.delete(listener); }; },
    snapshot: () => state,
    setPreference: (next: ThemePreference) => {
      choice += 1; preference = next; resolve();
      try { localStorage.setItem(THEME_PREFERENCE_KEY, JSON.stringify({ version: 1, preference })); } catch { /* Server save remains independently reported. */ }
      void persist(next).catch(() => undefined);
    },
    resetPreference: (next: ThemePreference) => {
      choice += 1; reset += 1; save += 1; preference = next; resolve(); publish({ saveState: 'unsaved' });
    },
    // Explicit Save intent predating a later draft-only reset cannot save that reset.
    captureSave: () => { const capturedReset = reset; return () => capturedReset === reset ? persist() : Promise.resolve(); },
    persistCurrent: () => persist(),
    start: () => {
      active = true; const capturedGeneration = ++generation, capturedChoice = choice;
      const query = media(); const changed = () => { if (preference === 'system') resolve(); };
      try { if (query?.addEventListener) query.addEventListener('change', changed); else query?.addListener(changed); } catch { /* Missing browser listener is nonfatal. */ }
      resolve();
      void client.fetchQuery({ queryKey: settingsKeys.items(), queryFn: () => settingsAPI.listItems(), staleTime: 0 }).then(items => {
        if (!active || capturedGeneration !== generation || capturedChoice !== choice) return;
        const raw = items.find(item => item.key === 'theme');
        // Only the actual stored wire enum is intent; a mapper default never enters here.
        const saved = raw && ['auto', 'light', 'dark'].includes(String(raw.value)) ? normalizeTheme(raw.value) : null;
        if (saved) { preference = saved; resolve(); }
      }).catch(() => { /* Keep local/system fallback; failed GET is not persisted intent. */ });
      return () => {
        active = false; generation += 1;
        try { if (query?.removeEventListener) query.removeEventListener('change', changed); else query?.removeListener(changed); } catch { /* Optional browser API. */ }
      };
    },
  };
}
