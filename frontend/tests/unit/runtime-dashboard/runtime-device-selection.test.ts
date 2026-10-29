import { describe, expect, it } from 'vitest';
import { parseRuntimeDeviceSelection } from '../../../src/features/datalink/runtime-dashboard/runtimeDeviceSelection';

describe('strict raw runtime device selection', () => {
  it.each(['', '?', '?keep=a%20b', '?different=device-A', '?keep=Server+mode&x=OK%2FNG'])('retains absent selector: %s', search => {
    expect(parseRuntimeDeviceSelection(search)).toEqual({ kind: 'absent' });
  });
  it.each(['?device_id=device-A', '?%64evice_id=device-A', '?keep=x&device_id=%64evice-A', 'device_id=device-A&keep=y'])('accepts exactly one safe opaque identity: %s', search => {
    expect(parseRuntimeDeviceSelection(search)).toEqual({ kind: 'valid', deviceId: 'device-A' });
  });
  it.each(['?device_id', '?device_id=', '?device_id=+', '?device_id=%ZZ', '?device_id=%', '?device_id=%C0%AF', '?device_id=%E0%A4%A', '?%ZZ=device-A', '?keep=%ZZ', '?device_id=device-A&device_id=device-A', '?device_id=device-A&%64evice_id=device-B', '?device_id=device-A&device%5Fid=', '?device_id=host.example', '?device_id=192.168.1.1', '?device_id=%5B%3A%3A1%5D%3A5432', '?device_id=%2Fprivate%2Fdb', '?device_id=password%3Dsecret', '?device_id=device-A%00', '?device_id=device-A%0A', '?device_id=device-A;device_id=device-B', `?device_id=${'a'.repeat(129)}`])('rejects malformed, ambiguous or unsafe explicit input: %s', search => {
    expect(parseRuntimeDeviceSelection(search)).toEqual({ kind: 'invalid' });
  });
});
