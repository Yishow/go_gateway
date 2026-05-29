import { useTranslation } from 'react-i18next';

/**
 * SaveBar 元件屬性
 */
interface SaveBarProps {
  onReset: () => void;
  onSave: () => void;
}

/**
 * 系統設定底部 Sticky 儲存控制列元件
 * 落地需求：「Save bar actions」與「Sticky 儲存列定位」
 */
export function SaveBar({ onReset, onSave }: SaveBarProps) {
  const { t } = useTranslation('workbench-v2');

  const handleResetClick = () => {
    // 彈出原生 confirm 確認
    const ok = window.confirm(
      t('settings.reset_confirm_msg', '確定要將所有系統設定重設為預設值嗎？這將會清除您目前所有的連接器設定與變更。')
    );
    if (ok) {
      onReset();
    }
  };

  return (
    <div className="sticky bottom-4 z-10 flex flex-wrap items-center justify-between gap-4 p-4 rounded-2xl border border-blue-500/10 bg-slate-950/80 backdrop-blur-md shadow-xl shadow-black/20">
      {/* 提示說明文案 */}
      <div className="flex items-center gap-2">
        <span className="text-base select-none">💡</span>
        <span className="text-xs text-gray-400">
          {t('settings.save_bar_info', '設定變更會立即套用至本地運行，並於後端服務下次重啟後持續生效。')}
        </span>
      </div>

      {/* 按鈕操作區 */}
      <div className="flex items-center gap-3">
        {/* 重設為預設 */}
        <button
          type="button"
          onClick={handleResetClick}
          className="px-4 py-2 rounded-xl text-xs font-semibold border border-transparent hover:border-gray-800 text-gray-400 hover:text-gray-200 bg-transparent hover:bg-gray-900/40 cursor-pointer transition-all"
        >
          {t('settings.reset_to_defaults_btn', '重設為預設')}
        </button>

        {/* 儲存所有設定 */}
        <button
          type="button"
          onClick={onSave}
          className="px-5 py-2.5 rounded-xl text-xs font-semibold bg-emerald-600 hover:bg-emerald-500 text-white cursor-pointer transition-all shadow-md shadow-emerald-600/10 hover:shadow-emerald-500/20 active:scale-[0.98]"
        >
          {t('settings.save_all_btn', '儲存所有設定')}
        </button>
      </div>
    </div>
  );
}
