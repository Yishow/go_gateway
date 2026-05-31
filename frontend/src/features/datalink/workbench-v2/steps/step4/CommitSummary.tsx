import { useTranslation } from 'react-i18next';
import type { DbConnector } from '../../state/types';

/**
 * CommitSummary 元件屬性
 */
interface CommitSummaryProps {
  deviceCount: number;
  ruleCount: number;
  pointCount: number;
  mappingCount: number;
  connector: DbConnector;
  enabledTargetCount: number;
  hasConflict: boolean;
  onActivate: () => void;
}

/**
 * 提交前設定摘要與啟動按鈕元件
 * 落地需求：「Commit sequence and animation」之觸發與摘要區塊
 */
export function CommitSummary({
  deviceCount,
  ruleCount,
  pointCount,
  mappingCount,
  connector,
  enabledTargetCount,
  hasConflict,
  onActivate,
}: CommitSummaryProps) {
  const { t } = useTranslation('workbench-v2');

  const { kind, schema, table } = connector;

  // 格式化資料庫寫入目標路徑
  const dbTargetText = kind === 'sqlite' || !schema
    ? `${kind} → ${table}`
    : `${kind} → ${schema}.${table}`;

  // 提交按鈕不可用條件：有衝突，或者啟用的寫入欄位為 0
  const isSubmitDisabled = hasConflict || enabledTargetCount === 0;

  const summaryItems = [
    { label: t('step4.summary_devices', '採集裝置'), value: `${deviceCount} 台` },
    { label: t('step4.summary_rules', '點位規則'), value: `${ruleCount} 組` },
    { label: t('step4.summary_points', '採集點位'), value: `${pointCount} 個` },
    { label: t('step4.summary_mappings', '點位映射'), value: `${mappingCount} 筆` },
    { label: t('step4.summary_database', '資料庫寫入'), value: dbTargetText, isDb: true }
  ];

  return (
    <div className="bg-gray-900/10 border border-gray-800 rounded-2xl p-6 flex flex-col justify-between h-full backdrop-blur-sm">
      <div className="space-y-6">
        {/* 標題 */}
        <div>
          <h3 className="text-lg font-semibold text-white">
            {t('step4.summary_title', '確認設定並第一次啟動')}
          </h3>
          <p className="text-xs text-gray-500 mt-1">
            {t('step4.summary_subtitle', 'Step 4 只會啟動目前合法、可用、尚未啟動的設備。')}
          </p>
        </div>

        {/* 摘要列 */}
        <div className="divide-y divide-gray-800/40 border border-gray-800/60 rounded-xl overflow-hidden bg-gray-950/40">
          {summaryItems.map((item, idx) => (
            <div key={idx} className="flex justify-between items-center p-3 text-xs">
              <span className="text-gray-400 select-none">{item.label}</span>
              <span className={`font-semibold ${item.isDb ? 'font-mono text-blue-400 max-w-[200px] truncate' : 'text-gray-200'}`}>
                {item.value}
              </span>
            </div>
          ))}
        </div>
      </div>

      <div className="mt-8 space-y-3">
        {/* 提示訊息 */}
        {enabledTargetCount === 0 && !hasConflict && (
          <p className="text-xs text-amber-500 text-center select-none">
            ⚠️ {t('step4.no_enabled_targets_warning', '尚未啟用任何資料表欄位寫入')}
          </p>
        )}

        {/* 寬版啟動按鈕 */}
        <button
          type="button"
          disabled={isSubmitDisabled}
          onClick={onActivate}
          className={`
            w-full py-3 px-4 rounded-xl font-medium text-sm text-center transition-all duration-200 flex items-center justify-center gap-2
            ${isSubmitDisabled
              ? 'bg-gray-800 text-gray-500 cursor-not-allowed border border-gray-800/20'
              : 'bg-blue-600 hover:bg-blue-500 text-white cursor-pointer shadow-lg shadow-blue-600/10 hover:shadow-blue-500/20 active:scale-[0.98]'
            }
          `}
        >
          <span>🚀</span>
          {t('step4.activate_btn', '第一次啟動設備')}
        </button>

        <p className="text-[10px] text-gray-500 text-center leading-normal select-none">
          {t('step4.activate_info', '系統會逐台啟動符合條件的設備，並保留每台成功或失敗結果。')}
        </p>
      </div>
    </div>
  );
}
