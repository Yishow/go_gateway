import { parsePreciseInstant, comparePreciseInstants } from '../../../../utils/preciseInstant';
import { useRef } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { WriteGroupDelivery } from '../../../../types/studioV2WriteGroupDelivery';
import type { RecordingStartOperation } from '../../../../types/studioV2RecordingStart';
import { useWriteGroupDeliveryQuery } from '../../../../hooks/datalink/useStudioV2WriteGroups';
import { recordingHomeProjection } from '../recording-home/recordingHomeProjection';
import { recordingHomeResultSummary } from '../recording-home/recordingHomeResultSummary';
/** Reuses the existing receipt projection; an operation ACK, timestamp or old receipt is not this start's SQL evidence. */
export function savedStartHasMatchingReceipt(group: WriteGroup, workspaceRevision: string, operation: RecordingStartOperation, delivery: WriteGroupDelivery | undefined): boolean {
  const progress = operation.groups.find(value => value.group_id === group.id);
  if (!delivery || operation.workspace_id !== group.workspace_id || operation.status !== 'succeeded' || operation.stage !== 'complete' ||
    !progress?.applied || progress.group_revision !== group.revision || progress.applied_revision !== group.applied_revision) return false;
  const scope = { workspaceId: group.workspace_id, workspaceRevision, epoch: 0 };
  const projection = recordingHomeProjection({ scope,
    groups: { state: 'available', scope, value: { workspace_id: group.workspace_id, workspace_revision: workspaceRevision, groups: [group] } },
    delivery: { [group.id]: { state: 'available', scope: { ...scope, groupId: group.id, groupRevision: group.revision, appliedRevision: group.applied_revision }, value: delivery } },
  });
  if (projection.groups.state !== 'available') return false;
  const facts = recordingHomeResultSummary(projection.groups.value[0]);
  const started = parsePreciseInstant(operation.created_at), committed = parsePreciseInstant(facts.lastCommittedAt);
  return facts.sql === 'revision_receipt' && started !== null && committed !== null && comparePreciseInstants(committed, started) >= 0;
}
/** Polls only the existing revision-scoped delivery query, never starts or repairs anything. */
export function SavedRecordingStartEvidence({ group, workspaceRevision, operation }: { group: WriteGroup; workspaceRevision: string; operation: RecordingStartOperation }) {
  const { t } = useTranslation('recording-editor');
  const identity = JSON.stringify([group.workspace_id, workspaceRevision, group.id, group.revision, group.applied_revision, operation.operation_id]);
  const lease = useRef({ identity, epoch: 0 });
  if (lease.current.identity !== identity) lease.current = { identity, epoch: lease.current.epoch + 1 };
  const readScope = { workspaceId: group.workspace_id, workspaceRevision, epoch: lease.current.epoch, groupRevision: group.revision, appliedRevision: group.applied_revision };
  const delivery = useWriteGroupDeliveryQuery(group.id, Boolean(group.applied_revision), 5000, readScope);
  const committed = !delivery.isError && savedStartHasMatchingReceipt(group, workspaceRevision, operation, delivery.data);
  return <div className="space-y-2">
    <p role="status" data-testid={committed ? 'saved-recording-start-committed' : 'saved-recording-start-waiting'}>{t(committed ? 'start.committed' : 'start.waiting')}</p>
    {committed && <p className="text-sm text-[var(--gw-muted)]">{t('start.notReadback')}</p>}
    {delivery.isError && <p role="alert">{t('start.deliveryUnavailable')}</p>}
  </div>;
}
