import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type {
  WriteGroup,
  WriteGroupCreateRequest,
  WriteGroupMutationResponse,
  WriteGroupUpdateRequest,
} from '../../../../../types/studioV2WriteGroup';
import { parseStrictProfileIntervalInput } from '../../state/circuitProfileBasic';
import {
  type BatchPreviewItem,
  buildBatchCircuitItems,
  buildBatchPreview,
  hasPlanGroupRevisionDrift,
} from '../../state/circuitProfileBatch';
import type { FrozenBatchPlan } from './circuitProfileBatchExecutor';
import { CircuitProfileBatchPreview } from './CircuitProfileBatchPreview';
import { useCircuitProfileBatchRun } from './useCircuitProfileBatchRun';

export interface CircuitProfileBatchEditorProps {
  workspaceId: string;
  workspaceRevision: string;
  connectorRevision?: string;
  /** Identity of the Basic scope that owns this editor; a reviewed plan is void once it changes. */
  scopeKey?: string;
  savedLongGroups: WriteGroup[];
  readonly: boolean;
  onClose: () => void;
  /** Resolves once the parent has reloaded the server authority; the editor stays locked until then. */
  onSuccess: (updatedGroups: WriteGroup[], nextWorkspaceRevision: string) => void | Promise<void>;
  /** Read-only reload of the saved state, used to check an outcome that could not be confirmed. */
  onVerify?: () => void | Promise<void>;
  updateWriteGroup: (vars: { id: string; request: WriteGroupUpdateRequest }) => Promise<{ workspace_revision: string; group: WriteGroup }>;
  ensureCircuitProfile?: (vars: { request: WriteGroupCreateRequest }) => Promise<WriteGroupMutationResponse>;
  onOperationActiveChange?: (active: boolean) => void;
  onBusyChange?: (busy: boolean) => void;
}

