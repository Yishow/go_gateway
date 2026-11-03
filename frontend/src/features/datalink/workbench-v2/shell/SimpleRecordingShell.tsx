import './recordingSimpleHeader.css';
import { RecordingProductShell } from './RecordingProductShell';
import { useLayoutEffect, useRef, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
/** Presentation only. The child is the original Step4 owner, never a second form. */
export function SimpleRecordingShell({ children, onBack, onAdvanced, workspaceId, mode = 'new' }: {
  mode?: 'new' | 'located' | 'copy'; workspaceId?: string; children: ReactNode; onBack: () => void; onAdvanced: () => void;
}) {
  const { t } = useTranslation('recording-editor');
  const main = useRef<HTMLElement>(null), heading = useRef<HTMLElement>(null);
  useLayoutEffect(() => {
    const root = main.current, element = heading.current;
    if (!root || !element) return;
    const update = () => root.style.setProperty('--recording-heading-height', `${Math.ceil(element.getBoundingClientRect().height)}px`);
    update(); const observer = typeof ResizeObserver === 'undefined' ? undefined : new ResizeObserver(update);
    observer?.observe(element); window.addEventListener('resize', update);
    return () => { observer?.disconnect(); window.removeEventListener('resize', update); root.style.removeProperty('--recording-heading-height'); };
  }, []);
  return <><RecordingProductShell workspaceId={workspaceId} /><main ref={main} className="gw-simple-recording-main mx-auto w-full max-w-[1440px] px-4 py-6 sm:px-6 lg:px-8" data-testid="simple-recording-shell">
    <header ref={heading} className="gw-simple-recording-heading">
      <div className="gw-simple-recording-title"><button type="button" className="gw-secondary" onClick={onBack}>{t('simple.back')}</button>
        <h1 className="text-2xl font-semibold text-[var(--gw-text)]">{t(mode === 'located' ? 'workflow.editTitle' : 'simple.create')}</h1></div>
      <button type="button" className="gw-secondary" onClick={onAdvanced}>{t('simple.advancedPurpose')}</button>
    </header>
    {children}
  </main></>;
}
