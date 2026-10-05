import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../../../../services/studioV2WorkspaceWriteGroups';
import type { DeliveryAttentionItem, DeliveryDecision } from '../../../../../../types/studioV2WriteGroupDelivery';

function DecisionRow({ groupId, item, readonly, onResolved }: { groupId: string; item: DeliveryAttentionItem; readonly: boolean; onResolved: () => void }) {
  const { t } = useTranslation('workbench-v2');
  const [reason, setReason] = useState('');
  const [confirmSkip, setConfirmSkip] = useState(false);
  const [busy, setBusy] = useState(false);
  const [failed, setFailed] = useState(false);
  const [resolved, setResolved] = useState(false);
  const requestRef = useRef<DeliveryDecision | undefined>(undefined);
  const active = useRef(true);
  const inFlight = useRef(false);
  useEffect(() => { active.current = true; return () => { active.current = false; }; }, []);
  const decide = async (resolution: 'retry' | 'skip') => {
    if (readonly || inFlight.current || busy || resolved || !reason.trim() || (resolution === 'skip' && !confirmSkip) || item.state === 'unknown') return;
    const draft = { effect_key: item.effect_key, expected_state: item.state, expected_state_revision: item.state_revision, payload_digest: item.payload_digest, resolution, reason: reason.trim(), confirm_skip: resolution === 'skip' };
    const old = requestRef.current;
    if (!old || JSON.stringify({ ...old, decision_id: undefined }) !== JSON.stringify(draft)) {
      requestRef.current = { ...draft, decision_id: crypto.randomUUID() };
    }
    inFlight.current = true; setBusy(true); setFailed(false);
    try {
      await studioV2WorkspaceWriteGroupsAPI.resolveDelivery(groupId, requestRef.current!);
      if (active.current) { setResolved(true); onResolved(); }
    } catch { if (active.current) setFailed(true); }
    finally { inFlight.current = false; if (active.current) setBusy(false); }
  };
  const actionable = item.state === 'blocked' || item.state === 'quarantined';
  const disabled = readonly || busy || resolved;
  return (
    <li className="space-y-2 rounded border border-amber-900/60 p-3">
      <p className="text-xs text-amber-200">{t(`step4.group.delivery.recovery.${item.state}`)}</p>
      <p className="text-xs text-slate-300">{item.table_schema ? `${item.table_schema}.` : ''}{item.table_name}</p>
      <details className="text-[11px] text-slate-500"><summary>{t('step4.group.delivery.recovery.diagnostics')}</summary><p>{item.effect_key}</p><p>{item.group_revision} · {item.connector_id} · {item.connector_revision}</p></details>
      {actionable && <>
        <label className="block text-xs text-slate-400">{t('step4.group.delivery.recovery.reason')}<input maxLength={200} value={reason} disabled={disabled} onChange={(event) => setReason(event.target.value)} className="mt-1 w-full rounded border border-slate-700 bg-slate-950 p-2" data-testid={`delivery-reason-${item.effect_key}`} /></label>
        <label className="flex items-start gap-2 text-xs text-amber-200"><input type="checkbox" checked={confirmSkip} disabled={disabled} onChange={(event) => setConfirmSkip(event.target.checked)} />{t('step4.group.delivery.recovery.confirm_skip')}</label>
        <div className="flex gap-2">
          <button type="button" onClick={() => void decide('retry')} disabled={disabled || !reason.trim()} className="rounded border border-slate-600 px-3 py-1 text-xs text-slate-200 disabled:opacity-40" data-testid={`delivery-retry-${item.effect_key}`}>{t('step4.group.delivery.recovery.retry')}</button>
          <button type="button" onClick={() => void decide('skip')} disabled={disabled || !reason.trim() || !confirmSkip} className="rounded border border-amber-700 px-3 py-1 text-xs text-amber-200 disabled:opacity-40" data-testid={`delivery-skip-${item.effect_key}`}>{t('step4.group.delivery.recovery.skip')}</button>
        </div>
      </>}
      {busy && <p role="status" className="text-xs text-slate-400">{t('step4.group.delivery.recovery.pending')}</p>}
      {resolved && <p role="status" className="text-xs text-slate-300">{t('step4.group.delivery.recovery.saved')}</p>}
      {failed && <p role="alert" className="text-xs text-amber-300">{t('step4.group.delivery.recovery.failed')}</p>}
    </li>
  );
}

export function GroupDeliveryRecovery({ groupId, items, readonly = false, onResolved }: { groupId: string; items: DeliveryAttentionItem[]; readonly?: boolean; onResolved: () => void }) {
  const { t } = useTranslation('workbench-v2');
  if (items.length === 0) return null;
  return <section className="space-y-2" data-testid="delivery-recovery"><h4 className="text-xs text-slate-200">{t('step4.group.delivery.recovery.title')}</h4><ul className="space-y-2">{items.map((item) => <DecisionRow key={`${groupId}|${item.effect_key}|${item.state}|${item.state_revision}|${item.payload_digest}`} groupId={groupId} item={item} readonly={readonly} onResolved={onResolved} />)}</ul></section>;
}
