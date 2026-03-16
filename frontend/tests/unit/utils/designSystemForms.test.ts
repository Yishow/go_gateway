import { describe, expect, it } from 'vitest';
import { designSystem } from '@/styles/designSystem';

describe('designSystem forms contract', () => {
  it('defines autocomplete tokens used by datalink search inputs', () => {
    expect(designSystem.forms.autocomplete.search).toBe('off');
  });
});
