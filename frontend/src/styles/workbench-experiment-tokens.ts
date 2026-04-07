const tokenSources = ['linear.app', 'sentry', 'clickhouse'] as const;

export const WORKBENCH_EXPERIMENT_TOKEN_SOURCES = tokenSources;

/**
 * Shared semantic tokens for the workbench operator-efficiency experiment.
 *
 * The values are intentionally kit-agnostic so v1/v2/v3 can map the same
 * semantics into shadcn/Radix, MUI, and Ant Design without drifting.
 */
export const workbenchExperimentTokens = {
  meta: {
    sources: tokenSources,
    audience: 'industrial-operator',
    intent: 'low-fatigue-high-density-ui',
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
  typography: {
    family: {
      sans: 'Inter, system-ui, sans-serif',
      mono: 'JetBrains Mono, monospace',
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
  v1: 'Map shared semantics through shadcn/Radix primitives and Tailwind classes.',
  v2: 'Map shared semantics through MUI theme tokens and component overrides.',
  v3: 'Map shared semantics through Ant Design theme tokens and component overrides.',
} as const;

export type WorkbenchExperimentTokens = typeof workbenchExperimentTokens;
