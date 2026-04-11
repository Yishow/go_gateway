import { useTranslation } from 'react-i18next';
import type { DatabaseTableColumn, DatabaseTableInfo, DatabaseWriteMode } from '../../../types/datalink';
import type { WorkbenchOutputCandidate } from './workbenchOutputTypes';
import type { DatabaseRowPlan, DatabaseRowPlanIssueCode, DatabaseRowPlanMember } from './databaseRowPlannerModel';
import { buildTableKey } from './databaseTargetBoardUtils';

type DatabaseGroupedRowPlannerProps = {
  rowPlans: ReadonlyArray<DatabaseRowPlan>;
  selectedCandidate: WorkbenchOutputCandidate | null;
  selectedConnectorId: string;
  tables: ReadonlyArray<DatabaseTableInfo>;
  tableKey: string;
  writeMode: DatabaseWriteMode;
  timestampColumn: string;
  tableColumns: ReadonlyArray<DatabaseTableColumn>;
  isBusy: boolean;
  onChangeTable: (value: string) => void;
  onChangeWriteMode: (value: DatabaseWriteMode) => void;
  onChangeTimestampColumn: (value: string) => void;
  onChangeRowGroup: (rowPlan: DatabaseRowPlan, value: string) => void;
  onChangeRowInterval: (rowPlan: DatabaseRowPlan, value: string) => void;
  onChangeMemberColumn: (
    rowPlan: DatabaseRowPlan,
    member: DatabaseRowPlanMember,
    value: string,
  ) => void;
  onApplyRow: (rowPlan: DatabaseRowPlan) => void;
};

function getIssueCopy(code: DatabaseRowPlanIssueCode, t: ReturnType<typeof useTranslation>['t']) {
  switch (code) {
    case 'candidate_blocked':
      return t('workbench.output.database.rowPlanner.issue.candidate_blocked', {
        defaultValue: 'A member is already blocked by the current review set.',
      });
    case 'tag_pending':
      return t('workbench.output.database.rowPlanner.issue.tag_pending', {
        defaultValue: 'Apply the Tag review decisions before publishing this row.',
      });
    case 'connector_required':
      return t('workbench.output.database.rowPlanner.issue.connector_required', {
        defaultValue: 'Choose a connector before applying this row.',
      });
    case 'table_required':
      return t('workbench.output.database.rowPlanner.issue.table_required', {
        defaultValue: 'Choose a target table before applying this row.',
      });
    case 'timestamp_required':
      return t('workbench.output.database.rowPlanner.issue.timestamp_required', {
        defaultValue: 'Choose a timestamp column for upsert rows.',
      });
    case 'column_required':
      return t('workbench.output.database.rowPlanner.issue.column_required', {
        defaultValue: 'Every row member needs a column name before apply.',
      });
    case 'connector_mismatch':
      return t('workbench.output.database.rowPlanner.issue.connector_mismatch', {
        defaultValue: 'Members point at different connectors.',
      });
    case 'table_mismatch':
      return t('workbench.output.database.rowPlanner.issue.table_mismatch', {
        defaultValue: 'Members point at different tables.',
      });
    case 'write_mode_mismatch':
      return t('workbench.output.database.rowPlanner.issue.write_mode_mismatch', {
        defaultValue: 'Members disagree on write mode.',
      });
    case 'timestamp_mismatch':
      return t('workbench.output.database.rowPlanner.issue.timestamp_mismatch', {
        defaultValue: 'Members disagree on timestamp column.',
      });
    case 'interval_mismatch':
      return t('workbench.output.database.rowPlanner.issue.interval_mismatch', {
        defaultValue: 'Members must share the same effective interval.',
      });
  }
}

