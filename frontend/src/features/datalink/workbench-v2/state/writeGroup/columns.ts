import type { TargetType } from '../types';
import type { DbColumn } from '../dbSchemas';
import type { GroupCandidate } from './candidates';

export type ColumnCompatibility = 'compatible' | 'incompatible' | 'unknown';

type ValueFamily = 'bool' | 'int' | 'float' | 'text';

/** SQL dialects whose exact-value rules are implemented by the backend codec. */
export type WriteGroupSqlDialect = 'portable' | 'sqlite' | 'postgres';

function familyOf(target: TargetType): ValueFamily {
  if (target === 'bool') return 'bool';
  if (target === 'string') return 'text';
  if (target === 'float32' || target === 'float64') return 'float';
  return 'int';
}

/** Column families by SQL type text. Matching is by substring so `varchar(64)` and `double precision` both resolve. */
function columnFamilies(type: string): ValueFamily[] | null {
  const text = type.trim().toLowerCase();
  if (!text) return null;
  if (/(double|float|real|numeric|decimal)/.test(text)) return ['float', 'int'];
  if (/(bigint|int8|smallint|int2|tinyint|integer|int4|^int\b|serial)/.test(text)) return ['int', 'bool'];
  if (/(bool|bit)/.test(text)) return ['bool'];
  if (/(char|text|string|clob|citext)/.test(text)) return ['text'];
  return null;
}

function normalizedType(type: string): string {
  return type.trim().toLowerCase().replace(/\s+/g, ' ');
}

/**
 * uint64 cannot use a signed integer or floating-point column. TEXT is the
 * portable exact representation; PostgreSQL NUMERIC is safe only when it has
 * at least twenty integer digits and no scale. SQLite text-affinity aliases
 * are safe when the actual dialect is known.
 */
function exactUint64Column(type: string, dialect: WriteGroupSqlDialect): boolean {
  const text = normalizedType(type);
  if (text === 'text') return true;
  // Keep this in step with sqliteAffinity: any valid declaration whose base
  // contains CHAR, CLOB or TEXT has text affinity. SQLite does not enforce
  // the optional length/scale annotation, so the codec stores the full text.
  if (dialect === 'sqlite') {
    const sqliteDeclaration = /^([a-z][a-z0-9 ]*)(?:\(\s*\d+\s*(?:,\s*\d+\s*)?\))?$/.exec(text);
    const base = sqliteDeclaration?.[1] ?? '';
    // sqliteAffinity checks INT before CHAR/CLOB/TEXT. Preserve that order so
    // declarations such as CHARINT or INTEGER TEXT remain integer-backed.
    if (base && !base.includes('int') && /(?:char|clob|text)/.test(base)) return true;
  }
  if (dialect !== 'postgres') return false;
  const numeric = /^(?:numeric|decimal)\s*(?:\(\s*(\d+)(?:\s*,\s*(\d+))?\s*\))?$/.exec(text);
  if (!numeric) return false;
  const precision = numeric[1] === undefined ? Number.POSITIVE_INFINITY : Number(numeric[1]);
  const scale = numeric[2] === undefined ? 0 : Number(numeric[2]);
  return scale === 0 && precision >= 20;
}

/**
 * Whether a column can hold a tag type exactly. `unknown` means the type text
 * is not one the editor recognises; the backend readiness check decides then.
 */
export function columnCompatibility(
  target: TargetType,
  column: DbColumn,
  dialect: WriteGroupSqlDialect = 'portable',
): ColumnCompatibility {
  const families = columnFamilies(column.type);
  if (!families) return 'unknown';
  if (target === 'uint64') return exactUint64Column(column.type, dialect) ? 'compatible' : 'incompatible';
  const wanted = familyOf(target);
  if (!families.includes(wanted)) return 'incompatible';
  // A float column only counts for integers when it is numeric/decimal; double and real are lossy for 64-bit values.
  if (wanted === 'int' && /(double|float|real)/.test(column.type.toLowerCase())) return 'incompatible';
  return 'compatible';
}

export type AssignmentStatus = 'confirmed' | 'suggested' | 'repair';

export interface ColumnAssignment {
  column: string;
  status: AssignmentStatus;
}

export interface ColumnSuggestions {
  assignments: Record<string, ColumnAssignment>;
  /** Candidate keys that have neither a confirmed nor a name-based suggestion. */
  unmatched: string[];
}

function tail(tagKey: string): string {
  return (tagKey.split('.').pop() ?? tagKey).toLowerCase();
}

function nameMatches(column: string, token: string): boolean {
  if (!token) return false;
  const name = column.toLowerCase();
  return name === token || name.endsWith(`_${token}`) || name.startsWith(`${token}_`);
}

/**
 * Reviewable column proposals from real metadata. A confirmed assignment is kept
 * while its column still exists and fits the tag, otherwise it is marked for
 * repair and never silently rebound. A proposal comes only from the tag name
 * and a compatible column; there is no index fallback, no wraparound and no
 * column is proposed twice. Primary key columns are never proposed.
 */
export function suggestColumns(
  candidates: GroupCandidate[],
  columns: DbColumn[],
  confirmed: Record<string, string>,
  dialect: WriteGroupSqlDialect = 'portable',
): ColumnSuggestions {
  const assignments: Record<string, ColumnAssignment> = {};
  const used = new Set<string>();
  const byName = new Map(columns.map((column) => [column.name, column]));

  for (const candidate of candidates) {
    const chosen = confirmed[candidate.key];
    if (!chosen) continue;
    const column = byName.get(chosen);
    const valid = column !== undefined && columnCompatibility(candidate.target_type, column, dialect) !== 'incompatible';
    assignments[candidate.key] = { column: chosen, status: valid ? 'confirmed' : 'repair' };
    used.add(chosen);
  }
  const unmatched: string[] = [];
  for (const candidate of candidates) {
    if (assignments[candidate.key]) continue;
    const token = tail(candidate.tag_key);
    const match = columns.find((column) => !column.primary_key && !used.has(column.name) &&
      nameMatches(column.name, token) && columnCompatibility(candidate.target_type, column, dialect) !== 'incompatible');
    if (match) {
      assignments[candidate.key] = { column: match.name, status: 'suggested' };
      used.add(match.name);
    } else {
      unmatched.push(candidate.key);
    }
  }
  return { assignments, unmatched };
}

export interface ColumnConflict {
  column: string;
  member_keys: string[];
}

interface ConflictMember {
  key: string;
  column: string;
  entity_key: string;
}

/**
 * Members sharing a column collide unless every one of them carries its own,
 * different entity key (distinct rows). A missing entity key is unproven
 * identity, so it never makes a shared column legal.
 */
export function findColumnConflicts(members: ConflictMember[]): ColumnConflict[] {
  const byColumn = new Map<string, ConflictMember[]>();
  for (const member of members) {
    if (!member.column) continue;
    const column = member.column.trim().toLowerCase();
    const list = byColumn.get(column) ?? [];
    list.push(member);
    byColumn.set(column, list);
  }
  const conflicts: ColumnConflict[] = [];
  for (const list of byColumn.values()) {
    if (list.length < 2) continue;
    const entities = list.map((member) => member.entity_key.trim());
    const distinct = entities.every((entity) => entity !== '') && new Set(entities).size === entities.length;
    if (!distinct) conflicts.push({ column: list[0].column, member_keys: list.map((member) => member.key) });
  }
  return conflicts;
}
