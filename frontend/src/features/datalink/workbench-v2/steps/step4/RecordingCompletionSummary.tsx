import { useId, useLayoutEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { safeRuntimeLabel } from '../../../../../utils/runtimeOperatorSafety';
import type { RecordingCompletionView } from './recordingCompletion';

export interface RecordingCompletionSummaryProps {
  view: RecordingCompletionView;
  deviceNames?: Readonly<Record<string, string>>;
  groupNames?: Readonly<Record<string, string>>;
  primary?: boolean;
  onContinue?: (deviceId: string) => void;
}

/** Read-only start results; the owner supplies verified scope and retains all operation details. */
export function RecordingCompletionSummary({ view, deviceNames, groupNames, primary, onContinue }: RecordingCompletionSummaryProps) {
  const { t } = useTranslation('workbench-v2-recording-guided');
  const headingId = useId();
  const identity = JSON.stringify(view);
  const owner = useRef({ identity, epoch: 0 });
  if (owner.current.identity !== identity) owner.current = { identity, epoch: owner.current.epoch + 1 };
  const epoch = owner.current.epoch;
  const mounted = useRef(true), navigated = useRef<number | undefined>(undefined);
  const current = useRef({ view, onContinue }); current.current = { view, onContinue };
  useLayoutEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const handoff = view.state === 'available' && (view.status === 'partial' || view.status === 'succeeded') &&
    view.activatedDeviceIds.length === 1 && view.handoffDeviceId === view.activatedDeviceIds[0] &&
    view.handoffDeviceId && view.deviceIds.includes(view.handoffDeviceId) ? view.handoffDeviceId : undefined;
  const navigate = () => {
    if (!mounted.current || owner.current.epoch !== epoch || owner.current.identity !== identity || navigated.current === epoch ||
      JSON.stringify(current.current.view) !== identity || !handoff || !current.current.onContinue) return;
    navigated.current = epoch;
    current.current.onContinue(handoff);
  };
  const label = (id: string, names: Readonly<Record<string, string>> | undefined, scope: readonly string[], kind: 'device' | 'group') => {
    const alias = names && Object.hasOwn(names, id) ? safeRuntimeLabel(names[id]) : '';
    return alias || t(`completion.${kind}Fallback`, { index: scope.indexOf(id) + 1 });
  };
  const list = (ids: readonly string[], kind: 'device' | 'group', title: string) => {
    if (view.state !== 'available' || ids.length === 0) return null;
    return <div className="space-y-1">
      <h4 className="font-medium text-slate-200">{title}</h4>
      <ul aria-label={title} className="max-h-40 list-inside list-disc overflow-y-auto break-words text-slate-300">
        {ids.map(id => <li key={id}>{label(id, kind === 'device' ? deviceNames : groupNames, kind === 'device' ? view.deviceIds : view.groupIds, kind)}</li>)}
      </ul>
    </div>;
  };
  return <section aria-labelledby={headingId} data-testid="recording-completion-summary" className="min-w-0 space-y-3 rounded border border-slate-700 bg-slate-950/60 p-3 text-sm">
    <h3 id={headingId} className="font-semibold text-slate-100">{t('completion.heading')}</h3>
    {view.state === 'unavailable' ? <p role="status" className="text-amber-300">{t(`reasons.${view.reason}`)}</p> : <>
      <p role="status" className="text-slate-200">{t(`states.${view.status}`)}</p>
      {view.activatedDeviceIds.length ? list(view.activatedDeviceIds, 'device', t('completion.activated'))
        : <p className="text-slate-300">{t('completion.noneActivated')}</p>}
      {list(view.incompleteDeviceIds, 'device', t('completion.incompleteDevices'))}
      {list(view.incompleteGroupIds, 'group', t('completion.incompleteGroups'))}
      <p className="text-xs text-slate-400">{t('completion.deliverySeparate')}</p>
      {handoff && onContinue && <button type="button" data-testid="basic-recording-continue-runtime" data-primary={primary ? 'true' : undefined}
        className={primary ? 'rounded bg-cyan-600 px-4 py-2 font-semibold text-white' : 'rounded border border-emerald-700 px-3 py-2 text-emerald-200'}
        onClick={navigate}>{t('completion.continue')}</button>}
    </>}
  </section>;
}
