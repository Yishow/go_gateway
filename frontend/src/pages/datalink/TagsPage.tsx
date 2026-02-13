/**
 * TagsPage - 標籤字典管理頁面
 *
 * 使用 TanStack Query 進行 Server State 管理。
 */
import { useState, useMemo } from 'react';
import type { Tag, CreateTagRequest, UpdateTagRequest } from '../../types/datalink';
import TagTable from '../../components/datalink/TagTable';
import TagForm from '../../components/datalink/TagForm';
import ConfirmDialog from '../../components/datalink/ConfirmDialog';
import { useToast } from '../../contexts/ToastContext';
import {
  useTagsQuery,
  useCreateTagMutation,
  useUpdateTagMutation,
  useDeleteTagMutation,
  useActivateTagMutation,
  useRetireTagMutation,
} from '../../hooks/datalink';

export default function TagsPage() {
  // 篩選狀態（本地 UI 狀態）
  const [searchKey, setSearchKey] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('');
  const { showSuccess, showError } = useToast();

  // 建立篩選參數物件（useMemo 避免不必要的重新查詢）
  const filters = useMemo(
    () => ({
      key_prefix: searchKey || undefined,
      status: filterStatus || undefined,
    }),
    [searchKey, filterStatus]
  );

  // 使用 Query Hook 取得資料
  const { data: tags = [], isLoading, error } = useTagsQuery(filters);

  // Mutation Hooks
  const createMutation = useCreateTagMutation();
  const updateMutation = useUpdateTagMutation();
  const deleteMutation = useDeleteTagMutation();
  const activateMutation = useActivateTagMutation();
  const retireMutation = useRetireTagMutation();

  // Modal 狀態
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTag, setEditingTag] = useState<Tag | null>(null);
  
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
   * 刪除標籤
   */
  const handleDelete = (id: string) => {
    showConfirm(
      '刪除標籤',
      '確定要刪除此標籤嗎？此操作無法復原。',
      async () => {
        try {
          await deleteMutation.mutateAsync(id);
          showSuccess('標籤已成功刪除');
          closeConfirm();
        } catch (err: unknown) {
          const message = err instanceof Error ? err.message : 'Unknown error';
          showError(`刪除標籤失敗: ${message}`);
          closeConfirm();
        }
      },
      'danger'
    );
  };

  /**
   * 啟用標籤（樂觀更新）
   */
  const handleActivate = (id: string) => {
    activateMutation.mutate(id, {
      onSuccess: () => {
        showSuccess('標籤已成功啟用');
      },
      onError: (err) => {
        showError(`啟用標籤失敗: ${err.message}`);
      },
    });
  };

  /**
   * 退役標籤（樂觀更新）
   */
  const handleRetire = (id: string) => {
    showConfirm(
      '退役標籤',
      '確定要退役此標籤嗎？退役的標籤無法用於新的映射。',
      () => {
        retireMutation.mutate(id, {
          onSuccess: () => {
            showSuccess('標籤已成功退役');
            closeConfirm();
          },
          onError: (err) => {
            showError(`退役標籤失敗: ${err.message}`);
            closeConfirm();
          },
        });
      },
      'warning'
    );
  };

  /**
   * 開啟編輯 Modal
   */
  const handleEdit = (tag: Tag) => {
    setEditingTag(tag);
    setIsModalOpen(true);
  };

  /**
   * 開啟新增 Modal
   */
  const handleCreate = () => {
    setEditingTag(null);
    setIsModalOpen(true);
  };

  /**
   * 提交表單
   */
  const handleSubmit = async (data: CreateTagRequest | UpdateTagRequest) => {
    if (editingTag) {
      await updateMutation.mutateAsync({
        id: editingTag.id,
        data: data as UpdateTagRequest,
      });
      showSuccess('標籤已成功更新');
    } else {
      await createMutation.mutateAsync(data as CreateTagRequest);
      showSuccess('標籤已成功建立');
    }
    setIsModalOpen(false);
  };

  return (
    <div className="space-y-8">
      {/* Header with Glow */}
      <div className="flex justify-between items-center relative z-0">
        <div className="absolute -left-10 -top-10 w-32 h-32 bg-purple-600/20 blur-3xl pointer-events-none"></div>
        <div>
          <h2 className="text-3xl font-bold text-slate-100 tracking-tight">Tag Dictionary</h2>
          <p className="text-slate-400 mt-1">Manage global standard tags and metadata definitions</p>
        </div>
        <button
          onClick={handleCreate}
          className="px-5 py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/20 hover:shadow-blue-500/40 flex items-center space-x-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          <span>Create Tag</span>
        </button>
      </div>

      {/* Filters */}
      <div className="flex space-x-4 bg-slate-800/80 backdrop-blur-sm p-4 rounded-xl border border-slate-700 shadow-sm">
        <div className="flex-1 relative">
          <span className="absolute left-3 top-2.5 text-slate-500">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </span>
          <input
            type="text"
            placeholder="Search by tag keys..."
            value={searchKey}
            onChange={e => setSearchKey(e.target.value)}
            className="w-full bg-slate-900 border border-slate-700 rounded-lg pl-10 pr-4 py-2 text-slate-200 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all"
          />
        </div>
        <div>
          <select
            value={filterStatus}
            onChange={e => setFilterStatus(e.target.value)}
            className="bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-slate-200 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all cursor-pointer"
          >
            <option value="">All Status</option>
            <option value="draft">Draft</option>
            <option value="active">Active</option>
            <option value="retired">Retired</option>
          </select>
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3, 4, 5].map(i => (
            <div key={i} className="h-16 bg-slate-800/50 rounded-lg animate-pulse border border-slate-700/50"></div>
          ))}
        </div>
      ) : error ? (
        <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-6 text-red-400 flex items-center space-x-3">
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{error instanceof Error ? error.message : 'Failed to fetch tags'}</span>
        </div>
      ) : tags.length === 0 ? (
        <div
          className="text-center py-20 bg-slate-800/30 rounded-2xl border-2 border-slate-700/50 border-dashed group hover:border-purple-500/30 transition-colors cursor-pointer"
          onClick={handleCreate}
        >
          <div className="w-16 h-16 bg-slate-700/50 rounded-full mx-auto flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
            <svg className="w-8 h-8 text-slate-500 group-hover:text-purple-400 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
            </svg>
          </div>
          <h3 className="text-xl font-bold text-slate-200 mb-2">No tags found</h3>
          <p className="text-slate-500 max-w-sm mx-auto mb-6">
            Create your first standardized data tag to map device points to business logic.
          </p>
          <button className="text-purple-400 hover:text-purple-300 font-medium hover:underline">
            Create your first tag &rarr;
          </button>
        </div>
      ) : (
        <div className="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden shadow-sm">
          <TagTable
            tags={tags}
            onEdit={handleEdit}
            onDelete={handleDelete}
            onActivate={handleActivate}
            onRetire={handleRetire}
          />
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
              <span>{editingTag ? 'Edit Tag' : 'New Tag'}</span>
              <button onClick={() => setIsModalOpen(false)} className="text-slate-400 hover:text-white transition-colors">
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </h3>

            <TagForm
              tag={editingTag || undefined}
              onSubmit={handleSubmit}
              onCancel={() => setIsModalOpen(false)}
            />
          </div>
        </div>
      )}
    </div>
  );
}
