import type { ReactNode } from 'react';
import Box from '@mui/material/Box';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import { MuiWorkbenchContextBar } from './MuiWorkbenchContextBar';
import { MuiWorkbenchBottomSummaryBar } from './MuiWorkbenchBottomSummaryBar';
import { WorkbenchInspectorPanel } from './WorkbenchInspectorPanel';
import { useWorkbench } from './WorkbenchProvider';

type MuiWorkbenchFrameProps = {
  children: ReactNode;
};

const s = tokens.archetype.sentry;

/**
 * Sentry Incident Desk — three-region command center shell.
 *
 * ┌────────────────────────────────────────────┐
 * │  ▌COMMAND STRIP (context + phase chips)    │
 * ├─────────────────────────┬──────────────────┤
 * │  OPERATIONS DESK        │   INTEL PANEL    │
 * │  (step content)         │   (inspector)    │
 * ├─────────────────────────┴──────────────────┤
 * │  ▌SYSTEM TELEMETRY BAR                     │
 * └────────────────────────────────────────────┘
 *
 * All surfaces use sentry archetype palette, ambient purple depth,
 * and inset glass treatment. NOT generic MUI Paper.
 */
export function MuiWorkbenchFrame({ children }: MuiWorkbenchFrameProps) {
  const { sourceStepInspectorBanner } = useWorkbench();

  return (
    <Box
      data-testid="workbench-frame"
      data-variant="v2-mui"
      sx={{
        display: 'grid',
        height: '100dvh',
        maxHeight: '100dvh',
        minHeight: 0,
        overflow: 'hidden',
        gap: '6px',
        p: '6px',
        bgcolor: s.canvas,
        color: tokens.text.primary,
        fontFamily: s.uiFont,
        gridTemplateRows: 'auto 1fr auto',
        gridTemplateColumns: '1fr 272px',
        gridTemplateAreas: `
          "context    context"
          "main       inspector"
          "summary    summary"
        `,
      }}
    >
      {/* COMMAND STRIP */}
      <Box sx={{ gridArea: 'context' }}>
        <MuiWorkbenchContextBar />
      </Box>

      {/* OPERATIONS DESK */}
      <Box
        component="main"
        data-testid="workbench-primary-work-area"
        sx={{
          gridArea: 'main',
          display: 'flex',
          flexDirection: 'column',
          minHeight: 0,
          overflow: 'hidden',
          p: 2,
          bgcolor: s.panel,
          border: `1px solid ${s.border}`,
          borderRadius: tokens.radius.md,
          boxShadow: tokens.treatment.ambientPurple,
          position: 'relative',
          '&::before': {
            content: '""',
            position: 'absolute',
            inset: 0,
            background: `linear-gradient(180deg, ${s.accent}08 0%, transparent 40%)`,
            borderRadius: 'inherit',
            pointerEvents: 'none',
          },
        }}
      >
        {children}
      </Box>

      {/* INTEL PANEL */}
      <Box
        sx={{
          gridArea: 'inspector',
          display: 'flex',
          flexDirection: 'column',
          minHeight: 0,
          minWidth: 0,
          overflow: 'hidden',
          bgcolor: s.elevated,
          border: `1px solid ${s.border}`,
          borderRadius: tokens.radius.md,
          boxShadow: tokens.treatment.ambientPurple,
        }}
      >
        {sourceStepInspectorBanner ? (
          <>
            <Box
              sx={{
                flexShrink: 0,
                borderBottom: `1px solid ${s.border}`,
                px: 1.5,
                py: 1,
              }}
            >
              {sourceStepInspectorBanner}
            </Box>
            <Box sx={{ display: 'flex', flexDirection: 'column', flex: 1, minHeight: 0, overflow: 'hidden' }}>
              <WorkbenchInspectorPanel variant="embedded" />
            </Box>
          </>
        ) : (
          <Box sx={{ display: 'flex', flexDirection: 'column', flex: 1, minHeight: 0, overflow: 'hidden' }}>
            <WorkbenchInspectorPanel variant="standalone" />
          </Box>
        )}
      </Box>

      {/* SYSTEM TELEMETRY */}
      <Box sx={{ gridArea: 'summary' }}>
        <MuiWorkbenchBottomSummaryBar />
      </Box>
    </Box>
  );
}
