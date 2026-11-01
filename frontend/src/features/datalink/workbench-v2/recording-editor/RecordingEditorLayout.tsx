import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import './recordingEditorLayout.css';

export interface RecordingEditorLayoutProps {
  presentation: 'guided' | 'simple';
  identity: ReactNode;
  sources: ReactNode;
  method: ReactNode;
  destination: ReactNode;
  policy: ReactNode;
  summary?: ReactNode;
  memberCount: number;
  perPoint: boolean;
}
/** Presentation only: same-page navigation never changes history, intent or command owners. */
export function RecordingEditorLayout({ presentation, identity, sources, method, destination, policy, summary, memberCount, perPoint }: RecordingEditorLayoutProps) {
  const { t } = useTranslation('recording-editor');
  if (presentation === 'guided') return <>{method}{identity}{destination}{policy}{sources}</>;
  const locate = (name: string) => {
    const section = document.getElementById(`recording-editor-${name}`);
    section?.scrollIntoView({ block: 'start' });
    section?.focus({ preventScroll: true });
  };
  return <div className="recording-editor-layout" data-testid="recording-editor-layout">
    <nav className="recording-editor-steps" aria-label={t('simple.sections')}>
      {(['sources', 'method', 'destination'] as const).map((name, index) => <button key={name} type="button" onClick={() => locate(name)}>{index + 1} {t(`simple.section.${name}`)}</button>)}
    </nav>
    <div className="recording-editor-columns">
      <div className="recording-editor-main">
        <section className="recording-editor-panel" id="recording-editor-sources" tabIndex={-1} aria-labelledby="recording-editor-sources-title" data-testid="recording-editor-sources">
          {identity}
          <header><h4 id="recording-editor-sources-title">{t('simple.section.sources')}</h4><p>{t(perPoint ? 'simple.sourcesHint' : 'simple.legacySourcesHint')}</p></header>
          <h5 className="recording-editor-point-count">{t('simple.sources', { count: memberCount })}</h5>
          {sources}
          <p className="recording-editor-note">{t('simple.sourceCadenceNote')}</p>
          <button type="button" className="recording-editor-next" onClick={() => locate('method')}>{t('simple.nextMethod')} →</button>
        </section>
        <section className="recording-editor-panel" id="recording-editor-method" tabIndex={-1} aria-labelledby="recording-editor-method-title" data-testid="recording-editor-method-note">
          <header><h4 id="recording-editor-method-title">{t('simple.section.method')}</h4><p>{t('simple.methodHint')}</p></header>
          {method}
          <button type="button" className="recording-editor-next" onClick={() => locate('destination')}>{t('simple.nextDestination')} →</button>
        </section>
        <section className="recording-editor-panel" id="recording-editor-destination" tabIndex={-1} aria-labelledby="recording-editor-destination-title" data-testid="recording-editor-destination">
          <header><h4 id="recording-editor-destination-title">{t('simple.section.destination')}</h4><p>{t('simple.destinationHint')}</p></header>
          {destination}{policy}
        </section>
      </div>
      {summary && <aside className="recording-editor-summary" aria-label={t('simple.summary.title')}>{summary}</aside>}
    </div>
  </div>;
}
