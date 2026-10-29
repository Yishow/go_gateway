import { useEffect, useRef, useState, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { recordingTestWriteRecovery, type RecordingTestWriteReference } from '../../state/recordingTestWriteRecovery';
import type { WriteGroupTestWriteOperation } from '../../../../../types/studioV2WriteGroupTestWrite';
import { RecordingTestWriteOperationResult } from './RecordingTestWriteOperationResult';
import { useRecordingTestWriteRecovery } from '../../../../../hooks/datalink/useRecordingTestWriteRecovery';
import { useWriteGroupTestWriteOperationQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroupTestWrite';

interface Observation {
  reference: RecordingTestWriteReference;
  operation: WriteGroupTestWriteOperation;
}
/** A read-only entry has no token and never exposes a repair mutation. */
function RecoveryEntry({ reference, observation, onObserved }: {
  reference: RecordingTestWriteReference;
  observation?: WriteGroupTestWriteOperation;
  onObserved: (value: Observation) => void;
}) {
  const { t } = useTranslation('workbench-v2');
  const query = useWriteGroupTestWriteOperationQuery(reference.operation_id, reference.workspace_id);
  const active = useRef(true);
  useEffect(() => { active.current = true; return () => { active.current = false; }; }, []);
  const check = async () => {
    const result = await query.refetch();
    if (!active.current || result.isError || result.data?.operation_id !== reference.operation_id) return;
    onObserved({ reference, operation: result.data });
    recordingTestWriteRecovery.resolve(result.data, reference.workspace_id);
  };
  return <div>
    <span>{reference.group_id} / {reference.group_revision} / {reference.operation_id}</span>
    <button type="button" onClick={() => void check()} disabled={query.isFetching}
      data-testid={`group-test-write-recovery-query-${reference.operation_id}`} className="ml-2 text-xs underline">
      {t('step4.group.test_write.check')}
    </button>
    {query.isError && <p role="alert">{t('step4.group.test_write.unconfirmed')}</p>}
    {observation && <>
      <p className="text-xs text-slate-400">{t('step4.group.test_write.last_checked_result')}</p>
      <RecordingTestWriteOperationResult operation={observation} testIdPrefix={`recovery-${reference.operation_id}`}
        manualGuidance={observation.cleanup_status === 'failed' || observation.cleanup_status === 'unknown'} />
    </>}
  </div>;
}
interface EntriesProps {
  scopeKey: string;
  references: RecordingTestWriteReference[];
  testId?: string;
  visible?: boolean;
  children?: ReactNode;
}
/** A scope-keyed mount bounds display-only snapshots; it never stores or authorizes an operation. */
export function RecoveryEntries(props: EntriesProps) {
  return <RecoveryEntriesScope key={props.scopeKey} {...props} />;
}
function RecoveryEntriesScope({ references, testId = 'group-test-write-recovery', visible = false, children }: EntriesProps) {
  const { t } = useTranslation('workbench-v2');
  const [observations, setObservations] = useState<Record<string, Observation>>({});
  const entries = new Map(Object.values(observations).map((value) => [value.reference.operation_id, value.reference]));
  references.forEach((reference) => entries.set(reference.operation_id, reference));
  if (!visible && entries.size === 0) return null;
  return <aside data-testid={testId} className="space-y-2 rounded border border-amber-700 p-3">
    {children ?? <p>{t('step4.group.test_write.recovery_results_title')}</p>}
    {[...entries.values()].map((reference) => <RecoveryEntry key={`${reference.workspace_id}|${reference.operation_id}`}
      reference={reference} observation={observations[reference.operation_id]?.operation}
      onObserved={(value) => setObservations((current) => ({ ...current, [reference.operation_id]: value }))} />)}
  </aside>;
}

/** Workspace-level discovery must not depend on a surviving group or open editor. */
export function RecordingTestWriteRecoveryPanel({ workspaceId }: { workspaceId: string }) {
  const { t } = useTranslation('workbench-v2');
  const list = useRecordingTestWriteRecovery(workspaceId);
  const [recovery, setRecovery] = useState(() => recordingTestWriteRecovery.read());
  useEffect(() => {
    const refresh = () => setRecovery(recordingTestWriteRecovery.read());
    refresh();
    const unsubscribe = recordingTestWriteRecovery.subscribe(refresh);
    window.addEventListener('storage', refresh);
    return () => { unsubscribe(); window.removeEventListener('storage', refresh); };
  }, []);
  const references = recovery.references.filter((reference) => reference.workspace_id === workspaceId);
  return <RecoveryEntries scopeKey={workspaceId} references={references} testId="workspace-test-write-recovery" visible={!list.available || recovery.unavailable}>
    <h3>{t('step4.group.test_write.recovery_results_title')}</h3>
    {!list.available && <p role="alert">{t('step4.group.test_write.recovery_list_unavailable')}</p>}
    {recovery.unavailable && <p role="alert">{t('step4.group.test_write.recovery_storage_error')}</p>}
    <button type="button" disabled={list.isFetching} onClick={() => void list.refetch()}>{t('step4.group.test_write.check')}</button>
  </RecoveryEntries>;
}
