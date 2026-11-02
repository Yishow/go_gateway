import type { MappingQueueObservation, RuleQueueObservation } from './sourceQueueObservation';
/** Original source owner is always notified first. View observers cannot fail it. */
export function composeRuleObservation(source: RuleQueueObservation, view?: Pick<RuleQueueObservation, 'submitted' | 'saved' | 'failed'>): RuleQueueObservation {
  return { ...source,
    submitted: (rule, generation) => { source.submitted?.(rule, generation); try { view?.submitted?.(rule, generation); } catch { /* A view cannot fail the save pipeline. */ } },
    saved: (rule, saved) => { source.saved(rule, saved); try { view?.saved?.(rule, saved); } catch { /* Same non-authoritative notification boundary. */ } },
    failed: (rule, error) => { source.failed(rule, error); try { view?.failed?.(rule, error); } catch { /* Preserve the original failure. */ } },
  };
}
export function composeMappingObservation(source: MappingQueueObservation, name: MappingQueueObservation): MappingQueueObservation {
  return {
    acknowledged: (point, mapping, saved) => {
      if (!source.acknowledged(point, mapping, saved)) { name.failed(point); return false; }
      return name.acknowledged(point, mapping, saved);
    },
    saved: (point, mapping, saved) => { source.saved(point, mapping, saved); name.saved(point, mapping, saved); },
    failed: point => { source.failed(point); name.failed(point); },
  };
}
