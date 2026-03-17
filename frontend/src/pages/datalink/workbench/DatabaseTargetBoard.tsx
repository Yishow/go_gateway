import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { dbTargetAPI } from '../../../services/datalink';
import type {
  DatabaseConnector,
  DatabaseConnectorKind,
  DatabaseTableColumn,
  DatabaseTableInfo,
  DatabaseTargetMapping,
  DatabaseTargetValidationResult,
  DatabaseWriteMode,
} from '../../../types/datalink';
import type { WorkbenchOutputCandidate } from './workbenchOutputTypes';

type DatabaseTargetBoardProps = {
  candidates: WorkbenchOutputCandidate[];
  selectedTagId: string;
  onSelectedTagChange: (tagId: string) => void;
};

type ConnectorDraft = {
  name: string;
  kind: DatabaseConnectorKind;
  enabled: boolean;
  sqliteDsn: string;
  postgresHost: string;
  postgresPort: string;
  postgresUser: string;
  postgresPassword: string;
  postgresDatabase: string;
  postgresSSLMode: string;
};

function createEmptyDraft(): ConnectorDraft {
  return {
    name: '',
    kind: 'sqlite',
    enabled: true,
    sqliteDsn: '',
    postgresHost: '127.0.0.1',
    postgresPort: '5432',
    postgresUser: '',
    postgresPassword: '',
    postgresDatabase: '',
    postgresSSLMode: 'disable',
  };
}

function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}

function buildTableKey(table: DatabaseTableInfo): string {
  return `${table.schema}.${table.name}`;
}

function parseTableKey(key: string): { schema: string; name: string } | null {
  const [schemaName, ...nameParts] = key.split('.');
  if (!schemaName || nameParts.length === 0) {
    return null;
  }
  return {
    schema: schemaName,
    name: nameParts.join('.'),
  };
}

function defaultValueColumn(columns: DatabaseTableColumn[]): DatabaseTableColumn | null {
  return columns.find((column) => !column.primary_key) ?? columns[0] ?? null;
}

function defaultTimestampColumn(columns: DatabaseTableColumn[]): DatabaseTableColumn | null {
  return (
    columns.find((column) =>
      ['ts', 'timestamp', 'created_at', 'updated_at'].includes(
        column.name.toLowerCase(),
      ) &&
      (column.primary_key || column.unique),
    ) ?? null
  );
}

function applyConnectorDraft(connector: DatabaseConnector | null): ConnectorDraft {
  if (!connector) {
    return createEmptyDraft();
  }

  const config = connector.connection_config ?? {};

  return {
    name: connector.name,
    kind: connector.kind,
    enabled: connector.enabled,
    sqliteDsn: String(config.dsn ?? config.path ?? ''),
    postgresHost: String(config.host ?? '127.0.0.1'),
    postgresPort: String(config.port ?? '5432'),
    postgresUser: String(config.user ?? ''),
    postgresPassword: String(config.password ?? ''),
    postgresDatabase: String(config.database ?? config.dbname ?? ''),
    postgresSSLMode: String(config.sslmode ?? 'disable'),
  };
}

function buildConnectionConfig(draft: ConnectorDraft): Record<string, unknown> {
  if (draft.kind === 'sqlite') {
    return {
      dsn: draft.sqliteDsn.trim(),
    };
  }

  return {
    host: draft.postgresHost.trim(),
    port: draft.postgresPort.trim(),
    user: draft.postgresUser.trim(),
    password: draft.postgresPassword,
    database: draft.postgresDatabase.trim(),
    sslmode: draft.postgresSSLMode.trim() || 'disable',
  };
}

function statusBadgeClasses(status: DatabaseConnector['status']): string {
  switch (status) {
    case 'ready':
      return 'bg-emerald-500/10 text-emerald-200';
    case 'auth_failed':
      return 'bg-amber-500/10 text-amber-100';
    case 'unreachable':
      return 'bg-rose-500/10 text-rose-100';
    default:
      return 'bg-slate-800 text-slate-200';
  }
}

