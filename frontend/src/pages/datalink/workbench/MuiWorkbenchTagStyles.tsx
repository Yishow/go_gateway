import GlobalStyles from '@mui/material/GlobalStyles';
import { alpha } from '@mui/material/styles';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';

const sentry = tokens.archetype.sentry;
const variantRoot = '[data-variant="v2-mui"]';

/* ── Review surface — the primary Tag panel ── */
const reviewSurface = `${variantRoot} [data-testid="source-rule-tag-review-surface"]`;
const reviewRows = `${variantRoot} [data-testid^="source-rule-tag-review-row-"]`;
const reviewSelects = [
  `${variantRoot} [data-testid="source-rule-tag-review-rule-select"]`,
  `${variantRoot} [data-testid^="source-rule-tag-review-override-select-"]`,
].join(', ');
const reviewInputs = `${variantRoot} [data-testid^="source-rule-tag-review-rename-input-"]`;
const reviewButtons = [
  `${variantRoot} [data-testid^="source-rule-tag-review-rename-save-"]`,
  `${variantRoot} [data-testid^="source-rule-tag-review-skip-"]`,
  `${variantRoot} [data-testid^="source-rule-tag-review-override-save-"]`,
  `${variantRoot} [data-testid="source-rule-tag-review-refresh"]`,
  `${variantRoot} [data-testid="source-rule-tag-review-retry"]`,
].join(', ');
const reviewChips = [
  `${variantRoot} [data-testid^="source-rule-tag-review-status-"]`,
  `${variantRoot} [data-testid^="source-rule-tag-review-mapping-intent-"]`,
  `${variantRoot} [data-testid^="source-rule-tag-review-decision-"]`,
  `${variantRoot} [data-testid="source-rule-tag-review-set-status"]`,
].join(', ');

/* ── Board-level panels (candidate board, master list, batch diff, etc.) ── */
const boardPanels = [
  `${variantRoot} [data-testid="tag-board-surface"]`,
  `${variantRoot} [data-testid="tag-candidate-board"]`,
  `${variantRoot} [data-testid="tag-master-surface"]`,
  `${variantRoot} [data-testid="batch-diff-preview"]`,
  `${variantRoot} [data-testid="tag-exception-tools"]`,
].join(', ');

const summaryCards = [
  `${variantRoot} [data-testid="tag-review-summary"]`,
  `${variantRoot} [data-testid="tag-review-generated"]`,
  `${variantRoot} [data-testid="tag-review-needs-review"]`,
  `${variantRoot} [data-testid="tag-review-selected-exceptions"]`,
].join(', ');

const statusPanels = [
  `${variantRoot} [data-testid="source-rule-tag-review-loading"]`,
  `${variantRoot} [data-testid="source-rule-tag-review-error"]`,
  `${variantRoot} [data-testid="source-rule-tag-review-decision-error"]`,
  `${variantRoot} [data-testid="source-rule-tag-review-stale"]`,
  `${variantRoot} [data-testid="source-rule-tag-review-empty"]`,
].join(', ');

const feedbackBanner = `${variantRoot} [data-testid="source-rule-tag-review-feedback"]`;

/**
 * Sentry Incident Desk GlobalStyles for the **Tag** step.
 *
 * Follows the same CSS-override pattern as {@link MuiWorkbenchSourceStyles}:
 * targets `data-testid` selectors inside `[data-variant="v2-mui"]` so the
 * existing Tailwind components render with the warm purple-black command
 * center aesthetic without structural changes.
 */
