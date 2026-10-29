import { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useWriteGroupTestWriteOperationQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroupTestWrite';
import type { WriteGroupTestWriteOperation } from '../../../../../types/studioV2WriteGroupTestWrite';
import { recordingTestWriteRecovery, type RecordingTestWriteReference } from '../../state/recordingTestWriteRecovery';
import { RecordingTestWriteOperationResult } from './RecordingTestWriteOperationResult';

export interface RecordingPreparationTestWriteQueryReport {
  readonly operationId: string;
  readonly check: () => Promise<void>;
}
export interface RecordingPreparationTestWriteQueryProps {
  reference: RecordingTestWriteReference;
  ownerIdentity: string;
  onQuery: (report: RecordingPreparationTestWriteQueryReport | null) => void;
}
const fields = ['workspace_id', 'group_id', 'group_revision', 'operation_id', 'operation_kind', 'created_at'] as const;
function retainsOriginal(references: RecordingTestWriteReference[], original: Readonly<RecordingTestWriteReference>) {
  const matching = references.filter(reference => reference.operation_id === original.operation_id);
  return matching.length === 1 && fields.every(field => matching[0][field] === original[field]);
}
interface Observation {
  key: string; epoch: number; failed: boolean; unreadable: boolean;
  operation?: WriteGroupTestWriteOperation;
}

/** Explicit readonly lookup of the immutable retained locator, never a new test write or repair. */
export function RecordingPreparationTestWriteQuery({ reference, ownerIdentity, onQuery }: RecordingPreparationTestWriteQueryProps) {
  const { t } = useTranslation('workbench-v2');
  const original = useMemo(() => Object.freeze({ workspace_id: reference.workspace_id, group_id: reference.group_id,
    group_revision: reference.group_revision, operation_id: reference.operation_id, operation_kind: reference.operation_kind, created_at: reference.created_at }),
  [reference.workspace_id, reference.group_id, reference.group_revision, reference.operation_id, reference.operation_kind, reference.created_at]);
  const key = JSON.stringify([ownerIdentity, ...fields.map(field => original[field])]);
  const fence = useRef({ key, epoch: 0 });
  if (fence.current.key !== key) fence.current = { key, epoch: fence.current.epoch + 1 };
  const epoch = fence.current.epoch;
  const lease = useMemo(() => ({ active: true, epoch, report: onQuery }), [epoch, onQuery]);
  const currentLease = useRef(lease); currentLease.current = lease;
  useLayoutEffect(() => { lease.active = true; return () => { lease.active = false; }; }, [lease]);
  const valid = Boolean(ownerIdentity && original.operation_kind === 'test_write' && fields.every(field =>
    typeof original[field] === 'string' && original[field].length > 0));
  const current = useCallback(() => valid && lease.active && currentLease.current === lease && fence.current.epoch === epoch, [valid, lease, epoch]);
  const query = useWriteGroupTestWriteOperationQuery(original.operation_id, original.workspace_id);
  const refetch = query.refetch;
  const inFlight = useRef<{ lease: typeof lease; promise: Promise<void> } | undefined>(undefined);
  const [observation, setObservation] = useState<Observation>();
  const check = useCallback((): Promise<void> => {
    if (!current()) return Promise.resolve();
    if (inFlight.current?.lease === lease) return inFlight.current.promise;
    const before = recordingTestWriteRecovery.read();
    if (!retainsOriginal(before.references, original)) {
      setObservation({ key, epoch, failed: true, unreadable: before.unreadable });
      return Promise.resolve();
    }
    setObservation({ key, epoch, failed: false, unreadable: before.unreadable });
    const pending: Promise<void> = (async () => {
      const result = await refetch({ cancelRefetch: false });
      if (!current()) return;
      const after = recordingTestWriteRecovery.read();
      if (!retainsOriginal(after.references, original) || result.isError || result.data?.operation_id !== original.operation_id) {
        setObservation({ key, epoch, failed: true, unreadable: after.unreadable });
        return;
      }
      setObservation({ key, epoch, failed: false, unreadable: after.unreadable, operation: result.data });
      // Known inaccessible storage still permits inspecting O1, never dropping it.
      // A later explicit healthy read may resolve it through the unchanged store.
      if (!before.unreadable && !after.unreadable) recordingTestWriteRecovery.resolve(result.data, original.workspace_id);
    })().finally(() => { if (inFlight.current?.promise === pending) inFlight.current = undefined; });
    inFlight.current = { lease, promise: pending };
    return pending;
  }, [current, lease, original, key, epoch, refetch]);
  useEffect(() => {
    onQuery({ operationId: original.operation_id, check });
    return () => onQuery(null);
  }, [original.operation_id, check, onQuery]);
  const visible = observation?.key === key && observation.epoch === epoch ? observation : undefined;
  return <div data-testid="preparation-test-write-query" aria-busy={query.isFetching}>
    {visible?.unreadable && <p role="alert">{t('step4.group.test_write.recovery_storage_error')}</p>}
    {visible?.failed && <p role="alert">{t('step4.group.test_write.unconfirmed')}</p>}
    {visible?.operation && <>
      <p className="text-xs text-slate-400">{t('step4.group.test_write.last_checked_result')}</p>
      <RecordingTestWriteOperationResult operation={visible.operation} testIdPrefix="preparation-test-write"
        manualGuidance={visible.operation.cleanup_status === 'failed' || visible.operation.cleanup_status === 'unknown'} />
    </>}
  </div>;
}
