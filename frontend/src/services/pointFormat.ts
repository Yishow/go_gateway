import { api } from './datalinkClient';
import type { PointFormatIntent, PointFormatOperation, PointFormatOrdinaryPatch, PointFormatScope, PointFormatSnapshot } from '../types/pointFormat';
import { PointFormatError, record, safePointFormatError } from './pointFormatErrors';
import { parseMutation, parseSnapshot, responseData, text, wordFormats } from './pointFormatParser';
function invalid(): never { throw new PointFormatError('point_format_invalid', 'rejected'); }
function scopeParams(scope: PointFormatScope) {
  text(scope.point_id, false);
  return { workspace_id: text(scope.workspace_id, false), device_id: text(scope.device_id, false), ...(scope.rule_id ? { rule_id: text(scope.rule_id, false) } : {}) };
}
function copyPatch(patch?: PointFormatOrdinaryPatch): Readonly<PointFormatOrdinaryPatch> | undefined {
  if (patch === undefined) return undefined;
  const source = record(patch) ?? invalid();
  const copy: Record<string, string | boolean> = {};
  for (const [key, value] of Object.entries(source)) {
    if (!['name', 'description', 'address', 'data_type', 'enabled', 'polling_group_id'].includes(key)) return invalid();
    if (key === 'enabled') { if (typeof value !== 'boolean') return invalid(); copy[key] = value; }
    else { if (typeof value !== 'string') return invalid(); copy[key] = value; }
  }
  return Object.freeze(copy);
}
/** Capture once when editing starts. Save/inspect never refresh or rebase this evidence. */
export function createPointFormatOperation(initial: PointFormatSnapshot, intent: PointFormatIntent, requestId: string, patch?: PointFormatOrdinaryPatch): PointFormatOperation {
  const snapshot = parseSnapshot(initial, initial);
  if (!snapshot.editable || !snapshot.resolution.resolved || !/^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$/.test(requestId)) return invalid();
  if (intent.format_mode === 'override') {
    if (!snapshot.resolution.applicable || !wordFormats.includes(intent.configured_format)) return invalid();
  } else if (intent.format_mode !== 'inherit' || intent.configured_format !== '') return invalid();
  const copiedPatch = copyPatch(patch);
  return Object.freeze({ workspace_id: snapshot.workspace_id, device_id: snapshot.device_id, point_id: snapshot.point_id,
    ...(snapshot.rule_id ? { rule_id: snapshot.rule_id, expected_rule_revision: snapshot.rule_revision } : {}),
    expected_workspace_revision: snapshot.workspace_revision, expected_point_revision: snapshot.point_revision,
    request_id: requestId, format_mode: intent.format_mode, configured_format: intent.configured_format,
    ...(copiedPatch ? { patch: copiedPatch } : {}),
  } as PointFormatOperation);
}
export const pointFormatAPI = {
  async read(scope: PointFormatScope): Promise<PointFormatSnapshot> {
    try {
      const response = await api.get(`/points/${encodeURIComponent(scope.point_id)}/format`, { params: scopeParams(scope) });
      return parseSnapshot(responseData(response.data), scope);
    } catch (error) { throw safePointFormatError(error); }
  },
  async save(operation: PointFormatOperation) {
    try {
      const response = await api.put(`/points/${encodeURIComponent(operation.point_id)}/format`, operation);
      return parseMutation(responseData(response.data), operation);
    } catch (error) { throw safePointFormatError(error, operation.request_id); }
  },
  /** Explicit recovery only: a missing receipt does not authorize automatic resubmission. */
  async inspect(operation: PointFormatOperation) {
    try {
      const response = await api.get(`/points/${encodeURIComponent(operation.point_id)}/format/operations/${encodeURIComponent(operation.request_id)}`, { params: scopeParams(operation) });
      return parseMutation(responseData(response.data), operation);
    } catch (error) { throw safePointFormatError(error, operation.request_id); }
  },
};
