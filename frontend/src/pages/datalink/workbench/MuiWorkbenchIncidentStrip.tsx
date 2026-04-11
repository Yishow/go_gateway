import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import { useTranslation } from 'react-i18next';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import { useWorkbench } from './WorkbenchProvider';
import { useWorkbenchOutputMainline } from './useWorkbenchOutputMainline';
import { SENTRY_SX } from './sentrySurfaceStyles';
import { useWorkbenchShellDiagnostics } from './useWorkbenchShellDiagnostics';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import {
  getWorkbenchShellIncidentModel,
  getWorkbenchShellReadinessStats,
} from './workbenchShellIncidentModel';

const s = tokens.archetype.sentry;

function toneColor(tone: 'ok' | 'warning' | 'critical') {
  switch (tone) {
    case 'ok':
      return s.highlight;
    case 'warning':
      return s.warm;
    case 'critical':
      return tokens.status.error;
  }
}

export function MuiWorkbenchIncidentStrip() {
  const { t } = useTranslation();
  const { activeStep, selectedDeviceId, setActiveStep } = useWorkbench();
  const { selectedDevice, sourceReady, sourceReadiness, tagReady, tagReadiness } =
    useWorkbenchSummary();
  const { outputReadiness } = useWorkbenchOutputMainline({
    hasSelectedDevice: Boolean(selectedDeviceId),
    sourceReady,
    tagReady,
  });
  const diagnostics = useWorkbenchShellDiagnostics();

  const readiness = getWorkbenchShellReadinessStats([
    sourceReadiness,
    tagReadiness,
    outputReadiness,
  ]);
  const incident = getWorkbenchShellIncidentModel({
    hasSelectedDevice: Boolean(selectedDevice),
    activeStep,
    sourceReadiness,
    tagReadiness,
    outputReadiness,
  });
  const accent = toneColor(incident.tone);

  const handleReturnAction = () => {
    if (incident.returnStep === activeStep && incident.returnStep === 'output') {
      const anchor = document.querySelector<HTMLElement>('[data-testid="output-primary-anchor"]');
      anchor?.focus();
      anchor?.scrollIntoView?.({ block: 'nearest' });
      return;
    }

    setActiveStep(incident.returnStep);
  };

  return (
    <Box
      data-testid="shell-incident-strip"
      sx={{
        display: 'flex',
        minWidth: 300,
        maxWidth: 440,
        flex: '1 1 360px',
        alignItems: 'stretch',
        gap: 1,
        px: 1.25,
        py: 1,
        borderRadius: tokens.radius.md,
        border: `1px solid ${accent}44`,
        bgcolor: `${s.panel}`,
        boxShadow: tokens.treatment.insetButton,
      }}
    >
      <Box sx={{ minWidth: 0, flex: 1, display: 'flex', flexDirection: 'column', gap: 0.75 }}>
        <Box sx={{ display: 'flex', alignItems: 'baseline', gap: 0.75, minWidth: 0 }}>
          <Typography sx={{ ...SENTRY_SX.sectionLabel, color: accent }}>
            {t('workbench.shell.eyebrow')}
          </Typography>
          <Typography
            data-testid="shell-readiness-summary"
            sx={{ ...SENTRY_SX.monoData, minWidth: 0, fontSize: '11px', color: tokens.text.secondary }}
          >
            {t('workbench.shell.readinessSummary', readiness)}
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.25, minWidth: 0 }}>
          <Typography sx={{ ...SENTRY_SX.sectionLabel, color: tokens.text.muted }}>
            {t('workbench.shell.labels.blocker')}
          </Typography>
          <Typography
            data-testid="shell-incident-blocker"
            sx={{ fontSize: '12px', lineHeight: 1.45, color: tokens.text.primary }}
          >
            {t(incident.blockerKey)}
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.75, minWidth: 0 }}>
          <Typography sx={{ ...SENTRY_SX.sectionLabel, color: tokens.text.muted }}>
            {t('workbench.shell.labels.refresh')}
          </Typography>
          <Typography
            data-testid="shell-refresh-status"
            sx={{
              ...SENTRY_SX.monoData,
              minWidth: 0,
              fontSize: '11px',
              color: diagnostics.isRefreshing || diagnostics.statusKey === 'workbench.shell.refresh.error'
                ? accent
                : tokens.text.muted,
            }}
          >
            {t(diagnostics.statusKey)}
          </Typography>
        </Box>
      </Box>

      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.75, justifyContent: 'space-between' }}>
        <Button
          data-testid="shell-refresh-action"
          size="small"
          disabled={diagnostics.isRefreshing}
          onClick={() => {
            void diagnostics.refreshDiagnostics();
          }}
          sx={{
            ...SENTRY_SX.insetBtn,
            minWidth: 120,
            color: tokens.text.primary,
          }}
        >
          {t(diagnostics.actionKey)}
        </Button>
        <Button
          data-testid="shell-return-action"
          size="small"
          onClick={handleReturnAction}
          sx={{
            ...SENTRY_SX.insetBtn,
            minWidth: 120,
            bgcolor: `${accent}22`,
            borderColor: `${accent}44`,
            color: accent,
            '&:hover': { bgcolor: `${accent}33`, borderColor: accent },
          }}
        >
          {t(incident.returnLabelKey)}
        </Button>
      </Box>
    </Box>
  );
}
