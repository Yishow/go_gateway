import * as React from 'react';
import { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useWriteGroupLifecycleMutation, useWriteGroupDeliveryQuery } from '../../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';
import { normalizeTypedEnvelope } from '../../../../../../utils/safeJson';
import { getSafeErrorMessage } from '../../../../../../utils/typedErrors';

export interface GroupLifecycleBarProps {
  group: WriteGroup;
  workspaceId: string;
  workspaceRevision: string;
  /** The server said this exact revision is ready to apply and the editor has no unsaved edits. */
  canApply: boolean;
  applyBlockedReason?: string;
  readonly: boolean;
  /** Called with the group a lifecycle action returned, so the editor adopts the new revision. */
  onChanged: (group: WriteGroup) => void;
  onConflict: () => void;
}

/**
 * Apply, disable and delete of a saved group. Each action carries the
 * revisions it was shown with, runs once per click, and ends with what the
 * server returned; deleting first shows what stays behind.
 */
export const GroupLifecycleBar: React.FC<GroupLifecycleBarProps> = ({
  group, workspaceId, workspaceRevision, canApply, applyBlockedReason, readonly, onChanged, onConflict,
}) => {
  const { t } = useTranslation('workbench-v2');
  const apply = useWriteGroupLifecycleMutation('apply');
  const disable = useWriteGroupLifecycleMutation('disable');
  const remove = useWriteGroupLifecycleMutation('remove');
  const [confirmingDelete, setConfirmingDelete] = useState(false);
  const [error, setError] = useState<unknown>(null);
  const busyRef = useRef(false);
  const delivery = useWriteGroupDeliveryQuery(group.id, confirmingDelete);
  const deleted = group.status === 'deleted';
  const busy = apply.isPending || disable.isPending || remove.isPending;

  const run = async (mutation: typeof apply) => {
    if (busyRef.current) return;
    busyRef.current = true;
    setError(null);
    try {
      const result = await mutation.mutateAsync({
        id: group.id,
        request: {
          workspace_id: workspaceId, expected_workspace_revision: workspaceRevision,
          expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision,
        },
      });
      setConfirmingDelete(false);
      onChanged(result.group);
    } catch (cause) {
      setError(cause);
      if (normalizeTypedEnvelope(cause).code === 'revision_mismatch') onConflict();
    } finally {
      busyRef.current = false;
    }
  };

  const pending = delivery.data
    ? delivery.data.stages.queued + delivery.data.stages.retrying + delivery.data.stages.blocked +
      delivery.data.stages.quarantined + delivery.data.stages.unknown + delivery.data.stages.collecting
    : null;

  return (
    <div className="space-y-3" data-testid="group-lifecycle">
      <div className="flex flex-wrap gap-2">
        <button
          type="button" onClick={() => void run(apply)} disabled={readonly || busy || deleted || !canApply}
          className="rounded bg-emerald-600 px-3 py-1 text-xs text-white disabled:opacity-40" data-testid="group-apply"
        >
          {t('step4.group.lifecycle.apply')}
        </button>
        <button
          type="button" onClick={() => void run(disable)} disabled={readonly || busy || deleted || group.status === 'disabled'}
          className="rounded border border-slate-600 px-3 py-1 text-xs text-slate-100 disabled:opacity-40" data-testid="group-disable"
        >
          {t('step4.group.lifecycle.disable')}
        </button>
        <button
          type="button" onClick={() => setConfirmingDelete(true)} disabled={readonly || busy || deleted || confirmingDelete}
          className="rounded border border-red-700 px-3 py-1 text-xs text-red-200 disabled:opacity-40" data-testid="group-delete"
        >
          {t('step4.group.lifecycle.delete')}
        </button>
      </div>
      {!canApply && !deleted && applyBlockedReason && (
        <p className="text-xs text-amber-300" data-testid="group-apply-blocked">{applyBlockedReason}</p>
      )}
      <p className="text-xs text-slate-500">{t('step4.group.lifecycle.apply_note')}</p>

      {confirmingDelete && !deleted && (
        <div role="alertdialog" aria-label={t('step4.group.lifecycle.delete_title')} className="space-y-2 rounded border border-red-800 p-3 text-xs text-slate-200" data-testid="group-delete-confirm">
          <p className="font-semibold">{t('step4.group.lifecycle.delete_title')}</p>
          <p>{t('step4.group.lifecycle.delete_keeps')}</p>
          <p data-testid="group-delete-impact">
            {delivery.isError
              ? t('step4.group.lifecycle.delete_impact_unknown')
              : pending === null
                ? t('step4.group.lifecycle.delete_impact_loading')
                : t('step4.group.lifecycle.delete_impact', { count: pending })}
          </p>
          <div className="flex gap-2">
            <button type="button" onClick={() => void run(remove)} disabled={busy} className="rounded bg-red-700 px-3 py-1 text-white disabled:opacity-40" data-testid="group-delete-confirm-button">
              {t('step4.group.lifecycle.delete_confirm')}
            </button>
            <button type="button" onClick={() => setConfirmingDelete(false)} className="rounded border border-slate-600 px-3 py-1" data-testid="group-delete-cancel">
              {t('step4.group.lifecycle.cancel')}
            </button>
          </div>
        </div>
      )}
      {error !== null && (
        <p role="alert" className="text-xs text-red-300" data-testid="group-lifecycle-error">{getSafeErrorMessage(error, t).message}</p>
      )}
    </div>
  );
};
