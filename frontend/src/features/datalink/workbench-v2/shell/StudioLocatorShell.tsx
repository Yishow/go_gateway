import type { ReactNode } from 'react';
import type { TFunction } from 'i18next';
import { useTranslation } from 'react-i18next';
import type { StudioV2WorkspaceReadinessIssue } from '../../../../types/studioV2WorkspaceReadiness';

/** A precise editor never advertises unloaded workflow steps. Leaving is explicit. */
export function StudioLocatorShell({ canLeave, navigateTo, children, unavailable, onRetry, recordContext }: {
  recordContext?: { name: string; target: string };
  canLeave: boolean; navigateTo?: (target: string) => void; children: ReactNode; unavailable?: boolean; onRetry?: () => void;
}) {
  const { t } = useTranslation('recording-home');
  return <section data-testid="studio-local-editor" className="mx-auto w-full max-w-screen-xl space-y-4 p-4">
    <header className="flex flex-wrap items-start justify-between gap-3">
      <div><h1 className="text-xl font-semibold">{recordContext ? <>{recordContext.name} · {t('locator.recordTitle')}</> : t('locator.localTitle')}</h1><p className="text-sm text-[var(--gw-muted)]">{t('locator.localDescription')}</p></div>
      <div className="flex flex-wrap gap-2">
        <button type="button" className="gw-secondary" disabled={!canLeave || !navigateTo} onClick={() => { if (canLeave) navigateTo?.(recordContext?.target ?? '/studio/v2?view=recordings'); }}>{t(recordContext ? 'locator.backRecord' : 'locator.backRecordings')}</button>
        <button type="button" className="gw-secondary" disabled={!canLeave || !navigateTo} onClick={() => { if (canLeave) navigateTo?.('/studio/v2?view=advanced'); }}>{t('locator.openAdvanced')}</button>
      </div>
    </header>
    {!canLeave && <p role="status" className="text-sm text-[var(--gw-muted)]">{t('locator.finishBeforeLeaving')}</p>}
    {unavailable && <div role="alert"><p>{t('locator.readUnavailable')}</p><button type="button" className="gw-secondary" onClick={onRetry}>{t('locator.retryRead')}</button></div>}
    <fieldset disabled={unavailable} inert={unavailable} className="min-w-0 border-0 p-0">{children}</fieldset>
  </section>;
}

export function runtimeReturnGuidanceFor(
  issue: StudioV2WorkspaceReadinessIssue,
  t: TFunction<'workbench-v2'>,
): string | null {
  switch (issue.code) {
    case 'point-missing':
      return t('runtime_return.point_missing_action');
    default:
      return null;
  }
}
