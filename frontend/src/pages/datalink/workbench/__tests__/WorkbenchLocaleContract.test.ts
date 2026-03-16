import enCommon from '../../../../i18n/locales/en/common.json';
import zhTWCommon from '../../../../i18n/locales/zh-TW/common.json';
import { describe, expect, it } from 'vitest';

describe('Workbench locale contract', () => {
  it('ships localized zh-TW strings for the foundation shell landmarks and steps', () => {
    expect(zhTWCommon.workbench.title).not.toBe(enCommon.workbench.title);
    expect(zhTWCommon.workbench.foundation.label).not.toBe(enCommon.workbench.foundation.label);
    expect(zhTWCommon.workbench.header.ariaLabel).not.toBe(enCommon.workbench.header.ariaLabel);
    expect(zhTWCommon.workbench.header.eyebrow).not.toBe(enCommon.workbench.header.eyebrow);
    expect(zhTWCommon.workbench.stepNavigator.ariaLabel).not.toBe(
      enCommon.workbench.stepNavigator.ariaLabel,
    );
    expect(zhTWCommon.workbench.actionDock.ariaLabel).not.toBe(
      enCommon.workbench.actionDock.ariaLabel,
    );

    expect(zhTWCommon.workbench.steps.device).not.toBe(enCommon.workbench.steps.device);
    expect(zhTWCommon.workbench.steps.source).not.toBe(enCommon.workbench.steps.source);
    expect(zhTWCommon.workbench.steps.tag).not.toBe(enCommon.workbench.steps.tag);
    expect(zhTWCommon.workbench.steps.output).not.toBe(enCommon.workbench.steps.output);
  });
});
