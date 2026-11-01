import { ParameterHelp } from '../../../components/ParameterHelp';
import * as React from 'react';
import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { DbColumn } from '../../../state/dbSchemas';
import type { Step4TargetMetadataStatus } from '../useStep4TargetColumns';
import type { ExcludedCandidate, GroupCandidate } from '../../../state/writeGroup/candidates';
import { columnCompatibility, type ColumnSuggestions, type WriteGroupSqlDialect } from '../../../state/writeGroup/columns';
import { memberKey, type DraftIssue, type EditorMember } from '../../../state/writeGroup/draft';

export interface GroupMemberTableProps {
  presentation?: 'guided' | 'simple';
  candidates: GroupCandidate[];
  excluded: ExcludedCandidate[];
  members: EditorMember[];
  columns: DbColumn[];
  metadataStatus: Step4TargetMetadataStatus;
  suggestions: ColumnSuggestions;
  issues: DraftIssue[];
  showEntityKey: boolean;
  disabled: boolean;
  /** Exact-value SQL dialect used when checking persisted metadata. */
  dialect?: WriteGroupSqlDialect;
  managed?: boolean;
  /** Fresh managed intent has no server-owned column assignment until save. Presentation only. */
  columnsAssignedOnSave?: boolean;
  /** Picker owns membership when supplied; this table edits only already selected rows. */
  selectionControls?: { removeLabel: (member: EditorMember) => string; onRemove: (member: EditorMember) => void };
  /** Optional visible result column; the editor remains the calculation owner. */
  renderCalculation?: (member: EditorMember) => React.ReactNode;
  renderSource?: (member: EditorMember) => React.ReactNode;
  onChange: (members: EditorMember[]) => void;
}

interface Row {
  candidate: GroupCandidate;
  member: EditorMember | undefined;
  suggestion: { column: string; status: string } | undefined;
  problem: string | null;
}

function emptyMember(candidate: GroupCandidate, column: string): EditorMember {
  return {
    key: memberKey(candidate.point_id, candidate.tag_id), device_id: candidate.device_id,
    point_id: candidate.point_id, tag_id: candidate.tag_id, entity_key: '', target_column: column, required: true,
  };
}

/**
 * Chooses which saved Tags belong to the group and which real column each one
 * writes to. Search and the problems-only filter narrow what is shown; every
 * bulk action names its scope and count and touches only the rows shown.
 */
