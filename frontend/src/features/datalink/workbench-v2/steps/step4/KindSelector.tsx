import { ParameterHelp } from '../../components/ParameterHelp';
import { supportsWriteGroupKind } from './writeGroupCapabilities';
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
      <label className="text-sm font-medium text-[var(--gw-muted)]">
        {t('step4.connector_kind')}
      </label><ParameterHelp help="db.kind" />
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {options.map((opt) => {
          const isActive = value === opt.kind;
          const unavailable = !supportsWriteGroupKind(opt.kind);
          return (
            <button
              key={opt.kind}
              type="button"
              disabled={disabled || unavailable}
              onClick={() => { if (!unavailable) onChange(opt.kind); }}
              className={`
                flex flex-col items-center justify-center p-4 rounded-xl border text-center transition-all duration-200
                ${disabled
                  ? 'opacity-50 cursor-not-allowed border-[var(--gw-control)] bg-[var(--gw-surface)]'
                  : 'cursor-pointer'
                }
                ${isActive && !disabled
                  ? 'border-[var(--gw-focus)] bg-[var(--gw-active)] shadow-[0_0_12px_rgba(59,130,246,0.25)] text-[var(--gw-text)]'
                  : 'border-[var(--gw-control)] bg-[var(--gw-surface)] hover:border-[var(--gw-control)] hover:bg-[var(--gw-surface)] text-[var(--gw-muted)] hover:text-[var(--gw-text)]'
                }
              `}
            >
              <span className="text-3xl mb-2 select-none" role="img" aria-label={opt.label}>
                {opt.emoji}
              </span>
              <span className="font-semibold text-sm block">
                {t(`step4.kind_${opt.kind}`, opt.label)}
              </span>
              <span className="text-xs text-[var(--gw-muted)] mt-1 block leading-tight">
                {unavailable ? t('step4.group_kind_unsupported') : opt.desc}
              </span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
