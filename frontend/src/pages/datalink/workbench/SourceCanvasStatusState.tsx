import { useTranslation } from 'react-i18next';
import { Spinner } from '../../../components/ui/spinner';

type SourceCanvasStatusStateProps = {
  state: 'loading' | 'error';
  error?: unknown;
  onRetry?: () => void;
};

export function SourceCanvasStatusState({
  error,
  onRetry,
  state,
}: SourceCanvasStatusStateProps) {
  const { t } = useTranslation();
  if (state === 'loading')
    return (
      <section className="min-h-0 flex-1 rounded-2xl border border-slate-800 bg-slate-950/40 p-6">
        <div
          className="flex items-center justify-center gap-3 rounded-2xl border border-slate-800 bg-slate-950/40 p-8"
          data-testid="source-loading-state"
        >
          <Spinner size="lg" />
          <span className="text-sm text-slate-400">{t('workbench.source.loading')}</span>
        </div>
      </section>
    );

  return (
    <section className="min-h-0 flex-1 rounded-2xl border border-rose-500/30 bg-rose-500/10 p-6">
      <div className="space-y-2" data-testid="source-error-state">
        <p className="font-medium text-rose-100">{t('workbench.source.loadFailed')}</p>
        <p className="text-sm text-rose-100/80">{error instanceof Error ? error.message : ''}</p>
        {onRetry ? (
          <button
            className="inline-flex items-center rounded-xl border border-rose-400/40 bg-rose-500/10 px-3 py-2 text-sm font-medium text-rose-100 transition hover:border-rose-300/50 hover:bg-rose-500/20"
            onClick={onRetry}
            type="button"
          >
            {t('workbench.source.retry')}
          </button>
        ) : null}
      </div>
    </section>
  );
}
