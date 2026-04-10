import Alert from '@mui/material/Alert';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Chip from '@mui/material/Chip';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import { alpha } from '@mui/material/styles';
import { useTranslation } from 'react-i18next';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import { useSourceRulesQuery } from '../../../hooks/datalink/useSourceRules';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import { LocalModbusBoard } from './LocalModbusBoard';
import { SENTRY_SX } from './sentrySurfaceStyles';
import { useWorkbench } from './WorkbenchProvider';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import type { OutputTarget } from './workbenchTypes';

const sentry = tokens.archetype.sentry;

type DeskCommand = {
  label: string;
  targetId: string;
};

const TARGETS: readonly { key: OutputTarget; labelKey: string }[] = [
  { key: 'modbus', labelKey: 'workbench.output.targetSwitcher.modbus' },
  { key: 'database', labelKey: 'workbench.output.targetSwitcher.database' },
] as const;

function focusSurface(testId: string) {
  const element = document.querySelector<HTMLElement>(`[data-testid="${testId}"]`);
  if (!element) {
    return;
  }
  element.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
  element.focus?.();
}

function triggerWorkbenchAction(testId: string) {
  const element = document.querySelector<HTMLElement>(`[data-testid="${testId}"]`);
  if (!element) {
    return;
  }
  if (element instanceof HTMLButtonElement && element.disabled) {
    return;
  }
  element.focus?.();
  element.click();
}

