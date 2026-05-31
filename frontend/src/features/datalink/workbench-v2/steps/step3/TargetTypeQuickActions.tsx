import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { TargetType } from '../../state/types';

const QUICK_TARGET_TYPES: TargetType[] = ['bool', 'int16', 'float64', 'string'];

export interface TargetTypeQuickActionsProps {
  currentType: TargetType;
  onApplyAll: () => void;
  onSelect: (value: TargetType) => void;
}

export const TargetTypeQuickActions: React.FC<TargetTypeQuickActionsProps> = ({
  currentType,
  onApplyAll,
  onSelect,
}) => {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="space-y-3" data-testid="target-type-quick-actions">
      <div className="flex flex-wrap items-center gap-2">
        {QUICK_TARGET_TYPES.map((value) => {
          const active = value === currentType;
          return (
            <button
              key={value}
              type="button"
              onClick={() => onSelect(value)}
              className={`rounded-full border px-2.5 py-1 text-[11px] font-mono transition-colors ${
                active
                  ? 'border-blue-500/70 bg-blue-500/15 text-blue-200'
                  : 'border-slate-800 bg-slate-950/60 text-slate-400 hover:border-slate-700 hover:text-slate-200'
              }`}
              data-testid={`target-type-chip-${value}`}
            >
              {value}
            </button>
          );
        })}
      </div>
      <button
        type="button"
        onClick={onApplyAll}
        className="text-xs font-semibold text-blue-400 transition-colors hover:text-blue-300"
        data-testid="btn-apply-target-type-all"
      >
        {t('step3.preview.applyAllTargetType', {
          defaultValue: '將目前目標型態套用到全部列',
        })}
      </button>
    </div>
  );
};
