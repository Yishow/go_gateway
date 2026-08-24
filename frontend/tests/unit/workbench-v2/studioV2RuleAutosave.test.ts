import { describe, expect, it } from 'vitest';
import {
  hydrateStudioV2Rule,
  toStudioV2RuleCreateRequest,
  toStudioV2RuleUpdateRequest,
} from '../../../src/features/datalink/workbench-v2/state/studioV2RuleAutosave';
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
});
