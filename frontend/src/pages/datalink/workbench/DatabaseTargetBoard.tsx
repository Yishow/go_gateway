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
import type {
  SourceRuleCandidateSetView,
  SourceRuleDatabaseOutputCandidateView,
} from '../../../types/sourceRuleCandidates';
import {
  DATABASE_OUTPUT_SCOPE_INITIAL,
  type DatabaseOutputScope,
} from './databaseOutputModel';
import { DatabaseOutputReviewPanel } from './DatabaseOutputReviewPanel';
import {
  type ConnectorDraft,
  applyConnectorDraft,
  buildConnectionConfig,
  buildTableKey,
  createEmptyDraft,
  defaultTimestampColumn,
  defaultValueColumn,
  getErrorMessage,
  parseTableKey,
  statusBadgeClasses,
} from './databaseTargetBoardUtils';
import { useRefreshSourceRuleCandidates } from './useRefreshSourceRuleCandidates';
import type { WorkbenchOutputCandidate } from './workbenchOutputTypes';

type DatabaseTargetBoardProps = {
  candidates: WorkbenchOutputCandidate[];
  selectedTagId: string;
  reviewRuleId?: string | null;
  reviewRevisionId?: string | null;
  reviewSet?: SourceRuleCandidateSetView<SourceRuleDatabaseOutputCandidateView> | null;
  reviewLoading?: boolean;
}

