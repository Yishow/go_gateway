import { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useUpdateWriteGroupMutation } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type { GroupCandidate } from '../../state/writeGroup/candidates';
import { basicMemberChange } from '../../state/basicRecordingMembers';
import { getSafeErrorMessage } from '../../../../../utils/typedErrors';

export function BasicMemberUpdate({ group, deviceId, candidates, workspaceRevision, disabled, onUpdated, onOpenAdvanced }: {
  group: WriteGroup; deviceId: string; candidates: GroupCandidate[]; workspaceRevision: string; disabled: boolean; onUpdated: () => void; onOpenAdvanced?: () => void;
}) {
  const { t } = useTranslation('workbench-v2');
  const mutation = useUpdateWriteGroupMutation();
  const [reviewed, setReviewed] = useState<string>();
  const [error, setError] = useState<unknown>();
  const inFlight = useRef(false);
  const change = basicMemberChange(group, deviceId, candidates);
  const signature = JSON.stringify([workspaceRevision, group, candidates]);
  const currentScope = useRef(signature); currentScope.current = signature;
  const [busy, setBusy] = useState(false);
  if (!change) return <button type="button" onClick={onOpenAdvanced} disabled={disabled} data-testid="basic-recording-members-review">{t('step4.basic.open_advanced')}</button>;
  const save = async () => {
    if (disabled || busy || inFlight.current || reviewed !== signature) return;
    inFlight.current = true; setBusy(true); setError(undefined);
    const requested = signature;
    try {
      await mutation.mutateAsync({ id: group.id, request: {
        workspace_id: group.workspace_id, expected_workspace_revision: workspaceRevision,
        expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision,
        group: change.draft,
      } });
      if (currentScope.current === requested) { setReviewed(undefined); onUpdated(); }
    } catch (cause) { if (currentScope.current === requested) setError(cause); }
    finally { inFlight.current = false; setBusy(false); }
  };
  return <section className="space-y-2 rounded border border-amber-800 p-3 text-xs" data-testid="basic-member-update">
    <button type="button" data-testid="basic-recording-members-review" disabled={disabled || busy} onClick={() => setReviewed(signature)}>{t('step4.basic.memberChange.review')}</button>
    {reviewed === signature && <>
      <p>{t('step4.basic.memberChange.note')}</p>
      <ul>{change.added.map(c => <li key={c.key}>+ {c.device_name} · {c.address} · {c.label}</li>)}{change.removed.map(m => <li key={`${m.point_id}|${m.tag_id}`}>− {m.point_id} · {m.tag_id} · {m.target_column}</li>)}</ul>
      <button type="button" data-testid="basic-members-confirm" disabled={disabled || busy} onClick={() => void save()}>{t('step4.basic.memberChange.confirm')}</button>
      {' · '}<button type="button" disabled={busy} onClick={() => setReviewed(undefined)}>{t('step3.bulk.cancel')}</button>
    </>}
    {error !== undefined && <p role="alert">{getSafeErrorMessage(error, t).message}</p>}
  </section>;
}
