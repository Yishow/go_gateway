import { useEffect, useLayoutEffect, useRef, useState, type Dispatch, type MutableRefObject, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import type { Rule, WorkbenchV2State } from '../state/types';
import { workbenchV2Reducer, type WorkbenchV2Action } from '../state/useWorkbenchV2State';
import { hydrateStudioV2Rule, isStudioV2RuleValid } from '../state/studioV2RuleAutosave';
import { RuleEditor } from '../steps/step2/RuleEditor';
import { readStudioLocatorEntry } from '../../../../services/studioV2LocatorEntry';
import { useRecordingDraftGuard } from '../recording-editor/RecordingDraftNavigation';
import type { HomeSourceTarget, HomeRuleObserver } from './RecordingHomeSourceHost';
export type HomeSourceLeave = MutableRefObject<((proceed: () => void) => void) | undefined>;
const editable = [ 'device_id', 'start_address', 'count', 'data_type', 'naming_prefix', 'enabled', 'skipped_addresses', 'scale_multiplier', 'scale_offset', 'data_format', 'share_enabled', 'share_start_register', 'share_stride'] as const satisfies readonly (keyof Rule)[];
/** A local form only: one explicit dispatch hands the original saved-rule queue
 * its unchanged update contract. No new API, retry, CAS or operation receipt. */
export function RecordingHomeSavedRangeEditor({ target, state, rule, dispatch, leave, suspended, list, onLeave, onDone, ruleObserverRef, shareEnabled, onRetryShare, onOpenShareSettings }: {
  shareEnabled?: boolean; onRetryShare: () => void; onOpenShareSettings: () => void;
  target: HomeSourceTarget; state: WorkbenchV2State; rule: Rule; dispatch: Dispatch<WorkbenchV2Action>; leave: HomeSourceLeave; suspended: boolean; list: ReactNode; onLeave: () => void; onDone: () => void; ruleObserverRef: MutableRefObject<HomeRuleObserver | undefined>;
}) {
  const { t } = useTranslation('recording-home');
  const [draft, setDraft] = useState(rule), base = useRef(rule);
  const [phase, setPhase] = useState<'editing' | 'checking' | 'saving' | 'saved' | 'unknown' | 'stale'>('editing');
  const [confirm, setConfirm] = useState<(() => void)>();
  const submitted = useRef<Partial<Rule> | undefined>(undefined), submittedRule = useRef<Rule | undefined>(undefined);
  const [currentRead, setCurrentRead] = useState<'available' | 'unavailable'>();
  const flight = useRef(false), mounted = useRef(true), current = useRef({ rule, suspended }); current.current = { rule, suspended };
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const dirty = editable.some(key => JSON.stringify(draft[key]) !== JSON.stringify(base.current[key]));
  const locked = phase === 'checking' || phase === 'saving' || phase === 'unknown';
  useLayoutEffect(() => {
    leave.current = proceed => { if (locked) return; if (dirty) setConfirm(() => proceed); else proceed(); };
    return () => { leave.current = undefined; };
  }, [dirty, locked, leave]);
  useRecordingDraftGuard({ version: JSON.stringify([target.epoch, phase, draft]), dirty: dirty && !locked, pending: locked, host: true,
    discard: () => { if (!locked) setDraft(base.current); } });
  useLayoutEffect(() => {
    const observer: HomeRuleObserver = {
      submitted: value => {
        if (flight.current && submitted.current && value.id === base.current.id && value.workspace_id === target.workspaceId && value.device_id === target.deviceId &&
          Object.entries(submitted.current).every(([key, expected]) => JSON.stringify(value[key as keyof Rule]) === JSON.stringify(expected))) submittedRule.current = value;
      },
      saved: (value, acknowledgement) => {
        if (!mounted.current || value !== submittedRule.current || !submitted.current) return;
        const saved = hydrateStudioV2Rule(acknowledgement);
        if (saved.id !== value.id || saved.workspace_id !== target.workspaceId || saved.device_id !== target.deviceId ||
          !saved.revision_id || saved.revision_id === base.current.revision_id ||
          !Object.entries(submitted.current).every(([key, expected]) => JSON.stringify(saved[key as keyof Rule]) === JSON.stringify(expected))) { setPhase('unknown'); return; }
        base.current = saved; setDraft(saved); setPhase('saved'); flight.current = false; submittedRule.current = undefined;
      },
      failed: value => { if (mounted.current && value === submittedRule.current) setPhase('unknown'); },
    };
    ruleObserverRef.current = observer;
    return () => { if (ruleObserverRef.current === observer) ruleObserverRef.current = undefined; };
  }, [ruleObserverRef, target.workspaceId, target.deviceId]);
  useEffect(() => {
    if (phase === 'saving' && (rule.save_state === 'save-error' || rule.save_state === 'draft-invalid')) setPhase('unknown');
  }, [phase, rule.save_state]);
  useEffect(() => { if (phase === 'saved' && !suspended) onDone(); }, [phase, suspended, onDone]);
  const localDispatch = (action: WorkbenchV2Action) => {
    if (locked || suspended || !('ruleId' in action) || action.ruleId !== rule.id) return;
    const next = workbenchV2Reducer({ ...state, rules: [draft] }, action).rules.find(value => value.id === rule.id);
    if (next && next.device_id === target.deviceId) { setDraft(next); setPhase('editing'); }
  };
  const save = async () => {
    if (flight.current || locked || suspended || !dirty || !isStudioV2RuleValid(draft)) return;
    flight.current = true; setPhase('checking');
    const original = base.current;
    let result;
    try { result = await readStudioLocatorEntry({ step: 2, workspace_id: target.workspaceId, device_id: target.deviceId, rule_id: original.id,
      mapping_id: null, group_id: null, group_revision: null }); } catch {
      if (mounted.current) { flight.current = false; setPhase('stale'); } return;
    }
    if (!mounted.current) return;
    const fresh = result.ok ? result.value.entry.source_rules?.find(item => item.id === original.id) : undefined;
    if (!fresh || fresh.workspace_id !== target.workspaceId || fresh.device_id !== target.deviceId || fresh.revision_id !== original.revision_id ||
      current.current.suspended || current.current.rule.revision_id !== original.revision_id || current.current.rule.id !== original.id) {
      flight.current = false; setPhase('stale'); return;
    }
    // Only edited fields: do not replay geometry keys and reset skips on a name-only save.
    const patch = Object.fromEntries(editable.filter(key => JSON.stringify(draft[key]) !== JSON.stringify(original[key])).map(key => [key, draft[key]]));
    submitted.current = patch; submittedRule.current = undefined; setPhase('saving'); dispatch({ type: 'updateRule', ruleId: original.id, patch });
  };
  const inspect = async () => {
    try {
      const result = await readStudioLocatorEntry({ step: 2, workspace_id: target.workspaceId, device_id: target.deviceId, rule_id: rule.id, mapping_id: null, group_id: null, group_revision: null });
      if (mounted.current) setCurrentRead(result.ok && result.value.entry.source_rules?.some(item => item.id === rule.id && item.workspace_id === target.workspaceId && item.device_id === target.deviceId) ? 'available' : 'unavailable');
    } catch { if (mounted.current) setCurrentRead('unavailable'); }
  };
  return <><div className="gw-modal-body space-y-4" data-testid="home-saved-range-form">
    {!confirm && list}
    {confirm ? <p>{t('source.discardDescription')}</p> : <>
      <fieldset className="gw-home-saved-range-form" disabled={locked || suspended}><RuleEditor hideLocalName simplePresentation rule={draft} devices={state.devices.filter(device => device.id === target.deviceId)}
        globalShareEnabled={state.settings.modbus_share.enabled} shareControls={{ enabled: shareEnabled, retry: onRetryShare, openSettings: () => leave.current?.(onOpenShareSettings) }} dispatch={localDispatch} /></fieldset>
      {phase === 'saved' && <p role="status" data-testid="home-range-saved">{t('source.saved')}</p>}
      {locked && <p role={phase === 'unknown' ? 'alert' : 'status'} data-testid={`home-range-${phase}`}>{t(`source.${phase}`)}</p>}
      {currentRead && <p role="status">{t(`source.currentRead.${currentRead}`)}</p>}
      {phase === 'stale' && <p role="alert">{t('source.stale')}</p>}
    </>}
  </div><footer className="gw-modal-footer" data-edit-actions style={{ justifyContent: 'flex-end' }}>
    {confirm ? <><button type="button" data-testid="home-source-keep" onClick={() => setConfirm(undefined)}>{t('source.keep')}</button>
      <button type="button" className="gw-primary" data-testid="home-source-discard" onClick={() => { const proceed = confirm; setDraft(base.current); setConfirm(undefined); proceed(); }}>{t('source.discard')}</button></>
      : <><button type="button" data-testid="home-range-done" disabled={locked || suspended} onClick={onDone}>{t('source.cancelEdit')}</button>{phase === 'unknown' && <><button type="button" onClick={() => { void inspect(); }}>{t('source.inspect')}</button><button type="button" onClick={onLeave}>{t('source.leaveSent')}</button></>}<button type="button" className="gw-primary" data-testid="home-range-save" disabled={locked || suspended || !dirty || !isStudioV2RuleValid(draft)} onClick={() => { void save(); }}>{t('source.saveAndReturn')}</button></>}
  </footer></>;
}
