import type { recordingHomeResultSummary } from './recordingHomeResultSummary';

type Summary = ReturnType<typeof recordingHomeResultSummary>;
export type WriteStageState = 'complete' | 'unconfirmed' | 'attention';

/** Only matching receipt evidence can complete the current SQL stage. */
export function recordingHomeWriteStageMapping(summary: Summary): Record<'configuration' | 'activation' | 'sql' | 'readback', WriteStageState> {
  return {
    configuration: summary.configuration === 'requested' ? 'complete' : 'unconfirmed',
    activation: summary.effective === 'requested_effective' ? 'complete' : summary.effective === 'requested_not_effective' ? 'attention' : 'unconfirmed',
    sql: summary.current === 'delivery_attention' ? 'attention' :
      summary.sql === 'revision_receipt' || summary.sql === 'effective_revision_receipt' ? 'complete' : 'unconfirmed',
    readback: summary.readback === 'test_write_verified' ? 'complete' : 'unconfirmed',
  };
}
