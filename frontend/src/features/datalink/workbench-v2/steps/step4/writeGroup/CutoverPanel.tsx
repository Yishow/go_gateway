import { ParameterHelp } from '../../../components/ParameterHelp';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useCutoversQuery } from '../../../../../../hooks/datalink/useStudioV2WriteGroupCutover';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';
import { groupProfileKind } from '../../../state/circuitProfileBasic';
import { instantFromLocal, isCutoverSource, isCutoverTarget, loadCutoverIntent, rollbackCandidates, type CutoverScope } from '../../../state/cutoverFlow';
import type { GroupCandidate } from '../../../state/writeGroup/candidates';
import { CutoverRun } from './CutoverRun';
import { CutoverHistory } from './CutoverOutcomeView';
import { useCutoverFlow } from './useCutoverFlow';

export interface CutoverPanelProps {
  /** The saved, not deleted groups of the workspace. */
  groups: WriteGroup[];
  workspaceId: string;
  workspaceRevision: string;
  connectorRevision: string;
  readonly: boolean;
  /** Opens one exact group in the existing editor, where its draft, schema and blockers are repaired. */
  onOpenGroup: (groupId: string) => void;
  candidates?: GroupCandidate[];
}

type Mode = 'forward' | 'rollback';

/**
 * Moves sources of a running wide or custom group to prepared long targets at a reviewed future boundary, or rolls such a
 * switch back at a new boundary. It selects from the canonical groups only and holds no configuration of its own.
 */
