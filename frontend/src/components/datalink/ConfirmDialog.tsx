/**
 * ConfirmDialog - 確認對話框組件
 * 
 * 替代 window.confirm，提供更好的 UI/UX
 */

interface ConfirmDialogProps {
  /** 是否顯示對話框 */
  isOpen: boolean;
  /** 標題 */
  title: string;
  /** 訊息內容 */
  message: string;
  /** 確認按鈕文字 */
  confirmText?: string;
  /** 取消按鈕文字 */
  cancelText?: string;
  /** 確認按鈕樣式（danger/warning/default） */
  variant?: 'danger' | 'warning' | 'default';
  /** 確認回調 */
  onConfirm: () => void;
  /** 取消回調 */
  onCancel: () => void;
}

export default function ConfirmDialog({
  isOpen,
  title,
  message,
  confirmText = '確認',
  cancelText = '取消',
  variant = 'default',
  onConfirm,
  onCancel,
}: ConfirmDialogProps) {
  if (!isOpen) return null;

  const variantStyles = {
    danger: 'bg-red-600 hover:bg-red-500 text-white',
    warning: 'bg-amber-600 hover:bg-amber-500 text-white',
    default: 'bg-blue-600 hover:bg-blue-500 text-white',
  };

  return (
    <div className="fixed inset-0 z-[10000] flex items-center justify-center p-4">
      {/* 背景遮罩 */}
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
        onClick={onCancel}
      />

      {/* 對話框 */}
      <div className="relative bg-slate-800 rounded-2xl border border-slate-700 shadow-2xl ring-1 ring-white/10 w-full max-w-md p-6 animate-in fade-in zoom-in duration-200">
        {/* 標題 */}
        <h3 className="text-xl font-bold text-slate-100 mb-3">{title}</h3>

        {/* 訊息 */}
        <p className="text-slate-300 mb-6 leading-relaxed">{message}</p>

        {/* 按鈕組 */}
        <div className="flex justify-end gap-3">
          <button
            onClick={onCancel}
            className="px-4 py-2 text-slate-300 hover:text-white transition-colors font-medium"
          >
            {cancelText}
          </button>
          <button
            onClick={onConfirm}
            className={`px-4 py-2 rounded-lg font-medium transition-colors ${variantStyles[variant]}`}
          >
            {confirmText}
          </button>
        </div>
      </div>
    </div>
  );
}
