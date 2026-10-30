/** Captured request identity. This does not grant mutation permission. */
export interface StudioV2RecordingReadScope {
  readonly workspaceId: string;
  readonly workspaceRevision: string;
  readonly epoch: number;
  /** Separate Share/bootstrap authority; never a database setup revision. */
  readonly bootstrapRevision?: string;
  readonly groupRevision?: string;
  readonly appliedRevision?: string;
}
export function validRecordingReadScope(scope: StudioV2RecordingReadScope): boolean {
  return Boolean(scope.workspaceId && typeof scope.workspaceRevision === 'string' && Number.isSafeInteger(scope.epoch) && scope.epoch >= 0);
}
/** Keep the canonical prefix so existing invalidation reaches scoped observers too. */
export function recordingReadKey(key: readonly unknown[], scope?: StudioV2RecordingReadScope): readonly unknown[] {
  return scope ? [...key, 'scoped-read', scope.workspaceId, scope.workspaceRevision, scope.epoch,
    scope.groupRevision ?? '', scope.appliedRevision ?? '', scope.bootstrapRevision ?? null] : key;
}
