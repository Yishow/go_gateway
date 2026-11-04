import { useLayoutEffect, useRef, useState, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import './recordingEditorLayout.css';
const sectionNames = ['sources', 'method', 'destination'] as const;
type SectionName = typeof sectionNames[number];

export interface RecordingEditorLayoutProps {
  presentation: 'guided' | 'simple';
  identity: ReactNode;
  sources: ReactNode;
  method: ReactNode;
  destination: ReactNode;
  policy: ReactNode;
  summary?: ReactNode;
  confirmation?: ReactNode;
  fieldsDisabled?: boolean;
  memberCount: number;
  selectedScope?: { shown: number; total: number };
  selectedScopeInHeading?: boolean;
  /** Kept for callers; the per-point and legacy hints were removed from the layout. */
  perPoint?: boolean;
}
/** Presentation only: same-page navigation never changes history, intent or command owners. */
export function RecordingEditorLayout({ presentation, identity, sources, method, destination, policy, summary, confirmation, fieldsDisabled = false, memberCount, selectedScope, selectedScopeInHeading = false }: RecordingEditorLayoutProps) {
  const { t } = useTranslation('recording-editor');
  const nav = useRef<HTMLElement>(null);
  const [activeSection, setActiveSection] = useState<SectionName>('sources');
  useLayoutEffect(() => {
    const element = nav.current, root = element?.closest<HTMLElement>('.gw-simple-recording-main');
    if (!root || !element) return;
    const sections = sectionNames.map(name => ({ name, element: document.getElementById(`recording-editor-${name}`) }));
    const sync = () => {
      const threshold = element.getBoundingClientRect().bottom + 12;
      let active: SectionName = 'sources';
      for (const section of sections) if (section.element && section.element.getBoundingClientRect().top <= threshold) active = section.name;
      if (window.scrollY > 0 && window.scrollY + window.innerHeight >= document.documentElement.scrollHeight - 2 &&
        sections[2].element && sections[2].element.getBoundingClientRect().top < window.innerHeight) active = 'destination';
      setActiveSection(active);
    };
    const update = () => { root.style.setProperty('--recording-nav-height', `${Math.ceil(element.getBoundingClientRect().height)}px`); sync(); };
    update(); const observer = typeof ResizeObserver === 'undefined' ? undefined : new ResizeObserver(update);
    observer?.observe(element); sections.forEach(section => { if (section.element) observer?.observe(section.element); });
    window.addEventListener('resize', update); window.addEventListener('scroll', sync, { passive: true });
    return () => { observer?.disconnect(); window.removeEventListener('resize', update); window.removeEventListener('scroll', sync); root.style.removeProperty('--recording-nav-height'); };
  }, [presentation]);
  if (presentation === 'guided') return <><fieldset disabled={fieldsDisabled} className="min-w-0 space-y-5" data-testid="group-save-fields">{method}{identity}{destination}{policy}{sources}</fieldset>{confirmation}</>;
  const locate = (name: SectionName) => {
    const section = document.getElementById(`recording-editor-${name}`);
    section?.scrollIntoView({ block: 'start' });
    section?.focus({ preventScroll: true });
  };
  return <div className="recording-editor-layout" data-testid="recording-editor-layout">
    <nav ref={nav} className="recording-editor-steps" aria-label={t('simple.sections')}>
      {sectionNames.map((name, index) => <button key={name} type="button" aria-current={activeSection === name ? 'step' : undefined} onClick={() => locate(name)}>{index + 1} {t(`simple.section.${name}`)}</button>)}
    </nav>
    <div className="recording-editor-columns">
      <div className="recording-editor-main">
        <section className="recording-editor-panel" id="recording-editor-sources" tabIndex={-1} aria-labelledby="recording-editor-sources-title" data-testid="recording-editor-sources">
          <fieldset disabled={fieldsDisabled} className="min-w-0 space-y-3" data-testid="group-save-fields">
          {identity}
          <header><h4 id="recording-editor-sources-title">{t('simple.section.sources')}</h4></header>
          <div className="recording-editor-source-selection"><h5 className="recording-editor-point-count" aria-live={selectedScopeInHeading ? 'polite' : undefined} data-testid={selectedScopeInHeading ? 'group-member-scope' : undefined}>{t('simple.sources', { count: memberCount })}{selectedScope && selectedScope.shown !== selectedScope.total && <> · {t('compact.selectedScope', selectedScope)}</>}</h5>
          {sources}</div>
          </fieldset>
        </section>
        <section className="recording-editor-panel" id="recording-editor-method" tabIndex={-1} aria-labelledby="recording-editor-method-title" data-testid="recording-editor-method-note">
          <header><h4 id="recording-editor-method-title">{t('simple.section.method')}</h4></header>
          <fieldset disabled={fieldsDisabled} className="min-w-0 space-y-3" data-testid="group-save-method-fields">{method}
          </fieldset>
        </section>
        <section className="recording-editor-panel" id="recording-editor-destination" tabIndex={-1} aria-labelledby="recording-editor-destination-title" data-testid="recording-editor-destination">
          <header><h4 id="recording-editor-destination-title">{t('simple.section.destination')}</h4></header>
          <fieldset disabled={fieldsDisabled} className="min-w-0 space-y-3" data-testid="group-save-destination-fields"><div className="recording-editor-destination-grid">{destination}</div>{policy}</fieldset>
          <div className="recording-editor-confirmation" data-testid="recording-editor-confirmation">{confirmation}</div>
        </section>
      </div>
      {summary && <aside className="recording-editor-summary" aria-label={t('simple.summary.title')}>{summary}</aside>}
    </div>
  </div>;
}
