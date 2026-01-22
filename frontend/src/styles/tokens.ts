export const tokens = {
  colors: {
    // 背景
    bgPrimary: "#0F172A", // Slate-900
    bgSecondary: "#1E293B", // Slate-800
    bgTertiary: "#334155", // Slate-700
    bgHover: "rgba(63, 79, 106, 0.7)", // Slate-600/70 roughly

    // 文字
    textPrimary: "#F1F5F9", // Slate-100
    textSecondary: "#94A3B8", // Slate-400
    textMuted: "#64748B", // Slate-500

    // 強調色
    primary: "#3B82F6", // Blue-500
    primaryHover: "#2563EB", // Blue-600
    primaryLight: "rgba(59, 130, 246, 0.12)", // Blue-500/12

    // 狀態
    success: "#22C55E", // Green-500
    warning: "#EAB308", // Yellow-500
    error: "#EF4444", // Red-500
    info: "#3B82F6", // Blue-500

    // 邊框
    borderDefault: "#334155", // Slate-700
    borderHover: "#475569", // Slate-600

    // 格子狀態
    cellUsed: "#22C55E", // Green-500
    cellAvailable: "#334155", // Slate-700
    cellSelected: "#3B82F6", // Blue-500
    cellPending: "#EAB308", // Yellow-500
  },

  spacing: {
    xs: "4px",
    sm: "8px",
    md: "16px",
    lg: "24px",
    xl: "32px",
    "2xl": "48px",
  },

  radius: {
    sm: "6px",
    md: "8px",
    lg: "12px",
    xl: "16px",
    full: "9999px",
  },

  shadow: {
    sm: "0 1px 2px rgba(0, 0, 0, 0.05)",
    md: "0 4px 6px -1px rgba(0, 0, 0, 0.1)",
    lg: "0 10px 15px -3px rgba(0, 0, 0, 0.1)",
    glow: {
      blue: "0 0 20px rgba(59, 130, 246, 0.3)",
      green: "0 0 20px rgba(34, 197, 94, 0.3)",
    },
  },

  transition: {
    fast: "150ms ease",
    normal: "200ms ease",
    slow: "300ms ease",
  },

  font: {
    family: {
      sans: "Inter, system-ui, sans-serif",
      mono: "JetBrains Mono, monospace",
    },
    size: {
      xs: "12px",
      sm: "14px",
      md: "16px",
      lg: "18px",
      xl: "24px",
    },
  },

  layout: {
    treeNav: {
      width: "240px",
      collapsedWidth: "64px",
    },
    quickActions: {
      width: "280px",
    },
    slidePanel: {
      width: "400px",
    },
    memoryGrid: {
      cellWidth: "64px",
      cellHeight: "48px",
      gap: "4px",
    },
  },
} as const;
