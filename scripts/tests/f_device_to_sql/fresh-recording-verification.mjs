// Pure F recording fixture and evidence checks. Callers own all I/O and pass
// independent SQL/API observations here; this module starts no process.
export const FRESH_REGISTERS = Object.freeze([
  215, 1013, 32, 0, 0, 1, 1, 0x3fc0, 0, 0x405e, 0xdd2f, 0x1a9f, 0xbe77,
  0, 32768, 0, 0, 0, 65535, 65535, 65535, 65535, 0x4131, 0, 0, 0, 0, 0, 0, 0, 0, 0,
]);

export const FRESH_RULES = Object.freeze([
  { name: 'A int16', start: 40001, count: 1, dataType: 'int16', prefix: 'A_INT16_' },
  { name: 'A uint16', start: 40002, count: 1, dataType: 'uint16', prefix: 'A_UINT16_' },
  { name: 'A uint64 precise', start: 40003, count: 1, dataType: 'uint64', prefix: 'A_U64A_' },
  { name: 'A bool', start: 40007, count: 1, dataType: 'bool', prefix: 'A_BOOL_' },
  { name: 'A float32', start: 40008, count: 1, dataType: 'float32', prefix: 'A_F32_' },
  { name: 'A float64', start: 40010, count: 1, dataType: 'float64', prefix: 'A_F64_' },
  { name: 'A uint64 high', start: 40015, count: 1, dataType: 'uint64', prefix: 'A_U64B_' },
  { name: 'A uint64 max', start: 40019, count: 1, dataType: 'uint64', prefix: 'A_U64C_' },
  { name: 'A string', start: 40023, count: 1, dataType: 'string', prefix: 'A_TEXT_' },
]);

export const FRESH_TARGET_TYPES = Object.freeze([
  'int16', 'uint16', 'uint64', 'bool', 'float32', 'float64', 'uint64', 'uint64', 'string',
]);

export const FRESH_EXPECTED_VALUES = Object.freeze([
  '215', '1013', '9007199254740993', '1', '1.5', '123.456',
  '9223372036854775808', '18446744073709551615', 'A1',
]);

export function expectedValuesByTag(tags, expectedValues = FRESH_EXPECTED_VALUES) {
  return Object.fromEntries(tags.map((tag, index) => [tag.tag_key, expectedValues[index]]));
}

function parseProvenance(value) {
  if (typeof value !== 'string') return null;
  try {
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed : null;
  } catch {
    return null;
  }
}

function memberIdentity(member) {
  return JSON.stringify([member?.device_id, member?.point_id, member?.tag_id]);
}

function sameIdentity(left, right) {
  return left !== undefined && right !== undefined && String(left) === String(right);
}

function rowEffect(outbox, row) {
  return outbox.find((item) => sameIdentity(item.record_id, row.record_id));
}

function rowReceipt(receipts, effect) {
  return effect && receipts.find((item) => item.effect_key === effect.effect_key);
}

function utcMillis(value) {
  const text = String(value ?? '');
  if (!/(?:Z|\+00:00)$/i.test(text)) return NaN;
  return Date.parse(text);
}

function sameBucketStart(left, right) {
  if (left === right) return true;
  const leftMillis = utcMillis(left);
  const rightMillis = utcMillis(right);
  return Number.isFinite(leftMillis) && leftMillis === rightMillis;
}

function preciseUTC(value) {
  const match = String(value ?? '').match(/^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(?:\.(\d{1,9}))?(?:Z|\+00:00)$/i);
  return match && Number.isFinite(utcMillis(value)) ? `${match[1]}.${(match[2] ?? '').padEnd(9, '0')}Z` : null;
}

function checkConsecutiveBuckets(rows, failures, label) {
  const bucketStarts = [];
  for (const row of rows) {
    if (!bucketStarts.includes(row.bucket_start)) bucketStarts.push(row.bucket_start);
    if (bucketStarts.length === 3) break;
  }
  const times = bucketStarts.map(utcMillis);
  if (bucketStarts.length !== 3 || times.some((time) => !Number.isFinite(time))) {
    failures.push(`three ${label} UTC buckets missing`);
  } else if (times[1] - times[0] !== 60_000 || times[2] - times[1] !== 60_000) {
    failures.push(`${label} buckets are not consecutive 60-second UTC buckets`);
  }
  return rows.filter((row) => bucketStarts.includes(row.bucket_start));
}

/**
 * Verify durable delivery independently of the transport used to observe it.
 * `rows`, `closedBuckets`, `outbox`, and `receipts` must already be read-only
 * observations from the destination/gateway stores.
 */
