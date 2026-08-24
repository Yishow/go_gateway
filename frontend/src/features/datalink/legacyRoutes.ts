export const LEGACY_DECOMMISSION_ROUTES = ['points', 'mappings', 'wizard'] as const;
export const DASHBOARD_SECTION_INTENTS = ['devices', 'settings'] as const;
export const DASHBOARD_MODAL_INTENTS = [
  'devices',
  'settings',
  'points',
  'mappings',
  'wizard',
  'polling-groups',
  'tags',
] as const;

export const WORKBENCH_COMPAT_STEPS = ['device', 'source', 'tag', 'output'] as const;
export const WORKBENCH_COMPAT_TARGETS = ['modbus', 'database'] as const;

export type LegacyDecommissionRoute = (typeof LEGACY_DECOMMISSION_ROUTES)[number];
export type DashboardSectionIntent = (typeof DASHBOARD_SECTION_INTENTS)[number];
export type DashboardModalIntent = (typeof DASHBOARD_MODAL_INTENTS)[number];
export type WorkbenchCompatStep = (typeof WORKBENCH_COMPAT_STEPS)[number];
export type WorkbenchCompatTarget = (typeof WORKBENCH_COMPAT_TARGETS)[number];

export function isLegacyDecommissionRoute(route: string | null): route is LegacyDecommissionRoute {
  if (!route) return false;
  return LEGACY_DECOMMISSION_ROUTES.includes(route as LegacyDecommissionRoute);
}

export function buildLegacyMigrationRedirect(route: LegacyDecommissionRoute): string {
  return `/datalink?legacy=${route}&modal=${route}`;
}

export function isDashboardSectionIntent(section: string | null): section is DashboardSectionIntent {
  if (!section) return false;
  return DASHBOARD_SECTION_INTENTS.includes(section as DashboardSectionIntent);
}

export function buildDashboardSectionRedirect(section: DashboardSectionIntent): string {
  return `/datalink?section=${section}`;
}

export function isDashboardModalIntent(modal: string | null): modal is DashboardModalIntent {
  if (!modal) return false;
  return DASHBOARD_MODAL_INTENTS.includes(modal as DashboardModalIntent);
}

export function buildDashboardModalRedirect(modal: DashboardModalIntent): string {
  return `/datalink?modal=${modal}`;
}

export function buildWorkbenchV2EntryRedirect(): string {
  return '/studio/v2';
}

function buildWorkbenchV2Query(options?: {
  step?: WorkbenchCompatStep;
  target?: WorkbenchCompatTarget;
  section?: string | null;
}): string {
  const params = new URLSearchParams();
  if (options?.step) {
    params.set('step', options.step);
  }
  if (options?.target) {
    params.set('target', options.target);
  }
  if (options?.section) {
    params.set('section', options.section);
  }
  const query = params.toString();
  return query ? `${buildWorkbenchV2EntryRedirect()}?${query}` : buildWorkbenchV2EntryRedirect();
}

export function buildLocalModbusCompatRedirect(section?: string | null): string {
  return buildWorkbenchV2Query({
    step: 'output',
    target: 'modbus',
    section,
  });
}
