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
            {t('step4.load_from_pool')}
          </span>
          <select
            aria-label={t('step4.load_from_pool')}
            disabled={disabled}
            defaultValue=""
            onChange={(e) => {
              const found = connectors.find((c) => c.id === e.target.value);
              if (found?.enabled && onSelectConnector) {
                onSelectConnector(found);
                e.target.value = '';
              }
            }}
            className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white text-xs focus:border-blue-500 outline-none"
          >
            <option value="" disabled>-- {t('step4.select_existing_connector')} --</option>
            {connectors.map((c) => (
              <option key={c.id} value={c.id} disabled={!c.enabled}>
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
            {t('step4.field_name')}
          </label>
          <input
            type="text"
            aria-label={t('step4.field_name')}
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
                {t('step4.field_host')}
              </label>
              <input
                type="text"
                aria-label={t('step4.field_host')}
                value={host}
                disabled={disabled}
                onChange={(e) => onUpdateConnector({ host: e.target.value })}
                className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                placeholder="127.0.0.1"
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-medium text-gray-400">
                {t('step4.field_port')}
              </label>
              <input
                type="number"
                aria-label={t('step4.field_port')}
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
            {isSqlite ? t('step4.field_db_sqlite') : t('step4.field_database')}
          </label>
          <input
            type="text"
            aria-label={isSqlite ? t('step4.field_db_sqlite') : t('step4.field_database')}
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
              {t('step4.field_username')}
            </label>
            <input
              type="text"
              aria-label={t('step4.field_username')}
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
              {t('step4.field_password')}
            </label>
            <input
              type="password"
              aria-label={t('step4.field_password')}
              value={password ?? ''}
              disabled={disabled}
              autoComplete="new-password"
              onChange={(e) => onUpdateConnector({ password: e.target.value, clear_password: false })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              placeholder={t('step4.field_password_placeholder')}
            />
            <label className="flex items-center gap-2 text-xs text-gray-400">
              <input type="checkbox" checked={connector.clear_password === true} disabled={disabled}
                onChange={(e) => onUpdateConnector({ clear_password: e.target.checked, password: '', password_required: false })} />
              {t('step4.clear_password')}
            </label>
          </div>
        )}

        {hasSchema && (
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('step4.field_schema')}
            </label>
            <input
              type="text"
              aria-label={t('step4.field_schema')}
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
            {t('step4.field_table')}
          </label>
          <input
            type="text"
            aria-label={t('step4.field_table')}
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
