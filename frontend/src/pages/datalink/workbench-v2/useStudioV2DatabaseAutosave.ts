import * as React from 'react';
import { useTranslation } from 'react-i18next';
import {
  useStudioV2DatabaseConfigQuery,
  useStudioV2DatabaseTargetsQuery,
  useUpdateStudioV2DatabaseConfigMutation,
  useUpsertStudioV2DatabaseTargetMutation,
} from '../../../hooks/datalink/useStudioV2WorkspaceDatabase';
import {
  hydrateStudioV2DatabaseConnector,
  hydrateStudioV2DatabaseRowGroups,
  hydrateStudioV2DatabaseTarget,
  isStudioV2DatabaseConnectorValid,
  isStudioV2DatabaseTargetValid,
  resolveStudioV2DatabaseTargetPointID,
  toStudioV2DatabaseConfigRequest,
  toStudioV2DatabaseTargetRequest,
} from '../../../features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import { workbenchV2Reducer, type WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { DbConnector, DbTarget, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import type { StudioV2WorkspaceDatabaseTargetRecord } from '../../../types/datalink';
import { getSafeErrorStateMessage } from '../../../utils/typedErrors';

type SaveMeta = {
  inFlight: boolean;
  pending: boolean;
};

function sameDatabaseTarget(current: DbTarget | undefined, next: DbTarget): boolean {
  return (
    current?.tag_id === next.tag_id &&
    current?.column_name === next.column_name &&
    current?.enabled === next.enabled &&
    current?.row_group_id === next.row_group_id &&
    current?.point_id === next.point_id &&
    current?.row_id === next.row_id &&
    current?.workspace_id === next.workspace_id &&
    current?.persisted === next.persisted &&
    current?.save_state === next.save_state &&
    current?.save_error === next.save_error
  );
}

function errorMessageOf(error: unknown, fallback: string): string {
  return getSafeErrorStateMessage(error, fallback);
}

function saveMetaFor(store: Record<string, SaveMeta>, key: string): SaveMeta {
  if (!store[key]) {
    store[key] = { inFlight: false, pending: false };
  }
  return store[key];
}

function buildPersistedPointAliases(state: WorkbenchV2State): Map<string, string> {
  const aliases = new Map<string, string>();
  Object.values(state.mappings).forEach((mapping) => {
    aliases.set(mapping.point_id, mapping.point_id);
    if (mapping.persisted_point_id) {
      aliases.set(mapping.persisted_point_id, mapping.point_id);
    }
  });
  return aliases;
}

function buildHydratedTargets(
  aliases: Map<string, string>,
  records: StudioV2WorkspaceDatabaseTargetRecord[] | undefined,
): Record<string, DbTarget> {
  return (records ?? []).reduce<Record<string, DbTarget>>((accumulator, record) => {
    const currentPointId = aliases.get(record.point_id);
    if (!currentPointId) {
      return accumulator;
    }
    accumulator[currentPointId] = hydrateStudioV2DatabaseTarget(
      { ...record, point_id: currentPointId },
      accumulator[currentPointId],
    );
    return accumulator;
  }, {});
}

