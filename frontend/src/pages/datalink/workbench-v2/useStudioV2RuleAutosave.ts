import { canonicalRulePatch, ruleSaveFields, type RuleFieldGenerations } from '../../../features/datalink/workbench-v2/state/studioV2RuleAcknowledgement';
import type { RuleQueueObservation } from './sourceQueueObservation';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useCreateStudioV2RuleMutation, useDeleteStudioV2RuleMutation, useUpdateStudioV2RuleMutation } from '../../../hooks/datalink/useStudioV2Rules';
import { workbenchV2Reducer, type WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { canHydrateStudioV2Rule, isStudioV2RuleValid, toStudioV2RuleCreateRequest, toStudioV2RuleUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
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

export function useStudioV2RuleAutosave(
  actions: { dispatch: (action: WorkbenchV2Action) => void },
  stateRef: React.MutableRefObject<WorkbenchV2State>,
  canMutateRule: CanMutateRule = allowAllRules,
  observation?: RuleQueueObservation,
  workspaceId?: string,
) {
  const owner = React.useRef({ workspaceId });
  if (owner.current.workspaceId !== workspaceId) owner.current = { workspaceId };
  const mounted = React.useRef(true);
  React.useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const fieldGenerationsRef = React.useRef<Record<string, RuleFieldGenerations>>({});
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

  const applyRulePatch = React.useCallback((ruleId: string, patch: Partial<Rule>, type: 'updateRule' | 'acknowledgeRuleSave' = 'updateRule') => {
    stateRef.current = workbenchV2Reducer(stateRef.current, {
      type,
      ruleId,
      patch,
    });
    actions.dispatch({
      type,
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
    const requestOwner = owner.current;
    const requestGeneration = generationsRef.current.get(ruleId) ?? 0;
    applyRulePatch(ruleId, {
      save_state: 'saving',
      save_error: null,
    });
    const requestScopeKey = persistedMCScopeKey(currentRule.mc_persisted_points);

    try {
      observationRef.current?.submitted?.(currentRule, requestGeneration);
      const savedRule = currentRule.persisted
        ? await updateRuleMutation.mutateAsync({
            ruleId,
            request: toStudioV2RuleUpdateRequest(currentRule),
          })
        : await createRuleMutation.mutateAsync(toStudioV2RuleCreateRequest(currentRule));

      if (!mounted.current || owner.current !== requestOwner) { meta.pending = false; return; }
      const validIdentity = canHydrateStudioV2Rule(savedRule) && savedRule.id === ruleId && savedRule.device_id === currentRule.device_id &&
        (!(requestOwner.workspaceId ?? currentRule.workspace_id) || savedRule.workspace_id === (requestOwner.workspaceId ?? currentRule.workspace_id));
      const acknowledged = observationRef.current?.acknowledged(currentRule, savedRule);
      if (!validIdentity || acknowledged === false) {
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
          ...canonicalRulePatch(currentRule, latestRule, savedRule, requestGeneration, fieldGenerationsRef.current[ruleId] ?? {}),
          workspace_id: savedRule.workspace_id,
          revision_id: savedRule.revision_id,
          mc_persisted_points: evidenceMatchesAck ? acknowledgedEvidence : markPersistedMCChecking(acknowledgedEvidence, savedRule.revision_id),
          persisted: true,
          save_state: meta.pending ? 'saving' : 'saved',
          save_error: null,
        }, 'acknowledgeRuleSave');
        observationRef.current?.saved(currentRule, savedRule);
      }
    } catch (error) {
      if (!mounted.current || owner.current !== requestOwner) { meta.pending = false; return; }
      observationRef.current?.failed(currentRule, error);
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
    const fields = ruleSaveFields(action);
    const ruleId = action.type === 'addRule' ? action.rule.id : 'ruleId' in action ? action.ruleId : undefined;
    if (!ruleId || (fields.length === 0 && action.type !== 'updateRule')) return;
    const generation = (generationsRef.current.get(ruleId) ?? 0) + 1;
    generationsRef.current.set(ruleId, generation);
    if (fields.length === 0) return;
    const edits = fieldGenerationsRef.current[ruleId] ??= {};
    fields.forEach(field => { edits[field] = generation; });
    const rule = stateRef.current.rules.find(row => row.id === ruleId);
    if (rule) observationRef.current?.queued?.(rule, generation);
    queueRuleSave(ruleId);
  }, [queueRuleSave, stateRef]);

  return {
    afterRuleAction,
    interceptRuleAction,
  };
}
