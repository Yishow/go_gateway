import { useTranslation } from 'react-i18next';
import type { SettingsConnector } from '../state/types';

/**
 * ConnectorRow 元件屬性
 */
interface ConnectorRowProps {
  connector: SettingsConnector;
  onUpdate: (patch: Partial<SettingsConnector>) => void;
  onRemove: () => void;
  onTest: () => void;
}

/**
 * 單一資料庫連接器行編輯與狀態展示元件
 * 落地需求：「Connector pool CRUD」之 Row 部分
 */
export function ConnectorRow({ connector, onUpdate, onRemove, onTest }: ConnectorRowProps) {
  const { t } = useTranslation('workbench-v2');

  const {
    id,
    name,
    kind,
    host,
    port,
    database,
    username,
    password,
    schema,
    table,
    status,
    last_check_at,
    last_check_error
  } = connector;

  const isSqlite = kind === 'sqlite';
  const hasSchema = kind === 'postgres' || kind === 'sqlserver';
  const isTesting = status === 'testing';

  // 取得 Emoji 標誌
  const getKindEmoji = (k: typeof kind) => {
    switch (k) {
      case 'sqlite': return '💾';
      case 'mysql': return '🐬';
      case 'sqlserver': return '🖥️';
      case 'postgres':
      default:
        return '🐘';
    }
  };

  // 取得 Status Chip 樣式與文字
  const getStatusConfig = (st: typeof status) => {
    switch (st) {
      case 'ready':
        return {
          text: t('settings.status_ready', '連線成功'),
          className: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
        };
      case 'testing':
        return {
          text: t('settings.status_testing', '正在測試…'),
          className: 'bg-blue-500/10 text-blue-400 border-blue-500/20'
        };
      case 'unreachable':
      case 'auth_failed':
        return {
          text: t('settings.status_failed', '無法連線'),
          className: 'bg-red-500/10 text-red-400 border-red-500/20'
        };
      case 'unknown':
      default:
        return {
          text: t('settings.status_unknown', '未測試'),
          className: 'bg-gray-500/10 text-gray-400 border-gray-500/20'
        };
    }
  };

  const statusConfig = getStatusConfig(status);

  return (
    <div className="p-5 space-y-4 bg-gray-950/20 hover:bg-gray-950/40 rounded-xl border border-gray-800/40 transition-all duration-200" data-testid={`connector-row-${id}`}>
      {/* 第一列：標題、類型、啟用 */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <span className="text-2xl select-none" role="img" aria-label={kind}>
            {getKindEmoji(kind)}
          </span>
          <div className="flex items-center gap-2">
            <input
              type="text"
              value={name}
              disabled={isTesting}
              onChange={(e) => onUpdate({ name: e.target.value })}
              className="bg-transparent text-sm font-semibold text-white border-b border-transparent hover:border-gray-700 focus:border-blue-500 outline-none px-1 py-0.5 transition-colors focus:ring-0"
              placeholder="連線名稱"
            />
            <span className="text-[10px] text-gray-500 px-1.5 py-0.5 rounded bg-gray-900 border border-gray-800 select-none">
              {kind.toUpperCase()}
            </span>
          </div>
        </div>

        {/* 總開關 */}
        <div className="flex items-center gap-2">
          <span className="text-xs text-gray-500 select-none">{t('settings.col_enabled', '啟用')}</span>
          <button
            type="button"
            disabled={isTesting}
            onClick={() => onUpdate({ enabled: !connector.enabled })}
            className={`
              relative inline-flex h-4 w-7 items-center rounded-full transition-colors duration-200 focus:outline-none cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed
              ${connector.enabled ? 'bg-blue-600' : 'bg-gray-800'}
            `}
          >
            <span
              className={`
                inline-block h-2.5 w-2.5 transform rounded-full bg-white transition-transform duration-200
                ${connector.enabled ? 'translate-x-3.5' : 'translate-x-0.5'}
              `}
            />
          </button>
        </div>
      </div>

      {/* 第二列：連線欄位 (12-col grid) */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-12 gap-3 text-xs">
        {/* 資料庫類型選擇 */}
        <div className="md:col-span-2 space-y-1">
          <label className="text-gray-500">{t('step4.connector_kind', '類型')}</label>
          <select
            value={kind}
            disabled={isTesting}
            onChange={(e) => onUpdate({ kind: e.target.value as 'postgres' | 'sqlite' | 'mysql' | 'sqlserver' })}
            className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white focus:border-blue-500 outline-none transition-all cursor-pointer"
          >
            <option value="postgres">PostgreSQL</option>
            <option value="sqlite">SQLite</option>
            <option value="mysql">MySQL</option>
            <option value="sqlserver">SQL Server</option>
          </select>
        </div>

        {/* Host */}
        {!isSqlite && (
          <div className="md:col-span-3 space-y-1">
            <label className="text-gray-500">{t('step4.field_host', '主機位址')}</label>
            <input
              type="text"
              value={host}
              disabled={isTesting}
              onChange={(e) => onUpdate({ host: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white focus:border-blue-500 outline-none transition-all font-mono"
            />
          </div>
        )}

        {/* Port */}
        {!isSqlite && (
          <div className="md:col-span-1.5 space-y-1">
            <label className="text-gray-500">{t('step4.field_port', '通訊埠')}</label>
            <input
              type="number"
              value={port}
              disabled={isTesting}
              onChange={(e) => onUpdate({ port: parseInt(e.target.value) || 0 })}
              className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white focus:border-blue-500 outline-none transition-all font-mono"
            />
          </div>
        )}

        {/* Database 名稱 */}
        <div className={isSqlite ? 'md:col-span-5 space-y-1' : 'md:col-span-2 space-y-1'}>
          <label className="text-gray-500">
            {isSqlite ? t('step4.field_db_sqlite', '檔案路徑') : t('step4.field_database', '資料庫')}
          </label>
          <input
            type="text"
            value={database}
            disabled={isTesting}
            onChange={(e) => onUpdate({ database: e.target.value })}
            className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white focus:border-blue-500 outline-none transition-all font-mono"
          />
        </div>

        {/* Username */}
        {!isSqlite && (
          <div className="md:col-span-2 space-y-1">
            <label className="text-gray-500">{t('step4.field_username', '使用者')}</label>
            <input
              type="text"
              value={username}
              disabled={isTesting}
              onChange={(e) => onUpdate({ username: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white focus:border-blue-500 outline-none transition-all"
            />
          </div>
        )}

        {!isSqlite && (
          <div className="md:col-span-2 space-y-1">
            <label className="text-gray-500">{t('step4.field_password', '密碼')}</label>
            <input
              type="password"
              value={password ?? ''}
              disabled={isTesting}
              onChange={(e) => onUpdate({ password: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white focus:border-blue-500 outline-none transition-all"
            />
          </div>
        )}

        {/* Schema (Postgres / SQL Server) */}
        {hasSchema && (
          <div className="md:col-span-1.5 space-y-1">
            <label className="text-gray-500">{t('step4.field_schema', '綱要')}</label>
            <input
              type="text"
              value={schema}
              disabled={isTesting}
              onChange={(e) => onUpdate({ schema: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white focus:border-blue-500 outline-none transition-all"
            />
          </div>
        )}

        {/* Table 名稱 */}
        <div className="md:col-span-2 space-y-1">
          <label className="text-gray-500">{t('step4.field_table', '資料表')}</label>
          <input
            type="text"
            value={table}
            disabled={isTesting}
            onChange={(e) => onUpdate({ table: e.target.value })}
            className="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 text-white focus:border-blue-500 outline-none transition-all"
          />
        </div>
      </div>

      {/* 第三列：跨欄狀態與控制按鈕 */}
      <div className="flex flex-wrap items-center justify-between gap-4 pt-2 border-t border-gray-900 text-xs">
        {/* 左側狀態 chip */}
        <div className="flex items-center gap-3">
          <span className={`px-2 py-0.5 rounded text-[10px] font-bold border select-none ${statusConfig.className}`}>
            {statusConfig.text}
          </span>
          {status === 'ready' && last_check_at && (
            <span className="text-[10px] text-gray-500 select-none">
              {t('settings.last_checked', '上次檢查')}：{new Date(last_check_at).toLocaleTimeString()}
            </span>
          )}
          {(status === 'unreachable' || status === 'auth_failed') && last_check_error && (
            <span className="text-[10px] text-red-500 font-mono">
              {last_check_error}
            </span>
          )}
        </div>

        {/* 右側按鈕 */}
        <div className="flex items-center gap-2">
          {/* 測試連線 */}
          <button
            type="button"
            disabled={isTesting}
            onClick={onTest}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-gray-800 hover:border-gray-700 bg-gray-900/60 hover:bg-gray-900 text-gray-300 font-medium cursor-pointer transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isTesting ? (
              <svg className="animate-spin h-3.5 w-3.5 text-blue-400" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
            ) : (
              <span>⚡</span>
            )}
            <span>{isTesting ? t('settings.btn_testing', '測試中…') : t('settings.btn_test', '測試連線')}</span>
          </button>

          {/* 刪除連接器 */}
          <button
            type="button"
            disabled={isTesting}
            onClick={onRemove}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-red-950/20 hover:border-red-900 bg-red-950/10 hover:bg-red-950/30 text-red-400 hover:text-red-300 font-medium cursor-pointer transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span>🗑️</span>
            <span>{t('settings.btn_delete', '刪除')}</span>
          </button>
        </div>
      </div>
    </div>
  );
}
