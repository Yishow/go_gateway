import * as React from 'react';
import { useQueries } from '@tanstack/react-query';
import type { SourceRuleRecord } from '../../../types/datalink';
import type { StudioV2WorkspaceDeviceRecord } from '../../../services/studioV2WorkspaceDevices';
import { sourceRuleCandidateAPI } from '../../../services/sourceRuleCandidates';
import {
  buildPersistedMCRuleState,
  type PersistedMCRuleState,
  type PersistedMCPointsScope,
} from '../../../features/datalink/workbench-v2/state/mcPersistedPoints';

export interface PersistedMCWorkspaceScope {
  workspace_id: string;
  workspace_revision: string;
}

interface PersistedMCQueryRule {
  id: string;
  revision_id: string;
}

interface LocalPersistedMCRule {
  id: string;
  device_id: string;
  persisted?: boolean;
  revision_id?: string;
}

function evidenceSignature(evidenceByRule: Readonly<Record<string, PersistedMCRuleState>>): string {
  return Object.entries(evidenceByRule)
    .sort(([left], [right]) => left.localeCompare(right))
    .map(([ruleID, evidence]) => `${ruleID}:${evidence.status}:${evidence.reason ?? ''}:${evidence.scope?.workspace_id ?? ''}:${evidence.scope?.workspace_revision ?? ''}:${evidence.scope?.revision_id ?? ''}:${evidence.points.map((point) => `${point.point_id}@${point.address}#${point.tag_id ?? ''}#${point.mapping_id ?? ''}`).join(',')}`)
    .join('|');
}

export function useStudioV2PersistedMCPoints(
  enabled: boolean,
  rules: readonly SourceRuleRecord[] | undefined,
  devices: readonly StudioV2WorkspaceDeviceRecord[] | undefined,
  workspaceScope?: PersistedMCWorkspaceScope,
  localRules: readonly LocalPersistedMCRule[] = [],
) {
  const mcDeviceIDs = React.useMemo(() => new Set((devices ?? []).filter((device) => device.protocol === 'mc_3e').map((device) => device.id)), [devices]);
  const serverMCRules = React.useMemo<SourceRuleRecord[]>(() => (rules ?? []).filter((rule) => mcDeviceIDs.has(rule.device_id)), [mcDeviceIDs, rules]);
  const ownedMCRuleIDs = React.useMemo(() => {
    const ids = new Set(serverMCRules.map((rule) => rule.id));
    localRules
      .filter((rule) => rule.persisted === true && mcDeviceIDs.has(rule.device_id))
      .forEach((rule) => ids.add(rule.id));
    return ids;
  }, [localRules, mcDeviceIDs, serverMCRules]);
  const mcRules = React.useMemo<PersistedMCQueryRule[]>(() => {
    return serverMCRules
      .filter((rule) => Boolean(rule.revision_id))
      .map((rule) => ({ id: rule.id, revision_id: rule.revision_id! }));
  }, [serverMCRules]);
  const scopeReady = Boolean(workspaceScope?.workspace_id && workspaceScope.workspace_revision);
  const queries = useQueries({
    queries: mcRules.map((rule) => ({
      queryKey: [
        'studio-v2',
        'persisted-mc-points',
        rule.id,
        workspaceScope?.workspace_id ?? '',
        workspaceScope?.workspace_revision ?? '',
        rule.revision_id,
      ],
      queryFn: () => sourceRuleCandidateAPI.get(rule.id, {
        workspace_id: workspaceScope!.workspace_id,
        expected_workspace_revision: workspaceScope!.workspace_revision,
        revision_id: rule.revision_id,
      }),
      enabled: enabled && scopeReady,
      retry: false,
      refetchOnWindowFocus: false,
      staleTime: 0,
    })),
  });

  const evidenceByRule = React.useMemo<Readonly<Record<string, PersistedMCRuleState>>>(() => {
    const result: Record<string, PersistedMCRuleState> = {};
    const queryByRule = new Map(queries.map((query, index) => [mcRules[index]?.id, query]));
    const sourceRules = new Map((rules ?? []).map((rule) => [rule.id, rule]));
    ownedMCRuleIDs.forEach((ruleID) => {
      const rule = sourceRules.get(ruleID);
      if (!rule) {
        result[ruleID] = { status: 'unavailable', points: [], reason: 'server_rule_missing' };
        return;
      }
      if (!rule.revision_id) {
        result[ruleID] = { status: 'unavailable', points: [], reason: 'server_rule_revision_missing' };
        return;
      }
      const scope: PersistedMCPointsScope | undefined = workspaceScope && scopeReady
        ? { ...workspaceScope, revision_id: rule.revision_id }
        : undefined;
      if (!scope) {
        result[ruleID] = { status: 'unavailable', points: [], reason: 'scope_unavailable' };
        return;
      }
      const query = queryByRule.get(rule.id);
      if (!query || query.isPending || (!query.data && !query.isError)) {
        result[ruleID] = { status: 'checking', scope, points: [] };
        return;
      }
      result[ruleID] = query.isError
        ? { status: 'unavailable', scope, points: [], reason: 'candidate_get_failed' }
        : buildPersistedMCRuleState(rule, scope, query.data);
    });
    return result;
  }, [mcRules, ownedMCRuleIDs, queries, rules, scopeReady, workspaceScope]);

  const signature = React.useMemo(() => evidenceSignature(evidenceByRule), [evidenceByRule]);
  const mcRuleIDs = ownedMCRuleIDs;
  const canMutateRule = React.useCallback((ruleID: string, revisionID?: string): boolean => {
    if (!mcRuleIDs.has(ruleID)) return false;
    const evidence = evidenceByRule[ruleID];
    return evidence?.status === 'compatible' && (!revisionID || evidence.scope?.revision_id === revisionID);
  }, [evidenceByRule, mcRuleIDs]);

  return {
    evidenceByRule,
    signature,
    canMutateRule,
    isChecking: Object.values(evidenceByRule).some((evidence) => evidence.status === 'checking'),
    hasConflict: Object.values(evidenceByRule).some((evidence) => evidence.status === 'conflict'),
    hasUnavailable: Object.values(evidenceByRule).some((evidence) => evidence.status === 'unavailable'),
  };
}
