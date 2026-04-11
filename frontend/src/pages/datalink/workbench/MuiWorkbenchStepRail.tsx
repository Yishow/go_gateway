import Box from '@mui/material/Box';
import ButtonBase from '@mui/material/ButtonBase';
import { useTranslation } from 'react-i18next';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import { useWorkbench } from './WorkbenchProvider';
import { useWorkbenchOutputMainline } from './useWorkbenchOutputMainline';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import { getWorkbenchReadinessExplanationKey } from './workbenchShellIncidentModel';
import {
  WORKBENCH_STEPS,
  WORKBENCH_STEP_META,
  type StepReadinessState,
  type WorkbenchReadiness,
  type WorkbenchStep,
} from './workbenchTypes';

const s = tokens.archetype.sentry;

function readinessGlyph(status: WorkbenchReadiness): string {
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
 * Sentry Phase Indicator — inline row of compact phase pills.
 *
 * Replaces generic MUI Tabs. Each phase is a tactile chip with:
 * - Ordinal number in a small box
 * - Uppercase label in Rubik
 * - Readiness dot with severity color
 * - Active phase: purple accent border + lime indicator glow
 */
export function MuiWorkbenchStepRail() {
  const { t } = useTranslation();
  const { activeStep, selectedDeviceId, setActiveStep } = useWorkbench();
  const { deviceReadiness, sourceReady, sourceReadiness, tagReady, tagReadiness } =
    useWorkbenchSummary();
  const { outputReadiness } = useWorkbenchOutputMainline({
    hasSelectedDevice: Boolean(selectedDeviceId),
    sourceReady,
    tagReady,
  });

  const stepReadiness: Record<WorkbenchStep, StepReadinessState> = {
    device: deviceReadiness,
    source: sourceReadiness,
    tag: tagReadiness,
    output: outputReadiness,
  };

  return (
    <Box
      component="nav"
      aria-label={t('workbench.stepRail.ariaLabel')}
      data-testid="workbench-step-rail"
      sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}
    >
      {WORKBENCH_STEPS.map((step, index) => {
        const readiness = stepReadiness[step];
        const isActive = step === activeStep;
        const dotColor = readinessGlyph(readiness.status);
        const explanationKey = getWorkbenchReadinessExplanationKey(readiness);
        const explanationLabel = explanationKey ? t(explanationKey) : null;
        return (
          <ButtonBase
            key={step}
            role="tab"
            onClick={() => setActiveStep(step)}
            data-readiness={readiness.status}
            aria-label={`${t(WORKBENCH_STEP_META[step].labelKey)} - ${t(`workbench.readiness.${readiness.status}`)}${explanationLabel ? ` - ${explanationLabel}` : ''}`}
            aria-current={isActive ? 'step' : undefined}
            aria-selected={isActive}
            title={explanationLabel ?? t(`workbench.readiness.${readiness.status}`)}
            sx={{
              display: 'inline-flex',
              alignItems: 'center',
              gap: 0.75,
              px: 1.25,
              py: 0.5,
              borderRadius: '13px',
              border: `1px solid ${isActive ? s.accent : s.border}`,
              bgcolor: isActive ? `${s.accent}18` : 'transparent',
              boxShadow: isActive ? tokens.treatment.insetButton : 'none',
              transition: 'all 0.15s ease',
              '&:hover': {
                bgcolor: `${s.accent}14`,
                borderColor: s.accentMuted,
              },
            }}
          >
            {/* Ordinal badge */}
            <Box
              component="span"
              sx={{
                display: 'inline-flex',
                alignItems: 'center',
                justifyContent: 'center',
                width: 18,
                height: 18,
                borderRadius: '3px',
                fontSize: '10px',
                fontWeight: 700,
                fontFamily: s.uiFont,
                bgcolor: isActive ? s.accent : s.elevated,
                color: isActive ? '#fff' : tokens.text.muted,
                border: `1px solid ${isActive ? s.accent : s.border}`,
              }}
            >
              {index + 1}
            </Box>
            {/* Phase label */}
            <Box
              component="span"
              sx={{
                fontSize: '11px',
                fontWeight: 700,
                fontFamily: s.uiFont,
                letterSpacing: '0.2px',
                textTransform: 'uppercase',
                color: isActive ? tokens.text.primary : tokens.text.muted,
                display: { xs: 'none', sm: 'inline' },
              }}
            >
              {t(WORKBENCH_STEP_META[step].labelKey)}
            </Box>
            {/* Readiness dot */}
            <Box
              component="span"
              data-testid={`step-readiness-${step}`}
              sx={{
                width: 6,
                height: 6,
                borderRadius: '50%',
                bgcolor: dotColor,
                flexShrink: 0,
                boxShadow: isActive ? `0 0 4px ${dotColor}66` : 'none',
              }}
            />
          </ButtonBase>
        );
      })}
    </Box>
  );
}
