import type { ModbusShareStatus } from '../../../types/datalink';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import { LocalModbusReviewPanel } from './LocalModbusReviewPanel';
import { useWorkbench } from './WorkbenchProvider';

type SourceRuleLocalModbusReviewSurfaceProps = {
  selectedTagId: string;
  status: ModbusShareStatus | null;
  conflictCount: number;
};

export function SourceRuleLocalModbusReviewSurface({
  selectedTagId,
  status,
  conflictCount,
}: SourceRuleLocalModbusReviewSurfaceProps) {
  const { crossStepContext } = useWorkbench();
  const reviewRuleId = crossStepContext.focusedRuleId;
  const reviewQuery = useSourceRuleCandidatesQuery(reviewRuleId);

  return (
    <LocalModbusReviewPanel
      reviewRuleId={reviewRuleId}
      reviewRevisionId={reviewQuery.data?.revision_id ?? null}
      reviewSet={reviewQuery.data?.local_modbus_outputs ?? null}
      reviewLoading={reviewQuery.isLoading || reviewQuery.isFetching}
      selectedTagId={selectedTagId}
      status={status}
      conflictCount={conflictCount}
    />
  );
}
