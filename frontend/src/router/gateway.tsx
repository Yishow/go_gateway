import type { ReactElement } from 'react';
import { Navigate } from 'react-router-dom';
import GatewayEntryPage from '../pages/gateway/GatewayEntryPage';
import GatewayQuickSetupPage from '../pages/gateway/GatewayQuickSetupPage';
import GatewayExpertWorkbenchPage from '../pages/gateway/GatewayExpertWorkbenchPage';
import { useGatewayDualEntryFlag } from '../features/gateway/useGatewayDualEntryFlag';
import { resolveGatewayCreateEntryPath } from '../features/gateway/dualEntryFlag';

function GatewayRouteLoading() {
  return <div className="p-6 text-sm text-slate-300">載入入口設定中...</div>;
}

function GatewayDualEntryGuard({ children }: { children: ReactElement }) {
  const { data: enabled, isLoading } = useGatewayDualEntryFlag();

  if (isLoading) return <GatewayRouteLoading />;
  if (!enabled) return <Navigate to="/datalink" replace />;

  return children;
}

export function GatewayCreateEntryRedirect() {
  const { data: enabled, isLoading } = useGatewayDualEntryFlag();

  if (isLoading) return <GatewayRouteLoading />;

  return <Navigate to={resolveGatewayCreateEntryPath(Boolean(enabled))} replace />;
}

export function GatewayEntryRoute() {
  return (
    <GatewayDualEntryGuard>
      <GatewayEntryPage />
    </GatewayDualEntryGuard>
  );
}

export function GatewayQuickSetupRoute() {
  return (
    <GatewayDualEntryGuard>
      <GatewayQuickSetupPage />
    </GatewayDualEntryGuard>
  );
}

export function GatewayExpertWorkbenchRoute() {
  return (
    <GatewayDualEntryGuard>
      <GatewayExpertWorkbenchPage />
    </GatewayDualEntryGuard>
  );
}
