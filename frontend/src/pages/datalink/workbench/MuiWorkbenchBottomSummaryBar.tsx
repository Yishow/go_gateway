import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { useTranslation } from 'react-i18next';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import { useWorkbench } from './WorkbenchProvider';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import { SENTRY_SX } from './sentrySurfaceStyles';
import type { StepReadinessState, WorkbenchReadiness, WorkbenchStep } from './workbenchTypes';

const s = tokens.archetype.sentry;

function readinessColor(status: WorkbenchReadiness): string {
  switch (status) {
    case 'ready':
    case 'applied':
      return s.highlight;
    case 'partial':
      return s.warm;
    case 'blocked':
      return tokens.status.error;
    case 'draft':
      return s.border;
  }
}

/**
 * Sentry System Telemetry Bar — bottom footer expressing command center
 * status metrics in monospace typography with severity chips.
 */
export function MuiWorkbenchBottomSummaryBar() {
  const { t } = useTranslation();
  const { activeOutputTarget, activeStep } = useWorkbench();
  const {
    pointCount,
    linkedTagCount,
    outputCandidateCount,
    deviceReadiness,
    sourceReadiness,
    tagReadiness,
    outputReadiness,
  } = useWorkbenchSummary();

  const readinessItems: ReadonlyArray<{
    step: WorkbenchStep;
    labelKey: string;
    state: StepReadinessState;
  }> = [
    { step: 'device', labelKey: 'workbench.bottomSummary.readiness.device', state: deviceReadiness },
    { step: 'source', labelKey: 'workbench.bottomSummary.readiness.source', state: sourceReadiness },
    { step: 'tag', labelKey: 'workbench.bottomSummary.readiness.tag', state: tagReadiness },
    { step: 'output', labelKey: 'workbench.bottomSummary.readiness.output', state: outputReadiness },
  ];

  return (
    <Box
      component="footer"
      aria-label={t('workbench.bottomSummary.ariaLabel')}
      data-testid="workbench-bottom-summary-bar"
      sx={{
        display: 'flex',
        flexWrap: 'wrap',
        alignItems: 'center',
        gap: 2,
        px: 2,
        py: 0.75,
        bgcolor: s.elevated,
        border: `1px solid ${s.border}`,
        borderRadius: tokens.radius.sm,
        boxShadow: tokens.treatment.insetButton,
        fontFamily: s.uiFont,
      }}
    >
      {/* Metric counters — monospace */}
      <Box component="dl" sx={{ display: 'flex', alignItems: 'center', gap: 2, m: 0 }}>
        <TelemetryMetric label={t('workbench.bottomSummary.points')} value={pointCount} testId="summary-point-count" />
        <TelemetryMetric label={t('workbench.bottomSummary.tags')} value={linkedTagCount} testId="summary-tag-count" />
        <TelemetryMetric label={t('workbench.bottomSummary.outputs')} value={outputCandidateCount} testId="summary-output-count" />
      </Box>

      {/* Separator */}
      <Box sx={{ height: 16, width: '1px', bgcolor: s.border, display: { xs: 'none', sm: 'block' } }} aria-hidden />

      {/* Readiness indicators with severity dots */}
      <Box
        component="ul"
        data-testid="readiness-indicators"
        sx={{ display: 'flex', alignItems: 'center', gap: 1, listStyle: 'none', m: 0, p: 0 }}
      >
        {readinessItems.map((item) => {
          const isActive = item.step === activeStep;
          const color = readinessColor(item.state.status);
          return (
            <Box
              component="li"
              key={item.step}
              data-testid={`readiness-${item.step}`}
              data-readiness={item.state.status}
              data-reason={item.state.reason ?? undefined}
              data-emphasis={isActive ? 'active' : 'compact'}
              sx={{
                display: 'inline-flex',
                alignItems: 'center',
                gap: 0.5,
                px: 0.75,
                py: 0.25,
                borderRadius: tokens.radius.sm,
                border: `1px solid ${isActive ? color + '44' : s.border}`,
                bgcolor: isActive ? `${color}14` : 'transparent',
              }}
            >
              <Box
                sx={{
                  width: 5,
                  height: 5,
                  borderRadius: '50%',
                  bgcolor: color,
                  flexShrink: 0,
                  boxShadow: isActive ? `0 0 4px ${color}66` : 'none',
                }}
              />
              <Typography
                sx={{
                  ...SENTRY_SX.sectionLabel,
                  fontSize: '10px',
                  color: isActive ? tokens.text.primary : tokens.text.muted,
                }}
              >
                {t(item.labelKey)}
              </Typography>
            </Box>
          );
        })}
      </Box>

      {/* Separator */}
      <Box sx={{ height: 16, width: '1px', bgcolor: s.border, display: { xs: 'none', sm: 'block' } }} aria-hidden />

      {/* Active output target */}
      <Box data-testid="active-output-target" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
        <Typography sx={{ ...SENTRY_SX.sectionLabel, fontSize: '10px' }}>
          {t('workbench.bottomSummary.target')}
        </Typography>
        <Typography sx={{ ...SENTRY_SX.monoData, fontSize: '11px', fontWeight: 600, color: s.highlight }}>
          {t(`workbench.bottomSummary.targets.${activeOutputTarget}`)}
        </Typography>
      </Box>
    </Box>
  );
}

function TelemetryMetric({ label, value, testId }: { label: string; value: number; testId: string }) {
  return (
    <Box sx={{ display: 'flex', alignItems: 'baseline', gap: 0.5 }}>
      <Typography component="dt" sx={{ ...SENTRY_SX.sectionLabel, fontSize: '10px' }}>
        {label}
      </Typography>
      <Typography
        component="dd"
        data-testid={testId}
        sx={{
          fontFamily: tokens.typography.family.code,
          fontWeight: 700,
          fontSize: '13px',
          color: s.highlight,
          m: 0,
          letterSpacing: '-0.02em',
        }}
      >
        {value}
      </Typography>
    </Box>
  );
}
