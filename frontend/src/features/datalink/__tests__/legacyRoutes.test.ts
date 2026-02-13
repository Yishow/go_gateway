import { describe, expect, it } from 'vitest';
import {
  DASHBOARD_SECTION_INTENTS,
  DASHBOARD_MODAL_INTENTS,
  LEGACY_DECOMMISSION_ROUTES,
  buildDashboardModalRedirect,
  buildDashboardSectionRedirect,
  buildLegacyMigrationRedirect,
  isDashboardModalIntent,
  isDashboardSectionIntent,
  isLegacyDecommissionRoute,
} from '../legacyRoutes';

describe('legacyRoutes', () => {
  it('locks decommission scope to points/mappings/wizard', () => {
    expect(LEGACY_DECOMMISSION_ROUTES).toEqual(['points', 'mappings', 'wizard']);
  });

  it('builds redirect target for legacy route', () => {
    expect(buildLegacyMigrationRedirect('points')).toBe('/datalink?legacy=points&modal=points');
    expect(buildLegacyMigrationRedirect('mappings')).toBe('/datalink?legacy=mappings&modal=mappings');
    expect(buildLegacyMigrationRedirect('wizard')).toBe('/datalink?legacy=wizard&modal=wizard');
  });

  it('detects legacy route values strictly', () => {
    expect(isLegacyDecommissionRoute('points')).toBe(true);
    expect(isLegacyDecommissionRoute('mappings')).toBe(true);
    expect(isLegacyDecommissionRoute('wizard')).toBe(true);
    expect(isLegacyDecommissionRoute('tags')).toBe(false);
    expect(isLegacyDecommissionRoute(null)).toBe(false);
  });

  it('locks dashboard section intents to devices/settings', () => {
    expect(DASHBOARD_SECTION_INTENTS).toEqual(['devices', 'settings']);
  });

  it('builds redirect target for dashboard section', () => {
    expect(buildDashboardSectionRedirect('devices')).toBe('/datalink?section=devices');
    expect(buildDashboardSectionRedirect('settings')).toBe('/datalink?section=settings');
  });

  it('detects dashboard section values strictly', () => {
    expect(isDashboardSectionIntent('devices')).toBe(true);
    expect(isDashboardSectionIntent('settings')).toBe(true);
    expect(isDashboardSectionIntent('points')).toBe(false);
    expect(isDashboardSectionIntent(null)).toBe(false);
  });

  it('locks dashboard modal intents to migrated sidebar features', () => {
    expect(DASHBOARD_MODAL_INTENTS).toEqual([
      'devices',
      'settings',
      'points',
      'mappings',
      'wizard',
      'polling-groups',
      'tags',
    ]);
  });

  it('builds redirect target for dashboard modal intent', () => {
    expect(buildDashboardModalRedirect('points')).toBe('/datalink?modal=points');
    expect(buildDashboardModalRedirect('polling-groups')).toBe('/datalink?modal=polling-groups');
    expect(buildDashboardModalRedirect('tags')).toBe('/datalink?modal=tags');
  });

  it('detects dashboard modal values strictly', () => {
    expect(isDashboardModalIntent('devices')).toBe(true);
    expect(isDashboardModalIntent('settings')).toBe(true);
    expect(isDashboardModalIntent('points')).toBe(true);
    expect(isDashboardModalIntent('wizard')).toBe(true);
    expect(isDashboardModalIntent('test')).toBe(false);
    expect(isDashboardModalIntent(null)).toBe(false);
  });
});
