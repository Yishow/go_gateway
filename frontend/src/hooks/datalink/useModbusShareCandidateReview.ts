import { useEffect, useMemo, useState } from 'react';
import { sourceRuleCandidateAPI } from '../../services/sourceRuleCandidates';
import type { SourceRuleLocalModbusOutputCandidateView } from '../../types/sourceRuleCandidates';
import type { ModbusShareDesiredMapping, ModbusShareStatus } from '../../types/modbusShare';

export type ModbusShareCandidateReviewState = 'not-requested' | 'loading' | 'ready' | 'error';

interface CandidateScope {
  source_rule_id: string;
  revision_id: string;
}

function uniqueCandidateScopes(mappings: ModbusShareDesiredMapping[]): CandidateScope[] {
  const scopes = new Map<string, CandidateScope>();
  mappings.forEach((mapping) => {
    if (!mapping.source_rule_id || !mapping.source_rule_revision || scopes.has(mapping.source_rule_id)) {
      return;
    }
    scopes.set(mapping.source_rule_id, {
      source_rule_id: mapping.source_rule_id,
      revision_id: mapping.source_rule_revision,
    });
  });
  return [...scopes.values()];
}

export function useModbusShareCandidateReview(
  status: ModbusShareStatus | null | undefined,
  mappings: ModbusShareDesiredMapping[] | null,
) {
  const workspaceID = status?.workspace_id ?? '';
  const workspaceRevision = status?.workspace_revision ?? '';
  const scopes = useMemo(() => uniqueCandidateScopes(mappings ?? []), [mappings]);
  const [candidateSnapshots, setCandidateSnapshots] = useState<SourceRuleLocalModbusOutputCandidateView[]>([]);
  const [reviewState, setReviewState] = useState<ModbusShareCandidateReviewState>('not-requested');
  const scopeKey = JSON.stringify({ workspaceID, workspaceRevision, scopes });

  useEffect(() => {
    let cancelled = false;
    if (!workspaceID || !workspaceRevision || scopes.length === 0) {
      setCandidateSnapshots([]);
      setReviewState('not-requested');
      return () => { cancelled = true; };
    }

    setReviewState('loading');
    void Promise.all(scopes.map((scope) => sourceRuleCandidateAPI.get(scope.source_rule_id, {
      workspace_id: workspaceID,
      expected_workspace_revision: workspaceRevision,
      revision_id: scope.revision_id,
    }))).then((snapshots) => {
      if (cancelled) return;
      setCandidateSnapshots(snapshots.flatMap((snapshot) => snapshot.local_modbus_outputs.candidates));
      setReviewState('ready');
    }).catch(() => {
      if (!cancelled) {
        setCandidateSnapshots([]);
        setReviewState('error');
      }
    });
    return () => { cancelled = true; };
  }, [scopeKey, scopes, workspaceID, workspaceRevision]);

  const candidatesByTagID = useMemo(() => {
    const candidates = new Map<string, SourceRuleLocalModbusOutputCandidateView>();
    candidateSnapshots.forEach((candidate) => {
      if (candidate.tag_id) {
        candidates.set(candidate.tag_id, candidate);
      }
    });
    return candidates;
  }, [candidateSnapshots]);

  return { candidatesByTagID, reviewState };
}
