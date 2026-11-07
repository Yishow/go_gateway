import { useTranslation } from 'react-i18next';
import { DeleteResponseError } from '@/services/studioV2WorkspaceDelete';

const knownStates = new Set(['pending', 'queued', 'sending', 'retrying', 'blocked', 'quarantined', 'unknown']);

/** Recovery actions never show server error text or combine different units. */
export function RecordingDeleteDialogError({ error, onStop, onOpenRecording, onRefresh }: {
  error: unknown; onStop?: () => void; onOpenRecording?: () => void; onRefresh?: () => void;
}) {
  const { t } = useTranslation('recording-home');
  if (!error) return null;
  const known = error instanceof DeleteResponseError ? error : undefined;
  const code = known?.code;
  const backlog = code === 'WRITE_GROUP_DELETE_BACKLOG' && known?.details && 'outbox' in known.details ? known.details : undefined;
  const successors = code === 'WRITE_GROUP_DELETE_HAS_SUCCESSOR' && known?.details && 'referenced_by' in known.details ? known.details : undefined;
  const message = code === 'WRITE_GROUP_DELETE_ACTIVE' ? 'active' : code === 'WRITE_GROUP_DELETE_BACKLOG' ? 'backlog' :
    code === 'WRITE_GROUP_DELETE_OPERATION_PENDING' ? 'operationPending' : code === 'WRITE_GROUP_DELETE_HAS_SUCCESSOR' ? 'hasSuccessor' : code === 'revision_mismatch' ? 'stale' : known?.outcome === 'failed' ? 'refused' : 'unknown';
  return <div className="gw-warning" role="alert">
    <p>{t(`deleteDialogs.recording.${message}`)}</p>
    {successors && <><ul>{successors.referenced_by.map(recording => <li key={recording.id}>{recording.name}</li>)}</ul>
      {successors.total > successors.referenced_by.length && <p data-testid="delete-successor-total">{t('deleteDialogs.recording.referenceTotal', { count: successors.total })} {successors.total}</p>}</>}
    {backlog && <dl>
      <div data-testid="delete-backlog-samples"><dt>{t('deleteDialogs.recording.samples')}</dt><dd>{backlog.samples_unconsumed}</dd></div>
      {Object.entries(backlog.outbox).map(([state, count], index) => <div key={state} data-testid={`delete-backlog-${knownStates.has(state) ? state : `other-${index}`}`}>
        <dt>{t(`deleteDialogs.recording.outbox.${knownStates.has(state) ? state : 'other'}`)}</dt><dd>{count}</dd>
      </div>)}
    </dl>}
    <div className="gw-delete-recovery-actions">
      {(code === 'WRITE_GROUP_DELETE_ACTIVE' || backlog) && onStop && <button type="button" onClick={onStop}>{t('deleteDialogs.recording.stop')}</button>}
      {(code === 'WRITE_GROUP_DELETE_OPERATION_PENDING' || code === 'WRITE_GROUP_DELETE_HAS_SUCCESSOR') && onOpenRecording &&
        <button type="button" onClick={onOpenRecording}>{t('deleteDialogs.recording.openRecording')}</button>}
      {onRefresh && <button type="button" onClick={onRefresh}>{t('deleteDialogs.recording.refresh')}</button>}
    </div>
  </div>;
}
