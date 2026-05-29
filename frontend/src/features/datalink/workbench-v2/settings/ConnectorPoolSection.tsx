import { useTranslation } from 'react-i18next';
import { ConnectorRow } from './ConnectorRow';
import type { SettingsConnector } from '../state/types';
import type { WorkbenchV2Action } from '../state/useWorkbenchV2State';

/**
 * ConnectorPoolSection 元件屬性
 */
interface ConnectorPoolSectionProps {
  connectors: SettingsConnector[];
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

/**
 * 資料庫連接器池設定區塊元件
 * 落地需求：「Connector pool CRUD」與「Connector mock test」
 */
export function ConnectorPoolSection({ connectors, dispatch }: ConnectorPoolSectionProps) {
  const { t } = useTranslation('workbench-v2');

  // 新增連接器 Handler
  const handleAddConnector = () => {
    dispatch({ type: 'addConnector' });
  };

  // 測試連接器連線 Handler (900ms 延遲模擬動畫)
  const handleTestConnector = (id: string) => {
    dispatch({ type: 'startConnectorTest', id });
    
    setTimeout(() => {
      // 85% 成功機率，15% 失敗機率
      const success = Math.random() > 0.15;
      
      dispatch({
        type: 'completeConnectorTest',
        id,
        result: success
          ? {
              status: 'ready',
              last_check_at: new Date().toISOString()
            }
          : {
              status: 'unreachable',
              last_check_at: new Date().toISOString(),
              last_check_error: 'connection refused'
            }
      });
    }, 900);
  };

  return (
    <div className="space-y-6 bg-gray-900/10 border border-gray-800 rounded-2xl p-6 backdrop-blur-sm" data-testid="connector-pool-section">
      {/* 區塊標題與新增按鈕 */}
      <div className="flex justify-between items-start">
        <div>
          <h3 className="text-lg font-semibold text-white">
            {t('settings.connector_pool_title', '資料庫連接器池')}
          </h3>
          <p className="text-xs text-gray-500 mt-1">
            {t('settings.connector_pool_subtitle', '設定多個目標寫入資料庫實體。可用於分流、備援或者是不同點位的多目標寫入設定。')}
          </p>
        </div>

        {/* 新增連接器按鈕 */}
        <button
          type="button"
          onClick={handleAddConnector}
          className="flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-500 text-white cursor-pointer transition-all shadow-md shadow-blue-600/10 hover:shadow-blue-500/20 active:scale-[0.98]"
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
            onUpdate={(patch) => dispatch({ type: 'updateConnector', id: connector.id, patch })}
            onRemove={() => dispatch({ type: 'removeConnector', id: connector.id })}
            onTest={() => handleTestConnector(connector.id)}
          />
        ))}

        {/* 空狀態 */}
        {connectors.length === 0 && (
          <div className="flex flex-col items-center justify-center p-12 border border-dashed border-gray-800 rounded-xl bg-gray-950/10 text-center space-y-3" data-testid="connector-pool-empty">
            <span className="text-3xl select-none" role="img" aria-label="empty">🔌</span>
            <div>
              <h4 className="font-semibold text-gray-400 text-sm">
                {t('settings.connector_pool_empty_title', '尚未建立任何資料庫連接器')}
              </h4>
              <p className="text-xs text-gray-600 mt-1">
                {t('settings.connector_pool_empty_desc', '點擊右上方「新增連接器」按鈕以新增一筆連線設定。')}
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
