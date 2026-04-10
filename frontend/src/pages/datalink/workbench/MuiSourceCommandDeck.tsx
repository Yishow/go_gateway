import { useState } from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Chip from '@mui/material/Chip';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import { alpha } from '@mui/material/styles';
import { useTranslation } from 'react-i18next';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import { SourceCanvasSection } from './SourceCanvasSection';
import { MuiWorkbenchSourceStyles } from './MuiWorkbenchSourceStyles';
import { useWorkbench } from './WorkbenchProvider';
import { getWorkbenchProtocolLabelKey } from './workbenchDeviceFormModel';
import { SENTRY_SX } from './sentrySurfaceStyles';
import { useWorkbenchSummary } from './useWorkbenchSummary';
const sentry = tokens.archetype.sentry;
/** Desk modes control which SourceCanvasSection panel gets focus. */
export type SentryDeskMode = 'inspect' | 'build' | 'triage';
const DESK_MODES: readonly { key: SentryDeskMode; labelKey: string }[] = [
  { key: 'inspect', labelKey: 'workbench.source.deskMode.inspect' },
  { key: 'build', labelKey: 'workbench.source.deskMode.build' },
  { key: 'triage', labelKey: 'workbench.source.deskMode.triage' },
] as const;

export function MuiSourceCommandDeck() {
  const { t } = useTranslation();
  const { selectedDeviceId, setActiveStep, sourceStepNotice } = useWorkbench();
  const { selectedDevice, sourceReady, pointCount } = useWorkbenchSummary();
  const [deskMode, setDeskMode] = useState<SentryDeskMode>('inspect');
  const hasDevice = Boolean(selectedDeviceId);
  const isBlocked = !hasDevice || !sourceReady;
  const incidentPriorityCopy =
    sourceStepNotice ??
    (hasDevice
      ? t('workbench.source.handoff.pointSummary', { count: pointCount })
      : t('workbench.source.handoff.noDevice'));

  return (
    <Box
      data-testid="sentry-source-deck"
      sx={{ display: 'flex', minHeight: 0, flex: 1, flexDirection: 'column', gap: 0 }}
    >
      <MuiWorkbenchSourceStyles deskMode={deskMode} />
      <Box
        data-testid="source-command-banner"
        sx={{
          ...SENTRY_SX.frostedPanel,
          px: 2,
          py: 1.5,
          borderRadius: '14px',
          background: `linear-gradient(135deg, ${alpha(sentry.accent, 0.22)} 0%, ${alpha(sentry.panel, 0.98)} 65%)`,
        }}
      >
        <Stack
          direction="row"
          spacing={1.5}
          alignItems="center"
          justifyContent="space-between"
          flexWrap="wrap"
        >
          <Box sx={{ minWidth: 0, flex: 1 }}>
            <Typography
              sx={{
                ...SENTRY_SX.sectionLabel,
                color: sentry.highlight,
                fontSize: '10px',
              }}
            >
              {t('workbench.source.sentryBanner.eyebrow')}
            </Typography>
            <Typography
              variant="subtitle1"
              sx={{
                mt: 0.25,
                fontWeight: 700,
                fontFamily: tokens.typography.family.sentryUi,
              }}
            >
              {t('workbench.source.sentryBanner.title')}
            </Typography>
          </Box>
          <Stack direction="row" spacing={0.75} alignItems="center" flexWrap="wrap">
            <Chip
              size="small"
              label={
                selectedDevice
                  ? t(getWorkbenchProtocolLabelKey(selectedDevice.protocol))
                  : t('workbench.contextBar.noDevice')
              }
              sx={{
                ...SENTRY_SX.modeChip,
                height: 24,
                color: tokens.text.primary,
                bgcolor: `${sentry.accent}20`,
                border: `1px solid ${sentry.border}`,
              }}
            />
            <Chip
              size="small"
              label={t(sourceReady ? 'workbench.source.sentryBanner.ready' : 'workbench.source.sentryBanner.blocked')}
              sx={{
                ...SENTRY_SX.modeChip,
                height: 24,
                color: sourceReady ? sentry.highlight : sentry.warm,
                bgcolor: sourceReady ? `${sentry.highlight}1f` : `${sentry.warm}1f`,
                border: `1px solid ${sourceReady ? `${sentry.highlight}44` : `${sentry.warm}44`}`,
              }}
            />
            {hasDevice ? (
              <Chip
                size="small"
                label={`${pointCount} pts`}
                sx={{
                  ...SENTRY_SX.modeChip,
                  height: 24,
                  color: tokens.text.secondary,
                  bgcolor: sentry.panel,
                  border: `1px solid ${sentry.border}`,
                  fontFamily: 'Monaco, Menlo, monospace',
                  fontSize: '10px',
                }}
              />
            ) : null}
          </Stack>
        </Stack>
      </Box>
      <Box
        data-testid="source-incident-desk-shell"
        sx={{
          display: 'grid',
          minHeight: 0,
          flex: 1,
          gap: 1.5,
          mt: 1.5,
          gridTemplateColumns: {
            xs: 'minmax(0, 1fr)',
            lg: 'minmax(16rem, 19rem) minmax(0, 1fr)',
          },
        }}
      >
        <Box
          data-testid="source-incident-command-panel"
          sx={{
            ...SENTRY_SX.commandStrip,
            display: 'flex',
            minHeight: 0,
            flexDirection: 'column',
            gap: 1.25,
            px: 1.5,
            py: 1.5,
            border: `1px solid ${sentry.border}`,
          }}
        >
          <Box
            data-testid="source-incident-priority-card"
            sx={{
              borderRadius: '12px',
              border: `1px solid ${sourceReady ? `${sentry.highlight}33` : `${sentry.warm}44`}`,
              bgcolor: sourceReady ? alpha(sentry.highlight, 0.08) : `${sentry.warm}14`,
              px: 1.25,
              py: 1,
            }}
          >
            <Typography
              sx={{
                ...SENTRY_SX.sectionLabel,
                color: sourceReady ? sentry.highlight : sentry.warm,
              }}
            >
              {t(sourceReady ? 'workbench.source.sentryBanner.ready' : 'workbench.source.sentryBanner.blocked')}
            </Typography>
            <Typography
              sx={{
                mt: 0.5,
                color: tokens.text.secondary,
                fontSize: '13px',
                lineHeight: 1.5,
              }}
            >
              {incidentPriorityCopy}
            </Typography>
          </Box>
          <Box sx={{ minHeight: 0 }}>
            <Typography
              sx={{
                ...SENTRY_SX.sectionLabel,
                mb: 0.75,
                color: tokens.text.muted,
              }}
            >
              {t('workbench.source.deskMode.tabListAria')}
            </Typography>
            <Stack
              role="tablist"
              aria-label={t('workbench.source.deskMode.tabListAria')}
              spacing={0.75}
            >
              {DESK_MODES.map(({ key, labelKey }) => (
                <Button
                  key={key}
                  role="tab"
                  fullWidth
                  aria-selected={deskMode === key}
                  data-testid={`source-desk-tab-${key}`}
                  onClick={() => setDeskMode(key)}
                  sx={{
                    ...SENTRY_SX.modeChip,
                    justifyContent: 'flex-start',
                    px: 1.5,
                    py: 1,
                    minHeight: 40,
                    borderRadius: '10px',
                    border: `1px solid ${deskMode === key ? alpha(sentry.highlight, 0.5) : sentry.border}`,
                    bgcolor:
                      deskMode === key
                        ? alpha(sentry.highlight, 0.14)
                        : alpha(sentry.panel, 0.72),
                    color:
                      deskMode === key ? sentry.highlight : tokens.text.secondary,
                    '&:hover': {
                      bgcolor:
                        deskMode === key
                          ? alpha(sentry.highlight, 0.18)
                          : alpha(sentry.panel, 0.92),
                      color:
                        deskMode === key
                          ? sentry.highlight
                          : tokens.text.primary,
                    },
                  }}
                >
                  {t(labelKey)}
                </Button>
              ))}
            </Stack>
          </Box>
          <Box
            data-testid="source-incident-handoff-panel"
            sx={{
              mt: 'auto',
              borderRadius: '12px',
              border: `1px solid ${sentry.border}`,
              bgcolor: alpha(sentry.elevated, 0.92),
              px: 1.25,
              py: 1.25,
            }}
          >
            <Typography sx={{ ...SENTRY_SX.sectionLabel, color: tokens.text.muted }}>
              {t('workbench.source.handoff.dockLabel')}
            </Typography>
            <Typography
              sx={{ ...SENTRY_SX.monoData, mt: 0.75, mb: 1, fontSize: '11px' }}
            >
              {hasDevice
                ? t('workbench.source.handoff.pointSummary', { count: pointCount })
                : t('workbench.source.handoff.noDevice')}
            </Typography>
            <Button
              fullWidth
              size="small"
              disabled={isBlocked}
              onClick={() => setActiveStep('tag')}
              data-testid="source-handoff-tag-btn"
              sx={SENTRY_SX.insetBtn}
            >
              {t('workbench.source.handoff.toTag')}
            </Button>
          </Box>
        </Box>
        <Box
          data-testid="source-incident-workboard"
          sx={{
            display: 'flex',
            minHeight: 0,
            minWidth: 0,
            overflow: 'hidden',
          }}
        >
          <Box
            data-testid="sentry-source-workspace"
            data-desk-mode={deskMode}
            sx={{
              display: 'flex',
              minHeight: 0,
              minWidth: 0,
              flex: 1,
              overflow: 'hidden',
              border: `1px solid ${sentry.border}`,
              borderRadius: '14px',
              bgcolor: alpha(sentry.canvas, 0.5),
            }}
          >
            <SourceCanvasSection />
          </Box>
        </Box>
      </Box>
    </Box>
  );
}
