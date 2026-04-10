import Box from '@mui/material/Box';
import ButtonBase from '@mui/material/ButtonBase';
import Typography from '@mui/material/Typography';
import { useTranslation } from 'react-i18next';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import type { Device } from '../../../types/datalink';
import {
  buildDeviceCapabilitySummary,
  buildDeviceEndpointSummary,
  getWorkbenchProtocolLabelKey,
  parseDeviceConnectionConfig,
} from './workbenchDeviceFormModel';
import { SENTRY_SX, severityColor } from './sentrySurfaceStyles';

type Props = {
  devices: Device[];
  selectedDeviceId: string | null;
  onSelect: (device: Device) => void;
};

/**
 * Sentry-style incident feed — narrow vertical rail of event cards.
 *
 * Each device is rendered as a compact incident event with a left severity bar
 * (green/red/neutral) based on last test health. This is fundamentally different
 * from the baseline card list: compact, severity-first, monospace endpoints.
 */
export function MuiDeviceList({ devices, selectedDeviceId, onSelect }: Props) {
  const { t } = useTranslation();

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5 }}>
      <Typography sx={{ ...SENTRY_SX.sectionLabel, px: 1, mb: 0.5 }}>
        {t('workbench.device.eyebrow')} ({devices.length})
      </Typography>
      {devices.map((device) => {
        const isSelected = device.id === selectedDeviceId;
        const connectionConfig = parseDeviceConnectionConfig(device.connection_config);
        const endpoint = buildDeviceEndpointSummary(device.protocol, connectionConfig);
        const severity = severityColor(device.last_test_success);
        const capabilities = buildDeviceCapabilitySummary(device.protocol, connectionConfig, t).slice(0, 2);
        const healthLabel = device.last_test_success === true
          ? t('workbench.device.card.testPassed')
          : device.last_test_success === false
            ? t('workbench.device.card.testFailed')
            : t('workbench.device.card.notTested');

        return (
          <ButtonBase
            key={device.id}
            component="div"
            data-testid={`device-row-${device.id}`}
            aria-label={device.name}
            aria-pressed={isSelected}
            onClick={() => onSelect(device)}
            sx={{
              display: 'flex',
              flexDirection: 'column',
              gap: 0.25,
              width: '100%',
              textAlign: 'left',
              px: 1.25,
              py: 1,
              transition: 'all 0.15s',
              bgcolor: isSelected ? `${tokens.archetype.sentry.accent}22` : 'transparent',
              ...SENTRY_SX.severityBar(isSelected ? tokens.archetype.sentry.accent : severity),
              '&:hover': {
                bgcolor: isSelected ? `${tokens.archetype.sentry.accent}22` : tokens.treatment.glassPanel.hoverFill,
              },
            }}
          >
            <Typography noWrap sx={{ fontWeight: 600, fontSize: '12px', lineHeight: 1.3, width: '100%' }}>
              {device.name}
            </Typography>
            <Typography noWrap sx={{ ...SENTRY_SX.monoData, fontSize: '10px', width: '100%' }}>
              {endpoint}
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, width: '100%' }}>
              <Typography sx={{ fontSize: '10px', textTransform: 'uppercase', letterSpacing: '0.2px', color: tokens.text.muted, fontWeight: 600 }}>
                {t(getWorkbenchProtocolLabelKey(device.protocol))}
              </Typography>
              <Box sx={{ flex: 1 }} />
              <Typography data-testid={`device-health-${device.id}`} sx={{ fontSize: '10px', color: severity, fontWeight: 600 }}>
                {healthLabel}
              </Typography>
              <Box sx={{ width: 6, height: 6, borderRadius: '50%', bgcolor: severity, flexShrink: 0 }} />
            </Box>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.75, mt: 0.25 }}>
              {capabilities.map((capability) => (
                <Typography key={capability.id} sx={{ fontSize: '10px', color: tokens.text.secondary }}>
                  {t(capability.labelKey)}
                </Typography>
              ))}
            </Box>
          </ButtonBase>
        );
      })}
    </Box>
  );
}
