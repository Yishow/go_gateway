import { useTranslation } from 'react-i18next';
import { KindSelector } from './KindSelector';
import { WriteStrategy } from './WriteStrategy';
import type { DbConnector, SettingsConnector } from '../../state/types';

/**
 * ConnectorSection 元件屬性
 */
interface ConnectorSectionProps {
  connector: DbConnector;
  connectors?: SettingsConnector[];
  onUpdateConnector: (patch: Partial<DbConnector>) => void;
  onKindChange: (kind: DbConnector['kind']) => void;
  onSelectConnector?: (connector: SettingsConnector) => void;
  disabled?: boolean;
}

/**
 * 資料庫連接器設定區塊元件
 * 落地需求：「Connector configuration form」之完整配置區塊
 */
export function ConnectorSection({
  connector,
  connectors,
  onUpdateConnector,
  onKindChange,
  onSelectConnector,
  disabled = false
}: ConnectorSectionProps) {
  const { t } = useTranslation('workbench-v2');

  const {
    kind,
    name,
    host,
    port,
    database,
    username,
    password,
    schema,
    table,
    write_mode,
    write_interval_seconds
  } = connector;

  // 判定當前資料庫是否需要特定欄位
  const isSqlite = kind === 'sqlite';
  const hasSchema = kind === 'postgres' || kind === 'sqlserver';

  return (
    <div className="space-y-6">
      {connectors && connectors.length > 0 && (
        <div className="flex items-center gap-2 text-xs bg-slate-900/60 p-2.5 rounded-lg border border-slate-800">
          <span className="text-slate-400 font-medium whitespace-nowrap">
            {t('step4.load_from_pool', '從連接器池載入：')}
          </span>
          <select
            aria-label={t('step4.load_from_pool', '從連接器池載入')}
            disabled={disabled}
            defaultValue=""
            onChange={(e) => {
              const found = connectors.find((c) => c.id === e.target.value);
              if (found && onSelectConnector) {
                onSelectConnector(found);
                e.target.value = '';
              }
            }}
            className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white text-xs focus:border-blue-500 outline-none"
          >
            <option value="" disabled>-- {t('step4.select_existing_connector', '選擇既有連線')} --</option>
            {connectors.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name} ({c.kind.toUpperCase()} - {c.database})
              </option>
            ))}
          </select>
        </div>
      )}

      <KindSelector
        value={kind}
        onChange={onKindChange}
        disabled={disabled}
      />


      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {t('step4.field_name', '連線名稱')}
          </label>
          <input
            type="text"
            aria-label={t('step4.field_name', '連線名稱')}
            value={name}
            disabled={disabled}
            onChange={(e) => onUpdateConnector({ name: e.target.value })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            placeholder="e.g. TimescaleDB Local"
          />
        </div>

        {!isSqlite && (
          <>
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-gray-400">
                {t('step4.field_host', '主機位址 (Host)')}
              </label>
              <input
                type="text"
                aria-label={t('step4.field_host', '主機位址 (Host)')}
                value={host}
                disabled={disabled}
                onChange={(e) => onUpdateConnector({ host: e.target.value })}
                className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                placeholder="127.0.0.1"
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-medium text-gray-400">
                {t('step4.field_port', '通訊埠 (Port)')}
              </label>
              <input
                type="number"
                aria-label={t('step4.field_port', '通訊埠 (Port)')}
                value={port}
                disabled={disabled}
                onChange={(e) => onUpdateConnector({ port: parseInt(e.target.value) || 0 })}
                className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                placeholder="5432"
              />
            </div>
          </>
        )}

        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {isSqlite ? t('step4.field_db_sqlite', '資料庫檔案路徑') : t('step4.field_database', '資料庫名稱')}
          </label>
          <input
            type="text"
            aria-label={isSqlite ? t('step4.field_db_sqlite', '資料庫檔案路徑') : t('step4.field_database', '資料庫名稱')}
            value={database}
            disabled={disabled}
            onChange={(e) => onUpdateConnector({ database: e.target.value })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            placeholder={isSqlite ? 'gateway.db' : 'gateway_metrics'}
          />
        </div>

        {!isSqlite && (
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('step4.field_username', '使用者名稱')}
            </label>
            <input
              type="text"
              aria-label={t('step4.field_username', '使用者名稱')}
              value={username}
              disabled={disabled}
              onChange={(e) => onUpdateConnector({ username: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              placeholder="postgres"
            />
          </div>
        )}

        {!isSqlite && (
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('step4.field_password', '密碼')}
            </label>
            <input
              type="password"
              aria-label={t('step4.field_password', '密碼')}
              value={password ?? ''}
              disabled={disabled}
              autoComplete="new-password"
              onChange={(e) => onUpdateConnector({ password: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              placeholder={t('step4.field_password_placeholder', '留空則沿用既有密碼')}
            />
          </div>
        )}

        {hasSchema && (
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('step4.field_schema', '綱要 (Schema)')}
            </label>
            <input
              type="text"
              aria-label={t('step4.field_schema', '綱要 (Schema)')}
              value={schema}
              disabled={disabled}
              onChange={(e) => onUpdateConnector({ schema: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              placeholder="public"
            />
          </div>
        )}

        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {t('step4.field_table', '資料表名稱')}
          </label>
          <input
            type="text"
            aria-label={t('step4.field_table', '資料表名稱')}
            value={table}
            disabled={disabled}
            onChange={(e) => onUpdateConnector({ table: e.target.value })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            placeholder="sensor_readings"
          />
        </div>
      </div>

      <WriteStrategy
        writeMode={write_mode}
        writeIntervalSeconds={write_interval_seconds}
        onWriteModeChange={(mode) => onUpdateConnector({ write_mode: mode })}
        onWriteIntervalChange={(seconds) => onUpdateConnector({ write_interval_seconds: seconds })}
        disabled={disabled}
      />
    </div>
  );
}
