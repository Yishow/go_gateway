/**
 * MappingsPage - 資料映射管理頁面
 *
 * 使用 TanStack Query 進行 Server State 管理。
 */
import { useState } from 'react';
import { Link } from 'react-router-dom';
import type { Mapping, CreateMappingRequest, UpdateMappingRequest } from '../../types/datalink';
import MappingCard from '../../components/datalink/MappingCard';
import MappingCanvas from '../../components/datalink/MappingCanvas';
import ConfirmDialog from '../../components/datalink/ConfirmDialog';
import { useToast } from '../../contexts/ToastContext';
import {
  useMappingsQuery,
  useCreateMappingMutation,
  useUpdateMappingMutation,
  useDeleteMappingMutation,
  useToggleMappingStatusMutation,
} from '../../hooks/datalink';

export default function MappingsPage() {
  // 使用 Query Hook 取得資料
  const { data: mappings = [], isLoading, error } = useMappingsQuery();
  const { showSuccess, showError } = useToast();

  // Mutation Hooks
  const createMutation = useCreateMappingMutation();
  const updateMutation = useUpdateMappingMutation();
  const deleteMutation = useDeleteMappingMutation();
  const toggleStatusMutation = useToggleMappingStatusMutation();

  // Modal 狀態
  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const [editingMapping, setEditingMapping] = useState<Mapping | null>(null);
  
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
   * 刪除映射
   */
  const handleDelete = (id: string) => {
    showConfirm(
      '刪除映射',
      '確定要刪除此映射嗎？此操作無法復原。',
      async () => {
        try {
          await deleteMutation.mutateAsync(id);
          showSuccess('映射已成功刪除');
          closeConfirm();
        } catch (err: unknown) {
          const message = err instanceof Error ? err.message : 'Unknown error';
          showError(`刪除映射失敗: ${message}`);
          closeConfirm();
        }
      },
      'danger'
    );
  };

  /**
   * 切換映射狀態（樂觀更新）
   */
  const handleToggleStatus = (id: string, enabled: boolean) => {
    const action = enabled ? '啟用' : '停用';
    
    toggleStatusMutation.mutate(
      { id, enabled },
      {
        onSuccess: () => {
          showSuccess(`映射已成功${action}`);
        },
        onError: (err) => {
          showError(`${action}映射失敗: ${err.message}`);
        },
      }
    );
  };

  /**
   * 開啟編輯 Modal
   */
  const handleEdit = (mapping: Mapping) => {
    setEditingMapping(mapping);
    setIsEditorOpen(true);
  };

  /**
   * 開啟新增 Modal
   */
  const handleCreate = () => {
    setEditingMapping(null);
    setIsEditorOpen(true);
  };

  /**
   * 儲存映射
   */
  const handleSave = async (data: CreateMappingRequest | UpdateMappingRequest) => {
    try {
      if (editingMapping) {
        await updateMutation.mutateAsync({
          id: editingMapping.id,
          data: data as UpdateMappingRequest,
        });
        showSuccess('映射已成功更新');
      } else {
        await createMutation.mutateAsync(data as CreateMappingRequest);
        showSuccess('映射已成功建立');
      }
      setIsEditorOpen(false);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      showError(`儲存映射失敗: ${message}`);
    }
  };

  return (
    <div className="space-y-8">
      {/* Header with Glow */}
      <div className="flex justify-between items-center relative z-0">
        <div className="absolute -left-10 -top-10 w-32 h-32 bg-emerald-600/20 blur-3xl pointer-events-none"></div>
        <div>
          <h2 className="text-3xl font-bold text-slate-100 tracking-tight">Data Mappings</h2>
          <p className="text-slate-400 mt-1">Design data transformation pipelines from Source to Destination</p>
        </div>
        <div className="flex items-center space-x-3">
          <Link
            to="/datalink/wizard"
            className="px-5 py-2.5 bg-slate-700 hover:bg-slate-600 text-white rounded-xl font-medium transition-all border border-slate-600 hover:border-slate-500 flex items-center space-x-2"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
            </svg>
            <span>Wizard</span>
          </Link>
          <button
            onClick={handleCreate}
            className="px-5 py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/20 hover:shadow-blue-500/40 flex items-center space-x-2"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            <span>New Mapping</span>
          </button>
        </div>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3].map(i => (
            <div key={i} className="h-40 bg-slate-800/50 rounded-xl border border-slate-700/50 animate-pulse"></div>
          ))}
        </div>
      ) : error ? (
        <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-6 text-red-400 flex items-center space-x-3">
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{error instanceof Error ? error.message : 'Failed to fetch mappings'}</span>
        </div>
      ) : mappings.length === 0 ? (
        <div
          className="text-center py-20 bg-slate-800/30 rounded-2xl border-2 border-slate-700/50 border-dashed group hover:border-emerald-500/30 transition-colors cursor-pointer"
          onClick={handleCreate}
        >
          <div className="w-16 h-16 bg-slate-700/50 rounded-full mx-auto flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
            <svg className="w-8 h-8 text-slate-500 group-hover:text-emerald-400 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
            </svg>
          </div>
          <h3 className="text-xl font-bold text-slate-200 mb-2">No mappings found</h3>
          <p className="text-slate-500 max-w-sm mx-auto mb-6">
            Create your first mapping pipeline to connect device data points to standardized tags.
          </p>
          <button className="text-emerald-400 hover:text-emerald-300 font-medium hover:underline">
            Start mapping &rarr;
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {mappings.map(mapping => (
            <MappingCard
              key={mapping.id}
              mapping={mapping}
              onEdit={handleEdit}
              onDelete={handleDelete}
              onToggleStatus={handleToggleStatus}
            />
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

      {/* Editor Modal - Full Screen Overlay */}
      {isEditorOpen && (
        <div className="fixed inset-0 bg-slate-900 z-50 overflow-hidden flex flex-col animate-in fade-in duration-200">
          <MappingCanvas
            initialMapping={editingMapping}
            onSave={handleSave}
            onCancel={() => setIsEditorOpen(false)}
          />
        </div>
      )}
    </div>
  );
}
