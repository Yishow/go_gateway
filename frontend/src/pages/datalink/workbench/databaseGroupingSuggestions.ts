import type {
  SourceRuleCandidateSetView,
  SourceRuleDatabaseOutputCandidateView,
  SourceRuleTagCandidateView,
} from '../../../types/sourceRuleCandidates';
import type { SourceRuleTagReviewDecision } from '../../../types/sourceRuleTagReviewDecisions';
import type { WorkbenchTagGroupingOverride } from './workbenchTypes';

export type DatabaseGroupingSuggestion = {
  groupKey: string | null;
  columnName: string;
  inferredGroupKey: string | null;
  inferredColumnName: string;
  memberTagKeys: ReadonlyArray<string>;
  writeIntervalSeconds: number | null;
  overridden: boolean;
};

type DatabaseGroupingOverridesByPointId = Readonly<
  Record<string, WorkbenchTagGroupingOverride>
>;

export function normalizeDatabaseGroupingColumnName(value: string): string {
  const replaced = value.trim().replaceAll('/', '_');
  return replaced === '' ? '' : replaced.toLowerCase();
}

export function inferDatabaseGroupAndColumnFromTagKey(tagKey: string): {
  groupKey: string | null;
  columnName: string;
} {
  const trimmed = tagKey.trim();
  const slashIndex = trimmed.indexOf('/');
  if (slashIndex === -1) {
    return { groupKey: null, columnName: '' };
  }

  const groupKey = trimmed.slice(0, slashIndex).trim();
  const columnName = normalizeDatabaseGroupingColumnName(trimmed.slice(slashIndex + 1));
  if (groupKey === '' || columnName === '') {
    return { groupKey: null, columnName: '' };
  }

  return { groupKey, columnName };
}

export function normalizeDatabaseGroupingOverrideDraft(
  groupKeyDraft: string,
  columnNameDraft: string,
): WorkbenchTagGroupingOverride {
  const trimmedGroupKey = groupKeyDraft.trim();
  if (trimmedGroupKey === '') {
    return { groupKey: null, columnName: '' };
  }

  return {
    groupKey: trimmedGroupKey,
    columnName: normalizeDatabaseGroupingColumnName(columnNameDraft),
  };
}

function getEffectiveTagKey(
  candidate: SourceRuleTagCandidateView,
  currentDecision: SourceRuleTagReviewDecision | null,
) {
  return currentDecision?.tag_key?.trim() || candidate.tag_key;
}

export function buildTagDatabaseGroupingSuggestions(
  tagCandidates: ReadonlyArray<SourceRuleTagCandidateView>,
  decisionByCandidateId: ReadonlyMap<string, SourceRuleTagReviewDecision | null>,
  databaseCandidates: ReadonlyArray<SourceRuleDatabaseOutputCandidateView>,
  overridesByPointId: DatabaseGroupingOverridesByPointId = {},
): Record<string, DatabaseGroupingSuggestion> {
  const databaseCandidateByPointId = new Map(
    databaseCandidates.map((candidate) => [candidate.point_id, candidate] as const),
  );

  const resolvedRows = tagCandidates.map((candidate) => {
    const currentDecision = decisionByCandidateId.get(candidate.id) ?? null;
    const effectiveTagKey = getEffectiveTagKey(candidate, currentDecision);
    const inferred = inferDatabaseGroupAndColumnFromTagKey(effectiveTagKey);
    const databaseCandidate = databaseCandidateByPointId.get(candidate.point_id);
    const override = overridesByPointId[candidate.point_id];
    const usesPersistedDatabaseScope = Boolean(
      databaseCandidate?.mapping_id || databaseCandidate?.connector_id,
    );
    const baseGroupKey = usesPersistedDatabaseScope
      ? databaseCandidate?.group_key ?? null
      : inferred.groupKey;
    const baseColumnName = usesPersistedDatabaseScope
      ? databaseCandidate?.column_name || inferred.columnName
      : inferred.columnName;
    return {
      candidateId: candidate.id,
      effectiveTagKey,
      writeIntervalSeconds: databaseCandidate?.write_interval_seconds ?? null,
      inferred,
      resolved: override
        ? {
            groupKey: override.groupKey,
            columnName: override.columnName,
          }
        : {
            groupKey: baseGroupKey ?? null,
            columnName: baseColumnName,
          },
      overridden: Boolean(override),
    };
  });

  const memberKeysByGroup = new Map<string, string[]>();
  for (const row of resolvedRows) {
    const bucketKey = row.resolved.groupKey ?? `__single__:${row.candidateId}`;
    const members = memberKeysByGroup.get(bucketKey) ?? [];
    members.push(row.effectiveTagKey);
    memberKeysByGroup.set(bucketKey, members);
  }

  return Object.fromEntries(
    resolvedRows.map((row) => {
      const bucketKey = row.resolved.groupKey ?? `__single__:${row.candidateId}`;
      return [
        row.candidateId,
        {
          groupKey: row.resolved.groupKey,
          columnName: row.resolved.columnName,
          inferredGroupKey: row.inferred.groupKey,
          inferredColumnName: row.inferred.columnName,
          memberTagKeys: memberKeysByGroup.get(bucketKey) ?? [row.effectiveTagKey],
          writeIntervalSeconds: row.writeIntervalSeconds,
          overridden: row.overridden,
        } satisfies DatabaseGroupingSuggestion,
      ];
    }),
  );
}

export function applyDatabaseGroupingOverridesToReviewSet(
  reviewSet: SourceRuleCandidateSetView<SourceRuleDatabaseOutputCandidateView> | null,
  overridesByPointId: DatabaseGroupingOverridesByPointId = {},
): SourceRuleCandidateSetView<SourceRuleDatabaseOutputCandidateView> | null {
  if (!reviewSet || Object.keys(overridesByPointId).length === 0) {
    return reviewSet;
  }

  return {
    ...reviewSet,
    candidates: reviewSet.candidates.map((candidate) => {
      const override = overridesByPointId[candidate.point_id];
      if (!override) {
        return candidate;
      }

      return {
        ...candidate,
        group_key: override.groupKey,
        column_name: override.columnName,
        write_interval_seconds:
          override.writeIntervalSeconds !== undefined
            ? override.writeIntervalSeconds
            : candidate.write_interval_seconds,
      };
    }),
  };
}