export function DatabaseGroupedRowPlanner({
  rowPlans,
  selectedCandidate,
  selectedConnectorId,
  tables,
  tableKey,
  writeMode,
  timestampColumn,
  tableColumns,
  isBusy,
  onChangeTable,
  onChangeWriteMode,
  onChangeTimestampColumn,
  onChangeRowGroup,
  onChangeRowInterval,
  onChangeMemberColumn,
  onApplyRow,
}: DatabaseGroupedRowPlannerProps) {
  const { t } = useTranslation();
  const readyRowCount = rowPlans.filter((rowPlan) => rowPlan.status === 'ready').length;

  return (
    <div className="space-y-4">
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
            onChange={(event) => onChangeTable(event.target.value)}
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
            onChange={(event) => onChangeWriteMode(event.target.value as DatabaseWriteMode)}
            className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
          >
            <option value="insert">{t('workbench.output.database.writeMode.insert')}</option>
            <option value="upsert">{t('workbench.output.database.writeMode.upsert')}</option>
          </select>
        </label>

        {writeMode === 'upsert' ? (
          <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400 lg:col-span-2">
            <span>{t('workbench.output.database.mapping.timestampColumn')}</span>
            <select
              aria-label={t('workbench.output.database.mapping.timestampColumn')}
              name="database-mapping-timestamp-column"
              value={timestampColumn}
              onChange={(event) => onChangeTimestampColumn(event.target.value)}
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

      <section
        className="space-y-4 rounded-2xl border border-violet-500/20 bg-slate-950/70 p-4"
        data-testid="database-grouped-row-planner"
      >
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div className="space-y-1">
            <p className="text-xs font-semibold uppercase tracking-[0.18em] text-violet-300">
              {t('workbench.output.database.rowPlanner.eyebrow', {
                defaultValue: 'Grouped row planner',
              })}
            </p>
            <h4 className="text-lg font-semibold text-slate-50">
              {t('workbench.output.database.rowPlanner.title', {
                defaultValue: 'Plan grouped database rows before apply',
              })}
            </h4>
            <p className="text-sm text-slate-300">
              {t('workbench.output.database.rowPlanner.description', {
                defaultValue:
                  'Use grouped rows as the first-pass workflow. Adjust group key, member columns, and interval before applying the whole row.',
              })}
            </p>
          </div>
          <div className="rounded-full border border-violet-500/30 bg-violet-500/10 px-3 py-1 text-xs font-semibold text-violet-100">
            {t('workbench.output.database.rowPlanner.summary', {
              defaultValue: '{{ready}} ready / {{total}} rows',
              ready: readyRowCount,
              total: rowPlans.length,
            })}
          </div>
        </div>

        {rowPlans.length === 0 ? (
          <p className="rounded-xl border border-dashed border-slate-700 bg-slate-900/60 px-3 py-4 text-sm text-slate-400">
            {t('workbench.output.database.rowPlanner.empty', {
              defaultValue: 'No database row plans are available for the current rule.',
            })}
          </p>
        ) : (
          <div className="grid gap-3">
            {rowPlans.map((rowPlan) => {
              const applyBlocked =
                rowPlan.status !== 'ready' || !selectedConnectorId || tableKey === '';
              return (
                <article
                  key={rowPlan.id}
                  data-testid={`database-row-plan-${rowPlan.id}`}
                  className={`space-y-3 rounded-2xl border p-4 ${
                    rowPlan.status === 'ready'
                      ? 'border-violet-500/30 bg-violet-500/10'
                      : 'border-amber-500/30 bg-amber-500/10'
                  }`}
                >
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <p className="text-sm font-semibold text-slate-50">
                        {rowPlan.displayKey}
                      </p>
                      <p className="text-xs text-slate-300">
                        {t('workbench.output.database.rowPlanner.memberCount', {
                          defaultValue: '{{count}} members',
                          count: rowPlan.members.length,
                        })}
                      </p>
                    </div>
                    <span className="rounded-full bg-slate-950/70 px-2 py-1 text-[11px] font-semibold text-slate-200">
                      {rowPlan.status === 'ready'
                        ? t('workbench.output.database.rowPlanner.status.ready', {
                            defaultValue: 'ready',
                          })
                        : t('workbench.output.database.rowPlanner.status.blocked', {
                            defaultValue: 'blocked',
                          })}
                    </span>
                  </div>

                  <div className="grid gap-3 md:grid-cols-[minmax(0,1fr)_160px]">
                    <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                      <span>
                        {t('workbench.output.database.rowPlanner.groupKey', {
                          defaultValue: 'Group key',
                        })}
                      </span>
                      <input
                        data-testid={`database-row-group-${rowPlan.id}`}
                        value={rowPlan.groupKey ?? ''}
                        onChange={(event) => onChangeRowGroup(rowPlan, event.target.value)}
                        className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                      />
                    </label>
                    <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                      <span>
                        {t('workbench.output.database.rowPlanner.interval', {
                          defaultValue: 'Interval (s)',
                        })}
                      </span>
                      <input
                        data-testid={`database-row-interval-${rowPlan.id}`}
                        type="number"
                        min="1"
                        value={rowPlan.intervalSeconds ?? ''}
                        onChange={(event) => onChangeRowInterval(rowPlan, event.target.value)}
                        className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                      />
                    </label>
                  </div>

                  {rowPlan.issueCodes.length > 0 ? (
                    <ul className="space-y-2">
                      {rowPlan.issueCodes.map((issueCode) => (
                        <li
                          key={`${rowPlan.id}-${issueCode}`}
                          className="rounded-xl border border-amber-500/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-100"
                        >
                          {getIssueCopy(issueCode, t)}
                        </li>
                      ))}
                    </ul>
                  ) : null}

                  <div className="grid gap-3 md:grid-cols-2">
                    {rowPlan.members.map((member) => (
                      <div
                        key={member.candidateId}
                        className="rounded-xl border border-slate-800 bg-slate-950/70 p-3"
                      >
                        <p className="text-sm font-semibold text-slate-50">{member.tagKey}</p>
                        <p className="mt-1 text-xs text-slate-400">
                          {member.displayName} · {member.address} · {member.dataType}
                        </p>
                        <label className="mt-3 block space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                          <span>
                            {t('workbench.output.database.rowPlanner.column', {
                              defaultValue: 'Column name',
                            })}
                          </span>
                          <input
                            data-testid={`database-row-column-${member.candidateId}`}
                            value={member.columnName}
                            onChange={(event) =>
                              onChangeMemberColumn(rowPlan, member, event.target.value)
                            }
                            className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                          />
                        </label>
                      </div>
                    ))}
                  </div>

                  <div className="flex justify-end">
                    <button
                      type="button"
                      data-testid={`database-row-apply-${rowPlan.id}`}
                      disabled={applyBlocked || isBusy}
                      onClick={() => onApplyRow(rowPlan)}
                      className="rounded-lg bg-violet-500 px-3 py-2 text-sm font-semibold text-slate-950 disabled:opacity-50"
                    >
                      {t('workbench.output.database.rowPlanner.apply', {
                        defaultValue: 'Apply row',
                      })}
                    </button>
                  </div>
                </article>
              );
            })}
          </div>
        )}
      </section>
    </div>
  );
}
