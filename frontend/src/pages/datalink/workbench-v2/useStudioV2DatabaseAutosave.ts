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
  isSameStudioV2DatabaseTargetEdit,
  isStudioV2DatabaseConnectorValid,
  isStudioV2DatabaseTargetValid,
  mergeSavedStudioV2DatabaseConnector,
  mergeSavedStudioV2DatabaseRowGroups,
  mergeSavedStudioV2DatabaseTarget,
  resolveStudioV2DatabaseTargetPointID,
  withLatestStudioV2DatabaseRevisions,
  toStudioV2DatabaseConfigRequest,
  toStudioV2DatabaseTargetRequest,
} from '../../../features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import { workbenchV2Reducer, type WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { DbConnector, DbTarget, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import { isPersistedMCRule, persistedMCScopeKey } from '../../../features/datalink/workbench-v2/state/mcPersistedPoints';
import type { StudioV2WorkspaceDatabaseTargetRecord } from '../../../types/datalink';
import { getSafeErrorStateMessage } from '../../../utils/typedErrors';
import { createDatabaseSaveQueue, type DatabaseSaveJob, type DatabaseSaveQueue } from './databaseSaveQueue';
import { hasPendingSave } from './studioV2AutosaveBarrier';

type CanMutateRule = (ruleId: string) => boolean;

const allowAllRules: CanMutateRule = () => true;

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

function databaseMCScopeKey(state: WorkbenchV2State): string {
  return state.rules
    .filter((rule) => rule.mc_persisted_points)
    .sort((left, right) => left.id.localeCompare(right.id))
    .map((rule) => `${rule.id}:${persistedMCScopeKey(rule.mc_persisted_points)}`)
    .join('|');
}

const connectorEditFields: readonly (keyof DbConnector)[] = [
  'kind', 'name', 'host', 'port', 'database', 'username', 'password', 'clear_password', 'password_required',
  'schema', 'table', 'write_mode', 'write_interval_seconds', 'timestamp_column',
];

