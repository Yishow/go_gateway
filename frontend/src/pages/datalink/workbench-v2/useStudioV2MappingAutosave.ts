import * as React from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from '../../../hooks/datalink/keys';
import type { StudioV2WorkspaceMappingRecord } from '../../../types/datalink';
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
import { mappingSaveError } from '../../../utils/mappingSaveErrors';

type RuleQueue = {
  running: boolean;
  paused: boolean;
  pending: Map<string, number>;
  deletions: Map<string, { point: Point; mapping: Mapping; restored?: boolean }>;
  failure?: Partial<Mapping>;
};

function mappingRowKey(ruleId: string, address: string): string {
  return `${ruleId}::${address}`;
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
  const queryClient = useQueryClient();
  const mappingsQuery = useStudioV2MappingsQuery(enabled);
  const createMappingMutation = useCreateStudioV2MappingMutation();
  const updateMappingMutation = useUpdateStudioV2MappingMutation();
  const deleteMappingMutation = useDeleteStudioV2MappingMutation();
  const queues = React.useRef(new Map<string, RuleQueue>());
  const mounted = React.useRef(true);
  const draftVersions = React.useRef(new Map<string, number>());
  const confirmedRecords = React.useRef(new Map<string, StudioV2WorkspaceMappingRecord>());
  const deletedRecords = React.useRef(new Set<string>());
  const refreshing = React.useRef(false);
  const refreshAgain = React.useRef(false);
  const generations = React.useRef(new Map<string, number>());
  React.useEffect(() => {
    mounted.current = true;
    const activeQueues = queues.current;
    return () => { mounted.current = false; activeQueues.clear(); };
  }, []);
  const ruleQueue = React.useCallback((ruleId: string) => {
    let queue = queues.current.get(ruleId);
    if (!queue) {
      queue = { running: false, paused: false, pending: new Map(), deletions: new Map() };
      queues.current.set(ruleId, queue);
    }
    return queue;
  }, []);

  const applyMappingPatch = React.useCallback((pointId: string, patch: Partial<Mapping>) => {
    if (!mounted.current || !stateRef.current.mappings[pointId]) return;
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

  const markCleanupIncomplete = React.useCallback((status: unknown) => {
    if (!mounted.current || status !== 'failed') return;
    const payload = { mapping_cleanup_incomplete: true };
    stateRef.current = workbenchV2Reducer(stateRef.current, { type: 'SET_STATE', payload });
    actions.dispatch({ type: 'SET_STATE', payload });
  }, [actions, stateRef]);

  const reconcilePersistedMappings = React.useCallback((points: Point[]) => {
    if (refreshing.current || !mappingsQuery.isSuccess || points.length === 0) {
      return;
    }

    const records = new Map(
      mappingsQuery.data.map((record) => [mappingRowKey(record.rule_id, record.address), record]),
    );

    points.forEach((point) => {
      const record = records.get(mappingRowKey(point.rule_id, point.address));
      const current = stateRef.current.mappings[point.id];
      if (!record || !current || queues.current.get(point.rule_id)?.deletions.has(record.id) || queues.current.get(point.rule_id)?.running || current.save_state === 'saving' || current.save_state === 'save-error') {
        return;
      }

      const hydrated = hydrateStudioV2Mapping(point, record, current);
      if (sameHydratedMapping(current, hydrated)) {
        return;
      }
      confirmedRecords.current.set(record.id, record);
      applyMappingPatch(point.id, hydrated);
    });
  }, [applyMappingPatch, mappingsQuery.data, mappingsQuery.isSuccess, stateRef]);

  React.useEffect(() => {
    reconcilePersistedMappings(stateRef.current.points);
  }, [reconcilePersistedMappings, stateRef]);

  const refreshAuthority = React.useCallback(function refreshAuthority() {
    if (!mounted.current) return;
    if (refreshing.current) { refreshAgain.current = true; return; }
    refreshing.current = true;
    const versions = new Map(draftVersions.current);
    const identities = new Map(generations.current);
    const acknowledgementsAtStart = new Map(confirmedRecords.current);
    void mappingsQuery.refetch({ cancelRefetch: false }).then((fresh) => {
      if (!mounted.current || fresh.isError) return;
      for (const point of stateRef.current.points) {
        const current = stateRef.current.mappings[point.id];
        const record = fresh.data?.find((item) => item.rule_id === point.rule_id && item.address === point.address);
        if (record && current?.save_state === 'saved' &&
          (versions.get(point.id) ?? 0) === (draftVersions.current.get(point.id) ?? 0) &&
          (identities.get(point.id) ?? 0) === (generations.current.get(point.id) ?? 0) &&
          !ruleQueue(point.rule_id).running && !ruleQueue(point.rule_id).paused) {
          confirmedRecords.current.set(record.id, record);
          applyMappingPatch(point.id, hydrateStudioV2Mapping(point, record, current));
        }
      }
      // Retain acknowledgements made after this GET began; later GETs remain
      // authoritative, including legitimate tag-only changes with equal stamps.
      const records = (fresh.data ?? []).filter((record) => !deletedRecords.current.has(record.id)).map((record) => confirmedRecords.current.get(record.id) ?? record);
      const presentIDs = new Set(records.map((record) => record.id));
      for (const [id, record] of confirmedRecords.current) {
        if (!presentIDs.has(id) && !deletedRecords.current.has(id) && acknowledgementsAtStart.get(id) !== record) records.push(record);
      }
      queryClient.setQueryData<StudioV2WorkspaceMappingRecord[]>(studioV2WorkspaceKeys.mappings(), records);
    }).finally(() => {
      refreshing.current = false;
      if (mounted.current && refreshAgain.current) {
        refreshAgain.current = false;
        refreshAuthority();
      }
    });
  }, [applyMappingPatch, mappingsQuery, queryClient, ruleQueue, stateRef]);

  const flushRule = React.useCallback(async (ruleId: string) => {
    const queue = ruleQueue(ruleId);
    if (queue.running || queue.paused || !mounted.current) return;
    queue.running = true;
    try {
      while (mounted.current && !queue.paused && (queue.pending.size || queue.deletions.size)) {
        const deletion = queue.deletions.entries().next().value;
        if (deletion) {
          const [mappingId, removed] = deletion;
          try {
            const result = await deleteMappingMutation.mutateAsync(mappingId);
            markCleanupIncomplete(result.cleanup_status);
            queue.deletions.delete(mappingId);
            confirmedRecords.current.delete(mappingId);
            deletedRecords.current.add(mappingId);
            if (mounted.current) {
              queryClient.setQueryData<StudioV2WorkspaceMappingRecord[]>(studioV2WorkspaceKeys.mappings(), (records = []) => records.filter((record) => record.id !== mappingId));
              if (removed.restored && stateRef.current.mappings[removed.point.id]?.mapping_id === mappingId) {
                const state = stateRef.current;
                const mappings = { ...state.mappings };
                delete mappings[removed.point.id];
                queue.pending.delete(removed.point.id);
                const payload = { points: state.points.filter((point) => point.id !== removed.point.id), mappings };
                stateRef.current = workbenchV2Reducer(state, { type: 'SET_STATE', payload });
                actions.dispatch({ type: 'SET_STATE', payload });
              }
            }
          } catch (error) {
            queue.paused = true;
            if (mounted.current) {
              queue.failure = { save_state: 'save-error', save_error: errorMessageOf(error, t('errors.autosave_failed')), save_error_detail: { ...mappingSaveError(error), operation: 'delete' } };
              for (const [pointId, generation] of queue.pending) {
                const point = stateRef.current.points.find((item) => item.id === pointId);
                if (point?.rule_id === ruleId && generation === (generations.current.get(pointId) ?? 0)) applyMappingPatch(pointId, queue.failure);
              }
              // Keep failed removal visible so its retry has a concrete row.
              const state = stateRef.current;
              const replacement = state.mappings[removed.point.id];
              if (replacement && replacement.mapping_id !== mappingId) {
                removed.point = { ...removed.point, id: `mapping-removal:${mappingId}` };
                removed.mapping = { ...removed.mapping, point_id: removed.point.id };
              }
              if (!state.mappings[removed.point.id]) {
                removed.restored = true;
                const payload = { points: [...state.points, removed.point], mappings: { ...state.mappings, [removed.point.id]: removed.mapping } };
                stateRef.current = workbenchV2Reducer(state, { type: 'SET_STATE', payload });
                actions.dispatch({ type: 'SET_STATE', payload });
              }
              applyMappingPatch(removed.point.id, { save_state: 'save-error', save_error: errorMessageOf(error, t('errors.autosave_failed')), save_error_detail: { ...mappingSaveError(error), operation: 'delete' } });
            }
          }
          continue;
        }
        const [pointId, queuedGeneration] = queue.pending.entries().next().value!;
        queue.pending.delete(pointId);
        const point = stateRef.current.points.find((item) => item.id === pointId);
        const draft = stateRef.current.mappings[pointId];
        if (!point || !draft || point.rule_id !== ruleId || queuedGeneration !== (generations.current.get(pointId) ?? 0)) continue;
        if (!isStudioV2MappingValid(draft) || !canAutosaveMappingRule(stateRef.current, point)) {
          applyMappingPatch(pointId, { save_state: 'draft-invalid', save_error: null, save_error_detail: null });
          continue;
        }
        const generation = generations.current.get(pointId) ?? 0;
        const request = toStudioV2MappingUpdateRequest(point, draft);
        applyMappingPatch(pointId, { save_state: 'saving', save_error: null, save_error_detail: null });
        try {
          const saved = draft.mapping_id
            ? await updateMappingMutation.mutateAsync({ mappingId: draft.mapping_id, request })
            : await createMappingMutation.mutateAsync(toStudioV2MappingCreateRequest(point, draft));
          markCleanupIncomplete(saved.cleanup_status);
          const latestPoint = stateRef.current.points.find((item) => item.id === pointId);
          const latest = stateRef.current.mappings[pointId];
          if (!mounted.current || generation !== (generations.current.get(pointId) ?? 0) || !latest || latestPoint?.rule_id !== point.rule_id || latestPoint.address !== point.address || latestPoint.device_id !== point.device_id || latestPoint.data_type !== point.data_type) {
            // A create completed after the row was removed: clean up its owned record.
            if (mounted.current && !draft.mapping_id) queue.deletions.set(saved.id, { point, mapping: { ...draft, mapping_id: saved.id } });
            continue;
          }
          if (mounted.current) {
            deletedRecords.current.delete(saved.id);
            confirmedRecords.current.set(saved.id, saved);
            queryClient.setQueryData<StudioV2WorkspaceMappingRecord[]>(studioV2WorkspaceKeys.mappings(), (records = []) => [...records.filter((record) => record.id !== saved.id).map((record) => confirmedRecords.current.get(record.id) ?? record), saved]);
          }
          const unchanged = sameMappingValue(latest, draft);
          const hydrated = hydrateStudioV2Mapping(point, saved, unchanged ? undefined : latest);
          applyMappingPatch(pointId, { ...hydrated, save_state: unchanged ? 'saved' : isStudioV2MappingValid(latest) ? 'saving' : 'draft-invalid' });
          if (!unchanged) queue.pending.set(pointId, generation);
        } catch (error) {
          const latestPoint = stateRef.current.points.find((item) => item.id === pointId);
          if (generation === (generations.current.get(pointId) ?? 0) && latestPoint?.rule_id === ruleId && latestPoint.address === point.address && latestPoint.device_id === point.device_id && latestPoint.data_type === point.data_type) {
            queue.pending.set(pointId, generation);
          }
          if (!queue.pending.size) continue;
          queue.paused = true;
          const detail = mappingSaveError(error);
          queue.failure = { save_state: 'save-error', save_error: errorMessageOf(error, t('errors.autosave_failed')), save_error_detail: detail };
          for (const [pendingId, pendingGeneration] of queue.pending) {
            const pendingPoint = stateRef.current.points.find((item) => item.id === pendingId);
            if (pendingPoint?.rule_id !== ruleId || pendingGeneration !== (generations.current.get(pendingId) ?? 0)) continue;
            applyMappingPatch(pendingId, { save_state: 'save-error', save_error: errorMessageOf(error, t('errors.autosave_failed')), save_error_detail: detail });
          }
        }
      }
    } finally {
      queue.running = false;
      if (mounted.current && !queue.paused) refreshAuthority();
    }
  }, [actions, applyMappingPatch, createMappingMutation, deleteMappingMutation, markCleanupIncomplete, queryClient, refreshAuthority, ruleQueue, stateRef, t, updateMappingMutation]);

  const queueMappingSave = React.useCallback((pointId: string) => {
    const point = stateRef.current.points.find((item) => item.id === pointId);
    if (!point) return;
    const queue = ruleQueue(point.rule_id);
    draftVersions.current.set(pointId, (draftVersions.current.get(pointId) ?? 0) + 1);
    queue.pending.set(pointId, generations.current.get(pointId) ?? 0);
    if (queue.paused) {
      applyMappingPatch(pointId, queue.failure ?? { save_state: 'save-error' });
    } else {
      applyMappingPatch(pointId, { save_state: 'saving', save_error: null, save_error_detail: null });
      void flushRule(point.rule_id);
    }
  }, [applyMappingPatch, flushRule, ruleQueue, stateRef]);

  const afterMappingAction = React.useCallback((previousState: WorkbenchV2State, action: WorkbenchV2Action, nextState: WorkbenchV2State) => {
    switch (action.type) {
      case 'initMappingsForPoints': {
        const nextPointIDs = new Set(nextState.points.map((point) => point.id));
        Object.entries(previousState.mappings).forEach(([pointId, mapping]) => {
          const point = previousState.points.find((item) => item.id === pointId);
          const nextPoint = nextState.points.find((item) => item.id === pointId);
          if (point && (!nextPointIDs.has(pointId) || nextPoint?.address !== point.address || nextPoint.rule_id !== point.rule_id || nextPoint.device_id !== point.device_id || nextPoint.data_type !== point.data_type)) {
            generations.current.set(pointId, (generations.current.get(pointId) ?? 0) + 1);
            const queue = ruleQueue(point.rule_id);
            queue.pending.delete(pointId);
            if (mapping.mapping_id) queue.deletions.set(mapping.mapping_id, { point, mapping });
            void flushRule(point.rule_id);
          }
        });
        reconcilePersistedMappings(nextState.points);
        break;
      }
      case 'retryMappingSave': {
        const point = stateRef.current.points.find((item) => item.id === action.pointId);
        if (point) {
          const queue = ruleQueue(point.rule_id);
          queue.paused = false;
          if (!queue.deletions.size) queue.pending.set(action.pointId, generations.current.get(action.pointId) ?? 0);
          void flushRule(point.rule_id);
        }
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
  }, [flushRule, queueMappingSave, reconcilePersistedMappings, ruleQueue, stateRef]);

  return {
    afterMappingAction,
    mappingsQuery,
  };
}
