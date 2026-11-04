import { ParameterHelp } from '../../../components/ParameterHelp';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { EditorDraft, IncompletePolicy } from '../../../state/writeGroup/draft';

const FIELD = 'mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1 text-sm text-[var(--gw-text)]';

export interface GroupPolicyFieldsProps {
  expanded?: boolean;
  draft: EditorDraft;
  readonly: boolean;
  /** A profile group's period lives in its recording profile and its rows are never skipped or completed per row. */
  profile: boolean;
  onChange: (patch: Partial<EditorDraft>) => void;
  onDedupe: (capability: string) => void;
}

/** Row policy of a group. Only the wide and custom layouts have a seconds interval and a completeness policy. */
export const GroupPolicyFields: React.FC<GroupPolicyFieldsProps> = ({ expanded = false, draft, readonly, profile, onChange, onDedupe }) => {
  const { t } = useTranslation('workbench-v2');
  const fields = (
      <div className="grid gap-3 sm:grid-cols-2">
        {!profile && (
          <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">{t('step4.group.editor.interval')}
            <input type="number" min={1} value={draft.interval_seconds} disabled={readonly || draft.recording_profile !== undefined} onChange={(event) => onChange({ interval_seconds: Number(event.target.value) })} className={FIELD} data-testid="group-interval" />
          </label><ParameterHelp help="policy.interval" /></div>
        )}
        <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">{t('step4.group.editor.lateness')}
          <input type="number" min={0} value={draft.allowed_lateness_seconds} disabled={readonly} onChange={(event) => onChange({ allowed_lateness_seconds: Number(event.target.value) })} className={FIELD} data-testid="group-lateness" />
        </label><ParameterHelp help="policy.lateness" /></div>
        {!profile && (
          <>
            <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">{t('step4.group.editor.incomplete')}
              <select value={draft.incomplete_policy} disabled={readonly} onChange={(event) => onChange({ incomplete_policy: event.target.value as IncompletePolicy })} className={FIELD} data-testid="group-incomplete">
                <option value="skip_row">{t('step4.group.editor.incomplete_skip')}</option>
                <option value="partial">{t('step4.group.editor.incomplete_partial')}</option>
              </select>
            </label><ParameterHelp help="policy.incomplete" /></div>
            <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">{t('step4.group.editor.dedupe')}
              <select value={draft.storage_strategy === 'managed' ? 'receipt' : draft.dedupe_capability} disabled={readonly || draft.storage_strategy === 'managed'} onChange={(event) => onDedupe(event.target.value)} className={FIELD} data-testid="group-dedupe">
                {draft.storage_strategy !== 'managed' && <option value="">{t('step4.group.editor.dedupe_none')}</option>}
                <option value="receipt">{t('step4.group.editor.dedupe_receipt')}</option>
              </select>
            </label><ParameterHelp help="policy.dedupe" /></div>
            <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">{t('step4.group.editor.entity_column')}
              <input type="text" value={draft.entity_key_column} disabled={readonly} onChange={(event) => onChange({ entity_key_column: event.target.value })} className={FIELD} data-testid="group-entity-column" />
            </label><ParameterHelp help="policy.entityColumn" /></div>
            <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">{t('step4.group.editor.provenance_column')}
              <input type="text" value={draft.provenance_column} disabled={readonly} onChange={(event) => onChange({ provenance_column: event.target.value })} className={FIELD} data-testid="group-provenance-column" />
            </label><ParameterHelp help="policy.provenanceColumn" /></div>
          </>
        )}
      </div>
  );
  return (
    expanded
      // The simple editor shows these rules as a plain section; only the guided flow folds them.
      ? <fieldset className="rounded border border-[var(--gw-control)] p-3" data-testid="group-advanced">
        <legend className="px-1 text-xs font-medium text-[var(--gw-text)]">{t('step4.group.editor.rules')}</legend>
        {fields}
      </fieldset>
      : <details className="rounded border border-[var(--gw-control)] p-3" data-testid="group-advanced">
        <summary className="cursor-pointer text-xs font-medium text-[var(--gw-text)]">{t('step4.group.editor.advanced')}</summary>
        <div className="mt-3">{fields}</div>
      </details>
  );
};
