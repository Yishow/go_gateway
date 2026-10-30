import * as React from 'react';
import { createContext, useContext, useEffect, useMemo, useSyncExternalStore } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { createThemePreference, type Theme, type ThemePreference, type ThemeSaveState } from './themePreference';
interface ThemeContextType {
  theme: Theme;
  preference: ThemePreference;
  saveState: ThemeSaveState;
  toggleTheme: () => void;
  setPreference: (preference: ThemePreference) => void;
  persistCurrent: () => Promise<unknown>;
  resetPreference: (preference: ThemePreference) => void;
  captureSave: () => () => Promise<unknown>;
}
const ThemeContext = createContext<ThemeContextType | undefined>(undefined);
export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const client = useQueryClient();
  const owner = useMemo(() => createThemePreference(client), [client]);
  const snapshot = useSyncExternalStore(owner.subscribe, owner.snapshot, owner.snapshot);
  useEffect(() => owner.start(), [owner]);
  useEffect(() => {
    const root = document.documentElement;
    root.classList.remove('light', 'dark'); root.classList.add(snapshot.theme);
    root.style.colorScheme = snapshot.theme;
  }, [snapshot.theme]);
  return <ThemeContext.Provider value={{ ...snapshot, setPreference: owner.setPreference, persistCurrent: owner.persistCurrent, resetPreference: owner.resetPreference, captureSave: owner.captureSave,
    toggleTheme: () => owner.setPreference(snapshot.theme === 'light' ? 'dark' : 'light') }}>{children}</ThemeContext.Provider>;
}
/** Optional only for independently mounted legacy settings components; production has one provider. */
export function useOptionalTheme() { return useContext(ThemeContext); }
export function useTheme() {
  const context = useOptionalTheme();
  if (!context) throw new Error('useTheme must be used within a ThemeProvider');
  return context;
}
