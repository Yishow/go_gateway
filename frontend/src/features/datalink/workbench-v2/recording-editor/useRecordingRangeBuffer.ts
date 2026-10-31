import { useEffect, useMemo, useRef, useState } from 'react';
import type { RangeSourceBuffer, RangeSourceOutcome, RangeSourcePort, RangeSourceRequest } from '../../../../pages/datalink/workbench-v2/rangeSourceSubmission';
import type { StudioV2WorkspaceDeviceRecord } from '../../../../services/studioV2WorkspaceDevices';
import { newRecordingRangeBuffer, parseRecordingRangeBuffer, type RecordingRangeBuffer } from './recordingRangeBuffer';
const raw = (buffer: RangeSourceBuffer): RecordingRangeBuffer => ({ name: buffer.name, address: buffer.start_address, count: String(buffer.count), dataType: buffer.data_type, scale: String(buffer.scale_multiplier), offset: String(buffer.scale_offset) });
export function useRecordingRangeBuffer(workspaceId: string, port?: RangeSourcePort, disabled = false) {
  const [buffer, setBuffer] = useState(newRecordingRangeBuffer), [dirty, setDirty] = useState(false);
  const [target, setTarget] = useState<{ id: string; name: string }>(), [device, setDevice] = useState<StudioV2WorkspaceDeviceRecord>();
  const [activationAllowed, setActivationAllowed] = useState(false);
  const [outcome, setOutcome] = useState<RangeSourceOutcome>(), [pending, setPending] = useState(false), [checking, setChecking] = useState(false), [preparing, setPreparing] = useState(false), [blocked, setBlocked] = useState(false);
  const owner = useMemo(() => ({ workspaceId }), [workspaceId]), ownerRef = useRef(owner); ownerRef.current = owner;
  const generation = useRef(0), session = useRef<string>(crypto.randomUUID()), requestRef = useRef<RangeSourceRequest | undefined>(undefined);
  const submitting = useRef(false), prepareGeneration = useRef(0), mounted = useRef(true);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; prepareGeneration.current += 1; }; }, []);
  const locked = pending || outcome?.status === 'pending' || outcome?.status === 'unknown' || outcome?.status === 'persisted-unavailable';
  const unresolved = port?.observations().find(item => item.request.workspaceId === workspaceId && ['pending', 'unknown', 'persisted-unavailable'].includes(item.status));
  const inactiveSources = port?.observations().filter(item => item.request.workspaceId === workspaceId && item.status === 'saved-inactive') ?? [];
  const matching = inactiveSources.filter(item => item.request.deviceId === target?.id);
  const inactive = inactiveSources.find(item => item.request.requestId === requestRef.current?.requestId) ?? (matching.length === 1 ? matching[0] : inactiveSources.length === 1 ? inactiveSources[0] : undefined);
  const recoverable = Boolean(locked || unresolved || inactiveSources.length);
  const unresolvedDeviceId = locked ? requestRef.current?.deviceId : unresolved?.request.deviceId;
  const canStart = (deviceId: string) => !(locked || unresolved) || unresolvedDeviceId === deviceId;
  const parsed = target && device ? parseRecordingRangeBuffer(buffer, target.id, device.protocol) : undefined;
  const change = (next: RecordingRangeBuffer) => { setBuffer(next); setDirty(true); generation.current += 1; };
  const reset = () => { if (locked) return; setDirty(false); setBuffer(newRecordingRangeBuffer()); setOutcome(undefined); requestRef.current = undefined; generation.current = 0; session.current = crypto.randomUUID(); };
  const recovery = (requestId?: string) => requestId ? inactiveSources.find(item => item.request.requestId === requestId) : unresolved ?? inactive;
  const canResume = (requestId?: string) => {
    if (locked || unresolved) return !requestId || requestId === (locked ? requestRef.current?.requestId : unresolved?.request.requestId);
    const previous = recovery(requestId);
    return Boolean(previous && (!dirty || previous.request.requestId === requestRef.current?.requestId));
  };
  const resume = (requestId?: string) => {
    if (!canResume(requestId)) return false;
    if (locked) return true;
    const previous = recovery(requestId);
    if (!previous) return false;
    if (outcome?.request.requestId === previous.request.requestId) { setOutcome(previous); return true; }
    prepareGeneration.current += 1; setDevice(undefined); setPreparing(false); setBlocked(false);
    setTarget({ id: previous.request.deviceId, name: previous.request.deviceId }); setBuffer(raw(previous.request.buffer)); setDirty(false);
    setOutcome(previous); requestRef.current = previous.request; generation.current = previous.request.generation; session.current = previous.request.sessionId; return true;
  };
  const start = async (next: { id: string; name: string }) => {
    if (disabled || !canStart(next.id)) return false;
    if ((locked || unresolved) && resume()) return true;
    reset(); setTarget(next); setDevice(undefined); setBlocked(false);
    if (!port) return true; // Local buffer exists, but no submit capability is claimed.
    const captured = owner, intent = ++prepareGeneration.current; setPreparing(true);
    try {
      const record = await port.prepare(next.id);
      if (!mounted.current || ownerRef.current !== captured || prepareGeneration.current !== intent) return false;
      setDevice(record?.device); setActivationAllowed(record?.activationAllowed === true); setBlocked(!record); return true;
    } catch { if (mounted.current && ownerRef.current === captured && prepareGeneration.current === intent) setBlocked(true); return false; }
    finally { if (mounted.current && ownerRef.current === captured && prepareGeneration.current === intent) setPreparing(false); }
  };
  const accept = (result: RangeSourceOutcome | undefined, captured: typeof owner) => {
    const request = requestRef.current;
    if (!mounted.current || ownerRef.current !== captured || !request || !result || result.request.workspaceId !== request.workspaceId || result.request.requestId !== request.requestId ||
      result.request.sessionId !== request.sessionId || result.request.generation !== request.generation) return;
    const verified = result.status === 'confirmed' ? Boolean(result.rule?.id && result.candidates?.candidates.length) : result.status === 'saved-inactive' ? Boolean(result.rule?.id && result.savedPoints?.length) : true;
    setOutcome(verified ? result : { ...result, status: 'unknown' }); submitting.current = false; setPending(false);
    if (verified && ['confirmed', 'saved-inactive'].includes(result.status) && generation.current === request.generation) setDirty(false);
  };
  const submit = async () => {
    if (!port || disabled || locked || !parsed || !dirty || submitting.current || ['confirmed', 'saved-inactive'].includes(outcome?.status ?? '')) return;
    const rule = { ...parsed.rule, enabled: activationAllowed };
    const request: RangeSourceRequest = { workspaceId, sessionId: session.current, requestId: crypto.randomUUID(), generation: generation.current, deviceId: rule.device_id,
      buffer: { name: rule.name, naming_prefix: rule.naming_prefix, start_address: rule.start_address, count: rule.count, data_type: rule.data_type, enabled: rule.enabled,
        scale_multiplier: rule.scale_multiplier, scale_offset: rule.scale_offset, data_format: rule.data_format, skipped_addresses: rule.skipped_addresses, share_enabled: rule.share_enabled, share_start_register: rule.share_start_register, share_stride: rule.share_stride } };
    requestRef.current = request; submitting.current = true; setPending(true); setOutcome(undefined); const captured = owner;
    try { accept(await port.submit(request), captured); } catch { accept({ request, status: 'unknown' }, captured); }
  };
  const prepareSaved = async () => {
    const source = outcome, original = requestRef.current;
    if (!port || disabled || pending || submitting.current || dirty || source?.status !== 'saved-inactive' || !original) return;
    const request = { ...original, requestId: crypto.randomUUID(), generation: generation.current, buffer: { ...original.buffer, enabled: true } };
    requestRef.current = request; submitting.current = true; setPending(true); const captured = owner;
    try { accept(await port.prepareSaved(original.requestId, request), captured); } catch { accept({ request, status: 'unknown' }, captured); }
  };
  const inspect = async () => {
    const request = requestRef.current; if (!port || !request || checking) return;
    setChecking(true); const captured = owner;
    try { accept(await port.inspect(request.requestId), captured); } catch { accept({ request, status: 'unknown' }, captured); }
    finally { if (mounted.current && ownerRef.current === captured) setChecking(false); }
  };
  return { buffer, dirty, target, device, outcome, pending, checking, preparing, blocked, activationAllowed, locked, recoverable, inactiveSources, canResume, canStart, parsed, change, reset, resume, start, submit, inspect, prepareSaved, submission: () => requestRef.current,
    canPrepareSaved: Boolean(port && !disabled && !pending && !dirty && outcome?.status === 'saved-inactive'),
    canSubmit: Boolean(port && !disabled && !locked && parsed && dirty && !['confirmed', 'saved-inactive'].includes(outcome?.status ?? '')) };
}
export type RecordingRangeController = ReturnType<typeof useRecordingRangeBuffer>;
