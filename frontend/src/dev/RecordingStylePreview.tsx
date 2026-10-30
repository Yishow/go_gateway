import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { semanticTokens } from '@/styles/tokens';
import { RecordingHome } from '@/features/datalink/workbench-v2/recording-home/RecordingHome';
import type { RecordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { SavedRecordingRemoveDialog } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog';

const PREVIEW_MARKER = 'RECORDING_STYLE_PREVIEW_SIMULATED_ONLY';
type PreviewState = 'loading' | 'error' | 'unknown' | 'missing_devices' | 'draft' | 'source_attention' | 'recordings';
/** Fixture props only: no fabricated bootstrap, network reads, operation controllers or persistence. */
function guidance(state: PreviewState): RecordingHomeGuidance {
  return { availability: ['loading', 'error', 'unknown'].includes(state)
    ? { state: 'unavailable', reason: state as 'loading' | 'error' | 'unknown' } : { state: 'available' },
    reason: ['loading', 'error', 'unknown'].includes(state) ? 'evidence_unavailable' : state as RecordingHomeGuidance['reason'],
    primaryAction: null, attentionItems: [], draftItems: [], recordingItems: [] };
}
export default function RecordingStylePreview() {
  const [theme, setTheme] = useState('light');
  const [state, setState] = useState<PreviewState>('unknown');
  const [open, setOpen] = useState(false);
  const [notice, setNotice] = useState('SIMULATED · No backend, credentials or saved preferences.');
  useEffect(() => {
    // This component is mounted only by the isolated development document, never the app provider.
    const root = document.documentElement; const previous = root.className;
    root.classList.remove('light', 'dark'); root.classList.add(theme);
    return () => { root.className = previous; };
  }, [theme]);
  const simulated = () => setNotice('SIMULATED · No backend action was performed.');
  return <div data-preview={PREVIEW_MARKER} style={{ background: semanticTokens.colors.ground, color: semanticTokens.colors.text, fontFamily: semanticTokens.font.sans }}>
    <section className="gw-recording-home" style={{ minHeight: 0 }} aria-label="Development preview controls">
      <h1>Development style preview · 模擬樣式驗收</h1>
      <p role="status">{notice}</p>
      <p>Fixture props only. Tab through controls to inspect focus. This is not an operator workflow.</p>
      <label>Preview theme <select value={theme} onChange={event => setTheme(event.target.value)}><option value="light">Light</option><option value="dark">Dark</option></select></label>
      <label>Preview state <select value={state} onChange={event => setState(event.target.value as PreviewState)}>
        {(['loading', 'error', 'unknown', 'missing_devices', 'draft', 'source_attention', 'recordings'] as const).map(value => <option key={value}>{value}</option>)}
      </select></label>
      <label>Sample input <input placeholder="Simulated recording name" /></label>
      <div className="gw-home-actions">
        <Button className="gw-primary" onClick={() => setOpen(true)}>Preview confirmation</Button>
        <Button disabled>Disabled action</Button>
      </div>
      <Alert variant="destructive"><AlertDescription>Simulated error · No configuration was written.</AlertDescription></Alert>
    </section>
    <RecordingHome guidance={guidance(state)} location={{ search: '?view=recordings', hash: '' }} onNavigate={simulated} onRetry={simulated} />
    <SavedRecordingRemoveDialog open={open} busy={false} objectLabel="SIMULATED recording · not a saved resource"
      onClose={() => setOpen(false)} onConfirm={() => { simulated(); setOpen(false); }} />
  </div>;
}
