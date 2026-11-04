import { useId, useMemo, useState, useSyncExternalStore, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { ParameterHelp } from '../components/ParameterHelp';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import type { EditorMember } from '../state/writeGroup/draft';
import type { Step4TargetMetadataStatus } from '../steps/step4/useStep4TargetColumns';
import { comparePointAddresses, sourceKey } from './recordingPickerState';
import type { PointOrderObservation, RecordingPointOrderPort } from './useRecordingPointOrder';
import './recordingCompactMembers.css';
interface CompactRow { candidate: GroupCandidate; member: EditorMember | undefined; problem: string | null }
interface Props {
  rows: CompactRow[]; total: number; query: string; setQuery: (query: string) => void; problemsOnly: boolean; setProblemsOnly: (value: boolean) => void; metadataStatus: Step4TargetMetadataStatus;
  onRetryMetadata?: () => void; metadataAtDestination?: boolean; scopeInHeading?: boolean;
  disabled: boolean; showEntityKey: boolean;
  controls: { removeLabel: (member: EditorMember) => string; onRemove: (member: EditorMember) => void };
  update: (candidate: GroupCandidate, patch: Partial<EditorMember>) => void;
  renderCalculation?: (member: EditorMember, tableCells?: boolean) => ReactNode;
  renderSource?: (member: EditorMember, inline?: boolean) => ReactNode;
  renderName?: (candidate: GroupCandidate) => ReactNode; nameLocked?: boolean;
  pointOrderSort?: { workspaceId: string; port: RecordingPointOrderPort };
}
type SortColumn = 'pointOnly' | 'device' | 'wordOrder' | 'addressOnly' | 'dataType' | 'calculation' | 'interval';
const noSubscribe = () => () => undefined;
const noObservations: readonly PointOrderObservation[] = [];
const collator = new Intl.Collator(undefined, { numeric: true, sensitivity: 'base' });
/** Presentation only. Rows, problems, exact tuple updates and membership all belong to GroupMemberTable. */
export function RecordingCompactMemberTable({ rows, total, query, setQuery, problemsOnly, setProblemsOnly, metadataStatus, onRetryMetadata, metadataAtDestination, scopeInHeading, disabled, showEntityKey, controls, update, renderCalculation, renderSource, renderName, nameLocked, pointOrderSort }: Props) {
  const { t } = useTranslation('workbench-v2');
  const { t: text } = useTranslation('recording-editor');
  const requiredHint = useId();
  const [sort, setSort] = useState<{ column: SortColumn; direction: 'ascending' | 'descending' }>({ column: 'addressOnly', direction: 'ascending' });
  const observedOrders = useSyncExternalStore(pointOrderSort?.port.subscribe ?? noSubscribe, () => pointOrderSort?.port.observations?.() ?? noObservations, () => noObservations);
  const orderedRows = useMemo(() => {
    const observed = new Map(observedOrders.map(({ scope, entry }) => [JSON.stringify([scope.workspace_id, scope.device_id, scope.point_id]), entry]));
    const orderValue = (row: CompactRow) => {
      if (!pointOrderSort || !row.member) return '';
      const scope = { workspace_id: pointOrderSort.workspaceId, device_id: row.member.device_id, point_id: row.member.point_id };
      const current = pointOrderSort.port.get(scope), matching = observed.get(JSON.stringify([scope.workspace_id, scope.device_id, scope.point_id]));
      const entry = current === matching ? matching : current;
      const resolution = entry?.saved?.resolution ?? entry?.snapshot?.resolution;
      if (entry?.intent) return entry.intent.format_mode === 'override' ? entry.intent.configured_format : 'inherit';
      return resolution?.resolved && resolution.applicable ? resolution.effective_format : '';
    };
    const value = (row: CompactRow): string => {
      switch (sort.column) {
        case 'pointOnly': return row.candidate.label;
        case 'device': return row.candidate.device_name;
        case 'wordOrder': return orderValue(row);
        case 'addressOnly': return row.candidate.address;
        case 'dataType': return row.candidate.target_type;
        case 'calculation': return row.member?.calculation ? text(`calculation.${row.member.calculation}`) : '';
        case 'interval': {
          const cadence = row.member?.cadence;
          return cadence && Number.isFinite(cadence.interval) ? String(cadence.interval * (cadence.unit === 'minute' ? 60 : 1)) : '';
        }
      }
    };
    return [...rows].sort((a, b) => {
      const order = sort.column === 'addressOnly' ? comparePointAddresses(value(a), value(b)) : collator.compare(value(a), value(b));
      return (sort.direction === 'ascending' ? order : -order) || (sourceKey(a.candidate) < sourceKey(b.candidate) ? -1 : sourceKey(a.candidate) > sourceKey(b.candidate) ? 1 : 0);
    });
  }, [rows, sort, pointOrderSort, observedOrders, text]);
  const toggleSort = (column: SortColumn) => setSort(current => ({ column, direction: current.column === column && current.direction === 'ascending' ? 'descending' : 'ascending' }));
  return <div className="recording-compact-members" data-testid="group-member-table" data-calculation-scope={renderCalculation ? 'member' : 'group'} data-show-entity-key={showEntityKey}>
    <div className="recording-compact-members__filters">
      <label className="recording-compact-members__search">{text('compact.search')}<input type="search" placeholder={text('compact.searchPlaceholder')} value={query} disabled={nameLocked} onChange={event => setQuery(event.target.value)} data-testid="group-member-search" /></label>
      <label className="recording-compact-members__problems"><input type="checkbox" checked={problemsOnly} disabled={nameLocked} onChange={event => setProblemsOnly(event.target.checked)} data-testid="group-member-problems-only" />{t('step4.group.members.problems_only')}</label>
      {!scopeInHeading && <span role="status" data-testid="group-member-scope">{text('compact.selectedScope', { shown: rows.length, total })}</span>}
    </div>
    <div className="recording-compact-members__guidance">{renderName && <p className="recording-compact-members__name-hint">{text('name.shared')}</p>}
      <div className="recording-compact-members__help"><span id={requiredHint}>{text('compact.requiredInfo')}</span><ParameterHelp compact help="member.required" />
        <span>{text('compact.orderInfo')}</span><ParameterHelp compact help="point.order" /></div></div>
    {!metadataAtDestination && metadataStatus !== 'exists' && <div className="recording-compact-members__metadata"><p role="status" data-testid="group-member-metadata-note">{metadataStatus === 'failed' ? text('compact.metadataFailed') : t(`step4.group.members.metadata_${metadataStatus}`)}</p>
      {metadataStatus === 'failed' && onRetryMetadata && <button type="button" className="gw-secondary" onClick={onRetryMetadata} data-testid="group-metadata-retry">{t('step4.group.retry')}</button>}</div>}
    {rows.length > 0 && <p className="recording-compact-members__horizontal-hint">{text('compact.horizontalHint')}</p>}
    <div className="recording-compact-members__scroll recording-selected-table-scroll" role="region" aria-label={t('step4.group.members.region')} tabIndex={0}>
      <table><thead><tr>{(['device', 'pointOnly', 'addressOnly', 'wordOrder', 'dataType', ...(renderCalculation ? ['calculation', 'interval'] as const : []), 'actions'] as const).map(key => <th key={key} scope="col" data-column={key} aria-sort={key !== 'actions' && sort.column === key ? sort.direction : undefined}>
        {key === 'actions' ? text('compact.actions') : <button type="button" className="recording-compact-members__sort" disabled={nameLocked} onClick={() => toggleSort(key)}>{text(`compact.${key}`)}<span aria-hidden="true">{sort.column === key ? sort.direction === 'ascending' ? '▲' : '▼' : '↕'}</span></button>}
      </th>)}</tr></thead>
        <tbody>{rows.length === 0 && <tr><td colSpan={renderCalculation ? 8 : 6} data-testid="group-member-empty">{t(total ? 'step4.group.members.none_shown' : 'step4.group.members.none')}</td></tr>}{orderedRows.map(({ candidate, member, problem }) => {
          if (!member) return null;
          const describedBy = problem ? `group-member-problem-${candidate.key}` : undefined;
          return <tr key={candidate.key} data-testid={`group-member-row-${candidate.point_id}`} data-problem={problem ?? ''}>
            <td data-column="device" title={candidate.device_name}>{candidate.device_name}</td>
            <td data-column="pointOnly">{renderName?.(candidate) ?? <strong>{candidate.label}</strong>}
              {problem && <p id={describedBy} role="alert" data-testid={`group-member-problem-${candidate.point_id}`} className="recording-compact-members__problem">{t(`step4.group.members.problem_${problem.replace('-', '_')}`)}</p>}
            </td>
            <td data-column="addressOnly" title={candidate.address}><span className="font-mono">{candidate.address}</span></td>
            <td data-column="wordOrder">{renderSource?.(member, true) ?? '—'}</td><td data-column="dataType" title={candidate.target_type}>{candidate.target_type}</td>
            {renderCalculation?.(member, true)}
            <td data-column="actions"><div className="recording-compact-members__actions"><button type="button" className="gw-secondary recording-compact-members__required" aria-label={t('step4.group.members.required_label', { tag: candidate.label })} aria-describedby={requiredHint} aria-pressed={member.required} disabled={disabled} onClick={() => update(candidate, { required: !member.required })}>{text(member.required ? 'compact.required' : 'compact.optional')}</button>
              <button type="button" className="gw-secondary" disabled={disabled || nameLocked} aria-label={`${candidate.label} · ${controls.removeLabel(member)}`} onClick={() => controls.onRemove(member)}>{text('compact.remove')}</button>
              {showEntityKey && <label className="recording-compact-members__entity">{t('step4.group.members.col_entity')}<input aria-label={t('step4.group.members.entity_label', { tag: candidate.label })} value={member.entity_key} disabled={disabled}
                onChange={event => update(candidate, { entity_key: event.target.value })} /></label>}
              </div>
            </td>
          </tr>;
        })}</tbody>
      </table>
    </div>
  </div>;
}
