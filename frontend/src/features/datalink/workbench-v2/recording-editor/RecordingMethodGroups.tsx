import { useTranslation } from 'react-i18next';
import { ParameterHelp } from '../components/ParameterHelp';
import type { EditorDraft } from '../state/writeGroup/draft';
import { recordingMethodGroups } from './recordingMethodGrouping';

/** Moves focus to a member's calculation control in the source table. When the row is
 * filtered out or paged away, the source section itself is shown instead. */
function focusMember(pointId: string) {
  const control = document.querySelector<HTMLElement>(`[data-testid="group-member-calculation-${CSS.escape(pointId)}"] select`);
  const target = control ?? document.getElementById('recording-editor-sources');
  target?.scrollIntoView({ block: 'center' });
  target?.focus({ preventScroll: true });
}

/** The per-point recording's methods, grouped as the summary rail shows them; edits stay in the source table. */
export function RecordingMethodGroups({ draft }: { draft: Pick<EditorDraft, 'members'> }) {
  const { t } = useTranslation('recording-editor');
  const groups = recordingMethodGroups(draft, t);
  return <div className="recording-method-groups" data-testid="recording-method-groups">
    <div className="recording-method-groups__head"><span>{t('simple.summary.results')}</span><ParameterHelp compact help="record.method" /></div>
    {groups.length ? <ul>{groups.map(group => <li key={group.key}>
      <span className="recording-method-groups__value">{group.method} · {group.interval} {group.unit}</span>
      <span className="recording-method-groups__count">{t('simple.summary.points', { count: group.count })}</span>
      <button type="button" className="recording-method-groups__go" onClick={() => focusMember(group.firstPointId)}>{t('simple.goToMember')}</button>
    </li>)}</ul> : <p>{t('simple.methodEmpty')}</p>}
  </div>;
}
