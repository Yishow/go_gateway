import { RecordingProductShell } from '../shell/RecordingProductShell';
import './recording-home.css';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import type { RecordingHomeGuidance, RecordingHomeGuidanceAction } from './recordingHomeGuidance';
import { buildRecordingRepairTarget, buildSimpleRecordingTarget, buildStudioPresentationTarget } from '../../../../services/studioRecordingView';
import { buildRecordingLogsTarget } from '../../../../services/runtimeLogsNavigation';

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
export function RecordingHome({ guidance, location, onNavigate, onRetry, collection, onConfigurePoints, workspaceId, mode = 'recordings', onCreateDevice, settling = false }: {
  /** First-load reads are still settling: hold one stable reading state instead of repainting each read batch. */
  settling?: boolean; mode?: 'recordings' | 'devices'; onCreateDevice?: () => void; workspaceId?: string; guidance: RecordingHomeGuidance; location: RecordingHomeLocation; onNavigate?: (target: string) => void; onRetry: () => void; collection?: ReactNode; onConfigurePoints?: (deviceId: string) => void;
}) {
  const { t } = useTranslation('recording-home');
  const { t: editorText } = useTranslation('recording-editor');
  const action = guidance.primaryAction;
  const logs = action?.target === 'recording' ? buildRecordingLogsTarget({ workspaceId: action.scope.workspaceId,
    groupId: action.groupId, groupRevision: action.groupRevision }, location.hash) : null;
  const navigate = (action: RecordingHomeGuidanceAction) => {
    if (action.target === 'configure_points') { onConfigurePoints?.(action.deviceId); return; }
    const target = recordingHomeActionTarget(action, location); if (target) onNavigate?.(target);
  };
  return <><RecordingProductShell workspaceId={workspaceId ?? guidance.primaryAction?.scope.workspaceId}
    tools={<RecordingHomeTools location={location} onNavigate={onNavigate} />} /><main data-workbench-v2="true" data-testid="recording-home" className="gw-recording-home" aria-busy={settling || undefined} data-home-settling={settling ? 'true' : 'false'}>
    <header className="gw-home-header">
      <h1>{t(mode === 'devices' ? 'devicesTitle' : 'title')}</h1>
      <button type="button" className={mode === 'devices' ? 'gw-home-create gw-primary' : 'gw-home-create'} onClick={() => { if (mode === 'devices') { onCreateDevice?.(); return; } const target = buildSimpleRecordingTarget({ kind: 'new' }, location.hash); if (target) onNavigate?.(target); }}>{mode === 'devices' ? t('createDevice') : editorText('simple.create')}</button>

    </header>
    <section aria-live="polite" className="gw-home-guidance" data-guidance-state={guidance.availability.state === 'unavailable' && guidance.availability.reason === 'error' ? 'error' : ['operation_recovery', 'delivery_attention', 'source_attention'].includes(guidance.reason) ? 'attention' : guidance.reason === 'evidence_unavailable' ? 'unknown' : 'neutral'} data-testid="recording-home-guidance">
      {settling && <p className="gw-guidance-heading gw-home-settling-status" role="status">{t('loading')}</p>}
      <p className="gw-guidance-heading">{t(`reason.${guidance.reason}`)}</p>
      {guidance.availability.state === 'unavailable' && <p role={guidance.availability.reason === 'error' ? 'alert' : 'status'}>
        {t(guidance.availability.reason === 'loading' ? 'loading' : guidance.availability.reason === 'error' ? 'error' : 'unknown')}
      </p>}
      <div className="gw-home-actions">{mode === 'recordings' && guidance.primaryAction && <button className="gw-primary" type="button" onClick={() => navigate(guidance.primaryAction!)}>{t('openSettings')}</button>}
      {logs && <button type="button" data-testid="recording-scoped-logs" onClick={() => onNavigate?.(logs)}>{t('logs')}</button>}
      <button type="button" onClick={onRetry}>{t('retry')}</button></div>
    </section>
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
