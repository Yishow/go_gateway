import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { useTranslation } from 'react-i18next';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import type { Device } from '../../../types/datalink';
import {
  buildDeviceCapabilitySummary,
  buildDeviceEndpointSummary,
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
  parseDeviceConnectionConfig,
} from './workbenchDeviceFormModel';
import { SENTRY_SX, severityColor } from './sentrySurfaceStyles';

type Props = {
  device: Device;
};

/**
 * Sentry "System Intel" sidebar — compact 260px right column.
 *
 * Displays capabilities as metric tiles, health readiness, and connection
 * summary in a dense instrument-panel layout. No CTA buttons (those are
 * in the Diagnostics Desk center panel).
 */
export function MuiDeviceDetail({ device }: Props) {
  const { t } = useTranslation();
  const connectionConfig = parseDeviceConnectionConfig(device.connection_config);
  const endpoint = buildDeviceEndpointSummary(device.protocol, connectionConfig);
  const capabilities = buildDeviceCapabilitySummary(device.protocol, connectionConfig, t);
  const severity = severityColor(device.last_test_success);

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
      {/* Health metric */}
      <Box sx={{ ...SENTRY_SX.frostedPanel, px: 2, py: 1.5, display: 'flex', alignItems: 'center', gap: 1.5 }}>
        <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: severity, flexShrink: 0, boxShadow: `0 0 6px ${severity}66` }} />
        <Box sx={{ flex: 1, minWidth: 0 }}>
          <Typography sx={SENTRY_SX.sectionLabel}>{t('workbench.device.inspector.healthLabel')}</Typography>
          <Typography sx={{ fontWeight: 600, fontSize: '12px', color: severity, mt: 0.25 }}>
            {device.last_test_success === true
              ? t('workbench.device.inspector.healthPassing')
              : device.last_test_success === false
                ? t('workbench.device.inspector.healthFailing')
                : t('workbench.device.inspector.healthUnknown')}
          </Typography>
        </Box>
      </Box>

      {/* Status + readiness */}
      <Box sx={{ ...SENTRY_SX.frostedPanel, px: 2, py: 1.5 }}>
        <Typography sx={{ ...SENTRY_SX.sectionLabel, mb: 1 }}>{t('workbench.device.inspector.statusLabel')}</Typography>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.75 }}>
          <IntelRow label={t('workbench.device.inspector.statusFieldLabel')} value={t(getWorkbenchDeviceStatusLabelKey(device.status))} />
          <IntelRow label={t('workbench.device.inspector.protocolFieldLabel')} value={t(getWorkbenchProtocolLabelKey(device.protocol))} />
          {device.readiness_status ? <IntelRow label={t('workbench.device.inspector.readinessLabel')} value={device.readiness_status} /> : null}
        </Box>
      </Box>

      {/* Connection summary */}
      <Box sx={{ ...SENTRY_SX.frostedPanel, px: 2, py: 1.5 }}>
        <Typography sx={{ ...SENTRY_SX.sectionLabel, mb: 0.75 }}>{t('workbench.device.inspector.connectionSummary')}</Typography>
        <Typography data-testid="device-detail-endpoint" sx={{ ...SENTRY_SX.monoData, wordBreak: 'break-all' }}>
          {endpoint || '—'}
        </Typography>
      </Box>

      {/* Capability tiles */}
      {capabilities.length > 0 ? (
        <Box sx={{ ...SENTRY_SX.frostedPanel, px: 2, py: 1.5 }}>
          <Typography sx={{ ...SENTRY_SX.sectionLabel, mb: 1 }}>
            {t('workbench.device.inspector.capabilitySummary')}
          </Typography>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.75 }}>
            {capabilities.map((item) => (
              <Box key={item.id} sx={{ display: 'flex', justifyContent: 'space-between', gap: 1 }}>
                <Typography sx={{ fontSize: '10px', textTransform: 'uppercase', letterSpacing: '0.2px', color: tokens.text.muted, fontWeight: 600 }}>
                  {t(item.labelKey)}
                </Typography>
                <Typography sx={{ ...SENTRY_SX.monoData, fontSize: '11px', fontWeight: 600, textAlign: 'right' }}>
                  {item.value}
                </Typography>
              </Box>
            ))}
          </Box>
        </Box>
      ) : null}
    </Box>
  );
}

function IntelRow({ label, value }: { label: string; value: string }) {
  return (
    <Box sx={{ display: 'flex', justifyContent: 'space-between', gap: 1 }}>
      <Typography sx={{ fontSize: '10px', color: tokens.text.muted }}>{label}</Typography>
      <Typography sx={{ fontSize: '11px', fontWeight: 600, color: tokens.text.primary }}>{value}</Typography>
    </Box>
  );
}
