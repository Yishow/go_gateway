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

/**
 * 預設快捷鍵列表 Hook
 */
export function useSmartDashboardShortcuts(handlers: {
  onBatchCreate?: () => void;
  onQuickMapping?: () => void;
  onClosePanel?: () => void;
  onToggleSidebar?: () => void;
  onSearch?: () => void;
  onUndo?: () => void;
  onRedo?: () => void;
  onImport?: () => void;
  onExport?: () => void;
}) {
  const shortcuts: KeyboardShortcut[] = [];

  if (handlers.onBatchCreate) {
    shortcuts.push({
      key: 'b',
      ctrl: true,
      description: '批量建立點位',
      action: handlers.onBatchCreate,
    });
  }

  if (handlers.onQuickMapping) {
    shortcuts.push({
      key: 'm',
      ctrl: true,
      description: '快速映射',
      action: handlers.onQuickMapping,
    });
  }

  if (handlers.onClosePanel) {
    shortcuts.push({
      key: 'Escape',
      description: '關閉面板',
      action: handlers.onClosePanel,
    });
  }

  if (handlers.onToggleSidebar) {
    shortcuts.push({
      key: '[',
      ctrl: true,
      description: '切換側邊欄',
      action: handlers.onToggleSidebar,
    });
  }

  if (handlers.onSearch) {
    shortcuts.push({
      key: 'k',
      ctrl: true,
      description: '搜尋',
      action: handlers.onSearch,
    });
  }

  if (handlers.onUndo) {
    shortcuts.push({
      key: 'z',
      ctrl: true,
      description: '撤銷',
      action: handlers.onUndo,
    });
  }

  if (handlers.onRedo) {
    shortcuts.push({
      key: 'y',
      ctrl: true,
      description: '重做',
      action: handlers.onRedo,
    });
  }

  if (handlers.onImport) {
    shortcuts.push({
      key: 'i',
      ctrl: true,
      description: '匯入點位',
      action: handlers.onImport,
    });
  }

  if (handlers.onExport) {
    shortcuts.push({
      key: 'e',
      ctrl: true,
      description: '匯出點位',
      action: handlers.onExport,
    });
  }

  useKeyboardShortcuts(shortcuts);

  return shortcuts;
}
