import type { WriteGroup, WriteGroupListResponse, WriteGroupMutationResponse, WriteGroupReadiness } from '../types/studioV2WriteGroup';
import type { StudioV2WorkspaceReadinessIssue, StudioV2WorkspaceReadinessSeverity, StudioV2WorkspaceReadinessStep } from '../types/studioV2WorkspaceReadiness';
import { MAX_SAFE_JSON_ARRAY_LENGTH, MAX_SAFE_JSON_DEPTH, normalizeTypedEnvelope, parseBoundedJson } from './safeJson';
import { isRecord, readOptionalText, readPresentText, readRequiredText } from './studioV2WriteGroupJsonPrimitives';
import { parseWriteGroupCalculationCapabilities } from './studioV2WriteGroupCalculationJson';
import { parseStudioV2WriteGroup } from './studioV2WriteGroupDefinitionJson';
export { parseStudioV2WriteGroup } from './studioV2WriteGroupDefinitionJson';
const READINESS_SEVERITIES = new Set<StudioV2WorkspaceReadinessSeverity>(['blocking', 'warning']);
const READINESS_STEPS = new Set<StudioV2WorkspaceReadinessStep>(['Step 1', 'Step 2', 'Step 3', 'Step 4']);
type ResponseOutcome = 'failed' | 'unconfirmed';

function parseWorkspaceRevision(value: unknown): string | null {
  if (typeof value !== 'string' || value.length > 256) return null;
  return value;
}

export function parseWriteGroupListData(value: unknown): WriteGroupListResponse | null {
  if (!isRecord(value)) return null;
  const workspaceId = readRequiredText(value, 'workspace_id');
  const workspaceRevision = parseWorkspaceRevision(value.workspace_revision);
  const capabilities = parseWriteGroupCalculationCapabilities(value.calculation_capabilities);
  if (capabilities === null) return null;
  const rawGroups = value.groups;
  if (!workspaceId || workspaceRevision === null || !Array.isArray(rawGroups) ||
    rawGroups.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const groups = rawGroups.map(parseStudioV2WriteGroup);
  const seenGroupIds = new Set<string>();
  for (const group of groups) {
    if (!group || group.workspace_id !== workspaceId || seenGroupIds.has(group.id)) return null;
    seenGroupIds.add(group.id);
  }
  return {
    workspace_id: workspaceId,
    workspace_revision: workspaceRevision,
    groups: groups as WriteGroup[],
    ...(capabilities !== undefined ? { calculation_capabilities: capabilities } : {}),
  };
}

export function parseWriteGroupMutationData(value: unknown): WriteGroupMutationResponse | null {
  if (!isRecord(value)) return null;
  const workspaceRevision = parseWorkspaceRevision(value.workspace_revision);
  const group = parseStudioV2WriteGroup(value.group);
  const capabilities = parseWriteGroupCalculationCapabilities(value.calculation_capabilities);
  if (workspaceRevision === null || !group || capabilities === null) return null;
  return { workspace_revision: workspaceRevision, group, ...(capabilities !== undefined ? { calculation_capabilities: capabilities } : {}) };
}

function parseWriteGroupReadinessIssue(value: unknown): StudioV2WorkspaceReadinessIssue | null {
  if (!isRecord(value)) return null;
  const code = readRequiredText(value, 'code');
  const severity = readRequiredText(value, 'severity');
  const step = readRequiredText(value, 'step');
  const scope = readPresentText(value, 'scope');
  const message = readRequiredText(value, 'message');
  if (!code || !severity || !READINESS_SEVERITIES.has(severity as StudioV2WorkspaceReadinessSeverity) ||
    !step || !READINESS_STEPS.has(step as StudioV2WorkspaceReadinessStep) || scope === undefined || scope === null ||
    !message) {
    return null;
  }
  return {
    code,
    severity: severity as StudioV2WorkspaceReadinessSeverity,
    step: step as StudioV2WorkspaceReadinessStep,
    scope,
    message,
  };
}

export function parseWriteGroupReadinessData(value: unknown): WriteGroupReadiness | null {
  if (!isRecord(value) || typeof value.config_ready !== 'boolean' ||
    typeof value.schema_ready !== 'boolean' || typeof value.ready !== 'boolean') {
    return null;
  }
  const workspaceId = readRequiredText(value, 'workspace_id');
  const workspaceRevision = readRequiredText(value, 'workspace_revision');
  const groupId = readRequiredText(value, 'group_id');
  const groupRevision = readRequiredText(value, 'group_revision');
  const appliedRevision = readPresentText(value, 'applied_revision');
  const schemaDigest = readOptionalText(value, 'schema_digest');
  const rawIssues = value.issues;
  if (!workspaceId || !workspaceRevision || !groupId || !groupRevision || appliedRevision === undefined ||
    appliedRevision === null || schemaDigest === null || !Array.isArray(rawIssues) ||
    rawIssues.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }
  const issues = rawIssues.map(parseWriteGroupReadinessIssue);
  if (issues.some((issue) => issue === null)) return null;
  const hasBlocking = issues.some((issue) => issue?.severity === 'blocking');
  if (value.ready !== (value.config_ready && value.schema_ready && !hasBlocking) ||
    (value.schema_ready && !value.config_ready) ||
    (value.schema_ready && !schemaDigest) ||
    (!value.schema_ready && schemaDigest)) {
    return null;
  }
  return {
    workspace_id: workspaceId,
    workspace_revision: workspaceRevision,
    group_id: groupId,
    group_revision: groupRevision,
    applied_revision: appliedRevision,
    config_ready: value.config_ready,
    schema_ready: value.schema_ready,
    ready: value.ready,
    ...(schemaDigest ? { schema_digest: schemaDigest } : {}),
    issues: issues as StudioV2WorkspaceReadinessIssue[],
  };
}

/** Error for a nominal success response that cannot safely be used by the UI. */
export class WriteGroupResponseError extends Error {
  readonly code?: string;
  readonly action?: string;
  readonly retryable?: boolean;
  readonly request_id?: string;
  readonly outcome: ResponseOutcome;

  constructor(operation: string, source?: unknown, outcome: ResponseOutcome = 'unconfirmed') {
    super(`${operation} response was invalid`);
    this.name = 'WriteGroupResponseError';
    this.outcome = outcome;
    const envelope = normalizeTypedEnvelope(source);
    this.code = envelope.code;
    this.action = envelope.action;
    this.retryable = envelope.retryable;
    this.request_id = envelope.requestId;
  }
}

export function parseWriteGroupAPIData<T>(
  value: unknown, operation: string, parseData: (data: unknown) => T | null, maxDepth: number = MAX_SAFE_JSON_DEPTH,
): T {
  const bounded = parseBoundedJson(value, maxDepth);
  if (!isRecord(bounded) || typeof bounded.success !== 'boolean') {
    throw new WriteGroupResponseError(operation, bounded);
  }
  if (!bounded.success) {
    throw new WriteGroupResponseError(operation, bounded, 'failed');
  }
  if (!Object.prototype.hasOwnProperty.call(bounded, 'data')) {
    throw new WriteGroupResponseError(operation, bounded);
  }
  const parsed = parseData(bounded.data);
  if (parsed === null) throw new WriteGroupResponseError(operation, bounded);
  return parsed;
}
