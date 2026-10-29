import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroupSampleReduction, WriteGroupStorageMode, WriteGroupValueRole } from '../../../../../types/studioV2WriteGroup';
import {
  candidateKey,
  type CandidateCounterStatus,
  type CircuitProfileLocalDraft,
  type CircuitPurpose,
  defaultMemberRole,
  deriveCircuitId,
  type ParsedProfileInterval,
} from '../../state/circuitProfileBasic';
import type { CircuitLayoutBlockReason, CircuitLayoutOptions } from '../../state/circuitProfileCapabilities';
import type { GroupCandidate } from '../../state/writeGroup/candidates';
import { CircuitCounterConfirmation } from './CircuitCounterConfirmation';
import { CircuitPurposeFields } from './CircuitPurposeFields';
import { CircuitProfileRoleOptions } from './CircuitProfileRoleOptions';
import { CircuitRequiredPolicy } from './CircuitRequiredPolicy';
import { DeviceScopeTiles } from './DeviceScopeTiles';
import { SourceReview } from './SourceReview';

export interface CircuitProfileConfigFieldsProps {
  workspaceId?: string;
  workspaceRevision: string;
  readonly: boolean;
  savedDevices: Array<{ id: string; name?: string }>;
  activeDeviceId: string;
  onSelectDevice: (deviceId: string) => void;
  deviceCandidates: GroupCandidate[];
  localDraft: CircuitProfileLocalDraft;
  updateLocalDraft: (updater: (prev: CircuitProfileLocalDraft) => CircuitProfileLocalDraft) => void;
  instParsed: ParsedProfileInterval;
  cumParsed: ParsedProfileInterval;
  isPostgres: boolean;
  isMySQL?: boolean;
  isSQLServer?: boolean;
  nativeMonthlyAvailable?: boolean;
  storageLocked: boolean;
  effectiveStorageMode: WriteGroupStorageMode;
  cumulativeCircuitMissing: boolean;
  /** A new counter definition needs a meaningful circuit label, not just a draft ID. */
  counterNameMissing?: boolean;
  circuitId: string;
  unconfirmedCounterStatuses: CandidateCounterStatus[];
  onCounterConfirmed: (nextWorkspaceRev?: string) => Promise<void> | void;
  savedCircuits?: Array<{ circuit_id: string; circuit_name: string }>;
  onSelectCircuit?: (circuitId: string) => void;
  onOpenBatch?: () => void;
  /** The layout this circuit records into, and whether a saved group already fixes it. */
  purpose?: CircuitPurpose;
  purposeLocked?: boolean;
  layoutOptions?: CircuitLayoutOptions;
  layoutBlockedReason?: CircuitLayoutBlockReason;
  /** A reviewed cutover prepares long targets only, so it offers no purpose to choose. */
  showPurpose?: boolean;
  /** `mean` for a new instantaneous profile; a saved profile shows the reduction it has (`latest` when none). */
  instantaneousReduction?: WriteGroupSampleReduction;
  /** Sources a saved instantaneous group already holds: their role is the saved one. */
  savedRoleKeys?: ReadonlySet<string>;
  memberRole?: (candidate: GroupCandidate) => WriteGroupValueRole;
  /** Saved members of this circuit that no current source matches, and the saved group that holds them. */
  unavailableSources?: { count: number; groupId?: string };
  /** The saved wide group's required sources; shown for a wide circuit only. */
  requiredPolicy?: { count: number; groupId: string };
  /** Opens one saved group in Advanced, the existing editor of its members and repairs. */
  onOpenAdvanced?: (groupId: string) => void;
}

const NO_SAVED_ROLES: ReadonlySet<string> = new Set();

