import Alert from '@mui/material/Alert';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Chip from '@mui/material/Chip';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import { useTranslation } from 'react-i18next';
import { useSourceRulesQuery } from '../../../hooks/datalink/useSourceRules';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import { useSourceRuleTagReviewDecisionsQuery } from '../../../hooks/datalink/useSourceRuleTagReviewDecisions';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import { TagBindingStudio } from './TagBindingStudio';
import { SENTRY_SX } from './sentrySurfaceStyles';
import { useWorkbench } from './WorkbenchProvider';
import { useWorkbenchSummary } from './useWorkbenchSummary';

const sentry = tokens.archetype.sentry;

/**
 * Incident-desk shell for Tag while preserving the shared TagBindingStudio
 * workboard and apply/recovery contract.
 */
export function MuiTagIncidentDesk() {
  const { t } = useTranslation();
  const { crossStepContext, selectedDeviceId, setActiveStep } = useWorkbench();
  const { tagReady } = useWorkbenchSummary();
  const { data: rules = [] } = useSourceRulesQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const ruleId = crossStepContext.focusedRuleId ?? rules[0]?.id ?? null;
  const activeRule = rules.find((rule) => rule.id === ruleId) ?? null;
  const candidatesQuery = useSourceRuleCandidatesQuery(ruleId);
  const { data: decisions = [] } = useSourceRuleTagReviewDecisionsQuery(ruleId);
  const tagCandidates = candidatesQuery.data?.tags.candidates ?? [];
  const decidedIds = new Set(decisions.map((decision) => decision.candidate_id));
  const pendingCount = tagCandidates.filter((candidate) => !decidedIds.has(candidate.id)).length;
  const readyCount = tagCandidates.length - pendingCount;
  const blockerReason = candidatesQuery.data?.tags.status === 'blocked'
    ? candidatesQuery.data.tags.reason ?? t('workbench.tag.empty.description')
    : null;
  const staleReview = Boolean(
    activeRule
      && candidatesQuery.data
      && activeRule.revision_id !== candidatesQuery.data.revision_id,
  );
  const handoffBlocked = !tagReady || tagCandidates.length === 0;

  const noRuleSelected = !selectedDeviceId || !ruleId;
  const priorityTitle = noRuleSelected
    ? t('workbench.tag.reviewSurface.setStatus.blocked')
    : blockerReason
      ? t('workbench.tag.reviewSurface.setStatus.blocked')
      : t('workbench.tag.reviewSurface.setStatus.ready');
  const priorityCopy = noRuleSelected
    ? t('workbench.tag.empty.description')
    : blockerReason ?? t('workbench.tag.reviewSurface.description');

  return (
    <Box sx={{ display: 'grid', minHeight: 0, flex: 1, gap: 1.5, gridTemplateColumns: { xs: 'minmax(0, 1fr)', lg: 'minmax(18rem, 20rem) minmax(0, 1fr)' } }}>
      <Box
        data-testid="tag-incident-command-panel"
        sx={{ ...SENTRY_SX.commandStrip, display: 'grid', gap: 1.25, alignContent: 'start', minHeight: 0 }}
      >
        <Box
          data-testid="tag-incident-priority-card"
          sx={{
            borderRadius: '12px',
            border: `1px solid ${blockerReason ? `${sentry.warm}55` : sentry.border}`,
            bgcolor: blockerReason ? `${sentry.warm}14` : `${sentry.highlight}12`,
            px: 1.5,
            py: 1.25,
          }}
        >
          <Typography sx={{ ...SENTRY_SX.sectionLabel, color: blockerReason ? sentry.warm : sentry.highlight }}>
            {priorityTitle}
          </Typography>
          <Typography sx={{ mt: 0.5, color: tokens.text.secondary, fontSize: '13px', lineHeight: 1.5 }}>
            {priorityCopy}
          </Typography>
        </Box>

        {candidatesQuery.isError ? (
          <Alert
            data-testid="tag-incident-retry-banner"
            severity="error"
            sx={{ ...SENTRY_SX.frostedPanel, borderColor: `${tokens.status.error}55` }}
            action={
              <Button
                data-testid="tag-incident-retry-btn"
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
              : t('workbench.tag.reviewSurface.loadFailed')}
          </Alert>
        ) : null}

        {staleReview ? (
          <Alert
            data-testid="tag-incident-stale-banner"
            severity="warning"
            sx={{ ...SENTRY_SX.frostedPanel, borderColor: `${sentry.warm}55` }}
            action={
              <Button
                data-testid="tag-incident-stale-refresh"
                size="small"
                onClick={() => void candidatesQuery.refetch()}
                sx={SENTRY_SX.insetBtn}
              >
                {t('workbench.source.templates.recoverAction', { defaultValue: 'Refresh' })}
              </Button>
            }
          >
            {t('workbench.tag.reviewSurface.stale', { defaultValue: 'Review data is stale. Refresh before applying.' })}
          </Alert>
        ) : null}

        <Box
          data-testid="tag-incident-summary-strip"
          sx={{ ...SENTRY_SX.frostedPanel, display: 'grid', gap: 0.75, px: 1.5, py: 1.25 }}
        >
          <Typography sx={{ ...SENTRY_SX.sectionLabel, color: sentry.highlight }}>
            {t('workbench.tag.reviewSurface.eyebrow')}
          </Typography>
          <Stack direction="row" spacing={0.75} flexWrap="wrap" useFlexGap>
            <Chip
              size="small"
              label={`${tagCandidates.length} ${t('workbench.tag.reviewSurface.metrics.generated')}`}
              sx={SENTRY_SX.modeChip}
            />
            <Chip
              size="small"
              label={`${pendingCount} ${t('workbench.tag.reviewSurface.metrics.pendingMappings')}`}
              sx={SENTRY_SX.modeChip}
            />
            <Chip
              size="small"
              label={`${readyCount} ${t('workbench.tag.reviewSurface.setStatus.ready')}`}
              sx={SENTRY_SX.modeChip}
            />
          </Stack>
        </Box>

        <Box
          data-testid="tag-incident-handoff-panel"
          sx={{ ...SENTRY_SX.frostedPanel, display: 'grid', gap: 0.75, px: 1.5, py: 1.25 }}
        >
          <Typography sx={{ ...SENTRY_SX.sectionLabel, color: handoffBlocked ? sentry.warm : sentry.highlight }}>
            {t('workbench.contextBar.actions.gotoOutput')}
          </Typography>
          <Typography data-testid="tag-incident-handoff-status" sx={{ color: tokens.text.secondary, fontSize: '13px', lineHeight: 1.5 }}>
            {handoffBlocked
              ? t('workbench.tag.results.outputBlocked', {
                  defaultValue: 'Resolve the review set before handing off to Output.',
                })
              : t('workbench.tag.reviewSurface.setStatus.ready')}
          </Typography>
          <Button
            data-testid="tag-incident-goto-output"
            size="small"
            disabled={handoffBlocked}
            onClick={() => setActiveStep('output')}
            sx={SENTRY_SX.insetBtn}
          >
            {t('workbench.contextBar.actions.gotoOutput')}
          </Button>
        </Box>
      </Box>

      <Box
        data-testid="tag-incident-workboard"
        sx={{ display: 'flex', minHeight: 0, minWidth: 0, overflow: 'hidden' }}
      >
        <TagBindingStudio />
      </Box>
    </Box>
  );
}
