import { useEffect, useMemo, useRef, useState } from 'react';
import type { DeviceSourceBuffer, DeviceSourceOutcome, DeviceSourcePort, DeviceSourceRequest } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { newDeviceSourceBuffer, parsedDeviceSourceBuffer } from './deviceSourceBuffer';
/** Ephemeral form state only. All submission and causal observations remain in
 * the existing workspace queue port. Lives above the modal's mount boundary. */
export function useRecordingDeviceBuffer(workspaceId: string, port?: DeviceSourcePort, disabled = false) {
  const [buffer, setBuffer] = useState(newDeviceSourceBuffer);
  const [generation, setGeneration] = useState(0), generationRef = useRef(0);
  const [dirty, setDirty] = useState(false);
  const [outcome, setOutcome] = useState<DeviceSourceOutcome>();
  const [pending, setPending] = useState(false), [checking, setChecking] = useState(false), [preparing, setPreparing] = useState(false);
  const [blocked, setBlocked] = useState(false), [deviceId, setDeviceId] = useState<string>();
  const [sessionId] = useState(() => crypto.randomUUID());
  const owner = useMemo(() => ({ workspaceId, sessionId }), [workspaceId, sessionId]);
  const ownerRef = useRef(owner); ownerRef.current = owner;
  const requestRef = useRef<DeviceSourceRequest | undefined>(undefined), submitting = useRef(false), prepareFlight = useRef(false);
  const mounted = useRef(true);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const locked = pending || outcome?.status === 'pending' || outcome?.status === 'unknown';
  const parsed = parsedDeviceSourceBuffer(buffer);
  const change = (next: DeviceSourceBuffer) => { setBuffer(next); generationRef.current += 1; setGeneration(generationRef.current); setDirty(true); };
  const reset = () => { if (locked) return; setBuffer(newDeviceSourceBuffer()); setGeneration(0); generationRef.current = 0; setDirty(false); setOutcome(undefined); setDeviceId(undefined); requestRef.current = undefined; };
  const accept = (result: DeviceSourceOutcome | undefined, captured: typeof owner) => {
    const request = requestRef.current;
    if (!mounted.current || ownerRef.current !== captured || !request || !result || result.request.workspaceId !== request.workspaceId ||
      result.request.sessionId !== request.sessionId || result.request.requestId !== request.requestId || result.request.generation !== request.generation) return;
    const verified = result.status !== 'confirmed' || Boolean(result.device?.id);
    setOutcome(verified ? result : { ...result, status: 'unknown' }); setPending(false); submitting.current = false;
    if (result.device?.id && (result.status === 'confirmed' || result.status === 'persisted-unavailable')) setDeviceId(result.device.id);
    if (verified && result.status === 'confirmed' && generationRef.current === request.generation) setDirty(false);
  };
  const unresolved = port?.observations().find(item => item.request.workspaceId === workspaceId && ['pending', 'unknown', 'persisted-unavailable'].includes(item.status));
  const localRecoverable = outcome?.request.workspaceId === workspaceId && (locked || outcome.status === 'persisted-unavailable');
  const resume = () => {
    if (localRecoverable) return true; // Keep any newer raw input in this owner.
    if (!unresolved) return false;
    requestRef.current = unresolved.request; setBuffer(unresolved.request.buffer); setGeneration(unresolved.request.generation);
    generationRef.current = unresolved.request.generation; setDirty(false); setOutcome(unresolved); setDeviceId(unresolved.device?.id); return true;
  };
  const start = async () => {
    if (!port || disabled || prepareFlight.current) return false;
    if (locked || dirty) return true;
    if (resume()) return true;
    prepareFlight.current = true; setPreparing(true); setBlocked(false); const captured = owner;
    let permitted: boolean;
    try { permitted = port.prepare ? await port.prepare(workspaceId) : true; } catch { permitted = false; }
    prepareFlight.current = false;
    if (!mounted.current || ownerRef.current !== captured) return false;
    setPreparing(false); setBlocked(!permitted);
    if (permitted && outcome?.status === 'confirmed') reset();
    return permitted;
  };
  const submit = async () => {
    if (!port || disabled || locked || !parsed || !dirty || submitting.current) return;
    const request: DeviceSourceRequest = { workspaceId, sessionId, requestId: crypto.randomUUID(), generation, ...(deviceId ? { deviceId } : {}), buffer: parsed };
    requestRef.current = request; submitting.current = true; setPending(true); setOutcome(undefined); const captured = owner;
    try { accept(await port.submit(request), captured); } catch { accept({ request, status: 'unknown' }, captured); }
  };
  const inspect = async () => {
    const request = requestRef.current; if (!port || !request || checking) return;
    setChecking(true); const captured = owner;
    try { accept(await port.inspect(request.requestId), captured); } catch { accept({ request, status: 'unknown' }, captured); }
    finally { if (mounted.current && ownerRef.current === captured) setChecking(false); }
  };
  return { recoverable: Boolean(localRecoverable || unresolved), resume, disabled, buffer, dirty, outcome, pending, checking, preparing, blocked, deviceId, locked, parsed, change, reset, start, submit, inspect };
}
export type RecordingDeviceBufferController = ReturnType<typeof useRecordingDeviceBuffer>;
