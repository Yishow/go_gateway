/**
 * PollingGroupsPage - 輪詢群組管理頁面
 *
 * 使用 TanStack Query 進行 Server State 管理。
 * 提供輪詢群組的 CRUD 操作介面。
 */
import { useState } from 'react';
import type {
  PollingGroup,
  CreatePollingGroupRequest,
  UpdatePollingGroupRequest,
} from '../../types/datalink';
import PollingGroupForm from '../../components/datalink/PollingGroupForm';
import ConfirmDialog from '../../components/datalink/ConfirmDialog';
import { useToast } from '../../contexts/ToastContext';
import {
  usePollingGroupsQuery,
  useCreatePollingGroupMutation,
  useUpdatePollingGroupMutation,
  useDeletePollingGroupMutation,
} from '../../hooks/datalink/usePollingGroups';

/**
 * 格式化間隔時間顯示
 */
function formatInterval(ms: number): string {
  if (ms < 1000) {
    return `${ms}ms`;
  }
  if (ms < 60000) {
    return `${(ms / 1000).toFixed(1)}s`;
  }
  return `${(ms / 60000).toFixed(1)}min`;
}

export default function PollingGroupsPage() {
  const { data: pollingGroups = [], isLoading, error } = usePollingGroupsQuery();
  const { showSuccess, showError } = useToast();

  // Mutation Hooks
  const createMutation = useCreatePollingGroupMutation();
  const updateMutation = useUpdatePollingGroupMutation();
  const deleteMutation = useDeletePollingGroupMutation();

  // Modal 狀態
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingGroup, setEditingGroup] = useState<PollingGroup | null>(null);

  // 確認對話框狀態
  const [confirmDialog, setConfirmDialog] = useState<{
    isOpen: boolean;
    title: string;
    message: string;
    onConfirm: () => void;
    variant?: 'danger' | 'warning' | 'default';
  }>({
    isOpen: false,
    title: '',
    message: '',
    onConfirm: () => {},
  });

  /**
   * 顯示確認對話框
   */
  const showConfirm = (
    title: string,
    message: string,
    onConfirm: () => void,
    variant: 'danger' | 'warning' | 'default' = 'default'
  ) => {
    setConfirmDialog({
      isOpen: true,
      title,
      message,
      onConfirm,
      variant,
    });
  };

  /**
   * 關閉確認對話框
   */
  const closeConfirm = () => {
    setConfirmDialog((prev) => ({ ...prev, isOpen: false }));
  };

  /**
   * 刪除輪詢群組
   */
  const handleDelete = (id: string, name: string) => {
    showConfirm(
      '刪除輪詢群組',
      `確定要刪除「${name}」嗎？此操作無法復原，且會影響使用此群組的點位。`,
      async () => {
        try {
          await deleteMutation.mutateAsync(id);
          showSuccess('輪詢群組已成功刪除');
          closeConfirm();
        } catch (err: unknown) {
          const message = err instanceof Error ? err.message : 'Unknown error';
          showError(`刪除輪詢群組失敗: ${message}`);
          closeConfirm();
        }
      },
      'danger'
    );
  };

  /**
   * 開啟編輯 Modal
   */
  const handleEdit = (group: PollingGroup) => {
    setEditingGroup(group);
    setIsModalOpen(true);
  };

  /**
   * 開啟新增 Modal
   */
  const handleCreate = () => {
    setEditingGroup(null);
    setIsModalOpen(true);
  };

  /**
   * 提交表單（新增或更新）
   */
  const handleSubmit = async (
    data: CreatePollingGroupRequest | UpdatePollingGroupRequest
  ) => {
    if (editingGroup) {
      await updateMutation.mutateAsync({
        id: editingGroup.id,
        data: data as UpdatePollingGroupRequest,
      });
      showSuccess('輪詢群組已成功更新');
    } else {
      await createMutation.mutateAsync(data as CreatePollingGroupRequest);
      showSuccess('輪詢群組已成功建立');
    }
    setIsModalOpen(false);
  };

  return (
    <div className="space-y-8">
      {/* Header with Glow */}
      <div className="flex justify-between items-center relative z-0">
        <div className="absolute -left-10 -top-10 w-32 h-32 bg-blue-600/20 blur-3xl pointer-events-none"></div>
        <div>
          <h2 className="text-3xl font-bold text-slate-100 tracking-tight">輪詢群組</h2>
          <p className="text-slate-400 mt-1">管理數據採集的輪詢排程與優先級</p>
        </div>
        <button
          onClick={handleCreate}
          className="px-5 py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/20 hover:shadow-blue-500/40 flex items-center space-x-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          <span>新增群組</span>
        </button>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {/* Skeleton Loaders */}
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="bg-slate-800/50 rounded-xl p-6 border border-slate-700/50 h-32 animate-pulse"
            >
              <div className="flex space-x-3 mb-4">
                <div className="h-6 bg-slate-700 rounded w-1/3"></div>
                <div className="h-4 bg-slate-700/50 rounded w-1/4"></div>
              </div>
              <div className="h-4 bg-slate-700/50 rounded w-full mb-2"></div>
              <div className="h-4 bg-slate-700/50 rounded w-2/3"></div>
            </div>
          ))}
        </div>
      ) : error ? (
        <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-6 text-red-400 flex items-center space-x-3">
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{error instanceof Error ? error.message : '載入輪詢群組失敗'}</span>
        </div>
      ) : pollingGroups.length === 0 ? (
        <div
          className="text-center py-20 bg-slate-800/30 rounded-2xl border-2 border-slate-700/50 border-dashed group hover:border-blue-500/30 transition-colors cursor-pointer"
          onClick={handleCreate}
        >
          <div className="w-16 h-16 bg-slate-700/50 rounded-full mx-auto flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
            <svg
              className="w-8 h-8 text-slate-500 group-hover:text-blue-400 transition-colors"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 6v6m0 0v6m0-6h6m-6 0H6"
              />
            </svg>
          </div>
          <h3 className="text-xl font-bold text-slate-200 mb-2">尚未建立輪詢群組</h3>
          <p className="text-slate-500 max-w-sm mx-auto mb-6">
            建立輪詢群組來管理數據採集的頻率與優先級。
          </p>
          <button className="text-blue-400 hover:text-blue-300 font-medium hover:underline">
            建立第一個輪詢群組 &rarr;
          </button>
        </div>
      ) : (
        <div className="space-y-4">
          {pollingGroups.map((group) => (
            <div
              key={group.id}
              className="bg-slate-800/50 rounded-xl p-6 border border-slate-700/50 hover:border-slate-600 transition-colors"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center space-x-3 mb-2">
                    <h3 className="text-lg font-semibold text-slate-100">{group.name}</h3>
                    {group.enabled ? (
                      <span className="px-2 py-1 text-xs font-medium bg-green-500/20 text-green-400 rounded-full border border-green-500/30">
                        啟用
                      </span>
                    ) : (
                      <span className="px-2 py-1 text-xs font-medium bg-slate-500/20 text-slate-400 rounded-full border border-slate-500/30">
                        停用
                      </span>
                    )}
                  </div>
                  {group.description && (
                    <p className="text-slate-400 text-sm mb-4">{group.description}</p>
                  )}
                  <div className="flex items-center space-x-6 text-sm text-slate-500">
                    <div className="flex items-center space-x-2">
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                      <span>間隔: {formatInterval(group.interval_ms)}</span>
                    </div>
                    <div className="flex items-center space-x-2">
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M13 10V3L4 14h7v7l9-11h-7z"
                        />
                      </svg>
                      <span>優先級: {group.priority}</span>
                    </div>
                  </div>
                </div>
                <div className="flex items-center space-x-2 ml-4">
                  <button
                    onClick={() => handleEdit(group)}
                    className="p-2 text-slate-400 hover:text-blue-400 hover:bg-blue-500/10 rounded-lg transition-colors"
                    title="編輯"
                  >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                      />
                    </svg>
                  </button>
                  <button
                    onClick={() => handleDelete(group.id, group.name)}
                    className="p-2 text-slate-400 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors"
                    title="刪除"
                  >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* 確認對話框 */}
      <ConfirmDialog
        isOpen={confirmDialog.isOpen}
        title={confirmDialog.title}
        message={confirmDialog.message}
        variant={confirmDialog.variant || 'default'}
        onConfirm={confirmDialog.onConfirm}
        onCancel={closeConfirm}
      />

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-slate-900/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-slate-800 rounded-2xl border border-slate-700 p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto shadow-2xl ring-1 ring-white/10">
            <h3 className="text-xl font-bold text-slate-100 mb-6 pb-4 border-b border-slate-700 flex items-center justify-between">
              <span>{editingGroup ? '編輯輪詢群組' : '新增輪詢群組'}</span>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-slate-400 hover:text-white transition-colors"
              >
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </h3>

            <PollingGroupForm
              pollingGroup={editingGroup || undefined}
              onSubmit={handleSubmit}
              onCancel={() => setIsModalOpen(false)}
            />
          </div>
        </div>
      )}
    </div>
  );
}
