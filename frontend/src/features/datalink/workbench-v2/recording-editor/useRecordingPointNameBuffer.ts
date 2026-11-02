import { useRef, useState } from 'react';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import type { PointNameOutcome, RecordingPointNamePort } from '../../../../pages/datalink/workbench-v2/recordingPointNameSubmission';
export function useRecordingPointNameBuffer(scope: string, port?: RecordingPointNamePort, disabled = false) {
  const current = useRef({ scope, port, disabled }); current.current = { scope, port, disabled };
  const sequence = useRef(0);
  const [value, setValue] = useState<{ scope: string; candidate: GroupCandidate; name: string; outcome?: PointNameOutcome; preparing: boolean; failed: boolean }>();
  const active = value?.scope === scope ? value : undefined;
  const begin = async (candidate: GroupCandidate) => {
    if (current.current.disabled || !port || active?.outcome?.status === 'pending') return false;
    const request = ++sequence.current, captured = scope;
    setValue({ scope, candidate, name: candidate.label, preparing: true, failed: false });
    try {
      const outcome = await port.prepare(candidate);
      if (sequence.current !== request || current.current.scope !== captured || current.current.port !== port) return false;
      setValue({ scope, candidate, name: outcome?.submittedName ?? outcome?.session.displayName ?? candidate.label, outcome, preparing: false, failed: !outcome });
      return true;
    } catch { if (sequence.current === request && current.current.scope === captured) setValue({ scope, candidate, name: candidate.label, preparing: false, failed: true }); return false; }
  };
  const save = async () => {
    if (current.current.disabled || !active?.outcome || active.outcome.status !== 'ready' || !active.name.trim() || !port) return;
    const request = ++sequence.current, captured = active;
    setValue({ ...captured, outcome: { ...captured.outcome!, status: 'pending' } });
    try {
      const outcome = await port.submit(captured.outcome!.session.id, captured.name);
      if (sequence.current !== request || current.current.scope !== captured.scope || current.current.port !== port) return;
      setValue({ ...captured, outcome: outcome ?? { ...captured.outcome!, status: 'blocked' } });
      return outcome?.status === 'confirmed' && outcome.currentName !== undefined ? { ...captured.candidate, label: outcome.currentName } : undefined;
    } catch { if (sequence.current === request && current.current.scope === captured.scope) setValue({ ...captured, failed: true }); }
  };
  const inspect = async () => {
    if (!active?.outcome || active.outcome.status !== 'unknown' || !port) return;
    const request = ++sequence.current, captured = active;
    try { const outcome = await port.inspect(active.outcome.session.id); if (outcome && sequence.current === request && current.current.scope === captured.scope) setValue({ ...captured, outcome }); } catch { /* Preserve unresolved write and the local name. */ }
  };
  const reset = () => { if (active?.outcome?.status === 'pending') return; sequence.current += 1; if (active?.outcome) port?.release(active.outcome.session.id); setValue(undefined); };
  return { value: active, begin, save, inspect, reset, change: (name: string) => { if (!current.current.disabled && active?.outcome?.status === 'ready') setValue({ ...active, name }); },
    pending: active?.preparing || active?.outcome?.status === 'pending', dirty: !!active && active.outcome?.status === 'ready' && active.name !== active.outcome.session.displayName };
}
