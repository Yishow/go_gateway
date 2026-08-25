import * as React from 'react';
import { useTranslation } from 'react-i18next';
import {
  useCreateStudioV2MappingMutation,
  useDeleteStudioV2MappingMutation,
  useStudioV2MappingsQuery,
  useUpdateStudioV2MappingMutation,
} from '../../../hooks/datalink/useStudioV2Mappings';
import {
  hydrateStudioV2Mapping,
  isStudioV2MappingValid,
  toStudioV2MappingCreateRequest,
  toStudioV2MappingUpdateRequest,
} from '../../../features/datalink/workbench-v2/state/studioV2MappingAutosave';
import { workbenchV2Reducer, type WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Mapping, MappingValue, Point, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import { getSafeErrorStateMessage } from '../../../utils/typedErrors';

type SaveMeta = {
  inFlight: boolean;
  pending: boolean;
};

function mappingRowKey(ruleId: string, address: string): string {
  return `${ruleId}::${address}`;
}

function saveMetaFor(store: Record<string, SaveMeta>, pointId: string): SaveMeta {
  if (!store[pointId]) {
    store[pointId] = { inFlight: false, pending: false };
  }

  return store[pointId];
}

function errorMessageOf(error: unknown, fallback: string): string {
  return getSafeErrorStateMessage(error, fallback);
}

function canAutosaveMappingRule(state: WorkbenchV2State, point: Point): boolean {
  const ownerRule = state.rules.find((rule) => rule.id === point.rule_id);
  return Boolean(ownerRule?.persisted && ownerRule.save_state === 'saved');
}

function shouldAutosaveMappingPatch(patch: Partial<Mapping>): boolean {
  return (
    patch.tag_key !== undefined ||
    patch.display_name !== undefined ||
    patch.unit !== undefined ||
    patch.target_type !== undefined ||
    patch.scale !== undefined ||
    patch.offset !== undefined ||
    patch.enabled !== undefined
  );
}

function sameMappingValue(left?: MappingValue, right?: MappingValue): boolean {
  return (
    left?.tag_key === right?.tag_key &&
    left?.display_name === right?.display_name &&
    left?.unit === right?.unit &&
    left?.target_type === right?.target_type &&
    left?.scale === right?.scale &&
    left?.offset === right?.offset &&
    left?.enabled === right?.enabled
  );
}

function sameHydratedMapping(current: Mapping, next: Mapping): boolean {
  return (
    current.tag_key === next.tag_key &&
    current.display_name === next.display_name &&
    current.unit === next.unit &&
    current.target_type === next.target_type &&
    current.scale === next.scale &&
    current.offset === next.offset &&
    current.enabled === next.enabled &&
    current.mapping_id === next.mapping_id &&
    current.workspace_id === next.workspace_id &&
    current.rule_id === next.rule_id &&
    current.device_id === next.device_id &&
    current.address === next.address &&
    current.tag_id === next.tag_id &&
    current.persisted_point_id === next.persisted_point_id &&
    current.persisted === next.persisted &&
    current.save_state === next.save_state &&
    current.save_error === next.save_error &&
    sameMappingValue(current.local_value, next.local_value) &&
    sameMappingValue(current.persisted_value, next.persisted_value)
  );
}

export function useStudioV2MappingAutosave(
  actions: { dispatch: (action: WorkbenchV2Action) => void },
  stateRef: React.MutableRefObject<WorkbenchV2State>,
  enabled: boolean,
) {
  const { t } = useTranslation('workbench-v2');
  const mappingsQuery = useStudioV2MappingsQuery(enabled);
  const createMappingMutation = useCreateStudioV2MappingMutation();
  const updateMappingMutation = useUpdateStudioV2MappingMutation();
  const deleteMappingMutation = useDeleteStudioV2MappingMutation();
  const saveMetaRef = React.useRef<Record<string, SaveMeta>>({});

  const applyMappingPatch = React.useCallback((pointId: string, patch: Partial<Mapping>) => {
    stateRef.current = workbenchV2Reducer(stateRef.current, {
      type: 'updateMapping',
      pointId,
      patch,
    });
    actions.dispatch({
      type: 'updateMapping',
      pointId,
      patch,
    });
  }, [actions, stateRef]);

  const reconcilePersistedMappings = React.useCallback((points: Point[]) => {
    if (!mappingsQuery.isSuccess || points.length === 0) {
      return;
    }

    const records = new Map(
      mappingsQuery.data.map((record) => [mappingRowKey(record.rule_id, record.address), record]),
    );

    points.forEach((point) => {
      const record = records.get(mappingRowKey(point.rule_id, point.address));
      const current = stateRef.current.mappings[point.id];
      if (!record || !current) {
        return;
      }

      const hydrated = hydrateStudioV2Mapping(point, record, current);
      if (sameHydratedMapping(current, hydrated)) {
        return;
      }
      applyMappingPatch(point.id, hydrated);
    });
  }, [applyMappingPatch, mappingsQuery.data, mappingsQuery.isSuccess, stateRef]);

  React.useEffect(() => {
    reconcilePersistedMappings(stateRef.current.points);
  }, [reconcilePersistedMappings, stateRef]);

  const flushMappingSave = React.useCallback(async (pointId: string) => {
    const meta = saveMetaFor(saveMetaRef.current, pointId);
    const point = stateRef.current.points.find((item) => item.id === pointId);
    const currentMapping = stateRef.current.mappings[pointId];
    if (!point || !currentMapping) {
      delete saveMetaRef.current[pointId];
      return;
    }

    if (!isStudioV2MappingValid(currentMapping)) {
      applyMappingPatch(pointId, {
        save_state: 'draft-invalid',
        save_error: null,
      });
      meta.inFlight = false;
      meta.pending = false;
      return;
    }

    if (!canAutosaveMappingRule(stateRef.current, point)) {
      applyMappingPatch(pointId, {
        save_state: 'draft-invalid',
        save_error: '請先完成來源規則保存',
      });
      meta.inFlight = false;
      meta.pending = false;
      return;
    }

    meta.inFlight = true;
    applyMappingPatch(pointId, {
      save_state: 'saving',
      save_error: null,
    });

    try {
      const savedMapping = currentMapping.mapping_id
        ? await updateMappingMutation.mutateAsync({
          mappingId: currentMapping.mapping_id,
          request: toStudioV2MappingUpdateRequest(point, currentMapping),
        })
        : await createMappingMutation.mutateAsync(toStudioV2MappingCreateRequest(point, currentMapping));

      applyMappingPatch(pointId, hydrateStudioV2Mapping(point, savedMapping));
    } catch (error) {
      applyMappingPatch(pointId, {
        save_state: 'save-error',
        save_error: errorMessageOf(error, t('errors.autosave_failed')),
      });
    } finally {
      meta.inFlight = false;
      if (meta.pending) {
        meta.pending = false;
        void flushMappingSave(pointId);
      }
    }
  }, [applyMappingPatch, createMappingMutation, stateRef, t, updateMappingMutation]);

  const queueMappingSave = React.useCallback((pointId: string) => {
    const meta = saveMetaFor(saveMetaRef.current, pointId);
    if (meta.inFlight) {
      meta.pending = true;
      return;
    }

    void flushMappingSave(pointId);
  }, [flushMappingSave]);

  const afterMappingAction = React.useCallback((previousState: WorkbenchV2State, action: WorkbenchV2Action, nextState: WorkbenchV2State) => {
    switch (action.type) {
      case 'initMappingsForPoints': {
        const nextPointIDs = new Set(nextState.points.map((point) => point.id));
        Object.entries(previousState.mappings).forEach(([pointId, mapping]) => {
          if (!nextPointIDs.has(pointId) && mapping.mapping_id) {
            void deleteMappingMutation.mutateAsync(mapping.mapping_id);
            delete saveMetaRef.current[pointId];
          }
        });
        reconcilePersistedMappings(nextState.points);
        break;
      }
      case 'updateMapping':
        if (shouldAutosaveMappingPatch(action.patch)) {
          queueMappingSave(action.pointId);
        }
        break;
      case 'toggleMappingEnabled':
        queueMappingSave(action.pointId);
        break;
      case 'setAllMappingsEnabled':
        Object.keys(nextState.mappings).forEach(queueMappingSave);
        break;
      case 'bulkApplyTransform':
        Object.keys(nextState.mappings).forEach(queueMappingSave);
        break;
      default:
        break;
    }
  }, [deleteMappingMutation, queueMappingSave, reconcilePersistedMappings]);

  return {
    afterMappingAction,
    mappingsQuery,
  };
}
