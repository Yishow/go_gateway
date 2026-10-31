import type { RuleQueueObservation } from './sourceQueueObservation';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useCreateStudioV2RuleMutation, useDeleteStudioV2RuleMutation, useUpdateStudioV2RuleMutation } from '../../../hooks/datalink/useStudioV2Rules';
import { workbenchV2Reducer, type WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { isStudioV2RuleValid, toStudioV2RuleCreateRequest, toStudioV2RuleUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { markPersistedMCChecking, persistedMCScopeKey } from '../../../features/datalink/workbench-v2/state/mcPersistedPoints';
import type { Rule, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import { getSafeErrorStateMessage } from '../../../utils/typedErrors';

type SaveMeta = {
  inFlight: boolean;
  pending: boolean;
};

type CanMutateRule = (ruleId: string) => boolean;

const allowAllRules: CanMutateRule = () => true;

function saveMetaFor(store: Record<string, SaveMeta>, ruleId: string): SaveMeta {
  if (!store[ruleId]) {
    store[ruleId] = { inFlight: false, pending: false };
  }

  return store[ruleId];
}

function errorMessageOf(error: unknown, fallback: string): string {
  return getSafeErrorStateMessage(error, fallback);
}

function shouldAutosaveRulePatch(patch: Partial<Rule>): boolean {
  return (
    patch.device_id !== undefined ||
    patch.start_address !== undefined ||
    patch.count !== undefined ||
    patch.data_type !== undefined ||
    patch.naming_prefix !== undefined ||
    patch.enabled !== undefined ||
    patch.skipped_addresses !== undefined ||
    patch.scale_multiplier !== undefined ||
    patch.scale_offset !== undefined ||
    patch.data_format !== undefined ||
    patch.share_enabled !== undefined ||
    patch.share_start_register !== undefined ||
    patch.share_stride !== undefined
  );
}

export function useStudioV2RuleAutosave(
  actions: { dispatch: (action: WorkbenchV2Action) => void },
  stateRef: React.MutableRefObject<WorkbenchV2State>,
  canMutateRule: CanMutateRule = allowAllRules,
  observation?: RuleQueueObservation,
) {
  const { t } = useTranslation('workbench-v2');
  const observationRef = React.useRef(observation);
  observationRef.current = observation;
  const createRuleMutation = useCreateStudioV2RuleMutation();
  const updateRuleMutation = useUpdateStudioV2RuleMutation();
  const deleteRuleMutation = useDeleteStudioV2RuleMutation();
  const saveMetaRef = React.useRef<Record<string, SaveMeta>>({});
  const canMutateRuleRef = React.useRef(canMutateRule);
  const flushRuleSaveRef = React.useRef<(ruleId: string) => void>(() => undefined);
  const generationsRef = React.useRef(new Map<string, number>());
  canMutateRuleRef.current = canMutateRule;

  const applyRulePatch = React.useCallback((ruleId: string, patch: Partial<Rule>) => {
    stateRef.current = workbenchV2Reducer(stateRef.current, {
      type: 'updateRule',
      ruleId,
      patch,
    });
    actions.dispatch({
      type: 'updateRule',
      ruleId,
      patch,
    });
  }, [actions, stateRef]);

  const flushRuleSave = React.useCallback(async (ruleId: string) => {
    const meta = saveMetaFor(saveMetaRef.current, ruleId);
    const currentRule = stateRef.current.rules.find((rule) => rule.id === ruleId);
    if (!currentRule) {
      delete saveMetaRef.current[ruleId];
      return;
    }

    if (!canMutateRuleRef.current(ruleId)) {
      meta.pending = true;
      applyRulePatch(ruleId, {
        save_state: 'draft-invalid',
        save_error: null,
      });
      meta.inFlight = false;
      return;
    }

    if (!isStudioV2RuleValid(currentRule)) {
      applyRulePatch(ruleId, {
        save_state: 'draft-invalid',
        save_error: null,
      });
      meta.inFlight = false;
      meta.pending = false;
      return;
    }

    meta.inFlight = true;
    applyRulePatch(ruleId, {
      save_state: 'saving',
      save_error: null,
    });
    const requestScopeKey = persistedMCScopeKey(currentRule.mc_persisted_points);

    try {
      const savedRule = currentRule.persisted
        ? await updateRuleMutation.mutateAsync({
            ruleId,
            request: toStudioV2RuleUpdateRequest(currentRule),
          })
        : await createRuleMutation.mutateAsync(toStudioV2RuleCreateRequest(currentRule));

      if (observationRef.current?.acknowledged(currentRule, savedRule) === false) {
        meta.pending = false;
        const latest = stateRef.current.rules.find(rule => rule.id === ruleId);
        if (latest && persistedMCScopeKey(latest.mc_persisted_points) === requestScopeKey) {
          applyRulePatch(ruleId, { save_state: 'save-error', save_error: t('errors.autosave_failed') });
        }
        return;
      }
      const latestRule = stateRef.current.rules.find((rule) => rule.id === ruleId);
      const currentScopeKey = persistedMCScopeKey(latestRule?.mc_persisted_points);
      const acknowledgedEvidence = latestRule?.mc_persisted_points;
      const requestScope = currentRule.mc_persisted_points?.scope;
      const acknowledgedScope = latestRule?.mc_persisted_points?.scope;
      const ownRevisionTransition = Boolean(
        acknowledgedScope &&
        acknowledgedScope.revision_id === savedRule.revision_id &&
        (!requestScope || requestScope.workspace_id === acknowledgedScope.workspace_id),
      );
      const scopeStillCurrent = requestScopeKey === currentScopeKey || ownRevisionTransition;
      const metadataOnlyAck = !acknowledgedEvidence && requestScopeKey === currentScopeKey;
      if (!latestRule || !scopeStillCurrent || (!canMutateRuleRef.current(ruleId) && !ownRevisionTransition && !metadataOnlyAck)) {
        if (latestRule && !meta.pending) {
          applyRulePatch(ruleId, { save_state: 'draft-invalid', save_error: null });
        }
        return;
      }
      if (latestRule) {
        const evidenceMatchesAck = ownRevisionTransition;
        applyRulePatch(ruleId, {
          device_id: latestRule.device_id,
          workspace_id: savedRule.workspace_id,
          revision_id: savedRule.revision_id,
          mc_persisted_points: evidenceMatchesAck ? acknowledgedEvidence : markPersistedMCChecking(acknowledgedEvidence, savedRule.revision_id),
          persisted: true,
          save_state: 'saved',
          save_error: null,
        });
        observationRef.current?.saved(currentRule, savedRule);
      }
    } catch (error) {
      observationRef.current?.failed(currentRule);
      const latestRule = stateRef.current.rules.find((rule) => rule.id === ruleId);
      if (!latestRule || persistedMCScopeKey(latestRule.mc_persisted_points) !== requestScopeKey || !canMutateRuleRef.current(ruleId)) {
        if (latestRule && !meta.pending) applyRulePatch(ruleId, { save_state: 'draft-invalid', save_error: null });
        return;
      }
      applyRulePatch(ruleId, {
        save_state: 'save-error',
        save_error: errorMessageOf(error, t('errors.autosave_failed')),
      });
    } finally {
      meta.inFlight = false;
      if (meta.pending && canMutateRuleRef.current(ruleId)) {
        meta.pending = false;
        void flushRuleSave(ruleId);
      }
    }
  }, [applyRulePatch, createRuleMutation, stateRef, t, updateRuleMutation]);

  flushRuleSaveRef.current = flushRuleSave;
  React.useEffect(() => {
    for (const [ruleId, meta] of Object.entries(saveMetaRef.current)) {
      if (meta.pending && !meta.inFlight && canMutateRuleRef.current(ruleId)) {
        meta.pending = false;
        void flushRuleSaveRef.current(ruleId);
      }
    }
  }, [canMutateRule]);

  const queueRuleSave = React.useCallback((ruleId: string) => {
    const meta = saveMetaFor(saveMetaRef.current, ruleId);
    if (!canMutateRuleRef.current(ruleId)) {
      meta.pending = true;
      applyRulePatch(ruleId, {
        save_state: 'draft-invalid',
        save_error: null,
      });
      return;
    }
    if (meta.inFlight) {
      meta.pending = true;
      return;
    }

    void flushRuleSave(ruleId);
  }, [applyRulePatch, flushRuleSave]);

  const interceptRuleAction = React.useCallback((previousState: { rules: Rule[] }, action: WorkbenchV2Action): boolean => {
    if (action.type === 'removeRule') {
      const targetRule = previousState.rules.find((rule) => rule.id === action.ruleId);
      if (!targetRule) {
        return true;
      }

      if (!targetRule.persisted) {
        return false;
      }

      if (!canMutateRuleRef.current(action.ruleId)) {
        applyRulePatch(action.ruleId, {
          save_state: 'draft-invalid',
          save_error: null,
        });
        return true;
      }

      const requestGeneration = (generationsRef.current.get(action.ruleId) ?? 0) + 1;
      generationsRef.current.set(action.ruleId, requestGeneration);
      const requestScopeKey = persistedMCScopeKey(targetRule.mc_persisted_points);
      applyRulePatch(action.ruleId, {
        save_state: 'saving',
        save_error: null,
      });
      void deleteRuleMutation.mutateAsync(action.ruleId)
        .then(() => {
          const latestRule = stateRef.current.rules.find((rule) => rule.id === action.ruleId);
          const currentScopeKey = persistedMCScopeKey(latestRule?.mc_persisted_points);
          if (
            !latestRule ||
            requestGeneration !== generationsRef.current.get(action.ruleId) ||
            requestScopeKey !== currentScopeKey ||
            !canMutateRuleRef.current(action.ruleId)
          ) {
            if (latestRule) applyRulePatch(action.ruleId, { save_state: 'draft-invalid', save_error: null });
            return;
          }
          stateRef.current = workbenchV2Reducer(stateRef.current, action);
          actions.dispatch(action);
          delete saveMetaRef.current[action.ruleId];
        })
        .catch((error) => {
          const latestRule = stateRef.current.rules.find((rule) => rule.id === action.ruleId);
          if (
            !latestRule ||
            requestGeneration !== generationsRef.current.get(action.ruleId) ||
            requestScopeKey !== persistedMCScopeKey(latestRule.mc_persisted_points) ||
            !canMutateRuleRef.current(action.ruleId)
          ) {
            if (latestRule) applyRulePatch(action.ruleId, { save_state: 'draft-invalid', save_error: null });
            return;
          }
          applyRulePatch(action.ruleId, {
            save_state: 'save-error',
            save_error: errorMessageOf(error, t('errors.autosave_failed')),
          });
        });
      return true;
    }
    return false;
  }, [actions, applyRulePatch, deleteRuleMutation, stateRef, t]);

  const afterRuleAction = React.useCallback((action: WorkbenchV2Action) => {
    switch (action.type) {
      case 'addRule':
        generationsRef.current.set(action.rule.id, (generationsRef.current.get(action.rule.id) ?? 0) + 1);
        queueRuleSave(action.rule.id);
        break;
      case 'updateRule':
        generationsRef.current.set(action.ruleId, (generationsRef.current.get(action.ruleId) ?? 0) + 1);
        if (shouldAutosaveRulePatch(action.patch)) {
          queueRuleSave(action.ruleId);
        }
        break;
      case 'toggleRuleEnabled':
      case 'updateRuleSkipped':
      case 'toggleRuleSkippedAddress':
      case 'toggleRuleShareEnabled':
      case 'updateRuleShareStart':
      case 'updateRuleShareStride':
        generationsRef.current.set(action.ruleId, (generationsRef.current.get(action.ruleId) ?? 0) + 1);
        queueRuleSave(action.ruleId);
        break;
      default:
        break;
    }
  }, [queueRuleSave]);

  return {
    afterRuleAction,
    interceptRuleAction,
  };
}
