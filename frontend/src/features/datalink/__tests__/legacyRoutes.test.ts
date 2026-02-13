import { describe, expect, it } from 'vitest';
import {
  LEGACY_DECOMMISSION_ROUTES,
  buildLegacyMigrationRedirect,
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
});
