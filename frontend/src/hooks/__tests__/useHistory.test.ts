import { describe, it, expect } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useHistory, usePointHistory, type HistoryAction } from '../useHistory';

describe('useHistory', () => {
  it('should initialize with empty state', () => {
    const { result } = renderHook(() => useHistory());

    expect(result.current.history).toHaveLength(0);
    expect(result.current.currentIndex).toBe(-1);
    expect(result.current.canUndo).toBe(false);
    expect(result.current.canRedo).toBe(false);
  });

  it('should push action to history', () => {
    const { result } = renderHook(() => useHistory<string>());

    act(() => {
      result.current.push({
        type: 'create',
        description: 'Created point',
        data: 'point-1',
      });
    });

    expect(result.current.history).toHaveLength(1);
    expect(result.current.currentIndex).toBe(0);
    expect(result.current.canUndo).toBe(true);
    expect(result.current.canRedo).toBe(false);
  });

  it('should undo action', () => {
    const { result } = renderHook(() => useHistory<string>());

    act(() => {
      result.current.push({ type: 'create', description: 'Action 1', data: 'data-1' });
    });

    let undoneAction: HistoryAction<string> | null = null;
    act(() => {
      undoneAction = result.current.undo();
    });

    expect(undoneAction).not.toBeNull();
    expect((undoneAction as any)?.data).toBe('data-1');
    expect(result.current.currentIndex).toBe(-1);
    expect(result.current.canUndo).toBe(false);
    expect(result.current.canRedo).toBe(true);
  });

  it('should redo action', () => {
    const { result } = renderHook(() => useHistory<string>());

    act(() => {
      result.current.push({ type: 'create', description: 'Action 1', data: 'data-1' });
    });
    
    act(() => {
      result.current.undo();
    });

    let redoneAction: HistoryAction<string> | null = null;
    act(() => {
      redoneAction = result.current.redo();
    });

    expect(redoneAction).not.toBeNull();
    expect((redoneAction as any)?.data).toBe('data-1');
    expect(result.current.currentIndex).toBe(0);
    expect(result.current.canUndo).toBe(true);
    expect(result.current.canRedo).toBe(false);
  });

  it('should return null when undo with no history', () => {
    const { result } = renderHook(() => useHistory());

    let undoneAction: HistoryAction<unknown> | null = null;
    act(() => {
      undoneAction = result.current.undo();
    });

    expect(undoneAction).toBeNull();
  });

  it('should return null when redo with nothing to redo', () => {
    const { result } = renderHook(() => useHistory<string>());

    act(() => {
      result.current.push({ type: 'create', description: 'Action', data: 'data' });
    });

    let redoneAction: HistoryAction<string> | null = null;
    act(() => {
      redoneAction = result.current.redo();
    });

    expect(redoneAction).toBeNull();
  });

  it('should clear history after new action following undo', () => {
    const { result } = renderHook(() => useHistory<string>());

    act(() => {
      result.current.push({ type: 'create', description: 'Action 1', data: 'data-1' });
      result.current.push({ type: 'create', description: 'Action 2', data: 'data-2' });
      result.current.push({ type: 'create', description: 'Action 3', data: 'data-3' });
    });

    act(() => {
      result.current.undo();
      result.current.undo();
    });

    act(() => {
      result.current.push({ type: 'create', description: 'Action 4', data: 'data-4' });
    });

    expect(result.current.history).toHaveLength(2);
    expect(result.current.history[1].data).toBe('data-4');
    expect(result.current.canRedo).toBe(false);
  });

  it('should limit history size', () => {
    const { result } = renderHook(() => useHistory<number>({ maxHistory: 3 }));

    act(() => {
      for (let i = 1; i <= 5; i++) {
        result.current.push({ type: 'create', description: `Action ${i}`, data: i });
      }
    });

    expect(result.current.history).toHaveLength(3);
    expect(result.current.history[0].data).toBe(3);
    expect(result.current.history[2].data).toBe(5);
  });

  it('should clear all history', () => {
    const { result } = renderHook(() => useHistory<string>());

    act(() => {
      result.current.push({ type: 'create', description: 'Action 1', data: 'data-1' });
      result.current.push({ type: 'create', description: 'Action 2', data: 'data-2' });
    });

    act(() => {
      result.current.clear();
    });

    expect(result.current.history).toHaveLength(0);
    expect(result.current.currentIndex).toBe(-1);
    expect(result.current.canUndo).toBe(false);
    expect(result.current.canRedo).toBe(false);
  });

  it('should getUndoAction without changing state', () => {
    const { result } = renderHook(() => useHistory<string>());

    act(() => {
      result.current.push({ type: 'create', description: 'Action', data: 'data' });
    });

    const action = result.current.getUndoAction();
    expect(action?.data).toBe('data');
    expect(result.current.currentIndex).toBe(0);
  });

  it('should getRedoAction without changing state', () => {
    const { result } = renderHook(() => useHistory<string>());

    act(() => {
      result.current.push({ type: 'create', description: 'Action', data: 'data' });
    });
    
    act(() => {
      result.current.undo();
    });

    const action = result.current.getRedoAction();
    expect(action?.data).toBe('data');
    expect(result.current.currentIndex).toBe(-1);
  });

  it('should return null from getUndoAction when empty', () => {
    const { result } = renderHook(() => useHistory());

    const action = result.current.getUndoAction();
    expect(action).toBeNull();
  });

  it('should return null from getRedoAction when at end', () => {
    const { result } = renderHook(() => useHistory<string>());

    act(() => {
      result.current.push({ type: 'create', description: 'Action', data: 'data' });
    });

    const action = result.current.getRedoAction();
    expect(action).toBeNull();
  });

  it('should add timestamp to actions', () => {
    const { result } = renderHook(() => useHistory<string>());
    const before = Date.now();

    act(() => {
      result.current.push({ type: 'create', description: 'Action', data: 'data' });
    });

    const after = Date.now();
    expect(result.current.history[0].timestamp).toBeGreaterThanOrEqual(before);
    expect(result.current.history[0].timestamp).toBeLessThanOrEqual(after);
  });
});

describe('usePointHistory', () => {
  it('should work with PointHistoryData type', () => {
    const { result } = renderHook(() => usePointHistory());

    act(() => {
      result.current.push({
        type: 'create',
        description: 'Created point D100',
        data: {
          pointId: 'p-1',
          newData: { name: 'D100', address: 'D100' },
        },
      });
    });

    expect(result.current.history).toHaveLength(1);
    expect(result.current.history[0].data.pointId).toBe('p-1');
  });

  it('should handle batch_create type', () => {
    const { result } = renderHook(() => usePointHistory());

    act(() => {
      result.current.push({
        type: 'batch_create',
        description: 'Created 5 points',
        data: {
          pointIds: ['p-1', 'p-2', 'p-3', 'p-4', 'p-5'],
        },
      });
    });

    expect(result.current.history[0].data.pointIds).toHaveLength(5);
  });

  it('should handle delete type with previous data', () => {
    const { result } = renderHook(() => usePointHistory());

    act(() => {
      result.current.push({
        type: 'delete',
        description: 'Deleted point D100',
        data: {
          pointId: 'p-1',
          previousData: { name: 'D100', address: 'D100', enabled: true },
        },
      });
    });

    let undone: ReturnType<typeof result.current.undo> = null;
    act(() => {
      undone = result.current.undo();
    });
    expect(undone?.data.previousData).toEqual({ name: 'D100', address: 'D100', enabled: true });
  });
});
