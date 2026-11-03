import './recordingSimpleHeader.css';
import { RecordingProductShell } from './RecordingProductShell';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
/** Presentation only. The child is the original Step4 owner, never a second form. */
export function SimpleRecordingShell({ children, onBack, onAdvanced, workspaceId, mode = 'new' }: {
  mode?: 'new' | 'located' | 'copy'; workspaceId?: string; children: ReactNode; onBack: () => void; onAdvanced: () => void;
}) {
  const { t } = useTranslation('recording-editor');
  return <><RecordingProductShell workspaceId={workspaceId} /><main className="gw-simple-recording-main mx-auto w-full max-w-[1440px] px-4 py-6 sm:px-6 lg:px-8" data-testid="simple-recording-shell">
    <header className="gw-simple-recording-heading">
      <div className="gw-simple-recording-title"><button type="button" className="gw-secondary" onClick={onBack}>{t('simple.back')}</button>
        <h1 className="text-2xl font-semibold text-[var(--gw-text)]">{t(mode === 'located' ? 'workflow.editTitle' : 'simple.create')}</h1></div>
      <button type="button" className="gw-secondary" onClick={onAdvanced}>{t('simple.advancedPurpose')}</button>
    </header>
    {children}
  </main></>;
}
