import { normalizeTypedEnvelope } from '../../../../../../../utils/safeJson';
import type { WriteGroup, WriteGroupMutationResponse } from '../../../../../../../types/studioV2WriteGroup';
import type { WriteGroupLifecycle, SavedRecordingActionIntent } from '../../../../../../../types/studioV2WriteGroupLifecycle';
import type { studioV2WriteGroupSuccessorAPI } from '../../../../../../../services/studioV2WriteGroupSuccessor';
import { sameLifecycleJSON } from '../../../../../../../services/studioV2LifecycleTransport';
import type * as S from '../../../../../../../types/studioV2WriteGroupSuccessor';
import { activationAction, currentStage, explicitRejection, isFutureBoundary, matchesIntent, operationScope, reviewRequest, sameBoundary, matchesRootStatus, canRepreviewAfterStatus } from './recordingChangeHelpers';
import { loadRecordingChanges, loadUnreadableRecordingScopes, saveRecordingChange } from './recordingChangeRecovery';
import type { RecordingChangeRecovery, RecordingChangeState } from './recordingChangeTypes';

export interface RecordingChangeDependencies {
  api: typeof studioV2WriteGroupSuccessorAPI;
  lifecycle: (groupId: string) => Promise<WriteGroupLifecycle>;
  getGroup: (groupId: string) => Promise<WriteGroupMutationResponse>;
  onChanged?: (group: WriteGroup, workspaceRevision: string) => void;
}
const initial = (): RecordingChangeState => ({
  isOpen: false, busy: false, phase: 'choose', boundary: '', unknown: false, canRetry: false, canRepreview: false,
  intent: undefined, lifecycle: undefined, target: undefined, review: undefined, recovery: undefined,
  preview: undefined, previewStage: undefined, activation: undefined, activationReview: undefined,
  unreadableScope: undefined, readback: undefined, error: undefined,
});
/** One workspace-scoped controller survives Basic/Advanced changes. Async completions never target a new selection. */
export class RecordingChangeController {
  state: RecordingChangeState;
  readonly workspaceId: string;
  readonly listeners = new Set<() => void>();
  private alive = true;
  private epoch = 0;
  private locked = false;
  readonly = false;
  deps: RecordingChangeDependencies;
  constructor(workspaceId: string, deps: RecordingChangeDependencies) {
    this.workspaceId = workspaceId; this.deps = deps;
    const records = loadRecordingChanges(workspaceId);
    const recovery = records.find((item) => item.phase !== 'scheduled' && item.outcome !== 'rejected') ?? records.at(-1);
    this.state = recovery ? { ...initial(), recovery, phase: recovery.phase, unknown: recovery.phase !== 'scheduled',
      preview: recovery.preview ? { ...recovery.preview, statements: [] } : undefined, previewStage: recovery.previewStage } : initial();
    const unreadableScope = loadUnreadableRecordingScopes(workspaceId)[0];
    if (unreadableScope) this.state = { ...initial(), unreadableScope, unknown: true, error: 'unreadable' };
  }
  subscribe = (listener: () => void) => { this.listeners.add(listener); return () => { this.listeners.delete(listener); }; };
  snapshot = () => this.state;
  activate = () => { this.alive = true; };
  dispose = () => { this.alive = false; this.epoch++; };
  private set(patch: Partial<RecordingChangeState>, epoch = this.epoch) {
    if (!this.alive || epoch !== this.epoch) return;
    this.state = { ...this.state, ...patch }; this.listeners.forEach((listener) => listener());
  }
  private async run(action: (epoch: number) => Promise<void>) {
    if (this.locked) return;
    this.locked = true; const epoch = this.epoch; this.set({ busy: true });
    try { await action(epoch); } catch { this.set({ error: 'failed' }, epoch); }
    finally { this.locked = false; this.set({ busy: false }); }
  }
  get unresolved() { return Boolean(this.state.unreadableScope || this.state.recovery && this.state.phase !== 'scheduled' && this.state.recovery.outcome !== 'rejected'); }
  dismiss = () => { this.epoch++; this.set({ isOpen: false }); };
  resume = () => {
    const unreadableScope = loadUnreadableRecordingScopes(this.workspaceId)[0];
    if (unreadableScope) { this.set({ ...initial(), isOpen: true, unreadableScope, unknown: true, error: 'unreadable' }); return; }
    const records = loadRecordingChanges(this.workspaceId);
    const record = records.find((item) => item.phase !== 'scheduled' && item.outcome !== 'rejected') ?? records.at(-1);
    if (record) this.set({ ...initial(), isOpen: true, recovery: record, phase: record.phase, unknown: record.phase !== 'scheduled',
      preview: record.preview ? { ...record.preview, statements: [] } : undefined, previewStage: record.previewStage });
  };
  requestAction = (intent: SavedRecordingActionIntent) => this.run(async (epoch) => {
    if (intent.workspaceId !== this.workspaceId || !['recreate', 'transition'].includes(intent.action)) return;
    if (this.unresolved) { this.resume(); return; }
    this.set({ ...initial(), busy: true, isOpen: true, intent }, epoch);
    if (this.readonly) { this.set({ error: 'unavailable' }, epoch); return; }
    const lifecycle = await this.deps.lifecycle(intent.groupId);
    if (!matchesIntent(lifecycle, intent)) { this.set({ error: 'stale' }, epoch); return; }
    if (!lifecycle.actions.some((item) => item.action === intent.action && item.allowed) || !lifecycle.successor_seed) {
      this.set({ error: 'unavailable' }, epoch); return;
    }
    this.set({ lifecycle, target: structuredClone(lifecycle.successor_seed.group) }, epoch);
  });
  setTarget = (target: WriteGroup) => {
    if (this.locked || this.readonly || this.state.recovery) return;
    this.set({ target, review: undefined, phase: 'choose' });
  };
  setBoundary = (boundary: string) => {
    if (this.locked || this.readonly || this.state.unknown || this.state.phase === 'scheduled') return;
    this.set({ boundary, activation: undefined, activationReview: undefined, phase: 'schema-confirmed' });
  };
  back = () => { if (!this.locked && !this.state.recovery) this.set({ phase: 'choose', review: undefined, error: undefined }); };
  review = () => this.run(async (epoch) => {
    const { lifecycle, target, intent } = this.state;
    if (!lifecycle || !target || !intent || this.readonly || this.unresolved) return;
    const fresh = await this.deps.lifecycle(intent.groupId);
    if (!matchesIntent(fresh, intent)) { this.set({ error: 'stale' }, epoch); return; }
    const request = reviewRequest(lifecycle, target, crypto.randomUUID());
    const review = await this.deps.api.review(intent.groupId, request);
    this.set({ review, phase: 'impact' }, epoch);
  });
  private persist(record: RecordingChangeRecovery, epoch: number): boolean {
    if (!saveRecordingChange(record)) { this.set({ recovery: record, unknown: true, error: 'storage' }, epoch); return false; }
    this.set({ recovery: record }, epoch); return true;
  }
  private async stage(record: RecordingChangeRecovery) {
    if (!record.preparation) throw new Error('No reservation');
    const [source, target] = await Promise.all([this.deps.lifecycle(record.envelope.predecessor_group_id), this.deps.lifecycle(record.preparation.group.id)]);
    return currentStage(record, source, target);
  }
  private async notify(record: RecordingChangeRecovery, epoch: number) {
    if (!record.preparation) return;
    try {
      const result = await this.deps.getGroup(record.preparation.group.id);
      if (this.alive && epoch === this.epoch) this.deps.onChanged?.(result.group, result.workspace_revision);
    } catch { /* A metadata refresh failure cannot undo a recorded operation. */ }
  }
  private async submit(record: RecordingChangeRecovery, epoch: number): Promise<RecordingChangeRecovery | undefined> {
    if (record.unreadable) return;
    const pending = { ...record, attempt: (record.attempt ?? 0) + 1, outcome: 'pending' as const, stageRejected: false,
      schemaConfirmationSubmitted: record.schemaConfirmationSubmitted || record.submitted.action === 'schemaConfirm' };
    if (!this.persist(pending, epoch)) return;
    const id = record.envelope.predecessor_group_id;
    try {
      let next = { ...pending, outcome: 'recorded' as const } as RecordingChangeRecovery;
      const submitted = record.submitted;
      if (submitted.action === 'reserve') {
        const preparation = await this.deps.api.reserve(id, submitted.request);
        next = { ...next, preparation, phase: 'reserved' };
      } else if (submitted.action === 'schemaPreview') {
        const preview = await this.deps.api.schemaPreview(id, submitted.request);
        const { statements: _sql, database: _database, ...safePreview } = preview;
        next = { ...next, phase: 'schema-plan', preview: safePreview, previewStage: submitted.request };
        this.set({ preview, previewStage: submitted.request }, epoch);
      } else if (submitted.action === 'schemaConfirm') {
        const operation = await this.deps.api.schemaConfirm(id, submitted.request);
        if (operation.status !== 'succeeded') { this.persist({ ...pending, outcome: 'unknown' }, epoch); this.set({ unknown: true, error: 'unknown' }, epoch); return; }
        next = { ...next, phase: 'schema-confirmed' };
      } else {
        const result = await this.deps.api[submitted.action](id, submitted.request);
        next = { ...next, result, phase: 'scheduled' };
      }
      this.persist(next, epoch); this.set({ phase: next.phase, unknown: false, canRetry: false, canRepreview: false, error: undefined }, epoch);
      if (submitted.action !== 'schemaPreview') await this.notify(next, epoch);
      return next;
    } catch (cause) {
      const outcome = explicitRejection(cause) ? 'rejected' : 'unknown';
      // A stage rejection cannot erase a committed reservation or release its lineage fence.
      this.persist({ ...pending, outcome: pending.preparation ? 'unknown' : outcome, stageRejected: outcome === 'rejected' }, epoch);
      this.set({ unknown: Boolean(pending.preparation) || outcome === 'unknown', canRetry: false, error: outcome === 'unknown' ? 'unknown' : 'stale' }, epoch);
    }
  }
  reserveAndPreview = () => this.run(async (epoch) => {
    const { review } = this.state;
    if (this.readonly || !review || review.blockers.length || this.state.unknown || this.state.recovery) return;
    const record: RecordingChangeRecovery = { version: 1, envelope: review.envelope, reviewDigest: review.review_digest,
      submitted: { action: 'reserve', request: { envelope: review.envelope, review_digest: review.review_digest } }, phase: 'impact', outcome: 'pending' };
    const reserved = await this.submit(record, epoch);
    if (!reserved || !this.alive || epoch !== this.epoch) return;
    await this.preview(reserved, epoch);
  });
  private async preview(record: RecordingChangeRecovery, epoch: number) {
    const request = await this.stage(record);
    if (!this.alive || epoch !== this.epoch) return;
    await this.submit({ ...record, submitted: { action: 'schemaPreview', request } }, epoch);
  }
  prepareSchema = () => this.run(async (epoch) => {
    const record = this.state.recovery;
    if (!this.readonly && record && this.state.phase === 'reserved' && !this.state.unknown) await this.preview(record, epoch);
  });
  confirmSchema = () => this.run(async (epoch) => {
    const { preview, previewStage, recovery } = this.state;
    if (this.readonly || this.state.unknown || !preview || !previewStage || !recovery || this.state.phase !== 'schema-plan') return;
    const fresh = await this.stage(recovery);
    if (!sameLifecycleJSON(fresh.current, previewStage.current)) { this.set({ error: 'stale', unknown: true }, epoch); return; }
    const request = { ...previewStage, token: preview.token, digest: preview.digest, schema_operation_id: preview.operation_id };
    await this.submit({ ...recovery, submitted: { action: 'schemaConfirm', request } }, epoch);
  });
  reviewActivation = () => this.run(async (epoch) => {
    const { recovery, boundary } = this.state;
    if (this.readonly || !recovery || this.state.unknown || this.state.phase !== 'schema-confirmed') return;
    if (!isFutureBoundary(boundary)) { this.set({ error: 'boundary' }, epoch); return; }
    const stage = await this.stage(recovery);
    const review = await this.deps.api.activationReview(recovery.envelope.predecessor_group_id, { stage, boundary });
    this.set({ activationReview: review }, epoch);
    if (review.blockers.length) { this.set({ error: 'unavailable' }, epoch); return; }
    if (!this.alive || epoch !== this.epoch) return;
    const activation = await this.deps.api.activationReadiness(recovery.envelope.predecessor_group_id, { envelope: review.envelope, review_digest: review.review_digest });
    this.set({ activation, phase: activation.readiness.ready && !activation.review.blockers.length ? 'activation-ready' : 'schema-confirmed',
      error: activation.readiness.ready ? undefined : 'unavailable' }, epoch);
  });
  confirmActivation = () => this.run(async (epoch) => {
    const { recovery, activation, boundary } = this.state;
    if (this.readonly || this.state.unknown || !recovery || !activation?.readiness.ready || activation.review.blockers.length || this.state.phase !== 'activation-ready') return;
    if (!isFutureBoundary(boundary) || !sameBoundary(boundary, activation.review.envelope.request.boundary)) { this.set({ error: 'boundary' }, epoch); return; }
    const fresh = await this.stage(recovery);
    if (!sameLifecycleJSON(fresh.current, activation.review.envelope.request.stage.current)) { this.set({ activation: undefined, phase: 'schema-confirmed', error: 'stale' }, epoch); return; }
    const request = { envelope: activation.review.envelope, review_digest: activation.review.review_digest, readiness_digest: activation.readiness_digest };
    await this.submit({ ...recovery, submitted: { action: activationAction(request.envelope), request } }, epoch);
  });
  checkStatus = () => this.run(async (epoch) => {
    const record = this.state.recovery;
    if (!record) {
      const scope = this.state.unreadableScope;
      if (!scope) return;
      try {
        const readback = await this.deps.api.operationStatus(scope.groupId, { workspace_id: scope.workspaceId, lineage_id: scope.lineageId, operation_id: scope.operationId });
        this.set({ readback, error: 'unreadable', unknown: true, canRetry: false, canRepreview: false }, epoch);
      } catch { this.set({ error: 'status_failed' }, epoch); }
      return;
    }
    try {
      const status = await this.deps.api.operationStatus(record.envelope.predecessor_group_id, operationScope(record));
      await this.acceptStatus(record, status, epoch);
    } catch (cause) {
      const absent = explicitRejection(cause) && ['WRITE_GROUP_NOT_FOUND'].includes(normalizeTypedEnvelope(cause).code ?? '');
      if (record.phase !== 'scheduled') this.set({ unknown: true, canRepreview: false, canRetry: !record.unreadable && absent && record.submitted.action === 'reserve', error: 'status_failed' }, epoch);
    }
  });
  private async acceptStatus(record: RecordingChangeRecovery, status: S.WriteGroupSuccessorOperationStatus, epoch: number) {
    if (!matchesRootStatus(record, status)) throw new Error('Wrong operation');
    const terminal = status.first_start.status === 'succeeded' ? status.first_start.result : status.cutover.status === 'succeeded' ? status.cutover.result : undefined;
    const next: RecordingChangeRecovery = { ...record, preparation: status.reserve.result, ...(terminal ? { result: terminal, phase: 'scheduled', outcome: 'recorded', terminalStatus: status } : {}) };
    if (terminal || record.phase === 'scheduled') { this.persist(next, epoch); this.set({ phase: 'scheduled', unknown: false, canRetry: false, canRepreview: false }, epoch); return; }
    if (record.unreadable) { this.set({ unknown: true, canRetry: false, error: 'unknown' }, epoch); return; }
    let phase = record.phase; let unknown = true; let canRetry = false;
    if (record.submitted.action === 'reserve') { phase = 'reserved'; unknown = false; }
    if (record.submitted.action === 'schemaPreview' && record.preview && record.previewStage) {
      const matching = status.schema_operations.find((item) => item.operation_id === record.preview?.operation_id);
      if (matching?.status === 'preview_not_submitted' && !record.stageRejected) {
        const fresh = await this.stage(next);
        phase = sameLifecycleJSON(fresh.current, record.previewStage.current) ? 'schema-plan' : 'reserved'; unknown = false;
      }
    }
    if (record.submitted.action === 'schemaPreview' && record.stageRejected && status.schema_operations.length === 0) { phase = 'reserved'; unknown = false; }
    if (record.submitted.action === 'schemaConfirm') {
      const schemaId = record.submitted.request.schema_operation_id;
      const operation = status.schema_operations.find((item) => item.operation_id === schemaId);
      canRetry = operation?.status === 'preview_not_submitted' || operation?.status === 'succeeded';
      if (operation?.status === 'preview_not_submitted' && record.stageRejected) { phase = 'reserved'; unknown = false; canRetry = false; }
      if (operation?.status === 'succeeded') {
        try {
          const result = await this.deps.getGroup(status.reserve.result.group.id);
          if (result.group.destination.schema_revision === operation.operation_id) { phase = 'schema-confirmed'; unknown = false; canRetry = false; }
        } catch { /* Retain exact schema confirmation until proof can be checked. */ }
      }
    }
    if (record.submitted.action === 'confirmCutover') canRetry = status.cutover.status === 'not_submitted';
    if (record.submitted.action === 'confirmFirstStart') canRetry = status.first_start.status === 'not_submitted';
    const canRepreview = canRepreviewAfterStatus(record, status);
    this.persist({ ...next, phase, priorPreviewStates: { ...record.priorPreviewStates, ...Object.fromEntries(status.schema_operations.map(item => [item.operation_id, item.status])) }, priorPreviewIds: [...new Set([...(record.priorPreviewIds ?? []), ...status.schema_operations.map(item => item.operation_id)])], outcome: unknown ? 'unknown' : 'recorded' }, epoch);
    this.set({ phase, unknown, canRetry, canRepreview, ...(phase === 'reserved' ? { preview: undefined, previewStage: undefined } : {}), error: unknown ? 'unknown' : undefined }, epoch);
  }
  repreview = () => this.run(async (epoch) => {
    const record = this.state.recovery;
    if (this.readonly || !record || !this.state.canRepreview) return;
    this.set({ canRepreview: false }, epoch);
    const status = await this.deps.api.operationStatus(record.envelope.predecessor_group_id, operationScope(record));
    if (!canRepreviewAfterStatus(record, status)) { this.set({ unknown: true, error: 'unknown' }, epoch); return; }
    const retained = { ...record, preparation: status.reserve.result,
      priorPreviewStates: { ...record.priorPreviewStates, ...Object.fromEntries(status.schema_operations.map(item => [item.operation_id, item.status])) }, priorPreviewIds: [...new Set([...(record.priorPreviewIds ?? []), ...status.schema_operations.map(item => item.operation_id)])] };
    await this.preview(retained, epoch);
  });
  retryExact = () => this.run(async (epoch) => {
    if (this.readonly || !this.state.unknown || !this.state.canRetry || !this.state.recovery) return;
    await this.submit(this.state.recovery, epoch);
  });
}
