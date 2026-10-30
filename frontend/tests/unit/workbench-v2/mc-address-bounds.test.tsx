import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { RangeSummary } from '../../../src/features/datalink/workbench-v2/steps/step2/RangeSummary';
import { addressParser } from '../../../src/utils/addressParser';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, value?: string | { defaultValue?: string }) =>
      typeof value === 'string' ? value : (value?.defaultValue ?? key),
  }),
}));

describe('MC 3E address bounds', () => {
  it('accepts the maximum client offset and its final one-step expansion', () => {
    expect(addressParser.validate('WFFFFFF', 'mc_3e').valid).toBe(true);
    expect(addressParser.offset('WFFFFFE', 1, 'mc_3e')).toBe('WFFFFFF');
  });

  it.each(['W1000000', 'WFFFFFFFF', 'W7FFFFFFFFFFFFFFF', 'D9223372036854775807'])(
    'rejects %s outside the client address range',
    (address) => {
      expect(addressParser.validate(address, 'mc_3e').valid).toBe(false);
      expect(addressParser.expand(address, 1, 'mc_3e')).toEqual([]);
    },
  );

  it('rejects expansion overflow instead of wrapping to W0', () => {
    expect(() => addressParser.offset('WFFFFFF', 1, 'mc_3e')).toThrow(/range|overflow/i);
    expect(addressParser.expand('WFFFFFE', 3, 'mc_3e')).toEqual([]);
  });

  it('renders a blocked summary for a valid maximum start with a multiword overflow', () => {
    render(
      <RangeSummary
        startAddress="WFFFFFF"
        count={1}
        dataType="int64"
        protocol="mc_3e"
      />,
    );

    expect(screen.getByTestId('range-summary')).toHaveAttribute('aria-invalid', 'true');
    expect(screen.getByTestId('range-summary-error')).toBeInTheDocument();
  });
});
