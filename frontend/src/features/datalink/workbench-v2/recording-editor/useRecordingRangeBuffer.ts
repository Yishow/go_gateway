import { useEffect, useMemo, useRef, useState } from 'react';
import type { RangeSourceBuffer, RangeSourceOutcome, RangeSourcePort, RangeSourceRequest } from '../../../../pages/datalink/workbench-v2/rangeSourceSubmission';
import type { StudioV2WorkspaceDeviceRecord } from '../../../../services/studioV2WorkspaceDevices';
import { newRecordingRangeBuffer, parseRecordingRangeBuffer, type RecordingRangeBuffer } from './recordingRangeBuffer';
const raw = (buffer: RangeSourceBuffer): RecordingRangeBuffer => ({ name: buffer.name, address: buffer.start_address, count: String(buffer.count), dataType: buffer.data_type, scale: String(buffer.scale_multiplier), offset: String(buffer.scale_offset) });
export function useRecordingRangeBuffer(workspaceId: string, port?: RangeSourcePort, disabled = false) {
  const [buffer, setBuffer] = useState(newRecordingRangeBuffer), [dirty, setDirty] = useState(false);
  const [target, setTarget] = useState<{ id: string; name: string }>(), [device, setDevice] = useState<StudioV2WorkspaceDeviceRecord>();
  const [outcome, setOutcome] = useState<RangeSourceOutcome>(), [pending, setPending] = useState(false), [checking, setChecking] = useState(false), [preparing, setPreparing] = useState(false), [blocked, setBlocked] = useState(false);
  const owner = useMemo(() => ({ workspaceId }), [workspaceId]), ownerRef = useRef(owner); ownerRef.current = owner;
  const generation = useRef(0), session = useRef<string>(crypto.randomUUID()), requestRef = useRef<RangeSourceRequest | undefined>(undefined);
  const submitting = useRef(false), prepareGeneration = useRef(0), mounted = useRef(true);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; prepareGeneration.current += 1; }; }, []);
  const locked = pending || outcome?.status === 'pending' || outcome?.status === 'unknown' || outcome?.status === 'persisted-unavailable';
  const unresolved = port?.observations().find(item => item.request.workspaceId === workspaceId && ['pending', 'unknown', 'persisted-unavailable'].includes(item.status));
  const recoverable = Boolean(locked || unresolved);
  const unresolvedDeviceId = locked ? requestRef.current?.deviceId : unresolved?.request.deviceId;
  const canStart = (deviceId: string) => !recoverable || unresolvedDeviceId === deviceId;
  const parsed = target && device ? parseRecordingRangeBuffer(buffer, target.id, device.protocol) : undefined;
  const change = (next: RecordingRangeBuffer) => { setBuffer(next); setDirty(true); generation.current += 1; };
  const reset = () => { if (locked) return; setDirty(false); setBuffer(newRecordingRangeBuffer()); setOutcome(undefined); requestRef.current = undefined; generation.current = 0; session.current = crypto.randomUUID(); };
  const resume = () => {
    if (locked) return true;
    if (!unresolved) return false;
    setTarget({ id: unresolved.request.deviceId, name: unresolved.request.deviceId }); setBuffer(raw(unresolved.request.buffer)); setDirty(false);
    setOutcome(unresolved); requestRef.current = unresolved.request; generation.current = unresolved.request.generation; session.current = unresolved.request.sessionId; return true;
  };
  const start = async (next: { id: string; name: string }) => {
    if (disabled || !canStart(next.id)) return false;
    if (resume()) return true;
    reset(); setTarget(next); setDevice(undefined); setBlocked(false);
    if (!port) return true; // Local buffer exists, but no submit capability is claimed.
    const captured = owner, intent = ++prepareGeneration.current; setPreparing(true);
    try {
      const record = await port.prepare(next.id);
      if (!mounted.current || ownerRef.current !== captured || prepareGeneration.current !== intent) return false;
      setDevice(record); setBlocked(!record); return true;
    } catch { if (mounted.current && ownerRef.current === captured && prepareGeneration.current === intent) setBlocked(true); return false; }
    finally { if (mounted.current && ownerRef.current === captured && prepareGeneration.current === intent) setPreparing(false); }
  };
  const accept = (result: RangeSourceOutcome | undefined, captured: typeof owner) => {
    const request = requestRef.current;
    if (!mounted.current || ownerRef.current !== captured || !request || !result || result.request.workspaceId !== request.workspaceId || result.request.requestId !== request.requestId ||
      result.request.sessionId !== request.sessionId || result.request.generation !== request.generation) return;
    const verified = result.status !== 'confirmed' || Boolean(result.rule?.id && result.candidates?.candidates.length);
    setOutcome(verified ? result : { ...result, status: 'unknown' }); submitting.current = false; setPending(false);
    if (verified && result.status === 'confirmed' && generation.current === request.generation) setDirty(false);
  };
  const submit = async () => {
    if (!port || disabled || locked || !parsed || !dirty || submitting.current || outcome?.status === 'confirmed') return;
    const rule = parsed.rule;
    const request: RangeSourceRequest = { workspaceId, sessionId: session.current, requestId: crypto.randomUUID(), generation: generation.current, deviceId: rule.device_id,
      buffer: { name: rule.name, naming_prefix: rule.naming_prefix, start_address: rule.start_address, count: rule.count, data_type: rule.data_type, enabled: rule.enabled,
        scale_multiplier: rule.scale_multiplier, scale_offset: rule.scale_offset, data_format: rule.data_format, skipped_addresses: rule.skipped_addresses, share_enabled: rule.share_enabled, share_start_register: rule.share_start_register, share_stride: rule.share_stride } };
    requestRef.current = request; submitting.current = true; setPending(true); setOutcome(undefined); const captured = owner;
    try { accept(await port.submit(request), captured); } catch { accept({ request, status: 'unknown' }, captured); }
  };
  const inspect = async () => {
    const request = requestRef.current; if (!port || !request || checking) return;
    setChecking(true); const captured = owner;
    try { accept(await port.inspect(request.requestId), captured); } catch { accept({ request, status: 'unknown' }, captured); }
    finally { if (mounted.current && ownerRef.current === captured) setChecking(false); }
  };
  return { buffer, dirty, target, device, outcome, pending, checking, preparing, blocked, locked, recoverable, canStart, parsed, change, reset, resume, start, submit, inspect, submission: () => requestRef.current,
    canSubmit: Boolean(port && !disabled && !locked && parsed && dirty && outcome?.status !== 'confirmed') };
}
export type RecordingRangeController = ReturnType<typeof useRecordingRangeBuffer>;
