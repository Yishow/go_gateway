import { useEffect, useCallback } from 'react';

/**
 * 鍵盤快捷鍵配置
 */
export interface KeyboardShortcut {
  key: string;
  ctrl?: boolean;
  alt?: boolean;
  shift?: boolean;
  description: string;
  action: () => void;
}

/**
 * 鍵盤快捷鍵 Hook
 * 
 * @param shortcuts 快捷鍵配置陣列
 * @param enabled 是否啟用快捷鍵
 */
export function useKeyboardShortcuts(
  shortcuts: KeyboardShortcut[],
  enabled: boolean = true
) {
  const handleKeyDown = useCallback((event: KeyboardEvent) => {
    if (!enabled) return;

    // 忽略在輸入框中的快捷鍵
    const target = event.target as HTMLElement | null;
    if (target) {
      if (
        target.tagName === 'INPUT' ||
        target.tagName === 'TEXTAREA' ||
        target.isContentEditable
      ) {
        // 但允許 Escape 鍵
        if (event.key !== 'Escape') {
          return;
        }
      }
    }

    for (const shortcut of shortcuts) {
      const keyMatch = event.key.toLowerCase() === shortcut.key.toLowerCase();
      const ctrlMatch = shortcut.ctrl ? event.ctrlKey || event.metaKey : !event.ctrlKey && !event.metaKey;
      const altMatch = shortcut.alt ? event.altKey : !event.altKey;
      const shiftMatch = shortcut.shift ? event.shiftKey : !event.shiftKey;

      if (keyMatch && ctrlMatch && altMatch && shiftMatch) {
        event.preventDefault();
        shortcut.action();
        return;
      }
    }
  }, [shortcuts, enabled]);

  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [handleKeyDown]);
}
