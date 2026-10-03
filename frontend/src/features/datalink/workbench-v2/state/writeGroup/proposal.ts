import type { TargetType } from '../types';
import type { GroupCandidate } from './candidates';
import type { WriteGroupSqlDialect } from './columns';

export interface ColumnProposal {
  key: string;
  /** Proposed column name; it does not exist in the table. */
  column: string;
  /** Portable SQL type text for the proposal. */
  sql_type: string;
}

const SQL_TYPE: Record<TargetType, string> = {
  bool: 'BOOLEAN', int16: 'INTEGER', int32: 'INTEGER', int64: 'BIGINT', uint16: 'INTEGER', uint32: 'BIGINT',
  uint64: 'TEXT', float32: 'DOUBLE PRECISION', float64: 'DOUBLE PRECISION', string: 'TEXT',
};

function sqlTypeFor(target: TargetType, dialect: WriteGroupSqlDialect): string {
  if (target === 'uint64' && dialect === 'postgres') return 'NUMERIC(20,0)';
  if (dialect === 'postgres' && ['int16', 'int32', 'int64', 'uint16', 'uint32'].includes(target)) return 'BIGINT';
  return SQL_TYPE[target];
}

function safeName(tagKey: string): string {
  const tail = (tagKey.split('.').pop() ?? tagKey).toLowerCase().replace(/[^a-z0-9_]+/g, '_').replace(/^_+|_+$/g, '');
  return /^[a-z_]/.test(tail) ? tail : `v_${tail || 'value'}`;
}

/**
 * Column proposals for Tags that have no column. They are suggestions for a
 * table change the operator decides on; nothing here is observed metadata, and
 * a proposal never becomes an assignment by itself.
 */
export function proposeColumns(
  unmatched: GroupCandidate[],
  existing: string[],
  dialect: WriteGroupSqlDialect = 'portable',
): ColumnProposal[] {
  const taken = new Set(existing.map((name) => name.toLowerCase()));
  return unmatched.map((candidate) => {
    let name = safeName(candidate.tag_key);
    let suffix = 2;
    while (taken.has(name)) name = `${safeName(candidate.tag_key)}_${suffix++}`;
    taken.add(name);
    return { key: candidate.key, column: name, sql_type: sqlTypeFor(candidate.target_type, dialect) };
  });
}
