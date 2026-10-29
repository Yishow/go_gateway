import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

/** Secondary recovery presentation only: the private owner verifies O1 and retains the unchanged request. */
export function BasicRecordingRetry({ context, eligible, disabled, verify, onRetry }: {
  context: string; eligible: boolean; disabled: boolean; verify: () => Promise<boolean>; onRetry: () => void;
}) {
  const { t } = useTranslation('workbench-v2');
  const [checking, setChecking] = useState(false), [rejectedAt, setRejectedAt] = useState<number>();
  const [online, setOnline] = useState(() => navigator.onLine);
  const mounted = useRef(true), inFlight = useRef(false);
  const identity = JSON.stringify([context, eligible, disabled, online]);
  const owner = useRef({ identity, epoch: 0 });
  if (owner.current.identity !== identity) owner.current = { identity, epoch: owner.current.epoch + 1 };
  const current = useRef({ disabled, eligible, onRetry }); current.current = { disabled, eligible, onRetry };
  useEffect(() => {
    mounted.current = true;
    const refresh = () => setOnline(navigator.onLine);
    window.addEventListener('online', refresh); window.addEventListener('offline', refresh);
    return () => { mounted.current = false; window.removeEventListener('online', refresh); window.removeEventListener('offline', refresh); };
  }, []);
  const epoch = owner.current.epoch;
  const retry = async () => {
    if (!mounted.current || owner.current.epoch !== epoch || inFlight.current || current.current.disabled || !current.current.eligible || !navigator.onLine) return;
    const active = () => mounted.current && owner.current.epoch === epoch && current.current.eligible && !current.current.disabled && navigator.onLine;
    inFlight.current = true; setChecking(true); setRejectedAt(undefined);
    try {
      const verified = await verify();
      if (!active()) return;
      if (verified) current.current.onRetry(); else setRejectedAt(epoch);
    } catch { if (active()) setRejectedAt(epoch); }
    finally { inFlight.current = false; if (mounted.current) setChecking(false); }
  };
  if (!eligible) return null;
  return <div className="space-y-1">
    <button type="button" data-testid="basic-recording-retry" disabled={disabled || checking || !online}
      className="rounded border border-slate-600 px-3 py-2 text-xs text-slate-200 disabled:opacity-40" onClick={() => { void retry(); }}>{t('step4.basic.retry')}</button>
    {rejectedAt === owner.current.epoch && <p role="status" className="text-xs text-amber-300">{t('step4.basic.action.retry_same_request')}</p>}
  </div>;
}
