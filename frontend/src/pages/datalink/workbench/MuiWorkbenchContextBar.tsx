import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import { useTranslation } from 'react-i18next';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import { useWorkbench } from './WorkbenchProvider';
import { WorkbenchSourceRuntimeCluster } from './WorkbenchSourceRuntimeCluster';
import { MuiWorkbenchIncidentStrip } from './MuiWorkbenchIncidentStrip';
import { MuiWorkbenchStepRail } from './MuiWorkbenchStepRail';
import {
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
} from './workbenchDeviceFormModel';
import { useWorkbenchOutputMainline } from './useWorkbenchOutputMainline';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import { SENTRY_SX, severityColor } from './sentrySurfaceStyles';
import type { StepReadinessState, WorkbenchStep } from './workbenchTypes';

const s = tokens.archetype.sentry;

/**
 * Sentry Command Strip — top-level navigation and context bar.
 *
 * NOT a standard MUI AppBar. This is a glass command strip with:
 * - Phase indicator chips (left)
 * - Active device identity with severity signal (center)
 * - Primary command button (right)
 * All expressed in Rubik uppercase with purple depth.
 */
export function MuiWorkbenchContextBar() {
  const { t } = useTranslation();
  const { activeStep, selectedDeviceId, setActiveStep, sourceStepNotice } = useWorkbench();
  const { selectedDevice, sourceReady, tagReady } = useWorkbenchSummary();
  const { outputReadiness } = useWorkbenchOutputMainline({
    hasSelectedDevice: Boolean(selectedDeviceId),
    sourceReady,
    tagReady,
  });
  const deviceLabel = selectedDevice?.name ?? t('workbench.contextBar.noDevice');
  const primaryAction = getPrimaryAction({
    activeStep,
    hasSelectedDevice: Boolean(selectedDeviceId),
    sourceReady,
    tagReady,
    outputReadiness,
    t,
    setActiveStep,
  });

  const severity = selectedDevice ? severityColor(selectedDevice.last_test_success) : s.border;

  return (
    <Box
      component="header"
      aria-label={t('workbench.contextBar.ariaLabel')}
      data-testid="workbench-context-bar"
      sx={{
        display: 'flex',
        flexWrap: 'wrap',
        alignItems: 'center',
        gap: 1.5,
        px: 2,
        py: 1,
        bgcolor: 'rgba(36, 26, 59, 0.8)',
        border: `1px solid ${s.border}`,
        borderRadius: tokens.radius.md,
        boxShadow: tokens.treatment.ambientPurple,
        backdropFilter: tokens.treatment.glassPanel.backdropFilter,
        position: 'relative',
        '&::before': {
          content: '""',
          position: 'absolute',
          inset: 0,
          background: `linear-gradient(90deg, ${s.accent}12 0%, transparent 50%, ${s.highlight}0a 100%)`,
          borderRadius: 'inherit',
          pointerEvents: 'none',
        },
      }}
    >
      {/* Phase chips */}
      <MuiWorkbenchStepRail />

      {/* Separator */}
      <Box
        sx={{
          display: { xs: 'none', lg: 'block' },
          height: 28,
          width: '1px',
          flexShrink: 0,
          bgcolor: s.border,
        }}
        aria-hidden
      />

      {/* Device identity + severity signal */}
      <Box
        sx={{ flex: 1, minWidth: 0, display: 'flex', flexWrap: 'wrap', alignItems: 'center', gap: 1 }}
        data-testid="context-bar-step-summary"
        data-active-step={activeStep}
      >
        {/* Severity dot */}
        <Box
          sx={{
            width: 8,
            height: 8,
            borderRadius: '50%',
            bgcolor: severity,
            flexShrink: 0,
            boxShadow: `0 0 6px ${severity}66`,
          }}
        />
        <Typography
          noWrap
          data-testid="context-bar-device-name"
          sx={{
            ...SENTRY_SX.sectionLabel,
            fontSize: '13px',
            color: tokens.text.primary,
          }}
        >
          {deviceLabel}
        </Typography>
        {selectedDevice ? (
          <>
            <Box
              sx={{
                px: 0.75,
                py: 0.25,
                borderRadius: tokens.radius.sm,
                border: `1px solid ${s.border}`,
                bgcolor: `${s.accent}14`,
              }}
            >
              <Typography sx={{ fontSize: '10px', fontWeight: 700, letterSpacing: '0.2px', textTransform: 'uppercase', color: s.accent }}>
                {t(getWorkbenchProtocolLabelKey(selectedDevice.protocol))}
              </Typography>
            </Box>
            {activeStep !== 'source' ? (
              <Box
                data-testid="context-bar-device-status-chip"
                sx={{
                  px: 0.75,
                  py: 0.25,
                  borderRadius: tokens.radius.sm,
                  border: `1px solid ${severity}44`,
                  bgcolor: `${severity}18`,
                }}
              >
                <Typography sx={{ fontSize: '10px', fontWeight: 700, letterSpacing: '0.2px', textTransform: 'uppercase', color: severity }}>
                  {t(getWorkbenchDeviceStatusLabelKey(selectedDevice.status))}
                </Typography>
              </Box>
            ) : null}
          </>
        ) : null}
        <WorkbenchSourceRuntimeCluster />
      </Box>

      <MuiWorkbenchIncidentStrip />

      {/* Source step notice */}
      {activeStep === 'source' && sourceStepNotice ? (
        <Box
          data-testid="source-step-notice"
          sx={{
            width: '100%',
            mt: 0.5,
            px: 1.5,
            py: 0.75,
            borderRadius: tokens.radius.sm,
            border: `1px solid ${s.warm}40`,
            bgcolor: `${s.warm}14`,
            fontSize: '11px',
            fontWeight: 600,
            letterSpacing: '0.2px',
            textTransform: 'uppercase',
            color: s.warm,
          }}
        >
          {sourceStepNotice}
        </Box>
      ) : null}

      {/* Primary CTA */}
      <Button
        data-testid="context-bar-primary-action"
        size="small"
        disabled={primaryAction.disabled}
        onClick={primaryAction.onClick}
        sx={{
          ...SENTRY_SX.insetBtn,
          flexShrink: 0,
          minWidth: 100,
          bgcolor: `${s.highlight}22`,
          borderColor: `${s.highlight}44`,
          color: s.highlight,
          '&:hover': { bgcolor: `${s.highlight}33`, borderColor: s.highlight },
          '&.Mui-disabled': { opacity: 0.4 },
        }}
      >
        {primaryAction.label}
      </Button>
    </Box>
  );
}

