import { useTranslation } from 'react-i18next';
import { KindSelector } from './KindSelector';
import { WriteStrategy } from './WriteStrategy';
import { DeliveryTruthStrip } from './DeliveryTruthStrip';
import type { DbConnector } from '../../state/types';

/**
 * ConnectorSection 元件屬性
 */
interface ConnectorSectionProps {
  connector: DbConnector;
  onUpdateConnector: (patch: Partial<DbConnector>) => void;
  onKindChange: (kind: DbConnector['kind']) => void;
  disabled?: boolean;
}

/**
 * 資料庫連接器設定區塊元件
 * 落地需求：「Connector configuration form」之完整配置區塊
 */
export function ConnectorSection({
  connector,
  onUpdateConnector,
  onKindChange,
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
    <div className="space-y-6 bg-gray-900/10 border border-gray-800 rounded-2xl p-6 backdrop-blur-sm">
      {/* 區塊標題 */}
      <div>
        <h3 className="text-lg font-semibold text-white">
          {t('step4.connector_title', '資料庫連接器')}
        </h3>
        <p className="text-xs text-gray-500 mt-1">
          {t('step4.connector_subtitle', '設定目標資料寫入的資料庫連線參數與策略')}
        </p>
      </div>

      <DeliveryTruthStrip connector={connector} />

      {/* 1. 資料庫種類選擇 */}
      <KindSelector
        value={kind}
        onChange={onKindChange}
        disabled={disabled}
      />

      {/* 2. 連線欄位設定 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* 連線名稱 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {t('step4.field_name', '連線名稱')}
          </label>
          <input
            type="text"
            value={name}
            disabled={disabled}
            onChange={(e) => onUpdateConnector({ name: e.target.value })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            placeholder="e.g. TimescaleDB Local"
          />
        </div>

        {/* SQLite 不需要 Host, Port, Username, Schema */}
        {!isSqlite && (
          <>
            {/* Host */}
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-gray-400">
                {t('step4.field_host', '主機位址 (Host)')}
              </label>
              <input
                type="text"
                value={host}
                disabled={disabled}
                onChange={(e) => onUpdateConnector({ host: e.target.value })}
                className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                placeholder="127.0.0.1"
              />
            </div>

            {/* Port */}
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-gray-400">
                {t('step4.field_port', '通訊埠 (Port)')}
              </label>
              <input
                type="number"
                value={port}
                disabled={disabled}
                onChange={(e) => onUpdateConnector({ port: parseInt(e.target.value) || 0 })}
                className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                placeholder="5432"
              />
            </div>
          </>
        )}

        {/* Database 名稱 / SQLite 為檔案路徑 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {isSqlite ? t('step4.field_db_sqlite', '資料庫檔案路徑') : t('step4.field_database', '資料庫名稱')}
          </label>
          <input
            type="text"
            value={database}
            disabled={disabled}
            onChange={(e) => onUpdateConnector({ database: e.target.value })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            placeholder={isSqlite ? 'gateway.db' : 'gateway_metrics'}
          />
        </div>

        {/* SQLite 不需要 Username */}
        {!isSqlite && (
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('step4.field_username', '使用者名稱')}
            </label>
            <input
              type="text"
              value={username}
              disabled={disabled}
              onChange={(e) => onUpdateConnector({ username: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              placeholder="postgres"
            />
          </div>
        )}

        {/* SQLite 不需要密碼；留空表示沿用既有密碼（更新時不覆寫） */}
        {!isSqlite && (
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('step4.field_password', '密碼')}
            </label>
            <input
              type="password"
              value={password ?? ''}
              disabled={disabled}
              autoComplete="new-password"
              onChange={(e) => onUpdateConnector({ password: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              placeholder={t('step4.field_password_placeholder', '留空則沿用既有密碼')}
            />
          </div>
        )}

        {/* 只有 Postgres / SQL Server 需要 Schema */}
        {hasSchema && (
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('step4.field_schema', '綱要 (Schema)')}
            </label>
            <input
              type="text"
              value={schema}
              disabled={disabled}
              onChange={(e) => onUpdateConnector({ schema: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              placeholder="public"
            />
          </div>
        )}

        {/* 資料表名稱 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {t('step4.field_table', '資料表名稱')}
          </label>
          <input
            type="text"
            value={table}
            disabled={disabled}
            onChange={(e) => onUpdateConnector({ table: e.target.value })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            placeholder="sensor_readings"
          />
        </div>
      </div>

      {/* 3. 寫入策略 */}
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
