import { useMemo } from 'react';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import { useSourceRulesQuery } from '../../../hooks/datalink/useSourceRules';
import type { Point } from '../../../types/datalink';
import { resolveActiveRuleId, scopePointsToActiveRule } from './sourceRuleSelection';
import { useWorkbench } from './WorkbenchProvider';

export function useTagBindingSourceScope(points: Point[]) {
  const { crossStepContext, selectedDeviceId, sourcePlanningState } = useWorkbench();
  const sourceRulesQuery = useSourceRulesQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const persistedRules = useMemo(
    () => (sourceRulesQuery.data ?? []).filter((rule) => rule.device_id === selectedDeviceId),
    [selectedDeviceId, sourceRulesQuery.data],
  );
  const activeRuleId = useMemo(
    () =>
      resolveActiveRuleId(persistedRules, [
        crossStepContext.focusedRuleId,
        sourcePlanningState.selectedRuleId,
      ]),
    [crossStepContext.focusedRuleId, persistedRules, sourcePlanningState.selectedRuleId],
  );
  const activeRuleCandidateQuery = useSourceRuleCandidatesQuery(activeRuleId);
  return useMemo(
    () => scopePointsToActiveRule(points, activeRuleId, activeRuleCandidateQuery.data),
    [activeRuleCandidateQuery.data, activeRuleId, points],
  );
}
