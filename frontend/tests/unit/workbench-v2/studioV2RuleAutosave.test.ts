import { describe, expect, it } from 'vitest';
import {
  canHydrateStudioV2Rule,
  hydrateStudioV2Rule,
  toStudioV2RuleCreateRequest,
  toStudioV2RuleUpdateRequest,
} from '../../../src/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { getRuleReadinessReason, isStep2Ready } from '../../../src/features/datalink/workbench-v2/state/sourceRule';
import type { Rule } from '../../../src/features/datalink/workbench-v2/state/types';
import type { SourceRuleRecord } from '../../../src/types/datalink';

const draftRule: Rule = {
  id: 'rule-share',
  device_id: 'device-share',
  name: 'Share rule',
  start_address: '40001',
  count: 4,
  data_type: 'int16',
  naming_prefix: 'SHARE_',
  enabled: true,
  scale_multiplier: 1,
  scale_offset: 0,
  data_format: '',
  skipped_addresses: [],
  share_enabled: true,
  share_start_register: 40001,
  share_stride: 1,
};

const persistedRecord: SourceRuleRecord = {
  id: 'rule-share',
  device_id: 'device-share',
  start_address: '40001',
  count: 4,
  data_type: 'int16',
  naming_prefix: 'SHARE_',
  enabled: true,
  locked: false,
  origin: 'manual',
  skipped_addresses: [],
  share_enabled: true,
  share_start_register: 40021,
  share_stride: 2,
  created_at: '',
  updated_at: '',
};

describe('Studio V2 source rule Share contract', () => {
  it('includes Share fields in create and update requests', () => {
    expect(toStudioV2RuleCreateRequest(draftRule)).toEqual(expect.objectContaining({
      share_enabled: true,
      share_start_register: 40001,
      share_stride: 1,
    }));
    expect(toStudioV2RuleUpdateRequest(draftRule)).toEqual(expect.objectContaining({
      share_enabled: true,
      share_start_register: 40001,
      share_stride: 1,
    }));
  });

  it('hydrates persisted Share fields instead of resetting them to defaults', () => {
    const hydrated = hydrateStudioV2Rule(persistedRecord);

    expect(hydrated.share_enabled).toBe(true);
    expect(hydrated.share_start_register).toBe(40021);
    expect(hydrated.share_stride).toBe(2);
  });

  it('keeps blank ownership and address fields for fail-closed readiness', () => {
    const blankDeviceRecord: SourceRuleRecord = {
      ...persistedRecord,
      id: 'rule-blank-device',
      device_id: '   ',
    };
    const blankAddressRecord: SourceRuleRecord = {
      ...persistedRecord,
      id: 'rule-blank-address',
      start_address: '   ',
    };

    expect(canHydrateStudioV2Rule(blankDeviceRecord)).toBe(true);
    expect(canHydrateStudioV2Rule(blankAddressRecord)).toBe(true);

    const hydratedBlankDeviceRule = hydrateStudioV2Rule(blankDeviceRecord);
    const hydratedBlankAddressRule = hydrateStudioV2Rule(blankAddressRecord);
    const hydratedValidRule = hydrateStudioV2Rule({
      ...persistedRecord,
      id: 'rule-valid-device',
    });
    const protocolMap = { 'device-share': 'modbus_tcp' as const };

    expect(hydratedBlankDeviceRule.device_id).toBe('   ');
    expect(getRuleReadinessReason(hydratedBlankDeviceRule, protocolMap)).toBe('unknown_device');
    expect(isStep2Ready([hydratedValidRule, hydratedBlankDeviceRule], protocolMap)).toBe(false);
    expect(hydratedBlankAddressRule.start_address).toBe('   ');
    expect(getRuleReadinessReason(hydratedBlankAddressRule, protocolMap)).toBe('invalid_address');
    expect(isStep2Ready([hydratedValidRule, hydratedBlankAddressRule], protocolMap)).toBe(false);
  });
});
