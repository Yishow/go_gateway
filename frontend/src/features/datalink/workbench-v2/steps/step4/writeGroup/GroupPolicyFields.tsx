import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { EditorDraft, IncompletePolicy } from '../../../state/writeGroup/draft';

const FIELD = 'mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-1 text-sm text-slate-100';

export interface GroupPolicyFieldsProps {
  draft: EditorDraft;
  readonly: boolean;
  /** A profile group's period lives in its recording profile and its rows are never skipped or completed per row. */
  profile: boolean;
  onChange: (patch: Partial<EditorDraft>) => void;
  onDedupe: (capability: string) => void;
}

/** Row policy of a group. Only the wide and custom layouts have a seconds interval and a completeness policy. */
export const GroupPolicyFields: React.FC<GroupPolicyFieldsProps> = ({ draft, readonly, profile, onChange, onDedupe }) => {
  const { t } = useTranslation('workbench-v2');
  return (
    <details className="rounded border border-slate-800 p-3" data-testid="group-advanced">
      <summary className="cursor-pointer text-xs font-medium text-slate-300">{t('step4.group.editor.advanced')}</summary>
      <div className="mt-3 grid gap-3 sm:grid-cols-2">
        {!profile && (
          <label className="text-xs text-slate-400">{t('step4.group.editor.interval')}
            <input type="number" min={1} value={draft.interval_seconds} disabled={readonly || draft.recording_profile !== undefined} onChange={(event) => onChange({ interval_seconds: Number(event.target.value) })} className={FIELD} data-testid="group-interval" />
          </label>
        )}
        <label className="text-xs text-slate-400">{t('step4.group.editor.lateness')}
          <input type="number" min={0} value={draft.allowed_lateness_seconds} disabled={readonly} onChange={(event) => onChange({ allowed_lateness_seconds: Number(event.target.value) })} className={FIELD} data-testid="group-lateness" />
        </label>
        {!profile && (
          <>
            <label className="text-xs text-slate-400">{t('step4.group.editor.incomplete')}
              <select value={draft.incomplete_policy} disabled={readonly} onChange={(event) => onChange({ incomplete_policy: event.target.value as IncompletePolicy })} className={FIELD} data-testid="group-incomplete">
                <option value="skip_row">{t('step4.group.editor.incomplete_skip')}</option>
                <option value="partial">{t('step4.group.editor.incomplete_partial')}</option>
              </select>
            </label>
            <label className="text-xs text-slate-400">{t('step4.group.editor.dedupe')}
              <select value={draft.storage_strategy === 'managed' ? 'receipt' : draft.dedupe_capability} disabled={readonly || draft.storage_strategy === 'managed'} onChange={(event) => onDedupe(event.target.value)} className={FIELD} data-testid="group-dedupe">
                {draft.storage_strategy !== 'managed' && <option value="">{t('step4.group.editor.dedupe_none')}</option>}
                <option value="receipt">{t('step4.group.editor.dedupe_receipt')}</option>
              </select>
            </label>
            <label className="text-xs text-slate-400">{t('step4.group.editor.entity_column')}
              <input type="text" value={draft.entity_key_column} disabled={readonly} onChange={(event) => onChange({ entity_key_column: event.target.value })} className={FIELD} data-testid="group-entity-column" />
            </label>
            <label className="text-xs text-slate-400">{t('step4.group.editor.provenance_column')}
              <input type="text" value={draft.provenance_column} disabled={readonly} onChange={(event) => onChange({ provenance_column: event.target.value })} className={FIELD} data-testid="group-provenance-column" />
            </label>
          </>
        )}
      </div>
    </details>
  );
};
