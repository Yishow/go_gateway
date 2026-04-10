import Alert from '@mui/material/Alert';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import EditIcon from '@mui/icons-material/Edit';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import HelpOutlineIcon from '@mui/icons-material/HelpOutline';
import { useTranslation } from 'react-i18next';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import type { ConnectionTestStageResult, Device } from '../../../types/datalink';
import type { DeviceTestHistoryEntry } from './workbenchTypes';
import {
  buildDeviceEndpointSummary,
  getWorkbenchProtocolLabelKey,
  parseDeviceConnectionConfig,
} from './workbenchDeviceFormModel';
import { SENTRY_SX } from './sentrySurfaceStyles';

type Props = {
  device: Device | null;
  recentTests: ReadonlyArray<DeviceTestHistoryEntry>;
  onContinue: () => void;
  onEdit?: () => void;
  onClone?: () => void;
};

function StageStatusIcon({ status }: { status: ConnectionTestStageResult['status'] | undefined }) {
  if (status === 'success') return <CheckCircleIcon sx={{ fontSize: 28, color: tokens.status.success }} />;
  if (status === 'failed') return <ErrorIcon sx={{ fontSize: 28, color: tokens.status.error }} />;
  return <HelpOutlineIcon sx={{ fontSize: 28, color: tokens.text.muted }} />;
}

function stageStatusColor(status: ConnectionTestStageResult['status'] | undefined): string {
  if (status === 'success') return tokens.status.success;
  if (status === 'failed') return tokens.status.error;
  return tokens.border.strong;
}

/**
 * Sentry Diagnostics Desk — first-class center-stage diagnostics.
 *
 * Two large stage tiles (Connect / Probe) displayed prominently with oversized
 * status indicators, latency metrics in monospace, and action commands below.
 * This replaces the old buried-in-editor diagnostics panel.
 */
