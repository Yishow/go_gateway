import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import type { DataType } from '../../../../../types/datalink';
import type { LayoutMode, MixedItem } from '../../../../../types/measurement';
import {
  SCATTERED_ITEM_DATA_TYPES,
  calculateScatteredRandomUsage,
  createDefaultScatteredItem,
  defaultSemanticKindForType,
  resolveScatteredRangeSupport,
  validateScatteredItems,
  type ScatteredDeviceLike,
} from '../../state/scatteredRange';

export interface ScatteredRangeSectionProps {
  layoutMode: LayoutMode;
  items: ReadonlyArray<MixedItem>;
  mixedItemsError?: string;
  device?: ScatteredDeviceLike | null;
  fallbackAddress?: string;
  fallbackDataType?: DataType;
  disabled?: boolean;
  onChangeMode: (nextMode: LayoutMode, nextItems?: MixedItem[]) => void;
  onChangeItems: (nextItems: MixedItem[]) => void;
}

export function ScatteredRangeSection({
  layoutMode,
  items,
  mixedItemsError,
  device,
  fallbackAddress,
  fallbackDataType = 'int16',
  disabled = false,
  onChangeMode,
  onChangeItems,
}: ScatteredRangeSectionProps) {
  const { t } = useTranslation('acquisition-settings');
  const sectionId = useId();
  const support = resolveScatteredRangeSupport(device);
  const isMixed = layoutMode === 'mixed';
  const usage = calculateScatteredRandomUsage(device?.protocol, items, support.effectiveLimit);
  const validation = validateScatteredItems(device?.protocol, items);

  const handleSelectMode = (nextMode: LayoutMode) => {
    if (disabled || nextMode === layoutMode) return;
    if (nextMode === 'mixed') {
      if (!support.supported) return;
      const nextItems =
        items.length > 0
          ? [...items]
          : [createDefaultScatteredItem(0, [], device?.protocol, fallbackAddress, fallbackDataType)];
      onChangeMode('mixed', nextItems);
      return;
    }
    onChangeMode('homogeneous');
  };

  const handleUpdateItem = (index: number, patch: Partial<MixedItem>) => {
    const next = items.map((item, idx) => {
      if (idx !== index) return item;
      const nextType = patch.data_type ?? item.data_type;
      return {
        ...item,
        ...patch,
        semantic_kind: patch.data_type ? defaultSemanticKindForType(nextType) : item.semantic_kind,
      };
    });
    onChangeItems(next);
  };

  const handleAddItem = () => {
    if (disabled) return;
    const newItem = createDefaultScatteredItem(items.length, items, device?.protocol, fallbackAddress, fallbackDataType);
    onChangeItems([...items, newItem]);
  };

  const handleRemoveItem = (index: number) => {
    if (disabled || items.length <= 1) return;
    onChangeItems(items.filter((_, idx) => idx !== index));
  };

  return (
    <div className="space-y-3" data-testid="scattered-range-section">
      <div className="gw-field">
        <div className="gw-label"><span>{t('scatteredRange.modeLabel')}</span></div>
        <div className="flex flex-wrap gap-2" role="group" aria-label={t('scatteredRange.modeLabel')}>
          <button
            type="button"
            data-testid="range-mode-homogeneous"
            disabled={disabled}
            aria-pressed={!isMixed}
            onClick={() => handleSelectMode('homogeneous')}
            className={`gw-btn gw-btn-sm ${!isMixed ? 'gw-btn-primary' : 'gw-btn-ghost'}`}
          >
            {t('scatteredRange.contiguousMode')}
          </button>
          <button
            type="button"
            data-testid="range-mode-mixed"
            disabled={disabled || !support.supported}
            aria-pressed={isMixed}
            onClick={() => handleSelectMode('mixed')}
            className={`gw-btn gw-btn-sm ${isMixed ? 'gw-btn-primary' : 'gw-btn-ghost'}`}
          >
            {t('scatteredRange.scatteredMode')}
          </button>
        </div>
        {!support.supported && support.disabledReason && (
          <p className="text-xs text-slate-500 mt-1" data-testid="scattered-mode-disabled-reason">
            {t(`scatteredRange.disabledReasons.${support.disabledReason}`, { profile: support.profile })}
          </p>
        )}
      </div>

      {isMixed && (
        <div className="space-y-3 rounded-lg border border-slate-200 dark:border-slate-700 p-3">
          {Boolean(mixedItemsError) && (
            <div
              role="alert"
              data-testid="scattered-corrupt-items-error"
              className="rounded border border-rose-300 bg-rose-50 px-2.5 py-1.5 text-xs text-rose-900 dark:border-rose-700 dark:bg-rose-950/40 dark:text-rose-200"
            >
              {t('scatteredRange.corruptItemsError')}
            </div>
          )}
          {support.effectiveLimit > 0 && (
            <div className="space-y-1.5" data-testid="scattered-usage-panel">
              <div className="flex items-center justify-between text-xs font-medium">
                <span>{t('scatteredRange.usageLabel')}</span>
                <span className="gw-mono" data-testid="scattered-usage-counter">
                  {t('scatteredRange.usageCounter', { used: usage.usedPoints, limit: usage.effectiveLimit })}
                </span>
              </div>
              <div className="h-1.5 w-full overflow-hidden rounded-full bg-slate-200 dark:bg-slate-700">
                <div
                  className={`h-full transition-all ${
                    usage.status === 'over_limit' ? 'bg-amber-600' : usage.status === 'warning' ? 'bg-amber-500' : 'bg-emerald-500'
                  }`}
                  style={{ width: `${Math.min(100, Math.round(usage.ratio * 100))}%` }}
                />
              </div>
              <p className="text-[11px] text-slate-500">
                {device?.protocol === 'fatek_fbs' ? t('scatteredRange.fatekRuleNote') : t('scatteredRange.mcRuleNote')}
              </p>
              {usage.status === 'warning' && (
                <div
                  role="status"
                  data-testid="scattered-usage-warning"
                  className="rounded border border-amber-300 bg-amber-50 px-2.5 py-1.5 text-xs text-amber-900 dark:border-amber-700 dark:bg-amber-950/40 dark:text-amber-200"
                >
                  {t('scatteredRange.warningNearLimit', { used: usage.usedPoints, limit: usage.effectiveLimit })}
                </div>
              )}
              {usage.status === 'over_limit' && (
                <div
                  role="status"
                  data-testid="scattered-usage-over-limit"
                  className="rounded border border-amber-400 bg-amber-50 px-2.5 py-1.5 text-xs text-amber-900 dark:border-amber-600 dark:bg-amber-950/50 dark:text-amber-200"
                >
                  {t('scatteredRange.overLimitSplit', { used: usage.usedPoints, limit: usage.effectiveLimit })}
                </div>
              )}
            </div>
          )}

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium">{t('scatteredRange.itemsTitle')}</span>
              <button
                type="button"
                data-testid="scattered-add-item"
                disabled={disabled}
                onClick={handleAddItem}
                className="gw-btn gw-btn-sm gw-btn-ghost"
              >
                + {t('scatteredRange.addItem')}
              </button>
            </div>

            <div className="space-y-2" data-testid="scattered-items-list">
              {items.length === 0 && (
                <p
                  role="alert"
                  data-testid="scattered-empty-items-error"
                  className="text-xs text-rose-600 dark:text-rose-400"
                >
                  {t('scatteredRange.emptyItemsError')}
                </p>
              )}
              {items.map((item, index) => {
                const itemId = (item.item_id ?? '').trim() || `item-${index + 1}`;
                const issue = validation.itemIssues[itemId];
                const errorId = `${sectionId}-item-error-${index}`;
                const isTypeError = issue?.code === 'bit_requires_bool' || issue?.code === 'unsupported_string';
                return (
                  <div key={itemId} className="space-y-1">
                    <div className="grid grid-cols-1 sm:grid-cols-12 gap-2 items-center">
                      <div className="sm:col-span-4">
                        <input
                          type="text"
                          aria-label={`${t('scatteredRange.addressHeader')} ${index + 1}`}
                          aria-invalid={Boolean(issue)}
                          aria-describedby={issue ? errorId : undefined}
                          data-testid={`scattered-item-address-${index}`}
                          disabled={disabled}
                          value={item.address}
                          placeholder="D0"
                          onChange={(e) =>
                            handleUpdateItem(index, { address: e.target.value.toUpperCase() })
                          }
                          className="gw-input gw-mono w-full"
                        />
                      </div>
                      <div className="sm:col-span-4">
                        <select
                          aria-label={`${t('scatteredRange.dataTypeHeader')} ${index + 1}`}
                          aria-invalid={Boolean(isTypeError)}
                          aria-describedby={issue ? errorId : undefined}
                          data-testid={`scattered-item-datatype-${index}`}
                          disabled={disabled}
                          value={item.data_type}
                          onChange={(e) =>
                            handleUpdateItem(index, {
                              data_type: e.target.value as MixedItem['data_type'],
                            })
                          }
                          className="gw-select gw-mono w-full"
                        >
                          {SCATTERED_ITEM_DATA_TYPES.map((dt) => (
                            <option key={dt} value={dt}>
                              {dt}
                            </option>
                          ))}
                        </select>
                      </div>
                      <div className="sm:col-span-3">
                        <input
                          type="text"
                          aria-label={`${t('scatteredRange.nameHeader')} ${index + 1}`}
                          data-testid={`scattered-item-name-${index}`}
                          disabled={disabled}
                          value={item.name ?? ''}
                          placeholder={t('scatteredRange.namePlaceholder')}
                          onChange={(e) => handleUpdateItem(index, { name: e.target.value })}
                          className="gw-input w-full"
                        />
                      </div>
                      <div className="sm:col-span-1 flex justify-end">
                        <button
                          type="button"
                          aria-label={`${t('scatteredRange.removeItem')} ${index + 1}`}
                          data-testid={`scattered-item-remove-${index}`}
                          disabled={disabled || items.length <= 1}
                          onClick={() => handleRemoveItem(index)}
                          className="gw-btn gw-btn-sm gw-btn-ghost"
                        >
                          ×
                        </button>
                      </div>
                    </div>
                    {issue && (
                      <p
                        id={errorId}
                        role="alert"
                        data-testid={`scattered-item-error-${index}`}
                        className="text-xs text-rose-600 dark:text-rose-400"
                      >
                        {t(`scatteredRange.itemErrors.${issue.code}`, {
                          address: issue.conflictWithAddress ?? '',
                        })}
                      </p>
                    )}
                  </div>
                );
              })}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
