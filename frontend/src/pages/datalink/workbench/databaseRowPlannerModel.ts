import type { DatabaseWriteMode } from '../../../types/datalink';
import type { SourceRuleDatabaseOutputCandidateView } from '../../../types/sourceRuleCandidates';
import { buildTableKey } from './databaseTargetBoardUtils';

export type DatabaseRowPlanIssueCode =
  | 'candidate_blocked'
  | 'tag_pending'
  | 'connector_required'
  | 'table_required'
  | 'timestamp_required'
  | 'column_required'
  | 'connector_mismatch'
  | 'table_mismatch'
  | 'write_mode_mismatch'
  | 'timestamp_mismatch'
  | 'interval_mismatch';

export type DatabaseRowPlanMember = {
  candidateId: string;
  pointId: string;
  tagId: string | null;
  tagKey: string;
  groupKey: string | null;
  displayName: string;
  address: string;
  dataType: SourceRuleDatabaseOutputCandidateView['data_type'];
  columnName: string;
  intervalSeconds: number | null;
  mappingId: string | null;
  status: SourceRuleDatabaseOutputCandidateView['status'];
  blockingReason: string | null;
  connectorId: string;
  tableKey: string;
  writeMode: DatabaseWriteMode;
  timestampColumn: string;
};

export type DatabaseRowPlan = {
  id: string;
  kind: 'grouped' | 'single';
  groupKey: string | null;
  displayKey: string;
  status: 'ready' | 'blocked';
  issueCodes: ReadonlyArray<DatabaseRowPlanIssueCode>;
  members: ReadonlyArray<DatabaseRowPlanMember>;
  intervalSeconds: number | null;
};

type BuildDatabaseRowPlansInput = {
  candidates: ReadonlyArray<SourceRuleDatabaseOutputCandidateView>;
  selectedConnectorId: string;
  selectedTableKey: string;
  selectedWriteMode: DatabaseWriteMode;
  selectedTimestampColumn: string;
  connectorDefaultWriteIntervalSeconds: number | null;
};

function normalizeTestIdSegment(value: string) {
  return value.replace(/[^a-zA-Z0-9_-]+/g, '-');
}

function resolveCandidateTableKey(candidate: SourceRuleDatabaseOutputCandidateView) {
  if (!candidate.table_schema || !candidate.table_name) {
    return '';
  }
  return buildTableKey({
    schema: candidate.table_schema,
    name: candidate.table_name,
    columns: [],
  });
}

function uniqueNonEmpty(values: ReadonlyArray<string>) {
  return Array.from(new Set(values.filter((value) => value !== '')));
}

function uniqueIntervals(values: ReadonlyArray<number | null>) {
  return Array.from(new Set(values.filter((value): value is number => value !== null)));
}

export function buildDatabaseRowPlans({
  candidates,
  selectedConnectorId,
  selectedTableKey,
  selectedWriteMode,
  selectedTimestampColumn,
  connectorDefaultWriteIntervalSeconds,
}: BuildDatabaseRowPlansInput): DatabaseRowPlan[] {
  const rowBuckets = new Map<string, DatabaseRowPlanMember[]>();

  for (const candidate of candidates) {
    const groupKey = candidate.group_key ?? null;
    const bucketKey = groupKey
      ? `group:${groupKey}`
      : `single:${candidate.id}`;
    const member: DatabaseRowPlanMember = {
      candidateId: candidate.id,
      pointId: candidate.point_id,
      tagId: candidate.tag_id ?? null,
      tagKey: candidate.tag_key,
      groupKey,
      displayName: candidate.display_name,
      address: candidate.address,
      dataType: candidate.data_type,
      columnName: candidate.column_name ?? '',
      intervalSeconds: candidate.write_interval_seconds ?? connectorDefaultWriteIntervalSeconds,
      mappingId: candidate.mapping_id ?? null,
      status: candidate.status,
      blockingReason: candidate.blocking_reason ?? null,
      connectorId: selectedConnectorId || candidate.connector_id || '',
      tableKey: selectedTableKey || resolveCandidateTableKey(candidate),
      writeMode:
        selectedConnectorId || selectedTableKey
          ? selectedWriteMode
          : candidate.write_mode || selectedWriteMode,
      timestampColumn:
        (selectedConnectorId || selectedTableKey) && selectedWriteMode === 'upsert'
          ? selectedTimestampColumn
          : candidate.write_mode === 'upsert'
            ? candidate.timestamp_column ?? ''
            : '',
    };
    const bucket = rowBuckets.get(bucketKey) ?? [];
    bucket.push(member);
    rowBuckets.set(bucketKey, bucket);
  }

  return Array.from(rowBuckets.entries())
    .map(([bucketKey, members]) => {
      const firstMember = members[0];
      const rowGroupKey = bucketKey.startsWith('group:') ? (firstMember.groupKey ?? null) : null;
      const rowId = bucketKey.startsWith('group:')
        ? `group-${members
            .map((member) => normalizeTestIdSegment(member.candidateId))
            .sort()
            .join('-')}`
        : `single-${normalizeTestIdSegment(firstMember.candidateId)}`;
      const issueCodes: DatabaseRowPlanIssueCode[] = [];
      const connectorIds = uniqueNonEmpty(members.map((member) => member.connectorId));
      const tableKeys = uniqueNonEmpty(members.map((member) => member.tableKey));
      const writeModes = uniqueNonEmpty(members.map((member) => member.writeMode));
      const timestampColumns = uniqueNonEmpty(members.map((member) => member.timestampColumn));
      const intervalValues = uniqueIntervals(members.map((member) => member.intervalSeconds));

      if (members.some((member) => member.status === 'blocked' || member.blockingReason)) {
        issueCodes.push('candidate_blocked');
      }
      if (members.some((member) => !member.tagId)) {
        issueCodes.push('tag_pending');
      }
      if (!members.every((member) => member.connectorId)) {
        issueCodes.push('connector_required');
      }
      if (!members.every((member) => member.tableKey)) {
        issueCodes.push('table_required');
      }
      if (members.some((member) => member.writeMode === 'upsert' && member.timestampColumn === '')) {
        issueCodes.push('timestamp_required');
      }
      if (members.some((member) => member.columnName.trim() === '')) {
        issueCodes.push('column_required');
      }
      if (connectorIds.length > 1) {
        issueCodes.push('connector_mismatch');
      }
      if (tableKeys.length > 1) {
        issueCodes.push('table_mismatch');
      }
      if (writeModes.length > 1) {
        issueCodes.push('write_mode_mismatch');
      }
      if (timestampColumns.length > 1) {
        issueCodes.push('timestamp_mismatch');
      }
      if (intervalValues.length > 1) {
        issueCodes.push('interval_mismatch');
      }

      return {
        id: rowId,
        kind: bucketKey.startsWith('group:') ? 'grouped' : 'single',
        groupKey: rowGroupKey,
        displayKey: rowGroupKey || firstMember.tagKey,
        status: issueCodes.length === 0 ? 'ready' : 'blocked',
        issueCodes,
        members: [...members].sort((left, right) => left.columnName.localeCompare(right.columnName)),
        intervalSeconds:
          intervalValues.length === 1
            ? intervalValues[0]
            : (members[0]?.intervalSeconds ?? connectorDefaultWriteIntervalSeconds ?? null),
      } satisfies DatabaseRowPlan;
    })
    .sort((left, right) => left.displayKey.localeCompare(right.displayKey));
}
