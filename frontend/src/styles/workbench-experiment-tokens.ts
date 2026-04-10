const tokenSources = ['linear.app', 'sentry', 'clickhouse'] as const;

export const WORKBENCH_EXPERIMENT_TOKEN_SOURCES = tokenSources;

/**
 * Shared semantic tokens for the workbench operator-efficiency experiment.
 *
 * Phase 1R widened the contract so v1/v2/v3 can express distinct archetypes
 * without inventing branch-only token files or version-specific palettes.
 */
export const workbenchExperimentTokens = {
  meta: {
    sources: tokenSources,
    audience: 'industrial-operator',
    intent: 'low-fatigue-high-density-ui',
    phase1r: 'three-archetypes-reboot',
  },
  surface: {
    app: '#0b1117',
    panel: '#111922',
    elevated: '#16212c',
    overlay: 'rgba(8, 12, 17, 0.82)',
  },
  text: {
    primary: '#edf3ff',
    secondary: '#b7c1d6',
    muted: '#7f8aa3',
    inverse: '#091018',
  },
  accent: {
    primary: '#4c8dff',
    subtle: 'rgba(76, 141, 255, 0.14)',
    strong: '#7bb0ff',
  },
  status: {
    success: '#32b36b',
    warning: '#e0a13a',
    error: '#e15659',
    info: '#4c8dff',
  },
  border: {
    default: 'rgba(147, 167, 194, 0.18)',
    muted: 'rgba(147, 167, 194, 0.1)',
    strong: 'rgba(147, 167, 194, 0.34)',
  },
  focus: {
    ring: '#8ab4ff',
    offset: '#0b1117',
    invalid: '#ff8b8d',
  },
  density: {
    compact: {
      controlHeight: 32,
      sectionGap: 12,
      rowGap: 4,
    },
    default: {
      controlHeight: 40,
      sectionGap: 16,
      rowGap: 8,
    },
    relaxed: {
      controlHeight: 48,
      sectionGap: 24,
      rowGap: 12,
    },
  },
  data: {
    rowHover: 'rgba(90, 134, 213, 0.12)',
    rowSelected: 'rgba(76, 141, 255, 0.2)',
    metricAccent: '#ff9b2f',
  },
  archetype: {
    linear: {
      canvas: '#08090a',
      panel: '#0f1011',
      elevated: '#191a1b',
      accent: '#7170ff',
      accentHover: '#828fff',
      border: 'rgba(255, 255, 255, 0.08)',
      labelTracking: '-0.02em',
      uiWeight: 510,
    },
    sentry: {
      canvas: '#1f1633',
      panel: '#150f23',
      elevated: '#241a3b',
      accent: '#6a5fc1',
      accentMuted: '#79628c',
      highlight: '#c2ef4e',
      warm: '#ffb287',
      border: '#362d59',
      labelTracking: '0.02em',
      uiFont: 'Rubik, Inter, system-ui, sans-serif',
    },
    clickhouse: {
      canvas: '#000000',
      panel: '#141414',
      elevated: '#1b1b1b',
      accent: '#faff69',
      accentPressed: '#f4f692',
      accentSecondary: '#166534',
      border: 'rgba(65, 65, 65, 0.8)',
      labelTracking: '0.14em',
      uiWeight: 700,
      displayWeight: 900,
    },
  },
  typography: {
    family: {
      sans: 'Inter, system-ui, sans-serif',
      mono: 'JetBrains Mono, monospace',
      linearMono: 'Berkeley Mono, JetBrains Mono, monospace',
      sentryUi: 'Rubik, Inter, system-ui, sans-serif',
      code: 'Inconsolata, JetBrains Mono, monospace',
    },
    size: {
      xs: '12px',
      sm: '14px',
      md: '16px',
      lg: '18px',
      xl: '24px',
      '2xl': '32px',
    },
    lineHeight: {
      compact: 1.3,
      default: 1.5,
      relaxed: 1.65,
    },
    weight: {
      regular: 400,
      linearUi: 510,
      medium: 600,
      strong: 700,
      heavy: 900,
    },
    tracking: {
      tight: '-0.03em',
      normal: '0em',
      label: '0.02em',
      overline: '0.14em',
    },
  },
  treatment: {
    glassPanel: {
      fill: 'rgba(255, 255, 255, 0.18)',
      hoverFill: 'rgba(54, 22, 107, 0.14)',
      backdropFilter: 'blur(18px) saturate(180%)',
    },
    insetButton: 'inset 0 1px 3px rgba(0, 0, 0, 0.35), 0 1px 0 rgba(255, 255, 255, 0.04)',
    ambientPurple: 'rgba(22, 15, 36, 0.9) 0px 4px 4px 9px',
    neonGlow: '0 0 12px rgba(250, 255, 105, 0.45)',
    cockpitInset: '0px 4px 25px rgba(0, 0, 0, 0.14) inset',
  },
  spacing: {
    xs: '4px',
    sm: '8px',
    md: '16px',
    lg: '24px',
    xl: '32px',
    '2xl': '48px',
  },
  radius: {
    sm: '6px',
    md: '10px',
    lg: '14px',
    xl: '18px',
    pill: '9999px',
  },
  elevation: {
    panel: '0 12px 30px -18px rgba(0, 0, 0, 0.55)',
    overlay: '0 24px 60px -28px rgba(0, 0, 0, 0.7)',
    focus: '0 0 0 3px rgba(138, 180, 255, 0.26)',
  },
} as const;

/**
 * Mapping boundaries are part of the contract: variants may change component
 * composition, but they may not fork token family names or values.
 */
export const WORKBENCH_EXPERIMENT_TOKEN_MAPPING_BOUNDARIES = {
  v1: 'Map shared semantics through shadcn/Radix primitives and Tailwind classes. May consume archetype.linear and shared typography/treatment semantics, but may not add branch-only token values.',
  v2: 'Map shared semantics through MUI theme tokens and component overrides. May consume archetype.sentry and shared typography/treatment semantics, but may not add branch-only token values.',
  v3: 'Map shared semantics through Ant Design theme tokens and component overrides. May consume archetype.clickhouse and shared typography/treatment semantics, but may not add branch-only token values.',
} as const;

export type WorkbenchExperimentTokens = typeof workbenchExperimentTokens;
