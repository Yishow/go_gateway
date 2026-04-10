import { type FormEvent } from 'react';
import Alert from '@mui/material/Alert';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import SendIcon from '@mui/icons-material/Send';
import CableIcon from '@mui/icons-material/Cable';
import CloseIcon from '@mui/icons-material/Close';
import IconButton from '@mui/material/IconButton';
import { useTranslation } from 'react-i18next';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import type { ConnectionTestResult, ProtocolType } from '../../../types/datalink';
import {
  createDefaultDeviceConnectionConfig,
  getWorkbenchProtocolLabelKey,
  WORKBENCH_PROTOCOLS,
  type DeviceConnectionConfig,
  type DeviceDraft,
} from './workbenchDeviceFormModel';
import type { DevicePanelState } from './workbenchTypes';
import { MuiDeviceConnectionFields } from './MuiDeviceConnectionFields';
import { SENTRY_SX } from './sentrySurfaceStyles';

type FieldErrorMap = Record<string, string>;

type Props = {
  panelState: NonNullable<DevicePanelState>;
  draft: DeviceDraft;
  setDraft: React.Dispatch<React.SetStateAction<DeviceDraft>>;
  fieldErrors: FieldErrorMap;
  setFieldErrors: React.Dispatch<React.SetStateAction<FieldErrorMap>>;
  draftTestResult: ConnectionTestResult | null;
  isSaving: boolean;
  isTesting: boolean;
  onClose: () => void;
  onSubmit: (e: FormEvent<HTMLFormElement>) => void;
  onTestConnection: () => void;
};

const s = tokens.archetype.sentry;

const MODE_META: Record<string, { labelKey: string; color: string }> = {
  create: { labelKey: 'workbench.device.panel.createTitle', color: s.highlight },
  edit: { labelKey: 'workbench.device.panel.editTitle', color: s.accent },
  clone: { labelKey: 'workbench.device.panel.cloneTitle', color: s.warm },
};

/**
 * Sentry "Task Workspace" — replaces center + right zones in task mode.
 *
 * Styled as a command form with uppercase section labels, inset buttons,
 * and an embedded diagnostics result area instead of the old Paper card.
 */
