import { SavedRecordingManagement, type SavedRecordingManagementProps } from './SavedRecordingManagement';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useWriteGroupDeliveryQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { CircuitProfileStatusView } from './CircuitProfileStatusView';

export interface CircuitProfileStatusProps {
  group: WriteGroup;
  lifecycle?: Omit<SavedRecordingManagementProps, 'group'>;
  sourceLabels?: Record<string, string>;
  /** Opens this exact saved group in Advanced. */
  onOpenAdvanced?: (groupId: string) => void;
  /** The Basic checkbox of this profile is unchecked: only the draft changed. */
  draftOff?: boolean;
}

const REFRESH_MS = 5000;

/** Reads the real delivery answer of one canonical profile group; an unapplied group is not queried. */
export const CircuitProfileStatus: React.FC<CircuitProfileStatusProps> = ({ group, sourceLabels, onOpenAdvanced, draftOff = false, lifecycle }) => {
  const { t } = useTranslation('workbench-v2');
  const kind = group.recording_profile?.kind ?? 'instantaneous';
  const applied = group.applied_revision !== '';
  const query = useWriteGroupDeliveryQuery(group.id, applied, REFRESH_MS);
  return (
    <section data-testid={`circuit-status-${kind}`} aria-label={t(`step4.basic.circuit.${kind}_label`)} className="space-y-1.5 rounded border border-slate-800/80 bg-slate-950/30 p-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <h4 className="text-xs font-semibold text-slate-200">{t(`step4.basic.circuit.${kind}_label`)} · {t('step4.basic.status.title')}</h4>
        {onOpenAdvanced && (
          <button type="button" data-testid={`circuit-status-open-advanced-${kind}`} onClick={() => onOpenAdvanced(group.id)}
            className="rounded border border-slate-700 px-2 py-0.5 text-[11px] text-slate-300 hover:bg-slate-800">
            {t('step4.basic.status.open_advanced')}
          </button>
        )}
      </div>
      {!applied && <p className="text-xs text-slate-400">{t('step4.basic.status.not_applied')}</p>}
      {draftOff && applied && group.status !== 'disabled' && group.status !== 'deleted' && (
        <div data-testid={`circuit-profile-draft-note-${kind}`} role="status" className="space-y-1 text-xs text-amber-300">
          <p>{t('step4.basic.circuit.draft_not_stop')}</p>
          {onOpenAdvanced && (
            <button type="button" data-testid={`circuit-profile-draft-lifecycle-${kind}`} onClick={() => onOpenAdvanced(group.id)} className="rounded border border-amber-600/60 px-2 py-1">
              {t('step4.basic.circuit.draft_open_lifecycle')}
            </button>
          )}
        </div>
      )}
      {applied && query.isError && <p role="status" className="text-xs text-amber-300">{t('step4.basic.status.query_failed')}</p>}
      {applied && !query.isError && <CircuitProfileStatusView group={group} delivery={query.data} sourceLabels={sourceLabels} onOpenAdvanced={onOpenAdvanced} />}
      {lifecycle?.lifecycleHost && <SavedRecordingManagement {...lifecycle} group={group} />}
    </section>
  );
};
