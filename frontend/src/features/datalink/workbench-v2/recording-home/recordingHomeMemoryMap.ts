import type { ProtocolType, SourceRuleRecord } from '../../../../types/datalink';
import { addressParser } from '../../../../utils/addressParser';
import { dataTypeWidth } from '../state/sourceRule';
import type { PersistedMCRuleState } from '../state/mcPersistedPoints';

/** Why an area cannot claim free cells. `string-width`: the frontend (10) and the
 * backend source-rule planner (1) disagree on string point spans. */
export type MemoryIssue = 'string-width' | 'bit-width' | 'mc-proof' | 'mc-scope' | 'list-scope' | 'invalid-rule' | 'out-of-range';
export interface MemoryRuleInput {
  id: string; label: string; start_address: string; count: number; data_type: SourceRuleRecord['data_type']; skipped_addresses: readonly string[];
  /** MC only: current persisted point addresses; undefined means no current proof. */
  proven?: readonly string[];
  /** MC only: the list had no valid workspace revision, so no proof could be bound to it. */
  scopeMissing?: boolean;
}
export interface MemoryRuleIssue { ruleId: string; label: string; issue: MemoryIssue }
interface Segment { ruleId: string; label: string; dataType: string; start: number; width: number; count: number; skipped: ReadonlySet<number>; partial: boolean }
export interface MemoryArea {
  key: string; protocol: ProtocolType; baseAddress: string; unit: 'bit' | 'word'; columns: number; base: number; max: number;
  issues: MemoryRuleIssue[]; ruleIds: string[]; blocked: boolean; segments: Segment[];
}
export interface MemoryOwner { ruleId: string; label: string; dataType: string; valueAddress: string; part: number; width: number; skipped: boolean; partial: boolean }
export interface MemoryCell { offset: number; address: string; state: 'empty' | 'used' | 'conflict' | 'unconfirmed'; owners: MemoryOwner[] }
/** `unplaced`: ranges no area can hold. `global`: issues that make every area unconfirmed
 * (unplaced ranges, and MC ranges read without a valid workspace revision). A missing proof
 * under a valid revision blocks only that range's own area. */
export interface MemoryMap { supported: boolean; areas: MemoryArea[]; unplaced: MemoryRuleIssue[]; global: MemoryRuleIssue[] }

const MAX_CELLS = 256;
// Area anchors only; unit, radix, base and limit are all read back from the parser.
const AREA_BASES: Partial<Record<ProtocolType, readonly string[]>> = {
  modbus_tcp: ['40001', '30001', '00001', '10001'], modbus_rtu: ['40001', '30001', '00001', '10001'], modbus_udp: ['40001', '30001', '00001', '10001'],
  mc_3e: ['D0', 'W0', 'M0', 'B0', 'X0', 'Y0'], fatek_fbs: ['D0', 'R0', 'M0', 'S0', 'X0', 'Y0', 'T0', 'C0'],
};

/** Exclusive protocol limits the frontend parser does not encode. FATEK mirrors
 * `ComponentType.AddressLimit()` in internal/protocol/fatek/address_native.go:
 * D/R 65536; X/Y/M/S/T/C 10^(FormatLen - len(Name)) = 10000. Used for this view's geometry only. */
const PROTOCOL_LIMITS: Partial<Record<ProtocolType, Readonly<Record<string, number>>>> = {
  fatek_fbs: { D: 65536, R: 65536, X: 10000, Y: 10000, M: 10000, S: 10000, T: 10000, C: 10000 },
};

function parse(text: string, protocol: ProtocolType) {
  try { return addressParser.parse(text.trim().toUpperCase(), protocol); } catch { return undefined; }
}
/** The parser-formatted address of `offset`, or undefined when the parser cannot express it. */
export function memoryAddress(area: Pick<MemoryArea, 'key' | 'protocol' | 'baseAddress' | 'base'>, offset: number): string | undefined {
  try {
    const text = addressParser.offset(area.baseAddress, offset - area.base, area.protocol);
    const parsed = addressParser.parse(text, area.protocol);
    return parsed.area === area.key && parsed.startNumber === offset ? text : undefined;
  } catch { return undefined; }
}
function areaLimit(area: Pick<MemoryArea, 'key' | 'protocol' | 'baseAddress' | 'base'>): number {
  let good = area.base, step = 1, bad: number | undefined;
  while (bad === undefined) {
    const next = Math.min(area.base + step, Number.MAX_SAFE_INTEGER);
    if (memoryAddress(area, next) === undefined) bad = next;
    else if (next === Number.MAX_SAFE_INTEGER) return next;
    else { good = next; step *= 2; }
  }
  while (bad - good > 1) { const mid = Math.floor((good + bad) / 2); if (memoryAddress(area, mid) === undefined) bad = mid; else good = mid; }
  return good;
}
function createArea(protocol: ProtocolType, baseAddress: string): MemoryArea | undefined {
  const parsed = parse(baseAddress, protocol);
  if (!parsed) return;
  const shape = { key: parsed.area, protocol, baseAddress, base: parsed.startNumber };
  const prefix = baseAddress.slice(0, -1);
  const columns = memoryAddress(shape, shape.base + 8) === `${prefix}10` ? 8 : memoryAddress(shape, shape.base + 16) === `${prefix}10` ? 16 : 10;
  const limit = PROTOCOL_LIMITS[protocol]?.[shape.key];
  return { ...shape, unit: parsed.type, columns, max: Math.min(areaLimit(shape), limit === undefined ? Number.MAX_SAFE_INTEGER : limit - 1), issues: [], ruleIds: [], blocked: false, segments: [] };
}
function runs(offsets: number[], width: number): { start: number; count: number }[] {
  const sorted = [...new Set(offsets)].sort((a, b) => a - b), result: { start: number; count: number }[] = [];
  for (const offset of sorted) {
    const last = result.at(-1);
    if (last && offset === last.start + last.count * width) last.count += 1; else result.push({ start: offset, count: 1 });
  }
  return result;
}

