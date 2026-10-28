import * as Dialog from '@radix-ui/react-dialog';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { RecordingChangeControllerView } from './useRecordingChangeController';
import type { RecordingDestinationChoice } from './recordingChangeTypes';
import { RecordingLifecycleHistory } from './RecordingLifecycleHistory';
import { activationAction, changeTargetLayout } from './recordingChangeHelpers';

const field = 'mt-1 w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-slate-900 disabled:bg-slate-100 disabled:text-slate-500';
const secondary = 'rounded-lg border border-slate-300 px-4 py-2 text-sm font-medium focus-visible:outline focus-visible:outline-2 focus-visible:outline-blue-600 disabled:opacity-50';
const primary = 'rounded-lg bg-blue-700 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 disabled:opacity-50';
export function RecordingChangeDialog({ controller: view, destinations = [] }: {
  controller: RecordingChangeControllerView;
  destinations?: RecordingDestinationChoice[];
}) {
  const { t } = useTranslation('recording-lifecycle');
  const { state, controller } = view;
  const { phase, target, review, recovery, activation } = state;
  const returnFocus = React.useRef<HTMLElement | null>(null);
  const step = phase === 'choose' ? 1 : ['impact', 'reserved', 'schema-plan'].includes(phase) ? 2 : 3;
  const locked = state.busy || controller.readonly || state.unknown;
  const title = t('dialog.title');
  const destinationName = (id: string) => destinations.find((item) => item.id === id)?.name ?? id;
  const choices = target && !destinations.some((item) => item.id === target.destination.connector_id)
    ? [{ id: target.destination.connector_id, revision: target.destination.connector_revision, name: t('dialog.original_destination') }, ...destinations] : destinations;
  const primaryAction = state.unknown && state.canRepreview ? { key: 'repreview', action: controller.repreview, disabled: state.busy || controller.readonly }
    : state.unknown ? { key: 'check_status', action: controller.checkStatus, disabled: state.busy }
    : phase === 'choose' ? { key: 'review', action: controller.review, disabled: locked || !target || !!state.error }
    : phase === 'impact' ? { key: 'reserve_preview', action: controller.reserveAndPreview, disabled: locked || !review || review.blockers.length > 0 }
    : phase === 'reserved' ? { key: 'prepare_schema', action: controller.prepareSchema, disabled: locked }
    : phase === 'schema-plan' ? { key: 'confirm_schema', action: controller.confirmSchema, disabled: locked || !state.preview }
    : phase === 'schema-confirmed' ? { key: 'review_activation', action: controller.reviewActivation, disabled: locked || !state.boundary }
    : phase === 'activation-ready' ? { key: activation && activationAction(activation.review.envelope) === 'confirmFirstStart' ? 'confirm_first_start' : 'confirm_cutover', action: controller.confirmActivation, disabled: locked }
    : undefined;
  return <>
    {!state.isOpen && (view.unresolved || state.phase === 'scheduled') && <div role="status" className="rounded-lg border border-amber-400 bg-amber-50 p-3 text-sm text-slate-900" data-testid="recording-change-recovery">
      <p>{t(state.phase === 'scheduled' ? 'dialog.scheduled' : 'dialog.recovery_available')}</p><button type="button" className={`${secondary} mt-2`} onClick={controller.resume}>{t('dialog.resume')}</button>
    </div>}
    <Dialog.Root open={state.isOpen} onOpenChange={(open) => { if (!open) controller.dismiss(); }}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-50 bg-slate-950/45 backdrop-blur-sm" />
        <Dialog.Content className="fixed left-1/2 top-1/2 z-50 flex max-h-[92dvh] w-[min(44rem,calc(100vw-1.5rem))] -translate-x-1/2 -translate-y-1/2 flex-col overflow-hidden rounded-2xl border border-slate-200 bg-white text-slate-900 shadow-2xl"
          onOpenAutoFocus={() => { returnFocus.current = document.activeElement instanceof HTMLElement ? document.activeElement : null; }}
          onCloseAutoFocus={(event) => { event.preventDefault(); returnFocus.current?.focus(); }}>
          <header className="border-b border-slate-200 px-5 py-4 sm:px-7">
            <Dialog.Title className="text-xl font-semibold">{title}</Dialog.Title>
            <Dialog.Description className="mt-1 text-sm text-slate-600">{t('dialog.description')}</Dialog.Description>
            <ol aria-label={t('dialog.steps_label')} className="mt-4 grid grid-cols-3 gap-2 text-xs sm:text-sm">
              {[1, 2, 3].map((number) => <li key={number} aria-current={step === number ? 'step' : undefined} className={`rounded-lg px-2 py-2 ${step === number ? 'bg-blue-50 font-semibold text-blue-800 ring-1 ring-blue-200' : 'bg-slate-50 text-slate-600'}`}>
                <span aria-hidden="true">{number}. </span>{t(`dialog.step_${number}`)}
              </li>)}
            </ol>
          </header>
          <div className="space-y-4 overflow-y-auto px-5 py-5 text-sm sm:px-7" aria-busy={state.busy}>
            {state.readback && <p role="status">{t('dialog.readback_received')}</p>}
            {state.error && <div role="alert" className="rounded-lg border border-amber-300 bg-amber-50 p-3">{t(`dialog.errors.${state.error}`)}</div>}
            {state.unknown && <p>{t(state.canRepreview ? 'dialog.preview_not_submitted' : 'dialog.unknown_no_replacement')}</p>}
            {!!recovery?.priorPreviewIds?.length && <details className="rounded-lg border border-slate-200 p-3"><summary>{t('dialog.previous_previews')}</summary><ul>{recovery.priorPreviewIds.map(id => <li key={id}><code>{id}</code> · {t(`dialog.operation_status.${recovery.priorPreviewStates?.[id] ?? 'unknown'}`)}</li>)}</ul></details>}
            {phase === 'choose' && target && <div className="space-y-4">
              <p>{t('dialog.saved_basis')}</p>
              <label className="block">{t('dialog.name')}<input className={field} value={target.name} disabled={locked} onChange={(event) => controller.setTarget({ ...target, name: event.target.value })} /></label>
              <label className="block">{t('dialog.layout')}<select className={field} value={target.layout_version ?? 'managed_wide_v1'} disabled={locked || !target.recording_profile || !!(target.storage_mode && target.storage_mode !== 'plain')} onChange={(event) => controller.setTarget(changeTargetLayout(target, event.target.value as typeof target.layout_version))}>
                {!target.recording_profile ? <option value="managed_wide_v1">{t('dialog.layout_basic')}</option> : <><option value="managed_wide_profile_v1">{t('dialog.layout_wide')}</option><option value="managed_long_v1">{t('dialog.layout_long')}</option></>}
              </select></label>
              <label className="block">{t('dialog.destination')}<select className={field} value={target.destination.connector_id} disabled={locked || !target.recording_profile} onChange={(event) => {
                const destination = choices.find((item) => item.id === event.target.value);
                if (destination) controller.setTarget({ ...target, destination: { ...target.destination, connector_id: destination.id, connector_revision: destination.revision, database: '' } });
              }}>{choices.map((item) => <option key={item.id} value={item.id}>{item.name}</option>)}</select></label>
              <p className="text-slate-600">{t(target.recording_profile ? 'dialog.compatibility_checked' : 'dialog.basic_same_shape')}</p>
            </div>}
            {step === 2 && <>
              <p className="rounded-lg bg-blue-50 p-3 text-blue-900">{t('dialog.history_kept')}</p>
              {review && <dl className="grid grid-cols-2 gap-2 rounded-lg bg-slate-50 p-3">
                <dt>{t('dialog.generation')}</dt><dd>{review.envelope.generation + 1}</dd>
                <dt>{t('dialog.sources')}</dt><dd>{review.envelope.request.group.members.length}</dd>
                <dt>{t('dialog.original_destination')}</dt><dd className="break-words">{state.lifecycle ? destinationName(state.lifecycle.connector_id) : '—'}</dd>
                <dt>{t('dialog.destination')}</dt><dd className="break-words">{destinationName(review.envelope.request.group.destination.connector_id)}</dd>
              </dl>}
              {review?.blockers.map((item, index) => <div role="alert" key={index} className="rounded border border-amber-300 bg-amber-50 p-3"><p>{t('dialog.blocked')}</p><p>{t(`reason.${item.reason_code}`, { defaultValue: t('dialog.blocked') })}</p><p>{t(`next.${item.next_action}`, { defaultValue: t('dialog.repair_source') })}</p></div>)}
              {review?.required_checks.length ? <p>{t('dialog.required_checks')}</p> : null}
              {recovery?.preparation && <p role="status">{t('dialog.reserved')}</p>}
              {state.preview && <section aria-label={t('dialog.schema_plan')} className="space-y-2 rounded-lg border border-blue-200 bg-blue-50/50 p-3">
                <h3 className="font-semibold">{t('dialog.schema_plan')}</h3><p>{t('dialog.schema_confirmation_required')}</p>
                {state.preview.tables.map((table) => <div key={table.name}><strong>{table.name}</strong> · {t(`dialog.schema_${table.action}`)}<p className="text-xs text-slate-600">{table.columns.join(', ')}</p></div>)}
                {state.preview.group_layout && <p>{t('dialog.columns', { count: state.preview.group_layout.columns.length })}</p>}
              </section>}
              <section className="rounded-lg border border-slate-200 p-3"><h3 className="font-medium">{t('dialog.history_details')}</h3><div className="mt-3"><RecordingLifecycleHistory history={review?.history ?? state.lifecycle?.history} /></div></section>
            </>}
            {step === 3 && <>
              <p role="status" className="rounded-lg bg-emerald-50 p-3 text-emerald-900">{t(phase === 'scheduled' ? 'dialog.scheduled' : 'dialog.schema_confirmed')}</p>
              <p>{t('dialog.observed_boundary')}</p>
              {phase !== 'scheduled' && <label className="block">{t('dialog.boundary')}<input type="datetime-local" step="1" className={field} disabled={locked} onChange={(event) => controller.setBoundary(event.target.value ? new Date(event.target.value).toISOString() : '')} /></label>}
              {(activation?.review.blockers ?? state.activationReview?.blockers ?? []).map((item, index) => <p role="alert" key={index}>{t(`reason.${item.reason_code}`, { defaultValue: t('dialog.blocked') })} {t(`next.${item.next_action}`, { defaultValue: t('dialog.repair_source') })}</p>)}
              {activation && <p>{t('dialog.reviewed_boundary')} <time>{activation.review.envelope.request.boundary}</time></p>}
              {recovery?.result && <p>{t('dialog.reviewed_boundary')} <time>{recovery.result.boundary}</time></p>}
              <p className="text-slate-600">{t('dialog.not_delivery_evidence')}</p>
              {activation && <p>{t(activationAction(activation.review.envelope) === 'confirmFirstStart' ? 'dialog.first_start_explanation' : 'dialog.cutover_explanation')}</p>}
            </>}
          </div>
          <footer className="flex flex-wrap items-center justify-between gap-3 border-t border-slate-200 bg-slate-50 px-5 py-4 sm:px-7">
            <div className="flex gap-2">
              <Dialog.Close asChild><button type="button" className={secondary}>{t(recovery ? 'dialog.close' : 'dialog.cancel')}</button></Dialog.Close>
              {phase === 'impact' && !recovery && <button type="button" className={secondary} disabled={state.busy} onClick={controller.back}>{t('dialog.back')}</button>}
              {state.unknown && state.canRepreview && <button type="button" className={secondary} disabled={state.busy} onClick={controller.checkStatus}>{t('dialog.check_status')}</button>}
              {state.unknown && state.canRetry && <button type="button" className={secondary} disabled={state.busy || controller.readonly} onClick={controller.retryExact}>{t('dialog.retry_exact')}</button>}
            </div>
            {primaryAction && <button type="button" className={primary} data-testid={`recording-change-${primaryAction.key}`} disabled={primaryAction.disabled} onClick={primaryAction.action}>{t(state.busy ? 'dialog.working' : `dialog.${primaryAction.key}`)}</button>}
          </footer>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  </>;
}
