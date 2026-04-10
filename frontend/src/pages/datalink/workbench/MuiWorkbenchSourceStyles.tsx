import GlobalStyles from '@mui/material/GlobalStyles';
import { alpha } from '@mui/material/styles';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import type { SentryDeskMode } from './MuiSourceCommandDeck';

const sentry = tokens.archetype.sentry;
const V = '[data-variant="v2-mui"]';
const W = '[data-testid="sentry-source-workspace"]';

interface Props {
  deskMode: SentryDeskMode;
}

/** Base sentry visual treatment for all source panels. */
function basePanelStyles() {
  const sourcePanels = [
    `${V} [data-testid="source-step-summary"]`,
    `${V} [data-testid="source-rule-layer"]`,
    `${V} [data-testid="source-canvas-workspace"]`,
    `${V} [data-testid="source-conflict-queue"]`,
    `${V} [data-testid="source-runtime-collection-panel"]`,
    `${V} [data-testid="source-loading-state"]`,
    `${V} [data-testid="source-error-state"]`,
  ].join(', ');

  const toolbars = [
    `${V} [data-testid="source-primary-toolbar"]`,
    `${V} [data-testid="source-selection-toolbar"]`,
  ].join(', ');

  const fields = [
    `${V} [data-testid="source-rule-layer"] input`,
    `${V} [data-testid="source-rule-layer"] select`,
    `${V} [data-testid="source-canvas-workspace"] input`,
    `${V} [data-testid="source-canvas-workspace"] select`,
  ].join(', ');

  const buttons = [
    `${V} [data-testid="source-rule-layer"] button:not([data-testid^="address-cell-"])`,
    `${V} [data-testid="source-canvas-workspace"] button:not([data-testid^="address-cell-"])`,
    `${V} [data-testid="source-error-state"] button`,
  ].join(', ');

  return {
    [sourcePanels]: {
      fontFamily: tokens.typography.family.sentryUi,
      border: `1px solid ${sentry.border}`,
      borderRadius: '14px',
      background: `linear-gradient(180deg, ${alpha(sentry.elevated, 0.98)} 0%, ${alpha(sentry.panel, 0.96)} 100%)`,
      boxShadow: tokens.treatment.ambientPurple,
    },
    [toolbars]: {
      border: `1px solid ${alpha(sentry.highlight, 0.18)}`,
      borderRadius: '14px',
      background: `linear-gradient(90deg, ${alpha(sentry.accentMuted, 0.22)} 0%, ${alpha(sentry.panel, 0.95)} 100%)`,
      boxShadow: tokens.treatment.ambientPurple,
    },
    [fields]: {
      borderRadius: '10px',
      borderColor: sentry.border,
      backgroundColor: alpha(sentry.panel, 0.92),
      color: tokens.text.primary,
      fontFamily: 'Monaco, Menlo, Ubuntu Mono, monospace',
    },
    [buttons]: {
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
  } as const;
}

/** Element-specific sentry overrides. */
function elementStyles() {
  return {
    [`${V} [data-testid="source-step-summary"]`]: {
      background: `linear-gradient(135deg, ${alpha(sentry.highlight, 0.14)} 0%, ${alpha(sentry.panel, 0.98)} 72%)`,
    },
    [`${V} [data-testid="source-rule-layer-description-hint"]`]: { color: sentry.warm },
    [`${V} [data-testid="source-rule-layer-tab-planner"], ${V} [data-testid="source-rule-layer-tab-rules"]`]: {
      borderRadius: '12px',
      borderColor: sentry.border,
      backgroundColor: alpha(sentry.panel, 0.9),
      color: tokens.text.secondary,
      fontFamily: tokens.typography.family.sentryUi,
      letterSpacing: '0.2px',
      textTransform: 'uppercase',
    },
    [`${V} [data-testid="source-rule-layer-tab-planner"][aria-selected="true"], ${V} [data-testid="source-rule-layer-tab-rules"][aria-selected="true"]`]:
      {
        backgroundColor: alpha(sentry.highlight, 0.16),
        borderColor: alpha(sentry.highlight, 0.55),
        color: sentry.highlight,
        boxShadow: `0 0 0 1px ${alpha(sentry.highlight, 0.3)}`,
      },
    [`${V} [data-testid="source-rule-layer"] [data-testid$="-card"]`]: {
      borderRadius: '14px',
      borderColor: alpha(sentry.highlight, 0.18),
      backgroundColor: alpha(sentry.panel, 0.92),
      boxShadow: tokens.treatment.ambientPurple,
    },
    [`${V} [data-testid="source-memory-scroll-region"]`]: {
      borderRadius: '14px',
      border: `1px solid ${sentry.border}`,
      background: `linear-gradient(180deg, ${alpha(sentry.panel, 0.98)} 0%, ${alpha(sentry.canvas, 0.98)} 100%)`,
    },
    [`${V} [data-testid="source-canvas-row-0"], ${V} [data-testid^="source-canvas-row-"]`]: {
      borderRadius: '12px',
      border: `1px solid ${sentry.border}`,
      backgroundColor: alpha(sentry.elevated, 0.96),
    },
    [`${V} [data-testid^="address-cell-"]`]: {
      borderRadius: '10px',
      borderWidth: '1px',
      borderStyle: 'solid',
      fontFamily: 'Monaco, Menlo, Ubuntu Mono, monospace',
      boxShadow: 'none',
    },
    [`${V} [data-testid^="address-cell-"][data-status="planned"]`]: {
      borderColor: alpha(sentry.highlight, 0.4),
      backgroundColor: alpha(sentry.highlight, 0.12),
      color: tokens.text.primary,
    },
    [`${V} [data-testid^="address-cell-"][data-status="used"]`]: {
      borderColor: alpha(sentry.accent, 0.5),
      backgroundColor: alpha(sentry.accent, 0.14),
      color: tokens.text.primary,
    },
    [`${V} [data-testid^="address-cell-"][data-status="unmanaged"]`]: {
      borderColor: alpha(sentry.warm, 0.44),
      backgroundColor: alpha(sentry.warm, 0.12),
      color: '#fff1e9',
    },
    [`${V} [data-testid^="address-cell-"][data-status="conflict"]`]: {
      borderColor: alpha(tokens.status.error, 0.55),
      backgroundColor: alpha(tokens.status.error, 0.16),
      color: '#ffe2e2',
    },
    [`${V} [data-testid^="address-cell-"][data-status="gap"]`]: {
      borderColor: sentry.border,
      backgroundColor: alpha(sentry.panel, 0.9),
      color: tokens.text.muted,
    },
    [`${V} [data-testid="source-coverage-overview"]`]: {
      borderRadius: '14px',
      border: `1px solid ${sentry.border}`,
      backgroundColor: alpha(sentry.panel, 0.94),
    },
    [`${V} [data-testid="source-conflict-queue"]`]: {
      borderColor: alpha(tokens.status.error, 0.32),
      background: `linear-gradient(180deg, ${alpha(tokens.status.error, 0.08)} 0%, ${alpha(sentry.panel, 0.98)} 100%)`,
    },
  } as const;
}

/** Loading/error status panel styles. */
function statusStyles() {
  const panels = [
    `${V} [data-testid="source-loading-state"]`,
    `${V} [data-testid="source-error-state"]`,
  ].join(', ');

  return {
    [panels]: {
      minHeight: '100%',
      display: 'flex',
      flexDirection: 'column' as const,
      justifyContent: 'center',
      gap: '12px',
      padding: '24px',
    },
    [`${V} [data-testid="source-loading-state"]`]: {
      background: `linear-gradient(135deg, ${alpha(sentry.highlight, 0.16)} 0%, ${alpha(sentry.panel, 0.96)} 72%)`,
      color: tokens.text.primary,
    },
    [`${V} [data-testid="source-error-state"]`]: {
      background: `linear-gradient(135deg, ${alpha(tokens.status.error, 0.16)} 0%, ${alpha(sentry.panel, 0.96)} 72%)`,
      color: '#ffe2e2',
    },
  } as const;
}

/**
 * Mode-driven layout restructuring.
 *
 * Overrides the SourceCanvasSection grid layout based on the active desk mode.
 * Uses flex ordering and proportions rather than display:none to keep all
 * panels DOM-accessible (required by tests and accessibility).
 *
 * - inspect: full-width canvas primary, aside compact below
 * - build:   full-width rule builder primary, canvas compact below
 * - triage:  conflict queue highlighted, canvas primary, aside compact
 */
function modeLayoutStyles() {
  /* The workspace wrapper contains SourceCanvasSection's root:
     section > div(grid) > aside(rules) + section(canvas) */
  const gridRoot = `${W} > section > div`;
  const ruleAside = `> section > div > aside`;
  const canvasPanel = `[data-testid="source-canvas-workspace"]`;
  const conflictQueue = `[data-testid="source-conflict-queue"]`;

  return {
    /* ── ALWAYS ── single-column stacked flex (replaces baseline 2-col grid) */
    [gridRoot]: {
      display: 'flex !important',
      flexDirection: 'column !important' as 'column',
      gap: '8px !important',
    },

    /* ── INSPECT MODE (default) ── canvas is primary, aside compacted below */
    [`${W}[data-desk-mode="inspect"] ${canvasPanel}`]: {
      flex: '1 1 auto !important',
      order: '-1 !important',
      borderColor: `${alpha(sentry.highlight, 0.4)} !important`,
    },
    [`${W}[data-desk-mode="inspect"] ${ruleAside}`]: {
      flex: '0 0 auto !important',
      maxHeight: '220px !important',
      overflow: 'auto !important',
      opacity: '0.72',
      borderTop: `1px solid ${sentry.border}`,
    },

    /* ── BUILD MODE ── rule builder is primary, canvas compacted below */
    [`${W}[data-desk-mode="build"] ${ruleAside}`]: {
      flex: '1 1 auto !important',
      order: '-1 !important',
    },
    [`${W}[data-desk-mode="build"] ${canvasPanel}`]: {
      flex: '0 0 auto !important',
      maxHeight: '180px !important',
      overflow: 'hidden !important',
      opacity: '0.65',
      borderTop: `1px solid ${sentry.border}`,
    },

    /* ── TRIAGE MODE ── conflict queue elevated, aside compacted */
    [`${W}[data-desk-mode="triage"] ${conflictQueue}`]: {
      order: '-2 !important',
      flex: '0 0 auto !important',
      minHeight: '160px !important',
      borderColor: `${alpha(tokens.status.error, 0.5)} !important`,
      boxShadow: `0 0 20px ${alpha(tokens.status.error, 0.15)}`,
    },
    [`${W}[data-desk-mode="triage"] ${ruleAside}`]: {
      flex: '0 0 auto !important',
      maxHeight: '120px !important',
      overflow: 'auto !important',
      opacity: '0.55',
    },
    [`${W}[data-desk-mode="triage"] ${canvasPanel}`]: {
      flex: '1 1 auto !important',
      order: '-1 !important',
    },
  } as Record<string, Record<string, string>>;
}

export function MuiWorkbenchSourceStyles({ deskMode: _deskMode }: Props) {
  return (
    <GlobalStyles
      styles={{
        ...basePanelStyles(),
        ...elementStyles(),
        ...statusStyles(),
        ...modeLayoutStyles(),
      }}
    />
  );
}