export function MuiOutputIncidentDesk() {
  const { t } = useTranslation();
  const {
    activeOutputTarget,
    crossStepContext,
    selectedDeviceId,
    setActiveOutputTarget,
    setActiveStep,
  } = useWorkbench();
  const { linkedTagCount } = useWorkbenchSummary();
  const { data: rules = [] } = useSourceRulesQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const ruleId = crossStepContext.focusedRuleId ?? rules[0]?.id ?? null;
  const activeRule = rules.find((rule) => rule.id === ruleId) ?? null;
  const candidatesQuery = useSourceRuleCandidatesQuery(ruleId);
  const reviewSet = activeOutputTarget === 'modbus'
    ? candidatesQuery.data?.local_modbus_outputs ?? null
    : candidatesQuery.data?.database_outputs ?? null;
  const targetCandidates = reviewSet?.candidates ?? [];
  const attentionCount = targetCandidates.filter((candidate) => candidate.status !== 'ready').length;
  const noRuleSelected = !selectedDeviceId || !ruleId;
  const staleReview = Boolean(
    activeRule
      && candidatesQuery.data
      && activeRule.revision_id !== candidatesQuery.data.revision_id,
  );
  const targetLabel = t(
    TARGETS.find((target) => target.key === activeOutputTarget)?.labelKey
      ?? 'workbench.output.targetSwitcher.modbus',
  );

  const priorityCopy = noRuleSelected
    ? t('workbench.output.incident.priority.noRule')
    : candidatesQuery.isError
      ? candidatesQuery.error instanceof Error
        ? candidatesQuery.error.message
        : t('workbench.output.incident.loadFailed')
      : reviewSet?.status === 'blocked'
        ? reviewSet.reason
          ?? t('workbench.output.incident.priority.blockedFallback', {
            target: targetLabel,
          })
        : reviewSet?.status === 'deferred'
          ? reviewSet.reason
            ?? t('workbench.output.incident.priority.deferredFallback', {
              target: targetLabel,
            })
          : attentionCount > 0
            ? t('workbench.output.incident.priority.blockedFallback', {
              target: targetLabel,
            })
            : t('workbench.output.incident.priority.readyCopy', {
              target: targetLabel,
            });
  const needsAttention = noRuleSelected
    || candidatesQuery.isError
    || reviewSet?.status === 'blocked'
    || reviewSet?.status === 'deferred'
    || attentionCount > 0;

  const commands: { primary: DeskCommand; secondary: DeskCommand; focus: DeskCommand } = activeOutputTarget === 'modbus'
    ? {
      primary: {
        label: t('workbench.output.incident.commands.runDryRun'),
        targetId: 'output-modbus-dry-run',
      },
      secondary: {
        label: t('workbench.output.incident.commands.syncMappings'),
        targetId: 'output-modbus-sync',
      },
      focus: {
        label: t('workbench.output.incident.commands.focusRegisterMap'),
        targetId: 'register-map-canvas',
      },
    }
    : {
      primary: {
        label: t('workbench.output.incident.commands.configureConnector'),
        targetId: 'output-database-toggle-connector',
      },
      secondary: {
        label: t('workbench.output.incident.commands.refreshValidation'),
        targetId: 'output-database-refresh-validation',
      },
      focus: {
        label: t('workbench.output.incident.commands.focusSchema'),
        targetId: 'schema-snapshot',
      },
    };

  const handoffCopy = noRuleSelected
    ? t('workbench.output.incident.handoff.noRule')
    : targetCandidates.length === 0
      ? t('workbench.output.incident.handoff.noCandidates', { target: targetLabel })
      : t('workbench.output.incident.handoff.ready');

  return (
    <Box
      sx={{
        display: 'grid',
        minHeight: 0,
        flex: 1,
        gap: 1.5,
        gridTemplateColumns: {
          xs: 'minmax(0, 1fr)',
          lg: 'minmax(18rem, 20rem) minmax(0, 1fr)',
        },
      }}
    >
      <Box
        data-testid="output-incident-command-panel"
        sx={{
          ...SENTRY_SX.commandStrip,
          display: 'grid',
          gap: 1.25,
          alignContent: 'start',
          minHeight: 0,
        }}
      >
        <Box
          data-testid="output-incident-priority-card"
          sx={{
            borderRadius: '12px',
            border: `1px solid ${needsAttention ? `${sentry.warm}55` : sentry.border}`,
            bgcolor: needsAttention ? `${sentry.warm}14` : `${sentry.highlight}12`,
            px: 1.5,
            py: 1.25,
          }}
        >
          <Typography
            sx={{
              ...SENTRY_SX.sectionLabel,
              color: needsAttention ? sentry.warm : sentry.highlight,
            }}
          >
            {t(
              needsAttention
                ? 'workbench.output.incident.priority.attention'
                : 'workbench.output.incident.priority.ready',
            )}
          </Typography>
          <Typography
            sx={{
              mt: 0.5,
              color: tokens.text.secondary,
              fontSize: '13px',
              lineHeight: 1.5,
            }}
          >
            {priorityCopy}
          </Typography>
        </Box>

        {candidatesQuery.isError ? (
          <Alert
            data-testid="output-incident-retry-banner"
            severity="error"
            sx={{ ...SENTRY_SX.frostedPanel, borderColor: `${tokens.status.error}55` }}
            action={
              <Button
                data-testid="output-incident-retry-btn"
                size="small"
                onClick={() => void candidatesQuery.refetch()}
                sx={SENTRY_SX.insetBtn}
              >
                {t('workbench.tag.reviewSurface.retry')}
              </Button>
            }
          >
            {candidatesQuery.error instanceof Error
              ? candidatesQuery.error.message
              : t('workbench.output.incident.loadFailed')}
          </Alert>
        ) : null}

        {staleReview ? (
          <Alert
            data-testid="output-incident-stale-banner"
            severity="warning"
            sx={{ ...SENTRY_SX.frostedPanel, borderColor: `${sentry.warm}55` }}
            action={
              <Button
                data-testid="output-incident-stale-refresh"
                size="small"
                onClick={() => void candidatesQuery.refetch()}
                sx={SENTRY_SX.insetBtn}
              >
                {t('workbench.source.templates.recoverAction', { defaultValue: 'Refresh' })}
              </Button>
            }
          >
            {t('workbench.output.incident.stale')}
          </Alert>
        ) : null}

        <Box sx={{ minHeight: 0 }}>
          <Typography
            sx={{
              ...SENTRY_SX.sectionLabel,
              mb: 0.75,
              color: tokens.text.muted,
            }}
          >
            {t('workbench.output.incident.tabListAria')}
          </Typography>
          <Stack
            role="tablist"
            aria-label={t('workbench.output.incident.tabListAria')}
            spacing={0.75}
          >
            {TARGETS.map(({ key, labelKey }) => (
              <Button
                key={key}
                role="tab"
                fullWidth
                aria-selected={activeOutputTarget === key}
                data-testid={`output-incident-target-${key}`}
                onClick={() => setActiveOutputTarget(key)}
                sx={{
                  ...SENTRY_SX.modeChip,
                  justifyContent: 'flex-start',
                  px: 1.5,
                  py: 1,
                  minHeight: 40,
                  borderRadius: '10px',
                  border: `1px solid ${activeOutputTarget === key ? alpha(sentry.highlight, 0.5) : sentry.border}`,
                  bgcolor:
                    activeOutputTarget === key
                      ? alpha(sentry.highlight, 0.14)
                      : alpha(sentry.panel, 0.72),
                  color:
                    activeOutputTarget === key ? sentry.highlight : tokens.text.secondary,
                  '&:hover': {
                    bgcolor:
                      activeOutputTarget === key
                        ? alpha(sentry.highlight, 0.18)
                        : alpha(sentry.panel, 0.92),
                    color:
                      activeOutputTarget === key ? sentry.highlight : tokens.text.primary,
                  },
                }}
              >
                {t(labelKey)}
              </Button>
            ))}
          </Stack>
        </Box>

        <Box
          data-testid="output-incident-summary-strip"
          sx={{ ...SENTRY_SX.frostedPanel, display: 'grid', gap: 0.75, px: 1.5, py: 1.25 }}
        >
          <Typography sx={{ ...SENTRY_SX.sectionLabel, color: sentry.highlight }}>
            {t('workbench.output.incident.summaryEyebrow')}
          </Typography>
          <Stack data-testid="output-incident-summary-primary" direction="row" spacing={0.75} flexWrap="wrap" useFlexGap>
            <Chip
              size="small"
              label={`${linkedTagCount} ${t('workbench.output.incident.metrics.linked')}`}
              sx={SENTRY_SX.modeChip}
            />
            <Chip
              size="small"
              label={`${targetCandidates.length} ${t('workbench.output.incident.metrics.scoped')}`}
              sx={SENTRY_SX.modeChip}
            />
          </Stack>
          <Stack data-testid="output-incident-summary-secondary" direction="row" spacing={0.75} flexWrap="wrap" useFlexGap>
            <Chip
              size="small"
              label={`${attentionCount} ${t('workbench.output.incident.metrics.attention')}`}
              sx={SENTRY_SX.modeChip}
            />
            {candidatesQuery.data?.revision_id ? (
              <Chip
                size="small"
                label={`rev ${candidatesQuery.data.revision_id}`}
                sx={{
                  ...SENTRY_SX.modeChip,
                  fontFamily: 'Monaco, Menlo, Ubuntu Mono, monospace',
                }}
              />
            ) : null}
          </Stack>
        </Box>

        <Box
          data-testid="output-incident-command-dock"
          sx={{ ...SENTRY_SX.frostedPanel, display: 'grid', gap: 0.75, px: 1.5, py: 1.25 }}
        >
          <Typography sx={{ ...SENTRY_SX.sectionLabel, color: sentry.highlight }}>
            {t('workbench.output.incident.commands.label')}
          </Typography>
          <Button
            data-testid="output-incident-primary-action"
            size="small"
            onClick={() => triggerWorkbenchAction(commands.primary.targetId)}
            sx={SENTRY_SX.insetBtn}
          >
            {commands.primary.label}
          </Button>
          <Button
            data-testid="output-incident-secondary-action"
            size="small"
            onClick={() => triggerWorkbenchAction(commands.secondary.targetId)}
            sx={SENTRY_SX.insetBtn}
          >
            {commands.secondary.label}
          </Button>
          <Button
            data-testid="output-incident-focus-workboard"
            size="small"
            onClick={() => focusSurface(commands.focus.targetId)}
            sx={SENTRY_SX.insetBtn}
          >
            {commands.focus.label}
          </Button>
        </Box>

        <Box
          data-testid="output-incident-handoff-panel"
          sx={{ ...SENTRY_SX.frostedPanel, display: 'grid', gap: 0.75, px: 1.5, py: 1.25 }}
        >
          <Typography
            sx={{
              ...SENTRY_SX.sectionLabel,
              color: noRuleSelected || targetCandidates.length === 0 ? sentry.warm : sentry.highlight,
            }}
          >
            {t('workbench.output.incident.handoff.label')}
          </Typography>
          <Typography sx={{ color: tokens.text.secondary, fontSize: '13px', lineHeight: 1.5 }}>
            {handoffCopy}
          </Typography>
          <Button
            data-testid="output-incident-goto-tag"
            size="small"
            onClick={() => setActiveStep('tag')}
            sx={SENTRY_SX.insetBtn}
          >
            {t('workbench.output.incident.handoff.gotoTag')}
          </Button>
        </Box>
      </Box>

      <Box
        data-testid="output-incident-workboard"
        sx={{ display: 'flex', minHeight: 0, minWidth: 0, overflow: 'hidden' }}
      >
        <LocalModbusBoard />
      </Box>
    </Box>
  );
}
