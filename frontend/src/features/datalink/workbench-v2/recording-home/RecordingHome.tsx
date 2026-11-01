import './recording-home.css';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import type { RecordingHomeGuidance, RecordingHomeGuidanceAction } from './recordingHomeGuidance';
import { buildStudioEntityLocatorTarget } from '../../../../services/studioEntityLocator';
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
  if (action.target === 'configure_points') return buildStudioEntityLocatorTarget({ step: 1, workspace_id: action.scope.workspaceId,
    device_id: action.deviceId, rule_id: null, mapping_id: null, group_id: null, group_revision: null });
  return buildStudioPresentationTarget(location.search, location.hash, 'advanced');
}
/** Actual readonly home; visual hierarchy never changes evidence or command ownership. */
export function RecordingHome({ guidance, location, onNavigate, onRetry, collection }: {
  guidance: RecordingHomeGuidance; location: RecordingHomeLocation; onNavigate?: (target: string) => void; onRetry: () => void; collection?: ReactNode;
}) {
  const { t } = useTranslation('recording-home');
  const { t: editorText } = useTranslation('recording-editor');
  const advanced = buildStudioPresentationTarget(location.search, location.hash, 'advanced');
  const action = guidance.primaryAction;
  const logs = action?.target === 'recording' ? buildRecordingLogsTarget({ workspaceId: action.scope.workspaceId,
    groupId: action.groupId, groupRevision: action.groupRevision }, location.hash) : null;
  const navigate = (action: RecordingHomeGuidanceAction) => {
    const target = recordingHomeActionTarget(action, location); if (target) onNavigate?.(target);
  };
  return <main data-workbench-v2="true" data-testid="recording-home" className="gw-recording-home">
    <header className="gw-home-header">
      <h1>{t('title')}</h1>
      <button type="button" className="gw-primary" onClick={() => { const target = buildSimpleRecordingTarget({ kind: 'new' }, location.hash); if (target) onNavigate?.(target); }}>{editorText('simple.create')}</button>
      <nav aria-label={t('navigation')}>
        {advanced && <button type="button" onClick={() => onNavigate?.(advanced)}>{t('advanced')}</button>}
        {[['/studio/runtime', 'runtime'], ['/studio/logs', 'logs'], ['/test', 'test']].map(([target, key]) =>
          <button key={target} type="button" onClick={() => onNavigate?.(target)}>{t(key)}</button>)}
      </nav>
    </header>
    <section aria-live="polite" className="gw-home-guidance" data-guidance-state={guidance.availability.state === 'unavailable' && guidance.availability.reason === 'error' ? 'error' : ['operation_recovery', 'delivery_attention', 'source_attention'].includes(guidance.reason) ? 'attention' : guidance.reason === 'evidence_unavailable' ? 'unknown' : 'neutral'} data-testid="recording-home-guidance">
      <h2 className="gw-guidance-heading">{t(`reason.${guidance.reason}`)}</h2>
      {guidance.availability.state === 'unavailable' && <p role={guidance.availability.reason === 'error' ? 'alert' : 'status'}>
        {t(guidance.availability.reason === 'loading' ? 'loading' : guidance.availability.reason === 'error' ? 'error' : 'unknown')}
      </p>}
      <div className="gw-home-actions">{guidance.primaryAction && <button className="gw-primary" type="button" onClick={() => navigate(guidance.primaryAction!)}>{t('openSettings')}</button>}
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
  </main>;
}
