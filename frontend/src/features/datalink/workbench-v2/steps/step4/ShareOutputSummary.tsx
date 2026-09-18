import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import type { ModbusShareDesiredMapping, ModbusShareStatus } from '../../../../../types/modbusShare';
import type { WorkbenchV2State } from '../../state/types';
import { useModbusShareCandidateReview } from '../../../../../hooks/datalink/useModbusShareCandidateReview';
import {
  allocationRowsFromCanonicalMappings,
  LocalModbusReviewSurface,
} from './LocalModbusReviewSurface';

interface ShareOutputSummaryProps {
  state: Pick<WorkbenchV2State, 'rules' | 'points' | 'mappings'>;
  shareStatus?: ModbusShareStatus | null;
}

interface ShareOutputRow {
  id: string;
  sourceRange: string;
  outputRange: string;
}

function canonicalMappings(status: ModbusShareStatus): ModbusShareDesiredMapping[] | null {
  const mappings = status.canonical_desired_mappings ?? status.canonical_plan?.desired_mappings;
  return Array.isArray(mappings) ? mappings : null;
}

function rowsFromCanonicalMappings(
  mappings: ModbusShareDesiredMapping[],
  state: ShareOutputSummaryProps['state'],
): ShareOutputRow[] {
  const rowsByRule = new Map<string, ModbusShareDesiredMapping[]>();
  mappings.forEach((mapping) => {
    if (!Number.isInteger(mapping.share_start_register) || !Number.isInteger(mapping.span_registers) || mapping.span_registers < 1) {
      return;
    }
    const existing = rowsByRule.get(mapping.source_rule_id) ?? [];
    existing.push(mapping);
    rowsByRule.set(mapping.source_rule_id, existing);
  });

  return Array.from(rowsByRule.entries()).flatMap(([ruleID, ruleMappings]) => {
    const sorted = [...ruleMappings].sort((left, right) => left.share_start_register - right.share_start_register);
    if (sorted.length === 0) {
      return [];
    }

    const sourceLabels = sorted.map((mapping) => {
      const browserMapping = Object.values(state.mappings).find((candidate) => candidate.tag_id === mapping.tag_id);
      const point = browserMapping
        ? state.points.find((candidate) => candidate.id === browserMapping.point_id)
        : undefined;
      const rule = state.rules.find((candidate) => candidate.id === mapping.source_rule_id);
      return point?.address ?? mapping.tag_key ?? mapping.display_name ?? rule?.name ?? mapping.tag_id;
    });
    const outputStart = sorted[0].share_start_register;
    const outputEnd = Math.max(...sorted.map((mapping) => mapping.share_start_register + mapping.span_registers - 1));
    return [{
      id: ruleID,
      sourceRange: sourceLabels.length === 1 ? sourceLabels[0] : `${sourceLabels[0]} ~ ${sourceLabels[sourceLabels.length - 1]}`,
      outputRange: `${outputStart} ~ ${outputEnd}`,
    }];
  });
}

