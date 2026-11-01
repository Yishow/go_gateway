import type { StudioV2LocatorReadScope } from '../../../../../../services/studioV2WorkspaceDatabase';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { GroupCandidate } from '../../../state/writeGroup/candidates';
import { memberKey, type EditorDraft, type EditorMember } from '../../../state/writeGroup/draft';
import { CounterPolicyEditor } from './CounterPolicyEditor';

/**
 * The saved recording profile of a canonical circuit group, exactly as Basic edits it. The period is shown in its
 * native unit and is never a second, independent seconds setting; layout and storage are fixed once created.
 */
export const GroupProfileSummary: React.FC<{ draft: EditorDraft }> = ({ draft }) => {
  const { t } = useTranslation('workbench-v2');
  const profile = draft.recording_profile;
  if (!profile) return null;
  return (
    <section data-testid="group-profile-summary" aria-label={t('step4.group.profile.title')} className="space-y-2 rounded border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3 text-xs text-[var(--gw-text)]">
      <h4 className="text-sm font-semibold text-[var(--gw-text)]">{t('step4.group.profile.title')}</h4>
      <dl className="grid grid-cols-2 gap-x-4 gap-y-1 sm:grid-cols-3">
        <dt className="text-[var(--gw-muted)]">{t('step4.group.profile.kind_label')}</dt>
        <dd data-testid="group-profile-kind" className="col-span-1 text-[var(--gw-text)]">{t(`step4.group.profile.kind.${profile.kind}`)}</dd>
        <dt className="text-[var(--gw-muted)]">{t('step4.group.profile.layout_label')}</dt>
        <dd data-testid="group-profile-layout" className="text-[var(--gw-text)]">{t(`step4.group.profile.layout.${draft.layout_version}`)}</dd>
        <dt className="text-[var(--gw-muted)]">{t('step4.group.profile.interval_label')}</dt>
        <dd data-testid="group-profile-interval" className="text-[var(--gw-text)]">{profile.interval} {t(`step4.group.profile.unit.${profile.unit}`)}</dd>
        {profile.kind === 'instantaneous' && (
          <>
            <dt className="text-[var(--gw-muted)]">{t('step4.group.profile.reduction_label')}</dt>
            <dd data-testid="group-profile-reduction" className="text-[var(--gw-text)]">{t(`step4.group.profile.reduction.${profile.sample_reduction ?? 'latest'}`)}</dd>
          </>
        )}
        {draft.layout_version === 'managed_long_v1' && (
          <>
            <dt className="text-[var(--gw-muted)]">{t('step4.group.profile.storage_label')}</dt>
            <dd data-testid="group-profile-storage" className="text-[var(--gw-text)]">{t(`step4.group.profile.storage.${draft.storage_mode ?? 'plain'}`)}</dd>
          </>
        )}
      </dl>
      <p className="text-[var(--gw-muted)]">{t('step4.group.profile.note')}</p>
    </section>
  );
};

/** The sources of a profile group. Their columns are fixed by the shared long layout, so there is nothing to map or require. */
export const GroupProfileMembers: React.FC<{ draft: EditorDraft; candidates: GroupCandidate[]; workspaceId: string; readonly: boolean; locatorScope?: StudioV2LocatorReadScope; renderSource?: (member: EditorMember) => React.ReactNode }> = ({ draft, candidates, workspaceId, readonly, locatorScope, renderSource }) => {
  const { t } = useTranslation('workbench-v2');
  const labels = new Map(candidates.map((candidate) => [candidate.key, candidate.label]));
  return (
    <section aria-label={t('step4.group.profile.members_title')} className="space-y-1 text-xs text-[var(--gw-text)]">
      <h4 className="text-sm font-semibold text-[var(--gw-text)]">{t('step4.group.profile.members_title')}</h4>
      <ul data-testid="group-profile-members" className="space-y-1">
        {draft.members.map((member) => (
          <li key={member.key} className="rounded border border-[var(--gw-border)] bg-[var(--gw-surface)] px-2.5 py-1.5">
            <span className="font-medium">{labels.get(memberKey(member.point_id, member.tag_id)) ?? member.tag_id}</span>
            {member.circuit_name && <span className="ml-2 text-[var(--gw-muted)]">{member.circuit_name}</span>}
            {member.value_role && <span className="ml-2 text-[var(--gw-muted)]">{t(`step4.group.profile.role.${member.value_role}`)}</span>}
            {renderSource?.(member)}
          </li>
        ))}
      </ul>
      {draft.layout_version === 'managed_wide_profile_v1'
        ? <p className="text-[var(--gw-muted)]">{t('step4.group.profile.members_note_wide')}</p>
        : <p className="text-[var(--gw-muted)]">{t('step4.group.profile.members_note')}</p>}
      {draft.recording_profile?.kind === 'cumulative' && (
        <CounterPolicyEditor locatorScope={locatorScope} workspaceId={workspaceId} members={draft.members} labels={Object.fromEntries(draft.members.map((m) => [m.tag_id, labels.get(memberKey(m.point_id, m.tag_id)) ?? m.tag_id]))} readonly={readonly} />
      )}
    </section>
  );
};