export function useStudioV2DatabaseAutosave(
  actions: { dispatch: (action: WorkbenchV2Action) => void },
  stateRef: React.MutableRefObject<WorkbenchV2State>,
  enabled: boolean,
) {
  const { t } = useTranslation('workbench-v2');
  const databaseConfigQuery = useStudioV2DatabaseConfigQuery(enabled);
  const databaseTargetsQuery = useStudioV2DatabaseTargetsQuery(enabled);
  const updateConfigMutation = useUpdateStudioV2DatabaseConfigMutation();
  const upsertTargetMutation = useUpsertStudioV2DatabaseTargetMutation();
  const lastConfigSignatureRef = React.useRef<string | null>(null);
  const lastTargetsSignatureRef = React.useRef<string | null>(null);
  const connectorSaveMetaRef = React.useRef<SaveMeta>({ inFlight: false, pending: false });
  const pendingConnectorStateRef = React.useRef<WorkbenchV2State | null>(null);
  const deferredTargetSaveRef = React.useRef<Set<string>>(new Set());
  const flushTargetSaveRef = React.useRef<(pointId: string) => void>(() => undefined);
  const targetSaveMetaRef = React.useRef<Record<string, SaveMeta>>({});

  const applyConnectorPatch = React.useCallback((patch: Partial<DbConnector>) => {
    stateRef.current = workbenchV2Reducer(stateRef.current, { type: 'updateDbConnector', patch });
    actions.dispatch({ type: 'updateDbConnector', patch });
  }, [actions, stateRef]);

  const applyTargetPatch = React.useCallback((pointId: string, patch: Partial<DbTarget>) => {
    const target = stateRef.current.db.targets[pointId];
    stateRef.current = workbenchV2Reducer(stateRef.current, target
      ? { type: 'updateDbTarget', pointId, patch }
      : {
        type: 'upsertDbTarget',
        pointId,
        target: {
          tag_id: '',
          column_name: '',
          enabled: true,
          ...patch,
        } as DbTarget,
      });
    if (target) {
      actions.dispatch({ type: 'updateDbTarget', pointId, patch });
      return;
    }
    actions.dispatch({
      type: 'upsertDbTarget',
      pointId,
      target: {
        tag_id: '',
        column_name: '',
        enabled: true,
        ...patch,
      } as DbTarget,
    });
  }, [actions, stateRef]);

  React.useEffect(() => {
    if (!enabled || !databaseConfigQuery.isSuccess || !databaseConfigQuery.data) {
      return;
    }
    const signature = `${databaseConfigQuery.data.id}:${databaseConfigQuery.data.updated_at}`;
    if (lastConfigSignatureRef.current === signature) {
      return;
    }
    lastConfigSignatureRef.current = signature;
    stateRef.current = workbenchV2Reducer(stateRef.current, {
      type: 'SET_STATE',
      payload: {
        db: {
          ...stateRef.current.db,
          connector: hydrateStudioV2DatabaseConnector(databaseConfigQuery.data, stateRef.current.db.connector),
          row_groups: hydrateStudioV2DatabaseRowGroups(databaseConfigQuery.data.row_groups),
        },
      },
    });
    actions.dispatch({
      type: 'SET_STATE',
      payload: {
        db: {
          ...stateRef.current.db,
        },
      },
    });
  }, [actions, databaseConfigQuery.data, databaseConfigQuery.isSuccess, enabled, stateRef]);

  const reconcilePersistedTargets = React.useCallback(() => {
    if (!databaseTargetsQuery.isSuccess) {
      return;
    }
    const aliases = buildPersistedPointAliases(stateRef.current);
    const aliasSignature = Array.from(aliases.entries())
      .sort(([left], [right]) => left.localeCompare(right))
      .map(([persistedPointId, currentPointId]) => `${persistedPointId}:${currentPointId}`)
      .join('|');
    const signature = databaseTargetsQuery.data
      .map((record) => `${record.id}:${record.updated_at}`)
      .sort()
      .join('|') + `|${aliasSignature}`;
    if (lastTargetsSignatureRef.current === signature) {
      return;
    }
    lastTargetsSignatureRef.current = signature;

    const nextTargets = buildHydratedTargets(aliases, databaseTargetsQuery.data);
    const currentTargets = stateRef.current.db.targets;
    const hasSameTargets = Object.keys(currentTargets).length === Object.keys(nextTargets).length &&
      Object.entries(nextTargets).every(([pointId, target]) => sameDatabaseTarget(currentTargets[pointId], target));
    if (hasSameTargets) {
      return;
    }

    stateRef.current = workbenchV2Reducer(stateRef.current, {
      type: 'SET_STATE',
      payload: {
        db: {
          ...stateRef.current.db,
          targets: nextTargets,
        },
      },
    });
    actions.dispatch({
      type: 'SET_STATE',
      payload: {
        db: {
          ...stateRef.current.db,
          targets: nextTargets,
        },
      },
    });
  }, [actions, databaseTargetsQuery.data, databaseTargetsQuery.isSuccess, stateRef]);

  React.useEffect(() => {
    reconcilePersistedTargets();
  }, [reconcilePersistedTargets]);

  const flushConnectorSave = React.useCallback(async (snapshot?: WorkbenchV2State) => {
    const meta = connectorSaveMetaRef.current;
    const saveState = snapshot ?? stateRef.current;
    const currentConnector = saveState.db.connector;
    if (!isStudioV2DatabaseConnectorValid(currentConnector)) {
      applyConnectorPatch({ save_state: 'draft-invalid', save_error: null });
      meta.inFlight = false;
      meta.pending = false;
      return;
    }

    meta.inFlight = true;
    applyConnectorPatch({ save_state: 'saving', save_error: null });
    let savedConfig = false;
    let pendingConnectorState: WorkbenchV2State | null | undefined;
    let pendingTargetIds: string[] = [];
    try {
      const savedConnector = await updateConfigMutation.mutateAsync(
        toStudioV2DatabaseConfigRequest(currentConnector, saveState.db.row_groups ?? []),
      );
      savedConfig = true;
      stateRef.current = workbenchV2Reducer(stateRef.current, {
        type: 'SET_STATE',
        payload: {
          db: {
            ...stateRef.current.db,
            connector: hydrateStudioV2DatabaseConnector(savedConnector, currentConnector),
            row_groups: hydrateStudioV2DatabaseRowGroups(savedConnector.row_groups),
          },
        },
      });
      actions.dispatch({
        type: 'SET_STATE',
        payload: {
          db: {
            ...stateRef.current.db,
          },
        },
      });
    } catch (error) {
      applyConnectorPatch({ save_state: 'save-error', save_error: errorMessageOf(error, t('errors.autosave_failed')) });
    } finally {
      meta.inFlight = false;
      if (meta.pending) {
        meta.pending = false;
        pendingConnectorState = pendingConnectorStateRef.current;
        pendingConnectorStateRef.current = null;
      } else if (savedConfig) {
        pendingTargetIds = Array.from(deferredTargetSaveRef.current);
        deferredTargetSaveRef.current.clear();
      }
    }
    if (pendingConnectorState !== undefined) {
      void flushConnectorSave(pendingConnectorState ?? undefined);
      return;
    }
    pendingTargetIds.forEach((pointId) => flushTargetSaveRef.current(pointId));
  }, [actions, applyConnectorPatch, stateRef, t, updateConfigMutation]);

  const queueConnectorSave = React.useCallback((snapshot?: WorkbenchV2State) => {
    const meta = connectorSaveMetaRef.current;
    if (meta.inFlight) {
      meta.pending = true;
      pendingConnectorStateRef.current = snapshot ?? stateRef.current;
      return;
    }
    void flushConnectorSave(snapshot ?? stateRef.current);
  }, [flushConnectorSave, stateRef]);

  const flushTargetSave = React.useCallback(async (pointId: string) => {
    const meta = saveMetaFor(targetSaveMetaRef.current, pointId);
    const currentTarget = stateRef.current.db.targets[pointId];
    const currentMapping = stateRef.current.mappings[pointId];
    if (!currentTarget) {
      delete targetSaveMetaRef.current[pointId];
      return;
    }
    if (connectorSaveMetaRef.current.inFlight) {
      deferredTargetSaveRef.current.add(pointId);
      meta.inFlight = false;
      meta.pending = false;
      return;
    }

    if (!isStudioV2DatabaseConnectorValid(stateRef.current.db.connector) || !isStudioV2DatabaseTargetValid(currentTarget, currentMapping)) {
      applyTargetPatch(pointId, { save_state: 'draft-invalid', save_error: null });
      meta.inFlight = false;
      meta.pending = false;
      return;
    }
    const requestPointId = resolveStudioV2DatabaseTargetPointID(pointId, currentMapping);
    if (!requestPointId) {
      applyTargetPatch(pointId, { save_state: 'draft-invalid', save_error: null });
      meta.inFlight = false;
      meta.pending = false;
      return;
    }

    meta.inFlight = true;
    applyTargetPatch(pointId, { save_state: 'saving', save_error: null });
    try {
      const savedTarget = await upsertTargetMutation.mutateAsync({
        pointId: requestPointId,
        request: toStudioV2DatabaseTargetRequest(currentTarget),
      });
      applyTargetPatch(pointId, hydrateStudioV2DatabaseTarget(savedTarget, currentTarget));
    } catch (error) {
      applyTargetPatch(pointId, { save_state: 'save-error', save_error: errorMessageOf(error, t('errors.autosave_failed')) });
    } finally {
      meta.inFlight = false;
      if (meta.pending) {
        meta.pending = false;
        void flushTargetSave(pointId);
      }
    }
  }, [applyTargetPatch, stateRef, t, upsertTargetMutation]);
  flushTargetSaveRef.current = flushTargetSave;

  const queueTargetSave = React.useCallback((pointId: string) => {
    const meta = saveMetaFor(targetSaveMetaRef.current, pointId);
    if (meta.inFlight) {
      meta.pending = true;
      return;
    }
    void flushTargetSave(pointId);
  }, [flushTargetSave]);

  const afterDatabaseAction = React.useCallback((action: WorkbenchV2Action, nextState: WorkbenchV2State) => {
    switch (action.type) {
      case 'initMappingsForPoints':
        reconcilePersistedTargets();
        break;
      case 'updateDbConnector':
      case 'setDbRowGroups':
        queueConnectorSave(nextState);
        break;
      case 'upsertDbTarget':
      case 'updateDbTarget':
        queueTargetSave(action.pointId);
        break;
      case 'setAllDbTargetsEnabled':
        Object.keys(nextState.db.targets).forEach(queueTargetSave);
        break;
      case 'autoAssignDbTargets':
        Object.keys(nextState.db.targets).forEach(queueTargetSave);
        break;
      default:
        break;
    }
  }, [queueConnectorSave, queueTargetSave, reconcilePersistedTargets]);

  return {
    afterDatabaseAction,
    databaseConfigQuery,
    databaseTargetsQuery,
  };
}
