import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import type { DbConnector, DbRowGroup, Mapping, Point } from '../../state/types';

interface RowGroupPlannerProps {
  connector: DbConnector;
  points: Point[];
  mappings: Record<string, Mapping>;
  rowGroups: DbRowGroup[];
  onSetRowGroups: (rowGroups: DbRowGroup[]) => void;
  disabled?: boolean;
}

function nextRowGroupID(rowGroups: DbRowGroup[]): string {
  let index = rowGroups.length + 1;
  const existing = new Set(rowGroups.map((group) => group.id));
  while (existing.has(`row-group-${index}`)) {
    index += 1;
  }
  return `row-group-${index}`;
}

export function RowGroupPlanner({
  connector,
  points,
  mappings,
  rowGroups,
  onSetRowGroups,
  disabled = false,
}: RowGroupPlannerProps) {
  const { t } = useTranslation('workbench-v2');

  const visiblePoints = useMemo(
    () => points.filter((point) => point.enabled && mappings[point.id]),
    [points, mappings],
  );
  const scopedRowGroups = useMemo(
    () => rowGroups.filter((group) => (
      group.table_name === connector.table &&
      (group.table_schema ?? '') === (connector.schema ?? '')
    )),
    [connector.schema, connector.table, rowGroups],
  );

  const createRowGroup = () => {
    if (visiblePoints.length === 0) {
      return;
    }
    const id = nextRowGroupID(rowGroups);
    const nextGroups = [
      ...rowGroups,
      {
        id,
        connector_id: connector.connector_id,
        table_schema: connector.schema,
        table_name: connector.table,
        member_point_ids: [],
        group_key_columns: [connector.timestamp_column || 'ts'],
      },
    ];
    onSetRowGroups(nextGroups);
  };

  return (
    <section className="rounded-xl border border-slate-800 bg-slate-950/35 p-4">
      <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
        <div>
          <h3 className="text-lg font-semibold text-white">
            {t('step4.row_group_title')}
          </h3>
          <p className="mt-1 text-xs text-slate-400">
            {t('step4.row_group_subtitle')}
          </p>
        </div>
        <button
          type="button"
          data-testid="step4-row-group-create"
          disabled={disabled || visiblePoints.length === 0}
          onClick={createRowGroup}
          className="rounded-md border border-cyan-700/70 px-3 py-1.5 text-xs text-cyan-200 transition-colors hover:bg-cyan-950/40 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {t('step4.row_group_create_btn')}
        </button>
      </div>

      <div data-testid="step4-row-group-list" className="mt-4 grid gap-3">
        {scopedRowGroups.length === 0 ? (
          <div className="rounded-lg border border-dashed border-slate-800 px-3 py-4 text-xs text-slate-500">
            {t('step4.row_group_empty')}
          </div>
        ) : scopedRowGroups.map((group) => (
          <div key={group.id} className="rounded-lg border border-slate-800 bg-slate-900/35 p-3">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div className="font-mono text-xs font-semibold text-slate-100">{group.id}</div>
              <div className="font-mono text-[11px] text-slate-500">
                {(group.table_schema ? `${group.table_schema}.` : '') + group.table_name}
              </div>
            </div>
            <div className="mt-2 flex flex-wrap gap-2">
              <span className="rounded-md border border-slate-800 px-2 py-1 text-[11px] text-slate-300">
                {t('step4.row_group_member_count', {
                  count: group.member_point_ids.length,
                })}
              </span>
              <span className="rounded-md border border-slate-800 px-2 py-1 font-mono text-[11px] text-slate-400">
                {(group.group_key_columns ?? ['ts']).join(', ')}
              </span>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
