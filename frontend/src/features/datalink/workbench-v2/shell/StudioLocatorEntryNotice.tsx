import { useTranslation } from 'react-i18next';
import type { StudioLocatorEntryRead } from '../../../../types/studioV2LocatorEntry';
import { buildStudioEntityLocatorTarget } from '../../../../services/studioEntityLocator';

export type StudioLocatorNoticeReason = 'loading' | 'invalid' | 'unavailable' | 'stale' | 'foreign' | 'removed' | 'draftBlocked' | 'unrequested';
/** Only current same-workspace proof can offer entity links; no unscoped bootstrap fallback. */
export function StudioLocatorEntryNotice({ reason, entry, navigateTo }: {
  reason: StudioLocatorNoticeReason; entry?: StudioLocatorEntryRead; navigateTo?: (target: string) => void;
}) {
  const { t } = useTranslation('runtime-dashboard');
  const owned = entry && entry.workspace.id === entry.locator.workspace_id ? entry.workspace : undefined;
  return <section role={reason === 'loading' ? 'status' : 'alert'} data-testid="studio-locator-notice" className="m-4 space-y-3 rounded border border-amber-600 p-4">
    <p>{t(`locatorEntry.${reason}`)}</p>
    {owned && <nav aria-label={t('locatorEntry.listTitle')}><p>{t('locatorEntry.listTitle')}</p>
      {owned.ordered_device_ids.map(id => {
        const target = buildStudioEntityLocatorTarget({ step: 1, workspace_id: owned.id, device_id: id,
          rule_id: null, mapping_id: null, group_id: null, group_revision: null });
        return target ? <a key={id} href={target} className="mr-3 underline" onClick={event => {
          if (navigateTo && !event.ctrlKey && !event.metaKey && !event.shiftKey && !event.altKey && event.button === 0) { event.preventDefault(); navigateTo(target); }
        }}>{t('locatorEntry.listDevice', { id })}</a> : null;
      })}
    </nav>}
  </section>;
}