function sameDatabaseConnectorEdit(left: DbConnector, right: DbConnector): boolean {
  return connectorEditFields.every((field) => left[field] === right[field]);
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
  canMutateRule: CanMutateRule = allowAllRules,
) {
  const { t } = useTranslation('workbench-v2');
  const databaseConfigQuery = useStudioV2DatabaseConfigQuery(enabled);
  const databaseTargetsQuery = useStudioV2DatabaseTargetsQuery(enabled);
  const updateConfigMutation = useUpdateStudioV2DatabaseConfigMutation();
  const upsertTargetMutation = useUpsertStudioV2DatabaseTargetMutation();
  const lastConfigSignatureRef = React.useRef<string | null>(null);
  const lastTargetsSignatureRef = React.useRef<string | null>(null);
  const runSaveRef = React.useRef<(job: DatabaseSaveJob<WorkbenchV2State>) => Promise<boolean>>(async () => false);
  const saveQueueRef = React.useRef<DatabaseSaveQueue<WorkbenchV2State> | null>(null);
  const canMutateRuleRef = React.useRef(canMutateRule);
  const pendingConnectorRef = React.useRef(false);
  const pendingTargetsRef = React.useRef(new Set<string>());
  canMutateRuleRef.current = canMutateRule;
  if (!saveQueueRef.current) {
    saveQueueRef.current = createDatabaseSaveQueue<WorkbenchV2State>((job) => runSaveRef.current(job));
  }

  const applyConnectorPatch = React.useCallback((patch: Partial<DbConnector>) => {
    stateRef.current = workbenchV2Reducer(stateRef.current, { type: 'updateDbConnector', patch });
    actions.dispatch({ type: 'updateDbConnector', patch });
  }, [actions, stateRef]);

  const applyConnectorStatus = React.useCallback((patch: Partial<DbConnector>) => {
    const state = stateRef.current;
    const payload = { db: { ...state.db, connector: { ...state.db.connector, ...patch } } };
    stateRef.current = workbenchV2Reducer(state, { type: 'SET_STATE', payload });
    actions.dispatch({ type: 'SET_STATE', payload });
  }, [actions, stateRef]);

  const applySetupRevision = React.useCallback((setupRevision: string | undefined) => {
    if (!setupRevision || stateRef.current.db.connector.setup_revision === setupRevision) {
      return;
    }
    const connector = { ...stateRef.current.db.connector, setup_revision: setupRevision };
    stateRef.current = workbenchV2Reducer(stateRef.current, {
      type: 'SET_STATE',
      payload: { db: { ...stateRef.current.db, connector } },
    });
    actions.dispatch({ type: 'SET_STATE', payload: { db: { ...stateRef.current.db } } });
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

  const applyTargetStatus = React.useCallback((pointId: string, patch: Partial<DbTarget>) => {
    const target = stateRef.current.db.targets[pointId];
    if (!target) return;
    const payload = {
      db: {
        ...stateRef.current.db,
        targets: { ...stateRef.current.db.targets, [pointId]: { ...target, ...patch } },
      },
    };
    stateRef.current = workbenchV2Reducer(stateRef.current, { type: 'SET_STATE', payload });
    actions.dispatch({ type: 'SET_STATE', payload });
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
    // 尚未存檔的本地編輯不得被後端快照覆蓋，否則排隊中的儲存會改送舊值並顯示已儲存。
    Object.entries(currentTargets).forEach(([pointId, target]) => {
      if (hasPendingSave(target.save_state, target.persisted)) {
        nextTargets[pointId] = target;
      }
    });
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
    if (enabled) reconcilePersistedTargets();
  }, [enabled, reconcilePersistedTargets]);

  const runConnectorSave = React.useCallback(async (snapshot: WorkbenchV2State): Promise<boolean> => {
    const currentState = stateRef.current;
    if (currentState.rules.some((rule) => isPersistedMCRule(rule, currentState.devices) && !canMutateRuleRef.current(rule.id))) {
      pendingConnectorRef.current = true;
      applyConnectorStatus({ save_state: 'draft-invalid', save_error: null });
      return false;
    }
    const requestScopeKey = databaseMCScopeKey(currentState);
    // 以執行當下的設定與連線身分版本送出：排隊期間前一筆回覆已更新版本，快照裡的版本可能過期。
    const currentConnector = withLatestStudioV2DatabaseRevisions(snapshot.db.connector, stateRef.current.db.connector);
    const submittedConnector = currentConnector;
    if (!isStudioV2DatabaseConnectorValid(currentConnector)) {
      applyConnectorPatch({ save_state: 'draft-invalid', save_error: null });
      return false;
    }

    applyConnectorPatch({ save_state: 'saving', save_error: null });
    try {
      const savedConnector = await updateConfigMutation.mutateAsync(
        toStudioV2DatabaseConfigRequest(submittedConnector, snapshot.db.row_groups ?? []),
      );
      const currentMCRules = stateRef.current.rules.filter((rule) => isPersistedMCRule(rule, stateRef.current.devices));
      if (databaseMCScopeKey(stateRef.current) !== requestScopeKey || currentMCRules.some((rule) => !canMutateRuleRef.current(rule.id))) {
        pendingConnectorRef.current = !sameDatabaseConnectorEdit(stateRef.current.db.connector, submittedConnector);
        applyConnectorStatus({ save_state: 'draft-invalid', save_error: null });
        return !pendingConnectorRef.current;
      }
      pendingConnectorRef.current = false;
      stateRef.current = workbenchV2Reducer(stateRef.current, {
        type: 'SET_STATE',
        payload: {
          db: {
            ...stateRef.current.db,
            connector: mergeSavedStudioV2DatabaseConnector(stateRef.current.db.connector, currentConnector, savedConnector),
            row_groups: mergeSavedStudioV2DatabaseRowGroups(stateRef.current.db.row_groups, snapshot.db.row_groups, savedConnector.row_groups),
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
      return true;
    } catch (error) {
      const currentMCRules = stateRef.current.rules.filter((rule) => isPersistedMCRule(rule, stateRef.current.devices));
      if (databaseMCScopeKey(stateRef.current) !== requestScopeKey || currentMCRules.some((rule) => !canMutateRuleRef.current(rule.id))) {
        pendingConnectorRef.current = !sameDatabaseConnectorEdit(stateRef.current.db.connector, submittedConnector);
        applyConnectorStatus({ save_state: 'draft-invalid', save_error: null });
        return !pendingConnectorRef.current;
      }
      applyConnectorPatch({ save_state: 'save-error', save_error: errorMessageOf(error, t('errors.autosave_failed')) });
      return false;
    }
  }, [actions, applyConnectorPatch, applyConnectorStatus, stateRef, t, updateConfigMutation]);

  const runTargetSave = React.useCallback(async (pointId: string): Promise<boolean> => {
    const currentTarget = stateRef.current.db.targets[pointId];
    const currentMapping = stateRef.current.mappings[pointId];
    if (!currentTarget) {
      return true;
    }
    const currentPoint = stateRef.current.points.find((point) => point.id === pointId);
    if (currentPoint && !canMutateRuleRef.current(currentPoint.rule_id)) {
      pendingTargetsRef.current.add(pointId);
      applyTargetStatus(pointId, { save_state: 'draft-invalid', save_error: null });
      return false;
    }
    const requestScopeKey = persistedMCScopeKey(stateRef.current.rules.find((rule) => rule.id === currentPoint?.rule_id)?.mc_persisted_points);
    if (!isStudioV2DatabaseConnectorValid(stateRef.current.db.connector) || !isStudioV2DatabaseTargetValid(currentTarget, currentMapping)) {
      applyTargetPatch(pointId, { save_state: 'draft-invalid', save_error: null });
      return false;
    }
    const requestPointId = resolveStudioV2DatabaseTargetPointID(pointId, currentMapping);
    if (!requestPointId) {
      applyTargetPatch(pointId, { save_state: 'draft-invalid', save_error: null });
      return false;
    }

    const submittedTarget = currentTarget;
    applyTargetPatch(pointId, { save_state: 'saving', save_error: null });
    try {
      const savedTarget = await upsertTargetMutation.mutateAsync({
        pointId: requestPointId,
        request: toStudioV2DatabaseTargetRequest(currentTarget, stateRef.current.db.connector.setup_revision),
      });
      const currentRule = stateRef.current.rules.find((rule) => rule.id === currentPoint?.rule_id);
      if (currentPoint && (!canMutateRuleRef.current(currentPoint.rule_id) || persistedMCScopeKey(currentRule?.mc_persisted_points) !== requestScopeKey)) {
        const hasNewerEdit = !isSameStudioV2DatabaseTargetEdit(stateRef.current.db.targets[pointId], submittedTarget);
        if (hasNewerEdit) pendingTargetsRef.current.add(pointId);
        applyTargetStatus(pointId, { save_state: 'draft-invalid', save_error: null });
        return !hasNewerEdit;
      }
      pendingTargetsRef.current.delete(pointId);
      applySetupRevision(savedTarget.setup_revision);
      applyTargetPatch(pointId, mergeSavedStudioV2DatabaseTarget(stateRef.current.db.targets[pointId], currentTarget, savedTarget));
      return true;
    } catch (error) {
      const currentRule = stateRef.current.rules.find((rule) => rule.id === currentPoint?.rule_id);
      if (currentPoint && (!canMutateRuleRef.current(currentPoint.rule_id) || persistedMCScopeKey(currentRule?.mc_persisted_points) !== requestScopeKey)) {
        const hasNewerEdit = !isSameStudioV2DatabaseTargetEdit(stateRef.current.db.targets[pointId], submittedTarget);
        if (hasNewerEdit) pendingTargetsRef.current.add(pointId);
        applyTargetStatus(pointId, { save_state: 'draft-invalid', save_error: null });
        return !hasNewerEdit;
      }
      // 已有較新的同列編輯排隊時，由那筆儲存回報結果，不把較新的值標成失敗。
      if (isSameStudioV2DatabaseTargetEdit(stateRef.current.db.targets[pointId], currentTarget)) {
        applyTargetPatch(pointId, { save_state: 'save-error', save_error: errorMessageOf(error, t('errors.autosave_failed')) });
      }
      return false;
    }
  }, [applySetupRevision, applyTargetPatch, applyTargetStatus, stateRef, t, upsertTargetMutation]);

  runSaveRef.current = (job) => (job.kind === 'connector' ? runConnectorSave(job.snapshot) : runTargetSave(job.pointId));

  const queueConnectorSave = React.useCallback((snapshot?: WorkbenchV2State) => {
    saveQueueRef.current?.enqueueConnector(snapshot ?? stateRef.current);
  }, [stateRef]);

  const queueTargetSave = React.useCallback((pointId: string) => {
    // 入列即標為儲存中：連線儲存失敗而暫停欄位儲存時，也不會繼續顯示「已儲存」。
    const target = stateRef.current.db.targets[pointId];
    const point = stateRef.current.points.find((item) => item.id === pointId);
    if (point && !canMutateRuleRef.current(point.rule_id)) {
      pendingTargetsRef.current.add(pointId);
      applyTargetPatch(pointId, { save_state: 'draft-invalid', save_error: null });
      return;
    }
    if (target && target.save_state !== 'saving') {
      applyTargetPatch(pointId, { save_state: 'saving', save_error: null });
    }
    saveQueueRef.current?.enqueueTarget(pointId);
  }, [applyTargetPatch, stateRef]);

  React.useEffect(() => {
    if (pendingConnectorRef.current && !stateRef.current.rules.some((rule) => isPersistedMCRule(rule, stateRef.current.devices) && !canMutateRuleRef.current(rule.id))) {
      pendingConnectorRef.current = false;
      queueConnectorSave(stateRef.current);
    }
    for (const pointId of pendingTargetsRef.current) {
      const point = stateRef.current.points.find((item) => item.id === pointId);
      if (!point || !canMutateRuleRef.current(point.rule_id)) continue;
      pendingTargetsRef.current.delete(pointId);
      queueTargetSave(pointId);
    }
  }, [canMutateRule, queueConnectorSave, queueTargetSave, stateRef]);

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
