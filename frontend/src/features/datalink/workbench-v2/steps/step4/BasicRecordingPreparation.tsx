import type { ReactNode } from 'react';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useUpdateWriteGroupMutation } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import { useWriteGroupLifecycleQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroupLifecycle';
import type { WriteGroup, WriteGroupDraft, WriteGroupMutationResponse } from '../../../../../types/studioV2WriteGroup';
import type { GroupCandidate } from '../../state/writeGroup/candidates';
import { basicManagedGroupDraft, isBasicManagedWriteGroup } from '../../state/basicRecording';
import { BasicSchemaPreparation } from './BasicSchemaPreparation';
import { GuidedRecordingPreparation } from './GuidedRecordingPreparation';
import { recordingStartObservation } from './useRecordingPreparationOwner';
import type { RecordingPreparationCommandFence } from './recordingPreparationPresentation';
import type { RenderRecordingScopeSummary } from './recordingScopeModel';
import { inspectSchemaOperation } from './writeGroup/schemaOperationRecovery';

interface BasicRecordingPreparationProps {
  fenceKey: string; onCommandFence: (fence: RecordingPreparationCommandFence | null) => void;
  workspaceId: string; workspaceRevision: string; deviceId: string; ownerKey: string;
  group?: WriteGroup; draft?: WriteGroupDraft; candidates: readonly GroupCandidate[];
  intervalSeconds: number; incompletePolicy: string; disabled: boolean; readonly: boolean; destinationMatches: boolean;
  ensureGroup: () => Promise<WriteGroupMutationResponse>; onPrepared: () => Promise<unknown>;
  start: () => Promise<void>; check: () => Promise<void>;
  live: { unresolved: boolean; operationId?: string; status?: string; queryFailed: boolean; busy: boolean };
  completedOperationId?: string; renderStartResult?: (primary: boolean) => ReactNode;
  summary?: RenderRecordingScopeSummary; isNewCommandBlocked?: () => boolean;
}

/** Presentation adapter for the existing legacy owner, not a second start path. */
export function BasicRecordingPreparation(props: BasicRecordingPreparationProps) {
  const { t } = useTranslation('workbench-v2');
  const { group, workspaceId, workspaceRevision, deviceId, intervalSeconds, incompletePolicy } = props;
  const mutation = useUpdateWriteGroupMutation();
  const lifecycle = useWriteGroupLifecycleQuery(workspaceId, group?.id, Boolean(group && group.workspace_id === workspaceId));
  const [busy, setBusy] = useState(false), [failed, setFailed] = useState(false);
  const inFlight = useRef(false), active = useRef(true);
  useEffect(() => { active.current = true; return () => { active.current = false; }; }, []);
  const draft = group ? basicManagedGroupDraft(group, intervalSeconds, incompletePolicy) : props.draft
    ? { ...props.draft, row_policy: { ...props.draft.row_policy, interval_seconds: intervalSeconds, incomplete_policy: incompletePolicy } } : undefined;
  const dirty = Boolean(group && (group.row_policy.interval_seconds !== intervalSeconds ||
    (group.row_policy.incomplete_policy === 'partial') !== (incompletePolicy === 'partial')));
  const observation = recordingStartObservation(props.ownerKey, props.live);
  const controls = JSON.stringify([props.ownerKey, draft, props.candidates, props.disabled, props.readonly, props.destinationMatches]);
  const identity = JSON.stringify([controls, workspaceRevision, group?.revision]);
  const current = useRef({ identity, ownerKey: props.ownerKey, epoch: 0 });
  if (current.current.identity !== identity) current.current = { identity, ownerKey: props.ownerKey, epoch: current.current.epoch + 1 };
  const epoch = current.current.epoch;
  const reportedFence = useRef<RecordingPreparationCommandFence | null>(null);
  const { onCommandFence, ownerKey, fenceKey } = props;
  const receiveFence = useCallback((fence: RecordingPreparationCommandFence | null) => {
    if (fence && (fence.ownerKey !== ownerKey || fence.scopeKey !== fenceKey)) return;
    reportedFence.current = fence; onCommandFence(fence);
  }, [ownerKey, fenceKey, onCommandFence]);
  const available = !props.disabled && !props.readonly && !busy && !props.live.busy && Boolean(workspaceId && workspaceRevision && deviceId && props.ownerKey && draft) &&
    (!group || group.workspace_id === workspaceId && isBasicManagedWriteGroup(group, deviceId) && !group.recording_profile && group.status !== 'deleted') && props.destinationMatches;
  const currentCommandsBlocked = () => {
    try { const fence = reportedFence.current; return props.isNewCommandBlocked?.() || !fence || fence.scopeKey !== props.fenceKey || fence.isBlocked(group?.id); } catch { return true; }
  };
  const prepare = async () => {
    if (!active.current || current.current.epoch !== epoch || !available || !navigator.onLine || inFlight.current || currentCommandsBlocked() ||
      recordingStartObservation(props.ownerKey, props.live).state !== 'observed' ||
      group && inspectSchemaOperation(workspaceId, group.id).state !== 'empty') return;
    const proof = lifecycle.data;
    const edit = proof?.actions.filter(action => action.action === 'edit');
    if (dirty && (!group || lifecycle.isFetching || lifecycle.isError || !proof || proof.workspace_id !== workspaceId ||
      proof.workspace_revision !== workspaceRevision || proof.group_id !== group.id || proof.group_revision !== group.revision ||
      proof.applied_revision !== group.applied_revision || proof.status !== group.status || proof.connector_id !== group.destination.connector_id ||
      proof.connector_revision !== group.destination.connector_revision || edit?.length !== 1 || edit[0].allowed !== true || edit[0].reason_code !== '')) return;
    inFlight.current = true; setBusy(true); setFailed(false);
    try {
      const result = dirty && group && draft ? await mutation.mutateAsync({ id: group.id, request: {
        workspace_id: workspaceId, expected_workspace_revision: workspaceRevision, expected_group_revision: group.revision,
        expected_connector_revision: group.destination.connector_revision, group: draft,
      } }) : await props.ensureGroup();
      if (!result.workspace_revision || result.group.workspace_id !== workspaceId || !result.group.revision ||
        group && result.group.id !== group.id || !isBasicManagedWriteGroup(result.group, deviceId) ||
        result.group.destination.connector_id !== draft?.destination.connector_id || result.group.destination.connector_revision !== draft?.destination.connector_revision) {
        throw new Error('Preparation scope unavailable');
      }
      // Canonical data remains the only input to the later preview and start review.
      if (active.current && current.current.ownerKey === props.ownerKey) await props.onPrepared();
    } catch { if (active.current && current.current.ownerKey === props.ownerKey) setFailed(true); }
    finally { inFlight.current = false; if (active.current) setBusy(false); }
  };
  return <>
    <GuidedRecordingPreparation isNewCommandBlocked={props.isNewCommandBlocked} ownerKey={props.ownerKey} fenceKey={props.fenceKey} onCommandFence={receiveFence} consentKey={controls} readonly={props.readonly}
      scope={{ workspaceId, workspaceRevision, deviceId, groupIds: group ? [group.id] : [], draft: {
        layout: 'wide', storageMode: 'plain', incompletePolicy: incompletePolicy === 'partial' ? 'partial' : 'skip_row',
        profiles: [{ kind: 'instantaneous', interval: String(intervalSeconds), unit: 'second', sources: props.candidates, reduction: 'latest' }] } }}
      groups={group ? [group] : []} drafts={draft ? [draft] : []} dirtyGroupIds={new Set(group && dirty ? [group.id] : [])}
      completedOperationId={props.completedOperationId} renderStartResult={props.renderStartResult}
      canPrepare={available} summary={props.summary} testIdPrefix="basic-recording"
      start={{ observation, busy: busy || props.live.busy, prepare, start: props.start, check: props.check }}
      renderPreparation={bridge => <BasicSchemaPreparation key={props.ownerKey} guided={bridge} group={group}
        workspaceId={workspaceId} workspaceRevision={workspaceRevision} disabled={!available} dirty={dirty}
        readonly={props.readonly || busy || props.live.busy} destinationMatches={props.destinationMatches}
        ensureGroup={props.ensureGroup} onApplied={() => { void props.onPrepared(); }} />} />
    {failed && <p role="alert" className="text-xs text-amber-300">{t('step4.basic.prepare_failed')}</p>}
  </>;
}
