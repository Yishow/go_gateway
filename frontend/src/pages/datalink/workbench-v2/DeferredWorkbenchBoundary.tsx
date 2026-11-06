import { Suspense, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { ErrorBoundary } from '../../../components/ErrorBoundary';

function DeferredWorkbenchState({ failed = false, onClose }: { failed?: boolean; onClose?: () => void }) {
  const { t } = useTranslation('recording-home');
  return <section className={onClose ? 'gw-modal-body gw-home-source-loading' : 'wbv2-bootstrap-state bg-canvas min-h-screen'}
    role={failed ? 'alert' : 'status'} data-testid={failed ? 'advanced-load-error' : 'advanced-load-pending'}>
    {!onClose && <h1>{t('advanced')}</h1>}
    <p>{t(failed ? 'error' : 'loading')}</p>
    {failed && <button type="button" onClick={() => window.location.reload()}>{t('retry')}</button>}
    {onClose && <button type="button" onClick={onClose}>{t('source.close')}</button>}
  </section>;
}

/** Keep this boundary inside the source host, never around its reading frame or draft guard.
 * Reload is an explicit action: resetting React.lazy would rethrow its cached rejection.
 */
export function DeferredWorkbenchBoundary({ children, onClose }: { children: ReactNode; onClose?: () => void }) {
  return <ErrorBoundary fallback={() => <DeferredWorkbenchState failed onClose={onClose} />}>
    <Suspense fallback={<DeferredWorkbenchState onClose={onClose} />}>{children}</Suspense>
  </ErrorBoundary>;
}
