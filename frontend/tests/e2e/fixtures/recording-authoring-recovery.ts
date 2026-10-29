import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';
import type { WriteGroupTestWriteOperation, WriteGroupTestWriteOutcome, WriteGroupTestCleanupStatus,
  WriteGroupTestWriteStatus, WriteGroupTestWritePreview } from '../../../src/types/studioV2WriteGroupTestWrite';
import type { RecordingTestWriteRecoveryOperation } from '../../../src/types/recordingTestWriteRecovery';
import { isRecord } from '../../../src/utils/studioV2LifecycleJsonPrimitives';
import { MockContractError } from './recording-lifecycle-successor';
import { AUTHORING_TIME, AUTHORING_HASH } from './recording-authoring-state';

/** Test-only durable operation oracle; it never inserts a row or queries a database. */
export class AuthoringRecoveryMock {
  constructor(private readonly dialect: string) {}
  private preview?: WriteGroupTestWritePreview;
  private operation?: RecordingTestWriteRecoveryOperation;
  private outcome: Pick<WriteGroupTestWriteOperation, 'status' | 'write_outcome' | 'cleanup_status'> = {
    status: 'unknown', write_outcome: 'unknown', cleanup_status: 'unknown',
  };
  insertCount = 0;
  repairCount = 0;
  previewCount = 0;
  setOutcome(status: WriteGroupTestWriteStatus, write_outcome?: WriteGroupTestWriteOutcome, cleanup_status?: WriteGroupTestCleanupStatus) {
    this.outcome = { status, ...(write_outcome ? { write_outcome } : {}), ...(cleanup_status ? { cleanup_status } : {}) };
    if (this.operation) this.operation = { ...this.identity(this.operation), ...this.result() };
  }
  private identity(operation: RecordingTestWriteRecoveryOperation) {
    const { workspace_id, workspace_revision, group_id, group_revision, connector_id, connector_revision, target_fingerprint } = operation;
    return { workspace_id, workspace_revision, group_id, group_revision, connector_id, connector_revision, target_fingerprint };
  }
  private result(): WriteGroupTestWriteOperation {
    return { operation_id: 'O1', action: 'test_write', ...this.outcome, created_at: AUTHORING_TIME, updated_at: AUTHORING_TIME,
      ...(!['pending', 'running'].includes(this.outcome.status) ? { completed_at: AUTHORING_TIME } : {}) };
  }
  unresolved() {
    return this.operation && !(this.operation.cleanup_status === 'cleaned' && ['succeeded', 'partial'].includes(this.operation.status))
      && !(this.operation.status === 'failed' && this.operation.cleanup_status === 'not_attempted');
  }
  list() { return { workspace_id: 'ws-1', complete: true, operations: this.unresolved() ? [structuredClone(this.operation!)] : [] }; }
  status(id: string) {
    if (id !== 'O1' || !this.operation) throw new MockContractError('MOCK_OPERATION_NOT_FOUND', 404);
    return structuredClone(this.operation);
  }
  snapshot() { return { operations: this.operation ? [structuredClone(this.operation)] : [], insertCount: this.insertCount,
    repairCount: this.repairCount, testWritePreviewCount: this.previewCount }; }
  makePreview(group: WriteGroup): WriteGroupTestWritePreview {
    if (this.unresolved()) throw new MockContractError('TEST_WRITE_UNRESOLVED');
    if (this.previewCount) throw new MockContractError('MOCK_REPLACEMENT_PREVIEW_FORBIDDEN');
    this.previewCount++;
    this.preview = { token: 'MOCK-TEST-WRITE-TOKEN', operation_id: 'O1', action: 'test_write', group_id: group.id,
      group_revision: group.revision, expires_at: '2099-01-01T00:00:00Z', target: { connector_id: group.destination.connector_id,
        dialect: this.dialect, database: 'mock-only', schema: 'main', table: group.destination.table_name },
      owner_column: '__mock_owner', owner_value: 'O1', dedupe: 'mock-effect-O1',
      values: [{ column: '__mock_owner', type: 'text', value: 'O1' }, { column: 'value', type: 'int64', value: '2' }],
      cleanup: 'Mock operation only; no SQL exists. Cleanup is restricted to O1.' };
    return structuredClone(this.preview);
  }
  confirm(group: WriteGroup, workspaceRevision: string, body: unknown) {
    if (!isRecord(body) || Object.keys(body).some(key => !['token', 'operation_id', 'action'].includes(key)) ||
      !this.preview || body.token !== this.preview.token || body.operation_id !== 'O1' || group.id !== this.preview.group_id ||
      group.revision !== this.preview.group_revision) throw new MockContractError('MOCK_TEST_WRITE_SCOPE_MISMATCH');
    if (body.action === 'cleanup_repair') {
      if (this.dialect !== 'sqlserver' || !this.operation || this.operation.write_outcome !== 'written_verified' ||
        !['failed', 'unknown'].includes(this.operation.cleanup_status ?? '')) throw new MockContractError('MOCK_CLEANUP_NOT_ALLOWED');
      this.repairCount++;
      this.setOutcome('succeeded', 'written_verified', 'cleaned');
    } else {
      if (body.action !== undefined || this.operation) throw new MockContractError('MOCK_DUPLICATE_INSERT_FORBIDDEN');
      this.insertCount++;
      this.operation = { ...this.result(), workspace_id: 'ws-1', workspace_revision: workspaceRevision,
        group_id: group.id, group_revision: group.revision, connector_id: group.destination.connector_id,
        connector_revision: group.destination.connector_revision, target_fingerprint: AUTHORING_HASH };
    }
    return structuredClone(this.operation);
  }
}
