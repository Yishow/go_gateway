export type GatewayUiVersion = 'dual_quick' | 'dual_expert';

export function resolveGatewayUiVersion(pathname?: string): GatewayUiVersion | null {
  const path = pathname ?? (typeof window !== 'undefined' ? window.location.pathname : '');

  if (path.startsWith('/gateway/quick-setup')) return 'dual_quick';
  if (path.startsWith('/gateway/expert-workbench')) return 'dual_expert';

  return null;
}
