/**
 * PollingGroup Query Hooks
 *
 * 使用 TanStack Query 管理 PollingGroup 資源的 Server State。
 */
import { useQuery } from '@tanstack/react-query';
import { pollingGroupAPI } from '../../services/datalink';
import { pollingGroupKeys } from './keys';

// =============================================================================
// 查詢 Hooks
// =============================================================================

/**
 * 取得 PollingGroup 列表
 *
 * @returns Query 結果
 */
export function usePollingGroupsQuery() {
  return useQuery({
    queryKey: pollingGroupKeys.lists(),
    queryFn: () => pollingGroupAPI.list(),
  });
}