export function DatabaseTargetBoard({
  candidates,
  selectedTagId,
  onSelectedTagChange: _onSelectedTagChange,
}: DatabaseTargetBoardProps) {
  const { t } = useTranslation();
  const tRef = useRef(t);
  tRef.current = t;

  const [connectors, setConnectors] = useState<DatabaseConnector[]>([]);
  const [mappings, setMappings] = useState<DatabaseTargetMapping[]>([]);
  const [tables, setTables] = useState<DatabaseTableInfo[]>([]);
  const [validation, setValidation] = useState<DatabaseTargetValidationResult | null>(
    null,
  );
  const [selectedConnectorId, setSelectedConnectorId] = useState('');
  const [isConnectorEditorOpen, setIsConnectorEditorOpen] = useState(false);
  const [draft, setDraft] = useState<ConnectorDraft>(() => createEmptyDraft());
  const [tableKey, setTableKey] = useState('');
  const [columnName, setColumnName] = useState('');
  const [writeMode, setWriteMode] = useState<DatabaseWriteMode>('insert');
  const [timestampColumn, setTimestampColumn] = useState('');
  const [clearSavedPassword, setClearSavedPassword] = useState(false);
  const [message, setMessage] = useState('');
  const [isBusy, setIsBusy] = useState(true);

  const selectedConnector =
    connectors.find((connector) => connector.id === selectedConnectorId) ?? null;
  const canClearSavedPassword =
    Boolean(selectedConnectorId) && draft.kind === 'postgres';
  const selectedCandidate =
    candidates.find((candidate) => candidate.tagId === selectedTagId) ?? null;
  const selectedMapping =
    mappings.find(
      (mapping) =>
        mapping.connector_id === selectedConnectorId && mapping.tag_id === selectedTagId,
    ) ?? null;
  const selectedTable =
    tables.find((table) => buildTableKey(table) === tableKey) ?? null;
  const tableColumns = selectedTable?.columns ?? [];
  const connectorMappings = useMemo(
    () =>
      mappings.filter((mapping) => mapping.connector_id === selectedConnectorId),
    [mappings, selectedConnectorId],
  );
  const validationErrorCount = useMemo(
    () =>
      validation?.issues.filter((issue) => issue.severity === 'error').length ?? 0,
    [validation],
  );

  const loadData = useCallback(async () => {
    setIsBusy(true);
    try {
      const [nextConnectors, nextMappings] = await Promise.all([
        dbTargetAPI.listConnectors(),
        dbTargetAPI.listMappings(),
      ]);
      setConnectors(nextConnectors);
      setMappings(nextMappings);
      setSelectedConnectorId((previous) => {
        if (previous && nextConnectors.some((connector) => connector.id === previous)) {
          return previous;
        }
        return nextConnectors[0]?.id ?? '';
      });
      setMessage('');
    } catch (error) {
      setMessage(
        getErrorMessage(error, tRef.current('workbench.output.database.results.load')),
      );
    } finally {
      setIsBusy(false);
    }
  }, []);

  const loadConnectorDetails = useCallback(async (connectorId: string) => {
    if (!connectorId) {
      setTables([]);
      setValidation(null);
      return;
    }

    try {
      const [nextTables, nextValidation] = await Promise.all([
        dbTargetAPI.listTables(connectorId),
        dbTargetAPI.validateConnector(connectorId),
      ]);
      setTables(nextTables);
      setValidation(nextValidation);
    } catch (error) {
      setTables([]);
      setValidation(null);
      setMessage(
        getErrorMessage(
          error,
          tRef.current('workbench.output.database.results.connectorDetail'),
        ),
      );
    }
  }, []);

  useEffect(() => {
    void loadData();
  }, [loadData]);

  useEffect(() => {
    setDraft(applyConnectorDraft(selectedConnector));
    setClearSavedPassword(false);
    void loadConnectorDetails(selectedConnectorId);
  }, [loadConnectorDetails, selectedConnector, selectedConnectorId]);

  useEffect(() => {
    if (draft.kind !== 'postgres' && clearSavedPassword) {
      setClearSavedPassword(false);
    }
  }, [clearSavedPassword, draft.kind]);

  useEffect(() => {
    if (!selectedConnectorId) {
      setTableKey('');
      setColumnName('');
      setWriteMode('insert');
      setTimestampColumn('');
      return;
    }

    if (selectedMapping) {
      setTableKey(`${selectedMapping.table_schema}.${selectedMapping.table_name}`);
      setColumnName(selectedMapping.column_name);
      setWriteMode(selectedMapping.write_mode);
      setTimestampColumn(selectedMapping.timestamp_column ?? '');
      return;
    }

    const firstTable = tables[0];
    if (!firstTable) {
      setTableKey('');
      setColumnName('');
      setWriteMode('insert');
      setTimestampColumn('');
      return;
    }

    const nextValueColumn = defaultValueColumn(firstTable.columns);
    const nextTimestampColumn = defaultTimestampColumn(firstTable.columns);
    setTableKey(buildTableKey(firstTable));
    setColumnName(nextValueColumn?.name ?? '');
    setWriteMode('insert');
    setTimestampColumn(nextTimestampColumn?.name ?? '');
  }, [selectedConnectorId, selectedMapping, tables]);

  useEffect(() => {
    if (!selectedTable) {
      return;
    }

    if (columnName && selectedTable.columns.some((column) => column.name === columnName)) {
      return;
    }

    const nextValueColumn = defaultValueColumn(selectedTable.columns);
    setColumnName(nextValueColumn?.name ?? '');
  }, [columnName, selectedTable]);

  useEffect(() => {
    if (writeMode !== 'upsert') {
      if (timestampColumn !== '') {
        setTimestampColumn('');
      }
      return;
    }
    if (!selectedTable) {
      return;
    }
    if (
      timestampColumn &&
      selectedTable.columns.some((column) => column.name === timestampColumn)
    ) {
      return;
    }

    const nextTimestampColumn = defaultTimestampColumn(selectedTable.columns);
    setTimestampColumn(nextTimestampColumn?.name ?? '');
  }, [selectedTable, timestampColumn, writeMode]);

  const handleNewConnector = useCallback(() => {
    setSelectedConnectorId('');
    setIsConnectorEditorOpen(true);
    setDraft(createEmptyDraft());
    setClearSavedPassword(false);
    setTables([]);
    setValidation(null);
    setMessage('');
  }, []);

  const handleConnectorField = useCallback(
    <K extends keyof ConnectorDraft,>(key: K, value: ConnectorDraft[K]) => {
      setDraft((previous) => ({ ...previous, [key]: value }));
    },
    [],
  );

  const handleSaveConnector = useCallback(async () => {
    const trimmedName = draft.name.trim();
    if (!trimmedName) {
      setMessage(t('workbench.output.database.results.connectorNameRequired'));
      return;
    }

    if (draft.kind === 'sqlite' && !draft.sqliteDsn.trim()) {
      setMessage(t('workbench.output.database.results.sqliteDsnRequired'));
      return;
    }

    setIsBusy(true);
    try {
      const payload = {
        name: trimmedName,
        kind: draft.kind,
        enabled: draft.enabled,
        connection_config: buildConnectionConfig(draft),
      };

      const connector = selectedConnectorId
        ? await dbTargetAPI.updateConnector(selectedConnectorId, {
            ...payload,
            clear_password: canClearSavedPassword && clearSavedPassword
              ? true
              : undefined,
          })
        : await dbTargetAPI.createConnector(payload);

      setSelectedConnectorId(connector.id);
      await loadData();
      setMessage(
        t('workbench.output.database.results.connectorSaved', {
          name: connector.name,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.database.results.connectorSaveFailed')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [canClearSavedPassword, clearSavedPassword, draft, loadData, selectedConnectorId, t]);

  const handleTestConnector = useCallback(async () => {
    if (!selectedConnectorId) {
      setMessage(t('workbench.output.database.results.connectorNotSaved'));
      return;
    }

    setIsBusy(true);
    try {
      const connector = await dbTargetAPI.testConnector(selectedConnectorId);
      setSelectedConnectorId(connector.id);
      await loadData();
      setMessage(
        t('workbench.output.database.results.connectorTested', {
          status: connector.status,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.database.results.connectorTestFailed')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadData, selectedConnectorId, t]);

  const handleDeleteConnector = useCallback(async () => {
    if (!selectedConnectorId || !selectedConnector) {
      return;
    }

    setIsBusy(true);
    try {
      await dbTargetAPI.deleteConnector(selectedConnectorId);
      await loadData();
      setMessage(
        t('workbench.output.database.results.connectorDeleted', {
          name: selectedConnector.name,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(
          error,
          t('workbench.output.database.results.connectorDeleteFailed'),
        ),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadData, selectedConnector, selectedConnectorId, t]);

  const handleSaveMapping = useCallback(async () => {
    if (!selectedConnectorId) {
      setMessage(t('workbench.output.database.results.connectorNotSaved'));
      return;
    }
    if (!selectedCandidate) {
      setMessage(t('workbench.output.database.results.noTagSelected'));
      return;
    }
    if (!selectedTable) {
      setMessage(t('workbench.output.database.results.tableRequired'));
      return;
    }
    if (!columnName) {
      setMessage(t('workbench.output.database.results.columnRequired'));
      return;
    }
    if (writeMode === 'upsert' && !timestampColumn.trim()) {
      setMessage(t('workbench.output.database.results.timestampRequired'));
      return;
    }

    setIsBusy(true);
    try {
        if (selectedMapping) {
          await dbTargetAPI.updateMapping(selectedMapping.id, {
            table_schema: selectedTable.schema,
            table_name: selectedTable.name,
            column_name: columnName,
            write_mode: writeMode,
            timestamp_column:
              writeMode === 'upsert' ? timestampColumn.trim() : '',
            enabled: true,
          });
        } else {
          await dbTargetAPI.createMapping({
            tag_id: selectedCandidate.tagId,
            connector_id: selectedConnectorId,
            table_schema: selectedTable.schema,
            table_name: selectedTable.name,
            column_name: columnName,
            write_mode: writeMode,
            timestamp_column:
              writeMode === 'upsert' ? timestampColumn.trim() : '',
            enabled: true,
          });
        }

      await loadData();
      setMessage(
        t('workbench.output.database.results.mappingSaved', {
          key: selectedCandidate.tagKey,
          table: selectedTable.name,
          column: columnName,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.database.results.mappingSaveFailed')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [
    columnName,
    loadData,
    selectedCandidate,
    selectedConnectorId,
    selectedMapping,
    selectedTable,
    t,
    timestampColumn,
    writeMode,
  ]);

  const handleDeleteMapping = useCallback(async () => {
    if (!selectedMapping || !selectedCandidate) {
      return;
    }

    setIsBusy(true);
    try {
      await dbTargetAPI.deleteMapping(selectedMapping.id);
      await loadData();
      setMessage(
        t('workbench.output.database.results.mappingDeleted', {
          key: selectedCandidate.tagKey,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(
          error,
          t('workbench.output.database.results.mappingDeleteFailed'),
        ),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadData, selectedCandidate, selectedMapping, t]);

  const handleRefreshValidation = useCallback(async () => {
    if (!selectedConnectorId) {
      setMessage(t('workbench.output.database.results.connectorNotSaved'));
      return;
    }

    setIsBusy(true);
    try {
      const nextValidation = await dbTargetAPI.validateConnector(selectedConnectorId);
      setValidation(nextValidation);
      setMessage(
        t('workbench.output.database.results.validationRefreshed', {
          ready: nextValidation.ready
            ? t('workbench.output.database.validation.ready')
            : t('workbench.output.database.validation.notReady'),
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(
          error,
          t('workbench.output.database.results.validationRefreshFailed'),
        ),
      );
    } finally {
      setIsBusy(false);
    }
  }, [selectedConnectorId, t]);

  const currentTableFromKey = parseTableKey(tableKey);
  const mappedTagCount = connectorMappings.length;
  const schemaColumns = selectedTable?.columns ?? [];

  return (
    <section className="space-y-6 rounded-2xl border border-slate-800 bg-slate-950/40 p-5">
      <div className="space-y-2">
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-violet-300">
          {t('workbench.output.database.eyebrow')}
        </p>
        <div>
          <h3 className="text-2xl font-semibold text-slate-50">
            {t('workbench.output.database.title')}
          </h3>
          <p className="mt-1 max-w-3xl text-sm text-slate-300">
            {t('workbench.output.database.description')}
          </p>
        </div>
      </div>

      <div className="grid gap-6 xl:grid-cols-[340px_minmax(0,1fr)]">
        <aside className="space-y-4 rounded-2xl border border-slate-800 bg-slate-900/70 p-4">
          <div className="grid gap-3 sm:grid-cols-3 xl:grid-cols-1">
            <article className="rounded-xl border border-slate-800 bg-slate-950/70 p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.output.database.summary.connector')}
              </p>
              <p className="mt-2 text-lg font-semibold text-slate-50">
                {selectedConnector
                  ? t(`workbench.output.database.status.${selectedConnector.status}`)
                  : t('workbench.output.database.summary.none')}
              </p>
            </article>
            <article className="rounded-xl border border-slate-800 bg-slate-950/70 p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.output.database.summary.mappingCount')}
              </p>
              <p className="mt-2 text-lg font-semibold text-slate-50">
                {mappedTagCount}
              </p>
            </article>
            <article className="rounded-xl border border-slate-800 bg-slate-950/70 p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.output.database.summary.validationIssues')}
              </p>
              <p className="mt-2 text-lg font-semibold text-slate-50">
                {validationErrorCount}
              </p>
            </article>
          </div>

          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.output.database.connector.listTitle')}
              </p>
              <button
                type="button"
                onClick={handleNewConnector}
                className="rounded-lg border border-slate-700 px-3 py-1 text-xs text-slate-200"
              >
                {t('workbench.output.database.actions.newConnector')}
              </button>
            </div>

            {connectors.length === 0 ? (
              <p className="rounded-xl border border-dashed border-slate-700 bg-slate-950/60 px-3 py-4 text-sm text-slate-400">
                {t('workbench.output.database.connector.empty')}
              </p>
            ) : (
              <div className="grid gap-2">
                {connectors.map((connector) => {
                  const active = connector.id === selectedConnectorId;
                  return (
                     <button
                       key={connector.id}
                       type="button"
                       aria-pressed={active}
                       onClick={() => {
                         setSelectedConnectorId(connector.id);
                         setIsConnectorEditorOpen(true);
                       }}
                       className={`flex items-center justify-between rounded-xl border px-3 py-3 text-left ${
                         active
                           ? 'border-violet-500/40 bg-violet-500/5'
                           : 'border-slate-800 bg-slate-950/70'
                      }`}
                    >
                      <span className="space-y-1">
                        <span className="block text-sm font-semibold text-slate-50">
                          {connector.name}
                        </span>
                        <span className="block text-xs text-slate-400">
                          {connector.kind}
                        </span>
                      </span>
                      <span
                        className={`rounded-full px-2 py-1 text-[11px] font-medium ${statusBadgeClasses(
                          connector.status,
                        )}`}
                      >
                        {t(`workbench.output.database.status.${connector.status}`)}
                      </span>
                    </button>
                  );
                })}
              </div>
            )}

            <button
              type="button"
              onClick={() => setIsConnectorEditorOpen((currentValue) => !currentValue)}
              className="rounded-lg border border-slate-700 px-3 py-2 text-xs text-slate-200"
            >
              {isConnectorEditorOpen
                ? t('workbench.output.database.actions.hideConnector')
                : t('workbench.output.database.actions.configureConnector')}
            </button>
          </div>

          {isConnectorEditorOpen ? (
            <div className="space-y-3 rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.output.database.connector.name')}</span>
                <input
                  aria-label={t('workbench.output.database.connector.name')}
                  value={draft.name}
                  onChange={(event) => handleConnectorField('name', event.target.value)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                />
              </label>

            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.output.database.connector.kind')}</span>
              <select
                aria-label={t('workbench.output.database.connector.kind')}
                value={draft.kind}
                onChange={(event) =>
                  handleConnectorField(
                    'kind',
                    event.target.value as DatabaseConnectorKind,
                  )
                }
                className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
              >
                <option value="sqlite">sqlite</option>
                <option value="postgres">postgres</option>
              </select>
            </label>

            {draft.kind === 'sqlite' ? (
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.output.database.connector.sqliteDsn')}</span>
                <input
                  aria-label={t('workbench.output.database.connector.sqliteDsn')}
                  value={draft.sqliteDsn}
                  onChange={(event) =>
                    handleConnectorField('sqliteDsn', event.target.value)
                  }
                  className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                />
              </label>
            ) : (
              <div className="grid gap-3 sm:grid-cols-2">
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.output.database.connector.postgresHost')}</span>
                  <input
                    aria-label={t('workbench.output.database.connector.postgresHost')}
                    value={draft.postgresHost}
                    onChange={(event) =>
                      handleConnectorField('postgresHost', event.target.value)
                    }
                    className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                  />
                </label>
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.output.database.connector.postgresPort')}</span>
                  <input
                    aria-label={t('workbench.output.database.connector.postgresPort')}
                    value={draft.postgresPort}
                    onChange={(event) =>
                      handleConnectorField('postgresPort', event.target.value)
                    }
                    className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                  />
                </label>
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.output.database.connector.postgresUser')}</span>
                  <input
                    aria-label={t('workbench.output.database.connector.postgresUser')}
                    value={draft.postgresUser}
                    onChange={(event) =>
                      handleConnectorField('postgresUser', event.target.value)
                    }
                    className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                  />
                </label>
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.output.database.connector.postgresPassword')}</span>
                  <input
                    aria-label={t('workbench.output.database.connector.postgresPassword')}
                    type="password"
                    value={draft.postgresPassword}
                    placeholder={
                      canClearSavedPassword
                        ? t('workbench.output.database.connector.postgresPasswordHint')
                        : undefined
                    }
                    disabled={clearSavedPassword}
                    onChange={(event) => {
                      setClearSavedPassword(false);
                      handleConnectorField('postgresPassword', event.target.value);
                    }}
                    className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                  />
                  {canClearSavedPassword ? (
                    <span className="block text-[11px] normal-case tracking-normal text-slate-500">
                      {t('workbench.output.database.connector.postgresPasswordHint')}
                    </span>
                  ) : null}
                  {canClearSavedPassword ? (
                    <label className="flex items-center gap-2 text-[11px] normal-case tracking-normal text-slate-400">
                      <input
                        type="checkbox"
                        checked={clearSavedPassword}
                        onChange={(event) => {
                          const nextChecked = event.target.checked;
                          setClearSavedPassword(nextChecked);
                          if (nextChecked) {
                            handleConnectorField('postgresPassword', '');
                          }
                        }}
                      />
                      <span>
                        {t('workbench.output.database.connector.clearSavedPassword')}
                      </span>
                    </label>
                  ) : null}
                </label>
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.output.database.connector.postgresDatabase')}</span>
                  <input
                    aria-label={t('workbench.output.database.connector.postgresDatabase')}
                    value={draft.postgresDatabase}
                    onChange={(event) =>
                      handleConnectorField('postgresDatabase', event.target.value)
                    }
                    className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                  />
                </label>
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.output.database.connector.postgresSslMode')}</span>
                  <input
                    aria-label={t('workbench.output.database.connector.postgresSslMode')}
                    value={draft.postgresSSLMode}
                    onChange={(event) =>
                      handleConnectorField('postgresSSLMode', event.target.value)
                    }
                    className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                  />
                </label>
              </div>
            )}

            <label className="flex items-center gap-3 rounded-xl border border-slate-800 bg-slate-950/80 px-3 py-3 text-sm text-slate-200">
              <input
                type="checkbox"
                checked={draft.enabled}
                onChange={(event) =>
                  handleConnectorField('enabled', event.target.checked)
                }
              />
              <span>{t('workbench.output.database.connector.enabled')}</span>
            </label>

              <div className="flex flex-wrap gap-2">
                <button
                  type="button"
                  onClick={() => void handleSaveConnector()}
                  disabled={isBusy}
                  className="rounded-lg bg-violet-500 px-3 py-2 text-sm font-semibold text-slate-950 disabled:opacity-50"
                >
                  {t('workbench.output.database.actions.saveConnector')}
                </button>
                <button
                  type="button"
                  onClick={() => void handleTestConnector()}
                  disabled={isBusy || !selectedConnectorId}
                  className="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-200 disabled:opacity-50"
                >
                  {t('workbench.output.database.actions.testConnector')}
                </button>
                <button
                  type="button"
                  onClick={() => void handleDeleteConnector()}
                  disabled={isBusy || !selectedConnectorId}
                  className="rounded-lg border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-100 disabled:opacity-50"
                >
                  {t('workbench.output.database.actions.deleteConnector')}
                </button>
              </div>
            </div>
          ) : null}
        </aside>

        <div className="space-y-4 rounded-2xl border border-slate-800 bg-slate-900/70 p-4">
          <div
            className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4"
            data-testid="database-selected-tag"
          >
            <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.output.database.mapping.selectedTag')}
            </p>
            <p className="mt-2 text-sm font-semibold text-slate-50">
              {selectedCandidate?.tagKey ?? t('workbench.output.database.summary.none')}
            </p>
            {selectedCandidate ? (
              <p className="mt-1 text-xs text-slate-400">
                {t('workbench.output.database.candidates.meta', {
                  point: selectedCandidate.pointName,
                  address: selectedCandidate.pointAddress,
                  dataType: selectedCandidate.dataType,
                })}
              </p>
            ) : null}
          </div>

          <div className="grid gap-4 rounded-2xl border border-slate-800 bg-slate-950/70 p-4 lg:grid-cols-2">

            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.output.database.mapping.table')}</span>
              <select
                aria-label={t('workbench.output.database.mapping.table')}
                value={tableKey}
                onChange={(event) => setTableKey(event.target.value)}
                className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
              >
                <option value="">{t('workbench.output.database.mapping.tablePlaceholder')}</option>
                {tables.map((table) => (
                  <option key={buildTableKey(table)} value={buildTableKey(table)}>
                    {buildTableKey(table)}
                  </option>
                ))}
              </select>
            </label>

            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.output.database.mapping.column')}</span>
              <select
                aria-label={t('workbench.output.database.mapping.column')}
                value={columnName}
                onChange={(event) => setColumnName(event.target.value)}
                className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
              >
                <option value="">
                  {t('workbench.output.database.mapping.columnPlaceholder')}
                </option>
                {tableColumns.map((column) => (
                  <option key={column.name} value={column.name}>
                    {column.name} ({column.data_type})
                  </option>
                ))}
              </select>
            </label>

            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.output.database.mapping.writeMode')}</span>
              <select
                aria-label={t('workbench.output.database.mapping.writeMode')}
                value={writeMode}
                onChange={(event) =>
                  setWriteMode(event.target.value as DatabaseWriteMode)
                }
                className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
              >
                <option value="insert">
                  {t('workbench.output.database.writeMode.insert')}
                </option>
                <option value="upsert">
                  {t('workbench.output.database.writeMode.upsert')}
                </option>
              </select>
            </label>

            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400 lg:col-span-2">
              <span>{t('workbench.output.database.mapping.timestampColumn')}</span>
              <select
                aria-label={t('workbench.output.database.mapping.timestampColumn')}
                value={timestampColumn}
                onChange={(event) => setTimestampColumn(event.target.value)}
                disabled={writeMode !== 'upsert'}
                className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 disabled:opacity-50"
              >
                <option value="">
                  {t('workbench.output.database.mapping.timestampPlaceholder')}
                </option>
                {tableColumns.map((column) => (
                  <option key={column.name} value={column.name}>
                    {column.name} ({column.data_type})
                  </option>
                ))}
              </select>
            </label>
          </div>

          <div className="flex flex-wrap gap-2">
            <button
              type="button"
              onClick={() => void handleSaveMapping()}
              disabled={isBusy || !selectedConnectorId || !selectedTagId}
              className="rounded-lg bg-violet-500 px-3 py-2 text-sm font-semibold text-slate-950 disabled:opacity-50"
            >
              {t('workbench.output.database.actions.saveMapping')}
            </button>
            <button
              type="button"
              onClick={() => void handleDeleteMapping()}
              disabled={isBusy || !selectedMapping}
              className="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-200 disabled:opacity-50"
            >
              {t('workbench.output.database.actions.deleteMapping')}
            </button>
            <button
              type="button"
              onClick={() => void handleRefreshValidation()}
              disabled={isBusy || !selectedConnectorId}
              className="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-200 disabled:opacity-50"
            >
              {t('workbench.output.database.actions.refreshValidation')}
            </button>
          </div>

          <div
            className="space-y-4"
            data-testid="database-secondary-panels"
            data-emphasis="supporting"
          >
            {schemaColumns.length > 0 ? (
              <div
                data-testid="schema-snapshot"
                className="space-y-3 rounded-2xl border border-slate-800 bg-slate-950/70 p-4"
              >
                <div className="space-y-1">
                  <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                    {t('workbench.output.database.schema.title')}
                  </p>
                  <p className="text-sm text-slate-300">
                    {t('workbench.output.database.schema.description')}
                  </p>
                </div>

                <div className="grid gap-2 md:grid-cols-2 xl:grid-cols-1">
                  {schemaColumns.map((column) => {
                    const isRequired =
                      !column.nullable
                      && !column.primary_key
                      && column.name !== columnName
                      && (writeMode !== 'upsert' || column.name !== timestampColumn);

                    return (
                      <article
                        key={column.name}
                        data-testid={`schema-column-${column.name}`}
                        data-required={isRequired ? 'true' : undefined}
                        className={`rounded-xl border p-3 ${
                          isRequired
                            ? 'border-amber-500/30 bg-amber-500/10'
                            : 'border-slate-800 bg-slate-900/60'
                        }`}
                      >
                        <div className="flex items-center justify-between gap-3">
                          <p className="text-sm font-semibold text-slate-50">{column.name}</p>
                          <div className="flex flex-wrap gap-2 text-[10px] uppercase tracking-[0.14em]">
                            <span className="rounded-full bg-slate-800 px-2 py-0.5 text-slate-300">
                              {column.data_type}
                            </span>
                            {column.primary_key ? (
                              <span className="rounded-full bg-violet-500/10 px-2 py-0.5 text-violet-200">
                                PK
                              </span>
                            ) : null}
                            {isRequired ? (
                              <span className="rounded-full bg-amber-500/10 px-2 py-0.5 text-amber-200">
                                {t('workbench.output.database.schema.required')}
                              </span>
                            ) : null}
                          </div>
                        </div>
                      </article>
                    );
                  })}
                </div>
              </div>
            ) : null}

            <div
              data-testid="write-row-preview"
              className="space-y-3 rounded-2xl border border-slate-800 bg-slate-950/70 p-4"
            >
              <div className="space-y-1">
                <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                  {t('workbench.output.database.preview.title')}
                </p>
                <p className="text-sm text-slate-300">
                  {t('workbench.output.database.preview.description')}
                </p>
              </div>

              <div className="grid gap-3 sm:grid-cols-3">
                <article className="rounded-xl border border-slate-800 bg-slate-900/60 p-4">
                  <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                    {t('workbench.output.database.preview.table')}
                  </p>
                  <p className="mt-2 text-sm font-semibold text-slate-50">
                    {currentTableFromKey
                      ? `${currentTableFromKey.schema}.${currentTableFromKey.name}`
                      : t('workbench.output.database.preview.empty')}
                  </p>
                </article>
                <article className="rounded-xl border border-slate-800 bg-slate-900/60 p-4">
                  <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                    {t('workbench.output.database.preview.column')}
                  </p>
                  <p className="mt-2 text-sm font-semibold text-slate-50">
                    {columnName || t('workbench.output.database.preview.empty')}
                  </p>
                </article>
                <article className="rounded-xl border border-slate-800 bg-slate-900/60 p-4">
                  <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                    {t('workbench.output.database.preview.value')}
                  </p>
                  <p className="mt-2 break-all text-sm font-semibold text-slate-50">
                    {selectedCandidate?.lastValue === null ||
                    selectedCandidate?.lastValue === undefined
                      ? t('workbench.output.database.preview.empty')
                      : String(selectedCandidate.lastValue)}
                  </p>
                </article>
              </div>
            </div>
          </div>

          {validation?.issues.length ? (
            <div className="space-y-2 rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
              <p className="text-sm font-semibold text-slate-50">
                {validation.ready
                  ? t('workbench.output.database.validation.ready')
                  : t('workbench.output.database.validation.notReady')}
              </p>
              <ul className="space-y-2 text-sm text-slate-300">
                {validation.issues.map((issue) => (
                  <li
                    key={`${issue.code}-${issue.mapping_id ?? 'global'}-${issue.tag_id ?? 'none'}`}
                    className={`rounded-xl border px-3 py-2 ${
                      issue.severity === 'error'
                        ? 'border-rose-500/20 bg-rose-500/10 text-rose-100'
                        : 'border-amber-500/20 bg-amber-500/10 text-amber-100'
                    }`}
                  >
                    {issue.message}
                  </li>
                ))}
              </ul>
            </div>
          ) : (
            <p className="rounded-xl border border-dashed border-slate-700 bg-slate-950/60 px-3 py-4 text-sm text-slate-400">
              {t('workbench.output.database.validation.empty')}
            </p>
          )}

          {message ? (
            <p
              role="status"
              aria-live="polite"
              className="rounded-xl border border-slate-800 bg-slate-950/80 px-3 py-2 text-sm text-slate-200"
            >
              {message}
            </p>
          ) : null}
        </div>
      </div>
    </section>
  );
}
