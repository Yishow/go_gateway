import { useEffect, useRef, useState, type ReactNode } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { studioV2WorkspaceKeys } from '../../../../hooks/datalink/keys';
import type { WriteGroup, WriteGroupReadiness } from '../../../../types/studioV2WriteGroup';
import type { EditorDraft } from '../state/writeGroup/draft';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import { inspectSchemaOperation } from '../steps/step4/writeGroup/schemaOperationRecovery';
import { recordingTestWriteRecovery } from '../state/recordingTestWriteRecovery';
interface Props {
  workspaceId: string; workspaceRevision: string; group?: WriteGroup; draft: EditorDraft;
  candidates: readonly GroupCandidate[]; blocked: boolean; busy: boolean; dirty: boolean;
  save: (stillCurrent: () => boolean) => Promise<void>;
}
/** Explicit consent to current-source Save semantics, not a source hash/CAS proof. */
export function useRecordingSourceRefresh(props: Props): ReactNode {
  const { t } = useTranslation('recording-editor');
  const client = useQueryClient();
  const mounted = useRef(false);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const query = client.getQueryState<WriteGroupReadiness>(studioV2WorkspaceKeys.writeGroupReadiness(props.group?.id ?? ''));
  const data = query?.data;
  const scoped = query?.status === 'success' && query.fetchStatus === 'idle' && !query.isInvalidated &&
    data?.workspace_id === props.workspaceId && data.workspace_revision === props.workspaceRevision &&
    data.group_id === props.group?.id && data.group_revision === props.group?.revision;
  const stale = Boolean(scoped && data?.issues.some(issue => issue.code === 'source-revision-stale' && issue.severity === 'blocking'));
  const recoverable = Boolean(props.group && props.group.status !== 'deleted');
  const rows = props.draft.members.map(member => props.candidates.filter(candidate => candidate.device_id === member.device_id && candidate.point_id === member.point_id && candidate.tag_id === member.tag_id));
  const noRetainedOperation = () => {
    const recovery = recordingTestWriteRecovery.read();
    return navigator.onLine && !recovery.unavailable && !recovery.references.some(reference => reference.workspace_id === props.workspaceId && reference.group_id === props.group?.id) && inspectSchemaOperation(props.workspaceId, props.group?.id).state === 'empty';
  };
  const canReview = noRetainedOperation() && stale && recoverable && !props.blocked && !props.dirty && rows.length > 0 && rows.every(row => row.length === 1);
  const identity = JSON.stringify([props.workspaceId, props.workspaceRevision, props.group?.id, props.group?.revision, props.draft, props.candidates, props.blocked, props.dirty, query?.dataUpdatedAt]);
  const epoch = useRef({ identity });
  if (epoch.current.identity !== identity) epoch.current = { identity };
  const current = useRef({ props, canReview, token: epoch.current, data }); current.current = { props, canReview, token: epoch.current, data };
  const [review, setReview] = useState<typeof epoch.current>();
  const valid = (token: typeof epoch.current) => {
    const now = current.current;
    const state = client.getQueryState<WriteGroupReadiness>(studioV2WorkspaceKeys.writeGroupReadiness(now.props.group?.id ?? ''));
    return mounted.current && noRetainedOperation() && now.canReview && now.token === token && state?.status === 'success' && state.fetchStatus === 'idle' && !state.isInvalidated && state.data === now.data;
  };
  if (!stale || props.dirty) return null;
  return <section className="space-y-2" data-testid="recording-source-refresh">
    <p>{t('sourceRefresh.impact')}</p>
    {!canReview && <p>{t('sourceRefresh.unavailable')}</p>}
    {rows.some(row => row.length !== 1) && <ul>{rows.map((row, index) => row.length !== 1 && <li key={props.draft.members[index].key}>{t('sourceRefresh.missingMember', { number: index + 1 })}</li>)}</ul>}
    {!canReview && <a href="#recording-editor-sources" className="underline">{t('sourceRefresh.inspect')}</a>}
    {canReview && <button type="button" className="gw-primary" data-testid="recording-source-refresh-review" disabled={props.busy} onClick={() => { if (valid(epoch.current)) setReview(epoch.current); }}>{t('sourceRefresh.review')}</button>}
    {review && valid(review) && <div role="region" aria-label={t('sourceRefresh.review')} className="space-y-2" data-testid="recording-source-refresh-confirmation">
      <p>{t('sourceRefresh.currentAtSave')}</p>
      <ul>{rows.map((row, index) => <li key={props.draft.members[index].key} data-testid="recording-source-refresh-member">{row[0].device_name} · {row[0].label} · {row[0].address}</li>)}</ul>
      <button type="button" className="gw-primary" data-testid="recording-source-refresh-confirm" disabled={props.busy} onClick={() => {
        const captured = review;
        if (valid(captured) && !current.current.props.busy) void current.current.props.save(() => valid(captured)).finally(() => { if (mounted.current) setReview(undefined); });
      }}>{t('sourceRefresh.confirm')}</button>
      <button type="button" className="gw-secondary" data-testid="recording-source-refresh-cancel" disabled={props.busy} onClick={() => setReview(undefined)}>{t('sourceRefresh.cancel')}</button>
    </div>}
  </section>;
}
