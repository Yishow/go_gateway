import { useTranslation } from 'react-i18next';
import type { SettingsConnector } from '../state/types';
import { getDbKindPatch } from '../state/dbSchemas';

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

  const handleKindChange = (nextKind: typeof kind) => {
    const patch = getDbKindPatch(nextKind, connector);
    // 換掉資料庫類型等於換掉連線身分，舊憑證不得沿用（SQLite 不需要密碼）。
    onUpdate({ ...patch, status: 'unknown', password_required: nextKind !== 'sqlite' });
  };

  const getKindEmoji = (k: typeof kind) => {
    switch (k) {
      case 'sqlite': return '💾';
      case 'mysql': return '🐬';
      case 'sqlserver': return '🖥️';
      case 'postgres':
      default: return '🐘';
    }
  };

  const getStatusConfig = (st: typeof status) => {
    switch (st) {
      case 'ready':
        return { text: t('settings.status_ready', '連線成功'), className: 'bg-[var(--gw-success-bg)] text-[var(--gw-success-fg)] border-[var(--gw-success-border)]' };
      case 'testing':
        return { text: t('settings.status_testing', '正在測試…'), className: 'bg-[var(--gw-active)] text-[var(--gw-accent)] border-[var(--gw-focus)]' };
      case 'unreachable':
      case 'auth_failed':
        return { text: t('settings.status_failed', '無法連線'), className: 'bg-[var(--gw-danger-bg)] text-[var(--gw-danger-fg)] border-[var(--gw-danger-border)]' };
      case 'unknown':
      default:
        return { text: t('settings.status_unknown', '未測試'), className: 'bg-[var(--gw-surface)] text-[var(--gw-muted)] border-[var(--gw-border)]' };
    }
  };

  const statusConfig = getStatusConfig(status);


  return (
    <div className="p-5 space-y-4 bg-[var(--gw-surface)] hover:bg-[var(--gw-surface)] rounded-xl border border-[var(--gw-border)] transition-all duration-200" data-testid={`connector-row-${id}`}>
      {/* 第一列：標題、類型、啟用 */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex min-w-0 w-full sm:w-auto items-center gap-3">
          <span className="shrink-0 text-2xl select-none" role="img" aria-label={kind}>
            {getKindEmoji(kind)}
          </span>
          <div className="flex min-w-0 flex-1 flex-wrap items-center gap-2">
            <input
              type="text"
              value={name}
              disabled={isTesting}
              onChange={(e) => onUpdate({ name: e.target.value })}
              className="min-w-0 w-full sm:w-auto bg-transparent text-sm font-semibold text-[var(--gw-text)] border-b border-transparent hover:border-[var(--gw-control)] focus:border-[var(--gw-focus)] outline-none px-1 py-0.5 transition-colors focus:ring-0"
              placeholder="連線名稱"
            />
            <span className="text-xs text-[var(--gw-muted)] px-1.5 py-0.5 rounded bg-[var(--gw-surface)] border border-[var(--gw-border)] select-none">
              {kind.toUpperCase()}
            </span>
          </div>
        </div>

        {/* 總開關 */}
        <div className="flex items-center gap-2">
          <span className="text-xs text-[var(--gw-muted)] select-none">{t('settings.col_enabled', '啟用')}</span>
          <button
            type="button"
            disabled={isTesting}
            onClick={() => onUpdate({ enabled: !connector.enabled })}
            role="switch"
            aria-label={`${t('settings.col_enabled', '啟用')} ${name}`}
            aria-checked={connector.enabled}
            className="inline-flex h-11 w-11 shrink-0 items-center justify-center bg-transparent cursor-pointer disabled:cursor-not-allowed"
          >
            <span className={`relative inline-flex h-4 w-7 shrink-0 items-center rounded-full transition-colors duration-200 ${connector.enabled ? 'bg-[var(--gw-primary)]' : 'bg-[var(--gw-control)]'}`}>
            <span
              className={`
                inline-block h-2.5 w-2.5 transform rounded-full bg-[var(--gw-on-primary)] transition-transform duration-200
                ${connector.enabled ? 'translate-x-3.5' : 'translate-x-0.5'}
              `}
            />
            </span>
          </button>
        </div>
      </div>

      {/* 第二列：連線欄位 (12-col grid) */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-12 gap-3 text-xs">
        {/* 資料庫類型選擇 */}
        <div className="md:col-span-2 space-y-1">
          <label className="text-[var(--gw-muted)]">{t('step4.connector_kind', '類型')}</label>
          <select
            value={kind}
            disabled={isTesting}
            onChange={(e) => handleKindChange(e.target.value as 'postgres' | 'sqlite' | 'mysql' | 'sqlserver')}
            className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded px-2 py-1 text-[var(--gw-text)] focus:border-[var(--gw-focus)] outline-none transition-all cursor-pointer"
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
            <label className="text-[var(--gw-muted)]">{t('step4.field_host', '主機位址')}</label>
            <input
              type="text"
              value={host}
              disabled={isTesting}
              onChange={(e) => onUpdate({ host: e.target.value })}
              className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded px-2 py-1 text-[var(--gw-text)] focus:border-[var(--gw-focus)] outline-none transition-all font-mono"
            />
          </div>
        )}

        {/* Port */}
        {!isSqlite && (
          <div className="md:col-span-1.5 space-y-1">
            <label className="text-[var(--gw-muted)]">{t('step4.field_port', '通訊埠')}</label>
            <input
              type="number"
              value={port}
              disabled={isTesting}
              onChange={(e) => onUpdate({ port: parseInt(e.target.value) || 0 })}
              className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded px-2 py-1 text-[var(--gw-text)] focus:border-[var(--gw-focus)] outline-none transition-all font-mono"
            />
          </div>
        )}

        {/* Database 名稱 */}
        <div className={isSqlite ? 'md:col-span-5 space-y-1' : 'md:col-span-2 space-y-1'}>
          <label className="text-[var(--gw-muted)]">
            {isSqlite ? t('step4.field_db_sqlite', '檔案路徑') : t('step4.field_database', '資料庫')}
          </label>
          <input
            type="text"
            value={database}
            disabled={isTesting}
            onChange={(e) => onUpdate({ database: e.target.value })}
            className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded px-2 py-1 text-[var(--gw-text)] focus:border-[var(--gw-focus)] outline-none transition-all font-mono"
          />
        </div>

        {/* Username */}
        {!isSqlite && (
          <div className="md:col-span-2 space-y-1">
            <label className="text-[var(--gw-muted)]">{t('step4.field_username', '使用者')}</label>
            <input
              type="text"
              value={username}
              disabled={isTesting}
              onChange={(e) => onUpdate({ username: e.target.value })}
              className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded px-2 py-1 text-[var(--gw-text)] focus:border-[var(--gw-focus)] outline-none transition-all"
            />
          </div>
        )}

        {!isSqlite && (
          <div className="md:col-span-2 space-y-1">
            <label className="text-[var(--gw-muted)]">{t('step4.field_password', '密碼')}</label>
            <input
              type="password"
              value={password ?? ''}
              disabled={isTesting}
              onChange={(e) => onUpdate({
                password: e.target.value,
                ...(e.target.value.trim() !== '' ? { password_required: false } : {}),
              })}
              className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded px-2 py-1 text-[var(--gw-text)] focus:border-[var(--gw-focus)] outline-none transition-all"
            />
          </div>
        )}

        {/* Schema (Postgres / SQL Server) */}
        {hasSchema && (
          <div className="md:col-span-1.5 space-y-1">
            <label className="text-[var(--gw-muted)]">{t('step4.field_schema', '綱要')}</label>
            <input
              type="text"
              value={schema}
              disabled={isTesting}
              onChange={(e) => onUpdate({ schema: e.target.value })}
              className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded px-2 py-1 text-[var(--gw-text)] focus:border-[var(--gw-focus)] outline-none transition-all"
            />
          </div>
        )}

        {/* Table 名稱 */}
        <div className="md:col-span-2 space-y-1">
          <label className="text-[var(--gw-muted)]">{t('step4.field_table', '資料表')}</label>
          <input
            type="text"
            value={table}
            disabled={isTesting}
            onChange={(e) => onUpdate({ table: e.target.value })}
            className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded px-2 py-1 text-[var(--gw-text)] focus:border-[var(--gw-focus)] outline-none transition-all"
          />
        </div>
      </div>

      {/* 第三列：跨欄狀態與控制按鈕 */}
      <div className="flex flex-wrap items-center justify-between gap-4 pt-2 border-t border-[var(--gw-border)] text-xs">
        {/* 左側狀態 chip */}
        <div className="flex items-center gap-3">
          <span className={`px-2 py-0.5 rounded text-xs font-bold border select-none ${statusConfig.className}`}>
            {statusConfig.text}
          </span>
          {status === 'ready' && last_check_at && (
            <span className="text-xs text-[var(--gw-muted)] select-none">
              {t('settings.last_checked', '上次檢查')}：{new Date(last_check_at).toLocaleTimeString()}
            </span>
          )}
          {(status === 'unreachable' || status === 'auth_failed') && last_check_error && (
            <span className="text-xs text-[var(--gw-danger-fg)] font-mono">
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
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-[var(--gw-control)] hover:border-[var(--gw-control)] bg-[var(--gw-surface)] hover:bg-[var(--gw-surface)] text-[var(--gw-text)] font-medium cursor-pointer transition-all disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)] disabled:cursor-not-allowed"
          >
            {isTesting ? (
              <svg className="animate-spin h-3.5 w-3.5 text-[var(--gw-accent)]" fill="none" viewBox="0 0 24 24">
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
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-[var(--gw-danger-border)] hover:border-[var(--gw-danger-border)] bg-[var(--gw-danger-bg)] hover:bg-[var(--gw-danger-bg)] text-[var(--gw-danger-fg)] hover:text-[var(--gw-danger-fg)] font-medium cursor-pointer transition-all disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)] disabled:cursor-not-allowed"
          >
            <span>🗑️</span>
            <span>{t('settings.btn_delete', '刪除')}</span>
          </button>
        </div>
      </div>
    </div>
  );
}
