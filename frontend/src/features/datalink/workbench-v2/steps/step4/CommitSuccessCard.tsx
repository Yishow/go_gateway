import { useTranslation } from 'react-i18next';
import type { StudioV2ActivationResponse } from '../../../../../types/studioV2Activation';
import { normalizeTypedEnvelope } from '../../../../../utils/safeJson';

/**
 * CommitSuccessCard 元件屬性
 */
interface CommitSuccessCardProps {
  response: StudioV2ActivationResponse;
  canContinue: boolean;
  onCommit: () => void;
  onReset: () => void;
}

/**
 * 提交完成成功展示卡片元件
 * 落地需求：「Commit completion card」之綠色完成卡片與導向按鈕
 */
export function CommitSuccessCard({
  response,
  canContinue,
  onCommit,
  onReset,
}: CommitSuccessCardProps) {
  const { t } = useTranslation('workbench-v2');
  const responseEnvelope = normalizeTypedEnvelope(response);

  return (
    <div className="bg-emerald-950/20 border border-emerald-500/30 rounded-2xl p-6 flex flex-col justify-between h-full backdrop-blur-sm shadow-[0_0_24px_rgba(16,185,129,0.05)]">
      <div className="space-y-6 flex-1">
        {/* 大綠勾 Icon */}
        <div className="space-y-2 text-center">
          <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full border border-emerald-500/30 bg-emerald-500/10 text-3xl text-emerald-400 shadow-[0_0_16px_rgba(16,185,129,0.15)]">
            ✓
          </div>
          <h3 className="text-xl font-bold text-emerald-300">
            {t('step4.results_title', '啟動結果')}
          </h3>
          <p className="text-sm text-emerald-400/80">
            {t('step4.results_subtitle', '每台設備的啟動結果會分開保留，不會整批回滾。')}
          </p>
        </div>

        {responseEnvelope.code && (
          <div
            data-testid="activation-empty-message"
            className="rounded-xl border border-amber-500/20 bg-amber-500/10 px-4 py-3 text-sm text-amber-100"
          >
            {t(`errors.${responseEnvelope.code}`, {
              defaultValue: t('errors.activation_failed', 'Workspace activation failed.'),
            })}
          </div>
        )}

        <div className="space-y-2">
          {response.results.length === 0 && (
            <div className="rounded-xl border border-emerald-500/20 bg-emerald-500/5 px-4 py-3 text-xs text-emerald-200">
              {t('step4.already_active_info', '工作區所有設備皆已在運行中，資料採集與儲存正常進行。')}
            </div>
          )}
          {response.results.map((result) => (
            <div
              key={result.device_id}
              data-testid={`activation-result-${result.device_id}`}
              className="rounded-xl border border-slate-800 bg-slate-950/40 px-4 py-3"
            >
              <div className="flex items-center justify-between gap-3">
                <span className="font-mono text-xs text-slate-200">{result.device_id}</span>
                <span
                  className={`rounded-full px-2 py-0.5 text-[11px] font-semibold ${
                    result.status === 'success'
                      ? 'bg-emerald-500/10 text-emerald-300'
                      : result.status === 'failed'
                        ? 'bg-rose-500/10 text-rose-300'
                        : 'bg-amber-500/10 text-amber-300'
                  }`}
                >
                  {result.status}
                </span>
              </div>
              <p className="mt-2 text-xs text-slate-400">
                {t(`step4.progress_status.${result.status}`, result.status)}
              </p>
            </div>
          ))}
        </div>
      </div>

      <div className="mt-8 space-y-3">
        {canContinue && (
          <button
            type="button"
            onClick={onCommit}
            className="w-full py-3 px-4 rounded-xl font-medium text-sm text-center bg-emerald-600 hover:bg-emerald-500 text-white cursor-pointer transition-all shadow-lg shadow-emerald-600/10 hover:shadow-emerald-500/20 active:scale-[0.98]"
          >
            {t('step4.go_to_dashboard_btn', '前往 Runtime Dashboard')}
          </button>
        )}

        <button
          type="button"
          onClick={onReset}
          className="w-full py-3 px-4 rounded-xl border border-slate-700 bg-slate-900/40 text-sm font-medium text-slate-200 transition-all hover:bg-slate-900/70"
        >
          {t('step4.reset_activation_btn', '重新檢查可啟動設備')}
        </button>

        <p className="text-[10px] text-gray-500 text-center select-none">
          {canContinue
            ? t('step4.success_info', '至少一台設備已成功啟動，可以前往 Runtime Dashboard 繼續觀察。')
            : t('step4.empty_info', '目前沒有成功啟動的設備，請回到前面步驟修正後再試一次。')}
        </p>
      </div>
    </div>
  );
}