export const CircuitProfileConfigFields: React.FC<CircuitProfileConfigFieldsProps> = ({
  workspaceId, workspaceRevision, readonly, savedDevices, activeDeviceId, onSelectDevice,
  deviceCandidates, localDraft, updateLocalDraft, instParsed, cumParsed, isPostgres, isMySQL = false, isSQLServer = false, nativeMonthlyAvailable = false,
  storageLocked, effectiveStorageMode, cumulativeCircuitMissing, counterNameMissing = false, circuitId,
  unconfirmedCounterStatuses, onCounterConfirmed, savedCircuits = [], onSelectCircuit, onOpenBatch,
  purpose = 'long', purposeLocked = false, layoutOptions, layoutBlockedReason, showPurpose = true, instantaneousReduction = 'latest', savedRoleKeys = NO_SAVED_ROLES, memberRole = defaultMemberRole,
  unavailableSources = { count: 0 }, requiredPolicy, onOpenAdvanced,
}) => {
  const { t } = useTranslation('workbench-v2');
  const circuitNameRef = React.useRef<HTMLInputElement>(null);
  const circuitNameErrorId = React.useId();
  const circuitNameMissing = cumulativeCircuitMissing || counterNameMissing;
  const [sourceExpandRequest, setSourceExpandRequest] = React.useState(0);
  const instantaneousSourcesMissing = localDraft.instantaneous_enabled &&
    !deviceCandidates.some((c) => localDraft.instantaneous_keys.includes(candidateKey(c)));
  const cumulativeSelectedCount = deviceCandidates.filter((c) => localDraft.cumulative_keys.includes(candidateKey(c))).length;
  const cumulativeSourcesMissing = localDraft.cumulative_enabled && cumulativeSelectedCount === 0;

  return (
    <>
      <div data-testid="circuit-profile-source-section" className="space-y-3 rounded border border-slate-800/80 bg-slate-950/20 p-3">
        <div className="flex items-center justify-between">
          <h4 className="text-xs font-semibold text-slate-200">{t('step4.basic.circuit.step_source')}</h4>
          {savedCircuits.length > 1 && onOpenBatch && (
            <button type="button" data-testid="circuit-profile-open-batch" disabled={readonly} onClick={onOpenBatch}
              className="rounded border border-cyan-800 bg-cyan-950/60 px-2 py-0.5 text-xs text-cyan-200 hover:bg-cyan-900/60">
              {t('step4.basic.circuit.batch_edit_btn') || '批次設定迴路'}
            </button>
          )}
        </div>
        <div className="grid gap-3 sm:grid-cols-3">
          <div className="block text-xs text-slate-300">
            <span>{t('step4.basic.device')}</span>
            <div className="mt-1">
              {/* Choosing a tile only changes the viewed scope, so a pending request or a read-only role never locks it. */}
              <DeviceScopeTiles devices={savedDevices} selectedId={activeDeviceId} onSelect={onSelectDevice} label={t('step4.basic.device')} />
            </div>
          </div>
          <label className="block text-xs text-slate-300">
            <span>{t('step4.basic.circuit.circuit_chooser_label') || '迴路選取'}</span>
            <select data-testid="circuit-profile-circuit-chooser" value={circuitId} disabled={readonly}
              onChange={(e) => onSelectCircuit?.(e.target.value)} className="mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-1.5 text-xs text-slate-100">
              {savedCircuits.map((c) => (
                <option key={c.circuit_id} value={c.circuit_id}>
                  {c.circuit_name || (c.circuit_id ? c.circuit_id : (t('step4.basic.circuit.unassigned') || '未分配迴路'))}
                </option>
              ))}
              {!savedCircuits.some((c) => c.circuit_id === circuitId) && circuitId && (
                <option value={circuitId}>{localDraft.circuit_name || (t('step4.basic.circuit.new_circuit_option') || '新增迴路')}</option>
              )}
              <option value="__new__">+ {t('step4.basic.circuit.new_circuit_option') || '新增迴路'}</option>
            </select>
          </label>
          <label className="block text-xs text-slate-300">
            <span id={`${circuitNameErrorId}-label`}>{t('step4.basic.circuit.circuit_name_label')}</span>
            <input ref={circuitNameRef} type="text" data-testid="circuit-profile-circuit-name" value={localDraft.circuit_name} disabled={readonly}
              aria-labelledby={`${circuitNameErrorId}-label`}
              aria-required={circuitNameMissing} aria-invalid={circuitNameMissing}
              aria-describedby={circuitNameMissing ? circuitNameErrorId : undefined}
              placeholder={t('step4.basic.circuit.circuit_name_placeholder')}
              onChange={(e) => updateLocalDraft((prev) => ({ ...prev, circuit_name: e.target.value, circuit_id: deriveCircuitId(e.target.value, prev.circuit_id) }))}
              className="mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-1.5 text-xs text-slate-100" />
            {circuitNameMissing && (
              <span id={circuitNameErrorId} role="alert" data-testid="circuit-profile-circuit-required-error" className="mt-1 block text-xs text-amber-300">
                {t(counterNameMissing ? 'step4.basic.circuit.counter_circuit_required' : 'step4.basic.circuit.cumulative_circuit_required')}
              </span>
            )}
          </label>
        </div>
        {showPurpose && <CircuitPurposeFields purpose={purpose} locked={purposeLocked} readonly={readonly}
          layoutOptions={layoutOptions} blockedReason={layoutBlockedReason}
          onChange={(next) => updateLocalDraft((prev) => ({ ...prev, purpose: next }))} />}
        <p className="text-[11px] text-slate-400">{t('step4.basic.circuit.unassigned_note')}</p>
        <SourceReview prefix="circuit-profile-sources" count={deviceCandidates.length + unavailableSources.count} unavailable={unavailableSources.count} expandRequest={sourceExpandRequest}
          onRepair={unavailableSources.groupId && onOpenAdvanced ? () => onOpenAdvanced(unavailableSources.groupId as string) : undefined}>
        <ul data-testid="circuit-profile-member-list" className="max-h-72 space-y-1.5 overflow-y-auto text-xs text-slate-300">
          {deviceCandidates.map((candidate) => {
            const key = candidateKey(candidate);
            const instChecked = localDraft.instantaneous_keys.includes(key);
            const cumChecked = localDraft.cumulative_keys.includes(key);
            return (
              <li key={key} className="flex flex-wrap items-center justify-between gap-2 rounded border border-slate-900 bg-slate-950/40 px-2.5 py-1.5">
                <div className="flex items-center gap-2">
                  <label className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      checked={instChecked || cumChecked}
                      disabled={readonly}
                      aria-label={candidate.label}
                      onChange={(event) => {
                        const checked = event.target.checked;
                        updateLocalDraft((prev) => {
                          const nextInst = checked
                            ? (prev.instantaneous_enabled ? Array.from(new Set([...prev.instantaneous_keys, key])) : prev.instantaneous_keys)
                            : prev.instantaneous_keys.filter((k) => k !== key);
                          const nextCum = checked
                            ? (prev.cumulative_enabled ? Array.from(new Set([...prev.cumulative_keys, key])) : prev.cumulative_keys)
                            : prev.cumulative_keys.filter((k) => k !== key);
                          return { ...prev, instantaneous_keys: nextInst, cumulative_keys: nextCum, selected_keys: Array.from(new Set([...nextInst, ...nextCum])) };
                        });
                      }}
                    />
                    <span className="font-medium">{candidate.label}</span>
                  </label>
                  <span className="text-[11px] text-slate-500">({candidate.address} · {candidate.target_type})</span>
                </div>
                <div className="flex items-center gap-3">
                  <label className="flex items-center gap-1 text-[11px] text-slate-400">
                    <input
                      type="checkbox"
                      data-testid={`circuit-profile-member-inst-${candidate.tag_id}`}
                      checked={instChecked}
                      disabled={readonly}
                      aria-label={`${candidate.label} ${t('step4.basic.circuit.instantaneous_label')}`}
                      onChange={(event) => {
                        const checked = event.target.checked;
                        updateLocalDraft((prev) => {
                          const nextInst = checked ? Array.from(new Set([...prev.instantaneous_keys, key])) : prev.instantaneous_keys.filter((k) => k !== key);
                          return { ...prev, instantaneous_keys: nextInst, selected_keys: Array.from(new Set([...nextInst, ...prev.cumulative_keys])) };
                        });
                      }}
                    />
                    <span>{t('step4.basic.circuit.instantaneous_label')}</span>
                  </label>
                  {instChecked && instantaneousReduction === 'mean' && (
                    <CircuitProfileRoleOptions
                      testId={`circuit-profile-role-${candidate.tag_id}`}
                      label={candidate.label}
                      value={memberRole(candidate)}
                      numeric={defaultMemberRole(candidate) === 'gauge'}
                      readonly={readonly}
                      saved={savedRoleKeys.has(key)}
                      onChange={(role) => updateLocalDraft((prev) => ({ ...prev, member_roles: { ...prev.member_roles, [key]: role } }))}
                    />
                  )}
                  <label className="flex items-center gap-1 text-[11px] text-slate-400">
                    <input
                      type="checkbox"
                      data-testid={`circuit-profile-member-cum-${candidate.tag_id}`}
                      checked={cumChecked}
                      disabled={readonly}
                      aria-label={`${candidate.label} ${t('step4.basic.circuit.cumulative_label')}`}
                      onChange={(event) => {
                        const checked = event.target.checked;
                        updateLocalDraft((prev) => {
                          const nextCum = checked ? Array.from(new Set([...prev.cumulative_keys, key])) : prev.cumulative_keys.filter((k) => k !== key);
                          return { ...prev, cumulative_keys: nextCum, selected_keys: Array.from(new Set([...prev.instantaneous_keys, ...nextCum])) };
                        });
                      }}
                    />
                    <span>{t('step4.basic.circuit.cumulative_label')}</span>
                  </label>
                </div>
              </li>
            );
          })}
        </ul>
        </SourceReview>
      </div>
      <div
        data-testid="circuit-profile-profiles-section"
        className="space-y-3 rounded border border-slate-800/80 bg-slate-950/20 p-3"
      >
        <h4 className="text-xs font-semibold text-slate-200">{t('step4.basic.circuit.step_profiles')}</h4>
        <div className="grid gap-3 sm:grid-cols-2">
          <div className="space-y-2 rounded border border-slate-800 p-2.5">
            <label className="flex items-center gap-2 text-xs font-medium text-slate-200">
              <input type="checkbox" data-testid="circuit-profile-toggle-instantaneous" checked={localDraft.instantaneous_enabled}
                disabled={readonly} onChange={(e) => updateLocalDraft((prev) => ({ ...prev, instantaneous_enabled: e.target.checked }))} />
              <span>{t('step4.basic.circuit.instantaneous_label')}</span>
            </label>
            <label className="block text-xs text-slate-300">
              <span>{t('step4.basic.circuit.instantaneous_interval')}</span>
              <input type="text" inputMode="numeric" data-testid="circuit-profile-interval-instantaneous"
                value={localDraft.instantaneous_interval_input} disabled={readonly || !localDraft.instantaneous_enabled}
                aria-invalid={localDraft.instantaneous_enabled && !instParsed.valid}
                aria-describedby={localDraft.instantaneous_enabled && !instParsed.valid ? 'circuit-profile-interval-instantaneous-err' : undefined}
                onChange={(e) => updateLocalDraft((prev) => ({ ...prev, instantaneous_interval_input: e.target.value }))}
                className="mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-1.5 text-xs text-slate-100" />
            </label>
            {instantaneousReduction === 'latest' || (localDraft.instantaneous_enabled && instParsed.valid) ? (
              <p data-testid="circuit-profile-reduction-instantaneous" className="text-[11px] text-slate-400">
                {instantaneousReduction === 'mean'
                  ? t('step4.basic.circuit.reduction_mean', { count: instParsed.intervalValue })
                  : t('step4.basic.circuit.reduction_latest')}
              </p>
            ) : null}
            {localDraft.instantaneous_enabled && !instParsed.valid && (
              <p id="circuit-profile-interval-instantaneous-err" role="alert" data-testid="circuit-profile-interval-instantaneous-error" className="text-xs text-rose-300">
                {t('step4.basic.errors.profile_interval_invalid')}
              </p>
            )}
            {localDraft.instantaneous_enabled && instParsed.valid && instantaneousSourcesMissing && (
              <p role="alert" data-testid="circuit-profile-instantaneous-sources-error" className="text-xs text-amber-300">
                {t('step4.basic.circuit.instantaneous_sources_required')}
              </p>
            )}
          </div>
          <div className="space-y-2 rounded border border-slate-800 p-2.5">
            <label className="flex items-center gap-2 text-xs font-medium text-slate-200">
              <input type="checkbox" data-testid="circuit-profile-toggle-cumulative" checked={localDraft.cumulative_enabled}
                disabled={readonly} onChange={(e) => {
                  const enabled = e.target.checked;
                  updateLocalDraft((prev) => ({ ...prev, cumulative_enabled: enabled }));
                  if (enabled) setSourceExpandRequest((request) => request + 1);
                }} />
              <span>{t('step4.basic.circuit.cumulative_label')}</span>
            </label>
            <label className="block text-xs text-slate-300">
              <span>{t('step4.basic.circuit.cumulative_interval')}</span>
              <input type="text" inputMode="numeric" data-testid="circuit-profile-interval-cumulative"
                value={localDraft.cumulative_interval_input} disabled={readonly || !localDraft.cumulative_enabled}
                aria-invalid={localDraft.cumulative_enabled && !cumParsed.valid}
                aria-describedby={localDraft.cumulative_enabled && !cumParsed.valid ? 'circuit-profile-interval-cumulative-err' : undefined}
                onChange={(e) => updateLocalDraft((prev) => ({ ...prev, cumulative_interval_input: e.target.value }))}
                className="mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-1.5 text-xs text-slate-100" />
            </label>
            <p data-testid="circuit-profile-reduction-cumulative" className="text-[11px] text-slate-400">{t('step4.basic.circuit.reduction_counter')}</p>
            {localDraft.cumulative_enabled && <p data-testid="circuit-profile-cumulative-selected" className="text-xs text-slate-300">
              {t('step4.basic.circuit.cumulative_selected', { count: cumulativeSelectedCount })}
            </p>}
            {localDraft.cumulative_enabled && !cumParsed.valid && (
              <p id="circuit-profile-interval-cumulative-err" role="alert" data-testid="circuit-profile-interval-cumulative-error" className="text-xs text-rose-300">
                {t('step4.basic.errors.profile_interval_invalid')}
              </p>
            )}
            {localDraft.cumulative_enabled && cumParsed.valid && cumulativeSourcesMissing && (
              <p role="alert" data-testid="circuit-profile-cumulative-sources-error" className="text-xs text-amber-300">
                {t('step4.basic.circuit.cumulative_sources_required')}
                <button type="button" data-testid="circuit-profile-cumulative-choose-sources" onClick={() => setSourceExpandRequest((request) => request + 1)}
                  className="ml-2 underline">{t('step4.basic.circuit.cumulative_choose_sources')}</button>
              </p>
            )}
          </div>
        </div>
        {localDraft.instantaneous_enabled && localDraft.cumulative_enabled && (
          <div data-testid="circuit-profile-interval-explainer" className="space-y-0.5 text-[11px] text-slate-400">
            <p>{t('step4.basic.circuit.interval_explainer')}</p>
            <p>{t('step4.basic.circuit.interval_not_report')}</p>
          </div>
        )}
        {requiredPolicy && purpose === 'wide' && (
          <CircuitRequiredPolicy count={requiredPolicy.count} onOpen={onOpenAdvanced ? () => onOpenAdvanced(requiredPolicy.groupId) : undefined} />
        )}
        {purpose !== 'wide' && (<>
        <fieldset data-testid="circuit-profile-storage-mode" className="space-y-1.5 text-xs text-slate-300">
          <legend className="font-medium text-slate-200">{t('step4.basic.circuit.storage_mode_label')}</legend>
          <div className="flex flex-wrap gap-4">
            <label className="flex items-center gap-2">
              <input type="radio" name="circuit-storage-mode" data-testid="circuit-profile-storage-plain"
                checked={effectiveStorageMode === 'plain'} disabled={readonly || storageLocked}
                onChange={() => updateLocalDraft((prev) => ({ ...prev, storage_mode: 'plain' }))} />
              <span>{t('step4.basic.circuit.storage_plain')}</span>
            </label>
            {(isMySQL || isSQLServer) ? (
              <label className="flex items-center gap-2">
                <input type="radio" name="circuit-storage-mode" data-testid="circuit-profile-storage-native-monthly"
                  checked={effectiveStorageMode === 'native_monthly'} disabled={readonly || storageLocked || !nativeMonthlyAvailable}
                  onChange={() => updateLocalDraft((prev) => ({ ...prev, storage_mode: 'native_monthly' }))} />
                <span>
                  {isSQLServer
                    ? t('step4.basic.circuit.storage_native_monthly_sqlserver')
                    : t('step4.basic.circuit.storage_native_monthly')}
                </span>
              </label>
            ) : (
              <label className="flex items-center gap-2">
                <input type="radio" name="circuit-storage-mode" data-testid="circuit-profile-storage-timescale"
                  checked={effectiveStorageMode === 'timescale'} disabled={readonly || !isPostgres || storageLocked}
                  onChange={() => updateLocalDraft((prev) => ({ ...prev, storage_mode: 'timescale' }))} />
                <span>{t('step4.basic.circuit.storage_timescale')}</span>
              </label>
            )}
          </div>
          {storageLocked && <p data-testid="circuit-profile-storage-locked" className="text-[11px] text-slate-400">{t('step4.basic.circuit.storage_locked')}</p>}
        </fieldset>
        <label className="flex items-center gap-2 text-xs text-slate-200">
          <input type="checkbox" data-testid="circuit-profile-recreate-missing" checked={localDraft.recreate_missing}
            disabled={readonly} onChange={(e) => updateLocalDraft((prev) => ({ ...prev, recreate_missing: e.target.checked }))} />
          <span>{t('step4.basic.circuit.recreate_missing_label')}</span>
        </label>
        </>)}
        {unconfirmedCounterStatuses.length > 0 && (
          <CircuitCounterConfirmation
            workspaceId={workspaceId ?? ''}
            workspaceRevision={workspaceRevision}
            circuitId={circuitId}
            unconfirmedStatuses={unconfirmedCounterStatuses}
            readonly={readonly}
            onConfirmed={onCounterConfirmed}
            namePrerequisite={counterNameMissing ? { descriptionId: circuitNameErrorId, onFocus: () => circuitNameRef.current?.focus() } : undefined}
          />
        )}
      </div>
    </>
  );
};
