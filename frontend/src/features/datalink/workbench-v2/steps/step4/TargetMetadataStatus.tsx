import { useTranslation } from 'react-i18next';
import type { Step4TargetMetadataStatus } from './useStep4TargetColumns';

const DEFAULT_MESSAGES: Record<Step4TargetMetadataStatus, string> = {
  not_checked: '儲存資料庫連線後才會查詢實際欄位。',
  checking: '正在查詢目標資料表的實際欄位…',
  exists: '已讀取目標資料表的實際欄位。',
  missing: '目標資料表尚不存在，請先預覽並確認建立。',
  forbidden: '此帳號沒有查詢目標資料表的權限，請調整權限後重新查詢。',
  failed: '無法查詢目標資料表，請確認連線後重新查詢。',
};

interface TargetMetadataStatusProps {
  status: Step4TargetMetadataStatus;
  onRetry: () => void;
}

/** 顯示目標資料表實際欄位的查詢狀態；失敗與無權限分開，並提供重新查詢。 */
export function TargetMetadataStatus({ status, onRetry }: TargetMetadataStatusProps) {
  const { t } = useTranslation('workbench-v2');
  const canRetry = status === 'failed' || status === 'forbidden';
  const tone = status === 'exists'
    ? 'text-emerald-400'
    : status === 'checking' || status === 'not_checked' ? 'text-slate-400' : 'text-amber-300';

  return (
    <div role="status" data-testid="step4-target-metadata-status" data-inspection-status={status} className={`text-xs ${tone}`}>
      {t(`step4.target_metadata_${status}`, { defaultValue: DEFAULT_MESSAGES[status] })}
      {canRetry && (
        <button type="button" className="ml-2 underline" onClick={onRetry} data-testid="step4-target-metadata-retry">
          {t('step4.target_metadata_retry')}
        </button>
      )}
    </div>
  );
}
