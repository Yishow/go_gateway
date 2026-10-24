import type { SchemaOperation } from '../../../../../../types/recordingPlan';
import { getSafeErrorMessage } from '../../../../../../utils/typedErrors';

type Translate = (key: string, options?: Record<string, unknown>) => string;

export function isTerminal(operation: SchemaOperation | null): boolean {
  return operation !== null && !['pending', 'running'].includes(operation.status);
}

export function scopeGroupId(scopeKey: string): string {
  return scopeKey.split('|', 1)[0] ?? '';
}

export function operationStatus(operation: SchemaOperation | null, t: Translate): string {
  return operation ? t(`step4.group.schema.operation_status.${operation.status}`) : t('step4.group.schema.operation_unconfirmed');
}

/** The typed code an operation reason stands for; permission_denied is the one reason that is not already a code. */
function operationCode(operation: SchemaOperation): string | undefined {
  return operation.reason === 'permission_denied' ? 'RECORDING_SCHEMA_PERMISSION_DENIED' : operation.reason;
}

export function operationDetail(operation: SchemaOperation, t: Translate): string | undefined {
  const safe = getSafeErrorMessage({ code: operationCode(operation), action: operation.next_action }, t);
  if (safe.code) return safe.message;
  if (safe.action) return safe.action;
  return operation.status === 'succeeded' ? undefined : t('step4.group.schema.operation_attention');
}

/** The stable backend code behind what the panel currently reports, whether a failed request or a failed operation. */
export function schemaIssueCode(error: unknown, operation: SchemaOperation | null, t: Translate): string | undefined {
  if (error !== null) return getSafeErrorMessage(error, t).code;
  return operation && operation.status !== 'succeeded' ? getSafeErrorMessage({ code: operationCode(operation) }, t).code : undefined;
}