export const CircuitProfileBatchEditor: React.FC<CircuitProfileBatchEditorProps> = ({
  workspaceId, workspaceRevision, connectorRevision, scopeKey, savedLongGroups, readonly,
  onClose, onSuccess, onVerify, updateWriteGroup, ensureCircuitProfile, onOperationActiveChange, onBusyChange,
}) => {
  const { t } = useTranslation('workbench-v2');
  const circuits = React.useMemo(() => buildBatchCircuitItems(savedLongGroups), [savedLongGroups]);
  const [selectedCircuitIds, setSelectedCircuitIds] = React.useState<Set<string>>(() => new Set(circuits.map((c) => c.circuit_id)));
  const [applyInst, setApplyInst] = React.useState(true);
  const [applyCum, setApplyCum] = React.useState(false);
  const [instIntervalInput, setInstIntervalInput] = React.useState('30');
  const [cumIntervalInput, setCumIntervalInput] = React.useState('10');
  const [showPreview, setShowPreview] = React.useState(false);
  const [frozenPlan, setFrozenPlan] = React.useState<FrozenBatchPlan | null>(null);
  const [previewItems, setPreviewItems] = React.useState<BatchPreviewItem[]>([]);

  const hasGroupRevisionDrift = React.useMemo(() => (
    frozenPlan ? hasPlanGroupRevisionDrift(frozenPlan.targets, savedLongGroups) : false
  ), [frozenPlan, savedLongGroups]);
  const isPlanStale = Boolean(
    frozenPlan && (
      frozenPlan.workspaceRevision !== workspaceRevision ||
      frozenPlan.workspaceId !== workspaceId ||
      frozenPlan.connectorRevision !== connectorRevision ||
      frozenPlan.scopeKey !== scopeKey ||
      hasGroupRevisionDrift
    ),
  );
  const run = useCircuitProfileBatchRun({
    plan: frozenPlan, stale: isPlanStale, readonly, connectorRevision, workspaceRevision, savedLongGroups,
    updateWriteGroup, ensureCircuitProfile, onSuccess, onClose, onVerify, onBusyChange, onOperationActiveChange,
  });

  const toggleCircuit = (cid: string) => {
    setSelectedCircuitIds((prev) => {
      const next = new Set(prev);
      if (next.has(cid)) next.delete(cid);
      else next.add(cid);
      return next;
    });
  };

  const instValidation = parseStrictProfileIntervalInput('instantaneous', instIntervalInput);
  const cumValidation = parseStrictProfileIntervalInput('cumulative', cumIntervalInput);
  const instError = applyInst && !instValidation.valid ? t('step4.basic.errors.profile_interval_invalid') : undefined;
  const cumError = applyCum && !cumValidation.valid ? t('step4.basic.errors.profile_interval_invalid') : undefined;
  const canPreview = !readonly && !run.saving && selectedCircuitIds.size > 0 && (applyInst || applyCum) && (!applyInst || instValidation.valid) && (!applyCum || cumValidation.valid);

  const handleOpenPreview = () => {
    if (!canPreview) return;
    const selectedCircuits = circuits.filter((c) => selectedCircuitIds.has(c.circuit_id));
    const targets = selectedCircuits.map((c) => ({
      circuit_id: c.circuit_id, circuit_name: c.circuit_name,
      instantaneousGroup: c.instantaneousGroup, cumulativeGroup: c.cumulativeGroup,
      anyGroup: (c.instantaneousGroup ?? c.cumulativeGroup)!,
    }));
    const items = buildBatchPreview(circuits, selectedCircuitIds, {
      applyInstantaneous: applyInst, instantaneousInterval: instValidation.intervalValue ?? 30,
      applyCumulative: applyCum, cumulativeInterval: cumValidation.intervalValue ?? 10,
    });
    setFrozenPlan({
      workspaceId, workspaceRevision, connectorRevision, scopeKey, applyInst,
      instInterval: instValidation.intervalValue ?? 30, applyCum,
      cumInterval: cumValidation.intervalValue ?? 10, targets,
    });
    setPreviewItems(items);
    run.reset();
    setShowPreview(true);
  };

  return (
    <div data-testid="circuit-profile-batch-editor" className="space-y-3 rounded-lg border border-cyan-800/80 bg-slate-950/80 p-4">
      <div className="flex items-center justify-between">
        <h4 className="text-xs font-semibold text-cyan-200">{t('step4.basic.circuit.batch_title')}</h4>
        <button type="button" data-testid="circuit-profile-batch-close" aria-label={t('step4.basic.circuit.batch_close')}
          onClick={onClose} disabled={run.saving} className="text-xs text-slate-400 hover:text-slate-200">✕</button>
      </div>

      {!showPreview ? (
        <div className="space-y-3 text-xs text-slate-300">
          <p className="text-[11px] text-slate-400">{t('step4.basic.circuit.batch_desc')}</p>
          <div className="space-y-1.5 rounded border border-slate-800 bg-slate-900/40 p-2.5">
            {circuits.map((c) => {
              const isSelected = selectedCircuitIds.has(c.circuit_id);
              return (
                <label key={c.circuit_id} className="flex items-center justify-between gap-2 py-1">
                  <div className="flex items-center gap-2">
                    <input type="checkbox" data-testid={`circuit-profile-batch-select-${c.circuit_id}`}
                      checked={isSelected} disabled={readonly || run.saving} onChange={() => toggleCircuit(c.circuit_id)} />
                    <span>{c.circuit_name || c.circuit_id}</span>
                  </div>
                  <div className="flex items-center gap-2 text-[11px] text-slate-400">
                    {c.currentInstantaneousInterval !== undefined && <span>{t('step4.basic.circuit.instantaneous_label')}: {c.currentInstantaneousInterval}s</span>}
                    {c.isOverrideInstantaneous && <span className="rounded bg-amber-950/60 px-1 text-amber-300">{t('step4.basic.circuit.batch_override')}</span>}
                  </div>
                </label>
              );
            })}
          </div>

          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1 rounded border border-slate-800 p-2">
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={applyInst} disabled={readonly || run.saving} onChange={(e) => setApplyInst(e.target.checked)} />
                <span className="font-medium text-slate-200">{t('step4.basic.circuit.batch_inst_interval')}</span>
              </label>
              <input id="circuit-profile-batch-interval-inst-input" type="text" data-testid="circuit-profile-batch-interval-inst"
                aria-label={t('step4.basic.circuit.instantaneous_interval')} aria-invalid={applyInst && !instValidation.valid}
                aria-describedby={applyInst && !instValidation.valid ? 'circuit-profile-batch-inst-error' : undefined}
                value={instIntervalInput} disabled={readonly || run.saving || !applyInst} onChange={(e) => setInstIntervalInput(e.target.value)}
                className="mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-1 text-xs text-slate-100 disabled:opacity-40" />
              {instError && <p id="circuit-profile-batch-inst-error" role="alert" className="text-xs text-rose-300">{instError}</p>}
            </div>
            <div className="space-y-1 rounded border border-slate-800 p-2">
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={applyCum} disabled={readonly || run.saving} onChange={(e) => setApplyCum(e.target.checked)} />
                <span className="font-medium text-slate-200">{t('step4.basic.circuit.batch_cum_interval')}</span>
              </label>
              <input id="circuit-profile-batch-interval-cum-input" type="text" data-testid="circuit-profile-batch-interval-cum"
                aria-label={t('step4.basic.circuit.cumulative_interval')} aria-invalid={applyCum && !cumValidation.valid}
                aria-describedby={applyCum && !cumValidation.valid ? 'circuit-profile-batch-cum-error' : undefined}
                value={cumIntervalInput} disabled={readonly || run.saving || !applyCum} onChange={(e) => setCumIntervalInput(e.target.value)}
                className="mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-1 text-xs text-slate-100 disabled:opacity-40" />
              {cumError && <p id="circuit-profile-batch-cum-error" role="alert" className="text-xs text-rose-300">{cumError}</p>}
            </div>
          </div>

          <div className="flex items-center justify-end gap-2 pt-2">
            <button type="button" data-testid="circuit-profile-batch-cancel" onClick={onClose} disabled={run.saving} className="rounded border border-slate-700 px-3 py-1.5 text-xs text-slate-300 hover:bg-slate-800">
              {t('step4.basic.circuit.batch_cancel')}
            </button>
            <button
              type="button"
              data-testid="circuit-profile-batch-preview-btn"
              disabled={!canPreview}
              onClick={handleOpenPreview}
              className="rounded bg-cyan-600 px-3 py-1.5 text-xs font-semibold text-white hover:bg-cyan-500 disabled:opacity-40"
            >
              {t('step4.basic.circuit.batch_preview')}
            </button>
          </div>
        </div>
      ) : (
        <div data-testid="circuit-profile-batch-preview" className="space-y-3 text-xs text-slate-300">
          <CircuitProfileBatchPreview
            items={previewItems} results={run.results} stale={isPlanStale && !run.finished} errorMessage={run.errorMessage}
            hasUnknown={run.hasUnknown} refreshFailed={run.refreshFailed} verifying={run.verifying}
            onVerify={onVerify ? () => void run.verify() : undefined}
          />
          <div className="flex items-center justify-end gap-2 pt-2">
            <button type="button" data-testid="circuit-profile-batch-cancel" disabled={run.saving || run.hasUnknown}
              onClick={() => { run.reset(); setShowPreview(false); }} className="rounded border border-slate-700 px-3 py-1.5 text-xs text-slate-300 hover:bg-slate-800 disabled:opacity-40">
              {t('step4.basic.circuit.batch_cancel')}
            </button>
            <button type="button" data-testid="circuit-profile-batch-confirm" disabled={readonly || run.saving || run.finished || isPlanStale}
              onClick={() => void run.confirm()} className="rounded bg-cyan-600 px-3 py-1.5 text-xs font-semibold text-white hover:bg-cyan-500 disabled:opacity-40">
              {run.saving ? t('step4.basic.circuit.batch_saving') : t('step4.basic.circuit.batch_confirm')}
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
