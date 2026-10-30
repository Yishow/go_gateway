import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import type { RecordingHomeGuidance, RecordingHomeGuidanceAction } from './recordingHomeGuidance';
import { buildStudioEntityLocatorTarget } from '../../../../services/studioEntityLocator';
import { buildStudioPresentationTarget } from '../../../../services/studioRecordingView';

export interface RecordingHomeLocation { search: string; hash: string }
export function recordingHomeActionTarget(action: RecordingHomeGuidanceAction, location: RecordingHomeLocation): string | null {
  if (action.target === 'recording') return buildStudioEntityLocatorTarget({ step: 4, workspace_id: action.scope.workspaceId,
    device_id: null, rule_id: null, mapping_id: null, group_id: action.groupId, group_revision: action.groupRevision });
  if (action.target === 'configure_points') return buildStudioEntityLocatorTarget({ step: 1, workspace_id: action.scope.workspaceId,
    device_id: action.deviceId, rule_id: null, mapping_id: null, group_id: null, group_revision: null });
  return buildStudioPresentationTarget(location.search, location.hash, 'advanced');
}
/** Initial readonly home. Management search/filter/paging is a later task. */
export function RecordingHome({ guidance, location, onNavigate, onRetry, collection }: {
  guidance: RecordingHomeGuidance; location: RecordingHomeLocation; onNavigate?: (target: string) => void; onRetry: () => void; collection?: ReactNode;
}) {
  const { t } = useTranslation('recording-home');
  const advanced = buildStudioPresentationTarget(location.search, location.hash, 'advanced');
  const navigate = (action: RecordingHomeGuidanceAction) => {
    const target = recordingHomeActionTarget(action, location); if (target) onNavigate?.(target);
  };
  return <main data-workbench-v2="true" data-testid="recording-home" className="min-h-screen p-6">
    <header className="flex flex-wrap items-center justify-between gap-3">
      <h1 className="text-2xl font-semibold">{t('title')}</h1>
      <nav aria-label={t('navigation')} className="flex flex-wrap gap-4">
        {advanced && <button type="button" onClick={() => onNavigate?.(advanced)}>{t('advanced')}</button>}
        {[['/studio/runtime', 'runtime'], ['/studio/logs', 'logs'], ['/test', 'test']].map(([target, key]) =>
          <button key={target} type="button" onClick={() => onNavigate?.(target)}>{t(key)}</button>)}
      </nav>
    </header>
    <section aria-live="polite" className="my-6 space-y-3" data-testid="recording-home-guidance">
      <h2 className="text-lg font-semibold">{t(`reason.${guidance.reason}`)}</h2>
      {guidance.availability.state === 'unavailable' && <p role={guidance.availability.reason === 'error' ? 'alert' : 'status'}>
        {t(guidance.availability.reason === 'loading' ? 'loading' : guidance.availability.reason === 'error' ? 'error' : 'unknown')}
      </p>}
      {guidance.primaryAction && <button type="button" onClick={() => navigate(guidance.primaryAction!)}>{t('openSettings')}</button>}
      <button type="button" onClick={onRetry} className="ml-4">{t('retry')}</button>
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
