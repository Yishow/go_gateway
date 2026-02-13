import { describe, expect, it } from 'vitest';
import {
  DASHBOARD_SECTION_INTENTS,
  LEGACY_DECOMMISSION_ROUTES,
  buildDashboardSectionRedirect,
  buildLegacyMigrationRedirect,
  isDashboardSectionIntent,
  isLegacyDecommissionRoute,
} from '../legacyRoutes';

describe('legacyRoutes', () => {
  it('locks decommission scope to points/mappings/wizard', () => {
    expect(LEGACY_DECOMMISSION_ROUTES).toEqual(['points', 'mappings', 'wizard']);
  });

  it('builds redirect target for legacy route', () => {
    expect(buildLegacyMigrationRedirect('points')).toBe('/datalink?legacy=points');
    expect(buildLegacyMigrationRedirect('mappings')).toBe('/datalink?legacy=mappings');
    expect(buildLegacyMigrationRedirect('wizard')).toBe('/datalink?legacy=wizard');
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
});
