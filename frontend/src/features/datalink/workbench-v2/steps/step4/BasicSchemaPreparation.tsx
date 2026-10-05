import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup, WriteGroupMutationResponse } from '../../../../../types/studioV2WriteGroup';
import { GroupSchemaPanel } from './writeGroup/GroupSchemaPanel';

/** Ensuring a draft changes local configuration only. DDL remains exclusively
 * behind the existing preview token, explicit confirmation and operation ledger. */
export function BasicSchemaPreparation({ group, workspaceId, workspaceRevision, disabled, readonly, destinationMatches, ensureGroup, onApplied }: {
  group?: WriteGroup; workspaceId: string; workspaceRevision: string; disabled: boolean;
  readonly: boolean; destinationMatches: boolean;
  ensureGroup: () => Promise<WriteGroupMutationResponse>; onApplied: () => void;
}) {
  const { t } = useTranslation('workbench-v2');
  const [prepared, setPrepared] = useState<WriteGroupMutationResponse>();
  const [open, setOpen] = useState(false);
  const [busy, setBusy] = useState(false);
  const [failed, setFailed] = useState(false);
  const active = useRef(true);
  const inFlight = useRef(false);
  useEffect(() => { active.current = true; return () => { active.current = false; }; }, []);
  const prepare = async () => {
    if (disabled || readonly || inFlight.current) return;
    inFlight.current = true; setBusy(true); setFailed(false);
    try {
      const result = await ensureGroup();
      if (active.current) { setPrepared(result); setOpen(true); }
    } catch { if (active.current) setFailed(true); }
    finally { inFlight.current = false; if (active.current) setBusy(false); }
  };
  const savedGroup = group ?? prepared?.group;
  const revision = group ? workspaceRevision : prepared?.workspace_revision ?? workspaceRevision;
  return <section className="space-y-3">
    <button type="button" disabled={disabled || readonly || busy} onClick={() => void prepare()} className="rounded border border-slate-600 px-3 py-2 text-xs text-slate-200 disabled:opacity-40" data-testid="basic-recording-prepare">{t('step4.basic.prepare_schema')}</button>
    {failed && <p role="alert" className="text-xs text-amber-300">{t('step4.basic.prepare_failed')}</p>}
    {open && <GroupSchemaPanel group={savedGroup} workspaceId={workspaceId} workspaceRevision={revision} destinationMatches={destinationMatches} dirty={false} readonly={readonly || busy} onApplied={onApplied} />}
  </section>;
}
