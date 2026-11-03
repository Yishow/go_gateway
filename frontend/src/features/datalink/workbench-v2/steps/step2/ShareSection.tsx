import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { Field } from '../../components/Field';
import { Input } from '../../components/inputs';
import { Toggle } from '../../components/Toggle';

export interface RuleShareControls { enabled: boolean | undefined; retry: () => void; openSettings: () => void }
export interface ShareSectionProps {
  shareControls?: RuleShareControls;
  ruleId: string;
  shareEnabled: boolean;
  shareStartRegister: number | null;
  shareStride: number | null;
  globalShareEnabled: boolean;
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

/**
 * Modbus Share 位址配置區塊元件
 * 
 * 落地設計決策：「Modbus Share 位址布局：自動接續 + 手動覆寫」
 * 提供展開式 details 區塊，當全域未啟用時顯示提示標籤。
 * 支援手動設定起點暫存器與 stride，並在手動時提供「改回自動分配」的重設重設連結。
 */
export const ShareSection: React.FC<ShareSectionProps> = ({
  ruleId,
  shareEnabled,
  shareStartRegister,
  shareStride,
  globalShareEnabled: inheritedShareEnabled,
  shareControls,
  dispatch,
}) => {
  const { t } = useTranslation('workbench-v2');
  const { t: sourceText } = useTranslation('recording-home');
  const globalShareEnabled = shareControls ? shareControls.enabled === true : inheritedShareEnabled;
  const globalShareKnown = !shareControls || shareControls.enabled !== undefined;

  return (
    <details className="group border-t border-slate-800/60 pt-4" data-testid="share-section">
      <summary className="flex items-center gap-2 cursor-pointer font-medium text-xs text-slate-400 select-none hover:text-slate-200 py-1.5 focus:outline-none">
        <svg
          viewBox="0 0 24 24"
          className="w-4 h-4 transition-transform group-open:rotate-90 text-slate-500"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.75"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M9 18l6-6-6-6" />
        </svg>
        <span>{t('step2.labels.share_details', 'Modbus Share 轉發配置 (選填)')}</span>

        {globalShareKnown && !globalShareEnabled && (
          <span
            className="ml-2 bg-slate-800 border border-slate-700/60 text-slate-500 text-[9px] font-sans font-normal px-1.5 py-0.5 rounded"
            data-testid="global-disabled-chip"
          >
            {t('step2.labels.global_disabled', '全域未啟用')}
          </span>
        )}
      </summary>

      <div className="space-y-4 pt-3 pl-6">
        {shareControls && !globalShareEnabled && <div role="status">
          <p>{sourceText(globalShareKnown ? 'source.shareOff' : 'source.shareUnknown')}</p>
          <button type="button" onClick={globalShareKnown ? shareControls.openSettings : shareControls.retry}>{sourceText(globalShareKnown ? 'source.openShareSettings' : 'source.retryShare')}</button>
        </div>}
        <div>
          <Toggle
            checked={shareEnabled && globalShareEnabled}
            disabled={!globalShareEnabled}
            onChange={() => dispatch({ type: 'toggleRuleShareEnabled', ruleId })}
            label="啟用此規則的 Modbus Share 記憶體對應"
            data-testid="share-enable-toggle"
          />
        </div>

        {shareEnabled && globalShareEnabled && (
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <Field
              label="轉發起點暫存器 (Share Start Register)"
              hint="留空則自動接續上一條規則"
            >
              <div className="relative">
                <Input
                  type="number"
                  placeholder="(自動)"
                  value={shareStartRegister ?? ''}
                  onChange={(e) => {
                    const val = parseInt(e.target.value, 10);
                    dispatch({
                      type: 'updateRuleShareStart',
                      ruleId,
                      shareStart: isNaN(val) ? null : val,
                    });
                  }}
                  data-testid="share-start-input"
                />
                {shareStartRegister !== null && (
                  <button
                    type="button"
                    onClick={() =>
                      dispatch({ type: 'updateRuleShareStart', ruleId, shareStart: null })
                    }
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] text-blue-400 hover:text-blue-300"
                    data-testid="share-reset-link"
                  >
                    改回自動分配
                  </button>
                )}
              </div>
            </Field>

            <Field
              label="點位間距 (Share Stride)"
              hint="留空則預設使用點位型別寬度"
            >
              <Input
                type="number"
                placeholder="(自動)"
                value={shareStride ?? ''}
                onChange={(e) => {
                  const val = parseInt(e.target.value, 10);
                  dispatch({
                    type: 'updateRuleShareStride',
                    ruleId,
                    shareStride: isNaN(val) ? null : val,
                  });
                }}
                data-testid="share-stride-input"
              />
            </Field>
          </div>
        )}
      </div>
    </details>
  );
};
