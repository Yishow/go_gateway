import { describe, expect, it } from 'vitest';
import {
  coerceBooleanFlag,
  resolveGatewayCreateEntryPath,
} from '@/features/gateway/dualEntryFlag';

describe('coerceBooleanFlag', () => {
  it('should parse truthy values', () => {
    expect(coerceBooleanFlag(true)).toBe(true);
    expect(coerceBooleanFlag(1)).toBe(true);
    expect(coerceBooleanFlag('true')).toBe(true);
    expect(coerceBooleanFlag('ENABLED')).toBe(true);
    expect(coerceBooleanFlag('on')).toBe(true);
  });

  it('should parse falsy values', () => {
    expect(coerceBooleanFlag(false, true)).toBe(false);
    expect(coerceBooleanFlag(0, true)).toBe(false);
    expect(coerceBooleanFlag('false', true)).toBe(false);
    expect(coerceBooleanFlag('disabled', true)).toBe(false);
    expect(coerceBooleanFlag('off', true)).toBe(false);
  });

  it('should fallback for unsupported values', () => {
    expect(coerceBooleanFlag('unknown', true)).toBe(true);
    expect(coerceBooleanFlag(undefined, false)).toBe(false);
    expect(coerceBooleanFlag(null, true)).toBe(true);
  });
});

describe('resolveGatewayCreateEntryPath', () => {
  it('routes to dual entry when enabled', () => {
    expect(resolveGatewayCreateEntryPath(true)).toBe('/gateway/entry');
  });

  it('routes to legacy datalink when disabled', () => {
    expect(resolveGatewayCreateEntryPath(false)).toBe('/datalink');
  });
});
