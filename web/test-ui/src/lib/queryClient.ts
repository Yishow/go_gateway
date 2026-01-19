/**
 * TanStack Query 客戶端配置
 *
 * 提供全域 QueryClient 實例，統一管理快取策略、重試邏輯與錯誤處理。
 */
import { QueryClient } from '@tanstack/react-query';

/**
 * 全域 QueryClient 實例
 *
 * 配置說明：
 * - staleTime: 30秒內資料視為新鮮，不會重新請求
 * - gcTime: 5分鐘後未使用的快取將被回收
 * - retry: 失敗時自動重試 2 次
 * - refetchOnWindowFocus: 視窗聚焦時重新驗證資料
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30 * 1000,
      gcTime: 5 * 60 * 1000,
      retry: 2,
      refetchOnWindowFocus: true,
    },
    mutations: {
      retry: 1,
    },
  },
});