export const GroupMemberTable: React.FC<GroupMemberTableProps> = ({
  presentation = 'guided', renderSource, renderCalculation, candidates, excluded, members, columns, metadataStatus, suggestions, issues, showEntityKey, disabled,
  dialect = 'portable', managed = false, columnsAssignedOnSave = false, onChange, selectionControls,
}) => {
  const { t } = useTranslation('workbench-v2');
  const { t: editorText } = useTranslation('recording-editor');
  const [query, setQuery] = useState('');
  const [problemsOnly, setProblemsOnly] = useState(false);
  const byKey = useMemo(() => new Map(members.map((member) => [member.key, member])), [members]);
  const metadataKnown = metadataStatus === 'exists';
  const columnByName = useMemo(() => new Map(columns.map((column) => [column.name, column])), [columns]);
  const issuesByMember = useMemo(() => {
    const map = new Map<string, DraftIssue[]>();
    for (const issue of issues) {
      if (!issue.member_key) continue;
      map.set(issue.member_key, [...(map.get(issue.member_key) ?? []), issue]);
    }
    return map;
  }, [issues]);
  const conflictColumns = useMemo(() => new Set(issues.filter((issue) => issue.code === 'column-conflict').map((issue) => issue.column)), [issues]);

  const selectedCandidates = selectionControls ? candidates.filter(candidate => members.some(member =>
    member.key === candidate.key && member.device_id === candidate.device_id && member.point_id === candidate.point_id && member.tag_id === candidate.tag_id)) : candidates;
  const rows: Row[] = selectedCandidates.map((candidate) => {
    const member = selectionControls ? members.find(value => value.key === candidate.key && value.device_id === candidate.device_id &&
      value.point_id === candidate.point_id && value.tag_id === candidate.tag_id) : byKey.get(candidate.key);
    const suggestion = suggestions.assignments[candidate.key];
    let problem: string | null = null;
    if (member) {
      const column = columnByName.get(member.target_column);
      if (issuesByMember.has(candidate.key)) problem = 'column-required';
      else if (conflictColumns.has(member.target_column)) problem = 'column-conflict';
      else if (metadataKnown && column && columnCompatibility(candidate.target_type, column, dialect, managed) === 'incompatible') problem = 'column-incompatible';
      else if (metadataKnown && member.target_column && !column) problem = 'column-missing';
    }
    return { candidate, member, suggestion, problem };
  });

  const needle = query.trim().toLowerCase();
  const shown = rows.filter((row) => {
    if (problemsOnly && !row.problem) return false;
    if (!needle) return true;
    const haystack = [row.candidate.label, row.candidate.device_name, row.candidate.address, row.member?.target_column ?? '']
      .join(' ').toLowerCase();
    return haystack.includes(needle);
  });
  const shownKeys = new Set(shown.map((row) => row.candidate.key));

  const update = (candidate: GroupCandidate, patch: Partial<EditorMember>) =>
    onChange(members.map(member => member.key === candidate.key && (!selectionControls ||
      member.device_id === candidate.device_id && member.point_id === candidate.point_id && member.tag_id === candidate.tag_id)
      ? { ...member, ...patch } : member));
  const include = (candidate: GroupCandidate, column = '') => {
    if (byKey.has(candidate.key)) return;
    onChange([...members, emptyMember(candidate, column)]);
  };
  const exclude = (key: string) => onChange(members.filter((member) => member.key !== key));

  const includeShown = () => {
    const additions = shown.filter((row) => !row.member).map((row) => emptyMember(row.candidate, ''));
    onChange([...members, ...additions]);
  };
  const excludeShown = () => onChange(members.filter((member) => !shownKeys.has(member.key)));
  const acceptShown = () => {
    const next = [...members];
    for (const row of shown) {
      if (row.suggestion?.status !== 'suggested') continue;
      const existing = next.findIndex((member) => member.key === row.candidate.key);
      if (existing >= 0) {
        if (next[existing].target_column === '') next[existing] = { ...next[existing], target_column: row.suggestion.column };
      } else {
        next.push(emptyMember(row.candidate, row.suggestion.column));
      }
    }
    onChange(next);
  };

  const includable = shown.filter((row) => !row.member).length;
  const includedShown = shown.filter((row) => row.member).length;
  const acceptable = shown.filter((row) => row.suggestion?.status === 'suggested' &&
    (!row.member || row.member.target_column === '')).length;

  if (presentation === 'simple' && selectionControls && members.length === 0) return <div data-testid="group-member-table">
    <p className="recording-editor-empty" data-testid="group-member-empty">{editorText('simple.noSelectedPoints')}</p>
    {metadataStatus !== 'exists' && <p role="status" data-testid="group-member-metadata-note">{t(`step4.group.members.metadata_${metadataStatus}`)}</p>}
  </div>;
  const classes = (value: string) => value;
  return (
    <div className={classes('space-y-3')} data-testid="group-member-table">
      <div className={classes('flex flex-wrap items-end gap-3')}>
        <label className={classes('flex flex-col text-sm text-[var(--gw-muted)]')}>
          {t('step4.group.members.search')}
          <input
            type="search" value={query} onChange={(event) => setQuery(event.target.value)}
            className={classes('mt-1 rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1 text-sm text-[var(--gw-text)]')}
            data-testid="group-member-search"
          />
        </label>
        <label className={classes('flex items-center gap-2 text-sm text-[var(--gw-text)]')}>
          <input type="checkbox" checked={problemsOnly} onChange={(event) => setProblemsOnly(event.target.checked)} data-testid="group-member-problems-only" />
          {t('step4.group.members.problems_only')}
        </label>
        {!selectionControls && <div className={classes('flex flex-wrap gap-2')} role="group" aria-label={t('step4.group.members.bulk_label')}>
          <button type="button" disabled={disabled || includable === 0} onClick={includeShown} className={classes('rounded border border-[var(--gw-control)] px-2 py-1 text-sm text-[var(--gw-text)] disabled:opacity-40')} data-testid="group-bulk-include">
            {t('step4.group.members.bulk_include', { count: includable })}
          </button>
          <button type="button" disabled={disabled || includedShown === 0} onClick={excludeShown} className={classes('rounded border border-[var(--gw-control)] px-2 py-1 text-sm text-[var(--gw-text)] disabled:opacity-40')} data-testid="group-bulk-exclude">
            {t('step4.group.members.bulk_exclude', { count: includedShown })}
          </button>
          <button type="button" disabled={disabled || acceptable === 0} onClick={acceptShown} className={classes('rounded border border-[var(--gw-focus)] px-2 py-1 text-sm text-[var(--gw-accent)] disabled:opacity-40')} data-testid="group-bulk-accept">
            {t('step4.group.members.bulk_accept', { count: acceptable })}
          </button>
        </div>}
        <span className={classes('text-sm text-[var(--gw-muted)]')} role="status" data-testid="group-member-scope">
          {t('step4.group.members.scope', { shown: shown.length, total: rows.length })}
        </span>
      </div>

      {columnsAssignedOnSave && <p className="text-sm text-[var(--gw-muted)]">{editorText('simple.columnsAssignedOnSave')}</p>}
      {metadataStatus !== 'exists' && (
        <p className={classes('text-sm text-[var(--gw-warning-fg)]')} data-testid="group-member-metadata-note">
          {t(`step4.group.members.metadata_${metadataStatus}`)}
        </p>
      )}

      <div className={classes('max-w-full overflow-x-auto rounded border border-[var(--gw-control)]')} role="region" aria-label={t('step4.group.members.region')} tabIndex={0}>
        <table className={classes('w-full min-w-[640px] border-collapse text-left text-sm')}>
          <thead className={classes('bg-[var(--gw-surface)] text-[var(--gw-muted)]')}>
            <tr>
              <th scope="col" className={classes('p-2')}>{t('step4.group.members.col_include')} <ParameterHelp help="member.include" /></th>
              <th scope="col" className={classes('p-2')}>{t('step4.group.members.col_tag')}</th>
              {renderCalculation ? <th scope="col" className={classes('p-2')}>{editorText('simple.calculationColumn')}</th> : <th scope="col" className={classes('p-2')}>{t('step4.group.members.col_type')}</th>}
              <th scope="col" className={classes('p-2')}>{t('step4.group.members.col_column')} <ParameterHelp help="member.targetColumn" /></th>
              <th scope="col" className={classes('p-2')}>{t('step4.group.members.col_required')} <ParameterHelp help="member.required" /></th>
              {showEntityKey && <th scope="col" className={classes('p-2')}>{t('step4.group.members.col_entity')} <ParameterHelp help="member.entityKey" /></th>}
              <th scope="col" className={classes('p-2')}>{t('step4.group.members.col_status')}</th>
            </tr>
          </thead>
          <tbody className={classes('divide-y divide-[var(--gw-control)] text-[var(--gw-text)]')}>
            {shown.length === 0 && (
              <tr><td colSpan={showEntityKey ? 7 : 6} className={classes('p-4 text-center text-[var(--gw-muted)]')} data-testid="group-member-empty">
                {rows.length === 0 ? t('step4.group.members.none') : t('step4.group.members.none_shown')}
              </td></tr>
            )}
            {shown.map(({ candidate, member, suggestion, problem }) => {
              const describedBy = problem ? `group-member-problem-${candidate.key}` : undefined;
              const chosenColumn = member?.target_column ?? '';
              const options = columns.filter((column) => !column.primary_key);

  return (
                <tr key={candidate.key} data-testid={`group-member-row-${candidate.point_id}`} data-problem={problem ?? ''}>
                  <td className={classes('p-2')}>
                    {selectionControls && member ? <button type="button" className={classes('inline-flex min-h-10 max-sm:min-h-11 items-center justify-center rounded-lg border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-2 text-sm text-[var(--gw-text)] hover:bg-[var(--gw-active)] focus-visible:outline-solid focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--gw-focus)] disabled:cursor-not-allowed disabled:opacity-50')} disabled={disabled}
                      onClick={() => selectionControls.onRemove(member)}>{selectionControls.removeLabel(member)}</button> : <input
                      type="checkbox" checked={Boolean(member)} disabled={disabled}
                      aria-label={t('step4.group.members.include_label', { tag: candidate.label })}
                      onChange={(event) => (event.target.checked ? include(candidate, suggestion?.status === 'confirmed' ? suggestion.column : '') : exclude(candidate.key))}
                      data-testid={`group-member-include-${candidate.point_id}`}
                    />}
                  </td>
                  <td className={classes('p-2')}>
                    <div className={classes('font-medium')}>{candidate.label}</div>
                    <div className={classes('text-[var(--gw-muted)]')}>{candidate.device_name} · {candidate.address}{renderCalculation && <> · {candidate.target_type}</>}</div>
                    {member && renderSource?.(member)}
                  </td>
                  {renderCalculation ? <td className={classes('p-2 align-top')} data-testid={`group-member-calculation-${candidate.point_id}`}>{member && renderCalculation(member)}</td> : <td className={classes('p-2 font-mono')}>{candidate.target_type}</td>}
                  <td className={classes('p-2')}>
                    {columnsAssignedOnSave ? <span className="text-sm text-[var(--gw-muted)]">{chosenColumn || editorText('simple.pendingColumn')}</span> : metadataKnown ? (
                      <select
                        value={chosenColumn} disabled={disabled || !member}
                        aria-label={t('step4.group.members.column_label', { tag: candidate.label })}
                        aria-invalid={Boolean(problem)} aria-describedby={describedBy}
                        onChange={(event) => update(candidate, { target_column: event.target.value })}
                        className={classes('rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1')}
                        data-testid={`group-member-column-${candidate.point_id}`}
                      >
                        <option value="">{t('step4.group.members.choose_column')}</option>
                        {chosenColumn && !columnByName.has(chosenColumn) && <option value={chosenColumn}>{chosenColumn}</option>}
                        {options.map((column) => (
                          <option key={column.name} value={column.name}>{column.name} ({column.type})</option>
                        ))}
                      </select>
                    ) : (
                      <input
                        type="text" value={chosenColumn} disabled={disabled || !member}
                        aria-label={t('step4.group.members.column_label', { tag: candidate.label })}
                        aria-invalid={Boolean(problem)} aria-describedby={describedBy}
                        onChange={(event) => update(candidate, { target_column: event.target.value })}
                        className={classes('rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1')}
                        data-testid={`group-member-column-${candidate.point_id}`}
                      />
                    )}
                  </td>
                  <td className={classes('p-2')}>
                    <input
                      type="checkbox" checked={member?.required ?? false} disabled={disabled || !member}
                      aria-label={t('step4.group.members.required_label', { tag: candidate.label })}
                      onChange={(event) => update(candidate, { required: event.target.checked })}
                    />
                  </td>
                  {showEntityKey && (
                    <td className={classes('p-2')}>
                      <input
                        type="text" value={member?.entity_key ?? ''} disabled={disabled || !member}
                        aria-label={t('step4.group.members.entity_label', { tag: candidate.label })}
                        onChange={(event) => update(candidate, { entity_key: event.target.value })}
                        className={classes('w-28 rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1')}
                      />
                    </td>
                  )}
                  <td className={classes('p-2')}>
                    {columnsAssignedOnSave && member ? <span className="text-sm text-[var(--gw-muted)]">{editorText('simple.pendingColumn')}</span> : problem ? (
                      <span id={describedBy} role="alert" className={classes('text-[var(--gw-danger-fg)]')} data-testid={`group-member-problem-${candidate.point_id}`}>
                        {t(`step4.group.members.problem_${problem.replace('-', '_')}`)}
                      </span>
                    ) : suggestion && (!member || member.target_column === suggestion.column || member.target_column === '') && suggestion.status === 'suggested' ? (
                      <span className={classes('text-[var(--gw-accent)]')} data-testid={`group-member-suggested-${candidate.point_id}`}>
                        {t('step4.group.members.suggested', { column: suggestion.column })}
                      </span>
                    ) : member?.target_column ? (
                      <span className={classes('text-[var(--gw-success-fg)]')}>{t('step4.group.members.confirmed')}</span>
                    ) : (
                      <span className={classes('text-[var(--gw-muted)]')}>{t('step4.group.members.unassigned')}</span>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {excluded.length > 0 && (
        <div className={classes('text-sm text-[var(--gw-warning-fg)]')} data-testid="group-member-excluded">
          <p>{t('step4.group.members.excluded', { count: excluded.length })}</p>
          <ul className={classes('mt-1 list-disc pl-5')}>
            {excluded.map((entry) => (
              <li key={`${entry.label}-${entry.reason}`}>{entry.label} — {t(`step4.group.members.excluded_${entry.reason.replace(/-/g, '_')}`)}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};