/** Projects saved ranges onto parser-defined areas. Unknown geometry never yields free cells.
 * `listScopeMissing`: the MC range list itself had no valid workspace revision, so even an
 * empty list cannot prove that any address is free. */
export function buildMemoryMap(protocol: ProtocolType, rules: readonly MemoryRuleInput[], options: { listScopeMissing?: boolean } = {}): MemoryMap {
  const bases = AREA_BASES[protocol];
  const areas = bases?.map(base => createArea(protocol, base)).filter((area): area is MemoryArea => Boolean(area)) ?? [];
  if (!areas.length) return { supported: false, areas: [], unplaced: [], global: [] };
  const unplaced: MemoryRuleIssue[] = [];
  for (const rule of rules) {
    const issueOf = (issue: MemoryIssue) => ({ ruleId: rule.id, label: rule.label, issue });
    const parsed = parse(rule.start_address, protocol), area = parsed && areas.find(item => item.key === parsed.area);
    if (!parsed || !area) { unplaced.push(issueOf('invalid-rule')); continue; }
    area.ruleIds.push(rule.id);
    const width = dataTypeWidth(rule.data_type), base = { ruleId: rule.id, label: rule.label, dataType: rule.data_type, width };
    const issue: MemoryIssue | undefined = !Number.isSafeInteger(rule.count) || rule.count <= 0 ? 'invalid-rule' : parsed.startNumber > area.max ? 'out-of-range' : rule.data_type === 'string' ? 'string-width'
      : area.unit === 'bit' && width > 1 ? 'bit-width' : protocol === 'mc_3e' && rule.scopeMissing ? 'mc-scope' : protocol === 'mc_3e' && !rule.proven ? 'mc-proof' : undefined;
    if (issue) {
      area.issues.push(issueOf(issue));
      // The saved start address is certain; the span after it is not.
      if (issue !== 'mc-proof' && issue !== 'mc-scope' && issue !== 'out-of-range') area.segments.push({ ...base, start: parsed.startNumber, width: 1, count: 1, skipped: new Set(), partial: true });
      continue;
    }
    // Never truncate a span past the parser limit and then call the rest free.
    const fits = (start: number, values: number) => Number.isSafeInteger(values * width) && Number.isSafeInteger(start + values * width) && start + values * width - 1 <= area.max;
    const skipped = new Set(rule.skipped_addresses.map(text => parse(text, protocol)).filter(item => item?.area === area.key).map(item => item!.startNumber));
    if (protocol !== 'mc_3e') {
      if (fits(parsed.startNumber, rule.count)) area.segments.push({ ...base, start: parsed.startNumber, count: rule.count, skipped, partial: false });
      else area.issues.push(issueOf('out-of-range'));
      continue;
    }
    const proven = rule.proven!.map(text => parse(text, protocol));
    if (proven.some(item => !item || item.area !== area.key)) { area.issues.push(issueOf('invalid-rule')); continue; }
    if (![...proven.map(item => item!.startNumber), ...skipped].every(start => fits(start, 1))) { area.issues.push(issueOf('out-of-range')); continue; }
    for (const run of runs([...proven.map(item => item!.startNumber), ...skipped], width)) area.segments.push({ ...base, ...run, skipped, partial: false });
  }
  const global = [...unplaced, ...areas.flatMap(area => area.issues.filter(item => item.issue === 'mc-scope'))];
  if (protocol === 'mc_3e' && options.listScopeMissing) global.push({ ruleId: '', label: '', issue: 'list-scope' });
  if (global.length) areas.forEach(area => { area.blocked = true; });
  return { supported: true, areas, unplaced, global };
}

