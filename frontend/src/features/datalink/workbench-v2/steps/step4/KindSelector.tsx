import { useTranslation } from 'react-i18next';
import type { DbConnector } from '../../state/types';

/**
 * KindSelector 元件屬性
 */
interface KindSelectorProps {
  value: DbConnector['kind'];
  onChange: (kind: DbConnector['kind']) => void;
  disabled?: boolean;
}

/**
 * 支援的資料庫選項定義
 */
interface DbOption {
  kind: DbConnector['kind'];
  label: string;
  emoji: string;
  desc: string;
}

/**
 * 選擇資料庫類型的卡片元件
 * 落地需求：「Connector configuration form」之 kind 選擇
 */
export function KindSelector({ value, onChange, disabled = false }: KindSelectorProps) {
  const { t } = useTranslation('workbench-v2');

  const options: DbOption[] = [
    { kind: 'postgres', label: 'PostgreSQL', emoji: '🐘', desc: 'PostgreSQL / TimescaleDB' },
    { kind: 'sqlite', label: 'SQLite', emoji: '💾', desc: 'SQLite Local File Database' },
    { kind: 'mysql', label: 'MySQL', emoji: '🐬', desc: 'MySQL / MariaDB' },
    { kind: 'sqlserver', label: 'SQL Server', emoji: '🖥️', desc: 'Microsoft SQL Server' }
  ];

  return (
    <div className="space-y-3">
      <label className="text-sm font-medium text-gray-400">
        {t('step4.connector_kind')}
      </label>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {options.map((opt) => {
          const isActive = value === opt.kind;
          return (
            <button
              key={opt.kind}
              type="button"
              disabled={disabled}
              onClick={() => onChange(opt.kind)}
              className={`
                flex flex-col items-center justify-center p-4 rounded-xl border text-center transition-all duration-200
                ${disabled 
                  ? 'opacity-50 cursor-not-allowed border-gray-800 bg-gray-900/20' 
                  : 'cursor-pointer'
                }
                ${isActive && !disabled
                  ? 'border-blue-500 bg-blue-500/10 shadow-[0_0_12px_rgba(59,130,246,0.25)] text-white'
                  : 'border-gray-800 bg-gray-900/40 hover:border-gray-700 hover:bg-gray-900/60 text-gray-400 hover:text-gray-200'
                }
              `}
            >
              <span className="text-3xl mb-2 select-none" role="img" aria-label={opt.label}>
                {opt.emoji}
              </span>
              <span className="font-semibold text-sm block">
                {t(`step4.kind_${opt.kind}`, opt.label)}
              </span>
              <span className="text-xs text-gray-500 mt-1 block leading-tight">
                {opt.desc}
              </span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
