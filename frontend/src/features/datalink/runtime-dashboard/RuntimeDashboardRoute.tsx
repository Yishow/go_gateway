import { useLocation, useNavigate } from 'react-router-dom';
import { RuntimeDashboardPage } from './RuntimeDashboardPage';
import { useRuntimeDashboardState } from './useRuntimeDashboardState';
import { createRuntimeEntityHistory, readRuntimeEntityHistory } from './runtimeEntityHistory';
import type { RuntimeSourceSelection } from '../../../types/studioV2RuntimeContext';

export function RuntimeDashboardRoute() {
  const state = useRuntimeDashboardState();
  const navigate = useNavigate();
  const location = useLocation();
  const scope = state.sourceEvidenceScope ?? null;
  const sourceSelection = readRuntimeEntityHistory(location.state, scope);
  const onSourceSelection = (selection: RuntimeSourceSelection | null) => {
    if (!scope) return;
    const history = createRuntimeEntityHistory(scope, selection);
    if (!history) return;
    // One explicit selection event updates this entry; no effects or history rewrite loops.
    void navigate({ pathname: location.pathname, search: location.search }, { replace: true, state: history });
  };
  return <RuntimeDashboardPage {...state} navigateTo={navigate} sourceSelection={sourceSelection} onSourceSelection={onSourceSelection} />;
}