export function MuiDeviceEditor({
  panelState, draft, setDraft, fieldErrors, setFieldErrors,
  draftTestResult, isSaving, isTesting,
  onClose, onSubmit, onTestConnection,
}: Props) {
  const { t } = useTranslation();
  const mode = panelState.mode;
  const meta = MODE_META[mode] ?? MODE_META.create;

  const onChangeValue = (key: string, value: DeviceConnectionConfig[string]) => {
    setDraft((d) => ({ ...d, connectionConfig: { ...d.connectionConfig, [key]: value } }));
  };
  const onChangeNumber = (key: string, raw: string) => {
    if (raw === '') { onChangeValue(key, undefined); return; }
    if (/^\d+$/.test(raw)) onChangeValue(key, Number.parseInt(raw, 10));
  };
  const onChangeTopics = (raw: string) => {
    onChangeValue('topics', raw.split(/\n|,/).map((s) => s.trim()).filter(Boolean));
  };

  return (
    <Box data-testid="device-inline-editor" component="form" onSubmit={onSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, height: '100%' }}>
      {/* Task header */}
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
        <Box sx={{ width: 4, height: 20, borderRadius: 1, bgcolor: meta.color, flexShrink: 0 }} />
        <Typography component="h2" variant="h6" sx={{ ...SENTRY_SX.modeChip, color: meta.color }}>
          {t(meta.labelKey)}
        </Typography>
        <Box sx={{ flex: 1 }} />
        <IconButton size="small" onClick={onClose} sx={SENTRY_SX.insetIconBtn}>
          <CloseIcon sx={{ fontSize: 16 }} />
        </IconButton>
      </Box>

      {/* Scrollable form body */}
      <Box sx={{ flex: 1, minHeight: 0, overflow: 'auto', display: 'flex', flexDirection: 'column', gap: 2 }}>
        {/* Identity fields */}
        <Box sx={{ display: 'grid', gap: 1.5, gridTemplateColumns: '1fr 1fr' }}>
          <TextField label={t('workbench.device.fields.name')} size="small" value={draft.name} onChange={(e) => setDraft((d) => ({ ...d, name: e.target.value }))} error={Boolean(fieldErrors.name)} helperText={fieldErrors.name} autoFocus />
          <TextField label={t('workbench.device.fields.protocol')} size="small" select SelectProps={{ native: true }} value={draft.protocol} disabled={mode !== 'create'} onChange={(e) => { const p = e.target.value as ProtocolType; setDraft((d) => ({ ...d, protocol: p, connectionConfig: createDefaultDeviceConnectionConfig(p) })); setFieldErrors({}); }}>
            {WORKBENCH_PROTOCOLS.map((p) => <option key={p} value={p}>{t(getWorkbenchProtocolLabelKey(p))}</option>)}
          </TextField>
        </Box>
        <TextField label={t('workbench.device.fields.description')} size="small" fullWidth multiline minRows={2} value={draft.description} onChange={(e) => setDraft((d) => ({ ...d, description: e.target.value }))} />

        {/* Connection section */}
        <Box sx={{ ...SENTRY_SX.frostedPanel, p: 2, display: 'flex', flexDirection: 'column', gap: 1.5 }}>
          <Typography sx={SENTRY_SX.sectionLabel}>
            {t('workbench.device.connection.title')}
          </Typography>
          <Alert severity="info" variant="outlined" sx={{ py: 0.25, fontSize: '11px' }}>
            {t('workbench.device.connection.backendHostHint')}
          </Alert>
          <MuiDeviceConnectionFields
            protocol={draft.protocol}
            connectionConfig={draft.connectionConfig}
            fieldErrors={fieldErrors}
            onChangeValue={onChangeValue}
            onChangeNumber={onChangeNumber}
            onChangeTopics={onChangeTopics}
          />
        </Box>

        {/* Embedded draft test result */}
        {draftTestResult ? (
          <Box sx={{ ...SENTRY_SX.frostedPanel, p: 2 }}>
            <Typography sx={{ ...SENTRY_SX.sectionLabel, mb: 1 }}>TEST RESULT</Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
              <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: draftTestResult.success ? tokens.status.success : tokens.status.error }} />
              <Typography sx={{ ...SENTRY_SX.metricValue, fontSize: '16px', color: draftTestResult.success ? tokens.status.success : tokens.status.error }}>
                {draftTestResult.success ? 'PASS' : 'FAIL'}
              </Typography>
              {draftTestResult.latency_ms != null ? (
                <Typography sx={SENTRY_SX.monoData}>{draftTestResult.latency_ms}ms</Typography>
              ) : null}
            </Box>
            {draftTestResult.error ? <Typography sx={{ fontSize: '11px', color: tokens.status.error, mt: 1 }}>{draftTestResult.error}</Typography> : null}
          </Box>
        ) : null}
      </Box>

      {/* Footer action strip */}
      <Box sx={{ display: 'flex', gap: 1, pt: 1, borderTop: `1px solid ${s.border}` }}>
        <Button size="small" disabled={isTesting} startIcon={isTesting ? <CircularProgress size={14} /> : <CableIcon sx={{ fontSize: 14 }} />} onClick={onTestConnection} type="button" data-testid="device-test-connection-btn" sx={SENTRY_SX.insetBtn}>
          {t('workbench.device.actions.testDraftConnection')}
        </Button>
        <Box sx={{ flex: 1 }} />
        <Button size="small" onClick={onClose} sx={{ ...SENTRY_SX.insetBtn, color: tokens.text.secondary }}>
          {t('common.cancel')}
        </Button>
        <Button size="small" disabled={isSaving} startIcon={isSaving ? <CircularProgress size={14} /> : <SendIcon sx={{ fontSize: 14 }} />} type="submit" sx={{ ...SENTRY_SX.insetBtn, bgcolor: s.accentMuted, borderColor: s.accent, color: tokens.text.primary }}>
          {t('workbench.device.actions.save')}
        </Button>
      </Box>
    </Box>
  );
}
