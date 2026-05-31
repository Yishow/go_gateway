import * as React from 'react';
import { useCreateStudioV2RuleMutation, useDeleteStudioV2RuleMutation, useUpdateStudioV2RuleMutation } from '../../../hooks/datalink/useStudioV2Rules';
import { workbenchV2Reducer, type WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { isStudioV2RuleValid, toStudioV2RuleCreateRequest, toStudioV2RuleUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
import type { Rule, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';

type SaveMeta = {
  inFlight: boolean;
  pending: boolean;
};

function saveMetaFor(store: Record<string, SaveMeta>, ruleId: string): SaveMeta {
  if (!store[ruleId]) {
    store[ruleId] = { inFlight: false, pending: false };
  }

  return store[ruleId];
}

function errorMessageOf(error: unknown): string {
  return error instanceof Error ? error.message : '儲存失敗';
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
    patch.data_format !== undefined
  );
}

export function useStudioV2RuleAutosave(
  actions: { dispatch: (action: WorkbenchV2Action) => void },
  stateRef: React.MutableRefObject<WorkbenchV2State>,
) {
  const createRuleMutation = useCreateStudioV2RuleMutation();
  const updateRuleMutation = useUpdateStudioV2RuleMutation();
  const deleteRuleMutation = useDeleteStudioV2RuleMutation();
  const saveMetaRef = React.useRef<Record<string, SaveMeta>>({});

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

    try {
      const savedRule = currentRule.persisted
        ? await updateRuleMutation.mutateAsync({
            ruleId,
            request: toStudioV2RuleUpdateRequest(currentRule),
          })
        : await createRuleMutation.mutateAsync(toStudioV2RuleCreateRequest(currentRule));

      const latestRule = stateRef.current.rules.find((rule) => rule.id === ruleId);
      if (latestRule) {
        applyRulePatch(ruleId, {
          device_id: latestRule.device_id,
          workspace_id: savedRule.workspace_id,
          revision_id: savedRule.revision_id,
          persisted: true,
          save_state: 'saved',
          save_error: null,
        });
      }
    } catch (error) {
      applyRulePatch(ruleId, {
        save_state: 'save-error',
        save_error: errorMessageOf(error),
      });
    } finally {
      meta.inFlight = false;
      if (meta.pending) {
        meta.pending = false;
        void flushRuleSave(ruleId);
      }
    }
  }, [applyRulePatch, createRuleMutation, stateRef, updateRuleMutation]);

  const queueRuleSave = React.useCallback((ruleId: string) => {
    const meta = saveMetaFor(saveMetaRef.current, ruleId);
    if (meta.inFlight) {
      meta.pending = true;
      return;
    }

    void flushRuleSave(ruleId);
  }, [flushRuleSave]);

  const interceptRuleAction = React.useCallback((previousState: { rules: Rule[] }, action: WorkbenchV2Action): boolean => {
    if (action.type === 'removeRule') {
      const targetRule = previousState.rules.find((rule) => rule.id === action.ruleId);
      if (!targetRule) {
        return true;
      }

      if (!targetRule.persisted) {
        return false;
      }

      applyRulePatch(action.ruleId, {
        save_state: 'saving',
        save_error: null,
      });
      void deleteRuleMutation.mutateAsync(action.ruleId)
        .then(() => {
          stateRef.current = workbenchV2Reducer(stateRef.current, action);
          actions.dispatch(action);
          delete saveMetaRef.current[action.ruleId];
        })
        .catch((error) => {
          applyRulePatch(action.ruleId, {
            save_state: 'save-error',
            save_error: errorMessageOf(error),
          });
        });
      return true;
    }
    return false;
  }, [actions, applyRulePatch, deleteRuleMutation, stateRef]);

  const afterRuleAction = React.useCallback((action: WorkbenchV2Action) => {
    switch (action.type) {
      case 'addRule':
        queueRuleSave(action.rule.id);
        break;
      case 'updateRule':
        if (shouldAutosaveRulePatch(action.patch)) {
          queueRuleSave(action.ruleId);
        }
        break;
      case 'toggleRuleEnabled':
      case 'updateRuleSkipped':
      case 'toggleRuleSkippedAddress':
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
