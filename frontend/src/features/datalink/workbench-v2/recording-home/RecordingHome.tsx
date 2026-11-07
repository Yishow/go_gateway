import { RecordingProductShell } from '../shell/RecordingProductShell';
import './recording-home.css';
import './recordingHomeStatus.css';
import { useId, type ReactNode, type RefObject } from 'react';
import { useTranslation } from 'react-i18next';
import type { RecordingHomeGuidance, RecordingHomeGuidanceAction } from './recordingHomeGuidance';
import { buildRecordingRepairTarget, buildSimpleRecordingTarget, buildStudioPresentationTarget } from '../../../../services/studioRecordingView';
import { buildRecordingLogsTarget } from '../../../../services/runtimeLogsNavigation';
import { ArrowRight, CircleHelp, CircleX, FilePenLine, Info, Plus, RefreshCw, ScrollText, TriangleAlert } from 'lucide-react';
import { RecordingHomeUnconfirmed } from './RecordingHomeUnconfirmed';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';

export interface RecordingHomeLocation { search: string; hash: string }
export function recordingHomeActionTarget(action: RecordingHomeGuidanceAction, location: RecordingHomeLocation): string | null {
  if (action.target === 'recording') {
    const intent = action.intent;
    if (intent && (intent.kind !== 'query_operation' && intent.kind !== 'recover_operation' || !intent.operationOwner ||
      intent.scope.workspaceId !== action.scope.workspaceId || intent.scope.workspaceRevision !== action.scope.workspaceRevision ||
      !('groupId' in intent.scope) || intent.scope.groupId !== action.groupId || intent.scope.groupRevision !== action.groupRevision)) return null;
    return buildRecordingRepairTarget({ workspaceId: action.scope.workspaceId, workspaceRevision: action.scope.workspaceRevision,
      groupId: action.groupId, groupRevision: action.groupRevision, section: action.section,
      ...(intent && intent.operationOwner && (intent.kind === 'query_operation' || intent.kind === 'recover_operation')
        ? { operation: { kind: intent.kind, owner: intent.operationOwner, deviceId: intent.scope.deviceId,
          ...(intent.scope.circuitId ? { circuitId: intent.scope.circuitId } : {}), ...(intent.operationId ? { operationId: intent.operationId } : {}) } } : {}) }, location.hash);
  }
  if (action.target === 'connect_device' || action.target === 'choose_sources') return buildSimpleRecordingTarget({ kind: 'new' }, location.hash);
  if (action.target === 'configure_points') return null; // Same-home device modal, not an advanced locator route.
  return buildStudioPresentationTarget(location.search, location.hash, 'advanced');
}
/** Product-header navigation shared by the home and its first-load frame, so the frame never jumps. */
export function RecordingHomeTools({ location, onNavigate }: { location: RecordingHomeLocation; onNavigate?: (target: string) => void }) {
  const { t } = useTranslation('recording-home');
  const advanced = buildStudioPresentationTarget(location.search, location.hash, 'advanced');
  return <>
    {advanced && <button type="button" onClick={() => onNavigate?.(advanced)}>{t('advanced')}</button>}
    {[['/studio/runtime', 'runtime'], ['/studio/logs', 'logs'], ['/test', 'test']].map(([target, key]) =>
      <button key={target} type="button" onClick={() => onNavigate?.(target)}>{t(key)}</button>)}
  </>;
}
/** Actual readonly home; visual hierarchy never changes evidence or command ownership. */
export function RecordingHome({ guidance, location, onNavigate, onRetry, collection, statusContent, updatedAt, onConfigurePoints, workspaceId, onDelete, settling = false, homeRef }: {
  /** First-load reads are still settling: hold one stable reading state instead of repainting each read batch. */
  settling?: boolean; homeRef?: RefObject<HTMLElement | null>; onDelete?: (group: RecordingReadonly<WriteGroup>) => void; workspaceId?: string; guidance: RecordingHomeGuidance; location: RecordingHomeLocation; onNavigate?: (target: string) => void; onRetry: () => void; collection?: ReactNode; statusContent?: ReactNode; updatedAt?: number; onConfigurePoints?: (deviceId: string) => void;
}) {
  const { t } = useTranslation('recording-home');
  const { t: editorText } = useTranslation('recording-editor');
  const action = guidance.primaryAction;
  const logs = action?.target === 'recording' ? buildRecordingLogsTarget({ workspaceId: action.scope.workspaceId,
    groupId: action.groupId, groupRevision: action.groupRevision }, location.hash) : null;
  const navigate = (action: RecordingHomeGuidanceAction) => {
    if (action.target === 'configure_points') { onConfigurePoints?.(action.deviceId); return; }
    // A removed item without a known operation opens its original management entry, never a deleted editor.
    const removed = action.target === 'recording' && !action.intent &&
      guidance.recordingItems.some(item => item.id === action.groupId && item.evidence.group.status === 'deleted');
    const target = removed ? buildStudioPresentationTarget('?step=output&target=database', location.hash, 'advanced')
      : recordingHomeActionTarget(action, location);
    if (target) onNavigate?.(target);
  };
  const guidanceState = guidance.availability.state === 'unavailable' && guidance.availability.reason === 'error' ? 'error'
    : ['operation_recovery', 'delivery_attention', 'source_attention'].includes(guidance.reason) ? 'attention' : guidance.reason === 'evidence_unavailable' ? 'unknown' : 'neutral';
  const GuidanceIcon = { error: CircleX, attention: TriangleAlert, unknown: CircleHelp, neutral: Info }[guidanceState];
  // Item-level gaps are named below; only an unreadable home says the whole state is unread.
  const heading = guidance.reason === 'evidence_unavailable' && guidance.availability.state === 'available' ? 'reason.evidence_unconfirmed' : `reason.${guidance.reason}`;
  const lastUpdated = updatedAt ? new Date(updatedAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false }) : '—';
  const draftId = useId();
  const drafts = guidance.draftItems;
  const cta = guidance.primaryAction && guidance.reason !== 'draft' ? guidance.primaryAction : null;
  const hint = guidance.availability.state === 'available' ? `deviceStatus.hints.${heading.slice('reason.'.length)}` : null;
  return <><RecordingProductShell workspaceId={workspaceId ?? guidance.primaryAction?.scope.workspaceId}
    tools={<RecordingHomeTools location={location} onNavigate={onNavigate} />} /><main ref={homeRef} data-workbench-v2="true" data-testid="recording-home" className="gw-recording-home" aria-busy={settling || undefined} data-home-settling={settling ? 'true' : 'false'}>
    <header className="gw-home-header">
      <div className="gw-home-title">
        <h1>{t('title')}</h1>
        <p className="gw-home-meta"><span className="gw-home-freshness">{t('deviceStatus.lastUpdated', { time: lastUpdated })}</span>
          <span className="gw-home-help"><button type="button" aria-label={t('deviceStatus.freshnessHelp')}><CircleHelp aria-hidden="true" size={15} /></button>
            <span role="tooltip">{t('sources.fresh')} {t('observationBoundary')}</span></span></p>
      </div>
      <div className="gw-home-toolbar">
        <div className="gw-home-toolgroup" role="group" aria-label={t('deviceStatus.tools')}>
          <button type="button" data-testid={logs ? 'recording-scoped-logs' : undefined} onClick={() => onNavigate?.(logs ?? '/studio/logs')}><ScrollText aria-hidden="true" size={16} /><span>{t('logs')}</span></button>
          <button type="button" onClick={onRetry}><RefreshCw aria-hidden="true" size={16} /><span>{t('retry')}</span></button>
        </div>
        <button type="button" className="gw-home-create gw-primary" onClick={() => { const target = buildSimpleRecordingTarget({ kind: 'new' }, location.hash); if (target) onNavigate?.(target); }}><Plus aria-hidden="true" size={16} /><span>{editorText('simple.create')}</span></button>
      </div>
    </header>
    {statusContent}
    <div className="gw-home-worklist">
    <section aria-live="polite" className="gw-home-guidance" data-guidance-state={guidanceState} data-testid="recording-home-guidance">
      {settling && <p className="gw-guidance-heading gw-home-settling-status" role="status">{t('loading')}</p>}
      <span className="gw-guidance-icon" aria-hidden="true"><GuidanceIcon size={20} /></span>
      <div className="gw-guidance-body">
        <p className="gw-guidance-eyebrow">{t('deviceStatus.nextStep')}</p>
        <p className="gw-guidance-heading">{t(heading)}</p>
        {hint && <p className="gw-guidance-hint">{t(hint, { count: drafts.length })}</p>}
        {guidance.availability.state === 'unavailable' && <p role={guidance.availability.reason === 'error' ? 'alert' : 'status'}>
          {t(guidance.availability.reason === 'loading' ? 'loading' : guidance.availability.reason === 'error' ? 'error' : 'unknown')}
        </p>}
        {drafts.length > 0 && <p className="gw-draft-label">{t('deviceStatus.drafts')}</p>}
        {drafts.length > 0 && <ul className="gw-draft-list" data-testid="home-draft-banner" aria-label={t('deviceStatus.drafts')}>
          {drafts.slice(0, 3).map((item, index) => <li key={item.id}>
            <span className="gw-draft-icon" aria-hidden="true"><FilePenLine size={16} /></span>
            <span className="gw-draft-name"><strong id={`${draftId}-${index}`}>{item.name}</strong><small>{t('deviceStatus.draftState')}</small></span>
            <button type="button" aria-describedby={`${draftId}-${index}`} onClick={() => navigate(item.action)}><span>{t('deviceStatus.resume')}</span><ArrowRight aria-hidden="true" size={15} /></button>
          </li>)}
          {drafts.length > 3 && <li className="gw-draft-more">{t('deviceStatus.moreDrafts', { count: drafts.length - 3 })}</li>}
        </ul>}
      </div>
      {cta && <button type="button" className="gw-guidance-cta" onClick={() => navigate(cta)}><span>{t('openSettings')}</span><ArrowRight aria-hidden="true" size={16} /></button>}
    </section>
    {/* Unresolved/unknown evidence stays traceable. */}
    <RecordingHomeUnconfirmed guidance={guidance} onDelete={onDelete} onInspect={item => navigate(item.action)}
      onLogs={item => { const target = buildRecordingLogsTarget({ workspaceId: item.action.scope.workspaceId, groupId: item.id,
        groupRevision: item.evidence.group.revision }, location.hash); if (target) onNavigate?.(target); }} />
    </div>
    {collection ?? <section aria-label={t('saved')}>
      <h2 className="text-lg font-semibold">{t('saved')}</h2>
      <ul>{guidance.recordingItems.map(item => <li key={item.id} className="my-2">
        <button type="button" onClick={() => navigate(item.action)}>{item.name}</button>
        <span className="ml-3">{guidance.attentionItems.find(attention => attention.id === item.id)?.reason === 'evidence_unavailable' ? t('statusUnknown') :
          guidance.attentionItems.some(attention => attention.id === item.id) ? t('attention') :
          guidance.draftItems.some(draft => draft.id === item.id) ? t('draft') : t('savedState')}</span>
      </li>)}</ul>
    </section>}
  </main></>;
}