function ownersAt(area: MemoryArea, offset: number): MemoryOwner[] {
  const owners: MemoryOwner[] = [];
  for (const segment of area.segments) {
    const index = Math.floor((offset - segment.start) / segment.width);
    if (offset < segment.start || index >= segment.count) continue;
    const valueStart = segment.start + index * segment.width;
    owners.push({ ruleId: segment.ruleId, label: segment.label, dataType: segment.dataType, valueAddress: memoryAddress(area, valueStart) ?? '', part: offset - valueStart + 1,
      width: segment.partial ? dataTypeWidth(segment.dataType) : segment.width, skipped: segment.skipped.has(valueStart), partial: segment.partial });
  }
  return owners;
}
/** At most MAX_CELLS cells from `from`, clipped to the parser's area limit. */
export function memoryWindow(area: MemoryArea, from: number, length: number, blocked: boolean): MemoryCell[] {
  const start = Math.max(area.base, from), end = Math.min(area.max, start + Math.min(length, MAX_CELLS) - 1), cells: MemoryCell[] = [];
  const unknown = blocked || area.blocked || area.issues.length > 0;
  for (let offset = start; offset <= end; offset += 1) {
    const owners = ownersAt(area, offset), address = memoryAddress(area, offset);
    if (address === undefined) continue;
    cells.push({ offset, address, owners, state: owners.length > 1 ? 'conflict' : owners.length ? 'used' : unknown ? 'unconfirmed' : 'empty' });
  }
  return cells;
}
export function memoryRowStart(area: MemoryArea, offset: number): number {
  return area.base + Math.floor((offset - area.base) / area.columns) * area.columns;
}
export function initialMemoryWindow(area: MemoryArea): { start: number; length: number } {
  const lastEnd = Math.max(area.base, ...area.segments.map(segment => segment.start + segment.count * segment.width));
  const rows = Math.min(Math.floor(MAX_CELLS / area.columns), Math.max(2, Math.ceil((lastEnd - area.base) / area.columns) + 1));
  return { start: area.base, length: rows * area.columns };
}
export function locateMemoryOffset(area: MemoryArea, text: string): number | undefined {
  const parsed = parse(text, area.protocol);
  return parsed && parsed.area === area.key && parsed.startNumber >= area.base && parsed.startNumber <= area.max ? parsed.startNumber : undefined;
}
/** Saved start/end/footprint for list rows; end and footprint are omitted when unconfirmed. */
export function describeMemoryRule(protocol: ProtocolType, rule: MemoryRuleInput) {
  const map = buildMemoryMap(protocol, [rule]), area = map.areas.find(item => item.ruleIds.includes(rule.id));
  if (!area) return { area: undefined, unit: undefined, start: rule.start_address, end: undefined, cells: undefined, issue: 'invalid-rule' as MemoryIssue };
  const start = memoryAddress(area, parse(rule.start_address, protocol)!.startNumber) ?? rule.start_address;
  const issue = area.issues[0]?.issue;
  if (issue) return { area: area.key, unit: area.unit, start, end: undefined, cells: undefined, issue };
  const first = Math.min(...area.segments.map(item => item.start)), last = Math.max(...area.segments.map(item => item.start + item.count * item.width)) - 1;
  return { area: area.key, unit: area.unit, start: memoryAddress(area, first) ?? start, end: memoryAddress(area, last), cells: area.segments.reduce((sum, item) => sum + item.count * item.width, 0), issue: undefined };
}
type ListedRule = Pick<SourceRuleRecord, 'id' | 'revision_id' | 'start_address' | 'count' | 'data_type' | 'naming_prefix' | 'skipped_addresses'>;
/** MC proof counts only for the list's workspace and workspace revision, the listed rule revision and compatible points. */
export function memoryRuleInputs(listed: readonly ListedRule[], protocol: ProtocolType, workspace: { id: string; revision?: string }, evidenceFor: (id: string) => PersistedMCRuleState | undefined): MemoryRuleInput[] {
  return listed.map(rule => {
    const evidence = protocol === 'mc_3e' ? evidenceFor(rule.id) : undefined;
    const current = evidence?.status === 'compatible' && evidence.points.length > 0 && Boolean(rule.revision_id) && evidence.scope?.workspace_id === workspace.id &&
      evidence.scope.revision_id === rule.revision_id && Boolean(workspace.revision) && evidence.scope.workspace_revision === workspace.revision;
    return { id: rule.id, label: rule.naming_prefix, start_address: rule.start_address, count: rule.count, data_type: rule.data_type,
      skipped_addresses: rule.skipped_addresses ?? [], proven: current ? evidence.points.map(point => point.address) : undefined, scopeMissing: protocol === 'mc_3e' && !workspace.revision || undefined };
  });
}
