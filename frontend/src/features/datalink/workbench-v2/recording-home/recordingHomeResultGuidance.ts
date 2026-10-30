import type { RecordingHomeGuidance, RecordingHomeAttentionItem } from './recordingHomeGuidance';
import { groupAction } from './recordingHomeGuidanceEvidence';
import type { recordingHomeResultSummary } from './recordingHomeResultSummary';
/** Current scoped negative evidence may guide inspection even when operation
 * owners remain unknown. It never grants a mutation or a healthy verdict. */
export function recordingHomeResultGuidance(guidance: RecordingHomeGuidance,
  summaries: Readonly<Record<string, ReturnType<typeof recordingHomeResultSummary>>>,
): RecordingHomeGuidance {
  const byId = new Map(guidance.attentionItems.filter(item => item.reason === 'operation_recovery' || item.reason === 'evidence_unavailable').map(item => [item.id, item]));
  for (const item of guidance.recordingItems) {
    const current = summaries[item.id]?.current;
    if (current !== 'source_attention' && current !== 'delivery_attention') {
      if (!byId.has(item.id) && item.evidence.group.status !== 'draft' && summaries[item.id]?.owners !== 'inspected')
        byId.set(item.id, { ...item, reason: 'evidence_unavailable' });
      continue;
    }
    const prior = byId.get(item.id);
    if (prior && prior.reason !== 'evidence_unavailable') continue;
    byId.set(item.id, { ...item, reason: current, action: groupAction(item.action.scope, item.evidence, current === 'source_attention' ? 'sources' : 'overview') });
  }
  const priority: Record<RecordingHomeAttentionItem['reason'], number> = { operation_recovery: 0, delivery_attention: 1, source_attention: 2, evidence_unavailable: 3 };
  const attentionItems = [...byId.values()].sort((a, b) => priority[a.reason] - priority[b.reason] || (a.id < b.id ? -1 : a.id > b.id ? 1 : 0));
  const first = attentionItems[0];
  const fallback = guidance.draftItems[0] ?? guidance.recordingItems[0];
  const withoutObsoleteFault = ['delivery_attention', 'source_attention', 'operation_recovery', 'evidence_unavailable'].includes(guidance.reason) && !first && guidance.availability.state === 'available' && fallback
    ? { reason: guidance.draftItems.length ? 'draft' as const : 'recordings' as const, primaryAction: fallback.action } : {};
  return { ...guidance, ...withoutObsoleteFault, attentionItems, ...(guidance.availability.state === 'available' && first ? { reason: first.reason, primaryAction: first.action } : {}) };
}
