import { describe, expect, it } from 'vitest';
import { createSettingsOperationOwnership } from '../../../src/features/datalink/workbench-v2/settings/settingsOperationOwnership';

describe('Settings operation ownership', () => {
  it('keeps an older request from owning the newer error/retry target', () => {
    const ownership = createSettingsOperationOwnership();
    const older = ownership.begin();
    const newer = ownership.begin();

    expect(ownership.isCurrent(older)).toBe(false);
    expect(ownership.isCurrent(newer)).toBe(true);
  });

  it('does not make unrelated operation scopes stale', () => {
    const ownership = createSettingsOperationOwnership();
    const connector = ownership.begin('connector:conn-prod');
    const settings = ownership.begin('settings-save');

    expect(ownership.isCurrent(connector, 'connector:conn-prod')).toBe(true);
    expect(ownership.isCurrent(settings, 'settings-save')).toBe(true);
    expect(ownership.isCurrent(connector, 'settings-save')).toBe(false);
  });
});
