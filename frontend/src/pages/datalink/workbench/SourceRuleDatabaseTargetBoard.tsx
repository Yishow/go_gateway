import type { ComponentProps } from 'react';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import { useSourceRulesQuery } from '../../../hooks/datalink/useSourceRules';
import { DatabaseTargetBoard } from './DatabaseTargetBoard';
import { useWorkbench } from './WorkbenchProvider';
import { resolvePersistedReviewRuleId } from './reviewRuleId';

export function SourceRuleDatabaseTargetBoard(
  props: ComponentProps<typeof DatabaseTargetBoard>,
) {
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
    <DatabaseTargetBoard
      {...props}
      reviewRuleId={reviewRuleId}
      reviewRevisionId={reviewQuery.data?.revision_id ?? null}
      reviewSet={reviewQuery.data?.database_outputs ?? null}
      reviewLoading={reviewQuery.isLoading || reviewQuery.isFetching}
    />
  );
}