export function MuiWorkbenchTagStyles() {
  return (
    <GlobalStyles
      styles={{
        /* ── Main review surface ── */
        [reviewSurface]: {
          fontFamily: tokens.typography.family.sentryUi,
          border: `1px solid ${sentry.border}`,
          borderRadius: '14px',
          background: `linear-gradient(180deg, ${alpha(sentry.elevated, 0.98)} 0%, ${alpha(
            sentry.panel,
            0.96,
          )} 100%)`,
          boxShadow: tokens.treatment.ambientPurple,
        },

        /* ── Candidate rows — frosted card with left accent ── */
        [reviewRows]: {
          fontFamily: tokens.typography.family.sentryUi,
          borderRadius: '14px',
          borderColor: sentry.border,
          backgroundColor: alpha(sentry.elevated, 0.92),
          boxShadow: 'rgba(0, 0, 0, 0.1) 0px 10px 15px -3px',
          borderLeft: `3px solid ${sentry.accent}`,
        },

        /* ── Selects & inputs — muted glass ── */
        [reviewSelects]: {
          borderRadius: '10px',
          borderColor: sentry.border,
          backgroundColor: alpha(sentry.panel, 0.92),
          color: tokens.text.primary,
          fontFamily: 'Monaco, Menlo, Ubuntu Mono, monospace',
        },
        [reviewInputs]: {
          borderRadius: '10px',
          borderColor: sentry.border,
          backgroundColor: alpha(sentry.panel, 0.92),
          color: tokens.text.primary,
          fontFamily: 'Monaco, Menlo, Ubuntu Mono, monospace',
        },

        /* ── Action buttons — inset Sentry style ── */
        [reviewButtons]: {
          borderRadius: '13px',
          border: `1px solid #584674`,
          backgroundColor: sentry.accentMuted,
          color: tokens.text.primary,
          fontFamily: tokens.typography.family.sentryUi,
          fontSize: '13px',
          fontWeight: 700,
          letterSpacing: '0.2px',
          textTransform: 'uppercase',
          boxShadow: tokens.treatment.insetButton,
        },

        /* ── Chips / badges — Sentry mono typography ── */
        [reviewChips]: {
          fontFamily: tokens.typography.family.sentryUi,
          fontWeight: 700,
          letterSpacing: '0.2px',
          textTransform: 'uppercase',
          fontSize: '11px',
        },

        /* ── Board panels ── */
        [boardPanels]: {
          fontFamily: tokens.typography.family.sentryUi,
          border: `1px solid ${sentry.border}`,
          borderRadius: '14px',
          background: `linear-gradient(180deg, ${alpha(sentry.elevated, 0.98)} 0%, ${alpha(
            sentry.panel,
            0.96,
          )} 100%)`,
          boxShadow: tokens.treatment.ambientPurple,
        },

        /* ── Summary cards — highlight gradient ── */
        [summaryCards]: {
          borderRadius: '14px',
          border: `1px solid ${alpha(sentry.highlight, 0.18)}`,
          backgroundColor: alpha(sentry.panel, 0.94),
          boxShadow: 'rgba(0, 0, 0, 0.1) 0px 10px 15px -3px',
        },

        /* ── Status / error / stale / empty panels ── */
        [statusPanels]: {
          borderRadius: '14px',
          fontFamily: tokens.typography.family.sentryUi,
        },
        [`${variantRoot} [data-testid="source-rule-tag-review-stale"]`]: {
          borderColor: `${sentry.warm}66`,
          background: `linear-gradient(135deg, ${alpha(sentry.warm, 0.16)} 0%, ${alpha(
            sentry.panel,
            0.96,
          )} 72%)`,
        },
        [`${variantRoot} [data-testid="source-rule-tag-review-error"]`]: {
          borderColor: `${tokens.status.error}55`,
          background: `linear-gradient(135deg, ${alpha(tokens.status.error, 0.16)} 0%, ${alpha(
            sentry.panel,
            0.96,
          )} 72%)`,
        },
        [`${variantRoot} [data-testid="source-rule-tag-review-loading"]`]: {
          background: `linear-gradient(135deg, ${alpha(sentry.highlight, 0.16)} 0%, ${alpha(
            sentry.panel,
            0.96,
          )} 72%)`,
        },

        /* ── Feedback toast — accent tint ── */
        [feedbackBanner]: {
          fontFamily: tokens.typography.family.sentryUi,
          borderRadius: '12px',
        },

        /* ── Eyebrow typography inside review surface ── */
        [`${reviewSurface} p[class*="tracking-"]`]: {
          color: sentry.highlight,
          fontFamily: tokens.typography.family.sentryUi,
        },

        /* ── dt labels ── */
        [`${reviewSurface} dt, ${reviewRows} dt`]: {
          fontFamily: tokens.typography.family.sentryUi,
          textTransform: 'uppercase',
          letterSpacing: '0.2px',
          fontWeight: 600,
          color: tokens.text.muted,
        },

        /* ── Metric dd values ── */
        [`${reviewSurface} dd`]: {
          fontFamily: 'Monaco, Menlo, Ubuntu Mono, monospace',
        },
      }}
    />
  );
}
