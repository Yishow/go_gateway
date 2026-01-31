/**
 * 撤銷/重做 Hook
 *
 * 提供操作歷史記錄和撤銷/重做功能
 */

import { useState, useCallback, useMemo } from 'react';

export interface HistoryAction<T = unknown> {
  type: string;
  description: string;
  data: T;
  timestamp: number;
}

export interface UseHistoryOptions {
  maxHistory?: number;
}

export interface UseHistoryReturn<T> {
  /** 當前歷史記錄 */
  history: HistoryAction<T>[];
  /** 當前位置索引 */
  currentIndex: number;
  /** 是否可以撤銷 */
  canUndo: boolean;
  /** 是否可以重做 */
  canRedo: boolean;
  /** 新增操作到歷史 */
  push: (action: Omit<HistoryAction<T>, 'timestamp'>) => void;
  /** 撤銷上一個操作 */
  undo: () => HistoryAction<T> | null;
  /** 重做下一個操作 */
  redo: () => HistoryAction<T> | null;
  /** 清除歷史記錄 */
  clear: () => void;
  /** 取得可撤銷的操作 */
  getUndoAction: () => HistoryAction<T> | null;
  /** 取得可重做的操作 */
  getRedoAction: () => HistoryAction<T> | null;
}

/**
 * 通用歷史記錄 Hook
 */
export function useHistory<T = unknown>(
  options: UseHistoryOptions = {}
): UseHistoryReturn<T> {
  const { maxHistory = 50 } = options;

  const [state, setState] = useState<{
    history: HistoryAction<T>[];
    currentIndex: number;
  }>({
    history: [],
    currentIndex: -1,
  });

  const canUndo = useMemo(() => state.currentIndex >= 0, [state.currentIndex]);
  const canRedo = useMemo(
    () => state.currentIndex < state.history.length - 1,
    [state.currentIndex, state.history.length]
  );

  const push = useCallback(
    (action: Omit<HistoryAction<T>, 'timestamp'>) => {
      setState((prev) => {
        // 如果不在歷史末端，移除後面的記錄
        const trimmed = prev.history.slice(0, prev.currentIndex + 1);
        const newAction: HistoryAction<T> = {
          ...action,
          timestamp: Date.now(),
        };

        // 限制歷史大小
        const updated = [...trimmed, newAction];
        const finalHistory = updated.length > maxHistory ? updated.slice(-maxHistory) : updated;
        
        return {
          history: finalHistory,
          currentIndex: finalHistory.length - 1,
        };
      });
    },
    [maxHistory]
  );

  const undo = useCallback((): HistoryAction<T> | null => {
    if (!canUndo) return null;

    const action = state.history[state.currentIndex];
    setState((prev) => ({
      ...prev,
      currentIndex: prev.currentIndex - 1,
    }));
    return action;
  }, [canUndo, state.currentIndex, state.history]);

  const redo = useCallback((): HistoryAction<T> | null => {
    if (!canRedo) return null;

    const action = state.history[state.currentIndex + 1];
    setState((prev) => ({
      ...prev,
      currentIndex: prev.currentIndex + 1,
    }));
    return action;
  }, [canRedo, state.currentIndex, state.history]);

  const clear = useCallback(() => {
    setState({ history: [], currentIndex: -1 });
  }, []);

  const getUndoAction = useCallback((): HistoryAction<T> | null => {
    if (!canUndo) return null;
    return state.history[state.currentIndex];
  }, [canUndo, state.currentIndex, state.history]);

  const getRedoAction = useCallback((): HistoryAction<T> | null => {
    if (!canRedo) return null;
    return state.history[state.currentIndex + 1];
  }, [canRedo, state.currentIndex, state.history]);

  return {
    history: state.history,
    currentIndex: state.currentIndex,
    canUndo,
    canRedo,
    push,
    undo,
    redo,
    clear,
    getUndoAction,
    getRedoAction,
  };
}

// =============================================================================
// 點位操作歷史類型
// =============================================================================

export type PointActionType = 'create' | 'update' | 'delete' | 'batch_create' | 'import';

export interface PointHistoryData {
  pointId?: string;
  pointIds?: string[];
  previousData?: Record<string, unknown>;
  newData?: Record<string, unknown>;
}

export type PointHistoryAction = HistoryAction<PointHistoryData>;

/**
 * 點位操作歷史 Hook
 */
export function usePointHistory(options?: UseHistoryOptions) {
  return useHistory<PointHistoryData>(options);
}
