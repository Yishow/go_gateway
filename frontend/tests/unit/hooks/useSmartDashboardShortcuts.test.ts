import { renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useSmartDashboardShortcuts } from '@/hooks/useKeyboardShortcuts';

describe('useSmartDashboardShortcuts', () => {
  it('includes flow action shortcuts when handlers are provided', () => {
    const handlers = {
      onValidateFlow: vi.fn(),
      onActivateFlow: vi.fn(),
      onRecoverFlow: vi.fn(),
    };

    const { result } = renderHook(() => useSmartDashboardShortcuts(handlers));
    const shortcuts = result.current;

    expect(
      shortcuts.some((s) => s.key === 'Enter' && s.ctrl && !s.shift)
    ).toBe(true);
    expect(
      shortcuts.some((s) => s.key === 'Enter' && s.ctrl && s.shift)
    ).toBe(true);
    expect(
      shortcuts.some((s) => s.key.toLowerCase() === 'r' && s.alt)
    ).toBe(true);
  });

  it('does not include recover shortcut when recover handler is absent', () => {
    const handlers = {
      onValidateFlow: vi.fn(),
      onActivateFlow: vi.fn(),
    };

    const { result } = renderHook(() => useSmartDashboardShortcuts(handlers));
    const shortcuts = result.current;

    expect(
      shortcuts.some((s) => s.key.toLowerCase() === 'r' && s.alt)
    ).toBe(false);
  });
});