export function ShareOutputSummary({ state, shareStatus }: ShareOutputSummaryProps) {
  const { t } = useTranslation('workbench-v2');
  const hydrationReady = shareStatus?.hydration_state === 'ready' && shareStatus.readiness === true;
  const configured = shareStatus?.configured_enabled ?? (
    shareStatus?.enabled === true || shareStatus?.bind_state !== 'disabled'
  );
  const running = shareStatus?.running ?? shareStatus?.enabled ?? false;
  const failed = shareStatus?.failed ?? (
    shareStatus?.bind_state === 'fail' || shareStatus?.lifecycle_state === 'failed'
  );
  const canonical = shareStatus ? canonicalMappings(shareStatus) : null;
  const candidateBlocked = shareStatus?.candidate_set_status === 'blocked' ||
    shareStatus?.candidate_statuses?.some((status) => status.startsWith('blocked')) === true;
  const projectionFailed = failed || Boolean(shareStatus?.error) || candidateBlocked || shareStatus?.outcome === 'invalidated_unknown' ||
    shareStatus?.outcome === 'dirty_unknown' || Boolean(shareStatus?.dirty_state || shareStatus?.recovery) || (shareStatus?.invalidated_count ?? 0) > 0;
  const rows = useMemo(() => {
    if (!hydrationReady || !configured || !running || !canonical || shareStatus?.bind_state === 'disabled' || projectionFailed) {
      return [];
    }

    return rowsFromCanonicalMappings(canonical, state);
  }, [canonical, configured, hydrationReady, projectionFailed, running, shareStatus?.bind_state, state]);
  const allocationRows = useMemo(() => {
    if (!hydrationReady || !canonical) {
      return [];
    }
    return allocationRowsFromCanonicalMappings(canonical, state);
  }, [canonical, hydrationReady, state]);
  const candidateReview = useModbusShareCandidateReview(
    shareStatus?.canonical_plan ? shareStatus : null,
    shareStatus?.canonical_plan?.desired_mappings ?? null,
  );

  if (!shareStatus) {
    return (
      <section
        data-testid="step4-share-output-blocked"
        className="rounded-2xl border border-amber-500/25 bg-amber-500/5 p-4"
        role="status"
      >
        <div className="text-xs font-semibold uppercase tracking-[0.18em] text-amber-200/80">
          {t('step4.share_output_title')}
        </div>
        <p className="mt-2 text-sm text-amber-100/80">
          {t('step4.share_output_unknown')}
        </p>
      </section>
    );
  }

  // Persisted global disablement is an absolute gate. It must win over
  // missing canonical data or hydration metadata so a disabled service is
  // never misreported as blocked/unknown.
  if (!configured) {
    return (
      <section data-testid="step4-share-output-disabled" className="rounded-2xl border border-slate-700/60 bg-slate-900/40 p-4" role="status">
        <div className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
          {t('step4.share_output_title')}
        </div>
        <p className="mt-2 text-sm text-slate-400">
          {t('step4.share_output_disabled')}
        </p>
      </section>
    );
  }

  if (!hydrationReady || canonical === null) {
    return (
      <section
        data-testid="step4-share-output-blocked"
        className="rounded-2xl border border-amber-500/25 bg-amber-500/5 p-4"
        role="status"
      >
        <div className="text-xs font-semibold uppercase tracking-[0.18em] text-amber-200/80">
          {t('step4.share_output_title')}
        </div>
        <p className="mt-2 text-sm text-amber-100/80">
          {t('step4.share_output_blocked')}
        </p>
      </section>
    );
  }

  if (projectionFailed) {
    return (
      <section
        data-testid="step4-share-output-failed"
        className="rounded-2xl border border-red-500/25 bg-red-500/5 p-4"
        role="alert"
      >
        <div className="text-xs font-semibold uppercase tracking-[0.18em] text-red-200/80">
          {t('step4.share_output_title')}
        </div>
        <p className="mt-2 text-sm text-red-100/80">
          {t('step4.share_output_failed')}
        </p>
        {shareStatus && <LocalModbusReviewSurface
          rows={allocationRows}
          status={shareStatus}
          showDiagnostic
          candidatesByTagID={candidateReview.candidatesByTagID}
          candidateReviewState={candidateReview.reviewState}
        />}
      </section>
    );
  }

  if (!running) {
    return (
      <section data-testid="step4-share-output-pending" className="rounded-2xl border border-amber-500/25 bg-amber-500/5 p-4" role="status">
        <div className="text-xs font-semibold uppercase tracking-[0.18em] text-amber-200/80">
          {t('step4.share_output_title')}
        </div>
        <p className="mt-2 text-sm text-amber-100/80">
          {t('step4.share_output_pending')}
        </p>
      </section>
    );
  }

  if (rows.length === 0) {
    return null;
  }

  return (
    <section
      data-testid="step4-share-output-summary"
      className="rounded-2xl border border-violet-500/25 bg-violet-500/5 p-4"
    >
      <div className="text-xs font-semibold uppercase tracking-[0.18em] text-violet-200/80">
        {t('step4.share_output_title')}
      </div>
      <div className="mt-1 text-sm text-slate-300">
        {t('step4.share_output_subtitle')}
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
      {shareStatus && <LocalModbusReviewSurface
        rows={allocationRows}
        status={shareStatus}
        candidatesByTagID={candidateReview.candidatesByTagID}
        candidateReviewState={candidateReview.reviewState}
      />}
    </section>
  );
}
