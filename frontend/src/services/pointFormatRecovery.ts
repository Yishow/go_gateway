import { api } from './datalinkClient';
import type { PointFormatMutationResult, PointFormatOperation, PointFormatSnapshot } from '../types/pointFormat';
import type { PointFormatRecoveryRequest, PointFormatRecoveryResult } from '../types/pointFormatRecovery';
import { PointFormatError, record, safePointFormatError } from './pointFormatErrors';
import { parseMutation, parseSnapshot, responseData, text } from './pointFormatParser';
const invalid = (): never => { throw new PointFormatError('point_format_invalid', 'rejected'); };
const unknown = (): never => { throw new PointFormatError('point_format_unknown', 'unknown'); };
const keys = ['workspace_id', 'device_id', 'point_id', 'rule_id', 'request_id', 'expected_workspace_revision', 'expected_point_revision', 'expected_rule_revision', 'expected_connection_revision', 'expected_settings_revision'];
function validate(request: PointFormatRecoveryRequest): void {
  if (Object.keys(request).some(key => !keys.includes(key))) invalid();
  for (const key of keys.filter(key => !['rule_id', 'expected_rule_revision'].includes(key))) text(request[key as keyof PointFormatRecoveryRequest], false);
  if (!/^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$/.test(request.request_id)) invalid();
  if (request.rule_id !== undefined) { text(request.rule_id, false); text(request.expected_rule_revision, false); }
  else if (request.expected_rule_revision !== undefined) invalid();
}
export function createPointFormatRecovery(initial: PointFormatSnapshot, operation: PointFormatOperation, saved: PointFormatMutationResult, settingsRevision: string): PointFormatRecoveryRequest {
  const snapshot = parseSnapshot(initial, operation), result = parseMutation(saved, operation);
  if (operation.expected_workspace_revision !== snapshot.workspace_revision || operation.expected_point_revision !== snapshot.point_revision ||
      operation.expected_rule_revision !== (snapshot.rule_id ? snapshot.rule_revision : undefined) || result.resolution.connection_revision !== snapshot.connection_revision) invalid();
  const request = { workspace_id: snapshot.workspace_id, device_id: snapshot.device_id, point_id: snapshot.point_id,
    ...(snapshot.rule_id ? { rule_id: snapshot.rule_id, expected_rule_revision: snapshot.rule_revision } : {}), request_id: operation.request_id,
    expected_workspace_revision: result.workspace_revision, expected_point_revision: result.point.revision_id,
    expected_connection_revision: snapshot.connection_revision, expected_settings_revision: settingsRevision };
  validate(request);
  return Object.freeze(request);
}
export const pointFormatRecoveryAPI = {
  async recover(request: PointFormatRecoveryRequest): Promise<PointFormatRecoveryResult> {
    try {
      validate(request);
      const response = await api.post('/modbus-share/reconcile', { format_recovery: request });
      const data = responseData(response.data), hydration = record(data.hydration) ?? unknown();
      if (!['recovered', 'replayed'].includes(String(data.outcome)) || !['installed', 'not_started'].includes(String(data.runtime_state)) || typeof data.share_applied !== 'boolean' ||
          hydration.state !== 'ready' || hydration.readiness !== true || hydration.workspace_id !== request.workspace_id || hydration.workspace_revision !== request.expected_workspace_revision || hydration.settings_revision !== request.expected_settings_revision) unknown();
      return { outcome: data.outcome as PointFormatRecoveryResult['outcome'], runtime_state: data.runtime_state as PointFormatRecoveryResult['runtime_state'], share_applied: data.share_applied as boolean,
        hydration: { state: 'ready', readiness: true, workspace_id: request.workspace_id, workspace_revision: request.expected_workspace_revision, settings_revision: request.expected_settings_revision } };
    } catch (error) { throw safePointFormatError(error, request.request_id); }
  },
};
