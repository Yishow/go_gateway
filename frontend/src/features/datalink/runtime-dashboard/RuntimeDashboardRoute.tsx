import { RuntimeDashboardPage } from './RuntimeDashboardPage';
import { useRuntimeDashboardState } from './useRuntimeDashboardState';

export function RuntimeDashboardRoute() {
  const state = useRuntimeDashboardState();
  return <RuntimeDashboardPage {...state} />;
}
