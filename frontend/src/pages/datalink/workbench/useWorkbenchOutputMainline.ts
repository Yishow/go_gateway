import { useMemo } from 'react';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import { useWorkbench } from './WorkbenchProvider';
import {
  buildWorkbenchOutputMainlineState,
  type WorkbenchOutputMainlineState,
} from './workbenchOutputMainlineModel';

export type UseWorkbenchOutputMainlineInput = {
  hasSelectedDevice: boolean;
  sourceReady: boolean;
  tagReady: boolean;
};

export type UseWorkbenchOutputMainlineResult = WorkbenchOutputMainlineState & {
  activeRuleId: string | null;
};

export function useWorkbenchOutputMainline(
  input: UseWorkbenchOutputMainlineInput,
): UseWorkbenchOutputMainlineResult {
  const { crossStepContext, sourcePlanningState } = useWorkbench();
  const activeRuleId =
    crossStepContext.focusedRuleId ?? sourcePlanningState.selectedRuleId ?? null;
  const { data: snapshot } = useSourceRuleCandidatesQuery(activeRuleId);

  return useMemo(
    () => ({
      activeRuleId,
      ...buildWorkbenchOutputMainlineState({
        ...input,
        activeRuleId,
        snapshot,
      }),
    }),
    [activeRuleId, input, snapshot],
  );
}