export function verifyFreshDeliveryEvidence({
  group,
  table,
  rows = [],
  mappings = [],
  expectedByTag = {},
  closedBuckets = [],
  outbox = [],
  receipts = [],
  sharedTargetColumns = [],
  entityColumn = '',
} = {}) {
  const failures = [];
  const groupRevision = group?.applied_revision;
  const members = Array.isArray(group?.members) ? group.members : [];
  const mappingByMember = members.map((member) => ({
    member,
    mapping: mappings.find((mapping) => (
      sameIdentity(mapping.point_id, member.point_id) && sameIdentity(mapping.tag_id, member.tag_id) &&
      sameIdentity(mapping.device_id, member.device_id)
    )),
  }));
  const destinationTable = group?.destination_table_name ?? group?.destination?.table_name;
  const sharedColumns = new Set(sharedTargetColumns);

  if (!groupRevision) failures.push('managed group applied revision is missing');
  if (destinationTable && destinationTable !== table) failures.push('managed destination table identity differs');
  if (members.length !== Object.keys(expectedByTag).length) failures.push(`managed member count ${members.length}`);
  if (mappingByMember.some(({ mapping }) => !mapping?.tag_key)) failures.push('managed member mapping identity missing');
  if (mappingByMember.some(({ member }) => !member.target_column)) failures.push('managed member target column missing');

  const firstThree = checkConsecutiveBuckets(rows, failures, 'destination');
  const closedThree = checkConsecutiveBuckets(closedBuckets, failures, 'gateway');
  if (closedThree.length === 0 || closedThree.some((bucket) => bucket.group_revision !== groupRevision)) {
    failures.push('first closed bucket revision is not current');
  }
  const expectedEntities = new Set(members.map((member) => String(entityColumn ? member.entity_key ?? '' : member.device_id ?? '')));
  const destinationStarts = new Set(firstThree.map((row) => utcMillis(row.bucket_start)));
  for (const start of destinationStarts) {
    for (const entity of expectedEntities) {
      const matching = firstThree.filter((row) => utcMillis(row.bucket_start) === start &&
        String(entityColumn ? row[entityColumn] ?? '' : row.device_id ?? '') === entity);
      if (!entity || matching.length !== 1) failures.push(`destination entity rows missing or duplicated for ${entity || '<missing>'}`);
    }
  }

  for (const row of firstThree) {
    const effect = rowEffect(outbox, row);
    const receipt = rowReceipt(receipts, effect);
    const bucket = closedBuckets.find((item) => (
      sameIdentity(item.record_id, row.record_id) && sameBucketStart(item.bucket_start, row.bucket_start)
    ));
    if (!effect || effect.state !== 'sql_committed' || effect.group_id !== group?.id ||
      effect.group_revision !== groupRevision || effect.table_name !== table || !sameBucketStart(effect.bucket_start, row.bucket_start)) {
      failures.push(`row ${row.record_id ?? '<missing>'} lacks current committed outbox evidence`);
    }
    if (!receipt || receipt.payload_digest !== effect?.payload_digest || !receipt.committed_at) {
      failures.push(`row ${row.record_id ?? '<missing>'} lacks matching destination receipt`);
    }
    if (!bucket || bucket.group_revision !== groupRevision) failures.push(`row ${row.record_id ?? '<missing>'} lacks current closed bucket evidence`);
    const rowEntity = entityColumn ? row[entityColumn] : row.device_id;
    if (!row.record_id || row.group_id !== group?.id || !rowEntity || !row.bucket_start || !row.provenance) {
      failures.push(`row ${row.record_id ?? '<missing>'} identity/provenance incomplete`);
    }

    const bucketStart = utcMillis(row.bucket_start);
    const rowMembers = members.filter((member) => entityColumn
      ? String(member.entity_key ?? '') === String(rowEntity ?? '')
      : (!row.device_id || !member.device_id || String(member.device_id) === String(row.device_id)));
    if (rowMembers.length === 0) failures.push(`row ${row.record_id ?? '<missing>'} has an unknown entity`);
    const expectedRowMembers = new Set(rowMembers.map(memberIdentity));
    const provenance = parseProvenance(row.provenance);
    const actualMemberIDs = new Set(provenance?.map((entry) => entry?.member).filter((member) => typeof member === 'string'));
    const provenanceMembersExact = actualMemberIDs.size === expectedRowMembers.size &&
      actualMemberIDs.size === (provenance?.length ?? 0) &&
      [...actualMemberIDs].every((member) => expectedRowMembers.has(member));
    if (!provenance || !provenanceMembersExact || provenance.some((entry) => {
      const observedText = String(entry?.observed_at ?? '');
      const observedAt = utcMillis(observedText);
      return entry?.status !== 'ok' || entry?.quality !== 'good' || !entry?.sample_id ||
        !Number.isFinite(observedAt) || !Number.isFinite(bucketStart) ||
        observedAt < bucketStart || observedAt >= bucketStart + 60_000;
    })) {
      failures.push(`row ${row.record_id ?? '<missing>'} quality/source provenance incomplete`);
    }
    const frozen = parseProvenance(bucket?.members);
    if (!frozen || !provenance || frozen.length !== provenance.length || provenance.some((entry) => {
      const original = frozen.find((member) => member.member === entry.member);
      return !original || original.sample_id !== entry.sample_id || original.status !== entry.status ||
        original.quality !== entry.quality || !preciseUTC(original.observed_at) ||
        preciseUTC(original.observed_at) !== preciseUTC(entry.observed_at);
    })) failures.push(`row ${row.record_id ?? '<missing>'} does not preserve frozen acquisition provenance`);

    for (const { member, mapping } of mappingByMember) {
      const value = row[member.target_column];
      const expected = expectedByTag[mapping?.tag_key];
      const memberIsInRow = expectedRowMembers.has(memberIdentity(member));
      if (memberIsInRow && (value === null || value === undefined || expected === undefined || String(value) !== expected)) {
        failures.push(`row ${row.record_id ?? '<missing>'} value mismatch for ${mapping?.tag_key ?? '<unknown>'}`);
      }
      if (!memberIsInRow && !sharedColumns.has(member.target_column) && value !== null && value !== undefined && value !== '') {
        failures.push(`row ${row.record_id ?? '<missing>'} has a value for another entity member ${member.target_column}`);
      }
    }
  }

  return {
    failures,
    mapping_by_member: mappingByMember,
    outbox,
    receipts,
    closed_buckets: closedBuckets.slice(0, 3),
  };
}