export function MuiDeviceDiagnosticsPanel({ device, recentTests, onContinue, onEdit, onClone }: Props) {
  const { t } = useTranslation();

  if (!device) {
    return (
      <Box sx={{ ...SENTRY_SX.frostedPanel, p: 4, textAlign: 'center', minHeight: 200, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
        <Typography sx={SENTRY_SX.sectionLabel}>{t('workbench.device.inspector.emptyTitle')}</Typography>
        <Typography variant="body2" sx={{ color: tokens.text.secondary, mt: 1.5, maxWidth: 280 }}>
          {t('workbench.device.inspector.emptyDescription')}
        </Typography>
      </Box>
    );
  }

  const connectionConfig = parseDeviceConnectionConfig(device.connection_config);
  const endpoint = buildDeviceEndpointSummary(device.protocol, connectionConfig);
  const lastTest = recentTests[0] ?? null;
  const connectStage = lastTest?.phaseDetails?.connect;
  const probeStage = lastTest?.phaseDetails?.probe;

  return (
    <Box data-testid="device-detail-panel" sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
      {/* Device identity bar */}
      <Box sx={{ ...SENTRY_SX.frostedPanel, px: 2, py: 1.5, display: 'flex', alignItems: 'center', gap: 1.5 }}>
        <Box sx={{ flex: 1, minWidth: 0 }}>
          <Typography sx={{ fontWeight: 700, fontSize: '14px' }}>{device.name}</Typography>
          <Typography data-testid="device-detail-endpoint" sx={SENTRY_SX.monoData}>{endpoint}</Typography>
        </Box>
        <Typography sx={{ ...SENTRY_SX.sectionLabel, flexShrink: 0 }}>
          {t(getWorkbenchProtocolLabelKey(device.protocol))}
        </Typography>
        {onEdit ? <Button size="small" aria-label={`${t('workbench.device.actions.edit')}-diagnostics`} startIcon={<EditIcon sx={{ fontSize: 14 }} />} sx={SENTRY_SX.insetBtn} onClick={onEdit}>{t('workbench.device.actions.edit')}</Button> : null}
        {onClone ? <Button size="small" aria-label={`${t('workbench.device.actions.clone')}-diagnostics`} startIcon={<FileCopyIcon sx={{ fontSize: 14 }} />} sx={SENTRY_SX.insetBtn} onClick={onClone}>{t('workbench.device.actions.clone')}</Button> : null}
      </Box>

      {/* Diagnostics stage tiles */}
      <Box sx={{ display: 'grid', gap: 1.5, gridTemplateColumns: '1fr 1fr' }}>
        <StageTile label={t('workbench.device.inspector.phases.connect')} stage={connectStage} />
        <StageTile label={t('workbench.device.inspector.phases.probe')} stage={probeStage} />
      </Box>

      {/* Recent test timeline */}
      {recentTests.length > 0 ? (
        <Box sx={{ ...SENTRY_SX.frostedPanel, px: 2, py: 1.5 }}>
          <Typography sx={{ ...SENTRY_SX.sectionLabel, mb: 1 }}>{t('workbench.device.inspector.recentTests')}</Typography>
          <Box sx={{ display: 'flex', gap: 1 }}>
            {recentTests.slice(0, 3).map((entry) => (
              <Box key={entry.id} sx={{ flex: 1, px: 1, py: 0.75, bgcolor: tokens.archetype.sentry.elevated, borderRadius: tokens.radius.sm, border: `1px solid ${entry.success ? tokens.status.success : tokens.status.error}33` }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                  <Box sx={{ width: 6, height: 6, borderRadius: '50%', bgcolor: entry.success ? tokens.status.success : tokens.status.error }} />
                  <Typography sx={{ fontSize: '10px', fontWeight: 600, color: entry.success ? tokens.status.success : tokens.status.error }}>
                    {entry.success ? t('workbench.commit.pass') : t('workbench.commit.fail')}
                  </Typography>
                </Box>
                {entry.latencyMs != null ? <Typography sx={{ ...SENTRY_SX.monoData, fontSize: '10px', mt: 0.25 }}>{entry.latencyMs}ms</Typography> : null}
              </Box>
            ))}
          </Box>
        </Box>
      ) : (
        <Box sx={{ ...SENTRY_SX.frostedPanel, px: 2, py: 1.5, textAlign: 'center' }}>
          <Typography sx={{ ...SENTRY_SX.sectionLabel, mb: 0.5 }}>DIAGNOSTICS</Typography>
          <Typography sx={{ fontSize: '11px', color: tokens.text.muted }}>
            {t('workbench.device.inspector.emptyDescription')}
          </Typography>
        </Box>
      )}

      {lastTest?.canActivate === false ? (
        <Alert severity="warning" variant="outlined">{t('workbench.device.inspector.activationBlocked')}</Alert>
      ) : null}

      {/* Actions */}
      <Box sx={{ display: 'flex', gap: 1 }}>
        <Button fullWidth startIcon={<ArrowForwardIcon />} onClick={onContinue} sx={{ ...SENTRY_SX.insetBtn, bgcolor: tokens.archetype.sentry.accentMuted, borderColor: tokens.archetype.sentry.accent, color: tokens.text.primary }}>
          {t('workbench.device.actions.continue')}
        </Button>
      </Box>
    </Box>
  );
}

function StageTile({ label, stage }: { label: string; stage: ConnectionTestStageResult | undefined }) {
  const { t } = useTranslation();
  const color = stageStatusColor(stage?.status);
  const hasResult = Boolean(stage);

  return (
    <Box data-testid={`diagnostics-stage-${label.toLowerCase()}`} sx={{ ...SENTRY_SX.stageTile, borderTopColor: color, borderTopWidth: 2 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <Typography sx={SENTRY_SX.sectionLabel}>{label}</Typography>
        <StageStatusIcon status={stage?.status} />
      </Box>
      {hasResult && stage ? (
        <>
          <Typography sx={{ ...SENTRY_SX.metricValue, color }}>
            {stage.latency_ms}ms
          </Typography>
          <Typography sx={{ fontSize: '11px', color: tokens.text.secondary }}>
            {stage.error || stage.message || '—'}
          </Typography>
        </>
      ) : (
        <Typography sx={{ ...SENTRY_SX.monoData, py: 1 }}>{t('workbench.device.inspector.awaitingStageTest')}</Typography>
      )}
    </Box>
  );
}
