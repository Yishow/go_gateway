import type { ReactNode } from 'react';
import { useRecordingDraftGuard } from './RecordingDraftNavigation';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../../components/ui/button';
import type { SchemaPreparationPresentation } from '../steps/step4/writeGroup/schemaPreparationPresentation';

export interface RecordingSavePresentation {
  saved: boolean; dirty: boolean; busy: boolean; disabled: boolean; unknown: boolean;
  invoke: () => Promise<void>;
}
export interface RecordingPreparationActionsProps {
  /** Simple editor: omits the generic "saving does not start recording" sentence; each stage label already says it. */
  compact?: boolean;
  statusContent?: ReactNode; sourceRefresh?: ReactNode; save: RecordingSavePresentation; schema?: SchemaPreparationPresentation;
  readonly: boolean; blocked: boolean; pendingTestWrite?: boolean; online: boolean; readinessReady: boolean; readinessBusy: boolean;
  checkReadiness: () => Promise<unknown>;
}
/** A single next-action projection. Commands and all operation identities stay in their original owners. */
export function RecordingPreparationActions({ compact = false, statusContent, sourceRefresh, save, schema, readonly, blocked, pendingTestWrite, online, readinessReady, readinessBusy, checkReadiness }: RecordingPreparationActionsProps) {
  useRecordingDraftGuard({ version: JSON.stringify([schema?.scope, schema?.observation]), dirty: false,
    pending: schema?.observation.state === 'unresolved' && Boolean(schema.observation.operationId) });
  const { t } = useTranslation('recording-editor');
  const { t: workbench } = useTranslation('workbench-v2');
  let phase = 'saved', label = 'check', testId = 'recording-preparation-check';
  let command: (() => Promise<unknown>) | undefined;
  if (schema?.observation.state === 'unresolved') {
    phase = schema.commands.check ? 'lookup' : 'checking'; label = 'lookup'; testId = 'recording-preparation-lookup'; command = schema.commands.check;
  } else if (save.busy || save.unknown || save.dirty || !save.saved || readonly) {
    phase = save.busy ? 'saving' : save.unknown ? 'saveUnknown' : !online ? 'offline' : 'save'; label = 'save'; testId = 'group-save';
    if (!save.disabled && online) command = save.invoke;
  } else if (!online) {
    phase = 'offline';
  } else if (sourceRefresh) {
    phase = 'sourceReview';
  } else if (schema?.preview === 'ready') {
    phase = 'review'; label = 'confirm'; testId = 'recording-preparation-confirm'; if (!blocked) command = schema.commands.confirm;
  } else if (readinessBusy) {
    phase = 'checking';
  } else if (blocked) {
    phase = 'blocked';
  } else if (readinessReady) {
    phase = 'ready'; label = 'recheck'; command = checkReadiness;
  } else if (schema?.observation.state === 'observed' && schema.observation.result === 'resolved' && schema.observation.status === 'succeeded') {
    phase = 'schemaConfirmed'; label = 'recheck'; command = checkReadiness;
  } else if (!blocked) {
    command = schema ? schema.commands.retryPreview ?? schema.commands.preview : checkReadiness;
  }
  if (!command && phase === 'saved') phase = 'blocked';
  return <section className="space-y-2 rounded border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3 text-[var(--gw-text)]" data-testid="recording-preparation-actions" data-phase={phase} data-command-available={Boolean(command)}>
    <p role="status" className="font-medium" data-testid="recording-preparation-phase">{t(`preparation.stage.${phase}`)}</p>
    {statusContent && <div id="recording-preparation-status" className="space-y-2">{statusContent}</div>}
    {!compact && phase !== 'sourceReview' && <p className="text-sm text-[var(--gw-muted)]">{t('preparation.noStart')}</p>}
    {pendingTestWrite && <p role="status">{t('preparation.testWritePending')}</p>}
    {schema?.blocker && !pendingTestWrite && !['dirty', 'busy', 'preview_missing'].includes(schema.blocker) && <p className="text-sm">{t(`preparation.blocker.${schema.blocker}`)}</p>}
    {phase === 'sourceReview' ? sourceRefresh : <Button type="button" variant="secondary" size="sm" disabled={!command} onClick={() => { void command?.(); }} data-testid={testId}
      className={`max-w-full h-auto min-h-9 whitespace-normal break-words py-2 border border-[var(--gw-control)] disabled:bg-[var(--gw-surface)] disabled:text-[var(--gw-muted)] focus-visible:outline-solid focus-visible:outline-2 focus-visible:outline-[var(--gw-focus)] focus-visible:outline-offset-2 focus-visible:ring-0 ${phase === 'ready' ? 'bg-[var(--gw-surface)] text-[var(--gw-text)] hover:bg-[var(--gw-active)]' : 'bg-[var(--gw-primary)] text-[var(--gw-on-primary)] hover:bg-[var(--gw-primary-hover)]'}`}>
      {t(`preparation.action.${label}`)}
    </Button>}
    {save.saved && !save.dirty && <span className="ml-3 text-sm" data-testid="group-saved-note">{workbench('step4.group.editor.saved')}</span>}
    {save.saved && save.dirty && <span className="ml-3 text-sm" data-testid="group-dirty-note">{workbench('step4.group.editor.unsaved')}</span>}
  </section>;
}
