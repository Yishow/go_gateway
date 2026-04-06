import type { ModbusShareStatus } from '../../../types/datalink';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import { useSourceRulesQuery } from '../../../hooks/datalink/useSourceRules';
import { LocalModbusReviewPanel } from './LocalModbusReviewPanel';
import { useWorkbench } from './WorkbenchProvider';
import { resolvePersistedReviewRuleId } from './reviewRuleId';

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
  const { crossStepContext, selectedDeviceId, sourcePlanningState } = useWorkbench();
  const persistedRulesQuery = useSourceRulesQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const reviewRuleId = resolvePersistedReviewRuleId(
    [crossStepContext.focusedRuleId, sourcePlanningState.selectedRuleId],
    persistedRulesQuery.data ?? [],
  );
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
