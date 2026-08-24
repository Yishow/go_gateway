import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { computeShareLayout } from '../../state/sourceRule';
import type { Rule, WorkbenchV2State } from '../../state/types';

interface ShareOutputSummaryProps {
  state: Pick<WorkbenchV2State, 'devices' | 'rules' | 'points' | 'settings'>;
}

interface ShareOutputRow {
  id: string;
  sourceRange: string;
  outputRange: string;
}

function sourceRangeForRule(rule: Rule, state: ShareOutputSummaryProps['state'], start: number, stride: number): ShareOutputRow | null {
  const points = state.points.filter((point) => (
    point.rule_id === rule.id && point.enabled && !point.skipped
  ));
  if (points.length === 0) {
    return null;
  }

  const outputEnd = start + (points.length - 1) * stride;
  return {
    id: rule.id,
    sourceRange: `${points[0].address} ~ ${points[points.length - 1].address}`,
    outputRange: `${start} ~ ${outputEnd}`,
  };
}

export function ShareOutputSummary({ state }: ShareOutputSummaryProps) {
  const { t } = useTranslation('workbench-v2');
  const rows = useMemo(() => {
    if (!state.settings.modbus_share.enabled) {
      return [];
    }

    const nonModbusRuleIDs = new Set(
      state.devices
        .filter((device) => !device.protocol.startsWith('modbus_'))
        .map((device) => device.id),
    );
    const layouts = computeShareLayout(state.rules, state.settings.modbus_share.base_register);

    return state.rules
      .filter((rule) => rule.enabled && rule.share_enabled && nonModbusRuleIDs.has(rule.device_id))
      .map((rule) => {
        const layout = layouts[rule.id];
        return layout ? sourceRangeForRule(rule, state, layout.start, layout.stride) : null;
      })
      .filter((row): row is ShareOutputRow => row !== null);
  }, [state]);

  if (rows.length === 0) {
    return null;
  }

  return (
    <section
      data-testid="step4-share-output-summary"
      className="rounded-2xl border border-violet-500/25 bg-violet-500/5 p-4"
    >
      <div className="text-xs font-semibold uppercase tracking-[0.18em] text-violet-200/80">
        {t('step4.share_output_title', 'Local Share output')}
      </div>
      <div className="mt-1 text-sm text-slate-300">
        {t('step4.share_output_subtitle', 'Non-Modbus source points are exposed through holding registers independently of database targets.')}
      </div>
      <div className="mt-3 space-y-2">
        {rows.map((row) => (
          <div key={row.id} className="flex flex-wrap items-center gap-2 text-sm" data-testid={`step4-share-output-row-${row.id}`}>
            <span className="font-mono text-slate-300">{row.sourceRange}</span>
            <span className="text-violet-200" aria-hidden="true">→</span>
            <span className="font-mono font-semibold text-violet-100">{row.outputRange}</span>
          </div>
        ))}
      </div>
    </section>
  );
}
