import { editStudioV2DatabaseConnector, hydrateStudioV2DatabaseConnector, hydrateStudioV2DatabaseRowGroups, mergeSavedStudioV2DatabaseConnector, mergeSavedStudioV2DatabaseRowGroups } from './studioV2DatabaseAutosave';
import type { WorkbenchV2State } from './types';
import type { WorkbenchV2Action } from './useWorkbenchV2State';

/**
 * 處理 Step 4 Database 相關狀態變更的子 Reducer
 * 落地設計決策：「Commit log 序列：純函式 + reducer 串聯」
 * 
 * @param state 當前工作台狀態
 * @param action 觸發的動作
 * @returns 變更後的狀態
 */
export function dbReducer(state: WorkbenchV2State, action: WorkbenchV2Action): WorkbenchV2State {
  switch (action.type) {
    case 'updateDbConnector':
      return {
        ...state,
        db: {
          ...state.db,
          connector: editStudioV2DatabaseConnector(state.db.connector, action.patch)
        }
      };

    case 'hydrateDbConnector': {
      const current = state.db.connector;
      // Refetch is canonical observation, not acknowledgement of a local draft.
      if (['saving', 'draft-invalid', 'save-error'].includes(current.save_state ?? '') ||
        current.workspace_id && current.workspace_id !== action.record.workspace_id) return state;
      return { ...state, db: { ...state.db,
        connector: hydrateStudioV2DatabaseConnector(action.record, current),
        row_groups: hydrateStudioV2DatabaseRowGroups(action.record.row_groups),
      } };
    }
    case 'acknowledgeDbConnectorSave': {
      const connector = mergeSavedStudioV2DatabaseConnector(state.db.connector, action.submitted, action.record);
      if (connector === state.db.connector) return state;
      const sameSelection = connector.connector_id === action.record.id &&
        (state.db.connector.field_generations?.connector_id ?? 0) <= (action.submitted.edit_generation ?? 0);
      return { ...state, db: { ...state.db, connector,
        row_groups: sameSelection ? mergeSavedStudioV2DatabaseRowGroups(state.db.row_groups, action.rowGroups, action.record.row_groups) : state.db.row_groups,
      } };
    }

    case 'replaceDbTargets':
      return { ...state, db: { ...state.db, targets: action.targets } };

    case 'upsertDbTarget':
      return {
        ...state,
        db: {
          ...state.db,
          targets: {
            ...state.db.targets,
            [action.pointId]: action.target
          }
        }
      };

    case 'updateDbTarget': {
      const existing = state.db.targets[action.pointId];
      if (!existing && !action.patch.row_group_id) return state;
      const baseTarget = existing ?? {
        tag_id: '',
        column_name: '',
        enabled: true,
      };
      return {
        ...state,
        db: {
          ...state.db,
          targets: {
            ...state.db.targets,
            [action.pointId]: {
              ...baseTarget,
              ...action.patch
            }
          }
        }
      };
    }

    case 'setAllDbTargetsEnabled':
      return {
        ...state,
        db: {
          ...state.db,
          targets: Object.fromEntries(
            Object.entries(state.db.targets).map(([pointId, target]) => [
              pointId,
              {
                ...target,
                enabled: action.enabled,
              },
            ]),
          ),
        },
      };

    case 'autoAssignDbTargets':
      return {
        ...state,
        db: {
          ...state.db,
          targets: {
            ...action.targets
          }
        }
      };

    case 'setDbRowGroups':
      return {
        ...state,
        db: {
          ...state.db,
          row_groups: action.rowGroups
        }
      };

    case 'startCommit':
      return {
        ...state,
        committed: false,
        commit: {
          status: 'committing',
          logs: [],
          started_at: new Date().toISOString()
        }
      };

    case 'appendCommitLog': {
      if (!state.commit) return state;
      // 避免重複新增相同的 log
      const logs = [...state.commit.logs];
      const alreadyExists = logs.some(l => l.label === action.log.label);
      if (alreadyExists) return state;

      return {
        ...state,
        commit: {
          ...state.commit,
          logs: [...logs, { ...action.log, status: 'success' }]
        }
      };
    }

    case 'completeCommit': {
      if (!state.commit) return state;
      const completed = new Set(state.completed);
      completed.add(4);
      return {
        ...state,
        completed,
        committed: true,
        commit: {
          ...state.commit,
          status: 'success',
          finished_at: new Date().toISOString()
        }
      };
    }

    case 'resetCommit':
      return {
        ...state,
        committed: false,
        commit: {
          status: 'idle',
          logs: []
        }
      };

    default:
      return state;
  }
}
