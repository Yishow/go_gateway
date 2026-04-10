import { createTheme } from '@mui/material/styles';
import { workbenchExperimentTokens as tokens } from './workbench-experiment-tokens';

const s = tokens.archetype.sentry;

/* ── Sentry design-doc values — v2 mapping layer only ── */
const BTN_BORDER = '#584674';
const BTN_RADIUS = 13;
const HOVER_SHADOW = 'rgba(0, 0, 0, 0.18) 0px 0.5rem 1.5rem';
const CARD_SHADOW = 'rgba(0, 0, 0, 0.1) 0px 10px 15px -3px';
const DEEP_VIOLET = '#422082';
const CORAL_FOCUS = '#ffb287';
const MONO = 'Monaco, Menlo, Ubuntu Mono, monospace';

/**
 * Sentry Incident Desk MUI theme — v2 archetype.
 *
 * All surfaces warm purple-black (never pure black). Typography is Rubik
 * with systematic uppercase + 0.2px letter-spacing on labels/buttons.
 * Primary buttons: muted purple (#79628c), tactile inset, 13px radius.
 * Hover: elevated shadow. Focus: purple ring.
 */
export const muiExperimentTheme = createTheme({
  palette: {
    mode: 'dark',
    background: {
      default: s.canvas,
      paper: s.panel,
    },
    primary: {
      main: s.accent,
      light: s.highlight,
      dark: s.accentMuted,
    },
    secondary: {
      main: s.warm,
    },
    success: { main: s.highlight },
    warning: { main: s.warm },
    error: { main: tokens.status.error },
    info: { main: s.accent },
    text: {
      primary: tokens.text.primary,
      secondary: tokens.text.secondary,
      disabled: tokens.text.muted,
    },
    divider: s.border,
    action: {
      hover: `${s.accent}18`,
      selected: `${s.accent}28`,
      focus: `${s.accent}22`,
    },
  },

  typography: {
    fontFamily: s.uiFont,
    fontSize: 14,
    h5: { fontWeight: 600, fontSize: '20px', letterSpacing: '0.2px', lineHeight: 1.25 },
    h6: { fontWeight: 600, fontSize: '16px', letterSpacing: '0.2px', lineHeight: 1.25 },
    subtitle1: { fontSize: '15px', color: tokens.text.secondary, fontWeight: 500, lineHeight: 1.4 },
    subtitle2: { fontSize: '12px', color: tokens.text.muted, textTransform: 'uppercase' as const, letterSpacing: '0.2px', fontWeight: 600 },
    body1: { fontSize: '14px', lineHeight: 1.5 },
    body2: { fontSize: '12px', lineHeight: 1.5 },
    caption: { fontSize: '12px', color: tokens.text.muted, letterSpacing: '0.2px', fontWeight: 600 },
    button: { textTransform: 'uppercase' as const, fontWeight: 700, letterSpacing: '0.2px', fontSize: '14px' },
  },

  shape: { borderRadius: 10 },
  spacing: (factor: number) => `${factor * 8}px`,

  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          backgroundColor: s.canvas,
          color: tokens.text.primary,
          fontFamily: s.uiFont,
        },
      },
    },
    MuiPaper: {
      defaultProps: { elevation: 0 },
      styleOverrides: {
        root: {
          backgroundImage: 'none',
          backgroundColor: s.panel,
          border: `1px solid ${s.border}`,
          borderRadius: 10,
          boxShadow: tokens.treatment.ambientPurple,
        },
      },
    },
    MuiButton: {
      defaultProps: { disableElevation: true },
      styleOverrides: {
        root: {
          borderRadius: BTN_RADIUS,
          fontSize: '14px',
          fontWeight: 700,
          letterSpacing: '0.2px',
          textTransform: 'uppercase' as const,
          minHeight: 34,
          boxShadow: tokens.treatment.insetButton,
          border: `1px solid ${BTN_BORDER}`,
          transition: 'all 0.15s ease',
          '&:hover': { boxShadow: HOVER_SHADOW },
          '&:focus-visible': {
            outline: `${s.accent} solid 0.125rem`,
            outlineOffset: '2px',
            backgroundColor: `${CORAL_FOCUS}18`,
          },
        },
        containedPrimary: {
          backgroundColor: s.accentMuted,
          color: '#fff',
          '&:hover': { backgroundColor: s.accent, boxShadow: HOVER_SHADOW },
        },
        outlined: {
          borderColor: s.border,
          color: tokens.text.secondary,
          '&:hover': { borderColor: s.accent, color: tokens.text.primary },
        },
      },
    },
    MuiIconButton: {
      styleOverrides: {
        root: {
          transition: 'all 0.15s ease',
          '&:focus-visible': {
            outline: `${s.accent} solid 0.125rem`,
            outlineOffset: '2px',
            backgroundColor: `${CORAL_FOCUS}18`,
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          height: 24,
          fontSize: '11px',
          fontWeight: 700,
          letterSpacing: '0.2px',
          textTransform: 'uppercase' as const,
          fontFamily: s.uiFont,
        },
      },
    },
    MuiTab: {
      styleOverrides: {
        root: {
          textTransform: 'uppercase' as const,
          fontWeight: 700,
          fontSize: '13px',
          letterSpacing: '0.2px',
          fontFamily: s.uiFont,
          minHeight: 36,
          color: tokens.text.muted,
          '&.Mui-selected': { color: s.highlight },
        },
      },
    },
    MuiTabs: {
      styleOverrides: {
        indicator: { backgroundColor: s.highlight, height: 2 },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            fontFamily: MONO,
            fontSize: '13px',
            borderRadius: 6,
            '& fieldset': { borderColor: s.border },
            '&:hover fieldset': { borderColor: s.accentMuted },
            '&.Mui-focused fieldset': { borderColor: s.accent },
          },
          '& .MuiInputLabel-root': {
            fontFamily: s.uiFont,
            fontSize: '12px',
            fontWeight: 600,
            letterSpacing: '0.2px',
            textTransform: 'uppercase' as const,
          },
        },
      },
    },
    MuiAlert: {
      styleOverrides: {
        root: {
          fontFamily: s.uiFont,
          fontSize: '12px',
          borderRadius: 8,
        },
      },
    },
    MuiSelect: {
      styleOverrides: {
        root: {
          '& .MuiSelect-select': {
            backgroundColor: DEEP_VIOLET,
            borderRadius: 8,
          },
        },
      },
    },
    MuiMenu: {
      styleOverrides: {
        paper: {
          backgroundColor: DEEP_VIOLET,
          border: `1px solid ${s.border}`,
          borderRadius: 8,
          boxShadow: CARD_SHADOW,
        },
      },
    },
    MuiTooltip: {
      styleOverrides: {
        tooltip: {
          backgroundColor: s.elevated,
          border: `1px solid ${s.border}`,
          color: tokens.text.primary,
          fontSize: '12px',
          fontFamily: s.uiFont,
          borderRadius: 8,
          boxShadow: tokens.treatment.ambientPurple,
        },
      },
    },
  },
});
