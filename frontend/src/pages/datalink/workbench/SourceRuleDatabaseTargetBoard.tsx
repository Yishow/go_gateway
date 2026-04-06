import type { ComponentProps } from 'react';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import { DatabaseTargetBoard } from './DatabaseTargetBoard';
import { useWorkbench } from './WorkbenchProvider';

export function SourceRuleDatabaseTargetBoard(
  props: ComponentProps<typeof DatabaseTargetBoard>,
) {
  const { crossStepContext } = useWorkbench();
  const reviewRuleId = crossStepContext.focusedRuleId;
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
