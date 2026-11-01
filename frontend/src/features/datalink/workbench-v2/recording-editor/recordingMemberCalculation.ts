import type { MeasurementDefinition } from '../../../../types/measurement';
import type { WriteGroupCalculationCapabilities } from '../../../../types/studioV2WriteGroupCalculation';
import type { WriteGroupMemberCalculation } from '../../../../types/studioV2WriteGroup';
import type { EditorMember } from '../state/writeGroup/draft';
import { isConfirmedCounterDefinition } from '../state/circuitCounterConfirmation';
export interface CalculationEditorMember extends EditorMember { /** UI-only lexical input; never serialized to the API. */ cadence_input?: string; gauge_confirmation?: string }
export interface MemberCalculationContext { workspaceId: string; recordId?: string; definitions: MeasurementDefinition[]; definitionsAvailable?: boolean; capabilities?: WriteGroupCalculationCapabilities; layout: string; dialect: string }
const numeric = new Set(['int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64']);
function definitionFor(member: EditorMember, context: MemberCalculationContext): MeasurementDefinition | undefined {
  const matches = context.definitions.filter(definition => definition.workspace_id === context.workspaceId && definition.device_id === member.device_id &&
    definition.point_id === member.point_id && definition.tag_id === member.tag_id && (!member.measurement_id || definition.id === member.measurement_id));
  if (matches.length !== 1) return undefined;
  const definition = matches[0];
  return [definition.id, definition.definition_revision, definition.source_binding_revision, definition.series_epoch].every(value => value?.trim()) ? definition : undefined;
}
const sourceDefinitions = (member: EditorMember, context: MemberCalculationContext) => context.definitions.filter(definition =>
  definition.workspace_id === context.workspaceId && definition.device_id === member.device_id && definition.point_id === member.point_id && definition.tag_id === member.tag_id);
const gaugeScope = (member: EditorMember, context: MemberCalculationContext) => JSON.stringify([context.workspaceId, context.recordId,
  member.device_id, member.point_id, member.tag_id, member.source_value_type, member.source_revision, member.mapping_revision]);
function canonicalGaugeIntent(member: CalculationEditorMember, context: MemberCalculationContext): boolean {
  if (member.value_role !== undefined && member.value_role !== 'gauge') return false;
  if (member.gauge_confirmation !== undefined) return member.gauge_confirmation === gaugeScope(member, context);
  return Boolean(member.source_revision?.trim() && member.mapping_revision?.trim() &&
    (member.value_role === 'gauge' || member.calculation === 'mean' && member.value_role === undefined));
}
/** Denial-only local explanation. Neither semantic eligibility nor capability proves a record ready to run. */
export function memberCalculationReason(member: EditorMember, method: WriteGroupMemberCalculation, context: MemberCalculationContext): string | undefined {
  const cap = context.capabilities;
  if (!cap?.editor_supported || !cap.draft_supported || !cap.methods.includes(method)) return 'capability';
  if (!cap.schema_dialects.some(value => value === context.dialect)) return 'dialect';
  if (!cap.schema_layouts.some(value => value === context.layout)) return 'layout';
  if (member.exact_recording_source) {
    const restriction = cap.unsupported_combinations.find(value => value.method === method && value.source_transform === 'exact_counter_scale' && (!value.layout || value.layout === context.layout));
    if (restriction) return restriction.reason;
  }
  if (method === 'latest') return undefined;
  const definition = definitionFor(member, context);
  const numericSource = numeric.has(member.source_value_type ?? '');
  if (method === 'mean') {
    if (!numericSource) return 'gauge_required';
    if (definition) return definition.semantic_kind === 'gauge' ? undefined : 'gauge_required';
    if (member.measurement_id?.trim() || member.exact_recording_source || sourceDefinitions(member, context).length ||
      context.definitions.some(value => value.workspace_id !== context.workspaceId)) return 'gauge_required';
    return canonicalGaugeIntent(member, context) ? undefined : 'gauge_required';
  }
  return definition && numericSource && isConfirmedCounterDefinition(definition) ? undefined : 'counter_required';
}
export function changeMemberCalculation(member: EditorMember, method: WriteGroupMemberCalculation, context: MemberCalculationContext): EditorMember {
  if (memberCalculationReason(member, method, context) || member.calculation === method) return member;
  const definition = definitionFor(member, context);
  const next: CalculationEditorMember = { ...member, calculation: method };
  if (method === 'latest' && member.calculation === 'mean' && !member.measurement_id?.trim() && canonicalGaugeIntent(member, context)) next.gauge_confirmation = gaugeScope(member, context);
  // Result SQL type belongs to the newly resolved intent, never the previous method's cached proof.
  delete next.result_value_type;
  if (method === 'mean') next.value_role = 'gauge';
  if (method !== 'latest' && definition) { next.measurement_id = definition.id; next.value_role = method === 'mean' ? 'gauge' : 'counter_raw'; }
  return next;
}
export function cadenceInputInvalid(member: CalculationEditorMember): boolean {
  return member.cadence_input !== undefined && (!member.cadence_input.trim() || !Number.isFinite(Number(member.cadence_input)) || Number(member.cadence_input) !== member.cadence?.interval);
}

export function initializeMixedMembers(members: EditorMember[], previous: EditorMember[], candidates: import('../state/writeGroup/candidates').GroupCandidate[], cap: WriteGroupCalculationCapabilities | undefined, initialCadence?: EditorMember['cadence']): EditorMember[] {
  const identity = (member: EditorMember) => JSON.stringify([member.device_id, member.point_id, member.tag_id]);
  const old = new Set(previous.map(identity));
  return members.map(member => {
    // Restored intent (including invalid text) is not a new source and must never receive defaults again.
    if (old.has(identity(member)) || member.calculation !== undefined || member.cadence !== undefined || 'cadence_input' in member) return member;
    const source = candidates.find(value => value.device_id === member.device_id && value.point_id === member.point_id && value.tag_id === member.tag_id);
    return { ...member, calculation: 'latest', cadence: initialCadence ? { ...initialCadence } : { interval: cap?.cadence.minimum ?? 15, unit: !cap || cap.cadence.units.includes('second') ? 'second' : 'minute', alignment: 'utc_epoch' }, ...(source ? { source_value_type: source.target_type } : {}) };
  });
}

export function canConfirmMemberGauge(member: CalculationEditorMember, context: MemberCalculationContext): boolean {
  const cap = context.capabilities;
  return Boolean([context.workspaceId, context.recordId, member.device_id, member.point_id, member.tag_id].every(value => value?.trim()) && context.definitionsAvailable === true && cap?.editor_supported && cap.draft_supported && cap.methods.includes('mean') &&
    cap.schema_dialects.some(value => value === context.dialect) && cap.schema_layouts.some(value => value === context.layout) &&
    numeric.has(member.source_value_type ?? '') && !member.measurement_id?.trim() && !member.exact_recording_source &&
    member.value_role !== 'counter_raw' && member.value_role !== 'state' && member.calculation !== 'delta' &&
    context.definitions.every(value => value.workspace_id === context.workspaceId) && sourceDefinitions(member, context).length === 0);
}
export function confirmMemberGauge(member: CalculationEditorMember, context: MemberCalculationContext): CalculationEditorMember {
  if (!canConfirmMemberGauge(member, context)) return member;
  return { ...member, value_role: 'gauge', gauge_confirmation: gaugeScope(member, context) };
}
