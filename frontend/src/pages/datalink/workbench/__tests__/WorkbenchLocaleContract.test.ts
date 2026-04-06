import enCommon from '../../../../i18n/locales/en/common.json';
import zhTWCommon from '../../../../i18n/locales/zh-TW/common.json';
import { describe, expect, it } from 'vitest';

describe('Workbench locale contract', () => {
  it('ships localized zh-TW strings for the foundation shell landmarks and steps', () => {
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

  it('positions /studio as the primary guided workflow in shell copy', () => {
    expect(zhTWCommon.workbench.title).toContain('Studio');
    expect(enCommon.workbench.title).toContain('Studio');
    expect(zhTWCommon.workbench.subtitle).toContain('/studio');
    expect(enCommon.workbench.subtitle).toContain('/studio');
    expect(zhTWCommon.workbench.header.eyebrow).toContain('Studio');
    expect(enCommon.workbench.header.eyebrow).toContain('Studio');
    expect(zhTWCommon.workbench.actionDock.nextAction.advanceToOutput).toContain('Studio');
    expect(enCommon.workbench.actionDock.nextAction.advanceToOutput).toContain('Studio');
    expect(zhTWCommon.workbench.actionDock.nextAction.configureOutput).toContain('資料庫');
    expect(enCommon.workbench.actionDock.nextAction.configureOutput).toContain('Database');
  });

  it('ships localized zh-TW strings for the five-region shell components', () => {
    expect(zhTWCommon.workbench.stepRail.ariaLabel).not.toBe(enCommon.workbench.stepRail.ariaLabel);
    expect(zhTWCommon.workbench.stepRail.heading).not.toBe(enCommon.workbench.stepRail.heading);
    expect(zhTWCommon.workbench.contextBar.ariaLabel).not.toBe(enCommon.workbench.contextBar.ariaLabel);
    expect(zhTWCommon.workbench.contextBar.noDevice).not.toBe(enCommon.workbench.contextBar.noDevice);
    expect(zhTWCommon.workbench.contextBar.actions.gotoTag).not.toBe(
      enCommon.workbench.contextBar.actions.gotoTag,
    );
    expect(zhTWCommon.workbench.inspector.ariaLabel).not.toBe(enCommon.workbench.inspector.ariaLabel);
    expect(zhTWCommon.workbench.inspector.heading).not.toBe(enCommon.workbench.inspector.heading);
    expect(zhTWCommon.workbench.inspector.placeholder).not.toBe(enCommon.workbench.inspector.placeholder);
    expect(zhTWCommon.workbench.bottomSummary.ariaLabel).not.toBe(enCommon.workbench.bottomSummary.ariaLabel);
    expect(zhTWCommon.workbench.bottomSummary.points).not.toBe(enCommon.workbench.bottomSummary.points);
  });

  it('ships localized zh-TW strings for the shared state shell (readiness, inspector, output target)', () => {
    // Readiness – device label
    expect(zhTWCommon.workbench.bottomSummary.readiness.device).not.toBe(
      enCommon.workbench.bottomSummary.readiness.device,
    );

    // Output target
    expect(zhTWCommon.workbench.bottomSummary.target).not.toBe(
      enCommon.workbench.bottomSummary.target,
    );
    expect(zhTWCommon.workbench.bottomSummary.targets.modbus).not.toBe(
      enCommon.workbench.bottomSummary.targets.modbus,
    );
    expect(zhTWCommon.workbench.bottomSummary.targets.database).not.toBe(
      enCommon.workbench.bottomSummary.targets.database,
    );

    // Inspector selection kinds
    expect(zhTWCommon.workbench.inspector.noSelection).not.toBe(
      enCommon.workbench.inspector.noSelection,
    );
    expect(zhTWCommon.workbench.inspector.selectionKind.device).not.toBe(
      enCommon.workbench.inspector.selectionKind.device,
    );
    expect(zhTWCommon.workbench.inspector.selectionKind.rule).not.toBe(
      enCommon.workbench.inspector.selectionKind.rule,
    );
    expect(zhTWCommon.workbench.inspector.selectionKind.span).not.toBe(
      enCommon.workbench.inspector.selectionKind.span,
    );
    expect(zhTWCommon.workbench.inspector.selectionKind.tag).not.toBe(
      enCommon.workbench.inspector.selectionKind.tag,
    );
    expect(zhTWCommon.workbench.inspector.selectionKind.outputCandidate).not.toBe(
      enCommon.workbench.inspector.selectionKind.outputCandidate,
    );
  });

  it('ships localized zh-TW strings for Step 1 capability, clone, and recent-test surfaces', () => {
    expect(zhTWCommon.workbench.device.actions.clone).not.toBe(
      enCommon.workbench.device.actions.clone,
    );
    expect(zhTWCommon.workbench.device.panel.cloneTitle).not.toBe(
      enCommon.workbench.device.panel.cloneTitle,
    );
    expect(zhTWCommon.workbench.device.panel.cloneDescription).not.toBe(
      enCommon.workbench.device.panel.cloneDescription,
    );
    expect(zhTWCommon.workbench.device.inspector.capabilitySummary).not.toBe(
      enCommon.workbench.device.inspector.capabilitySummary,
    );
    expect(zhTWCommon.workbench.device.inspector.recentTests).not.toBe(
      enCommon.workbench.device.inspector.recentTests,
    );
    expect(zhTWCommon.workbench.device.capability.labels.addressBase).not.toBe(
      enCommon.workbench.device.capability.labels.addressBase,
    );
    expect(zhTWCommon.workbench.device.capability.labels.wordOrder).not.toBe(
      enCommon.workbench.device.capability.labels.wordOrder,
    );
    expect(zhTWCommon.workbench.device.capability.labels.unitId).not.toBe(
      enCommon.workbench.device.capability.labels.unitId,
    );
    expect(zhTWCommon.workbench.device.capability.labels.protocolTraits).not.toBe(
      enCommon.workbench.device.capability.labels.protocolTraits,
    );
    expect(zhTWCommon.workbench.device.connection.dataFormats.abcd).not.toBe(
      enCommon.workbench.device.connection.dataFormats.abcd,
    );
    expect(zhTWCommon.workbench.device.connection.dataFormats.cdab).toBeTruthy();
    expect(enCommon.workbench.device.connection.dataFormats.cdab).toBeTruthy();
  });

  it('ships localized zh-TW strings for the compact source, tag, and output guidance additions', () => {
    expect(zhTWCommon.workbench.source.planner.addRule).not.toBe(
      enCommon.workbench.source.planner.addRule,
    );
    expect(zhTWCommon.workbench.source.ruleLayer.delete).not.toBe(
      enCommon.workbench.source.ruleLayer.delete,
    );
    expect(zhTWCommon.workbench.source.ruleLayer.descriptionHint).not.toBe(
      enCommon.workbench.source.ruleLayer.descriptionHint,
    );
    expect(zhTWCommon.workbench.source.toolbar.moreTools).not.toBe(
      enCommon.workbench.source.toolbar.moreTools,
    );
    expect(zhTWCommon.workbench.source.ruleLayer.protectPlan).not.toBe(
      enCommon.workbench.source.ruleLayer.protectPlan,
    );
    expect(zhTWCommon.workbench.source.ruleLayer.unprotectPlan).not.toBe(
      enCommon.workbench.source.ruleLayer.unprotectPlan,
    );
    expect(zhTWCommon.workbench.source.ruleLayer.protectHint).not.toBe(
      enCommon.workbench.source.ruleLayer.protectHint,
    );
    expect(zhTWCommon.workbench.source.conflictQueue.title).not.toBe(
      enCommon.workbench.source.conflictQueue.title,
    );
    expect(zhTWCommon.workbench.source.conflictQueue.step3Blocked).not.toBe(
      enCommon.workbench.source.conflictQueue.step3Blocked,
    );
    expect(zhTWCommon.workbench.source.planner.dataTypeUnsupported).not.toBe(
      enCommon.workbench.source.planner.dataTypeUnsupported,
    );
    expect(zhTWCommon.workbench.source.planner.advanced.heading).not.toBe(
      enCommon.workbench.source.planner.advanced.heading,
    );
    expect(zhTWCommon.workbench.source.planner.collapsedHint).not.toBe(
      enCommon.workbench.source.planner.collapsedHint,
    );
    expect(zhTWCommon.workbench.source.planner.expandSection).not.toBe(
      enCommon.workbench.source.planner.expandSection,
    );
    expect(zhTWCommon.workbench.source.ruleLayer.expandCard).not.toBe(
      enCommon.workbench.source.ruleLayer.expandCard,
    );
    expect(zhTWCommon.workbench.tag.board.flowModeHint.create).not.toBe(
      enCommon.workbench.tag.board.flowModeHint.create,
    );
    expect(zhTWCommon.workbench.tag.board.selectionHint.ready).not.toBe(
      enCommon.workbench.tag.board.selectionHint.ready,
    );
    expect(zhTWCommon.workbench.output.mapping.selectedTag).not.toBe(
      enCommon.workbench.output.mapping.selectedTag,
    );
    expect(zhTWCommon.workbench.output.database.mapping.selectedTag).not.toBe(
      enCommon.workbench.output.database.mapping.selectedTag,
    );
    expect(zhTWCommon.workbench.output.database.actions.configureConnector).not.toBe(
      enCommon.workbench.output.database.actions.configureConnector,
    );
  });

  it('labels lower-level point and tag tooling as secondary engineering surfaces', () => {
    expect(zhTWCommon.workbench.source.toolbar.moreMenu).toContain('工程');
    expect(enCommon.workbench.source.toolbar.moreMenu).toContain('Engineering');
    expect(zhTWCommon.workbench.source.toolbar.moreMenuAria).toContain('工程');
    expect(enCommon.workbench.source.toolbar.moreMenuAria).toContain('engineering');
    expect(zhTWCommon.workbench.source.actions.createSelectedPoints).toContain('工程');
    expect(enCommon.workbench.source.actions.createSelectedPoints).toContain('Engineering');
    expect(zhTWCommon.workbench.source.selectionToolbar.createPoints).toContain('工程');
    expect(enCommon.workbench.source.selectionToolbar.createPoints).toContain('Engineering');
    expect(zhTWCommon.workbench.tag.exception.title).toContain('工程工具');
    expect(enCommon.workbench.tag.exception.title).toContain('engineering');
    expect(zhTWCommon.workbench.tag.master.eyebrow).toContain('工程');
    expect(enCommon.workbench.tag.master.eyebrow).toContain('engineering');
  });

  it('keeps shell and navigation copy on studio review-first language', () => {
    expect(zhTWCommon.workbench.header.ariaLabel).toContain('Studio');
    expect(enCommon.workbench.header.ariaLabel).toContain('Studio');
    expect(zhTWCommon.workbench.stepRail.ariaLabel).toContain('Studio');
    expect(enCommon.workbench.stepRail.ariaLabel).toContain('Studio');
    expect(zhTWCommon.workbench.contextBar.ariaLabel).toContain('Studio');
    expect(enCommon.workbench.contextBar.ariaLabel).toContain('Studio');
    expect(zhTWCommon.workbench.actionDock.ariaLabel).toContain('Studio');
    expect(enCommon.workbench.actionDock.ariaLabel).toContain('Studio');

    expect(zhTWCommon.workbench.steps.tag).toContain('審查');
    expect(enCommon.workbench.steps.tag).toContain('Review');
    expect(zhTWCommon.workbench.steps.output).toContain('審查');
    expect(enCommon.workbench.steps.output).toContain('Review');
    expect(zhTWCommon.workbench.header.actions.gotoTag).toContain('審查');
    expect(enCommon.workbench.header.actions.gotoTag).toContain('Review');
    expect(zhTWCommon.workbench.contextBar.actions.gotoTag).toContain('審查');
    expect(enCommon.workbench.contextBar.actions.gotoTag).toContain('Review');
    expect(zhTWCommon.workbench.actionDock.readiness.tag).toContain('審查');
    expect(enCommon.workbench.actionDock.readiness.tag).toContain('Review');
    expect(zhTWCommon.workbench.actionDock.readiness.output).toContain('審查');
    expect(enCommon.workbench.actionDock.readiness.output).toContain('Review');
    expect(zhTWCommon.workbench.actionDock.nextAction.createPoints).toContain('規則');
    expect(enCommon.workbench.actionDock.nextAction.createPoints).toContain('rule');
    expect(zhTWCommon.workbench.actionDock.nextAction.bindTags).toContain('審查');
    expect(enCommon.workbench.actionDock.nextAction.bindTags).toContain('Review');
    expect(zhTWCommon.workbench.actionDock.nextAction.advanceToOutput).not.toContain('綁定');
    expect(enCommon.workbench.actionDock.nextAction.advanceToOutput).not.toContain('binding');
  });

  it('ships localized zh-TW strings for Step 3 empty-state guidance additions', () => {
    expect(zhTWCommon.workbench.tag.empty.eligibleSpans).not.toBe(
      enCommon.workbench.tag.empty.eligibleSpans,
    );
    expect(zhTWCommon.workbench.tag.empty.eligibleSpansNone).not.toBe(
      enCommon.workbench.tag.empty.eligibleSpansNone,
    );
    expect(zhTWCommon.workbench.tag.empty.goToSource).not.toBe(
      enCommon.workbench.tag.empty.goToSource,
    );
  });
});
