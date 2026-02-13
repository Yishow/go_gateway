export const LEGACY_DECOMMISSION_ROUTES = ['points', 'mappings', 'wizard'] as const;
export const DASHBOARD_SECTION_INTENTS = ['devices', 'settings'] as const;

export type LegacyDecommissionRoute = (typeof LEGACY_DECOMMISSION_ROUTES)[number];
export type DashboardSectionIntent = (typeof DASHBOARD_SECTION_INTENTS)[number];

export function isLegacyDecommissionRoute(route: string | null): route is LegacyDecommissionRoute {
  if (!route) return false;
  return LEGACY_DECOMMISSION_ROUTES.includes(route as LegacyDecommissionRoute);
}

export function buildLegacyMigrationRedirect(route: LegacyDecommissionRoute): string {
  return `/datalink?legacy=${route}`;
}

export function isDashboardSectionIntent(section: string | null): section is DashboardSectionIntent {
  if (!section) return false;
  return DASHBOARD_SECTION_INTENTS.includes(section as DashboardSectionIntent);
}

export function buildDashboardSectionRedirect(section: DashboardSectionIntent): string {
  return `/datalink?section=${section}`;
}