const CutoverPanelContent: React.FC<CutoverPanelProps> = ({ groups, workspaceId, workspaceRevision, connectorRevision, readonly, onOpenGroup, candidates: propCandidates }) => {
  const { t } = useTranslation('workbench-v2');
  const storedIntent = React.useMemo(() => loadCutoverIntent(workspaceId), [workspaceId]);
  const [mode, setMode] = React.useState<Mode>(() => storedIntent?.kind ?? 'forward');
  const [sourceId, setSourceId] = React.useState(() => storedIntent?.source_id ?? '');
  const [targetIds, setTargetIds] = React.useState<string[]>([]);
  const [boundaryLocal, setBoundaryLocal] = React.useState('');
  const [cutoverId, setCutoverId] = React.useState('');
  const sources = React.useMemo(
    () => groups.filter((group) => mode === 'forward' ? isCutoverSource(group) : groupProfileKind(group) === undefined),
    [groups, mode],
  );
  const source = sources.find((group) => group.id === sourceId) ?? sources[0];
  const targetChoices = React.useMemo(() => groups.filter(isCutoverTarget), [groups]);
  const cutovers = useCutoversQuery(source?.id, true);
  const candidates = rollbackCandidates(cutovers.data ?? []);
  const names = React.useMemo(() => new Map(groups.map((group) => [group.id, group.name])), [groups]);
  const groupName = (id: string) => names.get(id) ?? id;

  const selected = mode === 'forward' ? targetChoices.filter((group) => targetIds.includes(group.id)) : [];
  const chosen = candidates.find((item) => item.id === cutoverId);
  const ready = source && (mode === 'forward' ? selected.length > 0 : Boolean(chosen));
  const scope: CutoverScope | undefined = ready ? {
    workspaceId, workspaceRevision, connectorRevision, source, targets: selected, kind: mode,
    boundary: instantFromLocal(boundaryLocal), ...(mode === 'rollback' && chosen ? { cutoverId: chosen.id } : {}),
  } : undefined;

  const flow = useCutoverFlow(scope, readonly, workspaceId);
  const isLocked = readonly || flow.busy || flow.outcome.state === 'unknown';

  return (
    <section data-testid="cutover-panel" aria-labelledby="cutover-title" className="space-y-3 rounded-xl border border-[var(--gw-control)] bg-[var(--gw-surface)] p-4 text-xs text-[var(--gw-text)]">
      <div>
        <h4 id="cutover-title" className="text-sm font-semibold text-[var(--gw-text)]">{t('step4.cutover.title')}</h4>
        <p className="mt-1 text-[var(--gw-muted)]">{t('step4.cutover.subtitle')}</p>
      </div>
      <div role="group" aria-label={t('step4.cutover.title')} className="flex gap-2">
        <ParameterHelp help="cutover.mode" />
        {(['forward', 'rollback'] as const).map((item) => (
          <button key={item} type="button" data-testid={`cutover-mode-${item}`} aria-pressed={mode === item} disabled={isLocked} onClick={() => { setMode(item); setSourceId(''); flow.resetReview(); }} className={`rounded border px-3 py-1.5 ${mode === item ? 'border-[var(--gw-focus)] text-[var(--gw-accent)]' : 'border-[var(--gw-control)] text-[var(--gw-text)]'} disabled:opacity-40`}>
            {t(`step4.cutover.mode_${item}`)}
          </button>
        ))}
      </div>
      {sources.length === 0 ? (
        <p role="status" data-testid="cutover-no-source" className="text-[var(--gw-muted)]">{t('step4.cutover.no_source')}</p>
      ) : (
        <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block">
          <span>{t('step4.cutover.source')}</span>
          <select data-testid="cutover-source" value={source?.id ?? ''} disabled={isLocked} onChange={(event) => { setSourceId(event.target.value); flow.resetReview(); }} className="mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1.5 text-[var(--gw-text)]">
            {sources.map((group) => <option key={group.id} value={group.id}>{group.name}</option>)}
          </select>
        </label><ParameterHelp help="cutover.source" /></div>
      )}
      {source && mode === 'forward' && (
        <fieldset className="space-y-1">
          <legend className="font-medium text-[var(--gw-text)]">{t('step4.cutover.targets')}</legend>
          {targetChoices.length === 0 && <p data-testid="cutover-no-targets" className="text-[var(--gw-muted)]">{t('step4.cutover.no_targets')}</p>}
          {targetChoices.map((group) => (
            <div key={group.id} className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 flex items-center gap-2">
              <input type="checkbox" data-testid={`cutover-target-${group.id}`} checked={targetIds.includes(group.id)} disabled={isLocked}
                onChange={(event) => setTargetIds((current) => event.target.checked ? [...current, group.id] : current.filter((id) => id !== group.id))} />
              <span>{group.name} · {group.destination.table_name} · {group.members.length}</span>
            </label><ParameterHelp help="cutover.target" /></div>
          ))}
        </fieldset>
      )}
      {source && mode === 'rollback' && (
        <fieldset className="space-y-1">
          <legend className="font-medium text-[var(--gw-text)]">{t('step4.cutover.rollback_candidates')}</legend>
          {candidates.length === 0 && <p data-testid="cutover-no-rollback" className="text-[var(--gw-muted)]">{t('step4.cutover.no_rollback')}</p>}
          {candidates.map((item) => (
            <div key={item.id} className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 flex items-center gap-2">
              <input type="radio" name="cutover-rollback" data-testid={`cutover-rollback-${item.id}`} checked={cutoverId === item.id} disabled={isLocked} onChange={() => setCutoverId(item.id)} />
              <span><code>{item.boundary}</code> · {item.started.map((ref) => groupName(ref.group_id)).join(', ')}</span>
            </label><ParameterHelp help="cutover.rollback" /></div>
          ))}
          <p className="text-[var(--gw-muted)]">{t('step4.cutover.rollback_note')}</p>
        </fieldset>
      )}
      {source && (
        <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block">
          <span>{t(mode === 'forward' ? 'step4.cutover.boundary' : 'step4.cutover.rollback_boundary')}</span>
          <input type="datetime-local" step={1} data-testid="cutover-boundary" value={boundaryLocal} disabled={isLocked} onChange={(event) => setBoundaryLocal(event.target.value)} className="mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1.5 text-[var(--gw-text)]" />
          <span className="mt-1 block text-[var(--gw-muted)]">{t('step4.cutover.boundary_hint')}</span>
        </label><ParameterHelp help="cutover.boundary" /></div>
      )}
      {(source || flow.outcome.state === 'unknown') && <CutoverRun flow={flow} scope={scope} readonly={readonly} groupName={groupName} onOpenGroup={onOpenGroup} candidates={propCandidates} />}
      <CutoverHistory cutovers={cutovers.data ?? []} groupName={groupName} />
    </section>
  );
};

/** Each workspace owns its review and pending request; returning restores only that workspace's session intent. */
export const CutoverPanel: React.FC<CutoverPanelProps> = (props) => <CutoverPanelContent key={props.workspaceId} {...props} />;
