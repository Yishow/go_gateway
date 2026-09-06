import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { ConnectorSection } from './ConnectorSection';
import { DeliveryTruthStrip } from './DeliveryTruthStrip';
import { CollapsibleSupportCard } from './CollapsibleSupportCard';
import type { DbConnector, SettingsConnector } from '../../state/types';

interface Step4SupportPanelsProps {
  connector: DbConnector;
  connectors?: SettingsConnector[];
  onUpdateConnector: (patch: Partial<DbConnector>) => void;
  onKindChange: (kind: DbConnector['kind']) => void;
  onSelectConnector?: (connector: SettingsConnector) => void;
  tableSetup?: ReactNode;
  disabled?: boolean;
}

export function Step4SupportPanels({
  connector,
  connectors,
  onUpdateConnector,
  onKindChange,
  onSelectConnector,
  tableSetup,
  disabled = false,
}: Step4SupportPanelsProps) {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="grid grid-cols-1 gap-4 xl:grid-cols-2">
      <CollapsibleSupportCard
        title={t('step4.connector_title', '資料庫連接器')}
        subtitle={t('step4.connector_subtitle', '設定目標資料寫入的資料庫連線、schema/table 與寫入策略')}
        collapsedLabel={t('step4.connector_expand_btn', '展開連接器設定')}
        expandedLabel={t('step4.connector_collapse_btn', '收合連接器設定')}
        defaultOpen={true}
        summary={
          <div className="flex flex-wrap gap-2">
            <span className="rounded-full border border-slate-700 px-2 py-1 font-mono text-[11px] text-slate-200">
              {connector.kind}
            </span>
            <span className="rounded-full border border-slate-700 px-2 py-1 font-mono text-[11px] text-slate-200">
              {connector.database}
            </span>
            <span className="rounded-full border border-slate-700 px-2 py-1 font-mono text-[11px] text-slate-200">
              {connector.schema ? `${connector.schema}.${connector.table}` : connector.table}
            </span>
          </div>
        }
      >
        <ConnectorSection
          connector={connector}
          connectors={connectors}
          onUpdateConnector={onUpdateConnector}
          onKindChange={onKindChange}
          onSelectConnector={onSelectConnector}
          disabled={disabled}
        />

        {tableSetup ? <div className="mt-4 border-t border-slate-800 pt-4">{tableSetup}</div> : null}
      </CollapsibleSupportCard>

      <CollapsibleSupportCard
        title={t('step4.delivery_truth_title', '資料交付狀態')}
        subtitle={t('step4.delivery_truth_hint', '最近一次建表與寫入結果')}
        collapsedLabel={t('step4.delivery_expand_btn', '展開交付狀態')}
        expandedLabel={t('step4.delivery_collapse_btn', '收合交付狀態')}
        summary={
          <div className="flex flex-wrap gap-2">
            <span className="rounded-full border border-slate-700 px-2 py-1 text-[11px] text-slate-200">
              {connector.status}
            </span>
            <span className="rounded-full border border-slate-700 px-2 py-1 text-[11px] text-slate-200">
              {connector.last_write_status ?? 'unknown'}
            </span>
            <span className="rounded-full border border-slate-700 px-2 py-1 text-[11px] text-slate-200">
              {connector.last_flush_status ?? 'unknown'}
            </span>
          </div>
        }
      >
        <DeliveryTruthStrip connector={connector} />
      </CollapsibleSupportCard>
    </div>
  );
}
