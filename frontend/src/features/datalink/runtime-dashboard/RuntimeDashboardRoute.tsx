import { useNavigate } from 'react-router-dom';
import { RuntimeDashboardPage } from './RuntimeDashboardPage';
import { useRuntimeDashboardState } from './useRuntimeDashboardState';

export function RuntimeDashboardRoute() {
  const state = useRuntimeDashboardState();
  const navigate = useNavigate();

  return <RuntimeDashboardPage {...state} navigateTo={navigate} />;
}