function getPrimaryAction(input: {
  activeStep: WorkbenchStep;
  hasSelectedDevice: boolean;
  sourceReady: boolean;
  tagReady: boolean;
  outputReadiness: StepReadinessState;
  t: (key: string) => string;
  setActiveStep: (step: WorkbenchStep) => void;
}) {
  if (!input.hasSelectedDevice) {
    return {
      label: input.t('workbench.contextBar.actions.selectDevice'),
      onClick: () => input.setActiveStep('device'),
      disabled: false,
    };
  }
  if (input.activeStep === 'device') {
    return {
      label: input.t('workbench.contextBar.actions.gotoSource'),
      onClick: () => input.setActiveStep('source'),
      disabled: false,
    };
  }
  if (input.activeStep === 'source') {
    return {
      label: input.t('workbench.contextBar.actions.gotoTag'),
      onClick: () => input.setActiveStep('tag'),
      disabled: !input.sourceReady,
    };
  }
  if (input.activeStep === 'tag') {
    return {
      label: input.t('workbench.contextBar.actions.gotoOutput'),
      onClick: () => input.setActiveStep('output'),
      disabled: !input.tagReady,
    };
  }

  if (input.outputReadiness.reason === 'no-source-points') {
    return {
      label: input.t('workbench.shell.actions.returnSource'),
      onClick: () => input.setActiveStep('source'),
      disabled: false,
    };
  }

  if (
    input.outputReadiness.reason === 'no-tags-linked' ||
    input.outputReadiness.reason === 'output-rule-required'
  ) {
    return {
      label: input.t('workbench.shell.actions.returnTag'),
      onClick: () => input.setActiveStep('tag'),
      disabled: false,
    };
  }

  return {
    label: input.t('workbench.actionDock.nextAction.configureOutput'),
    onClick: () => {
      const anchor = document.querySelector<HTMLElement>('[data-testid="output-primary-anchor"]');
      anchor?.focus();
      anchor?.scrollIntoView?.({ block: 'nearest' });
    },
    disabled: false,
  };
}
