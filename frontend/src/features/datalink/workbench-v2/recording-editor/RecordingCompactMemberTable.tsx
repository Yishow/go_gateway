import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import type { DbColumn } from '../state/dbSchemas';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import type { EditorMember } from '../state/writeGroup/draft';
import type { Step4TargetMetadataStatus } from '../steps/step4/useStep4TargetColumns';
import './recordingCompactMembers.css';
interface CompactRow { candidate: GroupCandidate; member: EditorMember | undefined; problem: string | null; suggestion?: { column: string; status: string } }
interface Props {
  rows: CompactRow[]; total: number; query: string; setQuery: (query: string) => void; problemsOnly: boolean; setProblemsOnly: (value: boolean) => void; columns: DbColumn[]; metadataStatus: Step4TargetMetadataStatus;
  columnsAssignedOnSave: boolean; disabled: boolean; showEntityKey: boolean;
  controls: { removeLabel: (member: EditorMember) => string; onRemove: (member: EditorMember) => void };
  update: (candidate: GroupCandidate, patch: Partial<EditorMember>) => void;
  renderCalculation?: (member: EditorMember, tableCells?: boolean) => ReactNode;
  renderSource?: (member: EditorMember) => ReactNode;
}
/** Presentation only. Rows, problems, exact tuple updates and membership all belong to GroupMemberTable. */
export function RecordingCompactMemberTable({ rows, total, query, setQuery, problemsOnly, setProblemsOnly, metadataStatus, columnsAssignedOnSave, disabled, showEntityKey, controls, update, renderCalculation, renderSource }: Props) {
  const { t } = useTranslation('workbench-v2');
  const { t: text } = useTranslation('recording-editor');
  return <div className="recording-compact-members" data-testid="group-member-table">
    <details className="recording-compact-members__filters"><summary>{t('step4.group.members.search')} · <span role="status" data-testid="group-member-scope">{t('step4.group.members.scope', { shown: rows.length, total })}</span></summary>
      <label>{t('step4.group.members.search')}<input type="search" value={query} onChange={event => setQuery(event.target.value)} data-testid="group-member-search" /></label>
      <label><input type="checkbox" checked={problemsOnly} onChange={event => setProblemsOnly(event.target.checked)} data-testid="group-member-problems-only" />{t('step4.group.members.problems_only')}</label>
    </details>
    {metadataStatus !== 'exists' && <p role="status" data-testid="group-member-metadata-note">{t(`step4.group.members.metadata_${metadataStatus}`)}</p>}
    <div className="recording-compact-members__scroll" role="region" aria-label={t('step4.group.members.region')} tabIndex={0}>
      <table><thead><tr>{['point', 'address', 'calculation', 'interval', 'actions'].map(key => <th key={key} scope="col">{text(`compact.${key}`)}</th>)}</tr></thead>
        <tbody>{rows.length === 0 && <tr><td colSpan={5} data-testid="group-member-empty">{t(total ? 'step4.group.members.none_shown' : 'step4.group.members.none')}</td></tr>}{rows.map(({ candidate, member, problem, suggestion }) => {
          if (!member) return null;
          const describedBy = problem ? `group-member-problem-${candidate.key}` : undefined;
          const chosenColumn = member.target_column;
          return <tr key={candidate.key} data-testid={`group-member-row-${candidate.point_id}`} data-problem={problem ?? ''}>
            <td><strong>{candidate.label}</strong><small>{candidate.device_name}</small>{renderSource?.(member)}
              {problem && <p id={describedBy} role="alert" data-testid={`group-member-problem-${candidate.point_id}`} className="recording-compact-members__problem">{t(`step4.group.members.problem_${problem.replace('-', '_')}`)}</p>}
              {!problem && <small data-testid={`group-member-status-${candidate.point_id}`}>
                {columnsAssignedOnSave ? text('simple.pendingColumn') : suggestion?.status === 'suggested' && (member.target_column === suggestion.column || member.target_column === '')
                  ? t('step4.group.members.suggested', { column: suggestion.column }) : member.target_column ? t('step4.group.members.confirmed') : t('step4.group.members.unassigned')}
              </small>}
            </td>
            <td><span className="font-mono">{candidate.address}</span><small>{candidate.target_type}</small></td>
            {renderCalculation ? renderCalculation(member, true) : <><td data-testid={`group-member-calculation-${candidate.point_id}`}>{text('compact.groupRule')}</td><td>{text('compact.groupCadence')}</td></>}
            <td><button type="button" className="gw-secondary" disabled={disabled} aria-label={`${candidate.label} · ${controls.removeLabel(member)}`} onClick={() => controls.onRemove(member)}>{text('compact.remove')}</button>
              <details data-testid={`group-member-settings-${candidate.point_id}`} className="recording-compact-members__settings">
                <summary>{text('compact.settings')}</summary>
                <label>{t('step4.group.members.col_column')}
                  <span data-testid={`group-member-column-${candidate.point_id}`}>{text('compact.systemManaged')}</span>
                  <code className="recording-compact-members__binding">{chosenColumn || text('simple.pendingColumn')}</code>
                </label>
                <label><input type="checkbox" aria-label={t('step4.group.members.required_label', { tag: candidate.label })} checked={member.required} disabled={disabled} onChange={event => update(candidate, { required: event.target.checked })} />{t('step4.group.members.col_required')}</label>
                {showEntityKey && <label>{t('step4.group.members.col_entity')}<input aria-label={t('step4.group.members.entity_label', { tag: candidate.label })} value={member.entity_key} disabled={disabled}
                  onChange={event => update(candidate, { entity_key: event.target.value })} /></label>}
              </details>
            </td>
          </tr>;
        })}</tbody>
      </table>
    </div>
  </div>;
}
