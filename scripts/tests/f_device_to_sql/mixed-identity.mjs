// Point identity assertions for the seven-member mixed fixture. These checks
// deliberately use fixture literals, so a non-empty or wrong-shaped API result
// cannot make the acceptance witness look healthy.

export const MIXED_POINT_EXPECTATIONS = Object.freeze({
  A: Object.freeze([
    Object.freeze({ address: '40001', data_type: 'int16', last_value: '215' }),
    Object.freeze({ address: '40002', data_type: 'uint16', last_value: '1013' }),
    Object.freeze({ address: '40003', data_type: 'uint64', last_value: '9007199254740993' }),
    Object.freeze({ address: '40007', data_type: 'bool', last_value: 'true' }),
    Object.freeze({ address: '40008', data_type: 'float32', last_value: '1.5' }),
    Object.freeze({ address: '40010', data_type: 'float64', last_value: '123.456' }),
    Object.freeze({ address: '40014', data_type: 'string', last_value: 'A1' }),
  ]),
  B: Object.freeze([
    Object.freeze({ address: '40001', data_type: 'int16', last_value: '187' }),
    Object.freeze({ address: '40002', data_type: 'uint16', last_value: '777' }),
    Object.freeze({ address: '40003', data_type: 'uint64', last_value: '2' }),
    Object.freeze({ address: '40007', data_type: 'bool', last_value: 'false' }),
    Object.freeze({ address: '40008', data_type: 'float32', last_value: '-2.25' }),
    Object.freeze({ address: '40010', data_type: 'float64', last_value: '-654.321' }),
    Object.freeze({ address: '40014', data_type: 'string', last_value: 'B2' }),
  ]),
});

const GROUP_LINES = Object.freeze({ 'Line A': 'A', 'Line B': 'B' });

function pointIdentity(point) {
  return `${String(point?.address)}:${String(point?.data_type)}`;
}

function pointWitness(point) {
  return {
    id: point.id,
    device_id: point.device_id,
    address: point.address,
    data_type: point.data_type,
    last_value: point.last_value,
    last_error: point.last_error,
    last_read_at: point.last_read_at,
  };
}

/** Read and validate persisted points without treating length as identity. */
export async function collectAndValidatePointObservations(groups, readPoints) {
  const failures = [];
  const observations = [];
  const sourceGroups = Array.isArray(groups) ? groups : [];

  for (const [groupName, line] of Object.entries(GROUP_LINES)) {
    const group = sourceGroups.find((candidate) => candidate?.name === groupName);
    const expected = MIXED_POINT_EXPECTATIONS[line];
    if (!group) {
      failures.push(`${groupName}: persisted group is missing`);
      continue;
    }
    const members = Array.isArray(group.members) ? group.members : [];
    const deviceIDs = [...new Set(members.map((member) => member?.device_id).filter(Boolean))];
    if (deviceIDs.length !== 1) {
      failures.push(`${groupName}: expected one persisted device identity, got ${deviceIDs.length}`);
      continue;
    }
    const deviceID = deviceIDs[0];
    let points;
    try {
      points = await readPoints(deviceID);
    } catch (error) {
      failures.push(`${groupName}: point API failed for ${deviceID}: ${error.message.split('\n')[0]}`);
      observations.push({ device_id: deviceID, error: error.message.split('\n')[0] });
      continue;
    }
    if (!Array.isArray(points)) {
      failures.push(`${groupName}: point API returned non-array data`);
      observations.push({ device_id: deviceID, error: 'point API returned non-array data' });
      continue;
    }
    const memberIDs = new Set(members.map((member) => member?.point_id));
    observations.push(...points.filter((point) => memberIDs.has(point?.id)).map(pointWitness));
    const seen = new Set();
    for (const member of members) {
      const point = points.find((candidate) => candidate?.id === member?.point_id);
      if (!point) {
        failures.push(`${groupName}: point ${member?.point_id} is missing from device API result`);
        continue;
      }
      if (point.device_id !== deviceID) {
        failures.push(`${groupName}: point ${point.id} belongs to ${point.device_id}, expected ${deviceID}`);
      }
      const expectedPoint = expected.find((candidate) => pointIdentity(candidate) === pointIdentity(point));
      if (!expectedPoint) {
        failures.push(`${groupName}: point ${point.id} has unexpected address/type ${pointIdentity(point)}`);
        continue;
      }
      const identity = pointIdentity(expectedPoint);
      if (seen.has(identity)) failures.push(`${groupName}: duplicate point identity ${identity}`);
      seen.add(identity);
      if (String(point.last_value) !== expectedPoint.last_value) {
        failures.push(`${groupName}: ${identity} expected ${expectedPoint.last_value} got ${String(point.last_value)}`);
      }
    }
    for (const expectedPoint of expected) {
      const identity = pointIdentity(expectedPoint);
      if (!seen.has(identity)) failures.push(`${groupName}: expected point identity ${identity} was not observed`);
    }
  }
  return { observations, failures };
}
