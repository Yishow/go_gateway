import { useTranslation } from 'react-i18next';

/**
 * WriteStrategy 元件屬性
 */
interface WriteStrategyProps {
  writeMode: 'insert' | 'upsert';
  writeIntervalSeconds: number;
  onWriteModeChange: (mode: 'insert' | 'upsert') => void;
  onWriteIntervalChange: (seconds: number) => void;
  disabled?: boolean;
}

/**
 * 寫入策略設定元件
 * 落地需求：「Configure Database Write」之策略選擇與寫入時間間隔
 */
export function WriteStrategy({
  writeMode,
  writeIntervalSeconds,
  onWriteModeChange,
  onWriteIntervalChange,
  disabled = false
}: WriteStrategyProps) {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-4 rounded-xl border border-gray-800/60 bg-gray-900/20">
      {/* 寫入模式選擇 (INSERT / UPSERT) */}
      <div className="space-y-3">
        <label className="text-sm font-medium text-gray-400">
          {t('step4.write_mode_label', '寫入策略')}
        </label>
        <div className="flex gap-4">
          <label
            className={`
              flex items-center gap-3 px-4 py-3 rounded-lg border text-sm font-medium cursor-pointer transition-all flex-1
              ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
              ${writeMode === 'insert'
                ? 'border-blue-500 bg-blue-500/5 text-white'
                : 'border-gray-800 bg-gray-900/40 text-gray-400 hover:border-gray-700 hover:text-gray-200'
              }
            `}
          >
            <input
              type="radio"
              name="write_mode"
              checked={writeMode === 'insert'}
              disabled={disabled}
              onChange={() => onWriteModeChange('insert')}
              className="text-blue-500 focus:ring-blue-500 bg-gray-900 border-gray-800 cursor-pointer disabled:cursor-not-allowed"
            />
            <div className="flex flex-col">
              <span>{t('step4.write_mode_insert', '時序追加 (INSERT)')}</span>
              <span className="text-xs text-gray-500 font-normal">
                {t('step4.write_mode_insert_desc', '每一筆資料均新增一列，適合歷史趨勢')}
              </span>
            </div>
          </label>

          <label
            className={`
              flex items-center gap-3 px-4 py-3 rounded-lg border text-sm font-medium cursor-pointer transition-all flex-1
              ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
              ${writeMode === 'upsert'
                ? 'border-blue-500 bg-blue-500/5 text-white'
                : 'border-gray-800 bg-gray-900/40 text-gray-400 hover:border-gray-700 hover:text-gray-200'
              }
            `}
          >
            <input
              type="radio"
              name="write_mode"
              checked={writeMode === 'upsert'}
              disabled={disabled}
              onChange={() => onWriteModeChange('upsert')}
              className="text-blue-500 focus:ring-blue-500 bg-gray-900 border-gray-800 cursor-pointer disabled:cursor-not-allowed"
            />
            <div className="flex flex-col">
              <span>{t('step4.write_mode_upsert', '覆寫更新 (UPSERT)')}</span>
              <span className="text-xs text-gray-500 font-normal">
                {t('step4.write_mode_upsert_desc', '僅保留最新時間戳資料，節省空間')}
              </span>
            </div>
          </label>
        </div>
      </div>

      {/* 寫入間隔設定 */}
      <div className="space-y-3 flex flex-col justify-between">
        <div>
          <label className="text-sm font-medium text-gray-400 block mb-1">
            {t('step4.write_interval_label', '寫入間隔')}
          </label>
          <p className="text-xs text-gray-500">
            {t('step4.write_interval_desc', '設定資料庫批次排程寫入的週期時間')}
          </p>
        </div>
        <div className="flex items-center gap-3">
          <div className="relative flex-1">
            <input
              type="number"
              min={1}
              max={3600}
              value={writeIntervalSeconds}
              disabled={disabled}
              onChange={(e) => onWriteIntervalChange(Math.max(1, parseInt(e.target.value) || 1))}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2.5 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            />
            <span className="absolute right-3 top-2.5 text-sm text-gray-500 select-none">
              {t('step4.seconds_unit', '秒 (s)')}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
