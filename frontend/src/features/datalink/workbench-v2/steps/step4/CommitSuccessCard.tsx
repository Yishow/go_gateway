import { useTranslation } from 'react-i18next';

/**
 * CommitSuccessCard 元件屬性
 */
interface CommitSuccessCardProps {
  writeIntervalSeconds: number;
  onCommit: () => void;
}

/**
 * 提交完成成功展示卡片元件
 * 落地需求：「Commit completion card」之綠色完成卡片與導向按鈕
 */
export function CommitSuccessCard({ writeIntervalSeconds, onCommit }: CommitSuccessCardProps) {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="bg-emerald-950/20 border border-emerald-500/30 rounded-2xl p-6 flex flex-col justify-between h-full backdrop-blur-sm shadow-[0_0_24px_rgba(16,185,129,0.05)]">
      <div className="space-y-6 text-center flex-1 flex flex-col justify-center items-center">
        {/* 大綠勾 Icon */}
        <div className="w-16 h-16 rounded-full bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center text-emerald-400 text-3xl shadow-[0_0_16px_rgba(16,185,129,0.15)] animate-bounce select-none">
          ✓
        </div>

        {/* 標題與說明 */}
        <div className="space-y-2">
          <h3 className="text-xl font-bold text-emerald-300">
            {t('step4.success_title', '設定已套用')}
          </h3>
          <p className="text-sm text-emerald-400/80">
            {t('step4.success_subtitle', '採集管道與寫入任務已建立成功')}
          </p>
          <p className="text-xs text-gray-400 leading-normal max-w-xs mx-auto">
            {t('step4.scheduler_started_subline', 'Scheduler 已啟動 · 第一筆資料預計在 ~{{interval}}s 後寫入', { interval: writeIntervalSeconds })}
          </p>
        </div>
      </div>

      <div className="mt-8 space-y-3">
        {/* 前往 Runtime Dashboard 按鈕 */}
        <button
          type="button"
          onClick={onCommit}
          className="w-full py-3 px-4 rounded-xl font-medium text-sm text-center bg-emerald-600 hover:bg-emerald-500 text-white cursor-pointer transition-all shadow-lg shadow-emerald-600/10 hover:shadow-emerald-500/20 active:scale-[0.98]"
        >
          {t('step4.go_to_dashboard_btn', '前往 Runtime Dashboard')}
        </button>

        <p className="text-[10px] text-gray-500 text-center select-none">
          {t('step4.success_info', '您可以前往運行儀表板即時監看採集延遲、寫入吞吐量與錯誤日誌。')}
        </p>
      </div>
    </div>
  );
}
