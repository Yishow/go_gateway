import { modbusShareAPI } from '../../../../services/datalink';
import type { ModbusShareMapping } from '../../../../types/datalink';
import { computeShareLayout } from './sourceRule';
import type { Mapping, Rule, WorkbenchV2State } from './types';

interface StudioV2ShareDesiredMapping {
  tag_id: string;
  register: number;
}

interface StudioV2ShareSyncPlan {
  desired: StudioV2ShareDesiredMapping[];
  ownedTagIds: string[];
  issues: string[];
}

interface ModbusShareRuntimeAPI {
  listMappings: () => Promise<ModbusShareMapping[]>;
  deleteMapping: (tagId: string) => Promise<void>;
  upsertMapping: (tagId: string, register: number) => Promise<ModbusShareMapping>;
}

const HUMAN_HOLDING_REGISTER_BASE = 40001;
const MAX_API_REGISTER = 65535;

type SharePlanningState = Pick<WorkbenchV2State, 'rules' | 'points' | 'mappings' | 'settings'>;

function nonEmptyTagId(mapping: Mapping | undefined): string | null {
  const tagId = mapping?.tag_id?.trim();
  return tagId ? tagId : null;
}

function durableWorkspaceTagId(mapping: Mapping | undefined, rules: Rule[]): string | null {
  const tagId = nonEmptyTagId(mapping);
  const ruleID = mapping?.rule_id?.trim();
  const persistedPointID = mapping?.persisted_point_id?.trim();
  if (!tagId || mapping?.persisted !== true || !ruleID || !persistedPointID) {
    return null;
  }

  const rule = rules.find((candidate) => candidate.id === ruleID && candidate.persisted === true);
  if (!rule) {
    return null;
  }
  if (mapping.workspace_id && rule.workspace_id && mapping.workspace_id !== rule.workspace_id) {
    return null;
  }
  return tagId;
}

function toApiHoldingRegister(humanRegister: number): number | null {
  if (!Number.isInteger(humanRegister) || humanRegister < HUMAN_HOLDING_REGISTER_BASE) {
    return null;
  }

  const register = humanRegister - HUMAN_HOLDING_REGISTER_BASE;
  return register <= MAX_API_REGISTER ? register : null;
}

function shareOwnedTagIds(mappings: Record<string, Mapping>, rules: Rule[]): string[] {
  return [...new Set(Object.values(mappings)
    .map((mapping) => durableWorkspaceTagId(mapping, rules))
    .filter((tagId): tagId is string => tagId !== null))]
    .sort();
}

/**
 * Builds the runtime plan from the current workspace state without making API calls.
 * Share registers are presented as human holding registers in the workbench and
 * are converted to the runtime's zero-based register offsets here.
 */
export function buildStudioV2ShareSyncPlan(state: SharePlanningState): StudioV2ShareSyncPlan {
  const desired: StudioV2ShareDesiredMapping[] = [];
  const issues: string[] = [];
  if (!state.settings.modbus_share.enabled) {
    return { desired, ownedTagIds: [], issues };
  }

  const baseRegister = state.settings.modbus_share.base_register;
  const layouts = computeShareLayout(state.rules, baseRegister);
  const seenRegisters = new Map<number, string>();
  const seenTagIds = new Map<string, number>();

  if (!Number.isInteger(baseRegister) || baseRegister < HUMAN_HOLDING_REGISTER_BASE) {
    issues.push(`invalid Share base register: ${String(baseRegister)}`);
  }

  state.rules
    .filter((rule) => rule.enabled && rule.share_enabled)
    .forEach((rule) => {
      const layout = layouts[rule.id];
      const points = state.points.filter((point) => (
        point.rule_id === rule.id && point.enabled && !point.skipped
      ));

      if (!layout || !Number.isFinite(layout.start) || !Number.isFinite(layout.stride)) {
        issues.push(`missing Share layout for rule ${rule.id}`);
        return;
      }

      points.forEach((point, index) => {
        const mapping = state.mappings[point.id];
        const tagId = durableWorkspaceTagId(mapping, state.rules);
        if (!mapping?.enabled || !tagId) {
          issues.push(`missing persisted enabled tag_id for Share point ${point.id}`);
          return;
        }

        const humanRegister = layout.start + index * layout.stride;
        const register = toApiHoldingRegister(humanRegister);
        if (register === null) {
          issues.push(`invalid Share holding register ${String(humanRegister)} for point ${point.id}`);
          return;
        }

        const previousTagRegister = seenTagIds.get(tagId);
        if (previousTagRegister !== undefined) {
          issues.push(`duplicate tag_id ${tagId} for Share registers ${previousTagRegister} and ${register}`);
          return;
        }

        const previousTagId = seenRegisters.get(register);
        if (previousTagId) {
          issues.push(`duplicate register ${register} for Share tags ${previousTagId} and ${tagId}`);
          return;
        }

        seenTagIds.set(tagId, register);
        seenRegisters.set(register, tagId);
        desired.push({ tag_id: tagId, register });
      });
    });

  return {
    desired,
    ownedTagIds: shareOwnedTagIds(state.mappings, state.rules),
    issues,
  };
}

/** Synchronizes owned runtime Share mappings and fails before activation on any planning/API error. */
export async function syncStudioV2ShareMappings(
  state: SharePlanningState,
  runtimeAPI: ModbusShareRuntimeAPI = modbusShareAPI,
): Promise<void> {
  const plan = buildStudioV2ShareSyncPlan(state);
  if (plan.issues.length > 0) {
    throw new Error(`Modbus Share synchronization blocked: ${plan.issues.join('; ')}`);
  }
  if (!state.settings.modbus_share.enabled) {
    return;
  }

  const runtimeMappings = await runtimeAPI.listMappings();
  const desiredTagIds = new Set(plan.desired.map((mapping) => mapping.tag_id));
  const ownedTagIds = new Set(plan.ownedTagIds);
  const desiredByRegister = new Map(plan.desired.map((mapping) => [mapping.register, mapping.tag_id]));

  for (const runtimeMapping of runtimeMappings) {
    const desiredTagId = desiredByRegister.get(runtimeMapping.register);
    if (
      desiredTagId !== undefined &&
      runtimeMapping.tag_id !== desiredTagId &&
      !ownedTagIds.has(runtimeMapping.tag_id)
    ) {
      throw new Error(
        `Modbus Share synchronization blocked: register ${runtimeMapping.register} is occupied by external tag ${runtimeMapping.tag_id}`,
      );
    }
  }

  for (const runtimeMapping of runtimeMappings) {
    if (ownedTagIds.has(runtimeMapping.tag_id) && !desiredTagIds.has(runtimeMapping.tag_id)) {
      await runtimeAPI.deleteMapping(runtimeMapping.tag_id);
    }
  }

  for (const mapping of plan.desired) {
    await runtimeAPI.upsertMapping(mapping.tag_id, mapping.register);
  }
}

/** Projects Share mappings before activation so synchronization failures stop the workspace transition. */
export async function activateStudioV2WorkspaceWithShare<T>(
  state: SharePlanningState,
  activateWorkspace: () => Promise<T>,
  runtimeAPI: ModbusShareRuntimeAPI = modbusShareAPI,
): Promise<T> {
  await syncStudioV2ShareMappings(state, runtimeAPI);
  return activateWorkspace();
}
