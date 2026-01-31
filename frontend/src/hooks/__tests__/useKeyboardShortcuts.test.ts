import { renderHook } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useKeyboardShortcuts } from '../useKeyboardShortcuts';
import type { KeyboardShortcut } from '../useKeyboardShortcuts';

describe('useKeyboardShortcuts', () => {
  let originalAddEventListener: typeof window.addEventListener;
  let originalRemoveEventListener: typeof window.removeEventListener;
  let keydownHandler: ((event: KeyboardEvent) => void) | null = null;

  beforeEach(() => {
    originalAddEventListener = window.addEventListener;
    originalRemoveEventListener = window.removeEventListener;

    window.addEventListener = vi.fn((type, handler) => {
      if (type === 'keydown') {
        keydownHandler = handler as (event: KeyboardEvent) => void;
      }
    });

    window.removeEventListener = vi.fn();
  });

  afterEach(() => {
    window.addEventListener = originalAddEventListener;
    window.removeEventListener = originalRemoveEventListener;
    keydownHandler = null;
  });

  it('should register keyboard event listener', () => {
    const action = vi.fn();
    const shortcuts: KeyboardShortcut[] = [
      { key: 'b', ctrl: true, description: 'Test', action },
    ];

    renderHook(() => useKeyboardShortcuts(shortcuts));

    expect(window.addEventListener).toHaveBeenCalledWith('keydown', expect.any(Function));
  });

  it('should call action on matching shortcut', () => {
    const action = vi.fn();
    const shortcuts: KeyboardShortcut[] = [
      { key: 'b', ctrl: true, description: 'Batch', action },
    ];

    renderHook(() => useKeyboardShortcuts(shortcuts));

    // Simulate Ctrl+B
    const event = new KeyboardEvent('keydown', {
      key: 'b',
      ctrlKey: true,
    });
    Object.defineProperty(event, 'preventDefault', { value: vi.fn() });
    
    keydownHandler?.(event);

    expect(action).toHaveBeenCalled();
  });

  it('should handle Escape key', () => {
    const action = vi.fn();
    const shortcuts: KeyboardShortcut[] = [
      { key: 'Escape', description: 'Close', action },
    ];

    renderHook(() => useKeyboardShortcuts(shortcuts));

    const event = new KeyboardEvent('keydown', { key: 'Escape' });
    Object.defineProperty(event, 'preventDefault', { value: vi.fn() });
    
    keydownHandler?.(event);

    expect(action).toHaveBeenCalled();
  });

  it('should not trigger when disabled', () => {
    const action = vi.fn();
    const shortcuts: KeyboardShortcut[] = [
      { key: 'b', ctrl: true, description: 'Test', action },
    ];

    renderHook(() => useKeyboardShortcuts(shortcuts, false));

    const event = new KeyboardEvent('keydown', {
      key: 'b',
      ctrlKey: true,
    });
    
    keydownHandler?.(event);

    expect(action).not.toHaveBeenCalled();
  });

  it('should handle multiple shortcuts', () => {
    const action1 = vi.fn();
    const action2 = vi.fn();
    const shortcuts: KeyboardShortcut[] = [
      { key: 'b', ctrl: true, description: 'Batch', action: action1 },
      { key: 'm', ctrl: true, description: 'Mapping', action: action2 },
    ];

    renderHook(() => useKeyboardShortcuts(shortcuts));

    // Trigger Ctrl+M
    const event = new KeyboardEvent('keydown', {
      key: 'm',
      ctrlKey: true,
    });
    Object.defineProperty(event, 'preventDefault', { value: vi.fn() });
    
    keydownHandler?.(event);

    expect(action1).not.toHaveBeenCalled();
    expect(action2).toHaveBeenCalled();
  });

  it('should ignore shortcuts in input fields except Escape', () => {
    const action = vi.fn();
    const escapeAction = vi.fn();
    const shortcuts: KeyboardShortcut[] = [
      { key: 'b', ctrl: true, description: 'Batch', action },
      { key: 'Escape', description: 'Close', action: escapeAction },
    ];

    renderHook(() => useKeyboardShortcuts(shortcuts));

    // Create fake input element
    const input = document.createElement('input');
    document.body.appendChild(input);

    // Ctrl+B in input - should be ignored
    const ctrlBEvent = new KeyboardEvent('keydown', {
      key: 'b',
      ctrlKey: true,
    });
    Object.defineProperty(ctrlBEvent, 'target', { value: input });
    Object.defineProperty(ctrlBEvent, 'preventDefault', { value: vi.fn() });
    
    keydownHandler?.(ctrlBEvent);
    expect(action).not.toHaveBeenCalled();

    // Escape in input - should still work
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    Object.defineProperty(escEvent, 'target', { value: input });
    Object.defineProperty(escEvent, 'preventDefault', { value: vi.fn() });
    
    keydownHandler?.(escEvent);
    expect(escapeAction).toHaveBeenCalled();

    document.body.removeChild(input);
  });

  it('should handle shift modifier', () => {
    const action = vi.fn();
    const shortcuts: KeyboardShortcut[] = [
      { key: 's', ctrl: true, shift: true, description: 'Save All', action },
    ];

    renderHook(() => useKeyboardShortcuts(shortcuts));

    // Ctrl+Shift+S
    const event = new KeyboardEvent('keydown', {
      key: 's',
      ctrlKey: true,
      shiftKey: true,
    });
    Object.defineProperty(event, 'preventDefault', { value: vi.fn() });
    
    keydownHandler?.(event);

    expect(action).toHaveBeenCalled();
  });

  it('should handle alt modifier', () => {
    const action = vi.fn();
    const shortcuts: KeyboardShortcut[] = [
      { key: 'n', alt: true, description: 'New', action },
    ];

    renderHook(() => useKeyboardShortcuts(shortcuts));

    // Alt+N
    const event = new KeyboardEvent('keydown', {
      key: 'n',
      altKey: true,
    });
    Object.defineProperty(event, 'preventDefault', { value: vi.fn() });
    
    keydownHandler?.(event);

    expect(action).toHaveBeenCalled();
  });

  it('should cleanup on unmount', () => {
    const action = vi.fn();
    const shortcuts: KeyboardShortcut[] = [
      { key: 'b', ctrl: true, description: 'Test', action },
    ];

    const { unmount } = renderHook(() => useKeyboardShortcuts(shortcuts));
    
    unmount();

    expect(window.removeEventListener).toHaveBeenCalledWith('keydown', expect.any(Function));
  });
});