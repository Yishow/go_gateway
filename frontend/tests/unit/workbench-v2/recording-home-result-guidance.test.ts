import { describe, expect, it } from 'vitest';
import { recordingHomeResultGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeResultGuidance';
import { recordingHomeResultSummary } from '@/features/datalink/workbench-v2/recording-home/recordingHomeResultSummary';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture, group } from './recording-home-guidance-fixtures';
it.each(['delivery_attention', 'source_attention'] as const)('legacy %s cannot bypass the canonical result decision', reason => {
  const guidance = recordingHomeGuidance(fixture(['D1'], true, [group('G')]));
  const row = guidance.recordingItems[0];
  const legacy = { ...guidance, reason, attentionItems: [{ ...row, reason }] };
  const unknown = recordingHomeResultSummary(row.evidence);
  const result = recordingHomeResultGuidance(legacy, { G: unknown });
  expect(result.reason).toBe('evidence_unavailable');
  expect(result.attentionItems.map(item => item.reason)).toEqual(['evidence_unavailable']);
});
describe('recovery is independent of current delivery/source classification', () => {
  it('does not downgrade a known retained operation when runtime facts are unavailable', () => {
    const guidance = recordingHomeGuidance(fixture(['D1'], true, [group('G')])); const row = guidance.recordingItems[0];
    const result = recordingHomeResultGuidance({ ...guidance, reason: 'operation_recovery', attentionItems: [{ ...row, reason: 'operation_recovery' }] }, { G: recordingHomeResultSummary(row.evidence) });
    expect(result.reason).toBe('operation_recovery');
  });
});
