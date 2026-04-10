import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import type { SxProps, Theme } from '@mui/material/styles';

const sentry = tokens.archetype.sentry;

/* ── Sentry design-doc values — v2 mapping layer only ── */
const BTN_BORDER = '#584674';
const BTN_RADIUS = '13px';
const HOVER_SHADOW = 'rgba(0, 0, 0, 0.18) 0px 0.5rem 1.5rem';
const CARD_SHADOW = 'rgba(0, 0, 0, 0.1) 0px 10px 15px -3px';
const MONO = 'Monaco, Menlo, Ubuntu Mono, monospace';

/**
 * Sentry Incident Desk visual language — sx presets from shared tokens
 * plus supplementary Sentry DESIGN.md values (consumed only in v2 layer).
 */
export const SENTRY_SX = {
  /** Uppercase section label — Sentry's signature typography pattern */
  sectionLabel: {
    fontFamily: tokens.typography.family.sentryUi,
    textTransform: 'uppercase',
    letterSpacing: '0.2px',
    fontWeight: 600,
    fontSize: '11px',
    color: tokens.text.muted,
    lineHeight: 1.25,
  } satisfies SxProps<Theme>,

  /** Tactile muted-purple button — Sentry's signature 13px-radius inset button */
  insetBtn: {
    boxShadow: tokens.treatment.insetButton,
    textTransform: 'uppercase',
    letterSpacing: '0.2px',
    fontFamily: tokens.typography.family.sentryUi,
    fontWeight: 700,
    fontSize: '14px',
    bgcolor: sentry.accentMuted,
    color: tokens.text.primary,
    border: `1px solid ${BTN_BORDER}`,
    borderRadius: BTN_RADIUS,
    px: 2,
    py: 0.75,
    minHeight: 34,
    cursor: 'pointer',
    transition: 'all 0.15s ease',
    '&:hover': {
      bgcolor: sentry.accent,
      borderColor: sentry.accent,
      boxShadow: HOVER_SHADOW,
    },
  } satisfies SxProps<Theme>,

  /** Inset icon button variant */
  insetIconBtn: {
    boxShadow: tokens.treatment.insetButton,
    bgcolor: sentry.elevated,
    border: `1px solid ${sentry.border}`,
    borderRadius: '8px',
    width: 32,
    height: 32,
    color: tokens.text.secondary,
    transition: 'all 0.15s ease',
    '&:hover': { bgcolor: sentry.accentMuted, color: tokens.text.primary, boxShadow: HOVER_SHADOW },
  } satisfies SxProps<Theme>,

  /** Frosted glass panel — semi-transparent bg + backdrop blur */
  frostedPanel: {
    bgcolor: tokens.treatment.glassPanel.fill,
    backdropFilter: tokens.treatment.glassPanel.backdropFilter,
    border: `1px solid ${sentry.border}`,
    borderRadius: '12px',
    boxShadow: CARD_SHADOW,
  } satisfies SxProps<Theme>,

  /** Glass command strip — top/bottom bars with subtle glass and purple tint */
  commandStrip: {
    bgcolor: sentry.elevated,
    backdropFilter: tokens.treatment.glassPanel.backdropFilter,
    borderRadius: '10px',
    border: `1px solid ${sentry.border}`,
    boxShadow: tokens.treatment.ambientPurple,
  } satisfies SxProps<Theme>,

  /** Mode chip — uppercase, bold, letter-spaced */
  modeChip: {
    fontFamily: tokens.typography.family.sentryUi,
    fontWeight: 700,
    letterSpacing: '0.2px',
    textTransform: 'uppercase',
    fontSize: '11px',
  } satisfies SxProps<Theme>,

  /** Severity bar — colored left border on incident cards */
  severityBar: (color: string): SxProps<Theme> => ({
    borderLeft: `3px solid ${color}`,
    borderRadius: `0 ${tokens.radius.sm} ${tokens.radius.sm} 0`,
  }),

  /** Oversized monospace metric value — Sentry data display */
  metricValue: {
    fontFamily: MONO,
    fontWeight: tokens.typography.weight.medium,
    fontSize: tokens.typography.size.xl,
    lineHeight: 1.1,
  } satisfies SxProps<Theme>,

  /** Small monospace data label */
  monoData: {
    fontFamily: MONO,
    fontSize: '12px',
    color: tokens.text.secondary,
  } satisfies SxProps<Theme>,

  /** Stage tile — large diagnostic stage card */
  stageTile: {
    p: 2,
    bgcolor: sentry.elevated,
    border: `1px solid ${sentry.border}`,
    borderRadius: '12px',
    boxShadow: CARD_SHADOW,
    display: 'flex',
    flexDirection: 'column',
    gap: 1,
    position: 'relative',
    overflow: 'hidden',
    '&::before': {
      content: '""',
      position: 'absolute',
      inset: 0,
      background: `linear-gradient(135deg, ${sentry.accent}14 0%, transparent 60%)`,
      pointerEvents: 'none',
    },
  } satisfies SxProps<Theme>,
} as const;

/** Map device health to severity color */
export function severityColor(lastTestSuccess: boolean | null): string {
  if (lastTestSuccess === true) return sentry.highlight;
  if (lastTestSuccess === false) return sentry.warm;
  return sentry.border;
}
