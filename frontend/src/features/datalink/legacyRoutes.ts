export const LEGACY_DECOMMISSION_ROUTES = ['points', 'mappings', 'wizard'] as const;

export type LegacyDecommissionRoute = (typeof LEGACY_DECOMMISSION_ROUTES)[number];

export function isLegacyDecommissionRoute(route: string | null): route is LegacyDecommissionRoute {
  if (!route) return false;
  return LEGACY_DECOMMISSION_ROUTES.includes(route as LegacyDecommissionRoute);
}

export function buildLegacyMigrationRedirect(route: LegacyDecommissionRoute): string {
  return `/datalink?legacy=${route}`;
}
