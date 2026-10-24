import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useCreateStudioV2MeasurementMutation, useUpdateStudioV2MeasurementMutation } from '../../../../../hooks/datalink/useStudioV2WorkspaceMeasurements';
import type { MeasurementDefinition } from '../../../../../types/measurement';
import { getSafeErrorMessage } from '../../../../../utils/typedErrors';
import type { CandidateCounterStatus } from '../../state/circuitProfileBasic';
import { candidateKey } from '../../state/circuitProfileBasic';

export interface CircuitCounterConfirmationProps {
  workspaceId: string;
  workspaceRevision: string;
  circuitId: string;
  unconfirmedStatuses: CandidateCounterStatus[];
  readonly: boolean;
  onConfirmed: () => Promise<void> | void;
}

// Each source owns its unit and exact policy; repairing one must preserve others.
function CounterSourceConfirmation({ status, workspaceId, workspaceRevision, circuitId, readonly, onConfirmed }:
  Omit<CircuitCounterConfirmationProps, 'unconfirmedStatuses'> & { status: CandidateCounterStatus }) {
  const { t } = useTranslation('workbench-v2');
  const createMutation = useCreateStudioV2MeasurementMutation();
  const updateMutation = useUpdateStudioV2MeasurementMutation();
  const saved = status.measurement;
  const [semantics, setSemantics] = React.useState<'counter' | 'signed_counter'>(saved?.semantic_kind === 'signed_counter' ? 'signed_counter' : 'counter');
  const [unit, setUnit] = React.useState(saved?.unit ?? '');
  const [errorMessage, setErrorMessage] = React.useState<string>();
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const inFlight = React.useRef(false);
  const canSubmit = !readonly && !isSubmitting && Boolean(workspaceId && workspaceRevision && unit.trim() && (saved || circuitId));
  if (status.issueCode === 'counter-type-unsupported') {
    return <p role="status">{status.candidate.label.trim() || status.candidate.address} · {t('workspace_readiness.issue_labels.counter-type-unsupported')}</p>;
  }
  if (saved && ![saved.id, saved.definition_revision, saved.source_binding_revision, saved.series_epoch].every((part) => part?.trim())) {
    return <div role="status" data-testid="circuit-counter-binding-unconfirmed" className="space-y-2">
      <p>{status.candidate.label.trim() || status.candidate.address} · {t('step4.basic.circuit.binding_unconfirmed')}</p>
      <button type="button" data-testid="circuit-counter-reload-source" disabled={readonly} onClick={() => void onConfirmed()}>{t('step4.group.editor.reload')}</button>
    </div>;
  }

  const handleConfirm = async () => {
    if (!canSubmit || inFlight.current) return;
    inFlight.current = true;
    setErrorMessage(undefined);
    setIsSubmitting(true);
    try {
      const candidate = status.candidate;
      const def: Partial<MeasurementDefinition> = {
        ...saved,
        workspace_id: workspaceId,
        name: saved?.name || candidate.label || candidate.tag_key,
        equipment_id: saved?.equipment_id || circuitId,
        device_id: candidate.device_id,
        point_id: candidate.point_id,
        tag_id: candidate.tag_id,
        quantity: saved?.quantity ?? '',
        semantic_kind: semantics,
        unit: unit.trim(),
        counter_policy: { ...saved?.counter_policy, allow_negative: semantics === 'signed_counter' },
      };
      if (saved) await updateMutation.mutateAsync({ id: saved.id, def });
      else await createMutation.mutateAsync(def);
      await onConfirmed();
    } catch (error) {
      setErrorMessage(getSafeErrorMessage(error, t).message);
    } finally {
      inFlight.current = false;
      setIsSubmitting(false);
    }
  };

  return (
    <fieldset data-testid={`circuit-counter-confirmation-source-${status.candidate.tag_id}`} disabled={readonly || isSubmitting} className="space-y-2 rounded border border-amber-700/50 p-3">
      <legend className="px-1 font-medium">{status.candidate.label.trim() || status.candidate.address}</legend>
      <p>{t(`workspace_readiness.issue_labels.${status.issueCode ?? 'counter-semantics-required'}`)}</p>
      <div className="grid gap-2 sm:grid-cols-2">
        <label className="block">
          <span>{t('step4.basic.circuit.counter_semantics_label')}</span>
          <select data-testid="circuit-confirm-counter-semantics" value={semantics}
            onChange={(event) => setSemantics(event.target.value === 'signed_counter' ? 'signed_counter' : 'counter')}
            className="mt-1 w-full rounded border border-amber-700/60 bg-slate-950 px-2 py-1.5 text-slate-100">
            <option value="counter">{t('step4.basic.circuit.counter_semantics_counter')}</option>
            <option value="signed_counter">{t('step4.basic.circuit.counter_semantics_signed')}</option>
          </select>
        </label>
        <label className="block">
          <span>{t('step4.basic.circuit.counter_unit_label')}</span>
          <input type="text" data-testid="circuit-confirm-counter-unit" value={unit} onChange={(event) => setUnit(event.target.value)}
            className="mt-1 w-full rounded border border-amber-700/60 bg-slate-950 px-2 py-1.5 text-slate-100" />
        </label>
      </div>
      {errorMessage && <p role="alert" data-testid="circuit-confirm-counter-error" className="text-rose-300">{errorMessage}</p>}
      <button type="button" data-testid="circuit-confirm-counter-submit" disabled={!canSubmit} onClick={() => void handleConfirm()}
        className="rounded border border-amber-500/60 bg-amber-900/40 px-3 py-1.5 disabled:opacity-40">
        {isSubmitting ? t('step4.basic.circuit.counter_confirming') : t('step4.basic.circuit.counter_confirm_action')}
      </button>
    </fieldset>
  );
}

export const CircuitCounterConfirmation: React.FC<CircuitCounterConfirmationProps> = ({ unconfirmedStatuses, ...props }) => {
  const { t } = useTranslation('workbench-v2');
  if (!unconfirmedStatuses.length) return null;
  return (
    <section data-testid="circuit-counter-confirmation" role="region" aria-label={t('step4.basic.circuit.counter_confirm_title')}
      className="space-y-3 rounded border border-amber-500/40 bg-amber-950/20 p-3 text-xs text-amber-100">
      <p className="font-semibold">{t('step4.basic.circuit.counter_confirm_title')}</p>
      <p>{t('step4.basic.circuit.counter_confirm_note')}</p>
      {unconfirmedStatuses.map((status) => <CounterSourceConfirmation key={candidateKey(status.candidate)} status={status} {...props} />)}
    </section>
  );
};