export function DatabaseTargetBoard({
  candidates,
  selectedTagId,
  reviewRuleId = null,
  reviewRevisionId = null,
  reviewSet = null,
  reviewLoading = false,
}: DatabaseTargetBoardProps) {
  const refreshCandidateReview = useRefreshSourceRuleCandidates();
  const { t } = useTranslation();
  const tRef = useRef(t);
  tRef.current = t;

  const [connectors, setConnectors] = useState<DatabaseConnector[]>([]);
  const [mappings, setMappings] = useState<DatabaseTargetMapping[]>([]);
  const [tables, setTables] = useState<DatabaseTableInfo[]>([]);
  const [validation, setValidation] = useState<DatabaseTargetValidationResult | null>(
    null,
  );
  const [scope, setScope] = useState<DatabaseOutputScope>(
    DATABASE_OUTPUT_SCOPE_INITIAL,
  );
  const [isConnectorEditorOpen, setIsConnectorEditorOpen] = useState(false);
  const [draft, setDraft] = useState<ConnectorDraft>(() => createEmptyDraft());
  const [clearSavedPassword, setClearSavedPassword] = useState(false);
  const [dryRunResults, setDryRunResults] = useState<Array<{
    candidate_id: string;
    status: string;
    code?: string;
    reason?: string;
  }> | null>(null);
  const [message, setMessage] = useState('');
  const [isBusy, setIsBusy] = useState(true);

  const selectedConnectorId = scope.connectorId;
  const tableKey = scope.tableKey;
  const columnName = scope.columnName;
  const writeMode = scope.writeMode;
  const timestampColumn = scope.timestampColumn;

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
  const validationIssues = validation?.issues ?? [];
  const validationReady = validation?.ready ?? false;
  const validationErrorCount = validationIssues.filter(
    (issue) => issue.severity === 'error',
  ).length;

  const loadData = useCallback(async () => {
    setIsBusy(true);
    try {
      const [nextConnectors, nextMappings] = await Promise.all([
        dbTargetAPI.listConnectors(),
        dbTargetAPI.listMappings(),
      ]);
      setConnectors(nextConnectors);
      setMappings(nextMappings);
      setScope((previous) => {
        const nextConnectorId =
          previous.connectorId
          && nextConnectors.some((connector) => connector.id === previous.connectorId)
            ? previous.connectorId
            : (nextConnectors[0]?.id ?? '');
        if (nextConnectorId === previous.connectorId) {
          return previous;
        }
        return {
          ...DATABASE_OUTPUT_SCOPE_INITIAL,
          connectorId: nextConnectorId,
        };
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
    setDryRunResults(null);
  }, [selectedConnectorId, selectedTagId]);
  useEffect(() => {
    if (!selectedMapping) {
      return;
    }
    const nextTableKey = `${selectedMapping.table_schema}.${selectedMapping.table_name}`;
    const nextWriteMode = selectedMapping.write_mode;
    const nextTimestampColumn =
      nextWriteMode === 'upsert'
        ? (selectedMapping.timestamp_column ?? '')
        : '';

    setScope((previous) => {
      if (
        previous.tableKey === nextTableKey
        && previous.columnName === selectedMapping.column_name
        && previous.writeMode === nextWriteMode
        && previous.timestampColumn === nextTimestampColumn
      ) {
        return previous;
      }
      return {
        ...previous,
        tableKey: nextTableKey,
        columnName: selectedMapping.column_name,
        writeMode: nextWriteMode,
        timestampColumn: nextTimestampColumn,
      };
    });
  }, [selectedMapping]);
  useEffect(() => {
    setScope((previous) => {
      if (!selectedConnectorId) {
        if (
          previous.tableKey === ''
          && previous.columnName === ''
          && previous.writeMode === 'insert'
          && previous.timestampColumn === ''
        ) {
          return previous;
        }
        return {
          ...previous,
          tableKey: '',
          columnName: '',
          writeMode: 'insert',
          timestampColumn: '',
        };
      }

      if (selectedMapping) {
        return previous;
      }

      const nextTable =
        tables.find((table) => buildTableKey(table) === previous.tableKey) ?? tables[0] ?? null;
      const nextTableKey = nextTable ? buildTableKey(nextTable) : '';
      const nextColumnName = nextTable
        ? (previous.columnName
            && nextTable.columns.some((column) => column.name === previous.columnName)
            ? previous.columnName
            : (defaultValueColumn(nextTable.columns)?.name ?? ''))
        : '';
      const nextTimestampColumn =
        previous.writeMode === 'upsert'
          ? (nextTable
              ? (previous.timestampColumn
                  && nextTable.columns.some((column) => column.name === previous.timestampColumn)
                  ? previous.timestampColumn
                  : (defaultTimestampColumn(nextTable.columns)?.name ?? ''))
              : '')
          : '';

      if (
        previous.tableKey === nextTableKey
        && previous.columnName === nextColumnName
        && previous.timestampColumn === nextTimestampColumn
      ) {
        return previous;
      }

      return {
        ...previous,
        tableKey: nextTableKey,
        columnName: nextColumnName,
        timestampColumn: nextTimestampColumn,
      };
    });
  }, [columnName, selectedConnectorId, selectedMapping, tableKey, tables, timestampColumn, writeMode]);
  const handleNewConnector = useCallback(() => {
    setScope(DATABASE_OUTPUT_SCOPE_INITIAL);
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
      setScope(() => ({
        ...DATABASE_OUTPUT_SCOPE_INITIAL,
        connectorId: connector.id,
      }));
      await loadData();
      await refreshCandidateReview();
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
  }, [canClearSavedPassword, clearSavedPassword, draft, loadData, refreshCandidateReview, selectedConnectorId, t]);
  const handleTestConnector = useCallback(async () => {
    if (!selectedConnectorId) {
      setMessage(t('workbench.output.database.results.connectorNotSaved'));
      return;
    }
    setIsBusy(true);
    try {
      const connector = await dbTargetAPI.testConnector(selectedConnectorId);
      setScope(() => ({
        ...DATABASE_OUTPUT_SCOPE_INITIAL,
        connectorId: connector.id,
      }));
      await loadData();
      await refreshCandidateReview();
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
  }, [loadData, refreshCandidateReview, selectedConnectorId, t]);
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
  const handleRefreshValidation = useCallback(async () => {
    if (!selectedConnectorId) {
      setMessage(t('workbench.output.database.results.connectorNotSaved'));
      return;
    }
    setIsBusy(true);
    try {
      const nextValidation = await dbTargetAPI.validateConnector(selectedConnectorId);
      setValidation(nextValidation);
      if (nextValidation.ready) {
        await refreshCandidateReview();
      }
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
  }, [refreshCandidateReview, selectedConnectorId, t]);
  const handleGenerateSchema = useCallback(async () => {
    if (!selectedConnectorId) {
      setMessage(t('workbench.output.database.results.connectorNotSaved'));
      return;
    }
    setIsBusy(true);
    try {
      const result = await dbTargetAPI.generateSchema(selectedConnectorId, { dry_run: false });
      await loadConnectorDetails(selectedConnectorId);
      await refreshCandidateReview();
      setMessage(
        result.executed > 0
          ? t('workbench.output.database.results.schemaGenerated')
          : t('workbench.output.database.results.schemaGenerateNoChanges'),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.database.results.schemaGenerateFailed')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadConnectorDetails, refreshCandidateReview, selectedConnectorId, t]);
  const currentTableFromKey = parseTableKey(tableKey);
  const mappedTagCount = connectorMappings.length;
  const schemaColumns = selectedTable?.columns ?? [];

  const columnMappingByColumn = useMemo(() => {
    const map = new Map<string, { mapping: DatabaseTargetMapping; candidate: WorkbenchOutputCandidate | null }>();
    for (const mapping of connectorMappings) {
      const candidate = candidates.find((c) => c.tagId === mapping.tag_id) ?? null;
      map.set(mapping.column_name, { mapping, candidate });
    }
    return map;
  }, [connectorMappings, candidates]);

  const handleColumnClick = useCallback(async (column: DatabaseTableColumn) => {
    if (!selectedConnectorId) {
      setMessage(t('workbench.output.database.results.connectorNotSaved'));
      return;
    }

    const existingBinding = columnMappingByColumn.get(column.name);

    if (existingBinding) {
      setIsBusy(true);
      try {
        await dbTargetAPI.deleteMapping(existingBinding.mapping.id);
        await loadData();
        await refreshCandidateReview();
        setDryRunResults(null);
        setMessage(
          t('workbench.output.database.results.mappingDeleted', {
            key: existingBinding.candidate?.tagKey ?? existingBinding.mapping.tag_id,
          }),
        );
      } catch (error) {
        setMessage(
          getErrorMessage(error, t('workbench.output.database.results.mappingDeleteFailed')),
        );
      } finally {
        setIsBusy(false);
      }
      return;
    }

    if (!selectedTagId) {
      setMessage(t('workbench.output.database.results.noTagSelected'));
      return;
    }

    const selectedCandidateForBind = candidates.find((c) => c.tagId === selectedTagId) ?? null;
    if (!selectedCandidateForBind) {
      return;
    }

    if (!selectedTable) {
      setMessage(t('workbench.output.database.results.tableRequired'));
      return;
    }

    setIsBusy(true);
    try {
      const mapping = await dbTargetAPI.createMapping({
        tag_id: selectedTagId,
        connector_id: selectedConnectorId,
        table_schema: selectedTable.schema,
        table_name: selectedTable.name,
        column_name: column.name,
        write_mode: writeMode,
        timestamp_column:
          writeMode === 'upsert' ? timestampColumn.trim() : '',
        enabled: true,
      });

      await loadData();
      await refreshCandidateReview();
      const dryRun = await dbTargetAPI.dryRunMappings(selectedConnectorId, {
        candidate_ids: [mapping.id],
      });
      setDryRunResults(dryRun.results);
      setMessage(
        t('workbench.output.database.results.mappingSaved', {
          key: selectedCandidateForBind.tagKey,
          table: selectedTable.name,
          column: column.name,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.database.results.mappingSaveFailed')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [candidates, columnMappingByColumn, loadData, refreshCandidateReview, selectedConnectorId, selectedTable, selectedTagId, t, timestampColumn, writeMode]);

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

      <DatabaseOutputReviewPanel
        reviewRuleId={reviewRuleId}
        reviewRevisionId={reviewRevisionId}
        reviewSet={reviewSet}
        reviewLoading={reviewLoading}
        selectedConnectorId={selectedConnectorId}
        selectedConnectorName={selectedConnector?.name ?? null}
        selectedTableKey={tableKey}
        selectedTagId={selectedTagId}
      />

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
                          setScope((previous) =>
                            previous.connectorId === connector.id
                              ? previous
                              : {
                                  ...DATABASE_OUTPUT_SCOPE_INITIAL,
                                  connectorId: connector.id,
                                });
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
                data-testid="output-database-toggle-connector"
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
                name="database-connector-kind"
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
                  data-testid="output-database-save-connector"
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
              {selectedCandidate?.tagKey ?? '—'}
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
                  name="database-mapping-table"
                  value={tableKey}
                  onChange={(event) =>
                    setScope((previous) => ({
                      ...previous,
                      tableKey: event.target.value,
                      columnName: '',
                      timestampColumn: '',
                    }))
                  }
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
              <span>{t('workbench.output.database.mapping.writeMode')}</span>
                <select
                  aria-label={t('workbench.output.database.mapping.writeMode')}
                  name="database-mapping-write-mode"
                  value={writeMode}
                  onChange={(event) =>
                    setScope((previous) => ({
                      ...previous,
                      writeMode: event.target.value as DatabaseWriteMode,
                      timestampColumn: '',
                    }))
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

            {writeMode === 'upsert' ? (
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400 lg:col-span-2">
                <span>{t('workbench.output.database.mapping.timestampColumn')}</span>
                <select
                  aria-label={t('workbench.output.database.mapping.timestampColumn')}
                  name="database-mapping-timestamp-column"
                  value={timestampColumn}
                  onChange={(event) =>
                    setScope((previous) => ({
                      ...previous,
                      timestampColumn: event.target.value,
                    }))
                  }
                  className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
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
            ) : null}
          </div>

          <div className="flex flex-wrap gap-2">
            {selectedConnectorId && tables.length === 0 ? (
              <button
                type="button"
                data-testid="output-database-generate-schema"
                onClick={() => void handleGenerateSchema()}
                disabled={isBusy}
                className="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-200 disabled:opacity-50"
              >
                {t('workbench.output.database.actions.generateSchema')}
              </button>
            ) : null}
            <button
              type="button"
              data-testid="output-database-refresh-validation"
              onClick={() => void handleRefreshValidation()}
              disabled={isBusy || !selectedConnectorId}
              className="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-200 disabled:opacity-50"
            >
              {t('workbench.output.database.actions.refreshValidation')}
            </button>
          </div>

          {dryRunResults?.length ? (
            <div
              data-testid="database-dry-run-results"
              className="space-y-2 rounded-2xl border border-slate-800 bg-slate-950/70 p-4"
            >
              {dryRunResults.map((result) => (
                <p key={`${result.candidate_id}-${result.status}`} className="text-sm text-slate-300">
                  {result.status === 'ready'
                    ? t('workbench.output.database.dryRun.ready')
                    : t('workbench.output.database.dryRun.blocked', {
                        reason: result.reason ?? result.code ?? 'unknown',
                      })}
                </p>
              ))}
            </div>
          ) : null}

          <div
            className="space-y-4"
            data-testid="database-secondary-panels"
            data-emphasis="supporting"
          >
            {schemaColumns.length > 0 ? (
              <div
                data-testid="schema-snapshot"
                tabIndex={-1}
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

                <div className="grid gap-2 md:grid-cols-2 xl:grid-cols-1" data-testid="schema-column-surface">
                  {schemaColumns.map((column) => {
                    const isRequired =
                      !column.nullable
                      && !column.primary_key
                      && column.name !== columnName
                      && (writeMode !== 'upsert' || column.name !== timestampColumn);
                    const binding = columnMappingByColumn.get(column.name);
                    const isBound = Boolean(binding);
                    const isBindTarget = Boolean(selectedTagId) && !isBound && !column.primary_key;

                    return (
                      <button
                        type="button"
                        key={column.name}
                        data-testid={`schema-column-${column.name}`}
                        data-required={isRequired ? 'true' : undefined}
                        data-bound={isBound ? 'true' : undefined}
                        disabled={isBusy || column.primary_key}
                        onClick={() => void handleColumnClick(column)}
                        className={`rounded-xl border p-3 text-left transition ${
                          isBound
                            ? 'border-violet-500/40 bg-violet-500/10 hover:border-rose-400/40'
                            : isBindTarget
                              ? 'cursor-pointer border-slate-700 bg-slate-900/60 hover:border-violet-500/40 hover:bg-violet-500/5'
                              : isRequired
                                ? 'border-amber-500/30 bg-amber-500/10'
                                : 'border-slate-800 bg-slate-900/60'
                        } ${column.primary_key ? 'cursor-not-allowed opacity-60' : ''}`}
                      >
                        <div className="flex items-center justify-between gap-3">
                          <div className="space-y-1">
                            <p className="text-sm font-semibold text-slate-50">{column.name}</p>
                            {isBound ? (
                              <p className="text-xs text-violet-200">
                                {t('workbench.output.database.schema.boundTag', {
                                  key: binding?.candidate?.tagKey ?? binding?.mapping.tag_id ?? '',
                                })}
                              </p>
                            ) : isBindTarget ? (
                              <p className="text-xs text-slate-500">
                                {t('workbench.output.database.schema.clickToBind', {
                                  key: selectedCandidate?.tagKey ?? '',
                                  column: column.name,
                                })}
                              </p>
                            ) : null}
                          </div>
                          <div className="flex flex-wrap gap-2 text-[10px] uppercase tracking-[0.14em]">
                            <span className="rounded-full bg-slate-800 px-2 py-0.5 text-slate-300">
                              {column.data_type}
                            </span>
                            {column.primary_key ? (
                              <span className="rounded-full bg-violet-500/10 px-2 py-0.5 text-violet-200">
                                PK
                              </span>
                            ) : null}
                            {isBound ? (
                              <span className="rounded-full bg-violet-500/10 px-2 py-0.5 text-violet-200">
                                ✓
                              </span>
                            ) : null}
                            {isRequired ? (
                              <span className="rounded-full bg-amber-500/10 px-2 py-0.5 text-amber-200">
                                {t('workbench.output.database.schema.required')}
                              </span>
                            ) : null}
                          </div>
                        </div>
                      </button>
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

          {validationIssues.length ? (
            <div className="space-y-2 rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
              <p className="text-sm font-semibold text-slate-50">
                {validationReady
                  ? t('workbench.output.database.validation.ready')
                  : t('workbench.output.database.validation.notReady')}
              </p>
              <ul className="space-y-2 text-sm text-slate-300">
                {validationIssues.map((issue) => (
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
