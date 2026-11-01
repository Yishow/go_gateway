import { useTranslation } from 'react-i18next';
import { ConnectorRow } from './ConnectorRow';
import type { SettingsConnector } from '../state/types';

/**
 * ConnectorPoolSection 元件屬性
 */
interface ConnectorPoolSectionProps {
  connectors: SettingsConnector[];
  onAddConnector: () => void;
  onUpdateConnector: (id: string, patch: Partial<SettingsConnector>) => void;
  onRemoveConnector: (id: string) => void;
  onTestConnector: (id: string) => void;
}

/**
 * 資料庫連接器池設定區塊元件
 * 落地需求：「Connector pool CRUD」與「Connector backend test」
 */
export function ConnectorPoolSection({
  connectors,
  onAddConnector,
  onUpdateConnector,
  onRemoveConnector,
  onTestConnector,
}: ConnectorPoolSectionProps) {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="space-y-6 bg-[var(--gw-surface)] border border-[var(--gw-border)] rounded-2xl p-6" data-testid="connector-pool-section">
      {/* 區塊標題與新增按鈕 */}
      <div className="flex flex-col sm:flex-row justify-between items-start gap-3">
        <div>
          <h3 className="text-lg font-semibold text-[var(--gw-text)]">
            {t('settings.connector_pool_title', '資料庫連接器池')}
          </h3>
          <p className="text-xs text-[var(--gw-muted)] mt-1">
            {t('settings.connector_pool_subtitle', '設定多個目標寫入資料庫實體。可用於分流、備援或者是不同點位的多目標寫入設定。')}
          </p>
        </div>

        {/* 新增連接器按鈕 */}
        <button
          type="button"
          onClick={onAddConnector}
          className="flex shrink-0 whitespace-nowrap items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-semibold bg-[var(--gw-primary)] hover:bg-[var(--gw-primary-hover)] text-[var(--gw-on-primary)] cursor-pointer transition-all"
        >
          <span>＋</span>
          <span>{t('settings.add_connector_btn', '新增連接器')}</span>
        </button>
      </div>

      {/* 連接器列表 */}
      <div className="space-y-3">
        {connectors.map((connector) => (
          <ConnectorRow
            key={connector.id}
            connector={connector}
            onUpdate={(patch) => onUpdateConnector(connector.id, patch)}
            onRemove={() => onRemoveConnector(connector.id)}
            onTest={() => onTestConnector(connector.id)}
          />
        ))}

        {/* 空狀態 */}
        {connectors.length === 0 && (
          <div className="flex flex-col items-center justify-center p-12 border border-dashed border-[var(--gw-border)] rounded-xl bg-[var(--gw-surface)] text-center space-y-3" data-testid="connector-pool-empty">
            <span className="text-3xl select-none" role="img" aria-label="empty">🔌</span>
            <div>
              <h4 className="font-semibold text-[var(--gw-muted)] text-sm">
                {t('settings.connector_pool_empty_title', '尚未建立任何資料庫連接器')}
              </h4>
              <p className="text-xs text-[var(--gw-muted)] mt-1">
                {t('settings.connector_pool_empty_desc', '點擊右上方「新增連接器」按鈕以新增一筆連線設定。')}
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
